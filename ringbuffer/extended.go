// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package ringbuffer

import (
	"context"
	"math"

	"github.com/thanos-io/promql-engine/warnings"

	"github.com/prometheus/prometheus/model/histogram"
)

// ExtendedRingBuffer retains the newest samples at or before the range start as
// baselines for xrate, xincrease, xdelta and the anchored and smoothed
// modifiers. Samples are normally offered to Push in timestamp order, but
// baseline insertion also preserves ordering when a prefetched sample arrives
// after an in-window sample.
type ExtendedRingBuffer struct {
	*GenericRingBuffer

	extLookback      int64
	baselines        int
	metricAppearedTs int64
	// lastHistogramTs is the newest native histogram offered to Push, even if
	// it was not retained. Anchored and smoothed selectors reject the range when
	// it falls within the extended window, as Prometheus does.
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
		baselines:         1,
		metricAppearedTs:  math.MinInt64,
		lastHistogramTs:   math.MinInt64,
	}
}

// NewAnchored creates a buffer for a range selector with the anchored modifier.
// lookback is the maximum age in milliseconds of the sample that anchors the
// range start.
func NewAnchored(ctx context.Context, size int, selectRange, offset, lookback int64, call FunctionCall) *ExtendedRingBuffer {
	b := NewWithExtLookback(ctx, size, selectRange, offset, lookback, call)
	b.anchored = true
	// When no sample falls inside the range, Prometheus' pickFirstSampleIndex
	// selects the second newest sample at or before the range start, so changes
	// and resets need two baselines to compare.
	b.baselines = 2
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
// evaluation step. It keeps the suffix after mint and up to r.baselines of the
// newest samples at or before mint that are within extLookback.
func (r *ExtendedRingBuffer) Reset(mint int64, evalt int64) {
	r.currentStep = evalt
	r.currentRangeStart = mint

	var drop int
	for drop = 0; drop < len(r.items) && r.items[drop].T <= mint; drop++ {
	}
	for kept := 0; kept < r.baselines && drop > 0 && r.items[drop-1].T >= mint-r.extLookback; kept++ {
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
	if v.H != nil && t > r.lastHistogramTs {
		r.lastHistogramTs = t
	}

	if t > r.currentRangeStart {
		r.GenericRingBuffer.push(t, v)
		return
	}
	if r.currentRangeStart-t > r.extLookback {
		return
	}

	// Reset leaves at most r.baselines samples before the in-window suffix.
	// pos is where t belongs among them and n is how many there are.
	pos, n := 0, 0
	for ; n < len(r.items) && r.items[n].T <= r.currentRangeStart; n++ {
		if r.items[n].T < t {
			pos = n + 1
		}
	}
	if pos < n && r.items[pos].T == t {
		r.sampleCount -= valueSampleCount(r.items[pos].V)
		setSample(&r.items[pos], t, v)
		r.sampleCount += valueSampleCount(v)
		return
	}
	if pos == 0 && n >= r.baselines {
		// Older than every retained baseline.
		return
	}

	// Insert in timestamp order. Inserting before the end is only needed if a
	// prefetched baseline is supplied after a newer sample; it preserves the
	// ordering expected by MaxT and the range functions.
	r.items = append(r.items, Sample{})
	copy(r.items[pos+1:], r.items[pos:len(r.items)-1])
	r.items[pos] = Sample{}
	setSample(&r.items[pos], t, v)
	r.sampleCount += valueSampleCount(v)
	if n+1 > r.baselines {
		r.drop(1)
	}
}

func (r *ExtendedRingBuffer) Eval(ctx context.Context, scalarArg float64, scalarArg2 float64) (float64, *histogram.FloatHistogram, bool, warnings.Warnings, error) {
	if (r.anchored || r.smoothed) && r.lastHistogramTs >= r.currentRangeStart-r.extLookback {
		return 0, nil, false, 0, ErrExtendedRangeHistograms
	}
	return r.eval(scalarArg, scalarArg2, r.metricAppearedTs)
}
