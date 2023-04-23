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
	"github.com/thanos-community/promql-engine/logicalplan"
)

func BenchmarkChunkDecoding(b *testing.B) {
	test := setupStorage(b, 1000, 3, 720)
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
				series = append(series, ss.At().Iterator(nil))
			}
			for i := 0; i < len(series); i++ {
				for ts := start.UnixMilli(); ts <= end.UnixMilli(); ts += step.Milliseconds() {
					numIterations++
					if val := series[i].Seek(ts); val == chunkenc.ValNone {
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
				series = append(series, ss.At().Iterator(nil))
			}
			stepCount := 10
			ts := start.UnixMilli()
			for ts <= end.UnixMilli() {
				for i := 0; i < len(series); i++ {
					seriesTs := ts
					for currStep := 0; currStep < stepCount && seriesTs <= end.UnixMilli(); currStep++ {
						numIterations++
						if valType := series[i].Seek(seriesTs); valType == chunkenc.ValNone {
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
	test := setupStorage(b, 5000, 3, 720)
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
	samplesPerHour := 60 * 2
	sixHourDataset := setupStorage(b, 1000, 3, 6*samplesPerHour)
	defer sixHourDataset.Close()

	largeSixHourDataset := setupStorage(b, 10000, 10, 6*samplesPerHour)
	defer largeSixHourDataset.Close()

	sevenDaysAndTwoHoursDataset := setupStorage(b, 1000, 3, (7*24+2)*samplesPerHour)
	defer sevenDaysAndTwoHoursDataset.Close()

	start := time.Unix(0, 0)
	end := start.Add(2 * time.Hour)
	step := time.Second * 30

	cases := []struct {
		name  string
		query string
		test  *promql.Test
	}{
		{
			name:  "vector selector",
			query: "http_requests_total",
			test:  sixHourDataset,
		},
		{
			name:  "sum",
			query: "sum(http_requests_total)",
			test:  sixHourDataset,
		},
		{
			name:  "sum by pod",
			query: "sum by (pod) (http_requests_total)",
			test:  sixHourDataset,
		},
		{
			name:  "topk",
			query: "topk(2,http_requests_total)",
			test:  sixHourDataset,
		},
		{
			name:  "bottomk",
			query: "bottomk(2,http_requests_total)",
			test:  sixHourDataset,
		},
		{
			name:  "rate",
			query: "rate(http_requests_total[1m])",
			test:  sixHourDataset,
		},
		{
			name:  "rate with large range selection",
			query: "rate(http_requests_total[7d])",
			test:  sevenDaysAndTwoHoursDataset,
		},
		{
			name:  "rate with large number of series, 1m range",
			query: "rate(http_requests_total[1m])",
			test:  largeSixHourDataset,
		},
		{
			name:  "rate with large number of series, 5m range",
			query: "rate(http_requests_total[5m])",
			test:  largeSixHourDataset,
		},
		{
			name:  "sum rate",
			query: "sum(rate(http_requests_total[1m]))",
			test:  sixHourDataset,
		},
		{
			name:  "sum by rate",
			query: "sum by (pod) (rate(http_requests_total[1m]))",
			test:  sixHourDataset,
		},
		{
			name:  "quantile with variable parameter",
			query: "quantile by (pod) (scalar(min(http_requests_total)), http_requests_total)",
			test:  sixHourDataset,
		},
		{
			name:  "binary operation with one to one",
			query: `http_requests_total{container="c1"} / ignoring(container) http_responses_total`,
			test:  sixHourDataset,
		},
		{
			name:  "binary operation with many to one",
			query: `http_requests_total / on (pod) group_left http_responses_total`,
			test:  sixHourDataset,
		},
		{
			name:  "binary operation with vector and scalar",
			query: `http_requests_total * 10`,
			test:  sixHourDataset,
		},
		{
			name:  "unary negation",
			query: `-http_requests_total`,
			test:  sixHourDataset,
		},
		{
			name:  "vector and scalar comparison",
			query: `http_requests_total > 10`,
			test:  sixHourDataset,
		},
		{
			name:  "positive offset vector",
			query: "http_requests_total offset 5m",
			test:  sixHourDataset,
		},
		{
			name:  "at modifier ",
			query: "http_requests_total @ 600",
			test:  sixHourDataset,
		},
		{
			name:  "at modifier with positive offset vector",
			query: "http_requests_total @ 600 offset 5m",
			test:  sixHourDataset,
		},
		{
			name:  "clamp",
			query: `clamp(http_requests_total, 5, 10)`,
			test:  sixHourDataset,
		},
		{
			name:  "clamp_min",
			query: `clamp_min(http_requests_total, 10)`,
			test:  sixHourDataset,
		},
		{
			name:  "complex func query",
			query: `clamp(1 - http_requests_total, 10 - 5, 10)`,
			test:  sixHourDataset,
		},
		{
			name:  "func within func query",
			query: `clamp(irate(http_requests_total[30s]), 10 - 5, 10)`,
			test:  sixHourDataset,
		},
		{
			name:  "aggr within func query",
			query: `clamp(rate(http_requests_total[30s]), 10 - 5, 10)`,
			test:  sixHourDataset,
		},
		{
			name:  "histogram_quantile",
			query: `histogram_quantile(0.9, http_response_seconds_bucket)`,
			test:  sixHourDataset,
		},
		{
			name:  "sort",
			query: `sort(http_requests_total)`,
			test:  sixHourDataset,
		},
		{
			name:  "sort_desc",
			query: `sort_desc(http_requests_total)`,
			test:  sixHourDataset,
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
					qry, err := engine.NewRangeQuery(tc.test.Context(), tc.test.Queryable(), nil, tc.query, start, end, step)
					testutil.Ok(b, err)

					oldResult := qry.Exec(tc.test.Context())
					testutil.Ok(b, oldResult.Err)
				}
			})
			b.Run("new_engine", func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()

				for i := 0; i < b.N; i++ {
					newResult := executeRangeQuery(b, tc.query, tc.test, start, end, step)
					testutil.Ok(b, newResult.Err)
				}
			})
		})
	}
}

func BenchmarkNativeHistograms(b *testing.B) {
	test, err := promql.NewTest(b, "")
	testutil.Ok(b, err)
	defer test.Close()

	app := test.Storage().Appender(context.TODO())
	testutil.Ok(b, generateNativeHistogramSeries(app, 3000, false))
	testutil.Ok(b, app.Commit())
	testutil.Ok(b, test.Run())

	start := time.Unix(0, 0)
	end := start.Add(2 * time.Hour)
	step := time.Second * 30

	cases := []struct {
		name  string
		query string
	}{
		{
			name:  "selector",
			query: "native_histogram_series",
		},
		{
			name:  "sum",
			query: "sum(native_histogram_series)",
		},
		{
			name:  "rate",
			query: "rate(native_histogram_series[1m])",
		},
		{
			name:  "sum rate",
			query: "sum(rate(native_histogram_series[1m]))",
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
			query: "histogram_quantile(0.9, sum(native_histogram_series))",
		},
	}

	opts := promql.EngineOpts{
		Logger:               nil,
		Reg:                  nil,
		MaxSamples:           50000000,
		Timeout:              100 * time.Second,
		EnableAtModifier:     true,
		EnableNegativeOffset: true,
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.Run("old_engine", func(b *testing.B) {
				engine := promql.NewEngine(opts)

				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					qry, err := engine.NewRangeQuery(test.Context(), test.Queryable(), nil, tc.query, start, end, step)
					testutil.Ok(b, err)

					oldResult := qry.Exec(test.Context())
					testutil.Ok(b, oldResult.Err)
				}
			})
			b.Run("new_engine", func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()

				for i := 0; i < b.N; i++ {
					ng := engine.New(engine.Opts{EngineOpts: opts})

					qry, err := ng.NewRangeQuery(test.Context(), test.Queryable(), nil, tc.query, start, end, step)
					testutil.Ok(b, err)

					newResult := qry.Exec(context.Background())
					testutil.Ok(b, newResult.Err)
				}
			})
		})
	}
}

func BenchmarkInstantQuery(b *testing.B) {
	test := setupStorage(b, 1000, 3, 720)
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
			name:  "round",
			query: "round(http_requests_total)",
		},
		{
			name:  "round with argument",
			query: "round(http_requests_total, 0.5)",
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
		{
			name:  "sort",
			query: `sort(http_requests_total)`,
		},
		{
			name:  "sort_desc",
			query: `sort_desc(http_requests_total)`,
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
					qry, err := engine.NewInstantQuery(test.Context(), test.Queryable(), nil, tc.query, queryTime)
					testutil.Ok(b, err)

					res := qry.Exec(test.Context())
					testutil.Ok(b, res.Err)
				}
			})
			b.Run("new_engine", func(b *testing.B) {
				ng := engine.New(engine.Opts{EngineOpts: promql.EngineOpts{Timeout: 100 * time.Second}})
				b.ResetTimer()
				b.ReportAllocs()

				for i := 0; i < b.N; i++ {
					qry, err := ng.NewInstantQuery(test.Context(), test.Queryable(), nil, tc.query, queryTime)
					testutil.Ok(b, err)

					res := qry.Exec(context.Background())
					testutil.Ok(b, res.Err)
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
			opts := engine.Opts{
				LogicalOptimizers: logicalplan.NoOptimizers,
				EngineOpts:        promql.EngineOpts{Timeout: 100 * time.Second},
			}
			ng := engine.New(opts)
			ctx := context.Background()
			qry, err := ng.NewRangeQuery(ctx, db, nil, query, start, end, step)
			testutil.Ok(b, err)

			res := qry.Exec(ctx)
			testutil.Ok(b, res.Err)
		}
	})
	b.Run("withOptimizers", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ng := engine.New(engine.Opts{EngineOpts: promql.EngineOpts{Timeout: 100 * time.Second}})
			ctx := context.Background()
			qry, err := ng.NewRangeQuery(ctx, db, nil, query, start, end, step)
			testutil.Ok(b, err)

			res := qry.Exec(ctx)
			testutil.Ok(b, res.Err)
		}
	})

}

func executeRangeQuery(b *testing.B, q string, test *promql.Test, start time.Time, end time.Time, step time.Duration) *promql.Result {
	return executeRangeQueryWithOpts(b, q, test, start, end, step, engine.Opts{DisableFallback: true, EngineOpts: promql.EngineOpts{Timeout: 100 * time.Second}})
}

func executeRangeQueryWithOpts(b *testing.B, q string, test *promql.Test, start time.Time, end time.Time, step time.Duration, opts engine.Opts) *promql.Result {
	ng := engine.New(opts)
	ctx := context.Background()
	qry, err := ng.NewRangeQuery(ctx, test.Queryable(), nil, q, start, end, step)
	testutil.Ok(b, err)

	return qry.Exec(ctx)
}

func setupStorage(b *testing.B, numLabelsA int, numLabelsB int, numSteps int) *promql.Test {
	load := synthesizeLoad(numLabelsA, numLabelsB, numSteps)
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

func synthesizeLoad(numPods, numContainers, numSteps int) string {
	load := "load 30s\n"
	for i := 0; i < numPods; i++ {
		for j := 0; j < numContainers; j++ {
			load += fmt.Sprintf(`http_requests_total{pod="p%d", container="c%d"} %d+%dx%d%s`, i, j, i, j, numSteps, "\n")
		}
		load += fmt.Sprintf(`http_responses_total{pod="p%d"} %dx%d%s`, i, i, numSteps, "\n")
	}

	for i := 0; i < numPods; i++ {
		for j := 0; j < 10; j++ {
			load += fmt.Sprintf(`http_response_seconds_bucket{pod="p%d", le="%d"} %d+%dx%d%s`, i, j, i, j, numSteps, "\n")
		}
		load += fmt.Sprintf(`http_response_seconds_bucket{pod="p%d", le="+Inf"} %d+%dx%d%s`, i, i, i, numSteps, "\n")
	}

	return load
}
