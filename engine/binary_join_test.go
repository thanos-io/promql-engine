// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/engine"
	"github.com/thanos-io/promql-engine/execution/binary"
	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/query"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/promql/promqltest"
)

// These tests assert that binary vector-matching joins (group_left, group_right,
// one-to-one) return results identical to the Prometheus reference engine across
// a range of cardinality shapes. They cover the output-series identity logic in
// execution/binary, including the join-table de-duplication, which must be
// lossless: identical values, series identity, and ordering.
//
// The interesting case is temporal churn: many series on one side share the same
// join signature (e.g. the same namespace,pod) but do NOT overlap in time (e.g. a
// pod recreated under new uids). At any evaluation step only one is present, so
// the join is a legal many-to-one, yet the static series set seen at init time
// contains all of them.
//
// IMPORTANT: a `_` gap alone does NOT make two series non-overlapping, because the
// 5-minute lookback delta bridges the gap and resurrects the old sample at later
// steps. To genuinely terminate a series (as a real pod restart does) we write an
// explicit `stale` marker, which stops lookback. The replacement series then
// begins at the next step, so every step has at most one active series per join
// key and the join is a legal many-to-one at all times.

const joinQuery = `kube_pod_status_phase{phase="Failed"} * on (namespace, pod) group_left (node) kube_pod_info`

// joinRangeStart/End/Step span both temporal halves of the churn data (7 steps:
// t=0,30,60,90,120,150,180). The churn happens at t=90.
var (
	joinRangeStart = time.Unix(0, 0)
	joinRangeEnd   = time.Unix(180, 0)
	joinRangeStep  = 30 * time.Second
)

// assertEnginesAgree runs the same range query through the Thanos engine and the
// Prometheus reference engine over identical data and asserts the results match.
// It returns the Thanos result value so callers can make additional assertions
// (e.g. on the number of output series).
func assertEnginesAgree(t *testing.T, load, query string, start, end time.Time, step time.Duration) parser.Value {
	t.Helper()

	opts := promql.EngineOpts{
		Timeout:              1 * time.Hour,
		MaxSamples:           1e10,
		EnableNegativeOffset: true,
		EnableAtModifier:     true,
	}

	storage := promqltest.LoadedStorage(t, load)
	defer storage.Close()

	ctx := context.Background()

	thanosEngine := engine.New(engine.Opts{EngineOpts: opts})
	thanosQry, err := thanosEngine.NewRangeQuery(ctx, storage, nil, query, start, end, step)
	testutil.Ok(t, err)
	defer thanosQry.Close()
	thanosResult := thanosQry.Exec(ctx)
	testutil.Ok(t, thanosResult.Err, "thanos engine failed for query: %s", query)

	promEngine := promql.NewEngine(opts)
	promQry, err := promEngine.NewRangeQuery(ctx, storage, nil, query, start, end, step)
	testutil.Ok(t, err)
	defer promQry.Close()
	promResult := promQry.Exec(ctx)
	testutil.Ok(t, promResult.Err, "prometheus engine failed for query: %s", query)

	testutil.WithGoCmp(comparer).Equals(t, promResult, thanosResult, queryExplanation(thanosQry))

	return thanosResult.Value
}

// TestBinaryJoin_ManyToOneChurn is the core case: many low-card series share the
// SAME group_left value (node="n1"). Each pod is recreated under a new uid (a
// non-grouped label) over time, so several kube_pod_info series share the
// (namespace,pod) join signature but never overlap in time. The join must yield
// one output series per pod and match Prometheus.
func TestBinaryJoin_ManyToOneChurn(t *testing.T) {
	t.Parallel()

	// u1 runs t=0..60 then is terminated by a stale marker at t=90; u2 starts at
	// t=90. Both replicas share node="n1", so the join key resolves to a single
	// output series per pod across the whole range.
	load := `load 30s
	    kube_pod_status_phase{namespace="ns1", pod="web", phase="Failed"} 1 1 1 1 1 1 1
	    kube_pod_status_phase{namespace="ns1", pod="api", phase="Failed"} 1 1 1 1 1 1 1
	    kube_pod_info{namespace="ns1", pod="web", node="n1", uid="u1"} 1 1 1 stale
	    kube_pod_info{namespace="ns1", pod="web", node="n1", uid="u2"} _ _ _ 1 1 1 1
	    kube_pod_info{namespace="ns1", pod="api", node="n1", uid="u3"} 1 1 1 stale
	    kube_pod_info{namespace="ns1", pod="api", node="n1", uid="u4"} _ _ _ 1 1 1 1`

	value := assertEnginesAgree(t, load, joinQuery, joinRangeStart, joinRangeEnd, joinRangeStep)

	// Two pods, all churned replicas on node n1 -> exactly two output series.
	m, ok := value.(promql.Matrix)
	testutil.Assert(t, ok, "expected a matrix result, got %T", value)
	testutil.Equals(t, 2, m.Len(), "churned replicas sharing one group_left value must yield one series per pod")
}

// TestBinaryJoin_ManyToOneDistinctGroupValue ensures distinct group_left values
// are NOT collapsed: the same pod is seen on node n1 and then on node n2 (distinct
// values, non-overlapping in time). The two node values must produce two output
// series.
func TestBinaryJoin_ManyToOneDistinctGroupValue(t *testing.T) {
	t.Parallel()

	// node n1 runs t=0..60 then a stale marker at t=90 terminates it; node n2
	// starts at t=90. Distinct group_left values must remain two separate series.
	load := `load 30s
	    kube_pod_status_phase{namespace="ns1", pod="web", phase="Failed"} 1 1 1 1 1 1 1
	    kube_pod_info{namespace="ns1", pod="web", node="n1"} 1 1 1 stale
	    kube_pod_info{namespace="ns1", pod="web", node="n2"} _ _ _ 1 1 1 1`

	value := assertEnginesAgree(t, load, joinQuery, joinRangeStart, joinRangeEnd, joinRangeStep)

	// Distinct group_left values (node n1, node n2) must NOT be collapsed.
	m, ok := value.(promql.Matrix)
	testutil.Assert(t, ok, "expected a matrix result, got %T", value)
	testutil.Equals(t, 2, m.Len(), "distinct group_left values must be preserved as separate series")
}

// TestBinaryJoin_GroupLeftEmpty covers group_left() with an empty group list.
// With no copied labels, a churned pod resolves to a single output series
// carrying only the high-card labels.
func TestBinaryJoin_GroupLeftEmpty(t *testing.T) {
	t.Parallel()

	load := `load 30s
	    kube_pod_status_phase{namespace="ns1", pod="web", phase="Failed"} 1 1 1 1 1 1 1
	    kube_pod_info{namespace="ns1", pod="web", node="n1", uid="u1"} 1 1 1 stale
	    kube_pod_info{namespace="ns1", pod="web", node="n1", uid="u2"} _ _ _ 1 1 1 1`

	query := `kube_pod_status_phase{phase="Failed"} * on (namespace, pod) group_left () kube_pod_info`

	value := assertEnginesAgree(t, load, query, joinRangeStart, joinRangeEnd, joinRangeStep)

	m, ok := value.(promql.Matrix)
	testutil.Assert(t, ok, "expected a matrix result, got %T", value)
	testutil.Equals(t, 1, m.Len(), "empty group_left with one join key must yield a single series")
}

// TestBinaryJoin_OneToOne covers the plain one-to-one `a * on(...) b` path.
// Results must match Prometheus.
func TestBinaryJoin_OneToOne(t *testing.T) {
	t.Parallel()

	load := `load 30s
	    metric_a{namespace="ns1", pod="web"} 1 2 3 4 5 6
	    metric_a{namespace="ns1", pod="api"} 2 4 6 8 10 12
	    metric_b{namespace="ns1", pod="web"} 10 10 10 10 10 10
	    metric_b{namespace="ns1", pod="api"} 100 100 100 100 100 100`

	query := `metric_a * on (namespace, pod) metric_b`

	value := assertEnginesAgree(t, load, query, time.Unix(0, 0), time.Unix(150, 0), 30*time.Second)

	m, ok := value.(promql.Matrix)
	testutil.Assert(t, ok, "expected a matrix result, got %T", value)
	testutil.Equals(t, 2, m.Len(), "one-to-one join must yield one series per matched pair")
}

// TestBinaryJoin_GroupRight covers the group_right path (CardOneToMany), where the
// engine swaps sides so the LHS (kube_pod_info) is the de-duplicated side. Its
// churned replicas (same node, non-overlapping in time) must yield one output
// series per pod, matching the Prometheus engine.
func TestBinaryJoin_GroupRight(t *testing.T) {
	t.Parallel()

	load := `load 30s
	    kube_pod_status_phase{namespace="ns1", pod="web", phase="Failed"} 1 1 1 1 1 1 1
	    kube_pod_info{namespace="ns1", pod="web", node="n1", uid="u1"} 1 1 1 stale
	    kube_pod_info{namespace="ns1", pod="web", node="n1", uid="u2"} _ _ _ 1 1 1 1`

	query := `kube_pod_info * on (namespace, pod) group_right (node) kube_pod_status_phase{phase="Failed"}`

	value := assertEnginesAgree(t, load, query, joinRangeStart, joinRangeEnd, joinRangeStep)

	m, ok := value.(promql.Matrix)
	testutil.Assert(t, ok, "expected a matrix result, got %T", value)
	testutil.Equals(t, 1, m.Len(), "group_right churn must yield one series per pod")
}

// benchMockOperator is a minimal model.VectorOperator that replays a fixed set of
// series and pre-built step vectors. It is used to drive binary.NewVectorOperator
// in benchmarks without standing up a full execution pipeline.
type benchMockOperator struct {
	series   []labels.Labels
	stepVecs []model.StepVector
	cur      int
}

func (m *benchMockOperator) Series(context.Context) ([]labels.Labels, error) {
	return m.series, nil
}

func (m *benchMockOperator) Next(_ context.Context, buf []model.StepVector) (int, error) {
	if m.cur >= len(m.stepVecs) {
		return 0, nil
	}
	n := 0
	for m.cur < len(m.stepVecs) && n < len(buf) {
		buf[n] = m.stepVecs[m.cur]
		m.cur++
		n++
	}
	return n, nil
}

func (m *benchMockOperator) Explain() []model.VectorOperator { return nil }

func (m *benchMockOperator) String() string { return "benchMock" }

// BenchmarkBinaryJoinManyToOneDuplicated measures the heavy-duplication regime:
// every (namespace,pod) join key is backed by many low-card replica series that
// share the same group_left value (the pod's node) but are active at disjoint
// steps (temporal churn). The static low-card series set is pods*replicas while
// the result holds only one series per pod — the shape the join-table
// de-duplication targets.
func BenchmarkBinaryJoinManyToOneDuplicated(b *testing.B) {
	benchmarkBinaryJoin(b, 5000, 50, 10)
}

// BenchmarkBinaryJoinManyToOneDistinct measures the distinct-key regime: each
// (namespace,pod) has exactly one low-card series (and its own node), so no join
// signature repeats and the de-duplication is gated off. Paired with the
// duplicated benchmark it shows the de-duplication helps when it can and stays
// allocation-free when it cannot.
func BenchmarkBinaryJoinManyToOneDistinct(b *testing.B) {
	benchmarkBinaryJoin(b, 50000, 1, 10)
}

// benchmarkBinaryJoin drives a many-to-one group_left(node) join over `pods`
// high-card series and pods*replicas low-card series across `steps`, with one
// low-card replica active per pod per step (a legal many-to-one at every step).
// Each pod has its own node, so a pod's replicas collapse to one output series
// while distinct pods stay distinct.
func benchmarkBinaryJoin(b *testing.B, pods, replicas, steps int) {
	// High-card side: one status_phase series per pod, present at every step.
	highSeries := make([]labels.Labels, pods)
	for p := range pods {
		highSeries[p] = labels.FromStrings(
			labels.MetricName, "kube_pod_status_phase",
			"namespace", "ns1",
			"pod", "pod"+strconv.Itoa(p),
			"phase", "Failed",
		)
	}

	// Low-card side: pods*replicas info series. A pod's replicas share its
	// (namespace,pod) join key and its node; distinct pods have distinct nodes.
	lowSeries := make([]labels.Labels, pods*replicas)
	for p := range pods {
		for r := range replicas {
			lowSeries[p*replicas+r] = labels.FromStrings(
				labels.MetricName, "kube_pod_info",
				"namespace", "ns1",
				"pod", "pod"+strconv.Itoa(p),
				"node", "node"+strconv.Itoa(p),
				"uid", "uid"+strconv.Itoa(p)+"-"+strconv.Itoa(r),
			)
		}
	}

	// Pre-build step vectors. At each step exactly one replica per pod is active
	// (rotating by step), keeping the match a legal many-to-one at every step.
	highSteps := make([]model.StepVector, steps)
	lowSteps := make([]model.StepVector, steps)
	for s := range steps {
		ts := int64(s) * joinRangeStep.Milliseconds()

		hsv := model.StepVector{T: ts}
		for p := range pods {
			hsv.SampleIDs = append(hsv.SampleIDs, uint64(p))
			hsv.Samples = append(hsv.Samples, 1)
		}
		highSteps[s] = hsv

		lsv := model.StepVector{T: ts}
		active := s % replicas
		for p := range pods {
			lsv.SampleIDs = append(lsv.SampleIDs, uint64(p*replicas+active))
			lsv.Samples = append(lsv.Samples, 1)
		}
		lowSteps[s] = lsv
	}

	matching := &parser.VectorMatching{
		Card:           parser.CardManyToOne,
		On:             true,
		MatchingLabels: []string{"namespace", "pod"},
		Include:        []string{"node"},
	}
	opts := &query.Options{StepsBatch: steps}

	b.ReportAllocs()

	for b.Loop() {
		lhs := &benchMockOperator{series: highSeries, stepVecs: highSteps}
		rhs := &benchMockOperator{series: lowSeries, stepVecs: lowSteps}

		op, err := binary.NewVectorOperator(lhs, rhs, matching, parser.MUL, false, opts)
		if err != nil {
			b.Fatal(err)
		}

		if _, err := op.Series(context.Background()); err != nil {
			b.Fatal(err)
		}

		buf := make([]model.StepVector, steps)
		for {
			n, err := op.Next(context.Background(), buf)
			if err != nil {
				b.Fatal(err)
			}
			if n == 0 {
				break
			}
		}
	}
}
