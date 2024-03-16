package extscanners

import (
	"github.com/efficientgo/core/errors"
	"github.com/prometheus/prometheus/promql/parser"

	"github.com/thanos-io/promql-engine/execution/parse"
)

func UnwrapConstVal(e parser.Expr) (float64, error) {
	switch c := e.(type) {
	case *parser.NumberLiteral:
		return c.Val, nil
	case *parser.StepInvariantExpr:
		return UnwrapConstVal(c.Expr)
	}

	return 0, errors.Wrap(parse.ErrNotSupportedExpr, "argument must be a constant value")
}
