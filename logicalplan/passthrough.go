// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"github.com/thanos-io/promql-engine/api"
	"github.com/thanos-io/promql-engine/parser"
	"github.com/thanos-io/promql-engine/query"
)

// PassthroughOptimizer optimizes queries which can be simply passed
// through to a RemoteEngine.
type PassthroughOptimizer struct {
	Endpoints api.RemoteEndpoints
}

func (m PassthroughOptimizer) Optimize(plan parser.Expr, opts *query.Options) parser.Expr {
	engines := m.Endpoints.Engines()
	if len(engines) == 1 {
		return RemoteExecution{
			Engine:          engines[0],
			Query:           plan.String(),
			QueryRangeStart: opts.Start,
		}
	}

	return plan
}
