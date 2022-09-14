package aggregate

import (
	"fmt"
	"math"

	"github.com/fpetkovski/promql-engine/physicalplan/model"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

type aggregateResult struct {
	metric   labels.Labels
	sampleID uint64
}

type aggregateTable struct {
	timestamp    int64
	inputs       []uint64
	outputs      []*aggregateResult
	accumulators []*accumulator
}

func newAggregateTable(inputSampleIDs []uint64, outputs []*aggregateResult, makeAccumulator newAccumulatorFunc) *aggregateTable {
	accumulators := make([]*accumulator, len(outputs))
	for i := 0; i < len(outputs); i++ {
		accumulators[i] = makeAccumulator()
	}
	return &aggregateTable{
		inputs:       inputSampleIDs,
		outputs:      outputs,
		accumulators: accumulators,
	}
}

func (t *aggregateTable) addSample(ts int64, sample model.StepSample) {
	outputSampleID := t.inputs[sample.ID]
	output := t.outputs[outputSampleID]

	t.timestamp = ts
	t.accumulators[output.sampleID].AddFunc(sample.V)
}

func (t *aggregateTable) reset() {
	for i := range t.outputs {
		t.accumulators[i].Reset()
	}
}

func (t *aggregateTable) toVector(pool *model.VectorPool) model.StepVector {
	result := model.StepVector{
		Samples: pool.GetSamples(),
	}

	for i, v := range t.outputs {
		if t.accumulators[i].HasValue() {
			result.T = t.timestamp
			result.Samples = append(result.Samples, model.StepSample{
				Metric: v.metric,
				V:      t.accumulators[i].ValueFunc(),
				ID:     v.sampleID,
			})
		}
	}
	return result
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
