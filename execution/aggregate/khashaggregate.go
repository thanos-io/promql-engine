package aggregate

import (
	"container/heap"
	"context"
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/thanos-community/promql-engine/execution/model"
	"golang.org/x/exp/slices"
)

type kAggregate struct {
	next    model.VectorOperator
	paramOp model.VectorOperator

	vectorPool *model.VectorPool

	by          bool
	labels      []string
	aggregation parser.ItemType

	once          sync.Once
	series        []labels.Labels
	seriesAggHash map[uint64]uint64
	compare       func(float64, float64) bool
}

func NewKHashAggregate(
	points *model.VectorPool,
	next model.VectorOperator,
	paramOp model.VectorOperator,
	aggregation parser.ItemType,
	by bool,
	labels []string,
) (model.VectorOperator, error) {
	var compare func(float64, float64) bool

	if aggregation == parser.TOPK {
		compare = func(f float64, s float64) bool {
			return f < s
		}
	} else {
		compare = func(f float64, s float64) bool {
			return s < f
		}
	}
	// Grouping labels need to be sorted in order for metric hashing to work.
	// https://github.com/prometheus/prometheus/blob/8ed39fdab1ead382a354e45ded999eb3610f8d5f/model/labels/labels.go#L162-L181
	slices.Sort(labels)

	a := &kAggregate{
		next:          next,
		vectorPool:    points,
		by:            by,
		aggregation:   aggregation,
		labels:        labels,
		paramOp:       paramOp,
		seriesAggHash: map[uint64]uint64{},
		compare:       compare,
	}

	return a, nil
}

func (a *kAggregate) Next(ctx context.Context) ([]model.StepVector, error) {
	in, err := a.next.Next(ctx)
	if err != nil {
		return nil, err
	}
	if in == nil {
		return nil, nil
	}

	defer a.next.GetPool().PutVectors(in)

	var arg float64

	if a.paramOp != nil {
		p, err := a.paramOp.Next(ctx)
		if err != nil {
			return nil, err
		}
		arg = p[0].Samples[0]
	}

	a.once.Do(func() { err = a.init(ctx) })
	if err != nil {
		return nil, err
	}

	result := a.vectorPool.GetVectorBatch()
	for _, vector := range in {
		aggregatedVectors := a.byAggregation(vector)
		for _, av := range aggregatedVectors {
			if len(av.SampleIDs) == 0 {
				continue
			}
			s := a.aggregate(vector.T, int(arg), av.SampleIDs, av.Samples)
			result = append(result, s)
		}
		a.next.GetPool().PutStepVector(vector)
	}

	return result, nil
}

func (a *kAggregate) Series(ctx context.Context) ([]labels.Labels, error) {
	var err error
	a.once.Do(func() { err = a.init(ctx) })
	if err != nil {
		return nil, err
	}

	return a.series, nil
}

func (a *kAggregate) GetPool() *model.VectorPool {
	return a.vectorPool
}

func (a *kAggregate) Explain() (me string, next []model.VectorOperator) {
	if a.by {
		return fmt.Sprintf("[*kaggregate] %v by (%v)", a.aggregation.String(), a.labels), []model.VectorOperator{a.paramOp, a.next}
	}
	return fmt.Sprintf("[*kaggregate] %v without (%v)", a.aggregation.String(), a.labels), []model.VectorOperator{a.paramOp, a.next}
}

func (a *kAggregate) init(ctx context.Context) error {
	series, err := a.next.Series(ctx)
	if err != nil {
		return err
	}
	aggHash := make(map[uint64]uint64)
	buf := make([]byte, 1024)
	for i := 0; i < len(series); i++ {
		hash, _, _ := hashMetric(series[i], !a.by, a.labels, buf)
		if _, ok := aggHash[hash]; !ok {
			aggHash[hash] = uint64(len(aggHash))
		}
		a.seriesAggHash[uint64(i)] = aggHash[hash]
	}
	a.vectorPool.SetStepSize(len(series))
	a.series = series
	return nil
}

func (a *kAggregate) byAggregation(v model.StepVector) []model.StepVector {
	r := make([]model.StepVector, len(a.seriesAggHash))
	for i, id := range v.SampleIDs {
		r[a.seriesAggHash[id]].SampleIDs = append(r[a.seriesAggHash[id]].SampleIDs, id)
		r[a.seriesAggHash[id]].Samples = append(r[a.seriesAggHash[id]].Samples, v.Samples[i])
	}
	return r
}

func (a *kAggregate) aggregate(t int64, k int, SampleIDs []uint64, samples []float64) model.StepVector {
	result := a.vectorPool.GetStepVector(t)
	h := samplesHeap{compare: a.compare}
	for i, d := range SampleIDs {
		e := entry{sId: d, total: samples[i]}
		if h.Len() < k || h.compare(h.entries[0].total, e.total) || math.IsNaN(h.entries[0].total) {
			if k == 1 && h.Len() == 1 {
				h.entries[0].sId = e.sId
				h.entries[0].total = e.total
				continue
			}

			if h.Len() == k {
				heap.Pop(&h)
			}

			heap.Push(&h, &e)
		}
	}

	// The heap keeps the lowest value on top, so reverse it.
	if len(h.entries) > 1 {
		sort.Sort(sort.Reverse(h))
	}

	for _, e := range h.entries {
		result.SampleIDs = append(result.SampleIDs, e.sId)
		result.Samples = append(result.Samples, e.total)
	}
	return result
}

type entry struct {
	sId   uint64
	total float64
}

type samplesHeap struct {
	entries []entry
	compare func(float64, float64) bool
}

func (s samplesHeap) Len() int {
	return len(s.entries)
}

func (s samplesHeap) Less(i, j int) bool {
	if math.IsNaN(s.entries[i].total) {
		return true
	}
	return s.compare(s.entries[i].total, s.entries[j].total)
}

func (s samplesHeap) Swap(i, j int) {
	s.entries[i], s.entries[j] = s.entries[j], s.entries[i]
}

func (s *samplesHeap) Push(x interface{}) {
	s.entries = append(s.entries, *(x.(*entry)))
}

func (s *samplesHeap) Pop() interface{} {
	old := (*s).entries
	n := len(old)
	el := old[n-1]
	(*s).entries = old[0 : n-1]
	return el
}
