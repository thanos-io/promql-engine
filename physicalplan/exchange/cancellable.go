// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package exchange

import (
	"context"

	"github.com/prometheus/prometheus/model/labels"

	"github.com/thanos-community/promql-engine/physicalplan/model"
)

// TODO(bwplotka): Consider removing this and ensuring all operators check for context. It creates
// unnecessary cognitive and computation load, only to avoid one "if ctx.Err() != nil" in each operator.
// It is also inconsistently added (sometimes missing).
type CancellableOperator struct {
	next model.VectorOperator
}

func NewCancellable(next model.VectorOperator) *CancellableOperator {
	return &CancellableOperator{next: next}
}

func (c *CancellableOperator) Explain() (string, []model.VectorOperator) {
	return "[*CancellableOperator]", []model.VectorOperator{c.next}
}

func (c *CancellableOperator) Next(ctx context.Context) ([]model.StepVector, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return c.next.Next(ctx)
	}
}

func (c *CancellableOperator) Series(ctx context.Context) ([]labels.Labels, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return c.next.Series(ctx)
	}
}

func (c *CancellableOperator) GetPool() *model.VectorPool {
	return c.next.GetPool()
}
