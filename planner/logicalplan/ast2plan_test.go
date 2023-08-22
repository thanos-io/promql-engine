package logicalplan

import (
	"github.com/stretchr/testify/require"
	"github.com/thanos-io/promql-engine/parser"
	"math"
	"testing"
)

var ast2planTestCases = []struct {
	input    parser.Expr // The AST input.
	expected LogicalPlan // The expected logical plan.
}{
	{
		input:    &parser.NumberLiteral{Val: 1},
		expected: &NumberLiteral{Val: 1},
	},
	{
		input:    &parser.NumberLiteral{Val: math.Inf(1)},
		expected: &NumberLiteral{Val: math.Inf(1)},
	},
	{
		input:    &parser.NumberLiteral{Val: math.Inf(-1)},
		expected: &NumberLiteral{Val: math.Inf(-1)},
	},
	{
		input: &parser.BinaryExpr{
			Op:  parser.ADD,
			LHS: &parser.NumberLiteral{Val: 1},
			RHS: &parser.NumberLiteral{Val: 1},
		},
		expected: &BinaryExpr{
			Op:  parser.ADD,
			LHS: &NumberLiteral{Val: 1},
			RHS: &NumberLiteral{Val: 1},
		},
	},
	{
		input: &parser.BinaryExpr{
			Op:  parser.ADD,
			LHS: &parser.NumberLiteral{Val: 1},
			RHS: &parser.BinaryExpr{
				Op:  parser.DIV,
				LHS: &parser.NumberLiteral{Val: 2},
				RHS: &parser.ParenExpr{
					Expr: &parser.BinaryExpr{
						Op:  parser.MUL,
						LHS: &parser.NumberLiteral{Val: 3},
						RHS: &parser.NumberLiteral{Val: 1},
					},
				},
			},
		},
		expected: &BinaryExpr{
			Op:  parser.ADD,
			LHS: &NumberLiteral{Val: 1},
			RHS: &BinaryExpr{
				Op:  parser.DIV,
				LHS: &NumberLiteral{Val: 2},
				RHS: &ParenExpr{
					Expr: &BinaryExpr{
						Op:  parser.MUL,
						LHS: &NumberLiteral{Val: 3},
						RHS: &NumberLiteral{Val: 1},
					},
				},
			},
		},
	},
	// TODO add tests
}

func TestAST2Plan(t *testing.T) {
	for _, test := range ast2planTestCases {
		t.Run(test.input.String(), func(t *testing.T) {
			plan := NewLogicalPlan(&test.input)
			require.True(t, plan != nil, "could not convert AST to logical plan")
			require.Equal(t, test.expected, plan, "error on input '%s'", test.input.String())
		})
	}
}
