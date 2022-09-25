// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package aggregate

import (
	"context"
	"sync"

	"github.com/thanos-community/promql-engine/worker"

	"github.com/thanos-community/promql-engine/physicalplan/model"

	"github.com/prometheus/prometheus/model/labels"

	"github.com/prometheus/prometheus/promql/parser"
)

type aggregate struct {
	next model.VectorOperator

	hashBuf    []byte
	vectorPool *model.VectorPool

	by          bool
	labels      []string
	aggregation parser.ItemType

	once   sync.Once
	tables []aggregateTable
	series []labels.Labels

	stepsBatch int
	workers    worker.Group
}

func NewHashAggregate(
	points *model.VectorPool,
	next model.VectorOperator,
	aggregation parser.ItemType,
	by bool,
	labels []string,
	stepsBatch int,
) (model.VectorOperator, error) {
	a := &aggregate{
		next:       next,
		vectorPool: points,

		by:          by,
		aggregation: aggregation,
		labels:      labels,
		stepsBatch:  stepsBatch,
	}
	a.workers = worker.NewGroup(stepsBatch, a.workerTask)

	return a, nil
}

func (a *aggregate) Series(ctx context.Context) ([]labels.Labels, error) {
	var err error
	a.once.Do(func() { err = a.initializeTables(ctx) })
	if err != nil {
		return nil, err
	}

	return a.series, nil
}

func (a *aggregate) GetPool() *model.VectorPool {
	return a.vectorPool
}

func (a *aggregate) Next(ctx context.Context) ([]model.StepVector, error) {
	in, err := a.next.Next(ctx)
	if err != nil {
		return nil, err
	}
	if in == nil {
		return nil, nil
	}
	defer a.next.GetPool().PutVectors(in)

	a.once.Do(func() { err = a.initializeTables(ctx) })
	if err != nil {
		return nil, err
	}

	result := a.vectorPool.GetVectorBatch()
	for i, vector := range in {
		a.workers[i].Send(vector)
	}

	for i, vector := range in {
		result = append(result, a.workers[i].GetOutput())
		a.next.GetPool().PutStepVector(vector)
	}

	return result, nil
}

func (a *aggregate) initializeTables(ctx context.Context) error {
	series, err := a.next.Series(ctx)
	if err != nil {
		return err
	}

	a.workers.Start(ctx)

	if a.by && len(a.labels) == 0 {
		tables, err := newVectorizedTables(a.stepsBatch, a.aggregation)
		if err != nil {
			return err
		}
		a.tables = tables
		a.series = []labels.Labels{{}}
		a.vectorPool.SetStepSize(1)
		return nil
	}

	inputCache := make([]uint64, len(series))
	outputMap := make(map[uint64]*model.Series)
	outputCache := make([]*model.Series, 0)
	buf := make([]byte, 128)
	for i := 0; i < len(series); i++ {
		hash, _, lbls := hashMetric(series[i], !a.by, a.labels, buf)

		output, ok := outputMap[hash]
		if !ok {
			output = &model.Series{
				Metric: lbls,
				ID:     uint64(len(outputCache)),
			}
			outputMap[hash] = output
			outputCache = append(outputCache, output)
		}

		inputCache[i] = output.ID
	}

	a.series = make([]labels.Labels, len(outputCache))
	for i := 0; i < len(outputCache); i++ {
		a.series[i] = outputCache[i].Metric
	}

	a.tables = newScalarTables(a.stepsBatch, inputCache, outputCache, a.aggregation)
	a.vectorPool.SetStepSize(len(outputCache))

	return nil
}

func (a *aggregate) workerTask(workerID int, vector model.StepVector) model.StepVector {
	table := a.tables[workerID]
	table.aggregate(vector)
	return table.toVector(a.vectorPool)
}
