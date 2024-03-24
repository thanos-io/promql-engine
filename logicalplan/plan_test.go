// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"

	"github.com/prometheus/prometheus/promql/parser"

	"github.com/thanos-io/promql-engine/query"
)

var spaces = regexp.MustCompile(`\s+`)
var openParenthesis = regexp.MustCompile(`\(\s+`)
var closedParenthesis = regexp.MustCompile(`\s+\)`)

// renderExprTree renders the expression into a string. It is useful
// in tests to use strings for assertions in cases where the "String()"
// method might not yield enough information or would panic because of
// internal logical expression types. Implementations were largeley taken
// from upstream prometheus.
//
// TODO: maybe its better to Traverse the expression here and inject
// new nodes with prepared String methods? Like replacing MatrixSelector
// by testMatrixSelector that has a overridden string method?
func renderExprTree(expr Node) string {
	switch t := expr.(type) {
	case *NumberLiteral:
		return fmt.Sprint(t.Val)
	case *VectorSelector:
		var b strings.Builder
		base := t.VectorSelector.String()
		if t.BatchSize > 0 {
			base += fmt.Sprintf("[batch=%d]", t.BatchSize)
		}
		if len(t.Filters) > 0 {
			b.WriteString("filter(")
			b.WriteString(fmt.Sprintf("%s", t.Filters))
			b.WriteString(", ")
			b.WriteString(base)
			b.WriteRune(')')
			return b.String()
		}
		return base
	case *MatrixSelector:
		return t.String()
	case *Binary:
		var b strings.Builder
		b.WriteString(renderExprTree(t.LHS))
		b.WriteString(" ")
		b.WriteString(t.Op.String())
		b.WriteString(" ")
		if vm := t.VectorMatching; vm != nil && (len(vm.MatchingLabels) > 0 || vm.On) {
			vmTag := "ignoring"
			if vm.On {
				vmTag = "on"
			}
			matching := fmt.Sprintf("%s (%s)", vmTag, strings.Join(vm.MatchingLabels, ", "))

			if vm.Card == parser.CardManyToOne || vm.Card == parser.CardOneToMany {
				vmCard := "right"
				if vm.Card == parser.CardManyToOne {
					vmCard = "left"
				}
				matching += fmt.Sprintf(" group_%s (%s)", vmCard, strings.Join(vm.Include, ", "))
			}
			b.WriteString(matching)
			b.WriteString(" ")
		}
		b.WriteString(renderExprTree(t.RHS))
		return b.String()
	case *FunctionCall:
		var b strings.Builder
		b.Write([]byte(t.Func.Name))
		b.WriteRune('(')
		for i := range t.Args {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(renderExprTree(t.Args[i]))
		}
		b.WriteRune(')')
		return b.String()
	case *Aggregation:
		var b strings.Builder
		b.Write([]byte(t.Op.String()))
		switch {
		case t.Without:
			b.WriteString(fmt.Sprintf(" without (%s) ", strings.Join(t.Grouping, ", ")))
		case len(t.Grouping) > 0:
			b.WriteString(fmt.Sprintf(" by (%s) ", strings.Join(t.Grouping, ", ")))
		}
		b.WriteRune('(')
		if t.Param != nil {
			b.WriteString(renderExprTree(t.Param))
			b.WriteString(", ")
		}
		b.WriteString(renderExprTree(t.Expr))
		b.WriteRune(')')
		return b.String()
	case *StepInvariantExpr:
		return renderExprTree(t.Expr)
	case *CheckDuplicateLabels:
		return renderExprTree(t.Expr)
	default:
		return t.String()
	}
}

func TestDefaultOptimizers(t *testing.T) {
	cases := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "common selectors",
			expr:     `sum(metric{a="b", c="d"}) / sum(metric{a="b"})`,
			expected: `sum(filter([c="d"], metric{a="b"})) / sum(metric{a="b"})`,
		},
		{
			name:     "common selectors with duplicate matchers",
			expr:     `sum(metric{a="b", c="d", a="b"}) / sum(metric{a="b"})`,
			expected: `sum(filter([c="d"], metric{a="b"})) / sum(metric{a="b"})`,
		},
		{
			name:     "common selectors with different operators",
			expr:     `sum(metric{a="b"}) / sum(metric{a=~"b"})`,
			expected: `sum(metric{a="b"}) / sum(metric{a=~"b"})`,
		},
		{
			name:     "common selectors with regex",
			expr:     `http_requests_total / on () group_left sum(http_requests_total{pod=~"p1.+"})`,
			expected: `http_requests_total / on () group_left () sum(filter([pod=~"p1.+"], http_requests_total))`,
		},
		{
			name: "common selectors in different metrics",
			expr: `
	sum(metric_1{a="b", c="d"}) / sum(metric_1{a="b"}) +
	sum(metric_2{a="b", c="d"}) / sum(metric_2{a="b"})
`,
			expected: `
	sum(filter([c="d"], metric_1{a="b"})) / sum(metric_1{a="b"}) +
	sum(filter([c="d"], metric_2{a="b"})) / sum(metric_2{a="b"})`,
		},
		{
			name:     "different selectors",
			expr:     `sum(metric{a="b"}) / sum(metric{c="d"})`,
			expected: `sum(metric{a="b"}) / sum(metric{c="d"})`,
		},
		{
			name:     "different operator",
			expr:     `sum(metric{a="b"}) / sum(metric{a=~"b"})`,
			expected: `sum(metric{a="b"}) / sum(metric{a=~"b"})`,
		},
		{
			name:     "different metrics",
			expr:     `sum(metric_1{a="b"}) / sum(metric_2{a="b"})`,
			expected: `sum(metric_1{a="b"}) / sum(metric_2{a="b"})`,
		},
		{
			name:     "duplicate matchers",
			expr:     `metric_1{a="1", b="2", a="1"} / metric_2{a="1", b="2", a="1"}`,
			expected: `metric_1{a="1",a="1",b="2"} / metric_2{a="1",a="1",b="2"}`,
		},
		{
			name:     "duplicate matchers",
			expr:     `metric_1{a="1", b="2", a="1", e="f"} / metric_1{a="1", b="2", a="1"}`,
			expected: `filter([e="f"], metric_1{a="1",a="1",b="2"}) / metric_1{a="1",a="1",b="2"}`,
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			plan := New(expr, &query.Options{Start: time.Unix(0, 0), End: time.Unix(0, 0)}, PlanOptions{})
			optimizedPlan, _ := plan.Optimize(DefaultOptimizers)
			expectedPlan := strings.Trim(spaces.ReplaceAllString(tcase.expected, " "), " ")
			testutil.Equals(t, expectedPlan, renderExprTree(optimizedPlan.Root()))
		})
	}
}

func TestMatcherPropagation(t *testing.T) {
	cases := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "common matchers with same metric",
			expr:     `node_filesystem_files{host="$host", mountpoint="/"} - node_filesystem_files`,
			expected: `node_filesystem_files{host="$host",mountpoint="/"} - node_filesystem_files`,
		},
		{
			name:     "common matchers with same overlapping selectors",
			expr:     `node_filesystem_files{host="$host", mountpoint="/"} - node_filesystem_files{host!="$host"}`,
			expected: `node_filesystem_files{host="$host",mountpoint="/"} - node_filesystem_files{host!="$host"}`,
		},
		{
			name:     "common matchers with many-to-one",
			expr:     `node_filesystem_files{host="$host",mountpoint="/"} - on () group_left () node_filesystem_files_free`,
			expected: `node_filesystem_files{host="$host",mountpoint="/"} - on () group_left () node_filesystem_files_free`,
		},
		{
			name:     "common matchers",
			expr:     `node_filesystem_files{host="$host", mountpoint="/"} - node_filesystem_files_free`,
			expected: `node_filesystem_files{host="$host",mountpoint="/"} - node_filesystem_files_free{host="$host",mountpoint="/"}`,
		},
	}

	optimizers := []Optimizer{PropagateMatchersOptimizer{}}
	for _, tcase := range cases {
		tcase := tcase
		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()
			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			plan := New(expr, &query.Options{Start: time.Unix(0, 0), End: time.Unix(0, 0)}, PlanOptions{})
			optimizedPlan, _ := plan.Optimize(optimizers)
			expectedPlan := strings.Trim(spaces.ReplaceAllString(tcase.expected, " "), " ")
			testutil.Equals(t, expectedPlan, renderExprTree(optimizedPlan.Root()))
		})
	}
}

func TestTrimSorts(t *testing.T) {
	cases := []struct {
		name     string
		expr     string
		expected string
	}{
		// this test case is ok since the engine determines sorting order
		// before running optimziers
		{
			name:     "simple sort",
			expr:     "sort(X)",
			expected: "X",
		},
		{
			name:     "sort",
			expr:     "sum(sort(X))",
			expected: "sum(X)",
		},
		{
			name:     "nested",
			expr:     "sum(sort(rate(X[1m])))",
			expected: "sum(rate(X[1m]))",
		},
		{
			name:     "weirdly nested",
			expr:     "sum(sort(sqrt(sort(X))))",
			expected: "sum(sqrt(X))",
		},
		{
			name:     "sort in binary expression",
			expr:     "sort(sort(sqrt(X))/sort(sqrt(Y)))",
			expected: "sqrt(X) / sqrt(Y)",
		},
		{
			name:     "sort in argument to timestamp function",
			expr:     "timestamp(sort(X))",
			expected: "timestamp(sort(X))",
		},
	}
	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			exprPlan := New(expr, &query.Options{}, PlanOptions{})
			testutil.Equals(t, tcase.expected, exprPlan.Root().String())
		})
	}
}

func cleanUp(replacements map[string]*regexp.Regexp, expr string) string {
	for replacement, match := range replacements {
		expr = match.ReplaceAllString(expr, replacement)
	}
	return strings.Trim(expr, " ")
}
