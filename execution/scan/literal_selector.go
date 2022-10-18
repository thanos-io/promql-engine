// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package scan

import (
	"context"
	"fmt"
	"sync"

	"github.com/prometheus/prometheus/promql/parser"

	"github.com/thanos-community/promql-engine/execution/model"
	"github.com/thanos-community/promql-engine/query"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
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

	val      float64
	call     FunctionCall
	callName string
}

func NewNumberLiteralSelector(pool *model.VectorPool, opts *query.Options, val float64) *numberLiteralSelector {
	return &numberLiteralSelector{
		vectorPool:  pool,
		numSteps:    opts.NumSteps(),
		mint:        opts.Start.UnixMilli(),
		maxt:        opts.End.UnixMilli(),
		step:        opts.Step.Milliseconds(),
		currentStep: opts.Start.UnixMilli(),
		val:         val,
	}
}

func NewNumberLiteralSelectorWithFunc(pool *model.VectorPool, opts *query.Options, val float64, f *parser.Function) (model.VectorOperator, error) {
	call, err := NewFunctionCall(f)
	if err != nil {
		return nil, err
	}

	selector := NewNumberLiteralSelector(pool, opts, val)
	selector.call = call
	selector.callName = f.Name
	return selector, nil
}

func (o *numberLiteralSelector) Explain() (me string, next []model.VectorOperator) {
	if o.call != nil {
		return fmt.Sprintf("[*numberLiteralSelector] %v(%v)", o.callName, o.val), nil
	}
	return fmt.Sprintf("[*numberLiteralSelector] %v", o.val), nil
}

func (o *numberLiteralSelector) Series(context.Context) ([]labels.Labels, error) {
	o.loadSeries()
	return o.series, nil
}

func (o *numberLiteralSelector) GetPool() *model.VectorPool {
	return o.vectorPool
}

func (o *numberLiteralSelector) Next(context.Context) ([]model.StepVector, error) {
	if o.currentStep > o.maxt {
		return nil, nil
	}

	o.loadSeries()

	vectors := o.vectorPool.GetVectorBatch()
	ts := o.currentStep
	for currStep := 0; currStep < o.numSteps && ts <= o.maxt; currStep++ {
		if len(vectors) <= currStep {
			vectors = append(vectors, o.vectorPool.GetStepVector(ts))
		}

		result := promql.Sample{Point: promql.Point{T: ts, V: o.val}}
		if o.call != nil {
			result = o.call(o.series[0], []promql.Point{result.Point}, ts, 0)
		}

		vectors[currStep].T = result.T
		vectors[currStep].SampleIDs = append(vectors[currStep].SampleIDs, uint64(0))
		vectors[currStep].Samples = append(vectors[currStep].Samples, result.V)

		ts += o.step
	}

	// For instant queries, set the step to a positive value
	// so that the operator can terminate.
	if o.step == 0 {
		o.step = 1
	}
	o.currentStep += o.step * int64(o.numSteps)

	return vectors, nil
}

func (o *numberLiteralSelector) loadSeries() {
	// If number literal is included within function, []labels.labels must be initialized.
	o.once.Do(func() {
		o.series = make([]labels.Labels, 1)
		if o.call != nil {
			o.series = []labels.Labels{labels.New()}
		}
		o.vectorPool.SetStepSize(len(o.series))
	})
}
