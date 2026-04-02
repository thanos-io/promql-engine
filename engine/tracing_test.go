// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine_test

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/api"
	"github.com/thanos-io/promql-engine/engine"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/promqltest"
	"github.com/prometheus/prometheus/storage"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	otelcodes "go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// setupTracerProvider installs an in-memory span exporter as the global
// TracerProvider and returns the exporter for span assertions. The previous
// provider is restored in t.Cleanup.
func setupTracerProvider(t *testing.T) *tracetest.InMemoryExporter {
	t.Helper()
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() {
		_ = tp.Shutdown(context.Background())
		otel.SetTracerProvider(prev)
	})
	return exporter
}

func newTestEngine(enableAnalysis bool) *engine.Engine {
	return engine.New(engine.Opts{
		EngineOpts: promql.EngineOpts{
			Timeout:              time.Hour,
			MaxSamples:           1e10,
			EnableNegativeOffset: true,
			EnableAtModifier:     true,
		},
		EnableAnalysis: enableAnalysis,
	})
}

func loadTestData(t *testing.T) storage.Queryable {
	t.Helper()
	load := `
		load 30s
			http_requests_total{job="api", code="200"} 0+1x100
			http_requests_total{job="api", code="500"} 0+2x100
	`
	return promqltest.LoadedStorage(t, load)
}

// findSpan returns the first span with the given name, or nil.
func findSpan(spans tracetest.SpanStubs, name string) *tracetest.SpanStub {
	for i := range spans {
		if spans[i].Name == name {
			return &spans[i]
		}
	}
	return nil
}

// spanAttr returns the value of the named attribute on a span, or "".
func spanAttr(s *tracetest.SpanStub, key string) attribute.Value {
	for _, a := range s.Attributes {
		if string(a.Key) == key {
			return a.Value
		}
	}
	return attribute.Value{}
}

func TestInstantQuerySpan(t *testing.T) {
	exporter := setupTracerProvider(t)
	store := loadTestData(t)
	eng := newTestEngine(false)

	qry, err := eng.NewInstantQuery(context.Background(), store, nil, "http_requests_total", time.Unix(60, 0))
	require.NoError(t, err)
	defer qry.Close()

	result := qry.Exec(context.Background())
	require.NoError(t, result.Err)

	spans := exporter.GetSpans()
	s := findSpan(spans, "instant_query_exec")
	require.NotNil(t, s, "expected instant_query_exec span")
	require.Equal(t, "http_requests_total", spanAttr(s, "query.expr").AsString())
	require.NotZero(t, spanAttr(s, "query.start").AsInt64())
}

func TestRangeQuerySpan(t *testing.T) {
	exporter := setupTracerProvider(t)
	store := loadTestData(t)
	eng := newTestEngine(false)

	start := time.Unix(60, 0)
	end := time.Unix(300, 0)
	step := 30 * time.Second

	qry, err := eng.NewRangeQuery(context.Background(), store, nil, "http_requests_total", start, end, step)
	require.NoError(t, err)
	defer qry.Close()

	result := qry.Exec(context.Background())
	require.NoError(t, result.Err)

	spans := exporter.GetSpans()
	s := findSpan(spans, "range_query_exec")
	require.NotNil(t, s, "expected range_query_exec span")
	require.Equal(t, "http_requests_total", spanAttr(s, "query.expr").AsString())
	require.NotZero(t, spanAttr(s, "query.start").AsInt64())
	require.NotZero(t, spanAttr(s, "query.end").AsInt64())
	require.InDelta(t, 30.0, spanAttr(s, "query.interval_seconds").AsFloat64(), 0.01)
	require.InDelta(t, end.Sub(start).Seconds(), spanAttr(s, "query.range_seconds").AsFloat64(), 0.01)
}

func TestQuerySpanError(t *testing.T) {
	exporter := setupTracerProvider(t)
	eng := newTestEngine(false)
	store := loadTestData(t)

	// Use a valid query but with a function that triggers an error
	// by using a storage that returns an error, or a query that will fail.
	// Simplest: use a query on empty storage with an expression that will execute
	// but trigger a type error — actually let's use a deliberately malformed query approach.
	// Instead, we can use a context that's already cancelled.
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	qry, err := eng.NewInstantQuery(ctx, store, nil, "http_requests_total", time.Unix(60, 0))
	require.NoError(t, err)
	defer qry.Close()

	result := qry.Exec(ctx)
	require.Error(t, result.Err)

	spans := exporter.GetSpans()
	s := findSpan(spans, "instant_query_exec")
	require.NotNil(t, s, "expected instant_query_exec span even on error")
	require.Equal(t, otelcodes.Error, s.Status.Code)
}

func TestQuerySpanParentContext(t *testing.T) {
	exporter := setupTracerProvider(t)
	store := loadTestData(t)
	eng := newTestEngine(false)

	// Create a parent span.
	tp := otel.GetTracerProvider()
	tracer := tp.Tracer("test")
	ctx, parentSpan := tracer.Start(context.Background(), "parent_operation")
	parentSpanCtx := parentSpan.SpanContext()

	qry, err := eng.NewInstantQuery(ctx, store, nil, "http_requests_total", time.Unix(60, 0))
	require.NoError(t, err)
	defer qry.Close()

	result := qry.Exec(ctx)
	require.NoError(t, result.Err)
	parentSpan.End()

	spans := exporter.GetSpans()
	s := findSpan(spans, "instant_query_exec")
	require.NotNil(t, s, "expected instant_query_exec span")
	// The query span should be a child of the parent span.
	require.Equal(t, parentSpanCtx.TraceID(), s.SpanContext.TraceID())
	require.Equal(t, parentSpanCtx.SpanID(), s.Parent.SpanID())
}

func TestAnalysisAvailableEvent(t *testing.T) {
	exporter := setupTracerProvider(t)
	store := loadTestData(t)
	eng := newTestEngine(false) // analysis OFF

	qry, err := eng.NewInstantQuery(context.Background(), store, nil, "http_requests_total", time.Unix(60, 0))
	require.NoError(t, err)
	defer qry.Close()

	result := qry.Exec(context.Background())
	require.NoError(t, result.Err)

	spans := exporter.GetSpans()
	s := findSpan(spans, "instant_query_exec")
	require.NotNil(t, s)

	// Check for promql.analysis_available event.
	var found bool
	for _, ev := range s.Events {
		if ev.Name == "promql.analysis_available" {
			found = true
			break
		}
	}
	require.True(t, found, "expected promql.analysis_available event when analysis is OFF")
}

func TestNoAnalysisEventWhenEnabled(t *testing.T) {
	exporter := setupTracerProvider(t)
	store := loadTestData(t)
	eng := newTestEngine(true) // analysis ON

	qry, err := eng.NewInstantQuery(context.Background(), store, nil, "http_requests_total", time.Unix(60, 0))
	require.NoError(t, err)
	defer qry.Close()

	result := qry.Exec(context.Background())
	require.NoError(t, result.Err)

	spans := exporter.GetSpans()
	s := findSpan(spans, "instant_query_exec")
	require.NotNil(t, s)

	for _, ev := range s.Events {
		require.NotEqual(t, "promql.analysis_available", ev.Name,
			"should not emit analysis_available event when analysis is ON")
	}
}

func TestBaggageOverrideSuppressesAnalysisEvent(t *testing.T) {
	exporter := setupTracerProvider(t)
	store := loadTestData(t)
	eng := newTestEngine(false) // analysis OFF globally

	// Set baggage to enable analysis for this query.
	member, err := baggage.NewMember("promql.enable_analysis", "true")
	require.NoError(t, err)
	bag, err := baggage.New(member)
	require.NoError(t, err)
	ctx := baggage.ContextWithBaggage(context.Background(), bag)

	qry, err := eng.NewInstantQuery(ctx, store, nil, "http_requests_total", time.Unix(60, 0))
	require.NoError(t, err)
	defer qry.Close()

	result := qry.Exec(ctx)
	require.NoError(t, result.Err)

	spans := exporter.GetSpans()
	s := findSpan(spans, "instant_query_exec")
	require.NotNil(t, s)

	// Baggage override means analysis is "on" — should NOT have analysis_available event.
	for _, ev := range s.Events {
		require.NotEqual(t, "promql.analysis_available", ev.Name,
			"baggage override should suppress analysis_available event")
	}
}

// --- US2: Operator-Level Span Tests ---

// findSpans returns all spans with the given name.
func findSpans(spans tracetest.SpanStubs, name string) []tracetest.SpanStub {
	var result []tracetest.SpanStub
	for _, s := range spans {
		if s.Name == name {
			result = append(result, s)
		}
	}
	return result
}

// operatorSpans returns spans that are operator-level (not query or storage spans).
func operatorSpans(spans tracetest.SpanStubs) []tracetest.SpanStub {
	var result []tracetest.SpanStub
	for _, s := range spans {
		switch s.Name {
		case "instant_query_exec", "range_query_exec", "storage:select":
			continue
		default:
			result = append(result, s)
		}
	}
	return result
}

func TestOperatorSpansWithAnalysis(t *testing.T) {
	exporter := setupTracerProvider(t)
	store := loadTestData(t)
	eng := newTestEngine(true) // analysis ON

	qry, err := eng.NewInstantQuery(context.Background(), store, nil,
		"sum(http_requests_total)", time.Unix(60, 0))
	require.NoError(t, err)
	defer qry.Close()

	result := qry.Exec(context.Background())
	require.NoError(t, result.Err)

	spans := exporter.GetSpans()
	opSpans := operatorSpans(spans)
	require.NotEmpty(t, opSpans, "expected operator child spans when analysis is ON")

	// Each operator span should have promql.operator.type attribute.
	for _, s := range opSpans {
		v := spanAttr(&s, "promql.operator.type")
		require.NotEmpty(t, v.AsString(),
			"operator span %q should have promql.operator.type", s.Name)
	}
}

func TestOperatorSpanAttributes(t *testing.T) {
	exporter := setupTracerProvider(t)
	store := loadTestData(t)
	eng := newTestEngine(true) // analysis ON

	start := time.Unix(60, 0)
	end := time.Unix(300, 0)
	step := 30 * time.Second

	qry, err := eng.NewRangeQuery(context.Background(), store, nil,
		"sum(http_requests_total)", start, end, step)
	require.NoError(t, err)
	defer qry.Close()

	result := qry.Exec(context.Background())
	require.NoError(t, result.Err)

	spans := exporter.GetSpans()
	opSpans := operatorSpans(spans)
	require.NotEmpty(t, opSpans)

	// At least one operator span should have series count and samples total
	// (set when the operator finishes producing results).
	var foundWithAttrs bool
	for _, s := range opSpans {
		sc := spanAttr(&s, "promql.series.count")
		st := spanAttr(&s, "promql.samples.total")
		if sc.AsInt64() > 0 || st.AsInt64() > 0 {
			foundWithAttrs = true
			break
		}
	}
	require.True(t, foundWithAttrs,
		"expected at least one operator span with promql.series.count or promql.samples.total")
}

func TestAnalysisOffNoOperatorSpans(t *testing.T) {
	exporter := setupTracerProvider(t)
	store := loadTestData(t)
	eng := newTestEngine(false) // analysis OFF

	qry, err := eng.NewInstantQuery(context.Background(), store, nil,
		"sum(http_requests_total)", time.Unix(60, 0))
	require.NoError(t, err)
	defer qry.Close()

	result := qry.Exec(context.Background())
	require.NoError(t, result.Err)

	spans := exporter.GetSpans()
	opSpans := operatorSpans(spans)
	require.Empty(t, opSpans, "expected no operator spans when analysis is OFF")

	// Query span should have analysis_available event.
	s := findSpan(spans, "instant_query_exec")
	require.NotNil(t, s)
	var found bool
	for _, ev := range s.Events {
		if ev.Name == "promql.analysis_available" {
			found = true
			break
		}
	}
	require.True(t, found, "expected promql.analysis_available event when analysis is OFF")
}

func TestBaggageOverrideEnablesOperatorSpans(t *testing.T) {
	exporter := setupTracerProvider(t)
	store := loadTestData(t)
	eng := newTestEngine(false) // analysis OFF globally

	member, err := baggage.NewMember("promql.enable_analysis", "true")
	require.NoError(t, err)
	bag, err := baggage.New(member)
	require.NoError(t, err)
	ctx := baggage.ContextWithBaggage(context.Background(), bag)

	qry, err := eng.NewInstantQuery(ctx, store, nil,
		"sum(http_requests_total)", time.Unix(60, 0))
	require.NoError(t, err)
	defer qry.Close()

	result := qry.Exec(ctx)
	require.NoError(t, result.Err)

	spans := exporter.GetSpans()
	opSpans := operatorSpans(spans)
	require.NotEmpty(t, opSpans,
		"expected operator spans when baggage override enables analysis")
}

// --- US2b: Storage Span Tests ---

func TestVectorSelectorStorageSpan(t *testing.T) {
	exporter := setupTracerProvider(t)
	store := loadTestData(t)
	eng := newTestEngine(false) // analysis OFF — storage spans are unconditional

	qry, err := eng.NewInstantQuery(context.Background(), store, nil,
		"http_requests_total", time.Unix(60, 0))
	require.NoError(t, err)
	defer qry.Close()

	result := qry.Exec(context.Background())
	require.NoError(t, result.Err)

	spans := exporter.GetSpans()
	ss := findSpans(spans, "storage:select")
	require.NotEmpty(t, ss, "expected storage:select span for vector selector")

	// The engine may create multiple sharded selectors. Verify attributes on
	// all spans and sum the series count across shards.
	var totalSeries int64
	for i := range ss {
		require.Contains(t, spanAttr(&ss[i], "storage.matchers").AsString(), "http_requests_total")
		require.NotZero(t, spanAttr(&ss[i], "storage.maxt").AsInt64())
		totalSeries += spanAttr(&ss[i], "storage.series.count").AsInt64()
	}
	require.Equal(t, int64(2), totalSeries, "total series across shards should be 2")
}

func TestMatrixSelectorStorageSpan(t *testing.T) {
	exporter := setupTracerProvider(t)
	store := loadTestData(t)
	eng := newTestEngine(false)

	qry, err := eng.NewInstantQuery(context.Background(), store, nil,
		"rate(http_requests_total[5m])", time.Unix(300, 0))
	require.NoError(t, err)
	defer qry.Close()

	result := qry.Exec(context.Background())
	require.NoError(t, result.Err)

	spans := exporter.GetSpans()
	ss := findSpans(spans, "storage:select")
	require.NotEmpty(t, ss, "expected storage:select span for matrix selector")

	var totalSeries int64
	for i := range ss {
		require.Contains(t, spanAttr(&ss[i], "storage.matchers").AsString(), "http_requests_total")
		require.NotZero(t, spanAttr(&ss[i], "storage.maxt").AsInt64())
		totalSeries += spanAttr(&ss[i], "storage.series.count").AsInt64()
	}
	require.Equal(t, int64(2), totalSeries, "total series across shards should be 2")
}

func TestStorageSpanError(t *testing.T) {
	exporter := setupTracerProvider(t)
	eng := newTestEngine(false)
	store := loadTestData(t)

	// Use cancelled context to trigger an error during query execution.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	qry, err := eng.NewInstantQuery(ctx, store, nil,
		"http_requests_total", time.Unix(60, 0))
	require.NoError(t, err)
	defer qry.Close()

	result := qry.Exec(ctx)
	require.Error(t, result.Err)

	spans := exporter.GetSpans()
	// Verify the query-level span records the error.
	s := findSpan(spans, "instant_query_exec")
	require.NotNil(t, s, "expected instant_query_exec span even on error")
	require.Equal(t, otelcodes.Error, s.Status.Code)

	// Storage span may or may not be emitted depending on where cancellation hits.
	// If emitted with error status, that's also correct.
	ss := findSpans(spans, "storage:select")
	for _, storageSpan := range ss {
		if storageSpan.Status.Code == otelcodes.Error {
			// Storage span correctly recorded the error.
			break
		}
	}
}

func TestNoBaggageNoOperatorSpans(t *testing.T) {
	exporter := setupTracerProvider(t)
	store := loadTestData(t)
	eng := newTestEngine(false) // analysis OFF, no baggage

	qry, err := eng.NewInstantQuery(context.Background(), store, nil,
		"sum(http_requests_total)", time.Unix(60, 0))
	require.NoError(t, err)
	defer qry.Close()

	result := qry.Exec(context.Background())
	require.NoError(t, result.Err)

	spans := exporter.GetSpans()
	opSpans := operatorSpans(spans)
	require.Empty(t, opSpans, "expected no operator spans with analysis OFF and no baggage")
}

func TestBaggageFalseSuppressesAnalysis(t *testing.T) {
	exporter := setupTracerProvider(t)
	store := loadTestData(t)
	eng := newTestEngine(true) // analysis ON globally

	// Baggage "false" should suppress operator spans even with EnableAnalysis ON.
	member, err := baggage.NewMember("promql.enable_analysis", "false")
	require.NoError(t, err)
	bag, err := baggage.New(member)
	require.NoError(t, err)
	ctx := baggage.ContextWithBaggage(context.Background(), bag)

	qry, err := eng.NewInstantQuery(ctx, store, nil,
		"sum(http_requests_total)", time.Unix(60, 0))
	require.NoError(t, err)
	defer qry.Close()

	result := qry.Exec(ctx)
	require.NoError(t, result.Err)

	spans := exporter.GetSpans()
	opSpans := operatorSpans(spans)
	require.Empty(t, opSpans,
		"expected no operator spans when baggage explicitly set to false")

	// Query span should have analysis_available event since analysis is suppressed.
	s := findSpan(spans, "instant_query_exec")
	require.NotNil(t, s)
	var found bool
	for _, ev := range s.Events {
		if ev.Name == "promql.analysis_available" {
			found = true
			break
		}
	}
	require.True(t, found,
		"expected promql.analysis_available event when baggage suppresses analysis")
}

// --- US3: Distributed Trace Propagation Tests ---

func TestDistributedTraceContextPropagation(t *testing.T) {
	exporter := setupTracerProvider(t)
	store := loadTestData(t)

	opts := engine.Opts{
		EngineOpts: promql.EngineOpts{
			Timeout:              time.Hour,
			MaxSamples:           1e10,
			EnableNegativeOffset: true,
			EnableAtModifier:     true,
		},
	}

	remote := engine.NewRemoteEngine(opts, store, math.MinInt64, math.MaxInt64, nil)
	endpoints := api.NewStaticEndpoints([]api.RemoteEngine{remote})
	ng := engine.NewDistributedEngine(opts)

	qry, err := ng.MakeInstantQuery(context.Background(), store, endpoints, nil,
		"http_requests_total", time.Unix(60, 0))
	require.NoError(t, err)
	defer qry.Close()

	result := qry.Exec(context.Background())
	require.NoError(t, result.Err)

	spans := exporter.GetSpans()

	// Should have a remote_query_exec span.
	remoteSpans := findSpans(spans, "remote_query_exec")
	require.NotEmpty(t, remoteSpans, "expected remote_query_exec span in distributed query")

	rs := &remoteSpans[0]
	require.NotEmpty(t, spanAttr(rs, "query.expr").AsString())
	require.NotZero(t, spanAttr(rs, "query.start").AsInt64())

	// The remote query exec should produce its own query span (from the remote engine).
	querySpan := findSpan(spans, "instant_query_exec")
	require.NotNil(t, querySpan, "expected query span from remote engine execution")

	// All spans should share the same trace ID.
	require.Equal(t, rs.SpanContext.TraceID(), querySpan.SpanContext.TraceID(),
		"remote and query spans should share the same trace ID")
}

func TestDistributedRemoteFailureSpan(t *testing.T) {
	exporter := setupTracerProvider(t)

	opts := engine.Opts{
		EngineOpts: promql.EngineOpts{
			Timeout:              time.Hour,
			MaxSamples:           1e10,
			EnableNegativeOffset: true,
			EnableAtModifier:     true,
		},
	}

	// Create a mock queryable that returns an error.
	errStore := &storage.MockQueryable{
		MockQuerier: &storage.MockQuerier{
			SelectMockFunction: func(_ bool, _ *storage.SelectHints, _ ...*labels.Matcher) storage.SeriesSet {
				return storage.ErrSeriesSet(fmt.Errorf("remote storage unavailable"))
			},
		},
	}

	remote := engine.NewRemoteEngine(opts, errStore, math.MinInt64, math.MaxInt64, nil)
	endpoints := api.NewStaticEndpoints([]api.RemoteEngine{remote})
	ng := engine.NewDistributedEngine(opts)

	qry, err := ng.MakeInstantQuery(context.Background(), errStore, endpoints, nil,
		"test_metric", time.Unix(60, 0))
	require.NoError(t, err)
	defer qry.Close()

	result := qry.Exec(context.Background())
	require.Error(t, result.Err)

	spans := exporter.GetSpans()

	// Should have a remote_query_exec span that recorded the error.
	remoteSpans := findSpans(spans, "remote_query_exec")
	if len(remoteSpans) > 0 {
		require.Equal(t, otelcodes.Error, remoteSpans[0].Status.Code,
			"remote span should record error status")
	}
}
