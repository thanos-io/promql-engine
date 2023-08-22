package logicalplan

import (
	"github.com/prometheus/prometheus/model/labels"
	"github.com/thanos-io/promql-engine/parser"
	"time"
)

type LogicalPlan interface {
	Children() LogicalPlans
}

type LogicalPlans []LogicalPlan

// StepInvariantExpr represents a query which evaluates to the same result
// irrespective of the evaluation time given the raw samples from TSDB remain unchanged.
// Currently this is only used for engine optimisations and the parser does not produce this.
type StepInvariantExpr struct {
	Expr LogicalPlan
}

func (l *StepInvariantExpr) Children() LogicalPlans {
	return []LogicalPlan{l.Expr}
}

// VectorSelector represents a Vector selection.
type VectorSelector struct {
	Name string
	// OriginalOffset is the actual offset that was set in the query.
	// This never changes.
	OriginalOffset time.Duration
	// Offset is the offset used during the query execution
	// which is calculated using the original offset, at modifier time,
	// eval time, and subquery offsets in the AST tree.
	Offset        time.Duration
	Timestamp     *int64
	StartOrEnd    parser.ItemType // Set when @ is used with start() or end()
	LabelMatchers []*labels.Matcher
}

func (l *VectorSelector) Children() LogicalPlans {
	return []LogicalPlan{}
}

// MatrixSelector represents a Matrix selection.
type MatrixSelector struct {
	// It is safe to assume that this is an VectorSelector
	// if the parser hasn't returned an error.
	VectorSelector LogicalPlan
	Range          time.Duration
}

func (l *MatrixSelector) Children() LogicalPlans {
	return []LogicalPlan{l.VectorSelector}
}

// AggregateExpr represents an aggregation operation on a Vector.
type AggregateExpr struct {
	Op       parser.ItemType // The used aggregation operation.
	Expr     LogicalPlan     // The Vector expression over which is aggregated.
	Param    LogicalPlan     // Parameter used by some aggregators.
	Grouping []string        // The labels by which to group the Vector.
	Without  bool            // Whether to drop the given labels rather than keep them.
}

func (l *AggregateExpr) Children() LogicalPlans {
	return []LogicalPlan{l.Expr, l.Param}
}

// Call represents a function call.
type Call struct {
	Func *parser.Function // The function that was called.
	Args LogicalPlans     // Arguments used in the call.
}

func (l *Call) Children() LogicalPlans {
	return l.Args
}

// BinaryExpr represents a binary expression between two child expressions.
type BinaryExpr struct {
	Op       parser.ItemType // The operation of the expression.
	LHS, RHS LogicalPlan     // The operands on the respective sides of the operator.

	// The matching behavior for the operation if both operands are Vectors.
	// If they are not this field is nil.
	VectorMatching *parser.VectorMatching

	// If a comparison operator, return 0/1 rather than filtering.
	ReturnBool bool
}

func (l *BinaryExpr) Children() LogicalPlans {
	return []LogicalPlan{l.LHS, l.RHS}
}

// UnaryExpr represents a unary operation on another expression.
// Currently unary operations are only supported for Scalars.
type UnaryExpr struct {
	Op   parser.ItemType
	Expr LogicalPlan
}

func (l *UnaryExpr) Children() LogicalPlans {
	return []LogicalPlan{l.Expr}
}

// ParenExpr wraps an expression so it cannot be disassembled as a consequence
// of operator precedence.
type ParenExpr struct {
	Expr LogicalPlan
}

func (l *ParenExpr) Children() LogicalPlans {
	return []LogicalPlan{l.Expr}
}

// SubqueryExpr represents a subquery.
type SubqueryExpr struct {
	Expr  LogicalPlan
	Range time.Duration
	// OriginalOffset is the actual offset that was set in the query.
	// This never changes.
	OriginalOffset time.Duration
	// Offset is the offset used during the query execution
	// which is calculated using the original offset, at modifier time,
	// eval time, and subquery offsets in the AST tree.
	Offset     time.Duration
	Timestamp  *int64
	StartOrEnd parser.ItemType // Set when @ is used with start() or end()
	Step       time.Duration
}

func (l *SubqueryExpr) Children() LogicalPlans {
	return []LogicalPlan{l.Expr}
}

// literal types

// NumberLiteral represents a number.
type NumberLiteral struct {
	Val float64
}

func (l *NumberLiteral) Children() LogicalPlans {
	return []LogicalPlan{}
}

// StringLiteral represents a string.
type StringLiteral struct {
	Val string
}

func (l *StringLiteral) Children() LogicalPlans {
	return []LogicalPlan{}
}
