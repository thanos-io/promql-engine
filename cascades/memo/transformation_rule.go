package memo

import (
	"github.com/thanos-io/promql-engine/cascades/utils"
)

type TransformationRule interface {
	utils.Hashable
	Match(expr *GroupExpr) bool           // Check if the transformation can be applied to the expression
	Transform(expr *GroupExpr) *GroupExpr // Transform the expression
}
