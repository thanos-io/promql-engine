package executionplan

import (
	"context"
	"sync"

	"github.com/fpetkovski/promql-engine/model"
)

type concurrencyOperator struct {
	next   VectorOperator
	buffer chan []model.Vector
	once   sync.Once
}

func concurrent(next VectorOperator) VectorOperator {
	return &concurrencyOperator{
		next:   next,
		buffer: make(chan []model.Vector, 30),
	}
}

func (c *concurrencyOperator) Next(ctx context.Context) ([]model.Vector, error) {
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
