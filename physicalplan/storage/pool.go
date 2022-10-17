// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package storage

import (
	"fmt"
	"strconv"

	"github.com/cespare/xxhash/v2"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
)

var sep = []byte{'\xff'}

type SelectorPool struct {
	selectors map[uint64]*seriesSelector

	queryable storage.Queryable
}

func NewSelectorPool(queryable storage.Queryable) *SelectorPool {
	return &SelectorPool{
		selectors: make(map[uint64]*seriesSelector),
		queryable: queryable,
	}
}

func (p *SelectorPool) GetSelector(mint, maxt int64, matchers []*labels.Matcher) SeriesSelector {
	key := hashMatchers(matchers, mint, maxt)
	if _, ok := p.selectors[key]; !ok {
		p.selectors[key] = newSeriesSelector(p.queryable, mint, maxt, matchers)
	}
	return p.selectors[key]
}

func (p *SelectorPool) GetFilteredSelector(mint, maxt int64, matchers, filters []*labels.Matcher) SeriesSelector {
	key := hashMatchers(matchers, mint, maxt)
	if _, ok := p.selectors[key]; !ok {
		p.selectors[key] = newSeriesSelector(p.queryable, mint, maxt, matchers)
	}

	return NewFilteredSelector(p.selectors[key], NewFilter(filters))
}

func hashMatchers(matchers []*labels.Matcher, mint, maxt int64) uint64 {
	sb := xxhash.New()
	for _, m := range matchers {
		_, _ = sb.WriteString(m.Name)
		_, _ = sb.Write(sep)
		_, _ = sb.WriteString(strconv.Itoa(int(m.Type)))
		_, _ = sb.Write(sep)
		_, _ = sb.WriteString(m.Value)
		_, _ = sb.Write(sep)
	}
	_, _ = sb.WriteString(fmt.Sprintf("%d", mint))
	_, _ = sb.Write(sep)
	_, _ = sb.WriteString(fmt.Sprintf("%d", maxt))

	key := sb.Sum64()
	return key
}
