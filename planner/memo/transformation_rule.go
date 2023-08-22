package memo

import (
	"github.com/thanos-io/promql-engine/planner/utils"
)

type TransformationRule interface {
	utils.Hashable
	Match(expr *GroupExpr) bool                      // Check if the transformation can be applied to the expression
	Transform(memo Memo, expr *GroupExpr) *GroupExpr // Transform the expression
}
