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
	"github.com/thanos-io/promql-engine/ringbuffer"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/histogram"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
)

// rangeFunctions covers one function per ring buffer implementation.
var rangeFunctions = []string{
	"sum_over_time", // OverTimeBuffer
	"rate",          // RateBuffer
	"changes",       // GenericRingBuffer
}

// newTestMatrixSelector builds the operator directly so tests can inspect
// per-series scanner state.
func newTestMatrixSelector(t *testing.T, load, functionName string, selectRange, maxt time.Duration, stepsBatch int) (*matrixSelector, func()) {
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
	hints := storage.SelectHints{
		Start: -selectRange.Milliseconds(),
		End:   maxt.Milliseconds(),
		Range: selectRange.Milliseconds(),
	}
	querier, err := testStorage.Querier(hints.Start, hints.End)
	testutil.Ok(t, err)

	call, err := ringbuffer.NewRangeVectorFunc(functionName)
	testutil.Ok(t, err)

	o := &matrixSelector{
		storage:      newSeriesSelector(querier, matchers, hints),
		call:         call,
		functionName: functionName,
		fhReader:     &histogram.FloatHistogram{},

		opts:     opts,
		numSteps: opts.NumStepsPerBatch(),
		mint:     opts.Start.UnixMilli(),
		maxt:     opts.End.UnixMilli(),
		step:     opts.Step.Milliseconds(),

		selectRange: selectRange.Milliseconds(),
		currentStep: opts.Start.UnixMilli(),

		numShards: 1,
	}
	o.telemetry = telemetry.NewTelemetry(o, opts)

	return o, func() {
		testutil.Ok(t, querier.Close())
		testutil.Ok(t, testStorage.Close())
	}
}

// TestMatrixSelectorReleasesIteratorAndBuffer covers when a scanner's iterator
// and ring buffer are dropped, once per ring buffer implementation.
func TestMatrixSelectorReleasesIteratorAndBuffer(t *testing.T) {
	for _, fn := range rangeFunctions {
		t.Run(fn, func(t *testing.T) {
			t.Run("releases an exhausted series", testMatrixSelectorReleasesExhausted(fn))
			t.Run("keeps a series with a sample still to come", testMatrixSelectorKeepsGapped(fn))
		})
	}
}

func testMatrixSelectorReleasesExhausted(fn string) func(*testing.T) {
	// Samples at 0s..90s only; with a 2m range the last leaves at ts=210s.
	load := `load 30s
		metric 1+1x3`

	return func(t *testing.T) {
		o, cleanup := newTestMatrixSelector(t, load, fn, 2*time.Minute, 600*time.Second, 5)
		defer cleanup()

		ctx := context.Background()
		_, err := o.Series(ctx)
		testutil.Ok(t, err)

		series := &o.scanners[0]
		buf := make([]model.StepVector, 5)

		// Batch 1, ts 0s..120s: drains partway, so only the iterator goes.
		_, err = o.Next(ctx, buf)
		testutil.Ok(t, err)
		testutil.Assert(t, !series.exhausted, "must not be released while the window still covers its samples")
		testutil.Assert(t, series.buffer != nil, "buffer must stay live while the window still covers its samples")
		testutil.Assert(t, series.iterator == nil, "iterator must be dropped as soon as it is drained")

		// Batch 2, ts 150s..270s: the window moves past 90s.
		n, err := o.Next(ctx, buf)
		testutil.Ok(t, err)
		testutil.Assert(t, series.exhausted, "must be released once the window moves past the last sample")
		testutil.Assert(t, series.buffer == nil, "buffer must be released")
		testutil.Assert(t, series.iterator == nil, "iterator must be released")

		// mint is exclusive, so 90s drops out at ts=210s.
		lastValueAt := sampleTimestamps(buf[:n], series.signature)
		testutil.Assert(t, len(lastValueAt) > 0, "must still produce values in the steps before release, got none")
		for _, ts := range lastValueAt {
			testutil.Assert(t, ts <= 180000, "produced a value at %d, past its last window", ts)
		}

		for range 3 {
			n, err = o.Next(ctx, buf)
			testutil.Ok(t, err)
			testutil.Assert(t, series.buffer == nil, "released buffer must not be recreated")
			testutil.Assert(t, series.iterator == nil, "released iterator must not be recreated")
			testutil.Equals(t, 0, len(sampleTimestamps(buf[:n], series.signature)))
		}
	}
}

// A gap empties the buffer while the iterator still holds a sample beyond maxt.
func testMatrixSelectorKeepsGapped(fn string) func(*testing.T) {
	// Samples at 0s and 450s onwards; ts=150s..420s yield nothing.
	load := `load 30s
		metric 1 _ _ _ _ _ _ _ _ _ _ _ _ _ _ 5+1x5`

	return func(t *testing.T) {
		o, cleanup := newTestMatrixSelector(t, load, fn, 2*time.Minute, 600*time.Second, 5)
		defer cleanup()

		ctx := context.Background()
		_, err := o.Series(ctx)
		testutil.Ok(t, err)

		series := &o.scanners[0]
		buf := make([]model.StepVector, 5)

		// Across the gap the scanner must survive.
		for range 3 {
			_, err := o.Next(ctx, buf)
			testutil.Ok(t, err)
			testutil.Assert(t, !series.exhausted, "a series with a sample still to come must not be released")
			testutil.Assert(t, series.buffer != nil, "a series with a sample still to come must keep its buffer")
			testutil.Assert(t, series.iterator != nil, "a series with a sample still to come must keep its iterator")
		}

		// The series resumes at 450s.
		n, err := o.Next(ctx, buf)
		testutil.Ok(t, err)
		testutil.Assert(t, len(sampleTimestamps(buf[:n], series.signature)) > 0,
			"must produce values again after the gap")
	}
}
