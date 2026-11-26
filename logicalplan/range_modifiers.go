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

// hasInvalidRangeModifiers reports whether any function in the plan is called
// with a range selector whose modifier it does not support.
func hasInvalidRangeModifiers(root Node) bool {
	var invalid bool
	Traverse(&root, func(node *Node) {
		call, ok := (*node).(*FunctionCall)
		if !ok || invalid {
			return
		}
		for _, arg := range call.Args {
			if ms, ok := arg.(*MatrixSelector); ok && CheckRangeModifiers(call.Func.Name, ms.VectorSelector) != nil {
				invalid = true
			}
		}
	})
	return invalid
}
