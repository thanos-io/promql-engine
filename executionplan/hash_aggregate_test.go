package executionplan_test

import (
	"context"
	"fpetkovski/promql-engine/executionplan"
	"sort"
	"testing"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/stretchr/testify/require"
)

type stubOperator struct {
	data []promql.Vector
}

func (c *stubOperator) Next(ctx context.Context) (promql.Vector, error) {
	if len(c.data) == 0 {
		return nil, nil
	}
	r := c.data[0]
	c.data = c.data[1:]
	return r, nil
}

func TestAggregate(t *testing.T) {
	in := []promql.Vector{{{
		Metric: labels.FromStrings("__name__", "metric", "label", "v1"),
		Point:  promql.Point{T: 10, V: 20},
	}, {
		Metric: labels.FromStrings("__name__", "metric", "label", "v1"),
		Point:  promql.Point{T: 10, V: 40},
	}, {
		Metric: labels.FromStrings("__name__", "metric", "label", "v2"),
		Point:  promql.Point{T: 10, V: 50},
	}}}

	aggregation, err := executionplan.NewAggregate(&stubOperator{data: in}, parser.SUM, true, []string{"__name__", "label"})
	require.NoError(t, err)

	result := make([]promql.Vector, 0)
	for {
		r, err := aggregation.Next(context.Background())
		require.NoError(t, err)
		if r == nil {
			break
		}
		sort.Slice(r, func(i, j int) bool {
			return r[i].Metric.Hash() < r[j].Metric.Hash()
		})
		result = append(result, r)
	}

	expected := []promql.Vector{{{
		Metric: labels.FromStrings("__name__", "metric", "label", "v2"),
		Point:  promql.Point{T: 10, V: 50},
	}, {
		Metric: labels.FromStrings("__name__", "metric", "label", "v1"),
		Point:  promql.Point{T: 10, V: 60},
	}}}

	require.Equal(t, expected, result)
}
