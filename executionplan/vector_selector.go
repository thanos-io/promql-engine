package executionplan

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/fpetkovski/promql-engine/model"

	"github.com/fpetkovski/promql-engine/points"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/value"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"
)

type vectorScan struct {
	labels    labels.Labels
	signature string
	samples   *storage.MemoizedSeriesIterator
}

type vectorSelector struct {
	storage storage.Queryable
	series  []vectorScan
	once    sync.Once
	pool    *points.Pool

	matchers []*labels.Matcher
	hints    *storage.SelectHints

	mint        int64
	maxt        int64
	step        int64
	currentStep int64
}

func NewVectorSelector(pool *points.Pool, storage storage.Queryable, matchers []*labels.Matcher, hints *storage.SelectHints, mint, maxt time.Time, step time.Duration) VectorOperator {
	// TODO(fpetkovski): Add offset parameter.
	return &vectorSelector{
		storage: storage,
		pool:    pool,

		matchers: matchers,
		hints:    hints,

		mint:        mint.UnixMilli(),
		maxt:        maxt.UnixMilli(),
		step:        step.Milliseconds(),
		currentStep: mint.UnixMilli() - step.Milliseconds(),
	}
}

func (o *vectorSelector) Next(ctx context.Context) ([]model.Vector, error) {
	if o.currentStep > o.maxt {
		return nil, nil
	}

	var err error
	o.once.Do(func() { err = o.initializeSeries(ctx) })
	if err != nil {
		return nil, err
	}

	numSteps := 30
	vectors := make([]model.Vector, 0, numSteps)
	ts := o.currentStep
	for i := 0; i < len(o.series); i++ {
		var (
			series   = o.series[i]
			seriesTs = ts
		)

		for currStep := 0; currStep < numSteps && seriesTs <= o.maxt; currStep++ {
			_, v, ok := selectPoint(series.samples, seriesTs)
			if ok {
				if len(vectors) <= currStep {
					vectors = append(vectors, make(model.Vector, 0))
				}
				vectors[currStep] = append(vectors[currStep], model.Sample{
					Signature: series.signature,
					Sample: promql.Sample{
						Metric: series.labels,
						Point:  promql.Point{T: seriesTs, V: v},
					},
				})
			} else {
				break
			}
			seriesTs += o.step
		}
	}
	o.currentStep += o.step * int64(numSteps)

	return vectors, nil
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
	i := 0
	for seriesSet.Next() {
		s := seriesSet.At()
		series = append(series, vectorScan{
			signature: strconv.Itoa(i),
			labels:    s.Labels(),
			samples:   storage.NewMemoizedIterator(s.Iterator(), 5*time.Minute.Milliseconds()),
		})
		i++
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
