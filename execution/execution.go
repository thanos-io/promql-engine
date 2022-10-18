// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

// Copyright 2013 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package execution

import (
	"runtime"
	"time"

	"github.com/efficientgo/core/errors"

	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"

	"github.com/prometheus/prometheus/model/labels"

	"github.com/thanos-community/promql-engine/execution/aggregate"
	"github.com/thanos-community/promql-engine/execution/binary"
	"github.com/thanos-community/promql-engine/execution/exchange"
	"github.com/thanos-community/promql-engine/execution/model"
	"github.com/thanos-community/promql-engine/execution/parse"
	"github.com/thanos-community/promql-engine/execution/scan"
	"github.com/thanos-community/promql-engine/execution/step_invariant"
	engstore "github.com/thanos-community/promql-engine/execution/storage"
	"github.com/thanos-community/promql-engine/execution/unary"
	"github.com/thanos-community/promql-engine/logicalplan"
	"github.com/thanos-community/promql-engine/query"
)

const stepsBatch = 10

// New creates new physical query execution for a given query expression which represents logical plan.
// TODO(bwplotka): Add definition (could be parameters for each execution operator) we can optimize - it would represent physical plan.
func New(expr parser.Expr, storage storage.Queryable, mint, maxt time.Time, step, lookbackDelta time.Duration) (model.VectorOperator, error) {
	opts := &query.Options{
		Start:         mint,
		End:           maxt,
		Step:          step,
		LookbackDelta: lookbackDelta,
		StepsBatch:    stepsBatch,
	}
	selectorPool := engstore.NewSelectorPool(storage)
	return newCancellableOperator(expr, selectorPool, opts)
}

func newCancellableOperator(expr parser.Expr, selectorPool *engstore.SelectorPool, opts *query.Options) (*exchange.CancellableOperator, error) {
	operator, err := newOperator(expr, selectorPool, opts)
	if err != nil {
		return nil, err
	}

	return exchange.NewCancellable(operator), nil
}

func newOperator(expr parser.Expr, storage *engstore.SelectorPool, opts *query.Options) (model.VectorOperator, error) {
	switch e := expr.(type) {
	case *parser.NumberLiteral:
		return scan.NewNumberLiteralSelector(model.NewVectorPool(stepsBatch), opts, e.Val), nil

	case *parser.VectorSelector:
		start, end := getTimeRangesForVectorSelector(e, opts, 0)
		filter := storage.GetSelector(start, end, e.LabelMatchers)
		return newShardedVectorSelector(filter, opts, e.Offset)

	case *logicalplan.FilteredSelector:
		start, end := getTimeRangesForVectorSelector(e.VectorSelector, opts, 0)
		selector := storage.GetFilteredSelector(start, end, e.LabelMatchers, e.Filters)
		return newShardedVectorSelector(selector, opts, e.Offset)

	case *parser.Call:
		if len(e.Args) != 1 {
			return nil, errors.Wrapf(parse.ErrNotSupportedExpr, "got: %s", e)
		}

		switch t := e.Args[0].(type) {
		case *parser.MatrixSelector:
			vs, filters, err := unpackVectorSelector(t)
			if err != nil {
				return nil, err
			}
			call, err := scan.NewFunctionCall(e.Func)
			if err != nil {
				return nil, err
			}

			start, end := getTimeRangesForVectorSelector(vs, opts, t.Range)
			filter := storage.GetFilteredSelector(start, end, vs.LabelMatchers, filters)

			numShards := runtime.GOMAXPROCS(0) / 2
			if numShards < 1 {
				numShards = 1
			}

			operators := make([]model.VectorOperator, 0, numShards)
			for i := 0; i < numShards; i++ {
				operator := exchange.NewConcurrent(
					exchange.NewCancellable(
						scan.NewMatrixSelector(model.NewVectorPool(stepsBatch), filter, e, call, opts, t.Range, vs.Offset, i, numShards),
					), 2)
				operators = append(operators, operator)
			}

			return exchange.NewCoalesce(model.NewVectorPool(stepsBatch), operators...), nil

		case *parser.NumberLiteral:
			l, err := scan.NewNumberLiteralSelectorWithFunc(model.NewVectorPool(stepsBatch), opts, t.Val, e.Func)
			if err != nil {
				return nil, err
			}
			return exchange.NewCancellable(l), nil
		default:
			return nil, errors.Wrapf(parse.ErrNotSupportedExpr, "got: %s", t)
		}

	case *parser.AggregateExpr:
		next, err := newCancellableOperator(e.Expr, storage, opts)
		if err != nil {
			return nil, err
		}
		a, err := aggregate.NewHashAggregate(model.NewVectorPool(stepsBatch), next, e.Op, e.Param, !e.Without, e.Grouping, stepsBatch)
		if err != nil {
			return nil, err
		}

		return exchange.NewConcurrent(exchange.NewCancellable(a), 2), nil

	case *parser.BinaryExpr:
		if e.LHS.Type() == parser.ValueTypeScalar || e.RHS.Type() == parser.ValueTypeScalar {
			return newScalarBinaryOperator(e, storage, opts)
		}

		return newVectorBinaryOperator(e, storage, opts)

	case *parser.ParenExpr:
		return newCancellableOperator(e.Expr, storage, opts)

	case *parser.UnaryExpr:
		next, err := newCancellableOperator(e.Expr, storage, opts)
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

	case *parser.StepInvariantExpr:
		next, err := newCancellableOperator(e.Expr, storage, opts.WithEndTime(opts.Start))
		if err != nil {
			return nil, err
		}
		return step_invariant.NewStepInvariantOperator(model.NewVectorPool(stepsBatch), next, e.Expr, opts)

	default:
		return nil, errors.Wrapf(parse.ErrNotSupportedExpr, "got: %s", e)
	}
}

func unpackVectorSelector(t *parser.MatrixSelector) (*parser.VectorSelector, []*labels.Matcher, error) {
	switch t := t.VectorSelector.(type) {
	case *parser.VectorSelector:
		return t, nil, nil
	case *logicalplan.FilteredSelector:
		return t.VectorSelector, t.Filters, nil
	default:
		return nil, nil, parse.ErrNotSupportedExpr
	}
}

func newShardedVectorSelector(selector engstore.SeriesSelector, opts *query.Options, offset time.Duration) (model.VectorOperator, error) {
	numShards := runtime.GOMAXPROCS(0) / 2
	if numShards < 1 {
		numShards = 1
	}
	operators := make([]model.VectorOperator, 0, numShards)
	for i := 0; i < numShards; i++ {
		operator := exchange.NewConcurrent(
			exchange.NewCancellable(
				scan.NewVectorSelector(
					model.NewVectorPool(stepsBatch), selector, opts, offset, i, numShards)), 2)
		operators = append(operators, operator)
	}

	return exchange.NewCoalesce(model.NewVectorPool(stepsBatch), operators...), nil
}

func newVectorBinaryOperator(e *parser.BinaryExpr, selectorPool *engstore.SelectorPool, opts *query.Options) (model.VectorOperator, error) {
	leftOperator, err := newCancellableOperator(e.LHS, selectorPool, opts)
	if err != nil {
		return nil, err
	}
	rightOperator, err := newCancellableOperator(e.RHS, selectorPool, opts)
	if err != nil {
		return nil, err
	}
	return binary.NewVectorOperator(model.NewVectorPool(stepsBatch), leftOperator, rightOperator, e.VectorMatching, e.Op)
}

func newScalarBinaryOperator(e *parser.BinaryExpr, selectorPool *engstore.SelectorPool, opts *query.Options) (model.VectorOperator, error) {
	lhs, err := newCancellableOperator(e.LHS, selectorPool, opts)
	if err != nil {
		return nil, err
	}
	rhs, err := newCancellableOperator(e.RHS, selectorPool, opts)
	if err != nil {
		return nil, err
	}

	scalarSide := binary.ScalarSideRight
	if e.LHS.Type() == parser.ValueTypeScalar && e.RHS.Type() == parser.ValueTypeScalar {
		scalarSide = binary.ScalarSideBoth
	} else if e.LHS.Type() == parser.ValueTypeScalar {
		rhs, lhs = lhs, rhs
		scalarSide = binary.ScalarSideLeft
	}

	return binary.NewScalar(model.NewVectorPool(stepsBatch), lhs, rhs, e.Op, scalarSide)
}

// Copy from https://github.com/prometheus/prometheus/blob/v2.39.1/promql/engine.go#L791.
func getTimeRangesForVectorSelector(n *parser.VectorSelector, opts *query.Options, evalRange time.Duration) (int64, int64) {
	start := opts.Start.UnixMilli()
	end := opts.End.UnixMilli()
	if n.Timestamp != nil {
		start = *n.Timestamp
		end = *n.Timestamp
	}
	if evalRange == 0 {
		start -= opts.LookbackDelta.Milliseconds()
	} else {
		start -= evalRange.Milliseconds()
	}
	offset := n.OriginalOffset.Milliseconds()
	return start - offset, end - offset
}
