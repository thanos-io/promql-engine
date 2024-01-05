// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus_test

import (
	"testing"

	"github.com/prometheus/prometheus/model/labels"
	promstg "github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	"github.com/stretchr/testify/require"

	storage "github.com/thanos-io/promql-engine/storage/prometheus"
)

func TestFilter_MultipleMatcherWithSameName(t *testing.T) {
	f := storage.NewFilter([]*labels.Matcher{
		labels.MustNewMatcher(labels.MatchNotEqual, "phase", "Running"),
		labels.MustNewMatcher(labels.MatchNotEqual, "phase", "Succeeded"),
	})

	require.Equal(t, false, f.Matches(&mockLabelSeries{labels: labels.FromStrings("phase", "Running")}))
}

func TestFilter_Matches(t *testing.T) {
	testCases := []struct {
		name     string
		matchers []*labels.Matcher
		series   promstg.Series
		expected bool
	}{
		{
			name:     "empty matchers",
			matchers: []*labels.Matcher{},
			series:   &mockLabelSeries{labels: labels.FromStrings("foo", "bar")},
			expected: true,
		},
		{
			name:     "no match",
			matchers: []*labels.Matcher{labels.MustNewMatcher(labels.MatchEqual, "foo", "bar")},
			series:   &mockLabelSeries{labels: labels.FromStrings("foo", "baz")},
		},
		{
			name:     "regex match",
			matchers: []*labels.Matcher{labels.MustNewMatcher(labels.MatchRegexp, "foo", "ba.")},
			series:   &mockLabelSeries{labels: labels.FromStrings("foo", "bar")},
			expected: true,
		},
		{
			name:     "regex no match",
			matchers: []*labels.Matcher{labels.MustNewMatcher(labels.MatchRegexp, "foo", "ba.")},
			series:   &mockLabelSeries{labels: labels.FromStrings("foo", "nope")},
		},
		{
			name:     "multiple matchers",
			matchers: []*labels.Matcher{labels.MustNewMatcher(labels.MatchEqual, "foo", "bar"), labels.MustNewMatcher(labels.MatchEqual, "baz", "qux")},
			series:   &mockLabelSeries{labels: labels.FromStrings("foo", "bar", "baz", "qux")},
			expected: true,
		},
		{
			name:     "single regex matcher, with label name not present",
			matchers: []*labels.Matcher{labels.MustNewMatcher(labels.MatchRegexp, "foo", ".*")},
			series:   &mockLabelSeries{labels: labels.FromStrings("bar", "baz")},
			expected: true,
		},
		{
			name:     "single regex matcher, with label name not present, negative regex",
			matchers: []*labels.Matcher{labels.MustNewMatcher(labels.MatchNotRegexp, "foo", ".*")},
			series:   &mockLabelSeries{labels: labels.FromStrings("bar", "baz")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := storage.NewFilter(tc.matchers)
			if got := f.Matches(tc.series); got != tc.expected {
				if tc.expected {
					t.Errorf("expected %s to match %s, but it did not.", tc.series.Labels().String(), tc.matchers)
				} else {
					t.Errorf("expected %s to not match %s, but it did.", tc.series.Labels().String(), tc.matchers)
				}
			}
		})
	}
}

type mockLabelSeries struct {
	labels labels.Labels
}

func (s *mockLabelSeries) Labels() labels.Labels {
	return s.labels
}

func (s *mockLabelSeries) Iterator(chunkenc.Iterator) chunkenc.Iterator {
	return nil
}
