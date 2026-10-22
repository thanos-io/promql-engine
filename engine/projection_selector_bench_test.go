// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/logicalplan"
	"github.com/thanos-io/promql-engine/query"
	prometheusStorage "github.com/thanos-io/promql-engine/storage/prometheus"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/promqltest"
	"github.com/prometheus/prometheus/storage"
)

// BenchmarkProjectionSelector measures the memory savings from projection pushdown
// at the selector level. It uses a high-cardinality series set and compares
// memory allocations with and without projection.
func BenchmarkProjectionSelector(b *testing.B) {
	// Generate high-cardinality data: 1000 series with many labels
	load := synthesizeProjectionBenchData(1000, 8, 240)
	testStorage := promqltest.LoadedStorage(b, load)
	defer testStorage.Close()

	start := time.Unix(0, 0)
	end := start.Add(2 * time.Hour)
	step := 30 * time.Second

	engineOpts := promql.EngineOpts{
		Timeout:    1 * time.Hour,
		MaxSamples: 1e10,
	}

	queries := []struct {
		name  string
		query string
	}{
		{name: "sum", query: `sum(bench_metric)`},
		{name: "sum_by_job", query: `sum by (job) (bench_metric)`},
		{name: "sum_rate", query: `sum(rate(bench_metric[5m]))`},
		{name: "sum_by_job_rate", query: `sum by (job) (rate(bench_metric[5m]))`},
		{name: "count", query: `count(bench_metric)`},
		{name: "avg_by_region", query: `avg by (region) (bench_metric)`},
	}

	for _, tc := range queries {
		b.Run(tc.name+"/no_projection", func(b *testing.B) {
			ng := engine.New(engine.Opts{
				EngineOpts:        engineOpts,
				LogicalOptimizers: logicalplan.AllOptimizers,
			})
			ctx := context.Background()
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				qry, err := ng.NewRangeQuery(ctx, testStorage, nil, tc.query, start, end, step)
				testutil.Ok(b, err)
				result := qry.Exec(ctx)
				testutil.Ok(b, result.Err)
				qry.Close()
			}
		})

		b.Run(tc.name+"/with_projection", func(b *testing.B) {
			ng := newProjectionSelectorBenchEngine(b, testStorage, start, end, engineOpts)
			ctx := context.Background()
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				qry, err := ng.MakeRangeQuery(ctx, testStorage, &engine.QueryOpts{}, tc.query, start, end, step)
				testutil.Ok(b, err)
				result := qry.Exec(ctx)
				testutil.Ok(b, result.Err)
				qry.Close()
			}
		})
	}
}

// synthesizeProjectionBenchData generates test data with many labels per series.
// Each series has numLabels labels to make label materialization expensive.
func synthesizeProjectionBenchData(numSeries, numLabels, numSteps int) string {
	var sb strings.Builder
	sb.WriteString("load 30s\n")

	regions := []string{"us-west-2", "us-east-1", "eu-west-1", "ap-southeast-1"}
	jobs := []string{"api", "web", "worker", "scheduler"}
	envs := []string{"prod", "staging", "dev"}

	for i := range numSeries {
		labels := fmt.Sprintf(
			`bench_metric{job="%s", region="%s", env="%s", instance="i-%d", pod="pod-%d"`,
			jobs[i%len(jobs)],
			regions[i%len(regions)],
			envs[i%len(envs)],
			i,
			i,
		)
		// Add extra labels to increase cardinality and label materialization cost
		for j := 5; j < numLabels; j++ {
			labels += fmt.Sprintf(`, label_%d="value_%d_%d"`, j, i, j)
		}
		labels += "}"
		sb.WriteString(fmt.Sprintf("\t%s %d+%dx%d\n", labels, i%10+1, i%5+1, numSteps))
	}

	return sb.String()
}

func newProjectionSelectorBenchEngine(b *testing.B, queryable storage.Queryable, mint, maxt time.Time, engineOpts promql.EngineOpts) *engine.Engine {
	b.Helper()

	qOpts := &query.Options{
		Start:               mint,
		End:                 maxt,
		Step:                30 * time.Second,
		StepsBatch:          10,
		LookbackDelta:       5 * time.Minute,
		DecodingConcurrency: 1,
	}

	scanners, err := prometheusStorage.NewPrometheusScanners(
		queryable, qOpts, nil,
		prometheusStorage.WithSeriesHashLabel("__series_hash__"),
	)
	testutil.Ok(b, err)

	return engine.NewWithScanners(engine.Opts{
		EngineOpts: engineOpts,
		LogicalOptimizers: []logicalplan.Optimizer{
			logicalplan.SortMatchers{},
			logicalplan.ProjectionOptimizer{SeriesHashLabel: "__series_hash__"},
			logicalplan.MergeSelectsOptimizer{},
		},
	}, scanners)
}
