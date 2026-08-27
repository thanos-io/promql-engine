// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus

import "github.com/thanos-io/promql-engine/execution/model"

// collectSeriesIDs returns all unique series IDs seen across step vectors.
func collectSeriesIDs(vectors []model.StepVector) []uint64 {
	seen := make(map[uint64]bool)
	for _, v := range vectors {
		for _, id := range v.SampleIDs {
			seen[id] = true
		}
		for _, id := range v.HistogramIDs {
			seen[id] = true
		}
	}
	ids := make([]uint64, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}

// containsAll checks that all expected IDs are present in got.
func containsAll(got []uint64, expected []uint64) bool {
	set := make(map[uint64]bool, len(got))
	for _, id := range got {
		set[id] = true
	}
	for _, id := range expected {
		if !set[id] {
			return false
		}
	}
	return true
}

// containsAny checks that at least one of the IDs is present in got.
func containsAny(got []uint64, ids []uint64) bool {
	set := make(map[uint64]bool, len(got))
	for _, id := range got {
		set[id] = true
	}
	for _, id := range ids {
		if set[id] {
			return true
		}
	}
	return false
}
