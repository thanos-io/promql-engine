package engine_test

import (
	"context"
	"testing"
	"time"

	"github.com/fpetkovski/promql-engine/engine"

	"github.com/prometheus/prometheus/promql"
	"github.com/stretchr/testify/require"
)

func TestQueriesAgainstOldEngine(t *testing.T) {
	start := time.Unix(0, 0)
	end := time.Unix(120, 0)
	step := time.Second * 60
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
		expected []promql.Vector
	}{
		{
			name: "sum",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "sum (http_requests_total)",
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
			name: "sum rate",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x4
					http_requests_total{pod="nginx-2"} 1+2x4`,
			query: "sum(rate(http_requests_total[1m]))",
		},
		{
			name: "sum rate with stale series",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x4
					http_requests_total{pod="nginx-2"} 1+2x20`,
			query: "sum(rate(http_requests_total[1m]))",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			test, err := promql.NewTest(t, tc.load)
			require.NoError(t, err)
			defer test.Close()

			err = test.Run()
			require.NoError(t, err)

			if tc.start.UnixMilli() == 0 {
				tc.start = start
			}
			if tc.end.UnixMilli() == 0 {
				tc.end = end
			}

			newEngine := engine.New()
			q1, err := newEngine.NewRangeQuery(test.Storage(), nil, tc.query, tc.start, tc.end, step)
			require.NoError(t, err)
			newResult := q1.Exec(context.Background())

			oldEngine := promql.NewEngine(opts)
			q2, err := oldEngine.NewRangeQuery(test.Storage(), nil, tc.query, tc.start, tc.end, step)
			require.NoError(t, err)
			oldResult := q2.Exec(context.Background())

			require.Equal(t, oldResult, newResult)
		})
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
			name: "sum by pod",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x4
					http_requests_total{pod="nginx-2"} 1+2x4`,
			query: "sum by (pod) (http_requests_total)",
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
					http_requests_total{pod="nginx-1"} 1+1x4
					http_requests_total{pod="nginx-2"} 1+2x20`,
			query: "sum(rate(http_requests_total[1m]))",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			test, err := promql.NewTest(t, tc.load)
			require.NoError(t, err)
			defer test.Close()

			err = test.Run()
			require.NoError(t, err)

			newEngine := engine.New()
			q1, err := newEngine.NewInstantQuery(test.Storage(), nil, tc.query, queryTime)
			require.NoError(t, err)
			newResult := q1.Exec(context.Background())

			oldEngine := promql.NewEngine(opts)
			q2, err := oldEngine.NewInstantQuery(test.Storage(), nil, tc.query, queryTime)
			require.NoError(t, err)
			oldResult := q2.Exec(context.Background())

			require.Equal(t, oldResult, newResult)
		})
	}
}
