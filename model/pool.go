package model

import (
	"sync"
)

type VectorPool struct {
	pool sync.Pool
}

func NewPool() *VectorPool {
	return &VectorPool{
		pool: sync.Pool{
			New: func() any {
				return make([]StepVector, 0, 30)
			},
		},
	}
}

func (p *VectorPool) Get() []StepVector {
	return p.pool.Get().([]StepVector)
}

func (p *VectorPool) Put(vector []StepVector) {
	p.pool.Put(vector[:0])
}
