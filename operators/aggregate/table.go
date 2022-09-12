package aggregate

import (
	"github.com/fpetkovski/promql-engine/operators/model"
	"github.com/prometheus/prometheus/model/labels"
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
