// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine

import (
	"context"
	"sort"
	"sync"

	"github.com/thanos-community/promql-engine/physicalplan/model"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/util/stats"
)

type rangeQuery struct {
	once sync.Once
	plan model.VectorOperator
}

func newRangeQuery(plan model.VectorOperator) promql.Query {
	return &rangeQuery{
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
		series[i].Points = make([]promql.Point, 0, 121)
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
			for i, s := range vector.SampleIDs {
				series[s].Points = append(series[s].Points, promql.Point{
					T: vector.T,
					V: vector.Samples[i],
				})
			}
			q.plan.GetPool().PutStepVector(vector)
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
