// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/thanos-io/promql-engine/execution/model"

	"github.com/prometheus/prometheus/model/histogram"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/util/stats"
)

type OperatorTelemetry interface {
	fmt.Stringer

	MaxSeriesCount() int
	SetMaxSeriesCount(count int)
	ExecutionTimeTaken() time.Duration
	AddSeriesExecutionTime(time.Duration)
	SeriesExecutionTime() time.Duration
	AddNextExecutionTime(time.Duration)
	NextExecutionTime() time.Duration
	IncrementSamplesAtTimestamp(samples int, t int64) error
	Samples() *stats.QuerySamples
	IsSubquery() bool
	IsStepInvariant() bool
	UpdatePeak(count int)
}

func NewTelemetry(operator fmt.Stringer, enableAnalysis, enablePerStepStats bool, start, end int64, step time.Duration, limiter *SampleLimiter) OperatorTelemetry {
	if enableAnalysis {
		return newTrackedTelemetry(operator, enablePerStepStats, start, end, step, limiter, false, false)
	}
	return &NoopTelemetry{Stringer: operator, limiter: limiter}
}

func NewSubqueryTelemetry(operator fmt.Stringer, enableAnalysis, enablePerStepStats bool, start, end int64, step time.Duration, limiter *SampleLimiter) OperatorTelemetry {
	if enableAnalysis {
		return newTrackedTelemetry(operator, enablePerStepStats, start, end, step, limiter, true, false)
	}
	return &NoopTelemetry{Stringer: operator, limiter: limiter, isSubquery: true}
}

func NewStepInvariantTelemetry(operator fmt.Stringer, enableAnalysis, enablePerStepStats bool, start, end int64, step time.Duration, limiter *SampleLimiter) OperatorTelemetry {
	if enableAnalysis {
		return newTrackedTelemetry(operator, enablePerStepStats, start, end, step, limiter, false, true)
	}
	return &NoopTelemetry{Stringer: operator, limiter: limiter, isStepInvariant: true}
}

type NoopTelemetry struct {
	fmt.Stringer
	limiter         *SampleLimiter
	isSubquery      bool
	isStepInvariant bool
}

func (tm *NoopTelemetry) AddExecutionTimeTaken(t time.Duration) {}

func (tm *NoopTelemetry) ExecutionTimeTaken() time.Duration {
	return time.Duration(0)
}

func (tm *NoopTelemetry) AddSeriesExecutionTime(t time.Duration) {}

func (tm *NoopTelemetry) SeriesExecutionTime() time.Duration {
	return time.Duration(0)
}

func (tm *NoopTelemetry) AddNextExecutionTime(t time.Duration) {}

func (tm *NoopTelemetry) NextExecutionTime() time.Duration {
	return time.Duration(0)
}

func (tm *NoopTelemetry) IncrementSamplesAtTimestamp(samples int, t int64) error {
	return tm.limiter.Add(samples, t)
}

func (tm *NoopTelemetry) Samples() *stats.QuerySamples { return nil }

func (tm *NoopTelemetry) MaxSeriesCount() int { return 0 }

func (tm *NoopTelemetry) SetMaxSeriesCount(_ int) {}

func (tm *NoopTelemetry) IsSubquery() bool { return tm.isSubquery }

func (tm *NoopTelemetry) IsStepInvariant() bool { return tm.isStepInvariant }

func (tm *NoopTelemetry) UpdatePeak(_ int) {}

type TrackedTelemetry struct {
	fmt.Stringer

	Series          int
	ExecutionTime   time.Duration
	SeriesTime      time.Duration
	NextTime        time.Duration
	LoadedSamples   *stats.QuerySamples
	limiter         *SampleLimiter
	isSubquery      bool
	isStepInvariant bool
}

func newTrackedTelemetry(operator fmt.Stringer, enablePerStepStats bool, start, end int64, step time.Duration, limiter *SampleLimiter, isSubquery, isStepInvariant bool) *TrackedTelemetry {
	ss := stats.NewQuerySamples(enablePerStepStats)
	ss.InitStepTracking(start, end, StepTrackingInterval(step))
	return &TrackedTelemetry{
		Stringer:        operator,
		LoadedSamples:   ss,
		limiter:         limiter,
		isSubquery:      isSubquery,
		isStepInvariant: isStepInvariant,
	}
}

func StepTrackingInterval(step time.Duration) int64 {
	if step == 0 {
		return 1
	}
	return int64(step / (time.Millisecond / time.Nanosecond))
}

func (ti *TrackedTelemetry) AddExecutionTimeTaken(t time.Duration) { ti.ExecutionTime += t }

func (ti *TrackedTelemetry) ExecutionTimeTaken() time.Duration {
	return ti.ExecutionTime
}

func (ti *TrackedTelemetry) AddSeriesExecutionTime(t time.Duration) {
	ti.SeriesTime += t
	ti.ExecutionTime += t
}

func (ti *TrackedTelemetry) SeriesExecutionTime() time.Duration {
	return ti.SeriesTime
}

func (ti *TrackedTelemetry) AddNextExecutionTime(t time.Duration) {
	ti.NextTime += t
	ti.ExecutionTime += t
}

func (ti *TrackedTelemetry) NextExecutionTime() time.Duration {
	return ti.NextTime
}

func (ti *TrackedTelemetry) IncrementSamplesAtTimestamp(samples int, t int64) error {
	ti.LoadedSamples.IncrementSamplesAtTimestamp(t, int64(samples))
	return ti.limiter.Add(samples, t)
}

func (ti *TrackedTelemetry) IsSubquery() bool { return ti.isSubquery }

func (ti *TrackedTelemetry) IsStepInvariant() bool { return ti.isStepInvariant }

func (ti *TrackedTelemetry) Samples() *stats.QuerySamples { return ti.LoadedSamples }

func (ti *TrackedTelemetry) MaxSeriesCount() int { return ti.Series }

func (ti *TrackedTelemetry) SetMaxSeriesCount(count int) { ti.Series = count }

func (ti *TrackedTelemetry) UpdatePeak(count int) {
	ti.Samples().UpdatePeak(count)
}

type ObservableVectorOperator interface {
	model.VectorOperator
	OperatorTelemetry
}

// CalculateHistogramSampleCount returns the size of the FloatHistogram compared to the size of a Float.
// The total size is calculated considering the histogram timestamp (p.T - 8 bytes),
// and then a number of bytes in the histogram.
// This sum is divided by 16, as samples are 16 bytes.
// See: https://github.com/prometheus/prometheus/blob/2bf6f4c9dcbb1ad2e8fef70c6a48d8fc44a7f57c/promql/value.go#L178
func CalculateHistogramSampleCount(h *histogram.FloatHistogram) int {
	return (h.Size() + 8) / 16
}

func NewOperator(telemetry OperatorTelemetry, inner model.VectorOperator) model.VectorOperator {
	op := &Operator{
		inner: inner,
	}
	op.OperatorTelemetry = telemetry
	return op
}

// Operator wraps other inner operator to track its telemetry.
type Operator struct {
	OperatorTelemetry
	inner model.VectorOperator
}

func (t *Operator) Series(ctx context.Context) ([]labels.Labels, error) {
	start := time.Now()
	defer func() { t.OperatorTelemetry.AddSeriesExecutionTime(time.Since(start)) }()
	s, err := t.inner.Series(ctx)
	if err != nil {
		return nil, err
	}
	t.OperatorTelemetry.SetMaxSeriesCount(len(s))
	return s, err
}

func (t *Operator) Next(ctx context.Context, buf []model.StepVector) (int, error) {
	start := time.Now()
	var totalSamplesBeforeCount int64
	totalSamplesBefore := t.OperatorTelemetry.Samples()
	if totalSamplesBefore != nil {
		totalSamplesBeforeCount = totalSamplesBefore.TotalSamples
	} else {
		totalSamplesBeforeCount = 0
	}

	defer func() { t.OperatorTelemetry.AddNextExecutionTime(time.Since(start)) }()
	n, err := t.inner.Next(ctx, buf)
	if err != nil {
		return 0, err
	}

	var totalSamplesAfter int64
	totalSamplesAfterSamples := t.OperatorTelemetry.Samples()
	if totalSamplesAfterSamples != nil {
		totalSamplesAfter = totalSamplesAfterSamples.TotalSamples
	} else {
		totalSamplesAfter = 0
	}

	t.OperatorTelemetry.UpdatePeak(int(totalSamplesAfter) - int(totalSamplesBeforeCount))

	return n, err
}

func (t *Operator) Explain() []model.VectorOperator {
	return t.inner.Explain()
}

func (t *Operator) String() string {
	return t.inner.String()
}
