// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"testing"

	"github.com/thanos-io/promql-engine/query"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/promql/parser"
)

func TestSetBatchSize(t *testing.T) {
	cases := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "selector",
			expr:     `http_requests_total`,
			expected: `http_requests_total`,
		},
		{
			name:     "rate",
			expr:     `rate(http_requests_total[5m])`,
			expected: `rate(http_requests_total[5m0s])`,
		},
		{
			name:     "sum",
			expr:     `sum(http_requests_total)`,
			expected: `sum(http_requests_total[batch=10])`,
		},
		{
			name:     "quantile",
			expr:     `quantile(0.9, http_requests_total)`,
			expected: `quantile(0.9, http_requests_total)`,
		},
		{
			name:     "two-level aggregation",
			expr:     `max by (pod) (sum by (pod) (http_requests_total))`,
			expected: `max by (pod) (sum by (pod) (http_requests_total[batch=10]))`,
		},
		{
			name:     "aggregation of binary expression",
			expr:     `max by (pod) (metric_a / metric_b)`,
			expected: `max by (pod) (metric_a / metric_b)`,
		},
		{
			name:     "binary operation of aggregations",
			expr:     `max(metric_a) / max(metric_b)`,
			expected: `max(metric_a[batch=10]) / max(metric_b[batch=10])`,
		},
		{
			name:     "binary operation with same metric aggregations",
			expr:     `max(metric_a) / max(metric_a{code="foo"})`,
			expected: `max(metric_a[batch=10]) / max(filter([code="foo"], metric_a[batch=10]))`,
		},
		{
			name:     `histogram quantile`,
			expr:     `histogram_quantile(0.5, metric_bucket)`,
			expected: `histogram_quantile(0.5, metric_bucket)`,
		},
		{
			name:     "binary expression with time",
			expr:     `time() - max by (foo) (bar)`,
			expected: `time() - max by (foo) (bar[batch=10])`,
		},
		{
			name:     "binary expression with single aggregation",
			expr:     `metric_a - max by (foo) (bar)`,
			expected: `metric_a - max by (foo) (bar[batch=10])`,
		},
		{
			name:     "number literal",
			expr:     `1`,
			expected: `1`,
		},
		{
			name:     "absent",
			expr:     `absent(foo)`,
			expected: `absent(foo)`,
		},
		{
			name:     "histogram quantile with aggregation",
			expr:     `histogram_quantile(scalar(max(quantile)), http_requests_total)`,
			expected: `histogram_quantile(scalar(max(quantile[batch=10])), http_requests_total)`,
		},
		// Range vector functions should allow batching to propagate
		{
			name:     "sum with rate",
			expr:     `sum(rate(http_requests_total[5m]))`,
			expected: `sum(rate(http_requests_total[batch=10][5m0s]))`,
		},
		{
			name:     "avg with increase",
			expr:     `avg(increase(http_requests_total[5m]))`,
			expected: `avg(increase(http_requests_total[batch=10][5m0s]))`,
		},
		// Label-preserving functions should allow batching to propagate
		{
			name:     "sum with abs",
			expr:     `sum(abs(metric))`,
			expected: `sum(abs(metric[batch=10]))`,
		},
		{
			name:     "avg with ceil",
			expr:     `avg(ceil(metric))`,
			expected: `avg(ceil(metric[batch=10]))`,
		},
		{
			name:     "max with clamp_max",
			expr:     `max(clamp_max(metric, 100))`,
			expected: `max(clamp_max(metric[batch=10], 100))`,
		},
		{
			name:     "nested math functions",
			expr:     `sum(abs(floor(metric)))`,
			expected: `sum(abs(floor(metric[batch=10])))`,
		},
		{
			name: "aggregation with timestamp",
			expr: `sum(timestamp(metric))`,
			// timestamp() is pushed down into VectorSelector with SelectTimestamp=true
			expected: `sum(metric[batch=10][timestamp])`,
		},
		// Label-modifying functions should disable batching
		{
			name:     "label_replace disables batching",
			expr:     `sum(label_replace(metric, "dst", "$1", "src", "(.*)"))`,
			expected: `sum(label_replace(metric[batch=10], "dst", "$1", "src", "(.*)"))`,
		},
		{
			name:     "label_join disables batching",
			expr:     `sum(label_join(metric, "dst", ",", "src"))`,
			expected: `sum(label_join(metric[batch=10], "dst", ",", "src"))`,
		},
		{
			name:     "histogram_fraction disables batching",
			expr:     `sum(histogram_fraction(0, 0.5, metric))`,
			expected: `sum(histogram_fraction(0, 0.5, metric))`,
		},
	}

	optimizers := append([]Optimizer{SelectorBatchSize{Size: 10}}, DefaultOptimizers...)
	for _, tcase := range cases {
		t.Run(tcase.expr, func(t *testing.T) {
			t.Parallel()
			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			plan, _ := NewFromAST(expr, &query.Options{}, PlanOptions{})
			optimizedPlan, _ := plan.Optimize(optimizers)
			testutil.Equals(t, tcase.expected, renderExprTree(optimizedPlan.Root()))
		})
	}
}
