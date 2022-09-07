package executionplan

import (
	"context"

	"github.com/fpetkovski/promql-engine/model"
)

type coalesceOperator struct {
	operators []VectorOperator
}

func coalesce(operators ...VectorOperator) VectorOperator {
	return &coalesceOperator{operators: operators}
}

func (c coalesceOperator) Next(ctx context.Context) ([]model.StepVector, error) {
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
			out = make([]model.StepVector, len(in))
			for i := 0; i < len(in); i++ {
				size := len(in[i].Samples) * len(c.operators)
				out[i] = model.StepVector{
					T:       in[i].T,
					Samples: make([]model.StepSample, 0, size),
				}
			}
		}
		for i := 0; i < len(in); i++ {
			out[i].Samples = append(out[i].Samples, in[i].Samples...)
		}
	}
	if out == nil {
		return nil, nil
	}

	return out, nil
}
