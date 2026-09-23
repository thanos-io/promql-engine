// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package query

import (
	"fmt"
	"math"
	"sync/atomic"
)

type SampleTracker interface {
	Add(count int)
	Remove(count int)
	CheckLimit() error
	Limit() int64
}

// SampleLimitOvershoot is the fraction of maxSamples allowed as overshoot before enforcement.
const SampleLimitOvershoot = 0.05

type sampleTracker struct {
	current atomic.Int64
	limit   int64
}

func NewSampleTracker(maxSamples int) SampleTracker {
	if maxSamples <= 0 {
		return nopSampleTracker{}
	}
	return &sampleTracker{
		limit: int64(maxSamples),
	}
}

func (st *sampleTracker) Add(count int) {
	st.current.Add(int64(count))
}

func (st *sampleTracker) Remove(count int) {
	st.current.Add(-int64(count))
}

func (st *sampleTracker) CheckLimit() error {
	current := st.current.Load()
	if current > st.limit {
		return ErrMaxSamplesExceeded{Current: current, Limit: st.limit}
	}
	return nil
}

func (st *sampleTracker) Limit() int64 {
	return st.limit
}

type nopSampleTracker struct{}

func (nopSampleTracker) Add(int)           {}
func (nopSampleTracker) Remove(int)        {}
func (nopSampleTracker) CheckLimit() error { return nil }
func (nopSampleTracker) Limit() int64      { return math.MaxInt64 }

type ErrMaxSamplesExceeded struct {
	Current int64
	Limit   int64
}

func (e ErrMaxSamplesExceeded) Error() string {
	return fmt.Sprintf("query processing would load too many samples into memory: current=%d, limit=%d", e.Current, e.Limit)
}

// ComputeSampleLimitCheckThreshold returns the per-shard sample count before
// checking the global limit: floor(limit * SampleLimitOvershoot / decodingConcurrency).
func ComputeSampleLimitCheckThreshold(opts *Options) int {
	limit := opts.SampleTracker.Limit()
	if limit <= 0 || limit == math.MaxInt64 {
		return math.MaxInt64
	}
	concurrency := opts.DecodingConcurrency
	if concurrency <= 0 {
		concurrency = 1
	}
	threshold := int(math.Floor(float64(limit) * SampleLimitOvershoot / float64(concurrency)))
	if threshold <= 0 {
		threshold = 1
	}
	return threshold
}