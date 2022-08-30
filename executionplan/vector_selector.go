package executionplan

import (
	"context"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	"time"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"
)

type storageSeries struct {
	labels  labels.Labels
	samples chunkenc.Iterator
}

type vectorSelector struct {
	storage storage.Storage
	points  *points

	matchers []*labels.Matcher
	hints    *storage.SelectHints

	mint        int64
	maxt        int64
	step        int64
	selectRange int64
}

func NewSelector(storage storage.Storage, matchers []*labels.Matcher, hints *storage.SelectHints, mint, maxt time.Time, step time.Duration) VectorOperator {
	return &vectorSelector{
		storage: storage,
		points:  newPointPool(),

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

		series := make([]storageSeries, 0)
		seriesSet := querier.Select(true, o.hints, o.matchers...)
		for seriesSet.Next() {
			s := seriesSet.At()
			it := s.Iterator()
			it.Seek(o.mint)
			series = append(series, storageSeries{
				labels:  s.Labels(),
				samples: it,
			})
		}

		for stepTs := o.mint; stepTs <= o.maxt; stepTs += o.step {
			vector := make(promql.Vector, len(series))

			for i := 0; i < len(series); i++ {
				s := series[i]
				vector[i].Metric = s.labels
				for {
					t, v := s.samples.At()
					vector[i].T = t
					vector[i].V = v
					if t >= stepTs {
						break
					}
					if !s.samples.Next() {
						break
					}
				}
			}

			out <- vector
		}
	}()

	return out, nil
}
