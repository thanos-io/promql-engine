// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package api

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"

	"github.com/thanos-io/promql-engine/execution/model"
)

type RemoteQuery interface {
	fmt.Stringer
}

type RemoteEndpoints interface {
	Engines() []RemoteEngine
}

type RemoteEngine interface {
	MaxT() int64
	MinT() int64

	// The external labels of the remote engine. These are used to limit fanout. The engine uses these to
	// not distribute into remote engines that would return empty responses because their labelset is not matching.
	LabelSets() []labels.Labels

	// The external labels of the remote engine that form a logical partition. This is expected to be
	// a subset of the result of "LabelSets()". The engine uses these to compute how to distribute a query.
	// It is important that, for a given set of remote engines, these labels do not overlap meaningfully.
	PartitionLabelSets() []labels.Labels

	NewRangeQuery(ctx context.Context, opts promql.QueryOpts, plan RemoteQuery, start, end time.Time, interval time.Duration) (promql.Query, error)
}

// RemoteOperator is the wire-facing shape of a remote query result that is
// produced and consumed as a stream of columnar StepVectors, rather than
// being materialized as a full promql.Matrix first.
//
// It mirrors a subset of model.VectorOperator: callers must call Series first
// (or implicitly, before/inside Next) to learn the label set that the SampleIDs
// in subsequent StepVectors index into, then drive Next until it returns 0.
// Close must be called exactly once, after Next has returned 0 or after the
// caller decides to stop consuming early.
type RemoteOperator interface {
	io.Closer

	// Series returns the labels of all series that will appear in the StepVectors
	// produced by Next. The order is stable; SampleIDs in StepVectors are positional
	// indices into this slice.
	Series(ctx context.Context) ([]labels.Labels, error)

	// Next fills buf with as many consecutive StepVectors as it can produce, and
	// returns the count written. A return value of 0 indicates end of stream.
	Next(ctx context.Context, buf []model.StepVector) (int, error)
}

// StreamingRemoteEngine is an optional capability that a RemoteEngine may
// implement to produce a RemoteOperator directly. Compared to the legacy
// NewRangeQuery path which returns a promql.Query and forces a StepVector ->
// Matrix -> StepVector round-trip, the streaming path lets the local query
// consume the remote result columnar end-to-end.
//
// Instant queries are modelled as a single-step range query (start == end,
// interval == 0 or 1ms). Engines should accept that shape and return a
// RemoteOperator producing exactly one StepVector.
type StreamingRemoteEngine interface {
	RemoteEngine

	NewRangeQueryOperator(ctx context.Context, opts promql.QueryOpts, plan RemoteQuery, start, end time.Time, interval time.Duration) (RemoteOperator, error)
}

type staticEndpoints struct {
	engines []RemoteEngine
}

func (m staticEndpoints) Engines() []RemoteEngine {
	return m.engines
}

func NewStaticEndpoints(engines []RemoteEngine) RemoteEndpoints {
	return &staticEndpoints{engines: engines}
}
