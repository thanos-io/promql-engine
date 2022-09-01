package points

import (
	"sync"

	"github.com/prometheus/prometheus/promql"
)

type points struct {
	pool sync.Pool
}

func NewPool() *points {
	return &points{
		pool: sync.Pool{
			New: func() any {
				return &promql.Point{}
			},
		},
	}
}

func (p *points) get() *promql.Point {
	return p.pool.Get().(*promql.Point)
}

func (p *points) put(point *promql.Point) {
	p.pool.Put(point)
}
