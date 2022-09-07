package executionplan

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/fpetkovski/promql-engine/model"

	"github.com/fpetkovski/promql-engine/points"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
)

type aggregate struct {
	operator VectorOperator

	hashBuf []byte
	points  *points.Pool

	by          bool
	labels      []string
	aggregation parser.ItemType
	tables      []*aggregateTable
}

func NewAggregate(points *points.Pool, input VectorOperator, aggregation parser.ItemType, by bool, labels []string) (VectorOperator, error) {
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

		points: points,

		by:          by,
		aggregation: aggregation,
		labels:      labels,
	}, nil
}

func (a *aggregate) Next(ctx context.Context) ([]model.Vector, error) {
	in, err := a.operator.Next(ctx)
	if err != nil {
		return nil, err
	}
	if in == nil {
		return nil, nil
	}

	result := make([]model.Vector, len(in))
	var wg sync.WaitGroup
	for i, vector := range in {
		wg.Add(1)
		go func(i int, vector model.Vector) {
			defer wg.Done()
			table := a.tables[i]
			table.reset()
			for _, series := range vector {
				table.addSample(series)
			}
			result[i] = table.toVector()
		}(i, vector)
	}
	wg.Wait()
	return result, nil
}

type newAccumulatorFunc func() *accumulator

type accumulator struct {
	AddFunc   func(p promql.Point)
	ValueFunc func() float64
}

func newAccumulator(expr parser.ItemType) (*accumulator, error) {
	t := parser.ItemTypeStr[expr]
	switch t {
	case "sum":
		var value float64
		return &accumulator{
			AddFunc: func(p promql.Point) {
				value += p.V
			},
			ValueFunc: func() float64 {
				return value
			},
		}, nil
	case "max":
		var value float64
		return &accumulator{
			AddFunc: func(p promql.Point) {
				value = math.Max(value, p.V)
			},
			ValueFunc: func() float64 {
				return value
			},
		}, nil
	case "min":
		var value float64
		return &accumulator{
			AddFunc: func(p promql.Point) {
				value = math.Min(value, p.V)
			},
			ValueFunc: func() float64 {
				return value
			},
		}, nil
	case "count":
		var value float64
		return &accumulator{
			AddFunc: func(p promql.Point) {
				value += 1
			},
			ValueFunc: func() float64 {
				return value
			},
		}, nil
	case "avg":
		var count, sum float64
		return &accumulator{
			AddFunc: func(p promql.Point) {
				count += 1
				sum += p.V
			},
			ValueFunc: func() float64 {
				return sum / count
			},
		}, nil
	}
	return nil, fmt.Errorf("unknown aggregation function %s", t)
}
