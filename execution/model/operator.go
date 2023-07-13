// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package model

import (
	"context"
	"time"

	"github.com/prometheus/prometheus/model/labels"
)

type NoopTimingInformation struct{}

func (ti *NoopTimingInformation) AddCPUTimeTaken(t time.Duration) {}

type TimingInformation struct {
	CPUTime time.Duration
}

func (ti *TimingInformation) AddCPUTimeTaken(t time.Duration) {
	ti.CPUTime += t
}

type OperatorTelemetry interface {
	AddCPUTimeTaken(time.Duration)
	CPUTimeTaken() time.Duration
}

func (ti *NoopTimingInformation) CPUTimeTaken() time.Duration {
	return time.Duration(0)
}
func (ti *TimingInformation) CPUTimeTaken() time.Duration {
	return ti.CPUTime
}

type ObservableVectorOperator interface {
	VectorOperator
	OperatorTelemetry
	Analyze() (OperatorTelemetry, []ObservableVectorOperator)
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
	Explain() (me string, next []VectorOperator)
}
