// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package binary

import (
	"strings"
	"testing"

	"github.com/thanos-io/promql-engine/warnings"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/histogram"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

func TestHistogramTrimming(t *testing.T) {
	h := &histogram.FloatHistogram{Count: 4, ZeroCount: 4, ZeroThreshold: 1, CounterResetHint: histogram.NotCounterReset}
	original := h.Copy()
	for _, op := range []parser.ItemType{parser.TRIM_LOWER, parser.TRIM_UPPER} {
		t.Run(op.String(), func(t *testing.T) {
			_, got, keep, warn, err := binOp(op, 0, 0, h, nil)
			testutil.Ok(t, err)
			testutil.Assert(t, keep)
			testutil.Equals(t, warnings.Warnings(0), warn)
			testutil.Equals(t, 2.0, got.Count)
			testutil.Equals(t, original, h, "trimming must not mutate shared input")
			testutil.Assert(t, !shouldDropMetricName(op, false))

			for _, invalid := range [][2]*histogram.FloatHistogram{{nil, nil}, {nil, h}, {h, h}} {
				_, _, keep, warn, err := binOp(op, 0, 0, invalid[0], invalid[1])
				testutil.Ok(t, err)
				testutil.Assert(t, !keep)
				testutil.Equals(t, warnings.WarnIncompatibleTypesInBinOp, warn)
			}
		})
	}
}

func TestHistogramSubtractionProducesGauge(t *testing.T) {
	h := &histogram.FloatHistogram{Count: 4, Sum: 8, CounterResetHint: histogram.NotCounterReset}
	_, got, keep, _, err := binOp(parser.SUB, 0, 0, h, h)
	testutil.Ok(t, err)
	testutil.Assert(t, keep)
	testutil.Equals(t, histogram.GaugeType, got.CounterResetHint)
	testutil.Equals(t, histogram.NotCounterReset, h.CounterResetHint)
}

// TestErrManyToManyMatchDeterministicOrder verifies that errManyToManyMatch.Error()
// produces a byte-identical error string regardless of which of the two duplicate
// series happened to land in the "original" vs "duplicate" slot.
//
// Background: the choice of which series becomes original vs duplicate is driven
// by upstream StepVector ordering, which is in turn fed by map iteration (e.g.
// in execution/scan/subquery.go). That made the error message non-deterministic
// across runs and across parallel engine workers, breaking error-string equality
// checks in downstream consumers (see cortexproject/cortex#7546).
//
// Fix: lexicographically sort the two rendered label-set strings before
// formatting the message.
func TestErrManyToManyMatchDeterministicOrder(t *testing.T) {
	matching := &parser.VectorMatching{
		On:             true,
		MatchingLabels: []string{"job"},
	}

	for _, tc := range []struct {
		name string
		side binOpSide
		a    labels.Labels
		b    labels.Labels
	}{
		{
			name: "simple labels right side",
			side: rhBinOpSide,
			a:    labels.FromStrings("__name__", "requests_total", "instance", "a", "job", "api"),
			b:    labels.FromStrings("__name__", "requests_total", "instance", "b", "job", "api"),
		},
		{
			name: "simple labels left side",
			side: lhBinOpSide,
			a:    labels.FromStrings("zone", "us-east-1", "job", "api"),
			b:    labels.FromStrings("zone", "us-west-2", "job", "api"),
		},
		{
			// Adversarial values that include the same punctuation we use to
			// frame the rendered list ([, ], comma, braces, quotes). Sorting
			// the rendered label-set strings — not the labels themselves —
			// keeps the comparison robust against these characters.
			name: "adversarial label values",
			side: rhBinOpSide,
			a: labels.FromStrings(
				"detail", `left[,{}"]`,
				"instance", `a[,{}"]`,
				"job", `api[,{}"]`,
			),
			b: labels.FromStrings(
				"detail", `right[,{}"]`,
				"instance", `b[,{}"]`,
				"job", `api[,{}"]`,
			),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			aFirst := (&errManyToManyMatch{
				matching:  matching,
				side:      tc.side,
				original:  tc.a,
				duplicate: tc.b,
			}).Error()
			bFirst := (&errManyToManyMatch{
				matching:  matching,
				side:      tc.side,
				original:  tc.b,
				duplicate: tc.a,
			}).Error()

			testutil.Equals(t, aFirst, bFirst)

			// Sanity-check that both rendered label sets are still present and
			// that the existing prefix/format expected by downstream consumers
			// (e.g. Prometheus-compatible error matchers) is preserved.
			testutil.Assert(t, strings.Contains(aFirst, tc.a.String()), "expected %q to contain %q", aFirst, tc.a.String())
			testutil.Assert(t, strings.Contains(aFirst, tc.b.String()), "expected %q to contain %q", aFirst, tc.b.String())
			testutil.Assert(t, strings.Contains(aFirst, "found duplicate series for the match group"), "missing expected prefix in %q", aFirst)
			testutil.Assert(t, strings.Contains(aFirst, "many-to-many matching not allowed"), "missing expected suffix in %q", aFirst)
			testutil.Assert(t, strings.Contains(aFirst, string(tc.side)+" hand-side"), "missing expected side marker in %q", aFirst)
		})
	}
}
