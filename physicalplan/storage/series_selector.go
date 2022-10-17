// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package storage

import (
	"context"
	"sync"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
)

type SeriesSelector interface {
	GetSeries(ctx context.Context, shard, numShards int) ([]SignedSeries, error)
	Matchers() []*labels.Matcher
}

type SignedSeries struct {
	storage.Series
	Signature uint64
}

type seriesSelector struct {
	storage  storage.Queryable
	mint     int64
	maxt     int64
	matchers []*labels.Matcher

	once sync.Once

	series []SignedSeries
}

func newSeriesSelector(storage storage.Queryable, mint, maxt int64, matchers []*labels.Matcher) *seriesSelector {
	return &seriesSelector{
		storage:  storage,
		mint:     mint,
		maxt:     maxt,
		matchers: matchers,
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
	querier, err := o.storage.Querier(ctx, o.mint, o.maxt)
	if err != nil {
		return err
	}
	defer querier.Close()

	seriesSet := querier.Select(false, nil, o.matchers...)
	i := 0
	for seriesSet.Next() {
		s := seriesSet.At()
		o.series = append(o.series, SignedSeries{
			Series:    s,
			Signature: uint64(i),
		})
		i++
	}

	return nil
}

func seriesShard(series []SignedSeries, shard int, numShards int) []SignedSeries {
	start := shard * len(series) / numShards
	end := (shard + 1) * len(series) / numShards
	if end > len(series) {
		end = len(series)
	}
	return series[start:end]
}
