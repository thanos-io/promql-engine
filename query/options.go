// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package query

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/baggage"
)

type Options struct {
	Start                    time.Time
	End                      time.Time
	Step                     time.Duration
	StepsBatch               int
	LookbackDelta            time.Duration
	EnablePerStepStats       bool
	ExtLookbackDelta         time.Duration
	NoStepSubqueryIntervalFn func(time.Duration) time.Duration
	EnableAnalysis           bool
	DecodingConcurrency      int
}

// TotalSteps returns the total number of steps in the query, regardless of batching.
// This is useful for pre-allocating result slices.
func (o *Options) TotalSteps() int {
	// Instant evaluation is executed as a range evaluation with one step.
	if o.Step.Milliseconds() == 0 {
		return 1
	}
	return int((o.End.UnixMilli()-o.Start.UnixMilli())/o.Step.Milliseconds() + 1)
}

func (o *Options) NumStepsPerBatch() int {
	totalSteps := o.TotalSteps()
	if o.StepsBatch < totalSteps {
		return o.StepsBatch
	}
	return totalSteps
}

func (o *Options) IsInstantQuery() bool {
	return o.TotalSteps() == 1
}

func (o *Options) WithEndTime(end time.Time) *Options {
	result := *o
	result.End = end
	return &result
}

// ShouldEnableAnalysis returns true if analysis should be enabled for this
// query. The baggage key "promql.enable_analysis" acts as a per-query
// override: "true" forces analysis on, "false" forces it off, and any other
// value (including absent) defers to o.EnableAnalysis.
func (o *Options) ShouldEnableAnalysis(ctx context.Context) bool {
	if v := baggage.FromContext(ctx).Member("promql.enable_analysis").Value(); v != "" {
		return v == "true"
	}
	return o.EnableAnalysis
}

func NestedOptionsForSubquery(opts *Options, step, queryRange, offset time.Duration) *Options {
	nOpts := &Options{
		End:                      opts.End.Add(-offset),
		LookbackDelta:            opts.LookbackDelta,
		StepsBatch:               opts.StepsBatch,
		ExtLookbackDelta:         opts.ExtLookbackDelta,
		NoStepSubqueryIntervalFn: opts.NoStepSubqueryIntervalFn,
		EnableAnalysis:           opts.EnableAnalysis,
		DecodingConcurrency:      opts.DecodingConcurrency,
	}
	if step != 0 {
		nOpts.Step = step
	} else {
		nOpts.Step = opts.NoStepSubqueryIntervalFn(queryRange)
	}
	nOpts.Start = time.UnixMilli(nOpts.Step.Milliseconds() * (opts.Start.Add(-offset-queryRange).UnixMilli() / nOpts.Step.Milliseconds()))
	if nOpts.Start.Before(opts.Start.Add(-offset - queryRange)) {
		nOpts.Start = nOpts.Start.Add(nOpts.Step)
	}
	return nOpts
}
