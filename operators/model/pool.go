package model

import (
	"sync"
)

type VectorPool struct {
	vectors sync.Pool

	numSamples int
	samples    sync.Pool
}

func NewVectorPool() *VectorPool {
	pool := &VectorPool{}
	pool.vectors = sync.Pool{
		New: func() any {
			return make([]StepVector, 0, 10)
		},
	}
	pool.samples = sync.Pool{
		New: func() any {
			return make([]StepSample, 0, pool.numSamples)
		},
	}

	return pool
}

func (p *VectorPool) GetVectors() []StepVector {
	return p.vectors.Get().([]StepVector)
}

func (p *VectorPool) PutVectors(vector []StepVector) {
	p.vectors.Put(vector[:0])
}

func (p *VectorPool) GetSamples() []StepSample {
	return p.samples.Get().([]StepSample)
}

func (p *VectorPool) PutSamples(samples []StepSample) {
	p.samples.Put(samples[:0])
}

func (p *VectorPool) SetStepSamplesSize(n int) {
	p.numSamples = n
}
