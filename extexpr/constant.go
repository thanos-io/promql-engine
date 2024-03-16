// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package extexpr

import "github.com/prometheus/prometheus/promql/parser"

// IsConstantExpr reports if the expression evaluates to a constant.
func IsConstantExpr(expr parser.Expr) bool {
	// TODO: there are more possibilities for constant expressions
	switch texpr := expr.(type) {
	case *parser.NumberLiteral:
		return true
	case *parser.StepInvariantExpr:
		return IsConstantExpr(texpr.Expr)
	case *parser.ParenExpr:
		return IsConstantExpr(texpr.Expr)
	case *parser.Call:
		constArgs := true
		for _, arg := range texpr.Args {
			constArgs = constArgs && IsConstantExpr(arg)
		}
		return constArgs
	case *parser.BinaryExpr:
		return IsConstantExpr(texpr.LHS) && IsConstantExpr(texpr.RHS)
	default:
		return false
	}
}
