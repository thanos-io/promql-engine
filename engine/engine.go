package engine

import (
	"time"

	"github.com/fpetkovski/promql-engine/model"

	"github.com/fpetkovski/promql-engine/executionplan"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"
	v1 "github.com/prometheus/prometheus/web/api/v1"
)

type engine struct {
	logger promql.QueryLogger
	pool   *model.VectorPool

	lookbackDelta time.Duration
}

func New() v1.QueryEngine {
	return &engine{
		pool:          model.NewPool(),
		lookbackDelta: 5 * time.Minute,
	}
}

func (e *engine) SetQueryLogger(l promql.QueryLogger) {
	e.logger = l
}

func (e *engine) NewInstantQuery(q storage.Queryable, opts *promql.QueryOpts, qs string, ts time.Time) (promql.Query, error) {
	expr, err := parser.ParseExpr(qs)
	if err != nil {
		return nil, err
	}

	plan, err := executionplan.New(expr, q, ts, ts, 0)
	if err != nil {
		return nil, err
	}

	return newInstantQuery(plan), nil
}

func (e *engine) NewRangeQuery(q storage.Queryable, opts *promql.QueryOpts, qs string, start, end time.Time, interval time.Duration) (promql.Query, error) {
	expr, err := parser.ParseExpr(qs)
	if err != nil {
		return nil, err
	}

	plan, err := executionplan.New(expr, q, start, end, interval)
	if err != nil {
		return nil, err
	}

	return newRangeQuery(plan, e.pool), nil
}
