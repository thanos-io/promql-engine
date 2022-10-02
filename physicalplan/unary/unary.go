package unary

import (
	"context"
	"sync"

	"github.com/prometheus/prometheus/model/labels"
	"gonum.org/v1/gonum/floats"

	"github.com/thanos-community/promql-engine/physicalplan/model"
	"github.com/thanos-community/promql-engine/worker"
)

type unaryNegation struct {
	next model.VectorOperator
	once sync.Once

	series []labels.Labels

	stepsBatch int
	workers    worker.Group
}

func NewUnaryNegation(
	next model.VectorOperator,
	stepsBatch int,
) (model.VectorOperator, error) {
	u := &unaryNegation{
		next:       next,
		stepsBatch: stepsBatch,
	}

	u.workers = worker.NewGroup(stepsBatch, u.workerTask)
	return u, nil
}

func (u *unaryNegation) Series(ctx context.Context) ([]labels.Labels, error) {
	var err error
	u.once.Do(func() { err = u.loadSeries(ctx) })
	if err != nil {
		return nil, err
	}
	return u.series, nil
}

func (u *unaryNegation) loadSeries(ctx context.Context) error {
	vectorSeries, err := u.next.Series(ctx)
	if err != nil {
		return err
	}
	u.series = make([]labels.Labels, len(vectorSeries))
	for i := range vectorSeries {
		lbls := labels.NewBuilder(vectorSeries[i]).Del(labels.MetricName).Labels()
		u.series[i] = lbls
	}

	u.workers.Start(ctx)
	return nil
}

func (u *unaryNegation) GetPool() *model.VectorPool {
	return u.next.GetPool()
}

func (u *unaryNegation) Next(ctx context.Context) ([]model.StepVector, error) {
	in, err := u.next.Next(ctx)
	if err != nil {
		return nil, err
	}
	if in == nil {
		return nil, nil
	}
	defer u.next.GetPool().PutVectors(in)

	for i, vector := range in {
		if err := u.workers[i].Send(vector); err != nil {
			return nil, err
		}
	}

	result := u.next.GetPool().GetVectorBatch()
	for i, vector := range in {
		output, err := u.workers[i].GetOutput()
		if err != nil {
			return nil, err
		}
		result = append(result, output)
		u.next.GetPool().PutStepVector(vector)
	}

	return result, nil
}

func (u *unaryNegation) workerTask(_ int, vector model.StepVector) model.StepVector {
	floats.Scale(-1, vector.Samples)
	return vector
}
