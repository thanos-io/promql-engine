// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/internal/teststorage"
	"github.com/thanos-io/promql-engine/logicalplan"
	"github.com/thanos-io/promql-engine/query"
	engstorage "github.com/thanos-io/promql-engine/storage"
	promscan "github.com/thanos-io/promql-engine/storage/prometheus"

	"github.com/efficientgo/core/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/promql/promqltest"
	promstorage "github.com/prometheus/prometheus/storage"
)

// TestExtendedVectorsAcceptance runs Prometheus' extended_vectors.test, which
// TestPromqlAcceptance skips, adapted to this engine's limitations.
func TestExtendedVectorsAcceptance(t *testing.T) {
	content, err := os.ReadFile("testdata/extended_vectors.test")
	testutil.Ok(t, err)

	ng := engine.New(engine.Opts{
		EnableExperimentalFunctions: true,
		Logger:                      promslog.NewNopLogger(),
		MaxSamples:                  5e10,
		Timeout:                     time.Hour,
	})
	promqltest.RunTestWithStorage(t, string(content), ng, func(t testing.TB) promstorage.Storage {
		return teststorage.New(t)
	})
}

const (
	counterHistograms = `{{schema:0 sum:5 count:4 buckets:[1 2 1]}} {{schema:0 sum:6 count:5 buckets:[1 2 2]}} {{schema:0 sum:9 count:7 buckets:[2 3 2]}} {{schema:0 sum:2 count:1 buckets:[1]}} {{schema:0 sum:4 count:3 buckets:[1 1 1]}} {{schema:0 sum:8 count:6 buckets:[2 2 2]}}`
	gaugeHistograms   = `{{schema:0 sum:5 count:4 buckets:[1 2 1] counter_reset_hint:gauge}} {{schema:0 sum:3 count:2 buckets:[1 1] counter_reset_hint:gauge}} {{schema:0 sum:7 count:6 buckets:[2 2 2] counter_reset_hint:gauge}}`
)

// TestAnchoredSmoothedModifiers compares anchored and smoothed selectors with
// Prometheus, as range queries and as instant queries at every 5s, with the
// default lookback delta and with one short enough to expire samples.
func TestAnchoredSmoothedModifiers(t *testing.T) {
	cases := []struct {
		name  string
		load  string
		query string
		// start, end and step of the range query default to 0s, 100s and 10s.
		start, end time.Time
		step       time.Duration
	}{
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
		{
			// Smoothed interpolation must handle counter resets to avoid
			// negative increases.
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
		{
			// Both edges interpolate across the reset between 100@0s and 10@30s.
			name: "smoothed increase with reset after the range end",
			load: `load 30s
			    c_total 100 10 20`,
			query: `increase(c_total[10s] smoothed)`,
		},
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
			query: `sum(increase(http_total[30s] smoothed))`,
		},
		// Windows without samples on one side.
		{
			name: "smoothed rate before the first sample",
			load: `load 10s
			    metric _ _ _ _ 10 20 30 40`,
			query: `rate(metric[5s] smoothed)`,
		},
		{
			name: "anchored changes and resets after the last sample",
			load: `load 10s
			    metric 5 1 3 3 2`,
			query: `changes(metric[20s] anchored) + resets(metric[20s] anchored)`,
		},
		{
			name: "anchored rate with gaps",
			load: `load 10s
			    metric 0 _ _ _ 40 _ 60 _ _ _ _ 110`,
			query: `rate(metric[15s] anchored)`,
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
		{
			name: "timestamp of smoothed instant vector",
			load: `load 10s
			    metric 1 2 3 _ 5`,
			query: `timestamp(metric smoothed)`,
			step:  5 * time.Second,
		},
		// Offsets, @ and subqueries.
		{
			name: "smoothed rate with offset",
			load: `load 10s
			    metric 0 10 20 30 40 50 60 70 80 90 100`,
			query: `rate(metric[30s] smoothed offset 15s)`,
		},
		{
			name: "smoothed rate with negative offset",
			load: `load 10s
			    metric 0 10 20 30 40 50 60 70 80 90 100`,
			query: `rate(metric[30s] smoothed offset -15s)`,
		},
		{
			name: "anchored increase with negative offset",
			load: `load 10s
			    metric 0 10 20 5 15 25 35 45 55 65 75`,
			query: `increase(metric[30s] anchored offset -25s)`,
		},
		{
			name: "smoothed instant vector with offset",
			load: `load 10s
			    metric 1 2 stale 4 5`,
			query: `metric smoothed offset 5s`,
			start: time.Unix(20, 0),
			end:   time.Unix(30, 0),
			step:  5 * time.Second,
		},
		{
			name: "smoothed instant vector with negative offset",
			load: `load 10s
			    metric 1 2 3 4 5 6 7 8 9 10 11`,
			query: `metric smoothed offset -7s`,
			step:  5 * time.Second,
		},
		{
			name: "anchored rate at fixed time",
			load: `load 10s
			    metric 0 10 20 5 15 25 35 45 55 65 75`,
			query: `rate(metric[30s] anchored @ 45)`,
		},
		{
			name: "smoothed increase at fixed time",
			load: `load 10s
			    metric 0 10 20 5 15 25 35 45 55 65 75`,
			query: `increase(metric[30s] smoothed @ 45)`,
		},
		{
			name: "smoothed instant vector at fixed time",
			load: `load 10s
			    metric 1 2 3 4 5 6 7 8 9 10 11`,
			query: `metric smoothed @ 33`,
		},
		{
			name: "anchored rate in subquery",
			load: `load 10s
			    metric 0 10 20 5 15 25 35 45 55 65 75`,
			query: `max_over_time(rate(metric[30s] anchored)[40s:5s])`,
		},
		{
			name: "smoothed instant vector in subquery",
			load: `load 10s
			    metric 1 2 stale 4 5 6 7 8 9 10 11`,
			query: `avg_over_time((metric smoothed)[30s:5s])`,
		},
		// Native histograms.
		{
			name: "anchored rate on counter histograms",
			load: `load 10s
			    metric ` + counterHistograms,
			query: `rate(metric[30s] anchored)`,
		},
		{
			name: "smoothed increase on counter histograms",
			load: `load 10s
			    metric ` + counterHistograms,
			query: `increase(metric[30s] smoothed)`,
		},
		{
			name: "smoothed rate on counter histograms with short range",
			load: `load 10s
			    metric ` + counterHistograms,
			query: `rate(metric[5s] smoothed)`,
		},
		{
			name: "anchored delta on counter histograms",
			load: `load 10s
			    metric ` + counterHistograms,
			query: `delta(metric[30s] anchored)`,
		},
		{
			name: "anchored delta on gauge histograms",
			load: `load 10s
			    metric ` + gaugeHistograms,
			query: `delta(metric[20s] anchored)`,
		},
		{
			name: "smoothed delta on gauge histograms",
			load: `load 10s
			    metric ` + gaugeHistograms,
			query: `delta(metric[20s] smoothed)`,
		},
		{
			name: "smoothed rate on gauge histograms",
			load: `load 10s
			    metric ` + gaugeHistograms,
			query: `rate(metric[20s] smoothed)`,
		},
		{
			name: "anchored changes and resets on histograms",
			load: `load 10s
			    metric ` + counterHistograms,
			query: `changes(metric[30s] anchored) + resets(metric[30s] anchored)`,
		},
		{
			name: "smoothed instant vector on counter histograms",
			load: `load 10s
			    metric ` + counterHistograms,
			query: `metric smoothed`,
			step:  5 * time.Second,
		},
		{
			name: "smoothed instant vector on gauge histograms",
			load: `load 10s
			    metric ` + gaugeHistograms,
			query: `metric smoothed`,
			step:  5 * time.Second,
		},
		// Windows that mix floats and histograms.
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
			// The histogram precedes the float that anchors the range, so it is
			// only seen by the extended window.
			name: "anchored rate with histogram before the anchor",
			load: `load 10s
			    metric {{schema:0 sum:5 count:4 buckets:[1 2 1]}} 1 2 3 4 5 6 7 8 9 10`,
			query: `rate(metric[10s] anchored)`,
		},
		{
			name: "anchored changes over float then histogram",
			load: `load 10s
			    metric 1 2 {{schema:0 sum:5 count:4 buckets:[1 2 1]}} {{schema:0 sum:5 count:4 buckets:[1 2 1]}} 3`,
			query: `changes(metric[30s] anchored)`,
		},
		{
			name: "smoothed instant vector over histogram then floats",
			load: `load 10s
			    metric {{schema:0 sum:5 count:4 buckets:[1 2 1]}} 2 3`,
			query: `metric smoothed`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			storage := teststorage.LoadedStorage(t, tc.load)
			defer storage.Close()

			start, end, step := time.Unix(0, 0), time.Unix(100, 0), 10*time.Second
			if !tc.start.IsZero() {
				start = tc.start
			}
			if !tc.end.IsZero() {
				end = tc.end
			}
			if tc.step != 0 {
				step = tc.step
			}

			for _, lookbackDelta := range []time.Duration{0, 25 * time.Second} {
				for _, disableOptimizers := range []bool{false, true} {
					t.Run(fmt.Sprintf("lookbackDelta=%v/disableOptimizers=%v", lookbackDelta, disableOptimizers), func(t *testing.T) {
						optimizers := logicalplan.AllOptimizers
						if disableOptimizers {
							optimizers = logicalplan.NoOptimizers
						}
						opts := engine.Opts{
							Timeout:           time.Hour,
							MaxSamples:        1e10,
							LookbackDelta:     lookbackDelta,
							LogicalOptimizers: optimizers,
							SelectorBatchSize: 1,
						}
						newEngine := engine.New(opts)
						oldEngine := promql.NewEngine(prometheusEngineOpts(opts))
						ctx := context.Background()

						q1, err := newEngine.NewRangeQuery(ctx, storage, nil, tc.query, start, end, step)
						testutil.Ok(t, err)
						defer q1.Close()
						q2, err := oldEngine.NewRangeQuery(ctx, storage, nil, tc.query, start, end, step)
						testutil.Ok(t, err)
						defer q2.Close()
						want, got := q2.Exec(ctx), q1.Exec(ctx)
						// comparer discards PromQL annotations, so check them first.
						assertSameAnnotationsAndHints(t, want, got)
						testutil.WithGoCmp(comparer).Equals(t, want, got, queryExplanation(q1))

						for ts := start; !ts.After(end); ts = ts.Add(5 * time.Second) {
							q1, err := newEngine.NewInstantQuery(ctx, storage, nil, tc.query, ts)
							testutil.Ok(t, err)
							q2, err := oldEngine.NewInstantQuery(ctx, storage, nil, tc.query, ts)
							testutil.Ok(t, err)
							want, got := q2.Exec(ctx), q1.Exec(ctx)
							assertSameAnnotationsAndHints(t, want, got)
							testutil.WithGoCmp(comparer).Equals(t, want, got, "instant query at %v\n%s", ts.Unix(), queryExplanation(q1))
							q1.Close()
							q2.Close()
						}
					})
				}
			}
		})
	}
}

// assertSameAnnotationsAndHints checks what comparer ignores: PromQL
// annotations, without their source position which this engine does not
// propagate, and the counter reset hints of histogram results. comparer
// deletes PromQL annotations from both results, so call this before it.
func assertSameAnnotationsAndHints(t *testing.T, want, got *promql.Result) {
	t.Helper()

	annos := func(r *promql.Result) []string {
		var out []string
		for _, err := range r.Warnings {
			out = append(out, annotationPosition.ReplaceAllString(err.Error(), ""))
		}
		slices.Sort(out)
		return out
	}
	testutil.Equals(t, annos(want), annos(got), "annotations")

	hints := func(r *promql.Result) []string {
		var out []string
		switch v := r.Value.(type) {
		case promql.Vector:
			for _, s := range v {
				if s.H != nil {
					out = append(out, fmt.Sprintf("%s@%d %v", s.Metric, s.T, s.H.CounterResetHint))
				}
			}
		case promql.Matrix:
			for _, s := range v {
				for _, p := range s.Histograms {
					out = append(out, fmt.Sprintf("%s@%d %v", s.Metric, p.T, p.H.CounterResetHint))
				}
			}
		}
		slices.Sort(out)
		return out
	}
	testutil.Equals(t, hints(want), hints(got), "counter reset hints")
}

var annotationPosition = regexp.MustCompile(` \(\d+:\d+\)$`)

// TestAnchoredSmoothedValidation checks that modifiers are only accepted with
// the functions Prometheus supports, with the same error at evaluation time.
func TestAnchoredSmoothedValidation(t *testing.T) {
	storage := teststorage.LoadedStorage(t, `load 10s
		metric 0 10 20 30 40 50`)
	defer storage.Close()

	opts := engine.Opts{Timeout: time.Hour, MaxSamples: 1e10}
	ng := engine.New(opts)
	prom := promql.NewEngine(prometheusEngineOpts(opts))
	ctx := context.Background()

	for _, tc := range []struct {
		query string
		fails bool
	}{
		{query: `rate(metric[30s] anchored)`},
		{query: `increase(metric[30s] anchored)`},
		{query: `delta(metric[30s] anchored)`},
		{query: `resets(metric[30s] anchored)`},
		{query: `changes(metric[30s] anchored)`},
		{query: `rate(metric[30s] smoothed)`},
		{query: `increase(metric[30s] smoothed)`},
		{query: `delta(metric[30s] smoothed)`},
		{query: `avg_over_time(metric[30s] anchored)`, fails: true},
		{query: `sum_over_time(metric[30s] anchored)`, fails: true},
		{query: `max_over_time(metric[30s] anchored)`, fails: true},
		{query: `deriv(metric[30s] anchored)`, fails: true},
		{query: `absent_over_time(metric[30s] anchored)`, fails: true},
		{query: `predict_linear(metric[30s] anchored, 60)`, fails: true},
		{query: `quantile_over_time(0.5, metric[30s] anchored)`, fails: true},
		{query: `avg_over_time(metric[30s] smoothed)`, fails: true},
		{query: `resets(metric[30s] smoothed)`, fails: true},
		{query: `changes(metric[30s] smoothed)`, fails: true},
		{query: `absent_over_time(metric[30s] smoothed)`, fails: true},
		{query: `predict_linear(metric[30s] smoothed, 60)`, fails: true},
		{query: `quantile_over_time(0.5, metric[30s] smoothed)`, fails: true},
		// The error is raised even when no series is selected.
		{query: `sum_over_time(missing[30s] anchored)`, fails: true},
	} {
		t.Run(tc.query, func(t *testing.T) {
			q1, err := ng.NewInstantQuery(ctx, storage, nil, tc.query, time.Unix(50, 0))
			testutil.Ok(t, err)
			defer q1.Close()
			q2, err := prom.NewInstantQuery(ctx, storage, nil, tc.query, time.Unix(50, 0))
			testutil.Ok(t, err)
			defer q2.Close()

			want, got := q2.Exec(ctx), q1.Exec(ctx)
			if !tc.fails {
				testutil.Ok(t, want.Err)
				testutil.Ok(t, got.Err)
				return
			}
			testutil.NotOk(t, want.Err)
			testutil.NotOk(t, got.Err)
			testutil.Equals(t, want.Err.Error(), got.Err.Error())
		})
	}
}

// hintsSpy wraps real scanners and records the SelectHints passed to them.
type hintsSpy struct {
	inner engstorage.Scanners

	mu          sync.Mutex
	vectorHints []promstorage.SelectHints
	matrixHints []promstorage.SelectHints
}

func (s *hintsSpy) Close() error { return s.inner.Close() }

func (s *hintsSpy) NewVectorSelector(ctx context.Context, opts *query.Options, hints promstorage.SelectHints, selector logicalplan.VectorSelector) (model.VectorOperator, error) {
	s.mu.Lock()
	s.vectorHints = append(s.vectorHints, hints)
	s.mu.Unlock()
	return s.inner.NewVectorSelector(ctx, opts, hints, selector)
}

func (s *hintsSpy) NewMatrixSelector(ctx context.Context, opts *query.Options, hints promstorage.SelectHints, selector logicalplan.MatrixSelector, call logicalplan.FunctionCall) (model.VectorOperator, error) {
	s.mu.Lock()
	s.matrixHints = append(s.matrixHints, hints)
	s.mu.Unlock()
	return s.inner.NewMatrixSelector(ctx, opts, hints, selector, call)
}

// TestAnchoredSmoothedSelectHints checks the time range requested from storage.
// Modifiers widen Start and End by the lookback delta but keep Range at the
// selector range, which backends use to pick a resolution.
func TestAnchoredSmoothedSelectHints(t *testing.T) {
	storage := teststorage.LoadedStorage(t, `load 10s
		metric 0 10 20 30 40 50 60 70 80 90 100`)
	defer storage.Close()

	const (
		start    = int64(0)
		end      = int64(100_000)
		lookback = int64(5 * time.Minute / time.Millisecond) // engine default
	)
	cases := []struct {
		name   string
		query  string
		matrix bool
		// The selector excludes samples at exactly Start, so Start is 1ms
		// after the oldest timestamp it reads.
		wantStart, wantEnd, wantRange int64
	}{
		{name: "standard rate", query: `rate(metric[30s])`, matrix: true, wantStart: start - 30_000 + 1, wantEnd: end, wantRange: 30_000},
		{name: "anchored rate", query: `rate(metric[30s] anchored)`, matrix: true, wantStart: start - 30_000 - lookback + 1, wantEnd: end, wantRange: 30_000},
		{name: "smoothed rate", query: `rate(metric[30s] smoothed)`, matrix: true, wantStart: start - 30_000 - lookback + 1, wantEnd: end + lookback, wantRange: 30_000},
		{name: "smoothed delta with offset", query: `delta(metric[1m] smoothed offset 10s)`, matrix: true, wantStart: start - 60_000 - lookback + 1 - 10_000, wantEnd: end + lookback - 10_000, wantRange: 60_000},
		{name: "instant vector", query: `metric`, wantStart: start - lookback + 1, wantEnd: end},
		{name: "smoothed instant vector", query: `metric smoothed`, wantStart: start - lookback + 1, wantEnd: end + lookback},
		{name: "timestamp of smoothed instant vector", query: `timestamp(metric smoothed)`, wantStart: start - lookback + 1, wantEnd: end},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			qOpts := &query.Options{
				Start:            time.UnixMilli(start),
				End:              time.UnixMilli(end),
				Step:             10 * time.Second,
				LookbackDelta:    time.Duration(lookback) * time.Millisecond,
				ExtLookbackDelta: time.Hour,
			}
			expr, err := parser.NewParser(parser.Options{}).ParseExpr(tc.query)
			testutil.Ok(t, err)
			lplan, err := logicalplan.NewFromAST(expr, qOpts, logicalplan.PlanOptions{})
			testutil.Ok(t, err)
			optimizedPlan, _ := lplan.Optimize(logicalplan.AllOptimizers)
			realScanners, err := promscan.NewPrometheusScanners(storage, qOpts, optimizedPlan)
			testutil.Ok(t, err)
			spy := &hintsSpy{inner: realScanners}
			defer spy.Close()

			eng := engine.NewWithScanners(engine.Opts{Timeout: time.Hour, MaxSamples: 1e10}, spy)
			q, err := eng.NewRangeQuery(context.Background(), storage, nil, tc.query, qOpts.Start, qOpts.End, qOpts.Step)
			testutil.Ok(t, err)
			defer q.Close()
			testutil.Ok(t, q.Exec(context.Background()).Err)

			hints := spy.vectorHints
			if tc.matrix {
				hints = spy.matrixHints
			}
			testutil.Assert(t, len(hints) > 0, "no selector was created")
			for _, h := range hints {
				testutil.Equals(t, tc.wantStart, h.Start, "Start")
				testutil.Equals(t, tc.wantEnd, h.End, "End")
				testutil.Equals(t, tc.wantRange, h.Range, "Range")
			}
		})
	}
}

func FuzzAnchoredSmoothedModifiers(f *testing.F) {
	f.Add(int64(0), uint32(30), uint32(300), uint32(10), uint8(0))
	f.Add(int64(1), uint32(0), uint32(600), uint32(7), uint8(5))
	f.Add(int64(2), uint32(250), uint32(500), uint32(3), uint8(7))

	queries := []string{
		`rate(http_requests_total[30s] anchored)`,
		`rate(http_requests_total[30s] smoothed)`,
		`increase(http_requests_total[30s] anchored)`,
		`increase(http_requests_total[30s] smoothed)`,
		`delta(http_requests_total[30s] anchored)`,
		`delta(http_requests_total[30s] smoothed)`,
		`resets(http_requests_total[30s] anchored)`,
		`changes(http_requests_total[30s] anchored)`,
		`rate(http_requests_total[1m] anchored)`,
		`rate(http_requests_total[1m] smoothed)`,
		`sum(rate(http_requests_total[30s] anchored))`,
		`sum(rate(http_requests_total[30s] smoothed))`,
		`http_requests_total smoothed`,
		`increase(http_requests_total[45s] smoothed offset -20s)`,
	}

	f.Fuzz(func(t *testing.T, seed int64, startTS, endTS, intervalSeconds uint32, queryIdx uint8) {
		if intervalSeconds == 0 || endTS <= startTS || endTS > 3600 {
			return
		}
		query := queries[int(queryIdx)%len(queries)]
		load := randomCounterLoad(rand.New(rand.NewSource(seed)))

		storage := teststorage.LoadedStorage(t, load)
		defer storage.Close()

		opts := engine.Opts{
			Timeout:       time.Hour,
			MaxSamples:    1e10,
			LookbackDelta: time.Duration(15+uint64(seed)%4*15) * time.Second,
		}
		newEngine := engine.New(opts)
		oldEngine := promql.NewEngine(prometheusEngineOpts(opts))
		start, end := time.Unix(int64(startTS), 0), time.Unix(int64(endTS), 0)
		interval := time.Duration(intervalSeconds) * time.Second
		ctx := context.Background()

		q1, err := newEngine.NewRangeQuery(ctx, storage, nil, query, start, end, interval)
		testutil.Ok(t, err)
		defer q1.Close()
		q2, err := oldEngine.NewRangeQuery(ctx, storage, nil, query, start, end, interval)
		testutil.Ok(t, err)
		defer q2.Close()

		newResult, oldResult := q1.Exec(ctx), q2.Exec(ctx)
		if !cmp.Equal(oldResult, newResult, comparer) {
			t.Logf("load: %s", load)
			t.Logf("query: %s, start: %d, end: %d, interval: %v, lookback: %v", query, start.UnixMilli(), end.UnixMilli(), interval, opts.LookbackDelta)
			t.Errorf("result mismatch.\nnew: %s\nold: %s", newResult.String(), oldResult.String())
		}
	})
}

// randomCounterLoad returns a load command for two counters sampled every 10s
// with random delayed starts, gaps, stale markers and counter resets.
func randomCounterLoad(r *rand.Rand) string {
	var b strings.Builder
	b.WriteString("load 10s\n")
	for _, pod := range []string{"nginx-1", "nginx-2"} {
		fmt.Fprintf(&b, "  http_requests_total{pod=%q}", pod)
		value := math.Round(r.Float64() * 100)
		delay := r.Intn(6)
		for i := range 40 {
			switch p := r.Float64(); {
			case i < delay:
				b.WriteString(" _")
				continue
			case p < 0.1:
				b.WriteString(" _")
				continue
			case p < 0.15:
				b.WriteString(" stale")
				continue
			case p < 0.25:
				value = 0
			}
			value += math.Round(r.Float64() * 20)
			fmt.Fprintf(&b, " %v", value)
		}
		b.WriteString("\n")
	}
	return b.String()
}

// TestSmoothedSampleLimit checks that samples buffered for interpolation count
// towards the sample limit: peaks within a step or a batch are rejected, while
// buffers of series that are not held at the same time do not add up.
//
// The iterator only retains the newest sample before the evaluation time and
// the samples after it, while Prometheus buffers the whole window, so the
// limits are asserted for this engine and results are compared with Prometheus
// under a generous limit.
func TestSmoothedSampleLimit(t *testing.T) {
	wideHistogram := fmt.Sprintf("{{schema:0 sum:200 count:200 buckets:[%s]}}", strings.TrimSpace(strings.Repeat("1 ", 200)))
	cases := []struct {
		name       string
		load       string
		query      string // defaults to `metric smoothed`
		start, end time.Time
		step       time.Duration
		opts       engine.Opts
		// The query fails with maxSamples set to below and succeeds with
		// above, the most samples it holds at once.
		below, above int
	}{
		{
			// The newest sample before 300.5s, the 300 samples after it within
			// the lookback delta, and the output.
			name: "instant query",
			load: `load 1s
				metric 0+1x600`,
			start: time.Unix(300, 500_000_000),
			below: 301,
			above: 302,
		},
		{
			// At 0s, 301 samples ahead, the next one read beyond the window and
			// the output. They have drained by the last step of the batch.
			name: "buffer drained before the end of the batch",
			load: `load 1s
				metric 0+1x600`,
			start: time.Unix(0, 0),
			end:   time.Unix(1200, 0),
			step:  600 * time.Second,
			below: 302,
			above: 303,
		},
		{
			// At 300s, the previous sample, 301 samples ahead and the outputs of
			// the six steps of the batch so far.
			name: "more than one batch",
			load: `load 1s
				metric 0+1x600`,
			start: time.Unix(0, 0),
			end:   time.Unix(1200, 0),
			step:  60 * time.Second,
			below: 307,
			above: 308,
		},
		{
			// The returned copy of the histogram coexists with the buffered one.
			name: "histogram output and its buffer",
			load: `load 1s
				metric ` + wideHistogram,
			query: `count(metric smoothed)`,
			start: time.Unix(0, 0),
			below: 221,
			above: 222,
		},
		{
			// The histogram is still held when the float that replaces it is
			// read within the same step.
			name: "histogram replaced within a step",
			load: `load 1s
				metric ` + wideHistogram + ` 1`,
			start: time.Unix(2, 0),
			below: 111,
			above: 112,
		},
		{
			// The second series' peak and output, plus the first series'
			// outputs; the first series' buffer was released before.
			name: "buffers of series not held at the same time",
			load: `load 1s
				metric{id="a"} 0+1x600
				metric{id="b"} 0+1x600`,
			start: time.Unix(0, 0),
			end:   time.Unix(1200, 0),
			step:  600 * time.Second,
			opts:  engine.Opts{DecodingConcurrency: 1, SelectorBatchSize: 2},
			below: 304,
			above: 305,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			storage := teststorage.LoadedStorage(t, tc.load)
			defer storage.Close()

			query := tc.query
			if query == "" {
				query = `metric smoothed`
			}
			exec := func(ng promql.QueryEngine) *promql.Result {
				ctx := context.Background()
				var q promql.Query
				var err error
				if tc.end.IsZero() {
					q, err = ng.NewInstantQuery(ctx, storage, nil, query, tc.start)
				} else {
					q, err = ng.NewRangeQuery(ctx, storage, nil, query, tc.start, tc.end, tc.step)
				}
				testutil.Ok(t, err)
				defer q.Close()
				return q.Exec(ctx)
			}
			withLimit := func(maxSamples int) engine.Opts {
				opts := tc.opts
				opts.Timeout, opts.MaxSamples = time.Hour, maxSamples
				return opts
			}

			err := exec(engine.New(withLimit(tc.below))).Err
			testutil.NotOk(t, err)
			testutil.Assert(t, strings.Contains(err.Error(), "would load too many samples"), "unexpected error: %v", err)
			got := exec(engine.New(withLimit(tc.above)))
			testutil.Ok(t, got.Err)

			want := exec(promql.NewEngine(prometheusEngineOpts(withLimit(1e6))))
			testutil.Ok(t, want.Err)
			testutil.WithGoCmp(comparer).Equals(t, want, got)
		})
	}
}
