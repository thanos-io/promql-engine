// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"testing"

	"github.com/efficientgo/core/testutil"
	"github.com/thanos-io/promql-engine/parser"
)

func TestRegexMatchers(t *testing.T) {
	cases := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "valid regex match",
			expr:     `count(X{ding=~"abcd.+"})`,
			expected: `count(X{ding=~"abcd.+"})`,
		},
		{
			name:     "valid regex match with \\d",
			expr:     `count(X{ding=~"abcd\\d"})`,
			expected: `count(X{ding=~"abcd\\d"})`,
		},
		{
			name:     "part of query regex match",
			expr:     `my_awesome_metrics{has_observability=~"eveything",something=~"nothing.+"}`,
			expected: `my_awesome_metrics{has_observability="eveything",something=~"nothing.+"}`,
		},
	}
	optimizers := []Optimizer{RegexMatchersFunctions{}}
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
