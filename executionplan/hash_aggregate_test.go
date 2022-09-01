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

type chanOperator struct {
	data chan promql.Vector
}

func (c *chanOperator) Next(ctx context.Context) (<-chan promql.Vector, error) {
	return c.data, nil
}

func TestAggregate(t *testing.T) {
	in := make(chan promql.Vector, 1)
	in <- []promql.Sample{{
		Metric: labels.FromStrings("__name__", "metric", "label", "v1"),
		Point:  promql.Point{T: 10, V: 20},
	}, {
		Metric: labels.FromStrings("__name__", "metric", "label", "v1"),
		Point:  promql.Point{T: 10, V: 40},
	}, {
		Metric: labels.FromStrings("__name__", "metric", "label", "v2"),
		Point:  promql.Point{T: 10, V: 50},
	},
	}
	close(in)

	aggregation, err := executionplan.NewAggregate(&chanOperator{data: in}, parser.SUM, true, []string{"__name__", "label"})
	require.NoError(t, err)

	result := make([]promql.Vector, 0)
	out, err := aggregation.Next(context.Background())
	require.NoError(t, err)

	for r := range out {
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
