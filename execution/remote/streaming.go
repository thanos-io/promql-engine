// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package remote

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/thanos-io/promql-engine/api"
	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/execution/telemetry"
	"github.com/thanos-io/promql-engine/query"

	"github.com/prometheus/prometheus/model/labels"
)

// StreamingExecution is a model.VectorOperator that consumes an
// api.RemoteOperator end-to-end as a columnar StepVector stream.
//
// Contrast with Execution (operator.go), which executes a promql.Query
// up-front, materialises a promql.Matrix, and then re-transposes it back
// into StepVectors via a wrapped vectorSelector. StreamingExecution avoids
// both transposes: the remote engine produces StepVectors directly.
//
// Step-grid alignment: the parent query covers [opts.Start, opts.End] at
// opts.Step. The remote engine may have been scoped by the distributed
// optimizer to a narrower range (queryRangeStart..queryRangeEnd) because
// its data only covers part of the parent range. To keep downstream
// operators (notably exchange.coalesce, which zips per-step across multiple
// remote operators) correctly aligned, StreamingExecution pads the front
// and back of the remote's output with empty StepVectors at every step in
// the parent's grid that the remote did not produce.
//
// Lifecycle:
//   - Series(ctx) is cached on first call. It must be called (or implicitly
//     invoked from Next) before SampleIDs in produced StepVectors can be
//     interpreted by downstream operators.
//   - Close on the underlying api.RemoteOperator is invoked exactly once,
//     when Next returns n == 0 (end of stream). This mirrors the lifetime
//     handling in Execution: closing earlier can race with buffer recycling.
type StreamingExecution struct {
	remote          api.RemoteOperator
	opts            *query.Options
	engineLabels    []labels.Labels
	queryRangeStart time.Time
	queryRangeEnd   time.Time

	seriesOnce sync.Once
	series     []labels.Labels
	seriesErr  error

	// currentStep is the next timestamp the operator will emit (parent step grid).
	currentStep int64
	// remoteStartTS is the first timestamp the remote engine will produce.
	remoteStartTS int64
	// remoteEndTS is the timestamp after which the remote produces nothing.
	remoteEndTS int64
	// stepMs is the parent step in milliseconds.
	stepMs int64
	// endMs is the parent end timestamp in milliseconds (inclusive).
	endMs int64

	// pending buffers StepVectors fetched from the remote that have not yet
	// been delivered to the caller. We pull from the remote into our own
	// buffer (sized to opts.StepsBatch) and then hand vectors to the caller
	// one Next call at a time. This avoids passing a sliced buf to the
	// remote.Next, which would interact badly with the wrapping
	// exchange.concurrencyOperator (it silently drops vectors beyond
	// len(buf)).
	pending    []model.StepVector
	pendingIdx int
	remoteDone bool

	closeOnce sync.Once
}

// NewStreamingExecution constructs a StreamingExecution operator wrapped in
// the standard telemetry decorator.
func NewStreamingExecution(remote api.RemoteOperator, queryRangeStart, queryRangeEnd time.Time, engineLabels []labels.Labels, opts *query.Options) model.VectorOperator {
	op := &StreamingExecution{
		remote:          remote,
		opts:            opts,
		engineLabels:    engineLabels,
		queryRangeStart: queryRangeStart,
		queryRangeEnd:   queryRangeEnd,
		currentStep:     opts.Start.UnixMilli(),
		remoteStartTS:   queryRangeStart.UnixMilli(),
		remoteEndTS:     queryRangeEnd.UnixMilli(),
		stepMs:          stepMillis(opts.Step),
		endMs:           opts.End.UnixMilli(),
	}
	return telemetry.NewOperator(telemetry.NewTelemetry(op, opts), op)
}

// stepMillis returns step in milliseconds, treating zero (instant query) as 1.
// This mirrors how the engine treats instant queries internally — every
// timestamp comparison still uses a non-zero step so we never divide by zero
// when computing the next step.
func stepMillis(step time.Duration) int64 {
	if step <= 0 {
		return 1
	}
	return step.Milliseconds()
}

func (e *StreamingExecution) Series(ctx context.Context) ([]labels.Labels, error) {
	e.seriesOnce.Do(func() {
		e.series, e.seriesErr = e.remote.Series(ctx)
		// Some inner operators (pi(), time(), scalar(), and other scalar-returning
		// functions) intentionally report Series() == nil while still emitting
		// StepVectors with SampleID=0. compatibilityQuery.Exec adapts to that
		// convention by materialising a series slice from the first non-empty
		// Next call. Coalesce, dedup, and other downstream operators do not
		// adapt that way, so we normalise here: present a single empty-label
		// series instead, matching the implicit series that the scalar samples
		// belong to. If the remote turns out to emit no samples at all,
		// downstream operators see one declared series with no points, which
		// is harmless.
		if e.seriesErr == nil && len(e.series) == 0 {
			e.series = []labels.Labels{labels.EmptyLabels()}
		}
	})
	return e.series, e.seriesErr
}

func (e *StreamingExecution) Next(ctx context.Context, buf []model.StepVector) (int, error) {
	// Make sure Series has been resolved so downstream operators that build
	// output tables from Series before consuming Next see a consistent labelset.
	if _, err := e.Series(ctx); err != nil {
		return 0, err
	}

	n := 0
	for n < len(buf) {
		// End of parent grid: stop.
		if e.currentStep > e.endMs {
			break
		}

		// Phase 1: empty leading vectors below the remote's start.
		if e.currentStep < e.remoteStartTS {
			buf[n].Reset(e.currentStep)
			n++
			e.currentStep += e.stepMs
			continue
		}

		// Phase 2: forward StepVectors from the remote.
		if !e.remoteDone && e.currentStep <= e.remoteEndTS {
			if e.pendingIdx >= len(e.pending) {
				// Refill pending from the remote.
				if err := e.fillPending(ctx); err != nil {
					return 0, err
				}
				if e.remoteDone && len(e.pending) == 0 {
					continue
				}
			}
			if e.pendingIdx < len(e.pending) {
				// Swap the slice contents instead of copying samples.
				buf[n], e.pending[e.pendingIdx] = e.pending[e.pendingIdx], buf[n]
				e.currentStep = buf[n].T + e.stepMs
				n++
				e.pendingIdx++
				continue
			}
		}

		// Phase 3: empty trailing vectors after the remote's end.
		if e.currentStep <= e.endMs {
			buf[n].Reset(e.currentStep)
			n++
			e.currentStep += e.stepMs
			continue
		}

		break
	}

	if n == 0 {
		e.closeRemote()
	}
	return n, nil
}

// fillPending pulls one batch of StepVectors from the remote into e.pending.
// On EOF, sets e.remoteDone and leaves e.pending empty. Always resets
// pendingIdx to 0.
func (e *StreamingExecution) fillPending(ctx context.Context) error {
	// Restore e.pending to full capacity so we always pass a full-size buffer
	// to the remote. We must avoid passing a sub-slice because the wrapping
	// exchange.concurrencyOperator silently drops vectors past len(buf).
	if cap(e.pending) == 0 {
		size := e.opts.StepsBatch
		if size <= 0 {
			size = 1
		}
		e.pending = make([]model.StepVector, size)
	} else {
		e.pending = e.pending[:cap(e.pending)]
	}
	// Reset slice headers for safety; the producer will overwrite, but the
	// caller may have swapped previously-emitted vectors into these slots.
	e.pendingIdx = 0
	for i := range e.pending {
		e.pending[i] = model.StepVector{}
	}

	got, err := e.remote.Next(ctx, e.pending)
	if err != nil {
		e.closeRemote()
		return err
	}
	if got == 0 {
		e.remoteDone = true
		e.closeRemote()
		e.pending = e.pending[:0]
		return nil
	}
	e.pending = e.pending[:got]
	return nil
}

func (e *StreamingExecution) Explain() []model.VectorOperator {
	return nil
}

func (e *StreamingExecution) String() string {
	return fmt.Sprintf("[remoteStreamingExec] (%d, %d)", e.queryRangeStart.Unix(), e.queryRangeEnd.Unix())
}

func (e *StreamingExecution) closeRemote() {
	e.closeOnce.Do(func() {
		_ = e.remote.Close()
	})
}
