// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package ringbuffer

import (
	"testing"

	"github.com/efficientgo/core/testutil"
)

func TestRingBuffer(t *testing.T) {
	floats := newFloatReader([]Sample[float64]{{30, 1}, {60, 2}, {90, 3}})
	buffer := New[float64](4)

	buffer.ReadIntoNext(floats.ReadNext)
	testutil.Equals(t, []Sample[float64]{{30, 1}}, buffer.Samples())

	buffer.ReadIntoNext(floats.ReadNext)
	testutil.Equals(t, []Sample[float64]{{30, 1}, {60, 2}}, buffer.Samples())

	buffer.DropBefore(60)
	testutil.Equals(t, []Sample[float64]{{60, 2}}, buffer.Samples())

	buffer.ReadIntoNext(floats.ReadNext)
	testutil.Equals(t, []Sample[float64]{{60, 2}, {90, 3}}, buffer.Samples())
}

type floatReader struct {
	i     int
	items []Sample[float64]
}

func newFloatReader(items []Sample[float64]) *floatReader {
	return &floatReader{
		items: items,
	}
}

func (f *floatReader) ReadNext(item *Sample[float64]) {
	item.T = f.items[f.i].T
	item.V = f.items[f.i].V
	f.i++
}
