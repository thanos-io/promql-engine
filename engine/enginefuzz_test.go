package engine_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/promql"
	"github.com/thanos-community/promql-engine/engine"
)

func FuzzEngineInstantQueryAggregations(f *testing.F) {
	load := `load 30s
	http_requests_total{pod="nginx-1"} 1+1x4
	http_requests_total{pod="nginx-2"} 1+2x4`

	opts := promql.EngineOpts{
		Timeout:    1 * time.Hour,
		MaxSamples: 1e10,
	}

	test, err := promql.NewTest(f, load)
	testutil.Ok(f, err)
	defer test.Close()

	testutil.Ok(f, test.Run())

	f.Add(uint32(0), true)

	f.Fuzz(func(t *testing.T, ts uint32, by bool) {
		for _, funcName := range []string{
			"stddev", "sum", "max", "min", "avg", "group", "stdvar", "count",
		} {
			queryTime := time.Unix(int64(ts), 0)

			newEngine := engine.New(engine.Opts{EngineOpts: opts, DisableFallback: true})

			var byOp string
			if by {
				byOp = " by (pod)"
			}
			query := fmt.Sprintf("%s(http_requests_total)%s", funcName, byOp)
			t.Log("query is", query)
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
