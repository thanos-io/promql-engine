// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"testing"

	"github.com/efficientgo/core/testutil"

	"github.com/thanos-io/promql-engine/parser"
	"github.com/thanos-io/promql-engine/query"
)

func TestMergeSelects(t *testing.T) {
	cases := []struct {
		expr     string
		expected string
	}{
		{
			expr:     `X{a="b"}/X`,
			expected: `filter([a="b"], X) / X`,
		},
		{
			expr:     `floor(X{a="b"})/X`,
			expected: `floor(filter([a="b"], X)) / X`,
		},
		{
			expr:     `X/floor(X{a="b"})`,
			expected: `X / floor(filter([a="b"], X))`,
		},
		{
			expr:     `X{a="b"}/floor(X)`,
			expected: `filter([a="b"], X) / floor(X)`,
		},
	}
	optimizers := []Optimizer{MergeSelectsOptimizer{}}
	for _, tcase := range cases {
		t.Run(tcase.expr, func(t *testing.T) {
			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			plan := New(expr, &query.Options{})
			optimizedPlan := plan.Optimize(optimizers)
			testutil.Equals(t, tcase.expected, optimizedPlan.Expr().String())
		})
	}
}
