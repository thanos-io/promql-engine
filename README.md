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

In addition to implementing multi-threading, we would ultimately like to end up with a distributed execution model.

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

## Latest benchmarks

These are the latest benchmarks captured on an Apple M1 Pro processor.

Note that memory usage is higher when executing a query with parallelism greater than 1. This is due to the fact that the engine is able to execute multiple operations at once (e.g. decode chunks from multiple series at the same time), which requires using independent buffers for each parallel operation.

Single core benchmarks

```markdown
name                                                old time/op    new time/op    delta
RangeQuery/vector_selector                            35.5ms ± 2%    46.7ms ±13%  +31.42%  (p=0.008 n=5+5)
RangeQuery/sum                                        48.1ms ± 3%    35.6ms ± 2%  -25.97%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod                                  151ms ± 3%      47ms ± 4%  -68.77%  (p=0.008 n=5+5)
RangeQuery/topk                                       48.1ms ± 5%    38.8ms ± 5%  -19.31%  (p=0.008 n=5+5)
RangeQuery/bottomk                                    49.3ms ± 5%    39.1ms ± 3%  -20.65%  (p=0.008 n=5+5)
RangeQuery/rate                                       66.1ms ± 1%    75.5ms ± 2%  +14.26%  (p=0.008 n=5+5)
RangeQuery/sum_rate                                   78.4ms ± 3%    62.9ms ± 3%  -19.75%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate                                 184ms ± 5%      73ms ± 6%  -60.32%  (p=0.008 n=5+5)
RangeQuery/quantile_with_variable_parameter            269ms ± 2%     102ms ± 2%  -62.08%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one            126ms ±11%      30ms ± 3%  -76.01%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one           407ms ± 3%      65ms ± 3%  -84.03%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar     248ms ± 1%      52ms ± 3%  -78.86%  (p=0.008 n=5+5)
RangeQuery/unary_negation                             37.3ms ± 4%    48.0ms ± 1%  +28.63%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison                201ms ± 1%      52ms ± 1%  -74.09%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector                     33.1ms ± 2%    41.0ms ± 3%  +24.05%  (p=0.008 n=5+5)
RangeQuery/at_modifier_                               19.0ms ± 3%    17.9ms ± 1%   -6.12%  (p=0.008 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector    19.2ms ±15%    17.5ms ± 2%   -9.05%  (p=0.008 n=5+5)
RangeQuery/clamp                                       244ms ± 2%      60ms ± 7%  -75.25%  (p=0.008 n=5+5)
RangeQuery/clamp_min                                   241ms ± 2%      59ms ± 8%  -75.58%  (p=0.008 n=5+5)
RangeQuery/complex_func_query                          444ms ± 1%      65ms ± 4%  -85.35%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query                      260ms ± 7%      89ms ± 3%  -65.87%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query                      267ms ± 2%      90ms ± 3%  -66.13%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile                          551ms ± 3%     201ms ± 4%  -63.55%  (p=0.008 n=5+5)

name                                                old alloc/op   new alloc/op   delta
RangeQuery/vector_selector                            25.6MB ± 0%    38.3MB ± 0%  +49.55%  (p=0.016 n=4+5)
RangeQuery/sum                                        8.31MB ± 0%   10.35MB ± 0%  +24.50%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod                                 96.5MB ± 0%    22.5MB ± 0%  -76.67%  (p=0.008 n=5+5)
RangeQuery/topk                                       8.61MB ± 0%   12.92MB ± 0%  +49.96%  (p=0.008 n=5+5)
RangeQuery/bottomk                                    8.64MB ± 0%   12.91MB ± 0%  +49.50%  (p=0.008 n=5+5)
RangeQuery/rate                                       26.7MB ± 0%    41.2MB ± 0%  +54.32%  (p=0.008 n=5+5)
RangeQuery/sum_rate                                   9.37MB ± 0%   13.21MB ± 0%  +40.97%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate                                89.6MB ± 0%    25.3MB ± 0%  -71.77%  (p=0.008 n=5+5)
RangeQuery/quantile_with_variable_parameter            192MB ± 0%      39MB ± 0%  -79.48%  (p=0.016 n=5+4)
RangeQuery/binary_operation_with_one_to_one           17.3MB ± 0%    18.0MB ± 0%   +3.62%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one          73.5MB ± 0%    44.7MB ± 0%  -39.21%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar    40.3MB ± 0%    40.5MB ± 0%   +0.53%  (p=0.008 n=5+5)
RangeQuery/unary_negation                             26.7MB ± 0%    39.7MB ± 0%  +48.63%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison               38.8MB ± 0%    40.1MB ± 0%   +3.21%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector                     23.8MB ± 0%    36.5MB ± 0%  +53.29%  (p=0.008 n=5+5)
RangeQuery/at_modifier_                               40.6MB ± 0%    33.1MB ± 0%  -18.52%  (p=0.016 n=5+4)
RangeQuery/at_modifier_with_positive_offset_vector    40.4MB ± 0%    32.9MB ± 0%  -18.61%  (p=0.008 n=5+5)
RangeQuery/clamp                                      40.3MB ± 0%    38.7MB ± 0%   -3.82%  (p=0.008 n=5+5)
RangeQuery/clamp_min                                  40.2MB ± 0%    38.7MB ± 0%   -3.84%  (p=0.008 n=5+5)
RangeQuery/complex_func_query                         54.8MB ± 0%    40.8MB ± 0%  -25.57%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query                     41.2MB ± 0%    41.2MB ± 0%   +0.04%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query                     41.2MB ± 0%    41.2MB ± 0%   +0.04%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile                         51.7MB ± 0%    58.8MB ± 0%  +13.56%  (p=0.016 n=5+4)

name                                                old allocs/op  new allocs/op  delta
RangeQuery/vector_selector                              120k ± 0%      127k ± 0%   +6.12%  (p=0.008 n=5+5)
RangeQuery/sum                                          124k ± 0%      122k ± 0%   -2.19%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod                                   619k ± 0%      207k ± 0%  -66.62%  (p=0.008 n=5+5)
RangeQuery/topk                                         130k ± 0%      130k ± 0%   -0.52%  (p=0.008 n=5+5)
RangeQuery/bottomk                                      131k ± 0%      129k ± 0%   -1.33%  (p=0.008 n=5+5)
RangeQuery/rate                                         129k ± 0%      151k ± 0%  +17.32%  (p=0.016 n=5+4)
RangeQuery/sum_rate                                     133k ± 0%      146k ± 0%   +9.18%  (p=0.016 n=4+5)
RangeQuery/sum_by_rate                                  626k ± 0%      231k ± 0%  -63.15%  (p=0.008 n=5+5)
RangeQuery/quantile_with_variable_parameter            1.71M ± 0%     0.61M ± 0%  -64.33%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one            89.1k ± 0%     98.6k ± 0%  +10.65%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one            662k ± 0%      192k ± 0%  -70.92%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar      135k ± 0%      131k ± 0%   -2.95%  (p=0.008 n=5+5)
RangeQuery/unary_negation                               129k ± 0%      139k ± 0%   +8.25%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison                 126k ± 0%      128k ± 0%   +1.65%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector                      90.8k ± 0%     98.1k ± 0%   +7.99%  (p=0.016 n=5+4)
RangeQuery/at_modifier_                                92.0k ± 0%     74.6k ± 0%  -18.88%  (p=0.008 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector     86.0k ± 0%     68.6k ± 0%     ~     (p=0.079 n=4+5)
RangeQuery/clamp                                        135k ± 0%      131k ± 0%   -3.34%  (p=0.008 n=5+5)
RangeQuery/clamp_min                                    135k ± 0%      130k ± 0%   -3.38%  (p=0.008 n=5+5)
RangeQuery/complex_func_query                           150k ± 0%      135k ± 0%  -10.50%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query                       144k ± 0%      152k ± 0%   +5.19%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query                       144k ± 0%      152k ± 0%   +5.19%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile                           694k ± 0%      700k ± 0%   +0.81%  (p=0.008 n=5+5)
```

Multi-core (8 core) benchmarks

```markdown
name                                                  old time/op    new time/op    delta
RangeQuery/vector_selector-8                            32.4ms ± 1%    17.0ms ± 1%  -47.47%  (p=0.008 n=5+5)
RangeQuery/sum-8                                        50.9ms ± 4%    10.6ms ± 1%  -79.12%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod-8                                  145ms ± 3%      16ms ± 2%  -88.69%  (p=0.008 n=5+5)
RangeQuery/topk-8                                       50.5ms ± 3%    13.2ms ± 2%  -73.80%  (p=0.008 n=5+5)
RangeQuery/bottomk-8                                    49.2ms ± 1%    14.1ms ±13%  -71.28%  (p=0.008 n=5+5)
RangeQuery/rate-8                                       63.4ms ± 3%    24.8ms ± 4%  -60.98%  (p=0.008 n=5+5)
RangeQuery/sum_rate-8                                   78.1ms ± 3%    18.3ms ± 3%  -76.51%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate-8                                 169ms ± 2%      25ms ± 6%  -85.24%  (p=0.008 n=5+5)
RangeQuery/quantile_with_variable_parameter-8            254ms ± 6%      31ms ±17%  -87.95%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one-8            129ms ± 6%      15ms ± 2%  -88.70%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one-8           448ms ±25%      29ms ± 2%  -93.48%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar-8     257ms ± 6%      21ms ± 1%  -91.98%  (p=0.008 n=5+5)
RangeQuery/unary_negation-8                             33.6ms ± 1%    17.6ms ± 3%  -47.69%  (p=0.016 n=4+5)
RangeQuery/vector_and_scalar_comparison-8                201ms ± 1%      21ms ± 1%  -89.55%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector-8                     29.5ms ± 8%    15.4ms ± 1%  -47.75%  (p=0.008 n=5+5)
RangeQuery/at_modifier_-8                               14.3ms ±13%    12.6ms ± 1%  -11.89%  (p=0.008 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector-8    13.2ms ± 7%    12.0ms ± 1%   -9.50%  (p=0.008 n=5+5)
RangeQuery/clamp-8                                       251ms ± 2%      29ms ± 1%  -88.41%  (p=0.008 n=5+5)
RangeQuery/clamp_min-8                                   250ms ± 5%      26ms ± 1%  -89.68%  (p=0.008 n=5+5)
RangeQuery/complex_func_query-8                          448ms ± 1%      35ms ± 2%  -92.22%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query-8                      255ms ± 1%      35ms ±22%  -86.17%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query-8                      261ms ± 2%      32ms ± 6%  -87.62%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile-8                          541ms ± 2%     100ms ± 0%  -81.45%  (p=0.008 n=5+5)

name                                                  old alloc/op   new alloc/op   delta
RangeQuery/vector_selector-8                            25.6MB ± 0%    38.5MB ± 0%  +50.22%  (p=0.008 n=5+5)
RangeQuery/sum-8                                        8.30MB ± 0%   11.30MB ± 1%  +36.15%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod-8                                 96.5MB ± 0%    23.4MB ± 0%  -75.73%  (p=0.008 n=5+5)
RangeQuery/topk-8                                       8.66MB ± 0%   11.17MB ± 1%  +29.03%  (p=0.008 n=5+5)
RangeQuery/bottomk-8                                    8.85MB ± 0%   11.37MB ± 0%  +28.48%  (p=0.016 n=5+4)
RangeQuery/rate-8                                       26.7MB ± 0%    41.2MB ± 0%  +54.27%  (p=0.008 n=5+5)
RangeQuery/sum_rate-8                                   9.36MB ± 0%   13.87MB ± 1%  +48.11%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate-8                                89.6MB ± 0%    25.1MB ± 0%  -71.93%  (p=0.008 n=5+5)
RangeQuery/quantile_with_variable_parameter-8            192MB ± 0%      42MB ± 0%  -78.08%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one-8           17.3MB ± 0%    18.4MB ± 0%   +6.22%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one-8          73.5MB ± 0%    45.4MB ± 0%  -38.23%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar-8    40.2MB ± 0%    41.5MB ± 0%   +3.21%  (p=0.008 n=5+5)
RangeQuery/unary_negation-8                             26.7MB ± 0%    39.9MB ± 0%  +49.77%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison-8               38.8MB ± 0%    41.1MB ± 0%   +5.97%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector-8                     23.8MB ± 0%    36.6MB ± 0%  +53.54%  (p=0.008 n=5+5)
RangeQuery/at_modifier_-8                               40.6MB ± 0%    33.1MB ± 0%  -18.42%  (p=0.008 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector-8    40.4MB ± 0%    32.9MB ± 0%  -18.50%  (p=0.008 n=5+5)
RangeQuery/clamp-8                                      40.3MB ± 0%    38.7MB ± 0%   -3.96%  (p=0.008 n=5+5)
RangeQuery/clamp_min-8                                  40.2MB ± 0%    38.7MB ± 0%   -3.89%  (p=0.008 n=5+5)
RangeQuery/complex_func_query-8                         54.8MB ± 0%    41.7MB ± 0%  -23.88%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query-8                     41.2MB ± 0%    41.3MB ± 0%   +0.26%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query-8                     41.2MB ± 0%    41.2MB ± 0%     ~     (p=0.095 n=5+5)
RangeQuery/histogram_quantile-8                         51.8MB ± 0%    59.7MB ± 0%  +15.33%  (p=0.008 n=5+5)

name                                                  old allocs/op  new allocs/op  delta
RangeQuery/vector_selector-8                              120k ± 0%      129k ± 0%   +7.78%  (p=0.016 n=4+5)
RangeQuery/sum-8                                          124k ± 0%      123k ± 0%   -0.65%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod-8                                   619k ± 0%      209k ± 0%  -66.26%  (p=0.008 n=5+5)
RangeQuery/topk-8                                         131k ± 0%      131k ± 0%     ~     (p=0.683 n=5+5)
RangeQuery/bottomk-8                                      135k ± 0%      127k ± 0%   -5.91%  (p=0.008 n=5+5)
RangeQuery/rate-8                                         129k ± 0%      153k ± 0%  +18.80%  (p=0.008 n=5+5)
RangeQuery/sum_rate-8                                     133k ± 0%      147k ± 0%  +10.54%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate-8                                  626k ± 0%      232k ± 0%  -62.88%  (p=0.008 n=5+5)
RangeQuery/quantile_with_variable_parameter-8            1.71M ± 0%     0.61M ± 0%  -64.07%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one-8            89.0k ± 0%    102.4k ± 0%  +15.02%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one-8            662k ± 0%      196k ± 0%  -70.34%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar-8      135k ± 0%      133k ± 0%   -1.25%  (p=0.008 n=5+5)
RangeQuery/unary_negation-8                               129k ± 0%      141k ± 0%   +9.89%  (p=0.016 n=4+5)
RangeQuery/vector_and_scalar_comparison-8                 126k ± 0%      130k ± 0%   +3.47%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector-8                      90.7k ± 0%     99.9k ± 0%  +10.07%  (p=0.008 n=5+5)
RangeQuery/at_modifier_-8                                91.9k ± 0%     74.7k ± 0%  -18.74%  (p=0.029 n=4+4)
RangeQuery/at_modifier_with_positive_offset_vector-8     85.9k ± 0%     68.7k ± 0%  -20.05%  (p=0.000 n=4+5)
RangeQuery/clamp-8                                        135k ± 0%      133k ± 0%   -1.83%  (p=0.008 n=5+5)
RangeQuery/clamp_min-8                                    135k ± 0%      132k ± 0%   -1.90%  (p=0.008 n=5+5)
RangeQuery/complex_func_query-8                           150k ± 0%      137k ± 0%   -8.95%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query-8                       144k ± 0%      154k ± 0%   +6.66%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query-8                       144k ± 0%      154k ± 0%   +6.63%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile-8                           694k ± 0%      702k ± 0%   +1.05%  (p=0.008 n=5+5)
```
