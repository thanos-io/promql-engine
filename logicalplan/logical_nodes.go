// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/promql/parser/posrange"
)

// Projection has information on which series labels should be selected from storage.
type Projection struct {
	// Labels is a list of labels to be included or excluded from the selection result, depending on the value of Include.
	Labels []string
	// Include is true if only the provided list of labels should be retrieved from storage.
	// When set to false, the provided list of labels should be excluded from selection.
	Include bool
}

// VectorSelector is vector selector with additional configuration set by optimizers.
type VectorSelector struct {
	*parser.VectorSelector
	Filters         []*labels.Matcher
	BatchSize       int64
	SelectTimestamp bool
	Projection      Projection
}

func (f VectorSelector) String() string {
	if f.SelectTimestamp {
		// If we pushed down timestamp into the vector selector we need to render the proper
		// PromQL again.
		return fmt.Sprintf("timestamp(%s)", f.VectorSelector.String())
	}
	return f.VectorSelector.String()
}

func (f VectorSelector) Pretty(level int) string { return f.String() }

func (f VectorSelector) PositionRange() posrange.PositionRange { return posrange.PositionRange{} }

func (f VectorSelector) Type() parser.ValueType { return parser.ValueTypeVector }

func (f VectorSelector) PromQLExpr() {}

// MatrixSelector is matrix selector with additional configuration set by optimizers.
// It is used so we can get rid of VectorSelector in distributed mode too.
type MatrixSelector struct {
	VectorSelector parser.Expr
	Range          time.Duration

	// Needed because this operator is used in the distributed mode
	OriginalString string
}

func (f MatrixSelector) String() string {
	return f.OriginalString
}

func (f MatrixSelector) Pretty(level int) string { return f.String() }

func (f MatrixSelector) PositionRange() posrange.PositionRange { return posrange.PositionRange{} }

func (f MatrixSelector) Type() parser.ValueType { return parser.ValueTypeVector }

func (f MatrixSelector) PromQLExpr() {}

// CheckDuplicateLabels is a logical node that checks for duplicate labels in the same timestamp.
type CheckDuplicateLabels struct {
	Expr parser.Expr
}

func (c CheckDuplicateLabels) String() string {
	return c.Expr.String()
}

func (c CheckDuplicateLabels) Pretty(level int) string { return c.Expr.Pretty(level) }

func (c CheckDuplicateLabels) PositionRange() posrange.PositionRange { return c.Expr.PositionRange() }

func (c CheckDuplicateLabels) Type() parser.ValueType { return c.Expr.Type() }

func (c CheckDuplicateLabels) PromQLExpr() {}

// StringLiteral is a logical node representing a literal string.
type StringLiteral struct {
	Val string
}

func (c StringLiteral) String() string {
	return fmt.Sprintf("%q", c.Val)
}

func (c StringLiteral) Pretty(level int) string { return c.String() }

func (c StringLiteral) PositionRange() posrange.PositionRange { return posrange.PositionRange{} }

func (c StringLiteral) Type() parser.ValueType { return parser.ValueTypeString }

func (c StringLiteral) PromQLExpr() {}

// NumberLiteral is a logical node representing a literal number.
type NumberLiteral struct {
	Val float64
}

func (c NumberLiteral) String() string {
	return fmt.Sprint(c.Val)
}

func (c NumberLiteral) Pretty(level int) string { return c.String() }

func (c NumberLiteral) PositionRange() posrange.PositionRange { return posrange.PositionRange{} }

func (c NumberLiteral) Type() parser.ValueType { return parser.ValueTypeScalar }

func (c NumberLiteral) PromQLExpr() {}

// StepInvariantExpr is a logical node that expresses that the child expression
// returns the same value at every step in the evaluation.
type StepInvariantExpr struct {
	Expr parser.Expr
}

func (c StepInvariantExpr) String() string { return c.Expr.String() }

func (c StepInvariantExpr) Pretty(level int) string { return c.String() }

func (c StepInvariantExpr) PositionRange() posrange.PositionRange {
	return c.Expr.PositionRange()
}

func (c StepInvariantExpr) Type() parser.ValueType { return c.Expr.Type() }

func (c StepInvariantExpr) PromQLExpr() {}

// FunctionCall represents a PromQL function.
type FunctionCall struct {
	// The function that was called.
	Func *parser.Function
	// Arguments passed into the function.
	Args parser.Expressions
}

func (f FunctionCall) String() string {
	args := make([]string, 0, len(f.Args))
	for _, arg := range f.Args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("%s(%s)", f.Func.Name, strings.Join(args, ", "))
}

func (f FunctionCall) Pretty(level int) string { return f.String() }

func (f FunctionCall) PositionRange() posrange.PositionRange { return posrange.PositionRange{} }

func (f FunctionCall) Type() parser.ValueType { return f.Func.ReturnType }

func (f FunctionCall) PromQLExpr() {}

type Parens struct {
	Expr parser.Expr
}

func (p Parens) String() string {
	return fmt.Sprintf("(%s)", p.Expr.String())
}

func (p Parens) Pretty(level int) string { return p.String() }

func (p Parens) PositionRange() posrange.PositionRange { return p.Expr.PositionRange() }

func (p Parens) Type() parser.ValueType { return p.Expr.Type() }

func (p Parens) PromQLExpr() {}

type Unary struct {
	Op   parser.ItemType
	Expr parser.Expr
}

func (p Unary) String() string {
	return fmt.Sprintf("%s%s", p.Op.String(), p.Expr.String())
}

func (p Unary) Pretty(level int) string { return p.String() }

func (p Unary) PositionRange() posrange.PositionRange { return p.Expr.PositionRange() }

func (p Unary) Type() parser.ValueType { return p.Expr.Type() }

func (p Unary) PromQLExpr() {}

// Aggregation represents a PromQL aggregation.
type Aggregation struct {
	Op       parser.ItemType // The used aggregation operation.
	Expr     parser.Expr     // The Vector expression over which is aggregated.
	Param    parser.Expr     // Parameter used by some aggregators.
	Grouping []string        // The labels by which to group the Vector.
	Without  bool            // Whether to drop the given labels rather than keep them
}

func (f Aggregation) String() string {
	aggrString := f.getAggOpStr()
	aggrString += "("
	if f.Op.IsAggregatorWithParam() {
		aggrString += fmt.Sprintf("%s, ", f.Param)
	}
	aggrString += fmt.Sprintf("%s)", f.Expr)

	return aggrString
}

func (f Aggregation) Pretty(_ int) string { return f.String() }

func (f Aggregation) PositionRange() posrange.PositionRange { return posrange.PositionRange{} }

func (f Aggregation) Type() parser.ValueType { return parser.ValueTypeVector }

func (f Aggregation) PromQLExpr() {}

func (f Aggregation) getAggOpStr() string {
	aggrString := f.Op.String()

	switch {
	case f.Without:
		aggrString += fmt.Sprintf(" without (%s) ", strings.Join(f.Grouping, ", "))
	case len(f.Grouping) > 0:
		aggrString += fmt.Sprintf(" by (%s) ", strings.Join(f.Grouping, ", "))
	}

	return aggrString
}

type Binary struct {
	Op       parser.ItemType // The operation of the expression.
	LHS, RHS parser.Expr     // The operands on the respective sides of the operator.

	// The matching behavior for the operation if both operands are Vectors.
	// If they are not this field is nil.
	VectorMatching *parser.VectorMatching

	// If a comparison operator, return 0/1 rather than filtering.
	ReturnBool bool

	ValueType parser.ValueType
}

func (b Binary) Pretty(_ int) string { return b.String() }

func (b Binary) PositionRange() posrange.PositionRange { return posrange.PositionRange{} }

func (b Binary) Type() parser.ValueType {
	if b.LHS.Type() == parser.ValueTypeScalar && b.RHS.Type() == parser.ValueTypeScalar {
		return parser.ValueTypeScalar
	}
	return parser.ValueTypeVector
}

func (b Binary) PromQLExpr() {}

func (b Binary) String() string {
	returnBool := ""
	if b.ReturnBool {
		returnBool = " bool"
	}

	matching := b.getMatchingStr()
	return fmt.Sprintf("%s %s%s%s %s", b.LHS, b.Op, returnBool, matching, b.RHS)
}

func (b Binary) getMatchingStr() string {
	matching := ""
	vm := b.VectorMatching
	if vm != nil && (len(vm.MatchingLabels) > 0 || vm.On) {
		vmTag := "ignoring"
		if vm.On {
			vmTag = "on"
		}
		matching = fmt.Sprintf(" %s (%s)", vmTag, strings.Join(vm.MatchingLabels, ", "))

		if vm.Card == parser.CardManyToOne || vm.Card == parser.CardOneToMany {
			vmCard := "right"
			if vm.Card == parser.CardManyToOne {
				vmCard = "left"
			}
			matching += fmt.Sprintf(" group_%s (%s)", vmCard, strings.Join(vm.Include, ", "))
		}
	}
	return matching
}
