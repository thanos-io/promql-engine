package engine

import (
	"context"
	"sort"
	"sync"

	"github.com/fpetkovski/promql-engine/physicalplan/model"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/util/stats"
)

type rangeQuery struct {
	pool *model.VectorPool
	plan model.VectorOperator
	once sync.Once
}

func newRangeQuery(plan model.VectorOperator, pool *model.VectorPool) promql.Query {
	return &rangeQuery{
		pool: pool,
		plan: plan,
	}
}

func (q *rangeQuery) Exec(ctx context.Context) *promql.Result {
	resultSeries, err := q.plan.Series(ctx)
	if err != nil {
		return newErrResult(err)
	}

	series := make([]promql.Series, len(resultSeries))
	for i := 0; i < len(resultSeries); i++ {
		series[i].Metric = resultSeries[i]
	}
	for {
		r, err := q.plan.Next(ctx)
		if err != nil {
			return newErrResult(err)
		}
		if r == nil {
			break
		}

		for _, vector := range r {
			for _, sample := range vector.Samples {
				if len(series[sample.ID].Points) == 0 {
					series[sample.ID].Points = make([]promql.Point, 0, 121)
				}
				series[sample.ID].Points = append(series[sample.ID].Points, promql.Point{
					T: vector.T,
					V: sample.V,
				})
			}
			q.plan.GetPool().PutSamples(vector.Samples)
		}
		q.plan.GetPool().PutVectors(r)
	}

	result := make(promql.Matrix, 0, len(series))
	for _, s := range series {
		if len(s.Points) > 0 {
			result = append(result, s)
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
