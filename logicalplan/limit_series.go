// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"strconv"

	"github.com/thanos-io/promql-engine/query"

	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/util/annotations"
)

type LimitSeriesOptmizer struct{}

func (l LimitSeriesOptmizer) Optimize(plan Node, qOpts *query.Options) (Node, annotations.Annotations) {
	TraverseBottomUp(nil, &plan, func(parent, node *Node) (stop bool) {
		switch e := (*node).(type) {
		case *Binary:
			// we are traversing bottom-up so if we get Binary Expr, series to fetch can't be limited, stop the traversal
			return true
		case *Aggregation:
			switch e.Op {
			case parser.LIMITK:
				if len(e.Grouping) != 0 {
					return true
				}
				e.Limit, _ = strconv.Atoi(e.Param.String())
			default:
				// if limitk has any other aggregation at its downstream operator, then limit can't be imposed, else if limitk is itself downstream, then limiting is possible
				return true
			}
		}
		return false
	})
	return plan, nil
}
