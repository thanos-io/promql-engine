// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus

import (
	"testing"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/stretchr/testify/require"
)

func TestSelectorPoolSharesSelectorsAcrossRangeWindows(t *testing.T) {
	pool := NewSelectorPool(nil, nil)

	matchers := []*labels.Matcher{
		labels.MustNewMatcher(labels.MatchEqual, "__name__", "http_requests_total"),
		labels.MustNewMatcher(labels.MatchEqual, "pod", "nginx-1"),
	}

	maxt := int64(600000)
	step := int64(30000)

	// First selector: rate(metric[5m]) → mint = 600000 - 300000 = 300000
	hints1 := storage.SelectHints{Start: 300000, End: maxt, Step: step, Func: "rate", Range: 300000}
	sel1 := pool.GetFilteredSelector(300000, maxt, step, matchers, nil, hints1)

	// Second selector: rate(metric[2m]) → mint = 600000 - 120000 = 480000
	hints2 := storage.SelectHints{Start: 480000, End: maxt, Step: step, Func: "rate", Range: 120000}
	sel2 := pool.GetFilteredSelector(480000, maxt, step, matchers, nil, hints2)

	// Both should share the same underlying selector.
	require.Equal(t, sel1, sel2, "selectors with same matchers but different range windows should share one cache entry")

	// The cached selector should use the wider mint (300000, not 480000).
	key := DefaultSelectorHash(matchers, 300000, maxt, hints1)
	require.Equal(t, int64(300000), pool.selectors[key].hints.Start,
		"cached selector should use the widest (earliest) mint")
}

func TestSelectorPoolCustomHashSharesAcrossFunctions(t *testing.T) {
	// Custom hash that excludes Func — for storage layers where Func is a noop.
	customHash := func(matchers []*labels.Matcher, _, maxt int64, hints storage.SelectHints) uint64 {
		neutralHints := hints
		neutralHints.Func = ""
		return DefaultSelectorHash(matchers, 0, maxt, neutralHints)
	}

	pool := NewSelectorPool(nil, customHash)

	matchers := []*labels.Matcher{
		labels.MustNewMatcher(labels.MatchEqual, "__name__", "http_requests_total"),
	}

	maxt := int64(600000)
	step := int64(30000)

	// rate(metric[5m])
	hints1 := storage.SelectHints{Start: 300000, End: maxt, Step: step, Func: "rate", Range: 300000}
	sel1 := pool.GetFilteredSelector(300000, maxt, step, matchers, nil, hints1)

	// avg_over_time(metric[2m])
	hints2 := storage.SelectHints{Start: 480000, End: maxt, Step: step, Func: "avg_over_time", Range: 120000}
	sel2 := pool.GetFilteredSelector(480000, maxt, step, matchers, nil, hints2)

	// Custom hash excludes Func, so they should share.
	require.Equal(t, sel1, sel2, "custom hash excluding Func should share across functions")

	// Mint should be widened to the earlier value.
	for _, sel := range pool.selectors {
		require.Equal(t, int64(300000), sel.hints.Start)
	}
}
