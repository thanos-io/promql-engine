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

var IgnorePoint = promql.Point{T: -1, V: 0}

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
	case "stddev_over_time":
		return func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample {
			return promql.Sample{
				Point: promql.Point{
					T: stepTime.UnixMilli(),
					V: stddevOverTime(points),
				},
				Metric: labels,
			}
		}, nil
	case "stdvar_over_time":
		return func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample {
			return promql.Sample{
				Point: promql.Point{
					T: stepTime.UnixMilli(),
					V: stdvarOverTime(points),
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
	case "last_over_time":
		return func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample {
			return promql.Sample{
				Point: promql.Point{
					T: stepTime.UnixMilli(),
					V: points[len(points)-1].V,
				},
				Metric: labels,
			}
		}, nil
	case "present_over_time":
		return func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample {
			return promql.Sample{
				Point: promql.Point{
					T: stepTime.UnixMilli(),
					V: 1,
				},
				Metric: labels,
			}
		}, nil
	case "changes":
		return func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample {
			return promql.Sample{
				Point: promql.Point{
					T: stepTime.UnixMilli(),
					V: changes(points),
				},
				Metric: labels,
			}
		}, nil
	case "deriv":
		return func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample {
			if len(points) < 2 {
				return promql.Sample{
					Point: IgnorePoint,
				}
			}
			return promql.Sample{
				Point: promql.Point{
					T: stepTime.UnixMilli(),
					V: deriv(points),
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
	case "resets":
		return func(labels labels.Labels, points []promql.Point, stepTime time.Time) promql.Sample {
			resets := 0
			prev := points[0].V
			for _, sample := range points[1:] {
				current := sample.V
				if current < prev {
					resets++
				}
				prev = current
			}
			return promql.Sample{
				Point: promql.Point{
					T: stepTime.UnixMilli(),
					V: float64(resets),
				},
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
		mean, c = KahanSumInc(v.V/count-mean/count, mean, c)
	}

	if math.IsInf(mean, 0) {
		return mean
	}
	return mean + c
}

func sumOverTime(points []promql.Point) float64 {
	var sum, c float64
	for _, v := range points {
		sum, c = KahanSumInc(v.V, sum, c)
	}
	if math.IsInf(sum, 0) {
		return sum
	}
	return sum + c
}

func stddevOverTime(points []promql.Point) float64 {
	var count float64
	var mean, cMean float64
	var aux, cAux float64
	for _, v := range points {
		count++
		delta := v.V - (mean + cMean)
		mean, cMean = KahanSumInc(delta/count, mean, cMean)
		aux, cAux = KahanSumInc(delta*(v.V-(mean+cMean)), aux, cAux)
	}
	return math.Sqrt((aux + cAux) / count)
}

func stdvarOverTime(points []promql.Point) float64 {
	var count float64
	var mean, cMean float64
	var aux, cAux float64
	for _, v := range points {
		count++
		delta := v.V - (mean + cMean)
		mean, cMean = KahanSumInc(delta/count, mean, cMean)
		aux, cAux = KahanSumInc(delta*(v.V-(mean+cMean)), aux, cAux)
	}
	return (aux + cAux) / count
}

func changes(points []promql.Point) float64 {
	var count float64
	prev := points[0].V
	count = 0
	for _, sample := range points[1:] {
		current := sample.V
		if current != prev && !(math.IsNaN(current) && math.IsNaN(prev)) {
			count++
		}
		prev = current
	}
	return count
}

func deriv(points []promql.Point) float64 {
	// We pass in an arbitrary timestamp that is near the values in use
	// to avoid floating point accuracy issues, see
	// https://github.com/prometheus/prometheus/issues/2674
	slope, _ := linearRegression(points, points[0].T)
	return slope
}

func linearRegression(samples []promql.Point, interceptTime int64) (slope, intercept float64) {
	var (
		n          float64
		sumX, cX   float64
		sumY, cY   float64
		sumXY, cXY float64
		sumX2, cX2 float64
		initY      float64
		constY     bool
	)
	initY = samples[0].V
	constY = true
	for i, sample := range samples {
		// Set constY to false if any new y values are encountered.
		if constY && i > 0 && sample.V != initY {
			constY = false
		}
		n += 1.0
		x := float64(sample.T-interceptTime) / 1e3
		sumX, cX = KahanSumInc(x, sumX, cX)
		sumY, cY = KahanSumInc(sample.V, sumY, cY)
		sumXY, cXY = KahanSumInc(x*sample.V, sumXY, cXY)
		sumX2, cX2 = KahanSumInc(x*x, sumX2, cX2)
	}
	if constY {
		if math.IsInf(initY, 0) {
			return math.NaN(), math.NaN()
		}
		return 0, initY
	}
	sumX = sumX + cX
	sumY = sumY + cY
	sumXY = sumXY + cXY
	sumX2 = sumX2 + cX2

	covXY := sumXY - sumX*sumY/n
	varX := sumX2 - sumX*sumX/n

	slope = covXY / varX
	intercept = sumY/n - slope*sumX/n
	return slope, intercept
}

func KahanSumInc(inc, sum, c float64) (newSum, newC float64) {
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
