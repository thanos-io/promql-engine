// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"math"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"

	"github.com/thanos-io/promql-engine/api"
	"github.com/thanos-io/promql-engine/query"
)

func TestPassthrough(t *testing.T) {
	expr, err := parser.ParseExpr(`time()`)
	testutil.Ok(t, err)

	t.Run("optimized with one engine, in bounds", func(t *testing.T) {
		engines := []api.RemoteEngine{
			newEngineMock(math.MinInt64, math.MaxInt64, []labels.Labels{labels.FromStrings("region", "east"), labels.FromStrings("region", "south")}),
		}
		optimizers := []Optimizer{PassthroughOptimizer{Endpoints: api.NewStaticEndpoints(engines)}}

		plan := NewFromAST(expr, &query.Options{Start: time.Unix(0, 0), End: time.Unix(0, 0)}, PlanOptions{})
		optimizedPlan, _ := plan.Optimize(optimizers)

		testutil.Equals(t, "remote(time())", renderExprTree(optimizedPlan.Root()))
	})

	t.Run("not optimized with two engines", func(t *testing.T) {
		engines := []api.RemoteEngine{
			newEngineMock(math.MinInt64, math.MinInt64, []labels.Labels{labels.FromStrings("region", "east"), labels.FromStrings("region", "south")}),
			newEngineMock(math.MinInt64, math.MinInt64, []labels.Labels{labels.FromStrings("region", "west")}),
		}
		optimizers := []Optimizer{PassthroughOptimizer{Endpoints: api.NewStaticEndpoints(engines)}}

		plan := NewFromAST(expr, &query.Options{Start: time.Unix(0, 0), End: time.Unix(0, 0)}, PlanOptions{})
		optimizedPlan, _ := plan.Optimize(optimizers)

		testutil.Equals(t, "time()", renderExprTree(optimizedPlan.Root()))
	})

	t.Run("not optimized with one out of bound engine", func(t *testing.T) {
		engines := []api.RemoteEngine{
			newEngineMock(math.MinInt64, math.MinInt64, []labels.Labels{labels.FromStrings("region", "east"), labels.FromStrings("region", "south")}),
		}
		optimizers := []Optimizer{PassthroughOptimizer{Endpoints: api.NewStaticEndpoints(engines)}}

		plan := NewFromAST(expr, &query.Options{Start: time.Unix(0, 0), End: time.Unix(0, 0)}, PlanOptions{})
		optimizedPlan, _ := plan.Optimize(optimizers)

		testutil.Equals(t, "time()", renderExprTree(optimizedPlan.Root()))
	})

	t.Run("optimized with matching labels", func(t *testing.T) {
		selectorExpr, err := parser.ParseExpr(`{region="east"}`)
		testutil.Ok(t, err)

		engines := []api.RemoteEngine{
			newEngineMock(math.MinInt64, math.MaxInt64, []labels.Labels{labels.FromStrings("region", "east"), labels.FromStrings("region", "south")}),
			newEngineMock(math.MinInt64, math.MaxInt64, []labels.Labels{labels.FromStrings("region", "west")}),
		}
		optimizers := []Optimizer{PassthroughOptimizer{Endpoints: api.NewStaticEndpoints(engines)}}

		plan := NewFromAST(selectorExpr, &query.Options{Start: time.Unix(0, 0), End: time.Unix(0, 0)}, PlanOptions{})
		optimizedPlan, _ := plan.Optimize(optimizers)

		testutil.Equals(t, `remote({region="east"})`, renderExprTree(optimizedPlan.Root()))
	})

	t.Run("not optimized due to multiple engines", func(t *testing.T) {
		selectorExpr, err := parser.ParseExpr(`{region=~"east|west"}`)
		testutil.Ok(t, err)

		engines := []api.RemoteEngine{
			newEngineMock(math.MinInt64, math.MaxInt64, []labels.Labels{labels.FromStrings("region", "east"), labels.FromStrings("region", "south")}),
			newEngineMock(math.MinInt64, math.MaxInt64, []labels.Labels{labels.FromStrings("region", "west")}),
		}
		optimizers := []Optimizer{PassthroughOptimizer{Endpoints: api.NewStaticEndpoints(engines)}}

		plan := NewFromAST(selectorExpr, &query.Options{Start: time.Unix(0, 0), End: time.Unix(0, 0)}, PlanOptions{})
		optimizedPlan, _ := plan.Optimize(optimizers)

		testutil.Equals(t, `{region=~"east|west"}`, renderExprTree(optimizedPlan.Root()))
	})

	t.Run("optimized with matching labels on matrix selector", func(t *testing.T) {
		selectorExpr, err := parser.ParseExpr(`{region="east"}[5m]`)
		testutil.Ok(t, err)

		engines := []api.RemoteEngine{
			newEngineMock(math.MinInt64, math.MaxInt64, []labels.Labels{labels.FromStrings("region", "east"), labels.FromStrings("region", "south")}),
			newEngineMock(math.MinInt64, math.MaxInt64, []labels.Labels{labels.FromStrings("region", "west")}),
		}
		optimizers := []Optimizer{PassthroughOptimizer{Endpoints: api.NewStaticEndpoints(engines)}}

		plan := NewFromAST(selectorExpr, &query.Options{Start: time.Unix(0, 0), End: time.Unix(0, 0)}, PlanOptions{})
		optimizedPlan, _ := plan.Optimize(optimizers)

		testutil.Equals(t, `remote({region="east"}[5m])`, renderExprTree(optimizedPlan.Root()))
	})

	t.Run("not optimized with matching labels but not matching time", func(t *testing.T) {
		selectorExpr, err := parser.ParseExpr(`{region="east"}`)
		testutil.Ok(t, err)

		engines := []api.RemoteEngine{
			newEngineMock(math.MinInt64, math.MinInt64, []labels.Labels{labels.FromStrings("region", "east"), labels.FromStrings("region", "south")}),
			newEngineMock(math.MinInt64, math.MaxInt64, []labels.Labels{labels.FromStrings("region", "west")}),
		}
		optimizers := []Optimizer{PassthroughOptimizer{Endpoints: api.NewStaticEndpoints(engines)}}

		plan := NewFromAST(selectorExpr, &query.Options{Start: time.Unix(0, 0), End: time.Unix(0, 0)}, PlanOptions{})
		optimizedPlan, _ := plan.Optimize(optimizers)

		testutil.Equals(t, `{region="east"}`, renderExprTree(optimizedPlan.Root()))
	})

}
