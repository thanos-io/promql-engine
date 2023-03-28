// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/cortexproject/promqlsmith"
	"github.com/efficientgo/core/errors"
	"github.com/efficientgo/core/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"
	"github.com/stretchr/testify/require"

	"github.com/thanos-community/promql-engine/api"
	"github.com/thanos-community/promql-engine/engine"
	"github.com/thanos-community/promql-engine/execution/parse"
)

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

		var (
			q1    promql.Query
			query string
		)
		// Since we disabled fallback, keep trying until we find a query
		// that can be natively executed by the engine.
		for {
			expr := ps.WalkRangeQuery()
			query = expr.Pretty(0)
			q1, err = newEngine.NewRangeQuery(test.Storage(), nil, query, start, end, interval)
			if errors.Is(err, parse.ErrNotSupportedExpr) || errors.Is(err, parse.ErrNotImplemented) {
				continue
			} else {
				break
			}
		}

		testutil.Ok(t, err)
		t.Log(query)
		newResult := q1.Exec(context.Background())
		testutil.Ok(t, newResult.Err)

		oldEngine := promql.NewEngine(opts)
		q2, err := oldEngine.NewRangeQuery(test.Storage(), nil, query, start, end, interval)
		testutil.Ok(t, err)

		oldResult := q2.Exec(context.Background())
		testutil.Ok(t, oldResult.Err)

		emptyLabelsToNil(newResult)
		emptyLabelsToNil(oldResult)
		testutil.WithGoCmp(comparer).Equals(t, oldResult, newResult, query)
	})
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
		newEngine := engine.New(engine.Opts{EngineOpts: opts, DisableFallback: true})

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
		// Since we disabled fallback, keep trying until we find a query
		// that can be natively execute by the engine.
		for {
			expr := ps.WalkInstantQuery()
			query = expr.Pretty(0)
			q1, err = newEngine.NewInstantQuery(test.Storage(), nil, query, queryTime)
			if errors.Is(err, parse.ErrNotSupportedExpr) || errors.Is(err, parse.ErrNotImplemented) {
				continue
			} else {
				break
			}
		}

		testutil.Ok(t, err)
		t.Log(query)
		newResult := q1.Exec(context.Background())
		testutil.Ok(t, newResult.Err)

		oldEngine := promql.NewEngine(opts)
		q2, err := oldEngine.NewInstantQuery(test.Storage(), nil, query, queryTime)
		testutil.Ok(t, err)

		oldResult := q2.Exec(context.Background())
		testutil.Ok(t, oldResult.Err)

		emptyLabelsToNil(newResult)
		emptyLabelsToNil(oldResult)
		testutil.WithGoCmp(comparer).Equals(t, oldResult, newResult, query)
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
		// Since we disabled fallback, keep trying until we find a query
		// that can be natively execute by the engine.
		for {
			expr := ps.WalkRangeQuery()
			query = expr.Pretty(0)
			q1, err = distEngine.NewRangeQuery(mergeStore, nil, query, start, end, interval)
			if errors.Is(err, parse.ErrNotSupportedExpr) || errors.Is(err, parse.ErrNotImplemented) {
				continue
			} else {
				break
			}
		}

		testutil.Ok(t, err)
		t.Log(query)
		newResult := q1.Exec(context.Background())
		testutil.Ok(t, newResult.Err)

		oldEngine := promql.NewEngine(opts)
		q2, err := oldEngine.NewRangeQuery(mergeStore, nil, query, start, end, interval)
		testutil.Ok(t, err)

		oldResult := q2.Exec(context.Background())
		testutil.Ok(t, oldResult.Err)

		emptyLabelsToNil(newResult)
		emptyLabelsToNil(oldResult)
		testutil.WithGoCmp(comparer).Equals(t, oldResult, newResult, query)
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
		// Since we disabled fallback, keep trying until we find a query
		// that can be natively execute by the engine.
		for {
			expr := ps.Walk(parser.ValueTypeVector, parser.ValueTypeMatrix)
			query = expr.Pretty(0)
			q1, err = distEngine.NewInstantQuery(mergeStore, nil, query, queryTime)
			if errors.Is(err, parse.ErrNotSupportedExpr) || errors.Is(err, parse.ErrNotImplemented) {
				continue
			} else {
				break
			}
		}

		testutil.Ok(t, err)
		t.Log(query)
		newResult := q1.Exec(context.Background())
		testutil.Ok(t, newResult.Err)

		oldEngine := promql.NewEngine(opts)
		q2, err := oldEngine.NewInstantQuery(mergeStore, nil, query, queryTime)
		testutil.Ok(t, err)

		oldResult := q2.Exec(context.Background())
		testutil.Ok(t, oldResult.Err)

		emptyLabelsToNil(newResult)
		emptyLabelsToNil(oldResult)
		testutil.WithGoCmp(comparer).Equals(t, oldResult, newResult, query)
	})
}

var comparer = cmp.Comparer(func(x, y *promql.Result) bool {
	compareFloats := func(l, r float64) bool {
		const epsilon = 1e-6

		if math.IsNaN(l) && math.IsNaN(r) {
			return true
		}
		if math.IsNaN(l) || math.IsNaN(r) {
			return false
		}

		return math.Abs(l-r) < epsilon
	}

	if x.Err != y.Err {
		return false
	}

	vx, xvec := x.Value.(promql.Vector)
	vy, yvec := y.Value.(promql.Vector)

	if xvec && yvec {
		if len(vx) != len(vy) {
			return false
		}
		for i := 0; i < len(vx); i++ {
			if !cmp.Equal(vx[i].Metric, vy[i].Metric) {
				return false
			}
			if vx[i].T != vy[i].T {
				return false
			}
			if !compareFloats(vx[i].V, vy[i].V) {
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
		for i := 0; i < len(mx); i++ {
			mxs := mx[i]
			mys := my[i]

			if !cmp.Equal(mxs.Metric, mys.Metric) {
				return false
			}

			xps := mxs.Points
			yps := mys.Points

			if len(xps) != len(yps) {
				return false
			}
			for j := 0; j < len(xps); j++ {
				if xps[j].T != yps[j].T {
					return false
				}
				if !compareFloats(xps[j].V, yps[j].V) {
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
