package executionplan

import (
	"context"
	"fmt"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"math"
)

type aggregate struct {
	operator VectorOperator
	by       bool
	labels   []string

	aggregation parser.ItemType
}

func NewAggregate(input VectorOperator, aggregation parser.ItemType, by bool, labels []string) (VectorOperator, error) {
	return &aggregate{
		operator:    input,
		by:          by,
		aggregation: aggregation,
		labels:      labels,
	}, nil
}

func (a *aggregate) Next(ctx context.Context) (<-chan promql.Vector, error) {
	in, err := a.operator.Next(ctx)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 1024)
	table := newAggregateTable(
		newGroupingKeyGenerator(a.labels, !a.by, buf),
		func() *accumulator {
			f, err := newAccumulator(a.aggregation)
			if err != nil {
				panic(err)
			}
			return f
		},
	)

	out := make(chan promql.Vector, 121)
	go func() {
		defer close(out)

		for vector := range in {
			table.reset()
			for i, series := range vector {
				table.addSample(i, series)
			}

			a.send(table, out)
		}
	}()

	return out, nil
}

func (a *aggregate) send(table *aggregateTable, out chan promql.Vector) {
	result := make(promql.Vector, 0, len(table.table))
	for _, v := range table.table {
		result = append(result, promql.Sample{
			Metric: v.metric,
			Point: promql.Point{
				T: v.timestamp,
				V: v.accumulator.ValueFunc(),
			},
		})
	}

	out <- result
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
