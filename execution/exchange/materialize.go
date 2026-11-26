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

// materializeOperator drains its child, buffers the output, then replays it.
type materializeOperator struct {
	child model.VectorOperator

	materialized []model.StepVector
	cursor       int
	done         bool
	// bufferedSamples is what our own buffer counts towards the sample limit.
	// It lives on after the child is freed, so Next subtracts it at the end.
	bufferedSamples int

	series  []labels.Labels
	tempBuf []model.StepVector
	opts    *query.Options
}

// NewMaterialize wraps child so the first Next drains it and later calls replay.
func NewMaterialize(child model.VectorOperator, opts *query.Options) model.VectorOperator {
	m := &materializeOperator{
		child: child,
		opts:  opts,
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
		// Swap, not assign: drops our reference so the consumer's buffer is freed.
		buf[n], m.materialized[m.cursor] = m.materialized[m.cursor], buf[n]
		m.cursor++
		n++
	}

	if m.cursor >= len(m.materialized) {
		m.opts.SampleTracker.Remove(m.bufferedSamples)
		m.materialized = nil
	}

	return n, nil
}

func (m *materializeOperator) materialize(ctx context.Context) error {
	totalSteps := m.opts.TotalSteps()

	samplesBefore := m.opts.SampleTracker.Current()

	m.materialized = make([]model.StepVector, totalSteps)
	m.tempBuf = make([]model.StepVector, m.opts.StepsBatch)

	threshold := query.ComputeSampleLimitCheckThreshold(m.opts, 1)
	untracked := 0

	totalSamples := 0
	cursor := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := m.child.Next(ctx, m.tempBuf)
		if err != nil {
			return err
		}
		if n == 0 {
			break
		}

		batchSamples := 0
		for i := 0; i < n && cursor < totalSteps; i++ {
			// Swap: slice headers move, sample data does not.
			m.materialized[cursor], m.tempBuf[i] = m.tempBuf[i], m.materialized[cursor]
			batchSamples += stepVectorSamples(&m.materialized[cursor])
			cursor++
		}
		totalSamples += batchSamples
		untracked += batchSamples

		if untracked >= threshold {
			m.opts.SampleTracker.Add(untracked)
			untracked = 0
			if err := m.opts.SampleTracker.CheckLimit(); err != nil {
				return err
			}
		}
	}

	// Count the rest before measuring the child below.
	if untracked > 0 {
		m.opts.SampleTracker.Add(untracked)
		if err := m.opts.SampleTracker.CheckLimit(); err != nil {
			return err
		}
	}

	m.bufferedSamples = totalSamples
	m.materialized = m.materialized[:cursor]
	m.tempBuf = nil
	m.done = true

	model.ReleaseSubtree(m.child)

	// Subtract what the child used, but not our own buffer. Can be negative.
	subtreeSamples := m.opts.SampleTracker.Current() - samplesBefore - int64(totalSamples)
	if subtreeSamples > 0 {
		m.opts.SampleTracker.Remove(int(subtreeSamples))
	}

	return nil
}

// stepVectorSamples counts a step vector's samples, weighting histograms by size.
func stepVectorSamples(sv *model.StepVector) int {
	n := len(sv.SampleIDs)
	for _, h := range sv.Histograms {
		n += telemetry.CalculateHistogramSampleCount(h)
	}
	return n
}
