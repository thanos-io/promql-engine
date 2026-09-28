// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package ringbuffer

import (
	"github.com/thanos-io/promql-engine/warnings"

	"github.com/prometheus/prometheus/model/histogram"
)

// histogramRateState accumulates counter corrections for one evaluation window.
// Its retained histograms do not grow in number with the number of resets.
type histogramRateState struct {
	correction    *histogram.FloatHistogram
	hasCorrection bool
	err           error

	// A reset immediately after the first sample replaces the subtraction
	// baseline with an empty histogram using the second sample's layout.
	zeroBaseline *histogram.FloatHistogram
	ignoreFirst  bool

	minSchema          int32
	usingCustomBuckets bool
	mixedSchemas       bool
	warn               warnings.Warnings
}

func (s *histogramRateState) reset(first *histogram.FloatHistogram) {
	*s = histogramRateState{
		correction:         s.correction,
		zeroBaseline:       s.zeroBaseline,
		minSchema:          first.Schema,
		usingCustomBuckets: first.UsesCustomBuckets(),
	}
	if first.CounterResetHint == histogram.GaugeType {
		s.warn |= warnings.WarnNotCounter
	}
}

func (s *histogramRateState) push(previous, current *histogram.FloatHistogram, reset, secondSample bool) {
	if s.mixedSchemas {
		return
	}
	if reset && secondSample {
		s.ignoreFirst = true
		s.minSchema = current.Schema
		s.usingCustomBuckets = current.UsesCustomBuckets()
		if s.zeroBaseline == nil {
			s.zeroBaseline = &histogram.FloatHistogram{}
		}
		customValues := append(s.zeroBaseline.CustomValues[:0], current.CustomValues...)
		*s.zeroBaseline = histogram.FloatHistogram{
			Schema:       current.Schema,
			CustomValues: customValues,
		}
	}
	if current.CounterResetHint == histogram.GaugeType {
		s.warn |= warnings.WarnNotCounter
	}
	s.minSchema = min(s.minSchema, current.Schema)
	if current.UsesCustomBuckets() != s.usingCustomBuckets {
		s.mixedSchemas = true
		return
	}
	if !reset || secondSample || s.err != nil {
		return
	}
	if !s.hasCorrection {
		if s.correction == nil {
			s.correction = previous.Copy()
		} else {
			previous.CopyTo(s.correction)
		}
		s.hasCorrection = true
		return
	}
	_, _, _, s.err = s.correction.Add(previous)
}

func (s *histogramRateState) eval(first, last *histogram.FloatHistogram) (*histogram.FloatHistogram, warnings.Warnings, error) {
	// Match histogramRate's endpoint check before emitting other warnings.
	if last.UsesCustomBuckets() != s.usingCustomBuckets {
		return nil, warnings.WarnMixedExponentialCustomBuckets, nil
	}
	warn := s.warn
	if last.CounterResetHint == histogram.GaugeType {
		warn |= warnings.WarnNotCounter
	}
	if s.mixedSchemas {
		return nil, warn | warnings.WarnMixedExponentialCustomBuckets, nil
	}
	if s.err != nil {
		return nil, warn, s.err
	}

	baseline := first
	if s.ignoreFirst {
		baseline = s.zeroBaseline
	}
	h := last.CopyToSchema(s.minSchema)
	if _, _, reconciled, err := h.Sub(baseline); err != nil {
		return nil, warn, err
	} else if reconciled {
		warn |= warnings.WarnNHCBBoundsReconciled
	}
	if s.hasCorrection {
		if _, _, _, err := h.Add(s.correction); err != nil {
			return nil, warn, err
		}
	}
	h.CounterResetHint = histogram.GaugeType
	return h.Compact(0), warn, nil
}
