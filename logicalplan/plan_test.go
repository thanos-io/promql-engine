// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"

	"github.com/thanos-community/promql-engine/parser"
)

var spaces = regexp.MustCompile(`\s+`)
var openParenthesis = regexp.MustCompile(`\(\s+`)
var closedParenthesis = regexp.MustCompile(`\s+\)`)

func TestDefaultOptimizers(t *testing.T) {
	cases := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "common selectors",
			expr:     `sum(metric{a="b", c="d"}) / sum(metric{a="b"})`,
			expected: `sum(filter([c="d"], metric{a="b"})) / sum(metric{a="b"})`,
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

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			plan := New(expr, &Opts{Start: time.Unix(0, 0), End: time.Unix(0, 0)})
			optimizedPlan := plan.Optimize(DefaultOptimizers)
			expectedPlan := strings.Trim(spaces.ReplaceAllString(tcase.expected, " "), " ")
			testutil.Equals(t, expectedPlan, optimizedPlan.Expr().String())
		})
	}
}

func TestMatcherPropagation(t *testing.T) {
	cases := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "common matchers with same metric",
			expr:     `node_filesystem_files{host="$host", mountpoint="/"} - node_filesystem_files`,
			expected: `node_filesystem_files{host="$host",mountpoint="/"} - node_filesystem_files`,
		},
		{
			name:     "common matchers with same overlapping selectors",
			expr:     `node_filesystem_files{host="$host", mountpoint="/"} - node_filesystem_files{host!="$host"}`,
			expected: `node_filesystem_files{host="$host",mountpoint="/"} - node_filesystem_files{host!="$host"}`,
		},
		{
			name:     "common matchers with many-to-one",
			expr:     `node_filesystem_files{host="$host",mountpoint="/"} - on () group_left () node_filesystem_files_free`,
			expected: `node_filesystem_files{host="$host",mountpoint="/"} - on () group_left () node_filesystem_files_free`,
		},
		{
			name:     "common matchers",
			expr:     `node_filesystem_files{host="$host", mountpoint="/"} - node_filesystem_files_free`,
			expected: `node_filesystem_files{host="$host",mountpoint="/"} - node_filesystem_files_free{host="$host",mountpoint="/"}`,
		},
	}

	optimizers := []Optimizer{PropagateMatchersOptimizer{}}
	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			plan := New(expr, &Opts{Start: time.Unix(0, 0), End: time.Unix(0, 0)})
			optimizedPlan := plan.Optimize(optimizers)
			expectedPlan := strings.Trim(spaces.ReplaceAllString(tcase.expected, " "), " ")
			testutil.Equals(t, expectedPlan, optimizedPlan.Expr().String())
		})
	}
}

func cleanUp(replacements map[string]*regexp.Regexp, expr string) string {
	for replacement, match := range replacements {
		expr = match.ReplaceAllString(expr, replacement)
	}
	return strings.Trim(expr, " ")
}
