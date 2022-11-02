// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package storage

import (
	"context"
	"fmt"

	"github.com/prometheus/prometheus/model/labels"
)

type filteredSelector struct {
	selector SeriesSelector
	filter   Filter
}

func NewFilteredSelector(selector SeriesSelector, filter Filter) SeriesSelector {
	return &filteredSelector{
		selector: selector,
		filter:   filter,
	}
}

func (f *filteredSelector) Explain() string {
	return fmt.Sprintf("[*filteredSelector] {%v}: %v", f.filter.Matchers(), f.selector.Explain())
}

func (f *filteredSelector) Matchers() []*labels.Matcher {
	return append(f.selector.Matchers(), f.filter.Matchers()...)
}

func (f *filteredSelector) GetSeries(ctx context.Context, shard, numShards int) ([]SignedSeries, error) {
	series, err := f.selector.GetSeries(ctx, shard, numShards)
	if err != nil {
		return nil, err
	}

	i := uint64(0)
	ss := make([]SignedSeries, 0, len(series))
	for _, s := range series {
		if f.filter.Matches(s) {
			ss = append(ss, SignedSeries{
				Series:    s.Series,
				Signature: i,
			})
			i++
		}
	}
	return ss, nil
}
