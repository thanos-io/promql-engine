// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package remote

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/query"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/stretchr/testify/require"
)

// fakeRemote is a minimal api.RemoteOperator stub for unit tests.
type fakeRemote struct {
	series   []labels.Labels
	batches  [][]model.StepVector // batches to be returned by successive Next calls
	cursor   int
	nextErr  error
	seriesEr error
	closeCnt int32
}

func (f *fakeRemote) Series(ctx context.Context) ([]labels.Labels, error) {
	return f.series, f.seriesEr
}

func (f *fakeRemote) Next(ctx context.Context, buf []model.StepVector) (int, error) {
	if f.nextErr != nil {
		return 0, f.nextErr
	}
	if f.cursor >= len(f.batches) {
		return 0, nil
	}
	batch := f.batches[f.cursor]
	f.cursor++
	n := len(batch)
	if n > len(buf) {
		n = len(buf)
	}
	for i := 0; i < n; i++ {
		buf[i] = batch[i]
	}
	return n, nil
}

func (f *fakeRemote) Close() error {
	atomic.AddInt32(&f.closeCnt, 1)
	return nil
}

func newStepVector(t int64, ids []uint64, vals []float64) model.StepVector {
	sv := model.StepVector{T: t}
	for i := range ids {
		sv.AppendSample(ids[i], vals[i])
	}
	return sv
}

func runToEnd(t *testing.T, op model.VectorOperator) []model.StepVector {
	t.Helper()
	ctx := context.Background()
	var out []model.StepVector
	buf := make([]model.StepVector, 10)
	for {
		// Reset buf between calls because StreamingExecution may swap slice
		// contents from a previous emission back into our buf.
		for i := range buf {
			buf[i] = model.StepVector{}
		}
		n, err := op.Next(ctx, buf)
		require.NoError(t, err)
		if n == 0 {
			break
		}
		for i := 0; i < n; i++ {
			out = append(out, copyStepVector(buf[i]))
		}
	}
	return out
}

func copyStepVector(in model.StepVector) model.StepVector {
	out := model.StepVector{T: in.T}
	for i := range in.SampleIDs {
		out.AppendSample(in.SampleIDs[i], in.Samples[i])
	}
	return out
}

func TestStreamingExecution_ForwardsAlignedRemote(t *testing.T) {
	// Remote produces vectors at exactly the parent step grid; nothing to pad.
	opts := &query.Options{
		Start:      time.UnixMilli(0),
		End:        time.UnixMilli(60),
		Step:       30 * time.Millisecond,
		StepsBatch: 10,
	}
	remote := &fakeRemote{
		series: []labels.Labels{labels.FromStrings("foo", "a")},
		batches: [][]model.StepVector{
			{
				newStepVector(0, []uint64{0}, []float64{1}),
				newStepVector(30, []uint64{0}, []float64{2}),
				newStepVector(60, []uint64{0}, []float64{3}),
			},
		},
	}
	op := NewStreamingExecution(remote, opts.Start, opts.End, nil, opts)

	series, err := op.Series(context.Background())
	require.NoError(t, err)
	require.Equal(t, []labels.Labels{labels.FromStrings("foo", "a")}, series)

	got := runToEnd(t, op)
	require.Len(t, got, 3)
	require.Equal(t, int64(0), got[0].T)
	require.Equal(t, []float64{1}, got[0].Samples)
	require.Equal(t, int64(30), got[1].T)
	require.Equal(t, int64(60), got[2].T)
	require.Equal(t, int32(1), atomic.LoadInt32(&remote.closeCnt))
}

func TestStreamingExecution_PadsLeading(t *testing.T) {
	// Parent grid [0, 90] step 30. Remote produces starting at T=60.
	// Streaming op must emit empty StepVectors at T=0 and T=30 to align with
	// sibling operators in a coalesce.
	opts := &query.Options{
		Start:      time.UnixMilli(0),
		End:        time.UnixMilli(90),
		Step:       30 * time.Millisecond,
		StepsBatch: 10,
	}
	remote := &fakeRemote{
		series: []labels.Labels{labels.FromStrings("foo", "b")},
		batches: [][]model.StepVector{
			{
				newStepVector(60, []uint64{0}, []float64{10}),
				newStepVector(90, []uint64{0}, []float64{11}),
			},
		},
	}
	op := NewStreamingExecution(remote, time.UnixMilli(60), time.UnixMilli(90), nil, opts)

	got := runToEnd(t, op)
	require.Len(t, got, 4)
	require.Equal(t, int64(0), got[0].T)
	require.Empty(t, got[0].SampleIDs)
	require.Equal(t, int64(30), got[1].T)
	require.Empty(t, got[1].SampleIDs)
	require.Equal(t, int64(60), got[2].T)
	require.Equal(t, []float64{10}, got[2].Samples)
	require.Equal(t, int64(90), got[3].T)
	require.Equal(t, []float64{11}, got[3].Samples)
}

func TestStreamingExecution_PadsTrailing(t *testing.T) {
	// Parent grid [0, 90] step 30. Remote covers only [0, 30].
	opts := &query.Options{
		Start:      time.UnixMilli(0),
		End:        time.UnixMilli(90),
		Step:       30 * time.Millisecond,
		StepsBatch: 10,
	}
	remote := &fakeRemote{
		series: []labels.Labels{labels.FromStrings("foo", "c")},
		batches: [][]model.StepVector{
			{
				newStepVector(0, []uint64{0}, []float64{5}),
				newStepVector(30, []uint64{0}, []float64{6}),
			},
		},
	}
	op := NewStreamingExecution(remote, time.UnixMilli(0), time.UnixMilli(30), nil, opts)

	got := runToEnd(t, op)
	require.Len(t, got, 4)
	require.Equal(t, int64(0), got[0].T)
	require.Equal(t, int64(30), got[1].T)
	require.Equal(t, int64(60), got[2].T)
	require.Empty(t, got[2].SampleIDs)
	require.Equal(t, int64(90), got[3].T)
	require.Empty(t, got[3].SampleIDs)
}

func TestStreamingExecution_NormalisesEmptySeriesToSingleEmptyLabels(t *testing.T) {
	// Mirrors the in-engine convention used by scalar/pi/time operators:
	// Series() returns nil but Next emits StepVectors with SampleID=0.
	// StreamingExecution must present this as a single empty-label series so
	// downstream coalesce/dedup operators can interpret the SampleID.
	opts := &query.Options{
		Start:      time.UnixMilli(0),
		End:        time.UnixMilli(30),
		Step:       30 * time.Millisecond,
		StepsBatch: 10,
	}
	remote := &fakeRemote{
		series: nil,
		batches: [][]model.StepVector{
			{
				newStepVector(0, []uint64{0}, []float64{3.14}),
				newStepVector(30, []uint64{0}, []float64{3.14}),
			},
		},
	}
	op := NewStreamingExecution(remote, opts.Start, opts.End, nil, opts)

	series, err := op.Series(context.Background())
	require.NoError(t, err)
	require.Len(t, series, 1)
	require.True(t, series[0].IsEmpty())

	got := runToEnd(t, op)
	require.Len(t, got, 2)
	require.Equal(t, []float64{3.14}, got[0].Samples)
	require.Equal(t, []float64{3.14}, got[1].Samples)
}

func TestStreamingExecution_ClosesRemoteExactlyOnceOnEOF(t *testing.T) {
	opts := &query.Options{
		Start:      time.UnixMilli(0),
		End:        time.UnixMilli(30),
		Step:       30 * time.Millisecond,
		StepsBatch: 10,
	}
	remote := &fakeRemote{
		series: []labels.Labels{labels.FromStrings("foo", "d")},
		batches: [][]model.StepVector{
			{newStepVector(0, []uint64{0}, []float64{1}), newStepVector(30, []uint64{0}, []float64{2})},
		},
	}
	op := NewStreamingExecution(remote, opts.Start, opts.End, nil, opts)

	// Drain twice; second drain must not double-close.
	_ = runToEnd(t, op)
	_ = runToEnd(t, op)
	require.Equal(t, int32(1), atomic.LoadInt32(&remote.closeCnt))
}

func TestStreamingExecution_PropagatesRemoteError(t *testing.T) {
	opts := &query.Options{
		Start:      time.UnixMilli(0),
		End:        time.UnixMilli(30),
		Step:       30 * time.Millisecond,
		StepsBatch: 10,
	}
	wantErr := errors.New("boom")
	remote := &fakeRemote{
		series:  []labels.Labels{labels.FromStrings("foo", "e")},
		nextErr: wantErr,
	}
	op := NewStreamingExecution(remote, opts.Start, opts.End, nil, opts)

	buf := make([]model.StepVector, 10)
	_, err := op.Next(context.Background(), buf)
	require.ErrorIs(t, err, wantErr)
	require.Equal(t, int32(1), atomic.LoadInt32(&remote.closeCnt))
}
