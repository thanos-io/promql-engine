// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package exchange

import (
	"context"

	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/execution/telemetry"
	"github.com/thanos-io/promql-engine/query"

	"github.com/prometheus/prometheus/model/labels"
)

// materializeOperator runs its child to completion, buffers the output,
// then replays it. The child's iterators become GC-eligible after the
// drain, reducing peak memory in join queries.
type materializeOperator struct {
	child model.VectorOperator

	materialized []model.StepVector
	cursor       int
	done         bool
	batchSize    int

	series []labels.Labels
	opts   *query.Options
}

// NewMaterialize wraps a child operator so that the first Next() call
// drains it entirely. Subsequent calls replay from the buffer.
func NewMaterialize(child model.VectorOperator, batchSize int, opts *query.Options) model.VectorOperator {
	m := &materializeOperator{
		child:     child,
		batchSize: batchSize,
		opts:      opts,
	}
	return telemetry.NewOperator(telemetry.NewTelemetry(m, opts), m)
}

func (m *materializeOperator) String() string {
	return "[materialize]"
}

func (m *materializeOperator) Explain() []model.VectorOperator {
	if m.child == nil {
		return nil
	}
	return []model.VectorOperator{m.child}
}
func (m *materializeOperator) Series(ctx context.Context) ([]labels.Labels, error) {
	if m.series != nil {
		return m.series, nil
	}
	series, err := m.child.Series(ctx)
	if err != nil {
		return nil, err
	}
	m.series = series
	return m.series, nil
}

func (m *materializeOperator) Next(ctx context.Context, buf []model.StepVector) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	if !m.done {
		if err := m.materialize(ctx); err != nil {
			return 0, err
		}
	}

	if m.cursor >= len(m.materialized) {
		return 0, nil
	}

	n := 0
	for n < len(buf) && m.cursor < len(m.materialized) {
		buf[n] = m.materialized[m.cursor]
		released := len(m.materialized[m.cursor].SampleIDs) + len(m.materialized[m.cursor].HistogramIDs)
		m.opts.SampleTracker.Remove(released)
		m.materialized[m.cursor] = model.StepVector{}
		m.cursor++
		n++
	}

	if m.cursor >= len(m.materialized) {
		m.materialized = nil
	}

	return n, nil
}

func (m *materializeOperator) materialize(ctx context.Context) error {
	seriesCount := len(m.series)
	totalSteps := m.opts.TotalSteps()

	// Pre-allocate all step vectors upfront. The child writes directly
	// into these via tmpBuf slices, avoiding per-step allocations.
	m.materialized = make([]model.StepVector, totalSteps)
	for i := range m.materialized {
		m.materialized[i].SampleIDs = make([]uint64, 0, seriesCount)
		m.materialized[i].Samples = make([]float64, 0, seriesCount)
	}

	cursor := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		end := min(cursor+m.batchSize, totalSteps)
		if cursor >= end {
			break
		}
		tmpBuf := m.materialized[cursor:end]

		n, err := m.child.Next(ctx, tmpBuf)
		if err != nil {
			return err
		}
		if n == 0 {
			break
		}

		// Track buffered samples and enforce maxSamples limit.
		for i := range n {
			sv := m.materialized[cursor+i]
			sampleCount := len(sv.SampleIDs) + len(sv.HistogramIDs)
			m.opts.SampleTracker.Add(sampleCount)
		}
		if err := m.opts.SampleTracker.CheckLimit(); err != nil {
			return err
		}
		cursor += n
	}

	m.materialized = m.materialized[:cursor]
	m.done = true
	m.child = nil

	return nil
}
