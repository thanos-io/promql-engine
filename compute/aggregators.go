// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package compute

import (
	"math"
	"sort"

	"github.com/thanos-io/promql-engine/warnings"

	"github.com/efficientgo/core/errors"
	"github.com/prometheus/prometheus/model/histogram"
	"github.com/prometheus/prometheus/promql/parser/posrange"
	"github.com/prometheus/prometheus/util/annotations"
	"gonum.org/v1/gonum/floats"
)

type ValueType int

const (
	NoValue ValueType = iota
	SingleTypeValue
	MixedTypeValue
)

type Accumulator interface {
	Add(v float64, h *histogram.FloatHistogram) error
	Value() (float64, *histogram.FloatHistogram)
	ValueType() ValueType
	Reset(float64)
}

// CopyableAccumulator is an Accumulator that can be efficiently copied.
// This is used by sliding window buffers to create checkpoints.
type CopyableAccumulator interface {
	Accumulator
	Copy() CopyableAccumulator
}

// MergeableAccumulator is an Accumulator that can merge another accumulator's state.
// This enables O(1) checkpoint restoration by merging prefix + suffix accumulators.
type MergeableAccumulator interface {
	CopyableAccumulator
	// Merge combines another accumulator's state into this one.
	// The other accumulator should be of the same concrete type.
	Merge(other MergeableAccumulator) error
}

// SubtractableAccumulator is an Accumulator that supports removing values.
// This enables O(1) sliding window updates by subtracting removed samples.
type SubtractableAccumulator interface {
	CopyableAccumulator
	// Sub removes a value from the accumulator (inverse of Add).
	Sub(v float64, h *histogram.FloatHistogram) error
}

type VectorAccumulator interface {
	AddVector(vs []float64, hs []*histogram.FloatHistogram) error
	Value() (float64, *histogram.FloatHistogram)
	ValueType() ValueType
	Reset(float64)
}

type SumAcc struct {
	value        float64
	compensation float64
	histSum      *histogram.FloatHistogram
	hasFloatVal  bool
	hasHistError bool
}

func NewSumAcc() *SumAcc {
	return &SumAcc{}
}

func (s *SumAcc) AddVector(float64s []float64, histograms []*histogram.FloatHistogram) error {
	if len(float64s) > 0 {
		s.value, s.compensation = KahanSumInc(compensatedSum(float64s), s.value, s.compensation)
		s.hasFloatVal = true
	}

	var err error
	if len(histograms) > 0 {
		s.histSum, err = histogramSum(s.histSum, histograms)
	}
	return err
}

func (s *SumAcc) Add(v float64, h *histogram.FloatHistogram) error {
	if h == nil {
		s.hasFloatVal = true
		s.value, s.compensation = KahanSumInc(v, s.value, s.compensation)
		return nil
	}
	if s.hasHistError {
		return nil
	}
	if s.histSum == nil {
		s.histSum = h.Copy()
		return nil
	}
	// The histogram being added must have an equal or larger schema.
	// https://github.com/prometheus/prometheus/blob/57bcbf18880f7554ae34c5b341d52fc53f059a97/promql/engine.go#L2448-L2456
	var err error
	if h.Schema >= s.histSum.Schema {
		if s.histSum, err = s.histSum.Add(h); err != nil {
			s.histSum = nil
			s.hasHistError = true
			if errors.Is(err, histogram.ErrHistogramsIncompatibleSchema) {
				return annotations.MixedExponentialCustomHistogramsWarning
			}
			if errors.Is(err, histogram.ErrHistogramsIncompatibleBounds) {
				return annotations.IncompatibleCustomBucketsHistogramsWarning
			}
			return err
		}
	} else {
		t := h.Copy()
		if s.histSum, err = t.Add(s.histSum); err != nil {
			s.histSum = nil
			s.hasHistError = true
			if errors.Is(err, histogram.ErrHistogramsIncompatibleSchema) {
				return annotations.MixedExponentialCustomHistogramsWarning
			}
			if errors.Is(err, histogram.ErrHistogramsIncompatibleBounds) {
				return annotations.IncompatibleCustomBucketsHistogramsWarning
			}
			return err
		}
		s.histSum = t
	}
	return nil
}

func (s *SumAcc) Value() (float64, *histogram.FloatHistogram) {
	if s.histSum != nil {
		s.histSum.Compact(0)
	}
	return s.value + s.compensation, s.histSum
}

func (s *SumAcc) ValueType() ValueType {
	if s.hasFloatVal && s.histSum != nil {
		return MixedTypeValue
	}
	if s.hasFloatVal || s.histSum != nil {
		return SingleTypeValue
	}
	return NoValue
}

func (s *SumAcc) Reset(_ float64) {
	s.histSum = nil
	s.hasFloatVal = false
	s.hasHistError = false
	s.value = 0
	s.compensation = 0
}

func (s *SumAcc) Copy() CopyableAccumulator {
	cp := &SumAcc{
		value:        s.value,
		compensation: s.compensation,
		hasFloatVal:  s.hasFloatVal,
		hasHistError: s.hasHistError,
	}
	if s.histSum != nil {
		cp.histSum = s.histSum.Copy()
	}
	return cp
}

func (s *SumAcc) Merge(other MergeableAccumulator) error {
	o := other.(*SumAcc)
	s.value, s.compensation = KahanSumInc(o.value+o.compensation, s.value, s.compensation)
	s.hasFloatVal = s.hasFloatVal || o.hasFloatVal
	s.hasHistError = s.hasHistError || o.hasHistError
	if o.histSum != nil {
		var err error
		s.histSum, err = histogramSum(s.histSum, []*histogram.FloatHistogram{o.histSum})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SumAcc) Sub(v float64, h *histogram.FloatHistogram) error {
	if h == nil {
		// Subtract float value using Kahan summation (add negative)
		s.value, s.compensation = KahanSumInc(-v, s.value, s.compensation)
		return nil
	}
	// Subtract histogram by adding negated copy
	if s.histSum != nil && !s.hasHistError {
		negH := h.Copy().Mul(-1)
		var err error
		s.histSum, err = s.histSum.Add(negH)
		if err != nil {
			s.hasHistError = true
			return err
		}
	}
	return nil
}

func NewMaxAcc() *MaxAcc {
	return &MaxAcc{}
}

type MaxAcc struct {
	value    float64
	hasValue bool
}

func (c *MaxAcc) AddVector(vs []float64, hs []*histogram.FloatHistogram) error {
	var warn error
	if len(hs) > 0 {
		warn = annotations.NewHistogramIgnoredInAggregationInfo("max", posrange.PositionRange{})
	}
	if len(vs) == 0 {
		return warn
	}

	fst, rem := vs[0], vs[1:]
	warn = warnings.Coalesce(warn, c.Add(fst, nil))
	if len(rem) == 0 {
		return warn
	}
	return warnings.Coalesce(warn, c.Add(floats.Max(rem), nil))
}

func (c *MaxAcc) Add(v float64, h *histogram.FloatHistogram) error {
	if h != nil {
		if c.hasValue {
			return annotations.NewHistogramIgnoredInAggregationInfo("max", posrange.PositionRange{})
		}
		return nil
	}

	if !c.hasValue {
		c.value = v
		c.hasValue = true
		return nil
	}
	if c.value < v || math.IsNaN(c.value) {
		c.value = v
	}
	return nil
}

func (c *MaxAcc) Value() (float64, *histogram.FloatHistogram) {
	return c.value, nil
}

func (c *MaxAcc) ValueType() ValueType {
	if c.hasValue {
		return SingleTypeValue
	} else {
		return NoValue
	}
}

func (c *MaxAcc) Reset(_ float64) {
	c.hasValue = false
	c.value = 0
}

func NewMinAcc() *MinAcc {
	return &MinAcc{}
}

type MinAcc struct {
	value    float64
	hasValue bool
}

func (c *MinAcc) AddVector(vs []float64, hs []*histogram.FloatHistogram) error {
	var warn error
	if len(hs) > 0 {
		warn = annotations.NewHistogramIgnoredInAggregationInfo("min", posrange.PositionRange{})
	}
	if len(vs) == 0 {
		return warn
	}

	fst, rem := vs[0], vs[1:]
	warn = warnings.Coalesce(warn, c.Add(fst, nil))
	if len(rem) == 0 {
		return warn
	}

	return warnings.Coalesce(warn, c.Add(floats.Min(rem), nil))
}

func (c *MinAcc) Add(v float64, h *histogram.FloatHistogram) error {
	if h != nil {
		if c.hasValue {
			return annotations.NewHistogramIgnoredInAggregationInfo("min", posrange.PositionRange{})
		}
		return nil
	}

	if !c.hasValue {
		c.value = v
		c.hasValue = true
		return nil
	}
	if c.value > v || math.IsNaN(c.value) {
		c.value = v
	}
	return nil
}

func (c *MinAcc) Value() (float64, *histogram.FloatHistogram) {
	return c.value, nil
}

func (c *MinAcc) ValueType() ValueType {
	if c.hasValue {
		return SingleTypeValue
	} else {
		return NoValue
	}
}

func (c *MinAcc) Reset(_ float64) {
	c.hasValue = false
	c.value = 0
}

func NewGroupAcc() *GroupAcc {
	return &GroupAcc{}
}

type GroupAcc struct {
	value    float64
	hasValue bool
}

func (c *GroupAcc) AddVector(vs []float64, hs []*histogram.FloatHistogram) error {
	if len(vs) == 0 && len(hs) == 0 {
		return nil
	}
	c.hasValue = true
	c.value = 1
	return nil
}

func (c *GroupAcc) Add(v float64, h *histogram.FloatHistogram) error {
	c.hasValue = true
	c.value = 1
	return nil
}

func (c *GroupAcc) Value() (float64, *histogram.FloatHistogram) {
	return c.value, nil
}

func (c *GroupAcc) ValueType() ValueType {
	if c.hasValue {
		return SingleTypeValue
	} else {
		return NoValue
	}
}

func (c *GroupAcc) Reset(_ float64) {
	c.hasValue = false
	c.value = 0
}

type CountAcc struct {
	value    float64
	hasValue bool
}

func NewCountAcc() *CountAcc {
	return &CountAcc{}
}

func (c *CountAcc) AddVector(vs []float64, hs []*histogram.FloatHistogram) error {
	if len(vs) > 0 || len(hs) > 0 {
		c.hasValue = true
		c.value += float64(len(vs)) + float64(len(hs))
	}
	return nil
}

func (c *CountAcc) Add(v float64, h *histogram.FloatHistogram) error {
	c.hasValue = true
	c.value += 1
	return nil
}

func (c *CountAcc) Value() (float64, *histogram.FloatHistogram) {
	return c.value, nil
}

func (c *CountAcc) ValueType() ValueType {
	if c.hasValue {
		return SingleTypeValue
	} else {
		return NoValue
	}
}
func (c *CountAcc) Reset(_ float64) {
	c.hasValue = false
	c.value = 0
}

type AvgAcc struct {
	kahanSum    float64
	kahanC      float64
	avg         float64
	incremental bool
	count       int64
	hasValue    bool

	histSum        *histogram.FloatHistogram
	histScratch    *histogram.FloatHistogram
	histSumScratch *histogram.FloatHistogram
	histCount      float64
	hasHistError   bool
}

func NewAvgAcc() *AvgAcc {
	return &AvgAcc{}
}

func (a *AvgAcc) Add(v float64, h *histogram.FloatHistogram) error {
	if h != nil {
		if a.hasHistError {
			return nil
		}
		a.histCount++
		if a.histSum == nil {
			a.histSum = h.Copy()
			a.histScratch = &histogram.FloatHistogram{}
			a.histSumScratch = &histogram.FloatHistogram{}
			return nil
		}

		h.CopyTo(a.histScratch)
		left := a.histScratch.Div(a.histCount)
		a.histSum.CopyTo(a.histSumScratch)
		right := a.histSumScratch.Div(a.histCount)
		toAdd, err := left.Sub(right)
		if err != nil {
			a.histSum = nil
			a.histCount = 0
			a.hasHistError = true
			if errors.Is(err, histogram.ErrHistogramsIncompatibleSchema) {
				return annotations.MixedExponentialCustomHistogramsWarning
			}
			if errors.Is(err, histogram.ErrHistogramsIncompatibleBounds) {
				return annotations.IncompatibleCustomBucketsHistogramsWarning
			}
			return err
		}
		a.histSum, err = a.histSum.Add(toAdd)
		if err != nil {
			a.histSum = nil
			a.histCount = 0
			a.hasHistError = true
			if errors.Is(err, histogram.ErrHistogramsIncompatibleSchema) {
				return annotations.MixedExponentialCustomHistogramsWarning
			}
			if errors.Is(err, histogram.ErrHistogramsIncompatibleBounds) {
				return annotations.IncompatibleCustomBucketsHistogramsWarning
			}
			return err
		}
		return nil
	}

	a.count++
	if !a.hasValue {
		a.hasValue = true
		a.kahanSum = v
		return nil
	}

	a.hasValue = true

	if !a.incremental {
		newSum, newC := KahanSumInc(v, a.kahanSum, a.kahanC)

		if !math.IsInf(newSum, 0) {
			// The sum doesn't overflow, so we propagate it to the
			// group struct and continue with the regular
			// calculation of the mean value.
			a.kahanSum, a.kahanC = newSum, newC
			return nil
		}

		// If we are here, we know that the sum _would_ overflow. So
		// instead of continue to sum up, we revert to incremental
		// calculation of the mean value from here on.
		a.incremental = true
		a.avg = a.kahanSum / float64(a.count-1)
		a.kahanC /= float64(a.count) - 1
	}

	if math.IsInf(a.avg, 0) {
		if math.IsInf(v, 0) && (a.avg > 0) == (v > 0) {
			// The `floatMean` and `s.F` values are `Inf` of the same sign.  They
			// can't be subtracted, but the value of `floatMean` is correct
			// already.
			return nil
		}
		if !math.IsInf(v, 0) && !math.IsNaN(v) {
			// At this stage, the mean is an infinite. If the added
			// value is neither an Inf or a Nan, we can keep that mean
			// value.
			// This is required because our calculation below removes
			// the mean value, which would look like Inf += x - Inf and
			// end up as a NaN.
			return nil
		}
	}
	currentMean := a.avg + a.kahanC
	a.avg, a.kahanC = KahanSumInc(
		// Divide each side of the `-` by `group.groupCount` to avoid float64 overflows.
		v/float64(a.count)-currentMean/float64(a.count),
		a.avg,
		a.kahanC,
	)
	return nil
}

func (a *AvgAcc) AddVector(vs []float64, hs []*histogram.FloatHistogram) error {
	for _, v := range vs {
		if err := a.Add(v, nil); err != nil {
			return err
		}
	}
	for _, h := range hs {
		if err := a.Add(0, h); err != nil {
			if errors.Is(err, histogram.ErrHistogramsIncompatibleSchema) {
				// to make valueType NoValue
				a.histSum = nil
				a.histCount = 0
				return annotations.MixedExponentialCustomHistogramsWarning
			}
			if errors.Is(err, histogram.ErrHistogramsIncompatibleBounds) {
				// to make valueType NoValue
				a.histSum = nil
				a.histCount = 0
				return annotations.IncompatibleCustomBucketsHistogramsWarning
			}
			return err
		}
	}
	return nil
}

func (a *AvgAcc) Value() (float64, *histogram.FloatHistogram) {
	if a.histSum != nil {
		a.histSum.Compact(0)
	}
	if a.incremental {
		return a.avg + a.kahanC, a.histSum
	}
	return (a.kahanSum + a.kahanC) / float64(a.count), a.histSum
}

func (a *AvgAcc) ValueType() ValueType {
	hasFloat := a.count > 0
	hasHist := a.histCount > 0

	if hasFloat && hasHist {
		return MixedTypeValue
	}
	if hasFloat || hasHist {
		return SingleTypeValue
	}
	return NoValue
}

func (a *AvgAcc) Reset(_ float64) {
	a.hasValue = false
	a.incremental = false
	a.kahanSum = 0
	a.kahanC = 0
	a.count = 0

	a.histCount = 0
	a.histSum = nil
	a.hasHistError = false
}

func (a *AvgAcc) Copy() CopyableAccumulator {
	cp := &AvgAcc{
		kahanSum:     a.kahanSum,
		kahanC:       a.kahanC,
		avg:          a.avg,
		incremental:  a.incremental,
		count:        a.count,
		hasValue:     a.hasValue,
		histCount:    a.histCount,
		hasHistError: a.hasHistError,
	}
	if a.histSum != nil {
		cp.histSum = a.histSum.Copy()
	}
	if a.histScratch != nil {
		cp.histScratch = a.histScratch.Copy()
	}
	if a.histSumScratch != nil {
		cp.histSumScratch = a.histSumScratch.Copy()
	}
	return cp
}

func (a *AvgAcc) Merge(other MergeableAccumulator) error {
	o := other.(*AvgAcc)

	// Merge float values: combined_avg = (sum1 + sum2) / (count1 + count2)
	if o.count > 0 {
		// Get the sum from other accumulator
		var otherSum float64
		if o.incremental {
			otherSum = o.avg * float64(o.count)
		} else {
			otherSum = o.kahanSum + o.kahanC
		}

		if a.count == 0 {
			// This accumulator is empty, just copy other's state
			a.kahanSum = otherSum
			a.kahanC = 0
			a.count = o.count
			a.hasValue = true
			a.incremental = false
		} else {
			// Get our sum
			var ourSum float64
			if a.incremental {
				ourSum = a.avg * float64(a.count)
			} else {
				ourSum = a.kahanSum + a.kahanC
			}

			// Combine sums and counts
			totalCount := a.count + o.count
			totalSum := ourSum + otherSum

			// Check for overflow and switch to incremental if needed
			if math.IsInf(totalSum, 0) {
				a.incremental = true
				a.avg = totalSum / float64(totalCount)
				a.kahanC = 0
			} else {
				a.incremental = false
				a.kahanSum = totalSum
				a.kahanC = 0
			}
			a.count = totalCount
		}
	}

	// Merge histogram values
	if o.histCount > 0 && !a.hasHistError && !o.hasHistError {
		if a.histSum == nil {
			a.histSum = o.histSum.Copy()
			a.histCount = o.histCount
		} else {
			// For histogram average, we need to combine: (histSum1 + histSum2) / (count1 + count2)
			// But histSum stores the running average, not the sum. We need to recover sums first.
			// Actually looking at the Add code, histSum IS the running average, updated incrementally.
			// To merge, we compute: new_avg = (avg1 * count1 + avg2 * count2) / (count1 + count2)
			totalCount := a.histCount + o.histCount
			// Scale current average by its weight
			scaled1 := a.histSum.Copy().Mul(a.histCount / totalCount)
			scaled2 := o.histSum.Copy().Mul(o.histCount / totalCount)
			var err error
			a.histSum, err = scaled1.Add(scaled2)
			if err != nil {
				a.histSum = nil
				a.histCount = 0
				a.hasHistError = true
				return err
			}
			a.histCount = totalCount
		}
	}

	return nil
}

func (a *AvgAcc) Sub(v float64, h *histogram.FloatHistogram) error {
	if h != nil {
		// Subtract histogram from average
		if a.histSum != nil && !a.hasHistError && a.histCount > 0 {
			// Recover sum: sum = avg * count (histSum stores running average)
			// After subtraction: new_avg = (sum - h) / (count - 1)
			// But histSum is the average, so we need to adjust
			a.histCount--
			if a.histCount == 0 {
				a.histSum = nil
			} else {
				// new_sum = old_avg * old_count - h
				// new_avg = new_sum / new_count
				// This is complex for histograms, fall back to not supporting it
				// for now by invalidating
				a.histSum = nil
				a.hasHistError = true
			}
		}
		return nil
	}

	if a.count <= 0 {
		return nil
	}

	// Get current sum
	var currentSum float64
	if a.incremental {
		currentSum = a.avg * float64(a.count)
	} else {
		currentSum = a.kahanSum + a.kahanC
	}

	// Subtract value and decrement count
	newSum := currentSum - v
	a.count--

	if a.count == 0 {
		a.hasValue = false
		a.kahanSum = 0
		a.kahanC = 0
		a.avg = 0
		a.incremental = false
	} else {
		// Store as sum (non-incremental) for simplicity
		a.incremental = false
		a.kahanSum = newSum
		a.kahanC = 0
	}

	return nil
}

type statAcc struct {
	count    float64
	mean     float64
	cMean    float64
	value    float64
	cValue   float64
	hasValue bool
	hasNaN   bool
}

func (s *statAcc) ValueType() ValueType {
	if s.hasValue {
		return SingleTypeValue
	}
	return NoValue
}

func (s *statAcc) Reset(_ float64) {
	s.hasValue = false
	s.hasNaN = false
	s.count = 0
	s.mean = 0
	s.cMean = 0
	s.value = 0
	s.cValue = 0
}

func (s *statAcc) add(v float64) {
	s.hasValue = true
	s.count++
	if math.IsNaN(v) || math.IsInf(v, 0) {
		s.hasNaN = true
		return
	}
	delta := v - (s.mean + s.cMean)
	s.mean, s.cMean = KahanSumInc(delta/s.count, s.mean, s.cMean)
	s.value, s.cValue = KahanSumInc(delta*(v-(s.mean+s.cMean)), s.value, s.cValue)
}

func (s *statAcc) variance() float64 {
	if s.hasNaN {
		return math.NaN()
	}
	return (s.value + s.cValue) / s.count
}

type StdDevAcc struct {
	statAcc
}

func NewStdDevAcc() *StdDevAcc {
	return &StdDevAcc{}
}

func (s *StdDevAcc) Add(v float64, h *histogram.FloatHistogram) error {
	if h != nil {
		return annotations.NewHistogramIgnoredInAggregationInfo("stddev", posrange.PositionRange{})
	}
	s.add(v)
	return nil
}

func (s *StdDevAcc) Value() (float64, *histogram.FloatHistogram) {
	return math.Sqrt(s.variance()), nil
}

type StdVarAcc struct {
	statAcc
}

func NewStdVarAcc() *StdVarAcc {
	return &StdVarAcc{}
}

func (s *StdVarAcc) Add(v float64, h *histogram.FloatHistogram) error {
	if h != nil {
		return annotations.NewHistogramIgnoredInAggregationInfo("stdvar", posrange.PositionRange{})
	}
	s.add(v)
	return nil
}

func (s *StdVarAcc) Value() (float64, *histogram.FloatHistogram) {
	return s.variance(), nil
}

type statOverTimeAcc struct {
	statAcc
	seenHist   bool
	warnedHist bool
}

func (s *statOverTimeAcc) Reset(_ float64) {
	s.statAcc.Reset(0)
	s.seenHist = false
	s.warnedHist = false
}

func (s *statOverTimeAcc) addWithHist(v float64, h *histogram.FloatHistogram) error {
	if h != nil {
		s.seenHist = true
		if s.hasValue && !s.warnedHist {
			s.warnedHist = true
			return annotations.NewHistogramIgnoredInMixedRangeInfo("", posrange.PositionRange{})
		}
		return nil
	}

	if s.seenHist && !s.warnedHist {
		s.warnedHist = true
	}
	s.add(v)
	if s.seenHist && s.warnedHist {
		return annotations.NewHistogramIgnoredInMixedRangeInfo("", posrange.PositionRange{})
	}
	return nil
}

type StdDevOverTimeAcc struct {
	statOverTimeAcc
}

func NewStdDevOverTimeAcc() *StdDevOverTimeAcc {
	return &StdDevOverTimeAcc{}
}

func (s *StdDevOverTimeAcc) Add(v float64, h *histogram.FloatHistogram) error {
	return s.addWithHist(v, h)
}

func (s *StdDevOverTimeAcc) Value() (float64, *histogram.FloatHistogram) {
	return math.Sqrt(s.variance()), nil
}

type StdVarOverTimeAcc struct {
	statOverTimeAcc
}

func NewStdVarOverTimeAcc() *StdVarOverTimeAcc {
	return &StdVarOverTimeAcc{}
}

func (s *StdVarOverTimeAcc) Add(v float64, h *histogram.FloatHistogram) error {
	return s.addWithHist(v, h)
}

func (s *StdVarOverTimeAcc) Value() (float64, *histogram.FloatHistogram) {
	return s.variance(), nil
}

type QuantileAcc struct {
	arg      float64
	points   []float64
	hasValue bool
}

func NewQuantileAcc() Accumulator {
	return &QuantileAcc{}
}

func (q *QuantileAcc) Add(v float64, h *histogram.FloatHistogram) error {
	if h != nil {
		return annotations.NewHistogramIgnoredInAggregationInfo("quantile", posrange.PositionRange{})
	}

	q.hasValue = true
	q.points = append(q.points, v)
	return nil
}

func (q *QuantileAcc) Value() (float64, *histogram.FloatHistogram) {
	return Quantile(q.arg, q.points), nil
}

func (q *QuantileAcc) ValueType() ValueType {
	if q.hasValue {
		return SingleTypeValue
	} else {
		return NoValue
	}
}

func (q *QuantileAcc) Reset(f float64) {
	q.hasValue = false
	q.arg = f
	q.points = q.points[:0]
}

type HistogramAvgAcc struct {
	sum      *histogram.FloatHistogram
	count    int64
	hasFloat bool
}

func NewHistogramAvgAcc() *HistogramAvgAcc {
	return &HistogramAvgAcc{
		sum: &histogram.FloatHistogram{},
	}
}

func (acc *HistogramAvgAcc) Add(v float64, h *histogram.FloatHistogram) error {
	if h == nil {
		acc.hasFloat = true
	}
	if acc.count == 0 {
		h.CopyTo(acc.sum)
	}
	var err error
	if h.Schema >= acc.sum.Schema {
		if acc.sum, err = acc.sum.Add(h); err != nil {
			return err
		}
	} else {
		t := h.Copy()
		if _, err = t.Add(acc.sum); err != nil {
			return err
		}
		acc.sum = t
	}
	acc.count++
	return nil
}

func (acc *HistogramAvgAcc) Value() (float64, *histogram.FloatHistogram) {
	return 0, acc.sum.Mul(1 / float64(acc.count))
}

func (acc *HistogramAvgAcc) ValueType() ValueType {
	if acc.count > 0 && !acc.hasFloat {
		return SingleTypeValue
	}
	return NoValue
}

func (acc *HistogramAvgAcc) Reset(f float64) {
	acc.count = 0
}

// KahanSumInc implements kahan summation, see https://en.wikipedia.org/wiki/Kahan_summation_algorithm.
func KahanSumInc(inc, sum, c float64) (newSum, newC float64) {
	t := sum + inc
	switch {
	case math.IsInf(t, 0):
		c = 0

	// Using Neumaier improvement, swap if next term larger than sum.
	case math.Abs(sum) >= math.Abs(inc):
		c += (sum - t) + inc
	default:
		c += (inc - t) + sum
	}
	return t, c
}

func Quantile(q float64, points []float64) float64 {
	if len(points) == 0 || math.IsNaN(q) {
		return math.NaN()
	}
	if q < 0 {
		return math.Inf(-1)
	}
	if q > 1 {
		return math.Inf(+1)
	}
	sort.Float64s(points)

	n := float64(len(points))
	// When the quantile lies between two samples,
	// we use a weighted average of the two samples.
	rank := q * (n - 1)

	lowerIndex := math.Max(0, math.Floor(rank))
	upperIndex := math.Min(n-1, lowerIndex+1)

	weight := rank - math.Floor(rank)
	return points[int(lowerIndex)]*(1-weight) + points[int(upperIndex)]*weight
}

func histogramSum(current *histogram.FloatHistogram, histograms []*histogram.FloatHistogram) (*histogram.FloatHistogram, error) {
	if len(histograms) == 0 {
		return current, nil
	}
	if current == nil && len(histograms) == 1 {
		return histograms[0].Copy(), nil
	}
	var histSum *histogram.FloatHistogram
	if current != nil {
		histSum = current.Copy()
	} else {
		histSum = histograms[0].Copy()
		histograms = histograms[1:]
	}

	var err error
	for i := range histograms {
		if histograms[i].Schema >= histSum.Schema {
			histSum, err = histSum.Add(histograms[i])
			if err != nil {
				if errors.Is(err, histogram.ErrHistogramsIncompatibleSchema) {
					return nil, annotations.MixedExponentialCustomHistogramsWarning
				}
				if errors.Is(err, histogram.ErrHistogramsIncompatibleBounds) {
					return nil, annotations.IncompatibleCustomBucketsHistogramsWarning
				}
				return nil, err
			}
		} else {
			t := histograms[i].Copy()
			if histSum, err = t.Add(histSum); err != nil {
				if errors.Is(err, histogram.ErrHistogramsIncompatibleSchema) {
					return nil, annotations.NewMixedExponentialCustomHistogramsWarning("", posrange.PositionRange{})
				}
				if errors.Is(err, histogram.ErrHistogramsIncompatibleBounds) {
					return nil, annotations.NewIncompatibleCustomBucketsHistogramsWarning("", posrange.PositionRange{})
				}
				return nil, err
			}
		}
	}
	return histSum, nil
}

// compensatedSum returns the sum of the elements of the slice calculated with greater
// accuracy than Sum at the expense of additional computation.
func compensatedSum(s []float64) float64 {
	// compensatedSum uses an improved version of Kahan's compensated
	// summation algorithm proposed by Neumaier.
	// See https://en.wikipedia.org/wiki/Kahan_summation_algorithm for details.
	var sum, c float64
	for _, x := range s {
		// This type conversion is here to prevent a sufficiently smart compiler
		// from optimizing away these operations.
		t := sum + x
		switch {
		case math.IsInf(t, 0):
			c = 0

		// Using Neumaier improvement, swap if next term larger than sum.
		case math.Abs(sum) >= math.Abs(x):
			c += (sum - t) + x
		default:
			c += (x - t) + sum
		}
		sum = t
	}
	return sum + c
}
