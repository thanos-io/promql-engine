// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package ringbuffer

import (
	"context"
	"math"

	"github.com/thanos-io/promql-engine/execution/telemetry"
	"github.com/thanos-io/promql-engine/query"
	"github.com/thanos-io/promql-engine/warnings"

	"github.com/prometheus/prometheus/model/histogram"
)

// RateBuffer is a Buffer which can calculate rate, increase and delta for a
// series in a streaming manner, calculating the value incrementally for each
// step where the sample is used.
type RateBuffer struct {
	ctx context.Context
	// stepRanges contain the bounds, sample counts, and first sample for each
	// evaluation step.
	stepRanges []rateStepRange
	// lastSample is the lastSample sample in the current evaluation step.
	lastSample Sample

	currentMint int64
	selectRange int64
	step        int64
	offset      int64
	isCounter   bool
	isRate      bool

	evalTs int64
}

type stepRange struct {
	mint        int64
	maxt        int64
	numSamples  int
	sampleCount int
}

type rateStepRange struct {
	stepRange
	firstSample       Sample
	counterCorrection float64
	histogramState    *histogramRateState
	mixedSamples      bool
}

// NewRateBuffer creates a new RateBuffer.
func NewRateBuffer(ctx context.Context, opts query.Options, isCounter, isRate bool, selectRange, offset int64) *RateBuffer {
	var (
		step     = max(1, opts.Step.Milliseconds())
		numSteps = min(
			(selectRange-1)/step+1,
			querySteps(opts),
		)

		current    = opts.Start.UnixMilli()
		stepRanges = make([]rateStepRange, 0, numSteps)
	)
	for range int(numSteps) {
		var (
			maxt = current - offset
			mint = maxt - selectRange
		)
		stepRanges = append(stepRanges, rateStepRange{
			stepRange:   stepRange{mint: mint, maxt: maxt},
			firstSample: Sample{T: math.MaxInt64},
		})
		current += step
	}

	return &RateBuffer{
		ctx:         ctx,
		isCounter:   isCounter,
		isRate:      isRate,
		selectRange: selectRange,
		step:        step,
		offset:      offset,
		stepRanges:  stepRanges,
		lastSample:  Sample{T: math.MinInt64},
		currentMint: math.MaxInt64,
	}
}

func (r *RateBuffer) SampleCount() int {
	return r.stepRanges[0].sampleCount
}

func (r *RateBuffer) MaxT() int64 { return r.lastSample.T }

func (r *RateBuffer) Push(t int64, v Value) {
	if t <= r.currentMint {
		return
	}
	previousSample := r.lastSample
	floatReset := r.isCounter &&
		previousSample.T > r.currentMint &&
		previousSample.V.H == nil && v.H == nil &&
		previousSample.V.F > v.F

	histogramReset := r.isCounter &&
		previousSample.T > r.currentMint &&
		previousSample.V.H != nil && v.H != nil &&
		v.H.DetectReset(previousSample.V.H)
	sampleCount := 1
	if v.H != nil {
		sampleCount = telemetry.CalculateHistogramSampleCount(v.H)
	}

	// Accumulate each window's correction before overwriting the previous
	// histogram, whose storage is reused by lastSample.
	for i := 0; i < len(r.stepRanges) && t > r.stepRanges[i].mint && t <= r.stepRanges[i].maxt; i++ {
		step := &r.stepRanges[i]
		if step.numSamples == 0 {
			setSample(&step.firstSample, t, v)
			if r.isCounter && v.H != nil {
				if step.histogramState == nil {
					step.histogramState = &histogramRateState{}
				}
				step.histogramState.reset(v.H)
			}
		} else {
			step.mixedSamples = step.mixedSamples || (step.firstSample.V.H == nil) != (v.H == nil)
			if r.isCounter && v.H != nil && !step.mixedSamples {
				step.histogramState.push(previousSample.V.H, v.H, histogramReset && previousSample.T > step.mint, step.numSamples == 1)
			}
		}
		if floatReset && previousSample.T > step.mint {
			step.counterCorrection += previousSample.V.F
		}
		step.numSamples++
		step.sampleCount += sampleCount
	}
	setSample(&r.lastSample, t, v)
}

func (r *RateBuffer) Reset(mint int64, evalt int64) {
	r.currentMint, r.evalTs = mint, evalt
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
	r.stepRanges[lastSample].counterCorrection = 0
	r.stepRanges[lastSample].mixedSamples = false

	r.stepRanges[lastSample].firstSample.T = math.MaxInt64
}

func (r *RateBuffer) Eval(ctx context.Context, _, _ float64) (float64, *histogram.FloatHistogram, bool, warnings.Warnings, error) {
	step := &r.stepRanges[0]
	firstSample := step.firstSample
	if firstSample.T == math.MaxInt64 || firstSample.T == r.lastSample.T {
		return 0, nil, false, 0, nil
	}

	if step.mixedSamples {
		return 0, nil, false, warnings.WarnMixedFloatsHistograms, nil
	}
	var (
		f    float64
		h    *histogram.FloatHistogram
		warn warnings.Warnings
		err  error
	)
	if firstSample.V.H != nil {
		if r.isCounter {
			h, warn, err = step.histogramState.eval(firstSample.V.H, r.lastSample.V.H)
		} else {
			h, warn, err = histogramRate([]Sample{firstSample, r.lastSample}, false)
		}
		if err != nil || h == nil {
			return 0, nil, false, warn, err
		}
	} else {
		f = r.lastSample.V.F - firstSample.V.F + step.counterCorrection
	}
	f, h = extrapolateRate(firstSample, r.lastSample, step.numSamples, f, h, r.isCounter, r.isRate, r.evalTs, r.selectRange, r.offset)
	return f, h, true, warn, nil
}

func querySteps(o query.Options) int64 {
	// Instant evaluation is executed as a range evaluation with one step.
	if o.Step.Milliseconds() == 0 {
		return 1
	}

	return (o.End.UnixMilli()-o.Start.UnixMilli())/o.Step.Milliseconds() + 1
}
