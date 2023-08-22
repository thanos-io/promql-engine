package planner

import (
	"github.com/thanos-io/promql-engine/parser"
	"github.com/thanos-io/promql-engine/planner/cost"
	"github.com/thanos-io/promql-engine/planner/logicalplan"
	"github.com/thanos-io/promql-engine/planner/memo"
	"github.com/thanos-io/promql-engine/planner/physicalplan"
	"golang.org/x/exp/maps"
)

type Planner struct {
	expr parser.Expr
	memo memo.Memo
	// root
	root      logicalplan.LogicalPlan
	rootGroup *memo.Group
}

func New() *Planner {
	return &Planner{
		memo: memo.NewMemo(),
	}
}

func (o *Planner) MakeRoot(expr parser.Expr) {
	o.expr = expr
	o.root = logicalplan.NewLogicalPlan(&expr)
	o.rootGroup = o.memo.GetOrCreateGroup(o.root)
}

func (o *Planner) exploreGroup(rules []memo.TransformationRule, group *memo.Group, round memo.ExplorationRound) {
	for {
		if group.IsExplored(round) {
			break
		}
		group.SetExplore(round, true)
		for _, equivalentExpr := range maps.Values(group.Equivalents) {
			// if the equivalent expression is not yet explored, then we will explore it
			if !equivalentExpr.IsExplored(round) {
				equivalentExpr.SetExplore(round, true)
				for _, child := range equivalentExpr.Children {
					o.exploreGroup(rules, child, round)
					if equivalentExpr.IsExplored(round) && child.IsExplored(round) {
						equivalentExpr.SetExplore(round, true)
					} else {
						equivalentExpr.SetExplore(round, false)
					}
				}
			}
			// fire rules for more equivalent expressions
			for _, rule := range rules {
				if rule.Match(equivalentExpr) {
					if !equivalentExpr.AppliedTransformations.Contains(rule) {
						transformedExpr := rule.Transform(o.memo, equivalentExpr)
						group.Equivalents[transformedExpr.ID] = transformedExpr
						equivalentExpr.AppliedTransformations.Add(rule)
						// reset group exploration state
						transformedExpr.SetExplore(round, false)
						group.SetExplore(round, false)
					}
				}
			}
			if group.IsExplored(round) && equivalentExpr.IsExplored(round) {
				group.SetExplore(round, true)
			} else {
				group.SetExplore(round, false)
			}
		}
	}
}

func (o *Planner) Explore(rules []memo.TransformationRule, round memo.ExplorationRound) {
	o.exploreGroup(rules, o.rootGroup, round)
}

func (o *Planner) findBestImpl(costModel cost.CostModel, rules []memo.ImplementationRule, group *memo.Group) *memo.GroupImplementation {
	if group.Implementation != nil {
		return group.Implementation
	} else {
		var groupImpl *memo.GroupImplementation
		for _, expr := range group.Equivalents {
			// fire rules to find implementations for each equiv expr, returning un-calculated implementations
			var possibleImpls []physicalplan.PhysicalPlan
			for _, rule := range rules {
				possibleImpls = append(possibleImpls, rule.ListImplementations(expr)...)
			}
			// get the implementation of child groups
			var childImpls []physicalplan.PhysicalPlan
			for _, child := range expr.Children {
				childImpl := o.findBestImpl(costModel, rules, child)
				child.Implementation = childImpl
				childImpls = append(childImpls, childImpl.Implementation)
			}
			// calculate the implementation, and update the best cost for group
			for _, impl := range possibleImpls {
				impl.SetChildren(childImpls)
				calculatedCost := impl.Cost()
				if groupImpl != nil {
					if costModel.IsBetter(groupImpl.Cost, calculatedCost) {
						groupImpl.SelectedExpr = expr
						groupImpl.Implementation = impl
						groupImpl.Cost = calculatedCost
					}
				} else {
					groupImpl = &memo.GroupImplementation{
						SelectedExpr:   expr,
						Cost:           calculatedCost,
						Implementation: impl,
					}
				}
			}
		}
		return groupImpl
	}
}

func (o *Planner) FindBestImplementation(costModel cost.CostModel, rules []memo.ImplementationRule) {
	o.rootGroup.Implementation = o.findBestImpl(costModel, rules, o.rootGroup)
}
