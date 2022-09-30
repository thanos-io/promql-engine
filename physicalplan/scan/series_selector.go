// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package scan

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
)

type signedSeries struct {
	storage.Series
	signature uint64
}

type seriesSelector struct {
	storage     storage.Queryable
	mint        int64
	maxt        int64
	selectRange int64
	matchers    []*labels.Matcher

	once sync.Once

	series []signedSeries
}

func NewSeriesFilter(storage storage.Queryable, mint, maxt time.Time, selectRange, lookbackDelta time.Duration, matchers []*labels.Matcher) *seriesSelector {
	return &seriesSelector{
		storage: storage,

		mint:        mint.UnixMilli() - lookbackDelta.Milliseconds(),
		maxt:        maxt.UnixMilli(),
		selectRange: selectRange.Milliseconds(),
		matchers:    matchers,
	}
}

func (o *seriesSelector) getSeries(ctx context.Context, shard int, numShards int) ([]signedSeries, error) {
	var err error
	o.once.Do(func() { err = o.loadSeries(ctx) })
	if err != nil {
		return nil, err
	}

	start := shard * len(o.series) / numShards
	end := (shard + 1) * len(o.series) / numShards
	if end > len(o.series) {
		end = len(o.series)
	}
	return o.series[start:end], nil

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
		o.series = append(o.series, signedSeries{
			Series:    s,
			signature: uint64(i),
		})
		i++
	}

	return nil
}
