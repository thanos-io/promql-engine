// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package model

import (
	"sync"
)

type VectorPool struct {
	vectors sync.Pool

	stepSize  int
	samples   sync.Pool
	sampleIDs sync.Pool
}

func NewVectorPool(stepsBatch int) *VectorPool {
	pool := &VectorPool{}
	pool.vectors = sync.Pool{
		New: func() any {
			return make([]StepVector, 0, stepsBatch)
		},
	}
	pool.samples = sync.Pool{
		New: func() any {
			return make([]float64, 0, pool.stepSize)
		},
	}
	pool.sampleIDs = sync.Pool{
		New: func() any {
			return make([]uint64, 0, pool.stepSize)
		},
	}

	return pool
}

func (p *VectorPool) GetVectorBatch() []StepVector {
	return p.vectors.Get().([]StepVector)
}

func (p *VectorPool) PutVectors(vector []StepVector) {
	p.vectors.Put(vector[:0])
}

func (p *VectorPool) GetStepVector(t int64) StepVector {
	return StepVector{
		T:         t,
		SampleIDs: p.sampleIDs.Get().([]uint64),
		Samples:   p.samples.Get().([]float64),
	}
}

func (p *VectorPool) PutStepVector(v StepVector) {
	p.sampleIDs.Put(v.SampleIDs[:0])
	p.samples.Put(v.Samples[:0])
}

func (p *VectorPool) SetStepSize(n int) {
	p.stepSize = n
}
