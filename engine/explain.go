// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package engine

import (
	"slices"
	"sync"

	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/execution/telemetry"
	"github.com/thanos-io/promql-engine/logicalplan"

	"github.com/prometheus/prometheus/promql"
)

type ExplainableQuery interface {
	promql.Query

	Explain() *ExplainOutputNode
	Analyze() *AnalyzeOutputNode
}

type AnalyzeOutputNode struct {
	OperatorTelemetry telemetry.OperatorTelemetry `json:"telemetry,omitempty"`
	OperatorID        *uint64                     `json:"operatorId,omitempty"`
	Children          []*AnalyzeOutputNode        `json:"children,omitempty"`

	once                sync.Once
	perStepOnce         sync.Once
	totalSamples        int64
	peakSamples         int64
	totalSamplesPerStep []int64
}

type ExplainOutputNode struct {
	OperatorName string              `json:"name,omitempty"`
	Children     []ExplainOutputNode `json:"children,omitempty"`
}

var _ ExplainableQuery = &compatibilityQuery{}

func (a *AnalyzeOutputNode) TotalSamples() int64 {
	a.aggregateSamples()
	return a.totalSamples
}

func (a *AnalyzeOutputNode) TotalSamplesPerStep() []int64 {
	a.aggregatePerStep()
	return a.totalSamplesPerStep
}

func (a *AnalyzeOutputNode) PeakSamples() int64 {
	a.aggregateSamples()
	return a.peakSamples
}

func (a *AnalyzeOutputNode) aggregateSamples() {
	a.once.Do(func() {
		if nodeSamples := a.OperatorTelemetry.Samples(); nodeSamples != nil {
			a.totalSamples += nodeSamples.TotalSamples
			a.peakSamples += int64(nodeSamples.PeakSamples)
		}

		for _, child := range a.Children {
			childPeak := child.PeakSamples()
			a.peakSamples = max(a.peakSamples, childPeak)

			switch a.OperatorTelemetry.LogicalNode().(type) {
			case *logicalplan.Subquery:
				// Skip aggregating samples for subquery
			case *logicalplan.StepInvariantExpr:
				// A step invariant expression is evaluated once and its samples count at every step.
				// The total must not depend on per-step stats being enabled.
				a.totalSamples += child.TotalSamples() * stepsOf(a)
			default:
				a.totalSamples += child.TotalSamples()
			}
		}
	})
}

// aggregatePerStep builds the per-step samples. It is separate from aggregateSamples so that
// callers who only need totals or peak do not pay for copying the per-step slices.
func (a *AnalyzeOutputNode) aggregatePerStep() {
	a.perStepOnce.Do(func() {
		if nodeSamples := a.OperatorTelemetry.Samples(); nodeSamples != nil {
			// Copy the slice: the children are added into it below and the operator
			// telemetry must stay untouched, or repeated Analyze calls inflate the stats.
			a.totalSamplesPerStep = slices.Clone(nodeSamples.TotalSamplesPerStep)
		}

		for _, child := range a.Children {
			switch a.OperatorTelemetry.LogicalNode().(type) {
			case *logicalplan.Subquery:
				// Skip aggregating samples for subquery
			case *logicalplan.StepInvariantExpr:
				childSamples := child.TotalSamples()
				for i := range a.totalSamplesPerStep {
					a.totalSamplesPerStep[i] += childSamples
				}
			default:
				for i, s := range child.TotalSamplesPerStep() {
					a.totalSamplesPerStep[i] += s
				}
			}
		}
	})
}

func analyzeQuery(op model.VectorOperator) *AnalyzeOutputNode {
	var operatorID *uint64
	if ider, ok := op.(model.OperatorIDer); ok {
		id := ider.OperatorID()
		operatorID = &id
	}
	obsv, ok := model.Unwrap(op).(telemetry.ObservableVectorOperator)
	if !ok {
		return nil
	}

	children := obsv.Explain()
	var childTelemetry []*AnalyzeOutputNode
	for _, child := range children {
		if node := analyzeQuery(child); node != nil {
			childTelemetry = append(childTelemetry, node)
		}
	}

	return &AnalyzeOutputNode{
		OperatorTelemetry: obsv,
		OperatorID:        operatorID,
		Children:          childTelemetry,
	}
}

func explainVector(v model.VectorOperator) *ExplainOutputNode {
	vectors := v.Explain()

	var children []ExplainOutputNode
	for _, vector := range vectors {
		children = append(children, *explainVector(vector))
	}

	return &ExplainOutputNode{
		OperatorName: v.String(),
		Children:     children,
	}
}

// stepsOf returns the number of steps a step invariant node is evaluated for.
// The operator records its own count, which is correct inside subqueries.
func stepsOf(node *AnalyzeOutputNode) int64 {
	tm := node.OperatorTelemetry
	if op, ok := tm.(*telemetry.Operator); ok {
		tm = op.OperatorTelemetry
	}
	if sc, ok := tm.(telemetry.StepCounter); ok {
		return int64(sc.TotalSteps())
	}
	return int64(len(node.totalSamplesPerStep))
}
