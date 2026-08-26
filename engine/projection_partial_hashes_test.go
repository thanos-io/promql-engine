// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/logicalplan"

	"github.com/efficientgo/core/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/promqltest"
	promstorage "github.com/prometheus/prometheus/storage"
)

// partialHashQueryable trims label sets per the projection hints but exposes an
// origin hash only for series whose pod label is not in hashless. It models a
// fanout where some stores provide Series.original_labels_hash and others do
// not (a mixed-version leaf fleet).
type partialHashQueryable struct {
	promstorage.Queryable
	hashless []string
}

func (q *partialHashQueryable) Querier(mint, maxt int64) (promstorage.Querier, error) {
	inner, err := q.Queryable.Querier(mint, maxt)
	if err != nil {
		return nil, err
	}
	return &partialHashQuerier{Querier: inner, hashless: q.hashless}, nil
}

type partialHashQuerier struct {
	promstorage.Querier
	hashless []string
}

func (q *partialHashQuerier) Select(ctx context.Context, sorted bool, hints *promstorage.SelectHints, ms ...*labels.Matcher) promstorage.SeriesSet {
	set := q.Querier.Select(ctx, sorted, hints, ms...)
	if hints == nil || (!hints.ProjectionInclude && len(hints.ProjectionLabels) == 0) {
		return set
	}
	return &partialHashSeriesSet{SeriesSet: set, hints: hints, hashless: q.hashless}
}

type partialHashSeriesSet struct {
	promstorage.SeriesSet
	hints    *promstorage.SelectHints
	hashless []string
}

func (s *partialHashSeriesSet) At() promstorage.Series {
	series := s.SeriesSet.At()
	lb := labels.NewBuilder(series.Labels())
	if s.hints.ProjectionInclude {
		lb.Keep(s.hints.ProjectionLabels...)
	} else {
		lb.Del(s.hints.ProjectionLabels...)
	}
	return &projectedSeries{Series: series, lset: lb.Labels()}
}

func (s *partialHashSeriesSet) AtHash() uint64 {
	lset := s.SeriesSet.At().Labels()
	if slices.Contains(s.hashless, lset.Get("pod")) {
		return 0
	}
	return labels.StableHash(lset)
}

// TestProjectionPartialOriginHashes checks that series which do carry an origin
// hash keep being identified by it when other series in the same selection do
// not. Series are sharded across selectors by index, so a shard that holds only
// hashless series must not cost the other shards their identity.
func TestProjectionPartialOriginHashes(t *testing.T) {
	t.Parallel()

	// Sorted by label set, the first two series are the hashless ones and they
	// stay distinct under the projection; the remaining pairs collide.
	load := `load 30s
		http_requests_total{env="dev", job="a", pod="hashless-1"} 1+1x40
		http_requests_total{env="dev", job="b", pod="hashless-2"} 2+2x40
		http_requests_total{env="dev", job="c", pod="p1"} 3+3x40
		http_requests_total{env="dev", job="c", pod="p2"} 4+4x40
		http_requests_total{env="prod", job="c", pod="p1"} 5+5x40
		http_requests_total{env="prod", job="c", pod="p2"} 6+6x40
		http_requests_total{env="prod", job="d", pod="p1"} 7+7x40
		http_requests_total{env="prod", job="d", pod="p2"} 8+8x40`

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	engineOpts := promql.EngineOpts{Timeout: time.Minute, MaxSamples: 1e10}
	normalEngine := engine.New(engine.Opts{
		EngineOpts:        engineOpts,
		LogicalOptimizers: logicalplan.DefaultOptimizers,
	})
	projectionEngine := engine.New(engine.Opts{
		EngineOpts:          engineOpts,
		LogicalOptimizers:   append(append([]logicalplan.Optimizer{}, logicalplan.DefaultOptimizers...), logicalplan.ProjectionOptimizer{}),
		DecodingConcurrency: 4,
	})

	ctx := context.Background()
	queryTime := time.Unix(600, 0)
	query := `sum by (env, job) (abs(http_requests_total))`

	for _, tc := range []struct {
		name     string
		hashless []string
	}{
		{name: "all series carry a hash"},
		{name: "one shard is entirely hashless", hashless: []string{"hashless-1", "hashless-2"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			normalQuery, err := normalEngine.NewInstantQuery(ctx, storage, &engine.QueryOpts{}, query, queryTime)
			testutil.Ok(t, err)
			defer normalQuery.Close()
			normalResult := normalQuery.Exec(ctx)
			testutil.Ok(t, normalResult.Err)

			projectionQuery, err := projectionEngine.MakeInstantQuery(ctx, &partialHashQueryable{Queryable: storage, hashless: tc.hashless}, &engine.QueryOpts{}, query, queryTime)
			testutil.Ok(t, err)
			defer projectionQuery.Close()
			projectionResult := projectionQuery.Exec(ctx)
			if projectionResult.Err != nil {
				t.Fatalf("projected query failed while the hashless series stay distinct: %v", projectionResult.Err)
			}
			if diff := cmp.Diff(normalResult, projectionResult, comparer); diff != "" {
				t.Errorf("results differ:\nnormal:    %v\nprojected: %v", normalResult.Value, projectionResult.Value)
			}
		})
	}
}
