// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package parse

import (
	"fmt"

	"github.com/efficientgo/core/errors"
	"github.com/prometheus/prometheus/promql/parser"
)

var XFunctions = map[string]*parser.Function{
	"xdelta": {
		Name:       "xdelta",
		ArgTypes:   []parser.ValueType{parser.ValueTypeMatrix},
		ReturnType: parser.ValueTypeVector,
	},
	"xincrease": {
		Name:       "xincrease",
		ArgTypes:   []parser.ValueType{parser.ValueTypeMatrix},
		ReturnType: parser.ValueTypeVector,
	},
	"xrate": {
		Name:       "xrate",
		ArgTypes:   []parser.ValueType{parser.ValueTypeMatrix},
		ReturnType: parser.ValueTypeVector,
	},
}

// IsExtFunction is a convenience function to determine whether extended range calculations are required.
func IsExtFunction(functionName string) bool {
	_, ok := XFunctions[functionName]
	return ok
}

func UnknownFunctionError(name string) error {
	msg := fmt.Sprintf("unknown function: %s", name)
	if _, ok := parser.Functions[name]; ok {
		return errors.Wrap(ErrNotImplemented, msg)
	}

	return errors.Wrap(ErrNotSupportedExpr, msg)
}
