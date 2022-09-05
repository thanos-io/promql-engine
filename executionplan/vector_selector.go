package executionplan

import (
	"context"
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
	storage storage.Queryable
	series  []vectorScan

	matchers []*labels.Matcher
	hints    *storage.SelectHints

	mint        int64
	maxt        int64
	step        int64
	currentStep int64
}

func NewVectorSelector(storage storage.Queryable, matchers []*labels.Matcher, hints *storage.SelectHints, mint, maxt time.Time, step time.Duration) VectorOperator {
	// TODO(fpetkovski): Add offset parameter.
	return &vectorSelector{
		storage: storage,

		matchers: matchers,
		hints:    hints,

		mint:        mint.UnixMilli(),
		maxt:        maxt.UnixMilli(),
		step:        step.Milliseconds(),
		currentStep: mint.UnixMilli() - step.Milliseconds(),
	}
}

func (o *vectorSelector) Next(ctx context.Context) (promql.Vector, error) {
	o.currentStep += o.step
	if o.currentStep > o.maxt {
		return nil, nil
	}

	if o.series == nil {
		if err := o.initializeSeries(ctx); err != nil {
			return nil, err
		}
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
	mint := o.mint - 5*time.Minute.Milliseconds()
	querier, err := o.storage.Querier(ctx, mint, o.maxt)
	if err != nil {
		return err
	}
	defer querier.Close()

	series := make([]vectorScan, 0)
	seriesSet := querier.Select(true, o.hints, o.matchers...)
	for seriesSet.Next() {
		s := seriesSet.At()

		series = append(series, vectorScan{
			labels:  s.Labels(),
			samples: storage.NewMemoizedIterator(s.Iterator(), 5*time.Minute.Milliseconds()),
		})
	}
	o.series = series

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
