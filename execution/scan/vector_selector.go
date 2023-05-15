// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package scan

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/efficientgo/core/errors"
	"github.com/prometheus/prometheus/tsdb/chunkenc"

	"github.com/thanos-community/promql-engine/execution/model"
	engstore "github.com/thanos-community/promql-engine/execution/storage"
	"github.com/thanos-community/promql-engine/query"

	"github.com/prometheus/prometheus/model/histogram"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/value"

	"github.com/prometheus/prometheus/storage"
)

type vectorScanner struct {
	labels    labels.Labels
	signature uint64
	samples   *storage.MemoizedSeriesIterator
}

type vectorSelector struct {
	storage  engstore.SeriesSelector
	scanners []vectorScanner
	series   []labels.Labels

	once       sync.Once
	vectorPool *model.VectorPool

	numSteps      int
	mint          int64
	maxt          int64
	lookbackDelta int64
	step          int64
	currentStep   int64
	offset        int64

	shard     int
	numShards int
}

type point struct {
	t  int64
	v  float64
	fh *histogram.FloatHistogram
}

// NewVectorSelector creates operator which selects vector of series.
func NewVectorSelector(
	pool *model.VectorPool,
	selector engstore.SeriesSelector,
	queryOpts *query.Options,
	offset time.Duration,
	shard, numShards int,
) model.VectorOperator {
	return &vectorSelector{
		storage:    selector,
		vectorPool: pool,

		mint:          queryOpts.Start.UnixMilli(),
		maxt:          queryOpts.End.UnixMilli(),
		step:          queryOpts.Step.Milliseconds(),
		currentStep:   queryOpts.Start.UnixMilli(),
		lookbackDelta: queryOpts.LookbackDelta.Milliseconds(),
		offset:        offset.Milliseconds(),
		numSteps:      queryOpts.NumSteps(),

		shard:     shard,
		numShards: numShards,
	}
}

func (o *vectorSelector) Explain() (me string, next []model.VectorOperator) {
	return fmt.Sprintf("[*vectorSelector] {%v} %v mod %v", o.storage.Matchers(), o.shard, o.numShards), nil
}

func (o *vectorSelector) Series(ctx context.Context) ([]labels.Labels, error) {
	if err := o.loadSeries(ctx); err != nil {
		return nil, err
	}
	return o.series, nil
}

func (o *vectorSelector) GetPool() *model.VectorPool {
	return o.vectorPool
}

func (o *vectorSelector) Next(ctx context.Context) ([]model.StepVector, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

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
			series       = o.scanners[i]
			seriesTs     = ts
			lastSampleTs int64 // Added variable to store timestamp of the last sample in the lookback period.
		)

		for currStep := 0; currStep < o.numSteps && seriesTs <= o.maxt; currStep++ {
			if len(vectors) <= currStep {
				vectors = append(vectors, o.vectorPool.GetStepVector(seriesTs))
			}

			// Modify selectPoint call to retrieve timestamp of the last sample in the lookback period.
			p, ok, err := selectPoint(series.samples, seriesTs, o.lookbackDelta, o.offset)
			if err != nil {
				return nil, err
			}

			if ok {
				if p.fh != nil {
					vectors[currStep].AppendHistogram(o.vectorPool, series.signature, p.fh)
				} else {
					vectors[currStep].AppendSample(o.vectorPool, series.signature, p.v)
				}

				// Save the timestamp of the last sample in the lookback period.
				lastSampleTs = p.t
			}

			seriesTs += o.step
		}

		// Use the saved timestamp to compute timestamp of last sample in lookback period.
		if lastSampleTs > 0 {
			vectors[len(vectors)-1].T = lastSampleTs
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

func (o *vectorSelector) loadSeries(ctx context.Context) error {
	var err error
	o.once.Do(func() {
		series, loadErr := o.storage.GetSeries(ctx, o.shard, o.numShards)
		if loadErr != nil {
			err = loadErr
			return
		}

		o.scanners = make([]vectorScanner, len(series))
		o.series = make([]labels.Labels, len(series))
		for i, s := range series {
			o.scanners[i] = vectorScanner{
				labels:    s.Labels(),
				signature: s.Signature,
				samples:   storage.NewMemoizedIterator(s.Iterator(nil), o.lookbackDelta),
			}
			o.series[i] = s.Labels()
		}
		o.vectorPool.SetStepSize(len(series))
	})
	return err
}

// TODO(fpetkovski): Add max samples limit.
// To push down the timestamp function into this file and store the timestamp in the value for each series, you can modify the selectPoint function.
func selectPoint(it *storage.MemoizedSeriesIterator, ts, lookbackDelta, offset int64) (point, bool, error) {
	refTime := ts - offset
	var p point

	valueType := it.Seek(refTime)
	switch valueType {
	case chunkenc.ValNone:
		if it.Err() != nil {
			return p, false, it.Err()
		}
	case chunkenc.ValFloatHistogram, chunkenc.ValHistogram:
		t, fh := it.AtFloatHistogram()
		p = point{t: t, fh: fh}
	case chunkenc.ValFloat:
		t, v := it.At()
		p = point{t: t, v: v}
	default:
		panic(errors.Newf("unknown value type %v", valueType))
	}

	if valueType == chunkenc.ValNone || p.t > refTime {
		var ok bool
		p.t, p.v, _, p.fh, ok = it.PeekPrev()
		if !ok || p.t < refTime-lookbackDelta {
			return p, false, nil
		}
	}

	if value.IsStaleNaN(p.v) || (p.fh != nil && value.IsStaleNaN(p.fh.Sum)) {
		return p, false, nil
	}

	return p, true, nil
}
