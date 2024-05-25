// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package storage

import (
	"testing"

	"github.com/prometheus/prometheus/model/histogram"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	"github.com/prometheus/prometheus/tsdb/tsdbutil"
	"github.com/stretchr/testify/require"
)

func TestHistogramStatsIterator_Decode(t *testing.T) {
	numHistograms := 20
	chk, err := createHistogramChunk(numHistograms)
	require.NoError(t, err)

	decodedHistograms := make([]*histogram.Histogram, 0)
	decodedCounts := make([]uint64, 0)
	decodedSums := make([]float64, 0)

	it := chk.Iterator(nil)
	for it.Next() != chunkenc.ValNone {
		_, fh := it.AtHistogram(nil)
		decodedHistograms = append(decodedHistograms, fh)
	}

	statsIterator := NewHistogramStatsIterator(chk.Iterator(nil))
	for statsIterator.Next() != chunkenc.ValNone {
		_, h := statsIterator.AtHistogram(nil)
		decodedCounts = append(decodedCounts, h.Count)
		decodedSums = append(decodedSums, h.Sum)
	}
	for i := 0; i < len(decodedHistograms); i++ {
		require.Equal(t, decodedHistograms[i].Count, decodedCounts[i])
		require.Equal(t, decodedHistograms[i].Sum, decodedSums[i])
	}
}

func createHistogramChunk(n int) (*chunkenc.HistogramChunk, error) {
	chunk := chunkenc.NewHistogramChunk()
	appender, err := chunk.Appender()
	if err != nil {
		return nil, err
	}
	hAppender := appender.(*chunkenc.HistogramAppender)

	for i := 0; i < n; i++ {
		if _, _, _, err := hAppender.AppendHistogram(nil, int64(i*1000), tsdbutil.GenerateTestHistogram(i), true); err != nil {
			return nil, err
		}
	}
	return chunk, nil
}
