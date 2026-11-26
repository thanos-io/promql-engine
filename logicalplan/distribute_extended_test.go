// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/api"
	"github.com/thanos-io/promql-engine/query"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

func TestCalculateStartOffsetForExtendedRangeSelectors(t *testing.T) {
	lookbackDelta := 5 * time.Minute
	cases := []struct {
		expr        string
		startOffset time.Duration
	}{
		{expr: `rate(metric[1h])`, startOffset: time.Hour},
		{expr: `metric`, startOffset: lookbackDelta},
		{expr: `metric anchored`, startOffset: lookbackDelta},
		{expr: `increase(metric[1h] anchored)`, startOffset: time.Hour + lookbackDelta},
	}
	for _, tcase := range cases {
		t.Run(tcase.expr, func(t *testing.T) {
			expr, err := parser.NewParser(parser.Options{}).ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			opts := &query.Options{Start: time.Unix(0, 0), End: time.Unix(0, 0), LookbackDelta: lookbackDelta}
			plan, _ := NewFromAST(expr, opts, PlanOptions{})
			root := plan.Root()
			testutil.Equals(t, tcase.startOffset, calculateStartOffset(&root, lookbackDelta))
		})
	}
}

func TestDistributedExecutionWithExtendedRangeSelectors(t *testing.T) {
	// Two time-partitioned engines for the same partition label set overlapping by 2h.
	firstEngineOpts := engineOpts{
		minTime: time.Unix(0, 0),
		maxTime: time.Unix(0, 0).Add(8 * time.Hour),
	}
	secondEngineOpts := engineOpts{
		minTime: time.Unix(0, 0).Add(6 * time.Hour),
		maxTime: time.Unix(0, 0).Add(12 * time.Hour),
	}
	queryStart := time.Unix(0, 0)
	queryEnd := time.Unix(0, 0).Add(12 * time.Hour)

	cases := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name: "plain range selector is unaffected",
			expr: `rate(metric[1h])`,
			expected: `
dedup(
  remote(rate(metric[1h])),
  remote(rate(metric[1h])) [1970-01-01 07:00:00 +0000 UTC, 1970-01-01 12:00:00 +0000 UTC]
)`,
		},
		{
			name: "anchored range selector start offset includes lookback delta",
			expr: `increase(metric[1h] anchored)`,
			expected: `
dedup(
  remote(increase(metric[1h] anchored)),
  remote(increase(metric[1h] anchored)) [1970-01-01 07:05:00 +0000 UTC, 1970-01-01 12:00:00 +0000 UTC]
)`,
		},
		{
			name:     "smoothed range selector is not distributed",
			expr:     `rate(metric[1h] smoothed)`,
			expected: `rate(metric[1h] smoothed)`,
		},
		{
			name: "plain range selector fitting exactly into the engine overlap is distributed",
			expr: `rate(metric[2h])`,
			expected: `
dedup(
  remote(rate(metric[2h])),
  remote(rate(metric[2h])) [1970-01-01 08:00:00 +0000 UTC, 1970-01-01 12:00:00 +0000 UTC]
)`,
		},
		{
			name: "anchored range selector fitting exactly into the engine overlap is distributed",
			expr: `increase(metric[1h55m] anchored)`,
			expected: `
dedup(
  remote(increase(metric[1h55m] anchored)),
  remote(increase(metric[1h55m] anchored)) [1970-01-01 08:00:00 +0000 UTC, 1970-01-01 12:00:00 +0000 UTC]
)`,
		},
		{
			name:     "anchored range selector exceeding the engine overlap by its lookback is not distributed",
			expr:     `increase(metric[2h] anchored)`,
			expected: `increase(metric[2h] anchored)`,
		},
		{
			name:     "aggregation over a smoothed instant selector is not distributed",
			expr:     `sum by (region) (metric smoothed)`,
			expected: `sum by (region) (metric smoothed)`,
		},
		{
			name:     "smoothed selector in a subquery with negative offset is not distributed",
			expr:     `max_over_time((metric smoothed)[1m:1m] offset -1h)`,
			expected: `max_over_time(metric smoothed[1m:1m] offset -1h)`,
		},
		{
			name: "plain selector with negative offset is unaffected",
			expr: `sum(metric offset -1h55m)`,
			expected: `
sum(dedup(
  remote(sum by (region) (metric offset -1h55m)),
  remote(sum by (region) (metric offset -1h55m)) [1970-01-01 06:05:00 +0000 UTC, 1970-01-01 12:00:00 +0000 UTC]
))`,
		},
		{
			name:     "smoothed selector with negative offset is not distributed",
			expr:     `sum(metric offset -1h55m smoothed)`,
			expected: `sum(metric smoothed offset -1h55m)`,
		},
		{
			name: "plain selector with timestamp close to the end of an engine is distributed",
			expr: `sum(metric @ 28560)`,
			expected: `
sum(dedup(
  remote(sum by (region) (metric @ 28560.000)),
  remote(sum by (region) (metric @ 28560.000)) [1970-01-01 06:05:00 +0000 UTC, 1970-01-01 12:00:00 +0000 UTC]
))`,
		},
		{
			name:     "smoothed selector with timestamp is not distributed",
			expr:     `sum(metric @ 28560 smoothed)`,
			expected: `sum(metric @ 28560.000 smoothed)`,
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			engines := []api.RemoteEngine{
				newEngineMock(firstEngineOpts.mint(), firstEngineOpts.maxt(), []labels.Labels{labels.FromStrings("region", "east")}),
				newEngineMock(secondEngineOpts.mint(), secondEngineOpts.maxt(), []labels.Labels{labels.FromStrings("region", "east")}),
			}
			optimizers := []Optimizer{
				DistributedExecutionOptimizer{Endpoints: api.NewStaticEndpoints(engines)},
			}

			expr, err := parser.NewParser(parser.Options{}).ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			opts := &query.Options{Start: queryStart, End: queryEnd, Step: time.Minute, LookbackDelta: 5 * time.Minute}
			plan, _ := NewFromAST(expr, opts, PlanOptions{})
			optimizedPlan, _ := plan.Optimize(optimizers)
			testutil.Equals(t, cleanUp(replacements, tcase.expected), optimizedPlan.Root().String())
		})
	}
}
