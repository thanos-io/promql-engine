// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package physicalplan

import (
	"fmt"
	"runtime"
	"time"

	"github.com/thanos-community/promql-engine/physicalplan/model"

	"github.com/thanos-community/promql-engine/physicalplan/aggregate"
	"github.com/thanos-community/promql-engine/physicalplan/exchange"
	"github.com/thanos-community/promql-engine/physicalplan/scan"

	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"
)

const stepsBatch = 10

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
		a, err := aggregate.NewHashAggregate(model.NewVectorPool(), next, e.Op, !e.Without, e.Grouping, stepsBatch)
		if err != nil {
			return nil, err
		}
		return exchange.NewConcurrent(a, 2), nil

	case *parser.VectorSelector:
		filter := scan.NewSeriesFilter(storage, mint, maxt, 0, e.LabelMatchers)
		numShards := runtime.NumCPU() / 2
		operators := make([]model.VectorOperator, 0, numShards)
		for i := 0; i < numShards; i++ {
			operators = append(operators, exchange.NewConcurrent(scan.NewVectorSelector(model.NewVectorPool(), filter, mint, maxt, step, stepsBatch, i, numShards), 3))
		}
		return exchange.NewCoalesce(model.NewVectorPool(), operators...), nil

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
				operators = append(operators, exchange.NewConcurrent(
					scan.NewMatrixSelector(model.NewVectorPool(), filter, call, mint, maxt, stepsBatch, step, t.Range, i, numShards), 3),
				)
			}

			return exchange.NewCoalesce(model.NewVectorPool(), operators...), nil
		}
		return nil, fmt.Errorf("unsupported expression %s", e)

	default:
		return nil, fmt.Errorf("unsupported expression %s", e)
	}
}
