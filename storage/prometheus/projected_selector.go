// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus

import (
	"context"
	"slices"
	"strconv"
	"sync"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb/chunkenc"

	"github.com/thanos-io/promql-engine/logicalplan"
)

// projectedSelector wraps a SeriesSelector and applies label projection
// on GetSeries(). This reduces memory by only materializing the labels
// that downstream operators actually need (e.g. grouping labels for aggregations).
type projectedSelector struct {
	selector        SeriesSelector
	projection      *logicalplan.Projection
	seriesHashLabel string

	once   sync.Once
	series []SignedSeries
}

// NewProjectedSelector creates a new projected selector that wraps the given selector
// and applies label projection based on the provided projection hints.
// seriesHashLabel is the label name used to store the original series hash
// for identity preservation (typically "__series_hash__").
// If seriesHashLabel is empty, projection is not applied at the selector level
// (the remote storage or queryable wrapper is expected to handle it instead).
func NewProjectedSelector(selector SeriesSelector, projection *logicalplan.Projection, seriesHashLabel string) SeriesSelector {
	if projection == nil || isNoOpProjection(projection) || seriesHashLabel == "" {
		return selector
	}
	return &projectedSelector{
		selector:        selector,
		projection:      projection,
		seriesHashLabel: seriesHashLabel,
	}
}

// isNoOpProjection returns true if the projection would not actually remove any labels.
// An exclude-mode projection with no labels to exclude is a no-op.
func isNoOpProjection(p *logicalplan.Projection) bool {
	return !p.Include && len(p.Labels) == 0
}

func (p *projectedSelector) Matchers() []*labels.Matcher {
	return p.selector.Matchers()
}

func (p *projectedSelector) GetSeries(ctx context.Context, shard, numShards int) ([]SignedSeries, error) {
	var err error
	p.once.Do(func() { err = p.loadSeries(ctx) })
	if err != nil {
		return nil, err
	}

	return seriesShard(p.series, shard, numShards), nil
}

func (p *projectedSelector) loadSeries(ctx context.Context) error {
	series, err := p.selector.GetSeries(ctx, 0, 1)
	if err != nil {
		return err
	}

	p.series = make([]SignedSeries, 0, len(series))
	for i, s := range series {
		projected := p.projectSeries(s.Series)
		p.series = append(p.series, SignedSeries{
			Series:    projected,
			Signature: uint64(i),
		})
	}

	return nil
}

func (p *projectedSelector) projectSeries(s storage.Series) storage.Series {
	originalLabels := s.Labels()
	projectedLabels := p.projectLabels(originalLabels)

	if labels.Equal(originalLabels, projectedLabels) {
		return s
	}

	return &projectedSeries{
		Series: s,
		lset:   projectedLabels,
	}
}

// projectLabels applies the projection to a label set, keeping or removing
// labels based on the include/exclude mode. Always preserves __name__ is NOT
// done here — the projection optimizer decides which labels to include.
// The series hash label is appended to preserve original series identity.
func (p *projectedSelector) projectLabels(lset labels.Labels) labels.Labels {
	builder := labels.NewBuilder(labels.EmptyLabels())

	if p.projection.Include {
		// Include mode: only keep the labels in the projection list
		lset.Range(func(l labels.Label) {
			if slices.Contains(p.projection.Labels, l.Name) {
				builder.Set(l.Name, l.Value)
			}
		})
	} else {
		// Exclude mode: keep all labels except those in the projection list
		lset.Range(func(l labels.Label) {
			if !slices.Contains(p.projection.Labels, l.Name) {
				builder.Set(l.Name, l.Value)
			}
		})
	}

	// Append series hash label to preserve original series identity
	if p.seriesHashLabel != "" {
		builder.Set(p.seriesHashLabel, strconv.FormatUint(lset.Hash(), 10))
	}

	return builder.Labels()
}

// projectedSeries wraps a storage.Series but returns projected labels.
type projectedSeries struct {
	storage.Series
	lset labels.Labels
}

func (s *projectedSeries) Labels() labels.Labels {
	return s.lset
}

func (s *projectedSeries) Iterator(it chunkenc.Iterator) chunkenc.Iterator {
	return s.Series.Iterator(it)
}
