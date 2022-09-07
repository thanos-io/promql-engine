package engine_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb/chunkenc"

	"github.com/fpetkovski/promql-engine/engine"

	"github.com/prometheus/prometheus/promql"
	"github.com/stretchr/testify/require"
)

func BenchmarkChunkDecoding(b *testing.B) {
	test := setupStorage(b)
	defer test.Close()

	start := time.Unix(0, 0)
	end := start.Add(1 * time.Hour)
	step := time.Second * 30

	querier, err := test.Storage().Querier(test.Context(), start.UnixMilli(), end.UnixMilli())
	require.NoError(b, err)

	matcher, err := labels.NewMatcher(labels.MatchEqual, labels.MetricName, "http_requests_total")
	require.NoError(b, err)
	ss := querier.Select(false, nil, matcher)

	b.Run("iterate by series", func(b *testing.B) {
		series := make([]chunkenc.Iterator, 0)
		for ss.Next() {
			series = append(series, ss.At().Iterator())
		}
		b.ResetTimer()
		for c := 0; c < b.N; c++ {
			for i := 0; i < len(series); i++ {
				for ts := start.UnixMilli(); ts <= end.UnixMilli(); ts += step.Milliseconds() {
					series[i].Seek(ts)
				}
			}
		}
	})
	b.Run("iterate by time", func(b *testing.B) {
		series := make([]chunkenc.Iterator, 0)
		for ss.Next() {
			series = append(series, ss.At().Iterator())
		}
		b.ResetTimer()
		for c := 0; c < b.N; c++ {
			for ts := start.UnixMilli(); ts <= end.UnixMilli(); ts += step.Milliseconds() {
				for i := 0; i < len(series); i++ {
					series[i].Seek(ts)
				}
			}
		}
	})
}

func BenchmarkSingleQuery(b *testing.B) {
	test := setupStorage(b)
	defer test.Close()

	start := time.Unix(0, 0)
	end := start.Add(1 * time.Hour)
	step := time.Second * 30

	query := "sum by (pod) (rate(http_requests_total[1m]))"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		executeQuery(b, query, test, start, end, step)
	}
}

func BenchmarkOldEngine(b *testing.B) {
	test := setupStorage(b)
	defer test.Close()

	start := time.Unix(0, 0)
	end := start.Add(1 * time.Hour)
	step := time.Second * 30

	cases := []struct {
		name  string
		query string
	}{
		{
			name:  "vector selector",
			query: "http_requests_total",
		},
		{
			name:  "aggregation",
			query: "sum by (pod) (http_requests_total)",
		},
		{
			name:  "sum-rate",
			query: "sum by (pod) (rate(http_requests_total[1m]))",
		},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.Run("current_engine", func(b *testing.B) {
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
					qry, err := engine.NewRangeQuery(test.Queryable(), nil, tc.query, start, end, step)
					require.NoError(b, err)

					res := qry.Exec(test.Context())
					require.NoError(b, res.Err)
				}
			})
			b.Run("new_engine", func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()

				for i := 0; i < b.N; i++ {
					executeQuery(b, tc.query, test, start, end, step)
				}
			})
		})
	}
}

func executeQuery(b *testing.B, q string, test *promql.Test, start time.Time, end time.Time, step time.Duration) {
	ng := engine.New()
	qry, err := ng.NewRangeQuery(test.Queryable(), nil, q, start, end, step)
	require.NoError(b, err)

	qry.Exec(context.Background())
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
