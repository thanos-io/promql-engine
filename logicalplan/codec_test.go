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
	expr := `
sum(
  max_over_time(sum by (pod) (2 * -(rate(http_requests_total[1h])))[2m:1m]) 
  +
  http_requests_total{job="api-server"} @ end()
)`
	ast, err := parser.ParseExpr(expr)
	testutil.Ok(t, err)
	original := New(ast, &query.Options{}, PlanOptions{})
	original, _ = original.Optimize(DefaultOptimizers)

	bytes, err := Marshal(original.Root())
	testutil.Ok(t, err)

	clone, err := Unmarshal(bytes)
	testutil.Ok(t, err)
	testutil.Equals(t, original.Root().String(), clone.String())
}
