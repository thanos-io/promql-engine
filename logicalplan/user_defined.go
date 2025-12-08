// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"context"

	"github.com/thanos-io/promql-engine/execution/execopts"
	"github.com/thanos-io/promql-engine/execution/model"

	"github.com/prometheus/prometheus/storage"
)

// UserDefinedExpr is an extension point which allows users to define their execution operators.
type UserDefinedExpr interface {
	Node
	MakeExecutionOperator(
		ctx context.Context,
		vectors *model.VectorPool,
		opts *execopts.Options,
		hints storage.SelectHints,
	) (model.VectorOperator, error)
}
