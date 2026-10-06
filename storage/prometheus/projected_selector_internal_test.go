// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus

import (
	"context"
	"strconv"
	"testing"

	"github.com/thanos-io/promql-engine/logicalplan"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	"github.com/stretchr/testify/require"
)

// testSeriesHashLabel is the label name used to carry original series identity
// across projection in these tests.
const testSeriesHashLabel = "__series_hash__"

// TestProjectedSelector_HashCollisionDetection verifies that series which hash to
// the same value but have different label sets stay distinguishable after
// projection, while series that are genuinely identical keep a shared identity.
//
// Real collisions cannot be constructed against labels.Labels.Hash, so the hash
// function is injected. That is also why this test lives in the internal test
// package rather than in projected_selector_test.go.
func TestProjectedSelector_HashCollisionDetection(t *testing.T) {
	t.Parallel()

	// Every series hashes to the same value, so the deduper must fall back to
	// full label comparison for all of them.
	const collidingHash = uint64(42)
	alwaysCollide := func(labels.Labels) uint64 { return collidingHash }

	for _, tc := range []struct {
		name string
		// input label sets, in the order the inner selector returns them.
		input []labels.Labels
		// expected series hash label value per input series.
		expectedHashValues []string
	}{
		{
			name: "distinct series colliding on hash get distinct identities",
			input: []labels.Labels{
				labels.FromStrings("__name__", "m", "job", "app", "instance", "1"),
				labels.FromStrings("__name__", "m", "job", "app", "instance", "2"),
				labels.FromStrings("__name__", "m", "job", "app", "instance", "3"),
			},
			expectedHashValues: []string{"42", "42_1", "42_2"},
		},
		{
			name: "identical series keep a shared identity",
			input: []labels.Labels{
				labels.FromStrings("__name__", "m", "job", "app", "instance", "1"),
				labels.FromStrings("__name__", "m", "job", "app", "instance", "1"),
			},
			expectedHashValues: []string{"42", "42"},
		},
		{
			name: "repeats of a colliding series reuse its identity",
			input: []labels.Labels{
				labels.FromStrings("__name__", "m", "job", "app", "instance", "1"),
				labels.FromStrings("__name__", "m", "job", "app", "instance", "2"),
				labels.FromStrings("__name__", "m", "job", "app", "instance", "1"),
				labels.FromStrings("__name__", "m", "job", "app", "instance", "2"),
			},
			expectedHashValues: []string{"42", "42_1", "42", "42_1"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Include mode on "job" alone: every input projects down to the same
			// {job="app"} label set, so the series hash label is the only remaining
			// discriminator.
			selector := newProjectedSelectorWithHash(t,
				stubSelector(tc.input...),
				&logicalplan.Projection{Labels: []string{"job"}, Include: true},
				alwaysCollide,
			)

			series, err := selector.GetSeries(context.Background(), 0, 1)
			require.NoError(t, err)
			require.Len(t, series, len(tc.input))

			hashValues := make([]string, 0, len(series))
			for _, s := range series {
				lset := s.Series.Labels()
				require.Equal(t, "app", lset.Get("job"))
				hashValues = append(hashValues, lset.Get(testSeriesHashLabel))
			}
			require.Equal(t, tc.expectedHashValues, hashValues)

			// The projected label sets must agree with the input on which series are
			// the same: identical inputs project to equal labels, different inputs
			// project to different labels.
			for i := range tc.input {
				for j := i + 1; j < len(tc.input); j++ {
					sameInput := labels.Equal(tc.input[i], tc.input[j])
					sameProjection := labels.Equal(series[i].Series.Labels(), series[j].Series.Labels())
					require.Equal(t, sameInput, sameProjection,
						"series %d and %d: input equal=%v but projected equal=%v", i, j, sameInput, sameProjection)
				}
			}
		})
	}
}

// TestProjectedSelector_NoCollisionKeepsPlainHash guards the collision free path:
// with the real hash function the series hash label must still be the plain
// decimal hash of the original labels, unchanged by collision detection.
func TestProjectedSelector_NoCollisionKeepsPlainHash(t *testing.T) {
	t.Parallel()

	input := []labels.Labels{
		labels.FromStrings("__name__", "m", "job", "app", "instance", "1"),
		labels.FromStrings("__name__", "m", "job", "app", "instance", "2"),
	}

	selector := NewProjectedSelector(
		stubSelector(input...),
		&logicalplan.Projection{Labels: []string{"job"}, Include: true},
		testSeriesHashLabel,
	)

	series, err := selector.GetSeries(context.Background(), 0, 1)
	require.NoError(t, err)
	require.Len(t, series, len(input))

	for i, s := range series {
		require.Equal(t, strconv.FormatUint(input[i].Hash(), 10), s.Series.Labels().Get(testSeriesHashLabel))
	}
}

func TestOriginHashDeduper_VariantOf(t *testing.T) {
	t.Parallel()

	seriesA := labels.FromStrings("job", "app", "instance", "1")
	seriesB := labels.FromStrings("job", "app", "instance", "2")
	seriesC := labels.FromStrings("job", "app", "instance", "3")

	input := []SignedSeries{
		{Series: &stubSeries{lset: seriesA}},
		{Series: &stubSeries{lset: seriesB}},
		{Series: &stubSeries{lset: seriesA}},
		{Series: &stubSeries{lset: seriesC}},
	}

	const sharedHash = uint64(7)
	deduper := newOriginHashDeduper(len(input))

	// First distinct label set for the hash.
	require.Equal(t, 0, deduper.variantOf(input, 0, seriesA, sharedHash))
	// Second distinct label set for the same hash is a collision.
	require.Equal(t, 1, deduper.variantOf(input, 1, seriesB, sharedHash))
	// A repeat of the first label set reuses variant 0.
	require.Equal(t, 0, deduper.variantOf(input, 2, seriesA, sharedHash))
	// A third distinct label set continues the ordinals.
	require.Equal(t, 2, deduper.variantOf(input, 3, seriesC, sharedHash))
	// A different hash starts over at variant 0.
	require.Equal(t, 0, deduper.variantOf(input, 3, seriesC, sharedHash+1))
}

// TestOriginHashDeduper_NoCollisionAllocation asserts that the collision map is
// only allocated when a collision actually occurs, so the common path does not
// pay for it.
func TestOriginHashDeduper_NoCollisionAllocation(t *testing.T) {
	t.Parallel()

	seriesA := labels.FromStrings("job", "app", "instance", "1")
	seriesB := labels.FromStrings("job", "app", "instance", "2")
	input := []SignedSeries{
		{Series: &stubSeries{lset: seriesA}},
		{Series: &stubSeries{lset: seriesB}},
	}

	deduper := newOriginHashDeduper(len(input))
	require.Equal(t, 0, deduper.variantOf(input, 0, seriesA, 1))
	require.Equal(t, 0, deduper.variantOf(input, 1, seriesB, 2))
	require.Nil(t, deduper.collisions)

	// Now force a collision and confirm the map appears.
	require.Equal(t, 1, deduper.variantOf(input, 1, seriesB, 1))
	require.NotNil(t, deduper.collisions)
}

func TestSeriesHashValue(t *testing.T) {
	t.Parallel()

	require.Equal(t, "42", seriesHashValue(42, 0))
	require.Equal(t, "42_1", seriesHashValue(42, 1))
	require.Equal(t, "42_2", seriesHashValue(42, 2))

	// A disambiguated value must never be mistaken for a plain hash.
	require.NotEqual(t, seriesHashValue(421, 0), seriesHashValue(42, 1))
}

// newProjectedSelectorWithHash builds a projectedSelector with an injected hash
// function so that collision handling can be exercised deterministically.
func newProjectedSelectorWithHash(t *testing.T, inner SeriesSelector, projection *logicalplan.Projection, hash func(labels.Labels) uint64) SeriesSelector {
	t.Helper()

	selector := NewProjectedSelector(inner, projection, testSeriesHashLabel)
	projected, ok := selector.(*projectedSelector)
	require.True(t, ok, "expected NewProjectedSelector to return a *projectedSelector")
	projected.hashLabels = hash

	return projected
}

func stubSelector(lsets ...labels.Labels) SeriesSelector {
	series := make([]SignedSeries, 0, len(lsets))
	for i, lset := range lsets {
		series = append(series, SignedSeries{Series: &stubSeries{lset: lset}, Signature: uint64(i)})
	}

	return &stubSeriesSelector{series: series}
}

// stubSeriesSelector is a SeriesSelector returning a fixed set of series.
type stubSeriesSelector struct {
	series []SignedSeries
}

func (s *stubSeriesSelector) GetSeries(_ context.Context, _, _ int) ([]SignedSeries, error) {
	return s.series, nil
}

func (s *stubSeriesSelector) Matchers() []*labels.Matcher {
	return nil
}

// stubSeries is a storage.Series carrying only labels.
type stubSeries struct {
	lset labels.Labels
}

func (s *stubSeries) Labels() labels.Labels {
	return s.lset
}

func (s *stubSeries) Iterator(_ chunkenc.Iterator) chunkenc.Iterator {
	return chunkenc.NewNopIterator()
}

var _ storage.Series = &stubSeries{}
