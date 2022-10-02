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
	"github.com/prometheus/prometheus/promql"

	"github.com/thanos-community/promql-engine/engine"
)

func FuzzEngineInstantQueryAggregations(f *testing.F) {

	f.Add(uint32(0), true, 1.0, 1.0, 1.0, 2.0)

	f.Fuzz(func(t *testing.T, ts uint32, by bool, initialVal1, initialVal2, inc1, inc2 float64) {
		if math.IsNaN(initialVal1) || math.IsNaN(initialVal2) || math.IsNaN(inc1) || math.IsNaN(inc2) {
			return
		}
		if math.IsInf(initialVal1, 0) || math.IsInf(initialVal2, 0) || math.IsInf(inc1, 0) || math.IsInf(inc2, 0) {
			return
		}
		if inc1 < 0 || inc2 < 0 {
			return
		}
		for _, funcName := range []string{
			"stddev", "sum", "max", "min", "avg", "group", "stdvar", "count",
		} {
			load := fmt.Sprintf(`load 30s
			http_requests_total{pod="nginx-1"} %.2f+%.2fx4
			http_requests_total{pod="nginx-2"} %2.f+%.2fx4`, initialVal1, inc1, initialVal2, inc2)

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

			var byOp string
			if by {
				byOp = " by (pod)"
			}
			query := fmt.Sprintf("%s(http_requests_total)%s", funcName, byOp)
			q1, err := newEngine.NewInstantQuery(test.Storage(), nil, query, queryTime)
			testutil.Ok(t, err)
			newResult := q1.Exec(context.Background())
			testutil.Ok(t, newResult.Err)

			oldEngine := promql.NewEngine(opts)
			q2, err := oldEngine.NewInstantQuery(test.Storage(), nil, query, queryTime)
			testutil.Ok(t, err)

			oldResult := q2.Exec(context.Background())
			testutil.Ok(t, oldResult.Err)

			testutil.Equals(t, oldResult, newResult)
		}

	})
}
