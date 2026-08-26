// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/logicalplan"

	"github.com/cortexproject/promqlsmith"
	"github.com/efficientgo/core/errors"
	"github.com/efficientgo/core/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/promql/promqltest"
	promstorage "github.com/prometheus/prometheus/storage"
)

// TestProjectionAtHashFuzz is a differential fuzz test for the out-of-band
// origin-hash contract (AtHash, no __series_hash__ label) that Thanos uses.
// The load intentionally contains several series per aggregation group so that
// projection trimming collapses distinct series onto identical label sets.
func TestProjectionAtHashFuzz(t *testing.T) {
	t.Parallel()

	seed := int64(1755000000)
	if s := os.Getenv("FUZZ_SEED"); s != "" {
		v, err := strconv.ParseInt(s, 10, 64)
		testutil.Ok(t, err)
		seed = v
	}
	testRuns := 20000
	if s := os.Getenv("FUZZ_RUNS"); s != "" {
		v, err := strconv.Atoi(s)
		testutil.Ok(t, err)
		testRuns = v
	}
	rnd := rand.New(rand.NewSource(seed))

	// Two instances per (pod, job, env) plus two pods per env: any projection
	// narrower than {pod, instance} produces duplicate label sets.
	load := `load 30s
		http_requests_total{pod="nginx-1", job="app", env="prod", instance="1", cluster="a"} 1+1x40
		http_requests_total{pod="nginx-1", job="app", env="prod", instance="2", cluster="a"} 2+2x40
		http_requests_total{pod="nginx-2", job="app", env="prod", instance="1", cluster="b"} 3+3x40
		http_requests_total{pod="nginx-2", job="api", env="dev", instance="2", cluster="b"} 4+4x40
		http_requests_total{pod="nginx-3", job="api", env="dev", instance="1", cluster="a"} 5+5x40
		http_request_duration{pod="nginx-1", job="app", env="prod", instance="1", cluster="a"} {{schema:0 sum:5 count:4 buckets:[1 2 1]}}+{{schema:0 sum:5 count:4 buckets:[1 2 1]}}x40
		http_request_duration{pod="nginx-1", job="app", env="prod", instance="2", cluster="a"} {{schema:0 sum:7 count:6 buckets:[2 2 2]}}+{{schema:0 sum:7 count:6 buckets:[2 2 2]}}x40
		http_request_duration{pod="nginx-2", job="app", env="prod", instance="1", cluster="b"} {{schema:0 sum:9 count:8 buckets:[3 2 3]}}+{{schema:0 sum:9 count:8 buckets:[3 2 3]}}x40
		errors_total{pod="nginx-1", job="app", env="prod", instance="1", cluster="a"} 0.5+0.5x40
		errors_total{pod="nginx-1", job="app", env="prod", instance="2", cluster="a"} 1+1x40
		errors_total{pod="nginx-2", job="app", env="prod", instance="1", cluster="b"} 1.5+1.5x40
		errors_total{pod="nginx-3", job="api", env="dev", instance="1", cluster="a"} 2+2x40`

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	seriesSet, err := getSeries(context.Background(), storage, "http_requests_total")
	testutil.Ok(t, err)

	ps := promqlsmith.New(rnd, seriesSet,
		promqlsmith.WithEnableOffset(false),
		promqlsmith.WithEnableAtModifier(true),
		promqlsmith.WithEnableExperimentalPromQLFunctions(true),
		promqlsmith.WithEnabledAggrs([]parser.ItemType{
			parser.SUM, parser.MIN, parser.MAX, parser.AVG, parser.COUNT,
			parser.TOPK, parser.BOTTOMK, parser.GROUP, parser.STDDEV, parser.QUANTILE,
			parser.COUNT_VALUES, parser.STDVAR,
		}),
		promqlsmith.WithEnableVectorMatching(true),
	)

	engineOpts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}
	normalEngine := engine.New(engine.Opts{
		EngineOpts:        engineOpts,
		LogicalOptimizers: logicalplan.AllOptimizers,
	})
	// Thanos's configuration: DefaultOptimizers with the projection optimizer
	// appended last, and no SeriesHashLabel (identity travels out of band).
	projectionEngine := engine.New(engine.Opts{
		EngineOpts:        engineOpts,
		LogicalOptimizers: append(append([]logicalplan.Optimizer{}, logicalplan.DefaultOptimizers...), logicalplan.ProjectionOptimizer{}),
	})
	projectionStorage := &projectionQueryable{Queryable: storage, useAtHash: true}
	var rawStorage promstorage.Queryable = storage

	ctx := context.Background()
	queryTime := time.Unix(600, 0)
	rangeStart, rangeEnd, step := time.Unix(240, 0), time.Unix(900, 0), 90*time.Second

	type failure struct{ query, kind, detail string }
	var failures []failure
	var ties int
	seen := make(map[string]struct{})

	for i := range testRuns {
		var (
			query   string
			isRange = i%2 == 1
		)
		for {
			var expr parser.Expr
			if isRange {
				expr = ps.WalkRangeQuery()
			} else {
				expr = ps.WalkInstantQuery()
			}
			query = expr.Pretty(0)
			if !containsProjectionExprs(expr) {
				continue
			}
			break
		}
		if _, ok := seen[query]; ok {
			continue
		}
		seen[query] = struct{}{}

		newQuery := func(e *engine.Engine, projected bool) (promql.Query, error) {
			q := rawStorage
			if projected {
				q = projectionStorage
			}
			if isRange {
				return e.MakeRangeQuery(ctx, q, &engine.QueryOpts{}, query, rangeStart, rangeEnd, step)
			}
			return e.MakeInstantQuery(ctx, q, &engine.QueryOpts{}, query, queryTime)
		}

		normal, err := newQuery(normalEngine, false)
		if err != nil {
			continue
		}
		normalResult := normal.Exec(ctx)
		normal.Close()
		if normalResult.Err != nil {
			continue
		}

		projected, err := newQuery(projectionEngine, true)
		if err != nil {
			failures = append(failures, failure{query, "plan error", err.Error()})
			continue
		}
		projectedResult := projected.Exec(ctx)
		projected.Close()
		if projectedResult.Err != nil {
			failures = append(failures, failure{query, "exec error", projectedResult.Err.Error()})
			continue
		}
		if diff := cmp.Diff(normalResult, projectedResult, comparer); diff != "" {
			// topk/bottomk pick an arbitrary series among equal values, so a
			// different pick with an identical value multiset is not a
			// divergence the projection caused.
			if sameValues(normalResult, projectedResult) {
				ties++
				continue
			}
			failures = append(failures, failure{query, "result mismatch", diff})
		}
	}

	tokens := []string{"label_replace", "label_join", "timestamp(", "[5m:", "topk", "count_values", "absent", "histogram_quantile", "@ ", " on ", " ignoring ", "group_left", "group_right", "sort", "limitk", "quantile"}
	counts := make(map[string]int)
	for q := range seen {
		for _, tok := range tokens {
			if strings.Contains(q, tok) {
				counts[tok]++
			}
		}
	}
	for _, tok := range tokens {
		t.Logf("corpus contains %q: %d", tok, counts[tok])
	}
	t.Logf("seed=%d runs=%d unique=%d ties=%d failures=%d", seed, testRuns, len(seen), ties, len(failures))
	byKind := make(map[string]int)
	for _, f := range failures {
		byKind[f.kind+": "+f.detail]++
	}
	kinds := make([]string, 0, len(byKind))
	for k := range byKind {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)
	for _, k := range kinds {
		t.Logf("%4d x %s", byKind[k], k)
	}
	for i, f := range failures {
		if i >= 12 {
			t.Logf("... %d more failures", len(failures)-i)
			break
		}
		t.Errorf("[%s] %s\n%s", f.kind, f.query, truncate(f.detail, 1200))
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return fmt.Sprintf("%s... (%d bytes)", s[:n], len(s))
}

// sameValues reports whether both results carry the same multiset of sample
// values, i.e. they differ only in which of several equally ranked series was
// selected.
func sameValues(a, b *promql.Result) bool {
	va, err := values(a)
	if err != nil {
		return false
	}
	vb, err := values(b)
	if err != nil {
		return false
	}
	if len(va) != len(vb) || len(va) == 0 {
		return false
	}
	sort.Float64s(va)
	sort.Float64s(vb)
	return slices.Equal(va, vb)
}

func values(r *promql.Result) ([]float64, error) {
	var out []float64
	switch v := r.Value.(type) {
	case promql.Vector:
		for _, s := range v {
			out = append(out, s.F)
		}
	case promql.Matrix:
		for _, s := range v {
			for _, p := range s.Floats {
				out = append(out, p.F)
			}
		}
	default:
		return nil, errNotComparable
	}
	return out, nil
}

var errNotComparable = errors.New("result type not comparable by value")
