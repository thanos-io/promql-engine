// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/efficientgo/core/errors"
	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"

	"github.com/thanos-io/promql-engine/api"
	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/logicalplan"
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
		if ts < mint {
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

	makeSeriesWithName := func(name, zone, pod string) []string {
		return []string{labels.MetricName, name, "zone", zone, "pod", pod}
	}

	tests := []struct {
		name        string
		seriesSets  []partition
		timeOverlap partition
		rangeEnd    time.Time
	}{
		{
			name: "base case",
			seriesSets: []partition{{
				extLset: []labels.Labels{labels.FromStrings("zone", "east-1")},
				series: []*mockSeries{
					newMockSeries(makeSeries("east-1", "nginx-1"), []int64{30, 60, 90, 120}, []float64{2, 3, 4, 5}),
					newMockSeries(makeSeries("east-1", "nginx-2"), []int64{30, 60, 90, 120}, []float64{3, 4, 5, 6}),
				},
			}, {
				extLset: []labels.Labels{
					labels.FromStrings("zone", "west-1"),
					labels.FromStrings("zone", "west-2"),
				},
				series: []*mockSeries{
					newMockSeries(makeSeries("west-1", "nginx-1"), []int64{30, 60, 90, 120}, []float64{4, 5, 6, 7}),
					newMockSeries(makeSeries("west-1", "nginx-2"), []int64{30, 60, 90, 120}, []float64{5, 6, 7, 8}),
					newMockSeries(makeSeries("west-2", "nginx-1"), []int64{30, 60, 90, 120}, []float64{6, 7, 8, 9}),
				},
			}},
			timeOverlap: partition{
				extLset: []labels.Labels{
					labels.FromStrings("zone", "east-1"),
					labels.FromStrings("zone", "west-1"),
					labels.FromStrings("zone", "west-2"),
				},
				series: []*mockSeries{
					newMockSeries(makeSeries("east-1", "nginx-1"), []int64{30, 60}, []float64{2, 3}),
					newMockSeries(makeSeries("east-1", "nginx-2"), []int64{30, 60}, []float64{3, 4}),
					newMockSeries(makeSeries("west-1", "nginx-1"), []int64{30, 60}, []float64{4, 5}),
					newMockSeries(makeSeries("west-1", "nginx-2"), []int64{30, 60}, []float64{5, 6}),
					newMockSeries(makeSeries("west-2", "nginx-1"), []int64{30, 60}, []float64{6, 7}),
				},
			},
		},
		{
			// Repro for https://github.com/thanos-io/promql-engine/issues/187.
			name: "series with different ranges in a newer engine",
			seriesSets: []partition{{
				extLset: []labels.Labels{labels.FromStrings("zone", "east-1"), labels.FromStrings("zone", "east-1")},
				series: []*mockSeries{
					newMockSeries(makeSeries("east-1", "nginx-1"), []int64{60, 90, 120}, []float64{3, 4, 5}),
					newMockSeries(makeSeries("east-2", "nginx-1"), []int64{30, 60, 90, 120}, []float64{3, 4, 5, 6}),
				}},
			},
			timeOverlap: partition{
				extLset: []labels.Labels{labels.FromStrings("zone", "east-1"), labels.FromStrings("zone", "east-2")},
				series: []*mockSeries{
					newMockSeries(makeSeries("east-1", "nginx-1"), []int64{30, 60}, []float64{2, 3}),
					newMockSeries(makeSeries("east-2", "nginx-1"), []int64{30, 60}, []float64{3, 4}),
				},
			},
		},
		{
			name: "verify double lookback is not applied",
			seriesSets: []partition{{
				extLset: []labels.Labels{labels.FromStrings("zone", "east-2")},
				series: []*mockSeries{
					newMockSeries(makeSeries("east-2", "nginx-1"), []int64{30, 60, 90, 120}, []float64{3, 4, 5, 6}),
				}},
			},
			timeOverlap: partition{
				extLset: []labels.Labels{labels.FromStrings("zone", "east-2")},
				series: []*mockSeries{
					newMockSeries(makeSeries("east-2", "nginx-1"), []int64{30, 60}, []float64{3, 4}),
				},
			},
			rangeEnd: time.Unix(15000, 0),
		},
		{
			name: "count by __name__ label",
			seriesSets: []partition{{
				extLset: []labels.Labels{labels.FromStrings("zone", "east-2")},
				series: []*mockSeries{
					newMockSeries(makeSeriesWithName("foo", "east-2", "nginx-1"), []int64{30, 60, 90, 120}, []float64{3, 4, 5, 6}),
					newMockSeries(makeSeriesWithName("bar", "east-2", "nginx-1"), []int64{30, 60, 90, 120}, []float64{3, 4, 5, 6}),
				},
			}, {
				extLset: []labels.Labels{labels.FromStrings("zone", "east-2"), labels.FromStrings("zone", "west-1")},
				series: []*mockSeries{
					newMockSeries(makeSeriesWithName("xyz", "east-2", "nginx-1"), []int64{30, 60, 90, 120}, []float64{3, 4, 5, 6}),
				},
			}},
			timeOverlap: partition{
				series: []*mockSeries{
					newMockSeries(makeSeriesWithName("foo", "east-2", "nginx-1"), []int64{30, 60}, []float64{3, 4}),
					newMockSeries(makeSeriesWithName("bar", "east-2", "nginx-1"), []int64{30, 60}, []float64{3, 4}),
					newMockSeries(makeSeriesWithName("xyz", "east-2", "nginx-1"), []int64{30, 60}, []float64{3, 4}),
				},
			},
		},
		{
			name: "engines with different retentions",
			seriesSets: []partition{{
				extLset: []labels.Labels{labels.FromStrings("zone", "us-east1")},
				series: []*mockSeries{
					newMockSeries(makeSeries("us-east1", "nginx-1"), []int64{30, 60, 90, 120, 150}, []float64{3, 4, 5, 6, 9}),
				}}, {
				extLset: []labels.Labels{labels.FromStrings("zone", "us-east2")},
				series: []*mockSeries{
					newMockSeries(makeSeries("us-east2", "nginx-2"), []int64{90, 120, 150}, []float64{7, 9, 11}),
				},
			}},
			timeOverlap: partition{
				extLset: []labels.Labels{
					labels.FromStrings("zone", "us-east1"),
					labels.FromStrings("zone", "us-east2"),
				},
				series: []*mockSeries{
					newMockSeries(makeSeries("us-east1", "nginx-1"), []int64{30, 60, 90}, []float64{3, 4, 5}),
					newMockSeries(makeSeries("us-east2", "nginx-2"), []int64{30, 60, 90, 120}, []float64{2, 6, 7, 9}),
				},
			},
			rangeEnd: time.Unix(180, 0),
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
		{name: "count by __name__", query: `count by (__name__) ({__name__=~".+"})`},
		{name: "group", query: `group by (pod) (bar)`},
		{name: "topk", query: `topk by (pod) (1, bar)`},
		{name: "bottomk", query: `bottomk by (pod) (1, bar)`},
		{name: "label based pruning with no match", query: `sum by (pod) (bar{zone="north-2"})`},
		{name: "label based pruning with one match", query: `sum by (pod) (bar{zone="east-1"})`},
		{name: "double aggregation", query: `max by (pod) (sum by (pod) (bar))`},
		{name: "aggregation with function operand", query: `sum by (pod) (rate(bar[1m]))`},
		{name: "binary expression with constant operand", query: `sum by (region) (bar * 60)`},
		{name: "binary aggregation", query: `sum by (region) (bar) / sum by (pod) (bar)`},
		{name: "filtered selector interaction", query: `sum by (region) (bar{region="east"}) / sum by (region) (bar)`},
		{name: "unsupported aggregation", query: `count_values("pod", bar)`, expectFallback: true},
		{name: "absent_over_time for non-existing metric", query: `absent_over_time(foo[2m])`},
		{name: "absent_over_time for existing metric", query: `absent_over_time(bar{pod="nginx-1"}[2m])`},
		{name: "absent for non-existing metric", query: `absent(foo)`},
		{name: "absent for existing metric", query: `absent(bar{pod="nginx-1"})`},
	}

	optimizersOpts := map[string][]logicalplan.Optimizer{
		"none":    logicalplan.NoOptimizers,
		"default": logicalplan.DefaultOptimizers,
		"all":     logicalplan.AllOptimizers,
	}

	lookbackDeltas := []time.Duration{0, 30 * time.Second, 5 * time.Minute}
	allQueryOpts := []*promql.QueryOpts{nil}
	for _, l := range lookbackDeltas {
		allQueryOpts = append(allQueryOpts, &promql.QueryOpts{
			LookbackDelta: l,
		})
	}

	for _, test := range tests {
		for _, lookbackDelta := range lookbackDeltas {
			localOpts.LookbackDelta = lookbackDelta
			for _, queryOpts := range allQueryOpts {
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
					if len(test.timeOverlap.series) > 0 {
						remoteEngines = append(remoteEngines, engine.NewRemoteEngine(
							localOpts,
							storageWithMockSeries(test.timeOverlap.series...),
							test.timeOverlap.mint(),
							test.timeOverlap.maxt(),
							test.timeOverlap.extLset,
						))
						allSeries = append(allSeries, test.timeOverlap.series...)
					}
					completeSeriesSet := storageWithSeries(mergeWithSampleDedup(allSeries)...)

					ctx := context.Background()
					for _, query := range queries {
						t.Run(query.name, func(t *testing.T) {
							for o, optimizers := range optimizersOpts {
								t.Run(fmt.Sprintf("withOptimizers=%s", o), func(t *testing.T) {
									localOpts.LogicalOptimizers = optimizers
									distOpts := localOpts

									distOpts.DisableFallback = !query.expectFallback
									for _, instantTS := range instantTSs {
										t.Run(fmt.Sprintf("instant/ts=%d", instantTS.Unix()), func(t *testing.T) {
											distEngine := engine.NewDistributedEngine(distOpts,
												api.NewStaticEndpoints(remoteEngines),
											)
											distQry, err := distEngine.NewInstantQuery(ctx, completeSeriesSet, queryOpts, query.query, instantTS)
											testutil.Ok(t, err)

											distResult := distQry.Exec(ctx)
											promEngine := promql.NewEngine(localOpts.EngineOpts)
											promQry, err := promEngine.NewInstantQuery(ctx, completeSeriesSet, queryOpts, query.query, instantTS)
											testutil.Ok(t, err)
											promResult := promQry.Exec(ctx)

											roundValues(promResult)
											roundValues(distResult)

											// Instant queries have no guarantees on result ordering.
											sortByLabels(promResult)
											sortByLabels(distResult)

											testutil.Equals(t, promResult, distResult)
										})
									}

									t.Run("range", func(t *testing.T) {
										if test.rangeEnd == (time.Time{}) {
											test.rangeEnd = rangeEnd
										}
										distEngine := engine.NewDistributedEngine(distOpts,
											api.NewStaticEndpoints(remoteEngines),
										)
										distQry, err := distEngine.NewRangeQuery(ctx, completeSeriesSet, queryOpts, query.query, rangeStart, test.rangeEnd, rangeStep)
										testutil.Ok(t, err)

										distResult := distQry.Exec(ctx)
										promEngine := promql.NewEngine(localOpts.EngineOpts)
										promQry, err := promEngine.NewRangeQuery(ctx, completeSeriesSet, queryOpts, query.query, rangeStart, test.rangeEnd, rangeStep)
										testutil.Ok(t, err)
										promResult := promQry.Exec(ctx)

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
	}
}

func TestDistributedEngineWarnings(t *testing.T) {
	querier := &storage.MockQueryable{
		MockQuerier: &storage.MockQuerier{
			SelectMockFunction: func(sortSeries bool, hints *storage.SelectHints, matchers ...*labels.Matcher) storage.SeriesSet {
				return newWarningsSeriesSet(storage.Warnings{errors.New("test warning")})
			},
		},
	}

	opts := engine.Opts{
		EngineOpts: promql.EngineOpts{
			MaxSamples: math.MaxInt64,
			Timeout:    1 * time.Minute,
		},
	}
	remote := engine.NewRemoteEngine(opts, querier, math.MinInt64, math.MaxInt64, nil)
	ng := engine.NewDistributedEngine(opts, api.NewStaticEndpoints([]api.RemoteEngine{remote}))
	var (
		start = time.UnixMilli(0)
		end   = time.UnixMilli(600)
		step  = 30 * time.Second
	)
	q, err := ng.NewRangeQuery(context.Background(), nil, nil, "test", start, end, step)
	testutil.Ok(t, err)

	res := q.Exec(context.Background())
	testutil.Equals(t, 1, len(res.Warnings))
}
