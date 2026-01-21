// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package api

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
)

type RemoteQuery interface {
	fmt.Stringer
}

type RemoteEndpoints interface {
	Engines() []RemoteEngine
}

// RemoteEndpointsV2 describes endpoints that accept pruning hints when
// selecting remote engines.
//
// For example implementations may use the hints to prune the TSDBInfos, but
// also may safely ignore them and return all available remote engines.
//
// NOTE(Aleksandr Krivoshchekov):
// We add a new interface as a temporary backward compatibility.
// RemoteEndpoints will be replaced with it in a future breaking change.
type RemoteEndpointsV2 interface {
	EnginesV2(mint, maxt int64) []RemoteEngine
}

type RemoteEngine interface {
	MaxT() int64
	MinT() int64

	// The external labels of the remote engine. These are used to limit fanout. The engine uses these to
	// not distribute into remote engines that would return empty responses because their labelset is not matching.
	LabelSets() []labels.Labels

	// The external labels of the remote engine that form a logical partition. This is expected to be
	// a subset of the result of "LabelSets()". The engine uses these to compute how to distribute a query.
	// It is important that, for a given set of remote engines, these labels do not overlap meaningfully.
	PartitionLabelSets() []labels.Labels

	NewRangeQuery(ctx context.Context, opts promql.QueryOpts, plan RemoteQuery, start, end time.Time, interval time.Duration) (promql.Query, error)
}

type staticEndpoints struct {
	engines []RemoteEngine
}

func (m staticEndpoints) Engines() []RemoteEngine {
	return m.engines
}

func (m staticEndpoints) EnginesV2(mint, maxt int64) []RemoteEngine {
	return m.engines
}

func NewStaticEndpoints(engines []RemoteEngine) RemoteEndpoints {
	return &staticEndpoints{engines: engines}
}

type cachedEndpoints struct {
	endpoints RemoteEndpoints

	enginesOnce sync.Once
	engines     []RemoteEngine
}

func (l *cachedEndpoints) Engines() []RemoteEngine {
	const mint, maxt = math.MinInt64, math.MaxInt64
	return l.EnginesV2(mint, maxt)
}

func (l *cachedEndpoints) EnginesV2(mint, maxt int64) []RemoteEngine {
	l.enginesOnce.Do(func() {
		l.engines = getEngines(l.endpoints, mint, maxt)
	})
	return l.engines
}

func getEngines(endpoints RemoteEndpoints, mint, maxt int64) []RemoteEngine {
	if v2, ok := endpoints.(RemoteEndpointsV2); ok {
		return v2.EnginesV2(mint, maxt)
	}

	return endpoints.Engines()
}

// NewCachedEndpoints returns an endpoints wrapper that
// resolves and caches engines on first access.
//
// All subsequent Engines calls return cached engines, ignoring any query
// parameters.
func NewCachedEndpoints(endpoints RemoteEndpoints) RemoteEndpoints {
	if endpoints == nil {
		panic("api.NewCachedEndpoints: endpoints is nil")
	}

	if le, ok := endpoints.(*cachedEndpoints); ok {
		return le
	}

	return &cachedEndpoints{endpoints: endpoints}
}
