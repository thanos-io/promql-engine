package api

import (
	"github.com/prometheus/prometheus/promql"
	"time"
)

type ThanosEngine interface {
	NewInstantQuery(opts *promql.QueryOpts, qs string, ts time.Time) (promql.Query, error)
	NewRangeQuery(opts *promql.QueryOpts, qs string, start, end time.Time, interval time.Duration) (promql.Query, error)
}
