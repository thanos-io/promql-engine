// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/execution/telemetry"
	"github.com/thanos-io/promql-engine/extlabels"
	"github.com/thanos-io/promql-engine/query"
	"github.com/thanos-io/promql-engine/ringbuffer"

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
	smoothed  *smoothedIterator
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
	smoothed        bool

	opts               *query.Options
	lastTrackedSamples int
}

// NewVectorSelector creates operator which selects vector of series.
func NewVectorSelector(
	selector SeriesSelector,
	queryOpts *query.Options,
	offset time.Duration,
	batchSize int64,
	selectTimestamp bool,
	shard, numShards int,
	smoothed bool,
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
		smoothed:        smoothed,

		opts: queryOpts,
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
			series   = o.scanners[o.currentSeries]
			seriesTs = ts
		)
		for currStep := 0; currStep < n && seriesTs <= o.maxt; currStep++ {
			currStepSamples = 0
			var (
				t   int64
				v   float64
				h   *histogram.FloatHistogram
				ok  bool
				err error
			)
			if o.smoothed {
				t, v, ok, err = series.smoothed.selectPoint(seriesTs, o.lookbackDelta, o.offset)
			} else {
				t, v, h, ok, err = selectPoint(series.samples, seriesTs, o.lookbackDelta, o.offset)
			}
			if err != nil {
				return 0, err
			}
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

		if o.shouldCheckSampleLimit(fromSeries) {
			if err := o.updateSampleTracker(totalSamples); err != nil {
				return 0, err
			}
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
			}
			if o.smoothed {
				o.scanners[i].smoothed = newSmoothedIterator(s.Iterator(nil))
			} else {
				o.scanners[i].samples = storage.NewMemoizedIterator(s.Iterator(nil), o.lookbackDelta)
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

func (o *vectorSelector) updateSampleTracker(totalSamples int) error {
	delta := totalSamples - o.lastTrackedSamples
	if delta > 0 {
		o.opts.SampleTracker.Add(delta)
		o.lastTrackedSamples = totalSamples
		return o.opts.SampleTracker.CheckLimit()
	} else if delta < 0 {
		o.opts.SampleTracker.Remove(-delta)
	}
	o.lastTrackedSamples = totalSamples
	return nil
}

func (o *vectorSelector) shouldCheckSampleLimit(fromSeries int64) bool {
	seriesProcessed := o.currentSeries + 1 - fromSeries

	if seriesProcessed%sampleLimitCheckInterval == 0 {
		return true
	}

	isEndOfBatch := seriesProcessed >= o.seriesBatchSize
	isLastSeries := o.currentSeries+1 >= int64(len(o.scanners))

	return isEndOfBatch || isLastSeries
}

// smoothedIterator selects values for smoothed instant vector selectors. Like
// Prometheus' smoothSeries, it considers the float samples in
// (ts-lookbackDelta, ts+lookbackDelta], ignores stale markers, and fails if the
// window contains a native histogram. Evaluation timestamps must not decrease.
type smoothedIterator struct {
	it chunkenc.Iterator
	fh *histogram.FloatHistogram

	// prev is the newest float sample before the last reference time.
	prevT int64
	prevV float64
	// ahead holds the float samples read at or after the last reference time,
	// up to the end of its window, in timestamp order.
	ahead []fPoint
	// lastHistT is the newest native histogram read from the iterator.
	lastHistT int64

	// pending is the next non-stale sample read from the iterator but beyond
	// the window of the last reference time.
	pending     fPoint
	pendingHist bool
	hasPending  bool
	exhausted   bool
}

type fPoint struct {
	t int64
	v float64
}

func newSmoothedIterator(it chunkenc.Iterator) *smoothedIterator {
	return &smoothedIterator{
		it:        it,
		prevT:     math.MinInt64,
		lastHistT: math.MinInt64,
	}
}

// selectPoint returns the value at ts: the sample at the reference time if
// there is one, otherwise the linear interpolation between the surrounding
// samples, otherwise the previous sample carried forward.
func (s *smoothedIterator) selectPoint(ts, lookbackDelta, offset int64) (int64, float64, bool, error) {
	refTime := ts - offset
	windowEnd := refTime + lookbackDelta

	// Samples read for an earlier reference time may now precede this one.
	n := 0
	for ; n < len(s.ahead) && s.ahead[n].t < refTime; n++ {
		s.prevT, s.prevV = s.ahead[n].t, s.ahead[n].v
	}
	if n > 0 {
		s.ahead = append(s.ahead[:0], s.ahead[n:]...)
	}

	// Read the rest of the window so that a histogram anywhere in it is seen.
	for {
		if !s.hasPending {
			if err := s.readNext(); err != nil {
				return 0, 0, false, err
			}
			if !s.hasPending {
				break
			}
		}
		if s.pending.t > windowEnd {
			break
		}
		switch {
		case s.pendingHist:
			s.lastHistT = s.pending.t
		case s.pending.t < refTime:
			s.prevT, s.prevV = s.pending.t, s.pending.v
		default:
			s.ahead = append(s.ahead, s.pending)
		}
		s.hasPending = false
	}

	if s.lastHistT > refTime-lookbackDelta {
		return 0, 0, false, ringbuffer.ErrExtendedRangeHistograms
	}

	hasNext := len(s.ahead) > 0
	if hasNext && s.ahead[0].t == refTime {
		return ts, s.ahead[0].v, true, nil
	}
	if s.prevT <= refTime-lookbackDelta {
		return 0, 0, false, nil
	}
	if hasNext {
		next := s.ahead[0]
		v := s.prevV + (next.v-s.prevV)*float64(refTime-s.prevT)/float64(next.t-s.prevT)
		return ts, v, true, nil
	}
	return ts, s.prevV, true, nil
}

// readNext reads the next non-stale sample into pending.
func (s *smoothedIterator) readNext() error {
	for !s.exhausted {
		switch s.it.Next() {
		case chunkenc.ValNone:
			s.exhausted = true
			return s.it.Err()
		case chunkenc.ValFloat:
			t, v := s.it.At()
			if value.IsStaleNaN(v) {
				continue
			}
			s.pending, s.pendingHist, s.hasPending = fPoint{t: t, v: v}, false, true
			return nil
		case chunkenc.ValHistogram, chunkenc.ValFloatHistogram:
			var t int64
			t, s.fh = s.it.AtFloatHistogram(s.fh)
			if value.IsStaleNaN(s.fh.Sum) {
				continue
			}
			s.pending, s.pendingHist, s.hasPending = fPoint{t: t}, true, true
			return nil
		}
	}
	return nil
}

// TODO(fpetkovski): Add max samples limit.
func selectPoint(it *storage.MemoizedSeriesIterator, ts, lookbackDelta, offset int64) (int64, float64, *histogram.FloatHistogram, bool, error) {
	refTime := ts - offset
	var t int64
	var v float64
	var fh *histogram.FloatHistogram

	valueType := it.Seek(refTime)
	switch valueType {
	case chunkenc.ValNone:
		if it.Err() != nil {
			return 0, 0, nil, false, it.Err()
		}
	case chunkenc.ValFloatHistogram, chunkenc.ValHistogram:
		t, fh = it.AtFloatHistogram()
	case chunkenc.ValFloat:
		t, v = it.At()
	default:
		panic(errors.Newf("unknown value type %v", valueType))
	}
	if valueType == chunkenc.ValNone || t > refTime {
		var ok bool
		t, v, fh, ok = it.PeekPrev()
		if !ok || t <= refTime-lookbackDelta {
			return 0, 0, nil, false, nil
		}
	}
	if value.IsStaleNaN(v) || (fh != nil && value.IsStaleNaN(fh.Sum)) {
		return 0, 0, nil, false, nil
	}
	return t, v, fh, true, nil
}
