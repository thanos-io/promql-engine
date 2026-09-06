# Contract: Tracing API Surface

**Feature Branch**: `001-otel-tracing`
**Date**: 2026-03-20

This document defines the public API contract that consumers of the
thanos-promql-engine library observe after OTel tracing is added.

## Principle: No New Public API

The engine is a library. Tracing is added as a transparent behaviour change —
the engine emits OTel spans using the global TracerProvider. There are no new
exported types, functions, or configuration options exposed to callers beyond
what already exists.

**Callers control tracing by**:
1. Registering a `TracerProvider` via `otel.SetTracerProvider(tp)` before
   creating queries (embedding program responsibility).
2. Setting `EnableAnalysis: true` in `engine.Opts` to get operator-level spans
   (existing API).
3. Optionally setting `promql.enable_analysis` baggage on the context for
   per-query override (`"true"` to enable, `"false"` to suppress).

## Span Contract

The engine guarantees the following span emission contract when a non-noop
TracerProvider is registered:

### Query Spans (always emitted)

| Condition | Span Name | Required Attributes |
|-----------|-----------|-------------------|
| Instant query executed | `instant_query_exec` | `query.expr`, `query.start` |
| Range query executed | `range_query_exec` | `query.expr`, `query.start`, `query.end`, `query.interval_seconds`, `query.range_seconds` |

### Operator Spans (conditional)

| Condition | Span Name | Required Attributes |
|-----------|-----------|-------------------|
| `EnableAnalysis` ON (or baggage `"true"`) | `{operator_string}` | `promql.operator.type`, `promql.series.count`, `promql.samples.float`, `promql.samples.histogram`, `promql.samples.total` |
| `EnableAnalysis` ON, baggage `"false"` | (no operator spans) | Baggage explicitly suppresses operator spans |
| `EnableAnalysis` OFF, no baggage | (no operator spans) | Query span receives `promql.analysis_available` event |

### Storage Spans (always emitted)

| Condition | Span Name | Required Attributes |
|-----------|-----------|-------------------|
| Vector/matrix selector fetches series | `storage:select` | `storage.matchers`, `storage.mint`, `storage.maxt` |
| After series loaded | (same span) | `storage.series.count`, `storage.samples.float`, `storage.samples.histogram`, `storage.samples.total` |

### Remote Query Spans (distributed mode)

| Condition | Span Name | Required Attributes |
|-----------|-----------|-------------------|
| Remote engine query dispatched | `remote_query_exec` | `query.expr`, `query.start`, `query.end` |

### Error Contract

| Condition | Behaviour |
|-----------|-----------|
| Query returns error | `span.RecordError(err)` + `span.SetStatus(codes.Error, err.Error())` |
| Storage call fails | Same error recording on storage span |
| Remote query fails | Same error recording on remote span |

## Context Propagation Contract

| Guarantee | Description |
|-----------|-------------|
| Parent span respected | If `ctx` carries a parent span, the query span becomes its child |
| Trace context flows to operators | All operator spans are children of the query span |
| Trace context flows to remote engines | The `ctx` passed to `RemoteEngine.NewRangeQuery()` carries the active query span |

## No-Op Contract

| Guarantee | Description |
|-----------|-------------|
| Default behaviour unchanged | With no TracerProvider (or noop), query behaviour is identical to pre-tracing |
| No measurable overhead | No allocation or timing regression beyond noise margin |

## Backward Compatibility

This change is fully backward compatible:
- No existing public API signatures change.
- No new required configuration.
- The default (no TracerProvider) produces identical behaviour.
- `EnableAnalysis` retains its existing semantics; tracing adds span emission
  as a side effect.
