package executionplan

import (
	"context"
	"fmt"
	"math"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
)

type aggregate struct {
	operator VectorOperator
	hashBuf  []byte

	by          bool
	labels      []string
	aggregation parser.ItemType
	table       *aggregateTable
}

func NewAggregate(input VectorOperator, aggregation parser.ItemType, by bool, labels []string) (VectorOperator, error) {
	hashBuf := make([]byte, 1024)
	table := newAggregateTable(
		newGroupingKeyGenerator(labels, !by, hashBuf),
		func() *accumulator {
			f, err := newAccumulator(aggregation)
			if err != nil {
				panic(err)
			}
			return f
		},
	)

	return &aggregate{
		operator: input,
		table:    table,

		by:          by,
		aggregation: aggregation,
		labels:      labels,
	}, nil
}

func (a *aggregate) Next(ctx context.Context) (promql.Vector, error) {
	in, err := a.operator.Next(ctx)
	if err != nil {
		return nil, err
	}
	if in == nil {
		return nil, nil
	}

	a.table.reset()
	for i, series := range in {
		a.table.addSample(i, series)
	}

	return a.table.toVector(), nil
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
