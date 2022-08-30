package plan_test

import (
	"context"
	"fmt"
	"fpetkovski/promql-engine/plan"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func BenchmarkPlan(b *testing.B) {
	query := "http_requests_total"
	test := setupStorage(b)
	defer test.Close()

	start := time.Unix(0, 0)
	end := start.Add(1 * time.Hour)
	step := time.Second * 30

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		expr, err := parser.ParseExpr(query)
		require.NoError(b, err)

		p, err := plan.New(expr, test.Storage(), start, end, step)
		require.NoError(b, err)

		out, err := p.Next(context.Background())
		require.NoError(b, err)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range out {
			}
		}()
		wg.Wait()
	}
}

func BenchmarkCurrentEngine(b *testing.B) {
	query := "http_requests_total"
	test := setupStorage(b)
	defer test.Close()

	start := time.Unix(0, 0)
	end := start.Add(1 * time.Hour)
	step := time.Second * 30

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
		qry, err := engine.NewRangeQuery(test.Queryable(), nil, query, start, end, step)
		require.NoError(b, err)

		res := qry.Exec(test.Context())
		require.NoError(b, res.Err)
	}
}

func setupStorage(b *testing.B) *promql.Test {
	load := synthesizeLoad(50, 3)
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
  http_requests_total{pod="p%d", container="c%d"} %d+%dx200`, i, j, i, j)
		}
	}

	return load
}
