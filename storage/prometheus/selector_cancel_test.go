// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/query"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	"github.com/prometheus/prometheus/tsdb/chunks"
	"github.com/prometheus/prometheus/util/annotations"
	"github.com/stretchr/testify/require"
)

const numCancelTestSeries = 10 * ctxCheckInterval

func TestSelectorsStopReadingOnCancel(t *testing.T) {
	opts := &query.Options{
		Start:         time.UnixMilli(0),
		End:           time.UnixMilli(90),
		Step:          10 * time.Millisecond,
		LookbackDelta: 5 * time.Minute,
		StepsBatch:    10,
		SampleTracker: query.NewSampleTracker(0),
	}
	newVectorSelector := func(s SeriesSelector) (model.VectorOperator, error) {
		return NewVectorSelector(s, opts, 0, 0, false, 0, 1), nil
	}
	newMatrixSelector := func(s SeriesSelector) (model.VectorOperator, error) {
		return NewMatrixSelector(s, "rate", 0, 0, opts, 50*time.Millisecond, 0, 0, 0, 1)
	}

	for _, tcase := range []struct {
		name  string
		newOp func(SeriesSelector) (model.VectorOperator, error)
		// cancelDuringLoad cancels the query while series are being loaded
		// instead of while samples are read in Next.
		cancelDuringLoad bool
		// cancelOnIterator also cancels the query when an iterator is created,
		// not only when samples are read.
		cancelOnIterator bool
	}{
		{name: "vector selector next", newOp: newVectorSelector},
		{name: "vector selector load", newOp: newVectorSelector, cancelDuringLoad: true},
		{name: "matrix selector next", newOp: newMatrixSelector},
		{name: "matrix selector iterator creation", newOp: newMatrixSelector, cancelOnIterator: true},
	} {
		t.Run(tcase.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Once armed, cancel the query as soon as any series is used.
			var armed bool
			seriesUsed := make(map[int]struct{})
			series := make([]SignedSeries, numCancelTestSeries)
			for i := range series {
				onUse := func() {
					if armed {
						seriesUsed[i] = struct{}{}
						cancel()
					}
				}
				s := &readTrackingSeries{
					Series: storage.NewListSeries(labels.FromStrings("i", strconv.Itoa(i)), chunks.GenerateSamples(0, 100)),
					onRead: onUse,
				}
				if tcase.cancelOnIterator {
					s.onIterator = onUse
				}
				series[i] = SignedSeries{Series: s, Signature: uint64(i)}
			}

			op, err := tcase.newOp(staticSeriesSelector(series))
			require.NoError(t, err)

			if tcase.cancelDuringLoad {
				armed = true
				_, err = op.Series(ctx)
			} else {
				_, err = op.Series(ctx)
				require.NoError(t, err)
				armed = true
				_, err = op.Next(ctx, make([]model.StepVector, opts.StepsBatch))
			}
			require.ErrorIs(t, err, context.Canceled)
			require.LessOrEqual(t, len(seriesUsed), ctxCheckInterval)
		})
	}
}

func TestSeriesSelectorStopsLoadingOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	set := &cancelingSeriesSet{cancel: cancel}
	querier := &storage.MockQuerier{
		SelectMockFunction: func(bool, *storage.SelectHints, ...*labels.Matcher) storage.SeriesSet { return set },
	}

	_, err := newSeriesSelector(querier, nil, storage.SelectHints{}).GetSeries(ctx, 0, 1)
	require.ErrorIs(t, err, context.Canceled)
	require.LessOrEqual(t, set.n, ctxCheckInterval+1)
}

type staticSeriesSelector []SignedSeries

func (s staticSeriesSelector) GetSeries(context.Context, int, int) ([]SignedSeries, error) {
	return s, nil
}

func (staticSeriesSelector) Matchers() []*labels.Matcher { return nil }

type readTrackingSeries struct {
	storage.Series
	onRead func()
	// onIterator is called when an iterator is created, if set.
	onIterator func()
}

func (s *readTrackingSeries) Iterator(it chunkenc.Iterator) chunkenc.Iterator {
	if s.onIterator != nil {
		s.onIterator()
	}
	return &readTrackingIterator{Iterator: s.Series.Iterator(it), onRead: s.onRead}
}

type readTrackingIterator struct {
	chunkenc.Iterator
	onRead func()
}

func (it *readTrackingIterator) Next() chunkenc.ValueType {
	it.onRead()
	return it.Iterator.Next()
}

func (it *readTrackingIterator) Seek(t int64) chunkenc.ValueType {
	it.onRead()
	return it.Iterator.Seek(t)
}

// cancelingSeriesSet cancels the query when the first series is loaded.
type cancelingSeriesSet struct {
	cancel context.CancelFunc
	n      int
}

func (s *cancelingSeriesSet) Next() bool {
	s.n++
	if s.n == 1 {
		s.cancel()
	}
	return s.n <= numCancelTestSeries
}

func (s *cancelingSeriesSet) At() storage.Series {
	return storage.NewListSeries(labels.FromStrings("i", strconv.Itoa(s.n)), nil)
}

func (*cancelingSeriesSet) Err() error                        { return nil }
func (*cancelingSeriesSet) Warnings() annotations.Annotations { return nil }
