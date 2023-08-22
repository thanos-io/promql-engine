package memo

import (
	"github.com/stretchr/testify/assert"
	"github.com/thanos-io/promql-engine/parser"
	"github.com/thanos-io/promql-engine/planner/logicalplan"
	"golang.org/x/exp/maps"
	"reflect"
	"testing"
)

var memoInitRootTestCases = []struct {
	input    logicalplan.LogicalPlan // The initial logical plan.
	expected *Group                  // The expected group.
}{
	{
		input: &logicalplan.BinaryExpr{
			Op:  parser.ADD,
			LHS: &logicalplan.NumberLiteral{Val: 1},
			RHS: &logicalplan.BinaryExpr{
				Op:  parser.DIV,
				LHS: &logicalplan.NumberLiteral{Val: 2},
				RHS: &logicalplan.ParenExpr{
					Expr: &logicalplan.BinaryExpr{
						Op:  parser.MUL,
						LHS: &logicalplan.NumberLiteral{Val: 3},
						RHS: &logicalplan.NumberLiteral{Val: 1},
					},
				},
			},
		},
		expected: &Group{
			Equivalents: map[ID]*GroupExpr{
				0: {
					Expr: &logicalplan.BinaryExpr{ /* #1 */
						Op:  parser.ADD,
						LHS: &logicalplan.NumberLiteral{Val: 1}, /* #2 */
						RHS: &logicalplan.BinaryExpr{ /* #3 */
							Op:  parser.DIV,
							LHS: &logicalplan.NumberLiteral{Val: 2}, /* #4 */
							RHS: &logicalplan.ParenExpr{ /* #5 */
								Expr: &logicalplan.BinaryExpr{ /* #6 */
									Op:  parser.MUL,
									LHS: &logicalplan.NumberLiteral{Val: 3}, // #7 */
									RHS: &logicalplan.NumberLiteral{Val: 1}, // #8 */
								},
							},
						},
					},
					Children: []*Group{
						{
							Equivalents: map[ID]*GroupExpr{
								0: {Expr: &logicalplan.NumberLiteral{Val: 1}}, /* #2 */
							},
						},
						{
							Equivalents: map[ID]*GroupExpr{
								0: {
									Expr: &logicalplan.BinaryExpr{ /* #3 */
										Op:  parser.DIV,
										LHS: &logicalplan.NumberLiteral{Val: 2}, /* #4 */
										RHS: &logicalplan.ParenExpr{ /* #5 */
											Expr: &logicalplan.BinaryExpr{ /* #6 */
												Op:  parser.MUL,
												LHS: &logicalplan.NumberLiteral{Val: 3}, /* #7 */
												RHS: &logicalplan.NumberLiteral{Val: 1}, /* #8 */
											},
										},
									},
									Children: []*Group{
										{
											Equivalents: map[ID]*GroupExpr{
												0: {Expr: &logicalplan.NumberLiteral{Val: 2}}, /* #4 */
											},
										},
										{
											Equivalents: map[ID]*GroupExpr{
												0: {
													Expr: &logicalplan.ParenExpr{ /* #5 */
														Expr: &logicalplan.BinaryExpr{ /* #6 */
															Op:  parser.MUL,
															LHS: &logicalplan.NumberLiteral{Val: 3}, /* #7 */
															RHS: &logicalplan.NumberLiteral{Val: 1}, /* #8 */
														},
													},
													Children: []*Group{
														{
															Equivalents: map[ID]*GroupExpr{
																0: {
																	Expr: &logicalplan.BinaryExpr{ /* #6 */
																		Op:  parser.MUL,
																		LHS: &logicalplan.NumberLiteral{Val: 3}, /* #7 */
																		RHS: &logicalplan.NumberLiteral{Val: 1}, /* #8 */
																	},
																	Children: []*Group{
																		{
																			Equivalents: map[ID]*GroupExpr{
																				0: {Expr: &logicalplan.NumberLiteral{Val: 3}}, /* #7 */
																			},
																		},
																		{
																			Equivalents: map[ID]*GroupExpr{
																				0: {Expr: &logicalplan.NumberLiteral{Val: 1}}, /* #8 */
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
	// TODO add tests
}

// function to check if two root group are logically equals
func checkGroupEquals(left *Group, right *Group) bool {
	if left == right {
		return true
	}
	// since we're testing root group, all group only contains exactly one group expr
	if len(left.Equivalents) != len(right.Equivalents) || len(left.Equivalents) != 1 || len(right.Equivalents) != 1 {
		return false
	}
	leftExpr := maps.Values(left.Equivalents)[0]
	rightExpr := maps.Values(right.Equivalents)[0]
	return checkGroupExprEquals(leftExpr, rightExpr)
}

// function to check if two group exprs are logically equals
func checkGroupExprEquals(left *GroupExpr, right *GroupExpr) bool {
	if left == right {
		return true
	}
	if !reflect.DeepEqual(left.Expr, right.Expr) {
		return false
	}
	if len(left.Children) != len(right.Children) {
		return false
	}
	length := len(left.Children)
	for i := 0; i < length; i++ {
		if !checkGroupEquals(left.Children[i], right.Children[i]) {
			return false
		}
	}
	return true
}

func TestMemoInitRootGroup(t *testing.T) {
	for _, test := range memoInitRootTestCases {
		m := NewMemo()
		root := m.GetOrCreateGroup(test.input)
		assert.True(t, checkGroupEquals(root, test.expected), "error on input '%s'", test.input)
	}
}
