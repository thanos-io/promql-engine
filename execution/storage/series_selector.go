// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package storage

import (
	"context"
	"sync"

	"github.com/gogo/protobuf/types"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"

	"github.com/thanos-io/promql-engine/execution/warnings"
	"github.com/thanos-io/promql-engine/query"
	tquery "github.com/thanos-io/thanos/pkg/query"
	tstore "github.com/thanos-io/thanos/pkg/store"
	"github.com/thanos-io/thanos/pkg/store/hintspb"
)

type SeriesSelector interface {
	GetSeries(ctx context.Context, shard, numShards int) ([]SignedSeries, error)
	Matchers() []*labels.Matcher
	// TODO: add a method to get hints. Return nil if there are no hints.
	// Hints() map[string][]hintspb.SeriesHintsResponse
	GetSeriesHints() (map[string][]hintspb.SeriesResponseHints, error)
}

type SignedSeries struct {
	storage.Series
	Signature uint64
}

type seriesSelector struct {
	storage        storage.Queryable
	mint           int64
	maxt           int64
	step           int64
	matchers       []*labels.Matcher
	hints          storage.SelectHints
	opts           *query.Options
	hintsCollector *tstore.HintsCollector

	once   sync.Once
	series []SignedSeries
}

func newSeriesSelector(storage storage.Queryable, mint, maxt, step int64, matchers []*labels.Matcher, hints storage.SelectHints, opts *query.Options) *seriesSelector {
	return &seriesSelector{
		storage:  storage,
		maxt:     maxt,
		mint:     mint,
		step:     step,
		matchers: matchers,
		hints:    hints,
		opts:     opts,
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
	querier, err := o.storage.Querier(o.mint, o.maxt)
	if err != nil {
		return err
	}
	defer querier.Close()

	var seriesSet storage.SeriesSet
	var h *tstore.HintsCollector

	if o.opts.EnableAnalysis {
		if qwh, ok := querier.(tquery.QuerierWithHints); ok {
			seriesSet, h = qwh.SelectWithHints(ctx, false, &o.hints, o.matchers...)

			o.hintsCollector = h
		}
	}
	if seriesSet == nil {
		seriesSet = querier.Select(ctx, false, &o.hints, o.matchers...)
	}

	i := 0
	for seriesSet.Next() {
		s := seriesSet.At()
		o.series = append(o.series, SignedSeries{
			Series:    s,
			Signature: uint64(i),
		})
		i++
	}

	warnings.AddToContext(seriesSet.Warnings(), ctx)
	return seriesSet.Err()
}

func (o *seriesSelector) GetSeriesHints() (map[string][]hintspb.SeriesResponseHints, error) {
	if o.hintsCollector == nil {
		return nil, nil
	}

	hints := make(map[string][]hintspb.SeriesResponseHints)

	for key, value := range o.hintsCollector.Hints {
		
		for _,v := range value {
			h := hintspb.SeriesResponseHints{}
			if err := types.UnmarshalAny(v.GetHints(), &h); err != nil {
				return nil, err
			}

			hints[key] = append(hints[key], h)
		}
	}
	return hints, nil
}

func seriesShard(series []SignedSeries, index int, numShards int) []SignedSeries {
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
