// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"math/rand"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/logicalplan"
	"github.com/thanos-io/promql-engine/query"

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
	if !m.hints.ProjectionInclude && len(m.hints.ProjectionLabels) == 0 {
		return originalSeries
	}

	// Apply projection based on hints
	originalLabels := originalSeries.Labels()
	var projectedLabels labels.Labels

	if m.hints.ProjectionInclude {
		// Include mode: only keep the labels in the projection labels
		builder := labels.NewBuilder(labels.EmptyLabels())
		originalLabels.Range(func(l labels.Label) {
			if slices.Contains(m.hints.ProjectionLabels, l.Name) {
				builder.Set(l.Name, l.Value)
			}
		})
		builder.Set("__series_hash__", strconv.FormatUint(originalLabels.Hash(), 10))
		projectedLabels = builder.Labels()
	} else {
		// Exclude mode: keep all labels except those in the projection labels
		excludeMap := make(map[string]struct{})
		for _, groupLabel := range m.hints.ProjectionLabels {
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

func TestProjectionWithFuzz(t *testing.T) {
	t.Parallel()

	// Define test parameters
	seed := time.Now().UnixNano()
	rnd := rand.New(rand.NewSource(seed))
	testRuns := 10000

	// Create test data
	load := `load 30s
		http_requests_total{pod="nginx-1", job="app", env="prod", instance="1"} 1+1x40
		http_requests_total{pod="nginx-2", job="app", env="dev", instance="2"} 2+2x40
		http_requests_total{pod="nginx-3", job="api", env="prod", instance="3"} 3+3x40
		http_requests_total{pod="nginx-4", job="api", env="dev", instance="4"} 4+4x40
		http_requests_duration_seconds_bucket{pod="nginx-1", job="app", env="prod", instance="1", le="0.1"} 1+1x40
		http_requests_duration_seconds_bucket{pod="nginx-1", job="app", env="prod", instance="1", le="0.2"} 2+2x40
		http_requests_duration_seconds_bucket{pod="nginx-1", job="app", env="prod", instance="1", le="0.5"} 3+2x40
		http_requests_duration_seconds_bucket{pod="nginx-1", job="app", env="prod", instance="1", le="+Inf"} 4+2x40
		http_requests_duration_seconds_bucket{pod="nginx-2", job="api", env="dev", instance="2", le="0.1"} 1+1x40
		http_requests_duration_seconds_bucket{pod="nginx-2", job="api", env="dev", instance="2", le="0.2"} 2+2x40
		http_requests_duration_seconds_bucket{pod="nginx-2", job="api", env="dev", instance="2", le="0.5"} 3+2x40
		http_requests_duration_seconds_bucket{pod="nginx-2", job="api", env="dev", instance="2", le="+Inf"} 4+2x40
		errors_total{pod="nginx-1", job="app", env="prod", instance="1", cluster="us-west-2"} 0.5+0.5x40
		errors_total{pod="nginx-2", job="app", env="dev", instance="2", cluster="us-west-2"} 1+1x40
		errors_total{pod="nginx-3", job="api", env="prod", instance="3", cluster="us-east-2"} 1.5+1.5x40
		errors_total{pod="nginx-4", job="api", env="dev", instance="4", cluster="us-east-1"} 2+2x40`

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	// Get series for PromQLSmith
	seriesSet, err := getSeries(context.Background(), storage, "http_requests_total")
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

	projectionEngine := engine.New(engine.Opts{
		EngineOpts: engineOpts,
		// projection optimizer doesn't support merge selects optimizer
		// so disable it for now.
		LogicalOptimizers: []logicalplan.Optimizer{
			logicalplan.SortMatchers{},
			logicalplan.ProjectionOptimizer{SeriesHashLabel: "__series_hash__"},
			logicalplan.DetectHistogramStatsOptimizer{},
			logicalplan.MergeSelectsOptimizer{},
		},
		DisableDuplicateLabelChecks: false,
	})

	ctx := context.Background()
	queryTime := time.Unix(600, 0)

	t.Logf("Running %d fuzzy tests with seed %d", testRuns, seed)
	for i := range testRuns {
		var expr parser.Expr
		var query string

		// Generate a query that can be executed by the engine
		for {
			expr = ps.WalkInstantQuery()
			query = expr.Pretty(0)

			// Skip queries that don't benefit from projection pushdown
			if !containsProjectionExprs(expr) {
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

			projectionQuery, err := projectionEngine.MakeInstantQuery(ctx, projectionStorage, &engine.QueryOpts{}, query, queryTime)
			testutil.Ok(t, err)

			defer projectionQuery.Close()
			projectionResult := projectionQuery.Exec(ctx)
			testutil.Ok(t, projectionResult.Err, "query: %s", query)

			if diff := cmp.Diff(normalResult, projectionResult, comparer); diff != "" {
				t.Errorf("Results differ for query %s: %s", query, diff)
			}
		})
	}
}

// optimizedPlanHasBinaryWithProjection returns true if, after running the projection
// optimizer with PushDownBinaryProjection enabled, the logical plan has at least one
// Binary node with Projection set. This covers all cases where projection is pushed
// down to a binary (aggregation over binary, outer binary over inner binary, etc.).
func optimizedPlanHasBinaryWithProjection(queryStr string) bool {
	expr, err := parser.ParseExpr(queryStr)
	if err != nil {
		return false
	}
	// Skip if the query has nothing that can benefit from projection pushdown.
	if !containsProjectionExprs(expr) {
		return false
	}
	// Must contain at least one binary expression to push projection down to.
	if !containsBinaryExpr(expr) {
		return false
	}
	opts := &query.Options{Start: time.Unix(0, 0), End: time.Unix(0, 0)}
	plan, err := logicalplan.NewFromAST(expr, opts, logicalplan.PlanOptions{})
	if err != nil {
		return false
	}
	optimizer := logicalplan.ProjectionOptimizer{
		SeriesHashLabel:          "__series_hash__",
		PushDownBinaryProjection: true,
	}
	optimizedRoot, _ := optimizer.Optimize(plan.Root(), nil)
	return planHasBinaryWithProjection(optimizedRoot)
}

func containsBinaryExpr(expr parser.Expr) bool {
	var found bool
	parser.Inspect(expr, func(node parser.Node, _ []parser.Node) error {
		if _, ok := node.(*parser.BinaryExpr); ok {
			found = true
			return errors.New("stop")
		}
		return nil
	})
	return found
}

func planHasBinaryWithProjection(node logicalplan.Node) bool {
	if node == nil {
		return false
	}
	if b, ok := node.(*logicalplan.Binary); ok && b.Projection != nil {
		return true
	}
	for _, c := range node.Children() {
		if c != nil && planHasBinaryWithProjection(*c) {
			return true
		}
	}
	return false
}

// TestProjectionPushdownAggregationWithBinary runs generated queries and only exercises
// those where the logical plan optimizer actually pushes projection to a Binary node (i.e.
// the optimized plan has at least one Binary with Projection set). This covers aggregation
// over binary, outer binary over inner binary, and any other case where pushdown occurs.
// The projection engine has PushDownBinaryProjection enabled.
func TestProjectionPushdownAggregationWithBinary(t *testing.T) {
	t.Parallel()

	seed := time.Now().UnixNano()
	rnd := rand.New(rand.NewSource(seed))
	const testRuns = 10000

	load := `load 30s
		http_requests_total{pod="nginx-1", job="app", env="prod", instance="1"} 1+1x40
		http_requests_total{pod="nginx-2", job="app", env="dev", instance="2"} 2+2x40
		http_requests_total{pod="nginx-3", job="api", env="prod", instance="3"} 3+3x40
		http_requests_total{pod="nginx-4", job="api", env="dev", instance="4"} 4+4x40
		http_requests_duration_seconds_bucket{pod="nginx-1", job="app", env="prod", instance="1", le="0.1"} 1+1x40
		http_requests_duration_seconds_bucket{pod="nginx-1", job="app", env="prod", instance="1", le="0.2"} 2+2x40
		http_requests_duration_seconds_bucket{pod="nginx-1", job="app", env="prod", instance="1", le="0.5"} 3+2x40
		http_requests_duration_seconds_bucket{pod="nginx-1", job="app", env="prod", instance="1", le="+Inf"} 4+2x40
		http_requests_duration_seconds_bucket{pod="nginx-2", job="api", env="dev", instance="2", le="0.1"} 1+1x40
		http_requests_duration_seconds_bucket{pod="nginx-2", job="api", env="dev", instance="2", le="0.2"} 2+2x40
		http_requests_duration_seconds_bucket{pod="nginx-2", job="api", env="dev", instance="2", le="0.5"} 3+2x40
		http_requests_duration_seconds_bucket{pod="nginx-2", job="api", env="dev", instance="2", le="+Inf"} 4+2x40
		errors_total{pod="nginx-1", job="app", env="prod", instance="1", cluster="us-west-2"} 0.5+0.5x40
		errors_total{pod="nginx-2", job="app", env="dev", instance="2", cluster="us-west-2"} 1+1x40
		errors_total{pod="nginx-3", job="api", env="prod", instance="3", cluster="us-east-2"} 1.5+1.5x40
		errors_total{pod="nginx-4", job="api", env="dev", instance="4", cluster="us-east-1"} 2+2x40`

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	seriesSet, err := getSeries(context.Background(), storage, "http_requests_total")
	testutil.Ok(t, err)

	// Restrict to aggregation, binary, and vector selectors to get aggregation-over-binary queries.
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

	normalEngine := engine.New(engine.Opts{
		EngineOpts:                  engineOpts,
		LogicalOptimizers:           logicalplan.AllOptimizers,
		DisableDuplicateLabelChecks: false,
	})

	// Projection engine with PushDownBinaryProjection enabled to exercise aggregation-to-binary pushdown.
	projectionEngine := engine.New(engine.Opts{
		EngineOpts: engineOpts,
		LogicalOptimizers: []logicalplan.Optimizer{
			logicalplan.SortMatchers{},
			logicalplan.ProjectionOptimizer{
				SeriesHashLabel:          "__series_hash__",
				PushDownBinaryProjection: true,
			},
			logicalplan.DetectHistogramStatsOptimizer{},
			logicalplan.MergeSelectsOptimizer{},
		},
		DisableDuplicateLabelChecks: false,
	})

	ctx := context.Background()
	queryTime := time.Unix(600, 0)
	projectionStorage := &projectionQueryable{Queryable: storage}

	t.Logf("Running %d tests (queries where optimized plan has Binary with projection) with seed %d", testRuns, seed)
	for i := range testRuns {
		var query string

		for {
			expr := ps.WalkInstantQuery()
			query = expr.Pretty(0)
			if !optimizedPlanHasBinaryWithProjection(query) {
				continue
			}
			_, err := normalEngine.NewInstantQuery(ctx, storage, nil, query, queryTime)
			if err != nil {
				continue
			}
			break
		}

		t.Run(fmt.Sprintf("Query_%d", i), func(t *testing.T) {
			normalQuery, err := normalEngine.NewInstantQuery(ctx, storage, &engine.QueryOpts{}, query, queryTime)
			testutil.Ok(t, err)
			defer normalQuery.Close()
			normalResult := normalQuery.Exec(ctx)
			if normalResult.Err != nil {
				return
			}

			projectionQuery, err := projectionEngine.MakeInstantQuery(ctx, projectionStorage, &engine.QueryOpts{}, query, queryTime)
			testutil.Ok(t, err)
			defer projectionQuery.Close()
			projectionResult := projectionQuery.Exec(ctx)
			testutil.Ok(t, projectionResult.Err, "query: %s", query)

			if diff := cmp.Diff(normalResult, projectionResult, comparer); diff != "" {
				t.Errorf("Results differ for query %s: %s", query, diff)
			}
		})
	}
}

// containsProjectionExprs checks if the expression contains any expressions that might benefit from projection pushdown.
func containsProjectionExprs(expr parser.Expr) bool {
	found := false
	parser.Inspect(expr, func(node parser.Node, _ []parser.Node) error {
		switch n := node.(type) {
		case *parser.Call:
			if n.Func.Name == "histogram_quantile" || n.Func.Name == "absent_over_time" || n.Func.Name == "absent" || n.Func.Name == "scalar" {
				found = true
				return errors.New("found")
			}
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

func TestBinaryProjectionPushdown(t *testing.T) {
	load := `
		load 30s
			kube_pod_info{cluster="c1", node="n1", namespace="ns1", pod="p1", label_a="a1", label_b="b1", label_c="c1"} 1
			kube_pod_info{cluster="c1", node="n2", namespace="ns2", pod="p2", label_a="a2", label_b="b2", label_c="c2"} 1
			kube_pod_info{cluster="c2", node="n3", namespace="ns3", pod="p3", label_a="a3", label_b="b3", label_c="c3"} 1
			kube_node_labels{cluster="c1", node="n1", instance_type="t1"} 1
			kube_node_labels{cluster="c1", node="n2", instance_type="t2"} 1
			kube_node_labels{cluster="c2", node="n3", instance_type="t3"} 1
	`

	cases := []struct {
		name  string
		query string
	}{
		{
			name: "outer aggregation with group_left join",
			query: `
sum by (namespace, instance_type) (
  kube_pod_info * on (cluster, node) group_left (instance_type) kube_node_labels
)`,
		},
		{
			name: "outer aggregation with group_right join",
			query: `
sum by (namespace, instance_type) (
  kube_node_labels * on (cluster, node) group_right (namespace) kube_pod_info
)`,
		},
		{
			name: "nested binary with outer aggregation",
			query: `
sum by (namespace) (
  (kube_pod_info * on (cluster, node) group_left (instance_type) kube_node_labels) + 1
)`,
		},
		{
			name: "outer binary operation",
			query: `
  (kube_pod_info * on (cluster, node) group_left (instance_type) kube_node_labels)
+ on (namespace) group_left ()
  kube_pod_info`,
		},
		{
			name:  "one-to-one join with outer aggregation",
			query: `sum by (cluster) (kube_pod_info * on (cluster, node) kube_node_labels)`,
		},
		{
			name: "avg by with group_left join",
			query: `
avg by (namespace, instance_type) (
  kube_pod_info * on (cluster, node) group_left (instance_type) kube_node_labels
)`,
		},
		{
			name:  "max by with one-to-one join",
			query: `max by (cluster, node) (kube_pod_info * on (cluster, node) kube_node_labels)`,
		},
		{
			name: "aggregation over outer binary (binary op binary)",
			query: `
sum by (namespace, instance_type) (
    (kube_pod_info * on (cluster, node) group_left (instance_type) kube_node_labels)
  + on (namespace) group_left ()
    kube_pod_info
)`,
		},
	}

	opts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	start := time.Unix(0, 0)
	end := time.Unix(120, 0)
	interval := 30 * time.Second

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Test with projection pushdown disabled (baseline)
			engineBaseline := engine.New(engine.Opts{
				EngineOpts: opts,
				LogicalOptimizers: []logicalplan.Optimizer{
					logicalplan.ProjectionOptimizer{
						SeriesHashLabel:          "__series_hash__",
						PushDownBinaryProjection: false,
					},
				},
			})

			qBaseline, err := engineBaseline.NewRangeQuery(context.Background(), storage, nil, tc.query, start, end, interval)
			testutil.Ok(t, err)
			resultBaseline := qBaseline.Exec(context.Background())
			testutil.Ok(t, resultBaseline.Err)

			// Test with projection pushdown enabled
			engineOptimized := engine.New(engine.Opts{
				EngineOpts: opts,
				LogicalOptimizers: []logicalplan.Optimizer{
					logicalplan.ProjectionOptimizer{
						SeriesHashLabel:          "__series_hash__",
						PushDownBinaryProjection: true,
					},
				},
			})

			qOptimized, err := engineOptimized.NewRangeQuery(context.Background(), storage, nil, tc.query, start, end, interval)
			testutil.Ok(t, err)
			resultOptimized := qOptimized.Exec(context.Background())
			testutil.Ok(t, resultOptimized.Err)

			// Results should be identical
			testutil.Equals(t, resultBaseline.String(), resultOptimized.String())
		})
	}
}

func TestBinaryProjectionPushdownWithPrometheus(t *testing.T) {
	load := `
		load 30s
			kube_pod_info{cluster="c1", node="n1", namespace="ns1", pod="p1", label_a="a1", label_b="b1"} 1
			kube_pod_info{cluster="c1", node="n2", namespace="ns2", pod="p2", label_a="a2", label_b="b2"} 1
			kube_node_labels{cluster="c1", node="n1", instance_type="t1"} 1
			kube_node_labels{cluster="c1", node="n2", instance_type="t2"} 1
	`

	queries := []string{
		`sum by (namespace, instance_type) (kube_pod_info * on (cluster, node) group_left(instance_type) kube_node_labels)`,
		`sum by (cluster) (kube_pod_info * on (cluster, node) kube_node_labels)`,
		`sum by (namespace) ((kube_pod_info * on (cluster, node) group_left(instance_type) kube_node_labels) + 1)`,
	}

	opts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	start := time.Unix(0, 0)
	end := time.Unix(120, 0)
	interval := 30 * time.Second

	promEngine := promql.NewEngine(opts)
	thanosEngine := engine.New(engine.Opts{
		EngineOpts: opts,
		LogicalOptimizers: []logicalplan.Optimizer{
			logicalplan.ProjectionOptimizer{
				SeriesHashLabel:          "__series_hash__",
				PushDownBinaryProjection: true,
			},
		},
	})

	for i, query := range queries {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			qProm, err := promEngine.NewRangeQuery(context.Background(), storage, nil, query, start, end, interval)
			testutil.Ok(t, err)
			resultProm := qProm.Exec(context.Background())
			testutil.Ok(t, resultProm.Err)

			qThanos, err := thanosEngine.NewRangeQuery(context.Background(), storage, nil, query, start, end, interval)
			testutil.Ok(t, err)
			resultThanos := qThanos.Exec(context.Background())
			testutil.Ok(t, resultThanos.Err)

			testutil.Equals(t, resultProm.String(), resultThanos.String())
		})
	}
}
