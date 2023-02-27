// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"

	"github.com/thanos-community/promql-engine/engine"
	"github.com/thanos-community/promql-engine/execution/function"
)

func FuzzEngineQueryRangeMatrixFunctions(f *testing.F) {
	f.Add(uint32(0), uint32(120), uint32(30), 1.0, 1.0, 1.0, 2.0, 30)

	f.Fuzz(func(t *testing.T, startTS, endTS, intervalSeconds uint32, initialVal1, initialVal2, inc1, inc2 float64, stepRange int) {
		if math.IsNaN(initialVal1) || math.IsNaN(initialVal2) || math.IsNaN(inc1) || math.IsNaN(inc2) {
			return
		}
		if math.IsInf(initialVal1, 0) || math.IsInf(initialVal2, 0) || math.IsInf(inc1, 0) || math.IsInf(inc2, 0) {
			return
		}
		if inc1 < 0 || inc2 < 0 || stepRange <= 0 || intervalSeconds <= 0 || endTS < startTS {
			return
		}
		for funcName := range function.Funcs {
			// Skipping multi-arg functions in fuzz test for now.
			if len(parser.Functions[funcName].ArgTypes) != 1 || parser.Functions[funcName].ArgTypes[0] != parser.ValueTypeMatrix {
				continue
			}

			load := fmt.Sprintf(`load 30s
			http_requests_total{pod="nginx-1"} %.2f+%.2fx15
			http_requests_total{pod="nginx-2"} %2.f+%.2fx21`, initialVal1, inc1, initialVal2, inc2)

			opts := promql.EngineOpts{
				Timeout:              1 * time.Hour,
				MaxSamples:           1e10,
				EnableNegativeOffset: true,
				EnableAtModifier:     true,
			}

			test, err := promql.NewTest(t, load)
			testutil.Ok(t, err)
			defer test.Close()

			testutil.Ok(t, test.Run())

			start := time.Unix(int64(startTS), 0)
			end := time.Unix(int64(endTS), 0)
			interval := time.Duration(intervalSeconds) * time.Second
			query := fmt.Sprintf("%s(http_requests_total[%ds])", funcName, stepRange)
			if funcName == "vector" {
				query = fmt.Sprintf("vector(%d)", stepRange)
			}

			newEngine := engine.New(engine.Opts{EngineOpts: opts, DisableFallback: true})

			q1, err := newEngine.NewRangeQuery(test.Storage(), nil, query, start, end, interval)
			testutil.Ok(t, err)
			newResult := q1.Exec(context.Background())
			testutil.Ok(t, newResult.Err)

			oldEngine := promql.NewEngine(opts)
			q2, err := oldEngine.NewRangeQuery(test.Storage(), nil, query, start, end, interval)
			testutil.Ok(t, err)

			oldResult := q2.Exec(context.Background())
			testutil.Ok(t, oldResult.Err)

			testutil.Equals(t, oldResult, newResult, "inconsistent result for "+funcName)
		}
	})
}

func FuzzEngineInstantQueryAggregations(f *testing.F) {
	f.Add(uint32(0), 0, math.NaN(), 1.0, 1.0, 2.0)

	f.Fuzz(func(t *testing.T, ts uint32, groupingHash int, initialVal1, initialVal2, inc1, inc2 float64) {
		if inc1 < 0 || inc2 < 0 {
			return
		}
		for _, funcName := range []string{
			"stddev", "sum", "max", "min", "avg", "group", "stdvar", "count",
		} {
			load := fmt.Sprintf(`load 30s
			http_requests_total{pod="nginx-1", route="/"} %.2f+%.2fx4
			http_requests_total{pod="nginx-2", route="/"} %2.f+%.2fx4`, initialVal1, inc1, initialVal2, inc2)

			opts := promql.EngineOpts{
				Timeout:    1 * time.Hour,
				MaxSamples: 1e10,
			}

			test, err := promql.NewTest(t, load)
			testutil.Ok(t, err)
			defer test.Close()

			testutil.Ok(t, test.Run())

			queryTime := time.Unix(int64(ts), 0)

			newEngine := engine.New(engine.Opts{EngineOpts: opts, DisableFallback: true})

			var grouping string
			switch groupingHash % 3 {
			case 0:
				grouping = " by (pod)"
			case 1:
				grouping = " without (route)"
			default:
			}
			query := fmt.Sprintf("%s(http_requests_total)%s", funcName, grouping)
			q1, err := newEngine.NewInstantQuery(test.Storage(), nil, query, queryTime)
			testutil.Ok(t, err)
			newResult := q1.Exec(context.Background())
			testutil.Ok(t, newResult.Err)

			oldEngine := promql.NewEngine(opts)
			q2, err := oldEngine.NewInstantQuery(test.Storage(), nil, query, queryTime)
			testutil.Ok(t, err)

			oldResult := q2.Exec(context.Background())
			testutil.Ok(t, oldResult.Err)

			testutil.WithGoCmp(comparer).Equals(t, oldResult, newResult, query)
		}
	})
}

var comparer = cmp.Comparer(func(x, y *promql.Result) bool {
	compareFloats := func(l, r float64) bool {
		const epsilon = 1e-6

		if math.IsNaN(l) && math.IsNaN(r) {
			return true
		}
		if math.IsNaN(l) || math.IsNaN(r) {
			return false
		}

		return math.Abs(l-r) < epsilon
	}

	if x.Err != y.Err {
		return false
	}

	vx, xvec := x.Value.(promql.Vector)
	vy, yvec := y.Value.(promql.Vector)

	if xvec && yvec {
		if len(vx) != len(vy) {
			return false
		}
		for i := 0; i < len(vx); i++ {
			if !cmp.Equal(vx[i].Metric, vy[i].Metric) {
				return false
			}
			if vx[i].T != vy[i].T {
				return false
			}
			if !compareFloats(vx[i].V, vy[i].V) {
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
		for i := 0; i < len(mx); i++ {
			mxs := mx[i]
			mys := my[i]

			if !cmp.Equal(mxs.Metric, mys.Metric) {
				return false
			}

			xps := mxs.Points
			yps := mys.Points

			if len(xps) != len(yps) {
				return false
			}
			for j := 0; j < len(xps); j++ {
				if xps[j].T != yps[j].T {
					return false
				}
				if !compareFloats(xps[j].V, yps[i].V) {
					return false
				}
			}
		}
		return true
	}
	return false
})
