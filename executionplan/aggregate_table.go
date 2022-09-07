package executionplan

import (
	"github.com/fpetkovski/promql-engine/model"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
)

type groupingKey struct {
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

func newAggregateTable(g groupingKeyFunc, f newAccumulatorFunc) *aggregateTable {
	return &aggregateTable{
		groupingKeys:        make([]groupingKey, 100000),
		table:               make(map[uint64]*aggregateResult),
		makeAccumulatorFunc: f,
		groupingKeyFunc:     g,
	}
}

func (t *aggregateTable) addSample(sample model.Sample) {
	var (
		key      uint64
		sampleID uint64
		lbls     labels.Labels
	)

	cachedResult := t.groupingKeys[sample.ID]
	if cachedResult.labels != nil {
		key = cachedResult.hash
		lbls = cachedResult.labels
		sampleID = cachedResult.sampleID
	} else {
		key, _, lbls = t.groupingKeyFunc(sample.Metric)
		sampleID = key
		t.groupingKeys[sample.ID] = groupingKey{
			hash:     key,
			labels:   lbls,
			sampleID: sampleID,
		}
	}

	if _, ok := t.table[key]; !ok {
		t.table[key] = &aggregateResult{
			sampleID:    sampleID,
			metric:      lbls,
			timestamp:   sample.T,
			accumulator: t.makeAccumulatorFunc(),
		}
	}

	t.table[key].timestamp = sample.T
	t.table[key].accumulator.AddFunc(sample.Point)
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

func (t *aggregateTable) toVector() model.Vector {
	result := make(model.Vector, 0, len(t.table))
	for _, v := range t.table {
		result = append(result, model.Sample{
			Sample: promql.Sample{
				Metric: v.metric,
				Point: promql.Point{
					T: v.timestamp,
					V: v.accumulator.ValueFunc(),
				},
			},
			ID: v.sampleID,
		})
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
