// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"regexp"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"

	"github.com/thanos-community/promql-engine/api"
)

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
sum by (pod) (dedup(
  remote(sum by (pod, region) (rate(http_requests_total[5m]))), 
  remote(sum by (pod, region) (rate(http_requests_total[5m])))))`,
		},
		{
			name: "sum-rate without labels preserves engine labels",
			expr: `sum without (pod, region) (rate(http_requests_total[5m]))`,
			expected: `
sum without (pod, region) (
  dedup(
    remote(sum without (pod) (rate(http_requests_total[5m]))),
    remote(sum without (pod) (rate(http_requests_total[5m])))
  )
)`,
		},
		{
			name: "avg",
			expr: `avg by (pod) (http_requests_total)`,
			expected: `
avg by (pod) (
  dedup(
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
    dedup(
      remote(sum by (pod, region) (http_requests_total)),
      remote(sum by (pod, region) (http_requests_total))
    )
  )
)`,
		},
		{
			name: "aggregation of binary expression",
			expr: `max by (pod) (metric_a / metric_b)`,
			expected: `
max by (pod) (
  dedup(remote(metric_a), remote(metric_a)) 
  / 
  dedup(remote(metric_b), remote(metric_b))
)
`,
		},
		{
			name: "unsupported aggregation in the operand path",
			expr: `max by (pod) (sort(avg(http_requests_total)))`,
			expected: `
max by (pod) (sort(avg(
  dedup(
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
  dedup(remote(metric_a), remote(metric_a)) 
  / 
  dedup(remote(metric_b), remote(metric_b))
))`,
		},
		{
			name: "binary operation with aggregations",
			expr: `sum by (pod) (metric_a) / sum by (pod) (metric_b)`,
			expected: `
sum by (pod) (dedup(
  remote(sum by (pod, region) (metric_a)), 
  remote(sum by (pod, region) (metric_a)))
)
/ 
sum by (pod) (dedup(
  remote(sum by (pod, region) (metric_b)), 
  remote(sum by (pod, region) (metric_b))
))`,
		},
		{
			name: "function sharding",
			expr: `rate(http_requests_total[2m])`,
			expected: `
dedup(
  remote(rate(http_requests_total[2m])), 
  remote(rate(http_requests_total[2m]))
)`,
		},
		{
			name: `histogram quantile`,
			expr: `histogram_quantile(0.5, sum by (le) (rate(coredns_dns_request_duration_seconds_bucket[5m])))`,
			expected: `
histogram_quantile(0.5, sum by (le) (dedup(
  remote(sum by (le, region) (rate(coredns_dns_request_duration_seconds_bucket[5m]))), 
  remote(sum by (le, region) (rate(coredns_dns_request_duration_seconds_bucket[5m])))
)))`,
		},
		{
			name:     "binary expression with time",
			expr:     `time() - max by (foo) (bar)`,
			expected: `time() - max by (foo) (dedup(remote(max by (foo, region) (bar)), remote(max by (foo, region) (bar))))`,
		},
		{
			name:     "number literal",
			expr:     `1`,
			expected: `1`,
		},
		{
			name:     "aggregation with number literal",
			expr:     `max(foo) - 1`,
			expected: `max(dedup(remote(max by (region) (foo)), remote(max by (region) (foo)))) - 1`,
		},
		{
			name:     "absent",
			expr:     `absent(foo)`,
			expected: `remote(absent(foo)) * remote(absent(foo))`,
		},
	}

	engines := []api.RemoteEngine{
		newEngineMock(1, []labels.Labels{labels.FromStrings("region", "east")}),
		newEngineMock(2, []labels.Labels{labels.FromStrings("region", "west")}),
	}
	optimizers := []Optimizer{DistributedExecutionOptimizer{Endpoints: api.NewStaticEndpoints(engines)}}
	replacements := map[string]*regexp.Regexp{
		" ": spaces,
		"(": openParenthesis,
		")": closedParenthesis,
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			plan := New(expr, &Opts{Start: time.Unix(0, 0), End: time.Unix(0, 0)})
			optimizedPlan := plan.Optimize(optimizers)
			expectedPlan := cleanUp(replacements, tcase.expected)
			testutil.Equals(t, expectedPlan, optimizedPlan.Expr().String())
		})
	}
}

type engineMock struct {
	api.RemoteEngine
	minT      int64
	maxT      int64
	labelSets []labels.Labels
}

func (e engineMock) MaxT() int64 {
	return e.maxT
}

func (e engineMock) MinT() int64 {
	return e.minT
}

func (e engineMock) LabelSets() []labels.Labels {
	return e.labelSets
}

func newEngineMock(maxT int64, labelSets []labels.Labels) *engineMock {
	return &engineMock{maxT: maxT, labelSets: labelSets}
}
