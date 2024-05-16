// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"github.com/prometheus/prometheus/util/annotations"

	"github.com/thanos-io/promql-engine/query"
)

type DetectHistogramStatsOptimizer struct{}

func (d DetectHistogramStatsOptimizer) Optimize(plan Node, _ *query.Options) (Node, annotations.Annotations) {
	var decodeStats bool
	Traverse(&plan, func(node *Node) {
		switch n := (*node).(type) {
		case *VectorSelector:
			n.DecodeNativeHistogramStats = decodeStats
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
