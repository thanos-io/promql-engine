// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/internal/teststorage"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/promql"
)

// Results are compared against Prometheus in TestInstantQuery. This file holds
// what a result comparison cannot see: how many operators the chain became.

// Overlapping values, shared prefixes, a label on only some series and an
// empty-valued label: the shapes that make a chain's rules interact.
const relabelChainLoad = `load 30s
	http_requests_total{pod="api-1", job="web", env="prod", zone="a", tier="front"} 1+1x40
	http_requests_total{pod="api-2", job="web", env="dev", zone="a"} 2+2x40
	http_requests_total{pod="worker-1", job="batch", env="prod", zone="b", tier="back"} 3+3x40
	http_requests_total{pod="worker-2", job="batch", env="prod", zone="b", tier=""} 4+4x40
	http_requests_total{pod="odd", job="web", zone="c"} 5+5x40
`

func relabelChainEngines() (*promql.Engine, *engine.Engine) {
	opts := engine.Opts{
		Timeout:        1 * time.Hour,
		MaxSamples:     1e10,
		EnableAnalysis: true,
	}
	return promql.NewEngine(prometheusEngineOpts(opts)), engine.New(opts)
}

// Do not delete: this is the only test that sees the collapsing. Without it a
// chain silently reverts to one operator per call and results stay correct.
func TestRelabelChainCollapsesToOneOperator(t *testing.T) {
	t.Parallel()

	store := teststorage.LoadedStorage(t, relabelChainLoad)
	// Cleanup, not defer: the parallel subtests resume after this returns.
	t.Cleanup(func() { _ = store.Close() })

	_, ng := relabelChainEngines()
	ctx := context.Background()
	queryTime := time.Unix(600, 0)

	plan := func(t *testing.T, query string) string {
		t.Helper()
		q, err := ng.NewInstantQuery(ctx, store, nil, query, queryTime)
		testutil.Ok(t, err)
		defer q.Close()
		testutil.Ok(t, q.Exec(ctx).Err)
		return queryExplanation(q)
	}

	for _, tc := range []struct {
		name  string
		query string
		// rules is 0 when nothing collapses.
		operators int
		rules     int
	}{
		{
			name:      "two label_replace calls",
			query:     `label_replace(label_replace(http_requests_total, "a", "$1", "pod", "(.*)"), "b", "$1", "a", "(.*)")`,
			operators: 1,
			rules:     2,
		},
		{
			name:      "ten label_replace calls",
			query:     deepRelabelChain(10),
			operators: 1,
			rules:     10,
		},
		{
			name:      "two label_join calls",
			query:     `label_join(label_join(http_requests_total, "a", "-", "pod", "job"), "b", "/", "a", "env")`,
			operators: 1,
			rules:     2,
		},
		{
			name:      "mixed chain",
			query:     `label_join(label_replace(http_requests_total, "a", "$1", "pod", "(.*)"), "b", "-", "a", "job")`,
			operators: 1,
			rules:     2,
		},
		{
			name:      "single call is one operator with one rule",
			query:     `label_replace(http_requests_total, "a", "$1", "pod", "(.*)")`,
			operators: 1,
			rules:     0,
		},
		{
			// An aggregation between the calls breaks the chain.
			name: "chain broken by an aggregation",
			query: `
label_replace(
  sum by (a) (label_replace(http_requests_total, "a", "$1", "pod", "(.*)")),
  "b",
  "$1",
  "a",
  "(.*)"
)`,
			operators: 2,
			rules:     0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := plan(t, tc.query)
			// "[relabel] N rules" contains "[relabel]", so both forms count.
			testutil.Equals(t, tc.operators, strings.Count(got, "[relabel]"),
				"wrong number of relabel operators\nplan:\n%s", got)
			if tc.rules > 0 {
				testutil.Equals(t, 1, strings.Count(got, fmt.Sprintf("[relabel] %d rules", tc.rules)),
					"expected one operator carrying %d rules\nplan:\n%s", tc.rules, got)
			} else {
				testutil.Equals(t, 0, strings.Count(got, " rules"),
					"expected no collapsed operator\nplan:\n%s", got)
			}
		})
	}
}

// Collapsing drops the per-call duplicate-labelset checks, so pin the
// behaviour when those checks are off entirely too.
func TestRelabelChainWithDuplicateLabelChecksDisabled(t *testing.T) {
	t.Parallel()

	store := teststorage.LoadedStorage(t, relabelChainLoad)
	defer store.Close()

	opts := engine.Opts{
		Timeout:                     1 * time.Hour,
		MaxSamples:                  1e10,
		DisableDuplicateLabelChecks: true,
	}
	ng := engine.New(opts)
	promEngine := promql.NewEngine(prometheusEngineOpts(engine.Opts{
		Timeout:    1 * time.Hour,
		MaxSamples: 1e10,
	}))
	ctx := context.Background()
	queryTime := time.Unix(600, 0)

	for _, query := range []string{
		`label_replace(label_replace(http_requests_total, "a", "$1", "pod", "(.*)"), "b", "$1!", "a", "(.*)")`,
		`label_join(label_join(http_requests_total, "a", "-", "pod", "job"), "b", "/", "a", "env")`,
	} {
		t.Run(query, func(t *testing.T) {
			q1, err := ng.NewInstantQuery(ctx, store, nil, query, queryTime)
			testutil.Ok(t, err)
			defer q1.Close()
			got := q1.Exec(ctx)
			testutil.Ok(t, got.Err)

			q2, err := promEngine.NewInstantQuery(ctx, store, nil, query, queryTime)
			testutil.Ok(t, err)
			defer q2.Close()
			want := q2.Exec(ctx)
			testutil.Ok(t, want.Err)

			requireSameResult(t, want, got, query)
		})
	}
}

func deepRelabelChain(depth int) string {
	query := `label_replace(http_requests_total, "d0", "$1", "pod", "(.*)")`
	for i := 1; i < depth; i++ {
		query = fmt.Sprintf(`label_replace(%s, "d%d", "$1-%d", "d%d", "(.*)")`, query, i, i, i-1)
	}
	return query
}
