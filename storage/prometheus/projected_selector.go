// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus

import (
	"context"
	"slices"
	"strconv"
	"sync"

	"github.com/thanos-io/promql-engine/logicalplan"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
)

// seriesHashVariantSeparator separates the origin hash from the variant ordinal
// in a series hash label value. It cannot appear in a plain decimal hash, which
// keeps disambiguated values distinct from undisambiguated ones.
const seriesHashVariantSeparator = "_"

// projectedSelector wraps a SeriesSelector and applies label projection
// on GetSeries(). This reduces memory by only materializing the labels
// that downstream operators actually need (e.g. grouping labels for aggregations).
type projectedSelector struct {
	selector        SeriesSelector
	projection      *logicalplan.Projection
	seriesHashLabel string

	// hashLabels computes the identity hash of a series' original label set.
	// It is a field so that tests can inject a hash function which collides;
	// production code always uses labels.Labels.Hash.
	hashLabels func(labels.Labels) uint64

	once   sync.Once
	series []SignedSeries
}

// NewProjectedSelector creates a new projected selector that wraps the given selector
// and applies label projection based on the provided projection hints.
// seriesHashLabel is the label name used to store the original series hash
// for identity preservation (typically "__series_hash__").
// If seriesHashLabel is empty, projection is not applied at the selector level
// (the remote storage or queryable wrapper is expected to handle it instead).
func NewProjectedSelector(selector SeriesSelector, projection *logicalplan.Projection, seriesHashLabel string) SeriesSelector {
	if projection == nil || isNoOpProjection(projection) || seriesHashLabel == "" {
		return selector
	}
	return &projectedSelector{
		selector:        selector,
		projection:      projection,
		seriesHashLabel: seriesHashLabel,
		hashLabels:      labels.Labels.Hash,
	}
}

// isNoOpProjection returns true if the projection would not actually remove any labels.
// An exclude-mode projection with no labels to exclude is a no-op.
func isNoOpProjection(p *logicalplan.Projection) bool {
	return !p.Include && len(p.Labels) == 0
}

func (p *projectedSelector) Matchers() []*labels.Matcher {
	return p.selector.Matchers()
}

func (p *projectedSelector) GetSeries(ctx context.Context, shard, numShards int) ([]SignedSeries, error) {
	var err error
	p.once.Do(func() { err = p.loadSeries(ctx) })
	if err != nil {
		return nil, err
	}

	return seriesShard(p.series, shard, numShards), nil
}

func (p *projectedSelector) loadSeries(ctx context.Context) error {
	series, err := p.selector.GetSeries(ctx, 0, 1)
	if err != nil {
		return err
	}

	p.series = make([]SignedSeries, 0, len(series))

	// Once labels are projected away, the series hash label is the only thing that
	// keeps two distinct input series distinguishable. Two series with different
	// label sets that happen to hash to the same value would therefore become
	// indistinguishable and could be silently merged by downstream operators. The
	// deduper detects that case and hands out a disambiguated hash value.
	deduper := newOriginHashDeduper(len(series))

	for i, s := range series {
		lset := s.Series.Labels()
		originHash := p.hashLabels(lset)
		variant := deduper.variantOf(series, i, lset, originHash)

		p.series = append(p.series, SignedSeries{
			Series:    p.projectSeries(s.Series, lset, seriesHashValue(originHash, variant)),
			Signature: uint64(i),
		})
	}

	return nil
}

// originHashDeduper assigns a variant ordinal to every distinct original label
// set, keyed by the label set's hash. The first distinct label set seen for a
// hash gets variant 0, and each further distinct label set sharing that hash
// (i.e. a genuine hash collision) gets the next ordinal. Repeats of an
// already-seen label set get that label set's existing ordinal, so series which
// are truly identical continue to share one identity.
type originHashDeduper struct {
	// first maps an origin hash to the index of the first series that produced it.
	first map[uint64]int
	// collisions holds, per colliding hash, the indices of the additional distinct
	// label sets observed for that hash, ordered by variant. It stays nil unless a
	// genuine collision occurs, which for a 64 bit hash is astronomically rare.
	collisions map[uint64][]int
}

func newOriginHashDeduper(size int) *originHashDeduper {
	return &originHashDeduper{first: make(map[uint64]int, size)}
}

// variantOf returns the variant ordinal for the label set of series[index].
// Comparing against previously seen series reads their labels through the
// backing slice, so no additional copy of the original labels is retained.
func (d *originHashDeduper) variantOf(series []SignedSeries, index int, lset labels.Labels, originHash uint64) int {
	first, seen := d.first[originHash]
	if !seen {
		d.first[originHash] = index
		return 0
	}
	if labels.Equal(series[first].Series.Labels(), lset) {
		return 0
	}

	// Genuine hash collision: the same hash for a different label set.
	for ordinal, other := range d.collisions[originHash] {
		if labels.Equal(series[other].Series.Labels(), lset) {
			return ordinal + 1
		}
	}
	if d.collisions == nil {
		d.collisions = make(map[uint64][]int)
	}
	d.collisions[originHash] = append(d.collisions[originHash], index)

	return len(d.collisions[originHash])
}

// seriesHashValue renders the value stored in the series hash label. Variant 0
// is the plain decimal hash, so the collision free path is byte for byte what it
// was before collision detection existed. Higher variants get a suffix so that
// colliding series stay distinct; the value is only ever compared, never parsed
// back into a number.
func seriesHashValue(originHash uint64, variant int) string {
	hash := strconv.FormatUint(originHash, 10)
	if variant == 0 {
		return hash
	}
	return hash + seriesHashVariantSeparator + strconv.Itoa(variant)
}

func (p *projectedSelector) projectSeries(s storage.Series, originalLabels labels.Labels, hashValue string) storage.Series {
	projectedLabels := p.projectLabels(originalLabels, hashValue)

	if labels.Equal(originalLabels, projectedLabels) {
		return s
	}

	return &projectedSeries{
		Series: s,
		lset:   projectedLabels,
	}
}

// projectLabels applies the projection to a label set, keeping or removing
// labels based on the include/exclude mode. Always preserves __name__ is NOT
// done here — the projection optimizer decides which labels to include.
// hashValue identifies the original series and is stored in the series hash
// label to preserve that identity across projection.
func (p *projectedSelector) projectLabels(lset labels.Labels, hashValue string) labels.Labels {
	builder := labels.NewBuilder(labels.EmptyLabels())

	if p.projection.Include {
		// Include mode: only keep the labels in the projection list
		lset.Range(func(l labels.Label) {
			if slices.Contains(p.projection.Labels, l.Name) {
				builder.Set(l.Name, l.Value)
			}
		})
	} else {
		// Exclude mode: keep all labels except those in the projection list
		lset.Range(func(l labels.Label) {
			if !slices.Contains(p.projection.Labels, l.Name) {
				builder.Set(l.Name, l.Value)
			}
		})
	}

	// Append series hash label to preserve original series identity
	if p.seriesHashLabel != "" {
		builder.Set(p.seriesHashLabel, hashValue)
	}

	return builder.Labels()
}

// projectedSeries wraps a storage.Series but returns projected labels.
type projectedSeries struct {
	storage.Series
	lset labels.Labels
}

func (s *projectedSeries) Labels() labels.Labels {
	return s.lset
}

func (s *projectedSeries) Iterator(it chunkenc.Iterator) chunkenc.Iterator {
	return s.Series.Iterator(it)
}
