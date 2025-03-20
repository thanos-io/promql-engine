// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package function

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/prometheus/promql"

	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/execution/telemetry"
	"github.com/thanos-io/promql-engine/query"
)

type timestampOperator struct {
	next model.VectorOperator
	telemetry.OperatorTelemetry

	series                   []promql.Series
	once                     sync.Once
	enableDelayedNameRemoval bool
}

func newTimestampOperator(next model.VectorOperator, opts *query.Options) *timestampOperator {
	oper := &timestampOperator{
		next:                     next,
		enableDelayedNameRemoval: opts.EnableDelayedNameRemoval,
	}
	oper.OperatorTelemetry = telemetry.NewTelemetry(oper, opts)

	return oper
}

func (o *timestampOperator) Explain() (next []model.VectorOperator) {
	return []model.VectorOperator{o.next}
}

func (o *timestampOperator) Series(ctx context.Context) ([]promql.Series, error) {
	start := time.Now()
	defer func() { o.AddExecutionTimeTaken(time.Since(start)) }()

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
		o.series = make([]promql.Series, len(series))

		for i, s := range series {
			o.series[i] = s
			if !o.enableDelayedNameRemoval {
				o.series[i].Metric = o.series[i].Metric.DropMetricName()
			} else {
				o.series[i].DropName = true
			}
		}
	})

	return err
}

func (o *timestampOperator) GetPool() *model.VectorPool {
	return o.next.GetPool()
}

func (o *timestampOperator) Next(ctx context.Context) ([]model.StepVector, error) {
	start := time.Now()
	defer func() { o.AddExecutionTimeTaken(time.Since(start)) }()

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
