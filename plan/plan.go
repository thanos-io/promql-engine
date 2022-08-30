package plan

import (
	"context"
	"fmt"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"
	"time"
)

type Operator interface {
	Next(ctx context.Context) (<-chan promql.Matrix, error)
}

type plan struct {
	rootOperator Operator
}

func New(expr parser.Expr, storage storage.Storage, mint, maxt time.Time, step time.Duration) (Operator, error) {
	return newOperator(expr, storage, mint, maxt, step)
}

func newOperator(expr parser.Expr, storage storage.Storage, mint, maxt time.Time, step time.Duration) (Operator, error) {
	switch e := expr.(type) {
	case *parser.AggregateExpr:
		next, err := newOperator(e.Expr, storage, mint, maxt, step)
		if err != nil {
			return nil, err
		}
		return NewAggregate(next, e.Op, !e.Without, e.Grouping)
	case *parser.VectorSelector:
		return NewSelector(storage, e.LabelMatchers, nil, mint, maxt, step, 0), nil
	default:
		return nil, fmt.Errorf("unsupported expression %s", e)
	}
}
