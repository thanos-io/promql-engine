// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package ringbuffer

import (
	"testing"

	"github.com/thanos-io/promql-engine/warnings"

	"github.com/prometheus/prometheus/model/histogram"
	"github.com/stretchr/testify/require"
)

// newTestHistogram returns a well-formed FloatHistogram whose fields all scale
// with mult, so a series of increasing mult values behaves like a monotonic
// native-histogram counter with Sub-compatible spans.
func newTestHistogram(mult float64) *histogram.FloatHistogram {
	return newTestHistogramWithHint(mult, histogram.NotCounterReset)
}

func newTestHistogramWithHint(mult float64, hint histogram.CounterResetHint) *histogram.FloatHistogram {
	return &histogram.FloatHistogram{
		CounterResetHint: hint,
		Schema:           0,
		ZeroThreshold:    0.001,
		ZeroCount:        1 * mult,
		Count:            10 * mult,
		Sum:              20 * mult,
		PositiveSpans:    []histogram.Span{{Offset: 0, Length: 2}},
		PositiveBuckets:  []float64{4 * mult, 5 * mult},
	}
}

func histogramSamples(ts []int64, mults []float64) []Sample {
	return histogramSamplesWithHint(ts, mults, histogram.NotCounterReset)
}

func histogramSamplesWithHint(ts []int64, mults []float64, hint histogram.CounterResetHint) []Sample {
	samples := make([]Sample, len(ts))
	for i := range ts {
		samples[i] = Sample{T: ts[i], V: Value{H: newTestHistogramWithHint(mults[i], hint)}}
	}
	return samples
}

// TestExtendedRateHistogramPerSecond verifies the core fix: xrate divides the
// native-histogram delta by the range in seconds while xincrease does not, so
// xrate == xincrease / rangeSeconds and the two are no longer identical.
func TestExtendedRateHistogramPerSecond(t *testing.T) {
	const (
		selectRange = int64(300000)
		stepTime    = int64(300000)
		offset      = int64(0)
	)
	rangeSeconds := float64(selectRange) / 1000

	// Baseline at rangeStart (t=0) and last sample at rangeEnd (t=selectRange)
	// so the adjust-to-range factor is exactly 1 and the relationship is clean.
	ts := []int64{0, 60000, 120000, 180000, 240000, 300000}
	mults := []float64{1, 2, 3, 4, 5, 6}

	xrateF, xrateH, xrateOK, _, err := extendedRate(histogramSamples(ts, mults), true, true, stepTime, selectRange, offset, 0)
	require.NoError(t, err)
	require.True(t, xrateOK)
	require.NotNil(t, xrateH)
	require.Zero(t, xrateF, "xrate on a native histogram must not emit a float value")

	xincF, xincH, xincOK, _, err := extendedRate(histogramSamples(ts, mults), true, false, stepTime, selectRange, offset, 0)
	require.NoError(t, err)
	require.True(t, xincOK)
	require.NotNil(t, xincH)
	require.Zero(t, xincF, "xincrease on a native histogram must not emit a float value")

	require.InEpsilon(t, xincH.Sum/rangeSeconds, xrateH.Sum, 1e-9)
	require.InEpsilon(t, xincH.Count/rangeSeconds, xrateH.Count, 1e-9)
	require.False(t, xrateH.Equals(xincH), "xrate and xincrease must not return identical histograms")

	// Bucket-level checks: the delta is last(mult=6) - first(mult=1), so every
	// field is 5x the per-mult unit. A Sum/Count-only check would miss a
	// per-bucket regression from Sub/Mul/CopyToSchema/Compact.
	require.Equal(t, []float64{20, 25}, xincH.PositiveBuckets)
	require.Equal(t, 50.0, xincH.Count)
	require.Equal(t, 5.0, xincH.ZeroCount)
	require.Equal(t, xincH.PositiveSpans, xrateH.PositiveSpans)
	for i := range xincH.PositiveBuckets {
		require.InDelta(t, xincH.PositiveBuckets[i]/rangeSeconds, xrateH.PositiveBuckets[i], 1e-9)
	}
}

// TestExtendedRateHistogramSingleSample verifies xincrease's zero-injection for
// a newly-appeared single-sample histogram series, and that xrate emits nothing
// in that case.
func TestExtendedRateHistogramSingleSample(t *testing.T) {
	const (
		selectRange = int64(300000)
		stepTime    = int64(300000)
		offset      = int64(0)
	)
	samples := histogramSamples([]int64{300000}, []float64{3})

	_, xincH, xincOK, _, err := extendedRate(samples, true, false, stepTime, selectRange, offset, 0)
	require.NoError(t, err)
	require.True(t, xincOK)
	require.NotNil(t, xincH, "xincrease must inject the single sample as the increase")
	require.Equal(t, histogram.GaugeType, xincH.CounterResetHint)
	require.Equal(t, newTestHistogram(3).Sum, xincH.Sum)

	_, xrateH, xrateOK, _, err := extendedRate(histogramSamples([]int64{300000}, []float64{3}), true, true, stepTime, selectRange, offset, 0)
	require.NoError(t, err)
	require.False(t, xrateOK, "xrate needs at least two samples and must emit nothing")
	require.Nil(t, xrateH)
}

// TestExtendedRateHistogramMixed verifies that a range mixing a histogram and a
// float yields no sample instead of a bogus float zero, and warns. Both
// orderings must behave identically: the histogram-first case is caught inside
// histogramRate, while the float-first case is caught by the top-level mixed
// guard. Without that guard, float-first would fall into the float branch and
// read .V.F (== 0) from the histogram sample, emitting a silent bogus zero.
func TestExtendedRateHistogramMixed(t *testing.T) {
	const (
		selectRange = int64(300000)
		stepTime    = int64(300000)
		offset      = int64(0)
	)

	orderings := []struct {
		name    string
		samples []Sample
	}{
		{
			name: "histogram then float",
			samples: []Sample{
				{T: 0, V: Value{H: newTestHistogram(1)}},
				{T: 300000, V: Value{F: 5}},
			},
		},
		{
			name: "float then histogram",
			samples: []Sample{
				{T: 0, V: Value{F: 5}},
				{T: 300000, V: Value{H: newTestHistogram(1)}},
			},
		},
	}

	funcs := []struct {
		name              string
		isCounter, isRate bool
	}{
		{"xincrease", true, false},
		{"xrate", true, true},
		{"xdelta", false, false},
	}

	for _, tc := range orderings {
		t.Run(tc.name, func(t *testing.T) {
			for _, fn := range funcs {
				f, h, ok, warn, err := extendedRate(tc.samples, fn.isCounter, fn.isRate, stepTime, selectRange, offset, 0)
				require.NoError(t, err, fn.name)
				require.False(t, ok, "%s: a mixed float/histogram range must not emit a sample", fn.name)
				require.Nil(t, h, fn.name)
				require.Zero(t, f, fn.name)
				require.NotZero(t, warn&warnings.WarnMixedFloatsHistograms, "%s: a mixed range must warn", fn.name)
			}
		})
	}
}

// TestExtendedRateHistogramCases exercises the windowing branches of the
// histogram path: xdelta, counter-reset correction, the drop-too-far-point
// path, zero-injection edge cases, duplicate timestamps, and a non-zero offset.
func TestExtendedRateHistogramCases(t *testing.T) {
	// A series with a reset: mult drops from 3 back to 1 at index 3. The samples
	// carry UnknownCounterReset so DetectReset inspects the buckets (a
	// NotCounterReset hint would suppress detection).
	resetTs := []int64{0, 60000, 120000, 180000, 240000, 300000}
	resetMults := []float64{1, 2, 3, 1, 2, 3}
	resetSamples := func() []Sample {
		return histogramSamplesWithHint(resetTs, resetMults, histogram.UnknownCounterReset)
	}

	cases := []struct {
		name            string
		samples         []Sample
		isCounter       bool
		isRate          bool
		stepTime        int64
		selectRange     int64
		offset          int64
		metricTs        int64
		expectHist      bool
		expectSum       float64
		expectBuckets   []float64
		expectCount     float64
		expectZeroCount float64
	}{
		{
			name:            "xdelta returns the raw delta with no per-second scaling",
			samples:         histogramSamples([]int64{0, 150000, 300000}, []float64{1, 2, 3}),
			isCounter:       false,
			isRate:          false,
			stepTime:        300000,
			selectRange:     300000,
			expectHist:      true,
			expectSum:       40,
			expectBuckets:   []float64{8, 10},
			expectCount:     20,
			expectZeroCount: 2,
		},
		{
			name:            "xincrease corrects a counter reset",
			samples:         resetSamples(),
			isCounter:       true,
			isRate:          false,
			stepTime:        300000,
			selectRange:     300000,
			expectHist:      true,
			expectSum:       100,
			expectBuckets:   []float64{20, 25},
			expectCount:     50,
			expectZeroCount: 5,
		},
		{
			name:            "xdelta does not correct a counter reset",
			samples:         resetSamples(),
			isCounter:       false,
			isRate:          false,
			stepTime:        300000,
			selectRange:     300000,
			expectHist:      true,
			expectSum:       40,
			expectBuckets:   []float64{8, 10},
			expectCount:     20,
			expectZeroCount: 2,
		},
		{
			name:        "xrate drops a too-far pre-range baseline",
			samples:     histogramSamples([]int64{-500000, 250000, 300000}, []float64{1, 2, 3}),
			isCounter:   true,
			isRate:      true,
			stepTime:    300000,
			selectRange: 300000,
			expectHist:  true,
			expectSum:   20.0 / 300.0,
		},
		{
			name:        "xrate with all samples before the range emits nothing",
			samples:     histogramSamples([]int64{50000, 100000}, []float64{1, 2}),
			isCounter:   true,
			isRate:      true,
			stepTime:    300000,
			selectRange: 100000,
			expectHist:  false,
		},
		{
			name:        "xincrease injects a flat multi-sample series",
			samples:     histogramSamples([]int64{0, 150000, 300000}, []float64{2, 2, 2}),
			isCounter:   true,
			isRate:      false,
			stepTime:    300000,
			selectRange: 300000,
			expectHist:  true,
			expectSum:   40,
		},
		{
			name:        "xincrease past the injection window returns a zero delta",
			samples:     histogramSamples([]int64{0, 150000, 300000}, []float64{2, 2, 2}),
			isCounter:   true,
			isRate:      false,
			stepTime:    300000,
			selectRange: 300000,
			metricTs:    -400000,
			expectHist:  true,
			expectSum:   0,
		},
		{
			name:        "duplicate timestamps emit nothing",
			samples:     histogramSamples([]int64{100000, 100000}, []float64{1, 2}),
			isCounter:   true,
			isRate:      false,
			stepTime:    300000,
			selectRange: 300000,
			expectHist:  false,
		},
		{
			name:        "xincrease honors a non-zero offset",
			samples:     histogramSamples([]int64{0, 150000, 300000}, []float64{1, 2, 3}),
			isCounter:   true,
			isRate:      false,
			stepTime:    360000,
			selectRange: 300000,
			offset:      60000,
			expectHist:  true,
			expectSum:   40,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, h, ok, _, err := extendedRate(tc.samples, tc.isCounter, tc.isRate, tc.stepTime, tc.selectRange, tc.offset, tc.metricTs)
			require.NoError(t, err)
			if !tc.expectHist {
				require.False(t, ok, "expected no sample")
				require.Nil(t, h)
				return
			}
			require.True(t, ok)
			require.NotNil(t, h)
			require.InDelta(t, tc.expectSum, h.Sum, 1e-9)
			if tc.expectBuckets != nil {
				require.Equal(t, tc.expectBuckets, h.PositiveBuckets)
				require.InDelta(t, tc.expectCount, h.Count, 1e-9)
				require.InDelta(t, tc.expectZeroCount, h.ZeroCount, 1e-9)
			}
		})
	}
}

// TestExtendedRateHistogramWarnings verifies the warning bits reachable now that
// native histograms flow through the extended range functions.
func TestExtendedRateHistogramWarnings(t *testing.T) {
	const (
		selectRange = int64(300000)
		stepTime    = int64(300000)
	)
	ts := []int64{0, 150000, 300000}

	t.Run("xincrease on a gauge histogram warns it is not a counter", func(t *testing.T) {
		samples := histogramSamplesWithHint(ts, []float64{1, 2, 3}, histogram.GaugeType)
		_, _, _, warn, err := extendedRate(samples, true, false, stepTime, selectRange, 0, 0)
		require.NoError(t, err)
		require.NotZero(t, warn&warnings.WarnNotCounter)
	})

	t.Run("xdelta on a counter histogram warns it is not a gauge", func(t *testing.T) {
		_, _, _, warn, err := extendedRate(histogramSamples(ts, []float64{1, 2, 3}), false, false, stepTime, selectRange, 0, 0)
		require.NoError(t, err)
		require.NotZero(t, warn&warnings.WarnNotGauge)
	})
}

// TestExtendedRateHistogramSchemaChange exercises histogramRate's min-schema /
// CopyToSchema path through the extended functions: the last sample uses a finer
// schema than the baseline, so the delta must be down-converted to the coarser
// schema. xdelta is used to avoid the counter reset/null-out handling.
func TestExtendedRateHistogramSchemaChange(t *testing.T) {
	mk := func(mult float64, schema int32) *histogram.FloatHistogram {
		return &histogram.FloatHistogram{
			CounterResetHint: histogram.NotCounterReset,
			Schema:           schema,
			ZeroThreshold:    0.001,
			ZeroCount:        mult,
			Count:            5 * mult,
			Sum:              10 * mult,
			PositiveSpans:    []histogram.Span{{Offset: 0, Length: 2}},
			PositiveBuckets:  []float64{2 * mult, 2 * mult},
		}
	}
	samples := []Sample{
		{T: 0, V: Value{H: mk(1, 0)}},
		{T: 150000, V: Value{H: mk(2, 0)}},
		{T: 300000, V: Value{H: mk(3, 1)}},
	}

	_, h, ok, _, err := extendedRate(samples, false, false, 300000, 300000, 0, 0)
	require.NoError(t, err)
	require.True(t, ok)
	require.NotNil(t, h)
	require.Equal(t, int32(0), h.Schema, "result must be down-converted to the minimum schema in the range")
}

// TestExtendedRateHistogramZeroInjectionWarns is a regression test for the
// xincrease zero-injection path: it must scan every sample for a gauge hint,
// not just the first. sameHistogramValues compares values via
// FloatHistogram.Equals, which ignores CounterResetHint, so a range whose first
// sample is a counter and a later equal-valued sample is a gauge takes the
// injection path yet must still emit WarnNotCounter.
func TestExtendedRateHistogramZeroInjectionWarns(t *testing.T) {
	const (
		selectRange = int64(300000)
		stepTime    = int64(300000)
	)
	samples := []Sample{
		{T: 0, V: Value{H: newTestHistogramWithHint(3, histogram.NotCounterReset)}},
		{T: 150000, V: Value{H: newTestHistogramWithHint(3, histogram.GaugeType)}},
	}
	_, h, ok, warn, err := extendedRate(samples, true, false, stepTime, selectRange, 0, 0)
	require.NoError(t, err)
	require.True(t, ok, "equal-valued samples must take the zero-injection path")
	require.NotNil(t, h)
	require.Equal(t, newTestHistogram(3).Sum, h.Sum, "must return the injected sample value, confirming the injection path")
	require.NotZero(t, warn&warnings.WarnNotCounter, "a gauge sample after the first must still warn")
}

// TestExtendedRateHistogramSingleSampleUndefined pins finding #3's intentional
// behavior: unlike the float path (which emits 0), a single-sample histogram
// xrate/xdelta and a single-sample xincrease past the injection window emit
// nothing rather than a zero-shaped histogram.
func TestExtendedRateHistogramSingleSampleUndefined(t *testing.T) {
	const selectRange = int64(300000)

	t.Run("xdelta emits nothing", func(t *testing.T) {
		samples := histogramSamples([]int64{300000}, []float64{3})
		_, h, ok, _, err := extendedRate(samples, false, false, 300000, selectRange, 0, 0)
		require.NoError(t, err)
		require.False(t, ok)
		require.Nil(t, h)
	})

	t.Run("xincrease past the injection window emits nothing", func(t *testing.T) {
		samples := histogramSamples([]int64{600000}, []float64{3})
		_, h, ok, _, err := extendedRate(samples, true, false, 600000, selectRange, 0, 0)
		require.NoError(t, err)
		require.False(t, ok)
		require.Nil(t, h)
	})
}
