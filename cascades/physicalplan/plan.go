package physicalplan

import (
	"github.com/thanos-io/promql-engine/cascades/cost"
	"github.com/thanos-io/promql-engine/execution/model"
)

type Implementation interface {
	CalculateCost(children []Implementation) cost.Cost // Calculate cost based on provided child implementations, also update the actual implementation, cost, and children list
	Operator() model.VectorOperator                    // Return the saved physical operator set from the last CalculateCost
	Cost() cost.Cost                                   // Return the saved cost from the last CalculateCost call.
	Children() []Implementation                        // Return the saved child implementations from the last CalculateCost call.
}
