// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"math/rand"
	"testing"

	"github.com/cortexproject/promqlsmith"
	"github.com/efficientgo/core/testutil"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"

	"github.com/thanos-io/promql-engine/query"
)

const testRuns = 100

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

func FuzzNodesMarshalJSON(f *testing.F) {
	f.Add(int64(0))
	f.Fuzz(func(t *testing.T, seed int64) {
		lbls := []labels.Labels{
			labels.FromStrings("__name__", "http_requests_total"),
		}
		opts := []promqlsmith.Option{
			promqlsmith.WithEnableOffset(true),
			promqlsmith.WithEnableAtModifier(true),
		}
		rnd := rand.New(rand.NewSource(seed))
		pqSmith := promqlsmith.New(rnd, lbls, opts...)
		for i := 0; i < testRuns; i++ {
			qry := pqSmith.WalkRangeQuery()
			parser.Inspect(qry, func(node parser.Node, nodes []parser.Node) error {
				switch vs := (node).(type) {
				case *parser.VectorSelector:
					vs.Series = nil
					vs.UnexpandedSeriesSet = nil
				}
				return nil
			})

			original := NewFromAST(qry, &query.Options{}, PlanOptions{})
			original, _ = original.Optimize(DefaultOptimizers)

			bytes, err := Marshal(original.Root())
			testutil.Ok(t, err)

			clone, err := Unmarshal(bytes)
			testutil.Ok(t, err)
			testutil.Equals(t, original.Root().String(), clone.String())
		}
	})
}
