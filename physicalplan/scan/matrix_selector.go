package scan

import (
	"context"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/fpetkovski/promql-engine/physicalplan/model"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/value"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"
)

type matrixScanner struct {
	labels         labels.Labels
	signature      uint64
	previousPoints []promql.Point
	samples        *storage.BufferedSeriesIterator
}

type matrixSelector struct {
	call     FunctionCall
	selector *seriesSelector
	scanners []matrixScanner
	series   []labels.Labels
	once     sync.Once

	matchers   []*labels.Matcher
	hints      *storage.SelectHints
	vectorPool *model.VectorPool

	mint        int64
	maxt        int64
	step        int64
	selectRange int64
	currentStep int64
	stepsBatch  int

	shard     int
	numShards int
}

func NewMatrixSelector(
	pool *model.VectorPool,
	selector *seriesSelector,
	call FunctionCall,
	mint, maxt time.Time,
	stepsBatch int,
	step, selectRange time.Duration,
	shard, numShard int,
) model.VectorOperator {
	// TODO(fpetkovski): Add offset parameter.
	return &matrixSelector{
		selector:   selector,
		call:       call,
		vectorPool: pool,

		mint:       mint.UnixMilli(),
		maxt:       maxt.UnixMilli(),
		step:       step.Milliseconds(),
		stepsBatch: stepsBatch,

		selectRange: selectRange.Milliseconds(),
		currentStep: mint.UnixMilli(),

		shard:     shard,
		numShards: numShard,
	}
}

func (o *matrixSelector) Series(ctx context.Context) ([]labels.Labels, error) {
	if err := o.loadSeries(ctx); err != nil {
		return nil, err
	}
	return o.series, nil
}

func (o *matrixSelector) GetPool() *model.VectorPool {
	return o.vectorPool
}

func (o *matrixSelector) Next(ctx context.Context) ([]model.StepVector, error) {
	if o.currentStep > o.maxt {
		return nil, nil
	}

	if err := o.loadSeries(ctx); err != nil {
		return nil, err
	}

	totalSteps := (o.maxt+o.mint)/o.step + 1
	numSteps := int(math.Min(float64(o.stepsBatch), float64(totalSteps)))

	vectors := o.vectorPool.GetVectors()
	ts := o.currentStep
	for i := 0; i < len(o.scanners); i++ {
		var (
			series   = o.scanners[i]
			seriesTs = ts
		)

		for currStep := 0; currStep < numSteps && seriesTs <= o.maxt; currStep++ {
			if len(vectors) <= currStep {
				vectors = append(vectors, model.StepVector{
					T:       seriesTs,
					Samples: o.vectorPool.GetSamples(),
				})
			}
			maxt := seriesTs
			mint := maxt - o.selectRange

			rangePoints := selectPoints(series.samples, mint, maxt, o.scanners[i].previousPoints)
			result := o.call(series.labels, rangePoints, time.UnixMilli(seriesTs))
			if result.T >= 0 {
				vectors[currStep].T = result.T
				vectors[currStep].Samples = append(vectors[currStep].Samples, model.StepSample{
					ID:     series.signature,
					Metric: series.labels,
					V:      result.V,
				})
			}
			o.scanners[i].previousPoints = rangePoints

			// Only buffer stepRange milliseconds from the second step on.
			stepRange := o.selectRange
			if stepRange > o.step {
				stepRange = o.step
			}
			series.samples.ReduceDelta(stepRange)

			seriesTs += o.step
		}
	}
	o.currentStep += o.step * int64(numSteps)

	return vectors, nil

}

func (o *matrixSelector) loadSeries(ctx context.Context) error {
	var err error
	o.once.Do(func() {
		series, loadErr := o.selector.getSeries(ctx, o.shard, o.numShards)
		if loadErr != nil {
			err = loadErr
			return
		}

		o.scanners = make([]matrixScanner, len(series))
		o.series = make([]labels.Labels, len(series))
		for i, s := range series {
			lbls := dropMetricName(s.Labels())
			sort.Sort(lbls)

			o.scanners[i] = matrixScanner{
				labels:    lbls,
				signature: s.signature,
				samples:   storage.NewBufferIterator(s.Iterator(), o.selectRange),
			}
			o.series[i] = lbls
		}
	})
	return err
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
