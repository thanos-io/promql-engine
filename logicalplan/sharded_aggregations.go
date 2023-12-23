// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"fmt"
	"strings"

	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/promql/parser/posrange"
	"github.com/prometheus/prometheus/util/annotations"

	"github.com/thanos-io/promql-engine/query"
)

type Coalesce struct {
	Shards []parser.Expr
}

func (r Coalesce) String() string {
	parts := make([]string, len(r.Shards))
	for i, r := range r.Shards {
		parts[i] = r.String()
	}
	return fmt.Sprintf("coalesce(%s)", strings.Join(parts, ", "))
}

func (r Coalesce) Pretty(level int) string { return r.String() }

func (r Coalesce) PositionRange() posrange.PositionRange { return posrange.PositionRange{} }

func (r Coalesce) Type() parser.ValueType { return r.Shards[0].Type() }

func (r Coalesce) PromQLExpr() {}

type ShardedAggregations struct{ Shards int }

func (m ShardedAggregations) Optimize(plan parser.Expr, _ *query.Options) (parser.Expr, annotations.Annotations) {
	TraverseBottomUp(nil, &plan, func(parent, current *parser.Expr) (stop bool) {
		if parent == nil {
			return false
		}
		aggr, ok := (*parent).(*parser.AggregateExpr)
		if !ok {
			return false
		}
		// TODO: only care about sum now
		if aggr.Op == parser.COUNT {
			return false
		}
		vs, ok := (*current).(*parser.VectorSelector)
		if !ok {
			return false
		}

		coalesce := Coalesce{make([]parser.Expr, m.Shards)}
		for i := range coalesce.Shards {
			coalesce.Shards[i] = &parser.AggregateExpr{
				Op:       aggr.Op,
				Expr:     vectorSelectorForShard(vs, i, m.Shards),
				Param:    aggr.Param,
				Grouping: aggr.Grouping,
				Without:  aggr.Without,
				PosRange: aggr.PosRange,
			}
		}

		*parent = &parser.AggregateExpr{
			Op:       aggr.Op,
			Expr:     coalesce,
			Param:    aggr.Param,
			Grouping: aggr.Grouping,
			Without:  aggr.Without,
			PosRange: aggr.PosRange,
		}
		return true
	})
	return plan, nil
}

func vectorSelectorForShard(expr *parser.VectorSelector, n, shards int) parser.Expr {
	return &VectorSelector{
		VectorSelector: expr,
		N:              n,
		Shards:         shards,
	}
}
