// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus

import (
	"io"

	"github.com/prometheus/prometheus/tsdb/chunkenc"
	"github.com/thanos-io/promql-engine/ringbuffer"
)

type BatchIterator struct {
	it chunkenc.Iterator
}

func NewBatchIterator(it chunkenc.Iterator) *BatchIterator {
	return &BatchIterator{it: it}
}

func (b *BatchIterator) NextBatch(buf []ringbuffer.Sample) (int, error) {
	n := 0
	for n < len(buf) {
		valType := b.it.Next()
		if valType == chunkenc.ValNone {
			break
		}

		switch valType {
		case chunkenc.ValFloat:
			buf[n].T, buf[n].V.F = b.it.At()
			buf[n].V.H = nil
		case chunkenc.ValHistogram, chunkenc.ValFloatHistogram:
			buf[n].T, buf[n].V.H = b.it.AtFloatHistogram(buf[n].V.H)
			buf[n].V.F = 0
		}
		n++
	}
	if err := b.it.Err(); err != nil {
		return n, err
	}
	if n < len(buf) {
		return n, io.EOF
	}
	return n, nil
}
