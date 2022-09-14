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
			name: "sum",
			load: `load 30s
					http_requests_total{pod="nginx-1"} 1+1x15
					http_requests_total{pod="nginx-2"} 1+2x18`,
			query: "sum (http_requests_total)",
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
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			test, err := promql.NewTest(t, tc.load)
			require.NoError(t, err)
			defer test.Close()

			err = test.Run()
			require.NoError(t, err)

			if tc.start.Equal(time.Time{}) {
				tc.start = start
			}
			if tc.end.Equal(time.Time{}) {
				tc.end = end
			}
			if tc.step == 0 {
				tc.step = step
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
