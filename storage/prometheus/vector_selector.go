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
	"github.com/thanos-io/promql-engine/warnings"

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
	smoothed  *smoothedIterator
	// chargedBuffered is the number of samples buffered by smoothed that are
	// charged to the sample tracker.
	chargedBuffered int
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

	reusableChunkIter chunkenc.Iterator
	reusableMemoized  *storage.MemoizedSeriesIterator
	reusableSmoothed  *smoothedIterator

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
			series          = &o.scanners[o.currentSeries]
			seriesTs        = ts
			seriesExhausted = series.exhausted
		)
		// Initialize iterators lazily per series. An instant query visits each
		// series once, so its iterator and memoized wrapper can be handed to the
		// next one; a range query revisits every series and each keeps its own.
		if o.smoothed {
			if series.smoothed == nil && !seriesExhausted {
				if o.opts.IsInstantQuery() {
					o.reusableChunkIter = series.rawSeries.Iterator(o.reusableChunkIter)
					o.reusableSmoothed.reset(o.reusableChunkIter)
					series.smoothed = o.reusableSmoothed
				} else {
					series.smoothed = newSmoothedIterator(series.rawSeries.Iterator(nil))
				}
			}
		} else if series.samples == nil && !seriesExhausted {
			if o.opts.IsInstantQuery() {
				o.reusableChunkIter = series.rawSeries.Iterator(o.reusableChunkIter)
				// Reset rebinds the wrapper to new iterator.
				o.reusableMemoized.Reset(o.reusableChunkIter)
				series.samples = o.reusableMemoized
			} else {
				series.samples = storage.NewMemoizedIterator(series.rawSeries.Iterator(nil), o.lookbackDelta)
			}
		}
		// smoothedPeak is the most samples buffered by the smoothed iterator and
		// produced as outputs that existed at the same time.
		var smoothedPeak int
		for currStep := 0; !seriesExhausted && currStep < n && seriesTs <= o.maxt; currStep++ {
			currStepSamples = 0
			var (
				t     int64
				v     float64
				h     *histogram.FloatHistogram
				ok    bool
				atEnd bool
				warn  warnings.Warnings
				err   error
			)
			if o.smoothed {
				// Smoothed values are stamped with the evaluation time.
				t = seriesTs
				v, h, ok, atEnd, warn, err = series.smoothed.selectPoint(seriesTs, o.lookbackDelta, o.offset)
				// Samples read for this step coexist with the outputs produced so far.
				smoothedPeak = max(smoothedPeak, series.smoothed.takePeak()+o.uncheckedOutputs(totalSamples))
				if warn != 0 {
					emitRingbufferWarnings(ctx, warn, series.labels.Get(labels.MetricName))
				}
			} else {
				t, v, h, ok, atEnd, err = selectPoint(series.samples, seriesTs, o.lookbackDelta, o.offset)
			}
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
			if o.smoothed {
				// The retained samples coexist with all outputs, including this one.
				smoothedPeak = max(smoothedPeak, series.smoothed.buffered+o.uncheckedOutputs(totalSamples))
			}
			o.telemetry.IncrementSamplesAtTimestamp(currStepSamples, seriesTs)
			seriesTs += o.step
		}
		if series.smoothed != nil {
			released := o.opts.IsInstantQuery() || seriesExhausted
			if err := o.trackSmoothedBuffer(series, smoothedPeak, released); err != nil {
				return 0, err
			}
		}

		if o.shouldCheckSampleLimit(fromSeries) {
			if err := o.updateSampleTracker(totalSamples); err != nil {
				return 0, err
			}
		}

		if o.opts.IsInstantQuery() {
			series.samples = nil
			series.smoothed = nil
		} else if seriesExhausted {
			series.exhausted = true
			series.samples = nil
			series.smoothed = nil
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

		// One memoized wrapper or smoothed iterator per selector instead of
		// one per series.
		if o.opts.IsInstantQuery() {
			if o.smoothed {
				o.reusableSmoothed = newSmoothedIterator(nil)
			} else {
				o.reusableMemoized = storage.NewMemoizedEmptyIterator(o.lookbackDelta)
			}
		}
	})
	return err
}

// uncheckedOutputs returns how many of the batch's outputs, totalSamples so
// far, updateSampleTracker has not charged yet.
func (o *vectorSelector) uncheckedOutputs(totalSamples int) int {
	return max(totalSamples-o.lastTrackedSamples, 0)
}

// trackSmoothedBuffer charges the samples buffered by the smoothed iterator of
// series to the sample limit. peak, the most buffered samples and uncharged
// outputs that coexisted, is checked on top of the buffers of other series.
// The charge is then lowered to what the iterator still holds; a released
// iterator is reset so that it holds nothing. updateSampleTracker charges the
// outputs.
func (o *vectorSelector) trackSmoothedBuffer(series *vectorScanner, peak int, released bool) error {
	o.opts.SampleTracker.Add(peak - series.chargedBuffered)
	if err := o.opts.SampleTracker.CheckLimit(); err != nil {
		return err
	}
	if released {
		series.smoothed.reset(nil)
	}
	o.opts.SampleTracker.Remove(peak - series.smoothed.buffered)
	series.chargedBuffered = series.smoothed.buffered
	return nil
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
// Prometheus' smoothSeries, it considers the samples in
// (ts-lookbackDelta, ts+lookbackDelta] and ignores stale markers. Evaluation
// timestamps must not decrease.
type smoothedIterator struct {
	it chunkenc.Iterator

	// prev is the newest sample before the last reference time, if prevOK.
	prev   smoothedPoint
	prevOK bool
	// lastFloatT and lastHistT are the newest float and native histogram read
	// before the last reference time. They detect windows that mix the two
	// types even though only prev is retained.
	lastFloatT int64
	lastHistT  int64
	// ahead holds the samples read at or after the last reference time, up to
	// the end of its window, in timestamp order.
	ahead []smoothedPoint

	// pending is the next non-stale sample read from the iterator but beyond
	// the window of the last reference time.
	pending    smoothedPoint
	hasPending bool
	exhausted  bool

	// buffered is the number of samples held in prev, ahead and pending,
	// counting histograms by size, for sample limit accounting. peak is the
	// largest value of buffered since the last call to takePeak.
	buffered int
	peak     int
}

type smoothedPoint struct {
	t int64
	f float64
	// h is nil for float samples.
	h *histogram.FloatHistogram
}

func newSmoothedIterator(it chunkenc.Iterator) *smoothedIterator {
	s := &smoothedIterator{}
	s.reset(it)
	return s
}

// reset rebinds the iterator to a new series, keeping its buffer.
func (s *smoothedIterator) reset(it chunkenc.Iterator) {
	// Drop references to the previous series' histograms.
	clear(s.ahead)
	*s = smoothedIterator{
		it:         it,
		lastFloatT: math.MinInt64,
		lastHistT:  math.MinInt64,
		ahead:      s.ahead[:0],
	}
}

// takePeak returns the largest number of samples buffered since the last call.
func (s *smoothedIterator) takePeak() int {
	peak := max(s.peak, s.buffered)
	s.peak = s.buffered
	return peak
}

func smoothedPointSamples(p smoothedPoint) int {
	if p.h != nil {
		return telemetry.CalculateHistogramSampleCount(p.h)
	}
	return 1
}

// setPrev records p as the newest sample before the reference time. The caller
// transfers p's share of buffered.
func (s *smoothedIterator) setPrev(p smoothedPoint) {
	if s.prevOK {
		s.buffered -= smoothedPointSamples(s.prev)
	}
	s.prev, s.prevOK = p, true
	if p.h != nil {
		s.lastHistT = p.t
	} else {
		s.lastFloatT = p.t
	}
}

// selectPoint returns the value at ts: the sample at the reference time if
// there is one, otherwise the linear interpolation between the surrounding
// samples, otherwise the previous sample carried forward. exhausted reports
// that no later evaluation can return a value either.
func (s *smoothedIterator) selectPoint(ts, lookbackDelta, offset int64) (_ float64, _ *histogram.FloatHistogram, ok, exhausted bool, _ warnings.Warnings, _ error) {
	var (
		refTime     = ts - offset
		windowStart = refTime - lookbackDelta // exclusive
		windowEnd   = refTime + lookbackDelta
	)

	// Samples read for an earlier reference time may now precede this one.
	n := 0
	for ; n < len(s.ahead) && s.ahead[n].t < refTime; n++ {
		s.setPrev(s.ahead[n])
	}
	if n > 0 {
		kept := copy(s.ahead, s.ahead[n:])
		// Drop references to histograms in the vacated slots.
		clear(s.ahead[kept:])
		s.ahead = s.ahead[:kept]
	}

	// Read the rest of the window so that its sample types are all known.
	for {
		if !s.hasPending {
			if err := s.readNext(); err != nil {
				return 0, nil, false, false, 0, err
			}
			if !s.hasPending {
				break
			}
		}
		if s.pending.t > windowEnd {
			break
		}
		if s.pending.t < refTime {
			s.setPrev(s.pending)
		} else {
			s.ahead = append(s.ahead, s.pending)
		}
		s.pending, s.hasPending = smoothedPoint{}, false
	}

	prevOK := s.prevOK && s.prev.t > windowStart
	if !prevOK && len(s.ahead) == 0 {
		return 0, nil, false, s.exhausted && !s.hasPending, 0, nil
	}

	hasFloat, hasHist := s.lastFloatT > windowStart, s.lastHistT > windowStart
	for _, p := range s.ahead {
		hasFloat = hasFloat || p.h == nil
		hasHist = hasHist || p.h != nil
	}
	if hasFloat && hasHist {
		return 0, nil, false, false, warnings.WarnMixedFloatsHistograms, nil
	}

	switch {
	case len(s.ahead) > 0 && s.ahead[0].t == refTime:
		// Exact match.
		next := s.ahead[0]
		if next.h != nil {
			return 0, next.h.Copy(), true, false, 0, nil
		}
		return next.f, nil, true, false, 0, nil
	case !prevOK:
		// No earlier sample in the window to interpolate from.
		return 0, nil, false, false, 0, nil
	case len(s.ahead) > 0:
		prev, next := s.prev, s.ahead[0]
		if prev.h == nil {
			return prev.f + (next.f-prev.f)*float64(refTime-prev.t)/float64(next.t-prev.t), nil, true, false, 0, nil
		}
		if prev.h.UsesCustomBuckets() != next.h.UsesCustomBuckets() {
			return 0, nil, false, false, warnings.WarnMixedExponentialCustomBuckets, nil
		}
		// Treat the histograms as counters unless both carry the gauge hint.
		isCounter := prev.h.CounterResetHint != histogram.GaugeType || next.h.CounterResetHint != histogram.GaugeType
		h, warn, err := ringbuffer.InterpolateHistograms(prev.h, prev.t, next.h, next.t, refTime, isCounter)
		if err != nil {
			if errors.Is(err, histogram.ErrHistogramsIncompatibleSchema) {
				warn |= warnings.WarnMixedExponentialCustomBuckets
			}
			return 0, nil, false, false, warn, nil
		}
		return 0, h, true, false, warn, nil
	default:
		// No later sample yet; carry the previous value forward.
		if s.prev.h == nil {
			return s.prev.f, nil, true, false, 0, nil
		}
		// The hint describes the relation to the preceding sample, not the value.
		h := s.prev.h.Copy()
		h.CounterResetHint = histogram.UnknownCounterReset
		return 0, h, true, false, 0, nil
	}
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
			s.pending, s.hasPending = smoothedPoint{t: t, f: v}, true
			s.buffered++
			s.peak = max(s.peak, s.buffered)
			return nil
		case chunkenc.ValHistogram, chunkenc.ValFloatHistogram:
			// Each sample is retained, so it needs its own histogram.
			t, fh := s.it.AtFloatHistogram(nil)
			if value.IsStaleNaN(fh.Sum) {
				continue
			}
			s.pending, s.hasPending = smoothedPoint{t: t, h: fh}, true
			s.buffered += smoothedPointSamples(s.pending)
			s.peak = max(s.peak, s.buffered)
			return nil
		}
	}
	return nil
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
		_, t, v, fh, hasPrev = it.PeekPrev()
		if !hasPrev || t <= refTime-lookbackDelta {
			return 0, 0, nil, false, valueType == chunkenc.ValNone, nil
		}
	}
	if value.IsStaleNaN(v) || (fh != nil && value.IsStaleNaN(fh.Sum)) {
		return 0, 0, nil, false, false, nil
	}
	return t, v, fh, true, false, nil
}
