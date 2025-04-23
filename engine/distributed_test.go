// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/api"
	"github.com/thanos-io/promql-engine/engine"

	"github.com/efficientgo/core/errors"
	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/util/annotations"
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
	t.Parallel()

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
		name       string
		query      string
		rangeStart time.Time
	}{
		{name: "binop with selector and constant series", query: `bar or on () vector(0)`},
		{name: "binop with aggregation and constant series", query: `sum(bar) or on () vector(0)`},
		{name: "sum", query: `sum by (pod) (bar)`},
		{name: "sum by __name__", query: `sum by (__name__) ({__name__=~".+"})`},
		{name: "parenthesis", query: `sum by (pod) ((bar))`},
		{name: "avg", query: `avg(bar)`},
		{name: "avg by __name__", query: `avg by (__name__) ({__name__=~".+"})`},
		{name: "avg with by-grouping", query: `avg by (pod) (bar)`},
		{name: "avg with without-grouping", query: `avg without (pod) (bar)`},
		{name: "label_replace", query: `max by (instance) (label_replace(bar, "instance", "$1", "pod", ".*"))`},
		{name: "label_replace to ext label before aggregation", query: `max(sum(label_replace(bar, "zone", "hardcoded-zone", "zone", "(.*)")))`},
		{name: "label_replace to ext label after aggregation", query: `max(label_replace(sum by (zone) (bar), "zone", "hardcoded-zone", "zone", ".*"))`},
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
		{name: "binary expression with distributable pairing", query: `sum by (pod) (bar * bar)`},
		{name: "binary expression with non-distributable pairing", query: `sum by (pod) (foo * ignoring (region, zone) xyz)`},
		{name: "binary aggregation", query: `sum by (region) (bar) / sum by (pod) (bar)`},
		{name: "binary nested with constants", query: `(1 + 2) + (1 atan2 (-1 % -1))`},
		{name: "binary nested with functions", query: `(1 + exp(vector(1))) + (1 atan2 (-1 % -1))`},
		{name: "filtered selector interaction", query: `sum by (region) (bar{region="east"}) / sum by (region) (bar)`},
		{name: "unsupported aggregation", query: `count_values("pod", bar)`},
		{name: "absent_over_time for non-existing metric", query: `absent_over_time(foo[2m])`},
		{name: "absent_over_time for existing metric", query: `absent_over_time(bar{pod="nginx-1"}[2m])`},
		{name: "absent for non-existing metric", query: `absent(foo)`},
		{name: "absent for existing metric with aggregation", query: `sum(absent(foo))`},
		{name: "absent for existing metric", query: `absent(bar{pod="nginx-1"})`},
		{name: "absent for existing metric with aggregation", query: `sum(absent(bar{pod="nginx-1"}))`},
		{name: "subquery with sum/count", query: `max_over_time((sum(bar) / count(bar))[30s:15s])`},
		{name: "subquery with avg", query: `max_over_time(avg(bar)[30s:15s])`},
		{name: "subquery with window within engine range", query: `max_over_time(sum_over_time(bar[30s])[30s:15s])`},
		{name: "subquery with window outside of engine range", query: `max_over_time(sum_over_time(bar[1m])[10m:1m])`},
		{name: "subquery with misaligned ranges", rangeStart: time.Unix(7, 0), query: `max_over_time(sum(bar)[30s:15s])`},
		{name: "subquery with misaligned ranges", rangeStart: time.Unix(7, 0), query: `max_over_time(sum(sum(bar))[30s:15s])`},
		{name: "nested subqueries",
			rangeStart: time.Unix(7, 0),
			query:      `max_over_time(min_over_time(sum(bar)[15s:15s])[15s:15s])`,
		},
		{name: "subquery over distributed binary expression", query: `max_over_time((bar / bar)[30s:15s])`},
		{name: "timestamp", query: `timestamp(bar)`},
		{name: "timestamp - step invariant", query: `timestamp(bar @ 6000.000)`},
		{name: "query with @start() absolute timestamp", query: `sum(bar @ start())`},
		{name: "query with @end() timestamp", query: `sum(bar @ end())`},
		{name: "query with numeric timestamp", query: `sum(bar @ 140.000)`},
		{name: "query with range and @end() timestamp", query: `sum(count_over_time(bar[1h] @ end()))`},
		{name: `subquery with @end() timestamp`, query: `bar @ 100.000 - bar @ 150.000`},
		{name: "limitk", query: `limitk by (pod,zone) (1, bar)`}, // this is more of a condition when o/p will be consistent with prometheus, in certain cases
		// the root querier wouldn't have idea so as to which samples would've been first from the those coming from the leaf queriers
	}

	lookbackDeltas := []time.Duration{0, 30 * time.Second, 5 * time.Minute}
	allQueryOpts := []promql.QueryOpts{nil}
	for _, l := range lookbackDeltas {
		allQueryOpts = append(allQueryOpts, promql.NewPrometheusQueryOpts(false, l))
	}

	for _, query := range queries {
		t.Run(query.name, func(t *testing.T) {
			t.Parallel()
			for _, test := range tests {
				var allSeries []*mockSeries
				remoteEngines := make([]api.RemoteEngine, 0, len(test.seriesSets)+1)
				for _, s := range test.seriesSets {
					allSeries = append(allSeries, s.series...)
				}
				if len(test.timeOverlap.series) > 0 {
					allSeries = append(allSeries, test.timeOverlap.series...)
				}
				completeSeriesSet := storageWithSeries(mergeWithSampleDedup(allSeries)...)
				t.Run(test.name, func(t *testing.T) {
					for _, lookbackDelta := range lookbackDeltas {
						opts := engine.Opts{
							EngineOpts: promql.EngineOpts{
								Timeout:              1 * time.Hour,
								MaxSamples:           1e10,
								EnableNegativeOffset: true,
								EnableAtModifier:     true,
								LookbackDelta:        lookbackDelta,
							},
						}

						for _, s := range test.seriesSets {
							remoteEngines = append(remoteEngines, engine.NewRemoteEngine(
								opts,
								storageWithMockSeries(s.series...),
								s.mint(),
								s.maxt(),
								s.extLset,
							))
						}
						if len(test.timeOverlap.series) > 0 {
							remoteEngines = append(remoteEngines, engine.NewRemoteEngine(
								opts,
								storageWithMockSeries(test.timeOverlap.series...),
								test.timeOverlap.mint(),
								test.timeOverlap.maxt(),
								test.timeOverlap.extLset,
							))
						}
						endpoints := api.NewStaticEndpoints(remoteEngines)
						for _, queryOpts := range allQueryOpts {
							ctx := context.Background()
							for _, instantTS := range instantTSs {
								t.Run(fmt.Sprintf("instant/ts=%d", instantTS.Unix()), func(t *testing.T) {
									distEngine := engine.NewDistributedEngine(opts)
									distQry, err := distEngine.MakeInstantQuery(ctx, completeSeriesSet, endpoints, queryOpts, query.query, instantTS)
									testutil.Ok(t, err)

									distResult := distQry.Exec(ctx)
									promEngine := promql.NewEngine(opts.EngineOpts)
									promQry, err := promEngine.NewInstantQuery(ctx, completeSeriesSet, queryOpts, query.query, instantTS)
									testutil.Ok(t, err)
									promResult := promQry.Exec(ctx)

									testutil.WithGoCmp(comparer).Equals(t, promResult, distResult, queryExplanation(distQry))
								})
							}

							t.Run("range", func(t *testing.T) {
								if query.rangeStart == (time.Time{}) {
									query.rangeStart = rangeStart
								}
								if test.rangeEnd == (time.Time{}) {
									test.rangeEnd = rangeEnd
								}
								distEngine := engine.NewDistributedEngine(opts)
								distQry, err := distEngine.MakeRangeQuery(ctx, completeSeriesSet, endpoints, queryOpts, query.query, query.rangeStart, test.rangeEnd, rangeStep)
								testutil.Ok(t, err)

								distResult := distQry.Exec(ctx)
								promEngine := promql.NewEngine(opts.EngineOpts)
								promQry, err := promEngine.NewRangeQuery(ctx, completeSeriesSet, queryOpts, query.query, query.rangeStart, test.rangeEnd, rangeStep)
								testutil.Ok(t, err)
								promResult := promQry.Exec(ctx)

								testutil.WithGoCmp(comparer).Equals(t, promResult, distResult, queryExplanation(distQry))
							})
						}
					}
				})
			}
		})
	}
}

func TestDistributedEngineWarnings(t *testing.T) {
	t.Parallel()

	opts := engine.Opts{
		EngineOpts: promql.EngineOpts{
			MaxSamples: math.MaxInt64,
			Timeout:    1 * time.Minute,
		},
	}

	querier := &storage.MockQueryable{
		MockQuerier: &storage.MockQuerier{
			SelectMockFunction: func(sortSeries bool, hints *storage.SelectHints, matchers ...*labels.Matcher) storage.SeriesSet {
				return newWarningsSeriesSet(annotations.New().Add(errors.New("test warning")))
			},
		},
	}
	remote := engine.NewRemoteEngine(opts, querier, math.MinInt64, math.MaxInt64, nil)
	endpoints := api.NewStaticEndpoints([]api.RemoteEngine{remote})
	ng := engine.NewDistributedEngine(opts)
	q, err := ng.MakeInstantQuery(context.Background(), querier, endpoints, nil, "test", time.UnixMilli(0))
	testutil.Ok(t, err)

	res := q.Exec(context.Background())
	testutil.Equals(t, 1, len(res.Warnings))
}
