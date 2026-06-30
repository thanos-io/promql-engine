# Research: OpenTelemetry Tracing Support

**Feature Branch**: `001-otel-tracing`
**Date**: 2026-03-20

## R1: Prometheus OTel Tracing Patterns

**Decision**: Follow Prometheus's OTel conventions exactly.

**Rationale**: The spec mandates using the OpenTelemetry API exclusively and
adopting Prometheus tracing patterns. Prometheus v0.308.0 (our dependency)
already has mature OTel instrumentation that serves as the reference
implementation.

**Key patterns extracted from Prometheus source**:

1. **Tracer creation**: `otel.Tracer("")` — empty-string instrumentation scope.
   The engine never creates or manages a TracerProvider.
2. **Span creation**: `ctx, span := otel.Tracer("").Start(ctx, operationName)` /
   `defer span.End()`.
3. **SpanTimer pattern** (`util/stats/query_stats.go`): Combines span lifecycle
   with timer + Prometheus observer in one struct. `NewSpanTimer` starts span
   and timer; `Finish()` ends both and records the observation. Our engine
   already has per-operator timing in `TrackedTelemetry`; we can adopt a
   similar unified approach.
4. **Error recording**: Both `span.RecordError(err)` and
   `span.SetStatus(codes.Error, err.Error())` — always paired.
5. **Accessing parent spans**: `trace.SpanFromContext(ctx)` to add attributes
   or events to an existing span without creating a new child.
6. **Span events**: `span.AddEvent("message", trace.WithAttributes(...))` for
   lightweight in-span annotations.
7. **Propagation for distributed tracing**:
   `otel.GetTextMapPropagator().Inject(ctx, carrier)` to inject W3C
   TraceContext into outbound HTTP/gRPC headers.

**Alternatives considered**: OpenTracing (deprecated, spec explicitly forbids),
custom tracing abstraction (unnecessary given OTel API stability).

## R2: OTel Baggage API for Per-Query Analysis Override

**Decision**: Use `go.opentelemetry.io/otel/baggage` to support per-query
`EnableAnalysis` override via trace context (FR-006).

**Rationale**: The baggage package is already a transitive dependency (pulled
in by `otel v1.38.0`). Reading baggage from context is O(1) — a single
`context.Value()` lookup plus a map access. Zero allocation when the key is
absent.

**Implementation pattern**:
```go
import "go.opentelemetry.io/otel/baggage"

func shouldEnableAnalysis(ctx context.Context, globalSetting bool) bool {
    if globalSetting {
        return true
    }
    b := baggage.FromContext(ctx)
    return b.Member("promql.enable_analysis").Value() != ""
}
```

**Baggage key**: `promql.enable_analysis` — uses the `promql.*` namespace per
spec conventions. Any non-empty value (e.g. `"true"`, `"1"`) enables analysis
for that query.

**Alternatives considered**: Custom context key (not propagated across service
boundaries), query parameter (not accessible at engine level). Baggage is the
standard OTel mechanism for request-scoped feature toggles.

## R3: Telemetry Layer Integration Strategy

**Decision**: Extend the existing `execution/telemetry` package to optionally
create OTel spans, rather than creating a parallel tracing layer.

**Rationale**: The `telemetry.Operator` wrapper already instruments every
operator's `Series()` and `Next()` calls with timing and sample counting. Adding
span lifecycle management to this wrapper avoids duplicating the instrumentation
call sites across ~20 operator files.

**Integration points identified**:

| Layer | Location | Span Type |
|-------|----------|-----------|
| Query | `engine/engine.go:519` `Exec()` | `instant_query_exec` / `range_query_exec` |
| Operator | `execution/telemetry/telemetry.go:195` `Series()` / `Next()` | `operator:{name}` |
| Storage | `storage/prometheus/vector_selector.go:186` `GetSeries()` | `storage:select` |
| Storage | `storage/prometheus/matrix_selector.go:231` `GetSeries()` | `storage:select` |
| Remote | `execution/remote/operator.go:114` `query.Exec()` | `remote_query_exec` |

**Span gating strategy**:
- **Query-level spans**: Always emitted (negligible overhead with no-op tracer).
- **Operator-level spans**: Gated behind `EnableAnalysis` (or baggage override).
  When analysis is off, a span event is added to the query span noting that
  detail is available.
- **Storage spans**: Always emitted (these represent I/O boundaries and are
  critical for latency attribution regardless of analysis mode).

**Alternatives considered**: Separate `execution/tracing/` package (rejected —
duplicates wrapping logic), tracing only in `engine.go` without operator detail
(rejected — operator spans are a core requirement).

## R4: Operator Span Depth Limiting

**Decision**: Implement a configurable depth limit (default 32) for operator
spans to prevent span explosion in deeply nested subqueries.

**Rationale**: The spec edge case requires that the engine "MUST NOT create an
unbounded number of spans" for deep operator trees. A depth counter passed
through operator construction prevents this.

**Implementation pattern**:
- Track current depth in `query.Options` or as a context value during operator
  tree construction.
- Each call to `newOperator()` in `execution/execution.go` increments depth.
- When depth exceeds the limit, operator telemetry falls back to no-op span
  creation (timing/sample tracking continues).
- The depth limit applies only to operator spans; query and storage spans are
  unaffected.

**Alternatives considered**: Post-hoc span sampling (loses determinism), no
limit with documentation warning (violates spec MUST requirement).

## R5: No-Op Performance Overhead

**Decision**: Rely on the OTel SDK's built-in no-op tracer for zero-overhead
when tracing is disabled. Validate with benchmark comparison.

**Rationale**: When no TracerProvider is registered (or a `noop.TracerProvider`
is set), all OTel API calls (`otel.Tracer("").Start()`, `span.SetAttributes()`,
`span.End()`) resolve to no-op implementations. The overhead is one interface
method dispatch per call — no allocations, no synchronization.

**Validation approach**:
- Run existing benchmarks (`BenchmarkSingleQuery`, `BenchmarkRangeQuery`) before
  and after the change.
- Compare `benchstat` output for `alloc/op` and `ns/op`.
- Constitution Principle III requires no unexplained regression.

**Alternatives considered**: Compile-time build tags to exclude tracing code
(rejected — unnecessary complexity given no-op overhead is negligible).

## R6: Distributed Trace Context Propagation

**Decision**: Inject trace context into the `context.Context` passed to remote
engine `NewRangeQuery()` calls. The engine does NOT inject HTTP/gRPC headers
directly — it relies on the transport layer (provided by the embedding program)
to propagate context from the `ctx`.

**Rationale**: The remote engine interface (`api.RemoteEngine.NewRangeQuery`)
accepts a `context.Context`. OTel-instrumented gRPC/HTTP clients (which the
embedding program configures) automatically extract trace context from `ctx`
and inject it into outbound headers. The engine's responsibility is to ensure
the `ctx` carries the correct active span.

**Key insight**: The `execution/remote/operator.go` creates a query via
`e.NewRangeQuery(ctx, ...)` — the `ctx` already flows through. As long as the
query span is started before this call, the remote engine will receive the
trace context automatically if the transport is instrumented.

**Alternatives considered**: Direct header injection in the engine (rejected —
the engine doesn't own the transport; this is the embedding program's
responsibility per the spec's assumption section).

## R7: Required OTel Go Imports

All required packages are already transitive dependencies via `otel v1.38.0`:

```go
import (
    "go.opentelemetry.io/otel"              // otel.Tracer("")
    "go.opentelemetry.io/otel/attribute"     // attribute.String(), attribute.Int64()
    "go.opentelemetry.io/otel/baggage"       // baggage.FromContext()
    "go.opentelemetry.io/otel/codes"         // codes.Error
    "go.opentelemetry.io/otel/trace"         // trace.SpanFromContext()
)
```

These will need to be promoted from indirect to direct dependencies in `go.mod`.
