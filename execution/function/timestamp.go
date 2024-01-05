// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package function

import (
	"context"
	"time"

	"github.com/prometheus/prometheus/model/labels"

	"github.com/thanos-io/promql-engine/execution/model"
)

type timestampFunctionOperator struct {
	next model.VectorOperator
	model.OperatorTelemetry
}

func (o *timestampFunctionOperator) Analyze() (model.OperatorTelemetry, []model.ObservableVectorOperator) {
	o.SetName("[*timestampFunctionOperator]")
	next := make([]model.ObservableVectorOperator, 0, 1)
	if obsnext, ok := o.next.(model.ObservableVectorOperator); ok {
		next = append(next, obsnext)
	}
	return o, next
}

func (o *timestampFunctionOperator) Explain() (me string, next []model.VectorOperator) {
	return "[*timestampFunctionOperator]", []model.VectorOperator{o.next}
}

func (o *timestampFunctionOperator) Series(ctx context.Context) ([]labels.Labels, error) {
	return o.next.Series(ctx)
}

func (o *timestampFunctionOperator) GetPool() *model.VectorPool {
	return o.next.GetPool()
}

func (o *timestampFunctionOperator) Next(ctx context.Context) ([]model.StepVector, error) {
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
	for _, vector := range in {
		for i := range vector.Samples {
			vector.Samples[i] = float64(vector.T / 1000)
		}
	}
	o.OperatorTelemetry.AddExecutionTimeTaken(time.Since(start))
	return in, nil
}
