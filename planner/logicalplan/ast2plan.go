package logicalplan

import "github.com/thanos-io/promql-engine/parser"

func NewLogicalPlan(expr *parser.Expr) LogicalPlan {
	switch node := (*expr).(type) {
	case *parser.StepInvariantExpr:
		return &StepInvariantExpr{Expr: NewLogicalPlan(&node.Expr)}
	case *parser.VectorSelector:
		return &VectorSelector{
			Name:           node.Name,
			OriginalOffset: node.OriginalOffset,
			Offset:         node.Offset,
			Timestamp:      node.Timestamp,
			StartOrEnd:     node.StartOrEnd,
			LabelMatchers:  node.LabelMatchers,
		}
	case *parser.MatrixSelector:
		return &MatrixSelector{
			VectorSelector: NewLogicalPlan(&node.VectorSelector),
			Range:          node.Range,
		}
	case *parser.AggregateExpr:
		return &AggregateExpr{
			Op:       node.Op,
			Expr:     NewLogicalPlan(&node.Expr),
			Param:    NewLogicalPlan(&node.Param),
			Grouping: node.Grouping,
			Without:  node.Without,
		}
	case *parser.Call:
		var args []LogicalPlan
		for i := range node.Args {
			args = append(args, NewLogicalPlan(&node.Args[i]))
		}
		return &Call{
			Func: node.Func,
			Args: args,
		}
	case *parser.BinaryExpr:
		return &BinaryExpr{
			Op:             node.Op,
			LHS:            NewLogicalPlan(&node.LHS),
			RHS:            NewLogicalPlan(&node.RHS),
			VectorMatching: node.VectorMatching,
			ReturnBool:     node.ReturnBool,
		}
	case *parser.UnaryExpr:
		return &UnaryExpr{
			Op:   node.Op,
			Expr: NewLogicalPlan(&node.Expr),
		}
	case *parser.ParenExpr:
		return &ParenExpr{
			Expr: NewLogicalPlan(&node.Expr),
		}
	case *parser.SubqueryExpr:
		return &SubqueryExpr{
			Expr:           NewLogicalPlan(&node.Expr),
			Range:          node.Range,
			OriginalOffset: node.OriginalOffset,
			Offset:         node.Offset,
			Timestamp:      node.Timestamp,
			StartOrEnd:     node.StartOrEnd,
			Step:           node.Step,
		}
	// literal types
	case *parser.NumberLiteral:
		return &NumberLiteral{Val: node.Val}
	case *parser.StringLiteral:
		return &StringLiteral{Val: node.Val}
	}
	return nil // should never reach here
}
