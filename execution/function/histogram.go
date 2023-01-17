// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package function

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"

	"github.com/cespare/xxhash/v2"
	"github.com/prometheus/prometheus/model/histogram"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"

	"github.com/thanos-community/promql-engine/execution/model"
)

type histogramSeries struct {
	outputID   int
	upperBound float64
}

// histogramOperator is a function operator that calculates percentiles.
type histogramOperator struct {
	pool *model.VectorPool

	funcArgs parser.Expressions

	once     sync.Once
	series   []labels.Labels
	scalarOp model.VectorOperator
	vectorOp model.VectorOperator

	// scalarPoints is a reusable buffer for points from the first argument of histogram_quantile.
	scalarPoints []float64

	// outputIndex is a mapping from input series ID to the output series ID and its upper boundary value
	// parsed from the le label.
	// If outputIndex[i] is nil then series[i] has no valid `le` label.
	outputIndex []*histogramSeries

	// seriesBuckets are the buckets for each individual series.
	seriesBuckets []buckets
}

func NewHistogramOperator(pool *model.VectorPool, args parser.Expressions, nextOps []model.VectorOperator, stepsBatch int) (model.VectorOperator, error) {
	return &histogramOperator{
		pool:         pool,
		funcArgs:     args,
		once:         sync.Once{},
		scalarOp:     nextOps[0],
		vectorOp:     nextOps[1],
		scalarPoints: make([]float64, stepsBatch),
	}, nil
}

func (o *histogramOperator) Explain() (me string, next []model.VectorOperator) {
	next = []model.VectorOperator{o.scalarOp, o.vectorOp}
	return fmt.Sprintf("[*functionOperator] histogram_quantile(%v)", o.funcArgs), next
}

func (o *histogramOperator) Series(ctx context.Context) ([]labels.Labels, error) {
	var err error
	o.once.Do(func() { err = o.loadSeries(ctx) })
	if err != nil {
		return nil, err
	}

	return o.series, nil
}

func (o *histogramOperator) GetPool() *model.VectorPool {
	return o.pool
}

func (o *histogramOperator) Next(ctx context.Context) ([]model.StepVector, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var err error
	o.once.Do(func() { err = o.loadSeries(ctx) })
	if err != nil {
		return nil, err
	}

	scalars, err := o.scalarOp.Next(ctx)
	if err != nil {
		return nil, err
	}

	if len(scalars) == 0 {
		return nil, nil
	}

	vectors, err := o.vectorOp.Next(ctx)
	if err != nil {
		return nil, err
	}

	o.scalarPoints = o.scalarPoints[:0]
	for _, scalar := range scalars {
		if len(scalar.Samples) > 0 {
			o.scalarPoints = append(o.scalarPoints, scalar.Samples[0])
		}
		o.scalarOp.GetPool().PutStepVector(scalar)
	}
	o.scalarOp.GetPool().PutVectors(scalars)

	return o.processInputSeries(vectors)
}

func (o *histogramOperator) processInputSeries(vectors []model.StepVector) ([]model.StepVector, error) {
	out := o.pool.GetVectorBatch()
	for stepIndex, vector := range vectors {
		if len(vector.HistogramSamples) > 0 {
			// Deal with the sparse histograms.
			step := o.pool.GetStepVector(vector.T)
			for i, sample := range vector.HistogramSamples {
				val := histogramQuantile(o.scalarPoints[stepIndex], sample)
				step.SampleIDs = append(step.SampleIDs, uint64(i))
				step.Samples = append(step.Samples, val)
			}
			out = append(out, step)
			continue
		}
		o.resetBuckets()
		for i, seriesID := range vector.SampleIDs {
			outputSeries := o.outputIndex[seriesID]
			// This means that it has an invalid `le` label.
			if outputSeries == nil {
				continue
			}
			outputSeriesID := outputSeries.outputID
			bucket := le{
				upperBound: outputSeries.upperBound,
				count:      vector.Samples[i],
			}
			o.seriesBuckets[outputSeriesID] = append(o.seriesBuckets[outputSeriesID], bucket)
		}

		step := o.pool.GetStepVector(vector.T)
		for i, stepBuckets := range o.seriesBuckets {
			// It could be zero if multiple input series map to the same output series ID.
			if len(stepBuckets) == 0 {
				continue
			}
			// If there is only bucket or if we are after how many
			// scalar points we have then it needs to be NaN.
			if len(stepBuckets) == 1 || stepIndex >= len(o.scalarPoints) {
				step.SampleIDs = append(step.SampleIDs, uint64(i))
				step.Samples = append(step.Samples, math.NaN())
				continue
			}

			val := bucketQuantile(o.scalarPoints[stepIndex], stepBuckets)
			step.SampleIDs = append(step.SampleIDs, uint64(i))
			step.Samples = append(step.Samples, val)
		}
		out = append(out, step)
		o.vectorOp.GetPool().PutStepVector(vector)
	}

	o.vectorOp.GetPool().PutVectors(vectors)
	return out, nil
}

func (o *histogramOperator) loadSeries(ctx context.Context) error {
	series, err := o.vectorOp.Series(ctx)
	if err != nil {
		return err
	}

	var (
		hashBuf      = make([]byte, 0, 256)
		hasher       = xxhash.New()
		seriesHashes = make(map[uint64]int, len(series))
	)

	o.series = make([]labels.Labels, 0)
	o.outputIndex = make([]*histogramSeries, len(series))

	for i, s := range series {
		lbls, bucketLabel := dropLabel(s.Copy(), "le")
		value, err := strconv.ParseFloat(bucketLabel.Value, 64)
		if err != nil {
			continue
		}
		lbls, _ = DropMetricName(lbls)

		hasher.Reset()
		hashBuf = lbls.Bytes(hashBuf)
		if _, err := hasher.Write(hashBuf); err != nil {
			return err
		}

		seriesHash := hasher.Sum64()
		seriesID, ok := seriesHashes[seriesHash]
		if !ok {
			o.series = append(o.series, lbls)
			seriesID = len(o.series) - 1
			seriesHashes[seriesHash] = seriesID
		}

		o.outputIndex[i] = &histogramSeries{
			outputID:   seriesID,
			upperBound: value,
		}
	}
	o.seriesBuckets = make([]buckets, len(o.series))
	o.pool.SetStepSize(len(o.series))
	return nil
}

func (o *histogramOperator) resetBuckets() {
	for i := range o.seriesBuckets {
		o.seriesBuckets[i] = o.seriesBuckets[i][:0]
	}
}

// histogramQuantile calculates the quantile 'q' based on the given histogram.
//
// The quantile value is interpolated assuming a linear distribution within a
// bucket.
// TODO(beorn7): Find an interpolation method that is a better fit for
// exponential buckets (and think about configurable interpolation).
//
// A natural lower bound of 0 is assumed if the histogram has only positive
// buckets. Likewise, a natural upper bound of 0 is assumed if the histogram has
// only negative buckets.
// TODO(beorn7): Come to terms if we want that.
//
// There are a number of special cases (once we have a way to report errors
// happening during evaluations of AST functions, we should report those
// explicitly):
//
// If the histogram has 0 observations, NaN is returned.
//
// If q<0, -Inf is returned.
//
// If q>1, +Inf is returned.
//
// If q is NaN, NaN is returned.
func histogramQuantile(q float64, h *histogram.FloatHistogram) float64 {
	if q < 0 {
		return math.Inf(-1)
	}
	if q > 1 {
		return math.Inf(+1)
	}

	if h.Count == 0 || math.IsNaN(q) {
		return math.NaN()
	}

	var (
		bucket histogram.Bucket[float64]
		count  float64
		it     = h.AllBucketIterator()
		rank   = q * h.Count
	)
	for it.Next() {
		bucket = it.At()
		count += bucket.Count
		if count >= rank {
			break
		}
	}
	if bucket.Lower < 0 && bucket.Upper > 0 {
		if len(h.NegativeBuckets) == 0 && len(h.PositiveBuckets) > 0 {
			// The result is in the zero bucket and the histogram has only
			// positive buckets. So we consider 0 to be the lower bound.
			bucket.Lower = 0
		} else if len(h.PositiveBuckets) == 0 && len(h.NegativeBuckets) > 0 {
			// The result is in the zero bucket and the histogram has only
			// negative buckets. So we consider 0 to be the upper bound.
			bucket.Upper = 0
		}
	}
	// Due to numerical inaccuracies, we could end up with a higher count
	// than h.Count. Thus, make sure count is never higher than h.Count.
	if count > h.Count {
		count = h.Count
	}
	// We could have hit the highest bucket without even reaching the rank
	// (this should only happen if the histogram contains observations of
	// the value NaN), in which case we simply return the upper limit of the
	// highest explicit bucket.
	if count < rank {
		return bucket.Upper
	}

	rank -= count - bucket.Count
	// TODO(codesome): Use a better estimation than linear.
	return bucket.Lower + (bucket.Upper-bucket.Lower)*(rank/bucket.Count)
}
