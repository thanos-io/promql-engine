# Tasks: OpenTelemetry Tracing Support

**Input**: Design documents from `/specs/001-otel-tracing/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/

**Tests**: Included per Constitution Principle II (Test-Driven Correctness) and spec acceptance scenarios.

**Organization**: Tasks grouped by user story for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US2b, US3)
- Exact file paths included in descriptions

---

## Phase 1: Setup

**Purpose**: Promote OTel dependencies and capture performance baseline

- [x] T001 Promote OTel packages from indirect to direct in go.mod by running `go get go.opentelemetry.io/otel@v1.38.0 go.opentelemetry.io/otel/trace@v1.38.0 go.opentelemetry.io/otel/codes go.opentelemetry.io/otel/attribute go.opentelemetry.io/otel/baggage` and add OTel SDK test dependency `go.opentelemetry.io/otel/sdk`
- [x] T002 Capture baseline benchmark results by running `go test -bench=BenchmarkRangeQuery -benchmem -count=10 ./engine/` and saving output to `specs/001-otel-tracing/bench-baseline.txt`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core tracing infrastructure that all user stories depend on

**CRITICAL**: No user story work can begin until this phase is complete

- [x] T003 Add `ShouldEnableAnalysis(ctx context.Context) bool` helper to query/options.go that checks baggage key `promql.enable_analysis` via `baggage.FromContext(ctx)`, returning true if global `EnableAnalysis` is set OR baggage value is non-empty (per research R2)
- [x] T004 [P] Add span-aware telemetry to execution/telemetry/telemetry.go: add a `trace.Span` field to `TrackedTelemetry` (no interface change needed); in the `telemetry.Operator.Series(ctx)` wrapper, start a child span via `otel.Tracer("").Start(ctx, operatorName)` when analysis is enabled, store it on the inner telemetry, and pass the span-enriched ctx to the inner operator; in `Next(ctx)`, update span attributes and end the span on completion/error — span lifecycle is an internal detail driven by the context already flowing through `Series`/`Next`
**Checkpoint**: Foundation ready — user story implementation can begin

---

## Phase 3: User Story 1 — Query-Level Trace Visibility (Priority: P1) MVP

**Goal**: Every query execution produces a top-level OTel span with expression, timing attributes, and error status.

**Independent Test**: Register in-memory span exporter, execute query, assert span name and attributes.

### Tests for User Story 1

> **Write these tests FIRST, ensure they FAIL before implementation**

- [x] T006 [US1] Add test helper `setupTracerProvider(t)` in engine/engine_test.go that registers an in-memory span exporter (`tracetest.NewInMemoryExporter`) on the global TracerProvider and returns the exporter for span assertions; clean up in `t.Cleanup`
- [x] T007 [US1] Add test `TestInstantQuerySpan` in engine/engine_test.go: execute an instant query with tracing enabled, assert exporter captured a span named `instant_query_exec` with attributes `query.expr` and `query.start` (acceptance scenario 1)
- [x] T008 [US1] Add test `TestRangeQuerySpan` in engine/engine_test.go: execute a range query, assert span named `range_query_exec` with attributes `query.expr`, `query.start`, `query.end`, `query.interval_seconds`, `query.range_seconds` (acceptance scenario 2)
- [x] T009 [US1] Add test `TestQuerySpanError` in engine/engine_test.go: execute a query that returns an error, assert span has `RecordError` event and status `codes.Error` (acceptance scenario 3)
- [x] T010 [US1] Add test `TestQuerySpanParentContext` in engine/engine_test.go: create a parent span in the test, pass its context to query execution, assert query span is a child of the parent span (edge case: parent span respected)

### Implementation for User Story 1

- [x] T011 [US1] Add query-level span creation in engine/engine.go `compatibilityQuery.Exec()`: at the start of `Exec()`, call `otel.Tracer("").Start(ctx, spanName)` where spanName is `instant_query_exec` or `range_query_exec` based on query type; defer `span.End()`; set attributes `query.expr` (string), `query.start` (int64 Unix ms) on all queries; set `query.end`, `query.interval_seconds`, `query.range_seconds` on range queries only
- [x] T012 [US1] Add error recording to query span in engine/engine.go: when `Exec()` encounters an error, call `span.RecordError(err)` and `span.SetStatus(codes.Error, err.Error())` following Prometheus convention (FR-013)
- [x] T013 [US1] Add `promql.analysis_available` span event in engine/engine.go: when `EnableAnalysis` is OFF (and no baggage override per `ShouldEnableAnalysis`), call `span.AddEvent("promql.analysis_available", trace.WithAttributes(attribute.String("message", "Operator-level detail available by enabling analysis")))` on the query span (FR-005)

**Checkpoint**: US1 complete — instant/range queries produce spans with correct attributes and error handling

---

## Phase 4: User Story 2 — Operator-Level Spans (Priority: P2)

**Goal**: Each operator in the execution tree produces a child span (when analysis enabled) with operator name, series count, and sample counts.

**Independent Test**: Execute `sum(rate(m[5m]))` with analysis ON, assert child spans for each operator.

**Dependencies**: Requires US1 (query span exists as parent) and Phase 2 (telemetry extension)

### Tests for User Story 2

- [x] T014 [US2] Add test `TestOperatorSpansWithAnalysis` in execution/telemetry/telemetry_test.go: create a multi-operator query (`sum(rate(m[5m]))`) with `EnableAnalysis: true`, assert child spans exist for aggregation, function, and selector operators with correct parent-child nesting (acceptance scenario 1)
- [x] T015 [US2] Add test `TestOperatorSpanAttributes` in execution/telemetry/telemetry_test.go: with analysis ON, assert each operator span includes `promql.operator.type`, `promql.series.count`, `promql.samples.float`, `promql.samples.histogram`, `promql.samples.total` attributes (acceptance scenario 2)
- [x] T016 [US2] Add test `TestAnalysisOffNoOperatorSpans` in engine/engine_test.go: with `EnableAnalysis: false`, execute a query, assert no operator child spans exist but query span contains `promql.analysis_available` event (acceptance scenario 3)
- [x] T017 [US2] Add test `TestBaggageOverrideEnablesAnalysis` in engine/engine_test.go: with `EnableAnalysis: false` but `promql.enable_analysis` baggage set on context, assert operator child spans ARE emitted (acceptance scenario 4)
- [x] T018 [US2] Add test `TestBaggageOverrideDisabledByDefault` in engine/engine_test.go: with `EnableAnalysis: false` and NO baggage, confirm no operator spans and `promql.analysis_available` event present (reinforces T016 from query-span perspective)

### Implementation for User Story 2

- [x] T019 [US2] Wire analysis gating into operator span creation in execution/telemetry/telemetry.go: the span start logic added in T004 must check `ShouldEnableAnalysis(ctx)` — only create operator spans when analysis is enabled (global or baggage override); set `promql.operator.type` attribute on span creation
- [x] T020 [US2] Add sample count attributes to operator spans in execution/telemetry/telemetry.go: before `span.End()` in the `Next()` path, set `promql.series.count`, `promql.samples.float`, `promql.samples.histogram`, `promql.samples.total` from the tracked telemetry counters; record errors with `RecordError()` + `SetStatus(codes.Error, ...)`
- [x] T021 [US2] Wire baggage-based analysis check in engine/engine.go: before operator tree construction, call `ShouldEnableAnalysis(ctx)` to determine effective analysis state; pass effective state through `query.Options` so that telemetry wrapper respects per-query baggage override

**Checkpoint**: US2 complete — operator tree produces hierarchical child spans with performance attributes when analysis is enabled; baggage override works per-query

---

## Phase 5: User Story 2b — Storage and Series Selection Spans (Priority: P2)

**Goal**: Storage calls produce spans with request criteria and response volume attributes.

**Independent Test**: Execute query with vector selector, assert storage span with matchers and series/sample counts.

**Dependencies**: Requires US1 (query span as parent). Independent of US2 (operator spans).

### Tests for User Story 2b

- [x] T022 [P] [US2b] Add test `TestVectorSelectorStorageSpan` in storage/prometheus/selector_test.go: execute a query with a vector selector, assert a `storage:select` span with `storage.matchers`, `storage.mint`, `storage.maxt` attributes and response attributes `storage.series.count`, `storage.samples.float`, `storage.samples.histogram`, `storage.samples.total` (acceptance scenarios 1 & 2)
- [x] T023 [P] [US2b] Add test `TestMatrixSelectorStorageSpan` in storage/prometheus/selector_test.go: execute a query with a matrix selector (e.g. `rate(m[5m])`), assert same `storage:select` span with correct attributes (acceptance scenario 1)
- [x] T024 [US2b] Add test `TestStorageSpanError` in storage/prometheus/selector_test.go: trigger a storage error, assert storage span records error and sets status to `codes.Error` (acceptance scenario 3)

### Implementation for User Story 2b

- [x] T025 [P] [US2b] Add storage span to vector selector in storage/prometheus/vector_selector.go: wrap the `GetSeries()` call in `Series()` method with `otel.Tracer("").Start(ctx, "storage:select")`; set `storage.matchers` (string representation of label matchers), `storage.mint`, `storage.maxt` attributes before the call; after `GetSeries()` returns, set `storage.series.count`; on error, call `RecordError()` + `SetStatus(codes.Error, ...)`; defer `span.End()`
- [x] T026 [P] [US2b] Add storage span to matrix selector in storage/prometheus/matrix_selector.go: same pattern as T025 — wrap `GetSeries()` in `Series()` with a `storage:select` span, set request criteria and response volume attributes
- [x] T027 [US2b] Add sample count tracking to storage spans in storage/prometheus/vector_selector.go and storage/prometheus/matrix_selector.go: after all `Next()` calls complete, update the storage span with `storage.samples.float`, `storage.samples.histogram`, `storage.samples.total` attributes (FR-008)

**Checkpoint**: US2b complete — every storage call produces a span with request and response attributes

---

## Phase 6: User Story 3 — Context Propagation for Distributed Queries (Priority: P3)

**Goal**: Distributed queries produce a single trace spanning coordinator and all remote engines.

**Independent Test**: Configure distributed engine with in-process remote engines, assert remote spans share trace ID with coordinating span.

**Dependencies**: Requires US1 (query span context propagation)

### Tests for User Story 3

- [x] T028 [US3] Add test `TestDistributedTraceContextPropagation` in execution/remote/operator_test.go (or engine/engine_test.go): configure a distributed engine with two in-process remote engines, execute a query, assert remote query spans are children of the coordinating query span and share the same trace ID (acceptance scenario 1)
- [x] T029 [US3] Add test `TestDistributedRemoteFailureSpan` in execution/remote/operator_test.go: configure a distributed engine where one remote fails, assert failed remote span records error while parent span reflects partial failure (acceptance scenario 2)

### Implementation for User Story 3

- [x] T030 [US3] Add `remote_query_exec` span in execution/remote/operator.go: in the remote operator's execution path (before `query.Exec(ctx)`), start a span `otel.Tracer("").Start(ctx, "remote_query_exec")` with attributes `query.expr`, `query.start`, `query.end`; pass the span-enriched `ctx` to the remote query so trace context propagates; defer `span.End()`
- [x] T031 [US3] Add error recording to remote span in execution/remote/operator.go: when the remote query returns an error, call `span.RecordError(err)` + `span.SetStatus(codes.Error, err.Error())` on the remote span

**Checkpoint**: US3 complete — distributed queries produce end-to-end traces spanning coordinator and remote engines

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Verify performance, run full suite, validate documentation

- [x] T032 Run `go test ./...` to verify all existing tests still pass with tracing code added
- [x] T033 Run `go test -bench=BenchmarkRangeQuery -benchmem -count=10 ./engine/` and compare against baseline (`specs/001-otel-tracing/bench-baseline.txt`) using `benchstat`; verify no statistically significant regression in `alloc/op` or `ns/op` (SC-002, Constitution Principle III)
- [x] T034 Run `go vet ./...` and fix any issues introduced by tracing changes
- [x] T035 Verify quickstart.md code examples compile by reviewing against actual implementation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 (go.mod updated)
- **US1 (Phase 3)**: Depends on Phase 2 — MVP target
- **US2 (Phase 4)**: Depends on Phase 2 AND US1 (needs query span as parent)
- **US2b (Phase 5)**: Depends on Phase 2 AND US1; independent of US2
- **US3 (Phase 6)**: Depends on Phase 2 AND US1; independent of US2/US2b
- **Polish (Phase 7)**: Depends on all user stories being complete

### User Story Dependencies

```text
Phase 1: Setup
    │
Phase 2: Foundational (T003, T004)
    │
    ├── Phase 3: US1 - Query Spans (T006–T013) ← MVP
    │       │
    │       ├── Phase 4: US2 - Operator Spans (T014–T021)
    │       │
    │       ├── Phase 5: US2b - Storage Spans (T022–T027) [parallel with US2]
    │       │
    │       └── Phase 6: US3 - Distributed Propagation (T028–T031) [parallel with US2, US2b]
    │
Phase 7: Polish (T032–T035)
```

### Within Each User Story

- Tests written FIRST — verify they FAIL before implementation
- Implementation tasks in dependency order
- Story complete before moving to next priority

### Parallel Opportunities

- **Phase 2**: T003 and T004 can run in parallel (different files)
- **Phase 5**: T022 and T023 in parallel (same file but independent test cases); T025 and T026 in parallel (different files)
- **After US1**: US2, US2b, and US3 can proceed in parallel (independent stories)

---

## Parallel Example: User Story 2b

```text
# Storage span tests (parallel — independent test cases):
T022: "TestVectorSelectorStorageSpan in storage/prometheus/selector_test.go"
T023: "TestMatrixSelectorStorageSpan in storage/prometheus/selector_test.go"

# Storage span implementation (parallel — different files):
T025: "Add storage span to vector_selector.go"
T026: "Add storage span to matrix_selector.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001–T002)
2. Complete Phase 2: Foundational (T003–T005)
3. Complete Phase 3: User Story 1 (T006–T013)
4. **STOP and VALIDATE**: Run tests, verify spans in test output
5. Run benchmarks to confirm no regression

### Incremental Delivery

1. Setup + Foundational → Foundation ready
2. Add US1 → Test independently → **MVP!** (query-level tracing)
3. Add US2 + US2b (parallel) → Test independently → operator + storage visibility
4. Add US3 → Test independently → full distributed tracing
5. Polish → benchmarks, vet, docs

---

## Notes

- [P] tasks = different files, no dependencies on incomplete tasks
- [Story] label maps task to specific user story for traceability
- All OTel packages are already transitive deps — no new external dependencies
- Constitution Principle III: benchmark validation in T033 is mandatory
- Edge case (depth limit) deferred — can be added later as an optional enhancement
- Edge case (parent span) covered in T010
- Edge case (TracerProvider shutdown) handled by OTel SDK — no engine code needed
