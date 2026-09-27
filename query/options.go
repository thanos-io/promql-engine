// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package query

import (
	"sync"
	"time"
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
	SampleTracker            SampleTracker // Tracks current samples in memory
	// Workers tracks goroutines started by operators of the query, so that
	// Exec can wait for them before returning and storage is not closed
	// while they still read from it.
	Workers *sync.WaitGroup
}

// Go runs fn in a new goroutine that is tracked by Workers, if set.
func (o *Options) Go(fn func()) {
	if o.Workers == nil {
		go fn()
		return
	}
	o.Workers.Add(1)
	go func() {
		defer o.Workers.Done()
		fn()
	}()
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

func NestedOptionsForSubquery(opts *Options, step, queryRange, offset time.Duration) *Options {
	nOpts := &Options{
		End:                      opts.End.Add(-offset),
		LookbackDelta:            opts.LookbackDelta,
		StepsBatch:               opts.StepsBatch,
		ExtLookbackDelta:         opts.ExtLookbackDelta,
		NoStepSubqueryIntervalFn: opts.NoStepSubqueryIntervalFn,
		EnableAnalysis:           opts.EnableAnalysis,
		DecodingConcurrency:      opts.DecodingConcurrency,
		SampleTracker:            opts.SampleTracker,
		Workers:                  opts.Workers,
	}
	if nOpts.SampleTracker == nil {
		nOpts.SampleTracker = NewSampleTracker(0)
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
