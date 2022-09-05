package executionplan

import (
	"context"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/value"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"
)

type matrixScan struct {
	labels         labels.Labels
	previousPoints []promql.Point
	samples        *storage.BufferedSeriesIterator
}

type matrixSelector struct {
	call    FunctionCall
	storage storage.Queryable
	series  []matrixScan

	matchers []*labels.Matcher
	hints    *storage.SelectHints

	mint        int64
	maxt        int64
	step        int64
	selectRange int64
	currentStep int64
}

func NewMatrixSelector(storage storage.Queryable, matchers []*labels.Matcher, hints *storage.SelectHints, mint, maxt time.Time, step, selectRange time.Duration) VectorOperator {
	// TODO(fpetkovski): Add offset parameter.
	return &matrixSelector{
		storage: storage,
		call:    NewRate(selectRange),

		matchers:    matchers,
		hints:       hints,
		mint:        mint.UnixMilli(),
		maxt:        maxt.UnixMilli(),
		step:        step.Milliseconds(),
		selectRange: selectRange.Milliseconds(),
		currentStep: mint.UnixMilli() - step.Milliseconds(),
	}
}

func (o *matrixSelector) Next(ctx context.Context) (promql.Vector, error) {
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
		s := &o.series[i]
		vector[i].Metric = s.labels

		maxt := o.currentStep
		mint := maxt - o.selectRange

		rangePoints := selectPoints(s.samples, mint, maxt, o.series[i].previousPoints)
		result := o.call(rangePoints, time.UnixMilli(o.currentStep))
		if result != nil {
			vector[i].Point = *result
			o.series[i].previousPoints = rangePoints
		} else {
			vector[i].Point.T = -1
			continue
		}

		// Only buffer stepRange milliseconds from the second step on.
		stepRange := o.selectRange
		if stepRange > o.currentStep {
			stepRange = o.currentStep
		}
		s.samples.ReduceDelta(stepRange)
	}

	return vector, nil
}

func (o *matrixSelector) initializeSeries(ctx context.Context) error {
	mint := o.mint - 5*time.Minute.Milliseconds()
	querier, err := o.storage.Querier(ctx, mint, o.maxt)
	if err != nil {
		return err
	}
	defer querier.Close()

	series := make([]matrixScan, 0)
	seriesSet := querier.Select(true, o.hints, o.matchers...)
	for seriesSet.Next() {
		s := seriesSet.At()

		series = append(series, matrixScan{
			labels:         s.Labels(),
			previousPoints: make([]promql.Point, 0),
			samples:        storage.NewBufferIterator(s.Iterator(), o.selectRange),
		})
	}
	o.series = series

	return nil
}

// matrixIterSlice populates a matrix vector covering the requested range for a
// single time series, with points retrieved from an iterator.
//
// As an optimization, the matrix vector may already contain points of the same
// time series from the evaluation of an earlier step (with lower mint and maxt
// values). Any such points falling before mint are discarded; points that fall
// into the [mint, maxt] range are retained; only points with later timestamps
// are populated from the iterator.
// TODO(fpetkovski): Add error handling and max samples limit.
func selectPoints(it *storage.BufferedSeriesIterator, mint, maxt int64, out []promql.Point) []promql.Point {
	if len(out) > 0 && out[len(out)-1].T >= mint {
		// There is an overlap between previous and current ranges, retain common
		// points. In most such cases:
		//   (a) the overlap is significantly larger than the eval step; and/or
		//   (b) the number of samples is relatively small.
		// so a linear search will be as fast as a binary search.
		var drop int
		for drop = 0; out[drop].T < mint; drop++ {
		}
		copy(out, out[drop:])
		out = out[:len(out)-drop]
		// Only append points with timestamps after the last timestamp we have.
		mint = out[len(out)-1].T + 1
	} else {
		out = out[:0]
	}

	ok := it.Seek(maxt)
	buf := it.Buffer()
	for buf.Next() {
		t, v := buf.At()
		if value.IsStaleNaN(v) {
			continue
		}
		// Values in the buffer are guaranteed to be smaller than maxt.
		if t >= mint {
			out = append(out, promql.Point{T: t, V: v})
		}
	}
	// The seeked sample might also be in the range.
	if ok {
		t, v := it.At()
		if t == maxt && !value.IsStaleNaN(v) {
			out = append(out, promql.Point{T: t, V: v})
		}
	}
	return out
}
