// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"testing"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/promql/parser"

	"github.com/thanos-io/promql-engine/query"
)

func TestShardedAggregations(t *testing.T) {
	cases := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "sum",
			expr:     `topk(10, X)`,
			expected: `topk(10, coalesce(topk(10, X[shard=0/2]), topk(10, X[shard=1/2])))`,
		},
	}

	optimizers := []Optimizer{ShardedAggregations{Shards: 2}}
	for _, tcase := range cases {
		t.Run(tcase.expr, func(t *testing.T) {
			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			plan := New(expr, &query.Options{})
			optimizedPlan, _ := plan.Optimize(optimizers)
			testutil.Equals(t, tcase.expected, optimizedPlan.Expr().String())
		})
	}
}
