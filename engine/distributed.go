// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine

import (
	"context"
	"time"

	"github.com/thanos-io/promql-engine/api"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"
)

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

func (l remoteEngine) NewRangeQuery(ctx context.Context, opts promql.QueryOpts, plan api.RemoteQuery, start, end time.Time, interval time.Duration) (promql.Query, error) {
	return l.engine.NewRangeQuery(ctx, l.q, opts, plan.String(), start, end, interval)
}
