package executionplan

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/value"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"
)

type vectorScan struct {
	labels  labels.Labels
	samples *storage.MemoizedSeriesIterator
}

type vectorSelector struct {
	once     sync.Once
	selector *seriesSelector
	series   []vectorScan

	hints *storage.SelectHints

	mint        int64
	maxt        int64
	step        int64
	currentStep int64

	shard     int
	numShards int
}

func NewVectorSelector(
	selector *seriesSelector,
	hints *storage.SelectHints,
	mint, maxt time.Time,
	step time.Duration,
	shard, numShards int,
) VectorOperator {
	return &vectorSelector{
		selector: selector,

		hints: hints,

		// TODO(fpetkovski): Add offset parameter.
		mint:        mint.UnixMilli(),
		maxt:        maxt.UnixMilli(),
		step:        step.Milliseconds(),
		currentStep: mint.UnixMilli() - step.Milliseconds(),

		shard:     shard,
		numShards: numShards,
	}
}

func (o *vectorSelector) Next(ctx context.Context) (promql.Vector, error) {
	o.currentStep += o.step
	if o.currentStep > o.maxt {
		return nil, nil
	}

	var err error
	o.once.Do(func() { err = o.initializeSeries(ctx) })
	if err != nil {
		return nil, err
	}

	vector := make(promql.Vector, len(o.series))
	for i := 0; i < len(o.series); i++ {
		s := o.series[i]
		vector[i].Metric = s.labels
		_, v, ok := selectPoint(s.samples, o.currentStep)
		if ok {
			vector[i].V = v
			vector[i].T = o.currentStep
		}

	}
	return vector, nil
}

func (o *vectorSelector) initializeSeries(ctx context.Context) error {
	series, err := o.selector.Series(ctx, o.shard, o.numShards)
	if err != nil {
		return err
	}

	scanners := make([]vectorScan, 0)
	for _, s := range series {
		scanners = append(scanners, vectorScan{
			labels:  s.Labels(),
			samples: storage.NewMemoizedIterator(s.Iterator(), 5*time.Minute.Milliseconds()),
		})
	}
	o.series = scanners

	return nil
}

// TODO(fpetkovski): Add error handling and max samples limit.
func selectPoint(it *storage.MemoizedSeriesIterator, ts int64) (int64, float64, bool) {
	lookbackDelta := 5 * time.Minute.Milliseconds()
	refTime := ts
	var t int64
	var v float64

	ok := it.Seek(refTime)
	if ok {
		t, v = it.At()
	}

	if !ok || t > refTime {
		t, v, ok = it.PeekPrev()
		if !ok || t < refTime-lookbackDelta {
			return 0, 0, false
		}
	}
	if value.IsStaleNaN(v) {
		return 0, 0, false
	}
	return t, v, true
}
