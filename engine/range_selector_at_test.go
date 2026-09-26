// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/engine"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/promqltest"
)

// TestRangeSelectorWithAtModifier checks range selectors with the @ modifier in functions that are not step invariant,
// like predict_linear. The preprocessing does not wrap these in a step invariant expression, and Prometheus evaluates
// the function at every step, over the window at the timestamp of the selector.
func TestRangeSelectorWithAtModifier(t *testing.T) {
	t.Parallel()

	const load = `load 30s
		a{pod="p1"} 1 5 2 8 3 9 4 12 6 3 7 15 2 9 11 4 6 13 5 8 10
		a{pod="p2"} 9 3 7 2 8 1 6 4 10 2 5 9 3 7 1 8 4 6 2 9 5
		x{pod="p1"} 100+1x20
		x{pod="p2"} 200+1x20`

	queries := []string{
		// Short windows, which a window moving with the step would leave.
		`predict_linear(a[2m] @ 100, 1)`,
		`predict_linear(a[2m] @ 100 offset 30s, 1)`,
		`predict_linear(a[2m] @ 100 offset -30s, 1)`,
		`predict_linear(a[2m] @ start(), 1)`,
		`predict_linear(a[2m] @ end(), 1)`,
		`predict_linear(a[1m] @ 300, 60)`,
		// Long windows.
		`predict_linear(a[1h] @ 100, 1)`,
		`predict_linear(a[1h] @ 450 offset 1m, 3600)`,
		// In binary expressions.
		`x + on(pod) predict_linear(a[2m] @ 100, 1)`,
		`x unless on(pod) predict_linear(a[1m] @ 0, 1)`,
		`predict_linear(a[2m] @ end(), 1) - predict_linear(a[2m] @ start(), 1)`,
		// The inner range selector is evaluated at every step of the subquery.
		`max_over_time(predict_linear(a[2m] @ 100, 1)[5m:1m])`,
		// Step invariant functions are wrapped in a step invariant expression, and evaluated once.
		`rate(a[2m] @ 100)`,
		`quantile_over_time(0.5, a[2m] @ 100)`,
		`x + on(pod) max_over_time(a[2m] @ 100)`,
		// Without the @ modifier.
		`predict_linear(a[2m], 1)`,
		`predict_linear(a[2m] offset 30s, 1)`,
	}

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	opts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}
	ctx := context.Background()
	for _, query := range queries {
		t.Run(query, func(t *testing.T) {
			for _, rng := range []struct {
				start, end, step int64
			}{
				{start: 158, end: 158},
				{start: 0, end: 300, step: 65},
				{start: 90, end: 600, step: 20},
			} {
				newQuery := func(ng promql.QueryEngine) (promql.Query, error) {
					if rng.step == 0 {
						return ng.NewInstantQuery(ctx, storage, nil, query, time.Unix(rng.start, 0))
					}
					return ng.NewRangeQuery(ctx, storage, nil, query, time.Unix(rng.start, 0), time.Unix(rng.end, 0), time.Duration(rng.step)*time.Second)
				}

				q2, err := newQuery(promql.NewEngine(opts))
				testutil.Ok(t, err)
				defer q2.Close()
				oldResult := q2.Exec(ctx)
				testutil.Ok(t, oldResult.Err)

				q1, err := newQuery(engine.New(engine.Opts{EngineOpts: opts}))
				testutil.Ok(t, err)
				defer q1.Close()
				newResult := q1.Exec(ctx)
				testutil.WithGoCmp(comparer).Equals(t, oldResult, newResult, "range %+v", rng)
			}
		})
	}
}
