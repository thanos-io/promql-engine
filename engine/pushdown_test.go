package engine_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/cortexproject/promqlsmith"
	"github.com/efficientgo/core/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/promql/promqltest"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/util/annotations"
	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/logicalplan"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	pq "github.com/thanos-io/promql-engine/query"
)

// Mock storage.Querier that records the select hints it receives
type mockQuerier struct {
	storage.Querier
}

// Create a mock implementation that satisfies the storage.Querier interface
// but only records the hints it receives
func newMockQuerier(q storage.Querier) *mockQuerier {
	return &mockQuerier{
		Querier: q,
	}
}

// Mock SeriesSet that satisfies the interface but returns no series
type mockSeriesSet struct {
	storage.SeriesSet
	hints *storage.SelectHints
}

func (m mockSeriesSet) Next() bool { return m.SeriesSet.Next() }
func (m mockSeriesSet) At() storage.Series {
	// Get the original series
	originalSeries := m.SeriesSet.At()
	if originalSeries == nil {
		return nil
	}

	// If no projection hints, return the original series
	if m.hints == nil || len(m.hints.Grouping) == 0 {
		return originalSeries
	}

	// Apply projection based on hints
	originalLabels := originalSeries.Labels()
	var projectedLabels labels.Labels

	if m.hints.By {
		// Include mode: only keep the labels in the grouping
		builder := labels.NewBuilder(labels.EmptyLabels())
		for _, l := range originalLabels {
			for _, groupLabel := range m.hints.Grouping {
				if l.Name == groupLabel {
					builder.Set(l.Name, l.Value)
					break
				}
			}
		}
		builder.Set("__series_hash__", strconv.FormatUint(originalLabels.Hash(), 10))
		projectedLabels = builder.Labels()
	} else {
		// Exclude mode: keep all labels except those in the grouping
		excludeMap := make(map[string]struct{})
		for _, groupLabel := range m.hints.Grouping {
			excludeMap[groupLabel] = struct{}{}
		}

		builder := labels.NewBuilder(labels.EmptyLabels())
		for _, l := range originalLabels {
			if _, excluded := excludeMap[l.Name]; !excluded {
				builder.Set(l.Name, l.Value)
			}
		}
		builder.Set("__series_hash__", strconv.FormatUint(originalLabels.Hash(), 10))
		projectedLabels = builder.Labels()
	}

	// Return a projected series that wraps the original but with filtered labels
	return &projectedSeries{
		Series: originalSeries,
		lset:   projectedLabels,
	}
}

// projectedSeries wraps a storage.Series but returns projected labels
type projectedSeries struct {
	storage.Series
	lset labels.Labels
}

func (s *projectedSeries) Labels() labels.Labels {
	return s.lset
}

func (m mockSeriesSet) Err() error                        { return m.SeriesSet.Err() }
func (m mockSeriesSet) Warnings() annotations.Annotations { return m.SeriesSet.Warnings() }

// Implement the Querier interface methods
func (m *mockQuerier) Select(ctx context.Context, sortSeries bool, hints *storage.SelectHints, matchers ...*labels.Matcher) storage.SeriesSet {

	return mockSeriesSet{
		SeriesSet: m.Querier.Select(ctx, sortSeries, hints, matchers...),
		hints:     hints,
	}
}
func (m *mockQuerier) LabelValues(ctx context.Context, name string, _ *storage.LabelHints, matchers ...*labels.Matcher) ([]string, annotations.Annotations, error) {
	return nil, nil, nil
}
func (m *mockQuerier) LabelNames(ctx context.Context, _ *storage.LabelHints, matchers ...*labels.Matcher) ([]string, annotations.Annotations, error) {
	return nil, nil, nil
}
func (m *mockQuerier) Close() error { return nil }

// TestProjectionPushdownWithMockQuerier tests that projection information is correctly
// passed to the storage layer through SelectHints and produces the same results
func TestProjectionPushdownWithMockQuerier(t *testing.T) {
	cases := []struct {
		name  string
		load  string
		query string
	}{
		{
			name: "simple aggregation",
			load: `load 30s
				http_requests_total{pod="nginx-1", job="app", env="prod"} 1+1x15
				http_requests_total{pod="nginx-2", job="app", env="dev"} 1+2x18`,
			query: `sum (http_requests_total)`,
		},
		{
			name: "simple aggregation by",
			load: `load 30s
				http_requests_total{pod="nginx-1", job="app", env="prod"} 1+1x15
				http_requests_total{pod="nginx-2", job="app", env="dev"} 1+2x18`,
			query: `sum by (job) (http_requests_total)`,
		},
		{
			name: "simple aggregation without",
			load: `load 30s
				http_requests_total{pod="nginx-1", job="app", env="prod"} 1+1x15
				http_requests_total{pod="nginx-2", job="app", env="dev"} 1+2x18`,
			query: `sum without (env) (http_requests_total)`,
		},
		{
			name: "binary operation with vector matching",
			load: `load 30s
				metric1{pod="nginx-1", job="app", env="prod"} 1+1x15
				metric1{pod="nginx-2", job="app", env="dev"} 1+2x18
				metric2{pod="nginx-1", job="app", env="prod"} 2+1x15
				metric2{pod="nginx-2", job="app", env="dev"} 3+2x18`,
			query: `metric1 * on(job, env) metric2`,
		},
		{
			name: "complex query with multiple operations",
			load: `load 30s
				test_metric{pod="nginx-1", job="app", env="prod"} 1+1x15
				test_metric{pod="nginx-2", job="app", env="dev"} 1+2x18
				metric{pod="nginx-3", job="app", env="stage"} 2+1x15
				metric{pod="nginx-4", job="app", env="dev"} 3+2x18`,
			query: `sum by (job) (test_metric) / on(job) group_left count by (job) (metric)`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Load test data
			storage := promqltest.LoadedStorage(t, tc.load)
			defer storage.Close()

			ctx := context.Background()
			queryTime := time.Unix(60, 0) // Query at 60 seconds

			// Create engine with normal querier
			normalOpts := engine.Opts{
				EngineOpts: promql.EngineOpts{
					Timeout:              1 * time.Hour,
					MaxSamples:           1e10,
					EnableNegativeOffset: true,
					EnableAtModifier:     true,
				},
				LogicalOptimizers: []logicalplan.Optimizer{},
			}
			normalEngine := engine.New(normalOpts)

			// Create engine with projection pushdown optimizer
			pushdownOpts := engine.Opts{
				EngineOpts: promql.EngineOpts{
					Timeout:              1 * time.Hour,
					MaxSamples:           1e10,
					EnableNegativeOffset: true,
					EnableAtModifier:     true,
				},
				LogicalOptimizers: []logicalplan.Optimizer{
					&logicalplan.ProjectionPushdown{},
				},
			}
			pushdownEngine := engine.New(pushdownOpts)

			// Execute query with normal engine
			normalQuery, err := normalEngine.NewInstantQuery(ctx, storage, nil, tc.query, queryTime)
			testutil.Ok(t, err)
			defer normalQuery.Close()
			normalResult := normalQuery.Exec(ctx)
			testutil.Ok(t, normalResult.Err)

			// Create mock querier that wraps the original querier
			mockStorage := &mockQueryable{
				Queryable: storage,
			}

			// Execute query with pushdown engine and mock querier
			pushdownQuery, err := pushdownEngine.NewInstantQuery(ctx, mockStorage, nil, tc.query, queryTime)
			testutil.Ok(t, err)
			defer pushdownQuery.Close()
			pushdownResult := pushdownQuery.Exec(ctx)
			testutil.Ok(t, pushdownResult.Err)

			// Compare results
			testutil.WithGoCmp(comparer).Equals(t, normalResult, pushdownResult)
		})
	}
}

// mockQueryable is a storage.Queryable that applies projection to the querier
type mockQueryable struct {
	storage.Queryable
}

func (q *mockQueryable) Querier(mint, maxt int64) (storage.Querier, error) {
	querier, err := q.Queryable.Querier(mint, maxt)
	if err != nil {
		return nil, err
	}
	return &mockQuerier{
		Querier: querier,
	}, nil
}

// After the existing TestProjectionPushdown function

// FuzzProjectionPushdown tests the projection pushdown optimizer with randomly generated queries
// to ensure it produces the same results as without pushdown.
func TestProjectionPushdown(t *testing.T) {
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
		errors_total{pod="nginx-1", job="app", env="prod", instance="1"} 0.5+0.5x40
		errors_total{pod="nginx-2", job="app", env="dev", instance="2"} 1+1x40
		errors_total{pod="nginx-3", job="api", env="prod", instance="3"} 1.5+1.5x40
		errors_total{pod="nginx-4", job="api", env="dev", instance="4"} 2+2x40`

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	// Get series for PromQLSmith
	seriesSet, err := getSeries(context.Background(), storage)
	testutil.Ok(t, err)

	// Configure PromQLSmith
	psOpts := []promqlsmith.Option{
		promqlsmith.WithEnableOffset(true),
		promqlsmith.WithEnableAtModifier(true),
		// Focus on aggregations that benefit from projection pushdown
		promqlsmith.WithEnabledAggrs([]parser.ItemType{
			parser.SUM, parser.MIN, parser.MAX, parser.AVG, parser.COUNT,
		}),
	}
	ps := promqlsmith.New(rnd, seriesSet, psOpts...)

	// Engine options
	engineOpts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}

	// Create engines
	normalEngine := engine.New(engine.Opts{
		EngineOpts:        engineOpts,
		LogicalOptimizers: []logicalplan.Optimizer{},
	})

	pushdownEngine := engine.New(engine.Opts{
		EngineOpts:        engineOpts,
		LogicalOptimizers: append(logicalplan.AllOptimizers, &logicalplan.ProjectionPushdown{}),
	})

	// Test context and time
	ctx := context.Background()
	queryTime := time.Unix(600, 0)

	// Run tests
	t.Logf("Running %d fuzzy tests with seed %d", testRuns, seed)
	for i := 0; i < testRuns; i++ {
		var expr parser.Expr
		var query string

		// Generate a query that can be executed by the engine
		for {
			expr = ps.WalkInstantQuery()
			query = expr.Pretty(0)

			// Skip queries that don't benefit from projection pushdown
			if !containsAggregation(expr) {
				continue
			}

			// Skip queries that are too complex (depth > 4)
			if getExpressionDepth(expr) > 4 {
				continue
			}

			// Try to parse and execute the query
			_, err := normalEngine.NewInstantQuery(ctx, storage, nil, query, queryTime)
			if err != nil {
				continue
			}
			break
		}

		t.Run(fmt.Sprintf("Query_%d", i), func(t *testing.T) {
			query = `min without () ({__name__="http_requests_total"})`
			// Execute query with normal engine
			normalQuery, err := normalEngine.NewInstantQuery(ctx, storage, nil, query, queryTime)
			testutil.Ok(t, err)
			defer normalQuery.Close()
			normalResult := normalQuery.Exec(ctx)
			testutil.Ok(t, normalResult.Err)

			// Create mock querier that wraps the original querier
			mockStorage := &mockQueryable{
				Queryable: storage,
			}

			// Execute query with pushdown engine and mock querier
			pushdownQuery, err := pushdownEngine.NewInstantQuery(ctx, mockStorage, nil, query, queryTime)
			testutil.Ok(t, err)
			fmt.Println(query)

			plan := logicalplan.NewFromAST(expr, &pq.Options{Start: time.Unix(0, 0), End: time.Unix(0, 0)}, logicalplan.PlanOptions{})
			optimizer := logicalplan.ProjectionPushdown{}
			optimizedPlan, _ := optimizer.Optimize(plan.Root(), nil)
			fmt.Println(renderExprTree(optimizedPlan))

			defer pushdownQuery.Close()
			pushdownResult := pushdownQuery.Exec(ctx)
			if pushdownResult.Err != nil {
				fmt.Println("xxs")
			} else {
				testutil.Ok(t, pushdownResult.Err)
			}

			// Compare results
			if diff := cmp.Diff(normalResult, pushdownResult, comparer); diff != "" {
				t.Errorf("Results differ for query %s: %s", query, diff)
			}
		})
	}
}

// containsAggregation checks if the expression contains any aggregation operations
func containsAggregation(expr parser.Expr) bool {
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

// getExpressionDepth calculates the maximum depth of a PromQL expression tree
func getExpressionDepth(expr parser.Expr) int {
	maxDepth := 0
	var calculateDepth func(node parser.Node, currentDepth int)

	calculateDepth = func(node parser.Node, currentDepth int) {
		if node == nil {
			return
		}

		if currentDepth > maxDepth {
			maxDepth = currentDepth
		}

		switch n := node.(type) {
		case *parser.AggregateExpr:
			calculateDepth(n.Expr, currentDepth+1)
		case *parser.BinaryExpr:
			calculateDepth(n.LHS, currentDepth+1)
			calculateDepth(n.RHS, currentDepth+1)
		case *parser.Call:
			for _, arg := range n.Args {
				calculateDepth(arg, currentDepth+1)
			}
		case *parser.ParenExpr:
			calculateDepth(n.Expr, currentDepth+1)
		case *parser.UnaryExpr:
			calculateDepth(n.Expr, currentDepth+1)
		case *parser.SubqueryExpr:
			calculateDepth(n.Expr, currentDepth+1)
		}
	}

	calculateDepth(expr, 0)
	return maxDepth
}

// renderExprTree renders the expression into a string. It is useful
// in tests to use strings for assertions in cases where the "String()"
// method might not yield enough information or would panic because of
// internal logical expression types. Implementations were largeley taken
// from upstream prometheus.
//
// TODO: maybe its better to Traverse the expression here and inject
// new nodes with prepared String methods? Like replacing MatrixSelector
// by testMatrixSelector that has a overridden string method?
func renderExprTree(expr logicalplan.Node) string {
	switch t := expr.(type) {
	case *logicalplan.NumberLiteral:
		return fmt.Sprint(t.Val)
	case *logicalplan.VectorSelector:
		var b strings.Builder
		base := t.VectorSelector.String()
		if t.BatchSize > 0 {
			base += fmt.Sprintf("[batch=%d]", t.BatchSize)
		}
		if t.Projection.Labels != nil {
			sort.Strings(t.Projection.Labels)
			if t.Projection.Include {
				base += fmt.Sprintf("[projection=include(%s)]", strings.Join(t.Projection.Labels, ","))
			} else {
				base += fmt.Sprintf("[projection=exclude(%s)]", strings.Join(t.Projection.Labels, ","))
			}
		}
		if len(t.Filters) > 0 {
			b.WriteString("filter(")
			b.WriteString(fmt.Sprintf("%s", t.Filters))
			b.WriteString(", ")
			b.WriteString(base)
			b.WriteRune(')')
			return b.String()
		}
		return base
	case *logicalplan.MatrixSelector:
		// Render the inner vector selector first
		vsStr := renderExprTree(t.VectorSelector)
		// Then add the range
		return fmt.Sprintf("%s[%s]", vsStr, t.Range.String())
	case *logicalplan.Binary:
		var b strings.Builder
		b.WriteString(renderExprTree(t.LHS))
		b.WriteString(" ")
		b.WriteString(t.Op.String())
		b.WriteString(" ")
		if vm := t.VectorMatching; vm != nil && (len(vm.MatchingLabels) > 0 || vm.On) {
			vmTag := "ignoring"
			if vm.On {
				vmTag = "on"
			}
			matching := fmt.Sprintf("%s (%s)", vmTag, strings.Join(vm.MatchingLabels, ", "))

			if vm.Card == parser.CardManyToOne || vm.Card == parser.CardOneToMany {
				vmCard := "right"
				if vm.Card == parser.CardManyToOne {
					vmCard = "left"
				}
				matching += fmt.Sprintf(" group_%s (%s)", vmCard, strings.Join(vm.Include, ", "))
			}
			b.WriteString(matching)
			b.WriteString(" ")
		}
		b.WriteString(renderExprTree(t.RHS))
		return b.String()
	case *logicalplan.FunctionCall:
		var b strings.Builder
		b.Write([]byte(t.Func.Name))
		b.WriteRune('(')
		for i := range t.Args {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(renderExprTree(t.Args[i]))
		}
		b.WriteRune(')')
		return b.String()
	case *logicalplan.Aggregation:
		var b strings.Builder
		b.Write([]byte(t.Op.String()))
		switch {
		case t.Without:
			b.WriteString(fmt.Sprintf(" without (%s) ", strings.Join(t.Grouping, ", ")))
		case len(t.Grouping) > 0:
			b.WriteString(fmt.Sprintf(" by (%s) ", strings.Join(t.Grouping, ", ")))
		}
		b.WriteRune('(')
		if t.Param != nil {
			b.WriteString(renderExprTree(t.Param))
			b.WriteString(", ")
		}
		b.WriteString(renderExprTree(t.Expr))
		b.WriteRune(')')
		return b.String()
	case *logicalplan.StepInvariantExpr:
		return renderExprTree(t.Expr)
	case *logicalplan.CheckDuplicateLabels:
		return renderExprTree(t.Expr)
	case *logicalplan.Unary:
		return renderExprTree(t.Expr)
	case *logicalplan.Parens:
		return renderExprTree(t.Expr)
	case *logicalplan.Subquery:
		var b strings.Builder

		// Render the inner expression
		innerExpr := renderExprTree(t.Expr)
		b.WriteString(innerExpr)

		// Add the subquery range and step
		b.WriteString(fmt.Sprintf("[%s:%s]", t.Range.String(), t.Step.String()))
		return b.String()
	default:
		return t.String()
	}
}
