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

// TestProjectionAfterMergeSelects runs the projection optimizer in the order
// Thanos uses it: appended AFTER logicalplan.DefaultOptimizers, so
// MergeSelectsOptimizer rewrites selectors before projections are pushed.
func TestProjectionAfterMergeSelects(t *testing.T) {
	t.Parallel()

	load := `load 30s
		http_requests_total{pod="nginx-1", job="app", env="prod", instance="1"} 1+1x40
		http_requests_total{pod="nginx-2", job="app", env="dev", instance="2"} 2+2x40
		http_requests_total{pod="nginx-3", job="api", env="prod", instance="3"} 3+3x40
		http_requests_total{pod="nginx-4", job="api", env="dev", instance="4"} 4+4x40`

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	engineOpts := promql.EngineOpts{Timeout: time.Minute, MaxSamples: 1e10}
	normalEngine := engine.New(engine.Opts{
		EngineOpts:        engineOpts,
		LogicalOptimizers: logicalplan.DefaultOptimizers,
	})
	// Thanos's order: projections last, after MergeSelectsOptimizer.
	projectionEngine := engine.New(engine.Opts{
		EngineOpts:        engineOpts,
		LogicalOptimizers: append(append([]logicalplan.Optimizer{}, logicalplan.DefaultOptimizers...), logicalplan.ProjectionOptimizer{}),
	})
	projectionStorage := &projectionQueryable{Queryable: storage, useAtHash: true}

	ctx := context.Background()
	queryTime := time.Unix(600, 0)
	queries := []string{
		// The subset selector is rewritten to the superset plus a job="app"
		// filter; the projection then trims job away.
		`sum by (env) (http_requests_total{job="app"}) + on (env) sum by (env) (http_requests_total)`,
		`sum by (env) (http_requests_total{job="app"})`,
		`sum by (env) (http_requests_total{job="app"}) / scalar(sum(http_requests_total))`,
		`sum by (env, pod) (http_requests_total{job="app"}) + on (env) group_left () sum by (env) (http_requests_total)`,
		`sum by (env) (rate(http_requests_total{job="app"}[2m])) + on (env) sum by (env) (rate(http_requests_total[2m]))`,
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
			testutil.Ok(t, projectionResult.Err, "query: %s", query)

			if diff := cmp.Diff(normalResult, projectionResult, comparer); diff != "" {
				t.Errorf("Results differ for query %s:\nnormal:    %v\nprojected: %v", query, normalResult.Value, projectionResult.Value)
			}
		})
	}
}
