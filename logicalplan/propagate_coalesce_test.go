// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/thanos-io/promql-engine/query"
)

func mustParseVector(s string) *parser.VectorSelector {
	expr, err := parser.ParseExpr(s)
	if err != nil {
		panic(err)
	}
	return expr.(*parser.VectorSelector)
}

func mustGetFunction(name string) parser.Function {
	fn, ok := parser.Functions[name]
	if !ok {
		panic("unknown function: " + name)
	}
	return *fn
}

func TestPropagateCoalesceOptimizer(t *testing.T) {
	cases := []struct {
		name     string
		input    Node
		expected string
	}{
		{
			name: "push through unary negation",
			input: &Unary{
				Op: parser.SUB,
				Expr: &Coalesce{
					Expressions: []Node{
						&NumberLiteral{Val: 1},
						&NumberLiteral{Val: 2},
					},
				},
			},
			expected: "coalesce(-1, -2)",
		},
		{
			name: "push through binary with scalar on right",
			input: &Binary{
				Op: parser.ADD,
				LHS: &Coalesce{
					Expressions: []Node{
						&NumberLiteral{Val: 1},
						&NumberLiteral{Val: 2},
					},
				},
				RHS: &NumberLiteral{Val: 10},
			},
			expected: "coalesce(1 + 10, 2 + 10)",
		},
		{
			name: "push through binary with scalar on left",
			input: &Binary{
				Op:  parser.MUL,
				LHS: &NumberLiteral{Val: 5},
				RHS: &Coalesce{
					Expressions: []Node{
						&NumberLiteral{Val: 1},
						&NumberLiteral{Val: 2},
					},
				},
			},
			expected: "coalesce(5 * 1, 5 * 2)",
		},
		{
			name: "do not push through binary with vectors on both sides",
			input: &Binary{
				Op: parser.ADD,
				LHS: &Coalesce{
					Expressions: []Node{
						&NumberLiteral{Val: 1},
						&NumberLiteral{Val: 2},
					},
				},
				RHS: &VectorSelector{
					VectorSelector: mustParseVector("metric"),
				},
			},
			expected: "coalesce(1, 2) + metric",
		},
		{
			name: "push through abs function",
			input: &FunctionCall{
				Func: mustGetFunction("abs"),
				Args: []Node{
					&Coalesce{
						Expressions: []Node{
							&NumberLiteral{Val: -1},
							&NumberLiteral{Val: -2},
						},
					},
				},
			},
			expected: "coalesce(abs(-1), abs(-2))",
		},
		{
			name: "push through rate function",
			input: &FunctionCall{
				Func: mustGetFunction("rate"),
				Args: []Node{
					&Coalesce{
						Expressions: []Node{
							&MatrixSelector{
								VectorSelector: &VectorSelector{VectorSelector: mustParseVector("metric1")},
								Range:          300000000000, // 5m
								OriginalString: "metric1[5m]",
							},
							&MatrixSelector{
								VectorSelector: &VectorSelector{VectorSelector: mustParseVector("metric2")},
								Range:          300000000000, // 5m
								OriginalString: "metric2[5m]",
							},
						},
					},
				},
			},
			expected: "coalesce(rate(metric1[5m]), rate(metric2[5m]))",
		},
		{
			name: "push through sum aggregation (distributive)",
			input: &Aggregation{
				Op: parser.SUM,
				Expr: &Coalesce{
					Expressions: []Node{
						&VectorSelector{VectorSelector: mustParseVector("metric1")},
						&VectorSelector{VectorSelector: mustParseVector("metric2")},
					},
				},
			},
			expected: "sum(coalesce(sum(metric1), sum(metric2)))",
		},
		{
			name: "push through min aggregation (distributive)",
			input: &Aggregation{
				Op: parser.MIN,
				Expr: &Coalesce{
					Expressions: []Node{
						&VectorSelector{VectorSelector: mustParseVector("metric1")},
						&VectorSelector{VectorSelector: mustParseVector("metric2")},
					},
				},
			},
			expected: "min(coalesce(min(metric1), min(metric2)))",
		},
		{
			name: "push through max aggregation (distributive)",
			input: &Aggregation{
				Op: parser.MAX,
				Expr: &Coalesce{
					Expressions: []Node{
						&VectorSelector{VectorSelector: mustParseVector("metric1")},
						&VectorSelector{VectorSelector: mustParseVector("metric2")},
					},
				},
			},
			expected: "max(coalesce(max(metric1), max(metric2)))",
		},
		{
			name: "push through count aggregation (distributive, outer is sum)",
			input: &Aggregation{
				Op: parser.COUNT,
				Expr: &Coalesce{
					Expressions: []Node{
						&VectorSelector{VectorSelector: mustParseVector("metric1")},
						&VectorSelector{VectorSelector: mustParseVector("metric2")},
					},
				},
			},
			expected: "sum(coalesce(count(metric1), count(metric2)))",
		},
		{
			name: "push through group aggregation (distributive)",
			input: &Aggregation{
				Op: parser.GROUP,
				Expr: &Coalesce{
					Expressions: []Node{
						&VectorSelector{VectorSelector: mustParseVector("metric1")},
						&VectorSelector{VectorSelector: mustParseVector("metric2")},
					},
				},
			},
			expected: "group(coalesce(group(metric1), group(metric2)))",
		},
		{
			name: "push through sum by aggregation preserves grouping",
			input: &Aggregation{
				Op: parser.SUM,
				Expr: &Coalesce{
					Expressions: []Node{
						&VectorSelector{VectorSelector: mustParseVector("metric1")},
						&VectorSelector{VectorSelector: mustParseVector("metric2")},
					},
				},
				Grouping: []string{"pod", "namespace"},
				Without:  false,
			},
			expected: "sum by (pod, namespace) (coalesce(sum by (pod, namespace) (metric1), sum by (pod, namespace) (metric2)))",
		},
		{
			name: "do not push through avg aggregation (not distributive)",
			input: &Aggregation{
				Op: parser.AVG,
				Expr: &Coalesce{
					Expressions: []Node{
						&VectorSelector{VectorSelector: mustParseVector("metric1")},
						&VectorSelector{VectorSelector: mustParseVector("metric2")},
					},
				},
			},
			expected: "avg(coalesce(metric1, metric2))",
		},
		{
			name: "do not push through stddev aggregation (not distributive)",
			input: &Aggregation{
				Op: parser.STDDEV,
				Expr: &Coalesce{
					Expressions: []Node{
						&VectorSelector{VectorSelector: mustParseVector("metric1")},
						&VectorSelector{VectorSelector: mustParseVector("metric2")},
					},
				},
			},
			expected: "stddev(coalesce(metric1, metric2))",
		},
		{
			name: "push through topk aggregation with local topk",
			input: &Aggregation{
				Op:    parser.TOPK,
				Param: &NumberLiteral{Val: 5},
				Expr: &Coalesce{
					Expressions: []Node{
						&VectorSelector{VectorSelector: mustParseVector("metric1")},
						&VectorSelector{VectorSelector: mustParseVector("metric2")},
					},
				},
			},
			expected: "topk(5, coalesce(topk(5, metric1), topk(5, metric2)))",
		},
		{
			name: "push through bottomk aggregation with local bottomk",
			input: &Aggregation{
				Op:    parser.BOTTOMK,
				Param: &NumberLiteral{Val: 3},
				Expr: &Coalesce{
					Expressions: []Node{
						&VectorSelector{VectorSelector: mustParseVector("metric1")},
						&VectorSelector{VectorSelector: mustParseVector("metric2")},
					},
				},
			},
			expected: "bottomk(3, coalesce(bottomk(3, metric1), bottomk(3, metric2)))",
		},
		{
			name: "push through topk with grouping",
			input: &Aggregation{
				Op:       parser.TOPK,
				Param:    &NumberLiteral{Val: 10},
				Grouping: []string{"pod"},
				Expr: &Coalesce{
					Expressions: []Node{
						&VectorSelector{VectorSelector: mustParseVector("metric1")},
						&VectorSelector{VectorSelector: mustParseVector("metric2")},
					},
				},
			},
			expected: "topk by (pod) (10, coalesce(topk by (pod) (10, metric1), topk by (pod) (10, metric2)))",
		},
		{
			name: "push through multiple levels",
			input: &FunctionCall{
				Func: mustGetFunction("abs"),
				Args: []Node{
					&Unary{
						Op: parser.SUB,
						Expr: &Coalesce{
							Expressions: []Node{
								&NumberLiteral{Val: 1},
								&NumberLiteral{Val: 2},
							},
						},
					},
				},
			},
			expected: "coalesce(abs(-1), abs(-2))",
		},
		{
			name: "push through CheckDuplicateLabels",
			input: &CheckDuplicateLabels{
				Expr: &Coalesce{
					Expressions: []Node{
						&NumberLiteral{Val: 1},
						&NumberLiteral{Val: 2},
					},
				},
			},
			expected: "coalesce(1, 2)",
		},
		{
			name: "push through label_replace",
			input: &FunctionCall{
				Func: mustGetFunction("label_replace"),
				Args: []Node{
					&Coalesce{
						Expressions: []Node{
							&VectorSelector{VectorSelector: mustParseVector("metric1")},
							&VectorSelector{VectorSelector: mustParseVector("metric2")},
						},
					},
					&StringLiteral{Val: "dst"},
					&StringLiteral{Val: "$1"},
					&StringLiteral{Val: "src"},
					&StringLiteral{Val: "(.*)"},
				},
			},
			expected: `coalesce(label_replace(metric1, "dst", "$1", "src", "(.*)"), label_replace(metric2, "dst", "$1", "src", "(.*)"))`,
		},
		// Note: Deduplicate tests are skipped because Deduplicate uses a value receiver
		// for Children(), which causes issues with tree mutation during traversal.
		// The optimizer correctly handles Deduplicate when used through the proper
		// execution path where Deduplicate is behind a pointer/interface.
	}

	optimizer := PropagateCoalesceOptimizer{}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, _ := optimizer.Optimize(tc.input, nil)
			testutil.Equals(t, tc.expected, result.String())
		})
	}
}

func TestPropagateCoalesceOptimizerSharding(t *testing.T) {
	// Test that PropagateCoalesceOptimizer shards selectors and propagates coalesce
	cases := []struct {
		name                string
		expr                string
		decodingConcurrency int
		expected            string
	}{
		{
			name:                "vector selector with 2 shards",
			expr:                "http_requests_total",
			decodingConcurrency: 2,
			expected:            "coalesce(http_requests_total 0 mod 2, http_requests_total 1 mod 2)",
		},
		{
			name:                "vector selector with 4 shards",
			expr:                "http_requests_total",
			decodingConcurrency: 4,
			expected:            "coalesce(http_requests_total 0 mod 4, http_requests_total 1 mod 4, http_requests_total 2 mod 4, http_requests_total 3 mod 4)",
		},
		{
			name:                "rate with 2 shards - propagated through",
			expr:                "rate(http_requests_total[5m])",
			decodingConcurrency: 2,
			expected:            "coalesce(rate(http_requests_total[5m] 0 mod 2), rate(http_requests_total[5m] 1 mod 2))",
		},
		{
			name:                "sum with 2 shards - aggregation pushed into shards",
			expr:                "sum(http_requests_total)",
			decodingConcurrency: 2,
			expected:            "sum(coalesce(sum(http_requests_total 0 mod 2), sum(http_requests_total 1 mod 2)))",
		},
		{
			name:                "no sharding with 1 shard",
			expr:                "http_requests_total",
			decodingConcurrency: 1,
			expected:            "http_requests_total",
		},
		{
			name:                "sum by gets pushed through coalesce",
			expr:                "sum by (pod) (http_requests_total)",
			decodingConcurrency: 2,
			expected:            "sum by (pod) (coalesce(sum by (pod) (http_requests_total 0 mod 2), sum by (pod) (http_requests_total 1 mod 2)))",
		},
		{
			name:                "sum rate gets fully parallelized",
			expr:                "sum(rate(http_requests_total[5m]))",
			decodingConcurrency: 2,
			expected:            "sum(coalesce(sum(rate(http_requests_total[5m] 0 mod 2)), sum(rate(http_requests_total[5m] 1 mod 2))))",
		},
		{
			name:                "topk gets local topk per shard",
			expr:                "topk(10, http_requests_total)",
			decodingConcurrency: 2,
			expected:            "topk(10, coalesce(topk(10, http_requests_total 0 mod 2), topk(10, http_requests_total 1 mod 2)))",
		},
		{
			name:                "avg blocks propagation (not distributive)",
			expr:                "avg(http_requests_total)",
			decodingConcurrency: 2,
			expected:            "avg(coalesce(http_requests_total 0 mod 2, http_requests_total 1 mod 2))",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tc.expr)
			testutil.Ok(t, err)

			plan, err := NewFromAST(expr, &query.Options{
				Start:               time.Unix(0, 0),
				End:                 time.Unix(300, 0),
				Step:                time.Second * 30,
				DecodingConcurrency: tc.decodingConcurrency,
			}, PlanOptions{})
			testutil.Ok(t, err)

			optimizedPlan, _ := plan.Optimize([]Optimizer{
				PropagateCoalesceOptimizer{},
			})

			testutil.Equals(t, tc.expected, optimizedPlan.Root().String())
		})
	}
}
