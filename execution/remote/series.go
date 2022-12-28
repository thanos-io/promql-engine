package remote

import (
	"context"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	engstore "github.com/thanos-community/promql-engine/execution/storage"
)

type promqlSeries struct {
	series promql.Series
}

func (p promqlSeries) Labels() labels.Labels { return p.series.Metric }

func (p promqlSeries) Iterator() chunkenc.Iterator {
	return newSamplesIterator(p.series.Points)
}

type storageAdapter struct {
	series promql.Matrix
}

func newStorageAdapter(series promql.Matrix) *storageAdapter {
	return &storageAdapter{series: series}
}

func (s storageAdapter) GetSeries(_ context.Context, _, _ int) ([]engstore.SignedSeries, error) {
	result := make([]engstore.SignedSeries, len(s.series))
	for i := range s.series {
		result[i] = engstore.SignedSeries{
			Signature: uint64(i),
			Series:    &promqlSeries{series: s.series[i]},
		}
	}
	return result, nil
}

func (s storageAdapter) Matchers() []*labels.Matcher { return nil }
