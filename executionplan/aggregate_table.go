package executionplan

import (
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
)

type groupingKey struct {
	hash   uint64
	labels labels.Labels
}

type aggregateResult struct {
	metric      labels.Labels
	timestamp   int64
	accumulator *accumulator
}

type aggregateTable struct {
	// hashKeyCache is a map from series index to the cache key for the series.
	hashKeyCache        map[int]groupingKey
	table               map[uint64]*aggregateResult
	makeAccumulatorFunc newAccumulatorFunc
	groupingKeyFunc     groupingKeyFunc
}

func newAggregateTable(g groupingKeyFunc, f newAccumulatorFunc) *aggregateTable {
	return &aggregateTable{
		hashKeyCache:        make(map[int]groupingKey),
		table:               make(map[uint64]*aggregateResult),
		makeAccumulatorFunc: f,
		groupingKeyFunc:     g,
	}
}

func (t *aggregateTable) addSample(seriesID int, sample promql.Sample) {
	var (
		key  uint64
		lbls labels.Labels
	)
	if cachedResult, ok := t.hashKeyCache[seriesID]; ok {
		key = cachedResult.hash
		lbls = cachedResult.labels
	} else {
		key, lbls = t.groupingKeyFunc(sample.Metric)
		t.hashKeyCache[seriesID] = groupingKey{
			hash:   key,
			labels: lbls,
		}
	}

	if _, ok := t.table[key]; !ok {
		t.table[key] = &aggregateResult{
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
			metric:      v.metric,
			accumulator: t.makeAccumulatorFunc(),
		}
	}
}

type groupingKeyFunc func(metric labels.Labels) (uint64, labels.Labels)

// groupingKey builds and returns the grouping key and series labels
// for the given metric and grouping labels.
func newGroupingKeyGenerator(grouping []string, without bool, buf []byte) groupingKeyFunc {
	return func(metric labels.Labels) (uint64, labels.Labels) {
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
}
