// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/logicalplan"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/promqltest"
	"github.com/prometheus/prometheus/util/stats"
)

func TestLimitk(t *testing.T) {
	load := fmt.Sprintf(`load 30s
			http_requests_total{pod="nginx-1", route="/"} 5+3x10
			http_requests_total{pod="nginx-2", route="/"} 4+4x10
		  http_requests_total{pod="nginx-1", route="/health"} 3+5x10
			http_requests_total{pod="nginx-1", route="/test"} 5+2x10
			http_requests_total{pod="nginx-2", route="/health"} 1+4x10
	`)

	loadNativeHistograms := fmt.Sprintf(`load 2m
			http_request_duration_seconds{pod="nginx-1"} {{schema:0 count:3 sum:2 buckets:[1 2]}}+{{schema:1 count:5 buckets:[1 2 3]}}x20 
			http_request_duration_seconds{pod="nginx-2"} {{schema:-2 count:2 sum:1 buckets:[2]}}+{{schema:-1 count:9 buckets:[2 3 4]}}x20`)

	_ = fmt.Sprintln(loadNativeHistograms)

	opts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
		EnablePerStepStats:   true,
	}
	qOpts := promql.NewPrometheusQueryOpts(true, 0)

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	queryTime := time.Unix(int64(120), 0)
	newEngine := engine.New(engine.Opts{
		EngineOpts:        opts,
		LogicalOptimizers: logicalplan.AllOptimizers,
		EnableAnalysis:    true,
	})
	oldEngine := promql.NewEngine(opts)

	query := "limitk(2, http_request_duration_seconds)"

	q1, err := newEngine.NewInstantQuery(context.Background(), storage, qOpts, query, queryTime)

	testutil.Ok(t, err)
	newResult := q1.Exec(context.Background())
	newStats := q1.Stats()
	stats.NewQueryStats(newStats)

	fmt.Println(newResult)
	fmt.Println("Samples: ", newStats.Samples)

	fmt.Println("---")

	q2, err := oldEngine.NewInstantQuery(context.Background(), storage, qOpts, query, queryTime)
	testutil.Ok(t, err)

	oldResult := q2.Exec(context.Background())
	oldStats := q2.Stats()
	stats.NewQueryStats(oldStats)

	fmt.Println(oldResult)
	fmt.Println("Samples: ", oldStats.Samples)

	cases := make([]*testCase, 1)

	cases[0] = &testCase{
		query:           query,
		newRes:          newResult,
		newStats:        newStats,
		oldRes:          oldResult,
		oldStats:        oldStats,
		loads:           []string{load},
		start:           queryTime,
		end:             queryTime,
		validateSamples: true,
	}

	validateTestCases(t, cases[:1])
}
