// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package ringbuffer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtendedRingBufferAdd(t *testing.T) {
	tests := []struct {
		name        string
		extLookback int64
		samples     []Sample
		expected    []Sample
	}{
		{
			name:        "retains newest valid baseline",
			extLookback: 9,
			samples: []Sample{
				{T: 90, V: Value{F: 1}}, // Outside the lookback.
				{T: 91, V: Value{F: 2}},
				{T: 95, V: Value{F: 3}},
				{T: 100, V: Value{F: 4}},
				{T: 101, V: Value{F: 5}},
			},
			expected: []Sample{
				{T: 100, V: Value{F: 4}},
				{T: 101, V: Value{F: 5}},
			},
		},
		{
			name:        "inserts baseline before existing window",
			extLookback: 20,
			samples: []Sample{
				{T: 101, V: Value{F: 3}},
				{T: 95, V: Value{F: 1}},
				{T: 99, V: Value{F: 2}},
			},
			expected: []Sample{
				{T: 99, V: Value{F: 2}},
				{T: 101, V: Value{F: 3}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buffer := NewWithExtLookback(context.Background(), 4, 10, 0, test.extLookback, nil)
			buffer.Reset(100, 110)
			for _, sample := range test.samples {
				buffer.Push(sample.T, sample.V)
			}

			require.Equal(t, test.expected, buffer.items)
		})
	}
}

func TestExtendedRingBufferTracksMetricAppearance(t *testing.T) {
	buffer := NewWithExtLookback(context.Background(), 4, 10, 0, 10, nil)
	buffer.Reset(100, 110)

	// Track every candidate, including one rejected as too old for the current
	// baseline. Metric appearance is series-level state, not window state.
	buffer.Push(80, Value{F: 1})
	buffer.Push(101, Value{F: 2})
	buffer.Reset(200, 210)

	require.Equal(t, int64(80), buffer.metricAppearedTs)
}

func TestExtendedRingBufferRejectsStaleCandidatesAfterReset(t *testing.T) {
	buffer := NewWithExtLookback(context.Background(), 4, 10, 0, 10, nil)
	buffer.Reset(0, 10)
	buffer.Push(0, Value{F: 1})
	buffer.Push(5, Value{F: 2})

	buffer.Reset(100, 110)
	require.Empty(t, buffer.items)

	// These model a prefetched sample and a subsequently read iterator sample.
	// Neither may resurrect a baseline that Reset rejected as too old.
	buffer.Push(10, Value{F: 3})
	buffer.Push(20, Value{F: 4})
	buffer.Push(101, Value{F: 5})

	require.Equal(t, []Sample{{T: 101, V: Value{F: 5}}}, buffer.items)
}
