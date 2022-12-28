// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"github.com/thanos-community/promql-engine/api"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/promql/parser"
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

			plan := New(expr, time.Unix(0, 0), time.Unix(0, 0))
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

			plan := New(expr, time.Unix(0, 0), time.Unix(0, 0))
			optimizedPlan := plan.Optimize(optimizers)
			expectedPlan := strings.Trim(spaces.ReplaceAllString(tcase.expected, " "), " ")
			testutil.Equals(t, expectedPlan, optimizedPlan.Expr().String())
		})
	}
}

func TestDistributedExecution(t *testing.T) {
	cases := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name: "sum-rate",
			expr: `sum by (pod) (rate(http_requests_total[5m]))`,
			expected: `
sum by (pod) (
  coalesce(
    remote(sum by (pod) (rate(http_requests_total[5m]))),
    remote(sum by (pod) (rate(http_requests_total[5m])))
  )
)`,
		},
		{
			name: "avg",
			expr: `avg by (pod) (http_requests_total)`,
			expected: `
avg by (pod) (
  coalesce(
    remote(http_requests_total),
    remote(http_requests_total)
  )
)`,
		},
		{
			name: "two-level aggregation",
			expr: `max by (pod) (sum by (pod) (http_requests_total))`,
			expected: `
max by (pod) (
  sum by (pod) ( 
    coalesce(
      remote(sum by (pod) (http_requests_total)),
      remote(sum by (pod) (http_requests_total))
    )
  )
)`,
		},
		{
			name: "aggregation of binary expression",
			expr: `max by (pod) (metric_a / metric_b)`,
			expected: `
max by (pod) (
  coalesce(remote(metric_a), remote(metric_a)) 
  / 
  coalesce(remote(metric_b), remote(metric_b))
)
`,
		},
		{
			name: "unsupported aggregation in the operand path",
			expr: `max by (pod) (sort(avg(http_requests_total)))`,
			expected: `
max by (pod) (sort(avg(
  coalesce(
    remote(http_requests_total),
    remote(http_requests_total)
  )
)))`,
		},
		{
			name: "binary operation in the operand path",
			expr: `max by (pod) (sort(metric_a / metric_b))`,
			expected: `
max by (pod) (sort(
  coalesce(remote(metric_a), remote(metric_a)) 
  / 
  coalesce(remote(metric_b), remote(metric_b))
))`,
		},
		{
			name: "binary operation with aggregations",
			expr: `sum by (pod) (metric_a) / sum by (pod) (metric_b)`,
			expected: `
sum by (pod) (coalesce(
  remote(sum by (pod) (metric_a)), 
  remote(sum by (pod) (metric_a)))
) 
/ 
sum by (pod) (coalesce(
  remote(sum by (pod) (metric_b)), 
  remote(sum by (pod) (metric_b)))
)`,
		},
		{
			name: "function sharding",
			expr: `rate(http_requests_total[2m])`,
			expected: `
coalesce(
  remote(rate(http_requests_total[2m])), 
  remote(rate(http_requests_total[2m])))`,
		},
	}

	engines := make([]api.RemoteEngine, 2)
	optimizers := []Optimizer{DistributedExecutionOptimizer{Engines: engines}}
	replacements := map[string]*regexp.Regexp{
		" ": spaces,
		"(": openParenthesis,
		")": closedParenthesis,
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			plan := New(expr, time.Unix(0, 0), time.Unix(0, 0))
			optimizedPlan := plan.Optimize(optimizers)
			expectedPlan := cleanUp(replacements, tcase.expected)
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
