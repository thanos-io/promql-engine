// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"math"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"

	"github.com/thanos-io/promql-engine/api"
	"github.com/thanos-io/promql-engine/parser"
)

func TestPassthrough(t *testing.T) {
	expr, err := parser.ParseExpr(`time()`)
	testutil.Ok(t, err)

	t.Run("optimized with one engine", func(t *testing.T) {
		engines := []api.RemoteEngine{
			newEngineMock(math.MinInt64, math.MinInt64, []labels.Labels{labels.FromStrings("region", "east"), labels.FromStrings("region", "south")}),
		}
		optimizers := []Optimizer{PassthroughOptimizer{Endpoints: api.NewStaticEndpoints(engines)}}

		plan := New(expr, &Opts{Start: time.Unix(0, 0), End: time.Unix(0, 0)})
		optimizedPlan := plan.Optimize(optimizers)

		testutil.Equals(t, "remote(time())", optimizedPlan.Expr().String())
	})

	t.Run("not optimized with one engine", func(t *testing.T) {
		engines := []api.RemoteEngine{
			newEngineMock(math.MinInt64, math.MinInt64, []labels.Labels{labels.FromStrings("region", "east"), labels.FromStrings("region", "south")}),
			newEngineMock(math.MinInt64, math.MinInt64, []labels.Labels{labels.FromStrings("region", "west")}),
		}
		optimizers := []Optimizer{PassthroughOptimizer{Endpoints: api.NewStaticEndpoints(engines)}}

		plan := New(expr, &Opts{Start: time.Unix(0, 0), End: time.Unix(0, 0)})
		optimizedPlan := plan.Optimize(optimizers)

		testutil.Equals(t, "time()", optimizedPlan.Expr().String())
	})

}
