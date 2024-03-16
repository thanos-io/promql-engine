// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package extscanners

import (
	"fmt"

	"github.com/prometheus/prometheus/promql/parser"
)

type UnsupportedExprError struct {
	Type parser.ValueType
}

func newUnsupportedExprError(t parser.ValueType) *UnsupportedExprError {
	return &UnsupportedExprError{Type: t}
}

func (e *UnsupportedExprError) Error() string {
	return fmt.Sprintf("unsupported expression type %s", e.Type)
}

func UnwrapConstVal(e parser.Expr) (float64, error) {
	switch c := e.(type) {
	case *parser.NumberLiteral:
		return c.Val, nil
	case *parser.StepInvariantExpr:
		return UnwrapConstVal(c.Expr)
	}

	return 0, newUnsupportedExprError(e.Type())
}
