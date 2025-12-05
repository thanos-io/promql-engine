// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package ringbuffer

import (
	"context"
	"math"

	"github.com/thanos-io/promql-engine/compute"
	"github.com/thanos-io/promql-engine/execution/telemetry"
	"github.com/thanos-io/promql-engine/query"
	"github.com/thanos-io/promql-engine/warnings"

	"github.com/prometheus/prometheus/model/histogram"
	"github.com/prometheus/prometheus/util/annotations"
)

// =============================================================================
// SlidingDequeBuffer - min/max using monotonic deque (O(1) amortized)
// =============================================================================

type DequeComparator func(a, b float64) bool

type SlidingDequeBuffer struct {
	samples     []Sample
	tail        []Sample // scratch space for Reset
	lastT       int64
	deque       []int // indices into samples for monotonic deque
	floatCount  int
	sampleCount int
	shouldEvict DequeComparator
}

func NewSlidingMinOverTimeBuffer(_ query.Options, _, _ int64) *SlidingDequeBuffer {
	return &SlidingDequeBuffer{
		samples:     make([]Sample, 0, 8),
		lastT:       math.MinInt64,
		deque:       make([]int, 0, 8),
		shouldEvict: func(back, new float64) bool { return back >= new },
	}
}

func NewSlidingMaxOverTimeBuffer(_ query.Options, _, _ int64) *SlidingDequeBuffer {
	return &SlidingDequeBuffer{
		samples:     make([]Sample, 0, 8),
		lastT:       math.MinInt64,
		deque:       make([]int, 0, 8),
		shouldEvict: func(back, new float64) bool { return back <= new },
	}
}

func (r *SlidingDequeBuffer) SampleCount() int           { return r.sampleCount }
func (r *SlidingDequeBuffer) MaxT() int64                { return r.lastT }
func (r *SlidingDequeBuffer) ReadIntoLast(func(*Sample)) {}

func (r *SlidingDequeBuffer) Push(t int64, v Value) {
	r.lastT = t
	idx := len(r.samples)

	// Append sample (same as GenericRingBuffer)
	n := len(r.samples)
	if n < cap(r.samples) {
		r.samples = r.samples[:n+1]
	} else {
		r.samples = append(r.samples, Sample{})
	}
	r.samples[n].T = t
	r.samples[n].V.F = v.F
	if v.H != nil {
		r.sampleCount += telemetry.CalculateHistogramSampleCount(v.H)
		if r.samples[n].V.H == nil {
			r.samples[n].V.H = v.H.Copy()
		} else {
			v.H.CopyTo(r.samples[n].V.H)
		}
		return
	}
	r.samples[n].V.H = nil

	r.sampleCount++
	r.floatCount++

	if math.IsNaN(v.F) {
		return
	}

	// Maintain monotonic deque
	for len(r.deque) > 0 {
		backIdx := r.deque[len(r.deque)-1]
		backVal := r.samples[backIdx].V.F
		if r.shouldEvict(backVal, v.F) || math.IsNaN(backVal) {
			r.deque = r.deque[:len(r.deque)-1]
		} else {
			break
		}
	}
	r.deque = append(r.deque, idx)
}

func (r *SlidingDequeBuffer) Reset(mint int64, _ int64) {
	if len(r.samples) == 0 || r.samples[len(r.samples)-1].T <= mint {
		r.samples = r.samples[:0]
		r.deque = r.deque[:0]
		r.sampleCount = 0
		r.floatCount = 0
		return
	}

	// Find first sample to keep
	var drop int
	for drop = 0; drop < len(r.samples) && r.samples[drop].T <= mint; drop++ {
		s := &r.samples[drop]
		if s.V.H != nil {
			r.sampleCount -= telemetry.CalculateHistogramSampleCount(s.V.H)
		} else {
			r.sampleCount--
			r.floatCount--
		}
	}

	if drop == 0 {
		return
	}

	// Remove old indices from deque and adjust remaining
	var dequeDrop int
	for dequeDrop = 0; dequeDrop < len(r.deque) && r.deque[dequeDrop] < drop; dequeDrop++ {
	}
	keep := len(r.deque) - dequeDrop
	copy(r.deque, r.deque[dequeDrop:])
	r.deque = r.deque[:keep]
	for i := range r.deque {
		r.deque[i] -= drop
	}

	// Shift samples (same as GenericRingBuffer)
	keep = len(r.samples) - drop
	r.tail = resize(r.tail, drop)
	copy(r.tail, r.samples[:drop])
	copy(r.samples, r.samples[drop:])
	copy(r.samples[keep:], r.tail)
	r.samples = r.samples[:keep]
}

func (r *SlidingDequeBuffer) Eval(_ context.Context, _, _ float64, _ int64) (float64, *histogram.FloatHistogram, bool, error) {
	if len(r.samples) == 0 || r.floatCount == 0 {
		return 0, nil, false, nil
	}
	if len(r.deque) == 0 {
		return math.NaN(), nil, true, nil
	}
	return r.samples[r.deque[0]].V.F, nil, true, nil
}

// =============================================================================
// SlidingAccumulatorBuffer - sliding buffer with O(removed) updates
// =============================================================================

// SlidingAccumulatorBuffer wraps a compute.CheckpointableAccumulator with sliding window support.
// It provides O(removed) Reset by subtracting removed samples.
type SlidingAccumulatorBuffer struct {
	samples     []Sample
	tail        []Sample // scratch space for Reset
	lastT       int64
	acc         compute.CheckpointableAccumulator
	sampleCount int
	warn        error
}

func NewSlidingSumOverTimeBuffer(_ query.Options, _, _ int64) *SlidingAccumulatorBuffer {
	return &SlidingAccumulatorBuffer{
		samples: make([]Sample, 0, 8),
		lastT:   math.MinInt64,
		acc:     compute.NewSumAcc(),
	}
}

func NewSlidingAvgOverTimeBuffer(_ query.Options, _, _ int64) *SlidingAccumulatorBuffer {
	return &SlidingAccumulatorBuffer{
		samples: make([]Sample, 0, 8),
		lastT:   math.MinInt64,
		acc:     compute.NewAvgAcc(),
	}
}

func (r *SlidingAccumulatorBuffer) SampleCount() int           { return r.sampleCount }
func (r *SlidingAccumulatorBuffer) MaxT() int64                { return r.lastT }
func (r *SlidingAccumulatorBuffer) ReadIntoLast(func(*Sample)) {}

func (r *SlidingAccumulatorBuffer) Push(t int64, v Value) {
	r.lastT = t

	// Append sample (same as GenericRingBuffer)
	n := len(r.samples)
	if n < cap(r.samples) {
		r.samples = r.samples[:n+1]
	} else {
		r.samples = append(r.samples, Sample{})
	}
	r.samples[n].T = t
	r.samples[n].V.F = v.F
	if v.H != nil {
		r.sampleCount += telemetry.CalculateHistogramSampleCount(v.H)
		if r.samples[n].V.H == nil {
			r.samples[n].V.H = v.H.Copy()
		} else {
			v.H.CopyTo(r.samples[n].V.H)
		}
	} else {
		r.sampleCount++
		r.samples[n].V.H = nil
	}

	if err := r.acc.Add(v.F, v.H); err != nil {
		r.warn = err
	}
}

func (r *SlidingAccumulatorBuffer) Reset(mint int64, _ int64) {
	r.warn = nil

	if len(r.samples) == 0 || r.samples[len(r.samples)-1].T <= mint {
		r.samples = r.samples[:0]
		r.sampleCount = 0
		r.acc.Reset(0)
		return
	}

	// Find first sample to keep, subtracting removed samples
	var drop int
	for drop = 0; drop < len(r.samples) && r.samples[drop].T <= mint; drop++ {
		s := &r.samples[drop]
		if err := r.acc.Sub(s.V.F, s.V.H); err != nil {
			r.warn = err
		}
		if s.V.H != nil {
			r.sampleCount -= telemetry.CalculateHistogramSampleCount(s.V.H)
		} else {
			r.sampleCount--
		}
	}

	if drop == 0 {
		return
	}

	// Shift samples (same as GenericRingBuffer)
	keep := len(r.samples) - drop
	r.tail = resize(r.tail, drop)
	copy(r.tail, r.samples[:drop])
	copy(r.samples, r.samples[drop:])
	copy(r.samples[keep:], r.tail)
	r.samples = r.samples[:keep]
}

func (r *SlidingAccumulatorBuffer) Eval(ctx context.Context, _, _ float64, _ int64) (float64, *histogram.FloatHistogram, bool, error) {
	if r.warn != nil {
		warnings.AddToContext(r.warn, ctx)
	}

	if len(r.samples) == 0 {
		return 0, nil, false, nil
	}

	valueType := r.acc.ValueType()
	if valueType == compute.NoValue {
		return 0, nil, false, nil
	}
	if valueType == compute.MixedTypeValue {
		warnings.AddToContext(annotations.MixedFloatsHistogramsWarning, ctx)
		return 0, nil, false, nil
	}

	f, h := r.acc.Value()
	if h != nil {
		h = h.Copy()
		h.Compact(0)
	}
	return f, h, true, nil
}
