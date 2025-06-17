// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package function

import (
	"context"
	"fmt"

	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/execution/telemetry"
	"github.com/thanos-io/promql-engine/query"

	"github.com/prometheus/prometheus/model/labels"
)

// sortOperator only filters out native histogram samples.
// The actual sorting logic is handled by the engine at presentation time.
type sortOperator struct {
	sortFn string
	pool   *model.VectorPool
	next   model.VectorOperator
}

func newSortOperator(pool *model.VectorPool, sortFn string, next model.VectorOperator, opts *query.Options) model.VectorOperator {
	oper := &sortOperator{
		sortFn: sortFn,
		pool:   pool,
		next:   next,
	}

	return telemetry.NewOperator(telemetry.NewTelemetry(oper, opts), oper)
}

func (o *sortOperator) String() string {
	return fmt.Sprintf("[%s]", o.sortFn)
}

func (o *sortOperator) Explain() (next []model.VectorOperator) {
	return []model.VectorOperator{o.next}
}

func (o *sortOperator) Series(ctx context.Context) ([]labels.Labels, error) {
	return o.next.Series(ctx)
}

func (o *sortOperator) GetPool() *model.VectorPool {
	return o.pool
}

func (o *sortOperator) Next(ctx context.Context) ([]model.StepVector, error) {
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

	result := o.GetPool().GetVectorBatch()
	for _, vector := range in {
		vector.Histograms = nil
		sv := o.GetPool().GetStepVector(vector.T)

		sv.SampleIDs = vector.SampleIDs
		sv.Samples = vector.Samples
		result = append(result, sv)

		o.next.GetPool().PutStepVector(vector)
	}
	o.next.GetPool().PutVectors(in)

	return result, nil
}
