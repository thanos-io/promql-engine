// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"math"
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
sum by (pod) (
  dedup(coalesce(
    remote(sum by (pod, region) (rate(http_requests_total[5m]))),
    remote(sum by (pod, region) (rate(http_requests_total[5m])))
  ))
)`,
		},
		{
			name: "sum-rate without labels preserves engine labels",
			expr: `sum without (pod, region) (rate(http_requests_total[5m]))`,
			expected: `
sum without (pod, region) (
  dedup(coalesce(
    remote(sum without (pod) (rate(http_requests_total[5m]))),
    remote(sum without (pod) (rate(http_requests_total[5m])))
  ))
)`,
		},
		{
			name: "avg",
			expr: `avg by (pod) (http_requests_total)`,
			expected: `
avg by (pod) (
  dedup(coalesce(
    remote(http_requests_total),
    remote(http_requests_total)
  ))
)`,
		},
		{
			name: "two-level aggregation",
			expr: `max by (pod) (sum by (pod) (http_requests_total))`,
			expected: `
max by (pod) (
  sum by (pod) ( 
    dedup(coalesce(
      remote(sum by (pod, region) (http_requests_total)),
      remote(sum by (pod, region) (http_requests_total))
    ))
  )
)`,
		},
		{
			name: "aggregation of binary expression",
			expr: `max by (pod) (metric_a / metric_b)`,
			expected: `
max by (pod) (
  dedup(coalesce(remote(metric_a), remote(metric_a))) 
  / 
  dedup(coalesce(remote(metric_b), remote(metric_b)))
)
`,
		},
		{
			name: "unsupported aggregation in the operand path",
			expr: `max by (pod) (sort(avg(http_requests_total)))`,
			expected: `
max by (pod) (sort(avg(
  dedup(coalesce(
    remote(http_requests_total),
    remote(http_requests_total)
  ))
)))`,
		},
		{
			name: "binary operation in the operand path",
			expr: `max by (pod) (sort(metric_a / metric_b))`,
			expected: `
max by (pod) (sort(
  dedup(coalesce(remote(metric_a), remote(metric_a))) 
  / 
  dedup(coalesce(remote(metric_b), remote(metric_b)))
))`,
		},
		{
			name: "binary operation with aggregations",
			expr: `sum by (pod) (metric_a) / sum by (pod) (metric_b)`,
			expected: `
sum by (pod) (dedup(coalesce(
  remote(sum by (pod, region) (metric_a)), 
  remote(sum by (pod, region) (metric_a)))
))
/ 
sum by (pod) (dedup(coalesce(
  remote(sum by (pod, region) (metric_b)), 
  remote(sum by (pod, region) (metric_b))))
)`,
		},
		{
			name: "function sharding",
			expr: `rate(http_requests_total[2m])`,
			expected: `
dedup(coalesce(
  remote(rate(http_requests_total[2m])), 
  remote(rate(http_requests_total[2m]))))`,
		},
	}

	engines := []api.RemoteEngine{
		newEngineMock([]labels.Labels{labels.FromStrings("region", "east")}),
		newEngineMock([]labels.Labels{labels.FromStrings("region", "west")}),
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

			plan := New(expr, time.Unix(0, 0), time.Unix(0, 0))
			optimizedPlan := plan.Optimize(optimizers)
			expectedPlan := cleanUp(replacements, tcase.expected)
			testutil.Equals(t, expectedPlan, optimizedPlan.Expr().String())
		})
	}
}

type engineMock struct {
	api.RemoteEngine
	labelSets []labels.Labels
}

func (e engineMock) MaxT() int64 {
	return math.MaxInt64
}

func (e engineMock) LabelSets() []labels.Labels {
	return e.labelSets
}

func newEngineMock(labelSets []labels.Labels) *engineMock {
	return &engineMock{labelSets: labelSets}
}
