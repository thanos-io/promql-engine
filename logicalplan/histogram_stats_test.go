// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"testing"

	"github.com/thanos-io/promql-engine/query"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/promql/parser"
)

func TestDetectHistogramStatsOptimizer(t *testing.T) {
	cases := []struct {
		name string
		expr string
		want []bool
	}{
		{
			name: "histogram stats function",
			expr: `histogram_count(metric)`,
			want: []bool{true},
		},
		{
			name: "histogram stats through range function",
			expr: `histogram_count(rate(metric[5m]))`,
			want: []bool{true},
		},
		{
			name: "whole histogram function above stats function",
			expr: `histogram_quantile(0.5, histogram_count(metric))`,
			want: []bool{false},
		},
		{
			name: "histogram fraction above stats function",
			expr: `histogram_fraction(0, 1, histogram_count(metric))`,
			want: []bool{false},
		},
		{
			name: "whole histogram function with stats function in filtering branch",
			expr: `histogram_quantile(0.5, metric unless histogram_count(metric) == 0)`,
			want: []bool{false, false},
		},
		{
			name: "subquery under stats function",
			expr: `histogram_count(increase(metric[40m:9m]))`,
			want: []bool{false},
		},
		{
			name: "histogram stddev does not block upstream detector",
			expr: `histogram_count(histogram_stddev(metric) * metric)`,
			want: []bool{true, true},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := parser.NewParser(parser.Options{}).ParseExpr(tc.expr)
			testutil.Ok(t, err)

			plan, err := NewFromAST(expr, &query.Options{}, PlanOptions{})
			testutil.Ok(t, err)

			optimized, _ := DetectHistogramStatsOptimizer{}.Optimize(plan.Root(), nil)
			testutil.Equals(t, tc.want, histogramStatsDecoding(optimized))
		})
	}
}

func TestDetectHistogramStatsOptimizerHistogramQuantiles(t *testing.T) {
	// histogram_quantiles is not available in the currently pinned Prometheus
	// parser, but it is handled by the upstream detector. Build a minimal logical
	// tree directly to keep this behavior covered.
	metric := &VectorSelector{VectorSelector: &parser.VectorSelector{Name: "metric"}}
	root := &FunctionCall{
		Func: parser.Function{Name: "histogram_quantiles"},
		Args: []Node{
			&FunctionCall{
				Func: parser.Function{Name: "histogram_count"},
				Args: []Node{metric},
			},
			&StringLiteral{Val: "q"},
			&NumberLiteral{Val: 0.5},
		},
	}

	optimized, _ := DetectHistogramStatsOptimizer{}.Optimize(root, nil)
	testutil.Equals(t, []bool{false}, histogramStatsDecoding(optimized))
}

func histogramStatsDecoding(root Node) []bool {
	var got []bool
	Traverse(&root, func(node *Node) {
		vs, ok := (*node).(*VectorSelector)
		if !ok {
			return
		}
		got = append(got, vs.DecodeNativeHistogramStats)
	})
	return got
}
