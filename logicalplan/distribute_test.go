// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"math"
	"regexp"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"

	"github.com/thanos-io/promql-engine/api"
	"github.com/thanos-io/promql-engine/query"
)

var replacements = map[string]*regexp.Regexp{
	" ": spaces,
	"(": openParenthesis,
	")": closedParenthesis,
}

func TestDistributedExecution(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name              string
		expr              string
		skipBinopPushdown bool
		expectWarn        bool
		expected          string
	}{
		{
			name:     "selector",
			expr:     `http_requests_total`,
			expected: `dedup(remote(http_requests_total), remote(http_requests_total))`,
		},
		{
			name:     "parentheses",
			expr:     `(http_requests_total)`,
			expected: `dedup(remote((http_requests_total)), remote((http_requests_total)))`,
		},
		{
			name:     "scalar",
			expr:     `scalar(redis::shard_price_per_month)`,
			expected: `scalar(dedup(remote(redis::shard_price_per_month), remote(redis::shard_price_per_month)))`,
		},
		{
			name:     "rate",
			expr:     `rate(http_requests_total[5m])`,
			expected: `dedup(remote(rate(http_requests_total[5m])), remote(rate(http_requests_total[5m])))`,
		},
		{
			name: "sum-rate",
			expr: `sum by (pod) (rate(http_requests_total[5m]))`,
			expected: `
sum by (pod) (dedup(
  remote(sum by (pod, region) (rate(http_requests_total[5m]))), 
  remote(sum by (pod, region) (rate(http_requests_total[5m])))))`,
		},
		{
			name: "sum-rate without labels preserves engine labels",
			expr: `sum without (pod, region) (rate(http_requests_total[5m]))`,
			expected: `
sum without (pod, region) (
  dedup(
    remote(sum without (pod) (rate(http_requests_total[5m]))),
    remote(sum without (pod) (rate(http_requests_total[5m])))
  )
)`,
		},
		{
			name: "avg",
			expr: `avg(http_requests_total)`,
			expected: `
sum(
  dedup(
    remote(sum by (region) (http_requests_total)),
    remote(sum by (region) (http_requests_total))
  )
) / on ()
sum(
  dedup(
    remote(count by (region) (http_requests_total)),
    remote(count by (region) (http_requests_total))
  )
)`,
		},
		{
			name: "avg with by-grouping",
			expr: `avg by (pod) (http_requests_total)`,
			expected: `
sum by (pod) (
  dedup(
    remote(sum by (pod, region) (http_requests_total)),
    remote(sum by (pod, region) (http_requests_total))
  )
) / on (pod)
sum by (pod) (
  dedup(
    remote(count by (pod, region) (http_requests_total)),
    remote(count by (pod, region) (http_requests_total))
  )
)`,
		},
		{
			name: "avg with without-grouping",
			expr: `avg without (pod) (http_requests_total)`,
			expected: `
sum without (pod) (
  dedup(
    remote(sum without (pod) (http_requests_total)),
    remote(sum without (pod) (http_requests_total))
  )
) / ignoring (pod)
sum without (pod) (
  dedup(
    remote(count without (pod) (http_requests_total)),
    remote(count without (pod) (http_requests_total))
  )
)`,
		},
		{
			name: "avg with prior aggregation",
			expr: `avg by (pod) (sum by (pod) (http_requests_total))`,
			expected: `
avg by (pod) (
  sum by (pod) (
    dedup(
        remote(sum by (pod, region) (http_requests_total)),
        remote(sum by (pod, region) (http_requests_total))
    )
  )
)`,
		},
		{
			name: "avg with prior binary expression",
			expr: `avg by (pod) (metric_a / metric_b)`,
			expected: `
sum by (pod) (
  dedup(
    remote(sum by (pod, region) (metric_a / metric_b)),
    remote(sum by (pod, region) (metric_a / metric_b))
  )
)
/ on (pod)
sum by (pod) (
  dedup(
    remote(count by (pod, region) (metric_a / metric_b)),
    remote(count by (pod, region) (metric_a / metric_b))
  )
)`,
		},
		{
			name: "two-level aggregation",
			expr: `max by (pod) (sum by (pod) (http_requests_total))`,
			expected: `
max by (pod) (
  sum by (pod) ( 
    dedup(
      remote(sum by (pod, region) (http_requests_total)),
      remote(sum by (pod, region) (http_requests_total))
    )
  )
)`,
		},
		{
			name: "aggregation of binary expression",
			expr: `max by (pod) (metric_a / metric_b)`,
			expected: `
max by (pod) (
  dedup(
    remote(max by (pod, region) (metric_a / metric_b)),
    remote(max by (pod, region) (metric_a / metric_b))
  )
)
`,
		},
		{
			name: "unsupported aggregation in the operand path",
			expr: `max by (pod) (sort(quantile(0.9, http_requests_total)))`,
			expected: `
max by (pod) (quantile(0.9,
  dedup(
    remote(http_requests_total),
    remote(http_requests_total)
  )
))`,
		},
		{
			name: "label replace",
			expr: `label_replace(http_requests_total, "pod", "$1", "instance", "(.*)")`,
			expected: `
dedup(
  remote(label_replace(http_requests_total, "pod", "$1", "instance", "(.*)")), 
  remote(label_replace(http_requests_total, "pod", "$1", "instance", "(.*)"))
)`,
		},
		{
			name: "label replace to internal label before an aggregation",
			expr: `max by (instance) (label_replace(http_requests_total, "pod", "$1", "instance", "(.*)"))`,
			expected: `
max by (instance) (
  dedup(
    remote(max by (instance, region) (label_replace(http_requests_total, "pod", "$1", "instance", "(.*)"))), 
    remote(max by (instance, region) (label_replace(http_requests_total, "pod", "$1", "instance", "(.*)")))
  )
)`,
		},
		{
			name: "label replace to internal label before an aggregation",
			expr: `max by (location) (label_replace(http_requests_total, "zone", "$1", "location", "(.*)"))`,
			expected: `
max by (location) (dedup(
  remote(max by (location, region) (label_replace(http_requests_total, "zone", "$1", "location", "(.*)"))), 
  remote(max by (location, region) (label_replace(http_requests_total, "zone", "$1", "location", "(.*)")))
))`,
		},
		{
			name:       "label replace to external label before an aggregation",
			expr:       `max by (location) (label_replace(http_requests_total, "region", "$1", "location", "(.*)"))`,
			expected:   `max by (location) (label_replace(dedup(remote(http_requests_total), remote(http_requests_total)), "region", "$1", "location", "(.*)"))`,
			expectWarn: true,
		},
		{
			name:       "label replace to external label before an avg",
			expr:       `avg by (location) (label_replace(http_requests_total, "region", "$1", "location", "(.*)"))`,
			expected:   `avg by (location) (label_replace(dedup(remote(http_requests_total), remote(http_requests_total)), "region", "$1", "location", "(.*)"))`,
			expectWarn: true,
		},
		{
			name: "label replace to internal label before an avg",
			expr: `avg by (location) (label_replace(http_requests_total, "zone", "$1", "location", "(.*)"))`,
			expected: `
sum by (location) (
  dedup(
    remote(sum by (location, region) (label_replace(http_requests_total, "zone", "$1", "location", "(.*)"))),
    remote(sum by (location, region) (label_replace(http_requests_total, "zone", "$1", "location", "(.*)"))))) 
  / on (location)
sum by (location) (
  dedup(
    remote(count by (location, region) (label_replace(http_requests_total, "zone", "$1", "location", "(.*)"))),
    remote(count by (location, region) (label_replace(http_requests_total, "zone", "$1", "location", "(.*)"))))) 
`,
		},
		{
			name: "label replace after an aggregation",
			expr: `label_replace(max by (location) (http_requests_total), "region", "$1", "location", "(.*)")`,
			expected: `
label_replace(max by (location) (dedup(
  remote(max by (location, region) (http_requests_total)), 
  remote(max by (location, region) (http_requests_total))
)), "region", "$1", "location", "(.*)")`,
			expectWarn: true,
		},
		{
			name: "binary operation in the operand path",
			expr: `max by (pod) (metric_a / metric_b)`,
			expected: `
max by (pod) (
  dedup(
    remote(max by (pod, region) (metric_a / metric_b)),
    remote(max by (pod, region) (metric_a / metric_b))
  )
)
`,
		},
		{
			name: "binary operation with aggregations",
			expr: `sum by (pod) (metric_a) / sum by (pod) (metric_b)`,
			expected: `
sum by (pod) (dedup(
  remote(sum by (pod, region) (metric_a)), 
  remote(sum by (pod, region) (metric_a)))
)
/ 
sum by (pod) (dedup(
  remote(sum by (pod, region) (metric_b)), 
  remote(sum by (pod, region) (metric_b))
))`,
		},
		{
			name: "function sharding",
			expr: `rate(http_requests_total[2m])`,
			expected: `
dedup(
  remote(rate(http_requests_total[2m])), 
  remote(rate(http_requests_total[2m]))
)`,
		},
		{
			name:     "top level function with no args",
			expr:     `pi()`,
			expected: `pi()`,
		},
		{
			name:     "binary expression with no arg functions",
			expr:     `time() - pi()`,
			expected: `time() - pi()`,
		},
		{
			name: `histogram quantile`,
			expr: `histogram_quantile(0.5, sum by (le) (rate(coredns_dns_request_duration_seconds_bucket[5m])))`,
			expected: `
histogram_quantile(0.5, sum by (le) (dedup(
  remote(sum by (le, region) (rate(coredns_dns_request_duration_seconds_bucket[5m]))), 
  remote(sum by (le, region) (rate(coredns_dns_request_duration_seconds_bucket[5m])))
)))`,
		},
		{
			name:     "binary expression with time",
			expr:     `time() - max by (foo) (bar)`,
			expected: `time() - max by (foo) (dedup(remote(max by (foo, region) (bar)), remote(max by (foo, region) (bar))))`,
		},
		{
			name:     "number literal",
			expr:     `1`,
			expected: `1`,
		},
		{
			name:     "aggregation with number literal",
			expr:     `max(foo) - 1`,
			expected: `max(dedup(remote(max by (region) (foo)), remote(max by (region) (foo)))) - 1`,
		},
		{
			name:     "absent",
			expr:     `absent(foo)`,
			expected: `remote(absent(foo)) * remote(absent(foo))`,
		},
		{
			name:     "absent with aggregation",
			expr:     `sum(absent(foo))`,
			expected: `sum(remote(absent(foo)) * remote(absent(foo)))`,
		},
		{
			name: "binary expression with constant",
			expr: `sum by (pod) (rate(http_requests_total[2m]) * 60)`,
			expected: `sum by (pod) (dedup(
remote(sum by (pod, region) (rate(http_requests_total[2m]) * 60)), 
remote(sum by (pod, region) (rate(http_requests_total[2m]) * 60))))`,
		},
		{
			name:     "binary expression with no arg function",
			expr:     `time() - last_update_timestamp`,
			expected: `time() - dedup(remote(last_update_timestamp), remote(last_update_timestamp))`,
		},
		{
			name:     "subquery",
			expr:     `sum_over_time(http_requests_total[5m:1m])`,
			expected: `dedup(remote(sum_over_time(http_requests_total[5m:1m])), remote(sum_over_time(http_requests_total[5m:1m])))`,
		},
		{
			name:     "subquery over range function",
			expr:     `sum_over_time(rate(http_requests_total[5m])[5m:1m])`,
			expected: `dedup(remote(sum_over_time(rate(http_requests_total[5m])[5m:1m])), remote(sum_over_time(rate(http_requests_total[5m])[5m:1m])))`,
		},
		{
			name: "subquery over range aggregation",
			expr: `sum_over_time(max(http_requests_total)[5m:1m])`,
			expected: `
sum_over_time(max(dedup(
	remote(max by (region) (http_requests_total)) [1969-12-31 23:55:00 +0000 UTC, 1970-01-01 00:00:00 +0000 UTC], 
	remote(max by (region) (http_requests_total)) [1969-12-31 23:55:00 +0000 UTC, 1970-01-01 00:00:00 +0000 UTC])
)[5m:1m])`,
		},
		{
			name:     "label based pruning matches one engine",
			expr:     `sum by (pod) (rate(http_requests_total{region="west"}[2m]))`,
			expected: `sum by (pod) (dedup(remote(sum by (pod, region) (rate(http_requests_total{region="west"}[2m])))))`,
		},
		{
			name:     "label based pruning matches no engines",
			expr:     `http_requests_total{region="north"}`,
			expected: `noop`,
		},
		{
			name:     "label based pruning with grouping matches no engines",
			expr:     `sum by (pod) (rate(http_requests_total{region="north"}[2m]))`,
			expected: `sum by (pod) (noop)`,
		},
		{
			name:     "label based pruning with grouping matches single engine",
			expr:     `sum by (pod) (rate(http_requests_total{region="south"}[2m]))`,
			expected: `sum by (pod) (dedup(remote(sum by (pod, region) (rate(http_requests_total{region="south"}[2m])))))`,
		},
		{
			name:     "binary matching where hash contains partitioning label with on",
			expr:     `X * on (region) Y`,
			expected: `dedup(remote(X * on (region) Y), remote(X * on (region) Y))`,
		},
		{
			name:     "binary matching where hash contains partitioning label with ignoring",
			expr:     `X * ignoring (foo) Y`,
			expected: `dedup(remote(X * ignoring (foo) Y), remote(X * ignoring (foo) Y))`,
		},
		{
			name:     "binary matching where hash doesnt contain partitioning label with ignoring",
			expr:     `X * ignoring (region) Y`,
			expected: `dedup(remote(X), remote(X)) * ignoring (region) dedup(remote(Y), remote(Y))`,
		},
		{
			name:     "binary matching where hash doesnt contain partitioning label with on",
			expr:     `X * on (foo) Y`,
			expected: `dedup(remote(X), remote(X)) * on (foo) dedup(remote(Y), remote(Y))`,
		},

		{
			name: "binary matching and label replace with local label",
			expr: `
count by (cluster) (
	label_replace(up, "ns", "$0", "namespace", ".*")
	* on(region) group_left(project) label_replace(k8s_cluster_info, "k8s_cluster", "$0", "cluster", ".*")
)`,
			expected: `
sum by (cluster) (dedup(
	remote(count by (cluster, region) (label_replace(up, "ns", "$0", "namespace", ".*") * on (region) group_left (project) label_replace(k8s_cluster_info, "k8s_cluster", "$0", "cluster", ".*"))), 
	remote(count by (cluster, region) (label_replace(up, "ns", "$0", "namespace", ".*") * on (region) group_left (project) label_replace(k8s_cluster_info, "k8s_cluster", "$0", "cluster", ".*"))))
)`,
		},
		{
			name: "binary matching and label replace with engine label",
			expr: `
count by (cluster) (
    label_replace(up, "region", "$0", "k8s_region", ".*")
    * on(region) group_left(project) label_replace(k8s_cluster_info, "k8s_cluster", "$0", "cluster", ".*"))`,
			expected: `
count by (cluster) (
 	label_replace(dedup(remote(up), remote(up)), "region", "$0", "k8s_region", ".*")
	* on (region) group_left (project) dedup(
		remote(label_replace(k8s_cluster_info, "k8s_cluster", "$0", "cluster", ".*")),
		remote(label_replace(k8s_cluster_info, "k8s_cluster", "$0", "cluster", ".*"))
	)
)`,
			expectWarn: true,
		},
		{
			name:              "skip binary pushdown when configured",
			expr:              `metric_a / metric_b`,
			expected:          `dedup(remote(metric_a), remote(metric_a)) / dedup(remote(metric_b), remote(metric_b))`,
			skipBinopPushdown: true,
		},
	}

	engines := []api.RemoteEngine{
		newEngineMock(math.MinInt64, math.MaxInt64, []labels.Labels{labels.FromStrings("region", "east"), labels.FromStrings("region", "south")}),
		newEngineMock(math.MinInt64, math.MaxInt64, []labels.Labels{labels.FromStrings("region", "west")}),
	}
	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			optimizers := []Optimizer{
				DistributedExecutionOptimizer{
					Endpoints:          api.NewStaticEndpoints(engines),
					SkipBinaryPushdown: tcase.skipBinopPushdown,
				},
			}

			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			plan := NewFromAST(expr, &query.Options{Start: time.Unix(0, 0), End: time.Unix(0, 0)}, PlanOptions{})
			optimizedPlan, warns := plan.Optimize(optimizers)
			expectedPlan := cleanUp(replacements, tcase.expected)
			testutil.Equals(t, expectedPlan, optimizedPlan.Root().String())
			if tcase.expectWarn {
				testutil.Assert(t, len(warns) > 0, "expected warnings, got none")
			} else {
				testutil.Assert(t, len(warns) == 0, "expected no warnings, got some")
			}
		})
	}
}

type engineOpts struct {
	minTime time.Time
	maxTime time.Time
}

func (o engineOpts) mint() int64 {
	return o.minTime.UnixMilli()
}

func (o engineOpts) maxt() int64 {
	return o.maxTime.UnixMilli()
}

func TestDistributedExecutionWithLongSelectorRanges(t *testing.T) {
	sixHours := 6 * time.Hour
	eightHours := 8 * time.Hour
	twelveHours := 12 * time.Hour

	queryStart := time.Unix(0, 0)
	queryEnd := time.Unix(0, 0).Add(twelveHours)
	queryStep := time.Minute

	cases := []struct {
		name             string
		expr             string
		expected         string
		firstEngineOpts  engineOpts
		secondEngineOpts engineOpts
	}{
		{
			name: "sum over 5m adds a 5 minute offset to latest engine",
			firstEngineOpts: engineOpts{
				minTime: queryStart,
				maxTime: time.Unix(0, 0).Add(eightHours),
			},
			secondEngineOpts: engineOpts{
				minTime: time.Unix(0, 0).Add(sixHours),
				maxTime: queryEnd,
			},
			expr: `sum_over_time(metric[5m])`,
			expected: `
dedup(
  remote(sum_over_time(metric[5m])),
  remote(sum_over_time(metric[5m])) [1970-01-01 06:05:00 +0000 UTC, 1970-01-01 12:00:00 +0000 UTC]
)`,
		},
		{
			name: "sum over 2h adds a 2 hour offset to latest engine",
			firstEngineOpts: engineOpts{
				minTime: queryStart,
				maxTime: time.Unix(0, 0).Add(eightHours),
			},
			secondEngineOpts: engineOpts{
				minTime: time.Unix(0, 0).Add(sixHours),
				maxTime: queryEnd,
			},
			expr: `sum_over_time(metric[2h])`,
			expected: `
dedup(
  remote(sum_over_time(metric[2h])),
  remote(sum_over_time(metric[2h])) [1970-01-01 08:00:00 +0000 UTC, 1970-01-01 12:00:00 +0000 UTC]
)`,
		},
		{
			name: "subquery with a total 2h range is distributed with proper offsets",
			firstEngineOpts: engineOpts{
				minTime: queryStart,
				maxTime: time.Unix(0, 0).Add(eightHours),
			},
			secondEngineOpts: engineOpts{
				minTime: time.Unix(0, 0).Add(sixHours),
				maxTime: queryEnd,
			},
			expr: `sum_over_time(sum_over_time(metric[1h])[1h:30m])`,
			expected: `
dedup(
  remote(sum_over_time(sum_over_time(metric[1h])[1h:30m])), 
  remote(sum_over_time(sum_over_time(metric[1h])[1h:30m])) [1970-01-01 08:00:00 +0000 UTC, 1970-01-01 12:00:00 +0000 UTC]
)`,
		},
		{
			name: "multiple subqueries with a total 90m range get distributed with proper offsets",
			firstEngineOpts: engineOpts{
				minTime: queryStart,
				maxTime: time.Unix(0, 0).Add(eightHours),
			},
			secondEngineOpts: engineOpts{
				minTime: time.Unix(0, 0).Add(sixHours),
				maxTime: queryEnd,
			},
			expr: `max_over_time(sum_over_time(sum_over_time(metric[5m])[45m:10m])[15m:15m])`,
			expected: `dedup(
  remote(max_over_time(sum_over_time(sum_over_time(metric[5m])[45m:10m])[15m:15m])),
  remote(max_over_time(sum_over_time(sum_over_time(metric[5m])[45m:10m])[15m:15m])) [1970-01-01 07:05:00 +0000 UTC, 1970-01-01 12:00:00 +0000 UTC])`,
		},
		{
			name: "subquery with a total 4h range is cannot be distributed",
			firstEngineOpts: engineOpts{
				minTime: queryStart,
				maxTime: time.Unix(0, 0).Add(eightHours),
			},
			secondEngineOpts: engineOpts{
				minTime: time.Unix(0, 0).Add(sixHours),
				maxTime: queryEnd,
			},
			expr:     `sum_over_time(sum_over_time(metric[2h])[2h:30m])`,
			expected: `sum_over_time(sum_over_time(metric[2h])[2h:30m])`,
		},
		{
			name: "sum over 3h does not distribute the query due to insufficient engine overlap",
			firstEngineOpts: engineOpts{
				minTime: queryStart,
				maxTime: time.Unix(0, 0).Add(eightHours),
			},
			secondEngineOpts: engineOpts{
				minTime: time.Unix(0, 0).Add(sixHours),
				maxTime: queryEnd,
			},
			expr:     `sum_over_time(metric[3h])`,
			expected: `sum_over_time(metric[3h])`,
		},
		{
			name: "distribute queries with timestamp",
			firstEngineOpts: engineOpts{
				minTime: queryStart,
				maxTime: time.Unix(0, 0).Add(eightHours),
			},
			secondEngineOpts: engineOpts{
				minTime: time.Unix(0, 0).Add(sixHours),
				maxTime: queryEnd,
			},
			expr: `sum(metric @ 25200)`,
			expected: `
sum(dedup(
  remote(sum by (region) (metric @ 25200.000)), 
  remote(sum by (region) (metric @ 25200.000)) [1970-01-01 06:00:00 +0000 UTC, 1970-01-01 12:00:00 +0000 UTC]
))`,
		},
		{
			name: "skip distributing queries with timestamps outside of the range of an engine",
			firstEngineOpts: engineOpts{
				minTime: queryStart,
				maxTime: time.Unix(0, 0).Add(eightHours),
			},
			secondEngineOpts: engineOpts{
				minTime: time.Unix(0, 0).Add(sixHours),
				maxTime: queryEnd,
			},
			expr:     `sum(metric @ 18000)`,
			expected: `sum(sum by (region) (metric @ 18000.000))`,
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			engines := []api.RemoteEngine{
				newEngineMock(tcase.firstEngineOpts.mint(), tcase.firstEngineOpts.maxt(), []labels.Labels{labels.FromStrings("region", "east")}),
				newEngineMock(tcase.secondEngineOpts.mint(), tcase.secondEngineOpts.maxt(), []labels.Labels{labels.FromStrings("region", "east")}),
			}
			optimizers := []Optimizer{
				DistributedExecutionOptimizer{Endpoints: api.NewStaticEndpoints(engines)},
			}

			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			plan := NewFromAST(expr, &query.Options{Start: queryStart, End: queryEnd, Step: queryStep}, PlanOptions{})
			optimizedPlan, _ := plan.Optimize(optimizers)
			expectedPlan := cleanUp(replacements, tcase.expected)
			testutil.Equals(t, expectedPlan, optimizedPlan.Root().String())
		})
	}
}

func TestDistributedExecutionPruningByTime(t *testing.T) {
	firstEngineOpts := engineOpts{
		minTime: time.Unix(0, 0),
		maxTime: time.Unix(0, 0).Add(6 * time.Hour),
	}
	secondEngineOpts := engineOpts{
		minTime: time.Unix(0, 0).Add(4 * time.Hour),
		maxTime: time.Unix(0, 0).Add(8 * time.Hour),
	}

	cases := []struct {
		name       string
		expr       string
		expected   string
		queryStart time.Time
		queryEnd   time.Time
	}{
		{
			name:       "1 hour query at the end of the range prunes the first engine",
			expr:       `sum(metric)`,
			queryStart: time.Unix(0, 0).Add(7 * time.Hour),
			queryEnd:   time.Unix(0, 0).Add(8 * time.Hour),
			expected:   `sum(dedup(remote(sum by (region) (metric)) [1970-01-01 07:00:00 +0000 UTC, 1970-01-01 08:00:00 +0000 UTC]))`,
		},
		{
			name:       "1 hour range query at the start of the range prunes the second engine",
			expr:       `sum(metric)`,
			queryStart: time.Unix(0, 0).Add(1 * time.Hour),
			queryEnd:   time.Unix(0, 0).Add(2 * time.Hour),
			expected:   `sum(dedup(remote(sum by (region) (metric)) [1970-01-01 01:00:00 +0000 UTC, 1970-01-01 02:00:00 +0000 UTC]))`,
		},
		{
			name:       "instant query in the overlapping range queries both engines",
			expr:       `sum(metric)`,
			queryStart: time.Unix(0, 0).Add(6 * time.Hour),
			queryEnd:   time.Unix(0, 0).Add(6 * time.Hour),
			expected: `
sum(
  dedup(
    remote(sum by (region) (metric)) [1970-01-01 06:00:00 +0000 UTC, 1970-01-01 06:00:00 +0000 UTC], 
    remote(sum by (region) (metric)) [1970-01-01 06:00:00 +0000 UTC, 1970-01-01 06:00:00 +0000 UTC]
  )
)`,
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			engines := []api.RemoteEngine{
				newEngineMock(firstEngineOpts.mint(), firstEngineOpts.maxt(), []labels.Labels{labels.FromStrings("region", "east")}),
				newEngineMock(secondEngineOpts.mint(), secondEngineOpts.maxt(), []labels.Labels{labels.FromStrings("region", "east")}),
			}
			optimizers := []Optimizer{
				DistributedExecutionOptimizer{Endpoints: api.NewStaticEndpoints(engines)},
			}

			expr, err := parser.ParseExpr(tcase.expr)
			testutil.Ok(t, err)

			plan := NewFromAST(expr, &query.Options{Start: tcase.queryStart, End: tcase.queryEnd, Step: time.Minute}, PlanOptions{})
			optimizedPlan, _ := plan.Optimize(optimizers)
			expectedPlan := cleanUp(replacements, tcase.expected)
			testutil.Equals(t, expectedPlan, renderExprTree(optimizedPlan.Root()))
		})
	}
}

func TestDistributedExecutionClonesNodes(t *testing.T) {
	var (
		start    = time.Unix(0, 0)
		end      = time.Unix(0, 0).Add(6 * time.Hour)
		step     = time.Second
		expected = `
sum(dedup(
  remote(sum by (region) (metric{region="east"})), 
  remote(sum by (region) (metric{region="east"}))
))`
	)
	expr, err := parser.ParseExpr(`sum(metric{region="east"})`)
	testutil.Ok(t, err)

	engines := []api.RemoteEngine{
		newEngineMock(math.MinInt64, math.MaxInt64, []labels.Labels{labels.FromStrings("region", "east")}),
		newEngineMock(math.MinInt64, math.MaxInt64, []labels.Labels{labels.FromStrings("region", "east")}),
	}

	lplan := NewFromAST(expr, &query.Options{Start: start, End: end, Step: step}, PlanOptions{})
	optimizedPlan, _ := lplan.Optimize([]Optimizer{
		DistributedExecutionOptimizer{Endpoints: api.NewStaticEndpoints(engines)},
	})

	newMatcher := labels.MustNewMatcher(labels.MatchEqual, "region", "west")
	// Modify the original expression to ensure that changes to not leak into the optimized plan.
	originalVS := expr.(*parser.AggregateExpr).Expr.(*parser.VectorSelector)
	originalVS.LabelMatchers = append(originalVS.LabelMatchers, newMatcher)

	expectedPlan := cleanUp(replacements, expected)
	testutil.Equals(t, expectedPlan, renderExprTree(optimizedPlan.Root()))

	getSelector := func(i int) *VectorSelector {
		return optimizedPlan.Root().(*CheckDuplicateLabels).Expr.(*Aggregation).Expr.(Deduplicate).Expressions[i].Query.(*Aggregation).Expr.(*VectorSelector)
	}

	// Assert that modifying one subquery does not affect the other one.
	vs0 := getSelector(0)
	vs0.LabelMatchers = append(vs0.LabelMatchers, newMatcher)

	vs1 := getSelector(1)
	testutil.Assert(t, len(vs1.LabelMatchers) == len(vs0.LabelMatchers)-1, "expected %d label matchers, got %d", len(vs0.LabelMatchers)-1, len(vs1.LabelMatchers))
}

type engineMock struct {
	api.RemoteEngine
	minT      int64
	maxT      int64
	labelSets []labels.Labels
}

func (e engineMock) MaxT() int64 {
	return e.maxT
}

func (e engineMock) MinT() int64 {
	return e.minT
}

func (e engineMock) LabelSets() []labels.Labels {
	return e.labelSets
}

func newEngineMock(mint, maxt int64, labelSets []labels.Labels) *engineMock {
	return &engineMock{minT: mint, maxT: maxt, labelSets: labelSets}
}
