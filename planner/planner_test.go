package planner

import (
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/stretchr/testify/require"
	model2 "github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/parser"
	"github.com/thanos-io/promql-engine/planner/cost"
	"github.com/thanos-io/promql-engine/planner/logicalplan"
	"github.com/thanos-io/promql-engine/planner/memo"
	"github.com/thanos-io/promql-engine/planner/physicalplan"
	"testing"
)

// transformation rules

/*
a simple transformation rule to inject injection_foo to label matchers of VectorSelector
The transformation result might be invalid, but it's only mean to test the planner, so it's okay
*/
type dummyMatcherInjection struct{}

func (s *dummyMatcherInjection) HashCode() uint64 {
	return 1
}

func (s *dummyMatcherInjection) Match(expr *memo.GroupExpr) bool {
	_, ok := (expr.Expr).(*logicalplan.VectorSelector)
	return ok
}

func (s *dummyMatcherInjection) Transform(m memo.Memo, expr *memo.GroupExpr) *memo.GroupExpr {
	node, _ := (expr.Expr).(*logicalplan.VectorSelector)
	var newLabelMatchers []*labels.Matcher
	existingVals := make(map[string]bool)
	for _, matcher := range node.LabelMatchers {
		newLabelMatchers = append(newLabelMatchers, matcher)
		existingVals[matcher.Value] = true
	}
	toBeInjected := "injection_foo"
	if _, ok := existingVals[toBeInjected]; !ok {
		newLabelMatchers = append(newLabelMatchers, parser.MustLabelMatcher(labels.MatchEqual, model.MetricNameLabel, "injection_foo"))
		newNode := &logicalplan.VectorSelector{
			Name:           node.Name,
			OriginalOffset: node.OriginalOffset,
			Offset:         node.Offset,
			Timestamp:      node.Timestamp,
			StartOrEnd:     node.StartOrEnd,
			LabelMatchers:  newLabelMatchers,
		}
		newExpr := m.GetOrCreateGroupExpr(newNode)
		return newExpr
	} else {
		return expr // return the original
	}
}

// implementation rules

var implementationRules = []memo.ImplementationRule{ // implementation rule must cover all possible logical plan
	&mockVectorSelectorImplRule{},
	&mockBinaryExprImplRule{},
}

type mockVectorSelectorImplRule struct{}

func (m *mockVectorSelectorImplRule) ListImplementations(expr *memo.GroupExpr) []physicalplan.PhysicalPlan {
	if e, ok := (expr.Expr).(*logicalplan.VectorSelector); ok {
		return []physicalplan.PhysicalPlan{&mockVectorSelectorImpl{
			plan:   e,
			parent: expr,
		}}
	} else {
		return []physicalplan.PhysicalPlan{}
	}
}

type mockBinaryExprImplRule struct{}

func (m *mockBinaryExprImplRule) ListImplementations(expr *memo.GroupExpr) []physicalplan.PhysicalPlan {
	if e, ok := (expr.Expr).(*logicalplan.BinaryExpr); ok {
		return []physicalplan.PhysicalPlan{&mockBinaryExprImpl{
			plan:   e,
			parent: expr,
		}}
	} else {
		return []physicalplan.PhysicalPlan{}
	}
}

// the implementations

type mockVectorSelectorImpl struct {
	plan     *logicalplan.VectorSelector
	parent   *memo.GroupExpr
	children []physicalplan.PhysicalPlan
	cost     cost.Cost
}

func (m *mockVectorSelectorImpl) ParentExpr() *memo.GroupExpr {
	return m.parent
}

func (m *mockVectorSelectorImpl) SetChildren(children []physicalplan.PhysicalPlan) {
	m.children = children
	matcherValues := make(map[string]bool)
	for _, matcher := range m.plan.LabelMatchers {
		matcherValues[matcher.Value] = true
	}
	// we will bias the cost if there's a label with value "injection_foo" (because our transformation rule above)
	if _, ok := matcherValues["injection_foo"]; ok {
		m.cost = cost.Cost{
			CpuCost:    1,
			MemoryCost: 2,
		}
	} else {
		m.cost = cost.Cost{
			CpuCost:    4,
			MemoryCost: 3,
		}
	}
}

func (m *mockVectorSelectorImpl) Children() []physicalplan.PhysicalPlan {
	return m.children
}

func (m *mockVectorSelectorImpl) Operator() model2.VectorOperator {
	// FIXME this is just a test code to demo, not the real implementation, so I will return nil instead
	return nil
}

func (m *mockVectorSelectorImpl) Cost() cost.Cost {
	return m.cost
}

type mockBinaryExprImpl struct {
	plan     *logicalplan.BinaryExpr
	parent   *memo.GroupExpr
	children []physicalplan.PhysicalPlan
	cost     cost.Cost
}

func (m *mockBinaryExprImpl) SetChildren(children []physicalplan.PhysicalPlan) {
	m.children = children
	// this node always have 2 children
	leftChild := children[0]
	rightChild := children[1]
	// for simplicity, just combine their cost together (might be it's not true in real scenario, but it's just a demo)
	m.cost = cost.Cost{
		CpuCost:    leftChild.Cost().CpuCost + rightChild.Cost().CpuCost,
		MemoryCost: leftChild.Cost().MemoryCost + rightChild.Cost().MemoryCost,
	}
}

func (m *mockBinaryExprImpl) Children() []physicalplan.PhysicalPlan {
	return m.children
}

func (m *mockBinaryExprImpl) Operator() model2.VectorOperator {
	// FIXME this is just a test code to demo, not the real implementation, so I will return nil instead
	return nil
}

func (m *mockBinaryExprImpl) Cost() cost.Cost {
	return m.cost
}

// cost model

type mockCostModel struct{}

func (m *mockCostModel) IsBetter(currentCost cost.Cost, newCost cost.Cost) bool {
	if currentCost.CpuCost == newCost.CpuCost {
		return currentCost.MemoryCost > newCost.MemoryCost
	} else {
		return currentCost.CpuCost > newCost.CpuCost
	}
}

func TestPlanner0(t *testing.T) { // demo the new planner
	input := "foo / on(test,blub) group_left(bar) bar" // copied from parse_test.go
	/*
		&parser.BinaryExpr{
			Op: parser.DIV,
			LHS: &parser.VectorSelector{
				Name: "foo",
				LabelMatchers: []*labels.Matcher{
					parser.MustLabelMatcher(labels.MatchEqual, model.MetricNameLabel, "foo"),
				},
				PosRange: parser.PositionRange{
					Start: 0,
					End:   3,
				},
			},
			RHS: &parser.VectorSelector{
				Name: "bar",
				LabelMatchers: []*labels.Matcher{
					parser.MustLabelMatcher(labels.MatchEqual, model.MetricNameLabel, "bar"),
				},
				PosRange: parser.PositionRange{
					Start: 43,
					End:   46,
				},
			},
			VectorMatching: &parser.VectorMatching{
				Card:           parser.CardManyToOne,
				MatchingLabels: []string{"test", "blub"},
				Include:        []string{"blub"},
			},
		}
	*/
	expr, err := parser.ParseExpr(input)
	require.NoError(t, err)
	planner := New()
	planner.MakeRoot(expr)
	planner.Explore([]memo.TransformationRule{&dummyMatcherInjection{}}, 0) // now we have a dummy
	planner.FindBestImplementation(&mockCostModel{}, implementationRules)
	root := planner.rootGroup
	/*
		In the parsed expression, we have 2 VectorSelector expression,
		hence the logical plan also includes 2 VectorSelector nodes.

		Since we have a transformation rule `dummyMatcherInjection` to inject the "injection_foo" value into the matchers,
		and we have the implementation rules to bias the cost if "injection_foo" is presented
		({cpu_cost=1, mem_cost=2} if "injection_foo" is presented, {cpu_cost=4, mem_cost=3} otherwise).

		And the binary implementation rule will just sum up all of its child implementation cost.

		So the final cost is {cpu_cost=2, mem_cost=4}
	*/
	require.Equal(t, root.Implementation.Cost, cost.Cost{
		CpuCost:    2,
		MemoryCost: 4,
	})
}
