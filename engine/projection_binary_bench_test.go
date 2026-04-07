// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/thanos-io/promql-engine/execution/binary"
	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/logicalplan"
	"github.com/thanos-io/promql-engine/query"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

func BenchmarkBinaryProjectionPushdown(b *testing.B) {
	// Create mock series with many labels
	numSeries := 10000
	lhsSeries := make([]labels.Labels, numSeries)
	rhsSeries := make([]labels.Labels, 100)

	// LHS: 10000 series with 20 labels each
	for i := range numSeries {
		lhsSeries[i] = labels.FromStrings(
			"__name__", "kube_pod_info",
			"cluster", "c1",
			"node", fmt.Sprintf("n%d", i%100),
			"namespace", fmt.Sprintf("ns%d", i),
			"pod", fmt.Sprintf("p%d", i),
			"label_a", fmt.Sprintf("a%d", i),
			"label_b", fmt.Sprintf("b%d", i),
			"label_c", fmt.Sprintf("c%d", i),
			"label_d", fmt.Sprintf("d%d", i),
			"label_e", fmt.Sprintf("e%d", i),
			"label_f", fmt.Sprintf("f%d", i),
			"label_g", fmt.Sprintf("g%d", i),
			"label_h", fmt.Sprintf("h%d", i),
			"label_i", fmt.Sprintf("i%d", i),
			"label_j", fmt.Sprintf("j%d", i),
			"label_k", fmt.Sprintf("k%d", i),
			"label_l", fmt.Sprintf("l%d", i),
			"label_m", fmt.Sprintf("m%d", i),
			"label_n", fmt.Sprintf("n%d", i),
			"label_o", fmt.Sprintf("o%d", i),
		)
	}

	// RHS: 100 series with 3 labels each
	for i := range 100 {
		rhsSeries[i] = labels.FromStrings(
			"__name__", "kube_node_labels",
			"cluster", "c1",
			"node", fmt.Sprintf("n%d", i),
			"instance_type", fmt.Sprintf("t%d", i),
		)
	}

	opts := &query.Options{
		Start: time.Unix(0, 0),
		End:   time.Unix(3000, 0),
		Step:  30 * time.Second,
	}

	b.Run("without_projection", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; b.Loop(); i++ {
			op, err := binary.NewVectorOperator(
				&mockOperator{series: lhsSeries},
				&mockOperator{series: rhsSeries},
				&parser.VectorMatching{
					Card:           parser.CardManyToOne,
					MatchingLabels: []string{"cluster", "node"},
					On:             true,
					Include:        []string{"instance_type"},
				},
				parser.MUL,
				false,
				nil, // No projection
				opts,
			)
			testutil.Ok(b, err)
			series, _ := op.Series(context.Background())
			if i == 0 {
				b.Logf("Result series count: %d", len(series))
				if len(series) > 0 {
					b.Logf("First series labels: %v", series[0])
					b.Logf("First series label count: %d", series[0].Len())
				}
			}
		}
	})

	b.Run("with_projection", func(b *testing.B) {
		projection := &logicalplan.Projection{
			Labels:  []string{"namespace", "instance_type"},
			Include: true,
		}

		b.ReportAllocs()
		for i := 0; b.Loop(); i++ {
			op, err := binary.NewVectorOperator(
				&mockOperator{series: lhsSeries},
				&mockOperator{series: rhsSeries},
				&parser.VectorMatching{
					Card:           parser.CardManyToOne,
					MatchingLabels: []string{"cluster", "node"},
					On:             true,
					Include:        []string{"instance_type"},
				},
				parser.MUL,
				false,
				projection,
				opts,
			)
			testutil.Ok(b, err)
			series, _ := op.Series(context.Background())
			if i == 0 {
				b.Logf("Result series count: %d", len(series))
				if len(series) > 0 {
					b.Logf("First series labels: %v", series[0])
					b.Logf("First series label count: %d", series[0].Len())
				}
			}
		}
	})
}

func BenchmarkBinaryProjectionPushdownOneToOne(b *testing.B) {
	numSeries := 1000
	lhsSeries := make([]labels.Labels, numSeries)
	rhsSeries := make([]labels.Labels, numSeries)
	for i := range numSeries {
		lhsSeries[i] = labels.FromStrings(
			"__name__", "metric_a",
			"cluster", "c1",
			"node", fmt.Sprintf("n%d", i),
			"namespace", fmt.Sprintf("ns%d", i),
			"pod", fmt.Sprintf("p%d", i),
		)
		rhsSeries[i] = labels.FromStrings(
			"__name__", "metric_b",
			"cluster", "c1",
			"node", fmt.Sprintf("n%d", i),
			"env", fmt.Sprintf("e%d", i),
		)
	}
	opts := &query.Options{
		Start: time.Unix(0, 0),
		End:   time.Unix(3000, 0),
		Step:  30 * time.Second,
	}
	matching := &parser.VectorMatching{
		Card:           parser.CardOneToOne,
		MatchingLabels: []string{"cluster", "node"},
		On:             true,
	}

	b.Run("without_projection", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			op, err := binary.NewVectorOperator(
				&mockOperator{series: lhsSeries},
				&mockOperator{series: rhsSeries},
				matching,
				parser.MUL,
				false,
				nil,
				opts,
			)
			testutil.Ok(b, err)
			_, _ = op.Series(context.Background())
		}
	})
	b.Run("with_projection", func(b *testing.B) {
		projection := &logicalplan.Projection{
			Labels:  []string{"cluster", "node"},
			Include: true,
		}
		b.ReportAllocs()
		for b.Loop() {
			op, err := binary.NewVectorOperator(
				&mockOperator{series: lhsSeries},
				&mockOperator{series: rhsSeries},
				matching,
				parser.MUL,
				false,
				projection,
				opts,
			)
			testutil.Ok(b, err)
			_, _ = op.Series(context.Background())
		}
	})
}

func BenchmarkBinaryProjectionPushdownOneToMany(b *testing.B) {
	// One-to-many: LHS 100, RHS 10000 (RHS is high-card after swap in execution).
	lhsNum, rhsNum := 100, 10000
	lhsSeries := make([]labels.Labels, lhsNum)
	rhsSeries := make([]labels.Labels, rhsNum)
	for i := range lhsNum {
		lhsSeries[i] = labels.FromStrings(
			"__name__", "kube_node_labels",
			"cluster", "c1",
			"node", fmt.Sprintf("n%d", i),
			"instance_type", fmt.Sprintf("t%d", i),
		)
	}
	for i := range rhsNum {
		rhsSeries[i] = labels.FromStrings(
			"__name__", "kube_pod_info",
			"cluster", "c1",
			"node", fmt.Sprintf("n%d", i%lhsNum),
			"namespace", fmt.Sprintf("ns%d", i),
			"pod", fmt.Sprintf("p%d", i),
		)
	}
	opts := &query.Options{
		Start: time.Unix(0, 0),
		End:   time.Unix(3000, 0),
		Step:  30 * time.Second,
	}
	// group_right(namespace): RHS many, include namespace from RHS.
	matching := &parser.VectorMatching{
		Card:           parser.CardOneToMany,
		MatchingLabels: []string{"cluster", "node"},
		On:             true,
		Include:        []string{"namespace"},
	}

	b.Run("without_projection", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			op, err := binary.NewVectorOperator(
				&mockOperator{series: lhsSeries},
				&mockOperator{series: rhsSeries},
				matching,
				parser.MUL,
				false,
				nil,
				opts,
			)
			testutil.Ok(b, err)
			_, _ = op.Series(context.Background())
		}
	})
	b.Run("with_projection", func(b *testing.B) {
		projection := &logicalplan.Projection{
			Labels:  []string{"cluster", "node", "namespace"},
			Include: true,
		}
		b.ReportAllocs()
		for b.Loop() {
			op, err := binary.NewVectorOperator(
				&mockOperator{series: lhsSeries},
				&mockOperator{series: rhsSeries},
				matching,
				parser.MUL,
				false,
				projection,
				opts,
			)
			testutil.Ok(b, err)
			_, _ = op.Series(context.Background())
		}
	})
}

type mockOperator struct {
	series []labels.Labels
}

func (m *mockOperator) Series(context.Context) ([]labels.Labels, error) {
	return m.series, nil
}

func (m *mockOperator) Next(context.Context, []model.StepVector) (int, error) {
	return 0, nil
}

func (m *mockOperator) Explain() []model.VectorOperator {
	return nil
}

func (m *mockOperator) String() string {
	return "mock"
}
