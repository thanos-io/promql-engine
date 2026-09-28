// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package ringbuffer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/query"

	"github.com/prometheus/prometheus/model/histogram"
	"github.com/stretchr/testify/require"
)

func TestRateBufferAccumulatesFloatResetsPerStep(t *testing.T) {
	const second = int64(1000)
	const selectRange = 120 * second

	opts := query.Options{
		Start: time.UnixMilli(120 * second),
		End:   time.UnixMilli(300 * second),
		Step:  60 * time.Second,
	}
	streaming := NewRateBuffer(context.Background(), opts, true, false, selectRange, 0)
	call, err := NewRangeVectorFunc("increase")
	require.NoError(t, err)
	reference := New(context.Background(), 8, selectRange, 0, call)

	samples := []Sample{
		{T: 30 * second, V: Value{F: 10}},
		{T: 60 * second, V: Value{F: 20}},
		{T: 90 * second, V: Value{F: 5}},
		{T: 120 * second, V: Value{F: 15}},
		{T: 150 * second, V: Value{F: 3}},
		{T: 180 * second, V: Value{F: 8}},
		{T: 210 * second, V: Value{F: 2}},
		{T: 240 * second, V: Value{F: 12}},
		{T: 270 * second, V: Value{F: 22}},
		{T: 300 * second, V: Value{F: 1}},
	}

	nextSample := 0
	for evalT := 120 * second; evalT <= 300*second; evalT += 60 * second {
		mint := evalT - selectRange
		streaming.Reset(mint, evalT)
		reference.Reset(mint, evalT)

		for nextSample < len(samples) && samples[nextSample].T <= evalT {
			sample := samples[nextSample]
			streaming.Push(sample.T, sample.V)
			reference.Push(sample.T, sample.V)
			nextSample++
		}

		got, _, gotOK, gotWarnings, err := streaming.Eval(context.Background(), 0, 0)
		require.NoError(t, err)
		want, _, wantOK, wantWarnings, err := reference.Eval(context.Background(), 0, 0)
		require.NoError(t, err)
		require.Equal(t, wantOK, gotOK)
		require.Equal(t, wantWarnings, gotWarnings)
		require.InDelta(t, want, got, 1e-12)
		require.Equal(t, reference.SampleCount(), streaming.SampleCount())
	}
}

func TestRateBufferHistograms(t *testing.T) {
	hist := func(mult float64) Value {
		return Value{H: newTestHistogramWithHint(mult, histogram.UnknownCounterReset)}
	}
	hint := func(mult float64, resetHint histogram.CounterResetHint) Value {
		return Value{H: newTestHistogramWithHint(mult, resetHint)}
	}
	schema := func(mult float64, schema int32) Value {
		v := hist(mult)
		v.H.Schema = schema
		return v
	}
	custom := func(mult float64, bounds ...float64) Value {
		buckets := make([]float64, len(bounds)+1)
		for i := range buckets {
			buckets[i] = mult
		}
		return Value{H: &histogram.FloatHistogram{
			Schema:          histogram.CustomBucketsSchema,
			CustomValues:    bounds,
			Count:           mult * float64(len(buckets)),
			Sum:             mult * 10,
			PositiveSpans:   []histogram.Span{{Offset: 0, Length: uint32(len(buckets))}},
			PositiveBuckets: buckets,
		}}
	}
	for _, tc := range []struct {
		name   string
		values []Value
	}{
		{"no resets", []Value{hist(1), hist(2), hist(3), hist(4), hist(5), hist(6)}},
		{"repeated resets", []Value{hist(4), hist(1), hist(3), hist(1), hist(2), hist(1), hist(5), hist(1)}},
		{"explicit resets without a decrease", []Value{hist(1), hint(2, histogram.CounterReset), hist(3), hint(4, histogram.CounterReset), hist(5), hist(6)}},
		{"no reset hints with a decrease", []Value{hint(4, histogram.NotCounterReset), hint(1, histogram.NotCounterReset), hint(3, histogram.NotCounterReset), hist(4)}},
		{"gauge hint in the middle", []Value{hist(1), hint(2, histogram.GaugeType), hist(3), hist(4), hist(5)}},
		{"schema changes", []Value{schema(4, 1), schema(1, 0), schema(3, 0), schema(4, 1), schema(1, 0), schema(2, 1)}},
		{"custom buckets", []Value{custom(4, 1, 2), custom(1, 1, 2), custom(3, 1, 2), custom(1, 1, 2), custom(2, 1, 2)}},
		{"custom bounds changes", []Value{custom(4, 1, 2, 4), custom(1, 1, 2, 4), custom(3, 1, 4), custom(1, 1, 2), custom(2, 1, 2), custom(1, 1, 2, 4)}},
		{"first reset discards exponential layout", []Value{hist(4), custom(1, 1, 2), custom(3, 1, 2), custom(1, 1, 2), custom(2, 1, 2)}},
		{"first reset discards custom layout", []Value{custom(4, 1, 2), hist(1), hist(3), hist(1), hist(2)}},
		{"incompatible intermediate layout", []Value{hist(1), hist(2), custom(1, 1, 2), hist(4), hist(5), hist(6)}},
		{"intermediate float", []Value{hist(1), {F: 2}, hist(3), hist(4), hist(5), hist(6)}},
		{"intermediate histogram", []Value{{F: 1}, hist(2), {F: 3}, {F: 4}, {F: 5}, {F: 6}}},
	} {
		for _, function := range []string{"rate", "increase", "delta"} {
			for _, step := range []time.Duration{30 * time.Second, 60 * time.Second, 180 * time.Second} {
				for _, offset := range []time.Duration{0, 30 * time.Second} {
					t.Run(fmt.Sprintf("%s/%s/step=%s/offset=%s", tc.name, function, step, offset), func(t *testing.T) {
						const window = 120 * time.Second
						opts := query.Options{
							Start: time.Unix(0, 0).Add(step + offset),
							End:   time.Unix(0, 0).Add(12*time.Minute + offset),
							Step:  step,
						}
						streaming := NewRateBuffer(context.Background(), opts, function != "delta", function == "rate", window.Milliseconds(), offset.Milliseconds())
						call, err := NewRangeVectorFunc(function)
						require.NoError(t, err)
						reference := New(context.Background(), 8, window.Milliseconds(), offset.Milliseconds(), call)
						// Reuse the source histogram just as the TSDB iterator does.
						reader := &histogram.FloatHistogram{}
						next := 0
						for evalT := opts.Start.UnixMilli(); evalT <= opts.End.UnixMilli(); evalT += step.Milliseconds() {
							maxt := evalT - offset.Milliseconds()
							mint := maxt - window.Milliseconds()
							streaming.Reset(mint, evalT)
							reference.Reset(mint, evalT)
							for next < len(tc.values) && int64(next+1)*30_000 <= maxt {
								v := tc.values[next]
								if v.H != nil {
									v.H.CopyTo(reader)
									v.H = reader
								}
								ts := int64(next+1) * 30_000
								streaming.Push(ts, v)
								reference.Push(ts, v)
								next++
							}
							gotF, gotH, gotOK, gotWarn, gotErr := streaming.Eval(context.Background(), 0, 0)
							wantF, wantH, wantOK, wantWarn, wantErr := reference.Eval(context.Background(), 0, 0)
							require.Equal(t, wantErr, gotErr, "evaluation at %d", evalT)
							require.Equal(t, wantOK, gotOK, "evaluation at %d", evalT)
							require.Equal(t, wantWarn, gotWarn, "evaluation at %d", evalT)
							require.InDelta(t, wantF, gotF, 1e-12, "evaluation at %d", evalT)
							require.Equal(t, wantH, gotH, "evaluation at %d", evalT)
							require.Equal(t, reference.SampleCount(), streaming.SampleCount(), "evaluation at %d", evalT)
						}
					})
				}
			}
		}
	}
}
