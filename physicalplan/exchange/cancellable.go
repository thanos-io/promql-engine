package exchange

import (
	"context"

	"github.com/prometheus/prometheus/model/labels"

	"github.com/thanos-community/promql-engine/physicalplan/model"
)

type CancellableOperator struct {
	next model.VectorOperator
}

func NewCancellable(next model.VectorOperator) *CancellableOperator {
	return &CancellableOperator{next: next}
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
