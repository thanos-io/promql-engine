// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package model

import (
	"context"

	"github.com/prometheus/prometheus/model/labels"
)

// identifiedOperator wraps a VectorOperator with a deterministic fingerprint.
type identifiedOperator struct {
	VectorOperator
	id          uint64
	enrichedCtx context.Context
}

// WithID wraps op so that it implements OperatorIDer. The provided ctx is the
// query execution context; ContextWithOperatorID is called exactly once here and
// the result is reused on every Series and Next call.
func WithID(ctx context.Context, op VectorOperator, id uint64) VectorOperator {
	return &identifiedOperator{
		VectorOperator: op,
		id:             id,
		enrichedCtx:    ContextWithOperatorID(ctx, id),
	}
}

func (o *identifiedOperator) OperatorID() uint64 { return o.id }

func (o *identifiedOperator) Series(_ context.Context) ([]labels.Labels, error) {
	return o.VectorOperator.Series(o.enrichedCtx)
}

func (o *identifiedOperator) Next(_ context.Context, buf []StepVector) (int, error) {
	return o.VectorOperator.Next(o.enrichedCtx, buf)
}

func (o *identifiedOperator) Unwrap() VectorOperator { return o.VectorOperator }

type Unwrapper interface {
	Unwrap() VectorOperator
}

func Unwrap(op VectorOperator) VectorOperator {
	for {
		u, ok := op.(Unwrapper)
		if !ok {
			return op
		}
		op = u.Unwrap()
	}
}
