// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/logicalplan"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/promqltest"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/util/annotations"
)

// recordingQueryable records the projection hints each selector receives.
type recordingQueryable struct {
	storage.Queryable
	mu   sync.Mutex
	seen []string
}

func (r *recordingQueryable) Querier(mint, maxt int64) (storage.Querier, error) {
	q, err := r.Queryable.Querier(mint, maxt)
	if err != nil {
		return nil, err
	}
	return &recordingQuerier{Querier: q, rec: r}, nil
}

type recordingQuerier struct {
	storage.Querier
	rec *recordingQueryable
}

func (r *recordingQuerier) Select(ctx context.Context, sorted bool, hints *storage.SelectHints, ms ...*labels.Matcher) storage.SeriesSet {
	mode := "exclude"
	if hints != nil && hints.ProjectionInclude {
		mode = "include"
	}
	var lbls []string
	if hints != nil {
		lbls = append(lbls, hints.ProjectionLabels...)
		sort.Strings(lbls)
	}
	r.rec.mu.Lock()
	r.rec.seen = append(r.rec.seen, fmt.Sprintf("%v -> %s%v", ms, mode, lbls))
	r.rec.mu.Unlock()
	return r.Querier.Select(ctx, sorted, hints, ms...)
}

func (r *recordingQuerier) LabelValues(context.Context, string, *storage.LabelHints, ...*labels.Matcher) ([]string, annotations.Annotations, error) {
	return nil, nil, nil
}
func (r *recordingQuerier) LabelNames(context.Context, *storage.LabelHints, ...*labels.Matcher) ([]string, annotations.Annotations, error) {
	return nil, nil, nil
}

func TestProjectionHintsProbe(t *testing.T) {
	load := `load 30s
		http_requests_total{pod="nginx-1", job="app", env="prod", instance="1"} 1+1x40
		errors_total{pod="nginx-3", job="api", env="staging", instance="3"} 3+3x40`

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	projectionEngine := engine.New(engine.Opts{
		EngineOpts: promql.EngineOpts{Timeout: time.Minute, MaxSamples: 1e10},
		LogicalOptimizers: []logicalplan.Optimizer{
			logicalplan.SortMatchers{},
			logicalplan.ProjectionOptimizer{},
			logicalplan.DetectHistogramStatsOptimizer{},
			logicalplan.MergeSelectsOptimizer{},
		},
	})

	ctx := context.Background()
	for _, query := range []string{
		`sum by (env) (http_requests_total)`,
		`http_requests_total or on (env) errors_total`,
		`sum by (pod) (http_requests_total or on (env) errors_total)`,
		`sum by (pod) (http_requests_total and on (env) errors_total)`,
		`sum by (pod) (http_requests_total / on (env) group_left (job) errors_total)`,
		`sum by (pod) (http_requests_total / on (env, job) group_left (instance) errors_total)`,
		`sum by (env) (http_requests_total / on (env) errors_total)`,
		`sum by (env) (label_replace(http_requests_total, "env", "$1", "pod", "(.*)"))`,
		`sum by (env) (max_over_time(http_requests_total[2m:30s]))`,
		`count_values by (env) ("v", http_requests_total)`,
		`sum by (env) (sort_by_label(http_requests_total, "pod"))`,
		`sum by (env) (topk(1, http_requests_total))`,
		`sum by (env) (absent(http_requests_total{job="nope"}))`,
		`histogram_quantile(0.9, sum by (env, le) (http_requests_total))`,
	} {
		rec := &recordingQueryable{Queryable: &projectionQueryable{Queryable: storage, useAtHash: true}}
		q, err := projectionEngine.MakeInstantQuery(ctx, rec, &engine.QueryOpts{}, query, time.Unix(600, 0))
		if err != nil {
			t.Logf("%-70s PLAN ERROR %v", query, err)
			continue
		}
		res := q.Exec(ctx)
		q.Close()
		testutil.Ok(t, res.Err, "query: %s", query)
		sort.Strings(rec.seen)
		t.Logf("%s\n        %v", query, rec.seen)
	}
}
