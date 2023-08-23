// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"fmt"

	v1 "github.com/prometheus/prometheus/web/api/v1"
	"github.com/thanos-io/promql-engine/parser"
	"github.com/thanos-io/promql-engine/query"
)

type CountValuesExecutionRewriter struct {
	Engine v1.QueryEngine
}

// CountValuesExecutionRewriter makes sure that engines only have to deal with count_values at the
// top of their execution tree.
func (m CountValuesExecutionRewriter) Optimize(plan parser.Expr, _ *query.Options) parser.Expr {
	traverse(&plan, func(e *parser.Expr) {
		if e == nil {
			return
		}
		te, ok := (*e).(*parser.AggregateExpr)
		if !ok {
			return
		}
		if te.Op != parser.COUNT_VALUES {
			return
		}
		*e = CountValues{
			LocalExecution: LocalExecution{
				Engine: m.Engine,
				Query:  te.Expr.String(),
			},
			Grouping: te.Grouping,
			Param:    te.Param.(*parser.StringLiteral).Val,
		}
	})
	return plan
}

type CountValues struct {
	LocalExecution

	Param    string
	Grouping []string
}

func (r CountValues) String() string {
	return fmt.Sprintf("local(count_values(%s))", r.Query)
}

func (r CountValues) Pretty(level int) string { return r.String() }

func (r CountValues) PositionRange() parser.PositionRange { return parser.PositionRange{} }

func (r CountValues) Type() parser.ValueType { return parser.ValueTypeVector }

func (r CountValues) PromQLExpr() {}

type LocalExecution struct {
	Engine v1.QueryEngine
	Query  string
}

func (r LocalExecution) String() string {
	return fmt.Sprintf("local(%s)", r.Query)
}

func (r LocalExecution) Pretty(level int) string { return r.String() }

func (r LocalExecution) PositionRange() parser.PositionRange { return parser.PositionRange{} }

func (r LocalExecution) Type() parser.ValueType { return parser.ValueTypeMatrix }

func (r LocalExecution) PromQLExpr() {}
