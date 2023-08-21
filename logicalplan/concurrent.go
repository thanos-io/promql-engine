// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

// Copyright 2013 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logicalplan

import (
	"fmt"

	"github.com/thanos-io/promql-engine/parser"
	"github.com/thanos-io/promql-engine/query"
)

type Concurrent struct {
	Concurrency int

	Expr parser.Expr
}

func (r Concurrent) String() string {
	return fmt.Sprintf("concurrent(%d, %s)", r.Concurrency, r.Expr.String())
}

func (r Concurrent) Pretty(level int) string { return r.String() }

func (r Concurrent) PositionRange() parser.PositionRange { return parser.PositionRange{} }

func (r Concurrent) Type() parser.ValueType { return parser.ValueTypeMatrix }

func (r Concurrent) PromQLExpr() {}

type ConcurrentExecutionOptimizer struct {
}

// TODO: move over creation of concurrency operators completely to here
func (m ConcurrentExecutionOptimizer) Optimize(plan parser.Expr, _ *query.Options) parser.Expr {
	traverse(&plan, func(current *parser.Expr) {
		if current == nil {
			return
		}
		switch (*current).(type) {
		case *parser.AggregateExpr, Deduplicate, RemoteExecution:
			*current = Concurrent{
				Concurrency: 2,
				Expr:        *current,
			}
		}
	})
	return plan
}
