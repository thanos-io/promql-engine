// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package scan

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/value"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"

	"github.com/thanos-community/promql-engine/execution/function"
	"github.com/thanos-community/promql-engine/execution/model"
	engstore "github.com/thanos-community/promql-engine/execution/storage"
	"github.com/thanos-community/promql-engine/query"
)

type matrixScanner struct {
	labels         labels.Labels
	signature      uint64
	previousPoints []promql.Point
	samples        *storage.BufferedSeriesIterator
}

type matrixSelector struct {
	storage  engstore.SeriesSelector
	call     function.FunctionCall
	scanners []matrixScanner
	series   []labels.Labels
	once     sync.Once

	vectorPool *model.VectorPool

	numSteps    int
	mint        int64
	maxt        int64
	step        int64
	selectRange int64
	offset      int64
	currentStep int64

	shard     int
	numShards int
}

// NewMatrixSelector creates operator which selects vector of series over time.
func NewMatrixSelector(
	pool *model.VectorPool,
	selector engstore.SeriesSelector,
	call function.FunctionCall,
	opts *query.Options,
	selectRange, offset time.Duration,
	shard, numShard int,
) model.VectorOperator {
	// TODO(fpetkovski): Add offset parameter.
	return &matrixSelector{
		storage:    selector,
		call:       call,
		vectorPool: pool,

		numSteps: opts.NumSteps(),
		mint:     opts.Start.UnixMilli(),
		maxt:     opts.End.UnixMilli(),
		step:     opts.Step.Milliseconds(),

		selectRange: selectRange.Milliseconds(),
		offset:      offset.Milliseconds(),
		currentStep: opts.Start.UnixMilli(),

		shard:     shard,
		numShards: numShard,
	}
}

func (o *matrixSelector) Explain() (me string, next []model.VectorOperator) {
	r := time.Duration(o.selectRange) * time.Millisecond
	return fmt.Sprintf("[*matrixSelector] {%v}[%s] %v mod %v", o.storage.Matchers(), r, o.shard, o.numShards), nil
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

	vectors := o.vectorPool.GetVectorBatch()
	ts := o.currentStep
	for i := 0; i < len(o.scanners); i++ {
		var (
			series   = o.scanners[i]
			seriesTs = ts
		)

		for currStep := 0; currStep < o.numSteps && seriesTs <= o.maxt; currStep++ {
			if len(vectors) <= currStep {
				vectors = append(vectors, o.vectorPool.GetStepVector(seriesTs))
			}
			maxt := seriesTs - o.offset
			mint := maxt - o.selectRange
			rangePoints := selectPoints(series.samples, mint, maxt, o.scanners[i].previousPoints)

			// TODO(saswatamcode): Handle multi-arg functions for matrixSelectors via injectable.
			result := o.call(function.FunctionArgs{
				Labels:      series.labels,
				Points:      rangePoints,
				StepTime:    seriesTs,
				SelectRange: o.selectRange,
			})

			if result.Point != function.InvalidSample.Point {
				vectors[currStep].T = result.T
				vectors[currStep].Samples = append(vectors[currStep].Samples, result.V)
				vectors[currStep].SampleIDs = append(vectors[currStep].SampleIDs, series.signature)
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
	// For instant queries, set the step to a positive value
	// so that the operator can terminate.
	if o.step == 0 {
		o.step = 1
	}
	o.currentStep += o.step * int64(o.numSteps)

	return vectors, nil
}

func (o *matrixSelector) loadSeries(ctx context.Context) error {
	var err error
	o.once.Do(func() {
		series, loadErr := o.storage.GetSeries(ctx, o.shard, o.numShards)
		if loadErr != nil {
			err = loadErr
			return
		}

		o.scanners = make([]matrixScanner, len(series))
		o.series = make([]labels.Labels, len(series))
		for i, s := range series {
			lbls := s.Labels()
			sort.Sort(lbls)

			o.scanners[i] = matrixScanner{
				labels:    lbls,
				signature: s.Signature,
				samples:   storage.NewBufferIterator(s.Iterator(), o.selectRange),
			}
			o.series[i] = lbls
		}
		o.vectorPool.SetStepSize(len(series))
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
