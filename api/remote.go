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

// Deprecated: RemoteEndpoints will be replaced with
// RemoteEndpointsV2 / RemoteEndpointsV3 in a future breaking change.
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

type RemoteEndpointsQuery struct {
	MinT int64
	MaxT int64
}

// RemoteEndpointsV3 describes endpoints that accept pruning hints when
// selecting remote engines.
//
// For example implementations may use the hints to prune the TSDBInfos, but
// also may safely ignore them and return all available remote engines.
//
// NOTE(Aleksandr Krivoshchekov):
// We add a new interface as a temporary backward compatibility.
// RemoteEndpoints will be replaced with it in a future breaking change.
//
// Unlike RemoteEndpointsV2, this interface can be extended with more hints
// in the future, without making any breaking changes.
type RemoteEndpointsV3 interface {
	EnginesV3(query RemoteEndpointsQuery) []RemoteEngine
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

func (m staticEndpoints) EnginesV3(query RemoteEndpointsQuery) []RemoteEngine {
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
	return l.EnginesV3(RemoteEndpointsQuery{
		MaxT: math.MinInt64,
		MinT: math.MaxInt64,
	})
}

func (l *cachedEndpoints) EnginesV2(mint, maxt int64) []RemoteEngine {
	return l.EnginesV3(RemoteEndpointsQuery{
		MaxT: maxt,
		MinT: mint,
	})
}

func (l *cachedEndpoints) EnginesV3(query RemoteEndpointsQuery) []RemoteEngine {
	l.enginesOnce.Do(func() {
		l.engines = getEngines(l.endpoints, query)
	})
	return l.engines
}

func getEngines(endpoints RemoteEndpoints, query RemoteEndpointsQuery) []RemoteEngine {
	if v3, ok := endpoints.(RemoteEndpointsV3); ok {
		return v3.EnginesV3(query)
	}

	if v2, ok := endpoints.(RemoteEndpointsV2); ok {
		return v2.EnginesV2(query.MinT, query.MaxT)
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
