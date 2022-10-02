// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package parse

import (
	"fmt"

	"github.com/efficientgo/core/errors"
	"github.com/prometheus/prometheus/promql/parser"
)

var ErrNotSupportedExpr = errors.New("unsupported expression")

func NotSupportedOperationErr(op parser.ItemType) error {
	t := parser.ItemTypeStr[op]
	msg := fmt.Sprintf("operation not supported: %s", t)
	return errors.Wrap(ErrNotSupportedExpr, msg)
}
