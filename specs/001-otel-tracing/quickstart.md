# Quickstart: OpenTelemetry Tracing Support

**Feature Branch**: `001-otel-tracing`
**Date**: 2026-03-20

## Prerequisites

- Go 1.24+
- An OpenTelemetry-compatible tracing backend (Jaeger, Tempo, OTLP collector)
  or an in-memory exporter for testing

## Enabling Tracing

The engine emits OTel spans automatically. The embedding program controls
whether spans are collected by registering a TracerProvider.

### Production Setup (embedding program)

```go
import (
    "go.opentelemetry.io/otel"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

// Initialize once at program startup
exp, _ := otlptracehttp.New(ctx)
tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp))
otel.SetTracerProvider(tp)
defer tp.Shutdown(ctx)

// Create engine — no tracing configuration needed
eng := engine.New(engine.Opts{
    EngineOpts: promql.EngineOpts{...},
})

// Execute queries — spans are emitted automatically
qry, _ := eng.NewRangeQuery(ctx, queryable, nil, "sum(rate(m[5m]))", start, end, step)
result := qry.Exec(ctx)
```

### Testing Setup (in-memory exporter)

```go
import (
    "go.opentelemetry.io/otel"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    "go.opentelemetry.io/otel/sdk/trace/tracetest"
)

exporter := tracetest.NewInMemoryExporter()
tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
otel.SetTracerProvider(tp)

// Execute query...
qry, _ := eng.NewInstantQuery(ctx, queryable, nil, "up", ts)
result := qry.Exec(ctx)

// Assert spans
spans := exporter.GetSpans()
// spans[0].Name == "instant_query_exec"
// spans[0].Attributes includes "query.expr" = "up"
```

## Operator-Level Spans

Operator spans are gated behind `EnableAnalysis`:

```go
eng := engine.New(engine.Opts{
    EngineOpts:     promql.EngineOpts{...},
    EnableAnalysis: true, // enables operator-level child spans
})
```

When disabled (default), the query span includes a `promql.analysis_available`
event indicating that detail is available.

### Per-Query Override via Baggage

Enable analysis for a single query without changing engine configuration:

```go
import "go.opentelemetry.io/otel/baggage"

member, _ := baggage.NewMember("promql.enable_analysis", "true")
bag, _ := baggage.New(member)
ctx = baggage.ContextWithBaggage(ctx, bag)

// This query will emit operator-level spans even if EnableAnalysis is OFF
qry, _ := eng.NewInstantQuery(ctx, queryable, nil, expr, ts)
result := qry.Exec(ctx)
```

## What to Expect in Your Tracing Backend

### Basic Query Trace

```
instant_query_exec (query.expr="up", query.start=1711000000000)
└── storage:select (storage.matchers="{__name__=\"up\"}", storage.series.count=5)
```

### Multi-Operator Query with Analysis Enabled

```
range_query_exec (query.expr="sum(rate(m[5m]))", query.start=..., query.end=..., query.interval_seconds=15)
├── aggregate:sum (promql.series.count=10, promql.samples.total=7200)
│   └── functionCall:rate (promql.series.count=10, promql.samples.total=7200)
│       └── storage:select (storage.matchers="{__name__=\"m\"}", storage.series.count=10)
```

### Distributed Query Trace

```
range_query_exec (coordinating engine)
├── remote_query_exec (query.expr="sum(rate(m[5m]))", engine=remote-1)
│   └── range_query_exec (remote engine's own trace)
│       └── ...
├── remote_query_exec (query.expr="sum(rate(m[5m]))", engine=remote-2)
│   └── ...
└── aggregate:sum (final aggregation)
```

## Verifying No-Op Overhead

When no TracerProvider is registered (the default), the engine uses OTel's
built-in no-op tracer. Verify with benchmarks:

```bash
# Baseline (before tracing code)
go test -bench=BenchmarkRangeQuery -benchmem -count=10 ./engine/ > old.txt

# After tracing code
go test -bench=BenchmarkRangeQuery -benchmem -count=10 ./engine/ > new.txt

# Compare
benchstat old.txt new.txt
```

No statistically significant regression in `alloc/op` or `ns/op` is expected.
