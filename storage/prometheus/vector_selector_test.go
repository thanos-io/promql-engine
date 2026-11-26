// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/execution/telemetry"
	"github.com/thanos-io/promql-engine/internal/teststorage"
	"github.com/thanos-io/promql-engine/query"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
)

const testLookbackDelta = 5 * time.Minute

// sampleTimestamps returns the timestamps at which a series produced a sample.
func sampleTimestamps(vectors []model.StepVector, signature uint64) []int64 {
	var ts []int64
	for _, v := range vectors {
		for _, id := range v.SampleIDs {
			if id == signature {
				ts = append(ts, v.T)
			}
		}
	}
	return ts
}

// newTestVectorSelector builds the operator directly so tests can inspect
// per-series scanner state.
func newTestVectorSelector(t *testing.T, load string, maxt time.Duration, stepsBatch int) (*vectorSelector, func()) {
	t.Helper()

	testStorage := teststorage.LoadedStorage(t, load)

	opts := &query.Options{
		Start:               time.UnixMilli(0),
		End:                 time.UnixMilli(maxt.Milliseconds()),
		Step:                30 * time.Second,
		StepsBatch:          stepsBatch,
		LookbackDelta:       testLookbackDelta,
		DecodingConcurrency: 1,
		SampleTracker:       query.NewSampleTracker(0),
	}

	matchers := []*labels.Matcher{labels.MustNewMatcher(labels.MatchEqual, "__name__", "metric")}
	hints := storage.SelectHints{Start: -testLookbackDelta.Milliseconds(), End: maxt.Milliseconds()}
	querier, err := testStorage.Querier(hints.Start, hints.End)
	testutil.Ok(t, err)

	o := &vectorSelector{
		storage:         newSeriesSelector(querier, matchers, hints),
		mint:            opts.Start.UnixMilli(),
		maxt:            opts.End.UnixMilli(),
		step:            opts.Step.Milliseconds(),
		currentStep:     opts.Start.UnixMilli(),
		lookbackDelta:   opts.LookbackDelta.Milliseconds(),
		numSteps:        opts.NumStepsPerBatch(),
		seriesBatchSize: 0,
		numShards:       1,
		opts:            opts,
	}
	o.telemetry = telemetry.NewTelemetry(o, opts)

	return o, func() {
		testutil.Ok(t, querier.Close())
		testutil.Ok(t, testStorage.Close())
	}
}

// TestVectorSelectorReleasesIterator covers when a series iterator is dropped,
// including a series that only looks finished.
func TestVectorSelectorReleasesIterator(t *testing.T) {
	t.Run("releases an exhausted series", testVectorSelectorReleasesExhausted)
	t.Run("keeps a gapped series", testVectorSelectorKeepsGapped)
	t.Run("selectPoint reports exhaustion", testSelectPointExhaustion)
}

func testVectorSelectorReleasesExhausted(t *testing.T) {
	// Samples at 0s..90s only; with a 5m lookback it stops at ts=390s.
	load := `load 30s
		metric 1+1x3`

	o, cleanup := newTestVectorSelector(t, load, 600*time.Second, 10)
	defer cleanup()

	ctx := context.Background()
	_, err := o.Series(ctx)
	testutil.Ok(t, err)

	series := &o.scanners[0]

	buf := make([]model.StepVector, 10)

	// Batch 1, ts 0s..270s: still inside the lookback window.
	n, err := o.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 10, n)
	testutil.Assert(t, !series.exhausted, "must not be released while inside the lookback window")
	testutil.Assert(t, series.samples != nil, "iterator must stay live inside the lookback window")
	testutil.Equals(t,
		[]int64{0, 30000, 60000, 90000, 120000, 150000, 180000, 210000, 240000, 270000},
		sampleTimestamps(buf[:n], series.signature))

	// Batch 2, ts 300s..570s: the window is half-open, so the last point is
	// ts=360s and the release lands at ts=390s.
	n, err = o.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 10, n)
	testutil.Assert(t, series.exhausted, "must be released once past lastSample+lookbackDelta")
	testutil.Assert(t, series.samples == nil, "exhausted iterator must be released")
	testutil.Equals(t, []int64{300000, 330000, 360000}, sampleTimestamps(buf[:n], series.signature))

	// Batch 3, ts 600s: no rebuild.
	n, err = o.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 1, n)
	testutil.Assert(t, series.samples == nil, "released iterator must not be recreated")
	testutil.Equals(t, 0, len(sampleTimestamps(buf[:n], series.signature)))
}

// A gap wider than the lookback delta also makes selectPoint return ok=false,
// but the series is not finished.
func testVectorSelectorKeepsGapped(t *testing.T) {
	// Samples at 0s and 450s, nothing in between.
	load := `load 30s
		metric 1 _ _ _ _ _ _ _ _ _ _ _ _ _ _ 5+1x5`

	o, cleanup := newTestVectorSelector(t, load, 600*time.Second, 10)
	defer cleanup()

	ctx := context.Background()
	_, err := o.Series(ctx)
	testutil.Ok(t, err)

	series := &o.scanners[0]
	buf := make([]model.StepVector, 10)

	// Batch 1, ts 0s..270s: the first sample, then the gap.
	n, err := o.Next(ctx, buf)
	testutil.Ok(t, err)
	gotFirst := sampleTimestamps(buf[:n], series.signature)
	testutil.Equals(t, []int64{0, 30000, 60000, 90000, 120000, 150000, 180000, 210000, 240000, 270000}, gotFirst)
	testutil.Assert(t, !series.exhausted, "a gapped series must not be marked exhausted")
	testutil.Assert(t, series.samples != nil, "a gapped series must keep its iterator")

	// Batch 2, ts 300s..570s: nothing until the series resumes at 450s.
	n, err = o.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t,
		[]int64{450000, 480000, 510000, 540000, 570000},
		sampleTimestamps(buf[:n], series.signature))
	testutil.Assert(t, !series.exhausted, "a resumed series must not be marked exhausted")
}

// Only one of the three reasons for ok=false is terminal.
func testSelectPointExhaustion(t *testing.T) {
	lookback := testLookbackDelta.Milliseconds()

	t.Run("empty series", func(t *testing.T) {
		it := storage.NewMemoizedEmptyIterator(lookback)
		_, _, _, ok, exhausted, err := selectPoint(it, 0, lookback, 0)
		testutil.Ok(t, err)
		testutil.Equals(t, false, ok)
		testutil.Equals(t, true, exhausted)
	})

	t.Run("ended series", func(t *testing.T) {
		it := memoizedIteratorFor(t, `load 30s
			metric 1+1x3`, lookback)

		// Inside the lookback window the last sample still yields.
		_, _, _, ok, exhausted, err := selectPoint(it, 300_000, lookback, 0)
		testutil.Ok(t, err)
		testutil.Equals(t, true, ok)
		testutil.Equals(t, false, exhausted)

		// From refTime=390s the last sample leaves the window.
		_, _, _, ok, exhausted, err = selectPoint(it, 390_000, lookback, 0)
		testutil.Ok(t, err)
		testutil.Equals(t, false, ok)
		testutil.Equals(t, true, exhausted)

		_, _, _, ok, exhausted, err = selectPoint(it, 420_000, lookback, 0)
		testutil.Ok(t, err)
		testutil.Equals(t, false, ok)
		testutil.Equals(t, true, exhausted)
	})

	t.Run("gapped series", func(t *testing.T) {
		it := memoizedIteratorFor(t, `load 30s
			metric 1 _ _ _ _ _ _ _ _ _ _ _ _ _ _ 5`, lookback)

		// Inside the gap: no sample, but not finished.
		_, _, _, ok, exhausted, err := selectPoint(it, 360_000, lookback, 0)
		testutil.Ok(t, err)
		testutil.Equals(t, false, ok)
		testutil.Equals(t, false, exhausted)

		// It resumes at 450s.
		_, _, _, ok, exhausted, err = selectPoint(it, 450_000, lookback, 0)
		testutil.Ok(t, err)
		testutil.Equals(t, true, ok)
		testutil.Equals(t, false, exhausted)
	})
}

func memoizedIteratorFor(t *testing.T, load string, lookbackDelta int64) *storage.MemoizedSeriesIterator {
	t.Helper()

	testStorage := teststorage.LoadedStorage(t, load)
	t.Cleanup(func() { testutil.Ok(t, testStorage.Close()) })

	querier, err := testStorage.Querier(0, 600_000)
	testutil.Ok(t, err)
	t.Cleanup(func() { testutil.Ok(t, querier.Close()) })

	set := querier.Select(context.Background(), false, nil,
		labels.MustNewMatcher(labels.MatchEqual, "__name__", "metric"))
	testutil.Assert(t, set.Next(), "expected one series")
	return storage.NewMemoizedIterator(set.At().Iterator(nil), lookbackDelta)
}

func smoothedIteratorFor(t *testing.T, load string) *smoothedIterator {
	t.Helper()

	testStorage := teststorage.LoadedStorage(t, load)
	t.Cleanup(func() { testutil.Ok(t, testStorage.Close()) })

	querier, err := testStorage.Querier(0, 3_600_000)
	testutil.Ok(t, err)
	t.Cleanup(func() { testutil.Ok(t, querier.Close()) })

	set := querier.Select(context.Background(), false, nil,
		labels.MustNewMatcher(labels.MatchEqual, "__name__", "metric"))
	testutil.Assert(t, set.Next(), "expected one series")
	return newSmoothedIterator(set.At().Iterator(nil))
}

// TestSmoothedIteratorBuffered checks that the incremental count of buffered
// samples matches what the iterator holds, and that vacated buffer slots do
// not keep histograms alive.
func TestSmoothedIteratorBuffered(t *testing.T) {
	for _, load := range []string{
		`load 10s
			metric 0+1x60 _x300 5`,
		`load 10s
			metric {{schema:0 sum:5 count:4 buckets:[1 2 1]}}x60 _x300 {{schema:0 sum:5 count:4 buckets:[1 2 1]}}`,
		`load 10s
			metric 1 2 stale 4 {{schema:0 sum:5 count:4 buckets:[1 2 1]}} 6 _x20 7`,
	} {
		it := smoothedIteratorFor(t, load)
		lookback := testLookbackDelta.Milliseconds()
		for ts := int64(0); ts <= 3_600_000; ts += 15_000 {
			_, _, _, _, _, err := it.selectPoint(ts, lookback, 0)
			testutil.Ok(t, err)

			want := 0
			if it.prevOK {
				want += smoothedPointSamples(it.prev)
			}
			for _, p := range it.ahead {
				want += smoothedPointSamples(p)
			}
			if it.hasPending {
				want += smoothedPointSamples(it.pending)
			}
			testutil.Equals(t, want, it.buffered, "at %d", ts)
			for _, p := range it.ahead[len(it.ahead):cap(it.ahead)] {
				testutil.Assert(t, p.h == nil, "vacated slot at %d keeps a histogram", ts)
			}
		}
	}
}

// countingTracker is a SampleTracker without a limit that exposes its balance.
type countingTracker struct{ current int }

func (c *countingTracker) Add(n int)         { c.current += n }
func (c *countingTracker) Remove(n int)      { c.current -= n }
func (c *countingTracker) CheckLimit() error { return nil }
func (c *countingTracker) Limit() int64      { return math.MaxInt64 }

// TestSmoothedVectorSelectorReleasesIterator is TestVectorSelectorReleasesIterator
// for smoothed range queries: a finished series drops its iterator and its
// charge to the sample limit, a gapped series keeps both. The tracker holds
// the buffered samples plus the outputs of the last batch.
func TestSmoothedVectorSelectorReleasesIterator(t *testing.T) {
	t.Run("releases an exhausted series", func(t *testing.T) {
		// Samples at 0s..90s only. A smoothed point needs a sample after
		// ts-5m, so the last point is ts=360s and the release lands at ts=390s.
		o, cleanup := newTestVectorSelector(t, `load 30s
			metric 1+1x3`, 600*time.Second, 10)
		defer cleanup()
		o.smoothed = true
		tracker := &countingTracker{}
		o.opts.SampleTracker = tracker

		ctx := context.Background()
		_, err := o.Series(ctx)
		testutil.Ok(t, err)
		series := &o.scanners[0]
		buf := make([]model.StepVector, 10)

		// Batch 1, ts 0s..270s.
		n, err := o.Next(ctx, buf)
		testutil.Ok(t, err)
		testutil.Equals(t, 10, n)
		testutil.Equals(t,
			[]int64{0, 30000, 60000, 90000, 120000, 150000, 180000, 210000, 240000, 270000},
			sampleTimestamps(buf[:n], series.signature))
		testutil.Assert(t, !series.exhausted, "must not be released while inside the lookback window")
		testutil.Assert(t, series.smoothed != nil, "iterator must stay live inside the lookback window")
		testutil.Assert(t, series.chargedBuffered > 0, "the buffered samples must be charged")
		testutil.Equals(t, series.smoothed.buffered, series.chargedBuffered)
		testutil.Equals(t, 11, tracker.current, "10 outputs and 1 buffered sample")

		// Batch 2, ts 300s..570s.
		n, err = o.Next(ctx, buf)
		testutil.Ok(t, err)
		testutil.Equals(t, 10, n)
		testutil.Equals(t, []int64{300000, 330000, 360000}, sampleTimestamps(buf[:n], series.signature))
		testutil.Assert(t, series.exhausted, "must be released once past lastSample+lookbackDelta")
		testutil.Assert(t, series.smoothed == nil, "exhausted iterator must be released")
		testutil.Equals(t, 0, series.chargedBuffered)
		testutil.Equals(t, 3, tracker.current, "3 outputs and no buffered sample")

		// Batch 3, ts 600s: no rebuild.
		n, err = o.Next(ctx, buf)
		testutil.Ok(t, err)
		testutil.Equals(t, 1, n)
		testutil.Assert(t, series.smoothed == nil, "released iterator must not be recreated")
		testutil.Equals(t, 0, series.chargedBuffered)
		testutil.Equals(t, 0, len(sampleTimestamps(buf[:n], series.signature)))
		testutil.Equals(t, 0, tracker.current)
	})

	t.Run("keeps a gapped series", func(t *testing.T) {
		// Samples at 0s and 900s..1050s. The gap is wider than twice the
		// lookback, so some steps see no sample before or after them.
		o, cleanup := newTestVectorSelector(t, `load 30s
			metric 1 _x29 5+1x5`, 1200*time.Second, 10)
		defer cleanup()
		o.smoothed = true
		tracker := &countingTracker{}
		o.opts.SampleTracker = tracker

		ctx := context.Background()
		_, err := o.Series(ctx)
		testutil.Ok(t, err)
		series := &o.scanners[0]
		buf := make([]model.StepVector, 10)

		for _, want := range [][]int64{
			// ts 0s..270s: the first sample carried forward.
			{0, 30000, 60000, 90000, 120000, 150000, 180000, 210000, 240000, 270000},
			// ts 300s..570s: no sample within the lookback on either side.
			nil,
			// ts 600s..870s: only samples ahead, nothing to interpolate from.
			nil,
			// ts 900s..1170s: the series resumes.
			{900000, 930000, 960000, 990000, 1020000, 1050000, 1080000, 1110000, 1140000, 1170000},
		} {
			n, err := o.Next(ctx, buf)
			testutil.Ok(t, err)
			testutil.Equals(t, want, sampleTimestamps(buf[:n], series.signature))
			testutil.Assert(t, !series.exhausted, "a gapped series must not be marked exhausted")
			testutil.Assert(t, series.smoothed != nil, "a gapped series must keep its iterator")
			testutil.Equals(t, series.smoothed.buffered, series.chargedBuffered)
			testutil.Equals(t, series.smoothed.buffered+len(want), tracker.current)
		}
	})
}

// TestSmoothedInstantQueryReleasesBuffer checks that the smoothed iterator an
// instant query reuses across series holds no samples once a series is done, as
// its charge to the sample limit is dropped.
func TestSmoothedInstantQueryReleasesBuffer(t *testing.T) {
	o, cleanup := newTestVectorSelector(t, `load 1s
		metric{id="a"} 0+1x600
		metric{id="b"} 0+1x600`, 0, 10)
	defer cleanup()
	o.smoothed = true

	ctx := context.Background()
	_, err := o.Series(ctx)
	testutil.Ok(t, err)

	buf := make([]model.StepVector, 10)
	n, err := o.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 1, n)
	testutil.Equals(t, 2, len(buf[0].Samples))

	testutil.Equals(t, 0, o.reusableSmoothed.buffered)
	testutil.Equals(t, 0, len(o.reusableSmoothed.ahead))
	testutil.Assert(t, !o.reusableSmoothed.prevOK && !o.reusableSmoothed.hasPending, "released iterator keeps samples")
	for _, s := range o.scanners {
		testutil.Equals(t, 0, s.chargedBuffered)
	}
}
