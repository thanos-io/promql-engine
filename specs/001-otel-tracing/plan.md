# Implementation Plan: OpenTelemetry Tracing Support

**Branch**: `001-otel-tracing` | **Date**: 2026-03-20 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-otel-tracing/spec.md`

## Summary

Add OpenTelemetry tracing to the Thanos PromQL engine so that query execution
produces structured spans visible in any OTel-compatible tracing backend. The
engine emits spans at three levels — query, operator, and storage — using the
global TracerProvider (managed by the embedding program). Operator-level spans
are gated behind `EnableAnalysis` with a per-query baggage override. The
implementation follows Prometheus's OTel tracing patterns (`otel.Tracer("")`,
`SpanTimer`, `RecordError` + `SetStatus`).

## Technical Context

**Language/Version**: Go 1.24 (toolchain 1.24.4)
**Primary Dependencies**: `go.opentelemetry.io/otel v1.38.0` (trace, attribute, codes, baggage), `github.com/prometheus/prometheus v0.308.0`
**Storage**: N/A (engine is a library; storage is provided by the embedding program)
**Testing**: `github.com/stretchr/testify`, `go.opentelemetry.io/otel/sdk/trace/tracetest` (in-memory span exporter), existing `go test ./...` and benchmark suite
**Target Platform**: Any (library, cross-platform)
**Project Type**: Library
**Performance Goals**: No measurable regression in `alloc/op` or `ns/op` with no-op TracerProvider (Constitution Principle III)
**Constraints**: Must not accept/create/manage TracerProvider; must use OTel API exclusively; must not adopt legacy OpenTracing
**Scale/Scope**: ~8 files modified, ~3 new files, ~500 lines of production code + ~400 lines of test code

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

| Principle | Pre-Research | Post-Design | Notes |
|-----------|-------------|-------------|-------|
| I. Prometheus Compatibility | PASS | PASS | Tracing is additive; no PromQL semantics change |
| II. Test-Driven Correctness | PASS | PASS | In-memory exporter enables span assertion tests |
| III. Minimise Memory Use | PASS | PASS | No-op tracer = zero allocs; operator spans gated; depth limiter prevents explosion |

## Project Structure

### Documentation (this feature)

```text
specs/001-otel-tracing/
├── plan.md              # This file
├── research.md          # Phase 0 output — OTel patterns, baggage API, integration analysis
├── data-model.md        # Phase 1 output — span definitions, attributes, relationships
├── quickstart.md        # Phase 1 output — usage examples for embedding programs
├── contracts/
│   └── tracing-api.md   # Phase 1 output — public API contract (no new API surface)
└── tasks.md             # Phase 2 output (not created by /speckit.plan)
```

### Source Code (repository root)

```text
engine/
├── engine.go                          # MODIFY: Add query-level spans in Exec()
└── engine_test.go                     # MODIFY: Add tracing integration tests

execution/
├── execution.go                       # MODIFY: Pass tracing context through operator creation
└── telemetry/
    ├── telemetry.go                   # MODIFY: Extend OperatorTelemetry with span lifecycle
    └── telemetry_test.go              # MODIFY: Add tests for span-enabled telemetry

storage/prometheus/
├── vector_selector.go                 # MODIFY: Add storage spans around GetSeries()
├── matrix_selector.go                 # MODIFY: Add storage spans around GetSeries()
└── selector_test.go                   # NEW: Storage span assertion tests

execution/remote/
└── operator.go                        # MODIFY: Add remote query spans, ensure ctx propagation

query/
└── options.go                         # MODIFY: Add baggage-based analysis override helper
```

**Structure Decision**: This is a library with existing well-defined packages.
All changes are modifications to existing files. No new packages or directories
needed in production code. Test files are co-located with production code per
Go convention.

## Key Design Decisions

### D1: Extend existing telemetry vs. parallel tracing layer

**Chosen**: Extend `execution/telemetry` package.
**Why**: The `telemetry.Operator` wrapper already instruments every operator's
`Series()` and `Next()`. Adding span lifecycle to this wrapper avoids
duplicating call sites across ~20 operator files. See research R3.

### D2: Storage spans always emitted vs. gated

**Chosen**: Always emitted.
**Why**: Storage calls represent I/O boundaries — the dominant cost in query
execution. Visibility into storage latency is critical regardless of whether
operator-level analysis is enabled. The overhead of a single span per selector
is negligible compared to the I/O cost.

### D3: Operator span depth limiting

**Chosen**: Configurable depth limit, default 32.
**Why**: Spec edge case requires bounded span count. Depth counter in operator
construction prevents unbounded spans from deeply nested subqueries while
preserving timing/sample tracking. See research R4.

### D4: Trace context propagation to remote engines

**Chosen**: Rely on `context.Context` propagation — do not inject headers directly.
**Why**: The engine doesn't own the transport. OTel-instrumented gRPC/HTTP
clients in the embedding program automatically extract trace context from `ctx`.
The engine's responsibility is to ensure `ctx` carries the active span. See
research R6.

### D5: Baggage key for per-query analysis

**Chosen**: `promql.enable_analysis` baggage key, any non-empty value enables.
**Why**: Uses the `promql.*` namespace per spec. Zero-cost when absent (O(1)
map miss). Standard OTel mechanism for request-scoped feature toggles. See
research R2.

## Complexity Tracking

No constitution violations to justify. All design decisions reduce complexity:
- No new public API surface
- No new packages
- No new configuration options beyond existing `EnableAnalysis`
- No new dependencies (all OTel packages already transitive)
