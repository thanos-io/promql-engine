// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package noop

import (
	"context"

	"github.com/prometheus/prometheus/model/labels"

	"github.com/thanos-community/promql-engine/execution/model"
)

type operator struct{}

func NewOperator() model.VectorOperator { return &operator{} }

func (o operator) Next(_ context.Context, _ *model.OperatorTracer) ([]model.StepVector, error) {
	return nil, nil
}

func (o operator) Series(_ context.Context, _ *model.OperatorTracer) ([]labels.Labels, error) {
	return nil, nil
}

func (o operator) GetPool() *model.VectorPool { return nil }

func (o operator) Explain() (me string, next []model.VectorOperator) { return "noop", nil }
