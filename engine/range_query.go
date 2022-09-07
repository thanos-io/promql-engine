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
	seriesMap := make(map[uint64]*promql.Series)
	for {
		r, err := q.plan.Next(ctx)
		if err != nil {
			return newErrResult(err)
		}
		if r == nil {
			break
		}

		for _, vector := range r {
			for _, sample := range vector {
				if _, ok := seriesMap[sample.ID]; !ok {
					seriesMap[sample.ID] = &promql.Series{
						Metric: sample.Metric,
						Points: make([]promql.Point, 0),
					}
				}
				series := seriesMap[sample.ID]
				series.Points = append(seriesMap[sample.ID].Points, sample.Point)
			}
		}
	}

	result := make(promql.Matrix, 0, len(seriesMap))
	for _, series := range seriesMap {
		result = append(result, *series)
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
