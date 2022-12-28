package remote

import (
	"github.com/prometheus/prometheus/model/histogram"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
)

type samplesIterator struct {
	i      int
	points []promql.Point
}

func newSamplesIterator(points []promql.Point) *samplesIterator {
	return &samplesIterator{
		i:      -1,
		points: points,
	}
}

func (p *samplesIterator) Next() chunkenc.ValueType {
	p.i++
	if p.i < len(p.points) {
		return chunkenc.ValFloat
	}
	return chunkenc.ValNone
}

func (p *samplesIterator) Seek(t int64) chunkenc.ValueType {
	for {
		if p.AtT() >= t {
			return chunkenc.ValFloat
		}
		if p.Next() == chunkenc.ValNone {
			return chunkenc.ValNone
		}
	}
}

func (p *samplesIterator) At() (int64, float64) {
	return p.points[p.i].T, p.points[p.i].V
}

func (p *samplesIterator) AtHistogram() (int64, *histogram.Histogram) {
	return 0, nil
}

func (p *samplesIterator) AtFloatHistogram() (int64, *histogram.FloatHistogram) {
	return 0, nil
}

func (p *samplesIterator) AtT() int64 {
	return p.points[p.i].T
}

func (p *samplesIterator) Err() error { return nil }
