// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package exchange

import (
	"context"
	"sync"

	"github.com/prometheus/prometheus/model/labels"

	"github.com/thanos-community/promql-engine/execution/model"
)

type dedupSample struct {
	t int64
	v float64
}

// The dedupCache is an internal cache used to deduplicate samples inside a single step vector.
type dedupCache []dedupSample

type dedupOperator struct {
	once   sync.Once
	series []labels.Labels

	pool *model.VectorPool
	next model.VectorOperator
	// outputIndex is a slice that is used as an index from input sample ID to output sample ID.
	outputIndex []uint64
	outputCache dedupCache
}

func NewDedupOperator(pool *model.VectorPool, next model.VectorOperator) model.VectorOperator {
	return &dedupOperator{
		next: next,
		pool: pool,
	}
}

func (d *dedupOperator) Next(ctx context.Context) ([]model.StepVector, error) {
	var err error
	d.once.Do(func() { err = d.loadSeries(ctx) })
	if err != nil {
		return nil, err
	}

	in, err := d.next.Next(ctx)
	if err != nil {
		return nil, err
	}
	if in == nil {
		return nil, nil
	}

	result := d.pool.GetVectorBatch()
	for _, vector := range in {
		for i, inputSampleID := range vector.SampleIDs {
			outputSampleID := d.outputIndex[inputSampleID]
			d.outputCache[outputSampleID].t = vector.T
			d.outputCache[outputSampleID].v = vector.Samples[i]
		}

		out := d.pool.GetStepVector(vector.T)
		for outputSampleID, sample := range d.outputCache {
			if sample.t == vector.T {
				out.SampleIDs = append(out.SampleIDs, uint64(outputSampleID))
				out.Samples = append(out.Samples, sample.v)
			}
		}
		result = append(result, out)
	}

	return result, nil
}

func (d *dedupOperator) Series(ctx context.Context) ([]labels.Labels, error) {
	var err error
	d.once.Do(func() { err = d.loadSeries(ctx) })
	if err != nil {
		return nil, err
	}
	return d.series, nil
}

func (d *dedupOperator) GetPool() *model.VectorPool {
	return d.pool
}

func (d *dedupOperator) Explain() (me string, next []model.VectorOperator) {
	return "[*dedup]", []model.VectorOperator{d.next}
}

func (d *dedupOperator) loadSeries(ctx context.Context) error {
	series, err := d.next.Series(ctx)
	if err != nil {
		return err
	}

	outputIndex := make(map[uint64]uint64)
	inputIndex := make([]uint64, len(series))
	hashBuf := make([]byte, 0, 128)
	for inputSeriesID, inputSeries := range series {
		hash := hashSeries(hashBuf, inputSeries)

		inputIndex[inputSeriesID] = hash
		outputSeriesID, ok := outputIndex[hash]
		if !ok {
			outputSeriesID = uint64(len(d.series))
			d.series = append(d.series, inputSeries)
		}
		outputIndex[hash] = outputSeriesID
	}

	d.outputIndex = make([]uint64, len(inputIndex))
	for inputSeriesID, hash := range inputIndex {
		outputSeriesID := outputIndex[hash]
		d.outputIndex[inputSeriesID] = outputSeriesID
	}
	d.outputCache = make(dedupCache, len(outputIndex))
	for i := range d.outputCache {
		d.outputCache[i].t = -1
	}

	return nil
}

func hashSeries(hashBuf []byte, inputSeries labels.Labels) uint64 {
	hashBuf = hashBuf[:0]
	hash, _ := inputSeries.HashWithoutLabels(hashBuf)
	return hash
}
