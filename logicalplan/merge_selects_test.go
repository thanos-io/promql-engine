// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"testing"

	"github.com/thanos-io/promql-engine/query"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
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
		{
			expr:     `X{a!~"b",a=~"b",c="d"}/X{a=~"b"}`,
			expected: `filter([a!~"b" c="d"], X{a=~"b"}) / X{a=~"b"}`,
		},
		{
			expr:     `quantile by (pod) (scalar(min(http_requests_total)), http_requests_total)`,
			expected: `quantile by (pod) (scalar(min(http_requests_total)), http_requests_total)`,
		},
	}
	optimizers := []Optimizer{MergeSelectsOptimizer{}}
	for _, tcase := range cases {
		t.Run(tcase.expr, func(t *testing.T) {
			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			plan, _ := NewFromAST(expr, &query.Options{}, PlanOptions{})
			optimizedPlan, _ := plan.Optimize(optimizers)
			testutil.Equals(t, tcase.expected, renderExprTree(optimizedPlan.Root()))
		})
	}
}

func TestMergeSelectsWithProjections(t *testing.T) {
	cases := []struct {
		name     string
		plan     Node
		expected string
	}{
		{
			name: "no merge when left has projection",
			plan: &Binary{
				Op: parser.DIV,
				LHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "X",
						LabelMatchers: []*labels.Matcher{
							{Type: labels.MatchEqual, Name: labels.MetricName, Value: "X"},
							{Type: labels.MatchEqual, Name: "a", Value: "b"},
						},
					},
					Projection: &Projection{Include: true, Labels: []string{"a"}},
				},
				RHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "X",
						LabelMatchers: []*labels.Matcher{
							{Type: labels.MatchEqual, Name: labels.MetricName, Value: "X"},
						},
					},
				},
			},
			expected: `X{a="b"}[projection=include(a)] / X`,
		},
		{
			name: "no merge when right has projection",
			plan: &Binary{
				Op: parser.DIV,
				LHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "X",
						LabelMatchers: []*labels.Matcher{
							{Type: labels.MatchEqual, Name: labels.MetricName, Value: "X"},
						},
					},
				},
				RHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "X",
						LabelMatchers: []*labels.Matcher{
							{Type: labels.MatchEqual, Name: labels.MetricName, Value: "X"},

							{Type: labels.MatchEqual, Name: "a", Value: "b"},
						},
					},
					Projection: &Projection{Include: true, Labels: []string{"a"}},
				},
			},
			expected: `X / X{a="b"}[projection=include(a)]`,
		},
		{
			name: "no merge when both have projections",
			plan: &Binary{
				Op: parser.DIV,
				LHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "X",
						LabelMatchers: []*labels.Matcher{
							{Type: labels.MatchEqual, Name: labels.MetricName, Value: "X"},
						},
					},
					Projection: &Projection{Include: true, Labels: []string{"a"}},
				},
				RHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "X",
						LabelMatchers: []*labels.Matcher{
							{Type: labels.MatchEqual, Name: labels.MetricName, Value: "X"},
							{Type: labels.MatchEqual, Name: "c", Value: "d"},
						},
					},
					Projection: &Projection{Include: true, Labels: []string{"c"}},
				},
			},
			expected: `X[projection=include(a)] / X{c="d"}[projection=include(c)]`,
		},
		{
			name: "merge if empty projection",
			plan: &Binary{
				Op: parser.DIV,
				LHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "X",
						LabelMatchers: []*labels.Matcher{
							{Type: labels.MatchEqual, Name: labels.MetricName, Value: "X"},
							{Type: labels.MatchEqual, Name: "a", Value: "b"},
						},
					},
					Projection: &Projection{},
				},
				RHS: &VectorSelector{
					VectorSelector: &parser.VectorSelector{
						Name: "X",
						LabelMatchers: []*labels.Matcher{
							{Type: labels.MatchEqual, Name: labels.MetricName, Value: "X"},
						},
					},
				},
			},
			expected: `filter([a="b"], X) / X`,
		},
	}

	optimizer := MergeSelectsOptimizer{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			optimizedPlan, _ := optimizer.Optimize(tc.plan, &query.Options{})
			testutil.Equals(t, tc.expected, renderExprTree(optimizedPlan))
		})
	}
}
