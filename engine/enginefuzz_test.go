// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/logicalplan"

	"github.com/cortexproject/promqlsmith"
	"github.com/efficientgo/core/errors"
	"github.com/efficientgo/core/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/promql/promqltest"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/util/stats"
	"github.com/stretchr/testify/require"
)

const testRuns = 100

type testCase struct {
	query              string
	loads              []string
	oldRes, newRes     *promql.Result
	oldStats, newStats *stats.Statistics
	start, end         time.Time
	interval           time.Duration
	validateSamples    bool
}

type testType int

const (
	testTypeFloat           testType = iota // 0
	testTypeNativeHistogram                 // 1
)

// shouldValidateSamples checks if the samples can be compared for the expr.
// For certain known cases, Thanos engine returns less samples than Prometheus engine due to optimizations.
func shouldValidateSamples(expr parser.Expr) bool {
	valid := true

	parser.Inspect(expr, func(node parser.Node, path []parser.Node) error {
		switch n := node.(type) {
		case *parser.Call:
			switch n.Func.Name {
			case "scalar":
				// Optimized to step_invariant in Thanos engine.
				valid = false
				return errors.New("error")
			case "histogram_count", "histogram_sum", "histogram_avg":
				// Optimized using DetectHistogramStatsOptimizer and will return smaller samples than Prometheus engine.
				valid = false
				return errors.New("error")
			}
		}
		return nil
	})
	return valid
}

// validateExpr checks if the given expression is valid for fuzz tests.
// For certain known cases Thanos engine results do not match with Prometheus engine.
func validateExpr(expr parser.Expr, testType testType) bool {
	expr, _ = promql.PreprocessExpr(expr, time.Unix(0, 0), time.Unix(0, 0), 0)
	valid := true

	parser.Inspect(expr, func(node parser.Node, path []parser.Node) error {
		switch n := node.(type) {
		case *parser.Call:
			switch n.Func.Name {
			case "sort", "sort_desc", "sort_by_label", "sort_by_label_desc":
				if testType == testTypeNativeHistogram {
					// Prometheus engine filters out native histograms in nested sort().
					// Thanos engine implements sorting only at the presentation time and ignores nested sort().
					// See: https://github.com/thanos-io/promql-engine/pull/595
					valid = false
					return errors.New("error")
				}
			case "predict_linear":
				switch t := n.Args[0].(type) {
				case *parser.StepInvariantExpr:
					// Thanos engine cannot correctly handle a MatrixSelector wrapped by StepInvariant.
					// eg: predict_linear({__name__="http_request_duration_seconds"}[5m] @ end(), 0.5)
					// See: https://github.com/thanos-io/promql-engine/pull/527
					if _, ok := t.Expr.(*parser.MatrixSelector); ok {
						valid = false
						return errors.New("error")
					}
				}
			case "timestamp":
				if testType == testTypeNativeHistogram {
					// TODO(johrry): Remove after merging https://github.com/thanos-io/promql-engine/pull/598
					valid = false
					return errors.New("error")
				}
			}
		}
		return nil
	})
	return valid
}

func FuzzEnginePromQLSmithRangeQuery(f *testing.F) {
	f.Add(int64(0), uint32(0), uint32(120), uint32(30), 1.0, 1.0, 1.0, 2.0, 30)

	f.Fuzz(func(t *testing.T, seed int64, startTS, endTS, intervalSeconds uint32, initialVal1, initialVal2, inc1, inc2 float64, stepRange int) {
		if math.IsNaN(initialVal1) || math.IsNaN(initialVal2) || math.IsNaN(inc1) || math.IsNaN(inc2) {
			return
		}
		if math.IsInf(initialVal1, 0) || math.IsInf(initialVal2, 0) || math.IsInf(inc1, 0) || math.IsInf(inc2, 0) {
			return
		}
		if inc1 < 0 || inc2 < 0 || stepRange <= 0 || intervalSeconds <= 0 || endTS < startTS {
			return
		}
		rnd := rand.New(rand.NewSource(seed))

		load := fmt.Sprintf(`load 30s
			http_requests_total{pod="nginx-1"} %.2f+%.2fx15
			http_requests_total{pod="nginx-2"} %2.f+%.2fx21`, initialVal1, inc1, initialVal2, inc2)

		opts := promql.EngineOpts{
			Timeout:              1 * time.Hour,
			MaxSamples:           1e10,
			EnableNegativeOffset: true,
			EnableAtModifier:     true,
			EnablePerStepStats:   true,
		}
		qOpts := promql.NewPrometheusQueryOpts(true, 0)

		storage := promqltest.LoadedStorage(t, load)
		defer storage.Close()

		start := time.Unix(int64(startTS), 0)
		end := time.Unix(int64(endTS), 0)
		interval := time.Duration(intervalSeconds) * time.Second

		seriesSet, err := getSeries(context.Background(), storage, "http_requests_total")
		require.NoError(t, err)
		psOpts := []promqlsmith.Option{
			promqlsmith.WithEnableOffset(true),
			promqlsmith.WithEnableAtModifier(true),
			// bottomk and topk sometimes lead to random failures since their result on equal values is essentially random
			promqlsmith.WithEnabledAggrs([]parser.ItemType{parser.SUM, parser.MIN, parser.MAX, parser.AVG, parser.GROUP, parser.COUNT, parser.COUNT_VALUES, parser.QUANTILE}),
		}
		ps := promqlsmith.New(rnd, seriesSet, psOpts...)

		newEngine := engine.New(engine.Opts{EngineOpts: opts, EnableAnalysis: true})
		oldEngine := promql.NewEngine(opts)

		var (
			q1              promql.Query
			query           string
			validateSamples bool
		)
		cases := make([]*testCase, testRuns)
		for i := range testRuns {
			for {
				expr := ps.WalkRangeQuery()
				if !validateExpr(expr, testTypeFloat) {
					continue
				}
				validateSamples = shouldValidateSamples(expr)

				query = expr.Pretty(0)
				q1, err = newEngine.NewRangeQuery(context.Background(), storage, qOpts, query, start, end, interval)
				if engine.IsUnimplemented(err) || errors.As(err, &parser.ParseErrors{}) {
					continue
				} else {
					break
				}
			}

			testutil.Ok(t, err)
			newResult := q1.Exec(context.Background())
			newStats := q1.Stats()
			stats.NewQueryStats(newStats)

			q2, err := oldEngine.NewRangeQuery(context.Background(), storage, qOpts, query, start, end, interval)
			testutil.Ok(t, err)

			oldResult := q2.Exec(context.Background())
			oldStats := q2.Stats()
			stats.NewQueryStats(oldStats)

			cases[i] = &testCase{
				query:           query,
				newRes:          newResult,
				newStats:        newStats,
				oldRes:          oldResult,
				oldStats:        oldStats,
				loads:           []string{load},
				start:           start,
				end:             end,
				interval:        interval,
				validateSamples: validateSamples,
			}
		}
		validateTestCases(t, cases)
	})
}

func FuzzEnginePromQLSmithInstantQuery(f *testing.F) {
	f.Add(int64(0), uint32(0), 1.0, 1.0, 1.0, 2.0)

	f.Fuzz(func(t *testing.T, seed int64, ts uint32, initialVal1, initialVal2, inc1, inc2 float64) {
		t.Parallel()
		if inc1 < 0 || inc2 < 0 {
			return
		}
		rnd := rand.New(rand.NewSource(seed))

		load := fmt.Sprintf(`load 30s
			http_requests_total{pod="nginx-1", route="/"} %.2f+%.2fx40
			http_requests_total{pod="nginx-2", route="/"} %2.f+%.2fx40`, initialVal1, inc1, initialVal2, inc2)

		opts := promql.EngineOpts{
			Timeout:              1 * time.Hour,
			MaxSamples:           1e10,
			EnableNegativeOffset: true,
			EnableAtModifier:     true,
			EnablePerStepStats:   true,
		}
		qOpts := promql.NewPrometheusQueryOpts(true, 0)

		storage := promqltest.LoadedStorage(t, load)
		defer storage.Close()

		queryTime := time.Unix(int64(ts), 0)
		newEngine := engine.New(engine.Opts{
			EngineOpts:        opts,
			LogicalOptimizers: logicalplan.AllOptimizers,
			EnableAnalysis:    true,
		})
		oldEngine := promql.NewEngine(opts)

		seriesSet, err := getSeries(context.Background(), storage, "http_requests_total")
		require.NoError(t, err)
		psOpts := []promqlsmith.Option{
			promqlsmith.WithEnableOffset(true),
			promqlsmith.WithEnableAtModifier(true),
			promqlsmith.WithAtModifierMaxTimestamp(180 * 1000),
			// bottomk and topk sometimes lead to random failures since their result on equal values is essentially random
			promqlsmith.WithEnabledAggrs([]parser.ItemType{parser.SUM, parser.MIN, parser.MAX, parser.AVG, parser.GROUP, parser.COUNT, parser.COUNT_VALUES, parser.QUANTILE}),
		}
		ps := promqlsmith.New(rnd, seriesSet, psOpts...)

		var (
			q1              promql.Query
			query           string
			validateSamples bool
		)
		cases := make([]*testCase, testRuns)
		for i := range testRuns {
			// Since we disabled fallback, keep trying until we find a query
			// that can be natively execute by the engine.
			// Parsing experimental function, like mad_over_time, will lead to a parser.ParseErrors, so we also ignore those.
			for {
				expr := ps.WalkInstantQuery()
				if !validateExpr(expr, testTypeFloat) {
					continue
				}
				validateSamples = shouldValidateSamples(expr)
				query = expr.Pretty(0)
				q1, err = newEngine.NewInstantQuery(context.Background(), storage, qOpts, query, queryTime)
				if engine.IsUnimplemented(err) || errors.As(err, &parser.ParseErrors{}) {
					continue
				} else {
					break
				}
			}

			testutil.Ok(t, err)
			newResult := q1.Exec(context.Background())
			newStats := q1.Stats()
			stats.NewQueryStats(newStats)

			q2, err := oldEngine.NewInstantQuery(context.Background(), storage, qOpts, query, queryTime)
			testutil.Ok(t, err)

			oldResult := q2.Exec(context.Background())
			oldStats := q2.Stats()
			stats.NewQueryStats(oldStats)

			cases[i] = &testCase{
				query:           query,
				newRes:          newResult,
				newStats:        newStats,
				oldRes:          oldResult,
				oldStats:        oldStats,
				loads:           []string{load},
				start:           queryTime,
				end:             queryTime,
				validateSamples: validateSamples,
			}
		}
		validateTestCases(t, cases)
	})
}

func getSeries(ctx context.Context, q storage.Queryable, query string) ([]labels.Labels, error) {
	querier, err := q.Querier(0, time.Now().Unix())
	if err != nil {
		return nil, err
	}
	res := make([]labels.Labels, 0)
	ss := querier.Select(ctx, false, &storage.SelectHints{Func: "series"}, labels.MustNewMatcher(labels.MatchEqual, "__name__", query))
	for ss.Next() {
		lbls := ss.At().Labels()
		res = append(res, lbls)
	}
	if err := ss.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func validateTestCases(t *testing.T, cases []*testCase) {
	failures := 0
	logQuery := func(c *testCase) {
		for _, load := range c.loads {
			t.Log(load)
		}
		t.Logf("query: %s, start: %d, end: %d, interval: %v", c.query, c.start.UnixMilli(), c.end.UnixMilli(), c.interval)
	}
	for i, c := range cases {
		if !cmp.Equal(c.oldRes, c.newRes, comparer) {
			logQuery(c)
			t.Logf("case %d error mismatch.\nnew result: %s\nold result: %s\n", i, c.newRes.String(), c.oldRes.String())
			failures++
			continue
		}
		if !c.validateSamples || c.oldRes.Err != nil {
			// Skip sample comparison
			continue
		}
		if !cmp.Equal(c.oldStats.Samples, c.newStats.Samples, samplesComparer) {
			logQuery(c)
			t.Logf("case: %d, samples mismatch. total samples: old: %v, new: %v. samples per step: old: %v, new: %v", i, c.oldStats.Samples.TotalSamples, c.newStats.Samples.TotalSamples, c.oldStats.Samples.TotalSamplesPerStep, c.newStats.Samples.TotalSamplesPerStep)
			failures++
		}
	}
	if failures > 0 {
		t.Fatalf("failed %d test cases", failures)
	}
}

func normalizeBuckets(a, b, c uint64) []uint64 {
	// Ensure strictly increasing positive values (cumulative form).
	buckets := []uint64{a, b, c}
	sort.Slice(buckets, func(i, j int) bool { return buckets[i] < buckets[j] })

	// Avoid zero counts
	for i := range buckets {
		if buckets[i] == 0 {
			buckets[i] = 1
		}
	}

	return buckets
}

func sumBuckets(buckets []uint64) uint64 {
	var sum uint64
	for _, b := range buckets {
		sum += b
	}
	return sum
}

// Produces a realistic sum based on schema and bucket index.
func estimateSum(schema int8, buckets []uint64) float64 {
	base := math.Pow(2, math.Pow(2, float64(-schema)))
	var sum float64
	for i, count := range buckets {
		value := math.Pow(base, float64(i+1))
		sum += float64(count) * value
	}
	return sum
}

func FuzzNativeHistogramQuery(f *testing.F) {
	f.Add(int64(0), uint32(0), uint32(60), uint32(120), int8(0), int8(0), uint64(1), uint64(2), uint64(1))

	f.Fuzz(func(t *testing.T, seed int64, startTS, endTS, intervalSeconds uint32, schema1 int8, schema2 int8, b1, b2, b3 uint64) {
		t.Parallel()
		if endTS < startTS || intervalSeconds <= 0 {
			return
		}

		if schema1 < -4 || schema1 > 8 || schema2 < -4 || schema2 > 8 {
			return
		}

		// Ensure positive, increasing bucket counts (mimicking cumulative).
		bucket1 := normalizeBuckets(b1, b2, b3)
		bucket2 := normalizeBuckets(b1*2, b2+5, b3*2)

		count1 := sumBuckets(bucket1)
		count2 := sumBuckets(bucket2)

		if count1 == 0 || count2 == 0 {
			return
		}

		sum1 := estimateSum(schema1, bucket1)
		sum2 := estimateSum(schema2, bucket2)

		load := fmt.Sprintf(`load 2m
			http_request_duration_seconds{pod="nginx-1"} {{schema:%d count:%d sum:%.2f buckets:%v}}+{{schema:%d count:%d buckets:%v}}x20
			http_request_duration_seconds{pod="nginx-2"} {{schema:%d count:%d sum:%.2f buckets:%v}}+{{schema:%d count:%d buckets:%v}}x30`,
			schema1, count1, sum1, bucket1,
			schema1, count1, bucket1,
			schema2, count2, sum2, bucket2,
			schema2, count2, bucket2,
		)

		opts := promql.EngineOpts{
			Timeout:              1 * time.Hour,
			MaxSamples:           1e10,
			EnableNegativeOffset: true,
			EnableAtModifier:     true,
			EnablePerStepStats:   true,
		}

		qOpts := promql.NewPrometheusQueryOpts(true, 0)
		queryable := promqltest.LoadedStorage(t, load)
		t.Cleanup(func() {
			queryable.Close()
		})

		startTime := time.Unix(int64(startTS), 0)
		endTime := time.Unix(int64(endTS), 0)
		interval := time.Duration(intervalSeconds) * time.Second

		seriesSet, err := getSeries(context.Background(), queryable, "http_request_duration_seconds")
		require.NoError(t, err)

		psOpts := []promqlsmith.Option{
			promqlsmith.WithEnableOffset(true),
			promqlsmith.WithEnableAtModifier(true),
			promqlsmith.WithEnabledAggrs([]parser.ItemType{
				parser.SUM, parser.MIN, parser.MAX, parser.AVG, parser.GROUP, parser.COUNT, parser.COUNT_VALUES, parser.QUANTILE,
			}),
		}

		rnd := rand.New(rand.NewSource(seed))
		ps := promqlsmith.New(rnd, seriesSet, psOpts...)
		newEngine := engine.New(engine.Opts{EngineOpts: opts, EnableAnalysis: true})
		oldEngine := promql.NewEngine(opts)

		instantCases := make([]*testCase, 0, testRuns/2)
		rangeCases := make([]*testCase, 0, testRuns/2)

		var allQueries []promql.Query
		t.Cleanup(func() {
			for _, q := range allQueries {
				if q != nil {
					q.Close()
				}
			}
		})

		for range testRuns / 2 {
			var (
				qInstant                       promql.Query
				qRange                         promql.Query
				instantQuery                   string
				rangeQuery                     string
				validateSamplesForInstantQuery bool
				validateSamplesForRangeQuery   bool
			)

			for {
				expr := ps.WalkInstantQuery()
				if !validateExpr(expr, testTypeNativeHistogram) {
					continue
				}

				validateSamplesForInstantQuery = shouldValidateSamples(expr)
				instantQuery = expr.Pretty(0)

				qInstant, err = newEngine.NewInstantQuery(context.Background(), queryable, qOpts, instantQuery, startTime)
				allQueries = append(allQueries, qInstant)
				if engine.IsUnimplemented(err) || errors.As(err, &parser.ParseErrors{}) {
					continue
				} else {
					break
				}
			}
			testutil.Ok(t, err)

			for {
				expr := ps.WalkRangeQuery()
				if !validateExpr(expr, testTypeNativeHistogram) {
					continue
				}

				validateSamplesForRangeQuery = shouldValidateSamples(expr)
				rangeQuery = expr.Pretty(0)

				qRange, err = newEngine.NewRangeQuery(context.Background(), queryable, qOpts, rangeQuery, startTime, endTime, interval)
				allQueries = append(allQueries, qRange)
				if engine.IsUnimplemented(err) || errors.As(err, &parser.ParseErrors{}) {
					continue
				} else {
					break
				}
			}
			testutil.Ok(t, err)

			// Instant query
			newInstantResult := qInstant.Exec(context.Background())
			newInstantStats := qInstant.Stats()
			stats.NewQueryStats(newInstantStats)

			q2Instant, err := oldEngine.NewInstantQuery(context.Background(), queryable, qOpts, instantQuery, startTime)
			testutil.Ok(t, err)
			allQueries = append(allQueries, q2Instant)

			oldInstantResult := q2Instant.Exec(context.Background())
			oldInstantStats := q2Instant.Stats()
			stats.NewQueryStats(oldInstantStats)

			// Range query
			newRangeResult := qRange.Exec(context.Background())
			newRangeStats := qRange.Stats()
			stats.NewQueryStats(newRangeStats)

			q2Range, err := oldEngine.NewRangeQuery(context.Background(), queryable, qOpts, rangeQuery, startTime, endTime, interval)
			testutil.Ok(t, err)
			allQueries = append(allQueries, q2Range)

			oldRangeResult := q2Range.Exec(context.Background())
			oldRangeStats := q2Range.Stats()
			stats.NewQueryStats(oldRangeStats)

			instantCases = append(instantCases, &testCase{
				query:           instantQuery,
				newRes:          newInstantResult,
				newStats:        newInstantStats,
				oldRes:          oldInstantResult,
				oldStats:        oldInstantStats,
				loads:           []string{load},
				start:           startTime,
				end:             startTime,
				interval:        0,
				validateSamples: validateSamplesForInstantQuery,
			})

			rangeCases = append(rangeCases, &testCase{
				query:           rangeQuery,
				newRes:          newRangeResult,
				newStats:        newRangeStats,
				oldRes:          oldRangeResult,
				oldStats:        oldRangeStats,
				loads:           []string{load},
				start:           startTime,
				end:             endTime,
				interval:        interval,
				validateSamples: validateSamplesForRangeQuery,
			})
		}

		validateTestCases(t, instantCases)
		validateTestCases(t, rangeCases)
	})
}
