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
	var limitNode *Aggregation
	canLimit := true
	TraverseBottomUp(nil, &plan, func(parent, node *Node) (stop bool) {
		switch e := (*node).(type) {
		case *Binary:
			canLimit = false
			// we are traversing bottom-up so if we get Binary Expr, series to fetch can't be limited, stop the traversal
			return true
		case *Aggregation:
			switch e.Op {
			case parser.LIMITK:
				if len(e.Grouping) != 0 {
					canLimit = false
					return true
				}
				limitNode = e
			default:
				// if limitk has any other aggregation at its downstream operator, then limit can't be imposed, else if limitk is itself downstream, then limiting is possible
				canLimit = false
				return true
			}
		}
		return false
	})
	if canLimit && limitNode != nil {
		limitNode.Limit, _ = strconv.Atoi(limitNode.Param.String())
	}
	return plan, nil
}
