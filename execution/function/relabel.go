// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package function

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/execution/telemetry"
	"github.com/thanos-io/promql-engine/logicalplan"
	"github.com/thanos-io/promql-engine/query"

	"github.com/efficientgo/core/errors"
	prommodel "github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/model/labels"
)

// relabelRule is one label_replace or label_join call with its arguments
// unwrapped from the plan. Only the fields belonging to Func are set.
type relabelRule struct {
	Func string // "label_replace" or "label_join"
	Dst  string // destination label name, used by both

	// label_replace only.
	Repl  string
	Src   string
	Regex string

	// label_join only.
	Sep  string
	Srcs []string
}

// IsRelabelFunc reports whether name is a function the relabel operator handles.
func IsRelabelFunc(name string) bool {
	return name == "label_replace" || name == "label_join"
}

// relabelRuleFromCall converts a label_replace or label_join call into a rule.
// It fails when an argument is not a string literal.
func relabelRuleFromCall(call *logicalplan.FunctionCall) (relabelRule, error) {
	if len(call.Args) < 2 {
		return relabelRule{}, fmt.Errorf("%s takes at least 2 arguments, got %d", call.Func.Name, len(call.Args))
	}
	dst, err := logicalplan.UnwrapString(call.Args[1])
	if err != nil {
		return relabelRule{}, err
	}
	rule := relabelRule{Func: call.Func.Name, Dst: dst}

	switch call.Func.Name {
	case "label_replace":
		if len(call.Args) != 5 {
			return relabelRule{}, fmt.Errorf("label_replace takes 5 arguments, got %d", len(call.Args))
		}
		if rule.Repl, err = logicalplan.UnwrapString(call.Args[2]); err != nil {
			return relabelRule{}, err
		}
		if rule.Src, err = logicalplan.UnwrapString(call.Args[3]); err != nil {
			return relabelRule{}, err
		}
		if rule.Regex, err = logicalplan.UnwrapString(call.Args[4]); err != nil {
			return relabelRule{}, err
		}
	case "label_join":
		if len(call.Args) < 3 {
			return relabelRule{}, fmt.Errorf("label_join takes at least 3 arguments, got %d", len(call.Args))
		}
		if rule.Sep, err = logicalplan.UnwrapString(call.Args[2]); err != nil {
			return relabelRule{}, err
		}
		for _, arg := range call.Args[3:] {
			src, err := logicalplan.UnwrapString(arg)
			if err != nil {
				return relabelRule{}, err
			}
			rule.Srcs = append(rule.Srcs, src)
		}
	default:
		return relabelRule{}, fmt.Errorf("not a relabel function: %s", call.Func.Name)
	}
	return rule, nil
}

// compiledRule is a relabelRule with its regex compiled and destination
// validated, done once rather than once per series.
type compiledRule struct {
	funcName string
	dst      string

	// label_replace only.
	repl  string
	src   string
	regex *regexp.Regexp

	// label_join only.
	sep  string
	srcs []string
}

// relabelOperator applies a list of relabel rules in one pass, producing one
// output label set per series rather than one per rule.
type relabelOperator struct {
	next  model.VectorOperator
	rules []compiledRule

	// Reported from Series(), where an argument error has always surfaced.
	compileErr error

	once   sync.Once
	series []labels.Labels
}

// NewRelabelOperator builds an operator applying every call of a relabel chain
// in one pass. A single call is a chain of one. Calls are innermost-first.
func NewRelabelOperator(
	next model.VectorOperator,
	calls []*logicalplan.FunctionCall,
	opts *query.Options,
) model.VectorOperator {
	oper := &relabelOperator{next: next, rules: make([]compiledRule, 0, len(calls))}
	for _, call := range calls {
		rule, err := relabelRuleFromCall(call)
		if err != nil {
			oper.compileErr = errors.Wrap(err, "unable to unwrap string argument")
			break
		}
		compiled, err := compileRelabelRule(rule)
		if err != nil {
			oper.compileErr = err
			break
		}
		oper.rules = append(oper.rules, compiled)
	}
	return telemetry.NewOperator(telemetry.NewTelemetry(oper, opts), oper)
}

// compileRelabelRule validates the destination and compiles the regex, so
// neither is repeated per series.
func compileRelabelRule(rule relabelRule) (compiledRule, error) {
	if !prommodel.UTF8Validation.IsValidLabelName(rule.Dst) {
		return compiledRule{}, errors.Newf("invalid destination label name in %s: %s", rule.Func, rule.Dst)
	}
	c := compiledRule{funcName: rule.Func, dst: rule.Dst}
	switch rule.Func {
	case "label_replace":
		regex, err := regexp.Compile("^(?:" + rule.Regex + ")$")
		if err != nil {
			return compiledRule{}, errors.Newf("invalid regular expression in label_replace(): %s", rule.Regex)
		}
		c.repl, c.src, c.regex = rule.Repl, rule.Src, regex
	case "label_join":
		c.sep, c.srcs = rule.Sep, rule.Srcs
	default:
		return compiledRule{}, errors.Newf("invalid function name for relabel operator: %s", rule.Func)
	}
	return c, nil
}

func (o *relabelOperator) String() string {
	if len(o.rules) > 1 {
		return fmt.Sprintf("[relabel] %d rules", len(o.rules))
	}
	return "[relabel]"
}

func (o *relabelOperator) Explain() (next []model.VectorOperator) {
	return []model.VectorOperator{o.next}
}

func (o *relabelOperator) Series(ctx context.Context) ([]labels.Labels, error) {
	var err error
	o.once.Do(func() { err = o.loadSeries(ctx) })
	return o.series, err
}

func (o *relabelOperator) Next(ctx context.Context, buf []model.StepVector) (int, error) {
	return o.next.Next(ctx, buf)
}

func (o *relabelOperator) loadSeries(ctx context.Context) error {
	// Child first: a chain resolves bottom-up, so its error wins over ours.
	series, err := o.next.Series(ctx)
	if err != nil {
		return err
	}
	if o.compileErr != nil {
		return o.compileErr
	}
	o.series = make([]labels.Labels, len(series))

	// One builder and two scratch buffers for the whole pass. Builder.Reset
	// clears pending adds and deletes, so nothing leaks between series.
	builder := labels.NewBuilder(labels.EmptyLabels())
	var expanded []byte
	var srcVals []string

	for i, s := range series {
		builder.Reset(s)
		changed := false

		for _, rule := range o.rules {
			switch rule.funcName {
			case "label_replace":
				// Read through the builder, not the input: a later rule may
				// use a label an earlier rule just wrote.
				srcVal := builder.Get(rule.src)
				matches := rule.regex.FindStringSubmatchIndex(srcVal)
				if len(matches) == 0 {
					continue
				}
				expanded = rule.regex.ExpandString(expanded[:0], rule.repl, srcVal, matches)
				// Set("") delegates to Del, matching label_replace.
				builder.Set(rule.dst, string(expanded))
			case "label_join":
				srcVals = srcVals[:0]
				for _, src := range rule.srcs {
					srcVals = append(srcVals, builder.Get(src))
				}
				builder.Set(rule.dst, strings.Join(srcVals, rule.sep))
			}
			changed = true
		}

		if !changed {
			// No rule applied, so reuse the input rather than copying it.
			o.series[i] = s
			continue
		}
		o.series[i] = builder.Labels()
	}
	return nil
}
