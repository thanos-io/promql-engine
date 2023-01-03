package api

import (
	"time"

	"github.com/prometheus/prometheus/promql"
)

type RemoteEndpoints interface {
	Engines() []RemoteEngine
}

type RemoteEngine interface {
	NewInstantQuery(opts *promql.QueryOpts, qs string, ts time.Time) (promql.Query, error)
	NewRangeQuery(opts *promql.QueryOpts, qs string, start, end time.Time, interval time.Duration) (promql.Query, error)
}
