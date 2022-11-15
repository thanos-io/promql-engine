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
RangeQuery/vector_selector                            34.6ms ± 2%    45.0ms ± 2%  +29.94%  (p=0.008 n=5+5)
RangeQuery/sum                                        49.9ms ± 1%    34.6ms ± 2%  -30.60%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod                                  157ms ± 5%      46ms ± 2%  -70.53%  (p=0.008 n=5+5)
RangeQuery/rate                                       66.9ms ± 1%    76.3ms ± 3%  +14.05%  (p=0.008 n=5+5)
RangeQuery/sum_rate                                   78.8ms ± 2%    64.4ms ± 2%  -18.26%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate                                 178ms ± 0%      73ms ± 1%  -58.76%  (p=0.008 n=5+5)
RangeQuery/quantile_with_variable_parameter            275ms ± 1%     100ms ± 6%  -63.46%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one            123ms ± 3%      31ms ± 1%  -75.05%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one           402ms ± 2%      64ms ± 1%  -84.03%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar     242ms ± 1%      52ms ± 3%  -78.50%  (p=0.008 n=5+5)
RangeQuery/unary_negation                             37.4ms ± 1%    47.7ms ± 1%  +27.44%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison                196ms ± 1%      53ms ± 1%  -73.11%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector                     32.1ms ± 5%    41.3ms ± 2%  +28.72%  (p=0.008 n=5+5)
RangeQuery/at_modifier_                               18.9ms ± 2%    18.8ms ± 1%     ~     (p=0.222 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector    18.3ms ± 3%    17.9ms ± 1%     ~     (p=0.103 n=5+5)
RangeQuery/clamp                                       241ms ± 1%      60ms ± 2%  -75.05%  (p=0.008 n=5+5)
RangeQuery/clamp_min                                   239ms ± 0%      57ms ± 4%  -76.30%  (p=0.008 n=5+5)
RangeQuery/complex_func_query                          447ms ± 2%      64ms ± 1%  -85.63%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query                      256ms ± 2%      89ms ± 2%  -65.05%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query                      274ms ± 9%      91ms ± 1%  -66.63%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile                          561ms ± 2%     204ms ± 1%  -63.64%  (p=0.008 n=5+5)

name                                                old alloc/op   new alloc/op   delta
RangeQuery/vector_selector                            25.6MB ± 0%    36.9MB ± 0%  +43.67%  (p=0.016 n=4+5)
RangeQuery/sum                                        8.34MB ± 0%    8.88MB ± 0%   +6.50%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod                                 96.5MB ± 0%    21.0MB ± 0%  -78.21%  (p=0.008 n=5+5)
RangeQuery/rate                                       26.7MB ± 0%    39.7MB ± 0%  +48.67%  (p=0.008 n=5+5)
RangeQuery/sum_rate                                   9.40MB ± 0%   11.74MB ± 0%  +24.96%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate                                89.6MB ± 0%    23.8MB ± 0%  -73.44%  (p=0.008 n=5+5)
RangeQuery/quantile_with_variable_parameter            192MB ± 0%      36MB ± 0%  -81.03%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one           17.4MB ± 0%    17.0MB ± 0%   -2.10%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one          73.5MB ± 0%    42.7MB ± 0%  -41.91%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar    40.3MB ± 0%    39.0MB ± 0%   -3.18%  (p=0.008 n=5+5)
RangeQuery/unary_negation                             26.7MB ± 0%    38.0MB ± 0%  +42.09%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison               38.9MB ± 0%    38.6MB ± 0%   -0.55%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector                     23.8MB ± 0%    35.0MB ± 0%  +47.00%  (p=0.016 n=4+5)
RangeQuery/at_modifier_                               40.6MB ± 0%    44.4MB ± 0%   +9.47%  (p=0.016 n=4+5)
RangeQuery/at_modifier_with_positive_offset_vector    40.4MB ± 0%    44.2MB ± 0%   +9.51%  (p=0.016 n=4+5)
RangeQuery/clamp                                      40.3MB ± 0%    37.3MB ± 0%   -7.52%  (p=0.008 n=5+5)
RangeQuery/clamp_min                                  40.3MB ± 0%    37.2MB ± 0%   -7.56%  (p=0.008 n=5+5)
RangeQuery/complex_func_query                         54.8MB ± 0%    39.3MB ± 0%  -28.22%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query                     41.2MB ± 0%    39.8MB ± 0%   -3.50%  (p=0.016 n=5+4)
RangeQuery/aggr_within_func_query                     41.2MB ± 0%    39.8MB ± 0%   -3.50%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile                         51.6MB ± 0%    53.3MB ± 0%   +3.27%  (p=0.008 n=5+5)

name                                                old allocs/op  new allocs/op  delta
RangeQuery/vector_selector                              120k ± 0%      127k ± 0%   +5.51%  (p=0.008 n=5+5)
RangeQuery/sum                                          125k ± 0%      122k ± 0%   -2.75%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod                                   620k ± 0%      207k ± 0%  -66.66%  (p=0.008 n=5+5)
RangeQuery/rate                                         129k ± 0%      151k ± 0%  +16.69%  (p=0.008 n=5+5)
RangeQuery/sum_rate                                     134k ± 0%      146k ± 0%   +8.60%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate                                  627k ± 0%      231k ± 0%  -63.20%  (p=0.008 n=5+5)
RangeQuery/quantile_with_variable_parameter            1.71M ± 0%     0.61M ± 0%  -64.37%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one            89.1k ± 0%     98.3k ± 0%  +10.40%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one            662k ± 0%      192k ± 0%  -71.00%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar      136k ± 0%      131k ± 0%   -3.45%  (p=0.008 n=5+5)
RangeQuery/unary_negation                               129k ± 0%      136k ± 0%   +5.35%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison                 126k ± 0%      128k ± 0%   +1.10%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector                      91.3k ± 0%     98.0k ± 0%   +7.39%  (p=0.008 n=5+5)
RangeQuery/at_modifier_                                92.2k ± 0%     75.6k ± 0%  -18.05%  (p=0.029 n=4+4)
RangeQuery/at_modifier_with_positive_offset_vector     86.2k ± 0%     69.6k ± 0%  -19.31%  (p=0.000 n=5+4)
RangeQuery/clamp                                        136k ± 0%      131k ± 0%   -3.83%  (p=0.008 n=5+5)
RangeQuery/clamp_min                                    136k ± 0%      130k ± 0%   -3.89%  (p=0.008 n=5+5)
RangeQuery/complex_func_query                           151k ± 0%      135k ± 0%  -10.45%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query                       145k ± 0%      152k ± 0%   +5.16%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query                       145k ± 0%      152k ± 0%   +5.16%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile                           693k ± 0%      700k ± 0%   +0.98%  (p=0.008 n=5+5)
```

Multi-core (8 core) benchmarks

```markdown
name                                                  old time/op    new time/op    delta
RangeQuery/vector_selector-8                            32.6ms ± 2%    17.7ms ± 1%  -45.86%  (p=0.008 n=5+5)
RangeQuery/sum-8                                        51.5ms ± 1%    11.3ms ± 1%  -78.15%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod-8                                  147ms ± 1%      18ms ± 2%  -88.00%  (p=0.008 n=5+5)
RangeQuery/rate-8                                       64.9ms ± 1%    25.2ms ± 0%  -61.19%  (p=0.008 n=5+5)
RangeQuery/sum_rate-8                                   82.9ms ± 1%    19.1ms ± 2%  -76.94%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate-8                                 176ms ± 1%      25ms ± 2%  -85.66%  (p=0.008 n=5+5)
RangeQuery/quantile_with_variable_parameter-8            253ms ± 2%      29ms ± 1%  -88.44%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one-8            125ms ± 2%      14ms ± 2%  -88.58%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one-8           415ms ± 1%      30ms ± 2%  -92.80%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar-8     250ms ± 1%      21ms ± 0%  -91.54%  (p=0.008 n=5+5)
RangeQuery/unary_negation-8                             33.5ms ± 1%    18.8ms ± 4%  -43.85%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison-8                202ms ± 1%      23ms ± 9%  -88.52%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector-8                     29.6ms ± 2%    16.7ms ± 5%  -43.40%  (p=0.008 n=5+5)
RangeQuery/at_modifier_-8                               13.6ms ± 1%    14.1ms ± 3%   +3.90%  (p=0.016 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector-8    13.0ms ± 1%    13.1ms ± 1%     ~     (p=0.151 n=5+5)
RangeQuery/clamp-8                                       249ms ± 1%      30ms ± 1%  -88.07%  (p=0.008 n=5+5)
RangeQuery/clamp_min-8                                   248ms ± 1%      26ms ± 2%  -89.70%  (p=0.008 n=5+5)
RangeQuery/complex_func_query-8                          459ms ± 0%      35ms ± 1%  -92.44%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query-8                      264ms ± 1%      31ms ± 3%  -88.19%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query-8                      268ms ± 1%      31ms ± 2%  -88.58%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile-8                          588ms ± 1%     100ms ± 1%  -82.95%  (p=0.008 n=5+5)

name                                                  old alloc/op   new alloc/op   delta
RangeQuery/vector_selector-8                            25.6MB ± 0%    38.3MB ± 0%  +49.37%  (p=0.008 n=5+5)
RangeQuery/sum-8                                        8.33MB ± 0%   11.21MB ± 0%  +34.69%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod-8                                 96.5MB ± 0%    23.4MB ± 0%  -75.76%  (p=0.008 n=5+5)
RangeQuery/rate-8                                       26.7MB ± 0%    41.1MB ± 0%  +53.96%  (p=0.008 n=5+5)
RangeQuery/sum_rate-8                                   9.38MB ± 0%   13.80MB ± 0%  +47.09%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate-8                                89.6MB ± 0%    24.9MB ± 0%  -72.16%  (p=0.008 n=5+5)
RangeQuery/quantile_with_variable_parameter-8            192MB ± 0%      42MB ± 0%  -78.10%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one-8           17.4MB ± 0%    18.4MB ± 0%   +5.93%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one-8          73.5MB ± 0%    45.4MB ± 0%  -38.26%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar-8    40.3MB ± 0%    41.6MB ± 0%   +3.35%  (p=0.008 n=5+5)
RangeQuery/unary_negation-8                             26.7MB ± 0%    39.5MB ± 0%  +48.12%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison-8               38.8MB ± 0%    41.2MB ± 0%   +6.16%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector-8                     23.8MB ± 0%    36.5MB ± 0%  +53.07%  (p=0.008 n=5+5)
RangeQuery/at_modifier_-8                               40.6MB ± 0%    44.4MB ± 0%   +9.55%  (p=0.008 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector-8    40.4MB ± 0%    44.3MB ± 0%   +9.61%  (p=0.008 n=5+5)
RangeQuery/clamp-8                                      40.3MB ± 0%    38.7MB ± 0%   -3.81%  (p=0.008 n=5+5)
RangeQuery/clamp_min-8                                  40.3MB ± 0%    38.7MB ± 0%   -3.93%  (p=0.008 n=5+5)
RangeQuery/complex_func_query-8                         54.8MB ± 0%    41.8MB ± 0%  -23.67%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query-8                     41.2MB ± 0%    41.3MB ± 0%   +0.19%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query-8                     41.2MB ± 0%    41.2MB ± 0%     ~     (p=0.151 n=5+5)
RangeQuery/histogram_quantile-8                         51.7MB ± 0%    59.7MB ± 0%  +15.57%  (p=0.008 n=5+5)

name                                                  old allocs/op  new allocs/op  delta
RangeQuery/vector_selector-8                              120k ± 0%      129k ± 0%   +7.45%  (p=0.008 n=5+5)
RangeQuery/sum-8                                          124k ± 0%      123k ± 0%   -0.90%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod-8                                   619k ± 0%      209k ± 0%  -66.29%  (p=0.008 n=5+5)
RangeQuery/rate-8                                         129k ± 0%      153k ± 0%  +18.53%  (p=0.008 n=5+5)
RangeQuery/sum_rate-8                                     133k ± 0%      147k ± 0%  +10.30%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate-8                                  626k ± 0%      232k ± 0%  -62.91%  (p=0.008 n=5+5)
RangeQuery/quantile_with_variable_parameter-8            1.71M ± 0%     0.61M ± 0%  -64.09%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one-8            89.5k ± 0%    102.3k ± 0%  +14.31%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one-8            662k ± 0%      196k ± 0%  -70.36%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar-8      135k ± 0%      133k ± 0%   -1.43%  (p=0.008 n=5+5)
RangeQuery/unary_negation-8                               129k ± 0%      138k ± 0%   +7.24%  (p=0.016 n=4+5)
RangeQuery/vector_and_scalar_comparison-8                 126k ± 0%      130k ± 0%   +3.27%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector-8                      90.9k ± 0%     99.7k ± 0%   +9.77%  (p=0.008 n=5+5)
RangeQuery/at_modifier_-8                                92.0k ± 0%     75.6k ± 0%  -17.82%  (p=0.000 n=5+4)
RangeQuery/at_modifier_with_positive_offset_vector-8     86.0k ± 0%     69.6k ± 0%  -19.07%  (p=0.000 n=5+4)
RangeQuery/clamp-8                                        135k ± 0%      133k ± 0%   -2.01%  (p=0.008 n=5+5)
RangeQuery/clamp_min-8                                    135k ± 0%      132k ± 0%   -2.10%  (p=0.008 n=5+5)
RangeQuery/complex_func_query-8                           151k ± 0%      137k ± 0%   -8.67%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query-8                       144k ± 0%      154k ± 0%   +6.89%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query-8                       144k ± 0%      154k ± 0%   +6.88%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile-8                           693k ± 0%      702k ± 0%   +1.24%  (p=0.008 n=5+5)
```
