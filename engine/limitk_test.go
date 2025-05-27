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
		http_requests_total{pod="nginx-1", series="1"} 1+1.1x50
		http_requests_total{pod="nginx-2", series="1"} 2+2.3x50
		http_requests_total{pod="nginx-4", series="2"} 5+2.4x50
		http_requests_total{pod="nginx-5", series="2"} 8.4+2.3x50
	  http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50
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
		EnablePerStepStats:   false,
	}
	//qOpts := promql.NewPrometheusQueryOpts(false, 5*time.Minute)

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	queryTime := time.Unix(int64(0), 0)
	/* queryEnd := time.Unix(int64(1700), 0) */
	newEngine := engine.New(engine.Opts{
		EngineOpts:        opts,
		LogicalOptimizers: logicalplan.AllOptimizers,
		EnableAnalysis:    true,
	})
	oldEngine := promql.NewEngine(opts)

	query := "limitk(2, http_requests_total) or http_requests_total"

	q1, err := newEngine.NewInstantQuery(context.Background(), storage, nil, query, queryTime)

	testutil.Ok(t, err)
	newResult := q1.Exec(context.Background())
	newStats := q1.Stats()
	stats.NewQueryStats(newStats)

	fmt.Println(newResult)
	fmt.Println("Samples: ", newStats.Samples)

	fmt.Println("---")

	q2, err := oldEngine.NewInstantQuery(context.Background(), storage, nil, query, queryTime)
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
