// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus

import (
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/logicalplan"
	"github.com/thanos-io/promql-engine/query"

	"github.com/prometheus/prometheus/promql/parser"
	"github.com/stretchr/testify/require"
)

func TestScannersMinMaxTime(t *testing.T) {
	for _, tcase := range []struct {
		expr       string
		start, end time.Time
		step       time.Duration
		min, max   int64
	}{
		{
			expr:  "foo offset 5m",
			start: time.Unix(200, 0),
			end:   time.Unix(200, 0),
			step:  time.Second,

			min: -400000,
			max: -100000,
		},
		{
			expr:  `absent_over_time(http_requests_total @ 1800.000[1h:1m])`,
			start: time.Unix(200, 0),
			end:   time.Unix(200, 0),
			step:  time.Second,

			min: 1500000,
			max: 1800000,
		},
		{
			expr:  `rate(testcounter_zero_cutoff[20m])`,
			start: time.Unix(200, 0),
			end:   time.Unix(200, 0),
			step:  time.Second,

			min: -1000000,
			max: 200000,
		},
		{
			expr:  `rate(testcounter_zero_cutoff[20m])`,
			start: time.Unix(200, 0),
			end:   time.Unix(400, 0),
			step:  time.Second,

			min: -1000000,
			max: 400000,
		},
		{
			expr:  "foo @ 20",
			start: time.Unix(200, 0),
			end:   time.Unix(200, 0),
			step:  time.Second,

			min: -280000,
			max: 20000,
		},
	} {
		t.Run(tcase.expr, func(t *testing.T) {
			p, err := parser.ParseExpr(tcase.expr)
			require.NoError(t, err)

			qOpts := &query.Options{
				Start:         tcase.start,
				End:           tcase.end,
				Step:          tcase.step,
				LookbackDelta: 5 * time.Duration(time.Minute),
			}

			lplan := logicalplan.NewFromAST(p, qOpts, logicalplan.PlanOptions{})

			min, max := lplan.MinMaxTime(qOpts)

			require.Equal(t, tcase.min, min)
			require.Equal(t, tcase.max, max)
		})
	}
}
