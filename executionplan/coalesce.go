package executionplan

import (
	"context"
	"github.com/prometheus/prometheus/promql"
)

type coalesce struct {
	operators []VectorOperator
}

func newCoalesce(operators ...VectorOperator) *coalesce {
	return &coalesce{operators: operators}
}

func (c coalesce) Next(ctx context.Context) (promql.Vector, error) {
	out := make(promql.Vector, 0, len(c.operators))
	for _, o := range c.operators {
		r, err := o.Next(ctx)
		if err != nil {
			return nil, err
		}
		if r == nil {
			continue
		}
		out = append(out, r...)
	}
	if len(out) == 0 {
		return nil, nil
	}

	return out, nil
}
