// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package ringbuffer

import (
	"context"
	"math"

	"github.com/thanos-io/promql-engine/warnings"

	"github.com/prometheus/prometheus/model/histogram"
)

// ExtendedRingBuffer retains the newest sample at or before the range start as
// a baseline for xrate, xincrease, xdelta and the anchored and smoothed
// modifiers. Samples are normally offered to Push in timestamp order, but
// baseline insertion also preserves ordering when a prefetched sample arrives
// after an in-window sample.
type ExtendedRingBuffer struct {
	*GenericRingBuffer

	extLookback      int64
	metricAppearedTs int64
	anchored         bool
	smoothed         bool

	// lastFloatTs and lastHistogramTs are the newest float and native
	// histogram offered to Push, even if they were not retained as baselines.
	// Prometheus reads the whole extended window for anchored and smoothed
	// selectors, so a window holding both types must be detected even when
	// only the newest baseline is kept.
	lastFloatTs     int64
	lastHistogramTs int64
}

// NewWithExtLookback creates a buffer for an extended range function.
// extLookback is the maximum age in milliseconds of a baseline relative to the
// range start.
func NewWithExtLookback(
	ctx context.Context,
	size int,
	selectRange, offset, extLookback int64,
	call FunctionCall,
) *ExtendedRingBuffer {
	return &ExtendedRingBuffer{
		GenericRingBuffer: New(ctx, size, selectRange, offset, call),
		extLookback:       extLookback,
		metricAppearedTs:  math.MinInt64,
		lastFloatTs:       math.MinInt64,
		lastHistogramTs:   math.MinInt64,
	}
}

// NewAnchored creates a buffer for a range selector with the anchored modifier.
// lookback is the maximum age in milliseconds of the sample that anchors the
// range start.
func NewAnchored(ctx context.Context, size int, selectRange, offset, lookback int64, call FunctionCall) *ExtendedRingBuffer {
	b := NewWithExtLookback(ctx, size, selectRange, offset, lookback, call)
	b.anchored = true
	return b
}

// NewSmoothed creates a buffer for a range selector with the smoothed modifier.
// lookback is the maximum age in milliseconds of the sample used to
// interpolate the range start.
func NewSmoothed(ctx context.Context, size int, selectRange, offset, lookback int64, call FunctionCall) *ExtendedRingBuffer {
	b := NewWithExtLookback(ctx, size, selectRange, offset, lookback, call)
	b.smoothed = true
	return b
}

// Reset applies the extended-window rule to samples retained from the previous
// evaluation step. It keeps the suffix after mint and, when one exists within
// extLookback, the newest baseline at or before mint.
func (r *ExtendedRingBuffer) Reset(mint int64, evalt int64) {
	r.currentStep = evalt
	r.currentRangeStart = mint

	var drop int
	for drop = 0; drop < len(r.items) && r.items[drop].T <= mint; drop++ {
	}
	if drop > 0 && r.items[drop-1].T >= mint-r.extLookback {
		drop--
	}
	r.drop(drop)
}

// Push applies the same extended-window rule to samples supplied after Reset.
// It also records the earliest candidate observed for xincrease's initial-zero
// injection. This covers both the scanner's prefetched sample and samples read
// directly from its iterator.
func (r *ExtendedRingBuffer) Push(t int64, v Value) {
	if r.metricAppearedTs == math.MinInt64 || t < r.metricAppearedTs {
		r.metricAppearedTs = t
	}
	if v.H != nil {
		r.lastHistogramTs = max(r.lastHistogramTs, t)
	} else {
		r.lastFloatTs = max(r.lastFloatTs, t)
	}

	if t > r.currentRangeStart {
		r.GenericRingBuffer.push(t, v)
		return
	}
	if r.currentRangeStart-t > r.extLookback {
		return
	}

	// Reset leaves at most one baseline before the in-window suffix. Find and
	// replace it if one exists.
	baseline := -1
	for i := len(r.items) - 1; i >= 0; i-- {
		if r.items[i].T <= r.currentRangeStart {
			baseline = i
			break
		}
	}
	if baseline >= 0 {
		if t >= r.items[baseline].T {
			r.sampleCount -= valueSampleCount(r.items[baseline].V)
			setSample(&r.items[baseline], t, v)
			r.sampleCount += valueSampleCount(v)
		}
		return
	}

	// This path is only needed if a prefetched baseline is supplied after an
	// in-window sample. Insert it before the in-window suffix to preserve the
	// timestamp ordering expected by MaxT and the range functions.
	r.items = append(r.items, Sample{})
	copy(r.items[1:], r.items[:len(r.items)-1])
	r.items[0] = Sample{}
	setSample(&r.items[0], t, v)
	r.sampleCount += valueSampleCount(v)
}

func (r *ExtendedRingBuffer) Eval(ctx context.Context, scalarArg float64, scalarArg2 float64) (float64, *histogram.FloatHistogram, bool, warnings.Warnings, error) {
	var mixed bool
	if r.anchored || r.smoothed {
		// Samples are only pushed up to the current step's maxt, so the newest
		// sample of each type lies in the extended window if it is not older
		// than the oldest baseline the window admits.
		windowStart := r.currentRangeStart - r.extLookback
		mixed = r.lastFloatTs >= windowStart && r.lastHistogramTs >= windowStart
	}
	args := r.functionArgs(scalarArg, scalarArg2)
	args.MetricAppearedTs = r.metricAppearedTs
	args.Anchored = r.anchored
	args.Smoothed = r.smoothed
	args.MixedWindow = mixed
	return r.call(args)
}
