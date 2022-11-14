// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package step_invariant

import (
	"context"
	"sync"

	"github.com/efficientgo/core/errors"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"

	"github.com/thanos-community/promql-engine/execution/model"
	"github.com/thanos-community/promql-engine/query"
)

type stepInvariantOperator struct {
	vectorPool       *model.VectorPool
	next             model.VectorOperator
	once             sync.Once
	duplicateResults bool

	series []labels.Labels

	mint int64
	maxt int64
	step int64
}

func (u *stepInvariantOperator) Explain() (me string, next []model.VectorOperator) {
	return "[*stepInvariantOperator]", []model.VectorOperator{u.next}
}

func NewStepInvariantOperator(
	pool *model.VectorPool,
	next model.VectorOperator,
	expr parser.Expr,
	opts *query.Options,
) (model.VectorOperator, error) {
	interval := opts.Step.Milliseconds()
	// We set interval to be at least 1.
	if interval == 0 {
		interval = 1
	}
	u := &stepInvariantOperator{
		vectorPool:       pool,
		next:             next,
		mint:             opts.Start.UnixMilli(),
		maxt:             opts.End.UnixMilli(),
		step:             interval,
		duplicateResults: true,
	}
	// We do not duplicate results for range selectors since result is a matrix
	// with their unique timestamps which does not depend on the step.
	switch expr.(type) {
	case *parser.MatrixSelector, *parser.SubqueryExpr:
		u.duplicateResults = false
	}

	return u, nil
}

func (u *stepInvariantOperator) Series(ctx context.Context) ([]labels.Labels, error) {
	var err error
	u.once.Do(func() {
		u.series, err = u.next.Series(ctx)
	})
	if err != nil {
		return nil, err
	}
	return u.series, nil
}

func (u *stepInvariantOperator) GetPool() *model.VectorPool {
	return u.vectorPool
}

func (u *stepInvariantOperator) Next(ctx context.Context) ([]model.StepVector, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	in, err := u.next.Next(ctx)
	if err != nil {
		return nil, err
	}
	if len(in) == 0 || len(in[0].Samples) == 0 {
		return nil, nil
	}
	if !u.duplicateResults {
		return in, nil
	}
	// Make sure we only have one step vector.
	if len(in) != 1 {
		return nil, errors.New("unexpected number of samples")
	}
	defer u.next.GetPool().PutVectors(in)

	result := u.vectorPool.GetVectorBatch()

	// Copy the evaluated step vector.
	sv := u.vectorPool.GetStepVector(in[0].T)
	sv.Samples = append(sv.Samples, in[0].Samples...)
	sv.SampleIDs = append(sv.SampleIDs, in[0].SampleIDs...)
	result = append(result, sv)
	u.next.GetPool().PutStepVector(in[0])

	for ts := u.mint + u.step; ts <= u.maxt; ts += u.step {
		sv := u.vectorPool.GetStepVector(ts)
		sv.Samples = append(sv.Samples, result[0].Samples...)
		sv.SampleIDs = append(sv.SampleIDs, result[0].SampleIDs...)
		result = append(result, sv)
	}

	return result, nil
}
