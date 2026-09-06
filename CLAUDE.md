# thanos-promql-engine Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-03-20

## Active Technologies
- Go 1.26 (toolchain 1.26.4) + `go.opentelemetry.io/otel v1.38.0` (trace, attribute, codes, baggage), `github.com/prometheus/prometheus v0.308.0` (001-otel-tracing)
- no database (engine is a library; storage is provided by the embedding program backed by Parquet or TSDB) (001-otel-tracing)

- (001-otel-tracing)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for

## Code Style

: Follow standard conventions

Check for gofmt compliance in changed files:

    gofmt -w $(git diff --stat --name-only origin/main ':*.go')

Check golangci compliance on each commit:

    golangci-lint run --new

## Recent Changes
- 001-otel-tracing: Added Go 1.24 (toolchain 1.24.4) + `go.opentelemetry.io/otel v1.38.0` (trace, attribute, codes, baggage), `github.com/prometheus/prometheus v0.308.0`

- 001-otel-tracing: Added

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
