package executionplan_test

import (
	"context"
	"fmt"
	"fpetkovski/promql-engine/executionplan"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMatrixSelector(t *testing.T) {
	testCases := []struct {
		name        string
		load        string
		start       time.Time
		end         time.Time
		interval    time.Duration
		selectRange time.Duration
		expected    []promql.Matrix
	}{
		{
			name: "timestamps match with step",
			load: `load 30s
              bar 0 1 10 100`,
			start:       time.Unix(0, 0),
			end:         time.Unix(120, 0),
			interval:    30 * time.Second,
			selectRange: 30 * time.Second,
			expected: []promql.Matrix{
				{
					seriesWithPoints("bar", promql.Point{T: 0, V: 0}),
				},
				{
					seriesWithPoints("bar", []promql.Point{{T: 0, V: 0}, {T: 30000, V: 1}}...),
				},
				{
					seriesWithPoints("bar", []promql.Point{{T: 30000, V: 1}, {T: 60000, V: 10}}...),
				},
				{
					seriesWithPoints("bar", []promql.Point{{T: 60000, V: 10}, {T: 90000, V: 100}}...),
				},
				{
					seriesWithPoints("bar", []promql.Point{{T: 90000, V: 100}}...),
				},
			},
		},
		{
			name: "timestamps match with step",
			load: `load 29s
              bar 0 1 10 100`,
			start:       time.Unix(0, 0),
			end:         time.Unix(120, 0),
			interval:    30 * time.Second,
			selectRange: 30 * time.Second,
			expected: []promql.Matrix{
				{
					seriesWithPoints("bar", promql.Point{T: 0, V: 0}),
				},
				{
					seriesWithPoints("bar", []promql.Point{{T: 0, V: 0}, {T: 29000, V: 1}}...),
				},
				{
					seriesWithPoints("bar", []promql.Point{{T: 58000, V: 10}}...),
				},
				{
					seriesWithPoints("bar", []promql.Point{{T: 87000, V: 100}}...),
				},
				{
					seriesWithPoints("bar", nil...),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			test, err := promql.NewTest(t, tc.load)
			require.NoError(t, err)
			defer test.Close()

			err = test.Run()
			require.NoError(t, err)

			ng, err := test.QueryEngine().NewRangeQuery(test.Storage(), nil, "rate(bar[30s])", tc.start, tc.end, tc.interval)
			require.NoError(t, err)
			fmt.Println(ng.Exec(context.Background()))

			nameMatcher, err := labels.NewMatcher(labels.MatchEqual, labels.MetricName, "bar")
			require.NoError(t, err)
			matchers := []*labels.Matcher{nameMatcher}

			selector := executionplan.NewMatrixSelector(test.Storage(), matchers, nil, tc.start, tc.end, tc.interval, tc.selectRange)
			out, err := selector.Next(context.Background())
			require.NoError(t, err)

			result := make([]promql.Vector, 0, len(out))
			for r := range out {
				result = append(result, r)
			}
			fmt.Println(result)
			require.Equal(t, tc.expected, result)
		})
	}
}

func seriesWithPoints(name string, points ...promql.Point) promql.Series {
	return promql.Series{
		Metric: labels.FromStrings("__name__", name), Points: points,
	}
}
