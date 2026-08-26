// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"slices"

	"github.com/thanos-io/promql-engine/query"

	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/util/annotations"
)

type ProjectionOptimizer struct {
	SeriesHashLabel string
}

func (p ProjectionOptimizer) Optimize(plan Node, _ *query.Options) (Node, annotations.Annotations) {
	p.pushProjection(&plan, nil)
	return plan, nil
}

// pushProjection recursively traverses the tree and pushes projection information down.
func (p ProjectionOptimizer) pushProjection(node *Node, projection *Projection) {
	switch n := (*node).(type) {
	case *VectorSelector:
		proj := &Projection{}
		if projection != nil {
			// Copy: the same projection is handed to sibling selectors.
			proj = &Projection{Labels: slices.Clone(projection.Labels), Include: projection.Include}
		}
		// MergeSelectsOptimizer may have replaced the matchers with a less
		// selective superset plus compensating filters, which are matched
		// against the series the storage returns, so the labels they read must
		// survive the projection.
		for _, f := range n.Filters {
			if proj.Include {
				if !slices.Contains(proj.Labels, f.Name) {
					proj.Labels = append(proj.Labels, f.Name)
				}
			} else {
				proj.Labels = slices.DeleteFunc(proj.Labels, func(s string) bool { return s == f.Name })
			}
		}
		n.Projection = proj

	case *Aggregation:
		// Special handling for aggregation functions that need all labels
		// regardless of grouping (topk, bottomk, limitk, limit_ratio)
		switch n.Op {
		case parser.TOPK, parser.BOTTOMK, parser.LIMITK, parser.LIMIT_RATIO:
			// These functions need all labels, so clear any requirements
			p.pushProjection(&n.Expr, nil)
			return
		}

		// For aggregations, we directly use the grouping labels
		grouping := n.Grouping
		groupingProjection := &Projection{
			Labels:  grouping,
			Include: !n.Without,
		}
		// Note that we don't push projection to Aggregation.Param as they are not
		// selecting data for the aggregation.
		p.pushProjection(&n.Expr, groupingProjection)

		if p.SeriesHashLabel != "" && n.Without {
			n.Grouping = append(grouping, p.SeriesHashLabel)
		}

	case *Binary:
		var highCard, lowCard = n.LHS, n.RHS

		// "or" passes series from both sides through with their label sets
		// unchanged, so trimming either side would corrupt the output.
		if n.Op == parser.LOR {
			p.pushProjection(&highCard, nil)
			p.pushProjection(&lowCard, nil)
			return
		}

		if n.VectorMatching == nil || (!n.VectorMatching.On && len(n.VectorMatching.MatchingLabels) == 0) {
			// A scalar operand takes part in no label matching, so the vector
			// side can keep the outer projection. With two vector operands and
			// default matching the signature is the whole label set, so
			// trimming either side would change which series match: constant
			// but vector-valued operands such as vector() or day_of_week()
			// match on their empty label set.
			if lowCard.ReturnType() == parser.ValueTypeScalar {
				p.pushProjection(&highCard, projection)
			} else {
				p.pushProjection(&highCard, nil)
			}

			if highCard.ReturnType() == parser.ValueTypeScalar {
				p.pushProjection(&lowCard, projection)
			} else {
				p.pushProjection(&lowCard, nil)
			}
			return
		}

		if n.VectorMatching.Card == parser.CardOneToOne {
			proj := &Projection{
				Labels:  n.VectorMatching.MatchingLabels,
				Include: n.VectorMatching.On,
			}

			for _, child := range n.Children() {
				p.pushProjection(child, proj)
			}

			if !n.VectorMatching.On && p.SeriesHashLabel != "" {
				n.VectorMatching.MatchingLabels = append(n.VectorMatching.MatchingLabels, p.SeriesHashLabel)
			}
			return
		}

		if n.VectorMatching.Card == parser.CardOneToMany {
			highCard, lowCard = lowCard, highCard
		}

		// Handle high card side projection. Only ignoring mode is supported.
		hcProjection := &Projection{}
		// Only push projection for high card side if there is an outer projection available
		// to remove series hash
		if projection != nil && projection.Include {
			// Include labels are from low card side so we don't need to fetch
			// them from high card side if include labels are not used as join keys.
			hcProjection.Labels = n.VectorMatching.Include
			if !n.VectorMatching.On {
				hcProjection.Labels = intersect(hcProjection.Labels, n.VectorMatching.MatchingLabels)
			}
		}
		if len(hcProjection.Labels) > 1 {
			p.pushProjection(&highCard, hcProjection)
		} else {
			// If there is only 1 label to project then it is not worth to push projection
			// down to high card side as calculating hash might be more expensive.
			p.pushProjection(&highCard, nil)
		}

		// Handle low card side projection.
		lcProjection := extendProjection(Projection{
			Include: n.VectorMatching.On,
			Labels:  n.VectorMatching.MatchingLabels,
		}, n.VectorMatching.Include)
		p.pushProjection(&lowCard, &lcProjection)

		if !n.VectorMatching.On && p.SeriesHashLabel != "" {
			n.VectorMatching.MatchingLabels = append(n.VectorMatching.MatchingLabels, p.SeriesHashLabel)
		}
		return

	case *FunctionCall:
		// Handle function-specific label requirements.
		updatedProjection := getFunctionLabelRequirements(n.Func.Name, n.Args, projection)
		for _, child := range n.Children() {
			p.pushProjection(child, updatedProjection)
		}

	default:
		// For other node types, propagate to children
		for _, child := range (*node).Children() {
			p.pushProjection(child, projection)
		}
	}
}

func extendProjection(projection Projection, lbls []string) Projection {
	var extendedLabels []string
	if projection.Include {
		extendedLabels = union(projection.Labels, lbls)
	} else {
		extendedLabels = subtract(projection.Labels, lbls)
	}
	return Projection{
		Include: projection.Include,
		Labels:  extendedLabels,
	}
}

// unwrapStepInvariantExpr recursively unwraps step invariant expressions to get to the underlying node.
func unwrapStepInvariantExpr(node Node) Node {
	if stepInvariant, ok := node.(*StepInvariantExpr); ok {
		return unwrapStepInvariantExpr(stepInvariant.Expr)
	}
	return node
}

// getFunctionLabelRequirements returns an updated projection based on function-specific requirements.
func getFunctionLabelRequirements(funcName string, args []Node, projection *Projection) *Projection {
	if projection == nil {
		projection = &Projection{}
	}
	result := &Projection{
		Labels:  make([]string, len(projection.Labels)),
		Include: projection.Include,
	}
	copy(result.Labels, projection.Labels)

	// Add function-specific required labels
	switch funcName {
	case "absent_over_time", "absent", "scalar":
		return &Projection{
			Labels:  []string{},
			Include: true,
		}
	case "histogram_quantile":
		// Unsafe to push projection down for histogram_quantile as it requires le label.
		return nil
	case "label_replace":
		dstArg := unwrapStepInvariantExpr(args[1])
		if dstLit, ok := dstArg.(*StringLiteral); ok {
			dstLabel := dstLit.Val
			needed := slices.Contains(result.Labels, dstLabel)
			needSourceLabels := (result.Include && needed) || (!result.Include && !needed)
			if !needSourceLabels {
				return result
			}

			srcArg := unwrapStepInvariantExpr(args[3])
			if strLit, ok := srcArg.(*StringLiteral); ok {
				if result.Include && needed {
					result.Labels = append(result.Labels, strLit.Val)
				} else {
					result.Labels = slices.DeleteFunc(result.Labels, func(s string) bool {
						return s == strLit.Val
					})
				}
			}
		}
	case "label_join":
		dstArg := unwrapStepInvariantExpr(args[1])
		if dstLit, ok := dstArg.(*StringLiteral); ok {
			dstLabel := dstLit.Val
			needed := slices.Contains(result.Labels, dstLabel)
			needSourceLabels := (result.Include && needed) || (!result.Include && !needed)
			if !needSourceLabels {
				return result
			}

			// Only if the destination label is needed, we need the source labels
			for i := 3; i < len(args); i++ {
				srcArg := unwrapStepInvariantExpr(args[i])
				if strLit, ok := srcArg.(*StringLiteral); ok {
					if result.Include && needed {
						result.Labels = append(result.Labels, strLit.Val)
					} else {
						result.Labels = slices.DeleteFunc(result.Labels, func(s string) bool {
							return s == strLit.Val
						})
					}
				}
			}
		}
	}

	return result
}

// union returns the sorted union of two string slices.
func union(l1 []string, l2 []string) []string {
	res := make([]string, 0, len(l1)+len(l2))
	res = append(res, l1...)
	res = append(res, l2...)
	slices.Sort(res)
	return slices.Compact(res)
}

// subtract returns l1 minus l2, sorted.
func subtract(l1 []string, l2 []string) []string {
	res := make([]string, 0, len(l1))
	for _, s := range l1 {
		if !slices.Contains(l2, s) {
			res = append(res, s)
		}
	}
	slices.Sort(res)
	return slices.Compact(res)
}

func intersect(l1 []string, l2 []string) []string {
	m := make(map[string]struct{})
	for _, s := range l1 {
		m[s] = struct{}{}
	}
	var result []string
	for _, s := range l2 {
		if _, ok := m[s]; ok {
			result = append(result, s)
		}
	}
	return result
}
