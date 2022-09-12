// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine

import (
	"time"

	"github.com/efficientgo/core/errors"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/thanos-community/promql-engine/physicalplan"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"
	v1 "github.com/prometheus/prometheus/web/api/v1"
)

type engine struct {
	logger promql.QueryLogger

	lookbackDelta time.Duration
}

type Opts struct {
	promql.EngineOpts

	// DisableFallback enables mode where engine returns error if some expression of feature is not yet implemented
	// in the new engine, instead of falling back to prometheus engine.
	DisableFallback bool
}

func New(opts Opts) v1.QueryEngine {
	if opts.Logger == nil {
		opts.Logger = log.NewNopLogger()
	}
	if opts.LookbackDelta == 0 {
		opts.LookbackDelta = 5 * time.Minute
		if l := opts.Logger; l != nil {
			level.Debug(l).Log("msg", "lookback delta is zero, setting to default value", "value", 5*time.Minute)
		}
	}

	core := &engine{
		lookbackDelta: opts.LookbackDelta,
	}
	if opts.DisableFallback {
		return core
	}

	return &compatibilityEngine{
		core: core,
		prom: promql.NewEngine(opts.EngineOpts),
	}
}

type compatibilityEngine struct {
	core *engine
	prom *promql.Engine
}

func (e *compatibilityEngine) SetQueryLogger(l promql.QueryLogger) {
	e.core.SetQueryLogger(l)
	e.prom.SetQueryLogger(l)
}

func (e *compatibilityEngine) NewInstantQuery(q storage.Queryable, opts *promql.QueryOpts, qs string, ts time.Time) (promql.Query, error) {
	ret, err := e.core.NewInstantQuery(q, opts, qs, ts)
	if triggerFallback(err) {
		return e.prom.NewInstantQuery(q, opts, qs, ts)
	}

	return ret, err
}

func (e *compatibilityEngine) NewRangeQuery(q storage.Queryable, opts *promql.QueryOpts, qs string, start, end time.Time, interval time.Duration) (promql.Query, error) {
	ret, err := e.core.NewRangeQuery(q, opts, qs, start, end, interval)
	if triggerFallback(err) {
		return e.prom.NewRangeQuery(q, opts, qs, start, end, interval)
	}

	return ret, err
}

func (e *engine) SetQueryLogger(l promql.QueryLogger) {
	e.logger = l
}

func triggerFallback(err error) bool {
	return errors.Is(err, physicalplan.ErrNotSupportedExpr) || errors.Is(err, errNotImplemented)
}

var errNotImplemented = errors.New("not implemented")

func (e *engine) NewInstantQuery(q storage.Queryable, opts *promql.QueryOpts, qs string, ts time.Time) (promql.Query, error) {
	return nil, errors.Wrap(errNotImplemented, "instant query")

	//expr, err := parser.ParseExpr(qs)
	//if err != nil {
	//	return nil, err
	//}

	//plan, err := physicalplan.New(expr, q, ts, ts, 0)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return newInstantQuery(plan), nil
}

func (e *engine) NewRangeQuery(q storage.Queryable, opts *promql.QueryOpts, qs string, start, end time.Time, interval time.Duration) (promql.Query, error) {
	expr, err := parser.ParseExpr(qs)
	if err != nil {
		return nil, err
	}

	plan, err := physicalplan.New(expr, q, start, end, interval)
	if err != nil {
		return nil, err
	}

	return newRangeQuery(plan), nil
}
