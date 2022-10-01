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
	cancel context.CancelFunc

	plan model.VectorOperator
	expr parser.Expr
	ts   time.Time
}

func newInstantQuery(plan model.VectorOperator, expr parser.Expr, ts time.Time) promql.Query {
	return &instantQuery{
		plan: plan,
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

	ctx, cancel := context.WithCancel(ctx)
	q.cancel = cancel

	defer q.Close()

	resultSeries, err := q.plan.Series(ctx)
	if err != nil {
		return newErrResult(err)
	}

	series := make([]promql.Series, len(resultSeries))
	for i := 0; i < len(resultSeries); i++ {
		series[i].Metric = resultSeries[i]
		series[i].Points = make([]promql.Point, 0, 1)
	}

	vs, err := q.plan.Next(ctx)
	if err != nil {
		return newErrResult(err)
	}

	if len(vs) == 0 {
		return &promql.Result{Value: promql.Vector{}}
	}

	for _, vector := range vs {
		for i, sample := range vector.SampleIDs {
			series[sample].Points = append(series[sample].Points, promql.Point{
				T: vector.T,
				V: vector.Samples[i],
			})
		}
		q.plan.GetPool().PutStepVector(vector)
	}
	q.plan.GetPool().PutVectors(vs)

	var result parser.Value
	switch q.expr.Type() {
	case parser.ValueTypeMatrix:
		result = promql.Matrix(series)
	case parser.ValueTypeVector:
		// Convert matrix with one value per series into vector.
		vector := make(promql.Vector, 0, len(resultSeries))
		for i := range series {
			if len(series[i].Points) == 0 {
				continue
			}
			// Point might have a different timestamp, force it to the evaluation
			// timestamp as that is when we ran the evaluation.
			vector = append(vector, promql.Sample{
				Metric: series[i].Metric,
				Point: promql.Point{
					V: series[i].Points[0].V,
					T: q.ts.UnixMilli(),
				},
			})
		}
		result = vector
	case parser.ValueTypeScalar:
		result = promql.Scalar{V: series[0].Points[0].V, T: q.ts.UnixMilli()}
	default:
		panic(errors.Newf("new.Engine.exec: unexpected expression type %q", q.expr.Type()))
	}

	return &promql.Result{
		Value: result,
	}
}

func (q *instantQuery) Statement() parser.Statement { return nil }

func (q *instantQuery) Stats() *stats.Statistics { return &stats.Statistics{} }

func (q *instantQuery) Close() { q.Cancel() }

func (q *instantQuery) String() string { return "" }

func (q *instantQuery) Cancel() {
	if q.cancel != nil {
		q.cancel()
	}
}
