// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package model

import (
	"context"

	"github.com/prometheus/prometheus/model/labels"
)

type VectorOperator interface {
	Next(ctx context.Context) ([]StepVector, error)
	Series(ctx context.Context) ([]labels.Labels, error)
	GetPool() *VectorPool

	Explain() (me string, next []VectorOperator)
}
