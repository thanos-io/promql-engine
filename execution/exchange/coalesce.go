// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package exchange

import (
	"context"
	"sync"

	"github.com/prometheus/prometheus/model/labels"

	"github.com/thanos-community/promql-engine/execution/model"
)

type errorChan chan error

func (c errorChan) getError() error {
	for err := range c {
		if err != nil {
			return err
		}
	}

	return nil
}

type coalesceOperator struct {
	once   sync.Once
	series []labels.Labels

	pool      *model.VectorPool
	mu        sync.Mutex
	wg        sync.WaitGroup
	operators []model.VectorOperator
}

func NewCoalesce(pool *model.VectorPool, operators ...model.VectorOperator) model.VectorOperator {
	return &coalesceOperator{
		pool:      pool,
		operators: operators,
	}
}

func (c *coalesceOperator) Explain() (me string, next []model.VectorOperator) {
	return "[*coalesceOperator]", c.operators
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

func (c *coalesceOperator) Next(ctx context.Context) ([]model.StepVector, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var out []model.StepVector = nil
	var errChan = make(errorChan, len(c.operators))
	for _, o := range c.operators {
		c.wg.Add(1)
		go func(o model.VectorOperator) {
			defer c.wg.Done()

			in, err := o.Next(ctx)
			if err != nil {
				errChan <- err
				return
			}
			if in == nil {
				return
			}
			c.mu.Lock()
			defer c.mu.Unlock()

			if len(in) > 0 && out == nil {
				out = c.pool.GetVectorBatch()
				for i := 0; i < len(in); i++ {
					out = append(out, c.pool.GetStepVector(in[i].T))
				}
			}

			for i := 0; i < len(in); i++ {
				out[i].Samples = append(out[i].Samples, in[i].Samples...)
				out[i].SampleIDs = append(out[i].SampleIDs, in[i].SampleIDs...)
				o.GetPool().PutStepVector(in[i])
			}
			o.GetPool().PutVectors(in)
		}(o)
	}
	c.wg.Wait()
	close(errChan)

	if err := errChan.getError(); err != nil {
		return nil, err
	}

	if out == nil {
		return nil, nil
	}

	return out, nil
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
	c.pool.SetStepSize(len(c.series))

	return nil
}
