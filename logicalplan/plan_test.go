// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/promql/parser"
)

func TestDefaultOptimizers(t *testing.T) {
	cases := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "common selectors",
			expr:     `sum(metric{a="b", c="d"}) / sum(metric{a="b"})`,
			expected: `sum(filter([a="b" c="d"], metric{a="b"})) / sum(metric{a="b"})`,
		},
		{
			name:     "different selectors",
			expr:     `sum(metric{a="b"}) / sum(metric{c="d"})`,
			expected: `sum(metric{a="b"}) / sum(metric{c="d"})`,
		},
		{
			name:     "different operator",
			expr:     `sum(metric{a="b"}) / sum(metric{a=~"b"})`,
			expected: `sum(metric{a="b"}) / sum(metric{a=~"b"})`,
		},
		{
			name:     "different metrics",
			expr:     `sum(metric_1{a="b"}) / sum(metric_2{a="b"})`,
			expected: `sum(metric_1{a="b"}) / sum(metric_2{a="b"})`,
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			plan := New(expr, time.Unix(0, 0), time.Unix(0, 0))
			optimizedPlan := plan.RunOptimizers(DefaultOptimizers)
			testutil.Equals(t, tcase.expected, optimizedPlan.String())
		})
	}
}
