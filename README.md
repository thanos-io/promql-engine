# PromQL Query Engine

A multi-threaded implementation of a PromQL Query Engine based on the [Volcano/Iterator model](https://paperhub.s3.amazonaws.com/dace52a42c07f7f8348b08dc2b186061.pdf).

The project is currently under active development.

## Roadmap

The engine intends to have full compatibility with the original engine used in Prometheus. Since implementing the full specification will take time, we aim to add support for most commonly used expressions while falling back to the original engine for operations that are not yet supported. This will allow us to have smaller and faster releases, and gather feedback on regular basis. Instructions on using the engine will be added after we have enough confidence on its correctness.

The following table shows operations which are currently supported by the engine

| Type                   | Supported                                                                                | Priority |
|------------------------|------------------------------------------------------------------------------------------|----------|
| Rate                   | Full support                                                                             |          |
| Binary expressions     | Full support                                                                             |          |
| Aggregations           | Partial support (sum, max, min, avg, count and group)                                    | Medium   |
| Aggregations over time | Partial support (sum, max, min, avg, count, stddev, stdvar, last and present) _over_time | Medium   |
| Functions              | No support                                                                               | Medium   |
| Quantiles              | No support                                                                               | High     |

In addition to implementing multi-threading, we would ultimately like to end up with a distributed execution model.

## Design

At the beginning of a PromQL query execution, the query engine computes a physical plan consisting of multiple independent operators, each responsible for calculating one part of the query expression.

Operators are assembled in a tree-like structure with every operator calling `Next()` on its dependants until there is no more data to be returned. The result of the `Next()` function is a *column vector* (also called a *step vector*) with elements in the vector representing samples with the same timestamp from different time series.

<p align="center">
  <img src="./assets/design.png"/>
</p>

This model allows for samples from individual time series to flow one execution step at a time from the left-most operators to the one at the very right. Since most PromQL expressions are aggregations, samples are reduced in number as they are pulled by the operators on the right. Because of this, samples from original timeseries can be decoded and kept in memory in batches instead of being fully expanded.

In addition to operators that have a one-to-one mapping with PromQL constructs, the Volcano model also describes so-called Exchange operators which can be used for flow control and optimizations, such as concurrency or batched selects. An example of an *Exchange* operator is described in the [Intra-operator parallelism](#intra-operator-parallelism) section.

### Inter-operator parallelism

Since operators are independent and rely on a common interface for pulling data, they can be run in parallel to each other. As soon as one operator has processed data from an evaluation step, it can pass the result onward so that its upstream can immediately start working on it.

<p align="center">
  <img src="./assets/promql-pipeline.png"/>
</p>

### Intra-operator parallelism

Parallelism can also be added within individual operators using a parallel coalesce exchange operator. Such exchange operators are indistinguishable from regular operators to their upstreams since they respect the same `Next()` interface.

<p align="center">
  <img src="./assets/parallel-coalesce.png"/>
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
RangeQuery/vector_selector                            29.6ms ± 1%    39.0ms ± 2%  +31.65%  (p=0.008 n=5+5)
RangeQuery/sum                                        50.2ms ± 1%    33.0ms ± 3%  -34.25%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod                                  141ms ± 2%      44ms ± 2%  -68.42%  (p=0.008 n=5+5)
RangeQuery/rate                                       56.3ms ± 1%    63.1ms ± 3%  +12.11%  (p=0.008 n=5+5)
RangeQuery/sum_rate                                   77.3ms ± 1%    58.3ms ± 3%  -24.58%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate                                 164ms ± 1%      68ms ± 3%  -58.63%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one            127ms ± 2%      27ms ± 5%  -78.56%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one           400ms ± 2%      56ms ± 1%  -85.87%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar     248ms ± 1%      45ms ± 4%  -81.94%  (p=0.008 n=5+5)
RangeQuery/unary_negation                             33.4ms ± 6%    41.3ms ± 4%  +23.52%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison                201ms ± 2%      45ms ± 2%  -77.74%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector                     28.3ms ± 3%    36.5ms ± 2%  +29.10%  (p=0.008 n=5+5)
RangeQuery/at_modifier_                               12.8ms ± 5%    14.3ms ± 4%  +11.70%  (p=0.008 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector    12.2ms ± 4%    14.1ms ± 4%  +15.96%  (p=0.008 n=5+5)
RangeQuery/clamp                                       240ms ± 1%      56ms ± 3%  -76.74%  (p=0.008 n=5+5)
RangeQuery/clamp_min                                   246ms ± 3%      52ms ± 6%  -78.89%  (p=0.008 n=5+5)
RangeQuery/complex_func_query                          451ms ± 2%      61ms ± 1%  -86.48%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query                      256ms ± 3%      73ms ± 1%  -71.29%  (p=0.016 n=5+4)
RangeQuery/aggr_within_func_query                      257ms ± 2%      79ms ± 4%  -69.40%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile                          545ms ± 2%     210ms ± 3%  -61.54%  (p=0.008 n=5+5)

name                                                old alloc/op   new alloc/op   delta
RangeQuery/vector_selector                            19.5MB ± 0%    29.1MB ± 0%  +49.19%  (p=0.008 n=5+5)
RangeQuery/sum                                        8.25MB ± 0%   10.28MB ± 0%  +24.59%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod                                 90.6MB ± 0%    19.4MB ± 0%  -78.60%  (p=0.008 n=5+5)
RangeQuery/rate                                       20.5MB ± 0%    30.6MB ± 0%  +49.22%  (p=0.008 n=5+5)
RangeQuery/sum_rate                                   9.31MB ± 0%   11.87MB ± 0%  +27.51%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate                                83.6MB ± 0%    20.9MB ± 0%  -75.04%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one           15.2MB ± 0%    14.9MB ± 0%   -2.43%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one          67.2MB ± 0%    35.4MB ± 0%  -47.37%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar    34.1MB ± 0%    30.9MB ± 0%   -9.22%  (p=0.008 n=5+5)
RangeQuery/unary_negation                             20.5MB ± 0%    30.1MB ± 0%  +46.77%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison               32.7MB ± 0%    30.9MB ± 0%   -5.48%  (p=0.016 n=5+4)
RangeQuery/positive_offset_vector                     17.7MB ± 0%    27.2MB ± 0%  +54.20%  (p=0.008 n=5+5)
RangeQuery/at_modifier_                               28.3MB ± 0%    35.1MB ± 0%  +24.20%  (p=0.016 n=4+5)
RangeQuery/at_modifier_with_positive_offset_vector    28.1MB ± 0%    35.0MB ± 0%  +24.36%  (p=0.008 n=5+5)
RangeQuery/clamp                                      34.1MB ± 0%    29.2MB ± 0%  -14.36%  (p=0.008 n=5+5)
RangeQuery/clamp_min                                  34.1MB ± 0%    29.1MB ± 0%  -14.41%  (p=0.008 n=5+5)
RangeQuery/complex_func_query                         48.5MB ± 0%    31.1MB ± 0%  -35.99%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query                     35.0MB ± 0%    30.6MB ± 0%  -12.51%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query                     35.0MB ± 0%    30.6MB ± 0%  -12.51%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile                         49.4MB ± 0%    96.1MB ± 0%  +94.43%  (p=0.008 n=5+5)

name                                                old allocs/op  new allocs/op  delta
RangeQuery/vector_selector                              120k ± 0%      127k ± 0%   +6.08%  (p=0.008 n=5+5)
RangeQuery/sum                                          124k ± 0%      121k ± 0%   -2.43%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod                                   619k ± 0%      206k ± 0%  -66.66%  (p=0.008 n=5+5)
RangeQuery/rate                                         129k ± 0%      148k ± 0%  +14.95%  (p=0.008 n=5+5)
RangeQuery/sum_rate                                     133k ± 0%      142k ± 0%   +6.71%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate                                  626k ± 0%      227k ± 0%  -63.67%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one            89.0k ± 0%     98.4k ± 0%  +10.64%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one            662k ± 0%      192k ± 0%  -70.94%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar      135k ± 0%      128k ± 0%   -5.21%  (p=0.008 n=5+5)
RangeQuery/unary_negation                               129k ± 0%      136k ± 0%   +5.69%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison                 126k ± 0%      128k ± 0%   +1.62%  (p=0.016 n=5+4)
RangeQuery/positive_offset_vector                      90.8k ± 0%     98.1k ± 0%   +7.97%  (p=0.008 n=5+5)
RangeQuery/at_modifier_                                92.0k ± 0%     75.5k ± 0%  -17.91%  (p=0.008 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector     86.0k ± 0%     69.5k ± 0%  -19.17%  (p=0.008 n=5+5)
RangeQuery/clamp                                        135k ± 0%      128k ± 0%   -5.59%  (p=0.008 n=5+5)
RangeQuery/clamp_min                                    135k ± 0%      127k ± 0%   -5.65%  (p=0.016 n=5+4)
RangeQuery/complex_func_query                           151k ± 0%      129k ± 0%  -14.05%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query                       144k ± 0%      146k ± 0%   +1.48%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query                       144k ± 0%      146k ± 0%   +1.48%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile                           694k ± 0%      690k ± 0%   -0.48%  (p=0.008 n=5+5)
```

Multi-core (8 core) benchmarks

```markdown
name                                                  old time/op    new time/op    delta
RangeQuery/vector_selector-8                            30.2ms ± 3%    16.1ms ± 2%  -46.68%  (p=0.008 n=5+5)
RangeQuery/sum-8                                        51.7ms ± 4%    11.1ms ± 2%  -78.59%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod-8                                  138ms ± 2%      16ms ± 0%  -88.21%  (p=0.008 n=5+5)
RangeQuery/rate-8                                       54.3ms ± 1%    22.9ms ±15%  -57.83%  (p=0.008 n=5+5)
RangeQuery/sum_rate-8                                   74.5ms ± 1%    18.1ms ±17%  -75.65%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate-8                                 155ms ± 1%      23ms ± 4%  -84.88%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one-8            120ms ± 0%      15ms ±22%  -87.26%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one-8           387ms ± 2%      28ms ± 0%  -92.72%  (p=0.016 n=5+4)
RangeQuery/binary_operation_with_vector_and_scalar-8     248ms ± 2%      19ms ± 0%  -92.21%  (p=0.016 n=5+4)
RangeQuery/unary_negation-8                             31.7ms ± 3%    18.3ms ±11%  -42.39%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison-8                203ms ± 3%      20ms ± 6%  -90.10%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector-8                     28.0ms ± 2%    14.9ms ± 0%  -46.97%  (p=0.008 n=5+5)
RangeQuery/at_modifier_-8                               11.5ms ± 3%    12.6ms ± 2%   +9.13%  (p=0.008 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector-8    10.8ms ± 2%    12.0ms ± 1%  +10.78%  (p=0.008 n=5+5)
RangeQuery/clamp-8                                       246ms ± 4%      28ms ± 3%  -88.60%  (p=0.008 n=5+5)
RangeQuery/clamp_min-8                                   242ms ± 4%      25ms ± 4%  -89.72%  (p=0.008 n=5+5)
RangeQuery/complex_func_query-8                          455ms ± 3%      33ms ± 4%  -92.75%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query-8                      248ms ± 0%      29ms ± 5%  -88.20%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query-8                      253ms ± 1%      29ms ± 3%  -88.36%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile-8                          538ms ± 2%     103ms ± 3%  -80.80%  (p=0.008 n=5+5)

name                                                  old alloc/op   new alloc/op   delta
RangeQuery/vector_selector-8                            19.5MB ± 0%    29.1MB ± 0%  +49.48%  (p=0.008 n=5+5)
RangeQuery/sum-8                                        8.27MB ± 0%   11.38MB ± 1%  +37.67%  (p=0.016 n=4+5)
RangeQuery/sum_by_pod-8                                 90.6MB ± 0%    20.3MB ± 0%  -77.64%  (p=0.008 n=5+5)
RangeQuery/rate-8                                       20.6MB ± 0%    30.6MB ± 0%  +49.02%  (p=0.008 n=5+5)
RangeQuery/sum_rate-8                                   9.32MB ± 0%   12.60MB ± 2%  +35.17%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate-8                                83.6MB ± 0%    20.7MB ± 1%  -75.26%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one-8           15.2MB ± 0%    15.3MB ± 0%   +0.54%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one-8          67.2MB ± 0%    36.1MB ± 0%  -46.27%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar-8    34.1MB ± 0%    32.1MB ± 0%   -5.91%  (p=0.008 n=5+5)
RangeQuery/unary_negation-8                             20.5MB ± 0%    30.4MB ± 0%  +48.12%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison-8               32.7MB ± 0%    32.0MB ± 0%   -2.07%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector-8                     17.7MB ± 0%    27.2MB ± 0%  +53.91%  (p=0.008 n=5+5)
RangeQuery/at_modifier_-8                               28.3MB ± 0%    35.2MB ± 0%  +24.31%  (p=0.008 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector-8    28.1MB ± 0%    35.0MB ± 0%  +24.49%  (p=0.008 n=5+5)
RangeQuery/clamp-8                                      34.1MB ± 0%    29.2MB ± 0%  -14.44%  (p=0.008 n=5+5)
RangeQuery/clamp_min-8                                  34.1MB ± 0%    29.2MB ± 0%  -14.20%  (p=0.008 n=5+5)
RangeQuery/complex_func_query-8                         48.6MB ± 0%    32.1MB ± 0%  -33.90%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query-8                     35.0MB ± 0%    30.7MB ± 0%  -12.38%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query-8                     35.0MB ± 0%    30.7MB ± 0%  -12.39%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile-8                         49.4MB ± 0%    97.5MB ± 0%  +97.15%  (p=0.008 n=5+5)

name                                                  old allocs/op  new allocs/op  delta
RangeQuery/vector_selector-8                              120k ± 0%      129k ± 0%   +7.66%  (p=0.008 n=5+5)
RangeQuery/sum-8                                          125k ± 0%      123k ± 0%   -0.83%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod-8                                   619k ± 0%      209k ± 0%  -66.28%  (p=0.008 n=5+5)
RangeQuery/rate-8                                         129k ± 0%      150k ± 0%  +16.39%  (p=0.016 n=4+5)
RangeQuery/sum_rate-8                                     134k ± 0%      144k ± 0%   +8.08%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate-8                                  626k ± 0%      229k ± 0%  -63.37%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one-8            89.0k ± 0%    102.4k ± 0%  +15.05%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one-8            662k ± 0%      197k ± 0%  -70.31%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar-8      135k ± 0%      130k ± 0%   -3.49%  (p=0.008 n=5+5)
RangeQuery/unary_negation-8                               129k ± 0%      138k ± 0%   +7.30%  (p=0.016 n=4+5)
RangeQuery/vector_and_scalar_comparison-8                 126k ± 0%      130k ± 0%   +3.45%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector-8                      90.9k ± 0%    100.0k ± 0%   +9.95%  (p=0.008 n=5+5)
RangeQuery/at_modifier_-8                                92.0k ± 0%     75.7k ± 0%  -17.74%  (p=0.008 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector-8     86.0k ± 0%     69.7k ± 0%  -18.98%  (p=0.016 n=4+5)
RangeQuery/clamp-8                                        135k ± 0%      130k ± 0%   -4.06%  (p=0.008 n=5+5)
RangeQuery/clamp_min-8                                    135k ± 0%      130k ± 0%   -4.12%  (p=0.008 n=5+5)
RangeQuery/complex_func_query-8                           151k ± 0%      132k ± 0%  -12.49%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query-8                       144k ± 0%      149k ± 0%   +2.90%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query-8                       144k ± 0%      149k ± 0%   +2.90%  (p=0.008 n=5+5)
RangeQuery/histogram_quantile-8                           693k ± 0%      692k ± 0%   -0.21%  (p=0.008 n=5+5)
```
