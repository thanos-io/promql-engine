// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package function

import (
	"context"
	"math"
	"time"

	"github.com/prometheus/prometheus/model/labels"

	"github.com/thanos-io/promql-engine/execution/model"
)

type scalarFunctionOperator struct {
	pool     *model.VectorPool
	next     model.VectorOperator
	ti       model.TimingInformation
	duration time.Duration
}

func (o *scalarFunctionOperator) Analyze() (*model.TimingInformation, []model.ObservableVectorOperator) {

	// Create a TimingInformation instance and add CPU time taken
	timingInfo := &model.TimingInformation{}
	timingInfo.AddCPUTimeTaken(o.duration)

	// Create a slice to store any observable vector operators found during analysis
	observableOperators := make([]model.ObservableVectorOperator, 0)

	//TODO: Add any observable vector operators found to the slice "observableOperators"

	return timingInfo, observableOperators
}

func (o *scalarFunctionOperator) Explain() (me string, next []model.VectorOperator) {
	return "[*scalarFunctionOperator]", []model.VectorOperator{}
}

func (o *scalarFunctionOperator) Series(ctx context.Context) ([]labels.Labels, error) {
	return nil, nil
}

func (o *scalarFunctionOperator) GetPool() *model.VectorPool {
	return o.pool
}

func (o *scalarFunctionOperator) Next(ctx context.Context) ([]model.StepVector, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	start := time.Now()
	in, err := o.next.Next(ctx)
	if err != nil {
		return nil, err
	}
	if len(in) == 0 {
		return nil, nil
	}

	result := o.GetPool().GetVectorBatch()
	for _, vector := range in {
		sv := o.GetPool().GetStepVector(vector.T)
		if len(vector.Samples) != 1 {
			sv.AppendSample(o.GetPool(), 0, math.NaN())
		} else {
			sv.AppendSample(o.GetPool(), 0, vector.Samples[0])
		}
		result = append(result, sv)
		o.next.GetPool().PutStepVector(vector)
	}
	o.next.GetPool().PutVectors(in)
	duration := time.Since(start)
	o.duration = duration
	return result, nil
}
