// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"

	"github.com/thanos-community/promql-engine/api"
	"github.com/thanos-community/promql-engine/engine"
	"github.com/thanos-community/promql-engine/logicalplan"
)

type partition struct {
	series  []*mockSeries
	extLset []labels.Labels
}

func (p partition) maxt() int64 {
	var maxt int64 = math.MinInt64
	for _, s := range p.series {
		ts := s.timestamps[len(s.timestamps)-1]
		if ts > maxt {
			maxt = ts
		}
	}

	return maxt
}

func (p partition) mint() int64 {
	mint := p.series[0].timestamps[0]
	for _, s := range p.series {
		ts := s.timestamps[0]
		if ts > mint {
			mint = ts
		}
	}

	return mint
}

func TestDistributedAggregations(t *testing.T) {
	localOpts := engine.Opts{
		EngineOpts: promql.EngineOpts{
			Timeout:              1 * time.Hour,
			MaxSamples:           1e10,
			EnableNegativeOffset: true,
			EnableAtModifier:     true,
		},
	}

	instantTSs := []time.Time{
		time.Unix(75, 0),
		time.Unix(121, 0),
		time.Unix(600, 0),
	}
	rangeStart := time.Unix(0, 0)
	rangeEnd := time.Unix(120, 0)
	rangeStep := time.Second * 30

	makeSeries := func(zone, pod string) []string {
		return []string{labels.MetricName, "bar", "zone", zone, "pod", pod}
	}

	tests := []struct {
		name        string
		seriesSets  []partition
		timeOverlap partition
	}{
		{
			name: "base case",
			seriesSets: []partition{
				{
					extLset: []labels.Labels{labels.FromStrings("zone", "east-1")},
					series: []*mockSeries{
						newMockSeries(makeSeries("east-1", "nginx-1"), []int64{30, 60, 90, 120}, []float64{2, 3, 4, 5}),
						newMockSeries(makeSeries("east-1", "nginx-2"), []int64{30, 60, 90, 120}, []float64{3, 4, 5, 6}),
					},
				},
				{
					extLset: []labels.Labels{
						labels.FromStrings("zone", "west-1"),
						labels.FromStrings("zone", "west-2"),
					},
					series: []*mockSeries{
						newMockSeries(makeSeries("west-1", "nginx-1"), []int64{30, 60, 90, 120}, []float64{4, 5, 6, 7}),
						newMockSeries(makeSeries("west-1", "nginx-2"), []int64{30, 60, 90, 120}, []float64{5, 6, 7, 8}),
						newMockSeries(makeSeries("west-2", "nginx-1"), []int64{30, 60, 90, 120}, []float64{6, 7, 8, 9}),
					},
				},
			},
			timeOverlap: partition{
				extLset: []labels.Labels{
					labels.FromStrings("zone", "east-1"),
					labels.FromStrings("zone", "west-1"),
					labels.FromStrings("zone", "west-2"),
				},
				series: []*mockSeries{
					newMockSeries(makeSeries("east-1", "nginx-1"), []int64{30, 60}, []float64{2, 3}),
					newMockSeries(makeSeries("west-1", "nginx-2"), []int64{30, 60}, []float64{5, 6}),
					newMockSeries(makeSeries("west-2", "nginx-1"), []int64{30, 60}, []float64{6, 7}),
				},
			},
		},
		{
			// Repro for https://github.com/thanos-community/promql-engine/issues/187.
			name: "series with different ranges in a newer engine",
			seriesSets: []partition{
				{
					series: []*mockSeries{
						newMockSeries(makeSeries("east-1", "nginx-1"), []int64{60, 90, 120}, []float64{3, 4, 5}),
						newMockSeries(makeSeries("east-2", "nginx-1"), []int64{30, 60, 90, 120}, []float64{3, 4, 5, 6}),
					},
				},
			},
			timeOverlap: partition{
				series: []*mockSeries{
					newMockSeries(makeSeries("east-1", "nginx-1"), []int64{30, 60}, []float64{2, 3}),
					newMockSeries(makeSeries("east-2", "nginx-1"), []int64{30, 60}, []float64{3, 4}),
				},
			},
		},
	}

	queries := []struct {
		name           string
		query          string
		expectFallback bool
	}{
		{name: "sum", query: `sum by (pod) (bar)`},
		{name: "avg", query: `avg by (pod) (bar)`},
		{name: "count", query: `count by (pod) (bar)`},
		{name: "group", query: `group by (pod) (bar)`},
		{name: "topk", query: `topk by (pod) (1, bar)`},
		{name: "bottomk", query: `bottomk by (pod) (1, bar)`},
		{name: "double aggregation", query: `max by (pod) (sum by (pod) (bar))`},
		// TODO(fpetkovski): This query fails because the range selector is longer than the
		// retention of one engine. Uncomment the test once the issue is fixed.
		// {name: "aggregation with function operand", query: `sum by (pod) (rate(bar[1m]))`},
		{name: "binary aggregation", query: `sum by (region) (bar) / sum by (pod) (bar)`},
		{name: "filtered selector interaction", query: `sum by (region) (bar{region="east"}) / sum by (region) (bar)`},
		{name: "unsupported aggregation", query: `count_values("pod", bar)`, expectFallback: true},
	}

	optimizersOpts := map[string][]logicalplan.Optimizer{
		"none":    logicalplan.NoOptimizers,
		"default": logicalplan.DefaultOptimizers,
		"all":     logicalplan.AllOptimizers,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var allSeries []*mockSeries
			remoteEngines := make([]api.RemoteEngine, 0, len(test.seriesSets)+1)
			for _, s := range test.seriesSets {
				remoteEngines = append(remoteEngines, engine.NewRemoteEngine(
					localOpts,
					storageWithMockSeries(s.series...),
					s.mint(),
					s.maxt(),
					s.extLset,
				))
				allSeries = append(allSeries, s.series...)
			}
			remoteEngines = append(remoteEngines, engine.NewRemoteEngine(
				localOpts,
				storageWithMockSeries(test.timeOverlap.series...),
				test.timeOverlap.mint(),
				test.timeOverlap.maxt(),
				test.timeOverlap.extLset,
			))
			allSeries = append(allSeries, test.timeOverlap.series...)
			completeSeriesSet := storageWithSeries(mergeWithSampleDedup(allSeries)...)

			for _, query := range queries {
				t.Run(query.name, func(t *testing.T) {
					for o, optimizers := range optimizersOpts {
						t.Run(fmt.Sprintf("withOptimizers=%s", o), func(t *testing.T) {
							localOpts.LogicalOptimizers = optimizers
							distOpts := localOpts
							distOpts.DisableFallback = !query.expectFallback
							distOpts.DebugWriter = os.Stdout
							for _, instantTS := range instantTSs {
								t.Run(fmt.Sprintf("instant/ts=%d", instantTS.Unix()), func(t *testing.T) {
									distEngine := engine.NewDistributedEngine(distOpts,
										api.NewStaticEndpoints(remoteEngines),
									)
									distQry, err := distEngine.NewInstantQuery(completeSeriesSet, nil, query.query, instantTS)
									testutil.Ok(t, err)

									distResult := distQry.Exec(context.Background())
									promEngine := promql.NewEngine(localOpts.EngineOpts)
									promQry, err := promEngine.NewInstantQuery(completeSeriesSet, nil, query.query, instantTS)
									testutil.Ok(t, err)
									promResult := promQry.Exec(context.Background())

									roundValues(promResult)
									roundValues(distResult)

									// Instant queries have no guarantees on result ordering.
									sortByLabels(promResult)
									sortByLabels(distResult)

									testutil.Equals(t, promResult, distResult)
								})
							}

							t.Run("range", func(t *testing.T) {
								distEngine := engine.NewDistributedEngine(distOpts,
									api.NewStaticEndpoints(remoteEngines),
								)
								distQry, err := distEngine.NewRangeQuery(completeSeriesSet, nil, query.query, rangeStart, rangeEnd, rangeStep)
								testutil.Ok(t, err)

								distResult := distQry.Exec(context.Background())
								promEngine := promql.NewEngine(localOpts.EngineOpts)
								promQry, err := promEngine.NewRangeQuery(completeSeriesSet, nil, query.query, rangeStart, rangeEnd, rangeStep)
								testutil.Ok(t, err)
								promResult := promQry.Exec(context.Background())

								roundValues(promResult)
								roundValues(distResult)
								testutil.Equals(t, promResult, distResult)
							})
						})
					}
				})
			}
		})
	}
}
