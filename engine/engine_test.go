// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"math"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"

	promparser "github.com/prometheus/prometheus/promql/parser"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/logicalplan"

	"github.com/efficientgo/core/testutil"
	"github.com/go-kit/log"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/prometheus/model/histogram"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/timestamp"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	"github.com/prometheus/prometheus/tsdb/tsdbutil"
	"github.com/prometheus/prometheus/util/stats"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestQueryExplain(t *testing.T) {
	opts := promql.EngineOpts{Timeout: 1 * time.Hour}
	series := storage.MockSeries(
		[]int64{240, 270, 300, 600, 630, 660},
		[]float64{1, 2, 3, 4, 5, 6},
		[]string{labels.MetricName, "foo"},
	)

	start := time.Unix(0, 0)
	end := time.Unix(1000, 0)

	// Calculate concurrencyOperators according to max available CPUs.
	totalOperators := runtime.GOMAXPROCS(0) / 2
	concurrencyOperators := []engine.ExplainOutputNode{}
	for i := 0; i < totalOperators; i++ {
		concurrencyOperators = append(concurrencyOperators, engine.ExplainOutputNode{
			OperatorName: "[*concurrencyOperator(buff=2)]", Children: []engine.ExplainOutputNode{
				{OperatorName: fmt.Sprintf("[*vectorSelector] {[__name__=\"foo\"]} %d mod %d", i, totalOperators)},
			},
		})
	}

	for _, tc := range []struct {
		query    string
		expected *engine.ExplainOutputNode
	}{
		{
			query:    "time()",
			expected: &engine.ExplainOutputNode{OperatorName: "[*noArgFunctionOperator] time()"},
		},
		{
			query:    "foo",
			expected: &engine.ExplainOutputNode{OperatorName: "[*coalesce]", Children: concurrencyOperators},
		},
		{
			query: "sum(foo) by (job)",
			expected: &engine.ExplainOutputNode{OperatorName: "[*concurrencyOperator(buff=2)]", Children: []engine.ExplainOutputNode{
				{OperatorName: "[*aggregate] sum by ([job])", Children: []engine.ExplainOutputNode{
					{OperatorName: "[*coalesce]", Children: concurrencyOperators},
				},
				},
			},
			},
		},
	} {
		{
			t.Run(tc.query, func(t *testing.T) {
				ng := engine.New(engine.Opts{EngineOpts: opts})
				ctx := context.Background()

				var (
					query promql.Query
					err   error
				)

				query, err = ng.NewInstantQuery(ctx, storageWithSeries(series), nil, tc.query, start)
				testutil.Ok(t, err)

				explainableQuery := query.(engine.ExplainableQuery)
				testutil.Equals(t, tc.expected, explainableQuery.Explain())

				query, err = ng.NewRangeQuery(ctx, storageWithSeries(series), nil, tc.query, start, end, 30*time.Second)
				testutil.Ok(t, err)

				explainableQuery = query.(engine.ExplainableQuery)
				testutil.Equals(t, tc.expected, explainableQuery.Explain())
			})
		}
	}
}

func TestVectorSelectorWithGaps(t *testing.T) {
	opts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}

	series := storage.MockSeries(
		[]int64{240, 270, 300, 600, 630, 660},
		[]float64{1, 2, 3, 4, 5, 6},
		[]string{labels.MetricName, "foo"},
	)

	query := "foo"
	start := time.Unix(0, 0)
	end := time.Unix(1000, 0)

	ctx := context.Background()
	newEngine := engine.New(engine.Opts{EngineOpts: opts})
	q1, err := newEngine.NewRangeQuery(ctx, storageWithSeries(series), nil, query, start, end, 30*time.Second)
	testutil.Ok(t, err)
	defer q1.Close()

	newResult := q1.Exec(ctx)
	testutil.Ok(t, newResult.Err)

	oldEngine := promql.NewEngine(opts)
	q2, err := oldEngine.NewRangeQuery(ctx, storageWithSeries(series), nil, query, start, end, 30*time.Second)
	testutil.Ok(t, err)
	defer q2.Close()

	oldResult := q2.Exec(context.Background())
	testutil.Ok(t, oldResult.Err)

	testutil.Equals(t, oldResult, newResult)

}

func TestQueriesAgainstOldEngine(t *testing.T) {
	start := time.Unix(0, 0)
	end := time.Unix(1800, 0)
	step := time.Second * 30
	// Negative offset and at modifier are enabled by default
	// since Prometheus v2.33.0 so we also enable them.
	opts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		Logger:               log.NewLogfmtLogger(os.Stderr),
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}

	cases := []struct {
		load  string
		name  string
		query string
		start time.Time
		end   time.Time
		step  time.Duration
	}{
		{
			name:  "nested unary negation",
			query: "1/(-(2*2))",
		},
		{
			name: "stddev with large values",
			load: `load 30s
              http_requests_total{pod="nginx-1", route="/"} 1e+181
              http_requests_total{pod="nginx-2", route="/"} 1e+80`,
			query: `stddev(http_requests_total)`,
		},
		{
			name: "stddev with NaN 1",
			load: `load 30s
				       http_requests_total{pod="nginx-1", route="/"} NaN
				       http_requests_total{pod="nginx-2", route="/"} 1`,
			query: "stddev by (route) (http_requests_total)",
		},
		{
			name: "stddev with NaN 2",
			load: `load 30s
				       http_requests_total{pod="nginx-1", route="/"} NaN
				       http_requests_total{pod="nginx-2", route="/"} 1`,
			query: "stddev by (pod) (http_requests_total)",
		},
		{
			name: "aggregate without",
			load: `load 30s
				       http_requests_total{pod="nginx-1"} 1+1.1x1
				       http_requests_total{pod="nginx-2"} 2+2.3x1`,
			start: time.Unix(0, 0),
			end:   time.Unix(60, 0),
			step:  30 * time.Second,
			query: "avg without (pod) (http_requests_total)",
		},
		{
			name: "func with scalar arg that selects storage, checks whether same series handled correctly",
			load: `load 30s
			    thanos_cache_redis_hits_total{name="caching-bucket",service="thanos-store"} 1+1x30`,
			query: `clamp_min(thanos_cache_redis_hits_total, scalar(max(thanos_cache_redis_hits_total) by (service))) + clamp_min(thanos_cache_redis_hits_total, scalar(max(thanos_cache_redis_hits_total) by (service)))`,
		},
		{
			name: "sum + rate divided by itself",
			load: `load 30s
			thanos_cache_redis_hits_total{name="caching-bucket",service="thanos-store"} 1+1x30`,
			query: `(sum(rate(thanos_cache_redis_hits_total{name="caching-bucket"}[2m])) by (service)) /
			(sum(rate(thanos_cache_redis_hits_total{name="caching-bucket"}[2m])) by (service))`,
		},
		{
			name: "stddev_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "stddev_over_time(http_requests_total[30s])",
		},
		{
			name: "stdvar_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "stdvar_over_time(http_requests_total[30s])",
		},
		{
			name: "changes",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "changes(http_requests_total[30s])",
		},
		{
			name: "deriv",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "deriv(http_requests_total[30s])",
		},
		{
			name: "abs",
			load: `load 30s
					http_requests_total{pod="nginx-1"} -5+1x15
					http_requests_total{pod="nginx-2"} -5+2x18`,
			query: "abs(http_requests_total)",
		},
		{
			name: "ceil",
			load: `load 30s
					http_requests_total{pod="nginx-1"} -5.5+1x15
					http_requests_total{pod="nginx-2"} -5.5+2x18`,
			query: "ceil(http_requests_total)",
		},
		{
			name: "exp",
			load: `load 30s
					http_requests_total{pod="nginx-1"} -5.5+1x15
					http_requests_total{pod="nginx-2"} -5.5+2x18`,
			query: "exp(http_requests_total)",
		},
		{
			name: "floor",
			load: `load 30s
					http_requests_total{pod="nginx-1"} -5.5+1x15
					http_requests_total{pod="nginx-2"} -5.5+2x18`,
			query: "floor(http_requests_total)",
		},
		{
			name: "floor with a filter",
			load: `load 30s
          http_requests_total{pod="nginx-1", route="/"} 1
          http_requests_total{pod="nginx-2", route="/"} 2
`,
			query: `floor(http_requests_total{pod="nginx-2"})/http_requests_total`,
		},
		{
			name: "sqrt",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 5.5+1x15
					http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: "sqrt(http_requests_total)",
		},
		{
			name: "ln",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 5.5+1x15
					http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: "ln(http_requests_total)",
		},
		{
			name: "log2",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 5.5+1x15
					http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: "log2(http_requests_total)",
		},
		{
			name: "log10",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 5.5+1x15
					http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: "log10(http_requests_total)",
		},
		{
			name: "sin",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 5.5+1x15
					http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: "sin(http_requests_total)",
		},
		{
			name: "cos",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 5.5+1x15
					http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: "cos(http_requests_total)",
		},
		{
			name: "tan",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 5.5+1x15
					http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: "tan(http_requests_total)",
		},
		{
			name: "asin",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 0
					http_requests_total{pod="nginx-2"} 1`,
			query: "asin(http_requests_total)",
		},
		{
			name: "acos",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 0
					http_requests_total{pod="nginx-2"} 1`,
			query: "acos(http_requests_total)",
		},
		{
			name: "atan",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 0
					http_requests_total{pod="nginx-2"} 1`,
			query: "atan(http_requests_total)",
		},
		{
			name: "sinh",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 0
					http_requests_total{pod="nginx-2"} 1`,
			query: "sinh(http_requests_total)",
		},
		{
			name: "cosh",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 0
					http_requests_total{pod="nginx-2"} 1`,
			query: "cosh(http_requests_total)",
		},
		{
			name: "tanh",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 5.5+1x15
					http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: "tanh(http_requests_total)",
		},
		{
			name: "asinh",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 5.5+1x15
					http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: "asinh(http_requests_total)",
		},
		{
			name: "acosh",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 5.5+1x15
					http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: "acosh(http_requests_total)",
		},
		{
			name: "atanh",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 0
					http_requests_total{pod="nginx-2"} 1`,
			query: "atanh(http_requests_total)",
		},
		{
			name: "rad",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 5.5+1x15
					http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: "rad(http_requests_total)",
		},
		{
			name: "deg",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 5.5+1x15
					http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: "deg(http_requests_total)",
		},
		{
			name:  "pi",
			load:  ``,
			query: "pi()",
		},
		{
			name: "sum",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "sum (http_requests_total)",
		},
		{
			name: "sum_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "sum_over_time(http_requests_total[30s])",
		},
		{
			name: "count",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "count(http_requests_total)",
		},
		{
			name: "count_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "count_over_time(http_requests_total[30s])",
		},
		{
			name: "average",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "avg(http_requests_total)",
		},
		{
			name: "avg_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "avg_over_time(http_requests_total[30s])",
		},
		{
			name: "abs",
			load: `load 30s
					http_requests_total{pod="nginx-1"} -10+1x15
					http_requests_total{pod="nginx-2"} -10+2x18`,
			query: "abs(http_requests_total)",
		},
		{
			name: "max",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "max(http_requests_total)",
		},
		{
			name: "max with only 1 sample",
			load: `load 30s
					http_requests_total{pod="nginx-1"} -1
					http_requests_total{pod="nginx-2"} 1`,
			query: "max(http_requests_total) by (pod)",
		},
		{
			name: "max_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "max_over_time(http_requests_total[30s])",
		},
		{
			name: "min",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "min(http_requests_total)",
		},
		{
			name: "min with only 1 sample",
			load: `load 30s
					http_requests_total{pod="nginx-1"} -1
					http_requests_total{pod="nginx-2"} 1`,
			query: "min(http_requests_total) by (pod)",
		},
		{
			name: "min_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "min_over_time(http_requests_total[30s])",
		},
		{
			name: "count_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "count_over_time(http_requests_total[30s])",
		},
		{
			name: "sum by pod",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-3"} 1+2x20
					http_requests_total{pod="nginx-4"} 1+2x50`,
			query: "sum by (pod) (http_requests_total)",
		},
		{
			name: "multi label grouping by",
			load: `load 30s
					http_requests_total{pod="nginx-1", ns="a"} 1+1x15
					http_requests_total{pod="nginx-2", ns="a"} 1+1x15`,
			query: `avg by (pod, ns) (avg_over_time(http_requests_total[2m]))`,
		},
		{
			name: "multi label grouping without",
			load: `load 30s
					http_requests_total{pod="nginx-1", ns="a"} 1+1x15
					http_requests_total{pod="nginx-2", ns="a"} 1+1x15`,
			query: `avg without (pod, ns) (avg_over_time(http_requests_total[2m]))`,
		},
		{
			name: "query in the future",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "sum by (pod) (http_requests_total)",
			start: time.Unix(400, 0),
			end:   time.Unix(3000, 0),
		},
		{
			name: "count_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1
					http_requests_total{pod="nginx-1"} 1+1x30
					http_requests_total{pod="nginx-2"} 1+2x600`,
			query: `count_over_time(http_requests_total[10m])`,
			start: time.Unix(60, 0),
			end:   time.Unix(600, 0),
		},
		{
			name: "rate",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "rate(http_requests_total[1m])",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "sum rate",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x4
					http_requests_total{pod="nginx-2"} 1+2x4`,
			query: "sum(rate(http_requests_total[1m]))",
		},
		{
			name: "sum rate with stale series",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x40
					http_requests_total{pod="nginx-2"} 1+2x50
					http_requests_total{pod="nginx-4"} 1+2x50
					http_requests_total{pod="nginx-5"} 1+2x50
					http_requests_total{pod="nginx-6"} 1+2x50`,
			query: "sum(rate(http_requests_total[1m]))",
			start: time.Unix(421, 0),
			end:   time.Unix(3230, 0),
			step:  28 * time.Second,
		},
		{
			name: "delta",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "delta(http_requests_total[1m])",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "increase",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "increase(http_requests_total[1m])",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "irate",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "irate(http_requests_total[1m])",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "idelta",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "idelta(http_requests_total[1m])",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name:  "number literal",
			load:  "",
			query: "34",
		},
		{
			name:  "vector",
			load:  "",
			query: "vector(24)",
		},
		{
			name: "binary operation atan2",
			load: `load 30s
         foo{} 10
         bar{} 2`,
			query: "foo atan2 bar",
		},
		{
			name: "binary operation atan2 with NaN",
			load: `load 30s
         foo{} 10
         bar{} NaN`,
			query: "foo atan2 bar",
		},
		{
			name: "binary operation with one-to-one matching",
			load: `load 30s
				foo{method="get", code="500"} 1+1x1
				foo{method="get", code="404"} 2+1x2
				foo{method="put", code="501"} 3+1x3
				foo{method="put", code="500"} 1+1x4
				foo{method="post", code="500"} 4+1x4
				foo{method="post", code="404"} 5+1x5
				bar{method="get"} 1+1x1
				bar{method="del"} 2+1x2
				bar{method="post"} 3+1x3`,
			query: `foo{code="500"} + ignoring(code) bar`,
			start: time.Unix(0, 0),
			end:   time.Unix(600, 0),
		},
		{
			// Example from https://prometheus.io/docs/prometheus/latest/querying/operators/#many-to-one-and-one-to-many-vector-matches
			name: "binary operation with group_left",
			load: `load 30s
				foo{method="get", code="500", path="/"} 1+1.1x30
				foo{method="get", code="404", path="/"} 1+2.2x20
				foo{method="put", code="501", path="/"} 4+3.4x60
				foo{method="post", code="500", path="/"} 1+5.1x40
				foo{method="post", code="404", path="/"} 2+3.7x40
				bar{method="get", path="/a"} 3+7.4x10
				bar{method="del", path="/b"} 8+6.1x30
				bar{method="post", path="/c"} 1+2.1x40`,
			query: `foo * ignoring(path, code) group_left bar`,
			start: time.Unix(0, 0),
			end:   time.Unix(600, 0),
		},
		{
			// Example from https://prometheus.io/docs/prometheus/latest/querying/operators/#many-to-one-and-one-to-many-vector-matches
			name: "binary operation with group_right",
			load: `load 30s
				foo{method="get", code="500"} 1+1.1x30
				foo{method="get", code="404"} 1+2.2x20
				foo{method="put", code="501"} 4+3.4x60
				foo{method="post", code="500"} 1+5.1x40
				foo{method="post", code="404"} 2+3.7x40
				bar{method="get", path="/a"} 3+7.4x10
				bar{method="del", path="/b"} 8+6.1x30
				bar{method="post", path="/c"} 1+2.1x40`,
			query: `bar * ignoring(code, path) group_right foo`,
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "binary operation with group_left and included labels",
			load: `load 30s
				foo{method="get", code="500"} 1+1.1x30
				foo{method="get", code="404"} 1+2.2x20
				foo{method="put", code="501"} 4+3.4x60
				foo{method="post", code="500"} 1+5.1x40
				foo{method="post", code="404"} 2+3.7x40
				bar{method="get", path="/a"} 3+7.4x10
				bar{method="del", path="/b"} 8+6.1x30
				bar{method="post", path="/c"} 1+2.1x40`,
			query: `foo * ignoring(code, path) group_left(path) bar`,
			start: time.Unix(0, 0),
			end:   time.Unix(600, 0),
		},
		{
			name: "binary operation with group_right and included labels",
			load: `load 30s
				foo{method="get", code="500"} 1+1.1x30
				foo{method="get", code="404"} 1+2.2x20
				foo{method="put", code="501"} 4+3.4x60
				foo{method="post", code="500"} 1+5.1x40
				foo{method="post", code="404"} 2+3.7x40
				bar{method="get", path="/a"} 3+7.4x10
				bar{method="del", path="/b"} 8+6.1x30
				bar{method="post", path="/c"} 1+2.1x40`,
			query: `bar * ignoring(code, path) group_right(path) foo`,
			start: time.Unix(0, 0),
			end:   time.Unix(600, 0),
		},
		{
			name: "binary operation with vector and scalar on the right",
			load: `load 30s
				foo{method="get", code="500"} 1+1.1x30
				foo{method="get", code="404"} 1+2.2x20`,
			query: `sum(foo) * 2`,
		},
		{
			name: "binary operation with vector and scalar on the left",
			load: `load 30s
				foo{method="get", code="500"} 1+1.1x30
				foo{method="get", code="404"} 1+2.2x20`,
			query: `2 * sum(foo)`,
		},
		{
			name: "complex binary operation",
			load: `load 30s
				foo{method="get", code="500"} 1+1.1x30
				foo{method="get", code="404"} 1+2.2x20`,
			query: `1 - (100 * sum(foo{method="get"}) / sum(foo))`,
		},
		{
			name: "binary operation with many-to-many matching",
			load: `load 30s
				foo{code="200", method="get"} 1+1x20
				foo{code="200", method="post"} 1+1x20
				bar{code="200", method="get"} 1+1x20
				bar{code="200", method="post"} 1+1x20`,
			query: `foo + on(code) bar`,
		},
		{
			name: "binary operation with many-to-many matching lhs high card",
			load: `load 30s
				foo{code="200", method="get"} 1+1x20
				foo{code="200", method="post"} 1+1x20
				bar{code="200", method="get"} 1+1x20
				bar{code="200", method="post"} 1+1x20`,
			query: `foo + on(code) group_left bar`,
		},
		{
			name: "binary operation with many-to-many matching rhs high card",
			load: `load 30s
				foo{code="200", method="get"} 1+1x20
				foo{code="200", method="post"} 1+1x20
				bar{code="200", method="get"} 1+1x20
				bar{code="200", method="post"} 1+1x20`,
			query: `foo + on(code) group_right bar`,
		},
		{
			name: "vector binary op ==",
			load: `load 30s
				foo{method="get", code="500"} 1+1x40
				bar{method="get", code="404"} 1+1.1x30`,
			query: `sum(foo) by (method) == sum(bar) by (method)`,
		},
		{
			name: "vector binary op !=",
			load: `load 30s
				foo{method="get", code="500"} 1+1x40
				bar{method="get", code="404"} 1+1.1x30`,
			query: `sum(foo) by (method) != sum(bar) by (method)`,
		},
		{
			name: "vector binary op >",
			load: `load 30s
				foo{method="get", code="500"} 1+1x40
				bar{method="get", code="404"} 1+1.1x30`,
			query: `sum(foo) by (method) > sum(bar) by (method)`,
		},
		{
			name: "vector binary op with name <",
			load: `load 30s
				foo{method="get", code="500"} 1+1x40
				bar{method="get", code="500"} 1+1.1x30`,
			query: `foo < bar`,
		},
		{
			name: "vector binary op with name < scalar",
			load: `load 30s
				foo{method="get", code="500"} 1+1x40
				bar{method="get", code="500"} 1+1.1x30`,
			query: `foo < 10`,
		},
		{
			name: "vector binary op with name < scalar and bool modifier",
			load: `load 30s
				foo{method="get", code="500"} 1+1x40
				bar{method="get", code="500"} 1+1.1x30`,
			query: `foo < bool 10`,
		},
		{
			name: "vector binary op > scalar",
			load: `load 30s
				foo{method="get", code="500"} 1+2x40
				bar{method="get", code="404"} 1+1x30`,
			query: `sum(foo) by (method) > 10`,
		},
		{
			name: "vector binary op > scalar and bool modifier",
			load: `load 30s
				foo{method="get", code="500"} 1+2x40
				bar{method="get", code="404"} 1+1x30`,
			query: `sum(foo) by (method) > bool 10`,
		},
		{
			name: "scalar < vector binary op",
			load: `load 30s
				foo{method="get", code="500"} 1+2x40
				bar{method="get", code="404"} 1+1x30`,
			query: `10 < sum(foo) by (method)`,
		},
		{
			name: "vector binary op <",
			load: `load 30s
				foo{method="get", code="500"} 1+1x40
				bar{method="get", code="404"} 1+1.1x30`,
			query: `sum(foo) by (method) < sum(bar) by (method)`,
		},
		{
			name: "vector binary op >=",
			load: `load 30s
				foo{method="get", code="500"} 1+1x40
				bar{method="get", code="404"} 1+1.1x30`,
			query: `sum(foo) by (method) >= sum(bar) by (method)`,
		},
		{
			name: "vector binary op <=",
			load: `load 30s
				foo{method="get", code="500"} 1+1x40
				bar{method="get", code="404"} 1+1.1x30`,
			query: `sum(foo) by (method) <= sum(bar) by (method)`,
		},
		{
			name: "vector binary op ^",
			load: `load 30s
				foo{method="get", code="500"} 1+1x40
				bar{method="get", code="404"} 1+1.1x30`,
			query: `sum(foo) by (method) ^ sum(bar) by (method)`,
		},
		{
			name: "vector binary op %",
			load: `load 30s
				foo{method="get", code="500"} 1+2x40
				bar{method="get", code="404"} 1+1x30`,
			query: `sum(foo) by (method) % sum(bar) by (method)`,
		},
		{
			name: "vector/vector binary op",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "(1 + rate(http_requests_total[30s])) > bool rate(http_requests_total[30s])",
		},
		{
			name: "vector/scalar binary op with a complicated expression on LHS",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "rate(http_requests_total[30s]) > bool 0",
		},
		{
			name: "vector/scalar binary op with a complicated expression on RHS",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "0 < bool rate(http_requests_total[30s])",
		},
		{
			name:  "scalar binary op == true",
			load:  ``,
			query: `1 == bool 1`,
		},
		{
			name:  "scalar binary op == false",
			load:  ``,
			query: `1 != bool 2`,
		},
		{
			name:  "scalar binary op !=",
			load:  ``,
			query: `1 != bool 1`,
		},
		{
			name:  "scalar binary op >",
			load:  ``,
			query: `1 > bool 0`,
		},
		{
			name:  "scalar binary op <",
			load:  ``,
			query: `1 < bool 2`,
		},
		{
			name:  "scalar binary op >=",
			load:  ``,
			query: `1 >= bool 0`,
		},
		{
			name:  "scalar binary op <=",
			load:  ``,
			query: `1 <= bool 2`,
		},
		{
			name:  "scalar binary op % 0",
			load:  ``,
			query: `2 % 2`,
		},
		{
			name:  "scalar binary op % 1",
			load:  ``,
			query: `1 % 2`,
		},
		{
			name:  "scalar binary op ^",
			load:  ``,
			query: `2 ^ 2`,
		},
		{
			name:  "empty series",
			load:  "",
			query: "http_requests_total",
		},
		{
			name:  "time function",
			load:  "",
			query: "time()",
		},
		{
			name:  "time function in binary expression",
			load:  "",
			query: "time() - 10",
		},
		{
			name:  "empty series with func",
			load:  "",
			query: "sum(http_requests_total)",
		},
		{
			name: "empty result",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total{pod="nginx-3"}`,
		},
		{
			name: "last_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `last_over_time(http_requests_total[30s])`,
		},
		{
			name: "group",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `group(http_requests_total)`,
		},
		{
			name: "resets",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 100-1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `resets(http_requests_total[5m])`,
		},
		{
			name: "present_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `present_over_time(http_requests_total[30s])`,
		},
		{
			name: "complex binary with aggregation",
			load: `load 30s
					grpc_server_handled_total{pod="nginx-1", grpc_method="Series", grpc_code="105"} 1+1x15
					grpc_server_handled_total{pod="nginx-2", grpc_method="Series", grpc_code="105"} 1+1x15
					grpc_server_handled_total{pod="nginx-3", grpc_method="Series", grpc_code="105"} 1+1x15
					prometheus_tsdb_head_samples_appended_total{pod="nginx-1", tenant="tenant-1"} 1+2x18
					prometheus_tsdb_head_samples_appended_total{pod="nginx-2", tenant="tenant-2"} 1+2x18
					prometheus_tsdb_head_samples_appended_total{pod="nginx-3", tenant="tenant-3"} 1+2x18`,
			query: `
	sum by (grpc_method, grpc_code) (
		sum by (pod, grpc_method, grpc_code) (
			rate(grpc_server_handled_total{grpc_method="Series", pod=~".+"}[1m])
		)
		+ on (pod) group_left() max by (pod) (
			prometheus_tsdb_head_samples_appended_total{pod=~".+"}
		)
	)`,
		},
		{
			name:  "unary sub operation for scalar",
			load:  ``,
			query: `-(1 + 5)`,
		},
		{
			name: "unary sub operation for vector",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `-http_requests_total`,
		},
		{
			name: "unary add operation for vector",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `+http_requests_total`,
		},
		{
			name: "vector positive offset",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total offset 30s`,
			start: time.Unix(600, 0),
			end:   time.Unix(1200, 0),
		},
		{
			name: "vector negative offset",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total offset -30s`,
			start: time.Unix(600, 0),
			end:   time.Unix(1200, 0),
		},
		{
			name: "matrix negative offset with sum_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x25
					http_requests_total{pod="nginx-2"} 1+2x28`,
			query: `sum_over_time(http_requests_total[5m] offset 5m)`,
			start: time.Unix(600, 0),
			end:   time.Unix(6000, 0),
		},
		{
			name: "matrix negative offset with count_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `count_over_time(http_requests_total[5m] offset -2m)`,
			start: time.Unix(600, 0),
			end:   time.Unix(6000, 0),
		},
		{
			name: "@ vector time 10s",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "http_requests_total @ 10",
		},
		{
			name: "@ vector time 120s",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "http_requests_total @ 120",
		},
		{
			name: "@ vector time 360s",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "http_requests_total @ 360",
		},
		{
			name: "@ vector start",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "http_requests_total @ start()",
		},
		{
			name: "@ vector end",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "http_requests_total @ end()",
		},
		{
			name: "count_over_time @ start",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "count_over_time(http_requests_total[5m] @ start())",
		},
		{
			name: "sum_over_time @ end",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "sum_over_time(http_requests_total[5m] @ start())",
		},
		{
			name: "avg_over_time @ 180s",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "avg_over_time(http_requests_total[4m] @ 180)",
		},
		{
			name: "@ vector 240s offset 2m",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "http_requests_total @ 240 offset 2m",
		},
		{
			name: "avg_over_time @ 120s offset -2m",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "http_requests_total @ 120 offset -2m",
		},
		{
			name: "sum_over_time @ 180s offset 2m",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "sum_over_time(http_requests_total[5m] @ 180 offset 2m)",
		},
		{
			name: "binop with @ end() modifier inside query range",
			load: `load 30s
					http_requests_total 2+3x100
					http_responses_total 2+4x100`,
			query: "max(http_requests_total @ end()) / max(http_responses_total)",
			end:   time.Unix(600, 0),
		},
		{
			name: "binop with @ end() modifier outside of query range",
			load: `load 30s
					http_requests_total 2+3x100
					http_responses_total 2+4x100`,
			query: "max(http_requests_total @ end()) / max(http_responses_total)",
			end:   time.Unix(60000, 0),
		},
		{
			name: "days_in_month with input",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "days_in_month(http_requests_total)",
		},
		{
			name: "days_in_month without input",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "days_in_month()",
		},
		{
			name: "day_of_month with input",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "day_of_month(http_requests_total)",
		},
		{
			name: "day_of_month without input",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "day_of_month()",
		},
		{
			name: "day_of_week with input",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "days_in_month(http_requests_total)",
		},
		{
			name: "day_of_week without input",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "days_in_month()",
		},
		{
			name: "day_of_year with input",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "day_of_year(http_requests_total)",
		},
		{
			name: "day_of_year without input",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "day_of_year()",
		},
		{
			name: "hour with input",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "hour(http_requests_total)",
		},
		{
			name: "hour without input",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "hour()",
		},
		{
			name: "minute with input",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "minute(http_requests_total)",
		},
		{
			name: "minute without input",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "minute()",
		},
		{
			name: "month with input",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "month(http_requests_total)",
		},
		{
			name: "month without input",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "month()",
		},
		{
			name: "year with input",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "year(http_requests_total)",
		},
		{
			name: "year without input",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "year()",
		},
		{
			name: "selector merge",
			load: `load 30s
					http_requests_total{pod="nginx-1", ns="nginx"} 1+1x15
					http_requests_total{pod="nginx-2", ns="nginx"} 1+2x18
					http_requests_total{pod="nginx-3", ns="nginx"} 1+2x21`,
			query: `http_requests_total{pod=~"nginx-1", ns="nginx"} / on() group_left() sum(http_requests_total{ns="nginx"})`,
		},
		{
			name: "selector merge with different ranges",
			load: `load 30s
					http_requests_total{pod="nginx-1", ns="nginx"} 2+2x16
					http_requests_total{pod="nginx-2", ns="nginx"} 2+4x18
					http_requests_total{pod="nginx-3", ns="nginx"} 2+6x20`,
			query: `
	rate(http_requests_total{pod=~"nginx-1", ns="nginx"}[2m])
	+ on() group_left()
	sum(http_requests_total{ns="nginx"})`,
		},
		{
			name: "binop with positive matcher using regex, only one side has data",
			load: `load 30s
					metric{} 1+2x5
					metric{} 1+2x20`,
			query: `sum(rate(metric{err=~".+"}[5m])) / sum(rate(metric{}[5m]))`,
		},
		{
			name: "binop with positive matcher using regex, both sides have data",
			load: `load 30s
					metric{} 1+2x5
					metric{err="FooBarKey"} 1+2x20`,
			query: `sum(rate(metric{err=~".+"}[5m])) / sum(rate(metric{}[5m]))`,
		},
		{
			name: "binop with negative matcher using regex, only one side has data",
			load: `load 30s
					metric{} 1+2x5
					metric{} 1+2x20`,
			query: `sum(rate(metric{err!~".+"}[5m])) / sum(rate(metric{}[5m]))`,
		},
		{
			name: "binop with negative matcher using regex, both sides have data",
			load: `load 30s
					metric{} 1+2x5
					metric{err="FooBarKey"} 1+2x20`,
			query: `sum(rate(metric{err!~".+"}[5m])) / sum(rate(metric{}[5m]))`,
		},
		{
			name: "scalar func with NaN",
			load: `load 30s
		 	http_requests_total{pod="nginx-1"} 1+1x15
		 	http_requests_total{pod="nginx-2"} NaN`,
			query: `scalar(http_requests_total)`,
		},
		{
			name: "scalar func with aggr",
			load: `load 30s
			http_requests_total{pod="nginx-1"} 1+1x15
			http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `scalar(max(http_requests_total))`,
		},
		{
			name: "scalar func with aggr and number on right",
			load: `load 30s
			http_requests_total{pod="nginx-1"} 1+1x15
			http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `scalar(max(http_requests_total)) + 10`,
		},
		{
			name: "scalar func with aggr and number on left",
			load: `load 30s
			http_requests_total{pod="nginx-1"} 1+1x15
			http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `10 + scalar(max(http_requests_total))`,
		},
		{
			name: "quantile",
			load: `load 30s
				       http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				       http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				       http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				       http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				       http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50	`,
			query: "quantile(scalar(sum(http_requests_total)), rate(http_requests_total[1m]))",
		},
		{
			name: "clamp",
			load: `load 30s
			http_requests_total{pod="nginx-1"} 1+1x15
			http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `clamp(http_requests_total, 5, 10)`,
		},
		{
			name: "clamp_min",
			load: `load 30s
			http_requests_total{pod="nginx-1"} 1+1x15
			http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `clamp_min(http_requests_total, 10)`,
		},
		{
			name: "complex func query",
			load: `load 30s
			http_requests_total{pod="nginx-1"} 1+1x15
			http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `clamp(1 - http_requests_total, 10 - 5, 10)`,
		},
		{
			name: "func within func query",
			load: `load 30s
			http_requests_total{pod="nginx-1"} 1+1x15
			http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `clamp(irate(http_requests_total[30s]), 10 - 5, 10)`,
		},
		{
			name: "aggr within func query",
			load: `load 30s
			http_requests_total{pod="nginx-1"} 1+1x15
			http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `clamp(rate(http_requests_total[30s]), 10 - 5, 10)`,
		},
		{
			name: "func with scalar arg that selects storage",
			load: `load 30s
			http_requests_total{pod="nginx-1"} 1+1x15
			http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `clamp_min(http_requests_total, scalar(max(http_requests_total)))`,
		},
		{
			name: "func with scalar arg that selects storage + number",
			load: `load 30s
			http_requests_total{pod="nginx-1"} 1+1x15
			http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `clamp_min(http_requests_total, scalar(max(http_requests_total)) + 10)`,
		},
		{
			name: "histogram quantile",
			load: `load 30s
			http_requests_total{pod="nginx-1", le="1"} 1+3x10
			http_requests_total{pod="nginx-2", le="1"} 2+3x10
			http_requests_total{pod="nginx-1", le="2"} 1+2x10
			http_requests_total{pod="nginx-2", le="2"} 2+2x10
			http_requests_total{pod="nginx-2", le="5"} 3+2x10
			http_requests_total{pod="nginx-1", le="+Inf"} 1+1x10
			http_requests_total{pod="nginx-2", le="+Inf"} 4+1x10`,
			query: `histogram_quantile(0.9, http_requests_total)`,
		},
		{
			name: "histogram quantile on malformed data",
			load: `load 30s
			http_requests_total{pod="nginx-1"} 1+3x10
			http_requests_total{pod="nginx-2"} 2+3x10`,
			query: `histogram_quantile(0.9, http_requests_total)`,
		},
		{
			name: "histogram quantile on partially malformed data",
			load: `load 30s
			http_requests_total{pod="nginx-1", le="1"} 1+3x10
			http_requests_total{pod="nginx-2", le="2"} 2+3x10
			http_requests_total{pod="nginx-3"} 3+3x10
			http_requests_total{pod="nginx-4"} 4+3x10`,
			query: `histogram_quantile(0.9, http_requests_total)`,
		},
		// TODO: uncomment once support for testing NaNs is added.
		{
			name: "histogram quantile on malformed, interleaved data",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+3x10
					http_requests_total{pod="nginx-2"} 2+3x10
					http_requests_total{pod="nginx-3", le="0.05"} 2+3x10
					http_requests_total{pod="nginx-4", le="0.1"} 2+3x10`,
			query: `histogram_quantile(0.9, http_requests_total)`,
		},
		{
			name: "histogram quantile on malformed, interleaved data 2",
			load: `load 30s
					http_requests_total{pod="nginx-1", le="0.01"} 1+3x10
					http_requests_total{pod="nginx-2", le="0.02"} 2+3x10
					http_requests_total{pod="nginx-3"} 2+3x10
					http_requests_total{pod="nginx-4"} 2+3x10`,
			query: `histogram_quantile(0.9, http_requests_total)`,
		},
		{
			name: "histogram quantile on malformed, interleaved data 3",
			load: `load 30s
					http_requests_total{pod="nginx-1", le="0.01"} 1+3x10
					http_requests_total{pod="nginx-2"} 2+3x10
					http_requests_total{pod="nginx-3"} 2+3x10
					http_requests_total{pod="nginx-4", le="0.03"} 2+3x10`,
			query: `histogram_quantile(0.9, http_requests_total)`,
		},
		{
			name: "histogram quantile on malformed, interleaved data 4",
			load: `load 30s
					http_requests_total{pod="nginx-1", le="0.01"} 1+3x10
					http_requests_total{pod="nginx-2"} 2+3x10
					http_requests_total{pod="nginx-2", le="0.05"} 2+3x10
					http_requests_total{pod="nginx-2", le="0.2"} 2+3x10
					http_requests_total{pod="nginx-3"} 2+3x10
					http_requests_total{pod="nginx-4", le="0.03"} 2+3x10`,
			query: `histogram_quantile(0.9, http_requests_total)`,
		},
		{
			name: "histogram quantile with sum",
			load: `load 30s
			http_requests_total{pod="nginx-1", le="1"} 1+3x10
			http_requests_total{pod="nginx-2", le="1"} 2+3x10
			http_requests_total{pod="nginx-1", le="2"} 1+2x10
			http_requests_total{pod="nginx-2", le="2"} 2+2x10
			http_requests_total{pod="nginx-2", le="5"} 3+2x10
			http_requests_total{pod="nginx-1", le="+Inf"} 1+1x10
			http_requests_total{pod="nginx-2", le="+Inf"} 4+1x10`,
			query: `histogram_quantile(0.9, sum by (pod, le) (rate(http_requests_total[2m])))`,
		},
		{
			name: "histogram quantile with scalar operator",
			load: `load 30s
			quantile{pod="nginx-1", le="1"} 1+1x2
			http_requests_total{pod="nginx-1", le="1"} 1+3x10
			http_requests_total{pod="nginx-2", le="1"} 2+3x10
			http_requests_total{pod="nginx-1", le="2"} 1+2x10
			http_requests_total{pod="nginx-2", le="2"} 2+2x10
			http_requests_total{pod="nginx-2", le="5"} 3+2x10
			http_requests_total{pod="nginx-1", le="+Inf"} 1+1x10
			http_requests_total{pod="nginx-2", le="+Inf"} 4+1x10`,
			query: `histogram_quantile(scalar(max(quantile)), http_requests_total)`,
		},
		{
			name: "topk",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="1"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="2"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="2"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "topk(2, http_requests_total)",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "topk with float64 parameter",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="1"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="2"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="2"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "topk(3.5, http_requests_total)",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "topk with float64 parameter that gets truncated to 0",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="1"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="2"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="2"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "topk(0.5, http_requests_total)",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "topk with float64 parameter that does not fit int64",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="1"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="2"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="2"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "topk(1e120, http_requests_total)",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "topk with NaN",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="1"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="2"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="2"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "topk(NaN, http_requests_total)",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name:  "topk with NaN and no matching series",
			query: "topk(NaN, not_there)",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "nested topk error that should not be skipped",
			load: `load 30s
				X 1+1x50`,
			query: "topk (0, topk (NaN, X))",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "topk wrapped by another aggregate",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="1"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="2"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="2"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "max(topk by (series) (2, http_requests_total))",
			end:   time.Unix(3000, 0),
		},
		{
			name: "topk on empty result",
			load: `load 30s
				metric_a 1+1x2`,
			query: "topk(2, histogram_quantile(0.1, metric_b))",
		},
		{
			name: "topk by",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="1"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="2"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="2"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "topk(2, http_requests_total) by (series)",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "topk with simple expression",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="1"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="2"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="2"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "topk(2 - 1, http_requests_total) by (series)",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "topk with expression",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="1"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="2"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="2"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "topk(scalar(min(http_requests_total)), http_requests_total) by (series)",
			start: time.Unix(0, 0),
			end:   time.Unix(500, 0),
			step:  2 * time.Second,
		},
		{
			name: "topk with expression as argument not returning any value",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="1"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="2"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="2"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "topk(scalar(min(non_existent_metric)), http_requests_total) by (series)",
			start: time.Unix(0, 0),
			end:   time.Unix(500, 0),
			step:  2 * time.Second,
		},
		{
			name: "bottomK",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="1"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="2"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="2"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "bottomk(2, http_requests_total)",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "bottomK by",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="1"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="2"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="2"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "bottomk(2, http_requests_total) by (series)",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "sgn",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="1"} -10+1x50`,
			query: "sgn(http_requests_total)",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "sort_desc",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="1"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="2"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="2"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "sort_desc(http_requests_total)",
		},
		{
			name: "count by __name__ label",
			load: `load 30s
				foo 1+1x5
				bar 2+2x5`,
			query: `count by (__name__) ({__name__=~".+"})`,
		},
		{
			name: "scalar with bool",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				http_requests_total{pod="nginx-3", series="3"} 6+0.8x60
				http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: `scalar(avg_over_time({__name__="http_requests_total"}[3m])) > bool 0.9464749352949011`,
		},
		{
			name: "repro https://github.com/thanos-io/promql-engine/issues/239",
			load: `load 30s
				storage_used{storage_index="1010"} 65x20
				storage_used{storage_index="1011"} 125x20
				storage_used{storage_index="1012"} 0x20
				storage_used{storage_index="20"} 2290380x20
				storage_used{storage_index="30"} 397304x20
				storage_used{storage_index="40"} 5590832x20
				storage_used{storage_index="41"} 65559832x20
				storage_used{storage_index="42"} 3516400x20
				storage_info{storage_info="Config", storage_index="40"} 1x20
				storage_info{storage_info="Log", storage_index="41"} 1x20
				storage_info{storage_info="Mem", storage_index="20"} 1x20
				storage_info{storage_info="Root", storage_index="42"} 1x20
				storage_info{storage_info="Swap", storage_index="30"} 1x20`,
			query: `avg by (storage_info) (storage_used * on (instance, storage_index) group_left(storage_info) (sum by (instance, storage_index, storage_info) (storage_info)))`,
		},
	}

	disableOptimizerOpts := []bool{true, false}
	lookbackDeltas := []time.Duration{0, 30 * time.Second, time.Minute, 5 * time.Minute, 10 * time.Minute}
	for _, lookbackDelta := range lookbackDeltas {
		opts.LookbackDelta = lookbackDelta
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				test, err := promql.NewTest(t, tc.load)
				testutil.Ok(t, err)
				defer test.Close()
				testutil.Ok(t, test.Run())

				if tc.start.Equal(time.Time{}) {
					tc.start = start
				}
				if tc.end.Equal(time.Time{}) {
					tc.end = end
				}
				if tc.step == 0 {
					tc.step = step
				}
				for _, disableOptimizers := range disableOptimizerOpts {
					t.Run(fmt.Sprintf("disableOptimizers=%v", disableOptimizers), func(t *testing.T) {
						for _, disableFallback := range []bool{false, true} {
							t.Run(fmt.Sprintf("disableFallback=%v", disableFallback), func(t *testing.T) {
								optimizers := logicalplan.AllOptimizers
								if disableOptimizers {
									optimizers = logicalplan.NoOptimizers
								}
								newEngine := engine.New(engine.Opts{
									EngineOpts:        opts,
									DisableFallback:   disableFallback,
									LogicalOptimizers: optimizers,
								})
								ctx := test.Context()
								q1, err := newEngine.NewRangeQuery(ctx, test.Storage(), nil, tc.query, tc.start, tc.end, tc.step)
								testutil.Ok(t, err)
								defer q1.Close()
								newResult := q1.Exec(ctx)

								oldEngine := promql.NewEngine(opts)
								q2, err := oldEngine.NewRangeQuery(ctx, test.Storage(), nil, tc.query, tc.start, tc.end, tc.step)
								testutil.Ok(t, err)
								defer q2.Close()
								oldResult := q2.Exec(ctx)

								if oldResult.Err != nil {
									testutil.NotOk(t, newResult.Err, "expected error "+oldResult.Err.Error())
									return
								}

								testutil.Ok(t, newResult.Err)
								if hasNaNs(oldResult) {
									t.Log("Applying comparison with NaN equality.")
									equalsWithNaNs(t, oldResult, newResult)
								} else {
									emptyLabelsToNil(oldResult)
									emptyLabelsToNil(newResult)
									testutil.Equals(t, oldResult, newResult)
								}
							})
						}
					})
				}
			})
		}
	}
}

func equalsWithNaNs(t *testing.T, oldResult, newResult interface{}) {
	if reflect.TypeOf(labels.Labels{}).Kind() == reflect.Struct {
		testutil.WithGoCmp(cmpopts.EquateNaNs(), cmp.AllowUnexported(labels.Labels{})).Equals(t, oldResult, newResult)
	} else {
		testutil.WithGoCmp(cmpopts.EquateNaNs()).Equals(t, oldResult, newResult)
	}
}

func hasNaNs(result *promql.Result) bool {
	switch result := result.Value.(type) {
	case promql.Matrix:
		for _, series := range result {
			for _, point := range series.Floats {
				if math.IsNaN(point.F) {
					return true
				}
			}
		}
	case promql.Vector:
		for _, sample := range result {
			if math.IsNaN(sample.F) {
				return true
			}
		}
	case promql.Scalar:
		return math.IsNaN(result.V)
	}

	return false
}

// mergeWithSampleDedup merges samples from series with the same labels,
// removing samples with identical timestamps.
func mergeWithSampleDedup(series []*mockSeries) []storage.Series {
	index := make(map[uint64]*mockSeries)
	for _, s := range series {
		hash := s.Labels().Hash()
		existing, ok := index[hash]
		if !ok {
			// Make a copy to avoid modifying the original series
			// when merging samples.
			index[hash] = &mockSeries{
				labels:     s.labels,
				timestamps: s.timestamps,
				values:     s.values,
			}
			continue
		}
		existing.timestamps = append(existing.timestamps, s.timestamps...)
		existing.values = append(existing.values, s.values...)
	}

	for _, s := range index {
		sort.Sort(byTimestamps(*s))
		// Remove exact timestamp duplicates.
		i := 1
		for i < len(s.timestamps) {
			if s.timestamps[i] == s.timestamps[i-1] {
				s.timestamps = append(s.timestamps[:i], s.timestamps[i+1:]...)
				s.values = append(s.values[:i], s.values[i+1:]...)
			} else {
				i++
			}
		}
	}

	sset := make([]storage.Series, 0, len(index))
	for _, s := range index {
		sset = append(sset, s)
	}
	return sset
}

func TestEdgeCases(t *testing.T) {
	testCases := []struct {
		name   string
		series []storage.Series
		query  string
		start  time.Time
		end    time.Time
	}{
		{
			name: "binop edge case",
			series: []storage.Series{
				newMockSeries(
					[]string{labels.MetricName, "foo"},
					[]int64{0, 30, 60, 1200, 1500, 1800},
					[]float64{1, 2, 3, 4, 5, 6},
				),
				newMockSeries(
					[]string{labels.MetricName, "bar", "id", "1"},
					[]int64{0, 30},
					[]float64{1, 2},
				),
				newMockSeries(
					[]string{labels.MetricName, "bar", "id", "2"},
					[]int64{1200, 1500},
					[]float64{3, 4},
				),
			},
			query: `foo * on () group_left bar`,
			start: time.Unix(0, 0),
			end:   time.Unix(30000, 0),
		},
		{
			name: "absent with gaps in series",
			series: []storage.Series{
				newMockSeries(
					[]string{labels.MetricName, "foo"},
					[]int64{30, 300, 3000, 6000, 12000, 18000},
					[]float64{1, 2, 3, 4, 5, 6},
				),
			},
			query: `absent(foo)`,
			start: time.Unix(0, 0),
			end:   time.Unix(30000, 0),
		},
	}

	opts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}
	step := time.Second * 30
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			oldEngine := promql.NewEngine(opts)
			q1, err := oldEngine.NewRangeQuery(ctx, storageWithSeries(tc.series...), nil, tc.query, tc.start, tc.end, step)
			testutil.Ok(t, err)

			newEngine := engine.New(engine.Opts{EngineOpts: opts})
			q2, err := newEngine.NewRangeQuery(ctx, storageWithSeries(tc.series...), nil, tc.query, tc.start, tc.end, step)
			testutil.Ok(t, err)

			oldResult := q1.Exec(ctx)
			newResult := q2.Exec(ctx)

			testutil.Equals(t, oldResult, newResult)
		})
	}
}

func TestDisabledXFunction(t *testing.T) {
	queryTime := time.Unix(50, 0)
	opts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}

	defaultLoad := `load 5s
	http_requests{path="/foo"}	0+10x10
	http_requests{path="/bar"}	0+10x5 0+10x4`

	cases := []struct {
		name     string
		load     string
		query    string
		expected []promql.Sample
	}{
		{
			name:  "xfunctions disable",
			load:  defaultLoad,
			query: "xincrease(http_requests[50s])",
			expected: []promql.Sample{
				createSample(queryTime.UnixMilli(), 100, labels.FromStrings("path", "/foo")),
				createSample(queryTime.UnixMilli(), 90, labels.FromStrings("path", "/bar")),
			},
		},
	}
	for _, tc := range cases {
		test, err := promql.NewTest(t, tc.load)
		testutil.Ok(t, err)
		defer test.Close()

		testutil.Ok(t, test.Run())
		optimizers := logicalplan.AllOptimizers

		newEngine := engine.New(engine.Opts{
			EngineOpts:        opts,
			DisableFallback:   true,
			LogicalOptimizers: optimizers,
		})
		_, err = newEngine.NewInstantQuery(test.Context(), test.Storage(), nil, tc.query, queryTime)
		testutil.NotOk(t, err)
	}
}

func TestXFunctions(t *testing.T) {
	defaultQueryTime := time.Unix(50, 0)
	// Negative offset and at modifier are enabled by default
	// since Prometheus v2.33.0, so we also enable them.
	opts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}

	defaultLoad := `load 5s
	http_requests{path="/foo"}	0+10x10
	http_requests{path="/bar"}	0+10x5 0+10x4`

	xDeltaLoad := `load 5m
	http_requests{path="/foo"}	0 50 300 150 200
	http_requests{path="/bar"}	200 150 300 50 0`

	cases := []struct {
		name         string
		load         string
		query        string
		queryTime    time.Time
		sortByLabels bool // if true, the series in the result between the old and new engine should be sorted before compared
		expected     []promql.Sample
		rangeQuery   bool
		startTime    time.Time
		endTime      time.Time
	}{
		// Tests for xIncrease
		{
			name:  "eval instant at 50s xincrease, with 50s lookback",
			load:  defaultLoad,
			query: "xincrease(http_requests[50s])",
			expected: []promql.Sample{
				createSample(defaultQueryTime.UnixMilli(), 100, labels.FromStrings("path", "/foo")),
				createSample(defaultQueryTime.UnixMilli(), 90, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:  "eval instant at 50s xincrease, with 5s lookback",
			load:  defaultLoad,
			query: "xincrease(http_requests[5s])",
			expected: []promql.Sample{
				createSample(defaultQueryTime.UnixMilli(), 20, labels.FromStrings("path", "/foo")),
				createSample(defaultQueryTime.UnixMilli(), 20, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:  "eval instant at 50s xincrease, with 3s lookback",
			load:  defaultLoad,
			query: "xincrease(http_requests[3s])",
			expected: []promql.Sample{
				createSample(defaultQueryTime.UnixMilli(), 10, labels.FromStrings("path", "/foo")),
				createSample(defaultQueryTime.UnixMilli(), 10, labels.FromStrings("path", "/bar")),
			},
		},
		// Additional tests
		{
			name:      "eval instant at 17s xincrease, with 5s lookback",
			load:      defaultLoad,
			query:     "xincrease(http_requests[5s])",
			queryTime: time.Unix(17, 0),
			expected: []promql.Sample{
				createSample(17000, 10, labels.FromStrings("path", "/foo")),
				createSample(17000, 10, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 17s xincrease, with 10s lookback",
			load:      defaultLoad,
			query:     "xincrease(http_requests[10s])",
			queryTime: time.Unix(17, 0),
			expected: []promql.Sample{
				createSample(17000, 20, labels.FromStrings("path", "/foo")),
				createSample(17000, 20, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:  "eval instant at 50s xrate, with 50s lookback",
			load:  defaultLoad,
			query: "xrate(http_requests[50s])",
			expected: []promql.Sample{
				createSample(defaultQueryTime.UnixMilli(), 2, labels.FromStrings("path", "/foo")),
				createSample(defaultQueryTime.UnixMilli(), 1.8, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:  "eval instant at 50s xrate, with 100s lookback",
			load:  defaultLoad,
			query: "xrate(http_requests[100s])",
			expected: []promql.Sample{
				createSample(defaultQueryTime.UnixMilli(), 1, labels.FromStrings("path", "/foo")),
				createSample(defaultQueryTime.UnixMilli(), 0.9, labels.FromStrings("path", "/bar")),
			},
		},
		// Test zero series injection.
		{
			name: "eval instant xincrease with only one point",
			load: `load 5m
			http_requests{path="/foo"}	stale stale stale 5`,
			query:     "xincrease(http_requests[1h15m])",
			queryTime: time.Unix(1*60*60+15*60, 0),
			expected: []promql.Sample{
				createSample(time.Unix(1*60*60+15*60, 0).UnixMilli(), 5, labels.FromStrings("path", "/foo")),
			},
		},
		{
			name:  "eval instant at 50s xrate, with 5s lookback",
			load:  defaultLoad,
			query: "xrate(http_requests[5s])",
			expected: []promql.Sample{
				createSample(defaultQueryTime.UnixMilli(), 2, labels.FromStrings("path", "/foo")),
				createSample(defaultQueryTime.UnixMilli(), 2, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:  "eval instant at 50s xrate, with 3s lookback",
			load:  defaultLoad,
			query: "xrate(http_requests[3s])",
			expected: []promql.Sample{
				createSample(defaultQueryTime.UnixMilli(), 2, labels.FromStrings("path", "/foo")),
				createSample(defaultQueryTime.UnixMilli(), 2, labels.FromStrings("path", "/bar")),
			},
		},
		// # Test for increase()/xincrease with counter reset.
		// # When the counter is reset, it always starts at 0.
		// # So the sequence 3 2 (decreasing counter = reset) is interpreted the same as 3 0 1 2.
		// # Prometheus assumes it missed the intermediate values 0 and 1.
		{
			name: "eval instant at 30m increase(http_requests[30m])",
			load: `load 5m
			http_requests{path="/foo"}	0 1 2 3 2 3 4`,
			query:     "increase(http_requests[30m])",
			queryTime: time.Unix(1800, 0),
			expected: []promql.Sample{
				createSample(1800000, 7, labels.FromStrings("path", "/foo")),
			},
		},
		{
			name: "eval instant at 30m xincrease(http_requests[30m])",
			load: `load 5m
			http_requests{path="/foo"}	0 1 2 3 2 3 4`,
			query:     "xincrease(http_requests[30m])",
			queryTime: time.Unix(1800, 0),
			expected: []promql.Sample{
				createSample(1800000, 7, labels.FromStrings("path", "/foo")),
			},
		},
		// Tests for xDelta
		{
			name:      "eval instant at 20m xdelta(http_requests[20m]), with 20m lookback",
			load:      xDeltaLoad,
			query:     "xdelta(http_requests[20m])",
			queryTime: time.Unix(1200, 0),
			expected: []promql.Sample{
				createSample(1200000, 200, labels.FromStrings("path", "/foo")),
				createSample(1200000, -200, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 20m xdelta(http_requests[19m]), with 19m lookback",
			load:      xDeltaLoad,
			query:     "xdelta(http_requests[19m])",
			queryTime: time.Unix(1200, 0),
			expected: []promql.Sample{
				createSample(1200000, 190, labels.FromStrings("path", "/foo")),
				createSample(1200000, -190, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 20m xdelta(http_requests[1m]), with 1m lookback",
			load:      xDeltaLoad,
			query:     "xdelta(http_requests[1m])",
			queryTime: time.Unix(1200, 0),
			expected: []promql.Sample{
				createSample(1200000, 10, labels.FromStrings("path", "/foo")),
				createSample(1200000, -10, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name: "eval instant at 4m xincrease(http_requests[2m]), with 1m lookback",
			load: `load 30s
					http_requests	0 0 0 0 1 1 1 1`,
			query:     "xincrease(http_requests[2m])",
			queryTime: time.Unix(240, 0),
			expected: []promql.Sample{
				createSample(240000, 1, labels.Labels{}),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			test, err := promql.NewTest(t, tc.load)
			testutil.Ok(t, err)
			defer test.Close()

			testutil.Ok(t, test.Run())
			queryTime := defaultQueryTime
			if tc.queryTime != (time.Time{}) {
				queryTime = tc.queryTime
			}

			optimizers := logicalplan.AllOptimizers

			ctx := test.Context()
			newEngine := engine.New(engine.Opts{
				EngineOpts:        opts,
				DisableFallback:   true,
				LogicalOptimizers: optimizers,
				EnableXFunctions:  true,
			})
			query, err := newEngine.NewInstantQuery(ctx, test.Storage(), nil, tc.query, queryTime)
			testutil.Ok(t, err)
			defer query.Close()

			engineResult := query.Exec(ctx)
			testutil.Ok(t, engineResult.Err)
			expectedResult := createVectorResult(tc.expected)

			testutil.Equals(t, expectedResult.Err, engineResult.Err)

			exR := expectedResult.Value.(promql.Vector)
			erR := engineResult.Value.(promql.Vector)

			sort.Slice(exR, func(i, j int) bool {
				return labels.Compare(exR[i].Metric, exR[j].Metric) < 0
			})

			sort.Slice(erR, func(i, j int) bool {
				return labels.Compare(erR[i].Metric, erR[j].Metric) < 0
			})

			testutil.Equals(t, exR, erR)
		})
	}
}

func TestRateVsXRate(t *testing.T) {
	defaultQueryTime := time.Unix(25, 0)
	// Negative offset and at modifier are enabled by default
	// since Prometheus v2.33.0, so we also enable them.
	opts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}

	defaultLoad := `load 5s
	http_requests{path="/foo"}  1 1 1 2 2 2 2 2 3 3 3
	http_requests{path="/bar"}  1 2 3 4 5 6 7 8 9 10 11`

	cases := []struct {
		name         string
		load         string
		query        string
		queryTime    time.Time
		sortByLabels bool // if true, the series in the result between the old and new engine should be sorted before compared
		expected     promql.Vector
		rangeQuery   bool
		startTime    time.Time
		endTime      time.Time
	}{
		// ### Timeseries starts insice range, (presumably) goes on after range end. ###
		// 1. Reference eval
		{
			name:      "eval instant at 25s rate, with 50s lookback",
			query:     "rate(http_requests[50s])",
			queryTime: time.Unix(25, 0),
			expected: []promql.Sample{
				createSample(defaultQueryTime.UnixMilli(), 0.022, labels.FromStrings("path", "/foo")),
				createSample(defaultQueryTime.UnixMilli(), 0.12, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 25s xrate, with 50s lookback",
			query:     "xrate(http_requests[50s])",
			queryTime: time.Unix(25, 0),
			expected: []promql.Sample{
				createSample(defaultQueryTime.UnixMilli(), 0.02, labels.FromStrings("path", "/foo")),
				createSample(defaultQueryTime.UnixMilli(), 0.1, labels.FromStrings("path", "/bar")),
			},
		},
		// 2. Eval 1 second earlier compared to (1).
		// * path="/foo" rate should be same or fractionally higher ("shorter" sample, same actual increase);
		// * path="/bar" rate should be same or fractionally lower (80% the increase, 80/96% range covered by sample).
		// XXX Seeing ~20% jump for path="/foo"
		{
			name:      "eval instant at 24s rate(http_requests[50s]), with 50s lookback",
			query:     "rate(http_requests[50s])",
			queryTime: time.Unix(24, 0),
			expected: []promql.Sample{
				createSample(24000, 0.0265, labels.FromStrings("path", "/foo")),
				createSample(24000, 0.116, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 24s xrate(http_requests[50s]), with 50s lookback",
			query:     "xrate(http_requests[50s])",
			queryTime: time.Unix(24, 0),
			expected: []promql.Sample{
				createSample(24000, 0.02, labels.FromStrings("path", "/foo")),
				createSample(24000, 0.08, labels.FromStrings("path", "/bar")),
			},
		},
		// 3. Eval 1 second later compared to (1)
		// * path="/foo" rate should be same or fractionally lower ("longer" sample, same actual increase).
		// * path="/bar" rate should be same or fractionally lower ("longer" sample, same actual increase).
		// XXX Higher instead of lower for both.
		{
			name:      "eval instant at 26s rate(http_requests[50s]), with 50s lookback",
			query:     "rate(http_requests[50s])",
			queryTime: time.Unix(26, 0),
			expected: []promql.Sample{
				createSample(26000, 0.02279999999, labels.FromStrings("path", "/foo")),
				createSample(26000, 0.124, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 26s xrate(http_requests[50s]), with 50s lookback",
			query:     "xrate(http_requests[50s])",
			queryTime: time.Unix(26, 0),
			expected: []promql.Sample{
				createSample(26000, 0.02, labels.FromStrings("path", "/foo")),
				createSample(26000, 0.1, labels.FromStrings("path", "/bar")),
			},
		},
		// ### Timeseries starts before range, ends within range. ###
		// 4. Reference eval
		{
			name:      "eval instant at 75s rate(http_requests[50s]), with 50s lookback",
			query:     "rate(http_requests[50s])",
			queryTime: time.Unix(75, 0),
			expected: []promql.Sample{
				createSample(75000, 0.022, labels.FromStrings("path", "/foo")),
				createSample(75000, 0.11, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 75s xrate(http_requests[50s]), with 50s lookback",
			query:     "xrate(http_requests[50s])",
			queryTime: time.Unix(75, 0),
			expected: []promql.Sample{
				createSample(75000, 0.02, labels.FromStrings("path", "/foo")),
				createSample(75000, 0.12, labels.FromStrings("path", "/bar")),
			},
		},
		// 5. Eval 1s earlier compared to (4)
		// * path="/foo" rate should be same or fractionally lower ("longer" sample, same actual increase).
		// * path="/bar" rate should be same or fractionally lower ("longer" sample, same actual increase).
		// # XXX Higher instead of lower for both.
		{
			name:      "eval instant at 74s rate(http_requests[50s]), with 50s lookback",
			query:     "rate(http_requests[50s])",
			queryTime: time.Unix(74, 0),
			expected: []promql.Sample{
				createSample(74000, 0.02279999999, labels.FromStrings("path", "/foo")),
				createSample(74000, 0.11399999999, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 74s xrate(http_requests[50s]), with 50s lookback",
			query:     "xrate(http_requests[50s])",
			queryTime: time.Unix(74, 0),
			expected: []promql.Sample{
				createSample(74000, 0.02, labels.FromStrings("path", "/foo")),
				createSample(74000, 0.12, labels.FromStrings("path", "/bar")),
			},
		},
		// 6. Eval 1s later compared to (4). Rate/increase (should be) fractionally smaller.
		// * path="/foo" rate should be same or fractionally higher ("shorter" sample, same actual increase)
		// * path="/bar" rate should be same or fractionally lower (80% the increase, 80/96% range covered by sample).
		// XXX Seeing ~20% jump for path="/foo", decrease instead of increase for path="/bar".
		{
			name:      "eval instant at 76s rate(http_requests[50s]), with 50s lookback",
			query:     "rate(http_requests[50s])",
			queryTime: time.Unix(76, 0),
			expected: []promql.Sample{
				createSample(76000, 0.0265, labels.FromStrings("path", "/foo")),
				createSample(76000, 0.106, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 76s xrate(http_requests[50s]), with 50s lookback",
			query:     "xrate(http_requests[50s])",
			queryTime: time.Unix(76, 0),
			expected: []promql.Sample{
				createSample(76000, 0.02, labels.FromStrings("path", "/foo")),
				createSample(76000, 0.1, labels.FromStrings("path", "/bar")),
			},
		},
		// Evaluation of 10 second rate every 10 seconds
		{
			name:      "eval instant at 9s rate(http_requests[10s]), with 10s lookback",
			query:     "rate(http_requests[10s])",
			queryTime: time.Unix(9, 0),
			expected: []promql.Sample{
				createSample(9000, 0, labels.FromStrings("path", "/foo")),
				createSample(9000, 0.2, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 19s rate(http_requests[10s]), with 10s lookback",
			query:     "rate(http_requests[10s])",
			queryTime: time.Unix(19, 0),
			expected: []promql.Sample{
				createSample(19000, 0.2, labels.FromStrings("path", "/foo")),
				createSample(19000, 0.2, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 29s rate(http_requests[10s]), with 10s lookback",
			query:     "rate(http_requests[10s])",
			queryTime: time.Unix(29, 0),
			expected: []promql.Sample{
				createSample(29000, 0, labels.FromStrings("path", "/foo")),
				createSample(29000, 0.2, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 39s rate(http_requests[10s]), with 10s lookback",
			query:     "rate(http_requests[10s])",
			queryTime: time.Unix(39, 0),
			expected: []promql.Sample{
				createSample(39000, 0, labels.FromStrings("path", "/foo")),
				createSample(39000, 0.2, labels.FromStrings("path", "/bar")),
			},
		},
		// XXX Missed an increase in path="/foo" between timestamps 35 and 40 (both in this eval and the one before).
		{
			name:      "eval instant at 49s rate(http_requests[10s]), with 10s lookback",
			query:     "rate(http_requests[10s])",
			queryTime: time.Unix(49, 0),
			expected: []promql.Sample{
				createSample(49000, 0, labels.FromStrings("path", "/foo")),
				createSample(49000, 0.2, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 9s xrate(http_requests[50s]), with 10s lookback",
			query:     "xrate(http_requests[10s])",
			queryTime: time.Unix(9, 0),
			expected: []promql.Sample{
				createSample(9000, 0, labels.FromStrings("path", "/foo")),
				createSample(9000, 0.1, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 19s xrate(http_requests[50s]), with 10s lookback",
			query:     "xrate(http_requests[10s])",
			queryTime: time.Unix(19, 0),
			expected: []promql.Sample{
				createSample(19000, 0.1, labels.FromStrings("path", "/foo")),
				createSample(19000, 0.2, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 29s xrate(http_requests[50s]), with 10s lookback",
			query:     "xrate(http_requests[10s])",
			queryTime: time.Unix(29, 0),
			expected: []promql.Sample{
				createSample(29000, 0, labels.FromStrings("path", "/foo")),
				createSample(29000, 0.2, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 39s xrate(http_requests[50s]), with 10s lookback",
			query:     "xrate(http_requests[10s])",
			queryTime: time.Unix(39, 0),
			expected: []promql.Sample{
				createSample(39000, 0, labels.FromStrings("path", "/foo")),
				createSample(39000, 0.2, labels.FromStrings("path", "/bar")),
			},
		},
		// Sees the increase in path="/foo" between timestamps 35 and 40.
		{
			name:      "eval instant at 49s xrate(http_requests[50s]), with 10s lookback",
			query:     "xrate(http_requests[10s])",
			queryTime: time.Unix(49, 0),
			expected: []promql.Sample{
				createSample(49000, 0.1, labels.FromStrings("path", "/foo")),
				createSample(49000, 0.2, labels.FromStrings("path", "/bar")),
			},
		},
		// xincrease injects a zero if there is only one sample in the given timerange.
		{
			name:      "eval instant at 1s xincrease(http_requests[50s]), with 5s lookback",
			query:     "xincrease(http_requests[5s])",
			queryTime: time.Unix(1, 0),
			expected: []promql.Sample{
				createSample(1000, 1, labels.FromStrings("path", "/foo")),
				createSample(1000, 1, labels.FromStrings("path", "/bar")),
			},
		},
		// xincrease injects a zero if there is only one sample in the given timerange.
		{
			name:      "eval instant at 1s xincrease(http_requests[50s]), with 5s lookback",
			query:     "xincrease(http_requests[5s])",
			queryTime: time.Unix(1, 0),
			expected: []promql.Sample{
				createSample(1000, 1, labels.FromStrings("path", "/foo")),
				createSample(1000, 1, labels.FromStrings("path", "/bar")),
			},
		},
		// xincrease does not inject anything at the end of the given timerange if there are two or more samples.
		{
			name:      "eval instant at 55s xincrease(http_requests[10s]), with 10s lookback",
			query:     "xincrease(http_requests[10s])",
			queryTime: time.Unix(55, 0),
			expected: []promql.Sample{
				createSample(55000, 0, labels.FromStrings("path", "/foo")),
				createSample(55000, 2, labels.FromStrings("path", "/bar")),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			load := defaultLoad
			if tc.load != "" {
				load = tc.load
			}

			test, err := promql.NewTest(t, load)
			testutil.Ok(t, err)
			defer test.Close()

			testutil.Ok(t, test.Run())
			queryTime := defaultQueryTime
			if tc.queryTime != (time.Time{}) {
				queryTime = tc.queryTime
			}

			optimizers := logicalplan.AllOptimizers

			newEngine := engine.New(engine.Opts{
				EngineOpts:        opts,
				DisableFallback:   true,
				LogicalOptimizers: optimizers,
				EnableXFunctions:  true,
			})
			query, err := newEngine.NewInstantQuery(test.Context(), test.Storage(), nil, tc.query, queryTime)
			testutil.Ok(t, err)
			defer query.Close()

			engineResult := query.Exec(test.Context())
			testutil.Ok(t, engineResult.Err)
			// Round engine result.
			roundValues(engineResult)
			expectedResult := createVectorResult(tc.expected)

			testutil.Equals(t, expectedResult.Err, engineResult.Err)

			exR := expectedResult.Value.(promql.Vector)
			erR := engineResult.Value.(promql.Vector)

			sort.Slice(exR, func(i, j int) bool {
				return labels.Compare(exR[i].Metric, exR[j].Metric) < 0
			})

			sort.Slice(erR, func(i, j int) bool {
				return labels.Compare(erR[i].Metric, erR[j].Metric) < 0
			})

			testutil.Equals(t, exR, erR)
		})
	}
}

func createSample(t int64, v float64, metric labels.Labels) promql.Sample {
	return promql.Sample{
		T:      t,
		F:      v,
		H:      nil,
		Metric: metric,
	}
}

func createVectorResult(vector promql.Vector) *promql.Result {
	return &promql.Result{
		Err:      nil,
		Value:    vector,
		Warnings: nil,
	}
}

func TestInstantQuery(t *testing.T) {
	defaultQueryTime := time.Unix(50, 0)
	// Negative offset and at modifier are enabled by default
	// since Prometheus v2.33.0, so we also enable them.
	opts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}

	cases := []struct {
		load         string
		name         string
		query        string
		queryTime    time.Time
		sortByLabels bool // if true, the series in the result between the old and new engine should be sorted before compared
	}{
		{
			name: "duplicate label set",
			load: `load 5m
  testmetric1{src="a",dst="b"} 0
  testmetric2{src="a",dst="b"} 1`,
			query: "changes({__name__=~'testmetric1|testmetric2'}[5m])",
		},
		{
			name:      "scalar",
			load:      ``,
			queryTime: time.Unix(160, 0),
			query:     "12 + 1",
		},
		{
			name:      "string literal",
			load:      ``,
			queryTime: time.Unix(160, 0),
			query:     "test-string-literal",
		},
		{
			name: "increase plus offset",
			load: `load 1s
			http_requests_total{pod="nginx-1"} 1+1x180`,
			queryTime: time.Unix(160, 0),
			query:     "increase(http_requests_total[1m] offset 1m)",
		},
		{
			name: "round",
			load: `load 1s
				       http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				       http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				       http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				       http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				       http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			queryTime: time.Unix(0, 0),
			query:     "round(http_requests_total)",
		},
		{
			name: "round with argument",
			load: `load 1s
				       http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				       http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				       http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				       http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				       http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			queryTime: time.Unix(0, 0),
			query:     "round(http_requests_total, 0.5)",
		},
		{
			name: "sort",
			load: `load 1s
				       http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				       http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				       http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				       http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				       http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			queryTime: time.Unix(0, 0),
			query:     "sort(http_requests_total)",
		},
		{
			name: "sort_desc",
			load: `load 1s
				       http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				       http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				       http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				       http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				       http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			queryTime: time.Unix(0, 0),
			query:     "sort_desc(http_requests_total)",
		},
		{
			name: "quantile by pod",
			load: `load 30s
				       http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				       http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				       http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				       http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				       http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "quantile by (pod) (0.9, rate(http_requests_total[1m]))",
		},
		{
			name: "quantile by pod with binary",
			load: `load 30s
				       http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				       http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				       http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				       http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				       http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "quantile by (pod) (1 - 0.1, rate(http_requests_total[1m]))",
		},
		{
			name: "quantile by pod with expression",
			load: `load 30s
				       http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				       http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				       http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				       http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				       http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "quantile by (pod) (scalar(min(http_requests_total)), rate(http_requests_total[1m]))",
		},
		{
			name: "quantile",
			load: `load 30s
				       http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				       http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				       http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				       http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				       http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50	`,
			query: "quantile(0.9, rate(http_requests_total[1m]))",
		},
		{
			name: "stdvar",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x4
						http_requests_total{pod="nginx-2"} 1+2x4`,
			query: "stdvar(http_requests_total)",
		},
		{
			name: "stdvar by pod",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1
						http_requests_total{pod="nginx-2"} 2
						http_requests_total{pod="nginx-3"} 8
						http_requests_total{pod="nginx-4"} 6`,
			query: "stdvar by (pod) (http_requests_total)",
		},
		{
			name: "stddev",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x4
						http_requests_total{pod="nginx-2"} 1+2x4`,
			query: "stddev(http_requests_total)",
		},
		{
			name: "stddev by pod",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1
						http_requests_total{pod="nginx-2"} 2
						http_requests_total{pod="nginx-3"} 8
						http_requests_total{pod="nginx-4"} 6`,
			query: "stddev by (pod) (http_requests_total)",
		},
		{
			name: "sum by pod",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x4
						http_requests_total{pod="nginx-2"} 1+2x4`,
			query: "sum by (pod) (http_requests_total)",
		},
		{
			name: "count",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "count(http_requests_total)",
		},
		{
			name: "average",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "avg(http_requests_total)",
		},
		{
			name: "label_join",
			load: `load 30s
						http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
						http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
						http_requests_total{pod="nginx-4", series="3"} 5+2.4x50`,
			queryTime: time.Unix(160, 0),
			query:     `label_join(http_requests_total{}, "label", "-", "pod", "series")`,
		},
		{
			name: "label_join with non-existing src labels",
			load: `load 30s
						http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
						http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
						http_requests_total{pod="nginx-4", series="3"} 5+2.4x50`,
			queryTime: time.Unix(160, 0),
			query:     `label_join(http_requests_total{}, "label", "-", "test", "fake")`,
		},
		{
			name: "label_join with overwrite dst label if exists",
			load: `load 30s
						http_requests_total{pod="nginx-1", series="1", label="test-1"} 1+1.1x40
						http_requests_total{pod="nginx-2", series="2", label="test-2"} 2+2.3x50
						http_requests_total{pod="nginx-4", series="3", label="test-3"} 5+2.4x50`,
			queryTime: time.Unix(160, 0),
			query:     `label_join(http_requests_total{}, "label", "-", "pod", "series")`,
		},
		{
			name: "label_join with no src labels provided",
			load: `load 30s
						http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
						http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
						http_requests_total{pod="nginx-4", series="3"} 5+2.4x50`,
			queryTime: time.Unix(160, 0),
			query:     `label_join(http_requests_total{}, "label", "-")`,
		},
		{
			name: "topk",
			load: `load 30s
						http_requests_total{pod="nginx-1", series="1"} 1
						http_requests_total{pod="nginx-2", series="1"} 2
						http_requests_total{pod="nginx-3", series="1"} 8
						http_requests_total{pod="nginx-4", series="2"} 6
						http_requests_total{pod="nginx-5", series="2"} 8
						http_requests_total{pod="nginx-6", series="3"} 15
						http_requests_total{pod="nginx-7", series="3"} 11
						http_requests_total{pod="nginx-8", series="4"} 22
						http_requests_total{pod="nginx-9", series="4"} 89`,
			query: "topk(2, http_requests_total)",
		},
		{
			name: "topk by series",
			load: `load 30s
						http_requests_total{pod="nginx-1", series="1"} 1
						http_requests_total{pod="nginx-2", series="1"} 2
						http_requests_total{pod="nginx-3", series="1"} 8
						http_requests_total{pod="nginx-4", series="2"} 6
						http_requests_total{pod="nginx-5", series="2"} 8
						http_requests_total{pod="nginx-6", series="3"} 15
						http_requests_total{pod="nginx-7", series="3"} 11
						http_requests_total{pod="nginx-8", series="4"} 22
						http_requests_total{pod="nginx-9", series="4"} 89`,
			query:        "topk(2, http_requests_total) by (series)",
			sortByLabels: true,
		},
		{
			name: "bottomK",
			load: `load 30s
						http_requests_total{pod="nginx-1", series="1"} 1
						http_requests_total{pod="nginx-2", series="1"} 2
						http_requests_total{pod="nginx-3", series="1"} 8
						http_requests_total{pod="nginx-4", series="2"} 6
						http_requests_total{pod="nginx-5", series="2"} 8
						http_requests_total{pod="nginx-6", series="3"} 15
						http_requests_total{pod="nginx-7", series="3"} 11
						http_requests_total{pod="nginx-8", series="4"} 22
						http_requests_total{pod="nginx-9", series="4"} 89`,
			query: "bottomk(2, http_requests_total)",
		},
		{
			name: "bottomk by series",
			load: `load 30s
						http_requests_total{pod="nginx-1", series="1"} 1
						http_requests_total{pod="nginx-2", series="1"} 2
						http_requests_total{pod="nginx-3", series="1"} 8
						http_requests_total{pod="nginx-4", series="2"} 6
						http_requests_total{pod="nginx-5", series="2"} 8
						http_requests_total{pod="nginx-6", series="3"} 15
						http_requests_total{pod="nginx-7", series="3"} 11
						http_requests_total{pod="nginx-8", series="4"} 22
						http_requests_total{pod="nginx-9", series="4"} 89`,
			query:        "bottomk(2, http_requests_total) by (series)",
			sortByLabels: true,
		},
		{
			name: "max",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "max(http_requests_total)",
		},
		{
			name: "max with only 1 sample",
			load: `load 30s
						http_requests_total{pod="nginx-1"} -1
						http_requests_total{pod="nginx-2"} 1`,
			query: "max(http_requests_total) by (pod)",
		},
		{
			name: "min",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "min(http_requests_total)",
		},
		{
			name: "min with only 1 sample",
			load: `load 30s
						http_requests_total{pod="nginx-1"} -1
						http_requests_total{pod="nginx-2"} 1`,
			query: "min(http_requests_total) by (pod)",
		},
		{
			name: "rate",
			load: `load 30s
					http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
					http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
					http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
					http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
					http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "rate(http_requests_total[1m])",
		},
		{
			name: "sum rate",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x4
						http_requests_total{pod="nginx-2"} 1+2x4`,
			query: "sum(rate(http_requests_total[1m]))",
		},
		{
			name: "sum rate with single sample series",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x4
						http_requests_total{pod="nginx-2"} 1+2x4
						http_requests_total{pod="nginx-3"} 0`,
			query: "sum by (pod) (rate(http_requests_total[1m]))",
		},
		{
			name: "sum rate with stale series",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x4
						http_requests_total{pod="nginx-2"} 1+2x20`,
			query: "sum(rate(http_requests_total[1m]))",
		},
		{
			name: "delta",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x4
						http_requests_total{pod="nginx-2"} 1+2x4`,
			query: "delta(http_requests_total[1m])",
		},
		{
			name: "increase",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x4
						http_requests_total{pod="nginx-2"} 1+2x4`,
			query: "increase(http_requests_total[1m])",
		},

		{
			name: "sum irate",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x4
						http_requests_total{pod="nginx-2"} 1+2x4`,
			query: "sum(irate(http_requests_total[1m]))",
		},
		{
			name: "sum irate with stale series",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x4
						http_requests_total{pod="nginx-2"} 1+2x20`,
			query: "sum(irate(http_requests_total[1m]))",
		},
		{
			name:  "number literal",
			load:  "",
			query: "34",
		},
		{
			name:  "vector",
			load:  "",
			query: "vector(24)",
		},
		{
			name: "binary operation with vector and scalar on the right",
			load: `load 30s
					foo{method="get", code="500"} 1+1.1x30
					foo{method="get", code="404"} 1+2.2x20`,
			query: `foo * 2`,
		},
		{
			name: "binary operation with vector and scalar on the left",
			load: `load 30s
					foo{method="get", code="500"} 1+1.1x30
					foo{method="get", code="404"} 1+2.2x20`,
			query: `2 * foo`,
		},
		{
			name: "complex binary operation",
			load: `load 30s
					foo{method="get", code="500"} 1+1.1x30
					foo{method="get", code="404"} 1+2.2x20`,
			query: `1 - (100 * sum(foo{method="get"}) / sum(foo))`,
		},
		{
			name: "vector binary op ==",
			load: `load 30s
					foo{method="get", code="500"} 1+1x40
					bar{method="get", code="404"} 1+1.1x30`,
			query: `sum(foo) by (method) == sum(bar) by (method)`,
		},
		{
			name: "vector binary op !=",
			load: `load 30s
					foo{method="get", code="500"} 1+1x40
					bar{method="get", code="404"} 1+1.1x30`,
			query: `sum(foo) by (method) != sum(bar) by (method)`,
		},
		{
			name: "vector binary op >",
			load: `load 30s
					foo{method="get", code="500"} 1+1x40
					bar{method="get", code="404"} 1+1.1x30`,
			query: `sum(foo) by (method) > sum(bar) by (method)`,
		},
		{
			name: "vector binary op <",
			load: `load 30s
					foo{method="get", code="500"} 1+1x40
					bar{method="get", code="404"} 1+1.1x30`,
			query: `sum(foo) by (method) < sum(bar) by (method)`,
		},
		{
			name: "vector binary op >=",
			load: `load 30s
					foo{method="get", code="500"} 1+1x40
					bar{method="get", code="404"} 1+1.1x30`,
			query: `sum(foo) by (method) >= sum(bar) by (method)`,
		},
		{
			name: "vector binary op <=",
			load: `load 30s
					foo{method="get", code="500"} 1+1x40
					bar{method="get", code="404"} 1+1.1x30`,
			query: `sum(foo) by (method) <= sum(bar) by (method)`,
		},
		{
			name: "vector binary op ^",
			load: `load 30s
					foo{method="get", code="500"} 1+1x40
					bar{method="get", code="404"} 1+1.1x30`,
			query: `sum(foo) by (method) ^ sum(bar) by (method)`,
		},
		{
			name: "vector binary op %",
			load: `load 30s
					foo{method="get", code="500"} 1+2x40
					bar{method="get", code="404"} 1+1x30`,
			query: `sum(foo) by (method) % sum(bar) by (method)`,
		},
		{
			name: "vector binary op > scalar",
			load: `load 30s
					foo{method="get", code="500"} 1+2x40
					bar{method="get", code="404"} 1+1x30`,
			query: `sum(foo) by (method) > 10`,
		},
		{
			name: "scalar < vector binary op",
			load: `load 30s
					foo{method="get", code="500"} 1+2x40
					bar{method="get", code="404"} 1+1x30`,
			query: `10 < sum(foo) by (method)`,
		},
		{
			name:  "scalar binary op == true",
			load:  ``,
			query: `1 == bool 1`,
		},
		{
			name:  "scalar binary op == false",
			load:  ``,
			query: `1 != bool 2`,
		},
		{
			name:  "scalar binary op !=",
			load:  ``,
			query: `1 != bool 1`,
		},
		{
			name:  "scalar binary op >",
			load:  ``,
			query: `1 > bool 0`,
		},
		{
			name:  "scalar binary op <",
			load:  ``,
			query: `1 > bool 2`,
		},
		{
			name:  "scalar binary op >=",
			load:  ``,
			query: `1 >= bool 0`,
		},
		{
			name:  "scalar binary op <=",
			load:  ``,
			query: `1 <= bool 2`,
		},
		{
			name:  "scalar binary op % 0",
			load:  ``,
			query: `2 % 2`,
		},
		{
			name:  "scalar binary op % 1",
			load:  ``,
			query: `1 % 2`,
		},
		{
			name:  "scalar binary op ^",
			load:  ``,
			query: `2 ^ 2`,
		},
		{
			name:  "empty series",
			load:  "",
			query: "http_requests_total",
		},
		{
			name:  "empty series with func",
			load:  "",
			query: "sum(http_requests_total)",
		},
		{
			name: "empty result",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total{pod="nginx-3"}`,
		},
		{
			name: "last_over_time",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `last_over_time(http_requests_total[30s])`,
		},
		{
			name: "group",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `group(http_requests_total)`,
		},
		{
			name: "reset",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 100-1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `resets(http_requests_total[5m])`,
		},
		{
			name: "present_over_time",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `present_over_time(http_requests_total[30s])`,
		},
		{
			name:  "unary sub operation for scalar",
			load:  ``,
			query: `-(1 + 5)`,
		},
		{
			name: "unary sub operation for vector",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `-http_requests_total`,
		},
		{
			name: "unary add operation for vector",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `+http_requests_total`,
		},
		{
			name: "vector positive offset",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total offset 30s`,
		},
		{
			name: "vector negative offset",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total offset -30s`,
		},
		{
			name: "matrix negative offset with sum_over_time",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x25
						http_requests_total{pod="nginx-2"} 1+2x28`,
			query: `sum_over_time(http_requests_total[5m] offset 5m)`,
		},
		{
			name: "matrix negative offset with count_over_time",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `count_over_time(http_requests_total[5m] offset -2m)`,
		},
		{
			name: "@ vector time 10s",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "http_requests_total @ 10",
		},
		{
			name: "@ vector time 120s",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "http_requests_total @ 120",
		},
		{
			name: "@ vector time 360s",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "http_requests_total @ 360",
		},
		{
			name: "@ vector start",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "http_requests_total @ start()",
		},
		{
			name: "@ vector end",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "http_requests_total @ end()",
		},
		{
			name: "count_over_time @ start",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "count_over_time(http_requests_total[5m] @ start())",
		},
		{
			name: "sum_over_time @ end",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "sum_over_time(http_requests_total[5m] @ start())",
		},
		{
			name: "avg_over_time @ 180s",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "avg_over_time(http_requests_total[4m] @ 180)",
		},
		{
			name: "@ vector 240s offset 2m",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "http_requests_total @ 240 offset 2m",
		},
		{
			name: "avg_over_time @ 120s offset -2m",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "http_requests_total @ 120 offset -2m",
		},
		{
			name: "sum_over_time @ 180s offset 2m",
			load: `load 30s
						http_requests_total{pod="nginx-1"} 1+1x15
						http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "sum_over_time(http_requests_total[5m] @ 180 offset 2m)",
		},
		{
			name:  "scalar func with non existent metric in scalar comparison",
			query: `scalar(non_existent_metric) < bool 0`,
		},
		{
			name: "scalar func with NaN",
			load: `load 30s
		 	http_requests_total{pod="nginx-1"} 1+1x15
		 	http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `scalar(http_requests_total)`,
		},
		{
			name: "scalar func with aggr",
			load: `load 30s
			http_requests_total{pod="nginx-1"} 1+1x15
			http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `scalar(max(http_requests_total))`,
		},
		{
			name: "scalar func with aggr and number on right",
			load: `load 30s
			http_requests_total{pod="nginx-1"} 1+1x15
			http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `scalar(max(http_requests_total)) + 10`,
		},
		{
			name: "scalar func with aggr and number on left",
			load: `load 30s
			http_requests_total{pod="nginx-1"} 1+1x15
			http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `10 + scalar(max(http_requests_total))`,
		},
		{
			name: "clamp",
			load: `load 30s
				http_requests_total{pod="nginx-1"} 1+1x15
				http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `clamp(http_requests_total, 5, 10)`,
		},
		{
			name: "clamp_min",
			load: `load 30s
				http_requests_total{pod="nginx-1"} 1+1x15
				http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `clamp_min(http_requests_total, 10)`,
		},
		{
			name: "complex func query",
			load: `load 30s
				http_requests_total{pod="nginx-1"} 1+1x15
				http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `clamp(1 - http_requests_total, 10 - 5, 10)`,
		},
		{
			name: "func within func query",
			load: `load 30s
				http_requests_total{pod="nginx-1"} 1+1x15
				http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `clamp(irate(http_requests_total[30s]), 10 - 5, 10)`,
		},
		{
			name: "aggr within func query",
			load: `load 30s
				http_requests_total{pod="nginx-1"} 1+1x15
				http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `clamp(rate(http_requests_total[30s]), 10 - 5, 10)`,
		},
		{
			name: "func with scalar arg that selects storage",
			load: `load 30s
				http_requests_total{pod="nginx-1"} 1+1x15
				http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `clamp_min(http_requests_total, scalar(max(http_requests_total)))`,
		},
		{
			name: "func with scalar arg that selects storage + number",
			load: `load 30s
				http_requests_total{pod="nginx-1"} 1+1x15
				http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `clamp_min(http_requests_total, scalar(max(http_requests_total)) + 10)`,
		},
		{
			name: "sgn",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="1"} -10+1x50
				http_requests_total{pod="nginx-3", series="1"} NaN`,
			query: "sgn(http_requests_total)",
		},
	}

	disableOptimizerOpts := []bool{true, false}
	lookbackDeltas := []time.Duration{0, 30 * time.Second, time.Minute, 5 * time.Minute, 10 * time.Minute}
	for _, disableOptimizers := range disableOptimizerOpts {
		t.Run(fmt.Sprintf("disableOptimizers=%t", disableOptimizers), func(t *testing.T) {
			for _, lookbackDelta := range lookbackDeltas {
				opts.LookbackDelta = lookbackDelta
				for _, tc := range cases {
					t.Run(tc.name, func(t *testing.T) {
						test, err := promql.NewTest(t, tc.load)
						testutil.Ok(t, err)
						defer test.Close()

						testutil.Ok(t, test.Run())

						for _, disableFallback := range []bool{false, true} {
							t.Run(fmt.Sprintf("disableFallback=%v", disableFallback), func(t *testing.T) {
								var queryTime time.Time = defaultQueryTime
								if tc.queryTime != (time.Time{}) {
									queryTime = tc.queryTime
								}

								optimizers := logicalplan.AllOptimizers
								if disableOptimizers {
									optimizers = logicalplan.NoOptimizers
								}
								newEngine := engine.New(engine.Opts{
									EngineOpts:        opts,
									DisableFallback:   disableFallback,
									LogicalOptimizers: optimizers,
								})

								ctx := test.Context()
								q1, err := newEngine.NewInstantQuery(ctx, test.Storage(), nil, tc.query, queryTime)
								testutil.Ok(t, err)
								defer q1.Close()

								newResult := q1.Exec(ctx)

								oldEngine := promql.NewEngine(opts)
								q2, err := oldEngine.NewInstantQuery(ctx, test.Storage(), nil, tc.query, queryTime)
								testutil.Ok(t, err)
								defer q2.Close()

								oldResult := q2.Exec(ctx)

								if tc.sortByLabels {
									sortByLabels(oldResult)
									sortByLabels(newResult)
								}

								if hasNaNs(oldResult) {
									t.Log("Applying comparison with NaN equality.")
									equalsWithNaNs(t, oldResult, newResult)
								} else if oldResult.Err != nil {
									testutil.Equals(t, oldResult.Err.Error(), newResult.Err.Error())
								} else {
									testutil.Equals(t, oldResult, newResult)
								}
							})
						}
					})
				}
			}
		})
	}
}

func TestQueryCancellation(t *testing.T) {
	twelveHours := int64(12 * time.Hour.Seconds())

	start := time.Unix(0, 0)
	end := time.Unix(twelveHours, 0)
	step := time.Second * 30
	query := `sum(rate(http_requests_total{pod="nginx-1"}[10s]))`
	load := `load 30s
				http_requests_total{pod="nginx-1"} 1+1x1
				http_requests_total{pod="nginx-2"} 1+2x40`

	test, err := promql.NewTest(t, load)
	testutil.Ok(t, err)
	defer test.Close()

	testutil.Ok(t, test.Run())

	querier := &storage.MockQueryable{
		MockQuerier: &storage.MockQuerier{
			SelectMockFunction: func(sortSeries bool, hints *storage.SelectHints, matchers ...*labels.Matcher) storage.SeriesSet {
				return newTestSeriesSet(&slowSeries{})
			},
		},
	}

	ctx := test.Context()
	newEngine := engine.New(engine.Opts{EngineOpts: promql.EngineOpts{Timeout: 1 * time.Hour}})
	q1, err := newEngine.NewRangeQuery(ctx, querier, nil, query, start, end, step)
	testutil.Ok(t, err)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-time.After(1000 * time.Millisecond)
		cancel()
	}()

	newResult := q1.Exec(ctx)
	testutil.Equals(t, context.Canceled, newResult.Err)
}

func TestQueryTimeout(t *testing.T) {
	end := time.Unix(120, 0)
	query := `http_requests_total{pod="nginx-1"}`
	load := `load 30s
				http_requests_total{pod="nginx-1"} 1+1x1
				http_requests_total{pod="nginx-2"} 1+2x1`

	opts := promql.EngineOpts{
		Timeout:    1 * time.Microsecond,
		MaxSamples: math.MaxInt64,
	}

	test, err := promql.NewTest(t, load)
	testutil.Ok(t, err)
	defer test.Close()

	newEngine := engine.New(engine.Opts{DisableFallback: true, EngineOpts: opts})

	q, err := newEngine.NewInstantQuery(test.Context(), test.Storage(), nil, query, end)
	testutil.Ok(t, err)

	res := q.Exec(context.Background())
	testutil.NotOk(t, res.Err, "expected timeout error but got none")
	testutil.Equals(t, context.DeadlineExceeded, res.Err)
}

type hintRecordingQuerier struct {
	storage.Querier
	mux   sync.Mutex
	hints []*storage.SelectHints
}

func (h *hintRecordingQuerier) Close() error { return nil }

func (h *hintRecordingQuerier) Select(sortSeries bool, hints *storage.SelectHints, matchers ...*labels.Matcher) storage.SeriesSet {
	h.mux.Lock()
	defer h.mux.Unlock()
	h.hints = append(h.hints, hints)
	return storage.EmptySeriesSet()
}

func TestSelectHintsSetCorrectly(t *testing.T) {
	for _, tc := range []struct {
		query string

		// All times are in milliseconds.
		start int64
		end   int64

		// TODO(bwplotka): Add support for better hints when subquerying.
		expected []*storage.SelectHints
	}{
		{
			query: "foo", start: 10000,
			expected: []*storage.SelectHints{
				{Start: 5000, End: 10000},
			},
		}, {
			query: "foo @ 15", start: 10000,
			expected: []*storage.SelectHints{
				{Start: 10000, End: 15000},
			},
		}, {
			query: "foo @ 1", start: 10000,
			expected: []*storage.SelectHints{
				{Start: -4000, End: 1000},
			},
		}, {
			query: "foo[2m]", start: 200000,
			expected: []*storage.SelectHints{
				{Start: 80000, End: 200000, Range: 120000},
			},
		}, {
			query: "foo[2m] @ 180", start: 200000,
			expected: []*storage.SelectHints{
				{Start: 60000, End: 180000, Range: 120000},
			},
		}, {
			query: "foo[2m] @ 300", start: 200000,
			expected: []*storage.SelectHints{
				{Start: 180000, End: 300000, Range: 120000},
			},
		}, {
			query: "foo[2m] @ 60", start: 200000,
			expected: []*storage.SelectHints{
				{Start: -60000, End: 60000, Range: 120000},
			},
		}, {
			query: "foo[2m] offset 2m", start: 300000,
			expected: []*storage.SelectHints{
				{Start: 60000, End: 180000, Range: 120000},
			},
		}, {
			query: "foo[2m] @ 200 offset 2m", start: 300000,
			expected: []*storage.SelectHints{
				{Start: -40000, End: 80000, Range: 120000},
			},
		}, {
			query: "foo[2m:1s]", start: 300000,
			expected: []*storage.SelectHints{
				{Start: 175000, End: 300000, Step: 1000},
			},
		}, {
			query: "count_over_time(foo[2m:1s])", start: 300000,
			expected: []*storage.SelectHints{
				{Start: 175000, End: 300000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: "count_over_time(foo[2m:1s] @ 300)", start: 200000,
			expected: []*storage.SelectHints{
				{Start: 175000, End: 300000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: "count_over_time(foo[2m:1s] @ 200)", start: 200000,
			expected: []*storage.SelectHints{
				{Start: 75000, End: 200000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: "count_over_time(foo[2m:1s] @ 100)", start: 200000,
			expected: []*storage.SelectHints{
				{Start: -25000, End: 100000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: "count_over_time(foo[2m:1s] offset 10s)", start: 300000,
			expected: []*storage.SelectHints{
				{Start: 165000, End: 290000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: "count_over_time((foo offset 10s)[2m:1s] offset 10s)", start: 300000,
			expected: []*storage.SelectHints{
				{Start: 155000, End: 280000, Func: "count_over_time", Step: 1000},
			},
		}, {
			// When the @ is on the vector selector, the enclosing subquery parameters
			// don't affect the hint ranges.
			query: "count_over_time((foo @ 200 offset 10s)[2m:1s] offset 10s)", start: 300000,
			expected: []*storage.SelectHints{
				{Start: 185000, End: 190000, Func: "count_over_time", Step: 1000},
			},
		}, {
			// When the @ is on the vector selector, the enclosing subquery parameters
			// don't affect the hint ranges.
			query: "count_over_time((foo @ 200 offset 10s)[2m:1s] @ 100 offset 10s)", start: 300000,
			expected: []*storage.SelectHints{
				{Start: 185000, End: 190000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: "count_over_time((foo offset 10s)[2m:1s] @ 100 offset 10s)", start: 300000,
			expected: []*storage.SelectHints{
				{Start: -45000, End: 80000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: "foo", start: 10000, end: 20000,
			expected: []*storage.SelectHints{
				{Start: 5000, End: 20000, Step: 1000},
			},
		}, {
			query: "foo @ 15", start: 10000, end: 20000,
			expected: []*storage.SelectHints{
				{Start: 10000, End: 15000, Step: 1000},
			},
		}, {
			query: "foo @ 1", start: 10000, end: 20000,
			expected: []*storage.SelectHints{
				{Start: -4000, End: 1000, Step: 1000},
			},
		}, {
			query: "rate(foo[2m] @ 180)", start: 200000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 60000, End: 180000, Range: 120000, Func: "rate", Step: 1000},
			},
		}, {
			query: "rate(foo[2m] @ 300)", start: 200000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 180000, End: 300000, Range: 120000, Func: "rate", Step: 1000},
			},
		}, {
			query: "rate(foo[2m] @ 60)", start: 200000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: -60000, End: 60000, Range: 120000, Func: "rate", Step: 1000},
			},
		}, {
			query: "rate(foo[2m])", start: 200000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 80000, End: 500000, Range: 120000, Func: "rate", Step: 1000},
			},
		}, {
			query: "rate(foo[2m] offset 2m)", start: 300000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 60000, End: 380000, Range: 120000, Func: "rate", Step: 1000},
			},
		}, {
			query: "rate(foo[2m:1s])", start: 300000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 175000, End: 500000, Func: "rate", Step: 1000},
			},
		}, {
			query: "count_over_time(foo[2m:1s])", start: 300000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 175000, End: 500000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: "count_over_time(foo[2m:1s] offset 10s)", start: 300000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 165000, End: 490000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: "count_over_time(foo[2m:1s] @ 300)", start: 200000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 175000, End: 300000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: "count_over_time(foo[2m:1s] @ 200)", start: 200000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 75000, End: 200000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: "count_over_time(foo[2m:1s] @ 100)", start: 200000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: -25000, End: 100000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: "count_over_time((foo offset 10s)[2m:1s] offset 10s)", start: 300000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 155000, End: 480000, Func: "count_over_time", Step: 1000},
			},
		}, {
			// When the @ is on the vector selector, the enclosing subquery parameters
			// don't affect the hint ranges.
			query: "count_over_time((foo @ 200 offset 10s)[2m:1s] offset 10s)", start: 300000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 185000, End: 190000, Func: "count_over_time", Step: 1000},
			},
		}, {
			// When the @ is on the vector selector, the enclosing subquery parameters
			// don't affect the hint ranges.
			query: "count_over_time((foo @ 200 offset 10s)[2m:1s] @ 100 offset 10s)", start: 300000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 185000, End: 190000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: "count_over_time((foo offset 10s)[2m:1s] @ 100 offset 10s)", start: 300000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: -45000, End: 80000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: "sum by (dim1) (foo)", start: 10000,
			expected: []*storage.SelectHints{
				{Start: 5000, End: 10000, Func: "sum", By: true, Grouping: []string{"dim1"}},
			},
		}, {
			query: "sum without (dim1) (foo)", start: 10000,
			expected: []*storage.SelectHints{
				{Start: 5000, End: 10000, Func: "sum", Grouping: []string{"dim1"}},
			},
		}, {
			query: "sum by (dim1) (avg_over_time(foo[1s]))", start: 10000,
			expected: []*storage.SelectHints{
				{Start: 9000, End: 10000, Func: "avg_over_time", Range: 1000},
			},
		}, {
			query: "sum by (dim1) (max by (dim2) (foo))", start: 10000,
			expected: []*storage.SelectHints{
				{Start: 5000, End: 10000, Func: "max", By: true, Grouping: []string{"dim2"}},
			},
		}, {
			query: "(max by (dim1) (foo))[5s:1s]", start: 10000,
			expected: []*storage.SelectHints{
				{Start: 0, End: 10000, Func: "max", By: true, Grouping: []string{"dim1"}, Step: 1000},
			},
		}, {
			query: "(sum(http_requests{group=~\"p.*\"})+max(http_requests{group=~\"c.*\"}))[20s:5s]", start: 120000,
			expected: []*storage.SelectHints{
				{Start: 95000, End: 120000, Func: "sum", By: true, Step: 5000},
				{Start: 95000, End: 120000, Func: "max", By: true, Step: 5000},
			},
		}, {
			query: "foo @ 50 + bar @ 250 + baz @ 900", start: 100000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 45000, End: 50000, Step: 1000},
				{Start: 245000, End: 250000, Step: 1000},
				{Start: 895000, End: 900000, Step: 1000},
			},
		}, {
			query: "foo @ 50 + bar + baz @ 900", start: 100000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 45000, End: 50000, Step: 1000},
				{Start: 95000, End: 500000, Step: 1000},
				{Start: 895000, End: 900000, Step: 1000},
			},
		}, {
			query: "rate(foo[2s] @ 50) + bar @ 250 + baz @ 900", start: 100000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 48000, End: 50000, Step: 1000, Func: "rate", Range: 2000},
				{Start: 245000, End: 250000, Step: 1000},
				{Start: 895000, End: 900000, Step: 1000},
			},
		}, {
			query: "rate(foo[2s:1s] @ 50) + bar + baz", start: 100000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 43000, End: 50000, Step: 1000, Func: "rate"},
				{Start: 95000, End: 500000, Step: 1000},
				{Start: 95000, End: 500000, Step: 1000},
			},
		}, {
			query: "rate(foo[2s:1s] @ 50) + bar + rate(baz[2m:1s] @ 900 offset 2m) ", start: 100000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 43000, End: 50000, Step: 1000, Func: "rate"},
				{Start: 95000, End: 500000, Step: 1000},
				{Start: 655000, End: 780000, Step: 1000, Func: "rate"},
			},
		}, { // Hints are based on the inner most subquery timestamp.
			query: `sum_over_time(sum_over_time(metric{job="1"}[100s])[100s:25s] @ 50)[3s:1s] @ 3000`, start: 100000,
			expected: []*storage.SelectHints{
				{Start: -150000, End: 50000, Range: 100000, Func: "sum_over_time", Step: 25000},
			},
		}, { // Hints are based on the inner most subquery timestamp.
			query: `sum_over_time(sum_over_time(metric{job="1"}[100s])[100s:25s] @ 3000)[3s:1s] @ 50`,
			expected: []*storage.SelectHints{
				{Start: 2800000, End: 3000000, Range: 100000, Func: "sum_over_time", Step: 25000},
			},
		},
	} {
		t.Run(tc.query, func(t *testing.T) {
			opts := promql.EngineOpts{
				Logger:           nil,
				Reg:              nil,
				MaxSamples:       10,
				Timeout:          10 * time.Second,
				LookbackDelta:    5 * time.Second,
				EnableAtModifier: true,
			}

			ng := engine.New(engine.Opts{EngineOpts: opts})
			hintsRecorder := &hintRecordingQuerier{}
			queryable := &storage.MockQueryable{MockQuerier: hintsRecorder}
			ctx := context.Background()

			var (
				query promql.Query
				err   error
			)
			if tc.end == 0 {
				query, err = ng.NewInstantQuery(ctx, queryable, nil, tc.query, timestamp.Time(tc.start))
			} else {
				query, err = ng.NewRangeQuery(ctx, queryable, nil, tc.query, timestamp.Time(tc.start), timestamp.Time(tc.end), time.Second)
			}
			testutil.Ok(t, err)

			res := query.Exec(context.Background())
			testutil.Ok(t, res.Err)

			// Selects are done in parallel so check that all hints are
			// present, but order does not matter.
			testutil.Equals(t, len(tc.expected), len(hintsRecorder.hints))
			for _, expected := range tc.expected {
				contains := false
				for _, hint := range hintsRecorder.hints {
					if reflect.DeepEqual(expected, hint) {
						contains = true
					}
				}
				testutil.Assert(t, contains, "hints did not contain contain %#v", expected)
			}
		})
	}
}

func TestFallback(t *testing.T) {
	start := time.Unix(0, 0)
	end := time.Unix(120, 0)
	step := time.Second * 30

	// TODO(fpetkovski): Update this expression once we add support for predict_linear.
	query := `predict_linear(http_requests_total{pod="nginx-1"}[5m], 10)`
	load := `load 30s
				http_requests_total{pod="nginx-1"} 1+1x1
				http_requests_total{pod="nginx-2"} 1+2x40`

	test, err := promql.NewTest(t, load)
	testutil.Ok(t, err)
	defer test.Close()

	testutil.Ok(t, test.Run())

	for _, disableFallback := range []bool{true, false} {
		t.Run(fmt.Sprintf("disableFallback=%t", disableFallback), func(t *testing.T) {
			opts := promql.EngineOpts{
				Timeout:    2 * time.Second,
				MaxSamples: math.MaxInt64,
			}
			newEngine := engine.New(engine.Opts{DisableFallback: disableFallback, EngineOpts: opts})
			q1, err := newEngine.NewRangeQuery(test.Context(), test.Storage(), nil, query, start, end, step)
			if disableFallback {
				testutil.NotOk(t, err)
			} else {
				testutil.Ok(t, err)
				newResult := q1.Exec(context.Background())
				testutil.Ok(t, newResult.Err)
			}
		})
	}
}

func TestQueryStats(t *testing.T) {
	start := time.Unix(0, 0)
	end := time.Unix(120, 0)
	step := time.Second * 30

	query := `http_requests_total{pod="nginx-1"}`
	load := `load 30s
				http_requests_total{pod="nginx-1"} 1+1x1
				http_requests_total{pod="nginx-2"} 1+2x1`
	opts := promql.EngineOpts{
		Timeout:    2 * time.Second,
		MaxSamples: math.MaxInt64,
	}

	test, err := promql.NewTest(t, load)
	testutil.Ok(t, err)
	defer test.Close()

	ctx := test.Context()
	newEngine := engine.New(engine.Opts{DisableFallback: true, EngineOpts: opts})
	q, err := newEngine.NewRangeQuery(ctx, test.Storage(), nil, query, start, end, step)
	testutil.Ok(t, err)
	stats.NewQueryStats(q.Stats())

	q, err = newEngine.NewInstantQuery(ctx, test.Storage(), nil, query, end)
	testutil.Ok(t, err)
	stats.NewQueryStats(q.Stats())
}

func storageWithMockSeries(mockSeries ...*mockSeries) *storage.MockQueryable {
	series := make([]storage.Series, 0, len(mockSeries))
	for _, mock := range mockSeries {
		series = append(series, storage.Series(mock))
	}
	return storageWithSeries(series...)
}

func storageWithSeries(series ...storage.Series) *storage.MockQueryable {
	return &storage.MockQueryable{
		MockQuerier: &storage.MockQuerier{
			SelectMockFunction: func(sortSeries bool, hints *storage.SelectHints, matchers ...*labels.Matcher) storage.SeriesSet {
				result := make([]storage.Series, 0)
			loopSeries:
				for _, s := range series {
					for _, m := range matchers {
						lbl := s.Labels().Get(m.Name)
						if lbl != "" && !m.Matches(lbl) {
							continue loopSeries
						}
					}
					result = append(result, s)
				}
				return newTestSeriesSet(result...)
			},
		},
	}
}

type byTimestamps mockSeries

func (b byTimestamps) Len() int {
	return len(b.timestamps)
}

func (b byTimestamps) Less(i, j int) bool {
	return b.timestamps[i] < b.timestamps[j]
}

func (b byTimestamps) Swap(i, j int) {
	b.timestamps[i], b.timestamps[j] = b.timestamps[j], b.timestamps[i]
	b.values[i], b.values[j] = b.values[j], b.values[i]
}

type mockSeries struct {
	labels     []string
	timestamps []int64
	values     []float64
}

func newMockSeries(labels []string, timestamps []int64, values []float64) *mockSeries {
	for i := range timestamps {
		timestamps[i] = timestamps[i] * 1000
	}
	return &mockSeries{labels: labels, timestamps: timestamps, values: values}
}

func (m mockSeries) Labels() labels.Labels {
	return labels.FromStrings(m.labels...)
}

func (m mockSeries) Iterator(chunkenc.Iterator) chunkenc.Iterator {
	return &mockIterator{
		i:          -1,
		timestamps: m.timestamps,
		values:     m.values,
	}
}

type mockIterator struct {
	i          int
	timestamps []int64
	values     []float64
}

func (m *mockIterator) Next() chunkenc.ValueType {
	m.i++
	if m.i >= len(m.values) {
		return chunkenc.ValNone
	}

	return chunkenc.ValFloat
}

func (m *mockIterator) Seek(t int64) chunkenc.ValueType {
	if m.i > -1 && m.i < len(m.timestamps) {
		currentTS := m.timestamps[m.i]
		if currentTS >= t {
			return chunkenc.ValFloat
		}
	}
	for {
		next := m.Next()
		if next == chunkenc.ValNone {
			return chunkenc.ValNone
		}

		if m.AtT() >= t {
			return next
		}
	}
}

func (m *mockIterator) At() (int64, float64) {
	return m.timestamps[m.i], m.values[m.i]
}

func (m *mockIterator) AtHistogram() (int64, *histogram.Histogram) { return 0, nil }

func (m *mockIterator) AtFloatHistogram() (int64, *histogram.FloatHistogram) { return 0, nil }

func (m *mockIterator) AtT() int64 { return m.timestamps[m.i] }

func (m *mockIterator) Err() error { return nil }

type testSeriesSet struct {
	i      int
	series []storage.Series
}

func newTestSeriesSet(series ...storage.Series) storage.SeriesSet {
	return &testSeriesSet{
		i:      -1,
		series: series,
	}
}

func (s *testSeriesSet) Next() bool                 { s.i++; return s.i < len(s.series) }
func (s *testSeriesSet) At() storage.Series         { return s.series[s.i] }
func (s *testSeriesSet) Err() error                 { return nil }
func (s *testSeriesSet) Warnings() storage.Warnings { return nil }

type slowSeries struct{}

func (d slowSeries) Labels() labels.Labels                        { return labels.FromStrings("foo", "bar") }
func (d slowSeries) Iterator(chunkenc.Iterator) chunkenc.Iterator { return &slowIterator{} }

type slowIterator struct {
	ts int64
}

func (d *slowIterator) AtHistogram() (int64, *histogram.Histogram) {
	panic("not implemented")
}

func (d *slowIterator) AtFloatHistogram() (int64, *histogram.FloatHistogram) {
	panic("not implemented")
}

func (d *slowIterator) AtT() int64 {
	return d.ts
}

func (d *slowIterator) At() (int64, float64) {
	return d.ts, 1
}

func (d *slowIterator) Next() chunkenc.ValueType {
	<-time.After(10 * time.Millisecond)
	d.ts += 30 * 1000
	return chunkenc.ValFloat
}

func (d *slowIterator) Seek(t int64) chunkenc.ValueType {
	<-time.After(10 * time.Millisecond)
	d.ts = t
	return chunkenc.ValFloat
}
func (d *slowIterator) Err() error { return nil }

type mockRuntimeErr struct{}

func (m *mockRuntimeErr) Error() string {
	return "panic!"
}

func (m *mockRuntimeErr) RuntimeError() {
}

func TestEngineRecoversFromPanic(t *testing.T) {
	t.Parallel()

	querier := &storage.MockQueryable{
		MockQuerier: &storage.MockQuerier{
			SelectMockFunction: func(sortSeries bool, hints *storage.SelectHints, matchers ...*labels.Matcher) storage.SeriesSet {
				panic(runtime.Error(&mockRuntimeErr{}))
			},
		},
	}
	t.Run("instant", func(t *testing.T) {
		newEngine := engine.New(engine.Opts{
			DisableFallback: true,
		})
		ctx := context.Background()
		q, err := newEngine.NewInstantQuery(ctx, querier, nil, "somequery", time.Time{})
		testutil.Ok(t, err)

		r := q.Exec(ctx)
		testutil.Assert(t, r.Err.Error() == "unexpected error: panic!")
	})

	t.Run("range", func(t *testing.T) {
		newEngine := engine.New(engine.Opts{
			DisableFallback: true,
		})
		ctx := context.Background()
		q, err := newEngine.NewRangeQuery(ctx, querier, nil, "somequery", time.Time{}, time.Time{}, 42)
		testutil.Ok(t, err)

		r := q.Exec(ctx)
		testutil.Assert(t, r.Err.Error() == "unexpected error: panic!")
	})

}

type histogramTestCase struct {
	name                   string
	query                  string
	wantEmptyForMixedTypes bool
}

type histogramGeneratorFunc func(app storage.Appender, numSeries int, withMixedTypes bool) error

func TestNativeHistograms(t *testing.T) {
	opts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}

	cases := []histogramTestCase{
		{
			name:  "plain selector",
			query: "native_histogram_series",
		},
		{
			name:  "irate()",
			query: "irate(native_histogram_series[1m])",
		},
		{
			name:  "rate()",
			query: "rate(native_histogram_series[1m])",
		},
		{
			name:  "increase()",
			query: "increase(native_histogram_series[1m])",
		},
		{
			name:  "delta()",
			query: "delta(native_histogram_series[1m])",
		},
		{
			name:                   "sum()",
			query:                  "sum(native_histogram_series)",
			wantEmptyForMixedTypes: true,
		},
		{
			name:                   "sum by (foo)",
			query:                  "sum by (foo) (native_histogram_series)",
			wantEmptyForMixedTypes: true,
		},
		{
			name:  "count",
			query: "count (native_histogram_series)",
		},
		{
			name:  "count by (foo)",
			query: "count by (foo) (native_histogram_series)",
		},
		// TODO(fpetkovski): The Prometheus engine returns an incorrect result for this case.
		// Uncomment once it gets fixed: https://github.com/prometheus/prometheus/issues/11973.
		// {
		//	name:  "max",
		//	query: "max (native_histogram_series)",
		// },
		{
			name:  "max by (foo)",
			query: "max by (foo) (native_histogram_series)",
		},
		// TODO(fpetkovski): The Prometheus engine returns an incorrect result for this case.
		// Uncomment once it gets fixed: https://github.com/prometheus/prometheus/issues/11973.
		// {
		//	name:  "min",
		//	query: "min (native_histogram_series)",
		// },
		{
			name:  "min by (foo)",
			query: "min by (foo) (native_histogram_series)",
		},
		{
			name:  "histogram_sum",
			query: "histogram_sum(native_histogram_series)",
		},
		{
			name:  "histogram_count",
			query: "histogram_count(native_histogram_series)",
		},
		{
			name:  "histogram_quantile",
			query: "histogram_quantile(0.7, native_histogram_series)",
		},
		{
			name:  "histogram_fraction",
			query: "histogram_fraction(0, 0.2, native_histogram_series)",
		},
	}

	t.Run("integer_histograms", func(t *testing.T) {
		testNativeHistograms(t, cases, opts, generateNativeHistogramSeries)
	})
	t.Run("float_histograms", func(t *testing.T) {
		testNativeHistograms(t, cases, opts, generateFloatHistogramSeries)
	})
}

func testNativeHistograms(t *testing.T, cases []histogramTestCase, opts promql.EngineOpts, generateHistograms histogramGeneratorFunc) {
	numHistograms := 100
	mixedTypesOpts := []bool{false, true}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, withMixedTypes := range mixedTypesOpts {
				t.Run(fmt.Sprintf("mixedTypes=%t", withMixedTypes), func(t *testing.T) {
					test, err := promql.NewTest(t, "")
					testutil.Ok(t, err)
					defer test.Close()

					app := test.Storage().Appender(context.TODO())
					err = generateHistograms(app, numHistograms, withMixedTypes)
					testutil.Ok(t, err)
					testutil.Ok(t, app.Commit())
					testutil.Ok(t, test.Run())

					// New Engine
					engine := engine.New(engine.Opts{
						EngineOpts:        opts,
						DisableFallback:   true,
						LogicalOptimizers: logicalplan.AllOptimizers,
					})

					t.Run("instant", func(t *testing.T) {
						ctx := test.Context()
						qry, err := engine.NewInstantQuery(ctx, test.Queryable(), nil, tc.query, time.Unix(50, 0))
						testutil.Ok(t, err)
						newResult := qry.Exec(test.Context())
						testutil.Ok(t, newResult.Err)
						newVector, err := newResult.Vector()
						testutil.Ok(t, err)

						promEngine := test.QueryEngine()
						qry, err = promEngine.NewInstantQuery(ctx, test.Queryable(), nil, tc.query, time.Unix(50, 0))
						testutil.Ok(t, err)
						promResult := qry.Exec(test.Context())
						testutil.Ok(t, promResult.Err)
						promVector, err := promResult.Vector()
						testutil.Ok(t, err)

						// Make sure we're not getting back empty results.
						if withMixedTypes && tc.wantEmptyForMixedTypes {
							testutil.Assert(t, len(promVector) == 0)
						}

						sortByLabels(promResult)
						sortByLabels(newResult)
						if hasNaNs(promResult) {
							t.Log("Applying comparison with NaN equality.")
							equalsWithNaNs(t, promVector, newVector)
						} else {
							testutil.Equals(t, promVector, newVector)
						}
					})

					t.Run("range", func(t *testing.T) {
						qry, err := engine.NewRangeQuery(test.Context(), test.Queryable(), nil, tc.query, time.Unix(50, 0), time.Unix(600, 0), 30*time.Second)
						testutil.Ok(t, err)
						res := qry.Exec(test.Context())
						testutil.Ok(t, res.Err)
						actual, err := res.Matrix()
						testutil.Ok(t, err)

						promEngine := test.QueryEngine()
						qry, err = promEngine.NewRangeQuery(test.Context(), test.Queryable(), nil, tc.query, time.Unix(50, 0), time.Unix(600, 0), 30*time.Second)
						testutil.Ok(t, err)
						res = qry.Exec(test.Context())
						testutil.Ok(t, res.Err)
						expected, err := res.Matrix()
						testutil.Ok(t, err)

						// Make sure we're not getting back empty results.
						if withMixedTypes && tc.wantEmptyForMixedTypes {
							testutil.Assert(t, len(expected) == 0)
						}
						testutil.Equals(t, len(expected), len(actual))
						if hasNaNs(res) {
							t.Log("Applying comparison with NaN equality.")
							equalsWithNaNs(t, expected, actual)
						} else {
							testutil.Equals(t, expected, actual)
						}
					})
				})
			}
		})
	}
}

func generateNativeHistogramSeries(app storage.Appender, numSeries int, withMixedTypes bool) error {
	commonLabels := []string{labels.MetricName, "native_histogram_series", "foo", "bar"}
	series := make([][]*histogram.Histogram, numSeries)
	for i := range series {
		series[i] = tsdbutil.GenerateTestHistograms(2000)
	}
	higherSchemaHist := &histogram.Histogram{
		Schema: 3,
		PositiveSpans: []histogram.Span{
			{Offset: -5, Length: 2}, // -5 -4
			{Offset: 2, Length: 3},  // -1 0 1
			{Offset: 2, Length: 2},  // 4 5
		},
		PositiveBuckets: []int64{1, 2, -2, 1, -1, 0, 3},
		Count:           13,
	}
	for sid, histograms := range series {
		lbls := append(commonLabels, "h", strconv.Itoa(sid))
		for i := range histograms {
			ts := time.Unix(int64(i*15), 0).UnixMilli()
			if i == 0 {
				// Inject a histogram with a higher schema.
				// Regression test for:
				// * https://github.com/thanos-io/promql-engine/pull/182
				// * https://github.com/thanos-io/promql-engine/pull/183.
				if _, err := app.AppendHistogram(0, labels.FromStrings(lbls...), ts, higherSchemaHist, nil); err != nil {
					return err
				}
			}
			if _, err := app.AppendHistogram(0, labels.FromStrings(lbls...), ts, histograms[i], nil); err != nil {
				return err
			}
			if withMixedTypes {
				if _, err := app.Append(0, labels.FromStrings(append(lbls, "le", "1")...), ts, float64(i)); err != nil {
					return err
				}
				if _, err := app.Append(0, labels.FromStrings(append(lbls, "le", "+Inf")...), ts, float64(i*2)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func generateFloatHistogramSeries(app storage.Appender, numSeries int, withMixedTypes bool) error {
	lbls := []string{labels.MetricName, "native_histogram_series", "foo", "bar"}
	h1 := tsdbutil.GenerateTestFloatHistograms(numSeries)
	h2 := tsdbutil.GenerateTestFloatHistograms(numSeries)
	for i := range h1 {
		ts := time.Unix(int64(i*15), 0).UnixMilli()
		if withMixedTypes {
			if _, err := app.Append(0, labels.FromStrings(append(lbls, "le", "1")...), ts, float64(i)); err != nil {
				return err
			}
			if _, err := app.Append(0, labels.FromStrings(append(lbls, "le", "+Inf")...), ts, float64(i*2)); err != nil {
				return err
			}
		}
		if _, err := app.AppendHistogram(0, labels.FromStrings(append(lbls, "h", "1")...), ts, nil, h1[i]); err != nil {
			return err
		}
		if _, err := app.AppendHistogram(0, labels.FromStrings(append(lbls, "h", "2")...), ts, nil, h2[i]); err != nil {
			return err
		}
	}
	return nil
}

func TestMixedNativeHistogramTypes(t *testing.T) {
	histograms := tsdbutil.GenerateTestHistograms(2)

	test, err := promql.NewTest(t, "")
	testutil.Ok(t, err)
	defer test.Close()

	lbls := []string{labels.MetricName, "native_histogram_series"}

	app := test.Storage().Appender(context.TODO())
	_, err = app.AppendHistogram(0, labels.FromStrings(lbls...), 0, nil, histograms[0].ToFloat())
	testutil.Ok(t, err)
	testutil.Ok(t, app.Commit())

	app = test.Storage().Appender(context.TODO())
	_, err = app.AppendHistogram(0, labels.FromStrings(lbls...), 30_000, histograms[1], nil)
	testutil.Ok(t, err)
	testutil.Ok(t, app.Commit())

	testutil.Ok(t, test.Run())

	opts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}

	engine := engine.New(engine.Opts{
		EngineOpts:        opts,
		DisableFallback:   true,
		LogicalOptimizers: logicalplan.AllOptimizers,
	})

	t.Run("vector_select", func(t *testing.T) {
		qry, err := engine.NewInstantQuery(test.Context(), test.Queryable(), nil, "sum(native_histogram_series)", time.Unix(30, 0))
		testutil.Ok(t, err)
		res := qry.Exec(context.Background())
		testutil.Ok(t, res.Err)
		actual, err := res.Vector()
		testutil.Ok(t, err)

		testutil.Equals(t, 1, len(actual), "expected vector with 1 element")
		expected := histograms[1].ToFloat()
		expected.CounterResetHint = histogram.UnknownCounterReset
		testutil.Equals(t, expected, actual[0].H)
	})

	t.Run("matrix_select", func(t *testing.T) {
		qry, err := engine.NewRangeQuery(test.Context(), test.Queryable(), nil, "rate(native_histogram_series[1m])", time.Unix(0, 0), time.Unix(60, 0), 60*time.Second)
		testutil.Ok(t, err)
		res := qry.Exec(context.Background())
		testutil.Ok(t, res.Err)
		actual, err := res.Matrix()
		testutil.Ok(t, err)

		testutil.Equals(t, 1, len(actual), "expected 1 series")
		testutil.Equals(t, 1, len(actual[0].Histograms), "expected 1 point")
		expected := histograms[1].ToFloat().Sub(histograms[0].ToFloat()).Mul(1 / float64(30))
		expected.CounterResetHint = histogram.GaugeType
		testutil.Equals(t, expected, actual[0].Histograms[0].H)
	})
}

func sortByLabels(r *promql.Result) {
	if r.Err != nil {
		return
	}
	switch r.Value.Type() {
	case promparser.ValueTypeVector:
		m, _ := r.Vector()
		sort.Sort(samplesByLabels(m))
		r.Value = m
	case promparser.ValueTypeMatrix:
		m, _ := r.Matrix()
		sort.Sort(seriesByLabels(m))
		r.Value = m
	}
}

type seriesByLabels []promql.Series

func (b seriesByLabels) Len() int           { return len(b) }
func (b seriesByLabels) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b seriesByLabels) Less(i, j int) bool { return labels.Compare(b[i].Metric, b[j].Metric) < 0 }

type samplesByLabels []promql.Sample

func (b samplesByLabels) Len() int           { return len(b) }
func (b samplesByLabels) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b samplesByLabels) Less(i, j int) bool { return labels.Compare(b[i].Metric, b[j].Metric) < 0 }

// roundValues rounds all values to 10 decimal points and
// can be used to eliminate floating point division errors
// when comparing two promql results.
func roundValues(r *promql.Result) {
	switch result := r.Value.(type) {
	case promql.Matrix:
		for i := range result {
			for j := range result[i].Floats {
				result[i].Floats[j].F = math.Floor(result[i].Floats[j].F*1e10) / 1e10
			}
		}
	case promql.Vector:
		for i := range result {
			result[i].F = math.Floor(result[i].F*10e10) / 10e10
		}
	}
}

// emptyLabelsToNil sets empty labelsets to nil to work around inconsistent
// results from the old engine depending on the literal type (e.g. number vs. compare).
func emptyLabelsToNil(result *promql.Result) {
	if value, ok := result.Value.(promql.Matrix); ok {
		for i, s := range value {
			if s.Metric.IsEmpty() {
				result.Value.(promql.Matrix)[i].Metric = labels.EmptyLabels()
			}
		}
	}
}
