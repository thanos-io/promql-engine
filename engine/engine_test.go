// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"runtime"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/efficientgo/core/errors"
	"github.com/efficientgo/core/testutil"
	"github.com/go-kit/log"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"golang.org/x/exp/maps"

	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/execution/warnings"
	"github.com/thanos-io/promql-engine/logicalplan"
	"github.com/thanos-io/promql-engine/query"
	"github.com/thanos-io/promql-engine/storage/prometheus"

	"github.com/prometheus/prometheus/model/histogram"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/timestamp"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/promqltest"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	"github.com/prometheus/prometheus/tsdb/tsdbutil"
	"github.com/prometheus/prometheus/util/annotations"
	"github.com/prometheus/prometheus/util/stats"
	"github.com/prometheus/prometheus/util/teststorage"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestPromqlAcceptance(t *testing.T) {
	t.Parallel()
	engine := engine.New(engine.Opts{
		EngineOpts: promql.EngineOpts{
			EnableAtModifier:         true,
			EnableNegativeOffset:     true,
			MaxSamples:               5e10,
			Timeout:                  1 * time.Hour,
			NoStepSubqueryIntervalFn: func(rangeMillis int64) int64 { return 30 * time.Second.Milliseconds() },
		}})
	promqltest.RunBuiltinTests(t, engine)
}

func TestVectorSelectorWithGaps(t *testing.T) {
	t.Parallel()
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

	testutil.WithGoCmp(comparer).Equals(t, oldResult, newResult, queryExplanation(q1))
}

func TestQueriesAgainstOldEngine(t *testing.T) {
	t.Parallel()
	start := time.Unix(0, 0)
	end := time.Unix(1800, 0)
	step := time.Second * 30
	// Negative offset and at modifier are enabled by default
	// since Prometheus v2.33.0 so we also enable them.
	opts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
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
			name: "predict_linear fuzz",
			load: `load 30s
			    http_requests_total{pod="nginx-1", route="/"} 48.00+9.17x40
			    http_requests_total{pod="nginx-2", route="/"} -108+173.00x40`,
			query: `predict_linear(http_requests_total{route="/"}[1h:1m] offset 1m, 60)`,
		},
		{
			name: "duplicate label fuzz",
			load: `load 30s
			    http_requests_total{pod="nginx-1", route="/"} 41.00+0.20x40
			    http_requests_total{pod="nginx-2", route="/"} 51+21.71x40`,
			query: `
-avg by (__name__) (
  (-group({__name__="http_requests_total"} @ 54.013) or {__name__="http_requests_total"} offset 1m32s)
)`,
		},
		{
			name: "timestamp fuzz 1",
			load: `load 30s
			    http_requests_total{pod="nginx-1", route="/"} 0.20+9.00x40
			    http_requests_total{pod="nginx-2", route="/"}  6+60.00x40`,
			query: `timestamp(last_over_time(http_requests_total{route="/"}[1h]))`,
		},
		{
			name: "timestamp fuzz 2",
			load: `load 30s
			    http_requests_total{pod="nginx-1", route="/"} 8.00+9.17x40
			    http_requests_total{pod="nginx-2", route="/"} -12+103.00x40`,
			query: `
timestamp(
  http_requests_total{pod="nginx-1"} >= bool (http_requests_total < 2 * http_requests_total)
)`,
		},
		{
			name: "timestamp with multiple parenthesis",
			load: `load 30s
			    http_requests_total{pod="nginx-1", route="/"} 8.00+9.17x40
			    http_requests_total{pod="nginx-2", route="/"} -12+103.00x40`,
			query: `timestamp((http_requests_total))`,
		},
		{
			name: "subqueries in binary expression",
			load: `load 30s
			    http_requests_total{pod="nginx-1", route="/"} 1.00+0.20x40
			    http_requests_total{pod="nginx-2", route="/"} -44+2.00x40`,
			query: `
  absent_over_time(http_requests_total @ end()[1h:1m])
or
  avg_over_time(http_requests_total @ end()[1h:1m])`,
		},

		{
			name:  "nested unary negation",
			query: `1 / (-(2 * 2))`,
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
			query: `stddev by (route) (http_requests_total)`,
		},
		{
			name: "stddev with NaN 2",
			load: `load 30s
			    http_requests_total{pod="nginx-1", route="/"} NaN
			    http_requests_total{pod="nginx-2", route="/"} 1`,
			query: `stddev by (pod) (http_requests_total)`,
		},
		{
			name: "aggregate without",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1.1x1
			    http_requests_total{pod="nginx-2"} 2+2.3x1`,
			start: time.Unix(0, 0),
			end:   time.Unix(60, 0),
			step:  30 * time.Second,
			query: `avg without (pod) (http_requests_total)`,
		},
		{
			name: "func with scalar arg that selects storage, checks whether same series handled correctly",
			load: `load 30s
			    thanos_cache_redis_hits_total{name="caching-bucket",service="thanos-store"} 1+1x30`,
			query: `
  clamp_min(thanos_cache_redis_hits_total, scalar(max by (service) (thanos_cache_redis_hits_total)))
+
  clamp_min(thanos_cache_redis_hits_total, scalar(max by (service) (thanos_cache_redis_hits_total)))`,
		},
		{
			name: "sum + rate divided by itself",
			load: `load 30s
			    thanos_cache_redis_hits_total{name="caching-bucket",service="thanos-store"} 1+1x30`,
			query: `
  (sum by (service) (rate(thanos_cache_redis_hits_total{name="caching-bucket"}[2m])))
/
  (sum by (service) (rate(thanos_cache_redis_hits_total{name="caching-bucket"}[2m])))`,
		},
		{
			name: "stddev_over_time",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `stddev_over_time(http_requests_total[30s])`,
		},
		{
			name: "stdvar_over_time",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `stdvar_over_time(http_requests_total[30s])`,
		},
		{
			name: "quantile_over_time",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `quantile_over_time(0.9, http_requests_total[1m])`,
		},
		{
			name: "quantile_over_time with subquery",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 41.00+0.20x40
			    http_requests_total{pod="nginx-2"} 51+21.71x40`,
			query: `quantile_over_time(0.5, http_requests_total{pod="nginx-1"}[5m:1m])`,
			start: start,
			end:   end,
		},
		{
			name: "quantile_over_time with subquery and non-constant param",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 41.00+0.20x40
			    http_requests_total{pod="nginx-2"} 51+21.71x40
			    param_series 0+0.01x40`,
			query: `quantile_over_time(scalar(param_series), http_requests_total{pod="nginx-1"}[5m:1m])`,
			start: start,
			end:   end,
		},
		{
			name: "predict_linear with subquery and non-constant param",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 41.00+0.20x40
			    http_requests_total{pod="nginx-2"} 51+21.71x40
			    param_series 1+1x40`,
			query: `predict_linear(http_requests_total{pod="nginx-1"}[5m:1m], scalar(param_series))`,
			start: start,
			end:   end,
		},
		{
			name: "predict_linear with subquery and non-existing param series",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 41.00+0.20x40
			    http_requests_total{pod="nginx-2"} 51+21.71x40`,
			query: `predict_linear(http_requests_total{pod="nginx-1"}[5m:1m], scalar(non_existent))`,
			start: start,
			end:   end,
		},
		{
			name: "changes",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18
			    http_requests_total{pod="nginx-2"} 1+2x18
			    http_requests_total{pod="nginx-2"} 1+2x18
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `changes(http_requests_total[30s])`,
		},
		{
			name: "deriv",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18
			    http_requests_total{pod="nginx-2"} 1+2x18
			    http_requests_total{pod="nginx-2"} 1+2x18
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `deriv(http_requests_total[30s])`,
		},
		{
			name: "abs",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} -5+1x15
			    http_requests_total{pod="nginx-2"} -5+2x18`,
			query: `abs(http_requests_total)`,
		},
		{
			name: "ceil",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} -5.5+1x15
			    http_requests_total{pod="nginx-2"} -5.5+2x18`,
			query: `ceil(http_requests_total)`,
		},
		{
			name: "exp",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} -5.5+1x15
			    http_requests_total{pod="nginx-2"} -5.5+2x18`,
			query: `exp(http_requests_total)`,
		},
		{
			name: "floor",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} -5.5+1x15
			    http_requests_total{pod="nginx-2"} -5.5+2x18`,
			query: `floor(http_requests_total)`,
		},
		{
			name: "floor with a filter",
			load: `load 30s
			    http_requests_total{pod="nginx-1", route="/"} 1
			    http_requests_total{pod="nginx-2", route="/"} 2`,
			query: `floor(http_requests_total{pod="nginx-2"}) / http_requests_total`,
		},
		{
			name: "sqrt",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 5.5+1x15
			    http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: `sqrt(http_requests_total)`,
		},
		{
			name: "ln",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 5.5+1x15
			    http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: `ln(http_requests_total)`,
		},
		{
			name: "log2",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 5.5+1x15
			    http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: `log2(http_requests_total)`,
		},
		{
			name: "log10",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 5.5+1x15
			    http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: `log10(http_requests_total)`,
		},
		{
			name: "sin",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 5.5+1x15
			    http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: `sin(http_requests_total)`,
		},
		{
			name: "cos",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 5.5+1x15
			    http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: `cos(http_requests_total)`,
		},
		{
			name: "tan",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 5.5+1x15
			    http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: `tan(http_requests_total)`,
		},
		{
			name: "asin",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 0
			    http_requests_total{pod="nginx-2"} 1`,
			query: `asin(http_requests_total)`,
		},
		{
			name: "acos",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 0
			    http_requests_total{pod="nginx-2"} 1`,
			query: `acos(http_requests_total)`,
		},
		{
			name: "atan",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 0
			    http_requests_total{pod="nginx-2"} 1`,
			query: `atan(http_requests_total)`,
		},
		{
			name: "sinh",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 0
			    http_requests_total{pod="nginx-2"} 1`,
			query: `sinh(http_requests_total)`,
		},
		{
			name: "cosh",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 0
			    http_requests_total{pod="nginx-2"} 1`,
			query: `cosh(http_requests_total)`,
		},
		{
			name: "tanh",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 5.5+1x15
			    http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: `tanh(http_requests_total)`,
		},
		{
			name: "asinh",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 5.5+1x15
			    http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: `asinh(http_requests_total)`,
		},
		{
			name: "acosh",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 5.5+1x15
			    http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: `acosh(http_requests_total)`,
		},
		{
			name: "atanh",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 0
			    http_requests_total{pod="nginx-2"} 1`,
			query: `atanh(http_requests_total)`,
		},
		{
			name: "rad",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 5.5+1x15
			    http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: `rad(http_requests_total)`,
		},
		{
			name: "deg",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 5.5+1x15
			    http_requests_total{pod="nginx-2"} 5.5+2x18`,
			query: `deg(http_requests_total)`,
		},
		{
			name:  "pi",
			load:  ``,
			query: `pi()`,
		},
		{
			name: "sum",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `sum(http_requests_total)`,
		},
		{
			name: "sum_over_time",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `sum_over_time(http_requests_total[30s])`,
		},
		{
			name: "count",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `count(http_requests_total)`,
		},
		{
			name: "count_over_time",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `count_over_time(http_requests_total[30s])`,
		},
		{
			name: "average",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `avg(http_requests_total)`,
		},
		{
			name: "avg_over_time",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `avg_over_time(http_requests_total[30s])`,
		},
		{
			name: "abs",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} -10+1x15
			    http_requests_total{pod="nginx-2"} -10+2x18`,
			query: `abs(http_requests_total)`,
		},
		{
			name: "max",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `max(http_requests_total)`,
		},
		{
			name: "max with only 1 sample",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} -1
			    http_requests_total{pod="nginx-2"} 1`,
			query: `max by (pod) (http_requests_total)`,
		},
		{
			name: "max_over_time",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `max_over_time(http_requests_total[30s])`,
		},
		{
			name: "min",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `min(http_requests_total)`,
		},
		{
			name: "min with only 1 sample",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} -1
			    http_requests_total{pod="nginx-2"} 1`,
			query: `min by (pod) (http_requests_total)`,
		},
		{
			name: "min_over_time",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `min_over_time(http_requests_total[30s])`,
		},
		{
			name: "count_over_time",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `count_over_time(http_requests_total[30s])`,
		},
		{
			name: "sum by pod",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18
			    http_requests_total{pod="nginx-3"} 1+2x20
			    http_requests_total{pod="nginx-4"} 1+2x50`,
			query: `sum by (pod) (http_requests_total)`,
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
			query: `sum by (pod) (http_requests_total)`,
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
			query: `rate(http_requests_total[1m])`,
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "sum rate",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x4
			    http_requests_total{pod="nginx-2"} 1+2x4`,
			query: `sum(rate(http_requests_total[1m]))`,
		},
		{
			name: "sum rate with stale series",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x40
			    http_requests_total{pod="nginx-2"} 1+2x50
			    http_requests_total{pod="nginx-4"} 1+2x50
			    http_requests_total{pod="nginx-5"} 1+2x50
			    http_requests_total{pod="nginx-6"} 1+2x50`,
			query: `sum(rate(http_requests_total[1m]))`,
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
			query: `delta(http_requests_total[1m])`,
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
			query: `increase(http_requests_total[1m])`,
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
			query: `irate(http_requests_total[1m])`,
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
			query: `idelta(http_requests_total[1m])`,
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name:  "number literal",
			load:  "",
			query: `34`,
		},
		{
			name:  "vector",
			load:  "",
			query: `vector(24)`,
		},
		{
			name: "binary operation atan2",
			load: `load 30s
			    foo{} 10
			    bar{} 2`,
			query: `foo atan2 bar`,
		},
		{
			name: "binary operation atan2 with NaN",
			load: `load 30s
			    foo{} 10
			    bar{} NaN`,
			query: `foo atan2 bar`,
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
			query: `foo{code="500"} + ignoring (code) bar`,
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
			query: `foo * ignoring (path, code) group_left () bar`,
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
			query: `bar * ignoring (code, path) group_right () foo`,
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
			query: `foo * ignoring (code, path) group_left (path) bar`,
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
			query: `bar * ignoring (code, path) group_right (path) foo`,
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
			query: `foo + on (code) bar`,
		},
		{
			name: "binary operation with many-to-many matching lhs high card",
			load: `load 30s
			    foo{code="200", method="get"} 1+1x20
			    foo{code="200", method="post"} 1+1x20
			    bar{code="200", method="get"} 1+1x20
			    bar{code="200", method="post"} 1+1x20`,
			query: `foo + on (code) group_left () bar`,
		},
		{
			name: "binary operation with many-to-many matching rhs high card",
			load: `load 30s
			    foo{code="200", method="get"} 1+1x20
			    foo{code="200", method="post"} 1+1x20
			    bar{code="200", method="get"} 1+1x20
			    bar{code="200", method="post"} 1+1x20`,
			query: `foo + on (code) group_right () bar`,
		},
		{
			name: "vector binary op ==",
			load: `load 30s
			    foo{method="get", code="500"} 1+1x40
			    bar{method="get", code="404"} 1+1.1x30`,
			query: `sum by (method) (foo) == sum by (method) (bar)`,
		},
		{
			name: "vector binary op !=",
			load: `load 30s
			    foo{method="get", code="500"} 1+1x40
			    bar{method="get", code="404"} 1+1.1x30`,
			query: `sum by (method) (foo) != sum by (method) (bar)`,
		},
		{
			name: "vector binary op >",
			load: `load 30s
			    foo{method="get", code="500"} 1+1x40
			    bar{method="get", code="404"} 1+1.1x30`,
			query: `sum by (method) (foo) > sum by (method) (bar)`,
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
			query: `sum by (method) (foo) > 10`,
		},
		{
			name: "vector binary op > scalar and bool modifier",
			load: `load 30s
			    foo{method="get", code="500"} 1+2x40
			    bar{method="get", code="404"} 1+1x30`,
			query: `sum by (method) (foo) > bool 10`,
		},
		{
			name: "scalar < vector binary op",
			load: `load 30s
			    foo{method="get", code="500"} 1+2x40
			    bar{method="get", code="404"} 1+1x30`,
			query: `10 < sum by (method) (foo)`,
		},
		{
			name: "vector binary op <",
			load: `load 30s
			    foo{method="get", code="500"} 1+1x40
			    bar{method="get", code="404"} 1+1.1x30`,
			query: `sum by (method) (foo) < sum by (method) (bar)`,
		},
		{
			name: "vector binary op >=",
			load: `load 30s
			    foo{method="get", code="500"} 1+1x40
			    bar{method="get", code="404"} 1+1.1x30`,
			query: `sum by (method) (foo) >= sum by (method) (bar)`,
		},
		{
			name: "vector binary op <=",
			load: `load 30s
			    foo{method="get", code="500"} 1+1x40
			    bar{method="get", code="404"} 1+1.1x30`,
			query: `sum by (method) (foo) <= sum by (method) (bar)`,
		},
		{
			name: "vector binary op ^",
			load: `load 30s
			    foo{method="get", code="500"} 1+1x40
			    bar{method="get", code="404"} 1+1.1x30`,
			query: `sum by (method) (foo) ^ sum by (method) (bar)`,
		},
		{
			name: "vector binary op %",
			load: `load 30s
			    foo{method="get", code="500"} 1+2x40
			    bar{method="get", code="404"} 1+1x30`,
			query: `sum by (method) (foo) % sum by (method) (bar)`,
		},
		{
			name: "vector/vector binary op",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18
			    http_requests_total{pod="nginx-2"} 1+2x18
			    http_requests_total{pod="nginx-2"} 1+2x18
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `(1 + rate(http_requests_total[30s])) > bool rate(http_requests_total[30s])`,
		},
		{
			name: "vector/scalar binary op with a complicated expression on LHS",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18
			    http_requests_total{pod="nginx-2"} 1+2x18
			    http_requests_total{pod="nginx-2"} 1+2x18
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `rate(http_requests_total[30s]) > bool 0`,
		},
		{
			name: "vector/scalar binary op with a complicated expression on RHS",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18
			    http_requests_total{pod="nginx-2"} 1+2x18
			    http_requests_total{pod="nginx-2"} 1+2x18
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `0 < bool rate(http_requests_total[30s])`,
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
			query: `http_requests_total`,
		},
		{
			name:  "time function",
			load:  "",
			query: `time()`,
		},
		{
			name:  "time function in binary expression",
			load:  "",
			query: `time() - 10`,
		},
		{
			name:  "empty series with func",
			load:  "",
			query: `sum(http_requests_total)`,
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
			    http_requests_total{pod="nginx-1"} 2+1x15
			    http_requests_total{pod="nginx-2"} 2+2x18`,
			query: `group(http_requests_total)`,
		},
		{
			name: "group by ",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 2+1x15
			    http_requests_total{pod="nginx-2"} 2+2x18`,
			query: `group by (pod) (http_requests_total)`,
		},
		{
			// Issue https://github.com/thanos-io/promql-engine/issues/326.
			name: "group by with NaN values",
			load: `load 30s
			    http_requests_total{pod="nginx-1", route="/"} 1.00+1.00x4
			    http_requests_total{pod="nginx-2", route="/"}  1+2.00x4`,
			query: `group by (pod, route) (atanh(-{__name__="http_requests_total"} offset -3m4s))`,
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
      rate(grpc_server_handled_total{grpc_method="Series",pod=~".+"}[1m])
    )
  + on (pod) group_left ()
    max by (pod) (prometheus_tsdb_head_samples_appended_total{pod=~".+"})
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
			query: `http_requests_total @ 10.000`,
		},
		{
			name: "@ vector time 120s",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total @ 120.000`,
		},
		{
			name: "@ vector time 360s",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total @ 360.000`,
		},
		{
			name: "@ vector start",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total @ start()`,
		},
		{
			name: "@ vector end",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total @ end()`,
		},
		{
			name: "count_over_time @ start",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `count_over_time(http_requests_total[5m] @ start())`,
		},
		{
			name: "sum_over_time @ end",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `sum_over_time(http_requests_total[5m] @ start())`,
		},
		{
			name: "avg_over_time @ 180s",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `avg_over_time(http_requests_total[4m] @ 180.000)`,
		},
		{
			name: "@ vector 240s offset 2m",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total @ 240.000 offset 2m`,
		},
		{
			name: "avg_over_time @ 120s offset -2m",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total @ 120.000 offset -2m`,
		},
		{
			name: "sum_over_time @ 180s offset 2m",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `sum_over_time(http_requests_total[5m] @ 180.000 offset 2m)`,
		},
		{
			name: "binop with @ end() modifier inside query range",
			load: `load 30s
			    http_requests_total 2+3x100
			    http_responses_total 2+4x100`,
			query: `max(http_requests_total @ end()) / max(http_responses_total)`,
			end:   time.Unix(600, 0),
		},
		{
			name: "binop with @ end() modifier outside of query range",
			load: `load 30s
			    http_requests_total 2+3x100
			    http_responses_total 2+4x100`,
			query: `max(http_requests_total @ end()) / max(http_responses_total)`,
			end:   time.Unix(60000, 0),
		},
		{
			name: "days_in_month with input",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `days_in_month(http_requests_total)`,
		},
		{
			name: "days_in_month without input",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `days_in_month()`,
		},
		{
			name: "day_of_month with input",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `day_of_month(http_requests_total)`,
		},
		{
			name: "day_of_month without input",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `day_of_month()`,
		},
		{
			name: "day_of_week with input",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `days_in_month(http_requests_total)`,
		},
		{
			name: "day_of_week without input",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `days_in_month()`,
		},
		{
			name: "day_of_year with input",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `day_of_year(http_requests_total)`,
		},
		{
			name: "day_of_year without input",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `day_of_year()`,
		},
		{
			name: "hour with input",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `hour(http_requests_total)`,
		},
		{
			name: "hour without input",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `hour()`,
		},
		{
			name: "minute with input",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `minute(http_requests_total)`,
		},
		{
			name: "minute without input",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `minute()`,
		},
		{
			name: "month with input",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `month(http_requests_total)`,
		},
		{
			name: "month without input",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `month()`,
		},
		{
			name: "year with input",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `year(http_requests_total)`,
		},
		{
			name: "year without input",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `year()`,
		},
		{
			name: "selector merge",
			load: `load 30s
			    http_requests_total{pod="nginx-1", ns="nginx"} 1+1x15
			    http_requests_total{pod="nginx-2", ns="nginx"} 1+2x18
			    http_requests_total{pod="nginx-3", ns="nginx"} 1+2x21`,
			query: `
  http_requests_total{ns="nginx",pod=~"nginx-1"}
/ on () group_left ()
  sum(http_requests_total{ns="nginx"})`,
		},
		{
			name: "selector merge with different ranges",
			load: `load 30s
			    http_requests_total{pod="nginx-1", ns="nginx"} 2+2x16
			    http_requests_total{pod="nginx-2", ns="nginx"} 2+4x18
			    http_requests_total{pod="nginx-3", ns="nginx"} 2+6x20`,
			query: `
  rate(http_requests_total{ns="nginx",pod=~"nginx-1"}[2m])
+ on () group_left ()
  sum(http_requests_total{ns="nginx"})`,
		},
		{
			name: "binop with positive matcher using regex, only one side has data",
			load: `load 30s
			    metric{} 1+2x5
			    metric{} 1+2x20`,
			query: `sum(rate(metric{err=~".+"}[5m])) / sum(rate(metric[5m]))`,
		},
		{
			name: "binop with positive matcher using regex, both sides have data",
			load: `load 30s
			    metric{} 1+2x5
			    metric{err="FooBarKey"} 1+2x20`,
			query: `sum(rate(metric{err=~".+"}[5m])) / sum(rate(metric[5m]))`,
		},
		{
			name: "binop with negative matcher using regex, only one side has data",
			load: `load 30s
			    metric{} 1+2x5
			    metric{} 1+2x20`,
			query: `sum(rate(metric{err!~".+"}[5m])) / sum(rate(metric[5m]))`,
		},
		{
			name: "binop with negative matcher using regex, both sides have data",
			load: `load 30s
			    metric{} 1+2x5
			    metric{err="FooBarKey"} 1+2x20`,
			query: `sum(rate(metric{err!~".+"}[5m])) / sum(rate(metric[5m]))`,
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
			name: "quantile with param series",
			load: `load 30s
			    http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
			    http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
			    http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50
			    param_series 0+0.1x50`,
			query: `quantile(scalar(param_series), rate(http_requests_total[1m]))`,
		},
		{
			name: "quantile with param series that evaluates to NaN",
			load: `load 30s
			    http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
			    http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
			    http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50
			    param_series NaN+0x50`,
			query: `quantile(scalar(param_series), rate(http_requests_total[1m]))`,
		},
		{
			name: "quantile with non-existing param series",
			load: `load 30s
			    http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
			    http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
			    http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: `quantile(scalar(non_existent), rate(http_requests_total[1m]))`,
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
			query: `topk(2, http_requests_total)`,
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
			query: `topk(3.5, http_requests_total)`,
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
			query: `topk(0.5, http_requests_total)`,
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
			query: `topk(1e+120, http_requests_total)`,
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
			query: `topk(NaN, http_requests_total)`,
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name:  "topk with NaN and no matching series",
			query: `topk(NaN, not_there)`,
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "topk with NaN comparison",
			load: `load 30s
			    http_requests_total{pod="nginx-1", route="/"} NaN
			    http_requests_total{pod="nginx-2", route="/"}  NaN`,
			query: `topk by (route) (1, http_requests_total)`,
		},
		{
			name: "nested topk error that should not be skipped",
			load: `load 30s
			    X 1+1x50`,
			query: `topk(0, topk(NaN, X))`,
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
			query: `max(topk by (series) (2, http_requests_total))`,
			end:   time.Unix(3000, 0),
		},
		{
			name: "topk on empty result",
			load: `load 30s
			    metric_a 1+1x2`,
			query: `topk(2, histogram_quantile(0.1, metric_b))`,
		},
		{
			name: "topk by",
			load: `load 30s
			    http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
			    http_requests_total{pod="nginx-2", series="1"} 2+2.3x50
			    http_requests_total{pod="nginx-4", series="2"} 5+2.4x50
			    http_requests_total{pod="nginx-5", series="2"} 8.4+2.3x50
			    http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: `topk by (series) (2, http_requests_total)`,
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
			query: `topk by (series) (2 - 1, http_requests_total)`,
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
			query: `topk by (series) (scalar(min(http_requests_total)), http_requests_total)`,
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
			query: `topk by (series) (scalar(min(non_existent_metric)), http_requests_total)`,
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
			query: `bottomk(2, http_requests_total)`,
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
			query: `bottomk by (series) (2, http_requests_total)`,
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "sgn",
			load: `load 30s
			    http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
			    http_requests_total{pod="nginx-2", series="1"} -10+1x50`,
			query: `sgn(http_requests_total)`,
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
			query: `sort_desc(http_requests_total)`,
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
			query: `
avg by (storage_info) (
    storage_used
  * on (instance, storage_index) group_left (storage_info)
    (sum by (instance, storage_index, storage_info) (storage_info))
)`,
		},
		{
			name: "absent with partial data in range",
			load: `load 30s
			    existent{job="myjob"} 1 1 1`,
			query: `absent(existent{job="myjob"})`,
		},
		{
			name:  "absent with no data in range",
			load:  `load 30s`,
			query: `absent(nonexistent{job="myjob"})`,
		},
		{
			name:  "absent_over_time with no data in range",
			query: `absent_over_time(non_existent[10m])`,
		},
		{
			name: "absent_over_time with data in range",
			load: `load 30s
			    X{a="b"}  1x10`,
			query: `absent_over_time(X{a="b"}[10m])`,
		},
		{
			name: "absent_over_time - present but out of range",
			load: `load 30s
			    X{a="b"}  1x10`,
			query: `absent_over_time(X{a="b"}[1m])`,
			start: time.Unix(600, 0),
		},
		{
			name: "absent_over_time - absent because of label",
			load: `load 30s
			    X{a="b"}  1x10`,
			query: `absent_over_time(X{a!="b"}[1m])`,
		},
		{
			name: "subquery in binary expression",
			load: `load 60s
			    http_requests_total{pod="nginx-1", series="1"} 1+1x40`,
			query: `http_requests_total * (sum_over_time(http_requests_total[5m:1m]) > 0)`,
		},
		{
			name: "sum_over_time subquery with outer step larger than inner step",
			load: `load 60s
			    http_requests_total{pod="nginx-1", series="1"} 1+1x40`,
			query: `sum_over_time(sum_over_time(http_requests_total[2m])[5m:1m])`,
		},
		{
			name: "sum_over_time subquery with outer step equal to inner step",
			load: `load 60s
			    http_requests_total{pod="nginx-1", series="1"} 1+1x40`,
			query: `sum_over_time(sum_over_time(http_requests_total[2m])[5m:30s])`,
		},
		{
			name: "sum_over_time subquery with outer step smaller than inner step",
			load: `load 60s
			    http_requests_total{pod="nginx-1", series="1"} 1+1x40`,
			query: `sum_over_time(sum_over_time(http_requests_total[2m])[5m:15s])`,
		},
		{
			name: "sum_over_time subquery with aggregation",
			load: `load 10s
			    http_requests_total{pod="nginx-1", series="1"} 1+1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2x50`,
			query: `sum_over_time(sum by (pod) (http_requests_total)[5m:1m])`,
		},
		{
			name: "rate subquery with outer @ modifier",
			load: `load 30s
			    http_requests_total{pod="nginx-1", series="1"} 1+1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2x50`,
			query: `rate(http_requests_total[20s:10s] @ 100.000)`,
		},
		{
			name: "rate subquery with offset",
			load: `load 10s
			    http_requests_total{pod="nginx-1", series="1"} 1+2x40`,
			query: `rate(http_requests_total[20s:10s] offset 20s)`,
		},
	}

	disableOptimizerOpts := []bool{true, false}
	lookbackDeltas := []time.Duration{0, 30 * time.Second, time.Minute, 5 * time.Minute, 10 * time.Minute}
	for _, lookbackDelta := range lookbackDeltas {
		opts.LookbackDelta = lookbackDelta
		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				storage := promqltest.LoadedStorage(t, tc.load)
				defer storage.Close()

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
									// Set to 1 to make sure batching is tested.
									SelectorBatchSize: 1,
								})
								ctx := context.Background()
								q1, err := newEngine.NewRangeQuery(ctx, storage, nil, tc.query, tc.start, tc.end, tc.step)
								testutil.Ok(t, err)
								defer q1.Close()
								newResult := q1.Exec(ctx)

								oldEngine := promql.NewEngine(opts)
								q2, err := oldEngine.NewRangeQuery(ctx, storage, nil, tc.query, tc.start, tc.end, tc.step)
								testutil.Ok(t, err)
								defer q2.Close()
								oldResult := q2.Exec(ctx)

								testutil.WithGoCmp(comparer).Equals(t, oldResult, newResult, queryExplanation(q1))
							})
						}
					})
				}
			})
		}
	}
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

func TestWarnings(t *testing.T) {
	querier := &storage.MockQueryable{
		MockQuerier: &storage.MockQuerier{
			SelectMockFunction: func(sortSeries bool, hints *storage.SelectHints, matchers ...*labels.Matcher) storage.SeriesSet {
				return newWarningsSeriesSet(annotations.New().Add(errors.New("test warning")))
			},
		},
	}

	var (
		start = time.UnixMilli(0)
		end   = time.UnixMilli(600)
		step  = 30 * time.Second
	)

	cases := []struct {
		name          string
		query         string
		expectedWarns annotations.Annotations
	}{
		{
			name:  "single select call",
			query: `http_requests_total`,
			expectedWarns: annotations.New().Add(
				errors.New("test warning"),
			),
		},
		{
			name:  "multiple select calls",
			query: `sum(http_requests_total) / sum(http_responses_total)`,
			expectedWarns: annotations.New().Add(
				errors.New("test warning"),
			),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			newEngine := engine.New(engine.Opts{EngineOpts: promql.EngineOpts{Timeout: 1 * time.Hour}})
			q1, err := newEngine.NewRangeQuery(context.Background(), querier, nil, tc.query, start, end, step)
			testutil.Ok(t, err)

			res := q1.Exec(context.Background())
			testutil.Ok(t, res.Err)
			testutil.WithGoCmp(cmp.Comparer(func(err1, err2 error) bool {
				return err1.Error() == err2.Error()
			})).Equals(t, tc.expectedWarns, res.Warnings)
		})
	}
}

type scannersWithWarns struct {
	warn         error
	promScanners *prometheus.Scanners
}

func newScannersWithWarns(warn error) *scannersWithWarns {
	return &scannersWithWarns{
		warn: warn,
		promScanners: prometheus.NewPrometheusScanners(&storage.MockQueryable{
			MockQuerier: storage.NoopQuerier(),
		}),
	}
}

func (s scannersWithWarns) NewVectorSelector(ctx context.Context, opts *query.Options, hints storage.SelectHints, selector logicalplan.VectorSelector) (model.VectorOperator, error) {
	warnings.AddToContext(annotations.New().Add(s.warn), ctx)
	return s.promScanners.NewVectorSelector(ctx, opts, hints, selector)
}

func (s scannersWithWarns) NewMatrixSelector(ctx context.Context, opts *query.Options, hints storage.SelectHints, selector logicalplan.MatrixSelector, call logicalplan.FunctionCall) (model.VectorOperator, error) {
	warnings.AddToContext(annotations.New().Add(s.warn), ctx)
	return s.promScanners.NewMatrixSelector(ctx, opts, hints, selector, call)
}

func TestWarningsPlanCreation(t *testing.T) {
	var (
		opts         = engine.Opts{EngineOpts: promql.EngineOpts{Timeout: 1 * time.Hour}}
		expectedWarn = errors.New("test warning")
	)
	newEngine := engine.NewWithScanners(opts, newScannersWithWarns(expectedWarn))
	q1, err := newEngine.NewRangeQuery(context.Background(), nil, nil, "http_requests_total", time.UnixMilli(0), time.UnixMilli(600), 30*time.Second)
	testutil.Ok(t, err)

	res := q1.Exec(context.Background())
	testutil.Ok(t, res.Err)
	testutil.WithGoCmp(cmp.Comparer(func(err1, err2 error) bool {
		return err1.Error() == err2.Error()
	})).Equals(t, annotations.New().Add(expectedWarn), res.Warnings)

}

func TestEdgeCases(t *testing.T) {
	t.Parallel()
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
			query: `foo * on () group_left () bar`,
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			oldEngine := promql.NewEngine(opts)
			q1, err := oldEngine.NewRangeQuery(ctx, storageWithSeries(tc.series...), nil, tc.query, tc.start, tc.end, step)
			testutil.Ok(t, err)

			newEngine := engine.New(engine.Opts{EngineOpts: opts})
			q2, err := newEngine.NewRangeQuery(ctx, storageWithSeries(tc.series...), nil, tc.query, tc.start, tc.end, step)
			testutil.Ok(t, err)

			oldResult := q1.Exec(ctx)
			newResult := q2.Exec(ctx)

			testutil.WithGoCmp(comparer).Equals(t, oldResult, newResult, queryExplanation(q1))
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
		storage := promqltest.LoadedStorage(t, tc.load)
		defer storage.Close()

		optimizers := logicalplan.AllOptimizers

		newEngine := engine.New(engine.Opts{
			EngineOpts:        opts,
			DisableFallback:   true,
			LogicalOptimizers: optimizers,
		})
		_, err := newEngine.NewInstantQuery(context.Background(), storage, nil, tc.query, queryTime)
		testutil.NotOk(t, err)
	}
}

func TestXFunctionsWithNativeHistograms(t *testing.T) {
	defaultQueryTime := time.Unix(50, 0)

	expr := "sum(xincrease(native_histogram_series[50s]))"

	// Negative offset and at modifier are enabled by default
	// since Prometheus v2.33.0, so we also enable them.
	opts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}

	lStorage := teststorage.New(t)
	defer lStorage.Close()

	app := lStorage.Appender(context.TODO())
	testutil.Ok(t, generateFloatHistogramSeries(app, 3000, false))
	testutil.Ok(t, app.Commit())

	optimizers := logicalplan.AllOptimizers

	ctx := context.Background()
	newEngine := engine.New(engine.Opts{
		EngineOpts:        opts,
		DisableFallback:   true,
		LogicalOptimizers: optimizers,
		EnableXFunctions:  true,
	})
	query, err := newEngine.NewInstantQuery(ctx, lStorage, nil, expr, defaultQueryTime)
	testutil.Ok(t, err)
	defer query.Close()

	engineResult := query.Exec(ctx)
	require.Error(t, engineResult.Err)
}

func TestXFunctionsWhenDisabled(t *testing.T) {
	var (
		query = "xincrease(http_requests[50s])"
		start = time.Unix(0, 0)
		end   = time.Unix(100, 0)
		step  = time.Second * 10
	)
	ng := engine.New(engine.Opts{})
	_, err := ng.NewRangeQuery(context.Background(), nil, nil, query, start, end, step)
	testutil.NotOk(t, err)
	testutil.Equals(t, `1:1: parse error: unknown function with name "xincrease"`, err.Error())

	_, err = ng.NewInstantQuery(context.Background(), nil, nil, query, start)
	testutil.NotOk(t, err)
	testutil.Equals(t, `1:1: parse error: unknown function with name "xincrease"`, err.Error())
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
		name       string
		load       string
		query      string
		queryTime  time.Time
		expected   []promql.Sample
		rangeQuery bool
		startTime  time.Time
		endTime    time.Time
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
			query:     `increase(http_requests[30m])`,
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
			storage := promqltest.LoadedStorage(t, tc.load)
			defer storage.Close()

			queryTime := defaultQueryTime
			if tc.queryTime != (time.Time{}) {
				queryTime = tc.queryTime
			}

			optimizers := logicalplan.AllOptimizers

			ctx := context.Background()
			newEngine := engine.New(engine.Opts{
				EngineOpts:        opts,
				DisableFallback:   true,
				LogicalOptimizers: optimizers,
				EnableXFunctions:  true,
			})
			query, err := newEngine.NewInstantQuery(ctx, storage, nil, tc.query, queryTime)
			testutil.Ok(t, err)
			defer query.Close()

			engineResult := query.Exec(ctx)
			testutil.Ok(t, engineResult.Err)
			expectedResult := createVectorResult(tc.expected)

			testutil.WithGoCmp(comparer).Equals(t, expectedResult, engineResult, queryExplanation(query))
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
		name       string
		load       string
		query      string
		queryTime  time.Time
		expected   promql.Vector
		rangeQuery bool
		startTime  time.Time
		endTime    time.Time
	}{
		// ### Timeseries starts insice range, (presumably) goes on after range end. ###
		// 1. Reference eval
		{
			name:      "eval instant at 25s rate, with 50s lookback",
			query:     `rate(http_requests[50s])`,
			queryTime: time.Unix(25, 0),
			expected: []promql.Sample{
				createSample(defaultQueryTime.UnixMilli(), 0.022, labels.FromStrings("path", "/foo")),
				createSample(defaultQueryTime.UnixMilli(), 0.11000000000000001, labels.FromStrings("path", "/bar")),
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
			query:     `rate(http_requests[50s])`,
			queryTime: time.Unix(24, 0),
			expected: []promql.Sample{
				createSample(24000, 0.0265, labels.FromStrings("path", "/foo")),
				createSample(24000, 0.106, labels.FromStrings("path", "/bar")),
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
			query:     `rate(http_requests[50s])`,
			queryTime: time.Unix(26, 0),
			expected: []promql.Sample{
				createSample(26000, 0.022799999999999997, labels.FromStrings("path", "/foo")),
				createSample(26000, 0.11399999999999999, labels.FromStrings("path", "/bar")),
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
			query:     `rate(http_requests[50s])`,
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
			query:     `rate(http_requests[50s])`,
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
			query:     `rate(http_requests[50s])`,
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
			query:     `rate(http_requests[10s])`,
			queryTime: time.Unix(9, 0),
			expected: []promql.Sample{
				createSample(9000, 0, labels.FromStrings("path", "/foo")),
				createSample(9000, 0.2, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 19s rate(http_requests[10s]), with 10s lookback",
			query:     `rate(http_requests[10s])`,
			queryTime: time.Unix(19, 0),
			expected: []promql.Sample{
				createSample(19000, 0.2, labels.FromStrings("path", "/foo")),
				createSample(19000, 0.2, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 29s rate(http_requests[10s]), with 10s lookback",
			query:     `rate(http_requests[10s])`,
			queryTime: time.Unix(29, 0),
			expected: []promql.Sample{
				createSample(29000, 0, labels.FromStrings("path", "/foo")),
				createSample(29000, 0.2, labels.FromStrings("path", "/bar")),
			},
		},
		{
			name:      "eval instant at 39s rate(http_requests[10s]), with 10s lookback",
			query:     `rate(http_requests[10s])`,
			queryTime: time.Unix(39, 0),
			expected: []promql.Sample{
				createSample(39000, 0, labels.FromStrings("path", "/foo")),
				createSample(39000, 0.2, labels.FromStrings("path", "/bar")),
			},
		},
		// XXX Missed an increase in path="/foo" between timestamps 35 and 40 (both in this eval and the one before).
		{
			name:      "eval instant at 49s rate(http_requests[10s]), with 10s lookback",
			query:     `rate(http_requests[10s])`,
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			load := defaultLoad
			if tc.load != "" {
				load = tc.load
			}

			storage := promqltest.LoadedStorage(t, load)
			defer storage.Close()

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
			query, err := newEngine.NewInstantQuery(context.Background(), storage, nil, tc.query, queryTime)
			testutil.Ok(t, err)
			defer query.Close()

			engineResult := query.Exec(context.Background())
			expectedResult := createVectorResult(tc.expected)

			testutil.WithGoCmp(comparer).Equals(t, expectedResult, engineResult, queryExplanation(query))
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
	t.Parallel()
	defaultQueryTime := time.Unix(50, 0)
	cases := []struct {
		load      string
		name      string
		query     string
		queryTime time.Time
	}{
		{
			name: "count_values fuzz",
			load: `load 30s
			    http_requests_total{pod="nginx-1", route="/"} 51.00+1.00x40
			    http_requests_total{pod="nginx-2", route="/"} -74+14.00x40`,
			query: `
count_values without () (
  "value",
    (atanh(http_requests_total{pod="nginx-1"}) > tanh(http_requests_total{route="/"}))
  or
    avg by (pod, __name__) (http_requests_total{route="/"})
)`,
		},
		{
			name: "sum evaluates to -0 fuzz",
			load: `load 30s
			    http_requests_total{pod="nginx-2", route="/"}  0`,
			query:     `sum by (pod) (-http_requests_total) atan2 -0`,
			queryTime: time.Unix(0, 0),
		},
		{
			name: "count_values",
			load: `load 30s
			    version{foo="bar"} 1
			    version{foo="baz"} 1
			    version{foo="quz"} 2`,
			query:     `count_values("val", version)`,
			queryTime: time.Unix(0, 0),
		},
		{
			name: "binary pairing early exit fuzz",
			load: `load 30s
			    http_requests_total{pod="nginx-1", route="/"} 33.00+1.00x40
			    http_requests_total{pod="nginx-2", route="/"}  1+2.00x40`,
			query: `
  avg without (route) (avg(http_requests_total) / http_requests_total)
<=
  sum by (__name__) (http_requests_total or avg(http_requests_total))`,
			queryTime: time.Unix(0, 0),
		},
		{
			name: "offset and @ modifiers",
			load: `load 30s
			    http_requests_total{pod="nginx-0", route="/"} 1+1x30`,
			query:     `http_requests_total @ end() offset 2m`,
			queryTime: time.Unix(300, 0),
		},
		{
			name: "timestamp - offset modifier",
			load: `load 30s
			    http_requests_total{pod="nginx-0", route="/"} 0x30`,
			query:     `timestamp(http_requests_total offset 2m)`,
			queryTime: time.Unix(300, 0),
		},
		{
			name: "timestamp - @ modifier",
			load: `load 30s
			    http_requests_total{pod="nginx-0", route="/"} 0x30`,
			query:     `timestamp(http_requests_total @ 60.000)`,
			queryTime: time.Unix(300, 0),
		},
		{
			name: "timestamp - nested functions with offset",
			load: `load 30s
			    http_requests_total{pod="nginx-0", route="/"} 0x30`,
			query:     `timestamp(timestamp(http_requests_total offset 2m))`,
			queryTime: time.Unix(300, 0),
		},
		{
			name:      "timestamp - nested functions without any scan",
			query:     `timestamp(vector(1))`,
			queryTime: time.Unix(300, 0),
		},
		{
			name: "timestamp - aggregation",
			load: `load 30s
			    http_requests_total{pod="nginx-0", route="/"} 0x30`,
			query:     `timestamp(sum(http_requests_total))`,
			queryTime: time.Unix(300, 0),
		},
		{
			name: "timestamp - fuzzing failure",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1.00+1.00x15
			    http_requests_total{pod="nginx-2"}  1+2.00x21`,
			query:     `timestamp(http_requests_total @ end() offset -2m23s)`,
			queryTime: time.Unix(300, 0),
		},
		{
			name: "fuzz - min with NaN",
			load: `load 30s
			    http_requests_total{pod="nginx-1", route="/"} 0
			    http_requests_total{pod="nginx-2", route="/"}  NaN`,
			query:     `min without (__name__, pod) (http_requests_total)`,
			queryTime: time.Unix(0, 0),
		},
		{
			name: "fuzz - max with NaN",
			load: `load 30s
			    http_requests_total{pod="nginx-1", route="/"} 0
			    http_requests_total{pod="nginx-2", route="/"}  NaN`,
			query:     `max without (__name__, pod) (http_requests_total)`,
			queryTime: time.Unix(0, 0),
		},
		{
			name: "fuzz - min with NaN",
			load: `load 30s
			    http_requests_total{pod="nginx-1", route="/"} 124.00+1.00x40
			    http_requests_total{pod="nginx-2", route="/"}  0+0.29x40`,
			query: `min by (route, pod) (sqrt(-http_requests_total))`,
		},
		{
			name: "fuzz - min with Inf",
			load: `load 30s
			    http_requests_total{pod="nginx-1", route="/"} 483.00+6035.00x40
			    http_requests_total{pod="nginx-2", route="/"}  2+47.14x40`,
			query: `
min without () (
  (
      {__name__="http_requests_total"} @ start() offset -2m40s
    ^
      {__name__="http_requests_total"} @ start() offset -2m49s
  )
)`,
		},
		/*
		   This is a known issue, we lose the signed 0 in the sum because we add to a
		   default element. Prometheus assigns the first element to the sum and preserves
		   the sign.
		       {
		         name:  "fuzz - signed zero",
		         query: `1/sum(-(absent(X)-1))`,
		       },
		*/
		{
			name: "sum_over_time with subquery",
			load: `load 10s
			    http_requests_total{pod="nginx-1", series="1"} 1+1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2x50
			    http_requests_total{pod="nginx-5", series="1"} 8+4x50
			    http_requests_total{pod="nginx-6", series="2"} 2+3x50`,
			queryTime: time.Unix(600, 0),
			query:     `sum_over_time(sum by (series) (http_requests_total)[5m:1m])`,
		},
		{
			name: "sum_over_time with subquery with default step",
			load: `load 10s
			    http_requests_total{pod="nginx-1", series="1"} 1+1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2x50
			    http_requests_total{pod="nginx-5", series="1"} 8+4x50
			    http_requests_total{pod="nginx-6", series="2"} 2+3x50`,
			queryTime: time.Unix(600, 0),
			query:     `sum_over_time(sum by (series) (http_requests_total)[5m:])`,
		},
		{
			name: "sum_over_time with subquery with resolution that doesnt divide step length",
			load: `load 10s
			    http_requests_total{pod="nginx-1", series="1"} 1+1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2x50
			    http_requests_total{pod="nginx-5", series="1"} 8+4x50
			    http_requests_total{pod="nginx-6", series="2"} 2+3x50`,
			queryTime: time.Unix(600, 0),
			query:     `sum_over_time(sum by (series) (http_requests_total)[5m:22s])`,
		},
		{
			name: "sum_over_time with subquery with offset",
			load: `load 10s
			    http_requests_total{pod="nginx-1", series="1"} 1+1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2x50
			    http_requests_total{pod="nginx-5", series="1"} 8+4x50
			    http_requests_total{pod="nginx-6", series="2"} 2+3x50`,
			queryTime: time.Unix(600, 0),
			query:     `sum_over_time(sum by (series) (http_requests_total)[5m:1m] offset 1m)`,
		},
		{
			name: "sum_over_time with subquery with inner offset",
			load: `load 10s
			    http_requests_total{pod="nginx-1", series="1"} 1+1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2x50
			    http_requests_total{pod="nginx-5", series="1"} 8+4x50
			    http_requests_total{pod="nginx-6", series="2"} 2+3x50`,
			queryTime: time.Unix(600, 0),
			query:     `sum_over_time(sum by (series) (http_requests_total offset 1m)[5m:1m])`,
		},
		{
			name: "sum_over_time with subquery with inner @ modifier",
			load: `load 10s
			    http_requests_total{pod="nginx-1", series="1"} 1+1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2x50
			    http_requests_total{pod="nginx-5", series="1"} 8+4x50
			    http_requests_total{pod="nginx-6", series="2"} 2+3x50`,
			queryTime: time.Unix(600, 0),
			query:     `sum_over_time(sum by (series) (http_requests_total @ 10.000)[5m:1m])`,
		},
		{
			name: "sum_over_time with nested subqueries with inner @ modifier",
			load: `load 10s
			    http_requests_total{pod="nginx-1", series="1"} 1+1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2x50
			    http_requests_total{pod="nginx-5", series="1"} 8+4x50
			    http_requests_total{pod="nginx-6", series="2"} 2+3x50`,
			queryTime: time.Unix(600, 0),
			query:     `sum_over_time(rate(sum by (series) (http_requests_total @ 10.000)[5m:1m] @ 0.000)[10m:1m])`,
		},
		{
			name: "sum_over_time with subquery should drop name label",
			load: `load 10s
			    http_requests_total{pod="nginx-1", series="1"} 1+1x40
			    http_requests_total{pod="nginx-2", series="1"} 2+2x50`,
			queryTime: time.Unix(0, 0),
			query:     `sum_over_time(http_requests_total{series="1"} offset 7s[1h:1m] @ 119.800)`,
		},
		{
			name: "duplicate label set",
			load: `load 5m
			    testmetric1{src="a",dst="b"} 0
			    testmetric2{src="a",dst="b"} 1`,
			query: `changes({__name__=~"testmetric1|testmetric2"}[5m])`,
		},
		{
			name:      "scalar",
			load:      ``,
			queryTime: time.Unix(160, 0),
			query:     `12 + 1`,
		},
		{
			name:      "string literal",
			load:      ``,
			queryTime: time.Unix(160, 0),
			query:     `test - string - literal`,
		},
		{
			name: "increase plus offset",
			load: `load 1s
			    http_requests_total{pod="nginx-1"} 1+1x180`,
			queryTime: time.Unix(160, 0),
			query:     `increase(http_requests_total[1m] offset 1m)`,
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
			query:     `round(http_requests_total)`,
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
			query:     `round(http_requests_total, 0.5)`,
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
			query:     `sort(http_requests_total)`,
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
			query:     `sort_desc(http_requests_total)`,
		},
		{
			name: "histogram_quantile with mock duplicate labels",
			load: `load 30s
			    http_requests_total{pod="nginx-2", route="/"}  0+0.14x40`,
			queryTime: time.Unix(600, 0),
			query:     `histogram_quantile(10, -http_requests_total or http_requests_total)`,
		},
		{
			name: "quantile by pod",
			load: `load 30s
			    http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
			    http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
			    http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: `quantile by (pod) (0.9, rate(http_requests_total[1m]))`,
		},
		{
			name: "quantile by pod with binary",
			load: `load 30s
			    http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
			    http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
			    http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: `quantile by (pod) (1 - 0.1, rate(http_requests_total[1m]))`,
		},
		{
			name: "quantile by pod with expression",
			load: `load 30s
			    http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
			    http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
			    http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: `quantile by (pod) (scalar(min(http_requests_total)), rate(http_requests_total[1m]))`,
		},
		{
			name: "quantile",
			load: `load 30s
			    http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
			    http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
			    http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: `quantile(0.9, rate(http_requests_total[1m]))`,
		},
		{
			name: "stdvar",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x4
			    http_requests_total{pod="nginx-2"} 1+2x4`,
			query: `stdvar(http_requests_total)`,
		},
		{
			name: "stdvar by pod",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1
			    http_requests_total{pod="nginx-2"} 2
			    http_requests_total{pod="nginx-3"} 8
			    http_requests_total{pod="nginx-4"} 6`,
			query: `stdvar by (pod) (http_requests_total)`,
		},
		{
			name: "stddev",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x4
			    http_requests_total{pod="nginx-2"} 1+2x4`,
			query: `stddev(http_requests_total)`,
		},
		{
			name: "stddev by pod",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1
			    http_requests_total{pod="nginx-2"} 2
			    http_requests_total{pod="nginx-3"} 8
			    http_requests_total{pod="nginx-4"} 6`,
			query: `stddev by (pod) (http_requests_total)`,
		},
		{
			name: "sum by pod",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x4
			    http_requests_total{pod="nginx-2"} 1+2x4`,
			query: `sum by (pod) (http_requests_total)`,
		},
		{
			name: "count",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `count(http_requests_total)`,
		},
		{
			name: "average",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `avg(http_requests_total)`,
		},
		{
			name: "label_join",
			load: `load 30s
			    http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2.4x50`,
			queryTime: time.Unix(160, 0),
			query:     `label_join(http_requests_total, "label", "-", "pod", "series")`,
		},
		{
			name: "label_join with non-existing src labels",
			load: `load 30s
			    http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2.4x50`,
			queryTime: time.Unix(160, 0),
			query:     `label_join(http_requests_total, "label", "-", "test", "fake")`,
		},
		{
			name: "label_join with overwrite dst label if exists",
			load: `load 30s
			    http_requests_total{pod="nginx-1", series="1", label="test-1"} 1+1.1x40
			    http_requests_total{pod="nginx-2", series="2", label="test-2"} 2+2.3x50
			    http_requests_total{pod="nginx-4", series="3", label="test-3"} 5+2.4x50`,
			queryTime: time.Unix(160, 0),
			query:     `label_join(http_requests_total, "label", "-", "pod", "series")`,
		},
		{
			name: "label_join with no src labels provided",
			load: `load 30s
			    http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2.4x50`,
			queryTime: time.Unix(160, 0),
			query:     `label_join(http_requests_total, "label", "-")`,
		},
		{
			name: "label_replace",
			load: `load 30s
			    http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2.4x50`,
			queryTime: time.Unix(160, 0),
			query:     `label_replace(http_requests_total, "foo", "$1", "series", ".*")`,
		},
		{
			name: "label_replace with bad regular expression",
			load: `load 30s
			    http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2.4x50`,
			queryTime: time.Unix(160, 0),
			query:     `label_replace(http_requests_total, "foo", "$1", "series", "]]")`,
		},
		{
			name: "label_replace non-existing src label",
			load: `load 30s
			    http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2.4x50`,
			queryTime: time.Unix(160, 0),
			query:     `label_replace(http_requests_total, "foo", "$1", "bar", ".*")`,
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
			query: `topk(2, http_requests_total)`,
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
			query: `topk by (series) (2, http_requests_total)`,
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
			query: `bottomk(2, http_requests_total)`,
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
			query: `bottomk by (series) (2, http_requests_total)`,
		},
		{
			name: "max",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `max(http_requests_total)`,
		},
		{
			name: "max with only 1 sample",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} -1
			    http_requests_total{pod="nginx-2"} 1`,
			query: `max by (pod) (http_requests_total)`,
		},
		{
			name: "min",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `min(http_requests_total)`,
		},
		{
			name: "min with only 1 sample",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} -1
			    http_requests_total{pod="nginx-2"} 1`,
			query: `min by (pod) (http_requests_total)`,
		},
		{
			name: "rate",
			load: `load 30s
			    http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
			    http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
			    http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
			    http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
			    http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: `rate(http_requests_total[1m])`,
		},
		{
			name: "sum rate",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x4
			    http_requests_total{pod="nginx-2"} 1+2x4`,
			query: `sum(rate(http_requests_total[1m]))`,
		},
		{
			name: "sum rate with single sample series",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x4
			    http_requests_total{pod="nginx-2"} 1+2x4
			    http_requests_total{pod="nginx-3"} 0`,
			query: `sum by (pod) (rate(http_requests_total[1m]))`,
		},
		{
			name: "sum rate with stale series",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x4
			    http_requests_total{pod="nginx-2"} 1+2x20`,
			query: `sum(rate(http_requests_total[1m]))`,
		},
		{
			name: "delta",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x4
			    http_requests_total{pod="nginx-2"} 1+2x4`,
			query: `delta(http_requests_total[1m])`,
		},
		{
			name: "increase",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x4
			    http_requests_total{pod="nginx-2"} 1+2x4`,
			query: `increase(http_requests_total[1m])`,
		},

		{
			name: "sum irate",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x4
			    http_requests_total{pod="nginx-2"} 1+2x4`,
			query: `sum(irate(http_requests_total[1m]))`,
		},
		{
			name: "sum irate with stale series",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x4
			    http_requests_total{pod="nginx-2"} 1+2x20`,
			query: `sum(irate(http_requests_total[1m]))`,
		},
		{
			name:  "number literal",
			load:  "",
			query: `34`,
		},
		{
			name:  "vector",
			load:  "",
			query: `vector(24)`,
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
			query: `sum by (method) (foo) == sum by (method) (bar)`,
		},
		{
			name: "vector binary op !=",
			load: `load 30s
			    foo{method="get", code="500"} 1+1x40
			    bar{method="get", code="404"} 1+1.1x30`,
			query: `sum by (method) (foo) != sum by (method) (bar)`,
		},
		{
			name: "vector binary op >",
			load: `load 30s
			    foo{method="get", code="500"} 1+1x40
			    bar{method="get", code="404"} 1+1.1x30`,
			query: `sum by (method) (foo) > sum by (method) (bar)`,
		},
		{
			name: "vector binary op <",
			load: `load 30s
			    foo{method="get", code="500"} 1+1x40
			    bar{method="get", code="404"} 1+1.1x30`,
			query: `sum by (method) (foo) < sum by (method) (bar)`,
		},
		{
			name: "vector binary op >=",
			load: `load 30s
			    foo{method="get", code="500"} 1+1x40
			    bar{method="get", code="404"} 1+1.1x30`,
			query: `sum by (method) (foo) >= sum by (method) (bar)`,
		},
		{
			name: "vector binary op <=",
			load: `load 30s
			    foo{method="get", code="500"} 1+1x40
			    bar{method="get", code="404"} 1+1.1x30`,
			query: `sum by (method) (foo) <= sum by (method) (bar)`,
		},
		{
			name: "vector binary op ^",
			load: `load 30s
			    foo{method="get", code="500"} 1+1x40
			    bar{method="get", code="404"} 1+1.1x30`,
			query: `sum by (method) (foo) ^ sum by (method) (bar)`,
		},
		{
			name: "vector binary op %",
			load: `load 30s
			    foo{method="get", code="500"} 1+2x40
			    bar{method="get", code="404"} 1+1x30`,
			query: `sum by (method) (foo) % sum by (method) (bar)`,
		},
		{
			name: "vector binary op and 1",
			load: `load 30s
			    foo{method="get", code="500"} 1+2x40
			    bar{method="get", code="404"} 1+1x30`,
			query: `sum by (method) (foo) and sum by (method) (bar)`,
		},
		{
			name: "vector binary op and 2",
			load: `load 30s
			    foo{method="get", code="500"} 1+2x40
			    bar{method="get", code="404"} 1+1x30`,
			query: `sum by (code) (foo) and sum by (code) (bar)`,
		},
		{
			name: "vector binary op unless 1",
			load: `load 30s
			    foo{method="get", code="500"} 1+2x40
			    bar{method="get", code="404"} 1+1x30`,
			query: `sum by (method) (foo) unless sum by (method) (bar)`,
		},
		{
			name: "vector binary op unless 2",
			load: `load 30s
			    foo{method="get", code="500"} 1+2x40
			    bar{method="get", code="404"} 1+1x30`,
			query: `sum by (code) (foo) unless sum by (code) (bar)`,
		},
		{
			name: "vector binary op unless 3",
			load: `load 30s
			    foo{method="get", code="500"} 1+2x40`,
			query: `sum by (code) (foo) unless nonexistent`,
		},
		{
			name: "vector binary op or 1",
			load: `load 30s
			    foo{A="1"} 1+1x40
			    foo{A="2"} 2+2x40`,
			query: `sinh(foo or exp(foo))`,
		},
		{
			name: "vector binary op one-to-one left multiple matches",
			load: `load 30s
			    foo{method="get", code="500"} 1
			    foo{method="get", code="200"} 1
			    bar{method="get", code="200"} 1`,
			query: `foo / ignoring (code) bar`,
		},
		{
			name: "vector binary operation with many-to-many matching rhs high card",
			load: `load 30s
			    foo{code="200", method="get"} 1+1x20
			    foo{code="200", method="post"} 1+1x20
			    bar{code="200", method="get"} 1+1x20
			    bar{code="200", method="post"} 1+1x20`,
			query: `foo + on (code) group_right () bar`,
		},
		{
			name: "vector binary op > scalar",
			load: `load 30s
			    foo{method="get", code="500"} 1+2x40
			    bar{method="get", code="404"} 1+1x30`,
			query: `sum by (method) (foo) > 10`,
		},
		{
			name: "scalar < vector binary op",
			load: `load 30s
			    foo{method="get", code="500"} 1+2x40
			    bar{method="get", code="404"} 1+1x30`,
			query: `10 < sum by (method) (foo)`,
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
			query: `http_requests_total`,
		},
		{
			name:  "empty series with func",
			load:  "",
			query: `sum(http_requests_total)`,
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
			query: `http_requests_total @ 10.000`,
		},
		{
			name: "@ vector time 120s",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total @ 120.000`,
		},
		{
			name: "@ vector time 360s",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total @ 360.000`,
		},
		{
			name: "@ vector start",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total @ start()`,
		},
		{
			name: "@ vector end",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total @ end()`,
		},
		{
			name: "count_over_time @ start",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `count_over_time(http_requests_total[5m] @ start())`,
		},
		{
			name: "sum_over_time @ end",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `sum_over_time(http_requests_total[5m] @ start())`,
		},
		{
			name: "avg_over_time @ 180s",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `avg_over_time(http_requests_total[4m] @ 180.000)`,
		},
		{
			name: "@ vector 240s offset 2m",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total @ 240.000 offset 2m`,
		},
		{
			name: "avg_over_time @ 120s offset -2m",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total @ 120.000 offset -2m`,
		},
		{
			name: "sum_over_time @ 180s offset 2m",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x15
			    http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `sum_over_time(http_requests_total[5m] @ 180.000 offset 2m)`,
		},
		{
			name: "scalar with nested binary operator with step invariant",
			load: `load 30s
			    http_requests_total{pod="nginx-1", route="/"} 53.33+56.00x40
			    http_requests_total{pod="nginx-2", route="/"} -26+2.00x40`,
			query: `vector(scalar((http_requests_total @ end() offset 5m > http_requests_total)))`,
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
			query: `sgn(http_requests_total)`,
		},
		{
			name:  "absent and series does not exist",
			load:  `load 30s`,
			query: `absent(nonexistent{job="myjob"})`,
		},
		{
			name: "absent and series exists",
			load: `load 30s
			    existent{job="myjob"} 1`,
			query: `absent(existent{job="myjob"})`,
		},
		{
			name:  "absent and regex matcher",
			load:  `load 30s`,
			query: `absent(nonexistent{instance=~".*",job="myjob"})`,
		},
		{
			name:  "absent and duplicate matchers",
			load:  `load 30s`,
			query: `absent(nonexistent{foo="bar",job="myjob",job="yourjob"})`,
		},
		{
			name:  "absent and nested function",
			load:  `load 30s`,
			query: `absent(sum(nonexistent{job="myjob"}))`,
		},
		{
			name: "absent and nested absent with existing series",
			load: `load 30s
			    existent{job="myjob"} 1`,
			query: `absent(absent(existent{job="myjob"}))`,
		},
		{
			name: "absent_over_time with subquery - present data",
			load: `load 30s
			    X{a="b"}  1x10`,
			query: `absent_over_time(sum_over_time(X{a="b"}[1m])[1m:30s])`,
		},
		{
			name: "absent_over_time with subquery - missing data",
			load: `load 30s
			    X{a="b"}  1x10`,
			query: `absent_over_time(sum_over_time(X{a!="b"}[1m])[1m:30s])`,
		},
		{
			name: "absent_over_time with subquery and fixed offset",
			load: `load 30s
			    http_requests_total{pod="nginx-1"} 1+1x10
			    http_requests_total{pod="nginx-2"} 1+2x10`,
			query: `absent_over_time(http_requests_total @ start()[1h:1m])`,
		},
		{
			name: "absent_over_time fuzzer findings",
			load: `load 30s
			    http_requests_total{pod="nginx-1", route="/"} 0.02+1.00x40
			    http_requests_total{pod="nginx-2", route="/"} -24+0.67x40`,
			query: `
  count without (route, pod) ({__name__="http_requests_total"} @ 153.689)
>
  absent_over_time({__name__="http_requests_total",route="/"}[3m] offset 1m45s)`,
		},
	}

	disableOptimizerOpts := []bool{true, false}
	lookbackDeltas := []time.Duration{0, 30 * time.Second, time.Minute, 5 * time.Minute, 10 * time.Minute}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testStorage := promqltest.LoadedStorage(t, tc.load)
			defer testStorage.Close()
			for _, disableOptimizers := range disableOptimizerOpts {
				disableOptimizers := disableOptimizers
				t.Run(fmt.Sprintf("disableOptimizers=%t", disableOptimizers), func(t *testing.T) {
					for _, lookbackDelta := range lookbackDeltas {
						lookbackDelta := lookbackDelta
						// Negative offset and at modifier are enabled by default
						// since Prometheus v2.33.0, so we also enable them.
						opts := promql.EngineOpts{
							Timeout:                  1 * time.Hour,
							MaxSamples:               1e10,
							EnableNegativeOffset:     true,
							EnableAtModifier:         true,
							NoStepSubqueryIntervalFn: func(rangeMillis int64) int64 { return 30 * time.Second.Milliseconds() },
							LookbackDelta:            lookbackDelta,
						}

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

								ctx := context.Background()
								q1, err := newEngine.NewInstantQuery(ctx, testStorage, nil, tc.query, queryTime)
								testutil.Ok(t, err)
								defer q1.Close()

								newResult := q1.Exec(ctx)

								oldEngine := promql.NewEngine(opts)
								q2, err := oldEngine.NewInstantQuery(ctx, testStorage, nil, tc.query, queryTime)
								testutil.Ok(t, err)
								defer q2.Close()

								oldResult := q2.Exec(ctx)

								testutil.WithGoCmp(comparer).Equals(t, oldResult, newResult, queryExplanation(q1))
							})
						}
					}
				})
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

	querier := &storage.MockQueryable{
		MockQuerier: &storage.MockQuerier{
			SelectMockFunction: func(sortSeries bool, hints *storage.SelectHints, matchers ...*labels.Matcher) storage.SeriesSet {
				return newTestSeriesSet(&slowSeries{})
			},
		},
	}

	ctx := context.Background()
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

func TestQueryConcurrency(t *testing.T) {
	const storageDelay = 200 * time.Millisecond
	queryable := &storage.MockQueryable{
		MockQuerier: &storage.MockQuerier{
			SelectMockFunction: func(sortSeries bool, hints *storage.SelectHints, matchers ...*labels.Matcher) storage.SeriesSet {
				return newSlowSeriesSet(storageDelay)
			},
		},
	}

	var (
		ctx          = context.Background()
		logger       = log.NewLogfmtLogger(os.Stdout)
		concurrency  = 2
		maxQueries   = 4
		responseChan = make(chan struct{}, maxQueries)
	)
	newEngine := engine.New(engine.Opts{
		EngineOpts: promql.EngineOpts{
			Timeout:            1 * time.Hour,
			MaxSamples:         math.MaxInt64,
			ActiveQueryTracker: promql.NewActiveQueryTracker(t.TempDir(), concurrency, logger),
		}},
	)
	for i := 0; i < maxQueries; i++ {
		go func() {
			qry, err := newEngine.NewRangeQuery(ctx, queryable, nil, `count(metric)`, time.Unix(0, 0), time.Unix(300, 0), time.Second*30)
			testutil.Ok(t, err)

			resp := qry.Exec(ctx)
			testutil.Ok(t, resp.Err)

			responseChan <- struct{}{}
		}()
	}

	var (
		i           = 0
		gracePeriod = storageDelay + 10*time.Millisecond
	)
	for i < concurrency {
		select {
		case <-time.After(gracePeriod):
			t.Errorf("expected query to complete within %f seconds", gracePeriod.Seconds())
		case <-responseChan:
		}
		i++
	}
	select {
	case <-responseChan:
		t.Error("Expected to block on a query but did not")
	case <-time.After(10 * time.Millisecond):
		break
	}
	for i < maxQueries {
		select {
		case <-time.After(gracePeriod):
			t.Errorf("expected query to complete within %f seconds", gracePeriod.Seconds())
		case <-responseChan:
		}
		i++
	}
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

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	newEngine := engine.New(engine.Opts{DisableFallback: true, EngineOpts: opts})

	q, err := newEngine.NewInstantQuery(context.Background(), storage, nil, query, end)
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

func (h *hintRecordingQuerier) Select(_ context.Context, sortSeries bool, hints *storage.SelectHints, matchers ...*labels.Matcher) storage.SeriesSet {
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
			query: `foo`, start: 10000,
			expected: []*storage.SelectHints{
				{Start: 5000, End: 10000},
			},
		}, {
			query: `foo @ 15.000`, start: 10000,
			expected: []*storage.SelectHints{
				{Start: 10000, End: 15000},
			},
		}, {
			query: `foo @ 1.000`, start: 10000,
			expected: []*storage.SelectHints{
				{Start: -4000, End: 1000},
			},
		}, {
			query: `foo[2m]`, start: 200000,
			expected: []*storage.SelectHints{
				{Start: 80000, End: 200000, Range: 120000},
			},
		}, {
			query: `foo[2m] @ 180.000`, start: 200000,
			expected: []*storage.SelectHints{
				{Start: 60000, End: 180000, Range: 120000},
			},
		}, {
			query: `foo[2m] @ 300.000`, start: 200000,
			expected: []*storage.SelectHints{
				{Start: 180000, End: 300000, Range: 120000},
			},
		}, {
			query: `foo[2m] @ 60.000`, start: 200000,
			expected: []*storage.SelectHints{
				{Start: -60000, End: 60000, Range: 120000},
			},
		}, {
			query: `foo[2m] offset 2m`, start: 300000,
			expected: []*storage.SelectHints{
				{Start: 60000, End: 180000, Range: 120000},
			},
		}, {
			query: `foo[2m] @ 200.000 offset 2m`, start: 300000,
			expected: []*storage.SelectHints{
				{Start: -40000, End: 80000, Range: 120000},
			},
		}, {
			query: `foo[2m:1s]`, start: 300000,
			expected: []*storage.SelectHints{
				{Start: 175000, End: 300000, Step: 1000},
			},
		}, {
			query: `count_over_time(foo[2m:1s])`, start: 300000,
			expected: []*storage.SelectHints{
				{Start: 175000, End: 300000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: `count_over_time(foo[2m:1s] @ 300.000)`, start: 200000,
			expected: []*storage.SelectHints{
				{Start: 175000, End: 300000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: `count_over_time(foo[2m:1s] @ 200.000)`, start: 200000,
			expected: []*storage.SelectHints{
				{Start: 75000, End: 200000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: `count_over_time(foo[2m:1s] @ 100.000)`, start: 200000,
			expected: []*storage.SelectHints{
				{Start: -25000, End: 100000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: `count_over_time(foo[2m:1s] offset 10s)`, start: 300000,
			expected: []*storage.SelectHints{
				{Start: 165000, End: 290000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: `count_over_time((foo offset 10s)[2m:1s] offset 10s)`, start: 300000,
			expected: []*storage.SelectHints{
				{Start: 155000, End: 280000, Func: "count_over_time", Step: 1000},
			},
		}, {
			// When the @ is on the vector selector, the enclosing subquery parameters
			// don't affect the hint ranges.
			query: `count_over_time((foo @ 200.000 offset 10s)[2m:1s] offset 10s)`, start: 300000,
			expected: []*storage.SelectHints{
				{Start: 185000, End: 190000, Func: "count_over_time", Step: 1000},
			},
		}, {
			// When the @ is on the vector selector, the enclosing subquery parameters
			// don't affect the hint ranges.
			query: `count_over_time((foo @ 200.000 offset 10s)[2m:1s] @ 100.000 offset 10s)`, start: 300000,
			expected: []*storage.SelectHints{
				{Start: 185000, End: 190000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: `count_over_time((foo offset 10s)[2m:1s] @ 100.000 offset 10s)`, start: 300000,
			expected: []*storage.SelectHints{
				{Start: -45000, End: 80000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: `foo`, start: 10000, end: 20000,
			expected: []*storage.SelectHints{
				{Start: 5000, End: 20000, Step: 1000},
			},
		}, {
			query: `foo @ 15.000`, start: 10000, end: 20000,
			expected: []*storage.SelectHints{
				{Start: 10000, End: 15000, Step: 1000},
			},
		}, {
			query: `foo @ 1.000`, start: 10000, end: 20000,
			expected: []*storage.SelectHints{
				{Start: -4000, End: 1000, Step: 1000},
			},
		}, {
			query: `rate(foo[2m] @ 180.000)`, start: 200000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 60000, End: 180000, Range: 120000, Func: "rate", Step: 1000},
			},
		}, {
			query: `rate(foo[2m] @ 300.000)`, start: 200000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 180000, End: 300000, Range: 120000, Func: "rate", Step: 1000},
			},
		}, {
			query: `rate(foo[2m] @ 60.000)`, start: 200000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: -60000, End: 60000, Range: 120000, Func: "rate", Step: 1000},
			},
		}, {
			query: `rate(foo[2m])`, start: 200000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 80000, End: 500000, Range: 120000, Func: "rate", Step: 1000},
			},
		}, {
			query: `rate(foo[2m] offset 2m)`, start: 300000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 60000, End: 380000, Range: 120000, Func: "rate", Step: 1000},
			},
		}, {
			query: `rate(foo[2m:1s])`, start: 300000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 175000, End: 500000, Func: "rate", Step: 1000},
			},
		}, {
			query: `count_over_time(foo[2m:1s])`, start: 300000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 175000, End: 500000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: `count_over_time(foo[2m:1s] offset 10s)`, start: 300000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 165000, End: 490000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: `count_over_time(foo[2m:1s] @ 300.000)`, start: 200000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 175000, End: 300000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: `count_over_time(foo[2m:1s] @ 200.000)`, start: 200000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 75000, End: 200000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: `count_over_time(foo[2m:1s] @ 100.000)`, start: 200000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: -25000, End: 100000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: `count_over_time((foo offset 10s)[2m:1s] offset 10s)`, start: 300000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 155000, End: 480000, Func: "count_over_time", Step: 1000},
			},
		}, {
			// When the @ is on the vector selector, the enclosing subquery parameters
			// don't affect the hint ranges.
			query: `count_over_time((foo @ 200.000 offset 10s)[2m:1s] offset 10s)`, start: 300000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 185000, End: 190000, Func: "count_over_time", Step: 1000},
			},
		}, {
			// When the @ is on the vector selector, the enclosing subquery parameters
			// don't affect the hint ranges.
			query: `count_over_time((foo @ 200.000 offset 10s)[2m:1s] @ 100.000 offset 10s)`, start: 300000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 185000, End: 190000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: `count_over_time((foo offset 10s)[2m:1s] @ 100.000 offset 10s)`, start: 300000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: -45000, End: 80000, Func: "count_over_time", Step: 1000},
			},
		}, {
			query: `sum by (dim1) (foo)`, start: 10000,
			expected: []*storage.SelectHints{
				{Start: 5000, End: 10000, Func: "sum", By: true, Grouping: []string{"dim1"}},
			},
		}, {
			query: `sum without (dim1) (foo)`, start: 10000,
			expected: []*storage.SelectHints{
				{Start: 5000, End: 10000, Func: "sum", Grouping: []string{"dim1"}},
			},
		}, {
			query: `sum by (dim1) (avg_over_time(foo[1s]))`, start: 10000,
			expected: []*storage.SelectHints{
				{Start: 9000, End: 10000, Func: "avg_over_time", Range: 1000},
			},
		}, {
			query: `sum by (dim1) (max by (dim2) (foo))`, start: 10000,
			expected: []*storage.SelectHints{
				{Start: 5000, End: 10000, Func: "max", By: true, Grouping: []string{"dim2"}},
			},
		}, {
			query: `(max by (dim1) (foo))[5s:1s]`, start: 10000,
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
			query: `foo @ 50.000 + bar @ 250.000 + baz @ 900.000`, start: 100000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 45000, End: 50000, Step: 1000},
				{Start: 245000, End: 250000, Step: 1000},
				{Start: 895000, End: 900000, Step: 1000},
			},
		}, {
			query: `foo @ 50.000 + bar + baz @ 900.000`, start: 100000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 45000, End: 50000, Step: 1000},
				{Start: 95000, End: 500000, Step: 1000},
				{Start: 895000, End: 900000, Step: 1000},
			},
		}, {
			query: `rate(foo[2s] @ 50.000) + bar @ 250.000 + baz @ 900.000`, start: 100000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 48000, End: 50000, Step: 1000, Func: "rate", Range: 2000},
				{Start: 245000, End: 250000, Step: 1000},
				{Start: 895000, End: 900000, Step: 1000},
			},
		}, {
			query: `rate(foo[2s:1s] @ 50.000) + bar + baz`, start: 100000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 43000, End: 50000, Step: 1000, Func: "rate"},
				{Start: 95000, End: 500000, Step: 1000},
				{Start: 95000, End: 500000, Step: 1000},
			},
		}, {
			query: `rate(foo[2s:1s] @ 50.000) + bar + rate(baz[2m:1s] @ 900.000 offset 2m)`, start: 100000, end: 500000,
			expected: []*storage.SelectHints{
				{Start: 43000, End: 50000, Step: 1000, Func: "rate"},
				{Start: 95000, End: 500000, Step: 1000},
				{Start: 655000, End: 780000, Step: 1000, Func: "rate"},
			},
		}, { // Hints are based on the inner most subquery timestamp.
			query: `sum_over_time(sum_over_time(metric{job="1"}[1m40s])[1m40s:25s] @ 50.000)[3s:1s] @ 3000.000`, start: 100000,
			expected: []*storage.SelectHints{
				{Start: -150000, End: 50000, Range: 100000, Func: "sum_over_time", Step: 25000},
			},
		}, { // Hints are based on the inner most subquery timestamp.
			query: `sum_over_time(sum_over_time(metric{job="1"}[1m40s])[1m40s:25s] @ 3000.000)[3s:1s] @ 50.000`,
			expected: []*storage.SelectHints{
				{Start: 2800000, End: 3000000, Range: 100000, Func: "sum_over_time", Step: 25000},
			},
		},
	} {
		tc := tc
		t.Run(tc.query, func(t *testing.T) {
			t.Parallel()
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

	cases := []struct {
		name  string
		query string
	}{
		{
			name:  "unsupported function with scalar",
			query: `quantile_over_time(scalar(sum(http_requests_total)), http_requests_total[2m])`,
		},
	}

	load := `load 30s
				http_requests_total{pod="nginx-1"} 1+1x1
				http_requests_total{pod="nginx-2"} 1+2x40`

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			for _, disableFallback := range []bool{true, false} {
				t.Run(fmt.Sprintf("disableFallback=%t", disableFallback), func(t *testing.T) {
					opts := promql.EngineOpts{
						Timeout:    2 * time.Second,
						MaxSamples: math.MaxInt64,
					}
					newEngine := engine.New(engine.Opts{DisableFallback: disableFallback, EngineOpts: opts})
					q1, err := newEngine.NewRangeQuery(context.Background(), storage, nil, tcase.query, start, end, step)
					if disableFallback {
						testutil.NotOk(t, err)
					} else {
						testutil.Ok(t, err)
						newResult := q1.Exec(context.Background())
						testutil.Ok(t, newResult.Err)
					}
				})
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

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	ctx := context.Background()
	newEngine := engine.New(engine.Opts{DisableFallback: true, EngineOpts: opts})
	q, err := newEngine.NewRangeQuery(ctx, storage, nil, query, start, end, step)
	testutil.Ok(t, err)
	stats.NewQueryStats(q.Stats())

	q, err = newEngine.NewInstantQuery(ctx, storage, nil, query, end)
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

func (m *mockIterator) AtHistogram(_ *histogram.Histogram) (int64, *histogram.Histogram) {
	return 0, nil
}

func (m *mockIterator) AtFloatHistogram(_ *histogram.FloatHistogram) (int64, *histogram.FloatHistogram) {
	return 0, nil
}

func (m *mockIterator) AtT() int64 { return m.timestamps[m.i] }

func (m *mockIterator) Err() error { return nil }

type slowSeriesSet struct {
	empty bool
	delay time.Duration
}

func newSlowSeriesSet(delay time.Duration) *slowSeriesSet {
	return &slowSeriesSet{delay: delay}
}

func (s *slowSeriesSet) Next() bool {
	if s.empty {
		return false
	}
	s.empty = true
	<-time.After(s.delay)
	return true
}

func (s slowSeriesSet) At() storage.Series {
	return storage.MockSeries([]int64{0}, []float64{0}, nil)
}

func (s slowSeriesSet) Err() error { return nil }

func (s slowSeriesSet) Warnings() annotations.Annotations { return nil }

type testSeriesSet struct {
	i      int
	series []storage.Series
	warns  annotations.Annotations
}

func newTestSeriesSet(series ...storage.Series) storage.SeriesSet {
	return &testSeriesSet{
		i:      -1,
		series: series,
	}
}

func newWarningsSeriesSet(warns annotations.Annotations) storage.SeriesSet {
	return &testSeriesSet{
		i:     -1,
		warns: warns,
	}
}

func (s *testSeriesSet) Next() bool                        { s.i++; return s.i < len(s.series) }
func (s *testSeriesSet) At() storage.Series                { return s.series[s.i] }
func (s *testSeriesSet) Err() error                        { return nil }
func (s *testSeriesSet) Warnings() annotations.Annotations { return s.warns }

type slowSeries struct{}

func (d slowSeries) Labels() labels.Labels                        { return labels.FromStrings("foo", "bar") }
func (d slowSeries) Iterator(chunkenc.Iterator) chunkenc.Iterator { return &slowIterator{} }

type slowIterator struct {
	ts int64
}

func (d *slowIterator) AtHistogram(_ *histogram.Histogram) (int64, *histogram.Histogram) {
	panic("not implemented")
}

func (d *slowIterator) AtFloatHistogram(_ *histogram.FloatHistogram) (int64, *histogram.FloatHistogram) {
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
	start                  time.Time
	wantEmptyForMixedTypes bool
}

type histogramGeneratorFunc func(app storage.Appender, numSeries int, withMixedTypes bool) error

func TestNativeHistograms(t *testing.T) {
	t.Parallel()
	opts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e16,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}

	cases := []histogramTestCase{
		{
			name:  "count_over_time() with different start time",
			query: `count_over_time(native_histogram_series[1m15s])`,
			start: time.Unix(400, 0),
		},
		{
			name:  "irate()",
			query: `irate(native_histogram_series[1m])`,
		},
		{
			name:  "rate()",
			query: `rate(native_histogram_series[1m])`,
		},

		{
			name:  "increase()",
			query: `increase(native_histogram_series[1m])`,
		},
		{
			name:  "delta()",
			query: `delta(native_histogram_series[1m])`,
		},
		{
			name:                   "sum()",
			query:                  `sum(native_histogram_series)`,
			wantEmptyForMixedTypes: true,
		},
		{
			name:                   "sum by (foo)",
			query:                  `sum by (foo) (native_histogram_series)`,
			wantEmptyForMixedTypes: true,
		},
		{
			name:  "count",
			query: `count(native_histogram_series)`,
		},
		{
			name:  "count by (foo)",
			query: `count by (foo) (native_histogram_series)`,
		},
		// TODO(fpetkovski): The Prometheus engine returns an incorrect result for this case.
		// Uncomment once it gets fixed: https://github.com/prometheus/prometheus/issues/11973.
		// {
		//	name:  "max",
		//	query: "max (native_histogram_series)",
		// },
		{
			name:  "max by (foo)",
			query: `max by (foo) (native_histogram_series)`,
		},
		// TODO(fpetkovski): The Prometheus engine returns an incorrect result for this case.
		// Uncomment once it gets fixed: https://github.com/prometheus/prometheus/issues/11973.
		// {
		//	name:  "min",
		//	query: "min (native_histogram_series)",
		// },
		{
			name:  "min by (foo)",
			query: `min by (foo) (native_histogram_series)`,
		},
		{
			name:  "absent",
			query: `absent(native_histogram_series)`,
		},
		{
			name:  "histogram_sum",
			query: `histogram_sum(native_histogram_series)`,
		},
		{
			name:  "histogram_count",
			query: `histogram_count(native_histogram_series)`,
		},
		{
			name:  "histogram_count of histogram product",
			query: `histogram_count(native_histogram_series * native_histogram_series)`,
		},
		{
			name:  "histogram_sum / histogram_count",
			query: `histogram_sum(native_histogram_series) / histogram_count(native_histogram_series)`,
		},
		{
			name:  "histogram_sum over histogram_fraction",
			query: `histogram_sum(scalar(histogram_quantile(1, sum(native_histogram_series))) * native_histogram_series)`,
		},
		{
			name: "histogram_sum over histogram_fraction",
			query: `
histogram_sum(
  scalar(histogram_fraction(-Inf, +Inf, sum(native_histogram_series))) * native_histogram_series
)`,
		},
		{
			name:  "histogram_quantile",
			query: `histogram_quantile(0.7, native_histogram_series)`,
		},
		{
			// Test strange query with a mix of histogram functions.
			name:  "histogram_quantile(histogram_sum)",
			query: `histogram_quantile(0.7, histogram_sum(native_histogram_series))`,
		},
		{
			name:  "histogram_count * histogram aggregation",
			query: `scalar(histogram_count(sum(native_histogram_series))) * sum(native_histogram_series)`,
		},
		{
			name:  "histogram_fraction",
			query: `histogram_fraction(0, 0.2, native_histogram_series)`,
		},
		{
			name:  "lhs multiplication",
			query: `native_histogram_series * 3`,
		},
		{
			name:  "rhs multiplication",
			query: `3 * native_histogram_series`,
		},
		{
			name:  "lhs division",
			query: `native_histogram_series / 2`,
		},
		{
			name:  "subqueries",
			query: `increase(rate(native_histogram_series[2m])[2m:15s])`,
		},
	}

	defer pprof.StopCPUProfile()
	t.Run("integer_histograms", func(t *testing.T) {
		t.Parallel()
		testNativeHistograms(t, cases, opts, generateNativeHistogramSeries)
	})
	t.Run("float_histograms", func(t *testing.T) {
		t.Parallel()
		testNativeHistograms(t, cases, opts, generateFloatHistogramSeries)
	})
}

func testNativeHistograms(t *testing.T, cases []histogramTestCase, opts promql.EngineOpts, generateHistograms histogramGeneratorFunc) {
	numHistograms := 50
	mixedTypesOpts := []bool{false, true}
	var (
		queryStart = time.Unix(50, 0)
		queryEnd   = time.Unix(600, 0)
		queryStep  = 30 * time.Second
	)
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			for _, withMixedTypes := range mixedTypesOpts {
				t.Run(fmt.Sprintf("mixedTypes=%t", withMixedTypes), func(t *testing.T) {
					storage := teststorage.New(t)
					defer storage.Close()

					app := storage.Appender(context.TODO())
					err := generateHistograms(app, numHistograms, withMixedTypes)
					testutil.Ok(t, err)
					testutil.Ok(t, app.Commit())

					promEngine := promql.NewEngine(opts)
					thanosEngine := engine.New(engine.Opts{
						EngineOpts:        opts,
						DisableFallback:   true,
						LogicalOptimizers: logicalplan.AllOptimizers,
					})

					t.Run("instant", func(t *testing.T) {
						ctx := context.Background()
						q1, err := thanosEngine.NewInstantQuery(ctx, storage, nil, tc.query, time.Unix(50, 0))
						testutil.Ok(t, err)
						newResult := q1.Exec(ctx)
						testutil.Ok(t, newResult.Err)

						q2, err := promEngine.NewInstantQuery(ctx, storage, nil, tc.query, time.Unix(50, 0))
						testutil.Ok(t, err)
						promResult := q2.Exec(ctx)
						testutil.Ok(t, promResult.Err)
						promVector, err := promResult.Vector()
						testutil.Ok(t, err)

						// Make sure we're not getting back empty results.
						if withMixedTypes && tc.wantEmptyForMixedTypes {
							testutil.Assert(t, len(promVector) == 0)
							testutil.Equals(t, len(promResult.Warnings), len(newResult.Warnings))
						}

						testutil.WithGoCmp(comparer).Equals(t, promResult, newResult, queryExplanation(q1))
					})

					t.Run("range", func(t *testing.T) {
						if tc.start == (time.Time{}) {
							tc.start = queryStart
						}
						ctx := context.Background()
						q1, err := thanosEngine.NewRangeQuery(ctx, storage, nil, tc.query, tc.start, queryEnd, queryStep)
						testutil.Ok(t, err)
						newResult := q1.Exec(ctx)
						testutil.Ok(t, newResult.Err)

						q2, err := promEngine.NewRangeQuery(ctx, storage, nil, tc.query, tc.start, queryEnd, queryStep)
						testutil.Ok(t, err)
						promResult := q2.Exec(ctx)
						testutil.Ok(t, promResult.Err)
						promMatrix, err := promResult.Matrix()
						testutil.Ok(t, err)

						// Make sure we're not getting back empty results.
						if withMixedTypes && tc.wantEmptyForMixedTypes {
							testutil.Assert(t, len(promMatrix) == 0)
							testutil.Equals(t, len(promResult.Warnings), len(newResult.Warnings))
							testutil.Equals(t, "PromQL warning: encountered a mix of histograms and floats for aggregation", newResult.Warnings.AsErrors()[0].Error())
						}
						testutil.WithGoCmp(comparer).Equals(t, promResult, newResult, queryExplanation(q1))
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
	t.Parallel()
	histograms := tsdbutil.GenerateTestHistograms(2)

	storage := teststorage.New(t)
	defer storage.Close()

	lbls := []string{labels.MetricName, "native_histogram_series"}

	app := storage.Appender(context.TODO())
	_, err := app.AppendHistogram(0, labels.FromStrings(lbls...), 0, nil, histograms[0].ToFloat(nil))
	testutil.Ok(t, err)
	testutil.Ok(t, app.Commit())

	app = storage.Appender(context.TODO())
	_, err = app.AppendHistogram(0, labels.FromStrings(lbls...), 30_000, histograms[1], nil)
	testutil.Ok(t, err)
	testutil.Ok(t, app.Commit())

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

	ctx := context.Background()

	t.Run("vector_select", func(t *testing.T) {
		qry, err := engine.NewInstantQuery(ctx, storage, nil, "sum(native_histogram_series)", time.Unix(30, 0))
		testutil.Ok(t, err)
		res := qry.Exec(context.Background())
		testutil.Ok(t, res.Err)
		actual, err := res.Vector()
		testutil.Ok(t, err)

		testutil.Equals(t, 1, len(actual), "expected vector with 1 element")
		expected := histograms[1].ToFloat(nil)
		expected.CounterResetHint = histogram.UnknownCounterReset
		testutil.Equals(t, expected, actual[0].H)
	})

	t.Run("matrix_select", func(t *testing.T) {
		qry, err := engine.NewRangeQuery(ctx, storage, nil, "rate(native_histogram_series[1m])", time.Unix(0, 0), time.Unix(60, 0), 60*time.Second)
		testutil.Ok(t, err)
		res := qry.Exec(context.Background())
		testutil.Ok(t, res.Err)
		actual, err := res.Matrix()
		testutil.Ok(t, err)

		testutil.Equals(t, 1, len(actual), "expected 1 series")
		testutil.Equals(t, 1, len(actual[0].Histograms), "expected 1 point")
		expected := histograms[1].ToFloat(nil).Sub(histograms[0].ToFloat(nil)).Mul(1 / float64(30))
		expected.CounterResetHint = histogram.GaugeType
		testutil.Equals(t, expected, actual[0].Histograms[0].H)
	})
}

type seriesByLabels []promql.Series

func (b seriesByLabels) Len() int           { return len(b) }
func (b seriesByLabels) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b seriesByLabels) Less(i, j int) bool { return labels.Compare(b[i].Metric, b[j].Metric) < 0 }

type samplesByLabels []promql.Sample

func (b samplesByLabels) Len() int           { return len(b) }
func (b samplesByLabels) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b samplesByLabels) Less(i, j int) bool { return labels.Compare(b[i].Metric, b[j].Metric) < 0 }

var (
	// comparer should be used to compare promql results between engines.
	comparer = cmp.Comparer(func(x, y *promql.Result) bool {
		compareFloats := func(l, r float64) bool {
			const epsilon = 1e-6
			return cmp.Equal(l, r, cmpopts.EquateNaNs(), cmpopts.EquateApprox(0, epsilon))
		}
		compareHistograms := func(l, r *histogram.FloatHistogram) bool {
			if l == nil && r == nil {
				return true
			}
			return l.Equals(r)
		}
		compareAnnotations := func(l, r annotations.Annotations) bool {
			// TODO: discard promql annotations for now, once we support them we should add them back
			discardPromqlAnnotations := func(k string, _ error) bool {
				hasInfoPrefix := strings.HasPrefix(k, annotations.PromQLInfo.Error())
				hasWarnPrefix := strings.HasPrefix(k, annotations.PromQLWarning.Error())
				return hasInfoPrefix || hasWarnPrefix
			}
			maps.DeleteFunc(l, discardPromqlAnnotations)
			maps.DeleteFunc(r, discardPromqlAnnotations)

			if len(l) != len(r) {
				return false
			}
			for k, v := range l {
				if !cmp.Equal(r[k], v) {
					return false
				}
			}
			for k, v := range r {
				if !cmp.Equal(l[k], v) {
					return false
				}
			}
			return true
		}
		compareMetrics := func(l, r labels.Labels) bool {
			return l.Hash() == r.Hash()
		}

		if x.Err != nil && y.Err != nil {
			return cmp.Equal(x.Err.Error(), y.Err.Error())
		} else if x.Err != nil || y.Err != nil {
			return false
		}

		if !compareAnnotations(x.Warnings, y.Warnings) {
			return false
		}

		vx, xvec := x.Value.(promql.Vector)
		vy, yvec := y.Value.(promql.Vector)

		if xvec && yvec {
			if len(vx) != len(vy) {
				return false
			}

			// Sort vector before comparing.
			sort.Sort(samplesByLabels(vx))
			sort.Sort(samplesByLabels(vy))

			for i := 0; i < len(vx); i++ {
				if !compareMetrics(vx[i].Metric, vy[i].Metric) {
					return false
				}
				if vx[i].T != vy[i].T {
					return false
				}
				if !compareFloats(vx[i].F, vy[i].F) {
					return false
				}
				if !compareHistograms(vx[i].H, vy[i].H) {
					return false
				}
			}
			return true
		}

		mx, xmat := x.Value.(promql.Matrix)
		my, ymat := y.Value.(promql.Matrix)

		if xmat && ymat {
			if len(mx) != len(my) {
				return false
			}
			// Sort matrix before comparing.
			sort.Sort(seriesByLabels(mx))
			sort.Sort(seriesByLabels(my))
			for i := 0; i < len(mx); i++ {
				mxs := mx[i]
				mys := my[i]

				if !compareMetrics(mxs.Metric, mys.Metric) {
					return false
				}

				xps := mxs.Floats
				yps := mys.Floats

				if len(xps) != len(yps) {
					return false
				}
				for j := 0; j < len(xps); j++ {
					if xps[j].T != yps[j].T {
						return false
					}
					if !compareFloats(xps[j].F, yps[j].F) {
						return false
					}
				}
				xph := mxs.Histograms
				yph := mys.Histograms

				if len(xph) != len(yph) {
					return false
				}
				for j := 0; j < len(xph); j++ {
					if xph[j].T != yph[j].T {
						return false
					}
					if !compareHistograms(xph[j].H, yph[j].H) {
						return false
					}
				}
			}
			return true
		}

		sx, xscalar := x.Value.(promql.Scalar)
		sy, yscalar := y.Value.(promql.Scalar)
		if xscalar && yscalar {
			if sx.T != sy.T {
				return false
			}
			return compareFloats(sx.V, sy.V)
		}
		return false
	})
)

func queryExplanation(q promql.Query) string {
	eq, ok := q.(engine.ExplainableQuery)
	if !ok {
		return ""
	}

	var explain func(w io.Writer, n engine.ExplainOutputNode, indent, indentNext string)

	explain = func(w io.Writer, n engine.ExplainOutputNode, indent, indentNext string) {
		next := n.Children
		me := n.OperatorName

		_, _ = w.Write([]byte(indent))
		_, _ = w.Write([]byte(me))
		if len(next) == 0 {
			_, _ = w.Write([]byte("\n"))
			return
		}

		if me == "[*CancellableOperator]" {
			_, _ = w.Write([]byte(": "))
			explain(w, next[0], "", indentNext)
			return
		}
		_, _ = w.Write([]byte(":\n"))

		for i, n := range next {
			if i == len(next)-1 {
				explain(w, n, indentNext+"└──", indentNext+"   ")
			} else {
				explain(w, n, indentNext+"├──", indentNext+"│  ")
			}
		}
	}

	var b bytes.Buffer
	explain(&b, *eq.Explain(), "", "")

	return fmt.Sprintf("Query: %s\nExplanation:\n%s\n", q.String(), b.String())
}
