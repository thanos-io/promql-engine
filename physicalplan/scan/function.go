// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package scan

import (
	"fmt"
	"math"
	"time"

	"github.com/efficientgo/core/errors"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"

	"github.com/thanos-community/promql-engine/physicalplan/parse"
)

type FunctionCall func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample

func NewFunctionCall(f *parser.Function, selectRange time.Duration) (FunctionCall, error) {
	switch f.Name {
	case "sum_over_time":
		return func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample {
			return promql.Sample{
				Point: promql.Point{
					T: stepTime.UnixMilli(),
					V: sumOverTime(points),
				},
				Metric: labels,
			}
		}, nil
	case "max_over_time":
		return func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample {
			return promql.Sample{
				Point: promql.Point{
					T: stepTime.UnixMilli(),
					V: maxOverTime(points),
				},
				Metric: labels,
			}
		}, nil
	case "min_over_time":
		return func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample {
			return promql.Sample{
				Point: promql.Point{
					T: stepTime.UnixMilli(),
					V: minOverTime(points),
				},
				Metric: labels,
			}
		}, nil
	case "avg_over_time":
		return func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample {
			return promql.Sample{
				Point: promql.Point{
					T: stepTime.UnixMilli(),
					V: avgOverTime(points),
				},
				Metric: labels,
			}
		}, nil
	case "count_over_time":
		return func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample {
			return promql.Sample{
				Point: promql.Point{
					T: stepTime.UnixMilli(),
					V: countOverTime(points),
				},
				Metric: labels,
			}
		}, nil
	case "rate":
		return func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample {
			point := extrapolatedRate(points, true, true, stepTime, selectRange)
			return promql.Sample{
				Point:  point,
				Metric: labels,
			}
		}, nil
	case "delta":
		return func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample {
			point := extrapolatedRate(points, false, false, stepTime, selectRange)
			return promql.Sample{
				Point:  point,
				Metric: labels,
			}
		}, nil
	case "increase":
		return func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample {
			point := extrapolatedRate(points, true, false, stepTime, selectRange)
			return promql.Sample{
				Point:  point,
				Metric: labels,
			}
		}, nil
	case "irate":
		return func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample {
			point := instantValue(points, true, stepTime, selectRange)
			return promql.Sample{
				Point:  point,
				Metric: labels,
			}
		}, nil
	case "idelta":
		return func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample {
			point := instantValue(points, false, stepTime, selectRange)
			return promql.Sample{
				Point:  point,
				Metric: labels,
			}
		}, nil
	case "vector":
		return func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample {
			return promql.Sample{
				Point:  points[0],
				Metric: labels,
			}
		}, nil
	default:
		msg := fmt.Sprintf("unknown function: %s", f.Name)
		return nil, errors.Wrap(parse.ErrNotSupportedExpr, msg)
	}
}

// extrapolatedRate is a utility function for rate/increase/delta.
// It calculates the rate (allowing for counter resets if isCounter is true),
// extrapolates if the first/last sample is close to the boundary, and returns
// the result as either per-second (if isRate is true) or overall.
func extrapolatedRate(samples []promql.Point, isCounter, isRate bool, stepTime time.Time, selectRange time.Duration) promql.Point {
	var (
		rangeStart = stepTime.UnixMilli() - selectRange.Milliseconds()
		rangeEnd   = stepTime.UnixMilli()
	)

	if len(samples) < 2 {
		return promql.Point{T: -1}
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

	return promql.Point{
		T: stepTime.UnixMilli(),
		V: resultValue,
	}
}

func instantValue(samples []promql.Point, isRate bool, stepTime time.Time, selectRange time.Duration) promql.Point {
	// No sense in trying to compute a rate without at least two points. Drop
	// this Vector element.
	if len(samples) < 2 {
		return promql.Point{T: -1}
	}

	lastSample := samples[len(samples)-1]
	previousSample := samples[len(samples)-2]

	var resultValue float64
	if isRate && lastSample.V < previousSample.V {
		// Counter reset.
		resultValue = lastSample.V
	} else {
		resultValue = lastSample.V - previousSample.V
	}

	sampledInterval := lastSample.T - previousSample.T
	if sampledInterval == 0 {
		// Avoid dividing by 0.
		return promql.Point{T: -1}
	}

	if isRate {
		// Convert to per-second.
		resultValue /= float64(sampledInterval) / 1000
	}

	return promql.Point{
		T: stepTime.UnixMilli(),
		V: resultValue,
	}
}

func maxOverTime(points []promql.Point) float64 {
	max := points[0].V
	for _, v := range points {
		if v.V > max || math.IsNaN(max) {
			max = v.V
		}
	}
	return max
}

func minOverTime(points []promql.Point) float64 {
	min := points[0].V
	for _, v := range points {
		if v.V < min || math.IsNaN(min) {
			min = v.V
		}
	}
	return min
}

func countOverTime(points []promql.Point) float64 {
	return float64(len(points))
}

func avgOverTime(points []promql.Point) float64 {
	var mean, count, c float64
	for _, v := range points {
		count++
		if math.IsInf(mean, 0) {
			if math.IsInf(v.V, 0) && (mean > 0) == (v.V > 0) {
				// The `mean` and `v.V` values are `Inf` of the same sign.  They
				// can't be subtracted, but the value of `mean` is correct
				// already.
				continue
			}
			if !math.IsInf(v.V, 0) && !math.IsNaN(v.V) {
				// At this stage, the mean is an infinite. If the added
				// value is neither an Inf or a Nan, we can keep that mean
				// value.
				// This is required because our calculation below removes
				// the mean value, which would look like Inf += x - Inf and
				// end up as a NaN.
				continue
			}
		}
		mean, c = kahanSumInc(v.V/count-mean/count, mean, c)
	}

	if math.IsInf(mean, 0) {
		return mean
	}
	return mean + c
}

func sumOverTime(points []promql.Point) float64 {
	var sum, c float64
	for _, v := range points {
		sum, c = kahanSumInc(v.V, sum, c)
	}
	if math.IsInf(sum, 0) {
		return sum
	}
	return sum + c
}

func kahanSumInc(inc, sum, c float64) (newSum, newC float64) {
	t := sum + inc
	// Using Neumaier improvement, swap if next term larger than sum.
	if math.Abs(sum) >= math.Abs(inc) {
		c += (sum - t) + inc
	} else {
		c += (inc - t) + sum
	}
	return t, c
}
func dropMetricName(l labels.Labels) labels.Labels {
	return labels.NewBuilder(l).Del(labels.MetricName).Labels()
}
