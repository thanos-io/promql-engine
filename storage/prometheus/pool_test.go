// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus

import (
	"testing"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
)

// TestSelectorPoolKeysOnRange checks that selectors with the same bounds but
// different ranges are not shared. An anchored 5m range with a 5m lookback
// delta reads the same window as a plain 10m range.
func TestSelectorPoolKeysOnRange(t *testing.T) {
	pool := NewSelectorPool(nil)
	matchers := []*labels.Matcher{labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, "m")}
	hints := func(rng int64) storage.SelectHints {
		return storage.SelectHints{Start: 0, End: 600_000, Range: rng, Func: "rate"}
	}

	anchored := pool.GetSelector(0, 600_000, 0, matchers, hints(300_000))
	plain := pool.GetSelector(0, 600_000, 0, matchers, hints(600_000))
	testutil.Assert(t, anchored != plain, "selectors with different ranges must not be shared")
	testutil.Assert(t, anchored == pool.GetSelector(0, 600_000, 0, matchers, hints(300_000)), "identical selectors should be shared")
}
