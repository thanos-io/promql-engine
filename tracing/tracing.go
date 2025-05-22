// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

// Package tracing adds OpenTracing support to the promql-engine library.
// It allows you to trace the execution of PromQL queries to understand performance bottlenecks and debug query issues.
// For now, we are using the OpenTracing API, but we plan to move to OpenTelemetry in the future (once Thanos also does the move).
// See https://github.com/thanos-io/thanos/issues/1972 for more details.
//
// # Setup
//
// To use tracing with promql-engine, you need to:
//
// 1. Initialize an OpenTracing compatible tracer (such as Jaeger, Zipkin, etc.)
// 2. Set the tracer using the SetTracer function
// 3. Use the engine normally - all query operations will automatically be traced
//
//	import (
//	    "github.com/opentracing/opentracing-go"
//	    "github.com/thanos-io/promql-engine/engine"
//	    "github.com/thanos-io/promql-engine/tracing"
//	    // Your tracer implementation, e.g., jaeger
//	)
//
//	func main() {
//	    // Initialize your tracer
//	    tracer := initTracer()
//
//	    // Set the tracer to be used by promql-engine
//	    tracing.SetTracer(tracer)
//
//	    // Initialize the engine as usual
//	    eng := engine.New(engine.Opts{
//	        // Your engine options here
//	    })
//
//	    // Use the engine as normal - all operations will now be traced
//	    // ...
//	}
//
// # Traced Operations
//
// The following operations are automatically traced:
//
// 1. Main query creation (MakeInstantQuery and MakeRangeQuery)
// 2. Query execution (Exec)
// 3. Plan optimization
// 4. Operator creation
// 5. Storage operations

package tracing

import (
	"context"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

var (
	// tracer is the global tracer used by the promql-engine.
	tracer = opentracing.GlobalTracer()
	mu     sync.RWMutex
)

// SetTracer sets the global tracer to be used by promql-engine.
// If nil is passed, it reverts to using the opentracing.GlobalTracer().
func SetTracer(t opentracing.Tracer) {
	mu.Lock()
	defer mu.Unlock()
	if t == nil {
		tracer = opentracing.GlobalTracer()
		return
	}
	tracer = t
}

// GetTracer returns the currently active tracer.
func GetTracer() opentracing.Tracer {
	mu.RLock()
	defer mu.RUnlock()
	return tracer
}

// StartSpan starts a new span with the given operation name.
// It uses the global tracer set by SetTracer or opentracing.GlobalTracer() if none was set.
func StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return GetTracer().StartSpan(operationName, opts...)
}

// StartSpanFromContext starts a new span from the given context.
// It returns the span and a context that includes the span.
func StartSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContextWithTracer(ctx, GetTracer(), operationName, opts...)
}

// InjectSpan injects a span into a carrier for transport across process boundaries.
func InjectSpan(span opentracing.Span, format interface{}, carrier interface{}) error {
	return GetTracer().Inject(span.Context(), format, carrier)
}

// ExtractSpan extracts a span from a carrier.
func ExtractSpan(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return GetTracer().Extract(format, carrier)
}

// LogError logs an error to the span and sets the span to error state.
func LogError(span opentracing.Span, err error) {
	ext.Error.Set(span, true)
	span.LogFields(log.Error(err))
}

// ChildSpan creates a child span with the given operation name from a parent span.
func ChildSpan(parent opentracing.Span, operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	childOpts := append([]opentracing.StartSpanOption{opentracing.ChildOf(parent.Context())}, opts...)
	return StartSpan(operationName, childOpts...)
}
