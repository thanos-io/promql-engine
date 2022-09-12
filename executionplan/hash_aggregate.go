package executionplan

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/prometheus/prometheus/model/labels"

	"github.com/fpetkovski/promql-engine/model"

	"github.com/prometheus/prometheus/promql/parser"
)

type aggregate struct {
	downstream VectorOperator

	hashBuf    []byte
	vectorPool *model.VectorPool

	by          bool
	labels      []string
	aggregation parser.ItemType

	once   sync.Once
	tables []*aggregateTable
	series []labels.Labels
}

func NewAggregate(points *model.VectorPool, downstream VectorOperator, aggregation parser.ItemType, by bool, labels []string) (VectorOperator, error) {
	return &aggregate{
		downstream: downstream,
		vectorPool: points,

		by:          by,
		aggregation: aggregation,
		labels:      labels,
	}, nil
}

func (a *aggregate) Series(ctx context.Context) ([]labels.Labels, error) {
	var err error
	a.once.Do(func() { err = a.initOutputBuffers(ctx) })
	if err != nil {
		return nil, err
	}

	return a.series, nil
}

func (a *aggregate) initOutputBuffers(ctx context.Context) error {
	series, err := a.downstream.Series(ctx)
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

	a.tables = tables
	return nil
}

func (a *aggregate) GetPool() *model.VectorPool {
	return a.vectorPool
}

func (a *aggregate) Next(ctx context.Context) ([]model.StepVector, error) {
	in, err := a.downstream.Next(ctx)
	if err != nil {
		return nil, err
	}
	if in == nil {
		return nil, nil
	}
	defer a.downstream.GetPool().PutVectors(in)

	a.once.Do(func() { err = a.initOutputBuffers(ctx) })
	if err != nil {
		return nil, err
	}

	result := make([]model.StepVector, len(in))
	var wg sync.WaitGroup
	for i, vector := range in {
		wg.Add(1)
		go func(i int, vector model.StepVector) {
			defer wg.Done()

			table := a.tables[i]
			table.reset()
			for _, series := range vector.Samples {
				table.addSample(vector.T, series)
			}
			result[i] = table.toVector(a.vectorPool)
			a.downstream.GetPool().PutSamples(vector.Samples)
		}(i, vector)
	}
	wg.Wait()
	return result, nil
}

type newAccumulatorFunc func() *accumulator

type accumulator struct {
	AddFunc   func(v float64)
	ValueFunc func() float64
	HasValue  func() bool
	Reset     func()
}

func newAccumulator(expr parser.ItemType) (*accumulator, error) {
	hasValue := false
	t := parser.ItemTypeStr[expr]
	switch t {
	case "sum":
		var value float64
		return &accumulator{
			AddFunc: func(v float64) {
				hasValue = true
				value += v
			},
			ValueFunc: func() float64 { return value },
			HasValue:  func() bool { return hasValue },
			Reset: func() {
				hasValue = false
				value = 0
			},
		}, nil
	case "max":
		var value float64
		return &accumulator{
			AddFunc: func(v float64) {
				hasValue = true
				value = math.Max(value, v)
			},
			ValueFunc: func() float64 { return value },
			HasValue:  func() bool { return hasValue },
			Reset: func() {
				hasValue = false
				value = 0
			},
		}, nil
	case "min":
		var value float64
		return &accumulator{
			AddFunc: func(v float64) {
				hasValue = true
				value = math.Min(value, v)
			},
			ValueFunc: func() float64 { return value },
			HasValue:  func() bool { return hasValue },
			Reset: func() {
				hasValue = false
				value = 0
			},
		}, nil
	case "count":
		var value float64
		return &accumulator{
			AddFunc: func(v float64) {
				hasValue = true
				value += 1
			},
			ValueFunc: func() float64 { return value },
			HasValue:  func() bool { return hasValue },
			Reset: func() {
				hasValue = false
				value = 0
			},
		}, nil
	case "avg":
		var count, sum float64
		return &accumulator{
			AddFunc: func(v float64) {
				hasValue = true
				count += 1
				sum += v
			},
			ValueFunc: func() float64 { return sum / count },
			HasValue:  func() bool { return hasValue },
			Reset: func() {
				hasValue = false
				sum = 0
				count = 0
			},
		}, nil
	}
	return nil, fmt.Errorf("unknown aggregation function %s", t)
}
