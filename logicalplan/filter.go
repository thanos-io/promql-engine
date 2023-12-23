// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"fmt"
	"strings"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/promql/parser/posrange"
)

// VectorSelector is vector selector with additional configuration set by optimizers.
// TODO(fpetkovski): Consider replacing all VectorSelector nodes with this one as the first step in the plan.
// This should help us avoid dealing with both types in the entire codebase.
type VectorSelector struct {
	*parser.VectorSelector

	Filters   []*labels.Matcher
	BatchSize int64

	Shards int
	N      int
}

func (f VectorSelector) String() string {
	var b strings.Builder
	var needComma bool
	b.WriteString(f.VectorSelector.String())
	b.WriteRune('[')
	if len(f.Filters) > 0 {
		b.WriteString(fmt.Sprintf("filters=%s", f.Filters))
		needComma = true
	}
	if f.BatchSize > 0 {
		if needComma {
			b.WriteRune(',')
		}
		b.WriteString(fmt.Sprintf("batch=%d", f.BatchSize))
		needComma = true
	}
	if f.Shards > 0 {
		if needComma {
			b.WriteRune(',')
		}
		b.WriteString(fmt.Sprintf("shard=%d/%d", f.N, f.Shards))
	}
	b.WriteRune(']')

	return b.String()
}

func (f VectorSelector) Pretty(level int) string { return f.String() }

func (f VectorSelector) PositionRange() posrange.PositionRange { return posrange.PositionRange{} }

func (f VectorSelector) Type() parser.ValueType { return parser.ValueTypeVector }

func (f VectorSelector) PromQLExpr() {}
