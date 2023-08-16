package engine_test

import (
	"context"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"

	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/execution/model"
	engstore "github.com/thanos-io/promql-engine/execution/storage"
	"github.com/thanos-io/promql-engine/logicalplan"
	"github.com/thanos-io/promql-engine/parser"
	"github.com/thanos-io/promql-engine/query"
)

func TestUserDefinedOperators(t *testing.T) {
	opts := promql.EngineOpts{
		Timeout:    1 * time.Hour,
		MaxSamples: 1e10,
	}

	load := `
load 30s
	http_requests_total{container="a"} 1x30
	http_requests_total{container="b"} 2x30`

	test, err := promql.NewTest(t, load)
	testutil.Ok(t, err)
	defer test.Close()
	testutil.Ok(t, test.Run())

	newEngine := engine.New(engine.Opts{
		EngineOpts:      opts,
		DisableFallback: true,
		LogicalOptimizers: []logicalplan.Optimizer{
			&injectVectorSelector{},
		},
	})
	query := "sum(http_requests_total)"
	qry, err := newEngine.NewRangeQuery(test.Context(), test.Storage(), nil, query, time.Unix(0, 0), time.Unix(90, 0), 30*time.Second)
	testutil.Ok(t, err)

	result := qry.Exec(context.Background())
	testutil.Ok(t, result.Err)

	expected := promql.Matrix{
		promql.Series{
			Metric: labels.EmptyLabels(),
			Floats: []promql.FPoint{{T: 0, F: 14}, {T: 30000, F: 14}, {T: 60000, F: 14}, {T: 90000, F: 14}},
		},
	}
	mat, err := result.Matrix()
	testutil.Ok(t, err)
	testutil.Equals(t, expected, mat)
}

type injectVectorSelector struct{}

func (i injectVectorSelector) Optimize(expr parser.Expr, _ *logicalplan.Opts) parser.Expr {
	logicalplan.TraverseBottomUp(nil, &expr, func(_, current *parser.Expr) bool {
		switch (*current).(type) {
		case *parser.VectorSelector:
			*current = &logicalVectorSelector{
				VectorSelector: (*current).(*parser.VectorSelector),
			}
		}
		return false
	})
	return expr
}

type logicalVectorSelector struct {
	*parser.VectorSelector
}

func (c logicalVectorSelector) MakeExecutionOperator(stepsBatch int, vectors *model.VectorPool, selectors *engstore.SelectorPool, opts *query.Options, hints storage.SelectHints) (model.VectorOperator, error) {
	return &vectorSelectorOperator{
		stepsBatch: stepsBatch,
		vectors:    vectors,

		mint:        opts.Start.UnixMilli(),
		maxt:        opts.End.UnixMilli(),
		step:        opts.Step.Milliseconds(),
		currentStep: opts.Start.UnixMilli(),
	}, nil
}

type vectorSelectorOperator struct {
	stepsBatch int
	vectors    *model.VectorPool

	mint        int64
	maxt        int64
	step        int64
	currentStep int64
}

func (c *vectorSelectorOperator) Next(ctx context.Context) ([]model.StepVector, error) {
	if c.currentStep > c.maxt {
		return nil, nil
	}
	vectors := c.vectors.GetVectorBatch()
	for i := 0; i < c.stepsBatch && c.currentStep <= c.maxt; i++ {
		vector := c.vectors.GetStepVector(c.currentStep)
		vector.AppendSample(c.vectors, 1, 7)
		vector.AppendSample(c.vectors, 2, 7)
		vectors = append(vectors, vector)
		c.currentStep += c.step
	}
	return vectors, nil
}

func (c *vectorSelectorOperator) Series(ctx context.Context) ([]labels.Labels, error) {
	return []labels.Labels{
		labels.FromStrings(labels.MetricName, "http_requests_total", "container", "a"),
		labels.FromStrings(labels.MetricName, "http_requests_total", "container", "b"),
	}, nil
}

func (c *vectorSelectorOperator) GetPool() *model.VectorPool {
	return c.vectors
}

func (c *vectorSelectorOperator) Explain() model.Explanation {
	return model.Explanation{Operator: "[*vectorSelectorOperator]"}
}
