// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb/chunkenc"

	"github.com/thanos-community/promql-engine/engine"
)

func TestQueriesAgainstOldEngine(t *testing.T) {
	start := time.Unix(0, 0)
	end := time.Unix(240, 0)
	step := time.Second * 30
	opts := promql.EngineOpts{
		Timeout:    1 * time.Hour,
		MaxSamples: 1e10,
	}

	cases := []struct {
		load     string
		name     string
		query    string
		start    time.Time
		end      time.Time
		step     time.Duration
		expected []promql.Vector
	}{
		{
			name: "stddev_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "stddev_over_time(http_requests_total[30s])",
		},
		{
			name: "stdvar_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "stdvar_over_time(http_requests_total[30s])",
		},
		{
			name: "changes",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "changes(http_requests_total[30s])",
		},
		{
			name: "deriv",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "deriv(http_requests_total[30s])",
		},
		{
			name: "sum",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "sum (http_requests_total)",
		},
		{
			name: "sum_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "sum_over_time(http_requests_total[30s])",
		},
		{
			name: "count",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "count(http_requests_total)",
		},
		{
			name: "count_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "count_over_time(http_requests_total[30s])",
		},
		{
			name: "average",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "avg(http_requests_total)",
		},
		{
			name: "avg_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "avg_over_time(http_requests_total[30s])",
		},
		{
			name: "max",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "max(http_requests_total)",
		},
		{
			name: "max with only 1 sample",
			load: `load 30s
					http_requests_total{pod="nginx-1"} -1
					http_requests_total{pod="nginx-2"} 1`,
			query: "max(http_requests_total) by (pod)",
		},
		{
			name: "max_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "max_over_time(http_requests_total[30s])",
		},
		{
			name: "min",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "min(http_requests_total)",
		},
		{
			name: "min with only 1 sample",
			load: `load 30s
					http_requests_total{pod="nginx-1"} -1
					http_requests_total{pod="nginx-2"} 1`,
			query: "min(http_requests_total) by (pod)",
		},
		{
			name: "min_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "min_over_time(http_requests_total[30s])",
		},
		{
			name: "count_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "count_over_time(http_requests_total[30s])",
		},
		{
			name: "sum by pod",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18
					http_requests_total{pod="nginx-3"} 1+2x20
					http_requests_total{pod="nginx-4"} 1+2x50`,
			query: "sum by (pod) (http_requests_total)",
		},
		{
			name: "query in the future",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "sum by (pod) (http_requests_total)",
			start: time.Unix(400, 0),
			end:   time.Unix(3000, 0),
		},
		{
			name: "rate",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "rate(http_requests_total[1m])",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "sum rate",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x4
					http_requests_total{pod="nginx-2"} 1+2x4`,
			query: "sum(rate(http_requests_total[1m]))",
		},
		{
			name: "sum rate with stale series",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x40
					http_requests_total{pod="nginx-2"} 1+2x50
					http_requests_total{pod="nginx-4"} 1+2x50
					http_requests_total{pod="nginx-5"} 1+2x50
					http_requests_total{pod="nginx-6"} 1+2x50`,
			query: "sum(rate(http_requests_total[1m]))",
			start: time.Unix(421, 0),
			end:   time.Unix(3230, 0),
			step:  28 * time.Second,
		},
		{
			name: "delta",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "delta(http_requests_total[1m])",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "increase",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "increase(http_requests_total[1m])",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "irate",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "irate(http_requests_total[1m])",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "idelta",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "idelta(http_requests_total[1m])",
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name:  "number literal",
			load:  "",
			query: "34",
		},
		{
			name:  "vector",
			load:  "",
			query: "vector(24)",
		},
		{
			name: "binary operation with one-to-one matching",
			load: `load 30s
				foo{method="get", code="500"} 1+1x1
				foo{method="get", code="404"} 2+1x2
				foo{method="put", code="501"} 3+1x3
				foo{method="put", code="500"} 1+1x4
				foo{method="post", code="500"} 4+1x4
				foo{method="post", code="404"} 5+1x5
				bar{method="get"} 1+1x1
				bar{method="del"} 2+1x2  
				bar{method="post"} 3+1x3`,
			query: `foo{code="500"} + ignoring(code) bar`,
			start: time.Unix(0, 0),
			end:   time.Unix(600, 0),
		},
		{
			// Example from https://prometheus.io/docs/prometheus/latest/querying/operators/#many-to-one-and-one-to-many-vector-matches
			name: "binary operation with group_left",
			load: `load 30s
				foo{method="get", code="500", path="/"} 1+1.1x30
				foo{method="get", code="404", path="/"} 1+2.2x20
				foo{method="put", code="501", path="/"} 4+3.4x60
				foo{method="post", code="500", path="/"} 1+5.1x40
				foo{method="post", code="404", path="/"} 2+3.7x40
				bar{method="get", path="/a"} 3+7.4x10
				bar{method="del", path="/b"} 8+6.1x30  
				bar{method="post", path="/c"} 1+2.1x40`,
			query: `foo * ignoring(code, path) group_left bar`,
			start: time.Unix(0, 0),
			end:   time.Unix(600, 0),
		},
		{
			// Example from https://prometheus.io/docs/prometheus/latest/querying/operators/#many-to-one-and-one-to-many-vector-matches
			name: "binary operation with group_right",
			load: `load 30s
				foo{method="get", code="500"} 1+1.1x30
				foo{method="get", code="404"} 1+2.2x20
				foo{method="put", code="501"} 4+3.4x60
				foo{method="post", code="500"} 1+5.1x40
				foo{method="post", code="404"} 2+3.7x40
				bar{method="get", path="/a"} 3+7.4x10
				bar{method="del", path="/b"} 8+6.1x30  
				bar{method="post", path="/c"} 1+2.1x40`,
			query: `bar * ignoring(code, path) group_right foo`,
			start: time.Unix(0, 0),
			end:   time.Unix(3000, 0),
			step:  2 * time.Second,
		},
		{
			name: "binary operation with group_left and included labels",
			load: `load 30s
				foo{method="get", code="500"} 1+1.1x30
				foo{method="get", code="404"} 1+2.2x20
				foo{method="put", code="501"} 4+3.4x60
				foo{method="post", code="500"} 1+5.1x40
				foo{method="post", code="404"} 2+3.7x40
				bar{method="get", path="/a"} 3+7.4x10
				bar{method="del", path="/b"} 8+6.1x30
				bar{method="post", path="/c"} 1+2.1x40`,
			query: `foo * ignoring(code, path) group_left(path) bar`,
			start: time.Unix(0, 0),
			end:   time.Unix(600, 0),
		},
		{
			name: "binary operation with group_right and included labels",
			load: `load 30s
				foo{method="get", code="500"} 1+1.1x30
				foo{method="get", code="404"} 1+2.2x20
				foo{method="put", code="501"} 4+3.4x60
				foo{method="post", code="500"} 1+5.1x40
				foo{method="post", code="404"} 2+3.7x40
				bar{method="get", path="/a"} 3+7.4x10
				bar{method="del", path="/b"} 8+6.1x30
				bar{method="post", path="/c"} 1+2.1x40`,
			query: `bar * ignoring(code, path) group_right(path) foo`,
			start: time.Unix(0, 0),
			end:   time.Unix(600, 0),
		},
		{
			name: "binary operation with vector and scalar on the right",
			load: `load 30s
				foo{method="get", code="500"} 1+1.1x30
				foo{method="get", code="404"} 1+2.2x20`,
			query: `sum(foo) * 2`,
		},
		{
			name: "binary operation with vector and scalar on the left",
			load: `load 30s
				foo{method="get", code="500"} 1+1.1x30
				foo{method="get", code="404"} 1+2.2x20`,
			query: `2 * sum(foo)`,
		},
		{
			name: "complex binary operation",
			load: `load 30s
				foo{method="get", code="500"} 1+1.1x30
				foo{method="get", code="404"} 1+2.2x20`,
			query: `1 - (100 * sum(foo{method="get"}) / sum(foo))`,
		},
		{
			name:  "empty series",
			load:  "",
			query: "http_requests_total",
		},
		{
			name:  "empty series with func",
			load:  "",
			query: "sum(http_requests_total)",
		},
		{
			name: "empty result",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total{pod="nginx-3"}`,
		},
		{
			name: "last_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `last_over_time(http_requests_total[30s])`,
		},
		{
			name: "group",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `group(http_requests_total)`,
		},
		{
			name: "resets",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 100-1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `resets(http_requests_total[5m])`,
		},
		{
			name: "present_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `present_over_time(http_requests_total[30s])`,
		},
		{
			name: "complex binary with aggregation",
			load: `load 30s
					grpc_server_handled_total{pod="nginx-1", grpc_method="Series", grpc_code="105"} 1+1x15
					grpc_server_handled_total{pod="nginx-2", grpc_method="Series", grpc_code="105"} 1+1x15
					grpc_server_handled_total{pod="nginx-3", grpc_method="Series", grpc_code="105"} 1+1x15
					prometheus_tsdb_head_samples_appended_total{pod="nginx-1", tenant="tenant-1"} 1+2x18
					prometheus_tsdb_head_samples_appended_total{pod="nginx-2", tenant="tenant-2"} 1+2x18
					prometheus_tsdb_head_samples_appended_total{pod="nginx-3", tenant="tenant-3"} 1+2x18`,
			query: `
	sum by (grpc_method, grpc_code) (
		sum by (pod, grpc_method, grpc_code) (
			rate(grpc_server_handled_total{grpc_method="Series", pod=~".+"}[1m])
		)
		+ on (pod) group_left() max by (pod) (
			prometheus_tsdb_head_samples_appended_total{pod=~".+"}
		)
	)`,
		},
	}

	lookbackDeltas := []time.Duration{30 * time.Second, time.Minute, 5 * time.Minute, 10 * time.Minute}
	for _, lookbackDelta := range lookbackDeltas {
		opts.LookbackDelta = lookbackDelta
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				test, err := promql.NewTest(t, tc.load)
				testutil.Ok(t, err)
				defer test.Close()

				testutil.Ok(t, test.Run())

				if tc.start.Equal(time.Time{}) {
					tc.start = start
				}
				if tc.end.Equal(time.Time{}) {
					tc.end = end
				}
				if tc.step == 0 {
					tc.step = step
				}

				for _, disableFallback := range []bool{false, true} {
					t.Run(fmt.Sprintf("disableFallback=%v", disableFallback), func(t *testing.T) {
						newEngine := engine.New(engine.Opts{EngineOpts: opts, DisableFallback: disableFallback})
						q1, err := newEngine.NewRangeQuery(test.Storage(), nil, tc.query, tc.start, tc.end, step)
						testutil.Ok(t, err)

						newResult := q1.Exec(context.Background())
						testutil.Ok(t, newResult.Err)

						oldEngine := promql.NewEngine(opts)
						q2, err := oldEngine.NewRangeQuery(test.Storage(), nil, tc.query, tc.start, tc.end, step)
						testutil.Ok(t, err)

						oldResult := q2.Exec(context.Background())
						testutil.Ok(t, oldResult.Err)

						testutil.Equals(t, oldResult, newResult)
					})
				}
			})
		}
	}
}

func TestInstantQuery(t *testing.T) {
	queryTime := time.Unix(50, 0)
	opts := promql.EngineOpts{
		Timeout:    1 * time.Hour,
		MaxSamples: 1e10,
	}

	cases := []struct {
		load     string
		name     string
		query    string
		expected []promql.Vector
	}{
		{
			name: "stdvar",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x4
					http_requests_total{pod="nginx-2"} 1+2x4`,
			query: "stdvar(http_requests_total)",
		},
		{
			name: "stdvar by pod",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1
					http_requests_total{pod="nginx-2"} 2
					http_requests_total{pod="nginx-3"} 8
					http_requests_total{pod="nginx-4"} 6`,
			query: "stdvar by (pod) (http_requests_total)",
		},
		{
			name: "stddev",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x4
					http_requests_total{pod="nginx-2"} 1+2x4`,
			query: "stddev(http_requests_total)",
		},
		{
			name: "stddev by pod",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1
					http_requests_total{pod="nginx-2"} 2
					http_requests_total{pod="nginx-3"} 8
					http_requests_total{pod="nginx-4"} 6`,
			query: "stddev by (pod) (http_requests_total)",
		},
		{
			name: "sum by pod",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x4
					http_requests_total{pod="nginx-2"} 1+2x4`,
			query: "sum by (pod) (http_requests_total)",
		},
		{
			name: "count",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "count(http_requests_total)",
		},
		{
			name: "average",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "avg(http_requests_total)",
		},
		{
			name: "max",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "max(http_requests_total)",
		},
		{
			name: "max with only 1 sample",
			load: `load 30s
					http_requests_total{pod="nginx-1"} -1
					http_requests_total{pod="nginx-2"} 1`,
			query: "max(http_requests_total) by (pod)",
		},
		{
			name: "min",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "min(http_requests_total)",
		},
		{
			name: "min with only 1 sample",
			load: `load 30s
					http_requests_total{pod="nginx-1"} -1
					http_requests_total{pod="nginx-2"} 1`,
			query: "min(http_requests_total) by (pod)",
		},
		{
			name: "rate",
			load: `load 30s
				http_requests_total{pod="nginx-1", series="1"} 1+1.1x40
				http_requests_total{pod="nginx-2", series="2"} 2+2.3x50
				http_requests_total{pod="nginx-4", series="3"} 5+2.4x50
				http_requests_total{pod="nginx-5", series="1"} 8.4+2.3x50
				http_requests_total{pod="nginx-6", series="2"} 2.3+2.3x50`,
			query: "rate(http_requests_total[1m])",
		},
		{
			name: "sum rate",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x4
					http_requests_total{pod="nginx-2"} 1+2x4`,
			query: "sum(rate(http_requests_total[1m]))",
		},
		{
			name: "sum rate with single sample series",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x4
					http_requests_total{pod="nginx-2"} 1+2x4
					http_requests_total{pod="nginx-3"} 0`,
			query: "sum by (pod) (rate(http_requests_total[1m]))",
		},
		{
			name: "sum rate with stale series",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x4
					http_requests_total{pod="nginx-2"} 1+2x20`,
			query: "sum(rate(http_requests_total[1m]))",
		},
		{
			name: "delta",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x4
					http_requests_total{pod="nginx-2"} 1+2x4`,
			query: "delta(http_requests_total[1m])",
		},
		{
			name: "increase",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x4
					http_requests_total{pod="nginx-2"} 1+2x4`,
			query: "increase(http_requests_total[1m])",
		},
		{
			name: "sum irate",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x4
					http_requests_total{pod="nginx-2"} 1+2x4`,
			query: "sum(irate(http_requests_total[1m]))",
		},
		{
			name: "sum irate with stale series",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x4
					http_requests_total{pod="nginx-2"} 1+2x20`,
			query: "sum(irate(http_requests_total[1m]))",
		},
		{
			name:  "string literal",
			load:  "",
			query: `"hello"`,
		},
		{
			name:  "number literal",
			load:  "",
			query: "34",
		},
		{
			name:  "vector",
			load:  "",
			query: "vector(24)",
		},
		{
			name: "binary operation with vector and scalar on the right",
			load: `load 30s
				foo{method="get", code="500"} 1+1.1x30
				foo{method="get", code="404"} 1+2.2x20`,
			query: `foo * 2`,
		},
		{
			name: "binary operation with vector and scalar on the left",
			load: `load 30s
				foo{method="get", code="500"} 1+1.1x30
				foo{method="get", code="404"} 1+2.2x20`,
			query: `2 * foo`,
		},
		{
			name: "complex binary operation",
			load: `load 30s
				foo{method="get", code="500"} 1+1.1x30
				foo{method="get", code="404"} 1+2.2x20`,
			query: `1 - (100 * sum(foo{method="get"}) / sum(foo))`,
		},
		{
			name:  "empty series",
			load:  "",
			query: "http_requests_total",
		},
		{
			name:  "empty series with func",
			load:  "",
			query: "sum(http_requests_total)",
		},
		{
			name: "empty result",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `http_requests_total{pod="nginx-3"}`,
		},
		{
			name: "last_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `last_over_time(http_requests_total[30s])`,
		},
		{
			name: "group",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `group(http_requests_total)`,
		},
		{
			name: "reset",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 100-1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `resets(http_requests_total[5m])`,
		},
		{
			name: "present_over_time",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: `present_over_time(http_requests_total[30s])`,
		},
	}

	lookbackDeltas := []time.Duration{30 * time.Second, time.Minute, 5 * time.Minute, 10 * time.Minute}
	for _, lookbackDelta := range lookbackDeltas {
		opts.LookbackDelta = lookbackDelta
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				test, err := promql.NewTest(t, tc.load)
				testutil.Ok(t, err)
				defer test.Close()

				testutil.Ok(t, test.Run())

				for _, disableFallback := range []bool{false, true} {
					t.Run(fmt.Sprintf("disableFallback=%v", disableFallback), func(t *testing.T) {
						newEngine := engine.New(engine.Opts{EngineOpts: opts, DisableFallback: disableFallback})
						q1, err := newEngine.NewInstantQuery(test.Storage(), nil, tc.query, queryTime)
						testutil.Ok(t, err)
						newResult := q1.Exec(context.Background())
						testutil.Ok(t, newResult.Err)

						oldEngine := promql.NewEngine(opts)
						q2, err := oldEngine.NewInstantQuery(test.Storage(), nil, tc.query, queryTime)
						testutil.Ok(t, err)

						oldResult := q2.Exec(context.Background())
						testutil.Ok(t, oldResult.Err)

						testutil.Equals(t, oldResult, newResult)
					})
				}
			})
		}
	}
}

func TestQueryCancellation(t *testing.T) {
	twelveHours := int64(12 * time.Hour.Seconds())

	start := time.Unix(0, 0)
	end := time.Unix(twelveHours, 0)
	step := time.Second * 30
	query := `sum(rate(http_requests_total{pod="nginx-1"}[10s]))`
	load := `load 30s
				http_requests_total{pod="nginx-1"} 1+1x1
				http_requests_total{pod="nginx-2"} 1+2x40`

	test, err := promql.NewTest(t, load)
	testutil.Ok(t, err)
	defer test.Close()

	testutil.Ok(t, test.Run())

	querier := &storage.MockQueryable{
		MockQuerier: &storage.MockQuerier{
			SelectMockFunction: func(sortSeries bool, hints *storage.SelectHints, matchers ...*labels.Matcher) storage.SeriesSet {
				return &testSeriesSet{
					series: &slowSeries{},
				}
			},
		},
	}

	newEngine := engine.New(engine.Opts{})
	q1, err := newEngine.NewRangeQuery(querier, nil, query, start, end, step)
	testutil.Ok(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-time.After(1000 * time.Millisecond)
		cancel()
	}()

	newResult := q1.Exec(ctx)
	testutil.Equals(t, context.Canceled, newResult.Err)
}

type testSeriesSet struct {
	i      int
	series storage.Series
}

func (s *testSeriesSet) Next() bool                 { s.i++; return s.i < 2 }
func (s *testSeriesSet) At() storage.Series         { return s.series }
func (s *testSeriesSet) Err() error                 { return nil }
func (s *testSeriesSet) Warnings() storage.Warnings { return nil }

type slowSeries struct{}

func (d slowSeries) Labels() labels.Labels       { return labels.FromStrings("foo", "bar") }
func (d slowSeries) Iterator() chunkenc.Iterator { return &slowIterator{} }

type slowIterator struct {
	ts int64
}

func (d *slowIterator) At() (int64, float64) { return d.ts, 1 }
func (d *slowIterator) Next() bool {
	<-time.After(10 * time.Millisecond)
	d.ts += 30 * 1000
	return true
}
func (d *slowIterator) Seek(t int64) bool {
	<-time.After(10 * time.Millisecond)
	d.ts = t
	return true
}
func (d *slowIterator) Err() error { return nil }
