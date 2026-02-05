// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package query

import (
	"fmt"
	"sync/atomic"
)

type SampleTracker struct {
	current atomic.Int64
	limit   int64
}

func NewSampleTracker(maxSamples int) *SampleTracker {
	return &SampleTracker{
		limit: int64(maxSamples),
	}
}

func (st *SampleTracker) Add(count int) {
	st.current.Add(int64(count))
}

func (st *SampleTracker) Remove(count int) {
	st.current.Add(-int64(count))
}

func (st *SampleTracker) CheckLimit() error {
	if st.limit <= 0 {
		return nil
	}
	current := st.current.Load()
	if current > st.limit {
		return ErrMaxSamplesExceeded{Current: current, Limit: st.limit}
	}
	return nil
}

type ErrMaxSamplesExceeded struct {
	Current int64
	Limit   int64
}

func (e ErrMaxSamplesExceeded) Error() string {
	return fmt.Sprintf("query processing would load too many samples into memory: current=%d, limit=%d", e.Current, e.Limit)
}
