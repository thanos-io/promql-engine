// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine

import (
	"context"
	"sort"

	"github.com/go-kit/log"

	"github.com/thanos-community/promql-engine/physicalplan/model"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/util/stats"
)

type rangeQuery struct {
	cancel context.CancelFunc
	plan   model.VectorOperator
	logger log.Logger
	expr   parser.Expr
}

func newRangeQuery(expr parser.Expr, logger log.Logger, plan model.VectorOperator) promql.Query {
	return &rangeQuery{
		logger: logger,
		plan:   plan,
		expr:   expr,
	}
}

func (q *rangeQuery) Exec(ctx context.Context) (ret *promql.Result) {
	ret = &promql.Result{
		Value: promql.Vector{},
	}
	defer recoverEngine(q.logger, q.expr, &ret.Err)

	ctx, cancel := context.WithCancel(ctx)
	q.cancel = cancel

	// TODO(bwplotka): This feels weird. Why we have public Close method in the first place if we don't let user to use it.
	defer q.Close()

	resultSeries, err := q.plan.Series(ctx)
	if err != nil {
		return newErrResult(ret, err)
	}

	series := make([]promql.Series, len(resultSeries))
	for i := 0; i < len(resultSeries); i++ {
		series[i].Metric = resultSeries[i]
		series[i].Points = make([]promql.Point, 0, 121)
	}

	if err := getAllSeries(ctx, q.plan, series); err != nil {
		return newErrResult(ret, err)
	}

	result := make(promql.Matrix, 0, len(series))
	for _, s := range series {
		if len(s.Points) == 0 {
			continue
		}
		result = append(result, s)
	}

	sort.Sort(result)
	ret.Value = result
	return ret
}

func getAllSeries(ctx context.Context, plan model.VectorOperator, series []promql.Series) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			r, err := plan.Next(ctx)
			if err != nil {
				return err
			}
			if r == nil {
				return nil
			}

			for _, vector := range r {
				for i, s := range vector.SampleIDs {
					series[s].Points = append(series[s].Points, promql.Point{
						T: vector.T,
						V: vector.Samples[i],
					})
				}
				plan.GetPool().PutStepVector(vector)
			}
			plan.GetPool().PutVectors(r)
		}
	}
}

func (q *rangeQuery) Statement() parser.Statement { return nil }

func (q *rangeQuery) Stats() *stats.Statistics { return &stats.Statistics{} }

func (q *rangeQuery) Close() { q.Cancel() }

func (q *rangeQuery) String() string { return "" }

func (q *rangeQuery) Cancel() {
	if q.cancel != nil {
		q.cancel()
	}
}

func newErrResult(r *promql.Result, err error) *promql.Result {
	if r == nil {
		r = &promql.Result{}
	}
	if r.Err == nil && err != nil {
		r.Err = err
	}
	return r
}
