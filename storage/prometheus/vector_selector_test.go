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

// TestSeriesBatchVectorSelectorCompletesAllStepsPerBatch verifies the same
// series-first behavior for the vectorSelector (instant vector, no range function).
func TestSeriesBatchVectorSelectorCompletesAllStepsPerBatch(t *testing.T) {
	// 4 series, 14 steps, stepsBatch=6, batchSize=2.
	// Same pattern as matrixSelector test: 3 calls per batch (6+6+2).

	load := `load 30s
		metric{s="0"} 1+1x30
		metric{s="1"} 2+1x30
		metric{s="2"} 3+1x30
		metric{s="3"} 4+1x30`

	testStorage := promqltest.LoadedStorage(t, load)
	defer testStorage.Close()

	ctx := context.Background()
	mint := int64(0)
	maxt := int64(390000) // 14 steps at 30s from 0

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

	// Create vectorSelector with batchSize=2.
	op := NewVectorSelector(
		selector,
		opts,
		0,     // offset
		2,     // batchSize = 2
		false, // selectTimestamp
		0, 1,  // shard, numShards
	)

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

	// Call 2: batch1 (series 0,1), steps 7-12.
	n2, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 6, n2)
	seriesInCall2 := collectSeriesIDs(buf[:n2])
	testutil.Assert(t, containsAll(seriesInCall2, []uint64{0, 1}),
		"call 2 should contain series 0,1 (same batch), got %v", seriesInCall2)
	testutil.Assert(t, !containsAny(seriesInCall2, []uint64{2, 3}),
		"call 2 should NOT contain series 2,3, got %v", seriesInCall2)

	// Call 3: batch1 (series 0,1), steps 13-14 (partial).
	n3, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 2, n3)
	seriesInCall3 := collectSeriesIDs(buf[:n3])
	testutil.Assert(t, containsAll(seriesInCall3, []uint64{0, 1}),
		"call 3 should contain series 0,1 (partial), got %v", seriesInCall3)

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
		"call 5 should contain series 2,3, got %v", seriesInCall5)

	// Call 6: batch2 (series 2,3), steps 13-14 (partial).
	n6, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 2, n6)
	seriesInCall6 := collectSeriesIDs(buf[:n6])
	testutil.Assert(t, containsAll(seriesInCall6, []uint64{2, 3}),
		"call 6 should contain series 2,3 (partial), got %v", seriesInCall6)

	// Call 7: all batches exhausted.
	n7, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 0, n7)
}

// TestVectorSelectorWithoutBatchingOriginalBehavior verifies original ordering
// for vectorSelector without batching.
func TestVectorSelectorWithoutBatchingOriginalBehavior(t *testing.T) {
	load := `load 30s
		metric{s="0"} 1+1x30
		metric{s="1"} 2+1x30
		metric{s="2"} 3+1x30
		metric{s="3"} 4+1x30`

	testStorage := promqltest.LoadedStorage(t, load)
	defer testStorage.Close()

	ctx := context.Background()
	mint := int64(0)
	maxt := int64(390000)

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

	// batchSize=0 → no batching.
	op := NewVectorSelector(selector, opts, 0, 0, false, 0, 1)

	_, err = op.Series(ctx)
	testutil.Ok(t, err)

	buf := make([]model.StepVector, 6)

	// Call 1: ALL series for steps 1-6.
	n1, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 6, n1)
	seriesInCall1 := collectSeriesIDs(buf[:n1])
	testutil.Assert(t, containsAll(seriesInCall1, []uint64{0, 1, 2, 3}),
		"call 1 should contain ALL series, got %v", seriesInCall1)

	// Call 2: ALL series for steps 7-12.
	n2, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 6, n2)
	seriesInCall2 := collectSeriesIDs(buf[:n2])
	testutil.Assert(t, containsAll(seriesInCall2, []uint64{0, 1, 2, 3}),
		"call 2 should contain ALL series, got %v", seriesInCall2)

	// Call 3: ALL series for steps 13-14 (partial).
	n3, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 2, n3)
	seriesInCall3 := collectSeriesIDs(buf[:n3])
	testutil.Assert(t, containsAll(seriesInCall3, []uint64{0, 1, 2, 3}),
		"call 3 should contain ALL series (partial), got %v", seriesInCall3)

	// Call 4: done.
	n4, err := op.Next(ctx, buf)
	testutil.Ok(t, err)
	testutil.Equals(t, 0, n4)
}
