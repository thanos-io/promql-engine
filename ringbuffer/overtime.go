// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package ringbuffer

import (
	"context"
	"math"

	"github.com/thanos-io/promql-engine/execution/telemetry"
	"github.com/thanos-io/promql-engine/query"

	"github.com/prometheus/prometheus/model/histogram"
)

// TODO: this is mostly copied over from the rate aligner and not at all special to
// count_over_time - we could of course generalize it to other over_time functions and we
// will do that soon.
// This is really just to get warmed up.
// The idea here is that we tile the range into "stepRanges" and when we push a sample
// and it happens to fall into a step range, we increment the state that we hold in a
// corresponding "resultSamples" slice.

// CountOverTimeBuffer is a Buffer which can calculate count_over_time for a series in a
// streaming manner, calculating the value incrementally for each step where the sample is used.
type CountOverTimeBuffer struct {
	// stepRanges contain the bounds and number of samples for each evaluation step.
	stepRanges []stepRange

	// resultscounts contains the resulting state for each evaluation step.
	resultCounts []float64

	// firstTimestamps contains the timestamp of the first sample for each evaluation step.
	firstTimestamps []int64

	// lastTimestamp is the timestamp of the lsat sample in the current evaluation step
	lastTimestamp int64

	step int64
}

// NewCountOverTime creates a new CountOverTimeBuffer.
func NewCountOverTimeBuffer(opts query.Options, selectRange, offset int64) *CountOverTimeBuffer {
	var (
		step     = max(1, opts.Step.Milliseconds())
		numSteps = min(
			(selectRange-1)/step+1,
			querySteps(opts),
		)

		current         = opts.Start.UnixMilli()
		resultCounts    = make([]float64, 0, numSteps)
		firstTimestamps = make([]int64, 0, numSteps)
		stepRanges      = make([]stepRange, 0, numSteps)
	)
	for range int(numSteps) {
		var (
			maxt = current - offset
			mint = maxt - selectRange
		)
		stepRanges = append(stepRanges, stepRange{mint: mint, maxt: maxt})
		resultCounts = append(resultCounts, 0.)
		firstTimestamps = append(firstTimestamps, math.MaxInt64)

		current += step
	}

	return &CountOverTimeBuffer{
		step:            step,
		stepRanges:      stepRanges,
		resultCounts:    resultCounts,
		firstTimestamps: firstTimestamps,
		lastTimestamp:   math.MinInt64,
	}
}

func (r *CountOverTimeBuffer) SampleCount() int {
	return r.stepRanges[0].sampleCount
}

func (r *CountOverTimeBuffer) MaxT() int64 { return r.lastTimestamp }

func (r *CountOverTimeBuffer) Push(t int64, v Value) {
	// Set the lastSample sample for the current evaluation step.
	r.lastTimestamp = t

	// Set the first sample for each evaluation step where the currently read sample is used.
	for i := 0; i < len(r.stepRanges) && t > r.stepRanges[i].mint && t <= r.stepRanges[i].maxt; i++ {
		r.stepRanges[i].numSamples++
		if v.H != nil {
			r.stepRanges[i].sampleCount += telemetry.CalculateHistogramSampleCount(v.H)
		} else {
			r.stepRanges[i].sampleCount++
		}

		// Count the current sample in its range
		r.resultCounts[i] += 1

		if fts := r.firstTimestamps[i]; t >= fts {
			continue
		}
		r.firstTimestamps[i] = t
	}
}

func (r *CountOverTimeBuffer) Reset(mint int64, evalt int64) {
	if r.stepRanges[0].mint == mint {
		return
	}

	lastSample := len(r.stepRanges) - 1
	var (
		nextMint = r.stepRanges[lastSample].mint + r.step
		nextMaxt = r.stepRanges[lastSample].maxt + r.step
	)

	nextStepRange := r.stepRanges[0]
	copy(r.stepRanges, r.stepRanges[1:])
	r.stepRanges[lastSample] = nextStepRange
	r.stepRanges[lastSample].mint = nextMint
	r.stepRanges[lastSample].maxt = nextMaxt
	r.stepRanges[lastSample].sampleCount = 0
	r.stepRanges[lastSample].numSamples = 0

	copy(r.resultCounts, r.resultCounts[1:])
	r.resultCounts[lastSample] = 0

	nextFirstTimestamp := r.firstTimestamps[0]
	copy(r.firstTimestamps, r.firstTimestamps[1:])
	r.firstTimestamps[lastSample] = nextFirstTimestamp
	r.firstTimestamps[lastSample] = math.MaxInt64
}

func (r *CountOverTimeBuffer) Eval(ctx context.Context, _, _ float64, _ int64) (float64, *histogram.FloatHistogram, bool, error) {
	if r.firstTimestamps[0] == math.MaxInt64 {
		return 0, nil, false, nil
	}

	return r.resultCounts[0], nil, true, nil
}
