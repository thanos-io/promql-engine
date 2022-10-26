// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"regexp"
	"strings"
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
			name:     "common selectors with one matching all",
			expr:     `sum(metric{c="d"}) / sum(metric{})`,
			expected: `sum(filter([c="d"], metric)) / sum(metric)`,
		},
		{
			name:     "common selectors",
			expr:     `sum(metric{a="b", c="d"}) / sum(metric{a="b"})`,
			expected: `sum(filter([c="d"], metric{a="b"})) / sum(metric{a="b"})`,
		},
		{
			name:     "common selectors with count",
			expr:     `count(metric{a="b", c="d"}) / count(metric{a="b"})`,
			expected: `count(filter([c="d"], metric{a="b"})) / count(metric{a="b"})`,
		},
		{
			name:     "common selectors with duplicate matchers",
			expr:     `sum(metric{a="b", c="d", a="b"}) / sum(metric{a="b"})`,
			expected: `sum(filter([c="d"], metric{a="b"})) / sum(metric{a="b"})`,
		},
		{
			name:     "common selectors with different operators",
			expr:     `sum(metric{a="b"}) / sum(metric{a=~"b"})`,
			expected: `sum(metric{a="b"}) / sum(metric{a=~"b"})`,
		},
		{
			name:     "common selectors with regex",
			expr:     `http_requests_total / on () group_left sum(http_requests_total{pod=~"p1.+"})`,
			expected: `http_requests_total / on () group_left () sum(filter([pod=~"p1.+"], http_requests_total))`,
		},
		{
			name: "common selectors in different metrics",
			expr: `
	sum(metric_1{a="b", c="d"}) / sum(metric_1{a="b"}) +
	sum(metric_2{a="b", c="d"}) / sum(metric_2{a="b"})
`,
			expected: `
	sum(filter([c="d"], metric_1{a="b"})) / sum(metric_1{a="b"}) +
	sum(filter([c="d"], metric_2{a="b"})) / sum(metric_2{a="b"})`,
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
		{
			name:     "duplicate matchers",
			expr:     `metric_1{a="1", b="2", a="1"} / metric_2{a="1", b="2", a="1"}`,
			expected: `metric_1{a="1",a="1",b="2"} / metric_2{a="1",a="1",b="2"}`,
		},
		{
			name:     "duplicate matchers",
			expr:     `metric_1{a="1", b="2", a="1", e="f"} / metric_1{a="1", b="2", a="1"}`,
			expected: `filter([e="f"], metric_1{a="1",a="1",b="2"}) / metric_1{a="1",a="1",b="2"}`,
		},
	}

	spaces := regexp.MustCompile(`\s+`)
	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			plan := New(expr, time.Unix(0, 0), time.Unix(0, 0))
			optimizedPlan := plan.Optimize(DefaultOptimizers, nil)
			expectedPlan := strings.Trim(spaces.ReplaceAllString(tcase.expected, " "), " ")
			testutil.Equals(t, expectedPlan, optimizedPlan.Expr().String())
		})
	}
}
