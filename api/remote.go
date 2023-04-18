// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package api

import (
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
)

type Opts struct {
	Query         string
	Start         time.Time
	End           time.Time
	Step          time.Duration
	LookbackDelta time.Duration
}

type RemoteEndpoints interface {
	Engines(opts *Opts) []RemoteEngine
}

type RemoteEngine interface {
	MaxT() int64
	MinT() int64
	LabelSets() []labels.Labels
	NewRangeQuery(opts *promql.QueryOpts, qs string, start, end time.Time, interval time.Duration) (promql.Query, error)
}

type staticEndpoints struct {
	engines []RemoteEngine
}

func (m staticEndpoints) Engines(opts *Opts) []RemoteEngine {
	return m.engines
}

func NewStaticEndpoints(engines []RemoteEngine) RemoteEndpoints {
	return &staticEndpoints{engines: engines}
}
