// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package aggregate

import (
	"fmt"
	"math"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/thanos-community/promql-engine/physicalplan/model"
)

type aggregateTable interface {
	aggregate(vector model.StepVector)
	toVector(pool *model.VectorPool) model.StepVector
	size() int
}

type scalarTable struct {
	timestamp    int64
	inputs       []uint64
	outputs      []*model.Series
	accumulators []*accumulator
}

func newScalarTables(stepsBatch int, inputCache []uint64, outputCache []*model.Series, a parser.ItemType) []aggregateTable {
	tables := make([]aggregateTable, stepsBatch)
	for i := 0; i < len(tables); i++ {
		tables[i] = newScalarTable(inputCache, outputCache, func() *accumulator {
			f, err := newAccumulator(a)
			if err != nil {
				panic(err)
			}
			return f
		})
	}
	return tables
}

func newScalarTable(inputSampleIDs []uint64, outputs []*model.Series, makeAccumulator newAccumulatorFunc) *scalarTable {
	accumulators := make([]*accumulator, len(outputs))
	for i := 0; i < len(outputs); i++ {
		accumulators[i] = makeAccumulator()
	}
	return &scalarTable{
		inputs:       inputSampleIDs,
		outputs:      outputs,
		accumulators: accumulators,
	}
}

func (t *scalarTable) aggregate(vector model.StepVector) {
	t.reset()
	for i := range vector.Samples {
		t.addSample(vector.T, vector.SampleIDs[i], vector.Samples[i])
	}
}

func (t *scalarTable) addSample(ts int64, sampleID uint64, sample float64) {
	outputSampleID := t.inputs[sampleID]
	output := t.outputs[outputSampleID]

	t.timestamp = ts
	t.accumulators[output.ID].AddFunc(sample)
}

func (t *scalarTable) reset() {
	for i := range t.outputs {
		t.accumulators[i].Reset()
	}
}

func (t *scalarTable) toVector(pool *model.VectorPool) model.StepVector {
	result := pool.GetStepVector(t.timestamp)
	for i, v := range t.outputs {
		if t.accumulators[i].HasValue() {
			result.SampleIDs = append(result.SampleIDs, v.ID)
			result.Samples = append(result.Samples, t.accumulators[i].ValueFunc())
		}
	}
	return result
}

func (t *scalarTable) size() int {
	return len(t.outputs)
}

func hashMetric(metric labels.Labels, without bool, grouping []string, buf []byte) (uint64, string, labels.Labels) {
	buf = buf[:0]
	if without {
		lb := labels.NewBuilder(metric)
		lb.Del(grouping...)
		key, bytes := metric.HashWithoutLabels(buf, grouping...)
		return key, string(bytes), lb.Labels()
	}

	if len(grouping) == 0 {
		return 0, "", labels.Labels{}
	}

	lb := labels.NewBuilder(metric)
	lb.Keep(grouping...)
	key, bytes := metric.HashForLabels(buf, grouping...)
	return key, string(bytes), lb.Labels()
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
