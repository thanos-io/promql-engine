package scan

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/fpetkovski/promql-engine/operators/model"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/value"

	"github.com/prometheus/prometheus/storage"
)

type vectorScan struct {
	labels    labels.Labels
	signature uint64
	samples   *storage.MemoizedSeriesIterator
}

type vectorSelector struct {
	storage  *seriesSelector
	scanners []vectorScan
	series   []labels.Labels

	once       sync.Once
	vectorPool *model.VectorPool

	mint        int64
	maxt        int64
	step        int64
	currentStep int64

	shard     int
	numShards int
}

func (o *vectorSelector) Series(ctx context.Context) ([]labels.Labels, error) {
	if err := o.loadSeries(ctx); err != nil {
		return nil, err
	}
	return o.series, nil
}

func NewVectorSelector(pool *model.VectorPool, storage *seriesSelector, mint, maxt time.Time, step time.Duration, shard, numShards int) model.Vector {
	// TODO(fpetkovski): Add offset parameter.
	return &vectorSelector{
		storage:    storage,
		vectorPool: pool,

		mint:        mint.UnixMilli(),
		maxt:        maxt.UnixMilli(),
		step:        step.Milliseconds(),
		currentStep: mint.UnixMilli(),

		shard:     shard,
		numShards: numShards,
	}
}

func (o *vectorSelector) GetPool() *model.VectorPool {
	return o.vectorPool
}

func (o *vectorSelector) Next(ctx context.Context) ([]model.StepVector, error) {
	if o.currentStep > o.maxt {
		return nil, nil
	}

	if err := o.loadSeries(ctx); err != nil {
		return nil, err
	}

	stepsBatch := 10
	totalSteps := (o.maxt+o.mint)/o.step + 1
	numSteps := int(math.Min(float64(stepsBatch), float64(totalSteps)))

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
			_, v, ok := selectPoint(series.samples, seriesTs)
			if ok {
				vectors[currStep].Samples = append(vectors[currStep].Samples, model.StepSample{
					ID:     series.signature,
					Metric: series.labels,
					V:      v,
				})
			}
			seriesTs += o.step
		}
	}
	o.currentStep += o.step * int64(numSteps)

	return vectors, nil
}

func (o *vectorSelector) loadSeries(ctx context.Context) error {
	var err error
	o.once.Do(func() {
		scanners, seriesShard, loadErr := o.storage.Series(ctx, o.shard, o.numShards)
		if loadErr != nil {
			err = loadErr
			return
		}

		o.scanners = scanners
		o.series = seriesShard
	})
	return err
}

// TODO(fpetkovski): Add error handling and max samples limit.
func selectPoint(it *storage.MemoizedSeriesIterator, ts int64) (int64, float64, bool) {
	lookbackDelta := 5 * time.Minute.Milliseconds()
	refTime := ts
	var t int64
	var v float64

	ok := it.Seek(refTime)
	if ok {
		t, v = it.At()
	}

	if !ok || t > refTime {
		t, v, ok = it.PeekPrev()
		if !ok || t < refTime-lookbackDelta {
			return 0, 0, false
		}
	}
	if value.IsStaleNaN(v) {
		return 0, 0, false
	}
	return t, v, true
}
