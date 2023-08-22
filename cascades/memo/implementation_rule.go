package memo

import "github.com/thanos-io/promql-engine/cascades/physicalplan"

type ImplementationRule interface {
	ListImplementations(expr *GroupExpr) []physicalplan.Implementation // List all implementation for the expression
}
