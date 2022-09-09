package executionplan

import (
	"sync"

	"github.com/fpetkovski/promql-engine/model"
	"github.com/prometheus/prometheus/model/labels"
)

type groupingKey struct {
	once     *sync.Once
	hash     uint64
	sampleID uint64
	labels   labels.Labels
}

type aggregateResult struct {
	metric      labels.Labels
	sampleID    uint64
	timestamp   int64
	accumulator *accumulator
}

type aggregateTable struct {
	// hashKeyCache is a map from series index to the cache key for the series.
	groupingKeys        []groupingKey
	table               map[uint64]*aggregateResult
	makeAccumulatorFunc newAccumulatorFunc
	groupingKeyFunc     groupingKeyFunc
}

func newAggregateTable(g groupingKeyFunc, f newAccumulatorFunc, groupingKeys []groupingKey) *aggregateTable {
	return &aggregateTable{
		groupingKeys:        groupingKeys,
		table:               make(map[uint64]*aggregateResult),
		makeAccumulatorFunc: f,
		groupingKeyFunc:     g,
	}
}

func (t *aggregateTable) addSample(ts int64, sample model.StepSample) {
	var (
		key      uint64
		sampleID uint64
		lbls     labels.Labels
	)

	once := t.groupingKeys[sample.ID].once
	once.Do(func() {
		key, _, lbls = t.groupingKeyFunc(sample.Metric)
		sampleID = key
		t.groupingKeys[sample.ID] = groupingKey{
			hash:     key,
			labels:   lbls,
			sampleID: sampleID,
			once:     once,
		}
	})

	cachedResult := t.groupingKeys[sample.ID]
	key = cachedResult.hash
	lbls = cachedResult.labels
	sampleID = cachedResult.sampleID

	if _, ok := t.table[key]; !ok {
		t.table[key] = &aggregateResult{
			sampleID:    sampleID,
			metric:      lbls,
			timestamp:   ts,
			accumulator: t.makeAccumulatorFunc(),
		}
	}

	t.table[key].timestamp = ts
	t.table[key].accumulator.AddFunc(sample.V)
}

func (t *aggregateTable) reset() {
	for k, v := range t.table {
		t.table[k] = &aggregateResult{
			sampleID:    v.sampleID,
			metric:      v.metric,
			accumulator: t.makeAccumulatorFunc(),
		}
	}
}

func (t *aggregateTable) toVector(pool *model.VectorPool) model.StepVector {
	result := model.StepVector{
		Samples: pool.GetSamples(),
	}
	for _, v := range t.table {
		if v.accumulator.HasValue() {
			result.T = v.timestamp
			result.Samples = append(result.Samples, model.StepSample{
				Metric: v.metric,
				V:      v.accumulator.ValueFunc(),
				ID:     v.sampleID,
			})
		}
	}
	return result
}

type groupingKeyFunc func(metric labels.Labels) (uint64, string, labels.Labels)

// groupingKey builds and returns the grouping key and the
// resulting labels value pairs for the given metric and grouping labels.
func newGroupingKeyGenerator(grouping []string, without bool, buf []byte) groupingKeyFunc {
	return func(metric labels.Labels) (uint64, string, labels.Labels) {
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
}
