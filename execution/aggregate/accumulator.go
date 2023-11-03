package aggregate

import (
	"math"

	"github.com/prometheus/prometheus/model/histogram"
)

type accumulator interface {
	Add(v float64, h *histogram.FloatHistogram)
	Value() (float64, *histogram.FloatHistogram)
	HasValue() bool
	Reset(float64)
}

type newAccumulatorFunc func() accumulator

type sumAcc struct {
	value       float64
	histSum     *histogram.FloatHistogram
	hasFloatVal bool
}

func newSumAcc() accumulator {
	return &sumAcc{}
}

func (s *sumAcc) Add(v float64, h *histogram.FloatHistogram) {
	if h == nil {
		s.hasFloatVal = true
		s.value += v
		return
	}
	if s.histSum == nil {
		s.histSum = h.Copy()
		return
	}
	// The histogram being added must have an equal or larger schema.
	// https://github.com/prometheus/prometheus/blob/57bcbf18880f7554ae34c5b341d52fc53f059a97/promql/engine.go#L2448-L2456
	if h.Schema >= s.histSum.Schema {
		s.histSum = s.histSum.Add(h)
	} else {
		t := h.Copy()
		t.Add(s.histSum)
		s.histSum = t
	}
}

func (s *sumAcc) Value() (float64, *histogram.FloatHistogram) {
	return s.value, s.histSum
}

// HasValue for sum returns an empty result when floats are histograms are aggregated.
func (s *sumAcc) HasValue() bool {
	return s.hasFloatVal != (s.histSum != nil)
}

func (s *sumAcc) Reset(_ float64) {
	s.histSum = nil
	s.hasFloatVal = false
	s.value = 0
}

type genericAcc struct {
	value     float64
	hasValue  bool
	aggregate func(float64, float64) float64
}

func maxAggregate(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
func minAggregate(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
func groupAggregate(_, _ float64) float64 { return 1 }

func newMaxAcc() accumulator {
	return &genericAcc{aggregate: maxAggregate}
}

func newMinAcc() accumulator {
	return &genericAcc{aggregate: minAggregate}
}

func newCountAcc() accumulator {
	return &countAcc{}
}

func newGroupAcc() accumulator {
	return &genericAcc{aggregate: groupAggregate}
}

func (g *genericAcc) Add(v float64, _ *histogram.FloatHistogram) {
	if !g.hasValue || math.IsNaN(g.value) {
		g.value = v
	}
	g.hasValue = true
	g.value = g.aggregate(g.value, v)
}

func (g *genericAcc) Value() (float64, *histogram.FloatHistogram) {
	return g.value, nil
}

func (g *genericAcc) HasValue() bool {
	return g.hasValue
}

func (g *genericAcc) Reset(_ float64) {
	g.hasValue = false
	g.value = 0
}

type countAcc struct {
	value    float64
	hasValue bool
}

func (c *countAcc) Add(v float64, h *histogram.FloatHistogram) {
	c.hasValue = true
	c.value += 1
}

func (c *countAcc) Value() (float64, *histogram.FloatHistogram) {
	return c.value, nil
}

func (c *countAcc) HasValue() bool {
	return c.hasValue
}

func (c *countAcc) Reset(_ float64) {
	c.hasValue = false
	c.value = 0
}

type avgAcc struct {
	count    float64
	sum      float64
	hasValue bool
}

func newAvgAcc() accumulator {
	return &avgAcc{}
}

func (a *avgAcc) Add(v float64, h *histogram.FloatHistogram) {
	a.hasValue = true
	a.count += 1
	a.sum += v
}

func (a *avgAcc) Value() (float64, *histogram.FloatHistogram) {
	return a.sum / a.count, nil
}

func (a *avgAcc) HasValue() bool {
	return a.hasValue
}

func (a *avgAcc) Reset(_ float64) {
	a.hasValue = false
	a.sum = 0
	a.count = 0
}

type statAcc struct {
	count    float64
	mean     float64
	value    float64
	hasValue bool
}

func (s *statAcc) Add(v float64, h *histogram.FloatHistogram) {
	s.hasValue = true
	s.count++

	delta := v - s.mean
	s.mean += delta / s.count
	s.value += delta * (v - s.mean)
}

func (s *statAcc) HasValue() bool {
	return s.hasValue
}

func (s *statAcc) Reset(_ float64) {
	s.hasValue = false
	s.count = 0
	s.mean = 0
	s.value = 0
}

type stdDevAcc struct {
	statAcc
}

func newStdDevAcc() accumulator {
	return &stdDevAcc{}
}

func (s *stdDevAcc) Value() (float64, *histogram.FloatHistogram) {
	if s.count == 1 {
		return 0, nil
	}
	return math.Sqrt(s.value / s.count), nil
}

type stdVarAcc struct {
	statAcc
}

func newStdVarAcc() accumulator {
	return &stdVarAcc{}
}

func (s *stdVarAcc) Value() (float64, *histogram.FloatHistogram) {
	if s.count == 1 {
		return 0, nil
	}
	return s.value / s.count, nil
}

type quantileAcc struct {
	arg      float64
	points   []float64
	hasValue bool
}

func newQuantileAcc() accumulator {
	return &quantileAcc{}
}

func (q *quantileAcc) Add(v float64, h *histogram.FloatHistogram) {
	q.hasValue = true
	q.points = append(q.points, v)
}

func (q *quantileAcc) Value() (float64, *histogram.FloatHistogram) {
	return quantile(q.arg, q.points), nil
}

func (q *quantileAcc) HasValue() bool {
	return q.hasValue
}

func (q *quantileAcc) Reset(f float64) {
	q.hasValue = false
	q.arg = f
	q.points = q.points[:0]
}
