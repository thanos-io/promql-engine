# PromQL Query Engine

A multi-threaded implementation of a PromQL Query Engine based on the [Volcano/Iterator model](https://paperhub.s3.amazonaws.com/dace52a42c07f7f8348b08dc2b186061.pdf).

The project is currently under active development.

## Roadmap

The engine intends to have full compatibility with the original engine used in Prometheus. Since implementing the full specification will take time, we aim to add support for most commonly used expressions while falling back to the original engine for operations that are not yet supported. This will allow us to have smaller and faster releases, and gather feedback on regular basis. Instructions on using the engine will be added after we have enough confidence on its correctness.

The following table shows operations which are currently supported by the engine

| Type                   | Supported                                                                 | Priority |
|------------------------|---------------------------------------------------------------------------|----------|
| Binary expressions     | Full support                                                              |          |
| Histogram quantile     | Partial support (no support for native histograms)                        | High     |
| Aggregations           | Full support except for `count_values`                                    | Medium   |
| Aggregations over time | Full support except for `absent_over_time` and `quantile_over_time`       | Medium   |
| Functions              | Partial support (`clamp_min`, `clamp_max`, `changes` and `rate` variants) | Medium   |

## Design

At the beginning of a PromQL query execution, the query engine computes a physical plan consisting of multiple independent operators, each responsible for calculating one part of the query expression.

Operators are assembled in a tree-like structure with every operator calling `Next()` on its dependants until there is no more data to be returned. The result of the `Next()` function is a *column vector* (also called a *step vector*) with elements in the vector representing samples with the same timestamp from different time series.

<p align="center">
  <img src="./docs/assets/design.png"/>
</p>

This model allows for samples from individual time series to flow one execution step at a time from the left-most operators to the one at the very right. Since most PromQL expressions are aggregations, samples are reduced in number as they are pulled by the operators on the right. Because of this, samples from original timeseries can be decoded and kept in memory in batches instead of being fully expanded.

In addition to operators that have a one-to-one mapping with PromQL constructs, the Volcano model also describes so-called Exchange operators which can be used for flow control and optimizations, such as concurrency or batched selects. An example of an *Exchange* operator is described in the [Intra-operator parallelism](#intra-operator-parallelism) section.

### Inter-operator parallelism

Since operators are independent and rely on a common interface for pulling data, they can be run in parallel to each other. As soon as one operator has processed data from an evaluation step, it can pass the result onward so that its upstream can immediately start working on it.

<p align="center">
  <img src="./docs/assets/promql-pipeline.png"/>
</p>

### Intra-operator parallelism

Parallelism can also be added within individual operators using a parallel coalesce exchange operator. Such exchange operators are indistinguishable from regular operators to their upstreams since they respect the same `Next()` interface.

<p align="center">
  <img src="./docs/assets/parallel-coalesce.png"/>
</p>

### Memory management

#### Step vector allocations

One challenge with the streamed execution model is knowing how much memory to allocate in each operator for each step.

To work around this issue, operators expose a `Series()` method which returns the labels for all time series that they will ever produce (for all `Next()` calls). Operators at the very bottom of the tree, like vector and matrix selectors, have this information since they are responsible for loading data from storage. Other operators can then call `Series()` on the downstream operator and pre-compute all possible outputs.

Even though this might look like an expensive operation, its cost is identical to just one evaluation step. Knowing sizes of input and output vectors also allows us to:
* allocate memory very precisely by properly sizing vector pools (see section bellow),
* use arrays instead of maps for indexing data, leading to faster execution times due to having less allocations and using index-based lookups, and
* use tight loops in operators by eliminating conditional statements associated with maps.

#### Vector pools

Since time series are decoded one step at a time, vectors between execution execution steps can be recycled manually instead of relying on the garbage collector. Each operator has its own pool that it uses to allocate new step vectors and send results to its upstream. Whenever the upstream operator is finished with processing a step vector, it will return that vector to the pool of its downstream so that it can be reused again for subsequent steps.

#### Memory limits

There are currently no mechanisms to apply memory limits to queries within the engine. This is a highly desirable feature, and we would like to explore ways in which we can support it.

### Concurrency control

The current implementation uses goroutines very liberally which means the query will use as many cores as possible. Limiting the number of cores which a query can use is not yet implemented but we would eventually like to have support for it.

### Plan optimization

The current implementation creates a physical plan directly from the PromQL abstract syntax tree. Plan optimizations not yet implemented and would require having a logical plan as an intermediary step.

## Distributed execution mode

The engine supports a distributed mode where aggregations can be delegated to multiple remote engines, each responsible for an independent dataset. This mode is currently implemented through an optimizer which rewrites a query as a combination of multiple remote and one local aggregation. For example, when two remote engines are available, a query like:

```
sum(rate(http_request_total[4m]))
```

would be rewritten as

```
sum(
  coalesce(
    sum(rate(http_request_total[4m])) # remote engine 1
    sum(rate(http_request_total[4m])) # remote engine 2
  )
)
```

The inner aggregations are forwarded to remote engines and the global result is completed in memory.

An engine using the distributed mode can be created through the `NewDistributedEngine` function. The user is expected to pass an implementation of `RemoteEndpoints` which has a single `Engines()` method. When invoked, `Engines()` should return all remote engines that can be used for a single query. The `Engines()` method is called separately for each individual query which allows the `RemoteEndpoints` implementation to do continuous service discovery and inject engines as they become available.

The interfaces used for remote execution can be found in [api](https://pkg.go.dev/github.com/thanos-community/promql-engine/api) package. Note that the `RemoteEngine` interface has a `NewRangeQuery` method, similar to the one in the Prometheus [v1.QueryEngine](https://pkg.go.dev/github.com/prometheus/prometheus@v0.42.0/web/api/v1#QueryEngine) interface. It is up to the user of the library to implement this method as they see fit. An example implementation could be to forward the query to an HTTP `/api/v1/query_range` endpoint of a Prometheus instance. In Thanos, this method is implemented as a gRPC call to a Thanos Querier.

For more details on the overall design, please refer to the [proposal](https://github.com/thanos-io/thanos/blob/main/docs/proposals-accepted/202301-distributed-query-execution.md) in the Thanos project.

## Continuous benchmark

If you are interested in the benchmark results captured by continuous benchmark, please check [here](https://thanos-community.github.io/promql-engine/dev/bench/).

## Latest benchmarks

These are the latest benchmarks captured on an Apple M1 Pro processor.

Note that memory usage is higher when executing a query with parallelism greater than 1. This is due to the fact that the engine is able to execute multiple operations at once (e.g. decode chunks from multiple series at the same time), which requires using independent buffers for each parallel operation.

Single core benchmarks

```markdown
name                                                old time/op    new time/op    delta
RangeQuery/vector_selector                            33.5ms ± 3%    43.4ms ± 3%  +29.59%  (p=0.008 n=5+5)
RangeQuery/sum                                        46.6ms ± 1%    34.3ms ± 2%  -26.37%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod                                  145ms ± 1%      46ms ± 3%  -68.36%  (p=0.008 n=5+5)
RangeQuery/topk                                       46.7ms ± 2%    37.0ms ± 6%  -20.79%  (p=0.008 n=5+5)
RangeQuery/bottomk                                    46.9ms ± 1%    35.3ms ± 7%  -24.72%  (p=0.008 n=5+5)
RangeQuery/rate                                       65.7ms ± 1%    72.2ms ± 2%   +9.88%  (p=0.008 n=5+5)
RangeQuery/sum_rate                                   76.6ms ± 1%    61.9ms ± 1%  -19.16%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate                                 180ms ± 3%      74ms ± 7%  -58.94%  (p=0.008 n=5+5)
RangeQuery/quantile_with_variable_parameter            263ms ± 6%      99ms ± 3%  -62.38%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one            119ms ± 1%      31ms ± 3%  -74.20%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one           396ms ± 1%      69ms ± 1%  -82.52%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar     241ms ± 1%      51ms ± 1%  -78.85%  (p=0.008 n=5+5)
RangeQuery/unary_negation                             35.6ms ± 2%    46.6ms ± 5%  +31.00%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison                205ms ± 3%      57ms ± 4%  -72.48%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector                     33.2ms ±10%    43.2ms ± 3%  +30.20%  (p=0.008 n=5+5)
RangeQuery/at_modifier_                               18.7ms ± 3%    18.0ms ± 5%     ~     (p=0.095 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector    17.9ms ± 2%    17.2ms ± 2%   -3.78%  (p=0.008 n=5+5)
RangeQuery/clamp                                       252ms ± 7%      62ms ± 4%  -75.28%  (p=0.008 n=5+5)
RangeQuery/clamp_min                                   253ms ± 5%      59ms ± 3%  -76.75%  (p=0.008 n=5+5)
RangeQuery/complex_func_query                          455ms ± 2%      68ms ± 3%  -85.06%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query                      265ms ± 3%      89ms ± 3%  -66.33%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query                      273ms ± 1%      91ms ± 2%  -66.43%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile                          579ms ± 2%     204ms ± 5%  -64.68%  (p=0.008 n=5+5)
RangeQuery/sort                                        299ms ± 1%      43ms ± 2%  -85.57%  (p=0.008 n=5+5)
RangeQuery/sort_desc                                   294ms ± 2%      44ms ± 3%  -84.97%  (p=0.008 n=5+5)
NativeHistograms/selector                              620ms ± 1%     662ms ± 5%   +6.79%  (p=0.008 n=5+5)
NativeHistograms/sum                                   1.21s ± 7%     1.01s ± 1%  -16.42%  (p=0.008 n=5+5)
NativeHistograms/rate                                  4.57s ± 3%     4.49s ± 1%     ~     (p=0.310 n=5+5)
NativeHistograms/sum_rate                              5.04s ± 1%     4.79s ± 1%   -4.99%  (p=0.008 n=5+5)
NativeHistograms/histogram_sum                         930ms ± 2%    1068ms ± 6%  +14.77%  (p=0.008 n=5+5)
NativeHistograms/histogram_count                       980ms ± 7%    1059ms ± 7%     ~     (p=0.095 n=5+5)
NativeHistograms/histogram_quantile                    1.20s ± 1%     1.02s ± 4%  -14.80%  (p=0.008 n=5+5)

name                                                old alloc/op   new alloc/op   delta
RangeQuery/vector_selector                            24.5MB ± 0%    38.1MB ± 0%  +55.28%  (p=0.008 n=5+5)
RangeQuery/sum                                        7.13MB ± 0%   10.20MB ± 0%  +43.11%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod                                 79.9MB ± 0%    22.7MB ± 0%  -71.61%  (p=0.008 n=5+5)
RangeQuery/topk                                       7.38MB ± 0%   12.68MB ± 0%  +71.84%  (p=0.008 n=5+5)
RangeQuery/bottomk                                    7.44MB ± 0%   12.72MB ± 0%  +71.02%  (p=0.029 n=4+4)
RangeQuery/rate                                       25.6MB ± 0%    41.0MB ± 0%  +60.30%  (p=0.008 n=5+5)
RangeQuery/sum_rate                                   8.19MB ± 0%   13.09MB ± 0%  +59.73%  (p=0.016 n=5+4)
RangeQuery/sum_by_rate                                80.7MB ± 0%    25.5MB ± 0%  -68.43%  (p=0.008 n=5+5)
RangeQuery/quantile_with_variable_parameter            174MB ± 0%      39MB ± 0%  -77.60%  (p=0.016 n=5+4)
RangeQuery/binary_operation_with_one_to_one           16.5MB ± 0%    21.6MB ± 0%  +30.83%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one          72.0MB ± 0%    55.8MB ± 0%  -22.54%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar    39.1MB ± 0%    40.2MB ± 0%   +2.80%  (p=0.008 n=5+5)
RangeQuery/unary_negation                             25.6MB ± 0%    39.4MB ± 0%  +54.13%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison               37.7MB ± 0%    39.9MB ± 0%   +5.63%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector                     23.0MB ± 0%    36.6MB ± 0%  +58.83%  (p=0.008 n=5+5)
RangeQuery/at_modifier_                               39.8MB ± 0%    33.1MB ± 0%  -16.75%  (p=0.008 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector    39.6MB ± 0%    32.9MB ± 0%  -16.83%  (p=0.016 n=4+5)
RangeQuery/clamp                                      39.2MB ± 0%    38.5MB ± 0%   -1.69%  (p=0.016 n=5+4)
RangeQuery/clamp_min                                  39.1MB ± 0%    38.5MB ± 0%   -1.75%  (p=0.008 n=5+5)
RangeQuery/complex_func_query                         53.8MB ± 0%    40.6MB ± 0%  -24.54%  (p=0.016 n=5+4)
RangeQuery/func_within_func_query                     40.2MB ± 0%    41.1MB ± 0%   +2.18%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query                     40.2MB ± 0%    41.1MB ± 0%   +2.19%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile                         47.5MB ± 0%    57.9MB ± 0%  +21.88%  (p=0.016 n=5+4)
RangeQuery/sort                                       37.8MB ± 0%    38.1MB ± 0%   +0.67%  (p=0.008 n=5+5)
RangeQuery/sort_desc                                  37.8MB ± 0%    38.1MB ± 0%   +0.67%  (p=0.008 n=5+5)
NativeHistograms/selector                              761MB ± 0%     774MB ± 0%   +1.72%  (p=0.016 n=4+5)
NativeHistograms/sum                                   943MB ± 0%     931MB ± 0%   -1.21%  (p=0.008 n=5+5)
NativeHistograms/rate                                 2.86GB ± 0%    2.87GB ± 0%   +0.53%  (p=0.029 n=4+4)
NativeHistograms/sum_rate                             3.04GB ± 0%    3.03GB ± 0%   -0.41%  (p=0.016 n=4+5)
NativeHistograms/histogram_sum                         786MB ± 0%     775MB ± 0%   -1.42%  (p=0.008 n=5+5)
NativeHistograms/histogram_count                       787MB ± 0%     774MB ± 0%   -1.63%  (p=0.016 n=5+4)
NativeHistograms/histogram_quantile                    942MB ± 0%     932MB ± 0%   -1.14%  (p=0.008 n=5+5)

name                                                old allocs/op  new allocs/op  delta
RangeQuery/vector_selector                             99.1k ± 0%    111.9k ± 0%  +12.96%  (p=0.016 n=5+4)
RangeQuery/sum                                          103k ± 0%      107k ± 0%   +3.41%  (p=0.016 n=5+4)
RangeQuery/sum_by_pod                                   598k ± 0%      202k ± 0%  -66.28%  (p=0.008 n=5+5)
RangeQuery/topk                                         108k ± 0%      114k ± 0%   +5.18%  (p=0.008 n=5+5)
RangeQuery/bottomk                                      109k ± 0%      116k ± 0%   +6.20%  (p=0.008 n=5+5)
RangeQuery/rate                                         111k ± 0%      136k ± 0%  +22.33%  (p=0.016 n=4+5)
RangeQuery/sum_rate                                     115k ± 0%      131k ± 0%  +13.45%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate                                  608k ± 0%      226k ± 0%  -62.89%  (p=0.008 n=5+5)
RangeQuery/quantile_with_variable_parameter            1.67M ± 0%     0.58M ± 0%  -65.23%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one            75.3k ± 0%     89.2k ± 0%  +18.55%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one            637k ± 0%      173k ± 0%  -72.86%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar      117k ± 0%      116k ± 0%   -1.24%  (p=0.008 n=5+5)
RangeQuery/unary_negation                               111k ± 0%      124k ± 0%  +11.82%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison                 105k ± 0%      113k ± 0%   +7.22%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector                      73.1k ± 0%     86.0k ± 0%  +17.69%  (p=0.008 n=5+5)
RangeQuery/at_modifier_                                74.1k ± 0%     62.6k ± 0%  -15.53%  (p=0.000 n=5+4)
RangeQuery/at_modifier_with_positive_offset_vector     68.1k ± 0%     56.6k ± 0%  -16.90%  (p=0.000 n=5+4)
RangeQuery/clamp                                        118k ± 0%      116k ± 0%   -1.69%  (p=0.016 n=5+4)
RangeQuery/clamp_min                                    117k ± 0%      115k ± 0%   -1.75%  (p=0.008 n=5+5)
RangeQuery/complex_func_query                           136k ± 0%      120k ± 0%  -11.97%  (p=0.016 n=5+4)
RangeQuery/func_within_func_query                       130k ± 0%      137k ± 0%   +5.40%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query                       130k ± 0%      137k ± 0%   +5.40%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile                           617k ± 0%      656k ± 0%   +6.29%  (p=0.016 n=5+4)
RangeQuery/sort                                         106k ± 0%      112k ± 0%   +6.10%  (p=0.008 n=5+5)
RangeQuery/sort_desc                                    106k ± 0%      112k ± 0%   +6.10%  (p=0.008 n=5+5)
NativeHistograms/selector                              9.63M ± 0%     9.64M ± 0%   +0.14%  (p=0.016 n=4+5)
NativeHistograms/sum                                   11.1M ± 0%     11.1M ± 0%   +0.02%  (p=0.008 n=5+5)
NativeHistograms/rate                                  34.1M ± 0%     34.1M ± 0%   +0.08%  (p=0.016 n=5+4)
NativeHistograms/sum_rate                              35.6M ± 0%     35.6M ± 0%   +0.05%  (p=0.008 n=5+5)
NativeHistograms/histogram_sum                         9.65M ± 0%     9.64M ± 0%   -0.04%  (p=0.008 n=5+5)
NativeHistograms/histogram_count                       9.65M ± 0%     9.64M ± 0%   -0.04%  (p=0.008 n=5+5)
NativeHistograms/histogram_quantile                    11.1M ± 0%     11.1M ± 0%   +0.02%  (p=0.008 n=5+5)
```

Multi-core (8 core) benchmarks

```markdown
name                                                  old time/op    new time/op    delta
RangeQuery/vector_selector-8                            31.1ms ± 1%    14.7ms ± 1%  -52.66%  (p=0.008 n=5+5)
RangeQuery/sum-8                                        49.3ms ± 2%    11.0ms ± 0%  -77.74%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod-8                                  138ms ± 4%      15ms ± 0%  -89.16%  (p=0.016 n=5+4)
RangeQuery/topk-8                                       47.7ms ± 4%    11.0ms ± 0%  -77.03%  (p=0.008 n=5+5)
RangeQuery/bottomk-8                                    48.3ms ± 3%    11.1ms ± 5%  -76.95%  (p=0.008 n=5+5)
RangeQuery/rate-8                                       61.4ms ± 2%    21.1ms ± 3%  -65.63%  (p=0.008 n=5+5)
RangeQuery/sum_rate-8                                   77.9ms ± 1%    19.0ms ± 4%  -75.63%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate-8                                 165ms ± 1%      22ms ± 3%  -86.80%  (p=0.008 n=5+5)
RangeQuery/quantile_with_variable_parameter-8            234ms ± 3%      25ms ± 1%  -89.17%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one-8            121ms ± 2%      14ms ± 1%  -88.53%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one-8           405ms ± 2%      30ms ± 1%  -92.49%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar-8     245ms ± 2%      20ms ± 0%  -91.88%  (p=0.008 n=5+5)
RangeQuery/unary_negation-8                             32.3ms ± 3%    15.5ms ± 2%  -52.10%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison-8                206ms ± 2%      21ms ± 2%  -89.78%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector-8                     27.6ms ± 1%    13.9ms ± 4%  -49.83%  (p=0.008 n=5+5)
RangeQuery/at_modifier_-8                               12.6ms ± 2%    10.0ms ± 2%  -20.88%  (p=0.008 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector-8    12.0ms ± 3%     9.6ms ± 1%  -19.73%  (p=0.008 n=5+5)
RangeQuery/clamp-8                                       246ms ± 4%      31ms ± 4%  -87.26%  (p=0.008 n=5+5)
RangeQuery/clamp_min-8                                   251ms ± 4%      27ms ±17%  -89.10%  (p=0.008 n=5+5)
RangeQuery/complex_func_query-8                          480ms ± 5%      38ms ± 5%  -92.15%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query-8                      279ms ± 1%      32ms ± 1%  -88.59%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query-8                      274ms ± 6%      32ms ± 2%  -88.28%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile-8                          545ms ± 5%      97ms ± 1%  -82.15%  (p=0.008 n=5+5)
RangeQuery/sort-8                                        301ms ± 7%      15ms ± 3%  -94.92%  (p=0.008 n=5+5)
RangeQuery/sort_desc-8                                   295ms ± 3%      15ms ± 1%  -94.88%  (p=0.008 n=5+5)
NativeHistograms/selector-8                              417ms ± 3%     217ms ± 3%  -47.95%  (p=0.008 n=5+5)
NativeHistograms/sum-8                                   897ms ± 1%     271ms ± 2%  -69.74%  (p=0.008 n=5+5)
NativeHistograms/rate-8                                  3.76s ± 1%     1.27s ± 2%  -66.30%  (p=0.008 n=5+5)
NativeHistograms/sum_rate-8                              4.24s ± 3%     1.27s ± 4%  -70.12%  (p=0.008 n=5+5)
NativeHistograms/histogram_sum-8                         683ms ± 1%     429ms ± 2%  -37.23%  (p=0.008 n=5+5)
NativeHistograms/histogram_count-8                       681ms ± 1%     423ms ± 1%  -37.97%  (p=0.008 n=5+5)
NativeHistograms/histogram_quantile-8                    903ms ± 3%     268ms ± 2%  -70.32%  (p=0.008 n=5+5)

name                                                  old alloc/op   new alloc/op   delta
RangeQuery/vector_selector-8                            24.5MB ± 0%    38.6MB ± 0%  +57.55%  (p=0.008 n=5+5)
RangeQuery/sum-8                                        7.12MB ± 0%    9.89MB ± 0%  +39.06%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod-8                                 79.9MB ± 0%    23.7MB ± 0%  -70.36%  (p=0.008 n=5+5)
RangeQuery/topk-8                                       7.43MB ± 0%   10.95MB ± 0%  +47.39%  (p=0.008 n=5+5)
RangeQuery/bottomk-8                                    7.46MB ± 0%   10.94MB ± 1%  +46.61%  (p=0.008 n=5+5)
RangeQuery/rate-8                                       25.6MB ± 0%    41.4MB ± 0%  +61.73%  (p=0.008 n=5+5)
RangeQuery/sum_rate-8                                   8.18MB ± 0%   12.67MB ± 0%  +54.81%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate-8                                80.7MB ± 0%    25.3MB ± 1%  -68.68%  (p=0.008 n=5+5)
RangeQuery/quantile_with_variable_parameter-8            174MB ± 0%      41MB ± 0%  -76.60%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one-8           16.5MB ± 0%    22.4MB ± 0%  +35.63%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one-8          72.0MB ± 0%    56.9MB ± 0%  -20.99%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar-8    39.1MB ± 0%    41.6MB ± 0%   +6.39%  (p=0.008 n=5+5)
RangeQuery/unary_negation-8                             25.6MB ± 0%    40.0MB ± 0%  +56.65%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison-8               37.7MB ± 0%    41.3MB ± 0%   +9.45%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector-8                     23.0MB ± 0%    37.1MB ± 0%  +61.21%  (p=0.008 n=5+5)
RangeQuery/at_modifier_-8                               39.8MB ± 0%    33.1MB ± 0%  -16.67%  (p=0.008 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector-8    39.6MB ± 0%    33.0MB ± 0%  -16.74%  (p=0.008 n=5+5)
RangeQuery/clamp-8                                      39.1MB ± 0%    39.1MB ± 0%   -0.22%  (p=0.008 n=5+5)
RangeQuery/clamp_min-8                                  39.1MB ± 0%    39.0MB ± 0%   -0.27%  (p=0.008 n=5+5)
RangeQuery/complex_func_query-8                         53.8MB ± 0%    41.9MB ± 0%  -22.09%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query-8                     40.2MB ± 0%    41.7MB ± 0%   +3.68%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query-8                     40.2MB ± 0%    41.6MB ± 0%   +3.66%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile-8                         47.5MB ± 0%    60.5MB ± 0%  +27.26%  (p=0.008 n=5+5)
RangeQuery/sort-8                                       37.8MB ± 0%    38.6MB ± 0%   +2.12%  (p=0.008 n=5+5)
RangeQuery/sort_desc-8                                  37.8MB ± 0%    38.6MB ± 0%   +2.12%  (p=0.016 n=4+5)
NativeHistograms/selector-8                              761MB ± 0%     775MB ± 0%   +1.90%  (p=0.008 n=5+5)
NativeHistograms/sum-8                                   940MB ± 0%     933MB ± 0%   -0.67%  (p=0.008 n=5+5)
NativeHistograms/rate-8                                 2.86GB ± 0%    2.88GB ± 0%   +0.61%  (p=0.008 n=5+5)
NativeHistograms/sum_rate-8                             3.04GB ± 0%    3.04GB ± 0%   -0.27%  (p=0.008 n=5+5)
NativeHistograms/histogram_sum-8                         786MB ± 0%     777MB ± 0%   -1.24%  (p=0.008 n=5+5)
NativeHistograms/histogram_count-8                       787MB ± 0%     777MB ± 0%   -1.29%  (p=0.008 n=5+5)
NativeHistograms/histogram_quantile-8                    940MB ± 0%     934MB ± 0%   -0.69%  (p=0.008 n=5+5)

name                                                  old allocs/op  new allocs/op  delta
RangeQuery/vector_selector-8                             98.8k ± 0%    113.6k ± 0%  +14.99%  (p=0.008 n=5+5)
RangeQuery/sum-8                                          103k ± 0%      108k ± 0%   +5.03%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod-8                                   598k ± 0%      204k ± 0%  -65.95%  (p=0.008 n=5+5)
RangeQuery/topk-8                                         109k ± 0%      116k ± 0%   +6.76%  (p=0.008 n=5+5)
RangeQuery/bottomk-8                                      110k ± 0%      115k ± 0%   +4.46%  (p=0.008 n=5+5)
RangeQuery/rate-8                                         111k ± 0%      138k ± 0%  +24.06%  (p=0.008 n=5+5)
RangeQuery/sum_rate-8                                     115k ± 0%      132k ± 0%  +14.87%  (p=0.016 n=4+5)
RangeQuery/sum_by_rate-8                                  608k ± 0%      227k ± 0%  -62.64%  (p=0.008 n=5+5)
RangeQuery/quantile_with_variable_parameter-8            1.66M ± 0%     0.58M ± 0%  -64.99%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one-8            75.2k ± 0%     93.4k ± 0%  +24.19%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one-8            637k ± 0%      177k ± 0%  -72.24%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar-8      117k ± 0%      118k ± 0%   +0.65%  (p=0.016 n=4+5)
RangeQuery/unary_negation-8                               111k ± 0%      126k ± 0%  +13.71%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison-8                 105k ± 0%      115k ± 0%   +9.37%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector-8                      72.9k ± 0%     87.8k ± 0%  +20.37%  (p=0.029 n=4+4)
RangeQuery/at_modifier_-8                                74.0k ± 0%     62.6k ± 0%  -15.40%  (p=0.008 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector-8     68.0k ± 0%     56.6k ± 0%  -16.77%  (p=0.008 n=5+5)
RangeQuery/clamp-8                                        117k ± 0%      117k ± 0%   +0.08%  (p=0.008 n=5+5)
RangeQuery/clamp_min-8                                    117k ± 0%      117k ± 0%     ~     (p=0.159 n=4+5)
RangeQuery/complex_func_query-8                           136k ± 0%      122k ± 0%  -10.31%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query-8                       129k ± 0%      139k ± 0%   +7.04%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query-8                       129k ± 0%      139k ± 0%   +7.03%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile-8                           617k ± 0%      658k ± 0%   +6.64%  (p=0.008 n=5+5)
RangeQuery/sort-8                                         105k ± 0%      114k ± 0%   +7.98%  (p=0.008 n=5+5)
RangeQuery/sort_desc-8                                    105k ± 0%      114k ± 0%   +7.98%  (p=0.016 n=4+5)
NativeHistograms/selector-8                              9.63M ± 0%     9.64M ± 0%   +0.16%  (p=0.008 n=5+5)
NativeHistograms/sum-8                                   11.1M ± 0%     11.1M ± 0%   +0.04%  (p=0.008 n=5+5)
NativeHistograms/rate-8                                  34.1M ± 0%     34.2M ± 0%   +0.09%  (p=0.008 n=5+5)
NativeHistograms/sum_rate-8                              35.6M ± 0%     35.6M ± 0%   +0.06%  (p=0.008 n=5+5)
NativeHistograms/histogram_sum-8                         9.65M ± 0%     9.65M ± 0%   -0.01%  (p=0.008 n=5+5)
NativeHistograms/histogram_count-8                       9.65M ± 0%     9.65M ± 0%   -0.01%  (p=0.008 n=5+5)
NativeHistograms/histogram_quantile-8                    11.1M ± 0%     11.1M ± 0%   +0.05%  (p=0.008 n=5+5)
```
