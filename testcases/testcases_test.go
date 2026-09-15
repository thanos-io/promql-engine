// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package testcases

import (
	"testing"

	"github.com/prometheus/prometheus/promql/parser"
	"github.com/stretchr/testify/require"
)

// TestRangeQueries guards the YAML against edits that other language bindings
// would only discover at runtime.
func TestRangeQueries(t *testing.T) {
	cases, err := RangeQueries()
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	parser.EnableExperimentalFunctions = true
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := parser.ParseExpr(tc.Query)
			require.NoError(t, err, "query does not parse")
			require.Positive(t, tc.StepMS, "step must be positive")
			require.GreaterOrEqual(t, tc.EndMS, tc.StartMS, "end must not be before start")
		})
	}
}
