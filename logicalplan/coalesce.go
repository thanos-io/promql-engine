// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/promql/parser/posrange"
	"github.com/prometheus/prometheus/util/annotations"

	"github.com/thanos-io/promql-engine/query"
)

type Coalesce struct {
	// We assume to always have at least one expression
	Exprs []parser.Expr
}

func (c Coalesce) String() string {
	return c.Exprs[0].String()
}

func (c Coalesce) Pretty(level int) string { return c.String() }

func (c Coalesce) PositionRange() posrange.PositionRange { return c.Exprs[0].PositionRange() }

func (c Coalesce) Type() parser.ValueType { return c.Exprs[0].Type() }

func (c Coalesce) PromQLExpr() {}

type CoalesceOptimizer struct{}

func (c CoalesceOptimizer) Optimize(expr parser.Expr, opts *query.Options) (parser.Expr, annotations.Annotations) {
	numShards := opts.NumShards()

	TraverseBottomUp(nil, &expr, func(parent, e *parser.Expr) bool {
		switch t := (*e).(type) {
		case *VectorSelector:
			if parent != nil {
				// we coalesce matrix selectors in a different branch
				if _, ok := (*parent).(MatrixSelector); ok {
					return false
				}
			}
			exprs := make([]parser.Expr, numShards)
			for i := 0; i < numShards; i++ {
				exprs[i] = &VectorSelector{
					VectorSelector:  t.VectorSelector,
					Filters:         t.Filters,
					BatchSize:       t.BatchSize,
					SelectTimestamp: t.SelectTimestamp,
					Shard:           i,
					NumShards:       numShards,
				}
			}
			*e = Coalesce{Exprs: exprs}
			return true
		case *MatrixSelector:
			// handled in *parser.Call branch
			return false
		case *parser.Call:
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

			exprs := make([]parser.Expr, numShards)
			for i := 0; i < numShards; i++ {
				aux := &MatrixSelector{
					MatrixSelector: ms.MatrixSelector,
					OriginalString: ms.OriginalString,
					Shard:          i,
					NumShards:      numShards,
				}
				f := &parser.Call{
					Func:     t.Func,
					Args:     t.Args,
					PosRange: t.PosRange,
				}
				f.Args[marg] = aux

				exprs[i] = f
			}
			*e = Coalesce{Exprs: exprs}
			return true
		default:
			return true
		}
	})
	return expr, nil
}
