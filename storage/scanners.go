// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package storage

import (
	"context"

	"github.com/thanos-io/promql-engine/execution/execopts"
	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/logicalplan"

	"github.com/prometheus/prometheus/storage"
)

type Scanners interface {
	Close() error
	NewVectorSelector(ctx context.Context, opts *execopts.Options, hints storage.SelectHints, selector logicalplan.VectorSelector) (model.VectorOperator, error)
	NewMatrixSelector(ctx context.Context, opts *execopts.Options, hints storage.SelectHints, selector logicalplan.MatrixSelector, call logicalplan.FunctionCall) (model.VectorOperator, error)
}
