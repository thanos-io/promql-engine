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
	storage storage.Storage

	matchers []*labels.Matcher
	hints    *storage.SelectHints

	mint int64
	maxt int64
	step int64
}

func NewVectorSelector(storage storage.Storage, matchers []*labels.Matcher, hints *storage.SelectHints, mint, maxt time.Time, step time.Duration) VectorOperator {
	// TODO(fpetkovski): Add offset parameter.
	return &vectorSelector{
		storage: storage,

		matchers: matchers,
		hints:    hints,
		mint:     mint.UnixMilli(),
		maxt:     maxt.UnixMilli(),
		step:     step.Milliseconds(),
	}
}

func (o *vectorSelector) Next(ctx context.Context) (<-chan promql.Vector, error) {
	querier, err := o.storage.Querier(ctx, o.mint, o.maxt)
	if err != nil {
		return nil, err
	}

	numSteps := (o.maxt-o.mint)/o.step + 1
	out := make(chan promql.Vector, numSteps)
	go func() {
		defer querier.Close()
		defer close(out)

		series := make([]vectorScan, 0)
		seriesSet := querier.Select(true, o.hints, o.matchers...)
		for seriesSet.Next() {
			s := seriesSet.At()

			series = append(series, vectorScan{
				labels:  s.Labels(),
				samples: storage.NewMemoizedIterator(s.Iterator(), 5*time.Minute.Milliseconds()),
			})
		}

		for stepTs := o.mint; stepTs <= o.maxt; stepTs += o.step {
			vector := make(promql.Vector, len(series))

			for i := 0; i < len(series); i++ {
				s := series[i]
				vector[i].Metric = s.labels
				_, v, ok := selectPoint(s.samples, stepTs)
				if ok {
					vector[i].V = v
					vector[i].T = stepTs
				}
			}

			out <- vector
		}
	}()

	return out, nil
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
