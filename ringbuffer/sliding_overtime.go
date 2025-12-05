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

type slidingSample struct {
	t int64
	f float64
	h *histogram.FloatHistogram
}

// estimateCapacity estimates the number of samples in a range window.
// Assumes ~15s scrape interval as a reasonable default.
func estimateCapacity(selectRange int64) int {
	// selectRange is in milliseconds, assume 15s scrape interval
	capacity := int(selectRange/15000) + 1
	if capacity < 8 {
		capacity = 8
	}
	return capacity
}

// =============================================================================
// SlidingDequeBuffer - min/max using monotonic deque (O(1) amortized)
// =============================================================================

type DequeComparator func(a, b float64) bool

type SlidingDequeBuffer struct {
	samples []slidingSample
	size    int
	lastT   int64

	deque       []int
	floatCount  int
	sampleCount int
	shouldEvict DequeComparator
}

func NewSlidingMinOverTimeBuffer(_ query.Options, selectRange, _ int64) *SlidingDequeBuffer {
	capacity := estimateCapacity(selectRange)
	return &SlidingDequeBuffer{
		samples:     make([]slidingSample, 0, capacity),
		lastT:       math.MinInt64,
		deque:       make([]int, 0, capacity),
		shouldEvict: func(back, new float64) bool { return back >= new },
	}
}

func NewSlidingMaxOverTimeBuffer(_ query.Options, selectRange, _ int64) *SlidingDequeBuffer {
	capacity := estimateCapacity(selectRange)
	return &SlidingDequeBuffer{
		samples:     make([]slidingSample, 0, capacity),
		lastT:       math.MinInt64,
		deque:       make([]int, 0, capacity),
		shouldEvict: func(back, new float64) bool { return back <= new },
	}
}

func (r *SlidingDequeBuffer) SampleCount() int           { return r.sampleCount }
func (r *SlidingDequeBuffer) MaxT() int64                { return r.lastT }
func (r *SlidingDequeBuffer) ReadIntoLast(func(*Sample)) {}

func (r *SlidingDequeBuffer) Push(t int64, v Value) {
	r.lastT = t
	idx := len(r.samples)

	if v.H != nil {
		r.sampleCount += telemetry.CalculateHistogramSampleCount(v.H)
		r.samples = append(r.samples, slidingSample{t: t, f: math.NaN(), h: v.H})
		r.size++
		return
	}

	r.sampleCount++
	r.floatCount++
	r.samples = append(r.samples, slidingSample{t: t, f: v.F})
	r.size++

	if math.IsNaN(v.F) {
		return
	}

	// Maintain monotonic deque - store absolute indices
	for len(r.deque) > 0 {
		backVal := r.samples[r.deque[len(r.deque)-1]].f
		if r.shouldEvict(backVal, v.F) || math.IsNaN(backVal) {
			r.deque = r.deque[:len(r.deque)-1]
		} else {
			break
		}
	}
	r.deque = append(r.deque, idx)
}

func (r *SlidingDequeBuffer) Reset(mint int64, _ int64) {
	// Find first sample with t > mint
	start := len(r.samples) - r.size
	newStart := start
	for i := start; i < len(r.samples); i++ {
		if r.samples[i].t > mint {
			break
		}
		s := &r.samples[i]
		if s.h != nil {
			r.sampleCount -= telemetry.CalculateHistogramSampleCount(s.h)
		} else {
			r.sampleCount--
			r.floatCount--
		}
		newStart = i + 1
	}
	r.size = len(r.samples) - newStart

	// Remove old indices from deque front
	for len(r.deque) > 0 && r.deque[0] < newStart {
		r.deque = r.deque[1:]
	}

	// Compact when start position exceeds size (more garbage than data)
	if newStart > r.size {
		r.compact(newStart)
	}
}

func (r *SlidingDequeBuffer) compact(start int) {
	n := copy(r.samples, r.samples[start:])
	r.samples = r.samples[:n]
	// Adjust deque indices
	for i := range r.deque {
		r.deque[i] -= start
	}
}

func (r *SlidingDequeBuffer) Eval(_ context.Context, _, _ float64, _ int64) (float64, *histogram.FloatHistogram, bool, error) {
	if r.size == 0 || r.floatCount == 0 {
		return 0, nil, false, nil
	}
	if len(r.deque) == 0 {
		return math.NaN(), nil, true, nil
	}
	return r.samples[r.deque[0]].f, nil, true, nil
}

// =============================================================================
// SlidingAccumulatorBuffer - sliding buffer with O(removed) updates
// =============================================================================

// SlidingAccumulatorBuffer wraps a compute.CheckpointableAccumulator with sliding window support.
// It provides O(removed) Reset by subtracting removed samples.
type SlidingAccumulatorBuffer struct {
	samples []slidingSample
	size    int
	lastT   int64

	acc         compute.CheckpointableAccumulator
	sampleCount int
	warn        error
}

func NewSlidingSumOverTimeBuffer(_ query.Options, selectRange, _ int64) *SlidingAccumulatorBuffer {
	capacity := estimateCapacity(selectRange)
	return &SlidingAccumulatorBuffer{
		samples: make([]slidingSample, 0, capacity),
		lastT:   math.MinInt64,
		acc:     compute.NewSumAcc(),
	}
}

func NewSlidingAvgOverTimeBuffer(_ query.Options, selectRange, _ int64) *SlidingAccumulatorBuffer {
	capacity := estimateCapacity(selectRange)
	return &SlidingAccumulatorBuffer{
		samples: make([]slidingSample, 0, capacity),
		lastT:   math.MinInt64,
		acc:     compute.NewAvgAcc(),
	}
}

func (r *SlidingAccumulatorBuffer) SampleCount() int           { return r.sampleCount }
func (r *SlidingAccumulatorBuffer) MaxT() int64                { return r.lastT }
func (r *SlidingAccumulatorBuffer) ReadIntoLast(func(*Sample)) {}

func (r *SlidingAccumulatorBuffer) Push(t int64, v Value) {
	r.lastT = t

	var h *histogram.FloatHistogram
	if v.H != nil {
		h = v.H.Copy()
		r.sampleCount += telemetry.CalculateHistogramSampleCount(v.H)
	} else {
		r.sampleCount++
	}

	r.samples = append(r.samples, slidingSample{t: t, f: v.F, h: h})
	r.size++

	if err := r.acc.Add(v.F, v.H); err != nil {
		r.warn = err
	}
}

func (r *SlidingAccumulatorBuffer) Reset(mint int64, _ int64) {
	r.warn = nil

	// Find first sample with t > mint, subtracting removed samples
	start := len(r.samples) - r.size
	newStart := start
	for i := start; i < len(r.samples); i++ {
		s := &r.samples[i]
		if s.t > mint {
			break
		}
		if err := r.acc.Sub(s.f, s.h); err != nil {
			r.warn = err
		}
		if s.h != nil {
			r.sampleCount -= telemetry.CalculateHistogramSampleCount(s.h)
		} else {
			r.sampleCount--
		}
		newStart = i + 1
	}
	r.size = len(r.samples) - newStart

	// Compact when start position exceeds size (more garbage than data)
	if newStart > r.size {
		n := copy(r.samples, r.samples[newStart:])
		r.samples = r.samples[:n]
	}
}

func (r *SlidingAccumulatorBuffer) Eval(ctx context.Context, _, _ float64, _ int64) (float64, *histogram.FloatHistogram, bool, error) {
	if r.warn != nil {
		warnings.AddToContext(r.warn, ctx)
	}

	if r.size == 0 {
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
