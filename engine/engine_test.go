package engine_test

import (
	"context"
	"fpetkovski/promql-engine/engine"
	"testing"
	"time"

	"github.com/prometheus/prometheus/promql"
	"github.com/stretchr/testify/require"
)

func TestEngine(t *testing.T) {
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
		expected []promql.Vector
	}{
		{
			name: "sum by pod",
			load: `
load 30s
  http_requests_total{pod="nginx-1"} 1+1x4
  http_requests_total{pod="nginx-2"} 1+2x4
`,
			query: "sum by (pod) (http_requests_total)",
		},

		{
			name: "sum rate",
			load: `
load 30s
  http_requests_total{pod="nginx-1"} 1+1x4
  http_requests_total{pod="nginx-2"} 1+2x4
`,
			query: "sum(rate(http_requests_total[1m]))",
		},
		{
			name: "sum rate with stale series",
			load: `
load 30s
  http_requests_total{pod="nginx-1"} 1+1x4
  http_requests_total{pod="nginx-2"} 1+2x20
`,
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
			q1, err := newEngine.NewRangeQuery(test.Storage(), nil, tc.query, start, end, step)
			require.NoError(t, err)
			newResult := q1.Exec(context.Background())

			oldEngine := promql.NewEngine(opts)
			q2, err := oldEngine.NewRangeQuery(test.Storage(), nil, tc.query, start, end, step)
			require.NoError(t, err)
			oldResult := q2.Exec(context.Background())

			require.Equal(t, oldResult, newResult)
		})
	}
}
