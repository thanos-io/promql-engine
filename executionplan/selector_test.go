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

func TestSelector(t *testing.T) {
	testCases := []struct {
		name        string
		load        string
		start       time.Time
		end         time.Time
		interval    time.Duration
		selectRange time.Duration
	}{
		{
			name: "sum_over_time with all values",
			load: `load 30s
              bar 0 1 10 100 1000`,
			start:       time.Unix(0, 0),
			end:         time.Unix(120, 0),
			interval:    60 * time.Second,
			selectRange: 30 * time.Second,
		},
		{
			name: "sum_over_time with all values",
			load: `load 30s
              bar 0 1 10 100 1000`,
			start:       time.Unix(0, 0),
			end:         time.Unix(120, 0),
			interval:    60 * time.Second,
			selectRange: 0 * time.Second,
		},
		{
			name: "sum_over_time with all values",
			load: `load 30s
              bar 1+1x4`,
			start:       time.Unix(0, 0),
			end:         time.Unix(120, 0),
			interval:    60 * time.Second,
			selectRange: 0 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			test, err := promql.NewTest(t, tc.load)
			require.NoError(t, err)
			defer test.Close()

			err = test.Run()
			require.NoError(t, err)

			nameMatcher, err := labels.NewMatcher(labels.MatchEqual, labels.MetricName, "bar")
			require.NoError(t, err)

			start := tc.start
			end := tc.end
			step := tc.interval
			matchers := []*labels.Matcher{nameMatcher}
			selector := executionplan.NewSelector(test.Storage(), matchers, nil, start, end, step)
			out, err := selector.Next(context.Background())
			require.NoError(t, err)

			result := make([]promql.Vector, 0, len(out))
			for r := range out {
				result = append(result, r)
			}
			fmt.Println(result)
		})
	}
}
