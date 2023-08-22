package memo

import (
	"github.com/thanos-io/promql-engine/cascades/logicalplan"
	"github.com/thanos-io/promql-engine/cascades/utils"
)

type Memo interface {
	GetOrCreateGroupExpr(node logicalplan.LogicalPlan) *GroupExpr
	GetOrCreateGroup(node logicalplan.LogicalPlan) *Group
}

type memo struct {
	Groups               map[ID]*Group                          // The ID-Group mapping, used to store all the groups.
	Parents              map[ID]*Group                          // The GroupExpr-Group mapping, mapped from expr ID to group, used to find the group containing the equivalent logical plans.
	GroupExprs           map[logicalplan.LogicalPlan]*GroupExpr // The LogicalPlan-GroupExpr mapping.
	groupIDGenerator     utils.Generator[ID]                    // The ID generator for groups
	groupExprIDGenerator utils.Generator[ID]                    // The ID generator for group exprs
}

func NewMemo() Memo {
	return &memo{
		Groups:               make(map[ID]*Group),
		Parents:              make(map[ID]*Group),
		GroupExprs:           make(map[logicalplan.LogicalPlan]*GroupExpr),
		groupIDGenerator:     NewIDGenerator(),
		groupExprIDGenerator: NewIDGenerator(),
	}
}

func (m *memo) GetOrCreateGroupExpr(node logicalplan.LogicalPlan) *GroupExpr {
	children := node.Children()
	var childGroups []*Group
	for _, child := range children {
		childGroups = append(childGroups, m.GetOrCreateGroup(child))
	}
	entry, ok := m.GroupExprs[node]
	if ok {
		return entry
	} else {
		id := m.groupExprIDGenerator.Generate()
		expr := &GroupExpr{
			ID:       id,
			Expr:     node,
			Children: childGroups,
		}
		m.GroupExprs[node] = expr
		return expr
	}
}

func (m *memo) GetOrCreateGroup(node logicalplan.LogicalPlan) *Group {
	groupExpr := m.GetOrCreateGroupExpr(node)
	entry, ok := m.Parents[groupExpr.ID]
	if ok {
		entry.Equivalents[groupExpr.ID] = groupExpr
		return entry
	} else {
		id := m.groupIDGenerator.Generate()
		group := &Group{
			ID:          id,
			Equivalents: map[ID]*GroupExpr{groupExpr.ID: groupExpr},
		}
		m.Groups[group.ID] = group
		m.Parents[groupExpr.ID] = group
		return group
	}
}
