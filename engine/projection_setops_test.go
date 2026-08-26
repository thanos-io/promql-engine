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

func TestProjectionSetOperations(t *testing.T) {
	t.Parallel()

	// errors_total has an env (staging) that http_requests_total lacks, so
	// `or on(env)` emits right-hand series into the output.
	load := `load 30s
		http_requests_total{pod="nginx-1", job="app", env="prod", instance="1"} 1+1x40
		http_requests_total{pod="nginx-2", job="app", env="dev", instance="2"} 2+2x40
		errors_total{pod="nginx-3", job="api", env="prod", instance="3"} 3+3x40
		errors_total{pod="nginx-4", job="api", env="staging", instance="4"} 4+4x40`

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	engineOpts := promql.EngineOpts{Timeout: time.Minute, MaxSamples: 1e10}
	normalEngine := engine.New(engine.Opts{
		EngineOpts:        engineOpts,
		LogicalOptimizers: logicalplan.AllOptimizers,
	})
	projectionEngine := engine.New(engine.Opts{
		EngineOpts: engineOpts,
		LogicalOptimizers: []logicalplan.Optimizer{
			logicalplan.SortMatchers{},
			logicalplan.ProjectionOptimizer{},
			logicalplan.DetectHistogramStatsOptimizer{},
			logicalplan.MergeSelectsOptimizer{},
		},
	})
	projectionStorage := &projectionQueryable{Queryable: storage, useAtHash: true}

	ctx := context.Background()
	queryTime := time.Unix(600, 0)
	queries := []string{
		`http_requests_total or on (env) errors_total`,
		`sum by (pod) (http_requests_total or on (env) errors_total)`,
		`sum by (pod) (http_requests_total or ignoring (instance, job, pod) errors_total)`,
		`sum by (pod) (http_requests_total and on (env) errors_total)`,
		`sum by (pod) (http_requests_total unless on (env) errors_total)`,
		`count by (job) (http_requests_total or on (env) errors_total)`,
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
