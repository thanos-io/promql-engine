// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package function

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/prometheus/model/labels"

	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/extlabels"
)

type timestampFunctionOperator struct {
	next model.VectorOperator
	model.OperatorTelemetry

	series []labels.Labels
	once   sync.Once
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
	if err := o.loadSeries(ctx); err != nil {
		return nil, err
	}
	return o.series, nil
}

func (o *timestampFunctionOperator) loadSeries(ctx context.Context) error {
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

func (o *timestampFunctionOperator) GetPool() *model.VectorPool {
	return o.next.GetPool()
}

func (o *timestampFunctionOperator) Next(ctx context.Context) ([]model.StepVector, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	if err := o.loadSeries(ctx); err != nil {
		return nil, err
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
