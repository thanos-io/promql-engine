package cascades

import (
	"github.com/thanos-io/promql-engine/cascades/cost"
	"github.com/thanos-io/promql-engine/cascades/logicalplan"
	"github.com/thanos-io/promql-engine/cascades/memo"
	"github.com/thanos-io/promql-engine/cascades/physicalplan"
	"github.com/thanos-io/promql-engine/parser"
	"golang.org/x/exp/maps"
)

type Optimize struct {
	expr parser.Expr
	memo memo.Memo
	// root
	root      logicalplan.LogicalPlan
	rootGroup *memo.Group
}

func New(expr parser.Expr) *Optimize {
	return &Optimize{
		expr: expr,
		memo: memo.NewMemo(),
	}
}

func (o *Optimize) SetRoot(root logicalplan.LogicalPlan) {
	o.root = root
	o.rootGroup = o.memo.GetOrCreateGroup(root)
}

func (o *Optimize) exploreGroup(rules []memo.TransformationRule, group *memo.Group, round memo.ExplorationRound) {
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
				if !equivalentExpr.AppliedTransformations.Contains(rule) {
					if rule.Match(equivalentExpr) {
						transformedExpr := rule.Transform(equivalentExpr)
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

func (o *Optimize) Explore(rules []memo.TransformationRule, round memo.ExplorationRound) {
	o.exploreGroup(rules, o.rootGroup, round)
}

func (o *Optimize) findBestImpl(costModel cost.CostModel, rules []memo.ImplementationRule, group *memo.Group) *memo.GroupImplementation {
	if group.Implementation != nil {
		return group.Implementation
	} else {
		var groupImpl *memo.GroupImplementation
		for _, expr := range group.Equivalents {
			// fire rules to find implementations for each equiv expr, returning un-calculated implementations
			var possibleImpls []physicalplan.Implementation
			for _, rule := range rules {
				possibleImpls = append(possibleImpls, rule.ListImplementations(expr)...)
			}
			// get the implementation of child groups
			var childImpls []physicalplan.Implementation
			for _, child := range expr.Children {
				childImpl := o.findBestImpl(costModel, rules, child)
				child.Implementation = childImpl
				childImpls = append(childImpls, childImpl.Implementation)
			}
			// calculate the implementation, and update the best cost for group
			for _, impl := range possibleImpls {
				calculatedCost := impl.CalculateCost(childImpls)
				if groupImpl != nil {
					if costModel.IsBetter(groupImpl, calculatedCost) {
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

func (o *Optimize) FindBestImplementation(costModel cost.CostModel, rules []memo.ImplementationRule) {
	o.rootGroup.Implementation = o.findBestImpl(costModel, rules, o.rootGroup)
}
