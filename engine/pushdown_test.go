// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/logicalplan"

	"slices"

	"github.com/cortexproject/promqlsmith"
	"github.com/efficientgo/core/errors"
	"github.com/efficientgo/core/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/promql/promqltest"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	"github.com/prometheus/prometheus/util/annotations"
)

type projectionQuerier struct {
	storage.Querier
}

type projectionSeriesSet struct {
	storage.SeriesSet
	hints *storage.SelectHints
}

func (m projectionSeriesSet) Next() bool { return m.SeriesSet.Next() }
func (m projectionSeriesSet) At() storage.Series {
	// Get the original series
	originalSeries := m.SeriesSet.At()
	if originalSeries == nil {
		return nil
	}
	// If no projection hints, return the original series
	if m.hints == nil {
		return originalSeries
	}
	if !m.hints.By && len(m.hints.Grouping) == 0 {
		return originalSeries
	}

	// Apply projection based on hints
	originalLabels := originalSeries.Labels()
	var projectedLabels labels.Labels

	if m.hints.By {
		// Include mode: only keep the labels in the grouping
		builder := labels.NewBuilder(labels.EmptyLabels())
		originalLabels.Range(func(l labels.Label) {
			if slices.Contains(m.hints.Grouping, l.Name) {
				builder.Set(l.Name, l.Value)
			}
		})
		builder.Set("__series_hash__", strconv.FormatUint(originalLabels.Hash(), 10))
		projectedLabels = builder.Labels()
	} else {
		// Exclude mode: keep all labels except those in the grouping
		excludeMap := make(map[string]struct{})
		for _, groupLabel := range m.hints.Grouping {
			excludeMap[groupLabel] = struct{}{}
		}

		builder := labels.NewBuilder(labels.EmptyLabels())
		originalLabels.Range(func(l labels.Label) {
			if _, excluded := excludeMap[l.Name]; !excluded {
				builder.Set(l.Name, l.Value)
			}
		})
		builder.Set("__series_hash__", strconv.FormatUint(originalLabels.Hash(), 10))
		projectedLabels = builder.Labels()
	}

	// Return a projected series that wraps the original but with filtered labels
	return &projectedSeries{
		Series: originalSeries,
		lset:   projectedLabels,
	}
}

// projectedSeries wraps a storage.Series but returns projected labels.
type projectedSeries struct {
	storage.Series
	lset labels.Labels
}

func (s *projectedSeries) Labels() labels.Labels {
	return s.lset
}

func (s *projectedSeries) Iterator(iter chunkenc.Iterator) chunkenc.Iterator {
	return s.Series.Iterator(iter)
}

func (m projectionSeriesSet) Err() error                        { return m.SeriesSet.Err() }
func (m projectionSeriesSet) Warnings() annotations.Annotations { return m.SeriesSet.Warnings() }

// Implement the Querier interface methods.
func (m *projectionQuerier) Select(ctx context.Context, sortSeries bool, hints *storage.SelectHints, matchers ...*labels.Matcher) storage.SeriesSet {
	return projectionSeriesSet{
		SeriesSet: m.Querier.Select(ctx, sortSeries, hints, matchers...),
		hints:     hints,
	}
}
func (m *projectionQuerier) LabelValues(ctx context.Context, name string, _ *storage.LabelHints, matchers ...*labels.Matcher) ([]string, annotations.Annotations, error) {
	return nil, nil, nil
}
func (m *projectionQuerier) LabelNames(ctx context.Context, _ *storage.LabelHints, matchers ...*labels.Matcher) ([]string, annotations.Annotations, error) {
	return nil, nil, nil
}
func (m *projectionQuerier) Close() error { return nil }

// projectionQueryable is a storage.Queryable that applies projection to the querier.
type projectionQueryable struct {
	storage.Queryable
}

func (q *projectionQueryable) Querier(mint, maxt int64) (storage.Querier, error) {
	querier, err := q.Queryable.Querier(mint, maxt)
	if err != nil {
		return nil, err
	}
	return &projectionQuerier{
		Querier: querier,
	}, nil
}

func TestProjectionPushdownWithFuzz(t *testing.T) {
	t.Parallel()

	// Define test parameters
	seed := time.Now().UnixNano()
	rnd := rand.New(rand.NewSource(seed))
	testRuns := 1000

	// Create test data
	load := `load 30s
		http_requests_total{pod="nginx-1", job="app", env="prod", instance="1"} 1+1x40
		http_requests_total{pod="nginx-2", job="app", env="dev", instance="2"} 2+2x40
		http_requests_total{pod="nginx-3", job="api", env="prod", instance="3"} 3+3x40
		http_requests_total{pod="nginx-4", job="api", env="dev", instance="4"} 4+4x40
		errors_total{pod="nginx-1", job="app", env="prod", instance="1", cluster="us-west-2"} 0.5+0.5x40
		errors_total{pod="nginx-2", job="app", env="dev", instance="2", cluster="us-west-2"} 1+1x40
		errors_total{pod="nginx-3", job="api", env="prod", instance="3", cluster="us-east-2"} 1.5+1.5x40
		errors_total{pod="nginx-4", job="api", env="dev", instance="4", cluster="us-east-1"} 2+2x40`

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	// Get series for PromQLSmith
	seriesSet, err := getSeries(context.Background(), storage)
	testutil.Ok(t, err)

	// Configure PromQLSmith
	psOpts := []promqlsmith.Option{
		promqlsmith.WithEnableOffset(false),
		promqlsmith.WithEnableAtModifier(false),
		// Focus on aggregations that benefit from projection pushdown
		promqlsmith.WithEnabledAggrs([]parser.ItemType{
			parser.SUM, parser.MIN, parser.MAX, parser.AVG, parser.COUNT, parser.TOPK, parser.BOTTOMK,
		}),
		promqlsmith.WithEnableVectorMatching(true),
	}
	ps := promqlsmith.New(rnd, seriesSet, psOpts...)

	// Engine options
	engineOpts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}

	normalEngine := engine.New(engine.Opts{
		EngineOpts:                  engineOpts,
		LogicalOptimizers:           logicalplan.AllOptimizers,
		DisableDuplicateLabelChecks: false,
	})

	pushdownEngine := engine.New(engine.Opts{
		EngineOpts: engineOpts,
		// ProjectionPushdown optimizer doesn't support merge selects optimizer
		// so disable it for now.
		LogicalOptimizers: []logicalplan.Optimizer{
			logicalplan.SortMatchers{},
			logicalplan.ProjectionPushdown{SeriesHashLabel: "__series_hash__"},
			logicalplan.DetectHistogramStatsOptimizer{},
			// logicalplan.MergeSelectsOptimizer{},
		},
		DisableDuplicateLabelChecks: false,
	})

	ctx := context.Background()
	queryTime := time.Unix(600, 0)

	t.Logf("Running %d fuzzy tests with seed %d", testRuns, seed)
	for i := 0; i < testRuns; i++ {
		var expr parser.Expr
		var query string

		// Generate a query that can be executed by the engine
		for {
			expr = ps.WalkInstantQuery()
			query = expr.Pretty(0)

			// Skip queries that don't benefit from projection pushdown
			if !containsAggregationOrBinaryOperation(expr) {
				continue
			}

			// Try to parse the query and see if it is valid.
			_, err := normalEngine.NewInstantQuery(ctx, storage, nil, query, queryTime)
			if err != nil {
				continue
			}
			break
		}

		t.Run(fmt.Sprintf("Query_%d", i), func(t *testing.T) {
			// Create projection querier that wraps the original querier
			projectionStorage := &projectionQueryable{
				Queryable: storage,
			}

			normalQuery, err := normalEngine.NewInstantQuery(ctx, storage, &engine.QueryOpts{}, query, queryTime)
			testutil.Ok(t, err)
			defer normalQuery.Close()
			normalResult := normalQuery.Exec(ctx)
			if normalResult.Err != nil {
				// Something wrong with the generated query so it even failed without projection pushdown, skipping.
				return
			}
			testutil.Ok(t, normalResult.Err, "query: %s", query)

			pushdownQuery, err := pushdownEngine.MakeInstantQuery(ctx, projectionStorage, &engine.QueryOpts{}, query, queryTime)
			testutil.Ok(t, err)

			defer pushdownQuery.Close()
			pushdownResult := pushdownQuery.Exec(ctx)
			testutil.Ok(t, pushdownResult.Err, "query: %s", query)

			if diff := cmp.Diff(normalResult, pushdownResult, comparer); diff != "" {
				t.Errorf("Results differ for query %s: %s", query, diff)
			}
		})
	}
}

// containsAggregationOrBinaryOperation checks if the expression contains any aggregation or binary operations.
func containsAggregationOrBinaryOperation(expr parser.Expr) bool {
	found := false
	parser.Inspect(expr, func(node parser.Node, _ []parser.Node) error {
		switch node.(type) {
		case *parser.AggregateExpr:
			found = true
			return errors.New("found")
		case *parser.BinaryExpr:
			found = true
			return errors.New("found")
		}
		return nil
	})
	return found
}
