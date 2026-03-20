# Feature Specification: OpenTelemetry Tracing Support

**Feature Branch**: `001-otel-tracing`
**Created**: 2026-03-20
**Status**: Draft
**Input**: User description: "Add OpenTelemetry tracing support to the thanos promql engine"

## Clarifications

### Session 2026-03-20

- Q: Should span/attribute names follow OTel semantic conventions, a custom namespace, or match Thanos conventions? → A: Match Thanos project attribute names where they exist (e.g. `query.expr`, `query.start`, `query.end`, `query.range_seconds`, `query.interval_seconds`); use `promql.*` namespace for engine-specific attributes with no Thanos equivalent.
- Q: Should the engine manage the TracerProvider or assume the embedding program provides one? → A: The embedding program initializes the TracerProvider. The engine MUST only emit spans (with attributes and events); it MUST NOT accept, create, or manage a TracerProvider.
- Q: Should operator-level spans have their own toggle independent of EnableAnalysis? → A: Keep coupled to `EnableAnalysis`. When analysis is disabled, emit a span event on the query span noting that operator detail is available by enabling analysis. If practical, support OTel trace baggage to enable analysis for a single query via the trace context.
- Q: Should `query.expr` be redacted or truncated in span attributes? → A: Include verbatim. PII is rarely a concern for PromQL, redaction is not practical, and traces should be treated as internal debug information. Matches existing Thanos convention.
- Implementation conventions: Use the OpenTelemetry API exclusively (ignore Thanos's deprecated OpenTracing stack). Adopt Prometheus tracing patterns: `otel.Tracer("")` for tracer creation, `span.AddEvent()` for in-span events, `span.RecordError()` + `span.SetStatus(codes.Error, ...)` for errors, `trace.SpanFromContext(ctx)` to access existing spans. Consider the Prometheus `SpanTimer` pattern for combining span timing with metric observers.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Query-Level Trace Visibility (Priority: P1)

As an operator running the Thanos PromQL engine in production, I want
each query execution to produce an OpenTelemetry span so that I can
see query latency, the PromQL expression, query type (instant/range),
and error status in my distributed tracing backend (Jaeger, Tempo,
etc.).

**Why this priority**: Without top-level query spans, there is no
entry point for any deeper tracing. This is the foundation all other
tracing builds on.

**Independent Test**: Can be verified by registering an in-memory span
exporter on the global TracerProvider, executing a query, and asserting
that the exporter captured a span with the expected name and attributes.

**Acceptance Scenarios**:

1. **Given** a global TracerProvider is registered,
   **When** an instant query is executed,
   **Then** a span named `instant_query_exec` is created with attributes
   `query.expr` (the PromQL expression) and the query timestamp.
2. **Given** a global TracerProvider is registered,
   **When** a range query is executed,
   **Then** a span named `range_query_exec` is created with attributes
   `query.expr`, `query.start`, `query.end`, and
   `query.interval_seconds`.
3. **Given** a global TracerProvider is registered,
   **When** a query returns an error,
   **Then** the span records the error and sets span status to Error.
4. **Given** a no-op TracerProvider (the default),
   **When** a query is executed,
   **Then** no tracing overhead is introduced beyond the no-op span
   check and query behaviour is identical.

---

### User Story 2 - Operator-Level Spans (Priority: P2)

As a performance engineer, I want each operator in the execution tree
to produce a child span so that I can see how time is distributed
across operators (e.g. vector selector, aggregation, binary operation)
within a single query trace.

**Why this priority**: Operator-level detail is what makes tracing
actionable for debugging slow queries, but it depends on the
query-level span from US1.

**Independent Test**: Can be verified by executing a multi-operator
query (e.g. `sum(rate(m[5m]))`) with tracing enabled and asserting
that child spans exist for each operator in the tree.

**Acceptance Scenarios**:

1. **Given** a global TracerProvider is registered,
   **When** a query with multiple operators is executed,
   **Then** each operator (selector, function, aggregation) produces a
   child span nested under the query span.
2. **Given** a global TracerProvider is registered and `EnableAnalysis`
   is ON,
   **When** operator spans are collected,
   **Then** each span includes the operator name, the number of series
   processed, and sample counts (float, histogram, total) as
   attributes.
3. **Given** a global TracerProvider is registered and `EnableAnalysis`
   is OFF,
   **When** a query is executed,
   **Then** the query span contains a span event noting that
   operator-level detail is available by enabling analysis.
4. **Given** a global TracerProvider is registered, `EnableAnalysis` is
   OFF, and the trace context carries baggage requesting analysis,
   **When** a query is executed,
   **Then** operator-level child spans are emitted for that query as
   if `EnableAnalysis` were ON.

---

### User Story 2b - Storage and Series Selection Spans (Priority: P2)

As a performance engineer, I want to see trace spans when the engine
calls out to a storage implementation or requests series from an
external source, so that I can distinguish storage latency from
compute latency and understand the volume of data fetched.

**Why this priority**: Storage calls are often the dominant cost in
query execution. Without visibility into selector/storage spans, it is
impossible to attribute latency correctly.

**Independent Test**: Can be verified by executing a query with a
vector or matrix selector and asserting that the resulting trace
includes a storage span with request criteria attributes and response
volume attributes.

**Acceptance Scenarios**:

1. **Given** a global TracerProvider is registered,
   **When** a query triggers a vector or matrix selector that fetches
   series from storage,
   **Then** a span is emitted covering the storage call, with
   attributes identifying the request criteria (label matchers, time
   range hints).
2. **Given** a global TracerProvider is registered,
   **When** the storage response has been processed,
   **Then** the storage span includes attributes for the number of
   result series, the number of float samples, the number of histogram
   samples, and the total sample count (float + histogram).
3. **Given** a global TracerProvider is registered,
   **When** a storage call returns an error,
   **Then** the storage span records the error and sets span status to
   Error.

---

### User Story 3 - Context Propagation for Distributed Queries (Priority: P3)

As a Thanos user running the distributed execution mode, I want trace
context to propagate from the coordinating engine to remote engines so
that I can see a single end-to-end trace spanning the local
aggregation and all remote query legs.

**Why this priority**: Distributed tracing is the highest-value use
case but only applies to the distributed execution mode. It naturally
builds on US1.

**Independent Test**: Can be verified by configuring a distributed
engine with two in-process remote engines, executing a query, and
asserting that remote query spans share the same trace ID as the
coordinating query span.

**Acceptance Scenarios**:

1. **Given** a distributed engine with a global TracerProvider
   registered,
   **When** a query is executed that fans out to remote engines,
   **Then** the remote query spans are children of the local
   aggregation span and share the same trace ID.
2. **Given** a distributed engine where one remote engine fails,
   **When** the coordinating query completes,
   **Then** the failed remote span records the error while the parent
   span reflects partial failure.

---

### Edge Cases

- What happens when the context passed to a query already carries a
  parent span? The query span MUST become a child of that existing
  span, preserving the caller's trace context.
- What happens when the TracerProvider is replaced or shut down
  mid-flight? The engine MUST NOT panic; in-flight queries MUST
  complete with whatever spans have already been started.
- What happens with very deep operator trees (e.g. deeply nested
  subqueries)? The engine MUST NOT create an unbounded number of
  spans; a reasonable default depth limit MUST be applied or
  configurable.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The engine MUST obtain a tracer from the global
  OpenTelemetry TracerProvider (via the OTel API); it MUST NOT accept,
  create, or manage a TracerProvider itself.
- **FR-002**: The engine MUST create a span for each query execution
  (instant and range).
- **FR-003**: Query spans MUST use Thanos-consistent attribute names:
  `query.expr`, `query.start`, `query.end`, `query.interval_seconds`,
  `query.range_seconds`. Engine-specific attributes MUST use a
  `promql.*` namespace.
- **FR-004**: When `EnableAnalysis` is active, operator-level child
  spans MUST be created for each operator in the execution tree. All
  spans (query, operator, and storage) MUST include attributes for
  the number of series and samples (float, histogram, and total)
  processed by that span's scope.
- **FR-005**: When `EnableAnalysis` is NOT active, the query span MUST
  include a span event indicating that operator-level detail is
  available by enabling analysis.
- **FR-006**: If practical, the engine SHOULD support an OTel trace
  baggage key that enables analysis for a single query via the trace
  context, overriding the engine-level `EnableAnalysis` setting.
- **FR-007**: The engine MUST emit a span when it calls out to a
  storage implementation or requests series from an external source
  (e.g. via vector/matrix selectors). The span MUST include attributes
  identifying the request criteria (label matchers, time range hints).
- **FR-008**: When a storage/selector response is processed, the span
  MUST be updated with attributes for the number of result series,
  float samples, histogram samples, and total samples (float +
  histogram).
- **FR-009**: The engine MUST propagate trace context through the
  existing `context.Context` so that callers' parent spans are
  respected.
- **FR-010**: The engine MUST inject trace context into all outbound
  requests (gRPC, HTTP, or any other transport) that it initiates,
  using the OTel propagation API. This ensures that downstream
  services (including remote engines) receive the active trace context
  and can participate in the same distributed trace.
- **FR-011**: When a no-op TracerProvider is active (the default), the
  engine MUST behave identically to today with negligible performance
  impact.
- **FR-012**: In distributed execution mode, trace context MUST be
  injected into requests sent to remote engines so that remote query
  spans are children of the coordinating engine's span.
- **FR-013**: Query spans MUST record errors using both
  `span.RecordError()` and `span.SetStatus(codes.Error, ...)` when a
  query fails, following the Prometheus error recording convention.

### Key Entities

- **Query Span**: Top-level span representing a single query
  execution, carrying expression and timing attributes.
- **Operator Span**: Child span representing a single operator's
  execution within the query plan tree.
- **Storage Span**: Span covering a storage call (series selection),
  carrying request criteria (matchers, time range) and response volume
  attributes (series count, float/histogram/total sample counts).

## Assumptions

- The embedding program (e.g. Thanos) is responsible for initializing
  and configuring the global OpenTelemetry TracerProvider. The engine
  only emits spans.
- The engine MUST use the OpenTelemetry API exclusively; the Thanos
  project's legacy OpenTracing instrumentation MUST NOT be adopted.
  Prometheus's OTel tracing patterns are the primary implementation
  reference.
- Operator-level spans are gated behind `EnableAnalysis` to avoid
  per-operator overhead when only query-level tracing is desired.
- PromQL expressions are included verbatim in span attributes;
  traces are internal debug information and do not require redaction.
- The engine is responsible for injecting trace context into all
  outbound requests it makes (including to remote engines). Remote
  engine implementations are responsible for extracting the injected
  trace context on the receiving side.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A query executed with tracing enabled produces a
  complete trace visible in a tracing backend within the same request
  lifecycle.
- **SC-002**: A query executed with a no-op TracerProvider shows no
  measurable regression in benchmark alloc/op or time/op compared to
  the baseline (within noise margin).
- **SC-003**: A distributed query produces a single trace spanning
  the coordinator and all remote engines.
- **SC-004**: 100% of query error conditions are reflected in span
  status and error attributes.

## Future Enhancements

- **Peak memory estimation**: A future enhancement SHOULD add
  estimated peak memory use as a span attribute for expensive
  operations such as label matching, sort operations, and
  aggregations. This is out of scope for the initial implementation
  but is a desired follow-on capability.
