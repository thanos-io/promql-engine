package remote

import (
	"context"
	"fmt"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"sort"
	"sync"

	"github.com/thanos-community/promql-engine/execution/model"
)

type Execution struct {
	query      promql.Query
	result     promql.Matrix
	timeStamps []int64

	once       sync.Once
	pool       *model.VectorPool
	stepsBatch int
}

func NewExecution(
	query promql.Query,
	pool *model.VectorPool,
	stepsBatch int,
) *Execution {
	return &Execution{
		query:      query,
		pool:       pool,
		stepsBatch: stepsBatch,
	}
}

func (e *Execution) Series(ctx context.Context) ([]labels.Labels, error) {
	var err error
	e.once.Do(func() { err = e.execute(ctx) })
	if err != nil {
		return nil, err
	}

	series := make([]labels.Labels, 0, len(e.result))
	for _, s := range e.result {
		series = append(series, s.Metric)
	}
	return series, nil
}

func (e *Execution) Next(ctx context.Context) ([]model.StepVector, error) {
	var err error
	e.once.Do(func() { err = e.execute(ctx) })
	if err != nil {
		return nil, err
	}

	if len(e.timeStamps) == 0 {
		return nil, nil
	}

	out := e.pool.GetVectorBatch()
	for i := 0; i < e.stepsBatch && i < len(e.timeStamps); i++ {
		ts := e.timeStamps[i]
		for seriesID := range e.result {
			if len(e.result[seriesID].Points) == 0 {
				continue
			}

			if e.result[seriesID].Points[0].T == ts {
				if i >= len(out) {
					out = append(out, e.pool.GetStepVector(ts))
				}
				n := len(out) - 1

				out[n].SampleIDs = append(out[n].SampleIDs, uint64(seriesID))
				out[n].Samples = append(out[n].Samples, e.result[seriesID].Points[0].V)
				e.result[seriesID].Points = e.result[seriesID].Points[1:]
			}
		}
	}
	e.timeStamps = e.timeStamps[len(out):]

	return out, nil
}

func (e *Execution) GetPool() *model.VectorPool {
	return e.pool
}

func (e *Execution) Explain() (me string, next []model.VectorOperator) {
	return fmt.Sprintf("[*remoteExec] %s", e.query), nil
}

func (e *Execution) execute(ctx context.Context) error {
	result := e.query.Exec(ctx)
	if result.Err != nil {
		return result.Err
	}

	matrix, err := result.Matrix()
	if err != nil {
		return err
	}

	timestampSet := make(map[int64]struct{})
	for _, series := range matrix {
		for _, p := range series.Points {
			timestampSet[p.T] = struct{}{}
		}
	}

	e.timeStamps = make([]int64, 0, len(timestampSet))
	for ts := range timestampSet {
		e.timeStamps = append(e.timeStamps, ts)
	}
	sort.Slice(e.timeStamps, func(i, j int) bool {
		return e.timeStamps[i] < e.timeStamps[j]
	})

	e.result = matrix
	return nil
}
