package engine

import (
	"context"
	"sort"

	"fpetkovski/promql-engine/executionplan"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/util/stats"
)

type rangeQuery struct {
	plan executionplan.VectorOperator
}

func newRangeQuery(plan executionplan.VectorOperator) promql.Query {
	return &rangeQuery{plan: plan}
}

func (q *rangeQuery) Exec(ctx context.Context) *promql.Result {
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

func (q *rangeQuery) Close() {
	//TODO implement me
	panic("implement me")
}

func (q *rangeQuery) Statement() parser.Statement {
	//TODO implement me
	panic("implement me")
}

func (q *rangeQuery) Stats() *stats.Statistics {
	//TODO implement me
	panic("implement me")
}

func (q *rangeQuery) Cancel() {
	//TODO implement me
	panic("implement me")
}

func (q *rangeQuery) String() string {
	//TODO implement me
	panic("implement me")
}

func newErrResult(err error) *promql.Result {
	return &promql.Result{Err: err}
}
