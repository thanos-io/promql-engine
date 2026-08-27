// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus

import (
	"context"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/query"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/promqltest"
	"github.com/prometheus/prometheus/storage"
)

// TestSeriesBatchMatrixSelectorCompletesAllStepsPerBatch verifies that with
// SelectorBatchSize configured, the matrixSelector exhausts all steps for a
// batch of series before moving to the next batch, and releases iterators
// for completed batches.
func TestSeriesBatchMatrixSelectorCompletesAllStepsPerBatch(t *testing.T) {
	// 4 series, 14 steps (step=30s, start=60s, end=450s).
	// batchSize=2 → batch1={series0, series1}, batch2={series2, series3}
	// stepsBatch=6 → each batch needs 3 Next() calls: 6 + 6 + 2 steps.
	//
	// Expected with series-first:
	//   Next() call 1: batch1 (series 0,1) for steps 1-6
	//   Next() call 2: batch1 (series 0,1) for steps 7-12
	//   Next() call 3: batch1 (series 0,1) for steps 13-14 (partial)
	//                   → batch1 exhausted all steps, release iterators
	//   Next() call 4: batch2 (series 2,3) for steps 1-6
	//   Next() call 5: batch2 (series 2,3) for steps 7-12
	//   Next() call 6: batch2 (series 2,3) for steps 13-14 (partial)
	//                   → batch2 exhausted all steps, release iterators
	//   Next() call 7: returns 0 (all batches exhausted)

	load := `load 30s
		metric{s="0"} 1+1x30
		metric{s="1"} 2+1x30
		metric{s="2"} 3+1x30
		metric{s="3"} 4+1x30`

	testStorage := promqltest.LoadedStorage(t, load)
	defer testStorage.Close()

	ctx := context.Background()
	mint := int64(60000)  // 60s in ms
	maxt := int64(450000) // 450s in ms → 14 steps at 30s from 60s

	opts := &query.Options{
		Start:               time.UnixMilli(mint),
		End:                 time.UnixMilli(maxt),
		Step:                30 * time.Second,
		StepsBatch:          6,
		LookbackDelta:       5 * time.Minute,
		DecodingConcurrency: 1,
		SampleTracker:       query.NewSampleTracker(0),
	}

	matchers := []*labels.Matcher{labels.MustNewMatcher(labels.MatchEqual, "__name__", "metric")}
	hints := storage.SelectHints{Start: mint - (5 * time.Minute).Milliseconds(), End: maxt}
	querier, err := testStorage.Querier(hints.Start, hints.End)
	testutil.Ok(t, err)
	defer querier.Close()

	selector := newSeriesSelector(querier, matchers, hints)

	// Create matrixSelector with batchSize=2.
	op, err := NewMatrixSelector(
		selector,
		"last_over_time",
		0, // scalarArg
		0, // scalarArg2
		opts,
		time.Minute, // selectRange
		0,           // offset
		2,           // batchSize = 2
		0,           // shard
		1,           // numShards
	)
	testutil.Ok(t, err)

	_, err = op.Series(ctx)
	testutil.Ok(t, err)

	buf := make([]model.StepVector, 6)

	// Call 1: batch1 (series 0,1), steps 1-6.
	n1, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 6, n1)
	seriesInCall1 := collectSeriesIDs(buf[:n1])
	testutil.Assert(t, containsAll(seriesInCall1, []uint64{0, 1}),
		"call 1 should contain series 0,1, got %v", seriesInCall1)
	testutil.Assert(t, !containsAny(seriesInCall1, []uint64{2, 3}),
		"call 1 should NOT contain series 2,3, got %v", seriesInCall1)

	// Call 2: batch1 (series 0,1), steps 7-12 (same series, next steps).
	n2, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 6, n2)
	seriesInCall2 := collectSeriesIDs(buf[:n2])
	testutil.Assert(t, containsAll(seriesInCall2, []uint64{0, 1}),
		"call 2 should contain series 0,1 (same batch), got %v", seriesInCall2)
	testutil.Assert(t, !containsAny(seriesInCall2, []uint64{2, 3}),
		"call 2 should NOT contain series 2,3, got %v", seriesInCall2)

	// Call 3: batch1 (series 0,1), steps 13-14 (partial, same batch).
	n3, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 2, n3) // only 2 steps remain
	seriesInCall3 := collectSeriesIDs(buf[:n3])
	testutil.Assert(t, containsAll(seriesInCall3, []uint64{0, 1}),
		"call 3 should contain series 0,1 (same batch, partial), got %v", seriesInCall3)
	testutil.Assert(t, !containsAny(seriesInCall3, []uint64{2, 3}),
		"call 3 should NOT contain series 2,3, got %v", seriesInCall3)

	// Call 4: batch2 (series 2,3), steps 1-6.
	n4, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 6, n4)
	seriesInCall4 := collectSeriesIDs(buf[:n4])
	testutil.Assert(t, containsAll(seriesInCall4, []uint64{2, 3}),
		"call 4 should contain series 2,3 (next batch), got %v", seriesInCall4)
	testutil.Assert(t, !containsAny(seriesInCall4, []uint64{0, 1}),
		"call 4 should NOT contain series 0,1, got %v", seriesInCall4)

	// Call 5: batch2 (series 2,3), steps 7-12.
	n5, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 6, n5)
	seriesInCall5 := collectSeriesIDs(buf[:n5])
	testutil.Assert(t, containsAll(seriesInCall5, []uint64{2, 3}),
		"call 5 should contain series 2,3 (same batch), got %v", seriesInCall5)

	// Call 6: batch2 (series 2,3), steps 13-14 (partial).
	n6, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 2, n6) // only 2 steps remain
	seriesInCall6 := collectSeriesIDs(buf[:n6])
	testutil.Assert(t, containsAll(seriesInCall6, []uint64{2, 3}),
		"call 6 should contain series 2,3 (same batch, partial), got %v", seriesInCall6)

	// Call 7: all batches exhausted.
	n7, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 0, n7)
}

// TestMatrixSelectorWithoutBatchingOriginalBehavior verifies that without
// SelectorBatchSize (batchSize=0), the matrixSelector uses the original
// ordering: all series for one step-batch, then advance steps.
func TestMatrixSelectorWithoutBatchingOriginalBehavior(t *testing.T) {
	// 4 series, 14 steps, stepsBatch=6, batchSize=0 (no batching).
	//
	// Expected (original ordering):
	//   Next() call 1: ALL series (0,1,2,3) for steps 1-6
	//   Next() call 2: ALL series (0,1,2,3) for steps 7-12
	//   Next() call 3: ALL series (0,1,2,3) for steps 13-14 (partial)
	//   Next() call 4: returns 0

	load := `load 30s
		metric{s="0"} 1+1x30
		metric{s="1"} 2+1x30
		metric{s="2"} 3+1x30
		metric{s="3"} 4+1x30`

	testStorage := promqltest.LoadedStorage(t, load)
	defer testStorage.Close()

	ctx := context.Background()
	mint := int64(60000)
	maxt := int64(450000)

	opts := &query.Options{
		Start:               time.UnixMilli(mint),
		End:                 time.UnixMilli(maxt),
		Step:                30 * time.Second,
		StepsBatch:          6,
		LookbackDelta:       5 * time.Minute,
		DecodingConcurrency: 1,
		SampleTracker:       query.NewSampleTracker(0),
	}

	matchers := []*labels.Matcher{labels.MustNewMatcher(labels.MatchEqual, "__name__", "metric")}
	hints := storage.SelectHints{Start: mint - (5 * time.Minute).Milliseconds(), End: maxt}
	querier, err := testStorage.Querier(hints.Start, hints.End)
	testutil.Ok(t, err)
	defer querier.Close()

	selector := newSeriesSelector(querier, matchers, hints)

	// Create matrixSelector with batchSize=0 (no batching).
	op, err := NewMatrixSelector(
		selector,
		"last_over_time",
		0, 0,
		opts,
		time.Minute, 0,
		0, // batchSize = 0 → no batching
		0, 1,
	)
	testutil.Ok(t, err)

	_, err = op.Series(ctx)
	testutil.Ok(t, err)

	buf := make([]model.StepVector, 6)

	// Call 1: ALL 4 series for steps 1-6.
	n1, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 6, n1)
	seriesInCall1 := collectSeriesIDs(buf[:n1])
	testutil.Assert(t, containsAll(seriesInCall1, []uint64{0, 1, 2, 3}),
		"call 1 should contain ALL series 0,1,2,3, got %v", seriesInCall1)

	// Call 2: ALL 4 series for steps 7-12.
	n2, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 6, n2)
	seriesInCall2 := collectSeriesIDs(buf[:n2])
	testutil.Assert(t, containsAll(seriesInCall2, []uint64{0, 1, 2, 3}),
		"call 2 should contain ALL series 0,1,2,3, got %v", seriesInCall2)

	// Call 3: ALL 4 series for steps 13-14 (partial).
	n3, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 2, n3)
	seriesInCall3 := collectSeriesIDs(buf[:n3])
	testutil.Assert(t, containsAll(seriesInCall3, []uint64{0, 1, 2, 3}),
		"call 3 should contain ALL series 0,1,2,3 (partial), got %v", seriesInCall3)

	// Call 4: done.
	n4, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 0, n4)
}
