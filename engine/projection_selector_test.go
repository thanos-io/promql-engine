// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/logicalplan"
	"github.com/thanos-io/promql-engine/query"
	prometheusStorage "github.com/thanos-io/promql-engine/storage/prometheus"

	"github.com/cortexproject/promqlsmith"
	"github.com/efficientgo/core/errors"
	"github.com/efficientgo/core/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/promql/promqltest"
	"github.com/prometheus/prometheus/storage"
)

// TestProjectionSelectorPushdown verifies that the projectedSelector path
// (triggered by ProjectionOptimizer with SeriesHashLabel + WithSeriesHashLabel on scanners)
// produces correct results compared to the standard Prometheus engine for various
// aggregation queries.
func TestProjectionSelectorPushdown(t *testing.T) {
	t.Parallel()

	load := `load 30s
		http_requests_total{pod="nginx-1", job="app", env="prod", instance="1", cluster="us-west-2"} 1+1x40
		http_requests_total{pod="nginx-2", job="app", env="dev", instance="2", cluster="us-west-2"} 2+2x40
		http_requests_total{pod="nginx-3", job="api", env="prod", instance="3", cluster="us-east-1"} 3+3x40
		http_requests_total{pod="nginx-4", job="api", env="dev", instance="4", cluster="us-east-1"} 4+4x40
		http_requests_total{pod="nginx-5", job="web", env="staging", instance="5", cluster="eu-west-1"} 5+5x40
		errors_total{pod="nginx-1", job="app", env="prod", instance="1", cluster="us-west-2"} 0.5+0.5x40
		errors_total{pod="nginx-2", job="app", env="dev", instance="2", cluster="us-west-2"} 1+1x40
		errors_total{pod="nginx-3", job="api", env="prod", instance="3", cluster="us-east-1"} 1.5+1.5x40
		errors_total{pod="nginx-4", job="api", env="dev", instance="4", cluster="us-east-1"} 2+2x40`

	testStorage := promqltest.LoadedStorage(t, load)
	defer testStorage.Close()

	start := time.Unix(0, 0)
	end := time.Unix(1200, 0)
	step := 30 * time.Second
	queryTime := time.Unix(600, 0)

	engineOpts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}

	// Reference engine: standard Thanos engine with all optimizers (no selector projection)
	referenceEngine := engine.New(engine.Opts{
		EngineOpts:        engineOpts,
		LogicalOptimizers: logicalplan.AllOptimizers,
	})

	cases := []struct {
		name  string
		query string
	}{
		// sum() — all labels dropped, projection is include=true, labels=[]
		{name: "sum all labels dropped", query: `sum(http_requests_total)`},
		// sum by (job) — only job label needed
		{name: "sum by single label", query: `sum by (job) (http_requests_total)`},
		// sum by multiple labels
		{name: "sum by multiple labels", query: `sum by (job, env) (http_requests_total)`},
		// sum without — exclude mode
		{name: "sum without", query: `sum without (instance, pod) (http_requests_total)`},
		// count — all labels dropped
		{name: "count all labels dropped", query: `count(http_requests_total)`},
		// count by (cluster)
		{name: "count by cluster", query: `count by (cluster) (http_requests_total)`},
		// avg by (job)
		{name: "avg by job", query: `avg by (job) (http_requests_total)`},
		// nested function with projection: count(rate(metric[5m]))
		{name: "count of rate", query: `count(rate(http_requests_total[5m]))`},
		// sum(rate())
		{name: "sum of rate", query: `sum(rate(http_requests_total[5m]))`},
		// sum by (job)(rate())
		{name: "sum by job of rate", query: `sum by (job) (rate(http_requests_total[5m]))`},
		// min/max aggregations
		{name: "min", query: `min(http_requests_total)`},
		{name: "max by job", query: `max by (job) (http_requests_total)`},
		// binary expression with shared metric, both sides projected
		{name: "binary sum or count", query: `sum(http_requests_total) or count(http_requests_total)`},
		// binary with different metrics
		{name: "sum ratio", query: `sum(errors_total) / sum(http_requests_total)`},
		// multiple levels of aggregation
		{name: "nested aggregation", query: `sum(avg by (job) (http_requests_total))`},
		// sum with increase
		{name: "sum of increase", query: `sum(increase(http_requests_total[5m]))`},
		// sum by with delta
		{name: "sum by job of delta", query: `sum by (job) (delta(http_requests_total[5m]))`},
		// avg without instance
		{name: "avg without instance", query: `avg without (instance) (http_requests_total)`},
		// sum by job of irate
		{name: "sum by job of irate", query: `sum by (job) (irate(http_requests_total[1m]))`},
	}

	ctx := context.Background()
	for _, tc := range cases {
		t.Run(tc.name+"/instant", func(t *testing.T) {
			// Create projection engine with scanners that have WithSeriesHashLabel
			projEngine := newProjectionSelectorEngine(t, testStorage, queryTime, queryTime, engineOpts)

			refQuery, err := referenceEngine.NewInstantQuery(ctx, testStorage, nil, tc.query, queryTime)
			testutil.Ok(t, err)
			defer refQuery.Close()
			refResult := refQuery.Exec(ctx)
			testutil.Ok(t, refResult.Err, "reference query failed: %s", tc.query)

			projQuery, err := projEngine.MakeInstantQuery(ctx, testStorage, &engine.QueryOpts{}, tc.query, queryTime)
			testutil.Ok(t, err)
			defer projQuery.Close()
			projResult := projQuery.Exec(ctx)
			testutil.Ok(t, projResult.Err, "projection query failed: %s", tc.query)

			if diff := cmp.Diff(refResult, projResult, comparer); diff != "" {
				t.Errorf("Instant results differ for query %q:\n%s", tc.query, diff)
			}
		})

		t.Run(tc.name+"/range", func(t *testing.T) {
			projEngine := newProjectionSelectorEngine(t, testStorage, start, end, engineOpts)

			refQuery, err := referenceEngine.NewRangeQuery(ctx, testStorage, nil, tc.query, start, end, step)
			testutil.Ok(t, err)
			defer refQuery.Close()
			refResult := refQuery.Exec(ctx)
			testutil.Ok(t, refResult.Err, "reference range query failed: %s", tc.query)

			projQuery, err := projEngine.MakeRangeQuery(ctx, testStorage, &engine.QueryOpts{}, tc.query, start, end, step)
			testutil.Ok(t, err)
			defer projQuery.Close()
			projResult := projQuery.Exec(ctx)
			testutil.Ok(t, projResult.Err, "projection range query failed: %s", tc.query)

			if diff := cmp.Diff(refResult, projResult, comparer); diff != "" {
				t.Errorf("Range results differ for query %q:\n%s", tc.query, diff)
			}
		})
	}
}

// TestProjectionSelectorWithFuzz uses promqlsmith to generate random queries with
// aggregations and verifies that the projection selector engine produces the same
// results as the reference engine.
func TestProjectionSelectorWithFuzz(t *testing.T) {
	t.Parallel()

	seed := time.Now().UnixNano()
	rnd := rand.New(rand.NewSource(seed))
	numRuns := 10000

	load := `load 30s
		http_requests_total{pod="nginx-1", job="app", env="prod", instance="1", cluster="us-west-2"} 1+1x40
		http_requests_total{pod="nginx-2", job="app", env="dev", instance="2", cluster="us-west-2"} 2+2x40
		http_requests_total{pod="nginx-3", job="api", env="prod", instance="3", cluster="us-east-1"} 3+3x40
		http_requests_total{pod="nginx-4", job="api", env="dev", instance="4", cluster="us-east-1"} 4+4x40
		http_requests_duration_seconds_bucket{pod="nginx-1", job="app", env="prod", instance="1", le="0.1"} 1+1x40
		http_requests_duration_seconds_bucket{pod="nginx-1", job="app", env="prod", instance="1", le="0.5"} 3+2x40
		http_requests_duration_seconds_bucket{pod="nginx-1", job="app", env="prod", instance="1", le="+Inf"} 4+2x40
		errors_total{pod="nginx-1", job="app", env="prod", instance="1", cluster="us-west-2"} 0.5+0.5x40
		errors_total{pod="nginx-2", job="app", env="dev", instance="2", cluster="us-west-2"} 1+1x40
		errors_total{pod="nginx-3", job="api", env="prod", instance="3", cluster="us-east-1"} 1.5+1.5x40
		errors_total{pod="nginx-4", job="api", env="dev", instance="4", cluster="us-east-1"} 2+2x40`

	testStorage := promqltest.LoadedStorage(t, load)
	defer testStorage.Close()

	seriesSet, err := getSeries(context.Background(), testStorage, "http_requests_total")
	testutil.Ok(t, err)

	psOpts := []promqlsmith.Option{
		promqlsmith.WithEnableOffset(false),
		promqlsmith.WithEnableAtModifier(false),
		promqlsmith.WithEnabledAggrs([]parser.ItemType{
			parser.SUM, parser.MIN, parser.MAX, parser.AVG, parser.COUNT,
		}),
		promqlsmith.WithEnableVectorMatching(true),
	}
	ps := promqlsmith.New(rnd, seriesSet, psOpts...)

	engineOpts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}

	referenceEngine := engine.New(engine.Opts{
		EngineOpts:        engineOpts,
		LogicalOptimizers: logicalplan.AllOptimizers,
	})

	ctx := context.Background()
	queryTime := time.Unix(600, 0)

	t.Logf("Running %d fuzz tests with seed %d", numRuns, seed)
	for i := range numRuns {
		var expr parser.Expr
		var query string

		for {
			expr = ps.WalkInstantQuery()
			query = expr.Pretty(0)

			if !containsProjectableAggregation(expr) {
				continue
			}

			_, err := referenceEngine.NewInstantQuery(ctx, testStorage, nil, query, queryTime)
			if err != nil {
				continue
			}
			break
		}

		t.Run(fmt.Sprintf("Query_%d", i), func(t *testing.T) {
			projEngine := newProjectionSelectorEngine(t, testStorage, queryTime, queryTime, engineOpts)

			refQuery, err := referenceEngine.NewInstantQuery(ctx, testStorage, &engine.QueryOpts{}, query, queryTime)
			testutil.Ok(t, err)
			defer refQuery.Close()

			refResult := refQuery.Exec(ctx)
			if refResult.Err != nil {
				// Query failed in reference engine, skip.
				return
			}

			projQuery, err := projEngine.MakeInstantQuery(ctx, testStorage, &engine.QueryOpts{}, query, queryTime)
			testutil.Ok(t, err)
			defer projQuery.Close()

			projResult := projQuery.Exec(ctx)
			testutil.Ok(t, projResult.Err, "query: %s", query)

			if diff := cmp.Diff(refResult, projResult, comparer); diff != "" {
				t.Errorf("Results differ for query %s:\n%s", query, diff)
			}
		})
	}
}

// containsProjectableAggregation checks if the expression contains aggregations
// that benefit from projection pushdown (sum, count, min, max, avg — NOT topk/bottomk).
func containsProjectableAggregation(expr parser.Expr) bool {
	found := false
	parser.Inspect(expr, func(node parser.Node, _ []parser.Node) error {
		switch n := node.(type) {
		case *parser.AggregateExpr:
			switch n.Op {
			case parser.SUM, parser.COUNT, parser.MIN, parser.MAX, parser.AVG:
				found = true
				return errors.New("found")
			}
		}
		return nil
	})
	return found
}

// newProjectionSelectorEngine creates a Thanos engine configured with the projectedSelector
// path. It uses NewWithScanners with WithSeriesHashLabel so that the ProjectionOptimizer
// triggers actual label projection at the selector level.
func newProjectionSelectorEngine(t testing.TB, queryable storage.Queryable, mint, maxt time.Time, engineOpts promql.EngineOpts) *engine.Engine {
	t.Helper()

	qOpts := &query.Options{
		Start:               mint,
		End:                 maxt,
		Step:                30 * time.Second,
		StepsBatch:          10,
		LookbackDelta:       5 * time.Minute,
		DecodingConcurrency: 1,
	}

	scanners, err := prometheusStorage.NewPrometheusScanners(
		queryable, qOpts, nil,
		prometheusStorage.WithSeriesHashLabel("__series_hash__"),
	)
	testutil.Ok(t, err)

	return engine.NewWithScanners(engine.Opts{
		EngineOpts: engineOpts,
		LogicalOptimizers: []logicalplan.Optimizer{
			logicalplan.SortMatchers{},
			logicalplan.ProjectionOptimizer{SeriesHashLabel: "__series_hash__"},
			logicalplan.MergeSelectsOptimizer{},
		},
	}, scanners)
}
