package memo

import "github.com/thanos-io/promql-engine/planner/physicalplan"

type ImplementationRule interface {
	ListImplementations(expr *GroupExpr) []physicalplan.PhysicalPlan // List all implementation for the expression
}
