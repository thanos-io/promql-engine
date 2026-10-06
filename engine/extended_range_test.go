// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/logicalplan"
	"github.com/thanos-io/promql-engine/query"
	engstorage "github.com/thanos-io/promql-engine/storage"
	promscan "github.com/thanos-io/promql-engine/storage/prometheus"

	"github.com/efficientgo/core/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/promql/promqltest"
	promstorage "github.com/prometheus/prometheus/storage"
)

func TestAnchoredSmoothedModifiers(t *testing.T) {
	parser.EnableExtendedRangeSelectors = true

	opts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}

	cases := []struct {
		name  string
		load  string
		query string
		start time.Time
		end   time.Time
		step  time.Duration
	}{
		// Anchored rate/increase/delta on a linear counter.
		{
			name: "anchored rate on linear counter",
			load: `load 10s
			    http_total 0 10 20 30 40 50 60 70 80 90 100`,
			query: `rate(http_total[30s] anchored)`,
		},
		{
			name: "anchored increase on linear counter",
			load: `load 10s
			    http_total 0 10 20 30 40 50 60 70 80 90 100`,
			query: `increase(http_total[30s] anchored)`,
		},
		{
			name: "anchored delta on gauge",
			load: `load 10s
			    temperature 20 22 21 23 25 24 26 28 27 29 30`,
			query: `delta(temperature[30s] anchored)`,
		},
		// Smoothed rate/increase/delta on a linear counter.
		{
			name: "smoothed rate on linear counter",
			load: `load 10s
			    http_total 0 10 20 30 40 50 60 70 80 90 100`,
			query: `rate(http_total[30s] smoothed)`,
		},
		{
			name: "smoothed increase on linear counter",
			load: `load 10s
			    http_total 0 10 20 30 40 50 60 70 80 90 100`,
			query: `increase(http_total[30s] smoothed)`,
		},
		{
			name: "smoothed delta on gauge",
			load: `load 10s
			    temperature 20 22 21 23 25 24 26 28 27 29 30`,
			query: `delta(temperature[30s] smoothed)`,
		},
		// Anchored with counter resets.
		{
			name: "anchored increase with counter reset",
			load: `load 10s
			    resets_total 0 10 20 5 15 25 10 20 30 40 50`,
			query: `increase(resets_total[30s] anchored)`,
		},
		{
			name: "anchored rate with counter reset",
			load: `load 10s
			    resets_total 0 10 20 5 15 25 10 20 30 40 50`,
			query: `rate(resets_total[30s] anchored)`,
		},
		// Smoothed with counter resets.
		{
			name: "smoothed increase with counter reset",
			load: `load 10s
			    resets_total 0 10 20 5 15 25 10 20 30 40 50`,
			query: `increase(resets_total[30s] smoothed)`,
		},
		{
			name: "smoothed rate with counter reset",
			load: `load 10s
			    resets_total 0 10 20 5 15 25 10 20 30 40 50`,
			query: `rate(resets_total[30s] smoothed)`,
		},
		// Anchored resets and changes (only supported for anchored).
		{
			name: "anchored resets",
			load: `load 10s
			    resets_total 0 10 20 5 15 25 10 20 30 40 50`,
			query: `resets(resets_total[30s] anchored)`,
		},
		{
			name: "anchored changes",
			load: `load 10s
			    metric 1 1 2 2 3 3 4 4 5 5 6`,
			query: `changes(metric[30s] anchored)`,
		},
		// Counter reset at range boundary (regression test: smoothed interpolation
		// must handle counter resets to avoid negative increase values).
		{
			name: "smoothed increase counter reset at boundary",
			load: `load 10s
			    counter_boundary 0 4 5 1 6 11`,
			query: `increase(counter_boundary[10s] smoothed)`,
		},
		{
			name: "smoothed rate counter reset at boundary",
			load: `load 10s
			    counter_boundary 0 4 5 1 6 11`,
			query: `rate(counter_boundary[10s] smoothed)`,
		},
		// Non-linear data.
		{
			name: "anchored rate on quadratic counter",
			load: `load 10s
			    quadratic 0 1 4 9 16 25 36 49 64 81 100`,
			query: `rate(quadratic[30s] anchored)`,
		},
		{
			name: "smoothed rate on quadratic counter",
			load: `load 10s
			    quadratic 0 1 4 9 16 25 36 49 64 81 100`,
			query: `rate(quadratic[30s] smoothed)`,
		},
		// Multiple series.
		{
			name: "anchored increase multiple series",
			load: `load 10s
			    http_total{path="/foo"} 0 5 10 15 20 25 30 35 40 45 50
			    http_total{path="/bar"} 0 10 20 30 40 50 60 70 80 90 100`,
			query: `increase(http_total[30s] anchored)`,
		},
		{
			name: "smoothed increase multiple series",
			load: `load 10s
			    http_total{path="/foo"} 0 5 10 15 20 25 30 35 40 45 50
			    http_total{path="/bar"} 0 10 20 30 40 50 60 70 80 90 100`,
			query: `increase(http_total[30s] smoothed)`,
		},
		// Stale markers are skipped, not treated as gaps.
		{
			name: "smoothed instant vector across stale marker",
			load: `load 10s
			    metric 1 2 stale 4 5 stale stale 8`,
			query: `metric smoothed`,
			step:  5 * time.Second,
		},
		{
			name: "smoothed instant vector ending with stale marker",
			load: `load 10s
			    metric 1 2 stale`,
			query: `metric smoothed`,
			step:  5 * time.Second,
		},
		{
			name: "smoothed rate across stale marker",
			load: `load 10s
			    http_total 0 10 stale 30 40 stale 60 70`,
			query: `rate(http_total[20s] smoothed)`,
		},
		{
			name: "anchored changes across stale marker",
			load: `load 10s
			    metric 1 2 stale 2 3 stale 3 4`,
			query: `changes(metric[20s] anchored)`,
		},
		// timestamp() ignores the smoothed modifier.
		{
			name: "timestamp of smoothed instant vector",
			load: `load 10s
			    metric 1 2 3 _ 5`,
			query: `timestamp(metric smoothed)`,
			step:  5 * time.Second,
		},
		// Native histograms are rejected wherever they appear in the window.
		{
			name: "anchored rate over float then histogram",
			load: `load 10s
			    metric 1 2 3 {{schema:0 sum:5 count:4 buckets:[1 2 1]}} {{schema:0 sum:6 count:5 buckets:[1 2 2]}}`,
			query: `rate(metric[30s] anchored)`,
		},
		{
			name: "smoothed increase over float then histogram",
			load: `load 10s
			    metric 1 2 3 {{schema:0 sum:5 count:4 buckets:[1 2 1]}} {{schema:0 sum:6 count:5 buckets:[1 2 2]}}`,
			query: `increase(metric[30s] smoothed)`,
		},
		{
			name: "anchored changes over histograms",
			load: `load 10s
			    metric {{schema:0 sum:5 count:4 buckets:[1 2 1]}} {{schema:0 sum:6 count:5 buckets:[1 2 2]}}`,
			query: `changes(metric[30s] anchored)`,
		},
		{
			name: "smoothed instant vector over histograms",
			load: `load 10s
			    metric {{schema:0 sum:5 count:4 buckets:[1 2 1]}} 2 3`,
			query: `metric smoothed`,
		},
	}

	start := time.Unix(0, 0)
	end := time.Unix(100, 0)
	step := 10 * time.Second

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			storage := promqltest.LoadedStorage(t, tc.load)
			defer storage.Close()

			tcStart := start
			if !tc.start.IsZero() {
				tcStart = tc.start
			}
			tcEnd := end
			if !tc.end.IsZero() {
				tcEnd = tc.end
			}
			tcStep := step
			if tc.step != 0 {
				tcStep = tc.step
			}

			for _, disableOptimizers := range []bool{false, true} {
				t.Run(fmt.Sprintf("disableOptimizers=%v", disableOptimizers), func(t *testing.T) {
					optimizers := logicalplan.AllOptimizers
					if disableOptimizers {
						optimizers = logicalplan.NoOptimizers
					}
					newEngine := engine.New(engine.Opts{
						EngineOpts:                   opts,
						LogicalOptimizers:            optimizers,
						SelectorBatchSize:            1,
						EnableExtendedRangeSelectors: true,
					})

					ctx := context.Background()
					q1, err := newEngine.NewRangeQuery(ctx, storage, nil, tc.query, tcStart, tcEnd, tcStep)
					testutil.Ok(t, err)
					defer q1.Close()
					newResult := q1.Exec(ctx)

					oldEngine := promql.NewEngine(opts)
					q2, err := oldEngine.NewRangeQuery(ctx, storage, nil, tc.query, tcStart, tcEnd, tcStep)
					testutil.Ok(t, err)
					defer q2.Close()
					oldResult := q2.Exec(ctx)

					testutil.WithGoCmp(comparer).Equals(t, oldResult, newResult, queryExplanation(q1))
				})
			}
		})
	}
}

func TestAnchoredSmoothedWhitelist(t *testing.T) {
	parser.EnableExtendedRangeSelectors = true

	opts := promql.EngineOpts{
		Timeout:    1 * time.Hour,
		MaxSamples: 1e10,
	}

	load := `load 10s
		metric 0 10 20 30 40 50`

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	newEngine := engine.New(engine.Opts{
		EngineOpts:                   opts,
		EnableExtendedRangeSelectors: true,
	})
	ctx := context.Background()

	// Functions not in the whitelist should error.
	unsupportedAnchored := []string{
		`avg_over_time(metric[30s] anchored)`,
		`sum_over_time(metric[30s] anchored)`,
		`max_over_time(metric[30s] anchored)`,
		`deriv(metric[30s] anchored)`,
	}
	for _, query := range unsupportedAnchored {
		t.Run("unsupported_anchored/"+query, func(t *testing.T) {
			q, err := newEngine.NewInstantQuery(ctx, storage, nil, query, time.Unix(50, 0))
			testutil.Ok(t, err)
			defer q.Close()
			res := q.Exec(ctx)
			testutil.NotOk(t, res.Err)
		})
	}

	unsupportedSmoothed := []string{
		`avg_over_time(metric[30s] smoothed)`,
		`resets(metric[30s] smoothed)`,
		`changes(metric[30s] smoothed)`,
	}
	for _, query := range unsupportedSmoothed {
		t.Run("unsupported_smoothed/"+query, func(t *testing.T) {
			q, err := newEngine.NewInstantQuery(ctx, storage, nil, query, time.Unix(50, 0))
			testutil.Ok(t, err)
			defer q.Close()
			res := q.Exec(ctx)
			testutil.NotOk(t, res.Err)
		})
	}

	// Functions in the whitelist should work.
	supportedAnchored := []string{
		`rate(metric[30s] anchored)`,
		`increase(metric[30s] anchored)`,
		`delta(metric[30s] anchored)`,
		`resets(metric[30s] anchored)`,
		`changes(metric[30s] anchored)`,
	}
	for _, query := range supportedAnchored {
		t.Run("supported_anchored/"+query, func(t *testing.T) {
			q, err := newEngine.NewInstantQuery(ctx, storage, nil, query, time.Unix(50, 0))
			testutil.Ok(t, err)
			q.Close()
		})
	}

	supportedSmoothed := []string{
		`rate(metric[30s] smoothed)`,
		`increase(metric[30s] smoothed)`,
		`delta(metric[30s] smoothed)`,
	}
	for _, query := range supportedSmoothed {
		t.Run("supported_smoothed/"+query, func(t *testing.T) {
			q, err := newEngine.NewInstantQuery(ctx, storage, nil, query, time.Unix(50, 0))
			testutil.Ok(t, err)
			q.Close()
		})
	}
}

// hintsSpy wraps real scanners and records the SelectHints passed to NewMatrixSelector.
type hintsSpy struct {
	inner          engstorage.Scanners
	mu             sync.Mutex
	matrixHintsCap []promstorage.SelectHints
}

func (s *hintsSpy) Close() error { return s.inner.Close() }

func (s *hintsSpy) NewVectorSelector(ctx context.Context, opts *query.Options, hints promstorage.SelectHints, selector logicalplan.VectorSelector) (model.VectorOperator, error) {
	return s.inner.NewVectorSelector(ctx, opts, hints, selector)
}

func (s *hintsSpy) NewMatrixSelector(ctx context.Context, opts *query.Options, hints promstorage.SelectHints, selector logicalplan.MatrixSelector, call logicalplan.FunctionCall) (model.VectorOperator, error) {
	s.mu.Lock()
	s.matrixHintsCap = append(s.matrixHintsCap, hints)
	s.mu.Unlock()
	return s.inner.NewMatrixSelector(ctx, opts, hints, selector, call)
}

func (s *hintsSpy) capturedHints() []promstorage.SelectHints {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.matrixHintsCap
}

func FuzzAnchoredSmoothedModifiers(f *testing.F) {
	f.Add(int64(0), uint32(30), uint32(300), uint32(10), 0.0, 10.0, 5.0, 20.0, uint8(0))

	f.Fuzz(func(t *testing.T, seed int64, startTS, endTS, intervalSeconds uint32, initialVal1, inc1, initialVal2, inc2 float64, funcIdx uint8) {
		if math.IsNaN(initialVal1) || math.IsNaN(initialVal2) || math.IsNaN(inc1) || math.IsNaN(inc2) {
			return
		}
		if math.IsInf(initialVal1, 0) || math.IsInf(initialVal2, 0) || math.IsInf(inc1, 0) || math.IsInf(inc2, 0) {
			return
		}
		if inc1 < 0 || inc2 < 0 || intervalSeconds == 0 || endTS <= startTS {
			return
		}
		// Cap values to avoid overflow.
		if initialVal1 > 1e12 || initialVal2 > 1e12 || inc1 > 1e8 || inc2 > 1e8 {
			return
		}
		// Ensure query range falls within data range (41 samples at 10s = 400s).
		// This avoids edge cases at data boundaries where Prometheus' extendFloats
		// synthesis for changes/resets differs from the thanos engine approach.
		const maxDataTS uint32 = 400
		if endTS > maxDataTS {
			endTS = maxDataTS
		}
		if startTS >= endTS {
			return
		}

		parser.EnableExtendedRangeSelectors = true

		type queryDef struct {
			name  string
			query string
		}

		// Queries covering all supported function+modifier combinations.
		allQueries := []queryDef{
			{"rate_anchored", `rate(http_requests_total[30s] anchored)`},
			{"rate_smoothed", `rate(http_requests_total[30s] smoothed)`},
			{"increase_anchored", `increase(http_requests_total[30s] anchored)`},
			{"increase_smoothed", `increase(http_requests_total[30s] smoothed)`},
			{"delta_anchored", `delta(http_requests_total[30s] anchored)`},
			{"delta_smoothed", `delta(http_requests_total[30s] smoothed)`},
			{"resets_anchored", `resets(http_requests_total[30s] anchored)`},
			{"changes_anchored", `changes(http_requests_total[30s] anchored)`},
			{"rate_anchored_1m", `rate(http_requests_total[1m] anchored)`},
			{"rate_smoothed_1m", `rate(http_requests_total[1m] smoothed)`},
			{"sum_rate_anchored", `sum(rate(http_requests_total[30s] anchored))`},
			{"sum_rate_smoothed", `sum(rate(http_requests_total[30s] smoothed))`},
		}

		// Pick one query per fuzz iteration to keep it fast.
		selected := allQueries[int(funcIdx)%len(allQueries)]

		load := fmt.Sprintf(`load 10s
			http_requests_total{pod="nginx-1"} %.2f+%.2fx40
			http_requests_total{pod="nginx-2"} %.2f+%.2fx40`, initialVal1, inc1, initialVal2, inc2)

		opts := promql.EngineOpts{
			Timeout:              1 * time.Hour,
			MaxSamples:           1e10,
			EnableNegativeOffset: true,
			EnableAtModifier:     true,
		}

		storage := promqltest.LoadedStorage(t, load)
		defer storage.Close()

		start := time.Unix(int64(startTS), 0)
		end := time.Unix(int64(endTS), 0)
		interval := time.Duration(intervalSeconds) * time.Second

		newEngine := engine.New(engine.Opts{
			EngineOpts:                   opts,
			EnableExtendedRangeSelectors: true,
		})
		oldEngine := promql.NewEngine(opts)

		ctx := context.Background()

		q1, err := newEngine.NewRangeQuery(ctx, storage, nil, selected.query, start, end, interval)
		if err != nil {
			return // Skip unsupported queries.
		}
		defer q1.Close()
		newResult := q1.Exec(ctx)

		q2, err := oldEngine.NewRangeQuery(ctx, storage, nil, selected.query, start, end, interval)
		if err != nil {
			t.Fatalf("prometheus engine error for %s: %v", selected.name, err)
		}
		defer q2.Close()
		oldResult := q2.Exec(ctx)

		if !cmp.Equal(oldResult, newResult, comparer) {
			t.Logf("load: %s", load)
			t.Logf("query: %s (%s), start: %d, end: %d, interval: %v", selected.query, selected.name, start.UnixMilli(), end.UnixMilli(), interval)
			t.Errorf("result mismatch.\nnew: %s\nold: %s", newResult.String(), oldResult.String())
		}
	})
}

func TestAnchoredSmoothedSelectHints(t *testing.T) {
	parser.EnableExtendedRangeSelectors = true

	load := `load 10s
		metric 0 10 20 30 40 50 60 70 80 90 100`
	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	lookbackDelta := 5 * time.Minute // engine default

	cases := []struct {
		name          string
		query         string
		expectedRange int64 // hints.Range in ms — should always be original selector range
	}{
		{
			name:          "standard rate — Range equals selector range",
			query:         `rate(metric[30s])`,
			expectedRange: 30000,
		},
		{
			name:          "anchored rate — Range equals selector range, not widened",
			query:         `rate(metric[30s] anchored)`,
			expectedRange: 30000,
		},
		{
			name:          "smoothed rate — Range equals selector range, not widened",
			query:         `rate(metric[30s] smoothed)`,
			expectedRange: 30000,
		},
		{
			name:          "anchored increase — Range equals selector range",
			query:         `increase(metric[1m] anchored)`,
			expectedRange: 60000,
		},
		{
			name:          "smoothed delta — Range equals selector range",
			query:         `delta(metric[1m] smoothed)`,
			expectedRange: 60000,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts := engine.Opts{
				EngineOpts: promql.EngineOpts{
					Timeout:    1 * time.Hour,
					MaxSamples: 1e10,
				},
				EnableExtendedRangeSelectors: true,
			}
			qOpts := &query.Options{
				Start:            time.Unix(0, 0),
				End:              time.Unix(100, 0),
				Step:             10 * time.Second,
				LookbackDelta:    lookbackDelta,
				ExtLookbackDelta: 1 * time.Hour,
			}

			// Parse and build the logical plan to create real scanners.
			expr, err := parser.NewParser(tc.query).ParseExpr()
			testutil.Ok(t, err)

			planOpts := logicalplan.PlanOptions{}
			lplan, err := logicalplan.NewFromAST(expr, qOpts, planOpts)
			testutil.Ok(t, err)
			optimizedPlan, _ := lplan.Optimize(logicalplan.AllOptimizers)

			realScanners, err := promscan.NewPrometheusScanners(storage, qOpts, optimizedPlan)
			testutil.Ok(t, err)
			defer realScanners.Close()

			spy := &hintsSpy{inner: realScanners}

			eng := engine.NewWithScanners(opts, spy)
			q, err := eng.NewRangeQuery(context.Background(), storage, nil, tc.query, time.Unix(0, 0), time.Unix(100, 0), 10*time.Second)
			testutil.Ok(t, err)
			defer q.Close()

			res := q.Exec(context.Background())
			testutil.Ok(t, res.Err)

			hints := spy.capturedHints()
			if len(hints) == 0 {
				t.Fatal("expected at least one NewMatrixSelector call, got none")
			}
			for i, h := range hints {
				if h.Range != tc.expectedRange {
					t.Errorf("hints[%d].Range = %d, want %d (original selector range)", i, h.Range, tc.expectedRange)
				}
			}
		})
	}
}

// TestSmoothedCounterResetAfterRangeEnd pins the right-edge interpolation
// across a counter reset. It follows Prometheus v0.310.0 and later, which model
// the post-reset counter as starting from zero on both edges. v0.308.0 instead
// adds the pre-reset value on the right edge and would return ~103.33 here.
func TestSmoothedCounterResetAfterRangeEnd(t *testing.T) {
	parser.EnableExtendedRangeSelectors = true

	storage := promqltest.LoadedStorage(t, `load 30s
		c_total 100 10 20`)
	defer storage.Close()

	ng := engine.New(engine.Opts{
		EngineOpts:                   promql.EngineOpts{Timeout: time.Hour, MaxSamples: 1e10},
		EnableExtendedRangeSelectors: true,
	})
	ctx := context.Background()

	// The range (10s, 20s] has no sample; both edges interpolate between
	// 100@0s and 10@30s, which is a reset: left = 10*10/30, right = 10*20/30.
	q, err := ng.NewInstantQuery(ctx, storage, nil, `increase(c_total[10s] smoothed)`, time.Unix(20, 0))
	testutil.Ok(t, err)
	defer q.Close()
	res := q.Exec(ctx)
	testutil.Ok(t, res.Err)
	vec, err := res.Vector()
	testutil.Ok(t, err)
	testutil.Equals(t, 1, len(vec))
	testutil.Assert(t, math.Abs(vec[0].F-10.0/3) < 1e-9, "got %v", vec[0].F)
}

func TestAnchoredSmoothedErrorMessages(t *testing.T) {
	parser.EnableExtendedRangeSelectors = true

	storage := promqltest.LoadedStorage(t, `load 10s
		metric 1 2 3 4 5 6
		hist 1 2 {{schema:0 sum:5 count:4 buckets:[1 2 1]}} {{schema:0 sum:6 count:5 buckets:[1 2 2]}}
		late 1 2 3 4 5 {{schema:0 sum:5 count:4 buckets:[1 2 1]}}`)
	defer storage.Close()

	opts := promql.EngineOpts{Timeout: time.Hour, MaxSamples: 1e10}
	ng := engine.New(engine.Opts{EngineOpts: opts, EnableExtendedRangeSelectors: true})
	prom := promql.NewEngine(opts)
	ctx := context.Background()

	for _, query := range []string{
		`absent_over_time(metric[30s] anchored)`,
		`absent_over_time(metric[30s] smoothed)`,
		`sum_over_time(metric[30s] anchored)`,
		`changes(metric[30s] smoothed)`,
		`rate(hist[30s] anchored)`,
		`increase(hist[30s] smoothed)`,
		`changes(hist[30s] anchored)`,
		`hist smoothed`,
		// A float at the evaluation time, then a histogram later in the window.
		`late smoothed`,
	} {
		t.Run(query, func(t *testing.T) {
			q1, err := ng.NewInstantQuery(ctx, storage, nil, query, time.Unix(40, 0))
			testutil.Ok(t, err)
			defer q1.Close()
			q2, err := prom.NewInstantQuery(ctx, storage, nil, query, time.Unix(40, 0))
			testutil.Ok(t, err)
			defer q2.Close()

			want := q2.Exec(ctx).Err
			testutil.NotOk(t, want)
			got := q1.Exec(ctx).Err
			testutil.NotOk(t, got)
			testutil.Equals(t, want.Error(), got.Error())
		})
	}
}

// TestSmoothedOffsetTimestamps pins output timestamps for smoothed selectors
// with an offset. Prometheus v0.308.0 stamps them with the data time rather
// than the evaluation time, which later versions fix, so this is not compared
// against the vendored engine.
func TestSmoothedOffsetTimestamps(t *testing.T) {
	parser.EnableExtendedRangeSelectors = true

	storage := promqltest.LoadedStorage(t, `load 10s
		metric 1 2 stale 4 5`)
	defer storage.Close()

	ng := engine.New(engine.Opts{
		EngineOpts:                   promql.EngineOpts{Timeout: time.Hour, MaxSamples: 1e10},
		EnableExtendedRangeSelectors: true,
	})
	ctx := context.Background()

	q, err := ng.NewRangeQuery(ctx, storage, nil, `metric smoothed offset 5s`, time.Unix(20, 0), time.Unix(30, 0), 5*time.Second)
	testutil.Ok(t, err)
	defer q.Close()
	res := q.Exec(ctx)
	testutil.Ok(t, res.Err)
	m, err := res.Matrix()
	testutil.Ok(t, err)
	testutil.Equals(t, 1, len(m))
	// Data at 15s, 20s and 25s, interpolated across the stale marker at 20s.
	testutil.Equals(t, []promql.FPoint{{T: 20000, F: 2.5}, {T: 25000, F: 3}, {T: 30000, F: 3.5}}, m[0].Floats)
}

// TestExtendedRangeSelectorsConcurrentQueries guards against writing the
// process-global parser flag on every query while other queries parse.
func TestExtendedRangeSelectorsConcurrentQueries(t *testing.T) {
	storage := promqltest.LoadedStorage(t, `load 10s
		http_total 0 10 20 30 40 50`)
	defer storage.Close()

	ng := engine.New(engine.Opts{
		EngineOpts:                   promql.EngineOpts{Timeout: time.Hour, MaxSamples: 1e10},
		EnableExtendedRangeSelectors: true,
	})
	ctx := context.Background()

	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			q, err := ng.NewInstantQuery(ctx, storage, nil, `rate(http_total[30s] anchored)`, time.Unix(50, 0))
			if err != nil {
				t.Error(err)
				return
			}
			defer q.Close()
			if res := q.Exec(ctx); res.Err != nil {
				t.Error(res.Err)
			}
		}()
	}
	wg.Wait()
}
