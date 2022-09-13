package exchange

import (
	"context"
	"sync"

	"github.com/fpetkovski/promql-engine/operators/model"
	"github.com/prometheus/prometheus/model/labels"
)

type coalesceOperator struct {
	once   sync.Once
	series []labels.Labels

	pool      *model.VectorPool
	operators []model.VectorOperator
}

func NewCoalesce(pool *model.VectorPool, operators ...model.VectorOperator) model.VectorOperator {
	return &coalesceOperator{
		pool:      pool,
		operators: operators,
	}
}

func (c *coalesceOperator) GetPool() *model.VectorPool {
	return c.pool
}

func (c *coalesceOperator) Series(ctx context.Context) ([]labels.Labels, error) {
	var err error
	c.once.Do(func() { err = c.loadSeries(ctx) })
	if err != nil {
		return nil, err
	}
	return c.series, nil
}

func (c *coalesceOperator) loadSeries(ctx context.Context) error {
	size := 0
	for i := 0; i < len(c.operators); i++ {
		series, err := c.operators[i].Series(ctx)
		if err != nil {
			return err
		}
		size += len(series)
	}

	idx := 0
	result := make([]labels.Labels, size)
	for _, o := range c.operators {
		series, err := o.Series(ctx)
		if err != nil {
			return err
		}
		for i := 0; i < len(series); i++ {
			result[idx] = series[i]
			idx++
		}
	}
	c.series = result
	c.pool.SetStepSamplesSize(len(c.series))

	return nil
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
