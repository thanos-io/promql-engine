package executionplan

import (
	"time"

	"github.com/prometheus/prometheus/promql"
)

type FunctionCall func(points []promql.Point, stepTime time.Time) *promql.Point

func NewRate(selectRange time.Duration) FunctionCall {
	return func(points []promql.Point, stepTime time.Time) *promql.Point {
		return extrapolatedRate(points, true, true, stepTime, selectRange)
	}
}

// extrapolatedRate is a utility function for rate/increase/delta.
// It calculates the rate (allowing for counter resets if isCounter is true),
// extrapolates if the first/last sample is close to the boundary, and returns
// the result as either per-second (if isRate is true) or overall.
func extrapolatedRate(samples []promql.Point, isCounter, isRate bool, stepTime time.Time, selectRange time.Duration) *promql.Point {
	var (
		rangeStart = stepTime.UnixMilli() - selectRange.Milliseconds()
		rangeEnd   = stepTime.UnixMilli()
	)

	if len(samples) < 2 {
		return nil
	}

	resultValue := samples[len(samples)-1].V - samples[0].V
	if isCounter {
		var lastValue float64
		for _, sample := range samples {
			if sample.V < lastValue {
				resultValue += lastValue
			}
			lastValue = sample.V
		}
	}

	// Duration between first/last samples and boundary of range.
	durationToStart := float64(samples[0].T-rangeStart) / 1000
	durationToEnd := float64(rangeEnd-samples[len(samples)-1].T) / 1000

	sampledInterval := float64(samples[len(samples)-1].T-samples[0].T) / 1000
	averageDurationBetweenSamples := sampledInterval / float64(len(samples)-1)

	if isCounter && resultValue > 0 && samples[0].V >= 0 {
		// Counters cannot be negative. If we have any slope at
		// all (i.e. resultValue went up), we can extrapolate
		// the zero point of the counter. If the duration to the
		// zero point is shorter than the durationToStart, we
		// take the zero point as the start of the series,
		// thereby avoiding extrapolation to negative counter
		// values.
		durationToZero := sampledInterval * (samples[0].V / resultValue)
		if durationToZero < durationToStart {
			durationToStart = durationToZero
		}
	}

	// If the first/last samples are close to the boundaries of the range,
	// extrapolate the result. This is as we expect that another sample
	// will exist given the spacing between samples we've seen thus far,
	// with an allowance for noise.
	extrapolationThreshold := averageDurationBetweenSamples * 1.1
	extrapolateToInterval := sampledInterval

	if durationToStart < extrapolationThreshold {
		extrapolateToInterval += durationToStart
	} else {
		extrapolateToInterval += averageDurationBetweenSamples / 2
	}
	if durationToEnd < extrapolationThreshold {
		extrapolateToInterval += durationToEnd
	} else {
		extrapolateToInterval += averageDurationBetweenSamples / 2
	}
	resultValue = resultValue * (extrapolateToInterval / sampledInterval)
	if isRate {
		resultValue = resultValue / selectRange.Seconds()
	}

	return &promql.Point{
		T: stepTime.UnixMilli(),
		V: resultValue,
	}
}
