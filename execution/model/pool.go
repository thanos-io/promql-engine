// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package model

import (
	"sync"

	"github.com/prometheus/prometheus/model/histogram"
)

type VectorPool struct {
	vectors sync.Pool

	stepSize         int
	samples          sync.Pool
	sampleIDs        sync.Pool
	histogramSamples sync.Pool
}

func NewVectorPool(stepsBatch int) *VectorPool {
	pool := &VectorPool{}
	pool.vectors = sync.Pool{
		New: func() any {
			sv := make([]StepVector, 0, stepsBatch)
			return &sv
		},
	}
	pool.samples = sync.Pool{
		New: func() any {
			samples := make([]float64, 0, pool.stepSize)
			return &samples
		},
	}
	pool.sampleIDs = sync.Pool{
		New: func() any {
			sampleIDs := make([]uint64, 0, pool.stepSize)
			return &sampleIDs
		},
	}
	pool.histogramSamples = sync.Pool{
		New: func() any {
			histogramSamples := make([]*histogram.FloatHistogram, 0, pool.stepSize)
			return &histogramSamples
		},
	}
	return pool
}

func (p *VectorPool) GetVectorBatch() []StepVector {
	return *p.vectors.Get().(*[]StepVector)
}

func (p *VectorPool) PutVectors(vector []StepVector) {
	vector = vector[:0]
	p.vectors.Put(&vector)
}

func (p *VectorPool) GetStepVector(t int64) StepVector {
	return StepVector{
		T:                t,
		SampleIDs:        *p.sampleIDs.Get().(*[]uint64),
		Samples:          *p.samples.Get().(*[]float64),
		HistogramSamples: *p.histogramSamples.Get().(*[]*histogram.FloatHistogram),
	}
}

func (p *VectorPool) PutStepVector(v StepVector) {
	v.SampleIDs = v.SampleIDs[:0]
	v.Samples = v.Samples[:0]
	v.HistogramSamples = v.HistogramSamples[:0]
	p.sampleIDs.Put(&v.SampleIDs)
	p.samples.Put(&v.Samples)
	p.histogramSamples.Put(&v.HistogramSamples)
}

func (p *VectorPool) SetStepSize(n int) {
	p.stepSize = n
}
