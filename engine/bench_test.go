// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/prometheus/prometheus/tsdb/chunkenc"

	"github.com/thanos-community/promql-engine/engine"
)

func BenchmarkChunkDecoding(b *testing.B) {
	test := setupStorage(b, 1000, 3)
	defer test.Close()

	start := time.Unix(0, 0)
	end := start.Add(6 * time.Hour)
	step := time.Second * 30

	querier, err := test.Storage().Querier(test.Context(), start.UnixMilli(), end.UnixMilli())
	testutil.Ok(b, err)

	matcher, err := labels.NewMatcher(labels.MatchEqual, labels.MetricName, "http_requests_total")
	testutil.Ok(b, err)

	b.Run("iterate by series", func(b *testing.B) {
		b.ResetTimer()
		for c := 0; c < b.N; c++ {
			numIterations := 0

			ss := querier.Select(false, nil, matcher)
			series := make([]chunkenc.Iterator, 0)
			for ss.Next() {
				series = append(series, ss.At().Iterator())
			}
			for i := 0; i < len(series); i++ {
				for ts := start.UnixMilli(); ts <= end.UnixMilli(); ts += step.Milliseconds() {
					numIterations++
					if ok := series[i].Seek(ts); !ok {
						break
					}
				}
			}
		}
	})
	b.Run("iterate by time", func(b *testing.B) {
		b.ResetTimer()
		for c := 0; c < b.N; c++ {
			numIterations := 0
			ss := querier.Select(false, nil, matcher)
			series := make([]chunkenc.Iterator, 0)
			for ss.Next() {
				series = append(series, ss.At().Iterator())
			}
			stepCount := 10
			ts := start.UnixMilli()
			for ts <= end.UnixMilli() {
				for i := 0; i < len(series); i++ {
					seriesTs := ts
					for currStep := 0; currStep < stepCount && seriesTs <= end.UnixMilli(); currStep++ {
						numIterations++
						if ok := series[i].Seek(seriesTs); !ok {
							break
						}
						seriesTs += step.Milliseconds()
					}
				}
				ts += step.Milliseconds() * int64(stepCount)
			}
		}
	})
}

func BenchmarkSingleQuery(b *testing.B) {
	test := setupStorage(b, 5000, 3)
	defer test.Close()

	start := time.Unix(0, 0)
	end := start.Add(6 * time.Hour)
	step := time.Second * 30

	query := "sum(rate(http_requests_total[2m]))"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		result := executeRangeQuery(b, query, test, start, end, step)
		testutil.Ok(b, result.Err)
	}
}

func BenchmarkRangeQuery(b *testing.B) {
	test := setupStorage(b, 1000, 3)
	defer test.Close()

	start := time.Unix(0, 0)
	end := start.Add(2 * time.Hour)
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
			name:  "sum",
			query: "sum(http_requests_total)",
		},
		{
			name:  "sum by pod",
			query: "sum by (pod) (http_requests_total)",
		},
		{
			name:  "rate",
			query: "rate(http_requests_total[1m])",
		},
		{
			name:  "sum rate",
			query: "sum(rate(http_requests_total[1m]))",
		},
		{
			name:  "sum by rate",
			query: "sum by (pod) (rate(http_requests_total[1m]))",
		},
		{
			name:  "binary operation with one to one",
			query: `http_requests_total{container="c1"} / ignoring(container) http_responses_total`,
		},
		{
			name:  "binary operation with many to one",
			query: `http_requests_total / on (pod) group_left http_responses_total`,
		},
		{
			name:  "binary operation with vector and scalar",
			query: `http_requests_total * 10`,
		},
		{
			name:  "unary negation",
			query: `-http_requests_total`,
		},
		{
			name:  "vector and scalar comparison",
			query: `http_requests_total > 10`,
		},
		{
			name:  "positive offset vector",
			query: "http_requests_total offset 5m",
		},
		{
			name:  "at modifier ",
			query: "http_requests_total @ 600",
		},
		{
			name:  "at modifier with positive offset vector",
			query: "http_requests_total @ 600 offset 5m",
		},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.Run("old_engine", func(b *testing.B) {
				opts := promql.EngineOpts{
					Logger:               nil,
					Reg:                  nil,
					MaxSamples:           50000000,
					Timeout:              100 * time.Second,
					EnableAtModifier:     true,
					EnableNegativeOffset: true,
				}
				engine := promql.NewEngine(opts)

				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					qry, err := engine.NewRangeQuery(test.Queryable(), nil, tc.query, start, end, step)
					testutil.Ok(b, err)

					oldResult := qry.Exec(test.Context())
					testutil.Ok(b, oldResult.Err)
				}
			})
			b.Run("new_engine", func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()

				for i := 0; i < b.N; i++ {
					newResult := executeRangeQuery(b, tc.query, test, start, end, step)
					testutil.Ok(b, newResult.Err)
				}
			})
		})
	}
}

func BenchmarkOldEngineInstant(b *testing.B) {
	test := setupStorage(b, 1000, 3)
	defer test.Close()

	queryTime := time.Unix(50, 0)

	cases := []struct {
		name  string
		query string
	}{
		{
			name:  "vector selector",
			query: "http_requests_total",
		},
		{
			name:  "count",
			query: "count(http_requests_total)",
		},
		{
			name:  "avg",
			query: "avg(http_requests_total)",
		},
		{
			name:  "sum",
			query: "sum(http_requests_total)",
		},
		{
			name:  "sum by pod",
			query: "sum by (pod) (http_requests_total)",
		},
		{
			name:  "rate",
			query: "rate(http_requests_total[1m])",
		},
		{
			name:  "sum rate",
			query: "sum(rate(http_requests_total[1m]))",
		},
		{
			name:  "sum by rate",
			query: "sum by (pod) (rate(http_requests_total[1m]))",
		},
		{
			name:  "binary operation with many to one",
			query: `http_requests_total / on (pod) group_left http_responses_total`,
		},
		{
			name:  "unary negation",
			query: `-http_requests_total`,
		},
		{
			name:  "vector and scalar comparison",
			query: `http_requests_total > 10`,
		},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.Run("current_engine", func(b *testing.B) {
				opts := promql.EngineOpts{
					Logger:               nil,
					Reg:                  nil,
					MaxSamples:           50000000,
					Timeout:              100 * time.Second,
					EnableAtModifier:     true,
					EnableNegativeOffset: true,
				}
				engine := promql.NewEngine(opts)

				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					qry, err := engine.NewInstantQuery(test.Queryable(), nil, tc.query, queryTime)
					testutil.Ok(b, err)

					res := qry.Exec(test.Context())
					testutil.Ok(b, res.Err)
				}
			})
			b.Run("new_engine", func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()

				for i := 0; i < b.N; i++ {
					executeInstantQuery(b, tc.query, test, queryTime)
				}
			})
		})
	}
}

func BenchmarkMergeSelectorsOptimizer(b *testing.B) {
	db := createRequestsMetricBlock(b, 10000, 9900)

	start := time.Unix(0, 0)
	end := start.Add(6 * time.Hour)
	step := time.Second * 30

	query := `sum(http_requests_total{code="200"}) / sum(http_requests_total)`
	b.Run("withoutOptimizers", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			opts := engine.Opts{EnableOptimizers: false}
			ng := engine.New(opts)
			qry, err := ng.NewRangeQuery(db, nil, query, start, end, step)
			testutil.Ok(b, err)

			res := qry.Exec(context.Background())
			testutil.Ok(b, res.Err)
		}
	})
	b.Run("withOptimizers", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			opts := engine.Opts{EnableOptimizers: true}
			ng := engine.New(opts)
			qry, err := ng.NewRangeQuery(db, nil, query, start, end, step)
			testutil.Ok(b, err)

			res := qry.Exec(context.Background())
			testutil.Ok(b, res.Err)
		}
	})

}

func executeRangeQuery(b *testing.B, q string, test *promql.Test, start time.Time, end time.Time, step time.Duration) *promql.Result {
	return executeRangeQueryWithOpts(b, q, test, start, end, step, engine.Opts{})
}

func executeRangeQueryWithOpts(b *testing.B, q string, test *promql.Test, start time.Time, end time.Time, step time.Duration, opts engine.Opts) *promql.Result {
	ng := engine.New(opts)
	qry, err := ng.NewRangeQuery(test.Queryable(), nil, q, start, end, step)
	testutil.Ok(b, err)

	return qry.Exec(context.Background())
}

func executeInstantQuery(b *testing.B, q string, test *promql.Test, start time.Time) {
	ng := engine.New(engine.Opts{})
	qry, err := ng.NewInstantQuery(test.Queryable(), nil, q, start)
	testutil.Ok(b, err)

	qry.Exec(context.Background())
}

func setupStorage(b *testing.B, numLabelsA int, numLabelsB int) *promql.Test {
	load := synthesizeLoad(numLabelsA, numLabelsB)
	test, err := promql.NewTest(b, load)
	testutil.Ok(b, err)
	testutil.Ok(b, test.Run())

	return test
}

func createRequestsMetricBlock(b *testing.B, numRequests int, numSuccess int) *tsdb.DB {
	dir := b.TempDir()

	db, err := tsdb.Open(dir, nil, nil, tsdb.DefaultOptions(), nil)
	testutil.Ok(b, err)
	appender := db.Appender(context.Background())

	sixHours := int64(6 * 60 * 2)

	for i := 0; i < numRequests; i++ {
		for t := int64(0); t < sixHours; t += 30 {
			code := "200"
			if numSuccess < i {
				code = "500"
			}
			lbls := labels.FromStrings(labels.MetricName, "http_requests_total", "code", code, "pod", strconv.Itoa(i))
			_, err = appender.Append(0, lbls, t, 1)
			testutil.Ok(b, err)
		}
	}

	testutil.Ok(b, appender.Commit())

	return db
}

func synthesizeLoad(numPods, numContainers int) string {
	load := `
load 30s`
	for i := 0; i < numPods; i++ {
		for j := 0; j < numContainers; j++ {
			load += fmt.Sprintf(`
  http_requests_total{pod="p%d", container="c%d"} %d+%dx720`, i, j, i, j)
		}
		load += fmt.Sprintf(`
  http_responses_total{pod="p%d"} %dx720`, i, i)
	}

	return load
}
