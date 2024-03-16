// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package extexpr

import (
	"github.com/efficientgo/core/errors"
	"github.com/prometheus/prometheus/promql/parser"
)

// UnwrapString recursively unwraps an parser.Expr until it reaches an StringLiteral.
func UnwrapString(expr parser.Expr) (string, error) {
	switch texpr := expr.(type) {
	case *parser.StringLiteral:
		return texpr.Val, nil
	case *parser.ParenExpr:
		return UnwrapString(texpr.Expr)
	case *parser.StepInvariantExpr:
		return UnwrapString(texpr.Expr)
	default:
		return "", errors.Newf("unexpected type: %T", texpr)
	}
}

// UnsafeUnwrapString is like UnwrapString but should only be used in cases where the parser
// guarantees success by already only allowing strings wrapped in parentheses.
func UnsafeUnwrapString(expr parser.Expr) string {
	v, _ := UnwrapString(expr)
	return v
}

// UnwrapString recursively unwraps an parser.Expr until it reaches an NumberLiteral.
func UnwrapFloat(expr parser.Expr) (float64, error) {
	switch texpr := expr.(type) {
	case *parser.NumberLiteral:
		return texpr.Val, nil
	case *parser.ParenExpr:
		return UnwrapFloat(texpr.Expr)
	case *parser.StepInvariantExpr:
		return UnwrapFloat(texpr.Expr)
	default:
		return 0, errors.Newf("unexpected type: %T", texpr)
	}
}

// UnwrapString recursively unwraps an parser.ParenExpr.
func UnwrapParens(expr parser.Expr) parser.Expr {
	switch t := expr.(type) {
	case *parser.ParenExpr:
		return UnwrapParens(t.Expr)
	default:
		return t
	}
}
