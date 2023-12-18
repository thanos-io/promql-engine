// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package scan

import (
	"context"
	"fmt"
	"sync"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"

	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/extlabels"
	"github.com/thanos-io/promql-engine/query"
	"github.com/thanos-io/promql-engine/ringbuffer"
)

// TODO: only instant subqueries for now.
type subqueryOperator struct {
	next        model.VectorOperator
	pool        *model.VectorPool
	call        FunctionCall
	mint        int64
	maxt        int64
	currentStep int64
	step        int64

	funcExpr *parser.Call
	subQuery *parser.SubqueryExpr

	onceSeries sync.Once
	series     []labels.Labels

	lastVectors   []model.StepVector
	lastCollected int
	buffers       []*ringbuffer.RingBuffer[Value]
}

func NewSubqueryOperator(pool *model.VectorPool, next model.VectorOperator, opts *query.Options, funcExpr *parser.Call, subQuery *parser.SubqueryExpr) (model.VectorOperator, error) {
	call, err := NewRangeVectorFunc(funcExpr.Func.Name)
	if err != nil {
		return nil, err
	}
	return &subqueryOperator{
		next:          next,
		call:          call,
		pool:          pool,
		funcExpr:      funcExpr,
		subQuery:      subQuery,
		mint:          opts.Start.UnixMilli(),
		maxt:          opts.End.UnixMilli(),
		currentStep:   opts.Start.UnixMilli(),
		step:          opts.Step.Milliseconds(),
		lastCollected: -1,
	}, nil
}

func (o *subqueryOperator) Explain() (me string, next []model.VectorOperator) {
	return fmt.Sprintf("[*subqueryOperator] %v()", o.funcExpr.Func.Name), []model.VectorOperator{o.next}
}

func (o *subqueryOperator) GetPool() *model.VectorPool { return o.pool }

func (o *subqueryOperator) Next(ctx context.Context) ([]model.StepVector, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	if o.currentStep > o.maxt {
		return nil, nil
	}
	if err := o.initSeries(ctx); err != nil {
		return nil, err
	}

	mint := o.currentStep - o.subQuery.Range.Milliseconds()
	for _, b := range o.buffers {
		b.DropBefore(mint)
	}
	if len(o.lastVectors) > 0 {
		for _, v := range o.lastVectors[o.lastCollected+1:] {
			if v.T > o.currentStep {
				break
			}
			o.collect(v)
			o.lastCollected++
		}
		if o.lastCollected == len(o.lastVectors)-1 {
			o.next.GetPool().PutVectors(o.lastVectors)
			o.lastVectors = nil
			o.lastCollected = -1
		}
	}

ACC:
	for len(o.lastVectors) == 0 {
		vectors, err := o.next.Next(ctx)
		if err != nil {
			return nil, err
		}
		if len(vectors) == 0 {
			break ACC
		}
		for i, vector := range vectors {
			if vector.T > o.currentStep {
				o.lastVectors = vectors
				break ACC
			}
			o.collect(vector)
			o.lastCollected = i
		}
		o.next.GetPool().PutVectors(vectors)
	}

	res := o.pool.GetVectorBatch()
	sv := o.pool.GetStepVector(o.currentStep)
	for sampleId, rangeSamples := range o.buffers {
		f, h, ok := o.call(FunctionArgs{
			Samples:     rangeSamples.Samples(),
			StepTime:    o.currentStep,
			SelectRange: o.subQuery.Range.Milliseconds(),
			Offset:      o.subQuery.Offset.Milliseconds(),
		})
		if ok {
			if h != nil {
				sv.AppendHistogram(o.pool, uint64(sampleId), h)
			} else {
				sv.AppendSample(o.pool, uint64(sampleId), f)
			}
		}
	}
	res = append(res, sv)

	o.currentStep += o.step
	return res, nil
}

func (o *subqueryOperator) collect(v model.StepVector) {
	for i, s := range v.Samples {
		buffer := o.buffers[v.SampleIDs[i]]
		if buffer.Len() > 0 && v.T <= buffer.MaxT() {
			continue
		}
		buffer.Push(v.T, Value{F: s})
	}
	for i, s := range v.Histograms {
		buffer := o.buffers[v.SampleIDs[i]]
		if buffer.Len() > 0 && v.T < buffer.MaxT() {
			continue
		}
		buffer.Push(v.T, Value{H: s})
	}
	o.next.GetPool().PutStepVector(v)
}

func (o *subqueryOperator) Series(ctx context.Context) ([]labels.Labels, error) {
	if err := o.initSeries(ctx); err != nil {
		return nil, err
	}
	return o.series, nil
}

func (o *subqueryOperator) initSeries(ctx context.Context) error {
	var err error
	o.onceSeries.Do(func() {
		var series []labels.Labels
		series, err = o.next.Series(ctx)
		if err != nil {
			return
		}

		o.series = make([]labels.Labels, len(series))
		o.buffers = make([]*ringbuffer.RingBuffer[Value], len(series))
		for i := range o.buffers {
			o.buffers[i] = ringbuffer.New[Value](8)
		}
		var b labels.ScratchBuilder
		for i, s := range series {
			lbls := s
			if o.funcExpr.Func.Name != "last_over_time" {
				lbls, _ = extlabels.DropMetricName(s, b)
			}
			o.series[i] = lbls
		}
	})
	return err
}
