package executionplan

import (
	"context"

	"github.com/fpetkovski/promql-engine/model"
)

type coalesceOperator struct {
	pool      *model.VectorPool
	operators []VectorOperator
}

func coalesce(pool *model.VectorPool, operators ...VectorOperator) VectorOperator {
	return &coalesceOperator{
		pool:      pool,
		operators: operators,
	}
}

func (c *coalesceOperator) GetPool() *model.VectorPool {
	return c.pool
}

func (c *coalesceOperator) Next(ctx context.Context) ([]model.StepVector, error) {
	var out []model.StepVector = nil
	for _, o := range c.operators {
		in, err := o.Next(ctx)
		if err != nil {
			return nil, err
		}
		if in == nil {
			continue
		}
		if len(in) > 0 && out == nil {
			c.pool.SetStepSamplesSize(len(in) * len(c.operators))
			out = c.pool.GetVectors()
			for i := 0; i < len(in); i++ {
				out = append(out, model.StepVector{
					T:       in[i].T,
					Samples: c.pool.GetSamples(),
				})
			}
		}
		for i := 0; i < len(in); i++ {
			out[i].Samples = append(out[i].Samples, in[i].Samples...)
			o.GetPool().PutSamples(in[i].Samples)
		}
		o.GetPool().PutVectors(in)
	}
	if out == nil {
		return nil, nil
	}

	return out, nil
}
