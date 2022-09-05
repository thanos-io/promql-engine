package engine

import (
	"context"
	"fpetkovski/promql-engine/executionplan"
	"sort"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/util/stats"
)

type query struct {
	plan executionplan.VectorOperator
}

func newQuery(plan executionplan.VectorOperator) promql.Query {
	return &query{plan: plan}
}

func (q *query) Exec(ctx context.Context) *promql.Result {
	var result promql.Matrix = nil
	for {
		r, err := q.plan.Next(ctx)
		if err != nil {
			return newErrResult(err)
		}
		if r == nil {
			break
		}
		if result == nil {
			result = make(promql.Matrix, len(r))
		}
		for i, v := range r {
			if v.Point.T == -1 {
				continue
			}
			result[i].Metric = v.Metric
			result[i].Points = append(result[i].Points, v.Point)
		}
	}

	sort.Sort(result)
	return &promql.Result{
		Value: result,
	}
}

func (q *query) Close() {
	//TODO implement me
	panic("implement me")
}

func (q *query) Statement() parser.Statement {
	//TODO implement me
	panic("implement me")
}

func (q *query) Stats() *stats.Statistics {
	//TODO implement me
	panic("implement me")
}

func (q *query) Cancel() {
	//TODO implement me
	panic("implement me")
}

func (q *query) String() string {
	//TODO implement me
	panic("implement me")
}

func newErrResult(err error) *promql.Result {
	return &promql.Result{Err: err}
}
