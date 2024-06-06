// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"testing"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/promql/parser"

	"github.com/thanos-io/promql-engine/query"
)

func TestNodesMarshalJSON(t *testing.T) {
	var cases = []struct {
		name  string
		query string
	}{
		{
			name: "complex query",
			query: `
sum(
  max_over_time(sum by (pod) (2 * -(rate(http_requests_total[1h])))[2m:1m]) 
  +
  http_requests_total{job="api-server"} @ end()
  + label_replace(metric, "new_label", "$1", "label", ".*")
)`,
		},
		{
			name:  "+Inf",
			query: "clamp_max(metric, +Inf)",
		},
		{
			name:  "NaN",
			query: "clamp_max(metric, NaN)",
		},
		{
			name:  "-Inf",
			query: "clamp_max(metric, -Inf)",
		},
	}
	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			ast, err := parser.ParseExpr(tcase.query)
			testutil.Ok(t, err)
			original := NewFromAST(ast, &query.Options{}, PlanOptions{})
			original, _ = original.Optimize(DefaultOptimizers)

			bytes, err := Marshal(original.Root())
			testutil.Ok(t, err)

			clone, err := Unmarshal(bytes)
			testutil.Ok(t, err)
			testutil.Equals(t, original.Root().String(), clone.String())
		})
	}
}

func TestUnmarshalMatchers(t *testing.T) {
	expr := `metric{name=~"value"}`
	ast, err := parser.ParseExpr(expr)
	testutil.Ok(t, err)

	original := NewFromAST(ast, &query.Options{}, PlanOptions{})
	bytes, err := Marshal(original.Root())
	testutil.Ok(t, err)
	clone, err := Unmarshal(bytes)
	testutil.Ok(t, err)
	testutil.Equals(t, original.Root().String(), clone.String())

	vs, ok := clone.(*VectorSelector)
	testutil.Assert(t, true, ok)
	testutil.Assert(t, true, vs.LabelMatchers[0].Matches("value"))
}
