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
RangeQuery/vector_selector                            28.1ms ± 2%    34.2ms ± 4%  +21.80%  (p=0.000 n=9+10)
RangeQuery/sum                                        46.4ms ± 4%    27.5ms ± 0%  -40.59%  (p=0.000 n=10+10)
RangeQuery/sum_by_pod                                  134ms ± 2%      37ms ± 0%  -72.72%  (p=0.000 n=10+8)
RangeQuery/rate                                       52.3ms ± 3%    59.2ms ± 3%  +13.27%  (p=0.000 n=10+10)
RangeQuery/sum_rate                                   69.8ms ± 3%    53.9ms ± 4%  -22.80%  (p=0.000 n=10+9)
RangeQuery/sum_by_rate                                 157ms ± 4%      62ms ± 4%  -60.23%  (p=0.000 n=10+10)
RangeQuery/binary_operation_with_one_to_one           10.9ms ± 3%     2.4ms ± 4%  -77.58%  (p=0.000 n=10+10)
RangeQuery/binary_operation_with_many_to_one           389ms ± 4%      51ms ± 5%  -86.80%  (p=0.000 n=10+10)
RangeQuery/binary_operation_with_vector_and_scalar     241ms ± 3%      37ms ± 2%  -84.55%  (p=0.000 n=9+10)

name                                                old alloc/op   new alloc/op   delta
RangeQuery/vector_selector                            16.1MB ± 0%    25.7MB ± 0%  +59.23%  (p=0.000 n=10+10)
RangeQuery/sum                                        4.93MB ± 0%    6.93MB ± 0%  +40.64%  (p=0.000 n=8+8)
RangeQuery/sum_by_pod                                 87.3MB ± 0%    16.0MB ± 0%  -81.63%  (p=0.000 n=10+10)
RangeQuery/rate                                       17.2MB ± 0%    28.2MB ± 0%  +63.65%  (p=0.000 n=10+10)
RangeQuery/sum_rate                                   5.98MB ± 0%    9.38MB ± 0%  +56.89%  (p=0.000 n=10+9)
RangeQuery/sum_by_rate                                80.3MB ± 0%    18.4MB ± 0%  -77.11%  (p=0.000 n=10+10)
RangeQuery/binary_operation_with_one_to_one           1.64MB ± 0%    2.83MB ± 0%  +72.80%  (p=0.000 n=9+10)
RangeQuery/binary_operation_with_many_to_one          63.3MB ± 0%    31.4MB ± 0%  -50.39%  (p=0.000 n=10+9)
RangeQuery/binary_operation_with_vector_and_scalar    30.8MB ± 0%    28.6MB ± 0%   -7.15%  (p=0.000 n=10+10)

name                                                old allocs/op  new allocs/op  delta
RangeQuery/vector_selector                             69.1k ± 0%     76.7k ± 0%  +11.02%  (p=0.000 n=10+10)
RangeQuery/sum                                         73.8k ± 0%     71.4k ± 0%   -3.31%  (p=0.000 n=10+9)
RangeQuery/sum_by_pod                                   569k ± 0%      156k ± 0%  -72.49%  (p=0.000 n=10+8)
RangeQuery/rate                                        78.2k ± 0%    103.7k ± 0%  +32.71%  (p=0.000 n=9+10)
RangeQuery/sum_rate                                    82.9k ± 0%     98.4k ± 0%  +18.72%  (p=0.000 n=8+8)
RangeQuery/sum_by_rate                                  576k ± 0%      183k ± 0%  -68.13%  (p=0.000 n=10+10)
RangeQuery/binary_operation_with_one_to_one            25.7k ± 0%     20.7k ± 0%  -19.42%  (p=0.000 n=10+10)
RangeQuery/binary_operation_with_many_to_one            594k ± 0%      126k ± 0%  -78.81%  (p=0.000 n=10+10)
RangeQuery/binary_operation_with_vector_and_scalar     84.4k ± 0%     86.6k ± 0%   +2.68%  (p=0.000 n=10+10)
```

Multi-core (8 core) benchmarks

```markdown
name                                                  old time/op    new time/op    delta
RangeQuery/vector_selector-8                            27.3ms ± 3%    13.7ms ± 3%  -49.94%  (p=0.000 n=10+10)
RangeQuery/sum-8                                        47.8ms ± 4%     8.9ms ± 7%  -81.44%  (p=0.000 n=10+9)
RangeQuery/sum_by_pod-8                                  131ms ± 3%      14ms ± 1%  -89.44%  (p=0.000 n=10+10)
RangeQuery/rate-8                                       49.6ms ± 0%    19.7ms ± 1%  -60.25%  (p=0.000 n=10+9)
RangeQuery/sum_rate-8                                   70.3ms ± 3%    15.1ms ± 1%  -78.52%  (p=0.000 n=9+10)
RangeQuery/sum_by_rate-8                                 150ms ± 2%      20ms ± 0%  -86.53%  (p=0.000 n=10+9)
RangeQuery/binary_operation_with_one_to_one-8           10.4ms ± 4%     1.9ms ± 2%  -82.13%  (p=0.000 n=10+10)
RangeQuery/binary_operation_with_many_to_one-8           383ms ± 2%      25ms ± 2%  -93.50%  (p=0.000 n=10+10)
RangeQuery/binary_operation_with_vector_and_scalar-8     236ms ± 2%      16ms ± 4%  -93.20%  (p=0.000 n=10+10)

name                                                  old alloc/op   new alloc/op   delta
RangeQuery/vector_selector-8                            16.1MB ± 0%    25.8MB ± 0%  +59.52%  (p=0.000 n=10+10)
RangeQuery/sum-8                                        4.93MB ± 0%    7.66MB ± 4%  +55.44%  (p=0.000 n=10+10)
RangeQuery/sum_by_pod-8                                 87.3MB ± 0%    16.8MB ± 0%  -80.73%  (p=0.000 n=10+10)
RangeQuery/rate-8                                       17.2MB ± 0%    28.1MB ± 0%  +63.21%  (p=0.000 n=10+10)
RangeQuery/sum_rate-8                                   5.98MB ± 0%    9.95MB ± 0%  +66.24%  (p=0.000 n=10+10)
RangeQuery/sum_by_rate-8                                80.3MB ± 0%    18.0MB ± 0%  -77.54%  (p=0.000 n=10+10)
RangeQuery/binary_operation_with_one_to_one-8           1.64MB ± 0%    2.61MB ± 1%  +58.84%  (p=0.000 n=10+10)
RangeQuery/binary_operation_with_many_to_one-8          63.3MB ± 0%    32.1MB ± 0%  -49.27%  (p=0.000 n=10+10)
RangeQuery/binary_operation_with_vector_and_scalar-8    30.8MB ± 0%    29.8MB ± 0%   -3.23%  (p=0.000 n=10+10)

name                                                  old allocs/op  new allocs/op  delta
RangeQuery/vector_selector-8                             69.1k ± 0%     79.2k ± 0%  +14.54%  (p=0.000 n=10+10)
RangeQuery/sum-8                                         73.8k ± 0%     73.8k ± 0%   -0.11%  (p=0.014 n=10+10)
RangeQuery/sum_by_pod-8                                   569k ± 0%      159k ± 0%  -71.99%  (p=0.000 n=10+10)
RangeQuery/rate-8                                        78.2k ± 0%    106.1k ± 0%  +35.73%  (p=0.000 n=10+10)
RangeQuery/sum_rate-8                                    82.9k ± 0%    100.7k ± 0%  +21.49%  (p=0.000 n=7+10)
RangeQuery/sum_by_rate-8                                  576k ± 0%      186k ± 0%  -67.73%  (p=0.000 n=10+10)
RangeQuery/binary_operation_with_one_to_one-8            25.7k ± 0%     21.1k ± 0%  -17.71%  (p=0.000 n=10+10)
RangeQuery/binary_operation_with_many_to_one-8            594k ± 0%      131k ± 0%  -77.91%  (p=0.000 n=9+10)
RangeQuery/binary_operation_with_vector_and_scalar-8     84.4k ± 0%     89.5k ± 0%   +6.02%  (p=0.000 n=9+10)
```
