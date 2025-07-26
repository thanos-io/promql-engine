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
	for i := range totalOperators {
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

func assertSeriesExecutionTimeNonZero(t *testing.T, got *engine.AnalyzeOutputNode) bool {
	if got != nil {
		if got.OperatorTelemetry.SeriesExecutionTime() <= 0 {
			t.Errorf("expected non-zero SeriesExecutionTime for Operator, got %s ", got.OperatorTelemetry.SeriesExecutionTime())
			return false
		}
		for i := range got.Children {
			child := got.Children[i]
			return got.OperatorTelemetry.SeriesExecutionTime() > 0 && assertSeriesExecutionTimeNonZero(t, child)
		}
	}
	return true
}

func assertNextExecutionTimeNonZero(t *testing.T, got *engine.AnalyzeOutputNode) bool {
	if got != nil {
		if got.OperatorTelemetry.NextExecutionTime() <= 0 {
			t.Errorf("expected non-zero NextExecutionTime for Operator, got %s ", got.OperatorTelemetry.NextExecutionTime())
			return false
		}
		for i := range got.Children {
			child := got.Children[i]
			return got.OperatorTelemetry.NextExecutionTime() > 0 && assertNextExecutionTimeNonZero(t, child)
		}
	}
	return true
}

// getMaxSeriesCount gets the max series count from the explain output node tree.
func getMaxSeriesCount(got *engine.AnalyzeOutputNode) int {
	maxSeriesCount := 0
	if got != nil {
		maxSeriesCount = got.OperatorTelemetry.MaxSeriesCount()
		for i := range got.Children {
			child := got.Children[i]
			maxSeriesCount = max(maxSeriesCount, getMaxSeriesCount(child))
		}
	}
	return maxSeriesCount
}

func TestQueryAnalyze(t *testing.T) {
	opts := promql.EngineOpts{Timeout: 1 * time.Hour}
	seriesList := []storage.Series{
		storage.MockSeries(
			[]int64{240, 270, 300, 600, 630, 660},
			[]float64{1, 2, 3, 4, 5, 6},
			[]string{labels.MetricName, "foo"},
		),
		storage.MockSeries(
			[]int64{240, 270, 300, 600, 630, 660},
			[]float64{1, 2, 3, 4, 5, 6},
			[]string{labels.MetricName, "http_requests_total", "pod", "nginx-1"},
		),
		storage.MockSeries(
			[]int64{240, 270, 300, 600, 630, 660},
			[]float64{1, 2, 3, 4, 5, 6},
			[]string{labels.MetricName, "http_requests_total", "pod", "nginx-2"},
		),
		storage.MockSeries(
			[]int64{240, 270, 300, 600, 630, 660},
			[]float64{1, 2, 3, 4, 5, 6},
			[]string{labels.MetricName, "http_requests_total", "pod", "nginx-3"},
		),
	}

	start := time.Unix(0, 0)
	end := time.Unix(1000, 0)

	for _, tc := range []struct {
		query          string
		maxSeriesCount int
	}{
		{
			query:          `foo`,
			maxSeriesCount: 1,
		},
		{
			query:          `time()`,
			maxSeriesCount: 0,
		},
		{
			query:          `sum by (job) (foo)`,
			maxSeriesCount: 1,
		},
		{
			query:          `rate(http_requests_total[30s]) > bool 0`,
			maxSeriesCount: 3,
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

				query, err = ng.NewInstantQuery(ctx, storageWithSeries(seriesList...), nil, tc.query, start)
				testutil.Ok(t, err)

				queryResults := query.Exec(context.Background())
				testutil.Ok(t, queryResults.Err)

				explainableQuery := query.(engine.ExplainableQuery)

				testutil.Assert(t, assertExecutionTimeNonZero(t, explainableQuery.Analyze()))
				testutil.Assert(t, assertSeriesExecutionTimeNonZero(t, explainableQuery.Analyze()))
				testutil.Assert(t, assertNextExecutionTimeNonZero(t, explainableQuery.Analyze()))

				testutil.Equals(t, tc.maxSeriesCount, getMaxSeriesCount(explainableQuery.Analyze()))

				query, err = ng.NewRangeQuery(ctx, storageWithSeries(seriesList...), nil, tc.query, start, end, 30*time.Second)
				testutil.Ok(t, err)

				queryResults = query.Exec(context.Background())
				testutil.Ok(t, queryResults.Err)

				explainableQuery = query.(engine.ExplainableQuery)
				testutil.Assert(t, assertExecutionTimeNonZero(t, explainableQuery.Analyze()))
				testutil.Assert(t, assertSeriesExecutionTimeNonZero(t, explainableQuery.Analyze()))
				testutil.Assert(t, assertNextExecutionTimeNonZero(t, explainableQuery.Analyze()))
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
	expected := `[duplicateLabelCheck]: max_series: 2 total_samples: 0 peak_samples: 0
|---[concurrent(buff=2)]: max_series: 2 total_samples: 0 peak_samples: 0
|   |---[aggregate] sum by ([pod]): max_series: 2 total_samples: 0 peak_samples: 0
|   |   |---[duplicateLabelCheck]: max_series: 2 total_samples: 0 peak_samples: 0
|   |   |   |---[coalesce]: max_series: 2 total_samples: 0 peak_samples: 0
|   |   |   |   |---[concurrent(buff=2)]: max_series: 1 total_samples: 0 peak_samples: 0
|   |   |   |   |   |---[matrixSelector] rate({[__name__="http_requests_total"]}[10m0s] 0 mod 2): max_series: 1 total_samples: 1010 peak_samples: 200
|   |   |   |   |---[concurrent(buff=2)]: max_series: 1 total_samples: 0 peak_samples: 0
|   |   |   |   |   |---[matrixSelector] rate({[__name__="http_requests_total"]}[10m0s] 1 mod 2): max_series: 1 total_samples: 1010 peak_samples: 200
`
	require.EqualValues(t, expected, result)
}

func renderAnalysisTree(node *engine.AnalyzeOutputNode, level int) string {
	var result strings.Builder

	totalSamples := int64(0)
	seriesCount := node.OperatorTelemetry.MaxSeriesCount()
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

	result.WriteString(fmt.Sprintf("%s: max_series: %d total_samples: %d peak_samples: %d\n", node.OperatorTelemetry.String(), seriesCount, totalSamples, peakSamples))
	for _, child := range node.Children {
		result.WriteString(renderAnalysisTree(child, level+1))
	}

	return result.String()
}

func TestAnalyzPeak(t *testing.T) {
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
		"sum(rate(http_requests_total[10m])) by (pod)",
		time.Unix(minT, 0),
		time.Unix(maxT, 0),
		60*time.Second,
	)
	testutil.Ok(t, err)
	queryResults = rangeQry.Exec(context.Background())
	testutil.Ok(t, queryResults.Err)

	explainableQuery = rangeQry.(engine.ExplainableQuery)
	analyzeOutput = explainableQuery.Analyze()

	t.Logf("value of peak = %v", analyzeOutput.PeakSamples())
	require.Equal(t, int64(200), analyzeOutput.PeakSamples())

	result := renderAnalysisTree(analyzeOutput, 0)
	expected := `[duplicateLabelCheck]: max_series: 2 total_samples: 0 peak_samples: 0
|---[concurrent(buff=2)]: max_series: 2 total_samples: 0 peak_samples: 0
|   |---[aggregate] sum by ([pod]): max_series: 2 total_samples: 0 peak_samples: 0
|   |   |---[duplicateLabelCheck]: max_series: 2 total_samples: 0 peak_samples: 0
|   |   |   |---[coalesce]: max_series: 2 total_samples: 0 peak_samples: 0
|   |   |   |   |---[concurrent(buff=2)]: max_series: 1 total_samples: 0 peak_samples: 0
|   |   |   |   |   |---[matrixSelector] rate({[__name__="http_requests_total"]}[10m0s] 0 mod 2): max_series: 1 total_samples: 1010 peak_samples: 200
|   |   |   |   |---[concurrent(buff=2)]: max_series: 1 total_samples: 0 peak_samples: 0
|   |   |   |   |   |---[matrixSelector] rate({[__name__="http_requests_total"]}[10m0s] 1 mod 2): max_series: 1 total_samples: 1010 peak_samples: 200
`
	require.EqualValues(t, expected, result)
}
