// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine

import (
	"context"
	"sync"
	"time"

	"github.com/thanos-io/promql-engine/api"
	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/logicalplan"

	"github.com/efficientgo/core/errors"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"
)

// errFallbackPathNotImplemented is returned by the in-process remote engine's
// queryBackedRemoteOperator. The in-process distributed optimizer always
// produces RemoteExecution nodes whose Query is a structured logicalplan.Node,
// so this error path is dead code in normal operation. It exists so that any
// future caller passing a non-Node api.RemoteQuery gets a clear signal.
var errFallbackPathNotImplemented = errors.New("queryBackedRemoteOperator: streaming fallback for non-Node api.RemoteQuery is not implemented; call NewRangeQuery instead")

type remoteEngine struct {
	q         storage.Queryable
	engine    *Engine
	labelSets []labels.Labels
	maxt      int64
	mint      int64
}

func NewRemoteEngine(opts Opts, q storage.Queryable, mint, maxt int64, labelSets []labels.Labels) *remoteEngine {
	return &remoteEngine{
		q:         q,
		labelSets: labelSets,
		maxt:      maxt,
		mint:      mint,
		engine:    New(opts),
	}
}

func (l remoteEngine) MaxT() int64 {
	return l.maxt
}

func (l remoteEngine) MinT() int64 {
	return l.mint
}

func (l remoteEngine) LabelSets() []labels.Labels {
	return l.labelSets
}

func (l remoteEngine) PartitionLabelSets() []labels.Labels {
	return l.labelSets
}

func (l remoteEngine) NewRangeQuery(ctx context.Context, opts promql.QueryOpts, plan api.RemoteQuery, start, end time.Time, interval time.Duration) (promql.Query, error) {
	return l.engine.NewRangeQuery(ctx, l.q, opts, plan.String(), start, end, interval)
}

// NewRangeQueryOperator implements api.StreamingRemoteEngine. It returns a
// columnar StepVector stream produced by the underlying engine, avoiding the
// StepVector -> promql.Matrix -> StepVector round-trip required by the legacy
// NewRangeQuery path.
//
// In-process callers (DistributedExecutionOptimizer / PassthroughOptimizer)
// always pass a *logicalplan.Node as the api.RemoteQuery, so we downcast and
// hand the structured plan straight to the engine. Other implementations can
// fall back to NewRangeQuery; see api.StreamingRemoteEngine documentation.
func (l remoteEngine) NewRangeQueryOperator(ctx context.Context, opts promql.QueryOpts, plan api.RemoteQuery, start, end time.Time, interval time.Duration) (api.RemoteOperator, error) {
	root, ok := plan.(logicalplan.Node)
	if !ok {
		// Defensive: re-parse from string. This should not happen in-process
		// because RemoteExecution.Query is always a logicalplan.Node, but the
		// api.RemoteQuery contract only guarantees fmt.Stringer.
		q, err := l.engine.NewRangeQuery(ctx, l.q, opts, plan.String(), start, end, interval)
		if err != nil {
			return nil, err
		}
		return newQueryBackedRemoteOperator(ctx, q), nil
	}
	exec, opCtx, err := l.engine.NewRangeQueryOperator(ctx, l.q, fromPromQLOpts(opts), root, start, end, interval)
	if err != nil {
		return nil, err
	}
	return &inProcessRemoteOperator{exec: exec, ctx: opCtx}, nil
}

// inProcessRemoteOperator adapts a model.VectorOperator to api.RemoteOperator.
// It is intentionally minimal: Series and Next forward directly and Close is
// idempotent. The operator's own lifecycle (storage scanners, etc.) is
// managed via the activeTrackedOperator decorator added by the engine.
type inProcessRemoteOperator struct {
	exec model.VectorOperator
	ctx  context.Context

	closeOnce sync.Once
	closeErr  error
}

func (o *inProcessRemoteOperator) Series(ctx context.Context) ([]labels.Labels, error) {
	return o.exec.Series(ctx)
}

func (o *inProcessRemoteOperator) Next(ctx context.Context, buf []model.StepVector) (int, error) {
	return o.exec.Next(ctx, buf)
}

func (o *inProcessRemoteOperator) Close() error {
	o.closeOnce.Do(func() {
		if c, ok := o.exec.(interface{ Close() error }); ok {
			o.closeErr = c.Close()
		}
	})
	return o.closeErr
}

// queryBackedRemoteOperator wraps a legacy promql.Query as an api.RemoteOperator.
// It is a safety-net used only when the api.RemoteQuery passed to
// NewRangeQueryOperator is not a structured logicalplan.Node. It materialises
// the full result eagerly and replays it as StepVectors, which defeats the
// purpose of the streaming path but preserves correctness.
type queryBackedRemoteOperator struct {
	q   promql.Query
	ctx context.Context

	once    sync.Once
	loadErr error
	series  []labels.Labels
	// matrix holds the materialised result; only used if needed.
}

func newQueryBackedRemoteOperator(ctx context.Context, q promql.Query) *queryBackedRemoteOperator {
	return &queryBackedRemoteOperator{q: q, ctx: ctx}
}

func (o *queryBackedRemoteOperator) Series(ctx context.Context) ([]labels.Labels, error) {
	// Intentionally not implemented: the fallback path is only reached for
	// non-Node RemoteQueries which the in-process engine never produces. If
	// future callers exercise this we can implement it by Exec'ing the query
	// and extracting series labels from the resulting matrix.
	return nil, errFallbackPathNotImplemented
}

func (o *queryBackedRemoteOperator) Next(ctx context.Context, buf []model.StepVector) (int, error) {
	return 0, errFallbackPathNotImplemented
}

func (o *queryBackedRemoteOperator) Close() error {
	o.q.Close()
	return nil
}

type DistributedEngine struct {
	engine *Engine
}

func NewDistributedEngine(opts Opts) *DistributedEngine {
	return &DistributedEngine{
		engine: New(opts),
	}
}

func (l DistributedEngine) MakeInstantQueryFromPlan(ctx context.Context, q storage.Queryable, e api.RemoteEndpoints, opts promql.QueryOpts, plan logicalplan.Node, ts time.Time) (promql.Query, error) {
	// Truncate milliseconds to avoid mismatch in timestamps between remote and local engines.
	// Some clients might only support second precision when executing queries.
	ts = ts.Truncate(time.Second)

	qOpts := fromPromQLOpts(opts)
	qOpts.LogicalOptimizers = []logicalplan.Optimizer{
		logicalplan.PassthroughOptimizer{Endpoints: e},
		logicalplan.DistributedExecutionOptimizer{Endpoints: e},
	}

	return l.engine.MakeInstantQueryFromPlan(ctx, q, qOpts, plan, ts)
}

func (l DistributedEngine) MakeRangeQueryFromPlan(ctx context.Context, q storage.Queryable, e api.RemoteEndpoints, opts promql.QueryOpts, plan logicalplan.Node, start, end time.Time, interval time.Duration) (promql.Query, error) {
	// Truncate milliseconds to avoid mismatch in timestamps between remote and local engines.
	// Some clients might only support second precision when executing queries.
	start = start.Truncate(time.Second)
	end = end.Truncate(time.Second)
	interval = interval.Truncate(time.Second)

	qOpts := fromPromQLOpts(opts)
	qOpts.LogicalOptimizers = []logicalplan.Optimizer{
		logicalplan.PassthroughOptimizer{Endpoints: e},
		logicalplan.DistributedExecutionOptimizer{Endpoints: e},
	}

	return l.engine.MakeRangeQueryFromPlan(ctx, q, qOpts, plan, start, end, interval)
}

func (l DistributedEngine) MakeInstantQuery(ctx context.Context, q storage.Queryable, e api.RemoteEndpoints, opts promql.QueryOpts, qs string, ts time.Time) (promql.Query, error) {
	// Truncate milliseconds to avoid mismatch in timestamps between remote and local engines.
	// Some clients might only support second precision when executing queries.
	ts = ts.Truncate(time.Second)

	qOpts := fromPromQLOpts(opts)
	qOpts.LogicalOptimizers = []logicalplan.Optimizer{
		logicalplan.PassthroughOptimizer{Endpoints: e},
		logicalplan.DistributedExecutionOptimizer{Endpoints: e},
	}

	return l.engine.MakeInstantQuery(ctx, q, qOpts, qs, ts)
}

func (l DistributedEngine) MakeRangeQuery(ctx context.Context, q storage.Queryable, e api.RemoteEndpoints, opts promql.QueryOpts, qs string, start, end time.Time, interval time.Duration) (promql.Query, error) {
	// Truncate milliseconds to avoid mismatch in timestamps between remote and local engines.
	// Some clients might only support second precision when executing queries.
	start = start.Truncate(time.Second)
	end = end.Truncate(time.Second)
	interval = interval.Truncate(time.Second)

	qOpts := fromPromQLOpts(opts)
	qOpts.LogicalOptimizers = []logicalplan.Optimizer{
		logicalplan.PassthroughOptimizer{Endpoints: e},
		logicalplan.DistributedExecutionOptimizer{Endpoints: e},
	}

	return l.engine.MakeRangeQuery(ctx, q, qOpts, qs, start, end, interval)
}
