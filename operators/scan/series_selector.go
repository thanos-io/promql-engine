package scan

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
)

type seriesSelector struct {
	storage  storage.Queryable
	mint     int64
	maxt     int64
	matchers []*labels.Matcher

	once     sync.Once
	scanners []vectorScan
	series   []labels.Labels
}

func NewSeriesFilter(storage storage.Queryable, mint time.Time, maxt time.Time, matchers []*labels.Matcher) *seriesSelector {
	return &seriesSelector{
		storage: storage,

		mint:     mint.UnixMilli() - 5*time.Minute.Milliseconds(),
		maxt:     maxt.UnixMilli(),
		matchers: matchers,
	}
}

func (o *seriesSelector) Series(ctx context.Context, shard int, numShards int) ([]vectorScan, []labels.Labels, error) {
	var err error
	o.once.Do(func() { err = o.loadSeries(ctx) })
	if err != nil {
		return nil, nil, err
	}

	start := shard * len(o.scanners) / numShards
	end := (shard + 1) * len(o.scanners) / numShards
	if end > len(o.scanners) {
		end = len(o.scanners)
	}
	return o.scanners[start:end], o.series[start:end], nil

}

func (o *seriesSelector) loadSeries(ctx context.Context) error {
	//start := time.Now()
	//defer func() {
	//	fmt.Println("Done fetching series", time.Since(start))
	//}()
	querier, err := o.storage.Querier(ctx, o.mint, o.maxt)
	if err != nil {
		return err
	}
	defer querier.Close()

	seriesSet := querier.Select(false, nil, o.matchers...)
	i := 0
	for seriesSet.Next() {
		series := seriesSet.At()
		o.scanners = append(o.scanners, vectorScan{
			labels:    series.Labels(),
			signature: uint64(i),
			samples:   storage.NewMemoizedIterator(series.Iterator(), 5*time.Minute.Milliseconds()),
		})
		o.series = append(o.series, series.Labels())
		i++
	}
	//fmt.Println("Total series fetched", len(o.series))

	return nil
}
