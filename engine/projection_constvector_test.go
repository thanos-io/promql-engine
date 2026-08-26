// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/logicalplan"

	"github.com/efficientgo/core/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/promqltest"
)

// TestProjectionConstantVectorMatching covers binary expressions with default
// (all-labels) vector matching where one side is a constant-but-vector
// expression such as vector() or day_of_week(): trimming the other side
// changes which label sets match.
func TestProjectionConstantVectorMatching(t *testing.T) {
	t.Parallel()

	load := `load 30s
		http_requests_total{pod="nginx-1", job="app", env="prod"} 1+1x40
		http_requests_total{pod="nginx-2", job="app", env="dev"} 2+2x40`

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	engineOpts := promql.EngineOpts{Timeout: time.Minute, MaxSamples: 1e10, EnableAtModifier: true}
	normalEngine := engine.New(engine.Opts{
		EngineOpts:        engineOpts,
		LogicalOptimizers: logicalplan.DefaultOptimizers,
	})
	projectionEngine := engine.New(engine.Opts{
		EngineOpts:        engineOpts,
		LogicalOptimizers: append(append([]logicalplan.Optimizer{}, logicalplan.DefaultOptimizers...), logicalplan.ProjectionOptimizer{}),
	})
	projectionStorage := &projectionQueryable{Queryable: storage, useAtHash: true}

	ctx := context.Background()
	queryTime := time.Unix(600, 0)
	queries := []string{
		// One matching series: projection include{} trims it to {} which then
		// matches vector()'s empty label set.
		`sum(vector(0) > bool http_requests_total{pod="nginx-1"})`,
		`sum(http_requests_total{pod="nginx-1"} + vector(1))`,
		`count(vector(1) == bool http_requests_total{pod="nginx-1"})`,
		// Several matching series: trimming makes their signatures collide.
		`sum by (job) (vector(0) > bool http_requests_total)`,
		`sum(day_of_week() * http_requests_total{pod="nginx-1"})`,
	}
	for _, query := range queries {
		t.Run(query, func(t *testing.T) {
			normalQuery, err := normalEngine.NewInstantQuery(ctx, storage, &engine.QueryOpts{}, query, queryTime)
			testutil.Ok(t, err)
			defer normalQuery.Close()
			normalResult := normalQuery.Exec(ctx)
			testutil.Ok(t, normalResult.Err, "query: %s", query)

			projectionQuery, err := projectionEngine.MakeInstantQuery(ctx, projectionStorage, &engine.QueryOpts{}, query, queryTime)
			testutil.Ok(t, err)
			defer projectionQuery.Close()
			projectionResult := projectionQuery.Exec(ctx)
			if projectionResult.Err != nil {
				t.Fatalf("projected query errored where normal returned %v: %v", normalResult.Value, projectionResult.Err)
			}

			if diff := cmp.Diff(normalResult, projectionResult, comparer); diff != "" {
				t.Errorf("Results differ for %s:\nnormal:    %v\nprojected: %v", query, normalResult.Value, projectionResult.Value)
			}
		})
	}
}
