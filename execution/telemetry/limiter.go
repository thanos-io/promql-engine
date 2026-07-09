// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package telemetry

import (
	"fmt"
	"sync/atomic"
)

// SampleLimiter enforces a per-step limit on the number of samples loaded
// into memory. It is shared across all operators in a query so that
// concurrent operators contribute to the same counters.
type SampleLimiter struct {
	maxSamples     int64
	samplesPerStep []atomic.Int64
	startTimestamp int64
	interval       int64
}

// NewSampleLimiter creates a limiter that tracks samples per evaluation step.
// When maxSamples is zero or negative, the limiter is a no-op.
func NewSampleLimiter(maxSamples int, start, end, interval int64) *SampleLimiter {
	if maxSamples <= 0 || interval <= 0 {
		return &SampleLimiter{}
	}
	numSteps := int((end-start)/interval) + 1
	return &SampleLimiter{
		maxSamples:     int64(maxSamples),
		samplesPerStep: make([]atomic.Int64, numSteps),
		startTimestamp: start,
		interval:       interval,
	}
}

// MaxSamples returns the configured sample limit.
func (sl *SampleLimiter) MaxSamples() int { return int(sl.maxSamples) }

// Add records samples loaded at timestamp t and returns an error if the
// per-step limit is exceeded. When no limit is configured, Add is a no-op.
func (sl *SampleLimiter) Add(samples int, t int64) error {
	if sl.maxSamples <= 0 {
		return nil
	}
	i := int((t - sl.startTimestamp) / sl.interval)
	if i < 0 || i >= len(sl.samplesPerStep) {
		return nil
	}
	current := sl.samplesPerStep[i].Add(int64(samples))
	if current > sl.maxSamples {
		return ErrMaxSamplesExceeded{Current: current, Limit: sl.maxSamples}
	}
	return nil
}

// ErrMaxSamplesExceeded is returned when a query would load more than the
// configured maximum number of samples into memory at a single evaluation step.
type ErrMaxSamplesExceeded struct {
	Current int64
	Limit   int64
}

func (e ErrMaxSamplesExceeded) Error() string {
	return fmt.Sprintf("query processing would load too many samples into memory: current=%d, limit=%d", e.Current, e.Limit)
}
