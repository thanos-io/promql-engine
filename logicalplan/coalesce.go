// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"github.com/prometheus/prometheus/util/annotations"

	"github.com/thanos-io/promql-engine/query"
)

type CoalesceOptimizer struct{}

func (c CoalesceOptimizer) Optimize(expr Node, opts *query.Options) (Node, annotations.Annotations) {
	numShards := opts.NumShards()

	TraverseBottomUp(nil, &expr, func(parent, e *Node) bool {
		switch t := (*e).(type) {
		case *VectorSelector:
			if parent != nil {
				// we coalesce matrix selectors in a different branch
				if _, ok := (*parent).(*MatrixSelector); ok {
					return false
				}
			}
			exprs := make([]Node, numShards)
			for i := range numShards {
				vs := t.Clone().(*VectorSelector)
				vs.Shard = i
				vs.NumShards = numShards
				exprs[i] = vs
			}
			*e = &Coalesce{Exprs: exprs}
		case *MatrixSelector:
			// handled in *parser.Call branch
			return false
		case *FunctionCall:
			// non-recursively handled in execution.go
			if t.Func.Name == "absent_over_time" {
				return true
			}
			var (
				ms   *MatrixSelector
				marg int
			)
			for i := range t.Args {
				if arg, ok := t.Args[i].(*MatrixSelector); ok {
					ms = arg
					marg = i
				}
			}
			if ms == nil {
				return false
			}
			exprs := make([]Node, numShards)
			for i := range numShards {
				aux := ms.Clone().(*MatrixSelector)
				aux.VectorSelector.Shard = i
				aux.VectorSelector.NumShards = numShards
				f := t.Clone().(*FunctionCall)
				f.Args[marg] = aux
				exprs[i] = f
			}
			*e = &Coalesce{Exprs: exprs}
		}
		return true
	})
	return expr, nil
}
