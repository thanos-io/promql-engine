// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/query"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/util/annotations"
)

type countingOptimizer struct{ calls int }

func (o *countingOptimizer) Optimize(plan Node, _ *query.Options) (Node, annotations.Annotations) {
	o.calls++
	return plan, nil
}

// TestOptimizeSkipsInvalidRangeModifiers checks that a query whose evaluation
// fails on an unsupported modifier is not optimized, so that no optimizer can
// prune the selector that raises the error.
func TestOptimizeSkipsInvalidRangeModifiers(t *testing.T) {
	for _, tc := range []struct {
		expr     string
		optimize bool
	}{
		{expr: `rate(metric[5m] anchored)`, optimize: true},
		{expr: `changes(metric[5m] anchored)`, optimize: true},
		{expr: `metric smoothed`, optimize: true},
		{expr: `sum_over_time(metric[5m] anchored)`},
		{expr: `changes(metric[5m] smoothed)`},
		{expr: `quantile_over_time(0.5, metric[5m] smoothed)`},
		{expr: `sum(rate(a[1m])) + absent_over_time(b[5m] anchored)`},
	} {
		t.Run(tc.expr, func(t *testing.T) {
			expr, err := parser.NewParser(parser.Options{}).ParseExpr(tc.expr)
			testutil.Ok(t, err)
			plan, err := NewFromAST(expr, &query.Options{Start: time.Unix(0, 0), End: time.Unix(0, 0)}, PlanOptions{})
			testutil.Ok(t, err)

			o := &countingOptimizer{}
			plan.Optimize([]Optimizer{o})
			testutil.Equals(t, tc.optimize, o.calls > 0)
		})
	}
}
