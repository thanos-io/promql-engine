package executionplan

import (
	"fmt"
	"time"

	"github.com/fpetkovski/promql-engine/operators/model"

	"github.com/fpetkovski/promql-engine/operators/aggregate"
	"github.com/fpetkovski/promql-engine/operators/exchange"
	"github.com/fpetkovski/promql-engine/operators/scan"

	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"
)

func New(expr parser.Expr, storage storage.Queryable, mint, maxt time.Time, step time.Duration) (model.Vector, error) {
	return newOperator(expr, storage, mint, maxt, step)
}

func newOperator(expr parser.Expr, storage storage.Queryable, mint, maxt time.Time, step time.Duration) (model.Vector, error) {
	switch e := expr.(type) {
	case *parser.AggregateExpr:
		next, err := newOperator(e.Expr, storage, mint, maxt, step)
		if err != nil {
			return nil, err
		}
		a, err := aggregate.NewHashAggregate(model.NewVectorPool(), next, e.Op, !e.Without, e.Grouping)
		if err != nil {
			return nil, err
		}
		return exchange.NewConcurrent(a, 2), nil

	case *parser.VectorSelector:
		filter := scan.NewSeriesFilter(storage, mint, maxt, e.LabelMatchers)
		numShards := 4
		operators := make([]model.Vector, 0, numShards)
		for i := 0; i < numShards; i++ {
			operators = append(operators, exchange.NewConcurrent(scan.NewVectorSelector(model.NewVectorPool(), filter, mint, maxt, step, i, numShards), 3))
		}
		return exchange.NewCoalesce(model.NewVectorPool(), operators...), nil

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
