// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"maps"
	"slices"
	"strings"

	"github.com/efficientgo/core/errors"
	"github.com/prometheus/prometheus/promql"
)

// CheckRangeModifiers returns Prometheus' evaluation error when the anchored or
// smoothed modifier is used with a function that does not support it.
func CheckRangeModifiers(funcName string, vs *VectorSelector) error {
	if vs.Anchored {
		if _, ok := promql.AnchoredSafeFunctions[funcName]; !ok {
			return errors.Newf("anchored modifier can only be used with: %s - not with %s", strings.Join(slices.Sorted(maps.Keys(promql.AnchoredSafeFunctions)), ", "), funcName)
		}
	}
	if vs.Smoothed {
		if _, ok := promql.SmoothedSafeFunctions[funcName]; !ok {
			return errors.Newf("smoothed modifier can only be used with: %s - not with %s", strings.Join(slices.Sorted(maps.Keys(promql.SmoothedSafeFunctions)), ", "), funcName)
		}
	}
	return nil
}

// ValidateRangeModifiers returns the error Prometheus raises when evaluating
// root, if a function in it is called with a range selector whose modifier the
// function does not support. When there are several, it returns the one
// Prometheus reaches first.
func ValidateRangeModifiers(root Node) error {
	switch e := root.(type) {
	case *FunctionCall:
		return validateCallRangeModifiers(e)
	case *Aggregation:
		// The parameter is evaluated before the aggregated expression.
		if e.Param != nil {
			if err := ValidateRangeModifiers(e.Param); err != nil {
				return err
			}
		}
		return ValidateRangeModifiers(e.Expr)
	}
	for _, c := range root.Children() {
		if err := ValidateRangeModifiers(*c); err != nil {
			return err
		}
	}
	return nil
}

// validateCallRangeModifiers follows the order of Prometheus' evaluation of a
// function call: the first range argument when it is a subquery, then the
// other arguments, then the modifier of a range selector argument.
func validateCallRangeModifiers(call *FunctionCall) error {
	rangeArg := -1
	for i, arg := range call.Args {
		switch arg := unwrapParens(arg).(type) {
		case *MatrixSelector:
			rangeArg = i
		case *Subquery:
			rangeArg = i
			if err := ValidateRangeModifiers(arg); err != nil {
				return err
			}
		}
		if rangeArg >= 0 {
			break
		}
	}
	for i, arg := range call.Args {
		if i == rangeArg {
			continue
		}
		if err := ValidateRangeModifiers(arg); err != nil {
			return err
		}
	}
	if rangeArg >= 0 {
		if ms, ok := unwrapParens(call.Args[rangeArg]).(*MatrixSelector); ok {
			return CheckRangeModifiers(call.Func.Name, ms.VectorSelector)
		}
	}
	return nil
}

// unwrapParens removes the parentheses around expr, as Prometheus does for
// function arguments before evaluation.
func unwrapParens(expr Node) Node {
	for {
		p, ok := expr.(*Parens)
		if !ok {
			return expr
		}
		expr = p.Expr
	}
}
