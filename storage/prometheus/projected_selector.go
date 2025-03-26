package prometheus

import (
	"context"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"sync"
)

type projectedSelector struct {
	selector   SeriesSelector
	dropLabels []string

	once   sync.Once
	series []SignedSeries
}

func NewProjectedSelector(selector SeriesSelector, dropLabels []string) SeriesSelector {
	return &projectedSelector{
		selector:   selector,
		dropLabels: dropLabels,
	}
}

func (f *projectedSelector) Matchers() []*labels.Matcher {
	return f.selector.Matchers()
}

func (f *projectedSelector) GetSeries(ctx context.Context, shard, numShards int) ([]SignedSeries, error) {
	var err error
	f.once.Do(func() { err = f.loadSeries(ctx) })
	if err != nil {
		return nil, err
	}

	return f.series, nil
}

func (f *projectedSelector) loadSeries(ctx context.Context) error {
	series, err := f.selector.GetSeries(ctx, 0, 1)
	if err != nil {
		return err
	}

	var i uint64
	f.series = make([]SignedSeries, 0, len(series))
	b := labels.NewBuilder(labels.EmptyLabels())
	for _, s := range series {
		b.Reset(s.Labels())
		b.Del(f.dropLabels...)
		f.series = append(f.series, SignedSeries{
			Series:    &projectedSeries{Series: s, lset: b.Labels()},
			Signature: i,
		})
		i++
	}

	return nil
}

// projectedSeries wraps a storage.Series but returns projected labels
type projectedSeries struct {
	storage.Series
	lset labels.Labels
}

func (s *projectedSeries) Labels() labels.Labels {
	return s.lset
}
