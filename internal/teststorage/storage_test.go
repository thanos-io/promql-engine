// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package teststorage

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/efficientgo/core/errors"
	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
)

func TestConcurrentQuerier(t *testing.T) {
	// Multiple chunks per series exercise both the shared index scratch buffer
	// and the head-chunk cache, not just independent within-chunk decoding.
	var input strings.Builder
	input.WriteString("load 1s\n")
	for i := range 4 {
		fmt.Fprintf(&input, "metric{id=\"%d\"} %d+1x1000\n", i, i*10000)
	}
	s := LoadedStorage(t, input.String())
	q, err := s.Querier(0, 1000000)
	testutil.Ok(t, err)
	t.Cleanup(func() { testutil.Ok(t, q.Close()) })

	var wg sync.WaitGroup
	start := make(chan struct{})
	for range 8 {
		wg.Go(func() {
			<-start
			var it chunkenc.Iterator
			for range 3 {
				set := q.Select(context.Background(), false, &storage.SelectHints{
					Start: 0, End: 1000000, Step: 1000,
				}, labels.MustNewMatcher(labels.MatchEqual, "__name__", "metric"))
				seriesCount := 0
				for set.Next() {
					series := set.At()
					id, err := strconv.Atoi(series.Labels().Get("id"))
					testutil.Ok(t, err)
					it = series.Iterator(it)
					samples := 0
					for typ := it.Next(); typ != chunkenc.ValNone; typ = it.Next() {
						testutil.Equals(t, chunkenc.ValFloat, typ)
						ts, value := it.At()
						testutil.Equals(t, int64(samples*1000), ts)
						testutil.Equals(t, float64(id*10000+samples), value)
						testutil.Equals(t, ts, it.AtT())
						samples++
					}
					testutil.Ok(t, it.Err())
					testutil.Equals(t, 1001, samples)

					// Reuse must unwrap the fixture iterator, not nest locks or
					// leave it bound to the previous series' underlying iterator.
					it = series.Iterator(it)
					testutil.Equals(t, chunkenc.ValFloat, it.Seek(500000))
					ts, value := it.At()
					testutil.Equals(t, int64(500000), ts)
					testutil.Equals(t, float64(id*10000+500), value)
					testutil.Ok(t, it.Err())
					seriesCount++
				}
				testutil.Ok(t, set.Err())
				testutil.Equals(t, 4, seriesCount)
			}
		})
	}
	close(start)
	wg.Wait()
}

func TestWrapQueryablePreservesErrors(t *testing.T) {
	want := errors.New("cannot open querier")
	q := WrapQueryable(storage.QueryableFunc(func(mint, maxt int64) (storage.Querier, error) {
		testutil.Equals(t, int64(10), mint)
		testutil.Equals(t, int64(20), maxt)
		return nil, want
	}))
	querier, err := q.Querier(10, 20)
	testutil.Assert(t, querier == nil)
	testutil.Assert(t, errors.Is(err, want))
}
