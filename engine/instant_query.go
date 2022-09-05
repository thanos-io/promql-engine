package engine

import (
	"context"
	"sort"

	"github.com/fpetkovski/promql-engine/executionplan"

	"github.com/prometheus/prometheus/model/labels"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/util/stats"
)

type instantQuery struct {
	plan executionplan.VectorOperator
}

func newInstantQuery(plan executionplan.VectorOperator) promql.Query {
	return &instantQuery{plan: plan}
}

func (q *instantQuery) Exec(ctx context.Context) *promql.Result {
	r, err := q.plan.Next(ctx)
	if err != nil {
		return newErrResult(err)
	}

	sort.Slice(r, func(i, j int) bool {
		return labels.Compare(r[i].Metric, r[j].Metric) < 0
	})

	return &promql.Result{
		Value: r,
	}
}

func (q *instantQuery) Close() {
	//TODO implement me
	panic("implement me")
}

func (q *instantQuery) Statement() parser.Statement {
	//TODO implement me
	panic("implement me")
}

func (q *instantQuery) Stats() *stats.Statistics {
	//TODO implement me
	panic("implement me")
}

func (q *instantQuery) Cancel() {
	//TODO implement me
	panic("implement me")
}

func (q *instantQuery) String() string {
	//TODO implement me
	panic("implement me")
}
