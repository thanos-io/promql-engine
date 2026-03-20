// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/promql/promqltest"

	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/logicalplan"
)

func TestAnchoredSmoothedModifiers(t *testing.T) {
	t.Parallel()
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
	}

	start := time.Unix(0, 0)
	end := time.Unix(100, 0)
	step := 10 * time.Second

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
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
	t.Parallel()
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
			_, err := newEngine.NewInstantQuery(ctx, storage, nil, query, time.Unix(50, 0))
			testutil.NotOk(t, err)
		})
	}

	unsupportedSmoothed := []string{
		`avg_over_time(metric[30s] smoothed)`,
		`resets(metric[30s] smoothed)`,
		`changes(metric[30s] smoothed)`,
	}
	for _, query := range unsupportedSmoothed {
		t.Run("unsupported_smoothed/"+query, func(t *testing.T) {
			_, err := newEngine.NewInstantQuery(ctx, storage, nil, query, time.Unix(50, 0))
			testutil.NotOk(t, err)
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
