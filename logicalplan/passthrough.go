// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"github.com/prometheus/prometheus/model/labels"

	"github.com/thanos-io/promql-engine/api"
	"github.com/thanos-io/promql-engine/parser"
	"github.com/thanos-io/promql-engine/query"
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

func (m PassthroughOptimizer) Optimize(plan parser.Expr, opts *query.Options) parser.Expr {
	engines := m.Endpoints.Engines()
	if len(engines) == 1 {
		if !matchingEngineTime(engines[0], opts) {
			return plan
		}
		return RemoteExecution{
			Engine:          engines[0],
			Query:           plan.String(),
			QueryRangeStart: opts.Start,
		}
	}

	if len(engines) == 0 {
		return plan
	}

	// Check matchers of each selector. If all of them match only one engine
	// then pass the query to it.
	var matchingEngine int = -1
	var matchingFound bool
	TraverseBottomUp(nil, &plan, func(parent, current *parser.Expr) (stop bool) {
		if vs, ok := (*current).(*parser.VectorSelector); ok {
			for i, e := range engines {
				if !labelSetsMatch(vs.LabelMatchers, e.LabelSets()...) {
					continue
				}

				if matchingEngine == -1 {
					matchingEngine = i
					matchingFound = true
				} else if matchingEngine != i {
					matchingFound = false
					return true
				}
			}
		}
		return false
	})

	if matchingFound && matchingEngineTime(engines[matchingEngine], opts) {
		return RemoteExecution{
			Engine:          engines[matchingEngine],
			Query:           plan.String(),
			QueryRangeStart: opts.Start,
		}
	}

	return plan
}
