package remote

import (
	"context"
	"fmt"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/thanos-community/promql-engine/execution/scan"
	"github.com/thanos-community/promql-engine/query"
	"sync"

	"github.com/thanos-community/promql-engine/execution/model"
)

type Execution struct {
	query promql.Query
	pool  *model.VectorPool
	opts  *query.Options

	once           sync.Once
	vectorSelector model.VectorOperator
}

func NewExecution(query promql.Query, pool *model.VectorPool, opts *query.Options) *Execution {
	return &Execution{
		query: query,
		pool:  pool,
		opts:  opts,
	}
}

func (e *Execution) Series(ctx context.Context) ([]labels.Labels, error) {
	var err error
	e.once.Do(func() { err = e.execute(ctx) })
	if err != nil {
		return nil, err
	}

	return e.vectorSelector.Series(ctx)
}

func (e *Execution) Next(ctx context.Context) ([]model.StepVector, error) {
	return e.vectorSelector.Next(ctx)
}

func (e *Execution) GetPool() *model.VectorPool {
	return e.pool
}

func (e *Execution) Explain() (me string, next []model.VectorOperator) {
	return fmt.Sprintf("[*remoteExec] %s", e.query), nil
}

func (e *Execution) execute(ctx context.Context) error {
	defer e.query.Close()

	result := e.query.Exec(ctx)
	if result.Err != nil {
		return result.Err
	}

	matrix, err := result.Matrix()
	if err != nil {
		return err
	}

	e.vectorSelector = scan.NewVectorSelector(e.pool, newStorageAdapter(matrix), e.opts, 0, 0, 1)
	return nil
}
