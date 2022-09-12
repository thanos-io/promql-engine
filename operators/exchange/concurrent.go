package exchange

import (
	"context"
	"sync"

	"github.com/fpetkovski/promql-engine/operators/model"

	"github.com/prometheus/prometheus/model/labels"
)

type concurrencyOperator struct {
	next   model.Vector
	buffer chan []model.StepVector
	once   sync.Once
}

func NewConcurrent(next model.Vector, bufferSize int) model.Vector {
	return &concurrencyOperator{
		next:   next,
		buffer: make(chan []model.StepVector, bufferSize),
	}
}

func (c *concurrencyOperator) Series(ctx context.Context) ([]labels.Labels, error) {
	return c.next.Series(ctx)
}

func (c *concurrencyOperator) GetPool() *model.VectorPool {
	return c.next.GetPool()
}

func (c *concurrencyOperator) Next(ctx context.Context) ([]model.StepVector, error) {
	c.once.Do(func() { c.pull(ctx) })

	r, ok := <-c.buffer
	if !ok {
		return nil, nil
	}
	return r, nil
}

func (c *concurrencyOperator) pull(ctx context.Context) {
	go func() {
		defer close(c.buffer)
		for {
			r, err := c.next.Next(ctx)
			if err != nil || r == nil {
				break
			}
			c.buffer <- r
		}
	}()
}
