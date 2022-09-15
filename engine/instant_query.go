// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine

import (
	"context"

	"github.com/thanos-community/promql-engine/physicalplan/model"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/util/stats"
)

type instantQuery struct {
	plan model.VectorOperator
}

func newInstantQuery(plan model.VectorOperator) promql.Query {
	return &instantQuery{plan: plan}
}

func (q *instantQuery) Exec(ctx context.Context) *promql.Result {
	return nil

	//vs, err := q.plan.Next(ctx)
	//if err != nil {
	//	return newErrResult(err)
	//}
	//
	//if len(vs) == 0 {
	//	return &promql.Result{}
	//}
	//r := vs[len(vs)-1]
	//
	//sort.Slice(r, func(i, j int) bool {
	//	return labels.Compare(r[i].Metric, r[j].Metric) < 0
	//})
	//
	//return &promql.Result{
	//	Value: r,
	//}
}

// TODO(fpetkovski): Check if any resources can be released.
func (q *instantQuery) Close() {}

func (q *instantQuery) Statement() parser.Statement {
	return nil
}

func (q *instantQuery) Stats() *stats.Statistics {
	return &stats.Statistics{}
}

func (q *instantQuery) Cancel() {}

func (q *instantQuery) String() string { return "" }
