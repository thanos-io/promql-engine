// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package aggregate

import (
	"context"
	"fmt"
	"sync"

	"github.com/efficientgo/core/errors"

	"github.com/thanos-community/promql-engine/physicalplan/parse"
	"github.com/thanos-community/promql-engine/worker"

	"github.com/thanos-community/promql-engine/physicalplan/model"

	"github.com/prometheus/prometheus/model/labels"

	"github.com/prometheus/prometheus/promql/parser"
)

type aggregate struct {
	next model.VectorOperator

	vectorPool *model.VectorPool

	by          bool
	labels      []string
	aggregation parser.ItemType

	once           sync.Once
	tables         []aggregateTable
	series         []labels.Labels
	newAccumulator newAccumulatorFunc
	param          parser.Expr
	stepsBatch     int
	workers        worker.Group
}

func NewHashAggregate(
	points *model.VectorPool,
	next model.VectorOperator,
	aggregation parser.ItemType,
	param parser.Expr,
	by bool,
	labels []string,
	stepsBatch int,
) (model.VectorOperator, error) {
	if err := validateFunc(aggregation, param); err != nil {
		return nil, err
	}
	newAccumulator, err := makeAccumulatorFunc(aggregation, param)
	if err != nil {
		return nil, err
	}
	a := &aggregate{
		next:           next,
		vectorPool:     points,
		by:             by,
		aggregation:    aggregation,
		labels:         labels,
		stepsBatch:     stepsBatch,
		param:          param,
		newAccumulator: newAccumulator,
	}
	a.workers = worker.NewGroup(stepsBatch, a.workerTask)

	return a, nil
}

func validateFunc(op parser.ItemType, param parser.Expr) error {
	if op == parser.TOPK || op == parser.BOTTOMK || op == parser.QUANTILE {
		if err := aggregateValidation(param, parser.ValueTypeScalar); err != nil {
			return err
		}
	}

	if op == parser.COUNT_VALUES {
		if err := aggregateValidation(param, parser.ValueTypeString); err != nil {
			return err
		}
	}

	return nil
}

func aggregateValidation(node parser.Expr, want parser.ValueType) error {
	var t parser.ValueType
	switch n := node.(type) {
	case *parser.NumberLiteral:
		t = parser.ValueTypeScalar
	case *parser.StringLiteral:
		t = parser.ValueTypeString
	default:
		return fmt.Errorf("aggregateValidation() expected type %s in aggregation, got %T", parser.DocumentedType(want), n)
	}
	if t != want {
		return fmt.Errorf("aggregateValidation() expected type %s in aggregation, got %s", parser.DocumentedType(want), parser.DocumentedType(t))
	}
	return nil
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
		if err := a.workers[i].Send(vector); err != nil {
			return nil, err
		}
	}

	for i, vector := range in {
		output, err := a.workers[i].GetOutput()
		if err != nil {
			return nil, err
		}
		result = append(result, output)
		a.next.GetPool().PutStepVector(vector)
	}

	return result, nil
}

func (a *aggregate) initializeTables(ctx context.Context) error {
	var (
		tables []aggregateTable
		series []labels.Labels
		err    error
	)

	if a.by && len(a.labels) == 0 {
		tables, series, err = a.initializeVectorizedTables(ctx)
	} else {
		tables, series, err = a.initializeScalarTables(ctx)
	}
	if err != nil {
		return err
	}
	a.tables = tables
	a.series = series
	a.workers.Start(ctx)

	return nil
}

func (a *aggregate) workerTask(workerID int, vector model.StepVector) model.StepVector {
	table := a.tables[workerID]
	table.aggregate(vector)
	return table.toVector(a.vectorPool)
}

func (a *aggregate) initializeVectorizedTables(ctx context.Context) ([]aggregateTable, []labels.Labels, error) {
	tables, err := newVectorizedTables(a.stepsBatch, a.aggregation, a.param)
	if errors.Is(err, parse.ErrNotSupportedExpr) {
		return a.initializeScalarTables(ctx)
	}

	if err != nil {
		return nil, nil, err
	}

	return tables, []labels.Labels{{}}, nil
}

func (a *aggregate) initializeScalarTables(ctx context.Context) ([]aggregateTable, []labels.Labels, error) {
	series, err := a.next.Series(ctx)
	if err != nil {
		return nil, nil, err
	}

	inputCache := make([]uint64, len(series))
	outputMap := make(map[uint64]*model.Series)
	outputCache := make([]*model.Series, 0)
	buf := make([]byte, 1024)
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
	a.vectorPool.SetStepSize(len(outputCache))
	tables := newScalarTables(a.stepsBatch, inputCache, outputCache, a.newAccumulator)

	series = make([]labels.Labels, len(outputCache))
	for i := 0; i < len(outputCache); i++ {
		series[i] = outputCache[i].Metric
	}

	return tables, series, nil
}
