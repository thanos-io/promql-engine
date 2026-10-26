// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus

import (
	"context"
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
