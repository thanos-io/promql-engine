package executionplan

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/fpetkovski/promql-engine/model"

	"github.com/prometheus/prometheus/promql/parser"
)

type aggregate struct {
	operator VectorOperator

	hashBuf    []byte
	vectorPool *model.VectorPool

	by          bool
	labels      []string
	aggregation parser.ItemType
	tables      []*aggregateTable
}

func NewAggregate(points *model.VectorPool, input VectorOperator, aggregation parser.ItemType, by bool, labels []string) (VectorOperator, error) {
	keys := make([]groupingKey, 100000)
	for i := 0; i < len(keys); i++ {
		keys[i] = groupingKey{
			once: &sync.Once{},
		}
	}
	tables := make([]*aggregateTable, 30)
	for i := 0; i < 30; i++ {
		hashBuf := make([]byte, 128)
		tables[i] = newAggregateTable(
			newGroupingKeyGenerator(labels, !by, hashBuf),
			func() *accumulator {
				f, err := newAccumulator(aggregation)
				if err != nil {
					panic(err)
				}
				return f
			},
			keys,
		)
	}

	return &aggregate{
		operator: input,
		tables:   tables,

		vectorPool: points,

		by:          by,
		aggregation: aggregation,
		labels:      labels,
	}, nil
}

func (a *aggregate) Next(ctx context.Context) ([]model.StepVector, error) {
	in, err := a.operator.Next(ctx)
	if err != nil {
		return nil, err
	}
	if in == nil {
		return nil, nil
	}

	defer a.vectorPool.Put(in)

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
			result[i] = table.toVector()
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
		}, nil
	}
	return nil, fmt.Errorf("unknown aggregation function %s", t)
}
