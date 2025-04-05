// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package unary

import (
	"context"
	"sync"
	"time"

	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/execution/telemetry"
	"github.com/thanos-io/promql-engine/query"

	"github.com/prometheus/prometheus/model/histogram"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"gonum.org/v1/gonum/floats"
)

type unaryNegation struct {
	next model.VectorOperator
	once sync.Once

	series []promql.Series
	telemetry.OperatorTelemetry
	enableDelayedNameRemoval bool
}

func NewUnaryNegation(next model.VectorOperator, opts *query.Options) (model.VectorOperator, error) {
	u := &unaryNegation{
		next:                     next,
		enableDelayedNameRemoval: opts.EnableDelayedNameRemoval,
	}
	u.OperatorTelemetry = telemetry.NewTelemetry(u, opts)

	return u, nil
}

func (u *unaryNegation) Explain() (next []model.VectorOperator) {
	return []model.VectorOperator{u.next}
}

func (u *unaryNegation) String() string {
	return "[unaryNegation]"
}

func (u *unaryNegation) Series(ctx context.Context) ([]promql.Series, error) {
	start := time.Now()
	defer func() { u.AddExecutionTimeTaken(time.Since(start)) }()

	if err := u.loadSeries(ctx); err != nil {
		return nil, err
	}
	return u.series, nil
}

func (u *unaryNegation) loadSeries(ctx context.Context) error {
	var err error
	u.once.Do(func() {
		var series []promql.Series
		series, err = u.next.Series(ctx)
		if err != nil {
			return
		}
		u.series = make([]promql.Series, len(series))
		for i := range series {
			if !u.enableDelayedNameRemoval {
				u.series[i].Metric = labels.NewBuilder(series[i].Metric).Del(labels.MetricName).Labels()
			} else {
				u.series[i] = promql.Series{
					Metric:   series[i].Metric,
					DropName: true,
				}
			}
		}
	})
	return err
}

func (u *unaryNegation) GetPool() *model.VectorPool {
	return u.next.GetPool()
}

func (u *unaryNegation) Next(ctx context.Context) ([]model.StepVector, error) {
	start := time.Now()
	defer func() { u.AddExecutionTimeTaken(time.Since(start)) }()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	in, err := u.next.Next(ctx)
	if err != nil {
		return nil, err
	}
	if in == nil {
		return nil, nil
	}
	for i := range in {
		floats.Scale(-1, in[i].Samples)
		negateHistograms(in[i].Histograms)
	}
	return in, nil
}

func negateHistograms(hists []*histogram.FloatHistogram) {
	for i := range hists {
		hists[i] = hists[i].Copy().Mul(-1)
	}
}
