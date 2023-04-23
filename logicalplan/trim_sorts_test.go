// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"testing"

	"github.com/efficientgo/core/testutil"

	"github.com/thanos-community/promql-engine/parser"
)

func TestTrimSorts(t *testing.T) {
	cases := []struct {
		name     string
		expr     string
		expected string
	}{
		// this test case is ok since the engine determines sorting order
		// before running optimziers
		{
			name:     "simple sort",
			expr:     "sort(X)",
			expected: "X",
		},
		{
			name:     "sort",
			expr:     "sum(sort(X))",
			expected: "sum(X)",
		},
		{
			name:     "nested",
			expr:     "sum(sort(rate(X[1m])))",
			expected: "sum(rate(X[1m]))",
		},
		{
			name:     "weirdly nested",
			expr:     "sum(sort(sqrt(sort(X))))",
			expected: "sum(sqrt(X))",
		},
		{
			name:     "sort in binary expression",
			expr:     "sort(sort(sqrt(X))/sort(sqrt(Y)))",
			expected: "sqrt(X) / sqrt(Y)",
		},
	}
	optimizers := []Optimizer{TrimSortFunctions{}}
	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			plan := New(expr, &Opts{})
			optimizedPlan := plan.Optimize(optimizers)
			testutil.Equals(t, tcase.expected, optimizedPlan.Expr().String())
		})
	}
}
