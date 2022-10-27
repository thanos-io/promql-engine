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
RangeQuery/vector_selector                            28.9ms ± 6%    35.0ms ± 1%  +21.12%  (p=0.008 n=5+5)
RangeQuery/sum                                        49.2ms ± 3%    29.0ms ± 2%  -40.96%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod                                  142ms ± 3%      39ms ± 1%  -72.81%  (p=0.008 n=5+5)
RangeQuery/rate                                       53.0ms ± 1%    58.3ms ± 1%  +10.12%  (p=0.008 n=5+5)
RangeQuery/sum_rate                                   72.1ms ± 1%    52.6ms ± 2%  -27.09%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate                                 158ms ± 2%      62ms ± 1%  -61.03%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one            120ms ± 2%      24ms ± 1%  -79.91%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one           388ms ± 0%      53ms ± 1%  -86.36%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar     237ms ± 1%      43ms ± 6%  -82.01%  (p=0.008 n=5+5)
RangeQuery/unary_negation                             28.6ms ± 1%    36.4ms ± 7%  +27.59%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison                192ms ± 2%      36ms ± 1%  -81.21%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector                     26.3ms ± 5%    31.7ms ± 1%  +20.57%  (p=0.008 n=5+5)
RangeQuery/at_modifier_                               11.1ms ± 0%    12.4ms ± 1%  +11.64%  (p=0.008 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector    10.6ms ± 1%    11.8ms ± 1%  +11.09%  (p=0.008 n=5+5)
RangeQuery/clamp                                       244ms ± 9%      50ms ± 4%  -79.62%  (p=0.008 n=5+5)
RangeQuery/clamp_min                                   238ms ± 1%      46ms ± 1%  -80.85%  (p=0.008 n=5+5)
RangeQuery/complex_func_query                          442ms ± 1%      54ms ± 1%  -87.79%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query                      250ms ± 1%      70ms ± 2%  -72.22%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query                      253ms ± 1%      71ms ± 2%  -71.86%  (p=0.008 n=5+5)

name                                                old alloc/op   new alloc/op   delta
RangeQuery/vector_selector                            16.5MB ± 0%    26.1MB ± 0%  +57.91%  (p=0.008 n=5+5)
RangeQuery/sum                                        5.32MB ± 0%    7.33MB ± 0%  +37.93%  (p=0.016 n=5+4)
RangeQuery/sum_by_pod                                 87.6MB ± 0%    16.4MB ± 0%  -81.25%  (p=0.008 n=5+5)
RangeQuery/rate                                       17.6MB ± 0%    27.7MB ± 0%  +57.42%  (p=0.008 n=5+5)
RangeQuery/sum_rate                                   6.37MB ± 0%    8.92MB ± 0%  +40.13%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate                                80.7MB ± 0%    17.9MB ± 0%  -77.79%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one           13.5MB ± 0%    13.2MB ± 0%   -2.73%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one          63.8MB ± 0%    31.9MB ± 0%  -49.95%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar    31.1MB ± 0%    29.0MB ± 0%   -6.75%  (p=0.008 n=5+5)
RangeQuery/unary_negation                             17.6MB ± 0%    27.2MB ± 0%  +54.56%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison               29.7MB ± 0%    27.6MB ± 0%   -7.17%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector                     15.7MB ± 0%    25.3MB ± 0%  +60.92%  (p=0.008 n=5+5)
RangeQuery/at_modifier_                               27.3MB ± 0%    34.2MB ± 0%  +25.05%  (p=0.008 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector    27.1MB ± 0%    34.0MB ± 0%  +25.23%  (p=0.016 n=4+5)
RangeQuery/clamp                                      31.1MB ± 0%    26.3MB ± 0%  -15.56%  (p=0.008 n=5+5)
RangeQuery/clamp_min                                  31.1MB ± 0%    26.2MB ± 0%  -15.72%  (p=0.008 n=5+5)
RangeQuery/complex_func_query                         45.6MB ± 0%    29.2MB ± 0%  -35.94%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query                     32.1MB ± 0%    27.7MB ± 0%  -13.62%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query                     32.1MB ± 0%    27.7MB ± 0%  -13.61%  (p=0.008 n=5+5)

name                                                old allocs/op  new allocs/op  delta
RangeQuery/vector_selector                             69.1k ± 0%     76.3k ± 0%  +10.40%  (p=0.008 n=5+5)
RangeQuery/sum                                         73.8k ± 0%     70.7k ± 0%   -4.21%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod                                   569k ± 0%      156k ± 0%  -72.60%  (p=0.008 n=5+5)
RangeQuery/rate                                        78.2k ± 0%     97.3k ± 0%  +24.51%  (p=0.008 n=5+5)
RangeQuery/sum_rate                                    82.9k ± 0%     91.7k ± 0%  +10.69%  (p=0.016 n=4+5)
RangeQuery/sum_by_rate                                  576k ± 0%      177k ± 0%  -69.29%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one            55.3k ± 0%     64.8k ± 0%  +17.17%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one            594k ± 0%      125k ± 0%  -79.00%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar     84.4k ± 0%     87.2k ± 0%   +3.37%  (p=0.008 n=5+5)
RangeQuery/unary_negation                              78.1k ± 0%     85.4k ± 0%   +9.26%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison                75.3k ± 0%     85.1k ± 0%  +12.99%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector                      57.1k ± 0%     64.3k ± 0%  +12.57%  (p=0.008 n=5+5)
RangeQuery/at_modifier_                                75.1k ± 0%     58.6k ± 0%  -21.97%  (p=0.000 n=5+4)
RangeQuery/at_modifier_with_positive_offset_vector     69.1k ± 0%     52.6k ± 0%  -23.88%  (p=0.000 n=5+4)
RangeQuery/clamp                                       84.7k ± 0%     78.4k ± 0%   -7.34%  (p=0.008 n=5+5)
RangeQuery/clamp_min                                   84.4k ± 0%     77.4k ± 0%   -8.29%  (p=0.008 n=5+5)
RangeQuery/complex_func_query                           100k ± 0%       89k ± 0%  -10.55%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query                      93.7k ± 0%     96.5k ± 0%   +2.95%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query                      93.7k ± 0%     96.5k ± 0%   +2.95%  (p=0.008 n=5+5)
```

Multi-core (8 core) benchmarks

```markdown
name                                                  old time/op    new time/op    delta
RangeQuery/vector_selector-8                            28.2ms ± 2%    13.9ms ± 1%  -50.60%  (p=0.008 n=5+5)
RangeQuery/sum-8                                        48.1ms ± 2%     8.9ms ± 2%  -81.50%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod-8                                  134ms ± 1%      14ms ± 0%  -89.53%  (p=0.008 n=5+5)
RangeQuery/rate-8                                       53.7ms ± 3%    19.8ms ± 3%  -63.05%  (p=0.008 n=5+5)
RangeQuery/sum_rate-8                                   72.3ms ± 1%    15.0ms ± 1%  -79.21%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate-8                                 156ms ± 2%      20ms ± 1%  -87.26%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one-8            125ms ± 1%      12ms ± 0%  -90.36%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one-8           391ms ± 4%      25ms ± 1%  -93.49%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar-8     234ms ± 4%      17ms ± 1%  -92.59%  (p=0.008 n=5+5)
RangeQuery/unary_negation-8                             28.9ms ± 5%    14.4ms ± 2%  -49.95%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison-8                194ms ± 5%      14ms ± 3%  -92.83%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector-8                     25.3ms ± 2%    13.0ms ± 0%  -48.62%  (p=0.008 n=5+5)
RangeQuery/at_modifier_-8                               9.66ms ± 2%   11.25ms ± 1%  +16.52%  (p=0.008 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector-8    9.17ms ± 2%   10.59ms ± 1%  +15.40%  (p=0.008 n=5+5)
RangeQuery/clamp-8                                       230ms ± 0%      27ms ± 1%  -88.49%  (p=0.016 n=5+4)
RangeQuery/clamp_min-8                                   231ms ± 1%      25ms ±21%  -89.10%  (p=0.008 n=5+5)
RangeQuery/complex_func_query-8                          439ms ± 1%      31ms ± 1%  -92.84%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query-8                      242ms ± 1%      32ms ±17%  -86.71%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query-8                      250ms ± 3%      27ms ± 3%  -89.10%  (p=0.008 n=5+5)

name                                                  old alloc/op   new alloc/op   delta
RangeQuery/vector_selector-8                            16.5MB ± 0%    26.1MB ± 0%  +58.16%  (p=0.008 n=5+5)
RangeQuery/sum-8                                        5.32MB ± 0%    7.98MB ± 1%  +50.22%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod-8                                 87.6MB ± 0%    17.2MB ± 0%  -80.36%  (p=0.008 n=5+5)
RangeQuery/rate-8                                       17.6MB ± 0%    27.7MB ± 0%  +57.26%  (p=0.016 n=4+5)
RangeQuery/sum_rate-8                                   6.37MB ± 0%    9.52MB ± 0%  +49.54%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate-8                                80.7MB ± 0%    17.6MB ± 0%  -78.22%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one-8           13.5MB ± 0%    13.6MB ± 0%   +0.47%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one-8          63.8MB ± 0%    32.7MB ± 0%  -48.78%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar-8    31.1MB ± 0%    30.2MB ± 0%   -3.05%  (p=0.008 n=5+5)
RangeQuery/unary_negation-8                             17.6MB ± 0%    27.4MB ± 0%  +55.85%  (p=0.008 n=5+5)
RangeQuery/vector_and_scalar_comparison-8               29.7MB ± 0%    27.9MB ± 0%   -6.01%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector-8                     15.7MB ± 0%    25.3MB ± 0%  +60.67%  (p=0.008 n=5+5)
RangeQuery/at_modifier_-8                               27.3MB ± 0%    34.2MB ± 0%  +25.17%  (p=0.008 n=5+5)
RangeQuery/at_modifier_with_positive_offset_vector-8    27.1MB ± 0%    34.0MB ± 0%  +25.38%  (p=0.008 n=5+5)
RangeQuery/clamp-8                                      31.1MB ± 0%    26.3MB ± 0%  -15.43%  (p=0.008 n=5+5)
RangeQuery/clamp_min-8                                  31.1MB ± 0%    26.3MB ± 0%  -15.51%  (p=0.008 n=5+5)
RangeQuery/complex_func_query-8                         45.6MB ± 0%    30.3MB ± 0%  -33.69%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query-8                     32.1MB ± 0%    27.8MB ± 0%  -13.46%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query-8                     32.1MB ± 0%    27.7MB ± 0%  -13.62%  (p=0.008 n=5+5)

name                                                  old allocs/op  new allocs/op  delta
RangeQuery/vector_selector-8                             69.1k ± 0%     78.0k ± 0%  +12.86%  (p=0.008 n=5+5)
RangeQuery/sum-8                                         73.8k ± 0%     72.3k ± 0%   -2.04%  (p=0.008 n=5+5)
RangeQuery/sum_by_pod-8                                   569k ± 0%      158k ± 0%  -72.24%  (p=0.008 n=5+5)
RangeQuery/rate-8                                        78.2k ± 0%     99.0k ± 0%  +26.65%  (p=0.016 n=4+5)
RangeQuery/sum_rate-8                                    82.9k ± 0%     93.3k ± 0%  +12.61%  (p=0.008 n=5+5)
RangeQuery/sum_by_rate-8                                  576k ± 0%      178k ± 0%  -69.01%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_one_to_one-8            55.3k ± 0%     68.6k ± 0%  +24.07%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_many_to_one-8            594k ± 0%      129k ± 0%  -78.35%  (p=0.008 n=5+5)
RangeQuery/binary_operation_with_vector_and_scalar-8     84.4k ± 0%     89.3k ± 0%   +5.82%  (p=0.008 n=5+5)
RangeQuery/unary_negation-8                              78.1k ± 0%     87.2k ± 0%  +11.61%  (p=0.016 n=4+5)
RangeQuery/vector_and_scalar_comparison-8                75.3k ± 0%     87.0k ± 0%  +15.50%  (p=0.008 n=5+5)
RangeQuery/positive_offset_vector-8                      57.1k ± 0%     66.0k ± 0%  +15.50%  (p=0.008 n=5+5)
RangeQuery/at_modifier_-8                                75.1k ± 0%     58.7k ± 0%  -21.85%  (p=0.000 n=5+4)
RangeQuery/at_modifier_with_positive_offset_vector-8     69.1k ± 0%     52.7k ± 0%  -23.75%  (p=0.016 n=4+5)
RangeQuery/clamp-8                                       84.7k ± 0%     80.2k ± 0%   -5.25%  (p=0.008 n=5+5)
RangeQuery/clamp_min-8                                   84.4k ± 0%     79.2k ± 0%   -6.17%  (p=0.008 n=5+5)
RangeQuery/complex_func_query-8                           100k ± 0%       91k ± 0%   -8.54%  (p=0.008 n=5+5)
RangeQuery/func_within_func_query-8                      93.7k ± 0%     98.3k ± 0%   +4.83%  (p=0.008 n=5+5)
RangeQuery/aggr_within_func_query-8                      93.7k ± 0%     98.3k ± 0%   +4.82%  (p=0.008 n=5+5)
```
