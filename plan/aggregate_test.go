package plan_test

import (
	"context"
	"fpetkovski/promql-engine/plan"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/stretchr/testify/require"
	"testing"
)

type chanOperator struct {
	data chan promql.Matrix
}

func (c *chanOperator) Next(ctx context.Context) (<-chan promql.Matrix, error) {
	return c.data, nil
}

func TestAggregate(t *testing.T) {
	in := make(chan promql.Matrix, 1)
	in <- promql.Matrix{
		promql.Series{
			Metric: labels.FromStrings("__name__", "metric", "label", "v1"),
			Points: []promql.Point{{T: 10, V: 20}},
		},
		promql.Series{
			Metric: labels.FromStrings("__name__", "metric", "label", "v1"),
			Points: []promql.Point{{T: 10, V: 40}},
		},
		promql.Series{
			Metric: labels.FromStrings("__name__", "metric", "label", "v2"),
			Points: []promql.Point{{T: 10, V: 50}},
		},
	}
	close(in)

	aggregation, err := plan.NewAggregate(&chanOperator{data: in}, parser.SUM, true, []string{"__name__", "label"})
	require.NoError(t, err)

	result := make([]promql.Matrix, 0)
	out, err := aggregation.Next(context.Background())
	require.NoError(t, err)

	for r := range out {
		result = append(result, r)
	}
	expected := []promql.Matrix{
		{
			{
				Metric: labels.FromStrings("__name__", "metric", "label", "v1"),
				Points: []promql.Point{
					{T: 10, V: 60},
				},
			},
			{
				Metric: labels.FromStrings("__name__", "metric", "label", "v2"),
				Points: []promql.Point{
					{T: 10, V: 50},
				},
			},
		},
	}

	require.Equal(t, expected, result)
}
