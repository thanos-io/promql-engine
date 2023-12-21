// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package function

import (
	"context"

	"github.com/prometheus/prometheus/model/labels"

	"github.com/thanos-io/promql-engine/execution/model"
)

type timestampFunctionOperator struct {
	next model.VectorOperator
}

func (o *timestampFunctionOperator) Explain() (me string, next []model.VectorOperator) {
	return "[*timestampFunctionOperator]", []model.VectorOperator{}
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
	in, err := o.next.Next(ctx)
	if err != nil {
		return nil, err
	}
	for _, vector := range in {
		for i := range vector.Samples {
			vector.Samples[i] = float64(vector.T / 1000)
		}
	}
	return in, nil
}
