// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package ringbuffer

import (
	"testing"

	"github.com/efficientgo/core/testutil"
)

func TestRingBuffer(t *testing.T) {
	floats := newReader([]Sample{{30, Value{F: 1}}, {60, Value{F: 2}}, {90, Value{F: 3}}})
	buffer := New(4)

	buffer.ReadIntoNext(floats.ReadNext)
	testutil.Equals(t, []Sample{{30, Value{F: 1}}}, buffer.Samples())

	buffer.ReadIntoNext(floats.ReadNext)
	testutil.Equals(t, []Sample{{30, Value{F: 1}}, {60, Value{F: 2}}}, buffer.Samples())

	buffer.DropBefore(60)
	testutil.Equals(t, []Sample{{60, Value{F: 2}}}, buffer.Samples())

	buffer.ReadIntoNext(floats.ReadNext)
	testutil.Equals(t, []Sample{{60, Value{F: 2}}, {90, Value{F: 3}}}, buffer.Samples())
}

type reader struct {
	i     int
	items []Sample
}

func newReader(items []Sample) *reader {
	return &reader{
		items: items,
	}
}

func (f *reader) ReadNext(item *Sample) bool {
	item.T = f.items[f.i].T
	item.V = f.items[f.i].V
	f.i++
	return true
}
