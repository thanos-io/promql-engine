package plan

import (
	"context"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
	"time"
)

func TestNewPlan(t *testing.T) {
	load := `
load 30s
  http_requests_total{pod="nginx-1"} 1+1x4
  http_requests_total{pod="nginx-2"} 1+2x4
`
	test, err := promql.NewTest(t, load)
	require.NoError(t, err)
	defer test.Close()

	err = test.Run()
	require.NoError(t, err)

	query := "sum by (pod) (http_requests_total)"
	expr, err := parser.ParseExpr(query)
	require.NoError(t, err)

	start := time.Unix(0, 0)
	end := time.Unix(120, 0)
	step := time.Second * 60
	plan, err := New(expr, test.Storage(), start, end, step)
	require.NoError(t, err)

	out, err := plan.Next(context.Background())
	result := make([]promql.Matrix, 0, len(out))
	for r := range out {
		for _, s := range r {
			sort.Sort(s.Metric)
		}
		result = append(result, r)
	}

	expected := []promql.Matrix{
		{
			{Metric: labels.FromStrings("pod", "nginx-1"), Points: []promql.Point{{T: 0, V: 1}}},
			{Metric: labels.FromStrings("pod", "nginx-2"), Points: []promql.Point{{T: 0, V: 1}}},
		},
		{
			{Metric: labels.FromStrings("pod", "nginx-1"), Points: []promql.Point{{T: 60000, V: 3}}},
			{Metric: labels.FromStrings("pod", "nginx-2"), Points: []promql.Point{{T: 60000, V: 5}}},
		},
		{
			{Metric: labels.FromStrings("pod", "nginx-1"), Points: []promql.Point{{T: 120000, V: 5}}},
			{Metric: labels.FromStrings("pod", "nginx-2"), Points: []promql.Point{{T: 120000, V: 9}}},
		},
	}

	require.Equal(t, expected, result)
}
