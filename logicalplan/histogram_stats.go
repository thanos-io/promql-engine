// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"github.com/prometheus/prometheus/util/annotations"

	"github.com/thanos-io/promql-engine/query"
)

type DetectHistogramStatsOptimizer struct{}

func (d DetectHistogramStatsOptimizer) Optimize(plan Node, _ *query.Options) (Node, annotations.Annotations) {
	return d.optimize(plan)
}

func (d DetectHistogramStatsOptimizer) optimize(plan Node) (Node, annotations.Annotations) {
	var (
		stop        bool
		decodeStats bool
	)
	Traverse(&plan, func(node *Node) {
		if stop {
			return
		}
		switch n := (*node).(type) {
		case *VectorSelector:
			n.DecodeNativeHistogramStats = decodeStats
		case *Binary:
			n.LHS, _ = d.optimize(n.LHS)
			n.RHS, _ = d.optimize(n.RHS)
			stop = true
			return
		case *FunctionCall:
			if n.Func.Name == "histogram_count" || n.Func.Name == "histogram_sum" {
				decodeStats = true
				return
			}
			if n.Func.Name == "histogram_quantile" || n.Func.Name == "histogram_fraction" {
				decodeStats = false
				return
			}
		}
	})
	return plan, nil
}
