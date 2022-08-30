package plan

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

type selector struct {
	storage storage.Storage
	points  *points

	matchers []*labels.Matcher
	hints    *storage.SelectHints

	mint        int64
	maxt        int64
	step        int64
	selectRange int64
}

func NewSelector(
	storage storage.Storage,
	matchers []*labels.Matcher,
	hints *storage.SelectHints,
	mint, maxt time.Time,
	step, selectRange time.Duration,
) Operator {
	return &selector{
		storage: storage,
		points:  newPointPool(),

		matchers:    matchers,
		hints:       hints,
		mint:        mint.UnixMilli(),
		maxt:        maxt.UnixMilli(),
		step:        step.Milliseconds(),
		selectRange: selectRange.Milliseconds(),
	}
}

func (o *selector) Next(ctx context.Context) (<-chan promql.Matrix, error) {
	querier, err := o.storage.Querier(ctx, o.mint, o.maxt)
	if err != nil {
		return nil, err
	}

	numSteps := (o.maxt-o.mint)/o.step + 1
	out := make(chan promql.Matrix, numSteps)
	go func() {
		defer querier.Close()
		defer close(out)

		series := make([]storageSeries, 0)
		seriesSet := querier.Select(true, o.hints, o.matchers...)
		for seriesSet.Next() {
			s := seriesSet.At()
			series = append(series, storageSeries{
				labels:  s.Labels(),
				samples: s.Iterator(),
			})
		}

		matrix := make(promql.Matrix, len(series))
		for i := 0; i < len(series); i++ {
			matrix[i] = promql.Series{
				Metric: series[i].labels,
				Points: make([]promql.Point, 0, numSteps),
			}
		}
		for stepTs := o.mint; stepTs <= o.maxt; stepTs += o.step {
			for i, s := range series {
				resultSeries := &matrix[i]

				maxt := stepTs
				mint := stepTs - o.selectRange
				if mint < 0 {
					mint = 0
				}

				if ok := s.samples.Seek(mint); !ok {
					continue
				}
				for {
					t, v := s.samples.At()
					if t > maxt {
						break
					}
					resultSeries.Points = append(resultSeries.Points, promql.Point{
						T: t,
						V: v,
					})
					if !s.samples.Next() {
						break
					}
				}
			}
			if len(matrix) > 0 {
				out <- matrix
			}
		}
	}()

	return out, nil
}
