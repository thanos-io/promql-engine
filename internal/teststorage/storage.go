// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

// Package teststorage adapts Prometheus TSDB for concurrent engine tests.
// It is a test fixture, not a production storage adapter. Its synchronization
// also affects TSDB-backed benchmarks; those do not measure lock-free decoding.
package teststorage

import (
	"context"
	"sync"
	"testing"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/promqltest"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	"github.com/prometheus/prometheus/util/annotations"
	promteststorage "github.com/prometheus/prometheus/util/teststorage"
)

// TestStorage retains the upstream fixture's loading, appending and cleanup
// behavior while serializing access to each querier's shared reader buffers.
type TestStorage struct {
	*promteststorage.TestStorage
}

func New(t testing.TB, opts ...promteststorage.Option) *TestStorage {
	t.Helper()
	return &TestStorage{TestStorage: promteststorage.New(t, opts...)}
}

func LoadedStorage(t testing.TB, input string) *TestStorage {
	t.Helper()
	return &TestStorage{TestStorage: promqltest.LoadedStorage(t, input)}
}

func (s *TestStorage) Querier(mint, maxt int64) (storage.Querier, error) {
	return wrapQuerier(s.TestStorage.Querier(mint, maxt))
}

// WrapQueryable adapts a directly opened TSDB used by a test or benchmark.
// The caller retains responsibility for closing the underlying storage.
func WrapQueryable(q storage.Queryable) storage.Queryable {
	return &queryable{Queryable: q}
}

type queryable struct {
	storage.Queryable
}

func (q *queryable) Querier(mint, maxt int64) (storage.Querier, error) {
	return wrapQuerier(q.Queryable.Querier(mint, maxt))
}

func wrapQuerier(q storage.Querier, err error) (storage.Querier, error) {
	if err != nil {
		return nil, err
	}
	return &querier{Querier: q}, nil
}

// Prometheus PR #18302 introduced querier-owned index/chunk scratch buffers.
// Different SeriesSets and sample iterators can reach the same buffers, so the
// mutex must belong to the querier, not to an individual selector or iterator.
// https://github.com/prometheus/prometheus/pull/18302
type querier struct {
	storage.Querier
	mtx sync.Mutex
}

func (q *querier) Select(ctx context.Context, sorted bool, hints *storage.SelectHints, matchers ...*labels.Matcher) storage.SeriesSet {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	return &seriesSet{SeriesSet: q.Querier.Select(ctx, sorted, hints, matchers...), mtx: &q.mtx}
}

func (q *querier) LabelValues(ctx context.Context, name string, hints *storage.LabelHints, matchers ...*labels.Matcher) ([]string, annotations.Annotations, error) {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	return q.Querier.LabelValues(ctx, name, hints, matchers...)
}

func (q *querier) LabelNames(ctx context.Context, hints *storage.LabelHints, matchers ...*labels.Matcher) ([]string, annotations.Annotations, error) {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	return q.Querier.LabelNames(ctx, hints, matchers...)
}

func (q *querier) Close() error {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	return q.Querier.Close()
}

type seriesSet struct {
	storage.SeriesSet
	mtx *sync.Mutex
}

func (s *seriesSet) Next() bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.SeriesSet.Next()
}

func (s *seriesSet) At() storage.Series {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return &series{Series: s.SeriesSet.At(), mtx: s.mtx}
}

type series struct {
	storage.Series
	mtx *sync.Mutex
}

func (s *series) Iterator(reuse chunkenc.Iterator) chunkenc.Iterator {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	it, ok := reuse.(*iterator)
	if !ok {
		it = &iterator{Iterator: reuse}
	}
	it.Iterator = s.Series.Iterator(it.Iterator)
	it.mtx = s.mtx
	return it
}

// Next and Seek can load chunks through the shared reader. At/AtHistogram/
// AtFloatHistogram/AtT/AtST/Err only inspect this iterator's state. As in the
// underlying API, callers must not use the same iterator concurrently.
type iterator struct {
	chunkenc.Iterator
	mtx *sync.Mutex
}

func (it *iterator) Next() chunkenc.ValueType {
	it.mtx.Lock()
	defer it.mtx.Unlock()
	return it.Iterator.Next()
}

func (it *iterator) Seek(t int64) chunkenc.ValueType {
	it.mtx.Lock()
	defer it.mtx.Unlock()
	return it.Iterator.Seek(t)
}
