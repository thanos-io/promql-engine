// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus_test

import (
	"context"
	"testing"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	"github.com/stretchr/testify/require"

	"github.com/thanos-io/promql-engine/logicalplan"
	storage "github.com/thanos-io/promql-engine/storage/prometheus"
)

func TestProjectedSelector_IncludeMode(t *testing.T) {
	t.Parallel()

	inner := &mockSeriesSelector{
		series: []storage.SignedSeries{
			{Series: &mockLabelSeries{labels: labels.FromStrings("__name__", "http_requests", "job", "app", "instance", "1", "env", "prod")}, Signature: 0},
			{Series: &mockLabelSeries{labels: labels.FromStrings("__name__", "http_requests", "job", "api", "instance", "2", "env", "dev")}, Signature: 1},
		},
	}

	projection := &logicalplan.Projection{
		Labels:  []string{"job"},
		Include: true,
	}

	selector := storage.NewProjectedSelector(inner, projection, "__series_hash__")

	series, err := selector.GetSeries(context.Background(), 0, 1)
	require.NoError(t, err)
	require.Len(t, series, 2)

	// First series: should only have "job" + "__series_hash__"
	lset0 := series[0].Series.Labels()
	require.Equal(t, "app", lset0.Get("job"))
	require.Equal(t, "", lset0.Get("instance"))
	require.Equal(t, "", lset0.Get("env"))
	require.Equal(t, "", lset0.Get("__name__"))
	require.NotEmpty(t, lset0.Get("__series_hash__"))

	// Second series: should only have "job" + "__series_hash__"
	lset1 := series[1].Series.Labels()
	require.Equal(t, "api", lset1.Get("job"))
	require.Equal(t, "", lset1.Get("instance"))
	require.Equal(t, "", lset1.Get("env"))
	require.NotEmpty(t, lset1.Get("__series_hash__"))

	// Series hash should differ for different original label sets
	require.NotEqual(t, lset0.Get("__series_hash__"), lset1.Get("__series_hash__"))

	// Signatures should be dense 0-based
	require.Equal(t, uint64(0), series[0].Signature)
	require.Equal(t, uint64(1), series[1].Signature)
}

func TestProjectedSelector_ExcludeMode(t *testing.T) {
	t.Parallel()

	inner := &mockSeriesSelector{
		series: []storage.SignedSeries{
			{Series: &mockLabelSeries{labels: labels.FromStrings("__name__", "http_requests", "job", "app", "instance", "1", "env", "prod")}, Signature: 0},
			{Series: &mockLabelSeries{labels: labels.FromStrings("__name__", "http_requests", "job", "api", "instance", "2", "env", "dev")}, Signature: 1},
		},
	}

	projection := &logicalplan.Projection{
		Labels:  []string{"instance"},
		Include: false,
	}

	selector := storage.NewProjectedSelector(inner, projection, "__series_hash__")

	series, err := selector.GetSeries(context.Background(), 0, 1)
	require.NoError(t, err)
	require.Len(t, series, 2)

	// First series: all labels except "instance", plus "__series_hash__"
	lset0 := series[0].Series.Labels()
	require.Equal(t, "http_requests", lset0.Get("__name__"))
	require.Equal(t, "app", lset0.Get("job"))
	require.Equal(t, "prod", lset0.Get("env"))
	require.Equal(t, "", lset0.Get("instance"))
	require.NotEmpty(t, lset0.Get("__series_hash__"))
}

func TestProjectedSelector_IncludeEmpty(t *testing.T) {
	t.Parallel()

	// Include mode with empty labels list means keep nothing except series hash.
	// This is the case for sum() without any by clause.
	inner := &mockSeriesSelector{
		series: []storage.SignedSeries{
			{Series: &mockLabelSeries{labels: labels.FromStrings("__name__", "metric", "job", "app", "instance", "1")}, Signature: 0},
			{Series: &mockLabelSeries{labels: labels.FromStrings("__name__", "metric", "job", "api", "instance", "2")}, Signature: 1},
		},
	}

	projection := &logicalplan.Projection{
		Labels:  []string{},
		Include: true,
	}

	selector := storage.NewProjectedSelector(inner, projection, "__series_hash__")

	series, err := selector.GetSeries(context.Background(), 0, 1)
	require.NoError(t, err)
	require.Len(t, series, 2)

	// Series should only have __series_hash__
	lset0 := series[0].Series.Labels()
	require.Equal(t, 1, lset0.Len(), "expected only __series_hash__ label, got: %s", lset0.String())
	require.NotEmpty(t, lset0.Get("__series_hash__"))
}

func TestProjectedSelector_NoOpProjection(t *testing.T) {
	t.Parallel()

	// Exclude mode with empty labels list is a no-op (exclude nothing).
	inner := &mockSeriesSelector{
		series: []storage.SignedSeries{
			{Series: &mockLabelSeries{labels: labels.FromStrings("job", "app", "instance", "1")}, Signature: 0},
		},
	}

	projection := &logicalplan.Projection{
		Labels:  []string{},
		Include: false,
	}

	// Should return the inner selector unchanged (no wrapping).
	selector := storage.NewProjectedSelector(inner, projection, "__series_hash__")

	series, err := selector.GetSeries(context.Background(), 0, 1)
	require.NoError(t, err)
	require.Len(t, series, 1)

	// Labels should be unchanged (no projection applied, no __series_hash__)
	lset := series[0].Series.Labels()
	require.Equal(t, "app", lset.Get("job"))
	require.Equal(t, "1", lset.Get("instance"))
	require.Equal(t, "", lset.Get("__series_hash__"))
}

func TestProjectedSelector_NilProjection(t *testing.T) {
	t.Parallel()

	inner := &mockSeriesSelector{
		series: []storage.SignedSeries{
			{Series: &mockLabelSeries{labels: labels.FromStrings("job", "app")}, Signature: 0},
		},
	}

	// Nil projection should return the inner selector unchanged.
	selector := storage.NewProjectedSelector(inner, nil, "__series_hash__")

	series, err := selector.GetSeries(context.Background(), 0, 1)
	require.NoError(t, err)
	require.Len(t, series, 1)

	lset := series[0].Series.Labels()
	require.Equal(t, "app", lset.Get("job"))
	require.Equal(t, "", lset.Get("__series_hash__"))
}

func TestProjectedSelector_Sharding(t *testing.T) {
	t.Parallel()

	inner := &mockSeriesSelector{
		series: []storage.SignedSeries{
			{Series: &mockLabelSeries{labels: labels.FromStrings("job", "a", "instance", "1")}, Signature: 0},
			{Series: &mockLabelSeries{labels: labels.FromStrings("job", "b", "instance", "2")}, Signature: 1},
			{Series: &mockLabelSeries{labels: labels.FromStrings("job", "c", "instance", "3")}, Signature: 2},
			{Series: &mockLabelSeries{labels: labels.FromStrings("job", "d", "instance", "4")}, Signature: 3},
		},
	}

	projection := &logicalplan.Projection{
		Labels:  []string{"job"},
		Include: true,
	}

	selector := storage.NewProjectedSelector(inner, projection, "__series_hash__")

	// Get first shard (2 shards total)
	shard0, err := selector.GetSeries(context.Background(), 0, 2)
	require.NoError(t, err)
	require.Len(t, shard0, 2)

	// Get second shard
	shard1, err := selector.GetSeries(context.Background(), 1, 2)
	require.NoError(t, err)
	require.Len(t, shard1, 2)

	// Signatures in each shard should be re-based to 0
	require.Equal(t, uint64(0), shard0[0].Signature)
	require.Equal(t, uint64(1), shard0[1].Signature)
	require.Equal(t, uint64(0), shard1[0].Signature)
	require.Equal(t, uint64(1), shard1[1].Signature)
}

func TestProjectedSelector_IdenticalProjectedLabels_DifferentHash(t *testing.T) {
	t.Parallel()

	// Two series with same projected labels but different original labels
	// must have different __series_hash__ values.
	inner := &mockSeriesSelector{
		series: []storage.SignedSeries{
			{Series: &mockLabelSeries{labels: labels.FromStrings("job", "app", "instance", "1", "env", "prod")}, Signature: 0},
			{Series: &mockLabelSeries{labels: labels.FromStrings("job", "app", "instance", "2", "env", "dev")}, Signature: 1},
		},
	}

	projection := &logicalplan.Projection{
		Labels:  []string{"job"},
		Include: true,
	}

	selector := storage.NewProjectedSelector(inner, projection, "__series_hash__")

	series, err := selector.GetSeries(context.Background(), 0, 1)
	require.NoError(t, err)
	require.Len(t, series, 2)

	// Both have job=app, but __series_hash__ should differ
	lset0 := series[0].Series.Labels()
	lset1 := series[1].Series.Labels()
	require.Equal(t, lset0.Get("job"), lset1.Get("job"))
	require.NotEqual(t, lset0.Get("__series_hash__"), lset1.Get("__series_hash__"))
}

func TestProjectedSelector_EmptySeriesHashLabel(t *testing.T) {
	t.Parallel()

	inner := &mockSeriesSelector{
		series: []storage.SignedSeries{
			{Series: &mockLabelSeries{labels: labels.FromStrings("job", "app", "instance", "1")}, Signature: 0},
		},
	}

	projection := &logicalplan.Projection{
		Labels:  []string{"job"},
		Include: true,
	}

	// Empty series hash label means projection is not applied at the selector level
	// (the remote storage/queryable wrapper handles it instead).
	selector := storage.NewProjectedSelector(inner, projection, "")

	series, err := selector.GetSeries(context.Background(), 0, 1)
	require.NoError(t, err)
	require.Len(t, series, 1)

	// Labels should be unchanged since projection is skipped with empty hash label
	lset := series[0].Series.Labels()
	require.Equal(t, "app", lset.Get("job"))
	require.Equal(t, "1", lset.Get("instance"))
}

func TestProjectedSelector_Matchers(t *testing.T) {
	t.Parallel()

	matchers := []*labels.Matcher{
		labels.MustNewMatcher(labels.MatchEqual, "job", "app"),
	}

	inner := &mockSeriesSelector{
		series:   nil,
		matchers: matchers,
	}

	projection := &logicalplan.Projection{
		Labels:  []string{"job"},
		Include: true,
	}

	selector := storage.NewProjectedSelector(inner, projection, "__series_hash__")
	require.Equal(t, matchers, selector.Matchers())
}

func TestProjectedSelector_Iterator(t *testing.T) {
	t.Parallel()

	// Verify that projected series delegates Iterator to the original
	inner := &mockSeriesSelector{
		series: []storage.SignedSeries{
			{Series: &mockIteratorSeries{
				labels: labels.FromStrings("job", "app", "instance", "1"),
			}, Signature: 0},
		},
	}

	projection := &logicalplan.Projection{
		Labels:  []string{"job"},
		Include: true,
	}

	selector := storage.NewProjectedSelector(inner, projection, "__series_hash__")

	series, err := selector.GetSeries(context.Background(), 0, 1)
	require.NoError(t, err)
	require.Len(t, series, 1)

	// Iterator should be delegated to the original
	iter := series[0].Series.Iterator(nil)
	require.NotNil(t, iter)
}

// mockSeriesSelector is a mock implementation of SeriesSelector for testing.
type mockSeriesSelector struct {
	series   []storage.SignedSeries
	matchers []*labels.Matcher
}

func (m *mockSeriesSelector) GetSeries(_ context.Context, shard, numShards int) ([]storage.SignedSeries, error) {
	if numShards <= 1 {
		return m.series, nil
	}
	start := shard * len(m.series) / numShards
	end := min((shard+1)*len(m.series)/numShards, len(m.series))
	return m.series[start:end], nil
}

func (m *mockSeriesSelector) Matchers() []*labels.Matcher {
	return m.matchers
}

// mockIteratorSeries is a mock series that returns a non-nil iterator.
type mockIteratorSeries struct {
	labels labels.Labels
}

func (s *mockIteratorSeries) Labels() labels.Labels {
	return s.labels
}

func (s *mockIteratorSeries) Iterator(_ chunkenc.Iterator) chunkenc.Iterator {
	return chunkenc.NewNopIterator()
}
