// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/api"
	"github.com/thanos-io/promql-engine/engine"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/promqltest"
	"github.com/prometheus/prometheus/storage"
	"github.com/stretchr/testify/require"
)

func TestQueryExplain(t *testing.T) {
	t.Parallel()
	opts := promql.EngineOpts{Timeout: 1 * time.Hour}
	series := storage.MockSeries(
		[]int64{240, 270, 300, 600, 630, 660},
		[]float64{1, 2, 3, 4, 5, 6},
		[]string{labels.MetricName, "foo"},
	)

	start := time.Unix(0, 0)
	end := time.Unix(1000, 0)

	// Calculate concurrencyOperators according to max available CPUs.
	totalOperators := runtime.GOMAXPROCS(0) / 2
	var concurrencyOperators []engine.ExplainOutputNode
	for i := 0; i < totalOperators; i++ {
		concurrencyOperators = append(concurrencyOperators, engine.ExplainOutputNode{
			OperatorName: "[concurrent(buff=2)]", Children: []engine.ExplainOutputNode{
				{OperatorName: fmt.Sprintf("[vectorSelector] {[__name__=\"foo\"]} %d mod %d", i, totalOperators)},
			},
		})
	}

	for _, tc := range []struct {
		query    string
		expected *engine.ExplainOutputNode
	}{
		{
			query: `time()`,
			expected: &engine.ExplainOutputNode{OperatorName: "[duplicateLabelCheck]", Children: []engine.ExplainOutputNode{
				{
					OperatorName: "[noArgFunction]",
					Children:     nil,
				},
			}},
		},
		{
			query:    `foo`,
			expected: &engine.ExplainOutputNode{OperatorName: "[coalesce]", Children: concurrencyOperators},
		},
		{
			query: `sum by (job) (foo)`,
			expected: &engine.ExplainOutputNode{
				OperatorName: "[duplicateLabelCheck]",
				Children: []engine.ExplainOutputNode{
					{
						OperatorName: "[concurrent(buff=2)]", Children: []engine.ExplainOutputNode{
							{
								OperatorName: "[aggregate] sum by ([job])", Children: []engine.ExplainOutputNode{
									{
										OperatorName: "[coalesce]",
										Children:     concurrencyOperators,
									},
								},
							},
						},
					},
				},
			},
		},
	} {
		{
			t.Run(tc.query, func(t *testing.T) {
				ng := engine.New(engine.Opts{EngineOpts: opts})
				ctx := context.Background()

				var (
					query promql.Query
					err   error
				)

				query, err = ng.NewInstantQuery(ctx, storageWithSeries(series), nil, tc.query, start)
				testutil.Ok(t, err)

				explainableQuery := query.(engine.ExplainableQuery)
				testutil.Equals(t, tc.expected, explainableQuery.Explain())

				query, err = ng.NewRangeQuery(ctx, storageWithSeries(series), nil, tc.query, start, end, 30*time.Second)
				testutil.Ok(t, err)

				explainableQuery = query.(engine.ExplainableQuery)
				testutil.Equals(t, tc.expected, explainableQuery.Explain())
			})
		}
	}
}

func assertExecutionTimeNonZero(t *testing.T, got *engine.AnalyzeOutputNode) bool {
	if got != nil {
		if got.OperatorTelemetry.ExecutionTimeTaken() <= 0 {
			t.Errorf("expected non-zero ExecutionTime for Operator, got %s ", got.OperatorTelemetry.ExecutionTimeTaken())
			return false
		}
		for i := range got.Children {
			child := got.Children[i]
			return got.OperatorTelemetry.ExecutionTimeTaken() > 0 && assertExecutionTimeNonZero(t, child)
		}
	}
	return true
}

func TestQueryAnalyze(t *testing.T) {
	opts := promql.EngineOpts{Timeout: 1 * time.Hour}
	series := storage.MockSeries(
		[]int64{240, 270, 300, 600, 630, 660},
		[]float64{1, 2, 3, 4, 5, 6},
		[]string{labels.MetricName, "foo"},
	)

	start := time.Unix(0, 0)
	end := time.Unix(1000, 0)

	for _, tc := range []struct {
		query string
	}{
		{
			query: `foo`,
		},
		{
			query: `time()`,
		},
		{
			query: `sum by (job) (foo)`,
		},
		{
			query: `rate(http_requests_total[30s]) > bool 0`,
		},
	} {
		{
			t.Run(tc.query, func(t *testing.T) {
				t.Parallel()
				ng := engine.New(engine.Opts{EngineOpts: opts, EnableAnalysis: true})
				ctx := context.Background()

				var (
					query promql.Query
					err   error
				)

				query, err = ng.NewInstantQuery(ctx, storageWithSeries(series), nil, tc.query, start)
				testutil.Ok(t, err)

				queryResults := query.Exec(context.Background())
				testutil.Ok(t, queryResults.Err)

				explainableQuery := query.(engine.ExplainableQuery)

				testutil.Assert(t, assertExecutionTimeNonZero(t, explainableQuery.Analyze()))

				query, err = ng.NewRangeQuery(ctx, storageWithSeries(series), nil, tc.query, start, end, 30*time.Second)
				testutil.Ok(t, err)

				queryResults = query.Exec(context.Background())
				testutil.Ok(t, queryResults.Err)

				explainableQuery = query.(engine.ExplainableQuery)
				testutil.Assert(t, assertExecutionTimeNonZero(t, explainableQuery.Analyze()))
			})
		}
	}
}

func TestAnalyzeOutputNode_Samples(t *testing.T) {
	t.Parallel()
	ng := engine.New(engine.Opts{EngineOpts: promql.EngineOpts{Timeout: 1 * time.Hour}, EnableAnalysis: true, DecodingConcurrency: 2})
	ctx := context.Background()

	load := `load 30s
				http_requests_total{pod="nginx-1"} 1+1x100
				http_requests_total{pod="nginx-2"} 1+1x100`

	tstorage := promqltest.LoadedStorage(t, load)
	defer tstorage.Close()
	minT := tstorage.Head().Meta().MinTime
	maxT := tstorage.Head().Meta().MaxTime

	query, err := ng.NewInstantQuery(ctx, tstorage, nil, "http_requests_total", time.Unix(0, 0))
	testutil.Ok(t, err)
	queryResults := query.Exec(context.Background())
	testutil.Ok(t, queryResults.Err)
	explainableQuery := query.(engine.ExplainableQuery)
	analyzeOutput := explainableQuery.Analyze()
	require.Greater(t, analyzeOutput.PeakSamples(), int64(0))
	require.Greater(t, analyzeOutput.TotalSamples(), int64(0))

	rangeQry, err := ng.NewRangeQuery(
		ctx,
		tstorage,
		promql.NewPrometheusQueryOpts(false, 0),
		"sum(rate(http_requests_total[10m])) by (pod)", // Increase range to 60 minutes
		time.Unix(minT, 0),
		time.Unix(maxT, 0),
		60*time.Second,
	)
	testutil.Ok(t, err)
	queryResults = rangeQry.Exec(context.Background())
	testutil.Ok(t, queryResults.Err)

	explainableQuery = rangeQry.(engine.ExplainableQuery)
	analyzeOutput = explainableQuery.Analyze()
	require.Greater(t, analyzeOutput.PeakSamples(), int64(0))
	require.Greater(t, analyzeOutput.TotalSamples(), int64(0))
	result := renderAnalysisTree(analyzeOutput, 0)
	expected := `[duplicateLabelCheck]: 0 peak: 0
|---[concurrent(buff=2)]: 0 peak: 0
|   |---[aggregate] sum by ([pod]): 0 peak: 0
|   |   |---[duplicateLabelCheck]: 0 peak: 0
|   |   |   |---[coalesce]: 0 peak: 0
|   |   |   |   |---[concurrent(buff=2)]: 0 peak: 0
|   |   |   |   |   |---[matrixSelector] rate({[__name__="http_requests_total"]}[10m0s] 0 mod 2): 1010 peak: 20
|   |   |   |   |---[concurrent(buff=2)]: 0 peak: 0
|   |   |   |   |   |---[matrixSelector] rate({[__name__="http_requests_total"]}[10m0s] 1 mod 2): 1010 peak: 20
`
	require.EqualValues(t, expected, result)
}

func renderAnalysisTree(node *engine.AnalyzeOutputNode, level int) string {
	var result strings.Builder

	totalSamples := int64(0)
	samples := node.OperatorTelemetry.Samples()
	if samples != nil {
		totalSamples = samples.TotalSamples
	}

	peakSamples := int64(0)
	if samples != nil {
		peakSamples = int64(samples.PeakSamples)
	}

	if level > 0 {
		result.WriteString(strings.Repeat("|   ", level-1) + "|---")
	}

	result.WriteString(fmt.Sprintf("%s: %d peak: %d\n", node.OperatorTelemetry.String(), totalSamples, peakSamples))
	for _, child := range node.Children {
		result.WriteString(renderAnalysisTree(child, level+1))
	}

	return result.String()
}

func TestQueryAnalyzeWithRemoteExecution(t *testing.T) {
	t.Parallel()
	opts := promql.EngineOpts{Timeout: 1 * time.Hour}

	var loadStrings []string
	loadStrings = append(loadStrings, "load 30s")

	clusters := []string{"cluster-0", "cluster-1", "cluster-2"}
	pods := []string{"pod-0", "pod-1", "pod-2", "pod-3", "pod-4"}
	namespaces := []string{"ns-0", "ns-1"}

	for _, cluster := range clusters {
		for _, pod := range pods {
			for _, ns := range namespaces {
				loadStrings = append(loadStrings, fmt.Sprintf(
					`kube_pod_info{k8s_cluster="%s", pod="%s", namespace="%s"} 1+1x14`,
					cluster, pod, ns,
				))
			}
		}
	}
	loadTestString := strings.Join(loadStrings, "\n")

	mockStorage := promqltest.LoadedStorage(t, loadTestString)
	defer mockStorage.Close()

	numRemoteEngines := 3
	remoteEngines := make([]api.RemoteEngine, 0, numRemoteEngines)
	for i := 0; i < numRemoteEngines; i++ {
		minT := mockStorage.Head().Meta().MinTime
		maxT := mockStorage.Head().Meta().MaxTime

		remoteEngines = append(remoteEngines, engine.NewRemoteEngine(
			engine.Opts{EngineOpts: opts, EnableAnalysis: true},
			mockStorage,
			minT,
			maxT,
			[]labels.Labels{labels.FromStrings("k8s_cluster", fmt.Sprintf("cluster-%d", i))},
		))
	}

	endpoints := api.NewStaticEndpoints(remoteEngines)
	distEngine := engine.NewDistributedEngine(engine.Opts{EngineOpts: opts, EnableAnalysis: true})

	start := time.Unix(240, 0)
	end := time.Unix(1000, 0)

	query := "sum by (k8s_cluster) (kube_pod_info)"

	t.Run("instant_query", func(t *testing.T) {
		ctx := context.Background()
		q, err := distEngine.MakeInstantQuery(ctx, mockStorage, endpoints, nil, query, start)
		testutil.Ok(t, err)

		results := q.Exec(ctx)
		testutil.Ok(t, results.Err)

		explainableQuery := q.(engine.ExplainableQuery)
		analysis := explainableQuery.Analyze()

		result := renderAnalysisTree(analysis, 0)
		t.Log(result)

		testutil.Assert(t, analysis != nil, "analysis should not be nil")
		testutil.Assert(t, assertExecutionTimeNonZero(t, analysis))

		testutil.Assert(t, analysis.TotalSamples() > 0, "total samples from root analysis node should be greater than 0")

		validateOperators(t, analysis, numRemoteEngines)
	})

	t.Run("range_query", func(t *testing.T) {
		ctx := context.Background()
		q, err := distEngine.MakeRangeQuery(ctx, mockStorage, endpoints, nil, query, start, end, 60*time.Second)
		testutil.Ok(t, err)

		results := q.Exec(ctx)
		testutil.Ok(t, results.Err)

		explainableQuery := q.(engine.ExplainableQuery)
		analysis := explainableQuery.Analyze()
		result := renderAnalysisTree(analysis, 0)
		t.Log(result)

		// Validate the structure
		testutil.Assert(t, analysis != nil, "analysis should not be nil")
		testutil.Assert(t, assertExecutionTimeNonZero(t, analysis))
		testutil.Assert(t, analysis.TotalSamples() >= 720, "total samples from root analysis node should be greater than 0")

		validateOperators(t, analysis, numRemoteEngines)
	})
}

// validateOperators checks that the plan contains the expected operators in the right structure
func validateOperators(t *testing.T, node *engine.AnalyzeOutputNode, numRemoteEngines int) {
	var remoteExecNodes []*engine.AnalyzeOutputNode
	var aggregateNodes []*engine.AnalyzeOutputNode
	var concurrentNodes []*engine.AnalyzeOutputNode
	var dedupNodes []*engine.AnalyzeOutputNode

	findOperatorsRecursive(node, "remoteExec", &remoteExecNodes)
	findOperatorsRecursive(node, "aggregate", &aggregateNodes)
	findOperatorsRecursive(node, "concurrent", &concurrentNodes)
	findOperatorsRecursive(node, "dedup", &dedupNodes)

	testutil.Assert(t, len(remoteExecNodes) > 0, "should find at least one remoteExec operator")
	for i, rn := range remoteExecNodes {
		t.Logf("RemoteExec node %d: %s, TotalSamples: %d, PeakSamples: %d", i, rn.OperatorTelemetry.String(), rn.TotalSamples(), rn.PeakSamples())
		testutil.Assert(t, rn.TotalSamples() > 0, fmt.Sprintf("remoteExec node %d (%s) should have > 0 total samples", i, rn.OperatorTelemetry.String()))
	}

	testutil.Assert(t, len(concurrentNodes) > 0, "should find concurrent operator")

	if numRemoteEngines > 1 {
		testutil.Assert(t, len(aggregateNodes) > 0, "should find aggregate operator when numRemoteEngines > 1")
		testutil.Assert(t, len(dedupNodes) > 0, "should find dedup operator when numRemoteEngines > 1")
	} else {
		testutil.Assert(t, len(aggregateNodes) == 0, "should NOT find separate aggregate operator when numRemoteEngines == 1")
		testutil.Assert(t, len(dedupNodes) == 0, "should not find dedup operator when numRemoteEngines == 1")
		// Check if remoteExec itself is doing the aggregation for the single remote engine case
		if len(remoteExecNodes) == 1 {
			testutil.Assert(t, strings.Contains(remoteExecNodes[0].OperatorTelemetry.String(), "sum by"), "remoteExec should be performing sum aggregation for single remote engine")
		}
	}
}

func findOperatorsRecursive(node *engine.AnalyzeOutputNode, operatorNameSubstring string, foundNodes *[]*engine.AnalyzeOutputNode) {
	if node == nil {
		return
	}
	if strings.Contains(node.OperatorTelemetry.String(), operatorNameSubstring) {
		*foundNodes = append(*foundNodes, node)
	}

	for _, child := range node.Children {
		findOperatorsRecursive(child, operatorNameSubstring, foundNodes)
	}
}
