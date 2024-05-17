// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"github.com/prometheus/prometheus/util/annotations"

	"github.com/thanos-io/promql-engine/query"
)

type DetectHistogramStatsOptimizer struct{}

func (d DetectHistogramStatsOptimizer) Optimize(plan Node, _ *query.Options) (Node, annotations.Annotations) {
	return d.optimize(plan, false)
}

func (d DetectHistogramStatsOptimizer) optimize(plan Node, decodeStats bool) (Node, annotations.Annotations) {
	var stop bool
	Traverse(&plan, func(node *Node) {
		if stop {
			return
		}
		switch n := (*node).(type) {
		case *VectorSelector:
			n.DecodeNativeHistogramStats = decodeStats
		case *FunctionCall:
			if n.Func.Name == "histogram_count" || n.Func.Name == "histogram_sum" {
				n.Args[0], _ = d.optimize(n.Args[0], true)
				stop = true
				return
			}
		}
	})
	return plan, nil
}
