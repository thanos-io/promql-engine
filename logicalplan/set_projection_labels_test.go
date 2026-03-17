package logicalplan

import (
	"fmt"
	"strings"
	"testing"

	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/util/annotations"
	"github.com/stretchr/testify/require"
	"github.com/thanos-io/promql-engine/query"
)

var reSpaces = strings.NewReplacer("\n", "", "\t", "")

func TestOptimizeSetProjectionLabels(t *testing.T) {
	cases := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "simple vector selector",
			expr:     `metric_a{job="api-server"}`,
			expected: `metric_a{job="api-server"}[exclude()]`,
		},
		{
			name:     "top-level label_replace",
			expr:     `label_replace(kube_node_info{node="gke-1"}, "instance", "$1", "node", "(.*)")`,
			expected: `label_replace(kube_node_info{node="gke-1"}[exclude()], "instance", "$1", "node", "(.*)")`,
		},
		{
			name:     "sum by all labels",
			expr:     `sum(label_replace(kube_node_info{node="gke-1"}, "instance", "$1", "node", "(.*)"))`,
			expected: `sum(label_replace(kube_node_info{node="gke-1"}[project()], "instance", "$1", "node", "(.*)"))`,
		},
		{
			name:     "sum by target label",
			expr:     `sum by (instance) (label_replace(kube_node_info{node="gke-1"}, "instance", "$1", "node", "(.*)"))`,
			expected: `sum by (instance) (label_replace(kube_node_info{node="gke-1"}[project(instance, node)], "instance", "$1", "node", "(.*)"))`,
		},
		{
			name:     "sum not including target label",
			expr:     `sum by (node, region) (label_replace(kube_node_info{node="gke-1"}, "instance", "$1", "node", "(.*)"))`,
			expected: `sum by (node, region) (label_replace(kube_node_info{node="gke-1"}[project(node, region)], "instance", "$1", "node", "(.*)"))`,
		},
		{
			name:     "sum by source and target label",
			expr:     `sum by (node, instance, region) (label_replace(kube_node_info{node="gke-1"}, "instance", "$1", "node", "(.*)"))`,
			expected: `sum by (node, instance, region) (label_replace(kube_node_info{node="gke-1"}[project(instance, node, region)], "instance", "$1", "node", "(.*)"))`,
		},
		{
			name: "multiple label replace calls",
			expr: `
sum by (instance, node, region) (
  label_replace(
    label_replace(kube_node_info{node="gke-1"}, "ip-addr", "$1", "ip", "(.*)"),
    "instance", "$1", "node", "(.*)"
  )
)`,
			expected: `sum by (instance, node, region) (label_replace(label_replace(kube_node_info{node="gke-1"}[project(instance, node, region)], "ip-addr", "$1", "ip", "(.*)"), "instance", "$1", "node", "(.*)"))`,
		},
		{
			name:     "sum without",
			expr:     `sum without (xyz) (label_replace(kube_node_info{node="gke-1"}, "instance", "$1", "node", "(.*)"))`,
			expected: `sum without (xyz) (label_replace(kube_node_info{node="gke-1"}[exclude(xyz)], "instance", "$1", "node", "(.*)"))`,
		},
		{
			name:     "absent",
			expr:     `absent(kube_node_info{node="gke-1"})`,
			expected: `absent(kube_node_info{node="gke-1"}[project()])`,
		},
		{
			name:     "aggregation with grouping",
			expr:     `sum by (pod) (kube_node_info{node="gke-1"})`,
			expected: `sum by (pod) (kube_node_info{node="gke-1"}[project(pod)])`,
		},
		{
			name:     "double aggregation with grouping",
			expr:     `max by (pod) (sum by (pod, target) (kube_node_info{node="gke-1"}))`,
			expected: `max by (pod) (sum by (pod, target) (kube_node_info{node="gke-1"}[project(pod, target)]))`,
		},
		{
			name:     "double aggregation with by and without grouping",
			expr:     `max by (pod) (sum without (pod, target) (kube_node_info{node="gke-1"}))`,
			expected: `max by (pod) (sum without (pod, target) (kube_node_info{node="gke-1"}[exclude(pod, target)]))`,
		},
		{
			name:     "double aggregation with by and without grouping",
			expr:     `max by (pod) (sum without (target) (kube_node_info{node="gke-1"}))`,
			expected: `max by (pod) (sum without (target) (kube_node_info{node="gke-1"}[exclude(target)]))`,
		},
		{
			name:     "aggregation without grouping",
			expr:     `sum without (pod) (kube_node_info{node="gke-1"})`,
			expected: `sum without (pod) (kube_node_info{node="gke-1"}[exclude(pod)])`,
		},
		{
			name:     "aggregation with binary expression",
			expr:     `sum without (pod) (metric_a * on (node) metric_b)`,
			expected: `sum without (pod) (metric_a[exclude()] * on (node) metric_b[project(__series__id, node)])`,
		},
		{
			name:     "binary expression with vector and constant",
			expr:     `sum(metric_a * 3)`,
			expected: `sum(metric_a[project()] * 3)`,
		},
		{
			name:     "binary expression with aggregation and constant",
			expr:     `sum(metric_a) * 3`,
			expected: `sum(metric_a[project()]) * 3`,
		},
		{
			name:     "binary expression with one to one matching",
			expr:     `metric_a - metric_b`,
			expected: `metric_a[exclude()] - metric_b[exclude()]`,
		},
		{
			name:     "binary expression with one to one matching on label",
			expr:     `metric_a - on (node) metric_b`,
			expected: `metric_a[exclude()] - on (node) metric_b[project(__series__id, node)]`,
		},
		{
			name:     "binary expression with one to one matching on label group_left",
			expr:     `metric_a - on (node) group_left (cluster) metric_b`,
			expected: `metric_a[exclude()] - on (node) group_left (cluster) metric_b[project(__series__id, cluster, node)]`,
		},
		{
			name:     "binary expression with one to one matching on label group_right",
			expr:     `metric_a - on (node) group_right (cluster) metric_b`,
			expected: `metric_a[project(__series__id, cluster, node)] - on (node) group_right (cluster) metric_b[exclude()]`,
		},
		{
			name:     "aggregation with binary expression and one to one matching",
			expr:     `max by (k8s_cluster) (metric_a * up)`,
			expected: `max by (k8s_cluster) (metric_a[exclude()] * up[exclude()])`,
		},
		{
			name:     "aggregation with binary expression with one to one matching on one label",
			expr:     `max by (k8s_cluster) (metric_a * on(node) up)`,
			expected: `max by (k8s_cluster) (metric_a[project(k8s_cluster, node)] * on (node) up[project(__series__id, node)])`,
		},
		{
			name:     "aggregation with binary expression with matching one label group_left",
			expr:     `max by (k8s_cluster) (metric_a * on(node) group_left(hostname) up)`,
			expected: `max by (k8s_cluster) (metric_a[project(k8s_cluster, node)] * on (node) group_left (hostname) up[project(__series__id, hostname, node)])`,
		},
		{
			name:     "aggregation with binary expression with matching one label group_right",
			expr:     `max by (k8s_cluster) (metric_a * on(node) group_right(hostname) up)`,
			expected: `max by (k8s_cluster) (metric_a[project(__series__id, hostname, node)] * on (node) group_right (hostname) up[project(k8s_cluster, node)])`,
		},
		{
			name: "binary expression with aggregation and label replace",
			expr: `
topk(5, 
    sum by (k8s_cluster) (
        max(metric_a) by (node) 
        * on(node) group_right(kubernetes_io_hostname) label_replace(label_replace(label_replace(up, "node", "$1", "kubernetes_io_hostname", "(.*)"),"node_role", "$1", "role", "(.*)"), "region", "$1", "topology_kubernetes_io_region", "(.*)")
        * on(k8s_cluster) group_left(project) label_replace(k8s_cluster_info, "k8s_cluster", "$0", "cluster", ".*")
    )
)`,
			expected: `
topk(5, 
	sum by (k8s_cluster) (
		max by (node) (metric_a[project(node)]) 
		* on (node) group_right (kubernetes_io_hostname) label_replace(label_replace(label_replace(up[project(k8s_cluster, kubernetes_io_hostname, node)], "node", "$1", "kubernetes_io_hostname", "(.*)"), "node_role", "$1", "role", "(.*)"), "region", "$1", "topology_kubernetes_io_region", "(.*)") 
		* on (k8s_cluster) group_left (project) label_replace(k8s_cluster_info[project(__series__id, cluster, k8s_cluster, project)], "k8s_cluster", "$0", "cluster", ".*")))`,
		},
		{
			name: "binary expression with aggregation and label replace",
			expr: `
count by (cluster) (
    label_replace(up, "region", "$0", "region", ".*")
    * on(cluster, region) group_left(project) label_replace(max by(project, region, cluster)(k8s_cluster_info), "k8s_cluster", "$0", "cluster", ".*")
)`,
			expected: `
count by (cluster) (
	label_replace(up[project(cluster, region)], "region", "$0", "region", ".*")
	 * on (cluster, region) group_left (project) label_replace(max by (project, region, cluster) (
		k8s_cluster_info[project(cluster, project, region)]), "k8s_cluster", "$0", "cluster", ".*"))`,
		},
	}

	for _, c := range cases {
		t.Run(c.expr, func(t *testing.T) {
			expr, err := parser.ParseExpr(c.expr)
			require.NoError(t, err)
			plan := NewFromAST(expr, &query.Options{}, PlanOptions{})
			optimized, annos := plan.Optimize(
				[]Optimizer{
					setProjectionLabels{},
					// This is a dummy optimizer that replaces VectorSelectors with a custom struct
					// which has a custom String() method.
					swapVectorSelectors{},
				})

			require.Equal(t, annotations.Annotations{}, annos)
			require.Equal(t, reSpaces.Replace(c.expected), reSpaces.Replace(optimized.Root().String()))
		})
	}
}

type swapVectorSelectors struct{}

func (s swapVectorSelectors) Optimize(plan Node, _ *query.Options) (Node, annotations.Annotations) {
	TraverseBottomUp(nil, &plan, func(_, expr *Node) bool {
		switch v := (*expr).(type) {
		case *VectorSelector:
			*expr = newVectorOutput(v)
			return true
		}
		return false
	})
	return plan, annotations.Annotations{}
}

type vectorOutput struct {
	*VectorSelector
}

func newVectorOutput(vectorSelector *VectorSelector) *vectorOutput {
	return &vectorOutput{
		VectorSelector: vectorSelector,
	}
}

func (vs vectorOutput) String() string {
	var projectionType string
	if vs.Projection.Include {
		projectionType = "project"
	} else {
		projectionType = "exclude"
	}
	return fmt.Sprintf("%s[%s(%s)]", vs.VectorSelector.String(), projectionType, strings.Join(vs.Projection.Labels, ", "))
}
