package memo

import (
	"github.com/thanos-io/promql-engine/planner/cost"
	"github.com/thanos-io/promql-engine/planner/logicalplan"
	"github.com/thanos-io/promql-engine/planner/physicalplan"
	"github.com/thanos-io/promql-engine/planner/utils"
	"sync/atomic"
)

// ID

type ID uint32

type idGenerator struct {
	counter uint64
}

func NewIDGenerator() utils.Generator[ID] {
	return &idGenerator{counter: 0}
}

func (g *idGenerator) Generate() ID {
	return ID(atomic.AddUint64(&g.counter, 1))
}

// Group

type Group struct {
	ID ID
	// logical
	Equivalents map[ID]*GroupExpr // The equivalent expressions.
	ExplorationMark
	// physical
	Implementation *GroupImplementation
}

type GroupImplementation struct {
	SelectedExpr   *GroupExpr
	Cost           cost.Cost
	Implementation physicalplan.PhysicalPlan
}

type GroupExpr struct {
	ID                     ID
	Expr                   logicalplan.LogicalPlan // The logical plan bind to the expression.
	Children               []*Group                // The children group of the expression, noted that it must be in the same order with LogicalPlan.Children().
	AppliedTransformations utils.Set[TransformationRule]
	ExplorationMark
}
