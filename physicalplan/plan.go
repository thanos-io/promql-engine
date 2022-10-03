// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package physicalplan

import (
	"runtime"
	"time"

	"github.com/thanos-community/promql-engine/physicalplan/aggregate"
	"github.com/thanos-community/promql-engine/physicalplan/binary"
	"github.com/thanos-community/promql-engine/physicalplan/exchange"
	"github.com/thanos-community/promql-engine/physicalplan/model"
	"github.com/thanos-community/promql-engine/physicalplan/parse"
	"github.com/thanos-community/promql-engine/physicalplan/scan"
	"github.com/thanos-community/promql-engine/physicalplan/unary"

	"github.com/efficientgo/core/errors"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"
)

const stepsBatch = 10

// New creates new physical query execution plan for a given query expression.
func New(expr parser.Expr, storage storage.Queryable, mint, maxt time.Time, step, lookbackDelta time.Duration) (model.VectorOperator, error) {
	return newCancellableOperator(expr, storage, mint, maxt, step, lookbackDelta)
}

func newCancellableOperator(expr parser.Expr, storage storage.Queryable, mint, maxt time.Time, step, lookbackDelta time.Duration) (*exchange.CancellableOperator, error) {
	operator, err := newOperator(expr, storage, mint, maxt, step, lookbackDelta)
	if err != nil {
		return nil, err
	}

	return exchange.NewCancellable(operator), nil
}

func newOperator(expr parser.Expr, storage storage.Queryable, mint time.Time, maxt time.Time, step, lookbackDelta time.Duration) (model.VectorOperator, error) {
	switch e := expr.(type) {
	case *parser.NumberLiteral:
		return scan.NewNumberLiteralSelector(model.NewVectorPool(stepsBatch), mint, maxt, step, stepsBatch, e.Val, nil), nil

	case *parser.VectorSelector:
		filter := scan.NewSeriesFilter(storage, mint, maxt, 0, lookbackDelta, e.LabelMatchers)
		numShards := runtime.GOMAXPROCS(0) / 2
		if numShards < 1 {
			numShards = 1
		}
		operators := make([]model.VectorOperator, 0, numShards)
		for i := 0; i < numShards; i++ {
			operator := exchange.NewConcurrent(
				exchange.NewCancellable(
					scan.NewVectorSelector(
						model.NewVectorPool(stepsBatch), filter, mint, maxt, step, lookbackDelta, stepsBatch, i, numShards)), 2)
			operators = append(operators, operator)
		}

		return exchange.NewCoalesce(model.NewVectorPool(stepsBatch), operators...), nil

	case *parser.Call:
		if len(e.Args) != 1 {
			return nil, errors.Wrapf(parse.ErrNotSupportedExpr, "got: %s", e)
		}

		switch t := e.Args[0].(type) {
		case *parser.MatrixSelector:
			vs := t.VectorSelector.(*parser.VectorSelector)
			call, err := scan.NewFunctionCall(e.Func)
			if err != nil {
				return nil, err
			}

			lookbackDelta = maxDuration(lookbackDelta, t.Range)
			filter := scan.NewSeriesFilter(storage, mint, maxt, t.Range, lookbackDelta, vs.LabelMatchers)
			numShards := runtime.GOMAXPROCS(0) / 2
			if numShards < 1 {
				numShards = 1
			}
			operators := make([]model.VectorOperator, 0, numShards)
			for i := 0; i < numShards; i++ {
				operator := exchange.NewConcurrent(
					exchange.NewCancellable(
						scan.NewMatrixSelector(
							model.NewVectorPool(stepsBatch), filter, e, call, mint, maxt, stepsBatch, step, t.Range, i, numShards)), 2)
				operators = append(operators, operator)
			}

			return exchange.NewCoalesce(model.NewVectorPool(stepsBatch), operators...), nil
		case *parser.NumberLiteral:
			call, err := scan.NewFunctionCall(e.Func)
			if err != nil {
				return nil, err
			}

			return exchange.NewCancellable(
				scan.NewNumberLiteralSelector(
					model.NewVectorPool(stepsBatch), mint, maxt, step, stepsBatch, t.Val, call)), nil
		default:
			return nil, errors.Wrapf(parse.ErrNotSupportedExpr, "got: %s", t)
		}

	case *parser.AggregateExpr:
		next, err := newCancellableOperator(e.Expr, storage, mint, maxt, step, lookbackDelta)
		if err != nil {
			return nil, err
		}

		a, err := aggregate.NewHashAggregate(model.NewVectorPool(stepsBatch), next, e.Op, !e.Without, e.Grouping, stepsBatch)
		if err != nil {
			return nil, err
		}

		return exchange.NewConcurrent(exchange.NewCancellable(a), 2), nil

	case *parser.BinaryExpr:
		if e.LHS.Type() == parser.ValueTypeScalar || e.RHS.Type() == parser.ValueTypeScalar {
			return newScalarBinaryOperator(e, storage, mint, maxt, step, lookbackDelta)
		}

		return newVectorBinaryOperator(e, storage, mint, maxt, step, lookbackDelta)

	case *parser.ParenExpr:
		return newCancellableOperator(e.Expr, storage, mint, maxt, step, lookbackDelta)

	case *parser.StringLiteral:
		return nil, nil

	case *parser.UnaryExpr:
		next, err := newOperator(e.Expr, storage, mint, maxt, step, lookbackDelta)
		if err != nil {
			return nil, err
		}
		switch e.Op {
		case parser.ADD:
			return next, nil
		case parser.SUB:
			return unary.NewUnaryNegation(next, stepsBatch)
		default:
			// This shouldn't happen as Op was validated when parsing already
			// https://github.com/prometheus/prometheus/blob/v2.38.0/promql/parser/parse.go#L573.
			return nil, errors.Wrapf(parse.ErrNotSupportedExpr, "got: %s", e)
		}

	default:
		return nil, errors.Wrapf(parse.ErrNotSupportedExpr, "got: %s", e)
	}
}

func newVectorBinaryOperator(e *parser.BinaryExpr, storage storage.Queryable, mint time.Time, maxt time.Time, step, lookbackDelta time.Duration) (model.VectorOperator, error) {
	leftOperator, err := newCancellableOperator(e.LHS, storage, mint, maxt, step, lookbackDelta)
	if err != nil {
		return nil, err
	}
	rightOperator, err := newCancellableOperator(e.RHS, storage, mint, maxt, step, lookbackDelta)
	if err != nil {
		return nil, err
	}
	return binary.NewVectorOperator(model.NewVectorPool(stepsBatch), leftOperator, rightOperator, e.VectorMatching, e.Op)
}

func newScalarBinaryOperator(e *parser.BinaryExpr, storage storage.Queryable, mint time.Time, maxt time.Time, step, lookbackDelta time.Duration) (model.VectorOperator, error) {
	lhs, err := newCancellableOperator(e.LHS, storage, mint, maxt, step, lookbackDelta)
	if err != nil {
		return nil, err
	}
	rhs, err := newCancellableOperator(e.RHS, storage, mint, maxt, step, lookbackDelta)
	if err != nil {
		return nil, err
	}

	if e.LHS.Type() == parser.ValueTypeScalar {
		return binary.NewScalar(model.NewVectorPool(stepsBatch), rhs, lhs, e.Op, true)
	}
	return binary.NewScalar(model.NewVectorPool(stepsBatch), lhs, rhs, e.Op, false)
}

func maxDuration(a, b time.Duration) time.Duration {
	if a.Milliseconds() >= b.Milliseconds() {
		return a
	}
	return b
}
