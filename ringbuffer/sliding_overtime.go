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
// Sliding buffer base with circular sample storage
// =============================================================================

type slidingSample struct {
	t int64
	f float64
	h *histogram.FloatHistogram
}

type slidingBase struct {
	samples  []slidingSample
	head     int
	size     int
	capacity int
	lastT    int64

	currentMint int64
	selectRange int64
	step        int64
	offset      int64
}

func newSlidingBase(opts query.Options, selectRange, offset int64) slidingBase {
	step := max(1, opts.Step.Milliseconds())
	capacity := int((selectRange / 15000) + 10)
	if capacity < 16 {
		capacity = 16
	}

	return slidingBase{
		samples:     make([]slidingSample, capacity),
		capacity:    capacity,
		lastT:       math.MinInt64,
		currentMint: math.MaxInt64,
		selectRange: selectRange,
		step:        step,
		offset:      offset,
	}
}

func (b *slidingBase) MaxT() int64                   { return b.lastT }
func (b *slidingBase) ReadIntoLast(func(*Sample))    {}
func (b *slidingBase) oldestIndex() int              { return (b.head - b.size + b.capacity) % b.capacity }
func (b *slidingBase) indexAt(i int) int             { return (b.head - b.size + i + b.capacity) % b.capacity }
func (b *slidingBase) sampleAt(i int) *slidingSample { return &b.samples[b.indexAt(i)] }

func (b *slidingBase) push(t int64, f float64, h *histogram.FloatHistogram) int {
	b.lastT = t
	if b.size >= b.capacity {
		b.grow()
	}
	idx := b.head
	b.samples[idx] = slidingSample{t: t, f: f, h: h}
	b.head = (b.head + 1) % b.capacity
	b.size++
	return idx
}

func (b *slidingBase) grow() {
	newCapacity := b.capacity * 2
	newSamples := make([]slidingSample, newCapacity)
	for i := 0; i < b.size; i++ {
		newSamples[i] = b.samples[b.indexAt(i)]
	}
	b.samples = newSamples
	b.capacity = newCapacity
	b.head = b.size
}

// removeOldSamples removes samples with t <= mint and returns the count removed.
func (b *slidingBase) removeOldSamples(mint int64) int {
	removed := 0
	for b.size > 0 {
		if b.samples[b.oldestIndex()].t > mint {
			break
		}
		b.size--
		removed++
	}
	return removed
}

// =============================================================================
// SlidingDequeBuffer - min/max using monotonic deque (O(1) amortized)
// =============================================================================

type DequeComparator func(a, b float64) bool

type SlidingDequeBuffer struct {
	slidingBase
	deque       []int
	floatCount  int
	sampleCount int
	shouldEvict DequeComparator
}

func NewSlidingMinOverTimeBuffer(opts query.Options, selectRange, offset int64) *SlidingDequeBuffer {
	base := newSlidingBase(opts, selectRange, offset)
	return &SlidingDequeBuffer{
		slidingBase: base,
		deque:       make([]int, 0, base.capacity),
		shouldEvict: func(back, new float64) bool { return back >= new },
	}
}

func NewSlidingMaxOverTimeBuffer(opts query.Options, selectRange, offset int64) *SlidingDequeBuffer {
	base := newSlidingBase(opts, selectRange, offset)
	return &SlidingDequeBuffer{
		slidingBase: base,
		deque:       make([]int, 0, base.capacity),
		shouldEvict: func(back, new float64) bool { return back <= new },
	}
}

func (r *SlidingDequeBuffer) SampleCount() int { return r.sampleCount }

func (r *SlidingDequeBuffer) Push(t int64, v Value) {
	r.lastT = t

	if v.H != nil {
		r.sampleCount += telemetry.CalculateHistogramSampleCount(v.H)
		if r.size >= r.capacity {
			r.growWithDeque()
		}
		r.samples[r.head] = slidingSample{t: t, f: math.NaN(), h: v.H}
		r.head = (r.head + 1) % r.capacity
		r.size++
		return
	}

	r.sampleCount++
	r.floatCount++

	if r.size >= r.capacity {
		r.growWithDeque()
	}
	idx := r.head
	r.samples[idx] = slidingSample{t: t, f: v.F}
	r.head = (r.head + 1) % r.capacity
	r.size++

	if math.IsNaN(v.F) {
		return
	}

	// Maintain monotonic deque
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

func (r *SlidingDequeBuffer) Reset(mint int64, evalt int64) {
	r.currentMint = mint

	for r.size > 0 {
		oldestIdx := r.oldestIndex()
		s := &r.samples[oldestIdx]
		if s.t > mint {
			break
		}
		// Decrement sample count
		if s.h != nil {
			r.sampleCount -= telemetry.CalculateHistogramSampleCount(s.h)
		} else {
			r.sampleCount--
			r.floatCount--
		}
		r.size--
		if len(r.deque) > 0 && r.deque[0] == oldestIdx {
			r.deque = r.deque[1:]
		}
	}
}

func (r *SlidingDequeBuffer) Eval(ctx context.Context, _, _ float64, _ int64) (float64, *histogram.FloatHistogram, bool, error) {
	if r.size == 0 || r.floatCount == 0 {
		return 0, nil, false, nil
	}
	if len(r.deque) == 0 {
		return math.NaN(), nil, true, nil
	}
	return r.samples[r.deque[0]].f, nil, true, nil
}

func (r *SlidingDequeBuffer) growWithDeque() {
	oldHead, oldCap, oldSize := r.head, r.capacity, r.size
	r.slidingBase.grow()
	newDeque := make([]int, 0, len(r.deque))
	for _, oldIdx := range r.deque {
		newDeque = append(newDeque, (oldIdx-(oldHead-oldSize)+oldCap)%oldCap)
	}
	r.deque = newDeque
}

// =============================================================================
// SlidingAccumulatorBuffer - Generic sliding buffer using compute.Accumulator
// with O(removed) updates for subtractable accumulators.
// =============================================================================

// SlidingAccumulatorBuffer wraps a compute.CheckpointableAccumulator with sliding window support.
// It provides O(removed) Reset by subtracting removed samples.
type SlidingAccumulatorBuffer struct {
	slidingBase

	acc        compute.CheckpointableAccumulator
	newAccFunc func() compute.CheckpointableAccumulator

	sampleCount int
	warn        error
}

func newSlidingAccumulatorBuffer(opts query.Options, selectRange, offset int64, newAcc func() compute.CheckpointableAccumulator) *SlidingAccumulatorBuffer {
	base := newSlidingBase(opts, selectRange, offset)

	return &SlidingAccumulatorBuffer{
		slidingBase: base,
		acc:         newAcc(),
		newAccFunc:  newAcc,
	}
}

func NewSlidingSumOverTimeBuffer(opts query.Options, selectRange, offset int64) *SlidingAccumulatorBuffer {
	return newSlidingAccumulatorBuffer(opts, selectRange, offset, func() compute.CheckpointableAccumulator { return compute.NewSumAcc() })
}

func NewSlidingAvgOverTimeBuffer(opts query.Options, selectRange, offset int64) *SlidingAccumulatorBuffer {
	return newSlidingAccumulatorBuffer(opts, selectRange, offset, func() compute.CheckpointableAccumulator { return compute.NewAvgAcc() })
}

func (r *SlidingAccumulatorBuffer) SampleCount() int { return r.sampleCount }

func (r *SlidingAccumulatorBuffer) Push(t int64, v Value) {
	r.lastT = t

	if r.size >= r.capacity {
		r.grow()
	}

	var h *histogram.FloatHistogram
	if v.H != nil {
		h = v.H.Copy()
		r.sampleCount += telemetry.CalculateHistogramSampleCount(v.H)
	} else {
		r.sampleCount++
	}

	r.samples[r.head] = slidingSample{t: t, f: v.F, h: h}
	r.head = (r.head + 1) % r.capacity
	r.size++

	// Add to main accumulator
	if err := r.acc.Add(v.F, v.H); err != nil {
		r.warn = err
	}

}

func (r *SlidingAccumulatorBuffer) Reset(mint int64, evalt int64) {
	r.currentMint = mint
	r.warn = nil

	// O(removed) update: subtract the samples being removed
	for r.size > 0 {
		s := r.sampleAt(0)
		if s.t > mint {
			break
		}
		if err := r.acc.Sub(s.f, s.h); err != nil {
			r.warn = err
		}
		// Decrement sample count
		if s.h != nil {
			r.sampleCount -= telemetry.CalculateHistogramSampleCount(s.h)
		} else {
			r.sampleCount--
		}
		r.size--
	}
}

func (r *SlidingAccumulatorBuffer) Eval(ctx context.Context, _, _ float64, _ int64) (float64, *histogram.FloatHistogram, bool, error) {
	if r.warn != nil {
		warnings.AddToContext(r.warn, ctx)
	}

	// No samples in buffer = no value
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

func (r *SlidingAccumulatorBuffer) grow() {
	r.slidingBase.grow()
	// Checkpoint remains valid - indices don't change, just internal buffer layout
}
