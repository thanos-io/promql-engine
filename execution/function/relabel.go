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

func IsRelabelFunc(name string) bool {
	return name == "label_replace" || name == "label_join"
}

func unwrapString(arg logicalplan.Node) (string, error) {
	s, err := logicalplan.UnwrapString(arg)
	if err != nil {
		return "", errors.Wrap(err, "unable to unwrap string argument")
	}
	return s, nil
}

func compiledRuleFromCall(call *logicalplan.FunctionCall) (compiledRule, error) {
	dst, err := unwrapString(call.Args[1])
	if err != nil {
		return compiledRule{}, err
	}
	if !prommodel.UTF8Validation.IsValidLabelName(dst) {
		return compiledRule{}, errors.Newf("invalid destination label name in %s: %s", call.Func.Name, dst)
	}
	rule := compiledRule{funcName: call.Func.Name, dst: dst}

	switch call.Func.Name {
	case "label_replace":
		if rule.repl, err = unwrapString(call.Args[2]); err != nil {
			return compiledRule{}, err
		}
		if rule.src, err = unwrapString(call.Args[3]); err != nil {
			return compiledRule{}, err
		}
		regex, err := unwrapString(call.Args[4])
		if err != nil {
			return compiledRule{}, err
		}
		if rule.regex, err = regexp.Compile("^(?:" + regex + ")$"); err != nil {
			return compiledRule{}, errors.Newf("invalid regular expression in label_replace(): %s", regex)
		}
	case "label_join":
		if rule.sep, err = unwrapString(call.Args[2]); err != nil {
			return compiledRule{}, err
		}
		for _, arg := range call.Args[3:] {
			src, err := unwrapString(arg)
			if err != nil {
				return compiledRule{}, err
			}
			rule.srcs = append(rule.srcs, src)
		}
	default:
		return compiledRule{}, errors.Newf("not a relabel function: %s", call.Func.Name)
	}
	return rule, nil
}

type compiledRule struct {
	funcName string // "label_replace" or "label_join"
	dst      string // destination label name, used by both

	// label_replace only.
	repl  string
	src   string
	regex *regexp.Regexp

	// label_join only.
	sep  string
	srcs []string
}

type relabelOperator struct {
	next  model.VectorOperator
	calls []*logicalplan.FunctionCall

	once   sync.Once
	series []labels.Labels
}

func NewRelabelOperator(
	next model.VectorOperator,
	calls []*logicalplan.FunctionCall,
	opts *query.Options,
) model.VectorOperator {
	oper := &relabelOperator{next: next, calls: calls}
	return telemetry.NewOperator(telemetry.NewTelemetry(oper, opts), oper)
}

func compileRules(calls []*logicalplan.FunctionCall) ([]compiledRule, error) {
	rules := make([]compiledRule, 0, len(calls))
	for _, call := range calls {
		rule, err := compiledRuleFromCall(call)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (o *relabelOperator) String() string {
	if len(o.calls) > 1 {
		return fmt.Sprintf("[relabel] %d rules", len(o.calls))
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
	series, err := o.next.Series(ctx)
	if err != nil {
		return err
	}
	rules, err := compileRules(o.calls)
	if err != nil {
		return err
	}
	o.series = make([]labels.Labels, len(series))

	builder := labels.NewBuilder(labels.EmptyLabels())
	var expanded []byte
	var srcVals []string

	for i, s := range series {
		builder.Reset(s)
		changed := false

		for _, rule := range rules {
			// Set("") delegates to Del, so an empty value deletes the destination.
			switch rule.funcName {
			case "label_replace":
				// Read the builder: a later rule may use what an earlier one wrote.
				srcVal := builder.Get(rule.src)
				matches := rule.regex.FindStringSubmatchIndex(srcVal)
				if len(matches) == 0 {
					continue
				}
				expanded = rule.regex.ExpandString(expanded[:0], rule.repl, srcVal, matches)
				if len(expanded) == 0 && builder.Get(rule.dst) == "" {
					continue
				}
				builder.Set(rule.dst, string(expanded))
			case "label_join":
				srcVals = srcVals[:0]
				for _, src := range rule.srcs {
					srcVals = append(srcVals, builder.Get(src))
				}
				joined := strings.Join(srcVals, rule.sep)
				if joined == "" && builder.Get(rule.dst) == "" {
					continue
				}
				builder.Set(rule.dst, joined)
			}
			changed = true
		}

		if !changed {
			o.series[i] = s
			continue
		}
		o.series[i] = builder.Labels()
	}
	return nil
}
