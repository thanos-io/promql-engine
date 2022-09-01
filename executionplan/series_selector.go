package executionplan

import (
	"context"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"sync"
	"time"
)

type seriesSelector struct {
	storage  storage.Storage
	mint     int64
	maxt     int64
	matchers []*labels.Matcher

	once   sync.Once
	series []storage.Series
}

func newSeriesFilter(storage storage.Storage, mint time.Time, maxt time.Time, matchers []*labels.Matcher) *seriesSelector {
	return &seriesSelector{
		storage:  storage,
		mint:     mint.UnixMilli(),
		maxt:     maxt.UnixMilli(),
		matchers: matchers,
	}
}

func (o *seriesSelector) Series(ctx context.Context, shard int, numShards int) ([]storage.Series, error) {
	var err error
	o.once.Do(func() {
		querier, qErr := o.storage.Querier(ctx, o.mint, o.maxt)
		if qErr != nil {
			err = qErr
			return
		}
		defer querier.Close()

		seriesSet := querier.Select(true, nil, o.matchers...)
		for seriesSet.Next() {
			o.series = append(o.series, seriesSet.At())
		}
	})
	if err != nil {
		return nil, err
	}

	result := make([]storage.Series, 0, len(o.series)/numShards)
	for _, s := range o.series {
		if numShards == 1 || s.Labels().Hash()%uint64(numShards) == uint64(shard) {
			result = append(result, s)
		}
	}

	return result, nil
}
