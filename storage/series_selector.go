// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package storage

import (
	"context"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
)

type SeriesSelector interface {
	GetSeries(ctx context.Context, shard, numShards int) ([]SignedSeries, error)
	Matchers() []*labels.Matcher
}

type SelectorPool interface {
	GetSelector(mint, maxt, step int64, matchers []*labels.Matcher, hints storage.SelectHints) SeriesSelector
	GetFilteredSelector(mint, maxt, step int64, matchers, filters []*labels.Matcher, hints storage.SelectHints) SeriesSelector
}

type SignedSeries struct {
	storage.Series
	Signature uint64
}

func SeriesShard(series []SignedSeries, index int, numShards int) []SignedSeries {
	start := index * len(series) / numShards
	end := (index + 1) * len(series) / numShards
	if end > len(series) {
		end = len(series)
	}

	slice := series[start:end]
	shard := make([]SignedSeries, len(slice))
	copy(shard, slice)

	for i := range shard {
		shard[i].Signature = uint64(i)
	}
	return shard
}

type histogramStatsSelector struct {
	SeriesSelector
}

func NewHistogramStatsSelector(seriesSelector SeriesSelector) SeriesSelector {
	return histogramStatsSelector{SeriesSelector: seriesSelector}
}

func (h histogramStatsSelector) GetSeries(ctx context.Context, shard, numShards int) ([]SignedSeries, error) {
	series, err := h.SeriesSelector.GetSeries(ctx, shard, numShards)
	if err != nil {
		return nil, err
	}
	for i := range series {
		series[i].Series = newHistogramStatsSeries(series[i].Series)
	}
	return series, nil
}

type histogramStatsSeries struct {
	storage.Series
}

func newHistogramStatsSeries(series storage.Series) histogramStatsSeries {
	return histogramStatsSeries{Series: series}
}

func (h histogramStatsSeries) Iterator(it chunkenc.Iterator) chunkenc.Iterator {
	return NewHistogramStatsIterator(h.Series.Iterator(it))
}
