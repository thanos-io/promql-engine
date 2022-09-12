// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package physicalplan

import (
	"runtime"
	"time"

	"github.com/thanos-community/promql-engine/physicalplan/model"

	"github.com/thanos-community/promql-engine/physicalplan/aggregate"
	"github.com/thanos-community/promql-engine/physicalplan/exchange"
	"github.com/thanos-community/promql-engine/physicalplan/scan"

	"github.com/efficientgo/core/errors"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"
)

const stepsBatch = 10

var ErrNotSupportedExpr = errors.New("unsupported expression")

// New creates new physical query execution plan for a given query expression.
func New(expr parser.Expr, storage storage.Queryable, mint, maxt time.Time, step time.Duration) (model.VectorOperator, error) {
	return newOperator(expr, storage, mint, maxt, step)
}

func newOperator(expr parser.Expr, storage storage.Queryable, mint, maxt time.Time, step time.Duration) (model.VectorOperator, error) {
	switch e := expr.(type) {
	case *parser.AggregateExpr:
		next, err := newOperator(e.Expr, storage, mint, maxt, step)
		if err != nil {
			return nil, err
		}
		a, err := aggregate.NewHashAggregate(model.NewVectorPool(stepsBatch), next, e.Op, !e.Without, e.Grouping, stepsBatch)
		if err != nil {
			return nil, err
		}
		return exchange.NewConcurrent(a, 2), nil

	case *parser.VectorSelector:
		filter := scan.NewSeriesFilter(storage, mint, maxt, 0, e.LabelMatchers)
		numShards := runtime.NumCPU() / 2
		operators := make([]model.VectorOperator, 0, numShards)
		for i := 0; i < numShards; i++ {
			operators = append(operators, exchange.NewConcurrent(scan.NewVectorSelector(model.NewVectorPool(stepsBatch), filter, mint, maxt, step, stepsBatch, i, numShards), 2))
		}
		return exchange.NewCoalesce(model.NewVectorPool(stepsBatch), operators...), nil

	case *parser.Call:
		switch t := e.Args[0].(type) {
		case *parser.MatrixSelector:
			vs := t.VectorSelector.(*parser.VectorSelector)
			call, err := scan.NewFunctionCall(e.Func, t.Range)
			if err != nil {
				return nil, err
			}

			filter := scan.NewSeriesFilter(storage, mint, maxt, t.Range, vs.LabelMatchers)
			numShards := runtime.NumCPU() / 2
			operators := make([]model.VectorOperator, 0, numShards)
			for i := 0; i < numShards; i++ {
				selector := scan.NewMatrixSelector(model.NewVectorPool(stepsBatch), filter, call, mint, maxt, stepsBatch, step, t.Range, i, numShards)
				operators = append(operators, exchange.NewConcurrent(selector, 2))
			}
			return exchange.NewCoalesce(model.NewVectorPool(stepsBatch), operators...), nil
		default:
			return nil, errors.Wrapf(ErrNotSupportedExpr, "got: %s", t)
		}
	default:
		return nil, errors.Wrapf(ErrNotSupportedExpr, "got: %s", e)
	}
}
