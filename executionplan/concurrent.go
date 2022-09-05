package executionplan

import (
	"context"

	"github.com/prometheus/prometheus/promql"
)

type concurrencyOperator struct {
	next    VectorOperator
	buffer  chan promql.Vector
	started bool
}

func concurrent(next VectorOperator) VectorOperator {
	return &concurrencyOperator{
		next:   next,
		buffer: make(chan promql.Vector, 128),
	}
}

func (c *concurrencyOperator) Next(ctx context.Context) (promql.Vector, error) {
	if !c.started {
		c.start(ctx)
	}

	r, ok := <-c.buffer
	if !ok {
		return nil, nil
	}
	return r, nil
}

func (c *concurrencyOperator) start(ctx context.Context) {
	c.started = true
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
