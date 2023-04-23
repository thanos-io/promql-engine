// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/cortexproject/promqlsmith"
	"github.com/efficientgo/core/errors"
	"github.com/efficientgo/core/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"
	"github.com/stretchr/testify/require"

	"github.com/thanos-community/promql-engine/api"
	"github.com/thanos-community/promql-engine/engine"
	"github.com/thanos-community/promql-engine/execution/parse"
	"github.com/thanos-community/promql-engine/logicalplan"
)

const testRuns = 100

type testCase struct {
	query          string
	load           string
	oldRes, newRes *promql.Result
}

func FuzzEnginePromQLSmithRangeQuery(f *testing.F) {
	f.Add(uint32(0), uint32(120), uint32(30), 1.0, 1.0, 1.0, 2.0, 30)

	f.Fuzz(func(t *testing.T, startTS, endTS, intervalSeconds uint32, initialVal1, initialVal2, inc1, inc2 float64, stepRange int) {
		if math.IsNaN(initialVal1) || math.IsNaN(initialVal2) || math.IsNaN(inc1) || math.IsNaN(inc2) {
			return
		}
		if math.IsInf(initialVal1, 0) || math.IsInf(initialVal2, 0) || math.IsInf(inc1, 0) || math.IsInf(inc2, 0) {
			return
		}
		if inc1 < 0 || inc2 < 0 || stepRange <= 0 || intervalSeconds <= 0 || endTS < startTS {
			return
		}

		load := fmt.Sprintf(`load 30s
			http_requests_total{pod="nginx-1"} %.2f+%.2fx15
			http_requests_total{pod="nginx-2"} %2.f+%.2fx21`, initialVal1, inc1, initialVal2, inc2)

		opts := promql.EngineOpts{
			Timeout:              1 * time.Hour,
			MaxSamples:           1e10,
			EnableNegativeOffset: true,
			EnableAtModifier:     true,
		}

		test, err := promql.NewTest(t, load)
		testutil.Ok(t, err)
		defer test.Close()

		testutil.Ok(t, test.Run())

		start := time.Unix(int64(startTS), 0)
		end := time.Unix(int64(endTS), 0)
		interval := time.Duration(intervalSeconds) * time.Second

		seriesSet, err := getSeries(context.Background(), test.Storage())
		require.NoError(t, err)
		rnd := rand.New(rand.NewSource(time.Now().Unix()))
		psOpts := []promqlsmith.Option{
			promqlsmith.WithEnableOffset(true),
			promqlsmith.WithEnableAtModifier(true),
		}
		ps := promqlsmith.New(rnd, seriesSet, psOpts...)

		newEngine := engine.New(engine.Opts{EngineOpts: opts, DisableFallback: true})
		oldEngine := promql.NewEngine(opts)

		var (
			q1    promql.Query
			query string
		)
		cases := make([]*testCase, testRuns)
		for i := 0; i < testRuns; i++ {
			// Since we disabled fallback, keep trying until we find a query
			// that can be natively executed by the engine.
			for {
				expr := ps.WalkRangeQuery()
				query = expr.Pretty(0)
				q1, err = newEngine.NewRangeQuery(test.Context(), test.Storage(), nil, query, start, end, interval)
				if errors.Is(err, parse.ErrNotSupportedExpr) || errors.Is(err, parse.ErrNotImplemented) {
					continue
				} else {
					break
				}
			}

			testutil.Ok(t, err)
			newResult := q1.Exec(context.Background())

			q2, err := oldEngine.NewRangeQuery(test.Context(), test.Storage(), nil, query, start, end, interval)
			testutil.Ok(t, err)

			oldResult := q2.Exec(context.Background())

			cases[i] = &testCase{
				query:  query,
				newRes: newResult,
				oldRes: oldResult,
				load:   load,
			}
		}
		validateTestCases(t, cases)
	})
}

func validateTestCases(t *testing.T, cases []*testCase) {
	failures := 0
	for i, c := range cases {
		emptyLabelsToNil(c.newRes)
		emptyLabelsToNil(c.oldRes)

		if !cmp.Equal(c.oldRes, c.newRes, comparer) {
			t.Logf(c.load)
			t.Logf(c.query)

			t.Logf("case %d error mismatch.\nnew result: %s\nold result: %s\n", i, c.newRes.String(), c.oldRes.String())
			failures++
		}
	}
	if failures > 0 {
		t.Fatalf("failed %d test cases", failures)
	}
}

func FuzzEnginePromQLSmithInstantQuery(f *testing.F) {
	f.Add(uint32(0), 1.0, 1.0, 1.0, 2.0)

	f.Fuzz(func(t *testing.T, ts uint32, initialVal1, initialVal2, inc1, inc2 float64) {
		if inc1 < 0 || inc2 < 0 {
			return
		}
		load := fmt.Sprintf(`load 30s
			http_requests_total{pod="nginx-1", route="/"} %.2f+%.2fx4
			http_requests_total{pod="nginx-2", route="/"} %2.f+%.2fx4`, initialVal1, inc1, initialVal2, inc2)

		opts := promql.EngineOpts{
			Timeout:              1 * time.Hour,
			MaxSamples:           1e10,
			EnableNegativeOffset: true,
			EnableAtModifier:     true,
		}

		test, err := promql.NewTest(t, load)
		testutil.Ok(t, err)
		defer test.Close()
		testutil.Ok(t, test.Run())

		queryTime := time.Unix(int64(ts), 0)
		newEngine := engine.New(engine.Opts{
			EngineOpts:        opts,
			DisableFallback:   true,
			LogicalOptimizers: logicalplan.AllOptimizers,
		})
		oldEngine := promql.NewEngine(opts)

		seriesSet, err := getSeries(context.Background(), test.Storage())
		require.NoError(t, err)
		rnd := rand.New(rand.NewSource(time.Now().Unix()))
		psOpts := []promqlsmith.Option{
			promqlsmith.WithEnableOffset(true),
			promqlsmith.WithEnableAtModifier(true),
		}
		ps := promqlsmith.New(rnd, seriesSet, psOpts...)

		var (
			q1    promql.Query
			query string
		)
		cases := make([]*testCase, testRuns)
		for i := 0; i < testRuns; i++ {
			// Since we disabled fallback, keep trying until we find a query
			// that can be natively execute by the engine.
			for {
				expr := ps.WalkInstantQuery()
				query = expr.Pretty(0)
				q1, err = newEngine.NewInstantQuery(test.Context(), test.Storage(), nil, query, queryTime)
				if errors.Is(err, parse.ErrNotSupportedExpr) || errors.Is(err, parse.ErrNotImplemented) {
					continue
				} else {
					break
				}
			}

			testutil.Ok(t, err)
			newResult := q1.Exec(context.Background())

			q2, err := oldEngine.NewInstantQuery(test.Context(), test.Storage(), nil, query, queryTime)
			testutil.Ok(t, err)

			oldResult := q2.Exec(context.Background())

			cases[i] = &testCase{
				query:  query,
				newRes: newResult,
				oldRes: oldResult,
				load:   load,
			}
		}
		validateTestCases(t, cases)
	})
}

func FuzzDistributedEnginePromQLSmithRangeQuery(f *testing.F) {
	f.Add(uint32(0), uint32(120), uint32(30), 1.0, 1.0, 1.0, 1.0, 1.0, 2.0, 30)

	f.Fuzz(func(t *testing.T, startTS, endTS, intervalSeconds uint32, initialVal1, initialVal2, initialVal3, initialVal4, inc1, inc2 float64, stepRange int) {
		if math.IsNaN(initialVal1) || math.IsNaN(initialVal2) || math.IsNaN(inc1) || math.IsNaN(inc2) {
			return
		}
		if math.IsInf(initialVal1, 0) || math.IsInf(initialVal2, 0) || math.IsInf(inc1, 0) || math.IsInf(inc2, 0) {
			return
		}
		if inc1 < 0 || inc2 < 0 || stepRange <= 0 || intervalSeconds <= 0 || endTS < startTS {
			return
		}
		load := fmt.Sprintf(`load 30s
			http_requests_total{pod="nginx-1", route="/"} %.2f+%.2fx4
			http_requests_total{pod="nginx-2", route="/"} %2.f+%.2fx4`, initialVal1, inc1, initialVal2, inc2)
		load2 := fmt.Sprintf(`load 30s
			http_requests_total{pod="nginx-1", route="/"} %.2f+%.2fx4
			http_requests_total{pod="nginx-2", route="/"} %2.f+%.2fx4`, initialVal3, inc1, initialVal4, inc2)

		opts := promql.EngineOpts{
			Timeout:              1 * time.Hour,
			MaxSamples:           1e10,
			EnableNegativeOffset: true,
			EnableAtModifier:     true,
		}
		engineOpts := engine.Opts{
			EngineOpts:        opts,
			DisableFallback:   true,
			LogicalOptimizers: logicalplan.AllOptimizers,
		}

		queryables := []*promql.Test{}
		test, err := promql.NewTest(t, load)
		testutil.Ok(t, err)
		defer test.Close()
		testutil.Ok(t, test.Run())
		queryables = append(queryables, test)

		test2, err := promql.NewTest(t, load2)
		testutil.Ok(t, err)
		defer test2.Close()
		testutil.Ok(t, test2.Run())
		queryables = append(queryables, test2)

		start := time.Unix(int64(startTS), 0)
		end := time.Unix(int64(endTS), 0)
		interval := time.Duration(intervalSeconds) * time.Second

		partitionLabels := [][]labels.Labels{
			{labels.FromStrings("zone", "west-1")},
			{labels.FromStrings("zone", "west-2")},
		}
		remoteEngines := make([]api.RemoteEngine, 0, 2)
		for i := 0; i < 2; i++ {
			e := engine.NewRemoteEngine(
				engineOpts,
				queryables[i].Storage(),
				queryables[i].TSDB().Head().MinTime(),
				queryables[i].TSDB().Head().MaxTime(),
				partitionLabels[i],
			)
			remoteEngines = append(remoteEngines, e)
		}
		distEngine := engine.NewDistributedEngine(engineOpts, api.NewStaticEndpoints(remoteEngines))
		oldEngine := promql.NewEngine(opts)

		mergeStore := storage.NewFanout(nil, test.Storage(), test2.Storage())
		seriesSet, err := getSeries(context.Background(), mergeStore)
		require.NoError(t, err)
		rnd := rand.New(rand.NewSource(time.Now().Unix()))
		psOpts := []promqlsmith.Option{
			promqlsmith.WithEnableOffset(true),
			promqlsmith.WithEnableAtModifier(true),
			promqlsmith.WithEnabledAggrs([]parser.ItemType{parser.SUM, parser.MIN, parser.MAX, parser.GROUP, parser.COUNT, parser.BOTTOMK, parser.TOPK}),
		}
		ps := promqlsmith.New(rnd, seriesSet, psOpts...)

		var (
			q1    promql.Query
			query string
		)
		cases := make([]*testCase, testRuns)
		ctx := context.Background()
		for i := 0; i < testRuns; i++ {
			// Since we disabled fallback, keep trying until we find a query
			// that can be natively execute by the engine.
			for {
				expr := ps.WalkRangeQuery()
				query = expr.Pretty(0)
				q1, err = distEngine.NewRangeQuery(ctx, mergeStore, nil, query, start, end, interval)
				if errors.Is(err, parse.ErrNotSupportedExpr) || errors.Is(err, parse.ErrNotImplemented) {
					continue
				} else {
					break
				}
			}

			testutil.Ok(t, err)
			newResult := q1.Exec(ctx)

			q2, err := oldEngine.NewRangeQuery(ctx, mergeStore, nil, query, start, end, interval)
			testutil.Ok(t, err)

			oldResult := q2.Exec(ctx)

			cases[i] = &testCase{
				query:  query,
				newRes: newResult,
				oldRes: oldResult,
				load:   load,
			}
		}
		validateTestCases(t, cases)
	})
}

func FuzzDistributedEnginePromQLSmithInstantQuery(f *testing.F) {
	f.Add(uint32(0), 1.0, 1.0, 1.0, 1.0, 1.0, 2.0)

	f.Fuzz(func(t *testing.T, ts uint32, initialVal1, initialVal2, initialVal3, initialVal4, inc1, inc2 float64) {
		if inc1 < 0 || inc2 < 0 {
			return
		}
		load := fmt.Sprintf(`load 30s
			http_requests_total{pod="nginx-1", route="/"} %.2f+%.2fx4
			http_requests_total{pod="nginx-2", route="/"} %2.f+%.2fx4`, initialVal1, inc1, initialVal2, inc2)
		load2 := fmt.Sprintf(`load 30s
			http_requests_total{pod="nginx-1", route="/"} %.2f+%.2fx4
			http_requests_total{pod="nginx-2", route="/"} %2.f+%.2fx4`, initialVal3, inc1, initialVal4, inc2)

		opts := promql.EngineOpts{
			Timeout:              1 * time.Hour,
			MaxSamples:           1e10,
			EnableNegativeOffset: true,
			EnableAtModifier:     true,
		}
		engineOpts := engine.Opts{EngineOpts: opts, DisableFallback: true}

		queryables := []*promql.Test{}
		test, err := promql.NewTest(t, load)
		testutil.Ok(t, err)
		defer test.Close()
		testutil.Ok(t, test.Run())
		queryables = append(queryables, test)

		test2, err := promql.NewTest(t, load2)
		testutil.Ok(t, err)
		defer test2.Close()
		testutil.Ok(t, test2.Run())
		queryables = append(queryables, test2)

		partitionLabels := [][]labels.Labels{
			{labels.FromStrings("zone", "west-1")},
			{labels.FromStrings("zone", "west-2")},
		}
		queryTime := time.Unix(int64(ts), 0)
		remoteEngines := make([]api.RemoteEngine, 0, 2)
		for i := 0; i < 2; i++ {
			e := engine.NewRemoteEngine(
				engineOpts,
				queryables[i].Storage(),
				queryables[i].TSDB().Head().MinTime(),
				queryables[i].TSDB().Head().MaxTime(),
				partitionLabels[i],
			)
			remoteEngines = append(remoteEngines, e)
		}
		distEngine := engine.NewDistributedEngine(engineOpts, api.NewStaticEndpoints(remoteEngines))
		oldEngine := promql.NewEngine(opts)

		mergeStore := storage.NewFanout(nil, test.Storage(), test2.Storage())
		seriesSet, err := getSeries(context.Background(), mergeStore)
		require.NoError(t, err)
		rnd := rand.New(rand.NewSource(time.Now().Unix()))
		psOpts := []promqlsmith.Option{
			promqlsmith.WithEnableOffset(true),
			promqlsmith.WithEnableAtModifier(true),
			promqlsmith.WithEnabledAggrs([]parser.ItemType{parser.SUM, parser.MIN, parser.MAX, parser.GROUP, parser.COUNT, parser.BOTTOMK, parser.TOPK}),
		}
		ps := promqlsmith.New(rnd, seriesSet, psOpts...)
		ctx := context.Background()

		var (
			q1    promql.Query
			query string
		)
		cases := make([]*testCase, testRuns)
		for i := 0; i < testRuns; i++ {
			// Since we disabled fallback, keep trying until we find a query
			// that can be natively execute by the engine.
			for {
				expr := ps.Walk(parser.ValueTypeVector, parser.ValueTypeMatrix)
				query = expr.Pretty(0)
				q1, err = distEngine.NewInstantQuery(ctx, mergeStore, nil, query, queryTime)
				if errors.Is(err, parse.ErrNotSupportedExpr) || errors.Is(err, parse.ErrNotImplemented) {
					continue
				} else {
					break
				}
			}

			testutil.Ok(t, err)
			newResult := q1.Exec(ctx)

			q2, err := oldEngine.NewInstantQuery(ctx, mergeStore, nil, query, queryTime)
			testutil.Ok(t, err)

			oldResult := q2.Exec(ctx)

			cases[i] = &testCase{
				query:  query,
				newRes: newResult,
				oldRes: oldResult,
				load:   load,
			}
		}
		validateTestCases(t, cases)
	})
}

var comparer = cmp.Comparer(func(x, y *promql.Result) bool {
	compareFloats := func(l, r float64) bool {
		const epsilon = 1e-6
		return cmp.Equal(l, r, cmpopts.EquateNaNs(), cmpopts.EquateApprox(0, epsilon))
	}

	if x.Err != nil && y.Err != nil {
		return cmp.Equal(x.Err.Error(), y.Err.Error())
	} else if x.Err != nil {
		return false
	}

	vx, xvec := x.Value.(promql.Vector)
	vy, yvec := y.Value.(promql.Vector)

	if xvec && yvec {
		if len(vx) != len(vy) {
			return false
		}
		// Sort vector before comparing.
		sort.Slice(vx, func(i, j int) bool {
			return labels.Compare(vx[i].Metric, vx[j].Metric) < 0
		})
		sort.Slice(vy, func(i, j int) bool {
			return labels.Compare(vy[i].Metric, vy[j].Metric) < 0
		})
		for i := 0; i < len(vx); i++ {
			if !cmp.Equal(vx[i].Metric, vy[i].Metric) {
				return false
			}
			if vx[i].T != vy[i].T {
				return false
			}
			if !compareFloats(vx[i].F, vy[i].F) {
				return false
			}
		}
		return true
	}

	mx, xmat := x.Value.(promql.Matrix)
	my, ymat := y.Value.(promql.Matrix)

	if xmat && ymat {
		if len(mx) != len(my) {
			return false
		}
		// Sort matrix before comparing.
		sort.Sort(mx)
		sort.Sort(my)
		for i := 0; i < len(mx); i++ {
			mxs := mx[i]
			mys := my[i]

			if !cmp.Equal(mxs.Metric, mys.Metric) {
				return false
			}

			xps := mxs.Floats
			yps := mys.Floats

			if len(xps) != len(yps) {
				return false
			}
			for j := 0; j < len(xps); j++ {
				if xps[j].T != yps[j].T {
					return false
				}
				if !compareFloats(xps[j].F, yps[j].F) {
					return false
				}
			}
		}
		return true
	}

	sx, xscalar := x.Value.(promql.Scalar)
	sy, yscalar := y.Value.(promql.Scalar)
	if xscalar && yscalar {
		if sx.T != sy.T {
			return false
		}
		return compareFloats(sx.V, sy.V)
	}
	return false
})

func getSeries(ctx context.Context, q storage.Queryable) ([]labels.Labels, error) {
	querier, err := q.Querier(ctx, 0, time.Now().Unix())
	if err != nil {
		return nil, err
	}
	res := make([]labels.Labels, 0)
	ss := querier.Select(false, &storage.SelectHints{Func: "series"}, labels.MustNewMatcher(labels.MatchEqual, "__name__", "http_requests_total"))
	for ss.Next() {
		lbls := ss.At().Labels()
		res = append(res, lbls)
	}
	if err := ss.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
