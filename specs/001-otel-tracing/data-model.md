# Data Model: OpenTelemetry Tracing Support

**Feature Branch**: `001-otel-tracing`
**Date**: 2026-03-20

## Entities

### 1. Span Definitions

No new persistent data structures are introduced. All entities below are
transient OTel spans emitted during query execution and consumed by the
configured tracing backend.

#### QuerySpan

Top-level span representing a single query execution.

| Field | OTel Concept | Value |
|-------|-------------|-------|
| Name | Span name | `instant_query_exec` or `range_query_exec` |
| Parent | Parent span | Inherited from caller's `context.Context` (if any) |
| `query.expr` | Attribute (string) | Verbatim PromQL expression |
| `query.start` | Attribute (int64) | Query start timestamp (Unix ms) |
| `query.end` | Attribute (int64) | Query end timestamp (Unix ms); instant: same as start |
| `query.interval_seconds` | Attribute (float64) | Step interval in seconds (range only) |
| `query.range_seconds` | Attribute (float64) | Query range duration in seconds (range only) |
| `promql.analysis_available` | Event | Emitted when `EnableAnalysis` is OFF (see below) |
| Error | Status + RecordError | On query failure: `span.RecordError(err)` + `span.SetStatus(codes.Error, msg)` |

**Lifecycle**: Created at the start of `compatibilityQuery.Exec()`, ended
after result formatting completes.

#### OperatorSpan

Child span representing a single operator's execution within the query plan.
Only emitted when `EnableAnalysis` is ON (or overridden via baggage).

| Field | OTel Concept | Value |
|-------|-------------|-------|
| Name | Span name | Operator string representation (e.g. `aggregate:sum`, `vectorSelector`) |
| Parent | Parent span | QuerySpan or parent OperatorSpan |
| `promql.operator.type` | Attribute (string) | Operator class name |
| `promql.series.count` | Attribute (int64) | Number of series processed |
| `promql.samples.float` | Attribute (int64) | Float sample count |
| `promql.samples.histogram` | Attribute (int64) | Histogram sample count |
| `promql.samples.total` | Attribute (int64) | Total sample count (float + histogram) |
| Error | Status + RecordError | On operator error |

**Lifecycle**: Created when the operator's `Series()` is called (tree
materialisation), ended after the last `Next()` call completes or errors.

**Depth limit**: Operator spans stop being created beyond a configurable depth
(default 32). Timing and sample tracking continue regardless.

#### StorageSpan

Span covering a storage call (series selection via vector/matrix selectors).
Always emitted regardless of `EnableAnalysis`.

| Field | OTel Concept | Value |
|-------|-------------|-------|
| Name | Span name | `storage:select` |
| Parent | Parent span | OperatorSpan (if analysis ON) or QuerySpan |
| `storage.matchers` | Attribute (string) | Label matcher representation |
| `storage.mint` | Attribute (int64) | Minimum time hint (Unix ms) |
| `storage.maxt` | Attribute (int64) | Maximum time hint (Unix ms) |
| `storage.series.count` | Attribute (int64) | Number of result series |
| `storage.samples.float` | Attribute (int64) | Float sample count |
| `storage.samples.histogram` | Attribute (int64) | Histogram sample count |
| `storage.samples.total` | Attribute (int64) | Total sample count |
| Error | Status + RecordError | On storage error |

**Lifecycle**: Created before `GetSeries()` call, ended after series loading
completes.

#### RemoteQuerySpan

Span covering a remote engine query in distributed execution mode. Child of the
coordinating query span.

| Field | OTel Concept | Value |
|-------|-------------|-------|
| Name | Span name | `remote_query_exec` |
| Parent | Parent span | QuerySpan of the coordinating engine |
| `query.expr` | Attribute (string) | Remote query expression |
| `query.start` | Attribute (int64) | Remote query start (Unix ms) |
| `query.end` | Attribute (int64) | Remote query end (Unix ms) |
| Error | Status + RecordError | On remote query error |

**Lifecycle**: Created before `query.Exec(ctx)` on the remote operator, ended
after the remote result is received.

### 2. Span Event: Analysis Available

When `EnableAnalysis` is OFF and no baggage override is present, a single span
event is added to the QuerySpan:

| Field | Value |
|-------|-------|
| Event name | `promql.analysis_available` |
| `message` | `"Operator-level detail available by enabling analysis"` |

### 3. Baggage Key

| Key | Purpose | Values |
|-----|---------|--------|
| `promql.enable_analysis` | Per-query analysis override via W3C Baggage | `"true"` enables analysis; `"false"` explicitly suppresses; absent or empty means no override (defers to `EnableAnalysis` option) |

## Relationships

```text
QuerySpan (instant_query_exec / range_query_exec)
├── OperatorSpan (aggregate:sum)          [only when analysis ON]
│   ├── OperatorSpan (functionCall:rate)  [only when analysis ON]
│   │   └── StorageSpan (storage:select)  [always]
│   └── ...
├── StorageSpan (storage:select)          [always, if no operator spans]
└── RemoteQuerySpan (remote_query_exec)   [distributed mode only]
    └── [remote engine's own trace tree]
```

## Validation Rules

1. QuerySpan `query.expr` MUST NOT be empty.
2. QuerySpan `query.start` MUST be set for all query types.
3. QuerySpan `query.end` and `query.interval_seconds` MUST be set for range
   queries only.
4. OperatorSpan sample count attributes MUST only be set when
   `EnableAnalysis` is active (they depend on telemetry tracking).
5. StorageSpan `storage.matchers` MUST faithfully represent the label matchers
   used in the storage call.
6. All error spans MUST use both `RecordError()` and `SetStatus(codes.Error, ...)`.

## State Transitions

Spans follow the standard OTel lifecycle: Created → Active → Ended.
No custom state machines are introduced.

| Trigger | From | To |
|---------|------|----|
| `otel.Tracer("").Start(ctx, name)` | (none) | Active |
| `span.End()` | Active | Ended |
| `span.RecordError(err)` | Active | Active (error recorded) |
| `span.SetStatus(codes.Error, msg)` | Active | Active (status set) |

Spans are immutable after `End()` is called. The OTel SDK handles export.
