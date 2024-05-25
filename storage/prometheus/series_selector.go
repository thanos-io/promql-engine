// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus

import (
	"context"
	"sync"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"

	"github.com/thanos-io/promql-engine/execution/warnings"
	engstorage "github.com/thanos-io/promql-engine/storage"
)

type seriesSelector struct {
	storage  storage.Queryable
	mint     int64
	maxt     int64
	step     int64
	matchers []*labels.Matcher
	hints    storage.SelectHints

	once   sync.Once
	series []engstorage.SignedSeries
}

func newSeriesSelector(storage storage.Queryable, mint, maxt, step int64, matchers []*labels.Matcher, hints storage.SelectHints) *seriesSelector {
	return &seriesSelector{
		storage:  storage,
		maxt:     maxt,
		mint:     mint,
		step:     step,
		matchers: matchers,
		hints:    hints,
	}
}

func (o *seriesSelector) Matchers() []*labels.Matcher {
	return o.matchers
}

func (o *seriesSelector) GetSeries(ctx context.Context, shard int, numShards int) ([]engstorage.SignedSeries, error) {
	var err error
	o.once.Do(func() { err = o.loadSeries(ctx) })
	if err != nil {
		return nil, err
	}

	return engstorage.SeriesShard(o.series, shard, numShards), nil
}

func (o *seriesSelector) loadSeries(ctx context.Context) error {
	querier, err := o.storage.Querier(o.mint, o.maxt)
	if err != nil {
		return err
	}
	defer querier.Close()

	seriesSet := querier.Select(ctx, false, &o.hints, o.matchers...)
	i := 0
	for seriesSet.Next() {
		s := seriesSet.At()
		o.series = append(o.series, engstorage.SignedSeries{
			Series:    s,
			Signature: uint64(i),
		})
		i++
	}

	warnings.AddToContext(seriesSet.Warnings(), ctx)
	return seriesSet.Err()
}
