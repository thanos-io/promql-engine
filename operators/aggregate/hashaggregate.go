package aggregate

import (
	"context"
	"sync"

	"github.com/fpetkovski/promql-engine/worker"

	"github.com/fpetkovski/promql-engine/operators/model"

	"github.com/prometheus/prometheus/model/labels"

	"github.com/prometheus/prometheus/promql/parser"
)

type aggregate struct {
	next model.Vector

	hashBuf    []byte
	vectorPool *model.VectorPool

	by          bool
	labels      []string
	aggregation parser.ItemType

	once   sync.Once
	tables []*aggregateTable
	series []labels.Labels

	workers worker.Group
}

func NewHashAggregate(
	points *model.VectorPool,
	next model.Vector,
	aggregation parser.ItemType,
	by bool,
	labels []string,
) (model.Vector, error) {
	a := &aggregate{
		next:       next,
		vectorPool: points,

		by:          by,
		aggregation: aggregation,
		labels:      labels,
	}
	a.workers = worker.NewGroup(10, a.workerTask)
	a.workers.Start()

	return a, nil
}

func (a *aggregate) Series(ctx context.Context) ([]labels.Labels, error) {
	var err error
	a.once.Do(func() { err = a.initOutputBuffers(ctx) })
	if err != nil {
		return nil, err
	}

	return a.series, nil
}

func (a *aggregate) GetPool() *model.VectorPool {
	return a.vectorPool
}

func (a *aggregate) Next(ctx context.Context) ([]model.StepVector, error) {
	in, err := a.next.Next(ctx)
	if err != nil {
		return nil, err
	}
	if in == nil {
		a.workers.Shutdown()
		return nil, nil
	}
	defer a.next.GetPool().PutVectors(in)

	a.once.Do(func() { err = a.initOutputBuffers(ctx) })
	if err != nil {
		return nil, err
	}

	result := make([]model.StepVector, len(in))
	for i, vector := range in {
		a.workers[i].Send(vector)
	}

	for i, vector := range in {
		a.workers[i].Done()
		result[i] = a.workers[i].GetOutput()
		a.next.GetPool().PutSamples(vector.Samples)

	}
	return result, nil
}

func (a *aggregate) shutdownWorkers() {
	for i := 0; i < len(a.workers); i++ {
		a.workers[i].Shutdown()
	}
}

func (a *aggregate) initOutputBuffers(ctx context.Context) error {
	series, err := a.next.Series(ctx)
	if err != nil {
		return err
	}

	inputCache := make([]uint64, len(series))
	outputMap := make(map[uint64]*aggregateResult)
	outputCache := make([]*aggregateResult, 0)

	for i := 0; i < len(series); i++ {
		buf := make([]byte, 128)
		hash, _, lbls := hashMetric(series[i], !a.by, a.labels, buf)

		output, ok := outputMap[hash]
		if !ok {
			output = &aggregateResult{
				metric:   lbls,
				sampleID: uint64(len(outputCache)),
			}
			outputMap[hash] = output
			outputCache = append(outputCache, output)
		}

		inputCache[i] = output.sampleID
	}

	var wg sync.WaitGroup
	a.series = make([]labels.Labels, len(outputCache))
	for i := 0; i < len(outputCache); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			a.series[i] = outputCache[i].metric
		}(i)
	}
	wg.Wait()

	tables := make([]*aggregateTable, 10)
	for i := 0; i < len(tables); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tables[i] = newAggregateTable(inputCache, outputCache, func() *accumulator {
				f, err := newAccumulator(a.aggregation)
				if err != nil {
					panic(err)
				}
				return f
			})
		}(i)
	}
	wg.Wait()

	a.vectorPool.SetStepSamplesSize(len(outputCache))
	a.tables = tables
	return nil
}

func (a *aggregate) workerTask(workerID int, vector model.StepVector) model.StepVector {
	table := a.tables[workerID]
	table.reset()
	for _, series := range vector.Samples {
		table.addSample(vector.T, series)
	}
	return table.toVector(a.vectorPool)
}
