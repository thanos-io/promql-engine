package physicalplan

import (
	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/planner/cost"
)

type PhysicalPlan interface {
	SetChildren(children []PhysicalPlan) // set child implementations, also update the operator and cost.
	Children() []PhysicalPlan            // Return the saved child implementations
	Operator() model.VectorOperator      // Return the saved physical operator
	Cost() cost.Cost                     // Return the saved cost
}
