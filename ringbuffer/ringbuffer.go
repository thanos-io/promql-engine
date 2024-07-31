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

	extLookback int64
}

func New(size int) *RingBuffer {
	return NewWithExtLookback(size, 0)
}

func NewWithExtLookback(size int, lookback int64) *RingBuffer {
	return &RingBuffer{
		items:       make([]Sample, 0, size),
		extLookback: lookback,
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

func (r *RingBuffer) ReadIntoLast(f func(*Sample)) {
	f(&r.items[len(r.items)-1])
}

func (r *RingBuffer) Push(t int64, v Value) {
	if n := len(r.items); n < cap(r.items) {
		r.items = r.items[:n+1]
		r.items[n].T = t
		r.items[n].V = v
	} else {
		r.items = append(r.items, Sample{T: t, V: v})
	}
}

func (r *RingBuffer) DropBefore(ts int64) {
	if len(r.items) == 0 || r.items[len(r.items)-1].T < ts {
		r.items = r.items[:0]
		return
	}
	var drop int
	for drop = 0; drop < len(r.items) && r.items[drop].T < ts; drop++ {
	}
	if r.extLookback > 0 && drop > 0 && r.items[drop-1].T >= ts-r.extLookback {
		drop--
	}

	keep := len(r.items) - drop

	r.tail = resize(r.tail, drop)
	copy(r.tail, r.items[:drop])
	copy(r.items, r.items[drop:])
	copy(r.items[keep:], r.tail)
	r.items = r.items[:keep]
}

func (r *RingBuffer) Samples() []Sample {
	return r.items
}

func resize(s []Sample, n int) []Sample {
	if cap(s) >= n {
		return s[:n]
	}
	return make([]Sample, n)
}
