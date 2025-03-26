package logicalplan

import (
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/util/annotations"
	"github.com/thanos-io/promql-engine/query"
)

type ProjectionPushdown struct{}

func (p ProjectionPushdown) Optimize(plan Node, _ *query.Options) (Node, annotations.Annotations) {
	// Single pass: top-down traversal to push projections directly
	pushProjection(&plan, nil, false)
	plan = insertDuplicateLabelChecks(plan)
	return insertRemoveSeriesHash(plan), nil
}

func insertRemoveSeriesHash(expr Node) Node {
	Traverse(&expr, func(node *Node) {
		switch t := (*node).(type) {
		case *VectorSelector:
			*node = &RemoveSeriesHash{Expr: t}
		}
	})
	return expr
}

// pushProjection recursively traverses the tree and pushes projection information down
// - requiredLabels: the set of labels required by parent nodes
// - isWithout: whether the projection should exclude (true) or include (false) the labels
func pushProjection(node *Node, requiredLabels map[string]struct{}, isWithout bool) {
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
			for _, child := range n.Children() {
				pushProjection(child, nil, false)
			}
			return
		}

		// For aggregations, we directly use the grouping labels
		grouping := n.Grouping
		if n.Without {
			grouping = append(grouping, labels.MetricName)
		}
		groupingLabels := stringSet(grouping)

		// Propagate to children using the aggregation's own grouping requirements
		for _, child := range n.Children() {
			pushProjection(child, groupingLabels, n.Without)
		}
		//if n.Without {
		//	n.Grouping = append(n.Grouping, "__series_hash__")
		//}
		return

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
						pushProjection(child, onLabels, false) // Always use include mode for "on"
					}
					return // Already propagated to children
				} else {
					// For "ignoring", we need to exclude only the matching labels
					// Don't consider parent requirements for binary operations
					ignoredLabels := make(map[string]struct{})
					for _, lbl := range n.VectorMatching.MatchingLabels {
						ignoredLabels[lbl] = struct{}{}
					}

					//// Also ignore the metric name label for "ignoring" mode
					//if len(n.VectorMatching.MatchingLabels) > 0 && !(n.Op == parser.LAND || n.Op == parser.LOR || n.Op == parser.LUNLESS) {
					//	ignoredLabels[labels.MetricName] = struct{}{}
					//}

					// Propagate to children
					for _, child := range n.Children() {
						pushProjection(child, ignoredLabels, true) // true for "without"
					}
					//n.VectorMatching.MatchingLabels = append(n.VectorMatching.MatchingLabels, "__series_hash__")
					return // Already propagated to children
				}
			} else {
				// For group_left/group_right with "on", we need matching labels and include labels
				if n.VectorMatching.On {
					// Don't consider parent requirements for binary operations
					for i, child := range n.Children() {
						childRequired := make(map[string]struct{})

						// Add the matching labels
						for _, lbl := range n.VectorMatching.MatchingLabels {
							childRequired[lbl] = struct{}{}
						}

						// For group_left, only the left side (i==0) needs the include labels
						// For group_right, only the right side (i==1) needs the include labels
						if (n.VectorMatching.Card == parser.CardManyToOne && i == 0) ||
							(n.VectorMatching.Card == parser.CardOneToMany && i == 1) {
							for _, lbl := range n.VectorMatching.Include {
								childRequired[lbl] = struct{}{}
							}
						}

						pushProjection(child, childRequired, false) // Always use include mode for "on"
					}
					return // Already propagated to children
				} else {
					// For "ignoring" with group_left/group_right
					for i, child := range n.Children() {
						// Don't consider parent requirements for binary operations
						ignoredLabels := make(map[string]struct{})
						for _, lbl := range n.VectorMatching.MatchingLabels {
							ignoredLabels[lbl] = struct{}{}
						}

						//// Also ignore the metric name label for "ignoring" mode
						//if len(n.VectorMatching.MatchingLabels) > 0 && !(n.Op == parser.LAND || n.Op == parser.LOR || n.Op == parser.LUNLESS) {
						//	ignoredLabels[labels.MetricName] = struct{}{}
						//}

						// For the many-to-one side, don't ignore include labels
						if (n.VectorMatching.Card == parser.CardManyToOne && i == 0) ||
							(n.VectorMatching.Card == parser.CardOneToMany && i == 1) {
							for _, lbl := range n.VectorMatching.Include {
								delete(ignoredLabels, lbl)
							}
						}

						pushProjection(child, ignoredLabels, true) // true for "without"
					}
					//n.VectorMatching.MatchingLabels = append(n.VectorMatching.MatchingLabels, "__series_hash__")
					return // Already propagated to children
				}
			}
		}

		// No vector matching, just propagate existing requirements
		for _, child := range n.Children() {
			pushProjection(child, requiredLabels, isWithout)
		}

	case *CheckDuplicateLabels:
		// These need all labels, so clear any requirements
		for _, child := range n.Children() {
			pushProjection(child, nil, false)
		}

	case *FunctionCall:
		// Check function requirements for labels
		updatedLabels := getFunctionLabelRequirements(n.Func.Name, n.Args, requiredLabels)
		for _, child := range n.Children() {
			pushProjection(child, updatedLabels, isWithout)
		}

	case *MatrixSelector:
		// Push projection to the inner vector selector
		var vs Node = n.VectorSelector
		pushProjection(&vs, requiredLabels, isWithout)

	case *Subquery:
		// Push projection to the inner expression
		pushProjection(&n.Expr, requiredLabels, isWithout)

	default:
		// For other node types, propagate to children
		for _, child := range (*node).Children() {
			pushProjection(child, requiredLabels, isWithout)
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
func getFunctionLabelRequirements(funcName string, args []Node, requiredLabels map[string]struct{}) map[string]struct{} {
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
					if _, needed := result[dstLabel]; needed {
						// Only if the destination label is needed, we need the source label
						// Unwrap any step invariant expressions
						srcArg := unwrapStepInvariantExpr(args[3])
						if strLit, ok := srcArg.(*StringLiteral); ok {
							// Add the source label to required labels
							result[strLit.Val] = struct{}{}
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
				}
			}
		}
	}

	return result
}
