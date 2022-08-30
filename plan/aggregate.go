package plan

import (
	"context"
	"fmt"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"math"
)

type aggregate struct {
	operator Operator
	by       bool
	labels   []string

	aggregation parser.ItemType
}

type aggregateTable struct {
	makeAccumulatorFunc newAccumulatorFunc
	table               map[uint64]*aggregateResult
}

func newAggregateTable(f newAccumulatorFunc) *aggregateTable {
	return &aggregateTable{
		makeAccumulatorFunc: f,
		table:               make(map[uint64]*aggregateResult),
	}
}

func (t *aggregateTable) addPoint(metric labels.Labels, p promql.Point, key uint64) {
	if _, ok := t.table[key]; !ok {
		t.table[key] = &aggregateResult{
			metric:      metric,
			timestamp:   p.T,
			accumulator: t.makeAccumulatorFunc(),
		}
	}
	t.table[key].accumulator.AddFunc(p)
}

type aggregateResult struct {
	metric      labels.Labels
	timestamp   int64
	accumulator *accumulator
}

func NewAggregate(input Operator, aggregation parser.ItemType, by bool, labels []string) (Operator, error) {
	return &aggregate{
		operator:    input,
		by:          by,
		aggregation: aggregation,
		labels:      labels,
	}, nil
}

func (a *aggregate) Next(ctx context.Context) (<-chan promql.Matrix, error) {
	in, err := a.operator.Next(ctx)
	if err != nil {
		return nil, err
	}

	out := make(chan promql.Matrix, 1024)
	go func() {
		defer close(out)
		buf := make([]byte, 1024)

		for matrix := range in {
			table := newAggregateTable(func() *accumulator {
				f, err := newAccumulator(a.aggregation)
				if err != nil {
					panic(err)
				}
				return f
			})
			for _, series := range matrix {
				buf := buf[:0]
				key, lbls := generateGroupingKey(series.Metric, a.labels, !a.by, buf)
				table.addPoint(lbls, series.Points[0], key)
			}

			a.send(table, out)
		}
	}()

	return out, nil
}

func (a *aggregate) send(table *aggregateTable, out chan promql.Matrix) {
	result := promql.Matrix{}
	for _, v := range table.table {
		series := promql.Series{
			Metric: v.metric,
			Points: []promql.Point{
				{
					T: v.timestamp,
					V: v.accumulator.ValueFunc(),
				},
			},
		}
		result = append(result, series)
	}

	out <- result
}

// groupingKey builds and returns the grouping key and series labels
// for the given metric and grouping labels.
func generateGroupingKey(metric labels.Labels, grouping []string, without bool, buf []byte) (uint64, labels.Labels) {
	if without {
		lb := labels.NewBuilder(metric)
		lb.Del(grouping...)
		key, _ := metric.HashWithoutLabels(buf, grouping...)
		return key, lb.Labels()
	}

	if len(grouping) == 0 {
		return 0, labels.Labels{}
	}

	lb := labels.NewBuilder(metric)
	lb.Keep(grouping...)
	key, _ := metric.HashForLabels(buf, grouping...)
	return key, lb.Labels()
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
