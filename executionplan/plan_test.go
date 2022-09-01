package executionplan

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/stretchr/testify/require"
)

func TestNewPlan(t *testing.T) {

	start := time.Unix(0, 0)
	end := time.Unix(120, 0)
	step := time.Second * 60

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
			expected: []promql.Vector{
				{
					{Metric: labels.FromStrings("pod", "nginx-1"), Point: promql.Point{T: 0, V: 1}},
					{Metric: labels.FromStrings("pod", "nginx-2"), Point: promql.Point{T: 0, V: 1}},
				},
				{
					{Metric: labels.FromStrings("pod", "nginx-1"), Point: promql.Point{T: 60000, V: 3}},
					{Metric: labels.FromStrings("pod", "nginx-2"), Point: promql.Point{T: 60000, V: 5}},
				},
				{
					{Metric: labels.FromStrings("pod", "nginx-1"), Point: promql.Point{T: 120000, V: 5}},
					{Metric: labels.FromStrings("pod", "nginx-2"), Point: promql.Point{T: 120000, V: 9}},
				},
			},
		},

		{
			name: "sum rate",
			load: `
load 30s
  http_requests_total{pod="nginx-1"} 1+1x4
  http_requests_total{pod="nginx-2"} 1+2x4
`,

			query: "sum(rate(http_requests_total[1m]))",
			expected: []promql.Vector{
				{
					{Metric: labels.Labels{}, Point: promql.Point{T: 0, V: 0}},
				},
				{
					{Metric: labels.Labels{}, Point: promql.Point{T: 60000, V: 0.1}},
				},
				{
					{Metric: labels.Labels{}, Point: promql.Point{T: 120000, V: 0.1}},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			test, err := promql.NewTest(t, tc.load)
			require.NoError(t, err)
			defer test.Close()

			err = test.Run()
			require.NoError(t, err)

			expr, err := parser.ParseExpr(tc.query)
			require.NoError(t, err)

			plan, err := New(expr, test.Storage(), start, end, step)
			require.NoError(t, err)

			result := make([]promql.Vector, 0)
			for {
				r, err := plan.Next(context.Background())
				require.NoError(t, err)
				if r == nil {
					break
				}
				sort.Slice(r, func(i, j int) bool {
					return r[i].Metric.Hash() < r[j].Metric.Hash()
				})
				for _, s := range r {
					sort.Sort(s.Metric)
				}
				result = append(result, r)
			}

			require.Equal(t, tc.expected, result)
		})
	}
}
