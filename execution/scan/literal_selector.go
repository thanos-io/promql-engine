// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package scan

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/query"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/prometheus/model/labels"
)

var (
	constantMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "operator_duration_seconds",
		Help: "Time taken by each operator",
	}, []string{"operator", "duration"})
)

// numberLiteralSelector returns []model.StepVector with same sample value across time range.
type numberLiteralSelector struct {
	vectorPool *model.VectorPool

	numSteps    int
	mint        int64
	maxt        int64
	step        int64
	currentStep int64
	series      []labels.Labels
	once        sync.Once

	val        float64
	recordTime bool
	duration   time.Duration
}

func NewNumberLiteralSelector(pool *model.VectorPool, opts *query.Options, val float64, recordTime bool) *numberLiteralSelector {
	return &numberLiteralSelector{
		vectorPool:  pool,
		numSteps:    opts.NumSteps(),
		mint:        opts.Start.UnixMilli(),
		maxt:        opts.End.UnixMilli(),
		step:        opts.Step.Milliseconds(),
		currentStep: opts.Start.UnixMilli(),
		val:         val,
		recordTime:  recordTime,
		duration:    time.Duration(0),
	}
}

func (o *numberLiteralSelector) Explain() (me string, next []model.VectorOperator) {
	return fmt.Sprintf("[*numberLiteralSelector] %v", o.val), nil
}

func (o *numberLiteralSelector) Series(context.Context) ([]labels.Labels, error) {
	o.loadSeries()
	return o.series, nil
}

func (o *numberLiteralSelector) GetPool() *model.VectorPool {
	return o.vectorPool
}

func (o *numberLiteralSelector) Next(ctx context.Context) ([]model.StepVector, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	start := time.Now()

	if o.currentStep > o.maxt {
		return nil, nil
	}

	o.loadSeries()

	ts := o.currentStep
	vectors := o.vectorPool.GetVectorBatch()
	for currStep := 0; currStep < o.numSteps && ts <= o.maxt; currStep++ {
		stepVector := o.vectorPool.GetStepVector(ts)
		stepVector.AppendSample(o.vectorPool, 0, o.val)
		vectors = append(vectors, stepVector)

		ts += o.step
	}

	// For instant queries, set the step to a positive value
	// so that the operator can terminate.
	if o.step == 0 {
		o.step = 1
	}
	o.currentStep += o.step * int64(o.numSteps)
	duration := time.Since(start)
	if o.recordTime {
		o.duration = duration
		const opName = "numberLiteralSelector"
		constantMetric.WithLabelValues(opName, fmt.Sprintf("%v", duration)).Add(duration.Seconds())
	}

	return vectors, nil
}

func (o *numberLiteralSelector) loadSeries() {
	// If number literal is included within function, []labels.labels must be initialized.
	o.once.Do(func() {
		o.series = make([]labels.Labels, 1)
		o.vectorPool.SetStepSize(len(o.series))
	})
}
