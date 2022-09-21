package binary

import (
	"context"
	"sync"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/thanos-community/promql-engine/physicalplan/model"
)

// scalarOperator evaluates expressions where one operand is a scalarOperator.
type scalarOperator struct {
	seriesOnce sync.Once
	series     []labels.Labels
	scalar     float64

	pool           *model.VectorPool
	numberSelector model.VectorOperator
	next           model.VectorOperator
	getOperands    getOperandsFunc
	operation      operation
}

func NewScalar(pool *model.VectorPool, next model.VectorOperator, numberSelector model.VectorOperator, op parser.ItemType, scalarSideLeft bool) (*scalarOperator, error) {
	binaryOperation, err := newOperation(op)
	if err != nil {
		return nil, err
	}
	var getOperands getOperandsFunc
	if scalarSideLeft {
		getOperands = getOperandsScalarLeft
	} else {
		getOperands = getOperandsScalarRight
	}

	// Cache the result of the number selector since it
	// will not change during execution.
	v, err := numberSelector.Next(context.Background())
	if err != nil {
		return nil, err
	}
	scalar := v[0].Samples[0]

	return &scalarOperator{
		pool:           pool,
		next:           next,
		scalar:         scalar,
		numberSelector: numberSelector,
		operation:      binaryOperation,
		getOperands:    getOperands,
	}, nil
}

func (o *scalarOperator) Series(ctx context.Context) ([]labels.Labels, error) {
	var err error
	o.seriesOnce.Do(func() { err = o.loadSeries(ctx) })
	if err != nil {
		return nil, err
	}
	return o.series, nil
}

func (o *scalarOperator) Next(ctx context.Context) ([]model.StepVector, error) {
	in, err := o.next.Next(ctx)
	if err != nil {
		return nil, err
	}
	if in == nil {
		return nil, nil
	}
	o.seriesOnce.Do(func() { err = o.loadSeries(ctx) })
	if err != nil {
		return nil, err
	}

	out := o.pool.GetVectorBatch()
	for _, vector := range in {
		step := o.pool.GetStepVector(vector.T)
		for i := range vector.Samples {
			lhs, rhs := o.getOperands(vector, i, o.scalar)
			val := o.operation(lhs, rhs)
			step.Samples = append(step.Samples, val)
			step.SampleIDs = append(step.SampleIDs, vector.SampleIDs[i])
		}
		out = append(out, step)
		o.next.GetPool().PutStepVector(vector)
	}
	o.next.GetPool().PutVectors(in)
	return out, nil
}

func (o *scalarOperator) GetPool() *model.VectorPool {
	return o.pool
}

func (o *scalarOperator) loadSeries(ctx context.Context) error {
	vectorSeries, err := o.next.Series(ctx)
	if err != nil {
		return err
	}
	series := make([]labels.Labels, len(vectorSeries))
	for i := range vectorSeries {
		lbls := labels.NewBuilder(vectorSeries[i]).Del(labels.MetricName).Labels()
		series[i] = lbls
	}

	o.series = series
	return nil
}

type getOperandsFunc func(v model.StepVector, i int, scalar float64) (float64, float64)

func getOperandsScalarLeft(v model.StepVector, i int, scalar float64) (float64, float64) {
	return scalar, v.Samples[i]
}

func getOperandsScalarRight(v model.StepVector, i int, scalar float64) (float64, float64) {
	return v.Samples[i], scalar
}
