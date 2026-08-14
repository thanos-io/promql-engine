// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package model

import (
	"context"
	"fmt"

	"github.com/prometheus/prometheus/model/labels"
)

// OperatorIDer is an optional interface for operators that carry a
// deterministic fingerprint derived from their logical plan subtree.
// Only data-fetching operators implement this; pure-computation operators
// do not.
type OperatorIDer interface {
	OperatorID() uint64
}

type operatorIDKey struct{}

// ContextWithOperatorID returns a copy of ctx carrying the given operator ID.
func ContextWithOperatorID(ctx context.Context, id uint64) context.Context {
	return context.WithValue(ctx, operatorIDKey{}, id)
}

// OperatorIDFromContext returns the operator ID stored in ctx, if any.
func OperatorIDFromContext(ctx context.Context) (uint64, bool) {
	id, ok := ctx.Value(operatorIDKey{}).(uint64)
	return id, ok
}

// VectorOperator performs operations on series in step by step fashion.
type VectorOperator interface {
	// Next yields vectors of samples from all series for one or more execution steps.
	// The caller provides a buffer (buf) to be filled with StepVectors.
	// Returns the number of StepVectors written to buf and any error encountered.
	// A return value of 0 indicates no more data is available.
	Next(ctx context.Context, buf []StepVector) (int, error)

	// Series returns all series that the operator will process during Next results.
	// The result can be used by upstream operators to allocate output tables and buffers
	// before starting to process samples.
	Series(ctx context.Context) ([]labels.Labels, error)

	// Explain returns human-readable explanation of the current operator and optional nested operators.
	Explain() (next []VectorOperator)

	fmt.Stringer
}

// OriginHashesProvider is an optional interface for operators whose output
// series map one-to-one onto storage series with a known original label set.
type OriginHashesProvider interface {
	// OriginHashes returns one hash per Series() entry, holding the hash of
	// the series' original label set before a projection trimmed it.
	// A nil slice or a zero hash means the origin of the series is unknown.
	OriginHashes(ctx context.Context) ([]uint64, error)
}

// OriginHashes returns the origin hashes of the operator's output series,
// or nil when the operator cannot provide them.
func OriginHashes(ctx context.Context, op VectorOperator) ([]uint64, error) {
	for {
		if p, ok := op.(OriginHashesProvider); ok {
			return p.OriginHashes(ctx)
		}
		u, ok := op.(Unwrapper)
		if !ok {
			return nil, nil
		}
		op = u.Unwrap()
	}
}
