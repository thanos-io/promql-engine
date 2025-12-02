// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package query

import (
	"sync/atomic"

	"github.com/efficientgo/core/errors"
)

// SampleTracker tracks the number of samples currently in memory during query execution.
// It enforces a maximum sample limit to prevent out-of-memory errors.
type SampleTracker struct {
	current atomic.Int64 // Current samples in memory
	limit   int64        // Maximum samples allowed (0 = no limit)
}

// NewSampleTracker creates a new sample tracker with the given limit.
func NewSampleTracker(maxSamples int) *SampleTracker {
	return &SampleTracker{
		limit: int64(maxSamples),
	}
}

// Add increments the current sample count and checks against the limit.
// Returns an error if adding these samples would exceed the limit.
func (st *SampleTracker) Add(count int) error {
	if count <= 0 {
		return nil
	}

	newCurrent := st.current.Add(int64(count))

	// Check limit
	if st.limit > 0 && newCurrent > st.limit {
		return errors.Newf("query exceeded maximum samples limit: current=%d, limit=%d", newCurrent, st.limit)
	}

	return nil
}

// Remove decrements the current sample count.
// This should be called when samples are freed from memory.
func (st *SampleTracker) Remove(count int) {
	if count <= 0 {
		return
	}
	st.current.Add(-int64(count))
}

// Current returns the current number of samples in memory.
func (st *SampleTracker) Current() int64 {
	return st.current.Load()
}
