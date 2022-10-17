// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"sort"
	"time"

	"github.com/efficientgo/core/errors"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/util/stats"
	v1 "github.com/prometheus/prometheus/web/api/v1"
	"github.com/thanos-community/promql-engine/executor"
	"github.com/thanos-community/promql-engine/executor/model"
	"github.com/thanos-community/promql-engine/executor/parse"
	"github.com/thanos-community/promql-engine/logicalplan"
)

type engine struct {
	logger           log.Logger
	debugWriter      io.Writer
	lookbackDelta    time.Duration
	enableOptimizers bool
}

type Opts struct {
	promql.EngineOpts

	// DisableOptimizers disables query optimizations using logicalPlan.DefaultOptimizers.
	DisableOptimizers bool

	// DisableFallback enables mode where engine returns error if some expression of feature is not yet implemented
	// in the new engine, instead of falling back to prometheus engine.
	DisableFallback bool

	// DebugWriter specifies output for debug (multi-line) information meant for humans debugging the engine.
	// If nil, nothing will be printed.
	// NOTE: Users will not check the errors, debug writing is best effort.
	DebugWriter io.Writer
}

func New(opts Opts) v1.QueryEngine {
	if opts.Logger == nil {
		opts.Logger = log.NewNopLogger()
	}
	if opts.LookbackDelta == 0 {
		opts.LookbackDelta = 5 * time.Minute
		level.Debug(opts.Logger).Log("msg", "lookback delta is zero, setting to default value", "value", 5*time.Minute)
	}

	core := &engine{
		debugWriter:      opts.DebugWriter,
		logger:           opts.Logger,
		lookbackDelta:    opts.LookbackDelta,
		enableOptimizers: !opts.DisableOptimizers,
	}
	return &compatibilityEngine{
		core:            core,
		disableFallback: opts.DisableFallback,
		prom:            promql.NewEngine(opts.EngineOpts),
		queries: promauto.With(opts.Reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "promql_engine_queries_total",
				Help: "Number of PromQL queries.",
			}, []string{"fallback"},
		),
	}
}

type compatibilityEngine struct {
	core            *engine
	prom            *promql.Engine
	queries         *prometheus.CounterVec
	disableFallback bool
}

func (e *compatibilityEngine) SetQueryLogger(l promql.QueryLogger) {
	e.core.SetQueryLogger(l)
	e.prom.SetQueryLogger(l)
}

func (e *compatibilityEngine) NewInstantQuery(q storage.Queryable, opts *promql.QueryOpts, qs string, ts time.Time) (promql.Query, error) {
	expr, err := parser.ParseExpr(qs)
	if err != nil {
		return nil, err
	}

	ret, err := e.core.NewExecutor(q, opts, expr, ts, ts, 0)
	if !e.disableFallback && triggerFallback(err) {
		e.queries.WithLabelValues("true").Inc()
		return e.prom.NewInstantQuery(q, opts, qs, ts)
	}
	e.queries.WithLabelValues("false").Inc()

	return &compatibilityQuery{
		executor: ret,
		expr:     expr,
		logger:   e.core.logger,
		ts:       ts,
	}, err
}

func (e *compatibilityEngine) NewRangeQuery(q storage.Queryable, opts *promql.QueryOpts, qs string, start, end time.Time, interval time.Duration) (promql.Query, error) {
	expr, err := parser.ParseExpr(qs)
	if err != nil {
		return nil, err
	}

	ret, err := e.core.NewExecutor(q, opts, expr, start, end, interval)
	if !e.disableFallback && triggerFallback(err) {
		e.queries.WithLabelValues("true").Inc()
		return e.prom.NewRangeQuery(q, opts, qs, start, end, interval)
	}
	e.queries.WithLabelValues("false").Inc()

	return &compatibilityQuery{
		executor: ret,
		expr:     expr,
		logger:   e.core.logger,
	}, err
}

type compatibilityQuery struct {
	executor model.VectorOperator
	logger   log.Logger
	expr     parser.Expr

	ts time.Time // Available if it's instant query.

	cancel context.CancelFunc
}

func (q *compatibilityQuery) Exec(ctx context.Context) (ret *promql.Result) {
	// Handle case with strings early on as this does not need us to process samples.
	// TODO(saswatamcode): Modify models.StepVector to support all types and check during executor creation.
	switch e := q.expr.(type) {
	case *parser.StringLiteral:
		return &promql.Result{Value: promql.String{V: e.Val, T: q.ts.UnixMilli()}}
	}

	ret = &promql.Result{
		Value: promql.Vector{},
	}
	defer recoverEngine(q.logger, q.expr, &ret.Err)

	ctx, cancel := context.WithCancel(ctx)
	q.cancel = cancel

	resultSeries, err := q.executor.Series(ctx)
	if err != nil {
		return newErrResult(ret, err)
	}

	series := make([]promql.Series, len(resultSeries))
	for i := 0; i < len(resultSeries); i++ {
		series[i].Metric = resultSeries[i]
		series[i].Points = make([]promql.Point, 0, 121) // Typically 1h of data.
	}

loop:
	for {
		select {
		case <-ctx.Done():
			return newErrResult(ret, ctx.Err())
		default:
			r, err := q.executor.Next(ctx)
			if err != nil {
				return newErrResult(ret, err)
			}
			if r == nil {
				break loop
			}

			for _, vector := range r {
				for i, s := range vector.SampleIDs {
					series[s].Points = append(series[s].Points, promql.Point{
						T: vector.T,
						V: vector.Samples[i],
					})
				}
				q.executor.GetPool().PutStepVector(vector)
			}
			q.executor.GetPool().PutVectors(r)
		}
	}

	// For range query we expect always a Matrix value type.
	if q.ts.Equal(time.Time{}) {
		resultMatrix := make(promql.Matrix, 0, len(series))
		for _, s := range series {
			if len(s.Points) == 0 {
				continue
			}
			resultMatrix = append(resultMatrix, s)
		}
		sort.Sort(resultMatrix)
		ret.Value = resultMatrix
		return ret
	}

	var result parser.Value
	switch q.expr.Type() {
	case parser.ValueTypeMatrix:
		result = promql.Matrix(series)
	case parser.ValueTypeVector:
		// Convert matrix with one value per series into vector.
		vector := make(promql.Vector, 0, len(resultSeries))
		for i := range series {
			if len(series[i].Points) == 0 {
				continue
			}
			// Point might have a different timestamp, force it to the evaluation
			// timestamp as that is when we ran the evaluation.
			vector = append(vector, promql.Sample{
				Metric: series[i].Metric,
				Point: promql.Point{
					V: series[i].Points[0].V,
					T: q.ts.UnixMilli(),
				},
			})
		}
		result = vector
	case parser.ValueTypeScalar:
		result = promql.Scalar{V: series[0].Points[0].V, T: q.ts.UnixMilli()}
	default:
		panic(errors.Newf("new.Engine.exec: unexpected expression type %q", q.expr.Type()))
	}

	ret.Value = result
	return ret
}

func newErrResult(r *promql.Result, err error) *promql.Result {
	if r == nil {
		r = &promql.Result{}
	}
	if r.Err == nil && err != nil {
		r.Err = err
	}
	return r
}

func (q *compatibilityQuery) Statement() parser.Statement { return nil }

func (q *compatibilityQuery) Stats() *stats.Statistics { return &stats.Statistics{} }

func (q *compatibilityQuery) Close() { q.Cancel() }

func (q *compatibilityQuery) String() string { return q.expr.String() }

func (q *compatibilityQuery) Cancel() {
	if q.cancel != nil {
		q.cancel()
		q.cancel = nil
	}
}

func (e *engine) SetQueryLogger(l promql.QueryLogger) {
	e.logger = l
}

func triggerFallback(err error) bool {
	return errors.Is(err, parse.ErrNotSupportedExpr) || errors.Is(err, errNotImplemented)
}

var errNotImplemented = errors.New("not implemented")

func (e *engine) NewExecutor(q storage.Queryable, _ *promql.QueryOpts, expr parser.Expr, start, end time.Time, step time.Duration) (model.VectorOperator, error) {
	// if step = 0 it is an instant query.
	// Use same check as Prometheus for range queries.
	if step > 0 && expr.Type() != parser.ValueTypeVector && expr.Type() != parser.ValueTypeScalar {
		return nil, errors.Newf("invalid expression type %q for range query, must be Scalar or instant Vector", parser.DocumentedType(expr.Type()))
	}

	// TODO(bwplotka): We plan to have single plan that can be optimized in both physical and logical context. Rename.
	logicalPlan := logicalplan.New(expr, start, end)
	if e.enableOptimizers {
		logicalPlan = logicalPlan.Optimize(logicalplan.DefaultOptimizers)
	}
	executor, err := executor.New(logicalPlan.Expr(), q, start, end, step, e.lookbackDelta)
	if err != nil {
		return nil, err
	}

	if e.debugWriter != nil {
		explain(e.debugWriter, executor, "", "")
	}
	return executor, nil
}

func recoverEngine(logger log.Logger, expr parser.Expr, errp *error) {
	e := recover()
	if e == nil {
		return
	}

	switch err := e.(type) {
	case runtime.Error:
		// Print the stack trace but do not inhibit the running application.
		buf := make([]byte, 64<<10)
		buf = buf[:runtime.Stack(buf, false)]

		level.Error(logger).Log("msg", "runtime panic in engine", "expr", expr.String(), "err", e, "stacktrace", string(buf))
		*errp = fmt.Errorf("unexpected error: %w", err)
	}
}

func explain(w io.Writer, o model.VectorOperator, indent, indentNext string) {
	me, next := o.Explain()
	_, _ = w.Write([]byte(indent))
	_, _ = w.Write([]byte(me))
	if len(next) == 0 {
		_, _ = w.Write([]byte("\n"))
		return
	}

	if me == "[*CancellableOperator]" {
		_, _ = w.Write([]byte(": "))
		explain(w, next[0], "", indentNext)
		return
	}
	_, _ = w.Write([]byte(":\n"))

	for i, n := range next {
		if i == len(next)-1 {
			explain(w, n, indentNext+"└──", indentNext+"   ")
		} else {
			explain(w, n, indentNext+"├──", indentNext+"│  ")
		}
	}
}
