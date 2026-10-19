// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/execution/telemetry"
	"github.com/thanos-io/promql-engine/extlabels"
	"github.com/thanos-io/promql-engine/query"

	"github.com/efficientgo/core/errors"
	"github.com/prometheus/prometheus/model/histogram"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/value"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
)

type vectorScanner struct {
	labels    labels.Labels
	signature uint64
	samples   *storage.MemoizedSeriesIterator
	rawSeries SignedSeries
	exhausted bool
}

type vectorSelector struct {
	telemetry telemetry.OperatorTelemetry

	storage  SeriesSelector
	scanners []vectorScanner
	series   []labels.Labels

	once sync.Once

	numSteps        int
	mint            int64
	maxt            int64
	lookbackDelta   int64
	step            int64
	offset          int64
	seriesBatchSize int64

	currentSeries int64
	currentStep   int64

	shard     int
	numShards int

	selectTimestamp bool

	opts *query.Options
	// lastTrackedSamples is the sample total at the last report; totalSamples restarts each Next.
	lastTrackedSamples        int
	sampleLimitCheckThreshold int
}

// NewVectorSelector creates operator which selects vector of series.
func NewVectorSelector(
	selector SeriesSelector,
	queryOpts *query.Options,
	offset time.Duration,
	batchSize int64,
	selectTimestamp bool,
	shard, numShards int,
) model.VectorOperator {
	o := &vectorSelector{
		storage: selector,

		mint:            queryOpts.Start.UnixMilli(),
		maxt:            queryOpts.End.UnixMilli(),
		step:            queryOpts.Step.Milliseconds(),
		currentStep:     queryOpts.Start.UnixMilli(),
		lookbackDelta:   queryOpts.LookbackDelta.Milliseconds(),
		offset:          offset.Milliseconds(),
		numSteps:        queryOpts.NumStepsPerBatch(),
		seriesBatchSize: batchSize,

		shard:     shard,
		numShards: numShards,

		selectTimestamp: selectTimestamp,

		opts: queryOpts,

		sampleLimitCheckThreshold: query.ComputeSampleLimitCheckThreshold(queryOpts, queryOpts.DecodingConcurrency),
	}

	// For instant queries, set the step to a positive value
	// so that the operator can terminate.
	if o.step == 0 {
		o.step = 1
	}

	o.telemetry = telemetry.NewTelemetry(o, queryOpts)
	return telemetry.NewOperator(o.telemetry, o)
}

func (o *vectorSelector) String() string {
	return fmt.Sprintf("[vectorSelector] {%v} %v mod %v", o.storage.Matchers(), o.shard, o.numShards)
}

func (o *vectorSelector) Explain() (next []model.VectorOperator) {
	return nil
}

func (o *vectorSelector) Series(ctx context.Context) ([]labels.Labels, error) {
	if err := o.loadSeries(ctx); err != nil {
		return nil, err
	}
	return o.series, nil
}

func (o *vectorSelector) Next(ctx context.Context, buf []model.StepVector) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}
	if o.currentStep > o.maxt {
		return 0, nil
	}

	if err := o.loadSeries(ctx); err != nil {
		return 0, err
	}

	ts := o.currentStep
	n := 0
	maxSteps := min(o.numSteps, len(buf))

	// Calculate expected samples per step: the actual number of series we'll process this batch.
	// This is min(seriesBatchSize, remaining series to process).
	remainingSeries := int64(len(o.scanners)) - o.currentSeries
	expectedSamples := int(min(o.seriesBatchSize, remainingSeries))
	if expectedSamples <= 0 {
		expectedSamples = len(o.scanners)
	}

	for currStep := 0; currStep < maxSteps && ts <= o.maxt; currStep++ {
		buf[n].Reset(ts)
		n++
		ts += o.step
	}

	var currStepSamples int
	var totalSamples int
	// Reset the current timestamp.
	ts = o.currentStep
	fromSeries := o.currentSeries

	for ; o.currentSeries-fromSeries < o.seriesBatchSize && o.currentSeries < int64(len(o.scanners)); o.currentSeries++ {
		var (
			series          = &o.scanners[o.currentSeries]
			seriesTs        = ts
			seriesExhausted = series.exhausted
		)
		if series.samples == nil && !seriesExhausted {
			series.samples = storage.NewMemoizedIterator(series.rawSeries.Iterator(nil), o.lookbackDelta)
		}
		for currStep := 0; !seriesExhausted && currStep < n && seriesTs <= o.maxt; currStep++ {
			currStepSamples = 0
			t, v, h, ok, atEnd, err := selectPoint(series.samples, seriesTs, o.lookbackDelta, o.offset)
			if err != nil {
				return 0, err
			}
			seriesExhausted = atEnd
			if o.selectTimestamp {
				v = float64(t) / 1000
			}
			if ok {
				if h != nil && !o.selectTimestamp {
					// Lazy pre-allocate histogram slices only when we actually have histograms
					buf[currStep].AppendHistogramWithSizeHint(series.signature, h, expectedSamples)
					currStepSamples += telemetry.CalculateHistogramSampleCount(h)
				} else {
					// Lazy pre-allocate sample slices with capacity hint
					buf[currStep].AppendSampleWithSizeHint(series.signature, v, expectedSamples)
					currStepSamples++
				}
				totalSamples += currStepSamples
			}
			o.telemetry.IncrementSamplesAtTimestamp(currStepSamples, seriesTs)
			seriesTs += o.step
		}

		if untracked := totalSamples - o.lastTrackedSamples; o.shouldCheckSampleLimit(fromSeries, untracked) {
			if err := o.updateSampleTracker(untracked); err != nil {
				return 0, err
			}
			o.lastTrackedSamples = totalSamples
		}

		if o.opts.IsInstantQuery() {
			series.samples = nil
		} else if seriesExhausted {
			series.exhausted = true
			series.samples = nil
		}
	}

	if o.currentSeries == int64(len(o.scanners)) {
		o.currentStep += o.step * int64(n)
		o.currentSeries = 0
	}
	return n, nil
}

func (o *vectorSelector) loadSeries(ctx context.Context) error {
	var err error
	o.once.Do(func() {
		series, loadErr := o.storage.GetSeries(ctx, o.shard, o.numShards)
		if loadErr != nil {
			err = loadErr
			return
		}

		b := labels.NewBuilder(labels.EmptyLabels())
		o.scanners = make([]vectorScanner, len(series))
		o.series = make([]labels.Labels, len(series))
		for i, s := range series {
			o.scanners[i] = vectorScanner{
				labels:    s.Labels(),
				signature: s.Signature,
				rawSeries: s,
			}
			b.Reset(s.Labels())
			// if we have pushed down a timestamp function into the scan we need to drop
			// the reserved labels (__name__, __type__, __unit__)
			if o.selectTimestamp {
				b.Del(labels.MetricName)
				b.Del(extlabels.MetricType)
				b.Del(extlabels.MetricUnit)
			}
			o.series[i] = b.Labels()
		}

		numSeries := int64(len(o.series))
		if o.seriesBatchSize == 0 || numSeries < o.seriesBatchSize {
			o.seriesBatchSize = numSeries
		}
	})
	return err
}

func (o *vectorSelector) updateSampleTracker(delta int) error {
	if delta > 0 {
		o.opts.SampleTracker.Add(delta)
		return o.opts.SampleTracker.CheckLimit()
	} else if delta < 0 {
		o.opts.SampleTracker.Remove(-delta)
	}
	return nil
}

func (o *vectorSelector) shouldCheckSampleLimit(fromSeries int64, untracked int) bool {
	if untracked > o.sampleLimitCheckThreshold {
		return true
	}
	if untracked == 0 {
		return false
	}
	seriesProcessed := o.currentSeries + 1 - fromSeries
	return seriesProcessed >= o.seriesBatchSize || o.currentSeries+1 >= int64(len(o.scanners))
}

// TODO(fpetkovski): Add max samples limit.
func selectPoint(it *storage.MemoizedSeriesIterator, ts, lookbackDelta, offset int64) (_ int64, _ float64, _ *histogram.FloatHistogram, ok, exhausted bool, _ error) {
	refTime := ts - offset
	var t int64
	var v float64
	var fh *histogram.FloatHistogram

	valueType := it.Seek(refTime)
	switch valueType {
	case chunkenc.ValNone:
		if it.Err() != nil {
			return 0, 0, nil, false, false, it.Err()
		}
	case chunkenc.ValFloatHistogram, chunkenc.ValHistogram:
		t, fh = it.AtFloatHistogram()
	case chunkenc.ValFloat:
		t, v = it.At()
	default:
		panic(errors.Newf("unknown value type %v", valueType))
	}
	if valueType == chunkenc.ValNone || t > refTime {
		var hasPrev bool
		t, v, fh, hasPrev = it.PeekPrev()
		if !hasPrev || t <= refTime-lookbackDelta {
			return 0, 0, nil, false, valueType == chunkenc.ValNone, nil
		}
	}
	if value.IsStaleNaN(v) || (fh != nil && value.IsStaleNaN(fh.Sum)) {
		return 0, 0, nil, false, false, nil
	}
	return t, v, fh, true, false, nil
}
