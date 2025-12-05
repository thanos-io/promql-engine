// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"github.com/thanos-io/promql-engine/query"

	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/util/annotations"
)

// PropagateCoalesceOptimizer shards selectors and moves Coalesce nodes as high
// as possible in the execution tree to maximize parallelism.
//
// The optimizer has two phases:
//  1. Shard vector and matrix selectors into Coalesce nodes based on
//     DecodingConcurrency. Each shard scans a portion of the series.
//  2. Propagate Coalesce nodes upward through parallelizable operations.
//
// Key insight: Since coalesce perfectly shards the vector (each shard has
// distinct series), distributive aggregations can be pushed through:
//
//	sum(coalesce(A, B)) = coalesce(sum(A), sum(B))
//
// This is the biggest performance win as aggregations are often expensive.
//
// For example, with 2 shards:
//
//	sum(rate(http_requests_total[5m]))
//
// Becomes:
//
//	sum(coalesce(
//	    sum(rate(http_requests_total shard 0 of 2[5m])),
//	    sum(rate(http_requests_total shard 1 of 2[5m]))
//	))
//
// The optimizer will NOT push coalesce through:
//   - Non-distributive aggregations (avg, stddev, etc.)
//   - Deduplicate - coalesce should stay inside dedup for correctness
//   - Binary operations with vectors on both sides - these need matching
//   - absent/absent_over_time functions - these have special semantics
type PropagateCoalesceOptimizer struct{}

func (p PropagateCoalesceOptimizer) Optimize(plan Node, opts *query.Options) (Node, annotations.Annotations) {
	// Phase 1: Shard selectors into Coalesce nodes
	if opts != nil && opts.DecodingConcurrency > 1 {
		plan = p.shardSelectors(plan, opts.DecodingConcurrency)
	}

	// Phase 2: Propagate Coalesce nodes upward
	changed := true
	for changed {
		plan, changed = p.propagateOnce(plan)
	}
	return plan, nil
}

// shardSelectors wraps vector and matrix selectors in Coalesce nodes.
func (p PropagateCoalesceOptimizer) shardSelectors(plan Node, numShards int) Node {
	TraverseBottomUp(nil, &plan, func(parent, current *Node) bool {
		// Skip selectors inside absent/absent_over_time functions.
		// These functions have special semantics that don't work with sharding.
		if parent != nil {
			if fn, ok := (*parent).(*FunctionCall); ok {
				if fn.Func.Name == "absent" || fn.Func.Name == "absent_over_time" {
					return false
				}
			}
		}

		switch node := (*current).(type) {
		case *VectorSelector:
			*current = shardVectorSelector(node, numShards)
			return false
		case *MatrixSelector:
			*current = shardMatrixSelector(node, numShards)
			return false
		}
		return false
	})
	return plan
}

func (p PropagateCoalesceOptimizer) propagateOnce(plan Node) (Node, bool) {
	changed := false

	TraverseBottomUp(nil, &plan, func(parent, current *Node) bool {
		if parent == nil || current == nil {
			return false
		}

		coalesce, ok := (*current).(*Coalesce)
		if !ok {
			return false
		}

		// Don't try to push if parent is something we can never push through
		// This prevents infinite loops when coalesce is already at its highest position
		switch p := (*parent).(type) {
		case Deduplicate:
			return false
		case *Aggregation:
			// Allow distributive aggregations and topk/bottomk
			if !isDistributiveAggregation(p.Op) && p.Op != parser.TOPK && p.Op != parser.BOTTOMK {
				return false
			}
		}

		// Try to push the coalesce through the parent
		newParent, pushed := p.tryPushThrough(parent, coalesce)
		if pushed {
			*parent = newParent
			changed = true
			return true // stop this traversal, we'll do another pass
		}

		return false
	})

	return plan, changed
}

// tryPushThrough attempts to push a coalesce node through its parent.
// Returns the new parent node and true if successful, or nil and false if not.
func (p PropagateCoalesceOptimizer) tryPushThrough(parent *Node, coalesce *Coalesce) (Node, bool) {
	switch parentNode := (*parent).(type) {
	case *FunctionCall:
		return p.pushThroughFunction(parentNode, coalesce)
	case *Binary:
		return p.pushThroughBinary(parentNode, coalesce)
	case *Unary:
		return p.pushThroughUnary(parentNode, coalesce)
	case *StepInvariantExpr:
		return p.pushThroughStepInvariant(parentNode, coalesce)
	case *CheckDuplicateLabels:
		return p.pushThroughCheckDuplicateLabels(parentNode, coalesce)
	case *Parens:
		return p.pushThroughParens(parentNode, coalesce)
	case Deduplicate:
		// Never push coalesce through Deduplicate - it needs all inputs
		// to properly deduplicate samples
		return nil, false
	case *Aggregation:
		return p.pushThroughAggregation(parentNode, coalesce)
	}
	return nil, false
}

// pushThroughFunction pushes coalesce through element-wise functions.
// These are functions that operate on each sample independently.
func (p PropagateCoalesceOptimizer) pushThroughFunction(fn *FunctionCall, coalesce *Coalesce) (Node, bool) {
	// Only element-wise functions can have coalesce pushed through them
	if !isElementWiseFunction(fn.Func.Name) {
		return nil, false
	}

	// Find which argument is the coalesce
	coalesceArgIdx := -1
	for i, arg := range fn.Args {
		if arg == coalesce {
			coalesceArgIdx = i
			break
		}
	}
	if coalesceArgIdx == -1 {
		return nil, false
	}

	// Create new expressions for each coalesce child
	newExprs := make([]Node, len(coalesce.Expressions))
	for i, expr := range coalesce.Expressions {
		// Clone the function and replace the coalesce argument with this expression
		newFn := fn.Clone().(*FunctionCall)
		newFn.Args[coalesceArgIdx] = expr.Clone()
		newExprs[i] = newFn
	}

	return &Coalesce{Expressions: newExprs}, true
}

// pushThroughBinary pushes coalesce through binary operations where
// one side is a scalar or constant expression.
func (p PropagateCoalesceOptimizer) pushThroughBinary(bin *Binary, coalesce *Coalesce) (Node, bool) {
	// Check if one side is a scalar/constant and the other is our coalesce
	lhsIsCoalesce := bin.LHS == coalesce
	rhsIsCoalesce := bin.RHS == coalesce

	if !lhsIsCoalesce && !rhsIsCoalesce {
		return nil, false
	}

	// We can push through if the other side is a constant scalar expression
	var otherSide Node
	if lhsIsCoalesce {
		otherSide = bin.RHS
	} else {
		otherSide = bin.LHS
	}

	if !IsConstantScalarExpr(otherSide) {
		return nil, false
	}

	// Create new expressions for each coalesce child
	newExprs := make([]Node, len(coalesce.Expressions))
	for i, expr := range coalesce.Expressions {
		newBin := bin.Clone().(*Binary)
		if lhsIsCoalesce {
			newBin.LHS = expr.Clone()
		} else {
			newBin.RHS = expr.Clone()
		}
		newExprs[i] = newBin
	}

	return &Coalesce{Expressions: newExprs}, true
}

// pushThroughUnary pushes coalesce through unary operations (negation).
func (p PropagateCoalesceOptimizer) pushThroughUnary(unary *Unary, coalesce *Coalesce) (Node, bool) {
	if unary.Expr != coalesce {
		return nil, false
	}

	newExprs := make([]Node, len(coalesce.Expressions))
	for i, expr := range coalesce.Expressions {
		newExprs[i] = &Unary{
			Op:   unary.Op,
			Expr: expr.Clone(),
		}
	}

	return &Coalesce{Expressions: newExprs}, true
}

// pushThroughStepInvariant pushes coalesce through step invariant expressions.
func (p PropagateCoalesceOptimizer) pushThroughStepInvariant(si *StepInvariantExpr, coalesce *Coalesce) (Node, bool) {
	if si.Expr != coalesce {
		return nil, false
	}

	newExprs := make([]Node, len(coalesce.Expressions))
	for i, expr := range coalesce.Expressions {
		newExprs[i] = &StepInvariantExpr{Expr: expr.Clone()}
	}

	return &Coalesce{Expressions: newExprs}, true
}

// pushThroughCheckDuplicateLabels pushes coalesce through duplicate label checks.
func (p PropagateCoalesceOptimizer) pushThroughCheckDuplicateLabels(cdl *CheckDuplicateLabels, coalesce *Coalesce) (Node, bool) {
	if cdl.Expr != coalesce {
		return nil, false
	}

	newExprs := make([]Node, len(coalesce.Expressions))
	for i, expr := range coalesce.Expressions {
		newExprs[i] = &CheckDuplicateLabels{Expr: expr.Clone()}
	}

	return &Coalesce{Expressions: newExprs}, true
}

// pushThroughParens pushes coalesce through parentheses.
func (p PropagateCoalesceOptimizer) pushThroughParens(parens *Parens, coalesce *Coalesce) (Node, bool) {
	if parens.Expr != coalesce {
		return nil, false
	}

	newExprs := make([]Node, len(coalesce.Expressions))
	for i, expr := range coalesce.Expressions {
		newExprs[i] = &Parens{Expr: expr.Clone()}
	}

	return &Coalesce{Expressions: newExprs}, true
}

// pushThroughAggregation pushes coalesce through distributive aggregations.
// Since coalesce perfectly shards the vector (each shard has distinct series),
// distributive aggregations can be computed on each shard independently:
//
//	sum(coalesce(A, B)) = sum(coalesce(sum(A), sum(B)))
//	min(coalesce(A, B)) = min(coalesce(min(A), min(B)))
//	max(coalesce(A, B)) = max(coalesce(max(A), max(B)))
//	count(coalesce(A, B)) = sum(coalesce(count(A), count(B)))  // Note: outer is SUM
//	group(coalesce(A, B)) = group(coalesce(group(A), group(B)))
//
// For topk/bottomk, we can compute local top/bottom K on each shard, then
// take the global top/bottom K from the combined results:
//
//	topk(K, coalesce(A, B)) = topk(K, coalesce(topk(K, A), topk(K, B)))
//	bottomk(K, coalesce(A, B)) = bottomk(K, coalesce(bottomk(K, A), bottomk(K, B)))
func (p PropagateCoalesceOptimizer) pushThroughAggregation(agg *Aggregation, coalesce *Coalesce) (Node, bool) {
	if agg.Expr != coalesce {
		return nil, false
	}

	// Only distributive aggregations (and topk/bottomk) can have coalesce pushed through
	if !isDistributiveAggregation(agg.Op) && agg.Op != parser.TOPK && agg.Op != parser.BOTTOMK {
		return nil, false
	}

	// Check if we've already pushed - if all children are already aggregations,
	// don't push again to avoid infinite loops.
	// For most aggregations, check if children have the same op.
	// For SUM with COUNT children, this is the result of a previous COUNT push-through,
	// so we should not push through again.
	if len(coalesce.Expressions) > 0 {
		allAggregations := true
		firstAggOp := parser.ItemType(0)
		for i, expr := range coalesce.Expressions {
			childAgg, ok := expr.(*Aggregation)
			if !ok {
				allAggregations = false
				break
			}
			if i == 0 {
				firstAggOp = childAgg.Op
			} else if childAgg.Op != firstAggOp {
				allAggregations = false
				break
			}
		}
		if allAggregations {
			// If children are the same aggregation as parent, don't push
			if firstAggOp == agg.Op {
				return nil, false
			}
			// If parent is SUM and children are COUNT, this is the result of
			// a previous COUNT push-through, don't push again
			if agg.Op == parser.SUM && firstAggOp == parser.COUNT {
				return nil, false
			}
		}
	}

	// Create new expressions with local aggregation for each coalesce child
	newExprs := make([]Node, len(coalesce.Expressions))
	for i, expr := range coalesce.Expressions {
		newExprs[i] = &Aggregation{
			Op:       agg.Op,
			Expr:     expr.Clone(),
			Param:    cloneIfNotNil(agg.Param),
			Grouping: agg.Grouping,
			Without:  agg.Without,
		}
	}

	// Wrap with outer aggregation to combine the partial results.
	// For most aggregations, the outer op is the same as the inner op.
	// However, for COUNT, the outer op should be SUM (count(A) + count(B) = total count).
	outerOp := agg.Op
	if agg.Op == parser.COUNT {
		outerOp = parser.SUM
	}

	return &Aggregation{
		Op:       outerOp,
		Expr:     &Coalesce{Expressions: newExprs},
		Param:    cloneIfNotNil(agg.Param),
		Grouping: agg.Grouping,
		Without:  agg.Without,
	}, true
}

// isDistributiveAggregation returns true if the aggregation is distributive,
// meaning it can be computed on shards independently and the results combined.
// For coalesce specifically, since each shard has distinct series, these
// aggregations can be pushed through directly.
func isDistributiveAggregation(op parser.ItemType) bool {
	switch op {
	case parser.SUM, parser.MIN, parser.MAX, parser.COUNT, parser.GROUP:
		return true
	default:
		return false
	}
}

func cloneIfNotNil(n Node) Node {
	if n == nil {
		return nil
	}
	return n.Clone()
}

// isElementWiseFunction returns true if the function operates element-wise
// on its input vector, meaning each output sample only depends on the
// corresponding input sample (not on other samples in the vector).
func isElementWiseFunction(name string) bool {
	switch name {
	// Math functions
	case "abs", "ceil", "floor", "exp", "sqrt", "ln", "log2", "log10",
		"sin", "cos", "tan", "asin", "acos", "atan",
		"sinh", "cosh", "tanh", "asinh", "acosh", "atanh",
		"deg", "rad", "sgn":
		return true
	// Functions with optional scalar argument
	case "round", "clamp", "clamp_min", "clamp_max":
		return true
	// Label manipulation functions - these operate per-series
	case "label_replace", "label_join":
		return true
	// Range vector functions that produce one output per input series
	// These operate on each series independently
	case "rate", "irate", "increase", "delta", "idelta",
		"deriv", "predict_linear",
		"avg_over_time", "min_over_time", "max_over_time",
		"sum_over_time", "count_over_time", "stddev_over_time", "stdvar_over_time",
		"quantile_over_time", "last_over_time", "present_over_time",
		"changes", "resets",
		// Extended range vector functions (from ringbuffer/functions.go)
		"xrate", "xincrease", "xdelta",
		"ts_of_max_over_time", "ts_of_min_over_time", "ts_of_last_over_time",
		"mad_over_time", "first_over_time",
		"double_exponential_smoothing":
		return true
	default:
		return false
	}
}

func shardVectorSelector(vs *VectorSelector, numShards int) *Coalesce {
	shards := make([]Node, numShards)
	for i := 0; i < numShards; i++ {
		shards[i] = &VectorSelector{
			VectorSelector:  vs.VectorSelector,
			Filters:         vs.Filters,
			Projection:      vs.Projection,
			SelectTimestamp: vs.SelectTimestamp,
			Shard:           ShardInfo{Index: i, Total: numShards},
		}
	}
	return &Coalesce{Expressions: shards}
}

func shardMatrixSelector(ms *MatrixSelector, numShards int) *Coalesce {
	shards := make([]Node, numShards)
	for i := 0; i < numShards; i++ {
		shards[i] = &MatrixSelector{
			VectorSelector: &VectorSelector{
				VectorSelector:  ms.VectorSelector.VectorSelector,
				Filters:         ms.VectorSelector.Filters,
				Projection:      ms.VectorSelector.Projection,
				SelectTimestamp: ms.VectorSelector.SelectTimestamp,
				Shard:           ShardInfo{Index: i, Total: numShards},
			},
			Range:          ms.Range,
			OriginalString: ms.OriginalString,
		}
	}
	return &Coalesce{Expressions: shards}
}
