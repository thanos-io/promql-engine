// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine

import (
	"context"
	"time"

	"github.com/efficientgo/core/errors"
	"github.com/thanos-community/promql-engine/physicalplan/model"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/util/stats"
)

type instantQuery struct {
	plan model.VectorOperator
	pool *model.VectorPool
	expr parser.Expr
	ts   time.Time
}

func newInstantQuery(plan model.VectorOperator, pool *model.VectorPool, expr parser.Expr, ts time.Time) promql.Query {
	return &instantQuery{
		plan: plan,
		pool: pool,
		expr: expr,
		ts:   ts,
	}
}

func (q *instantQuery) Exec(ctx context.Context) *promql.Result {
	// Handle case with strings early on as this does not need us to process samples.
	// TODO(saswatamcode): Modify models.StepVector to support all types and check during plan creation.
	switch e := q.expr.(type) {
	case *parser.StringLiteral:
		return &promql.Result{Value: promql.String{V: e.Val, T: q.ts.UnixMilli()}}
	}

	vs, err := q.plan.Next(ctx)
	if err != nil {
		return newErrResult(err)
	}

	if len(vs) == 0 {
		return &promql.Result{}
	}

	resultSeries, err := q.plan.Series(ctx)
	if err != nil {
		return newErrResult(err)
	}

	series := make([]promql.Series, len(resultSeries))
	for i := 0; i < len(resultSeries); i++ {
		series[i].Metric = resultSeries[i]
		series[i].Points = make([]promql.Point, 0, 1)
	}

	for _, vector := range vs {
		for _, sample := range vector.Samples {
			series[sample.ID].Points = append(series[sample.ID].Points, promql.Point{
				T: vector.T,
				V: sample.V,
			})
		}
		q.plan.GetPool().PutSamples(vector.Samples)
	}
	q.plan.GetPool().PutVectors(vs)

	var result parser.Value
	switch q.expr.Type() {
	case parser.ValueTypeMatrix:
		result = promql.Matrix(series)
	case parser.ValueTypeVector:
		// Convert matrix with one value per series into vector.
		vector := make(promql.Vector, len(resultSeries))
		for i := range series {
			// Point might have a different timestamp, force it to the evaluation
			// timestamp as that is when we ran the evaluation.
			vector[i] = promql.Sample{
				Metric: series[i].Metric,
				Point: promql.Point{
					V: series[i].Points[0].V,
					T: q.ts.UnixMilli(),
				},
			}
		}
		result = vector
	case parser.ValueTypeString:
		result = promql.Scalar{V: series[0].Points[0].V, T: q.ts.UnixMilli()}
	default:
		panic(errors.Newf("new.Engine.exec: unexpected expression type %q", q.expr.Type()))
	}

	return &promql.Result{
		Value: result,
	}
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
