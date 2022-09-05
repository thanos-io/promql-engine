package points

import (
	"sync"

	"github.com/prometheus/prometheus/promql"
)

type Pool struct {
	pool sync.Pool
}

func NewPool() *Pool {
	return &Pool{
		pool: sync.Pool{
			New: func() any {
				return make(promql.Vector, 0)
			},
		},
	}
}

func (p *Pool) Get() promql.Vector {
	return p.pool.Get().(promql.Vector)
}

func (p *Pool) Put(point promql.Vector) {
	p.pool.Put(point)
}
