// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus

import (
	"strconv"
	"strings"

	"github.com/cespare/xxhash/v2"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
)

var sep = []byte{'\xff'}

// SelectorHashFunc computes a cache key for label matchers and select hints.
type SelectorHashFunc func(matchers []*labels.Matcher, mint, maxt int64, hints storage.SelectHints) uint64

type SelectorPool struct {
	selectors map[uint64]*seriesSelector
	querier   storage.Querier
	hashFunc  SelectorHashFunc
}

func NewSelectorPool(querier storage.Querier, hashFunc SelectorHashFunc) *SelectorPool {
	if hashFunc == nil {
		hashFunc = DefaultSelectorHash
	}
	return &SelectorPool{
		selectors: make(map[uint64]*seriesSelector),
		querier:   querier,
		hashFunc:  hashFunc,
	}
}

func (p *SelectorPool) GetSelector(mint, maxt, step int64, matchers []*labels.Matcher, hints storage.SelectHints) SeriesSelector {
	key := p.hashFunc(matchers, mint, maxt, hints)
	if existing, ok := p.selectors[key]; ok {
		existing.hints.Start = min(existing.hints.Start, mint)
		return existing
	}
	p.selectors[key] = newSeriesSelector(p.querier, matchers, hints)
	return p.selectors[key]
}

func (p *SelectorPool) GetFilteredSelector(mint, maxt, step int64, matchers, filters []*labels.Matcher, hints storage.SelectHints) SeriesSelector {
	key := p.hashFunc(matchers, mint, maxt, hints)
	if existing, ok := p.selectors[key]; ok {
		existing.hints.Start = min(existing.hints.Start, mint)
		return NewFilteredSelector(existing, NewFilter(filters))
	}
	p.selectors[key] = newSeriesSelector(p.querier, matchers, hints)
	return NewFilteredSelector(p.selectors[key], NewFilter(filters))
}

// DefaultSelectorHash excludes mint so different range windows within a query share one Select().
func DefaultSelectorHash(matchers []*labels.Matcher, _, maxt int64, hints storage.SelectHints) uint64 {
	sb := xxhash.New()
	for _, m := range matchers {
		writeMatcher(sb, m)
	}
	writeInt64(sb, maxt)
	writeInt64(sb, hints.Step)
	writeString(sb, hints.Func)
	writeString(sb, strings.Join(hints.Grouping, ";"))
	writeBool(sb, hints.By)
	writeString(sb, strings.Join(hints.ProjectionLabels, ";"))
	writeBool(sb, hints.ProjectionInclude)

	key := sb.Sum64()
	return key
}

func writeMatcher(sb *xxhash.Digest, m *labels.Matcher) {
	writeString(sb, m.Name)
	writeString(sb, strconv.Itoa(int(m.Type)))
	writeString(sb, m.Value)
}

func writeInt64(sb *xxhash.Digest, val int64) {
	_, _ = sb.WriteString(strconv.FormatInt(val, 10))
	_, _ = sb.Write(sep)
}

func writeString(sb *xxhash.Digest, val string) {
	_, _ = sb.WriteString(val)
	_, _ = sb.Write(sep)
}

func writeBool(sb *xxhash.Digest, val bool) {
	_, _ = sb.WriteString(strconv.FormatBool(val))
	_, _ = sb.Write(sep)
}
