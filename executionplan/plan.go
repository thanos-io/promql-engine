package executionplan

import (
	"context"
	"fmt"
	"time"

	"github.com/fpetkovski/promql-engine/model"

	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"
)

type VectorOperator interface {
	Next(ctx context.Context) ([]model.StepVector, error)
}

func New(expr parser.Expr, storage storage.Queryable, mint, maxt time.Time, step time.Duration) (VectorOperator, error) {
	pool := model.NewPool()
	return newOperator(pool, expr, storage, mint, maxt, step)
}

func newOperator(pool *model.VectorPool, expr parser.Expr, storage storage.Queryable, mint, maxt time.Time, step time.Duration) (VectorOperator, error) {
	switch e := expr.(type) {
	case *parser.AggregateExpr:
		next, err := newOperator(pool, e.Expr, storage, mint, maxt, step)
		if err != nil {
			return nil, err
		}
		aggregate, err := NewAggregate(pool, next, e.Op, !e.Without, e.Grouping)
		if err != nil {
			return nil, err
		}
		return concurrent(aggregate), nil

	case *parser.VectorSelector:
		filter := newSeriesFilter(storage, mint, maxt, e.LabelMatchers)
		numShards := 7
		operators := make([]VectorOperator, 0, numShards)
		for i := 0; i < numShards; i++ {
			operators = append(operators, concurrent(NewVectorSelector(pool, filter, mint, maxt, step, i, numShards)))
		}
		return coalesce(operators...), nil

	//case *parser.Call:
	//	switch t := e.Args[0].(type) {
	//	case *parser.MatrixSelector:
	//		vs := t.VectorSelector.(*parser.VectorSelector)
	//		call, err := NewFunctionCall(e.Func, t.Range)
	//		if err != nil {
	//			return nil, err
	//		}
	//		selector := NewMatrixSelector(pool, storage, call, vs.LabelMatchers, nil, mint, maxt, step, t.Range)
	//		return concurrent(selector), nil
	//	}
	//	return nil, fmt.Errorf("unsupported expression %s", e)

	default:
		return nil, fmt.Errorf("unsupported expression %s", e)
	}
}
