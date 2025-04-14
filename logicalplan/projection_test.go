// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"slices"
	"sort"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/query"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/promql/parser"
)

func TestProjectionOptimizer(t *testing.T) {
	cases := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "simple aggregation by no labels",
			expr:     `sum (metric{instance="a", job="b", env="c"})`,
			expected: `sum(metric{env="c",instance="a",job="b"}[projection=include()])`,
		},
		{
			name:     "simple aggregation without no labels",
			expr:     `sum without() (metric{instance="a", job="b", env="c"})`,
			expected: `sum without (__series_hash__) (metric{env="c",instance="a",job="b"})`,
		},
		{
			name:     "simple aggregation",
			expr:     `sum by (job) (metric{instance="a", job="b", env="c"})`,
			expected: `sum by (job) (metric{env="c",instance="a",job="b"}[projection=include(job)])`,
		},
		{
			name:     "multiple aggregations",
			expr:     `sum by (job) (metric{instance="a", job="b", env="c"}) / count by (job) (metric{instance="a", job="b", env="c"})`,
			expected: `sum by (job) (metric{env="c",instance="a",job="b"}[projection=include(job)]) / count by (job) (metric{env="c",instance="a",job="b"}[projection=include(job)])`,
		},
		{
			name:     "binary operation with vector matching",
			expr:     `metric{instance="a", job="b"} * on(job) metric{instance="c", job="b"}`,
			expected: `metric{instance="a",job="b"}[projection=include(job)] * on (job) metric{instance="c",job="b"}[projection=include(job)]`,
		},
		{
			name:     "function call",
			expr:     `rate(metric{instance="a", job="b"}[5m])`,
			expected: `rate(metric{instance="a",job="b"}[5m0s])`,
		},
		{
			name:     "complex query with multiple operations",
			expr:     `sum by (job) (rate(test_metric{instance="a", job="b", env="c"}[5m])) / on(job) group_left count by (job) (metric{instance="d", job="b", env="e"})`,
			expected: `sum by (job) (rate(test_metric{env="c",instance="a",job="b"}[projection=include(job)][5m0s])) / on (job) group_left () count by (job) (metric{env="e",instance="d",job="b"}[projection=include(job)])`,
		},
		{
			name:     "aggregation with without",
			expr:     `sum without (instance) (metric{instance="a", job="b", env="c"})`,
			expected: `sum without (instance, __series_hash__) (metric{env="c",instance="a",job="b"}[projection=exclude(instance)])`,
		},
		{
			name:     "subquery with aggregation",
			expr:     `sum by (job) (count_over_time(up{job="prometheus"}[30m:1m]))`,
			expected: `sum by (job) (count_over_time(up{job="prometheus"}[projection=include(job)][30m0s:1m0s]))`,
		},
		{
			name:     "label_replace with required destination label",
			expr:     `sum by (new_job) (label_replace(metric{instance="a", job="b", env="c"}, "new_job", "$1", "job", "(.+)"))`,
			expected: `sum by (new_job) (label_replace(metric{env="c",instance="a",job="b"}[projection=include(job,new_job)], "new_job", "$1", "job", "(.+)"))`,
		},
		{
			name:     "label_replace with unrequired destination label",
			expr:     `sum by (instance) (label_replace(metric{instance="a", job="b", env="c"}, "new_job", "$1", "job", "(.+)"))`,
			expected: `sum by (instance) (label_replace(metric{env="c",instance="a",job="b"}[projection=include(instance)], "new_job", "$1", "job", "(.+)"))`,
		},
		{
			name:     "label_join with required destination label",
			expr:     `sum by (combined) (label_join(metric{instance="a", job="b", env="c"}, "combined", "-", "job", "env"))`,
			expected: `sum by (combined) (label_join(metric{env="c",instance="a",job="b"}[projection=include(combined,env,job)], "combined", "-", "job", "env"))`,
		},
		{
			name:     "label_join with unrequired destination label",
			expr:     `sum by (instance) (label_join(metric{instance="a", job="b", env="c"}, "combined", "-", "job", "env"))`,
			expected: `sum by (instance) (label_join(metric{env="c",instance="a",job="b"}[projection=include(instance)], "combined", "-", "job", "env"))`,
		},
		{
			name:     "histogram_quantile with aggregation inside",
			expr:     `histogram_quantile(0.9, sum by (le, job) (rate(http_request_duration_seconds_bucket{job="api-server", instance="localhost:9090"}[5m])))`,
			expected: `histogram_quantile(0.9, sum by (le, job) (rate(http_request_duration_seconds_bucket{instance="localhost:9090",job="api-server"}[projection=include(job,le)][5m0s])))`,
		},
		{
			name:     "label_replace with grouping on original label",
			expr:     `sum by (job) (label_replace(metric{instance="a", job="b", env="c"}, "new_job", "$1", "env", "(.+)"))`,
			expected: `sum by (job) (label_replace(metric{env="c",instance="a",job="b"}[projection=include(job)], "new_job", "$1", "env", "(.+)"))`,
		},
		{
			name:     "label_replace with grouping on different label",
			expr:     `sum by (instance) (label_replace(metric{instance="a", job="b", env="c"}, "new_job", "$1", "job", "(.+)"))`,
			expected: `sum by (instance) (label_replace(metric{env="c",instance="a",job="b"}[projection=include(instance)], "new_job", "$1", "job", "(.+)"))`,
		},
		{
			name:     "label_join with grouping on original label",
			expr:     `sum by (job) (label_join(metric{instance="a", job="b", env="c"}, "combined", "-", "env", "instance"))`,
			expected: `sum by (job) (label_join(metric{env="c",instance="a",job="b"}[projection=include(job)], "combined", "-", "env", "instance"))`,
		},
		{
			name:     "label_join with grouping on different label",
			expr:     `sum by (env) (label_join(metric{instance="a", job="b", env="c"}, "combined", "-", "job", "instance"))`,
			expected: `sum by (env) (label_join(metric{env="c",instance="a",job="b"}[projection=include(env)], "combined", "-", "job", "instance"))`,
		},
		{
			name:     "binary operation with ignoring",
			expr:     `metric{instance="a", job="b", env="c"} * ignoring(instance) metric{instance="d", job="b", env="c"}`,
			expected: `metric{env="c",instance="a",job="b"}[projection=exclude(instance)] * ignoring (instance, __series_hash__) metric{env="c",instance="d",job="b"}[projection=exclude(instance)]`,
		},
		{
			name:     "binary operation with ignoring and group_left",
			expr:     `metric{instance="a", job="b", env="c"} * ignoring(instance) group_left(env) metric{instance="d", job="b"}`,
			expected: `metric{env="c",instance="a",job="b"} * ignoring (instance, __series_hash__) group_left (env) metric{instance="d",job="b"}[projection=exclude(instance)]`,
		},
		{
			name:     "binary operation with ignoring and group_right",
			expr:     `metric{instance="a", job="b"} * ignoring(job) group_right(instance) metric{instance="d", job="b", env="e"}`,
			expected: `metric{instance="a",job="b"}[projection=exclude(job)] * ignoring (job, __series_hash__) group_right (instance) metric{env="e",instance="d",job="b"}`,
		},
		{
			name:     "aggregation with binary operation using on",
			expr:     `sum(metric1{instance="a", job="b", env="c"} * on(job) metric2{instance="d", job="b", env="e"})`,
			expected: `sum(metric1{env="c",instance="a",job="b"}[projection=include(job)] * on (job) metric2{env="e",instance="d",job="b"}[projection=include(job)])`,
		},
		{
			name:     "aggregation with binary operation using ignoring",
			expr:     `sum(metric1{instance="a", job="b", env="c"} * ignoring(instance) metric2{instance="d", job="b", env="c"})`,
			expected: `sum(metric1{env="c",instance="a",job="b"}[projection=exclude(instance)] * ignoring (instance, __series_hash__) metric2{env="c",instance="d",job="b"}[projection=exclude(instance)])`,
		},
		{
			name:     "aggregation by label with binary operation using on and group_left",
			expr:     `sum by (job) (metric1{instance="a", job="b", env="c"} * on(job) group_left(env) metric2{instance="d", job="b"})`,
			expected: `sum by (job) (metric1{env="c",instance="a",job="b"} * on (job) group_left (env) metric2{instance="d",job="b"}[projection=include(env,job)])`,
		},
		{
			name:     "aggregation by label with binary operation using on and group_right",
			expr:     `sum by (job) (metric1{instance="a", job="b"} * on(job) group_right(instance) metric2{instance="d", job="b", env="e"})`,
			expected: `sum by (job) (metric1{instance="a",job="b"}[projection=include(instance,job)] * on (job) group_right (instance) metric2{env="e",instance="d",job="b"})`,
		},
		{
			name:     "aggregation by label with binary operation using ignoring and group_left",
			expr:     `sum by (job) (metric1{instance="a", job="b", env="c"} * ignoring(instance) group_left(env) metric2{instance="d", job="b"})`,
			expected: `sum by (job) (metric1{env="c",instance="a",job="b"} * ignoring (instance, __series_hash__) group_left (env) metric2{instance="d",job="b"}[projection=exclude(instance)])`,
		},
		{
			name:     "aggregation by label with binary operation using ignoring and group_right",
			expr:     `sum by (job) (metric1{instance="a", job="b"} * ignoring(instance) group_right(env) metric2{instance="d", job="b", env="e"})`,
			expected: `sum by (job) (metric1{instance="a",job="b"}[projection=exclude(instance)] * ignoring (instance, __series_hash__) group_right (env) metric2{env="e",instance="d",job="b"})`,
		},
		{
			name:     "aggregation without label with binary operation using on",
			expr:     `sum without (instance) (metric1{instance="a", job="b", env="c"} * on(job) metric2{instance="d", job="b", env="e"})`,
			expected: `sum without (instance, __series_hash__) (metric1{env="c",instance="a",job="b"}[projection=include(job)] * on (job) metric2{env="e",instance="d",job="b"}[projection=include(job)])`,
		},
		{
			name:     "aggregation without label with binary operation using ignoring",
			expr:     `sum without (instance) (metric1{instance="a", job="b", env="c"} * ignoring(instance) metric2{instance="d", job="b", env="c"})`,
			expected: `sum without (instance, __series_hash__) (metric1{env="c",instance="a",job="b"}[projection=exclude(instance)] * ignoring (instance, __series_hash__) metric2{env="c",instance="d",job="b"}[projection=exclude(instance)])`,
		},
		{
			name:     "binary operation with on",
			expr:     `metric{instance="a", job="b", env="c"} * on(job) metric{instance="d", job="b", env="e"}`,
			expected: `metric{env="c",instance="a",job="b"}[projection=include(job)] * on (job) metric{env="e",instance="d",job="b"}[projection=include(job)]`,
		},
		{
			name:     "binary operation with on and group_right",
			expr:     `metric{instance="a", job="b"} * on(job) group_right(instance) metric{instance="d", job="b", env="e"}`,
			expected: `metric{instance="a",job="b"}[projection=include(instance,job)] * on (job) group_right (instance) metric{env="e",instance="d",job="b"}`,
		},
		{
			name:     "nested aggregation",
			expr:     `sum by (job) (avg by (job, instance) (metric{instance="a", job="b", env="c"}))`,
			expected: `sum by (job) (avg by (job, instance) (metric{env="c",instance="a",job="b"}[projection=include(instance,job)]))`,
		},
		{
			name:     "nested aggregation with without",
			expr:     `sum without (instance) (avg without (env) (metric{instance="a", job="b", env="c"}))`,
			expected: `sum without (instance, __series_hash__) (avg without (env, __series_hash__) (metric{env="c",instance="a",job="b"}[projection=exclude(env)]))`,
		},
		{
			name:     "nested aggregation with outer by and inner without",
			expr:     `sum by (job) (avg without (env) (metric{instance="a", job="b", env="c"}))`,
			expected: `sum by (job) (avg without (env, __series_hash__) (metric{env="c",instance="a",job="b"}[projection=exclude(env)]))`,
		},
		{
			name:     "nested aggregation with outer without and inner by",
			expr:     `sum without (env) (avg by (job, instance) (metric{instance="a", job="b", env="c"}))`,
			expected: `sum without (env, __series_hash__) (avg by (job, instance) (metric{env="c",instance="a",job="b"}[projection=include(instance,job)]))`,
		},
		{
			name:     "triple nested aggregation with mixed by and without",
			expr:     `sum by (job) (count without (env) (avg by (job, instance) (metric{instance="a", job="b", env="c"})))`,
			expected: `sum by (job) (count without (env, __series_hash__) (avg by (job, instance) (metric{env="c",instance="a",job="b"}[projection=include(instance,job)])))`,
		},
		{
			name:     "nested aggregation with different by labels",
			expr:     `sum by (job) (avg by (instance) (metric{instance="a", job="b", env="c"}))`,
			expected: `sum by (job) (avg by (instance) (metric{env="c",instance="a",job="b"}[projection=include(instance)]))`,
		},
		{
			name:     "aggregation with histogram_quantile",
			expr:     `sum by (job) (histogram_quantile(0.9, rate(http_request_duration_seconds_bucket{instance="a", job="b"}[5m])))`,
			expected: `sum by (job) (histogram_quantile(0.9, rate(http_request_duration_seconds_bucket{instance="a",job="b"}[5m0s])))`,
		},
		{
			name:     "topk aggregation",
			expr:     `topk(3, metric{instance="a", job="b", env="c"})`,
			expected: `topk(3, metric{env="c",instance="a",job="b"})`,
		},
		{
			name:     "bottomk aggregation",
			expr:     `bottomk(5, metric{instance="a", job="b", env="c"})`,
			expected: `bottomk(5, metric{env="c",instance="a",job="b"})`,
		},
		{
			name:     "topk with by clause",
			expr:     `topk by (job) (3, metric{instance="a", job="b", env="c"})`,
			expected: `topk by (job) (3, metric{env="c",instance="a",job="b"})`,
		},
		{
			name:     "bottomk with by clause",
			expr:     `bottomk by (job) (5, metric{instance="a", job="b", env="c"})`,
			expected: `bottomk by (job) (5, metric{env="c",instance="a",job="b"})`,
		},
		{
			name:     "topk with outer aggregation",
			expr:     `sum by (job) (topk(3, metric{instance="a", job="b", env="c"}))`,
			expected: `sum by (job) (topk(3, metric{env="c",instance="a",job="b"}))`,
		},
		{
			name:     "scalar function",
			expr:     `scalar(metric{instance="a", job="b", env="c"})`,
			expected: `scalar(metric{env="c",instance="a",job="b"})`,
		},
		{
			name:     "absent function",
			expr:     `absent(metric{instance="a", job="b", env="c"})`,
			expected: `absent(metric{env="c",instance="a",job="b"})`,
		},
		{
			name:     "absent_over_time function",
			expr:     `absent_over_time(metric{instance="a", job="b", env="c"}[5m])`,
			expected: `absent_over_time(metric{env="c",instance="a",job="b"}[5m0s])`,
		},
		{
			name:     "scalar function with aggregation",
			expr:     `sum by (job) (scalar(metric{instance="a", job="b", env="c"}))`,
			expected: `sum by (job) (scalar(metric{env="c",instance="a",job="b"}))`,
		},
		{
			name:     "absent function with aggregation",
			expr:     `sum by (job) (absent(metric{instance="a", job="b", env="c"}))`,
			expected: `sum by (job) (absent(metric{env="c",instance="a",job="b"}))`,
		},
		{
			name:     "absent_over_time function with aggregation",
			expr:     `sum by (job) (absent_over_time(metric{instance="a", job="b", env="c"}[5m]))`,
			expected: `sum by (job) (absent_over_time(metric{env="c",instance="a",job="b"}[5m0s]))`,
		},
	}

	for _, tc := range cases {
		if tc.name != "simple aggregation by no labels" {
			continue
		}
		t.Run(tc.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tc.expr)
			testutil.Ok(t, err)

			plan := NewFromAST(expr, &query.Options{Start: time.Unix(0, 0), End: time.Unix(0, 0)}, PlanOptions{})
			optimizer := ProjectionOptimizer{SeriesHashLabel: "__series_hash__"}
			optimizedPlan, _ := optimizer.Optimize(plan.Root(), nil)

			result := renderExprTree(optimizedPlan)
			testutil.Equals(t, tc.expected, result)
		})
	}
}

func TestGetFunctionLabelRequirements(t *testing.T) {
	tests := []struct {
		name       string
		funcName   string
		args       []Node
		projection *Projection
		expected   *Projection
	}{
		{
			name:     "label_replace with destination label needed",
			funcName: "label_replace",
			args: []Node{
				&VectorSelector{},
				&StringLiteral{Val: "new_label"},
				&StringLiteral{Val: "replacement"},
				&StringLiteral{Val: "src_label"},
				&StringLiteral{Val: "regex"},
			},
			projection: &Projection{
				Labels:  []string{"new_label"},
				Include: true,
			},
			expected: &Projection{
				Labels:  []string{"new_label", "src_label"},
				Include: true,
			},
		},
		{
			name:     "label_replace with destination label not needed",
			funcName: "label_replace",
			args: []Node{
				&VectorSelector{},
				&StringLiteral{Val: "new_label"},
				&StringLiteral{Val: "replacement"},
				&StringLiteral{Val: "src_label"},
				&StringLiteral{Val: "regex"},
			},
			projection: &Projection{
				Labels:  []string{"other_label"},
				Include: true,
			},
			expected: &Projection{
				Labels:  []string{"other_label"},
				Include: true,
			},
		},
		{
			name:     "label_replace with without clause",
			funcName: "label_replace",
			args: []Node{
				&VectorSelector{},
				&StringLiteral{Val: "new_label"},
				&StringLiteral{Val: "replacement"},
				&StringLiteral{Val: "src_label"},
				&StringLiteral{Val: "regex"},
			},
			projection: &Projection{
				Labels:  []string{"other_label"},
				Include: false,
			},
			expected: &Projection{
				Labels:  []string{"other_label"},
				Include: false,
			},
		},
		{
			name:     "label_join with destination label needed",
			funcName: "label_join",
			args: []Node{
				&VectorSelector{},
				&StringLiteral{Val: "new_label"},
				&StringLiteral{Val: "separator"},
				&StringLiteral{Val: "src_label1"},
				&StringLiteral{Val: "src_label2"},
			},
			projection: &Projection{
				Labels:  []string{"new_label"},
				Include: true,
			},
			expected: &Projection{
				Labels:  []string{"new_label", "src_label1", "src_label2"},
				Include: true,
			},
		},
		{
			name:     "label_join with without clause",
			funcName: "label_join",
			args: []Node{
				&VectorSelector{},
				&StringLiteral{Val: "new_label"},
				&StringLiteral{Val: "separator"},
				&StringLiteral{Val: "src_label1"},
				&StringLiteral{Val: "src_label2"},
			},
			projection: &Projection{
				Labels:  []string{"new_label"},
				Include: false,
			},
			expected: &Projection{
				Labels:  []string{"new_label"},
				Include: false,
			},
		},
		{
			name:     "scalar function returns empty projection",
			funcName: "scalar",
			args: []Node{
				&VectorSelector{},
			},
			projection: &Projection{
				Labels:  []string{"label1"},
				Include: true,
			},
			expected: &Projection{
				Labels:  []string{},
				Include: true,
			},
		},
		{
			name:     "absent function returns empty projection",
			funcName: "absent",
			args: []Node{
				&VectorSelector{},
			},
			projection: &Projection{
				Labels:  []string{"label1"},
				Include: true,
			},
			expected: &Projection{
				Labels:  []string{},
				Include: true,
			},
		},
		{
			name:     "absent_over_time function returns empty projection",
			funcName: "absent_over_time",
			args: []Node{
				&MatrixSelector{},
			},
			projection: &Projection{
				Labels:  []string{"label1"},
				Include: true,
			},
			expected: &Projection{
				Labels:  []string{},
				Include: true,
			},
		},
		{
			name:     "histogram_quantile function returns nil projection",
			funcName: "histogram_quantile",
			args: []Node{
				&NumberLiteral{Val: 0.9},
				&VectorSelector{},
			},
			projection: &Projection{
				Labels:  []string{"label1"},
				Include: true,
			},
			expected: nil,
		},
		{
			name:     "unknown function returns original labels",
			funcName: "unknown_function",
			args:     []Node{},
			projection: &Projection{
				Labels:  []string{"label1"},
				Include: true,
			},
			expected: &Projection{
				Labels:  []string{"label1"},
				Include: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getFunctionLabelRequirements(tt.funcName, tt.args, tt.projection)

			// Check if result is nil when expected is nil
			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil result, got %v", result)
				}
				return
			}

			// Sort labels for consistent comparison
			sort.Strings(result.Labels)
			sort.Strings(tt.expected.Labels)

			// Check if Include matches
			if result.Include != tt.expected.Include {
				t.Errorf("expected Include=%v, got %v", tt.expected.Include, result.Include)
			}

			// Check if all expected labels are in the result
			for _, label := range tt.expected.Labels {
				if !slices.Contains(result.Labels, label) {
					t.Errorf("expected label %s to be in result, but it wasn't", label)
				}
			}

			// Check if result doesn't have unexpected labels
			for _, label := range result.Labels {
				if !slices.Contains(tt.expected.Labels, label) {
					t.Errorf("unexpected label %s in result", label)
				}
			}
		})
	}
}
