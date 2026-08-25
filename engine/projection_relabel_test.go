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

// TestProjectionRelabelWithAtHash covers label_replace on top of a projected
// selector: the projection trims distinct series down to identical label sets,
// so the duplicate labelset check must keep using the original series identity
// across the relabel operator.
func TestProjectionRelabelWithAtHash(t *testing.T) {
	t.Parallel()

	// Two series share the projected label (pod) and differ only in a label the
	// projection drops (instance), mirroring several instances per host.
	load := `load 30s
		http_requests_total{pod="nginx-1", job="app", instance="1"} 1+1x40
		http_requests_total{pod="nginx-1", job="app", instance="2"} 2+2x40
		http_requests_total{pod="nginx-2", job="app", instance="3"} 3+3x40`

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	engineOpts := promql.EngineOpts{
		Timeout:    1 * time.Minute,
		MaxSamples: 1e10,
	}
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
		`sum by (datacenter) (label_replace(http_requests_total, "datacenter", "$1", "pod", "^(nginx).*"))`,
		`sum by (datacenter) (label_replace(label_replace(http_requests_total, "datacenter", "$1", "pod", "^(nginx).*"), "datacenter", "$1", "pod", "^(nginx).*"))`,
		`sum by (datacenter) (label_join(http_requests_total, "datacenter", "-", "pod"))`,
		// Other operators whose output series map one-to-one onto the
		// projected storage series.
		`sum by (pod) (timestamp(http_requests_total))`,
		`sum by (pod) (timestamp(abs(http_requests_total)))`,
		`sum by (pod) (timestamp(http_requests_total @ 600))`,
		`sum by (pod) (day_of_week(http_requests_total @ 600))`,
		`sum by (pod) (max_over_time(http_requests_total[2m:30s]))`,
		`sum by (pod) (abs(http_requests_total @ 600))`,
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
				t.Errorf("Results differ for query %s: %s", query, diff)
			}
		})
	}
}
