# Shared PromQL engine test cases

The YAML files in this directory hold engine test cases in a language agnostic
format. They are the source of truth: the Go tests in this repository read them
through [`testcases.go`](testcases.go), and engine implementations in other
languages can read the same files without going through Go.

## Files

| File                                       | Contents                                                    |
| ------------------------------------------ | ----------------------------------------------------------- |
| [`range_queries.yaml`](range_queries.yaml) | Range queries run against both this engine and Prometheus's. |

## Schema

```yaml
defaults:
  start_ms: 0 # start of the query range for cases that don't set one
  end_ms: 1800000 # end of the query range for cases that don't set one
  step_ms: 30000 # resolution step for cases that don't set one
tests:
  - name: group_left with included label changing across the window
    description: |
      Optional. Explains what the case covers.
    load: |
      load 30s
      metric_a{pod="x"} 1 1 1 1 1
      metric_b{pod="x", ns="a"} 1 1 _ _ _
    query: metric_a * on (pod) group_left (ns) metric_b
    start_ms: 0 # optional, defaults to defaults.start_ms
    end_ms: 120000 # optional, defaults to defaults.end_ms
    step_ms: 30000 # optional, defaults to defaults.step_ms
```

- `name` is required and identifies the case. It is used as the Go subtest name.
- `query` is required and holds a single PromQL expression.
- `load` is a Prometheus test script `load` block that defines the input series.
  It uses the same syntax as Prometheus's own `.test` files, so an existing
  parser for those files can read it. Cases that need no data omit it.
- All timestamps and durations are integer milliseconds. `start_ms` and `end_ms`
  are Unix timestamps and may be negative; `step_ms` is a duration.

## How the cases are used

Each case is a differential test: the query is executed against both this engine
and Prometheus's own engine over the same range, and the two results must be
equal. There are therefore no expected results in the file — the reference
implementation provides them. An engine in another language can either compare
against its own reference implementation the same way, or record the results of
a known good run as a golden file.

## Adding a case

Add an entry to the relevant file. Keep `name` unique and descriptive, and only
set `start_ms`, `end_ms` or `step_ms` when the case actually needs a range other
than the defaults.
