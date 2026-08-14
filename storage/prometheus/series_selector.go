// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus

import (
	"context"
	"sync"

	"github.com/thanos-io/promql-engine/warnings"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
)

type SeriesSelector interface {
	GetSeries(ctx context.Context, shard, numShards int) ([]SignedSeries, error)
	Matchers() []*labels.Matcher
}

type SignedSeries struct {
	storage.Series
	// OriginHash is the hash of the series' original label set, before a
	// projection trimmed it. It is zero when the selection did not request a
	// projection or the series set returned by the storage does not implement
	// AtHash() uint64.
	OriginHash uint64
}

type seriesSelector struct {
	storage  storage.Querier
	matchers []*labels.Matcher
	hints    storage.SelectHints

	once   sync.Once
	series []SignedSeries
}

func newSeriesSelector(storage storage.Querier, matchers []*labels.Matcher, hints storage.SelectHints) *seriesSelector {
	return &seriesSelector{
		storage:  storage,
		matchers: matchers,
		hints:    hints,
	}
}

func (o *seriesSelector) Matchers() []*labels.Matcher {
	return o.matchers
}

func (o *seriesSelector) GetSeries(ctx context.Context, shard int, numShards int) ([]SignedSeries, error) {
	var err error
	o.once.Do(func() { err = o.loadSeries(ctx) })
	if err != nil {
		return nil, err
	}

	return seriesShard(o.series, shard, numShards), nil
}

func (o *seriesSelector) loadSeries(ctx context.Context) error {
	seriesSet := o.storage.Select(ctx, false, &o.hints, o.matchers...)
	// Origin hashes are only needed to disambiguate series when a projection
	// could have trimmed distinct label sets down to identical ones. Without a
	// projection, series must keep being identified by their label sets so
	// that genuine duplicates are still detected.
	projected := o.hints.ProjectionInclude || len(o.hints.ProjectionLabels) > 0
	hashSet, hasHashes := seriesSet.(interface{ AtHash() uint64 })
	hasHashes = hasHashes && projected
	for seriesSet.Next() {
		s := SignedSeries{
			Series: seriesSet.At(),
		}
		if hasHashes {
			s.OriginHash = hashSet.AtHash()
		}
		o.series = append(o.series, s)
	}

	for _, w := range seriesSet.Warnings() {
		warnings.AddToContext(w, ctx)
	}
	return seriesSet.Err()
}

// originHashes returns the origin hashes of the given series, or nil when the
// storage did not provide any.
func originHashes(series []SignedSeries) []uint64 {
	var found bool
	for i := range series {
		if series[i].OriginHash != 0 {
			found = true
			break
		}
	}
	if !found {
		return nil
	}
	hashes := make([]uint64, len(series))
	for i, s := range series {
		hashes[i] = s.OriginHash
	}
	return hashes
}

func seriesShard(series []SignedSeries, index int, numShards int) []SignedSeries {
	start := index * len(series) / numShards
	end := min((index+1)*len(series)/numShards, len(series))

	slice := series[start:end]
	shard := make([]SignedSeries, len(slice))
	copy(shard, slice)

	return shard
}
