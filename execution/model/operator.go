// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package model

import (
	"context"

	"github.com/prometheus/prometheus/model/labels"
)

// VectorOperator performs operations on series in step by step fashion.
type VectorOperator interface {
	// Next yields stream of samples representing series (vector of multiple series) per one or more steps.
	Next(ctx context.Context) ([]StepVector, error)

	// Series returns all series that we will see in all Next results.
	Series(ctx context.Context) ([]labels.Labels, error)

	// GetPool returns pool of vectors that can be shared across operators.
	GetPool() *VectorPool

	// Explain returns human-readable explanation of the current operator and optional nested operators.
	Explain() (me string, next []VectorOperator)
}
