// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package model

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/util/stats"

	"github.com/thanos-io/promql-engine/query"
)

type metadataKey string

const key metadataKey = "promql-metadata"

func AddMetadataStorage(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, &FanoutMetadata{
		storage: make(map[string]map[string]any),
	})
}

func GetMetadataStorage(ctx context.Context) *FanoutMetadata {
	v := ctx.Value(key)
	if v == nil {
		return nil
	}
	return v.(*FanoutMetadata)
}

type FanoutMetadata struct {
	m sync.Mutex

	storage map[string]map[string]any
}

func (m *FanoutMetadata) SetMetadata(endpoint string, data map[string]any) {
	m.m.Lock()
	defer m.m.Unlock()

	m.storage[endpoint] = data
}

func (m *FanoutMetadata) GetMetadata() map[string]map[string]any {
	m.m.Lock()
	defer m.m.Unlock()

	return m.storage
}

type OperatorTelemetry interface {
	fmt.Stringer

	AddExecutionTimeTaken(time.Duration)
	ExecutionTimeTaken() time.Duration
	IncrementSamplesAtTimestamp(samples int, t int64)
	Samples() *stats.QuerySamples
	SubQuery() bool

	FanoutMetadata() *FanoutMetadata
}

func NewTelemetry(ctx context.Context, operator fmt.Stringer, opts *query.Options) OperatorTelemetry {
	if opts.EnableAnalysis {
		return NewTrackedTelemetry(ctx, operator, opts, false)
	}
	return NewNoopTelemetry(operator)
}

func NewSubqueryTelemetry(ctx context.Context, operator fmt.Stringer, opts *query.Options) OperatorTelemetry {
	if opts.EnableAnalysis {
		return NewTrackedTelemetry(ctx, operator, opts, true)
	}
	return NewNoopTelemetry(operator)
}

type NoopTelemetry struct {
	fmt.Stringer
}

func NewNoopTelemetry(operator fmt.Stringer) *NoopTelemetry {
	return &NoopTelemetry{Stringer: operator}
}

func (tm *NoopTelemetry) AddExecutionTimeTaken(t time.Duration) {}

func (tm *NoopTelemetry) ExecutionTimeTaken() time.Duration {
	return time.Duration(0)
}

func (tm *NoopTelemetry) FanoutMetadata() *FanoutMetadata {
	return nil
}

func (tm *NoopTelemetry) IncrementSamplesAtTimestamp(_ int, _ int64) {}

func (tm *NoopTelemetry) Samples() *stats.QuerySamples { return nil }
func (tm *NoopTelemetry) SubQuery() bool               { return false }

type TrackedTelemetry struct {
	fmt.Stringer

	ExecutionTime time.Duration
	LoadedSamples *stats.QuerySamples
	subquery      bool
	ctx           context.Context
}

func NewTrackedTelemetry(ctx context.Context, operator fmt.Stringer, opts *query.Options, subquery bool) *TrackedTelemetry {
	ss := stats.NewQuerySamples(opts.EnablePerStepStats)
	ss.InitStepTracking(opts.Start.UnixMilli(), opts.End.UnixMilli(), stepTrackingInterval(opts.Step))
	return &TrackedTelemetry{
		Stringer:      operator,
		LoadedSamples: ss,
		subquery:      subquery,
		ctx:           ctx,
	}
}

func stepTrackingInterval(step time.Duration) int64 {
	if step == 0 {
		return 1
	}
	return int64(step / (time.Millisecond / time.Nanosecond))
}

func (ti *TrackedTelemetry) AddExecutionTimeTaken(t time.Duration) { ti.ExecutionTime += t }

func (ti *TrackedTelemetry) ExecutionTimeTaken() time.Duration {
	return ti.ExecutionTime
}

func (ti *TrackedTelemetry) FanoutMetadata() *FanoutMetadata {
	return GetMetadataStorage(ti.ctx)
}

func (ti *TrackedTelemetry) IncrementSamplesAtTimestamp(samples int, t int64) {
	ti.updatePeak(samples)
	ti.LoadedSamples.IncrementSamplesAtTimestamp(t, int64(samples))
}

func (ti *TrackedTelemetry) SubQuery() bool {
	return ti.subquery
}

func (ti *TrackedTelemetry) updatePeak(samples int) {
	ti.LoadedSamples.UpdatePeak(samples)
}

func (ti *TrackedTelemetry) Samples() *stats.QuerySamples { return ti.LoadedSamples }

type ObservableVectorOperator interface {
	VectorOperator
	OperatorTelemetry
}

// VectorOperator performs operations on series in step by step fashion.
type VectorOperator interface {
	// Next yields vectors of samples from all series for one or more execution steps.
	Next(ctx context.Context) ([]StepVector, error)

	// Series returns all series that the operator will process during Next results.
	// The result can be used by upstream operators to allocate output tables and buffers
	// before starting to process samples.
	Series(ctx context.Context) ([]labels.Labels, error)

	// GetPool returns pool of vectors that can be shared across operators.
	GetPool() *VectorPool

	// Explain returns human-readable explanation of the current operator and optional nested operators.
	Explain() (next []VectorOperator)

	fmt.Stringer
}
