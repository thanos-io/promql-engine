package executionplan_test

import (
	"context"
	"fmt"
	"fpetkovski/promql-engine/executionplan"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var queries = []string{
	"http_requests_total",
	"sum by (pod) (http_requests_total)",
}

func BenchmarkExecutionPlan(b *testing.B) {
	test := setupStorage(b)
	defer test.Close()

	start := time.Unix(0, 0)
	end := start.Add(1 * time.Hour)
	step := time.Second * 30

	for _, q := range queries {
		b.Run("current_engine", func(b *testing.B) {
			b.Run(q, func(b *testing.B) {
				opts := promql.EngineOpts{
					Logger:     nil,
					Reg:        nil,
					MaxSamples: 50000000,
					Timeout:    100 * time.Second,
				}
				engine := promql.NewEngine(opts)

				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					qry, err := engine.NewRangeQuery(test.Queryable(), nil, q, start, end, step)
					require.NoError(b, err)

					res := qry.Exec(test.Context())
					require.NoError(b, res.Err)
				}
			})
		})
		b.Run("new_engine", func(b *testing.B) {
			b.Run(q, func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()

				for i := 0; i < b.N; i++ {
					expr, err := parser.ParseExpr(q)
					require.NoError(b, err)

					p, err := executionplan.New(expr, test.Storage(), start, end, step)
					require.NoError(b, err)

					out, err := p.Next(context.Background())
					require.NoError(b, err)
					for range out {
					}
				}
			})
		})
	}
}

func setupStorage(b *testing.B) *promql.Test {
	load := synthesizeLoad(1000, 3)
	test, err := promql.NewTest(b, load)
	require.NoError(b, err)
	require.NoError(b, test.Run())

	return test
}

func synthesizeLoad(numPods, numContainers int) string {
	load := `
load 30s`
	for i := 0; i < numPods; i++ {
		for j := 0; j < numContainers; j++ {
			load += fmt.Sprintf(`
  http_requests_total{pod="p%d", container="c%d"} %d+%dx720`, i, j, i, j)
		}
	}

	return load
}
