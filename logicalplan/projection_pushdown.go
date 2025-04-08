package logicalplan

import (
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/util/annotations"
	"github.com/thanos-io/promql-engine/query"
)

type ProjectionPushdown struct {
	seriesHashLabel string
}

func (p ProjectionPushdown) Optimize(plan Node, _ *query.Options) (Node, annotations.Annotations) {
	// Single pass: top-down traversal to push projections directly
	p.pushProjection(&plan, nil, false)
	return plan, nil
}

// pushProjection recursively traverses the tree and pushes projection information down
// - requiredLabels: the set of labels required by parent nodes
// - isWithout: whether the projection should exclude (true) or include (false) the labels
func (p ProjectionPushdown) pushProjection(node *Node, requiredLabels map[string]struct{}, isWithout bool) {
	switch n := (*node).(type) {
	case *VectorSelector:
		// Apply projection if we have required labels
		if requiredLabels != nil {
			projection := Projection{
				Labels:  make([]string, 0, len(requiredLabels)),
				Include: !isWithout, // For "without", we exclude the labels
			}
			for lbl := range requiredLabels {
				projection.Labels = append(projection.Labels, lbl)
			}
			n.Projection = projection
		}

	case *Aggregation:

		// Special handling for aggregation functions that need all labels
		// regardless of grouping (topk, bottomk, limitk, limit_ratio)
		switch n.Op {
		case parser.TOPK, parser.BOTTOMK, parser.LIMITK, parser.LIMIT_RATIO:
			// These functions need all labels, so clear any requirements
			p.pushProjection(&n.Expr, nil, false)
			return
		}

		// For aggregations, we directly use the grouping labels
		grouping := n.Grouping

		groupingLabels := stringSet(grouping)

		// Propagate to children using the aggregation's own grouping requirements
		p.pushProjection(&n.Expr, groupingLabels, n.Without)

		if p.seriesHashLabel != "" && n.Without {
			n.Grouping = append(grouping, p.seriesHashLabel)
		}

	case *Binary:
		// For binary operations with vector matching, we need the matching labels
		if n.VectorMatching != nil {
			if n.VectorMatching.Card == parser.CardOneToOne {
				if n.VectorMatching.On {
					// For "on", we need to include only the matching labels
					// Don't consider parent requirements for binary operations
					onLabels := make(map[string]struct{})

					// Add the matching labels
					for _, lbl := range n.VectorMatching.MatchingLabels {
						onLabels[lbl] = struct{}{}
					}

					// Propagate to children
					for _, child := range n.Children() {
						p.pushProjection(child, onLabels, false) // Always use include mode for "on"
					}
					return // Already propagated to children
				} else {
					// For "ignoring", we need to exclude only the matching labels
					// Don't consider parent requirements for binary operations
					ignoredLabels := make(map[string]struct{})
					for _, lbl := range n.VectorMatching.MatchingLabels {
						ignoredLabels[lbl] = struct{}{}
					}

					// Propagate to children
					for _, child := range n.Children() {
						p.pushProjection(child, ignoredLabels, true) // true for "without"
					}
					n.VectorMatching.MatchingLabels = append(n.VectorMatching.MatchingLabels, p.seriesHashLabel)
					return // Already propagated to children
				}
			}
		}

		// No vector matching, just propagate existing requirements
		for _, child := range n.Children() {
			p.pushProjection(child, nil, false)
		}

	case *FunctionCall:
		// Check function requirements for labels
		updatedLabels := getFunctionLabelRequirements(n.Func.Name, n.Args, requiredLabels, isWithout)
		for _, child := range n.Children() {
			p.pushProjection(child, updatedLabels, isWithout)
		}

	case *MatrixSelector:
		// Push projection to the inner vector selector
		var vs Node = n.VectorSelector
		p.pushProjection(&vs, requiredLabels, isWithout)

	case *Subquery:
		// Push projection to the inner expression
		p.pushProjection(&n.Expr, requiredLabels, isWithout)

	default:
		// For other node types, propagate to children
		for _, child := range (*node).Children() {
			p.pushProjection(child, requiredLabels, isWithout)
		}
	}
}

func stringSet(s []string) map[string]struct{} {
	set := make(map[string]struct{}, len(s))
	for _, v := range s {
		set[v] = struct{}{}
	}
	return set
}

// unwrapStepInvariantExpr recursively unwraps step invariant expressions to get to the underlying node
func unwrapStepInvariantExpr(node Node) Node {
	if stepInvariant, ok := node.(*StepInvariantExpr); ok {
		return unwrapStepInvariantExpr(stepInvariant.Expr)
	}
	return node
}

// getFunctionLabelRequirements ensures that specific labels required by functions are included
// in the projection
func getFunctionLabelRequirements(funcName string, args []Node, requiredLabels map[string]struct{}, isWithout bool) map[string]struct{} {
	if requiredLabels == nil {
		return nil
	}

	result := make(map[string]struct{}, len(requiredLabels))
	for k, v := range requiredLabels {
		result[k] = v
	}

	// Add function-specific required labels
	switch funcName {
	case "histogram_quantile":
		// histogram_quantile requires the "le" label
		result["le"] = struct{}{}
	case "label_replace":
		// label_replace(v instant-vector, dst_label string, replacement string, src_label string, regex string)
		if len(args) >= 4 {
			// Check if the destination label is in the required labels
			if len(args) >= 2 {
				// Unwrap any step invariant expressions
				dstArg := unwrapStepInvariantExpr(args[1])
				if dstLit, ok := dstArg.(*StringLiteral); ok {
					dstLabel := dstLit.Val
					if !isWithout {
						if _, needed := result[dstLabel]; needed {
							// Only if the destination label is needed, we need the source label
							// Unwrap any step invariant expressions
							srcArg := unwrapStepInvariantExpr(args[3])
							if strLit, ok := srcArg.(*StringLiteral); ok {
								// Add the source label to required labels
								result[strLit.Val] = struct{}{}
							}
						}
					} else {
						if _, needed := result[dstLabel]; !needed {
							// Only if the destination label is needed, we need the source labels
							for i := 3; i < len(args); i++ {
								// Unwrap any step invariant expressions
								srcArg := unwrapStepInvariantExpr(args[i])
								if strLit, ok := srcArg.(*StringLiteral); ok {
									delete(result, strLit.Val)
								}
							}
						}
					}
				}
			}
		}
	case "label_join":
		// label_join(v instant-vector, dst_label string, separator string, src_label_1 string, src_label_2 string, ...)
		if len(args) >= 4 {
			// Check if the destination label is in the required labels
			if len(args) >= 2 {
				// Unwrap any step invariant expressions
				dstArg := unwrapStepInvariantExpr(args[1])
				if dstLit, ok := dstArg.(*StringLiteral); ok {
					dstLabel := dstLit.Val
					if !isWithout {
						if _, needed := result[dstLabel]; needed {
							// Only if the destination label is needed, we need the source labels
							for i := 3; i < len(args); i++ {
								// Unwrap any step invariant expressions
								srcArg := unwrapStepInvariantExpr(args[i])
								if strLit, ok := srcArg.(*StringLiteral); ok {
									// Add each source label to required labels
									result[strLit.Val] = struct{}{}
								}
							}
						}
					} else {
						if _, needed := result[dstLabel]; !needed {
							// Only if the destination label is needed, we need the source labels
							for i := 3; i < len(args); i++ {
								// Unwrap any step invariant expressions
								srcArg := unwrapStepInvariantExpr(args[i])
								if strLit, ok := srcArg.(*StringLiteral); ok {
									delete(result, strLit.Val)
								}
							}
						}
					}
				}
			}
		}
	}

	return result
}
