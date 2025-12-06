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
	"github.com/prometheus/prometheus/promql/parser/posrange"
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

// =============================================================================
// SlidingCountBuffer - count_over_time with O(1) updates
// =============================================================================

// SlidingCountBuffer counts samples in a sliding window.
// It only stores timestamps (not values) since count doesn't need values.
type SlidingCountBuffer struct {
	timestamps  []int64
	tail        []int64 // scratch space for Reset
	lastT       int64
	count       int
	sampleCount int
}

func NewSlidingCountOverTimeBuffer(_ query.Options, _, _ int64) *SlidingCountBuffer {
	return &SlidingCountBuffer{
		timestamps: make([]int64, 0, 8),
		lastT:      math.MinInt64,
	}
}

func (r *SlidingCountBuffer) SampleCount() int           { return r.sampleCount }
func (r *SlidingCountBuffer) MaxT() int64                { return r.lastT }
func (r *SlidingCountBuffer) ReadIntoLast(func(*Sample)) {}

func (r *SlidingCountBuffer) Push(t int64, v Value) {
	r.lastT = t

	// Append timestamp
	n := len(r.timestamps)
	if n < cap(r.timestamps) {
		r.timestamps = r.timestamps[:n+1]
	} else {
		r.timestamps = append(r.timestamps, 0)
	}
	r.timestamps[n] = t

	r.count++
	if v.H != nil {
		r.sampleCount += telemetry.CalculateHistogramSampleCount(v.H)
	} else {
		r.sampleCount++
	}
}

func (r *SlidingCountBuffer) Reset(mint int64, _ int64) {
	if len(r.timestamps) == 0 || r.timestamps[len(r.timestamps)-1] <= mint {
		r.timestamps = r.timestamps[:0]
		r.count = 0
		r.sampleCount = 0
		return
	}

	// Find first sample to keep
	var drop int
	for drop = 0; drop < len(r.timestamps) && r.timestamps[drop] <= mint; drop++ {
		r.count--
		r.sampleCount-- // Approximate - we don't track histogram sample counts individually
	}

	if drop == 0 {
		return
	}

	// Shift timestamps
	keep := len(r.timestamps) - drop
	if cap(r.tail) < drop {
		r.tail = make([]int64, drop)
	} else {
		r.tail = r.tail[:drop]
	}
	copy(r.tail, r.timestamps[:drop])
	copy(r.timestamps, r.timestamps[drop:])
	copy(r.timestamps[keep:], r.tail)
	r.timestamps = r.timestamps[:keep]
}

func (r *SlidingCountBuffer) Eval(_ context.Context, _, _ float64, _ int64) (float64, *histogram.FloatHistogram, bool, error) {
	if r.count == 0 {
		return 0, nil, false, nil
	}
	return float64(r.count), nil, true, nil
}

// =============================================================================
// SlidingPresentBuffer - present_over_time with O(1) updates
// =============================================================================

// SlidingPresentBuffer returns 1 if any samples exist in the window.
// Reuses SlidingCountBuffer logic.
type SlidingPresentBuffer struct {
	SlidingCountBuffer
}

func NewSlidingPresentOverTimeBuffer(_ query.Options, _, _ int64) *SlidingPresentBuffer {
	return &SlidingPresentBuffer{
		SlidingCountBuffer: SlidingCountBuffer{
			timestamps: make([]int64, 0, 8),
			lastT:      math.MinInt64,
		},
	}
}

func (r *SlidingPresentBuffer) Eval(_ context.Context, _, _ float64, _ int64) (float64, *histogram.FloatHistogram, bool, error) {
	if r.count == 0 {
		return 0, nil, false, nil
	}
	return 1, nil, true, nil
}

// =============================================================================
// SlidingLastBuffer - last_over_time with O(1) updates
// =============================================================================

// SlidingLastBuffer returns the last sample value in the window.
// Only needs to store timestamps for window management, and the last value.
type SlidingLastBuffer struct {
	timestamps  []int64
	tail        []int64 // scratch space for Reset
	lastT       int64
	lastV       Value
	hasValue    bool
	sampleCount int
}

func NewSlidingLastOverTimeBuffer(_ query.Options, _, _ int64) *SlidingLastBuffer {
	return &SlidingLastBuffer{
		timestamps: make([]int64, 0, 8),
		lastT:      math.MinInt64,
	}
}

func (r *SlidingLastBuffer) SampleCount() int           { return r.sampleCount }
func (r *SlidingLastBuffer) MaxT() int64                { return r.lastT }
func (r *SlidingLastBuffer) ReadIntoLast(func(*Sample)) {}

func (r *SlidingLastBuffer) Push(t int64, v Value) {
	r.lastT = t
	r.hasValue = true

	// Store last value
	r.lastV.F = v.F
	if v.H != nil {
		if r.lastV.H == nil {
			r.lastV.H = v.H.Copy()
		} else {
			v.H.CopyTo(r.lastV.H)
		}
		r.sampleCount += telemetry.CalculateHistogramSampleCount(v.H)
	} else {
		r.lastV.H = nil
		r.sampleCount++
	}

	// Append timestamp
	n := len(r.timestamps)
	if n < cap(r.timestamps) {
		r.timestamps = r.timestamps[:n+1]
	} else {
		r.timestamps = append(r.timestamps, 0)
	}
	r.timestamps[n] = t
}

func (r *SlidingLastBuffer) Reset(mint int64, _ int64) {
	if len(r.timestamps) == 0 || r.timestamps[len(r.timestamps)-1] <= mint {
		r.timestamps = r.timestamps[:0]
		r.hasValue = false
		r.sampleCount = 0
		return
	}

	// Find first sample to keep
	var drop int
	for drop = 0; drop < len(r.timestamps) && r.timestamps[drop] <= mint; drop++ {
		r.sampleCount--
	}

	if drop == 0 {
		return
	}

	// Shift timestamps
	keep := len(r.timestamps) - drop
	if cap(r.tail) < drop {
		r.tail = make([]int64, drop)
	} else {
		r.tail = r.tail[:drop]
	}
	copy(r.tail, r.timestamps[:drop])
	copy(r.timestamps, r.timestamps[drop:])
	copy(r.timestamps[keep:], r.tail)
	r.timestamps = r.timestamps[:keep]

	// hasValue and lastV remain valid since last sample is still in window
}

func (r *SlidingLastBuffer) Eval(_ context.Context, _, _ float64, _ int64) (float64, *histogram.FloatHistogram, bool, error) {
	if !r.hasValue || len(r.timestamps) == 0 {
		return 0, nil, false, nil
	}
	if r.lastV.H != nil {
		return 0, r.lastV.H.Copy(), true, nil
	}
	return r.lastV.F, nil, true, nil
}

// =============================================================================
// SlidingRateBuffer - rate/increase/delta with sliding window optimization
// =============================================================================

// SlidingRateBuffer provides an optimized sliding window implementation for
// rate, increase, and delta functions. Instead of storing all counter reset
// samples, it tracks the cumulative counter correction incrementally.
//
// Key optimizations:
// - O(1) Push: just update last sample and detect resets
// - O(removed) Reset: only process samples leaving the window
// - Memory: O(window_size) for samples, no unbounded reset storage
type SlidingRateBuffer struct {
	ctx context.Context

	// samples stores all samples in the current window for counter reset tracking
	samples []Sample
	tail    []Sample // scratch space for Reset

	// Tracking for rate calculation
	lastT               int64
	counterCorrection   float64 // cumulative counter correction for current window
	numSamples          int
	sampleCount         int // for telemetry
	selectRange         int64
	offset              int64
	isCounter           bool
	isRate              bool
	evalTs              int64
	currentMint         int64
	hasHistograms       bool
	histogramResets     []Sample // only used when histograms are detected
	histogramResetsTail []Sample
}

// NewSlidingRateBuffer creates a new sliding rate buffer.
func NewSlidingRateBuffer(ctx context.Context, opts query.Options, isCounter, isRate bool, selectRange, offset int64) *SlidingRateBuffer {
	return &SlidingRateBuffer{
		ctx:         ctx,
		samples:     make([]Sample, 0, 8),
		lastT:       math.MinInt64,
		selectRange: selectRange,
		offset:      offset,
		isCounter:   isCounter,
		isRate:      isRate,
		currentMint: math.MinInt64,
	}
}

func (r *SlidingRateBuffer) SampleCount() int           { return r.sampleCount }
func (r *SlidingRateBuffer) MaxT() int64                { return r.lastT }
func (r *SlidingRateBuffer) ReadIntoLast(func(*Sample)) {}

func (r *SlidingRateBuffer) Push(t int64, v Value) {
	// Detect counter reset before updating last sample
	if r.isCounter && len(r.samples) > 0 {
		last := &r.samples[len(r.samples)-1]
		if last.T > r.currentMint {
			if v.H != nil && last.V.H != nil {
				// Histogram reset detection
				r.hasHistograms = true
				if v.H.DetectReset(last.V.H) {
					r.histogramResets = append(r.histogramResets, Sample{
						T: last.T,
						V: Value{H: last.V.H.Copy()},
					})
					r.histogramResets = append(r.histogramResets, Sample{
						T: t,
						V: Value{H: v.H.Copy()},
					})
				}
			} else if v.H == nil && last.V.H == nil && last.V.F > v.F {
				// Float counter reset - add the pre-reset value to correction
				r.counterCorrection += last.V.F
			}
		}
	}

	r.lastT = t

	// Append sample
	n := len(r.samples)
	if n < cap(r.samples) {
		r.samples = r.samples[:n+1]
	} else {
		r.samples = append(r.samples, Sample{})
	}
	r.samples[n].T = t
	r.samples[n].V.F = v.F
	if v.H != nil {
		r.hasHistograms = true
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
	r.numSamples++
}

func (r *SlidingRateBuffer) Reset(mint int64, evalt int64) {
	r.currentMint, r.evalTs = mint, evalt

	if len(r.samples) == 0 || r.samples[len(r.samples)-1].T <= mint {
		r.samples = r.samples[:0]
		r.counterCorrection = 0
		r.numSamples = 0
		r.sampleCount = 0
		r.histogramResets = r.histogramResets[:0]
		return
	}

	// Find first sample to keep and adjust counter correction
	var drop int
	for drop = 0; drop < len(r.samples) && r.samples[drop].T <= mint; drop++ {
		s := &r.samples[drop]
		if s.V.H != nil {
			r.sampleCount -= telemetry.CalculateHistogramSampleCount(s.V.H)
		} else {
			r.sampleCount--
		}
		r.numSamples--

		// If we're dropping a sample that was part of a counter reset,
		// we need to remove its contribution from counterCorrection.
		// A reset contributes the pre-reset value to counterCorrection.
		// We detect this by checking if the next sample (if exists and in window)
		// has a lower value than this one.
		if r.isCounter && !r.hasHistograms && drop+1 < len(r.samples) {
			next := &r.samples[drop+1]
			if next.T > mint && s.V.F > next.V.F {
				// This sample's value was added to counterCorrection during a reset
				r.counterCorrection -= s.V.F
			}
		}
	}

	if drop == 0 {
		return
	}

	// Drop histogram resets that are now outside the window
	if r.hasHistograms {
		dropResets := 0
		for ; dropResets < len(r.histogramResets) && r.histogramResets[dropResets].T <= mint; dropResets++ {
		}
		if dropResets > 0 {
			keep := len(r.histogramResets) - dropResets
			r.histogramResetsTail = resize(r.histogramResetsTail, dropResets)
			copy(r.histogramResetsTail, r.histogramResets[:dropResets])
			copy(r.histogramResets, r.histogramResets[dropResets:])
			copy(r.histogramResets[keep:], r.histogramResetsTail)
			r.histogramResets = r.histogramResets[:keep]
		}
	}

	// Shift samples
	keep := len(r.samples) - drop
	r.tail = resize(r.tail, drop)
	copy(r.tail, r.samples[:drop])
	copy(r.samples, r.samples[drop:])
	copy(r.samples[keep:], r.tail)
	r.samples = r.samples[:keep]
}

func (r *SlidingRateBuffer) Eval(ctx context.Context, _, _ float64, _ int64) (float64, *histogram.FloatHistogram, bool, error) {
	if len(r.samples) < 2 {
		return 0, nil, false, nil
	}

	first := &r.samples[0]
	last := &r.samples[len(r.samples)-1]

	// Check for mixed float/histogram
	if r.hasHistograms {
		var fd, hd bool
		for i := range r.samples {
			hd = hd || r.samples[i].V.H != nil
			fd = fd || r.samples[i].V.H == nil
		}
		if fd && hd {
			warnings.AddToContext(annotations.NewMixedFloatsHistogramsWarning("", posrange.PositionRange{}), ctx)
			return 0, nil, false, nil
		}
	}

	var (
		rangeStart      = r.evalTs - (r.selectRange + r.offset)
		rangeEnd        = r.evalTs - r.offset
		resultValue     float64
		resultHistogram *histogram.FloatHistogram
	)

	if first.V.H != nil {
		// Build samples slice for histogram rate calculation
		// We need all samples plus reset info
		samples := make([]Sample, 0, len(r.samples)+len(r.histogramResets))
		samples = append(samples, r.samples...)
		// histogramResets are stored as pairs (pre-reset, post-reset)
		// We need to insert them into the samples for proper rate calculation
		// For now, fall back to using histogramRate with all samples
		var err error
		resultHistogram, err = histogramRate(ctx, r.samples, r.isCounter)
		if err != nil {
			return 0, nil, false, err
		}
	} else {
		// Float rate calculation with pre-computed counter correction
		resultValue = last.V.F - first.V.F
		if r.isCounter {
			resultValue += r.counterCorrection
		}
	}

	// Extrapolation calculation (same as extrapolatedRate)
	durationToStart := float64(first.T-rangeStart) / 1000
	durationToEnd := float64(rangeEnd-last.T) / 1000
	sampledInterval := float64(last.T-first.T) / 1000
	if sampledInterval == 0 {
		return 0, nil, false, nil
	}
	averageDurationBetweenSamples := sampledInterval / float64(r.numSamples-1)

	extrapolationThreshold := averageDurationBetweenSamples * 1.1

	if durationToStart >= extrapolationThreshold {
		durationToStart = averageDurationBetweenSamples / 2
	}
	if r.isCounter {
		durationToZero := durationToStart
		if resultValue > 0 && first.V.F >= 0 {
			durationToZero = sampledInterval * (first.V.F / resultValue)
		} else if resultHistogram != nil && resultHistogram.Count > 0 && first.V.H != nil && first.V.H.Count >= 0 {
			durationToZero = sampledInterval * (first.V.H.Count / resultHistogram.Count)
		}
		if durationToZero < durationToStart {
			durationToStart = durationToZero
		}
	}

	if durationToEnd >= extrapolationThreshold {
		durationToEnd = averageDurationBetweenSamples / 2
	}

	factor := (sampledInterval + durationToStart + durationToEnd) / sampledInterval
	if r.isRate {
		factor /= float64(r.selectRange) / 1000
	}
	if resultHistogram == nil {
		resultValue *= factor
	} else {
		resultHistogram.Mul(factor)
	}

	if first.V.H != nil && resultHistogram == nil {
		return 0, nil, false, nil
	}

	return resultValue, resultHistogram, true, nil
}
