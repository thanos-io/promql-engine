// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"strings"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/thanos-io/promql-engine/parser"
)

const hasRegexChars = ".+|[]{}^$*?+()\\"

// RegexMatchersFunctions function will optimize the query by avoiding unnecessary regex match.
// For example: up{app=~"something"} will be optimized to up{app="something"}.
type RegexMatchersFunctions struct {
}

func (RegexMatchersFunctions) Optimize(expr parser.Expr, _ *Opts) parser.Expr {
	traverse(&expr, func(node *parser.Expr) {
		switch e := (*node).(type) {
		case *parser.VectorSelector:
			for matcherIndex, lm := range e.LabelMatchers {
				if lm.Type == labels.MatchRegexp {
					if !strings.ContainsAny(lm.Value, hasRegexChars) {
						nm, err := labels.NewMatcher(labels.MatchEqual, lm.Name, lm.Value)
						if err != nil {
							continue
						}
						e.LabelMatchers[matcherIndex] = nm
					}
				}
			}
		}
	})
	return expr
}
