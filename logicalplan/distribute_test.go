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

	"github.com/thanos-io/promql-engine/api"
	"github.com/thanos-io/promql-engine/query"
)

func TestDistributedExecution(t *testing.T) {
	cases := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "selector",
			expr:     `http_requests_total`,
			expected: `dedup(remote(http_requests_total), remote(http_requests_total))`,
		},
		{
			name:     "rate",
			expr:     `rate(http_requests_total[5m])`,
			expected: `dedup(remote(rate(http_requests_total[5m])), remote(rate(http_requests_total[5m])))`,
		},
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
			expr: `avg(http_requests_total)`,
			expected: `
sum(
  dedup(
    remote(sum by (region) (http_requests_total)),
    remote(sum by (region) (http_requests_total))
  )
) /
sum(
  dedup(
    remote(count by (region) (http_requests_total)),
    remote(count by (region) (http_requests_total))
  )
)`,
		},
		{
			name: "avg with grouping",
			expr: `avg by (pod) (http_requests_total)`,
			expected: `
sum by (pod) (
  dedup(
    remote(sum by (pod, region) (http_requests_total)),
    remote(sum by (pod, region) (http_requests_total))
  )
) /
sum by (pod) (
  dedup(
    remote(count by (pod, region) (http_requests_total)),
    remote(count by (pod, region) (http_requests_total))
  )
)`,
		},
		{
			name: "avg with prior aggregation",
			expr: `avg by (pod) (sum by (pod) (http_requests_total))`,
			expected: `
avg by (pod) (
  sum by (pod) (
	dedup(
      remote(sum by (pod, region) (http_requests_total)),
	  remote(sum by (pod, region) (http_requests_total))
	)
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
			expr: `max by (pod) (sort(quantile(0.9, http_requests_total)))`,
			expected: `
max by (pod) (quantile(0.9,
  dedup(
    remote(http_requests_total),
    remote(http_requests_total)
  )
))`,
		},
		{
			name: "binary operation in the operand path",
			expr: `max by (pod) (metric_a / metric_b)`,
			expected: `
max by (pod) (
  dedup(remote(metric_a), remote(metric_a)) 
  / 
  dedup(remote(metric_b), remote(metric_b))
)`,
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
		}, {
			name: "binary expression with constant",
			expr: `sum by (pod) (rate(http_requests_total[2m]) * 60)`,
			expected: `sum by (pod) (dedup(
remote(sum by (pod, region) (rate(http_requests_total[2m]) * 60)), 
remote(sum by (pod, region) (rate(http_requests_total[2m]) * 60))))`,
		},
		{
			name:     "label based pruning matches one engine",
			expr:     `sum by (pod) (rate(http_requests_total{region="west"}[2m]))`,
			expected: `sum by (pod) (dedup(remote(sum by (pod, region) (rate(http_requests_total{region="west"}[2m])))))`,
		},
		{
			name:     "label based pruning matches no engines",
			expr:     `http_requests_total{region="north"}`,
			expected: `noop`,
		},
		{
			name:     "label based pruning with grouping matches no engines",
			expr:     `sum by (pod) (rate(http_requests_total{region="north"}[2m]))`,
			expected: `sum by (pod) (noop)`,
		},
		{
			name:     "label based pruning with grouping matches single engine",
			expr:     `sum by (pod) (rate(http_requests_total{region="south"}[2m]))`,
			expected: `sum by (pod) (dedup(remote(sum by (pod, region) (rate(http_requests_total{region="south"}[2m])))))`,
		},
	}

	engines := []api.RemoteEngine{
		newEngineMock(math.MinInt64, math.MinInt64, []labels.Labels{labels.FromStrings("region", "east"), labels.FromStrings("region", "south")}),
		newEngineMock(math.MinInt64, math.MinInt64, []labels.Labels{labels.FromStrings("region", "west")}),
	}
	optimizers := []Optimizer{
		DistributeAvgOptimizer{},
		DistributedExecutionOptimizer{Endpoints: api.NewStaticEndpoints(engines)},
	}
	replacements := map[string]*regexp.Regexp{
		" ": spaces,
		"(": openParenthesis,
		")": closedParenthesis,
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			plan := New(expr, &query.Options{Start: time.Unix(0, 0), End: time.Unix(0, 0)})
			optimizedPlan := plan.Optimize(optimizers)
			expectedPlan := cleanUp(replacements, tcase.expected)
			testutil.Equals(t, expectedPlan, optimizedPlan.Expr().String())
		})
	}
}

type engineOpts struct {
	minTime time.Time
	maxTime time.Time
}

func (o engineOpts) mint() int64 {
	return o.minTime.UnixMilli()
}

func (o engineOpts) maxt() int64 {
	return o.maxTime.UnixMilli()
}

func TestDistributedExecutionWithLongSelectorRanges(t *testing.T) {
	replacements := map[string]*regexp.Regexp{
		" ": spaces,
		"(": openParenthesis,
		")": closedParenthesis,
	}

	sixHours := 6 * time.Hour
	eightHours := 8 * time.Hour
	twelveHours := 12 * time.Hour

	queryStart := time.Unix(0, 0)
	queryEnd := time.Unix(0, 0).Add(twelveHours)
	queryStep := time.Minute

	cases := []struct {
		name             string
		expr             string
		expected         string
		firstEngineOpts  engineOpts
		secondEngineOpts engineOpts
	}{
		{
			name: "sum over 5m adds a 5 minute offset to latest engine",
			firstEngineOpts: engineOpts{
				minTime: queryStart,
				maxTime: time.Unix(0, 0).Add(eightHours),
			},
			secondEngineOpts: engineOpts{
				minTime: time.Unix(0, 0).Add(sixHours),
				maxTime: queryEnd,
			},
			expr: `sum_over_time(metric[5m])`,
			expected: `
dedup(
  remote(sum_over_time(metric[5m])),
  remote(sum_over_time(metric[5m])) [1970-01-01 06:05:00 +0000 UTC]
)`,
		},
		{
			name: "sum over 2h adds a 2 hour offset to latest engine",
			firstEngineOpts: engineOpts{
				minTime: queryStart,
				maxTime: time.Unix(0, 0).Add(eightHours),
			},
			secondEngineOpts: engineOpts{
				minTime: time.Unix(0, 0).Add(sixHours),
				maxTime: queryEnd,
			},
			expr: `sum_over_time(metric[2h])`,
			expected: `
dedup(
  remote(sum_over_time(metric[2h])),
  remote(sum_over_time(metric[2h])) [1970-01-01 08:00:00 +0000 UTC]
)`,
		},
		{
			name: "sum over 3h does not distribute the query due to insufficient engine overlap",
			firstEngineOpts: engineOpts{
				minTime: queryStart,
				maxTime: time.Unix(0, 0).Add(eightHours),
			},
			secondEngineOpts: engineOpts{
				minTime: time.Unix(0, 0).Add(sixHours),
				maxTime: queryEnd,
			},
			expr:     `sum_over_time(metric[3h])`,
			expected: `sum_over_time(metric[3h])`,
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			engines := []api.RemoteEngine{
				newEngineMock(tcase.firstEngineOpts.mint(), tcase.firstEngineOpts.maxt(), []labels.Labels{labels.FromStrings("region", "east")}),
				newEngineMock(tcase.secondEngineOpts.mint(), tcase.secondEngineOpts.maxt(), []labels.Labels{labels.FromStrings("region", "east")}),
			}
			optimizers := []Optimizer{
				DistributeAvgOptimizer{},
				DistributedExecutionOptimizer{Endpoints: api.NewStaticEndpoints(engines)},
			}

			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			plan := New(expr, &query.Options{Start: queryStart, End: queryEnd, Step: queryStep})
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

func newEngineMock(mint, maxt int64, labelSets []labels.Labels) *engineMock {
	return &engineMock{minT: mint, maxT: maxt, labelSets: labelSets}
}
