// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package exchange

import (
	"context"
	"sync"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/thanos-community/promql-engine/physicalplan/model"
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

func (c *coalesceOperator) Next(ctx context.Context) ([]model.StepVector, error) {
	var out []model.StepVector = nil
	var wg sync.WaitGroup
	var mu sync.RWMutex
	var errChan = make(errorChan, len(c.operators))
	for _, o := range c.operators {
		wg.Add(1)
		go func(o model.VectorOperator) {
			defer wg.Done()

			in, err := o.Next(ctx)
			if err != nil {
				errChan <- err
				return
			}
			if in == nil {
				return
			}
			mu.RLock()
			if len(in) > 0 && out == nil {
				mu.RUnlock()
				mu.Lock()
				if len(in) > 0 && out == nil {
					out = c.pool.GetVectorBatch()
					for i := 0; i < len(in); i++ {
						out = append(out, c.pool.GetStepVector(in[i].T))
					}
				}
				mu.Unlock()
			} else {
				mu.RUnlock()
			}

			mu.Lock()
			for i := 0; i < len(in); i++ {
				out[i].Samples = append(out[i].Samples, in[i].Samples...)
				out[i].SampleIDs = append(out[i].SampleIDs, in[i].SampleIDs...)
				o.GetPool().PutStepVector(in[i])
			}
			mu.Unlock()
			o.GetPool().PutVectors(in)
		}(o)
	}
	wg.Wait()
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
