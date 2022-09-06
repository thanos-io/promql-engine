package executionplan

import (
	"context"
	"sync"

	"github.com/prometheus/prometheus/promql"
)

type concurrencyOperator struct {
	next   VectorOperator
	buffer chan promql.Vector
	once   sync.Once
}

func concurrent(next VectorOperator) VectorOperator {
	return &concurrencyOperator{
		next:   next,
		buffer: make(chan promql.Vector, 300),
	}
}

func (c *concurrencyOperator) Next(ctx context.Context) (promql.Vector, error) {
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
