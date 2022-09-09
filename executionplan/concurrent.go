package executionplan

import (
	"context"
	"sync"

	"github.com/fpetkovski/promql-engine/model"
)

type concurrencyOperator struct {
	next   VectorOperator
	buffer chan []model.StepVector
	once   sync.Once
}

func (c *concurrencyOperator) GetPool() *model.VectorPool {
	return c.next.GetPool()
}

func concurrent(next VectorOperator, bufferSize int) VectorOperator {
	return &concurrencyOperator{
		next:   next,
		buffer: make(chan []model.StepVector, bufferSize),
	}
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
