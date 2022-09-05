package engine

import (
	"fpetkovski/promql-engine/executionplan"
	"time"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"
	v1 "github.com/prometheus/prometheus/web/api/v1"
)

type engine struct {
	logger promql.QueryLogger
}

func New() v1.QueryEngine {
	return &engine{}
}

func (e *engine) SetQueryLogger(l promql.QueryLogger) {
	e.logger = e.logger
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

	return newQuery(plan), nil
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

	return newQuery(plan), nil
}
