// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package ringbuffer

import (
	"math"

	"github.com/prometheus/prometheus/model/histogram"
)

type Value struct {
	F float64
	H *histogram.FloatHistogram
}

type Sample struct {
	T int64
	V Value
}

type RingBuffer struct {
	items []Sample
	tail  []Sample

	currentStep int64
	offset      int64
	selectRange int64
	extLookback int64
	call        FunctionCall
}

func New(size int, selectRange, offset int64, call FunctionCall) *RingBuffer {
	return NewWithExtLookback(size, selectRange, offset, 0, call)
}

func NewWithExtLookback(size int, selectRange, offset int64, lookback int64, call FunctionCall) *RingBuffer {
	return &RingBuffer{
		items:       make([]Sample, 0, size),
		selectRange: selectRange,
		offset:      offset,
		extLookback: lookback,
		call:        call,
	}
}

func (r *RingBuffer) Len() int {
	return len(r.items)
}

// MaxT returns the maximum timestamp of the ring buffer.
// If the ring buffer is empty, it returns math.MinInt64.
func (r *RingBuffer) MaxT() int64 {
	if len(r.items) == 0 {
		return math.MinInt64
	}
	return r.items[len(r.items)-1].T
}

// ReadIntoNext can be used to read a sample into the next ring buffer slot through the passed in callback.
// If the callback function returns false, the sample is not kept in the buffer.
func (r *RingBuffer) ReadIntoNext(f func(*Sample) bool) {
	n := len(r.items)
	if cap(r.items) > len(r.items) {
		r.items = r.items[:n+1]
	} else {
		r.items = append(r.items, Sample{})
	}
	if keep := f(&r.items[n]); !keep {
		r.items = r.items[:n]
	}
}

// ReadIntoLast reads a sample into the last slot in the buffer, replacing the existing sample.
func (r *RingBuffer) ReadIntoLast(f func(*Sample)) {
	f(&r.items[len(r.items)-1])
}

// Push adds a new sample to the buffer.
func (r *RingBuffer) Push(t int64, v Value) {
	if n := len(r.items); n < cap(r.items) {
		r.items = r.items[:n+1]
		r.items[n].T = t
		r.items[n].V = v
	} else {
		r.items = append(r.items, Sample{T: t, V: v})
	}
}

func (r *RingBuffer) Reset(mint int64, evalt int64) {
	r.currentStep = evalt

	if len(r.items) == 0 || r.items[len(r.items)-1].T < mint {
		r.items = r.items[:0]
		return
	}
	var drop int
	for drop = 0; drop < len(r.items) && r.items[drop].T < mint; drop++ {
	}
	if r.extLookback > 0 && drop > 0 && r.items[drop-1].T >= mint-r.extLookback {
		drop--
	}

	keep := len(r.items) - drop
	r.tail = resize(r.tail, drop)
	copy(r.tail, r.items[:drop])
	copy(r.items, r.items[drop:])
	copy(r.items[keep:], r.tail)
	r.items = r.items[:keep]
}

func (r *RingBuffer) Eval(scalarArg float64, metricAppearedTs *int64) (float64, *histogram.FloatHistogram, bool, error) {
	return r.call(FunctionArgs{
		Samples:          r.items,
		StepTime:         r.currentStep,
		SelectRange:      r.selectRange,
		Offset:           r.offset,
		ScalarPoint:      scalarArg,
		MetricAppearedTs: metricAppearedTs,
	})
}

func resize(s []Sample, n int) []Sample {
	if cap(s) >= n {
		return s[:n]
	}
	return make([]Sample, n)
}
