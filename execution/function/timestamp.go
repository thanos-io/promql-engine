// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package function

import (
	"context"
	"sync"

	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/execution/telemetry"
	"github.com/thanos-io/promql-engine/extlabels"
	"github.com/thanos-io/promql-engine/query"

	"github.com/prometheus/prometheus/model/labels"
)

type timestampOperator struct {
	pool *model.VectorPool
	next model.VectorOperator

	series []labels.Labels
	once   sync.Once
}

func newTimestampOperator(pool *model.VectorPool, next model.VectorOperator, opts *query.Options) model.VectorOperator {
	oper := &timestampOperator{
		pool: pool,
		next: next,
	}
	return telemetry.NewOperator(telemetry.NewTelemetry(oper, opts), oper)
}

func (o *timestampOperator) Explain() (next []model.VectorOperator) {
	return []model.VectorOperator{o.next}
}

func (o *timestampOperator) Series(ctx context.Context) ([]labels.Labels, error) {
	if err := o.loadSeries(ctx); err != nil {
		return nil, err
	}
	return o.series, nil
}

func (o *timestampOperator) String() string {
	return "[timestamp]"
}

func (o *timestampOperator) loadSeries(ctx context.Context) error {
	var err error
	o.once.Do(func() {
		series, loadErr := o.next.Series(ctx)
		if loadErr != nil {
			err = loadErr
			return
		}
		o.series = make([]labels.Labels, len(series))

		b := labels.ScratchBuilder{}
		for i, s := range series {
			lbls, _ := extlabels.DropMetricName(s, b)
			o.series[i] = lbls
		}
	})

	return err
}

func (o *timestampOperator) GetPool() *model.VectorPool {
	return o.pool
}

func (o *timestampOperator) Next(ctx context.Context) ([]model.StepVector, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	in, err := o.next.Next(ctx)

	if err != nil {
		return nil, err
	}
	if len(in) == 0 {
		return nil, nil
	}
	result := o.pool.GetVectorBatch()
	for _, vector := range in {
		out := o.pool.GetStepVector(vector.T)
		value := float64(vector.T / 1000)

		for i := range vector.Samples {
			out.AppendSample(o.pool, vector.SampleIDs[i], value)
		}
		for i := range vector.Histograms {
			out.AppendSample(o.pool, vector.HistogramIDs[i], value)
		}

		result = append(result, out)
		o.next.GetPool().PutStepVector(vector)
	}

	o.next.GetPool().PutVectors(in)
	return result, nil
}
