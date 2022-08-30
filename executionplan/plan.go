package executionplan

import (
	"context"
	"fmt"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"
	"time"
)

type VectorOperator interface {
	Next(ctx context.Context) (<-chan promql.Vector, error)
}

func New(expr parser.Expr, storage storage.Storage, mint, maxt time.Time, step time.Duration) (VectorOperator, error) {
	return newOperator(expr, storage, mint, maxt, step)
}

func newOperator(expr parser.Expr, storage storage.Storage, mint, maxt time.Time, step time.Duration) (VectorOperator, error) {
	switch e := expr.(type) {
	case *parser.AggregateExpr:
		next, err := newOperator(e.Expr, storage, mint, maxt, step)
		if err != nil {
			return nil, err
		}
		return NewAggregate(next, e.Op, !e.Without, e.Grouping)
	case *parser.VectorSelector:
		return NewSelector(storage, e.LabelMatchers, nil, mint, maxt, step), nil
	default:
		return nil, fmt.Errorf("unsupported expression %s", e)
	}
}
