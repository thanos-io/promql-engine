// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"github.com/thanos-io/promql-engine/api"
	"github.com/thanos-io/promql-engine/query"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/util/annotations"
)

// PassthroughOptimizer optimizes queries which can be simply passed
// through to a RemoteEngine.
type PassthroughOptimizer struct {
	Endpoints api.RemoteEndpoints
}

// labelSetsMatch returns false if all label-set do not match the matchers (aka: OR is between all label-sets).
func labelSetsMatch(matchers []*labels.Matcher, lset ...labels.Labels) bool {
	if len(lset) == 0 {
		return true
	}

	for _, ls := range lset {
		notMatched := false
		for _, m := range matchers {
			if lv := ls.Get(m.Name); ls.Has(m.Name) && !m.Matches(lv) {
				notMatched = true
				break
			}
		}
		if !notMatched {
			return true
		}
	}
	return false
}

func matchingEngineTime(e api.RemoteEngine, opts *query.Options) bool {
	return !(opts.Start.UnixMilli() > e.MaxT() || opts.End.UnixMilli() < e.MinT())
}

func (m PassthroughOptimizer) Optimize(plan Node, opts *query.Options) (Node, annotations.Annotations) {
	engines := m.Endpoints.Engines(MinMaxTime(plan, opts))
	if len(engines) == 0 {
		return plan, nil
	}

	var (
		hasSelector         bool
		matchingEngines     int
		firstMatchingEngine api.RemoteEngine
	)
	TraverseBottomUp(nil, &plan, func(parent, current *Node) (stop bool) {
		if vs, ok := (*current).(*VectorSelector); ok {
			hasSelector = true

			for _, e := range engines {
				if !labelSetsMatch(vs.LabelMatchers, e.LabelSets()...) {
					continue
				}

				matchingEngines++
				if matchingEngines > 1 {
					return true
				}

				firstMatchingEngine = e
			}
		}
		return false
	})

	// Fallback to all engines.
	if !hasSelector && matchingEngines == 0 {
		matchingEngines = len(engines)
		firstMatchingEngine = engines[0]
	}

	if matchingEngines == 1 && matchingEngineTime(firstMatchingEngine, opts) {
		return RemoteExecution{
			Engine:          firstMatchingEngine,
			Query:           plan.Clone(),
			QueryRangeStart: opts.Start,
			QueryRangeEnd:   opts.End,
		}, nil
	}

	return plan, nil
}
