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

var InvalidSample = promql.Sample{Point: promql.Point{T: -1, V: 0}}

// FunctionCall represents functions as defined in https://prometheus.io/docs/prometheus/latest/querying/functions/
type FunctionCall func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample

var Funcs = map[string]FunctionCall{
	"sum_over_time": func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample {
		return ignoreEmpty(points, sample(labels, stepTime, points, sumOverTime))
	},
	"max_over_time": func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample {
		return ignoreEmpty(points, sample(labels, stepTime, points, maxOverTime))
	},
	"min_over_time": func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample {
		return ignoreEmpty(points, sample(labels, stepTime, points, minOverTime))
	},
	"avg_over_time": func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample {
		return ignoreEmpty(points, sample(labels, stepTime, points, avgOverTime))
	},
	"stddev_over_time": func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample {
		return ignoreEmpty(points, sample(labels, stepTime, points, stddevOverTime))
	},
	"stdvar_over_time": func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample {
		return ignoreEmpty(points, sample(labels, stepTime, points, stdvarOverTime))
	},
	"count_over_time": func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample {
		return ignoreEmpty(points, sample(labels, stepTime, points, countOverTime))
	},
	"last_over_time": func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample {
		return ignoreEmpty(points, sample(labels, stepTime, points, func(points []promql.Point) float64 {
			return points[len(points)-1].V
		}))
	},
	"present_over_time": func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample {
		return ignoreEmpty(points, sample(labels, stepTime, nil, func([]promql.Point) float64 { return 1 }))
	},
	"changes": func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample {
		return ignoreEmpty(points, sample(labels, stepTime, points, changes))
	},
	"deriv": func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample {
		return ignoreSingleValue(points, sample(labels, stepTime, points, deriv))
	},
	"irate": func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample {
		return ignoreSingleValue(points, maybeSample(labels, stepTime, points, func(points []promql.Point) (float64, bool) {
			return instantValue(points, true)
		}))
	},
	"idelta": func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample {
		return ignoreSingleValue(points, maybeSample(labels, stepTime, points, func(points []promql.Point) (float64, bool) {
			return instantValue(points, false)
		}))
	},
	"vector": func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample {
		return ignoreEmpty(points, sample(labels, stepTime, points, func(points []promql.Point) float64 {
			return points[0].V
		}))
	},
	"rate": func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample {
		return ignoreSingleValue(points, sample(labels, stepTime, points, func(points []promql.Point) float64 {
			return extrapolatedRate(points, true, true, stepTime, selectRange)
		}))
	},
	"delta": func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample {
		return ignoreSingleValue(points, sample(labels, stepTime, points, func(points []promql.Point) float64 {
			return extrapolatedRate(points, false, false, stepTime, selectRange)
		}))
	},
	"increase": func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample {
		return ignoreSingleValue(points, sample(labels, stepTime, points, func(points []promql.Point) float64 {
			return extrapolatedRate(points, true, false, stepTime, selectRange)
		}))
	},
	"resets": func(labels labels.Labels, points []promql.Point, stepTime time.Time, selectRange time.Duration) promql.Sample {
		return ignoreEmpty(points, sample(labels, stepTime, points, resets))
	},
}

func NewFunctionCall(f *parser.Function) (FunctionCall, error) {
	if call, ok := Funcs[f.Name]; ok {
		return call, nil
	}
	msg := fmt.Sprintf("unknown function: %s", f.Name)
	return nil, errors.Wrap(parse.ErrNotSupportedExpr, msg)
}

// extrapolatedRate is a utility function for rate/increase/delta.
// It calculates the rate (allowing for counter resets if isCounter is true),
// extrapolates if the first/last sample is close to the boundary, and returns
// the result as either per-second (if isRate is true) or overall.
func extrapolatedRate(samples []promql.Point, isCounter, isRate bool, stepTime time.Time, selectRange time.Duration) float64 {
	var (
		rangeStart = stepTime.UnixMilli() - selectRange.Milliseconds()
		rangeEnd   = stepTime.UnixMilli()
	)

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

	return resultValue
}

func instantValue(samples []promql.Point, isRate bool) (float64, bool) {
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
		return 0, false
	}

	if isRate {
		// Convert to per-second.
		resultValue /= float64(sampledInterval) / 1000
	}

	return resultValue, true
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

func resets(points []promql.Point) float64 {
	count := 0
	prev := points[0].V
	for _, sample := range points[1:] {
		current := sample.V
		if current < prev {
			count++
		}
		prev = current
	}

	return float64(count)
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
	return labels.NewBuilder(l).Del(labels.MetricName).Labels(nil)
}

type valFunc func([]promql.Point) float64

func sample(lbls labels.Labels, stepTime time.Time, points []promql.Point, valFunc valFunc) func() promql.Sample {
	return func() promql.Sample {
		return promql.Sample{
			Metric: lbls,
			Point: promql.Point{
				T: stepTime.UnixMilli(),
				V: valFunc(points),
			},
		}
	}
}

type maybeValFunc func([]promql.Point) (float64, bool)

func maybeSample(lbls labels.Labels, stepTime time.Time, points []promql.Point, valFunc maybeValFunc) func() promql.Sample {
	return func() promql.Sample {
		val, ok := valFunc(points)
		if !ok {
			return InvalidSample
		}
		return promql.Sample{
			Metric: lbls,
			Point: promql.Point{
				T: stepTime.UnixMilli(),
				V: val,
			},
		}
	}
}

func ignoreEmpty(points []promql.Point, f func() promql.Sample) promql.Sample {
	if len(points) == 0 {
		return InvalidSample
	}
	return f()
}

func ignoreSingleValue(points []promql.Point, f func() promql.Sample) promql.Sample {
	if len(points) < 2 {
		return InvalidSample
	}
	return f()
}
