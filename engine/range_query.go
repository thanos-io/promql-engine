package engine

import (
	"context"
	"sort"

	"github.com/fpetkovski/promql-engine/executionplan"

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

// TODO(fpetkovski): Check if any resources can be released.
func (q *rangeQuery) Close() {}

func (q *rangeQuery) Statement() parser.Statement {
	return nil
}

func (q *rangeQuery) Stats() *stats.Statistics {
	return &stats.Statistics{}
}

func (q *rangeQuery) Cancel() {}

func (q *rangeQuery) String() string { return "" }

func newErrResult(err error) *promql.Result {
	return &promql.Result{Err: err}
}
