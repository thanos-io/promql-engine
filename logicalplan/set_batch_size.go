// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package logicalplan

import (
	"github.com/thanos-io/promql-engine/query"

	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/util/annotations"
)

// SelectorBatchSize configures the batch size of selector based on
// aggregates present in the plan.
type SelectorBatchSize struct {
	// DefaultBatchSize is the series batch size for standard batching.
	// Applied to vector selectors under aggregations.
	DefaultBatchSize int64

	// EnableHighOverlapBatching reduces memory for queries with long lookback windows.
	EnableHighOverlapBatching bool

	// HighOverlapBatchSize is the series batch size for high-overlap queries. Defaults to 1000.
	HighOverlapBatchSize int64

	// HighOverlapThreshold is the overlap threshold that triggers the optimization. Defaults to 100.
	HighOverlapThreshold int64
}

// Optimize configures the batch size of selector based on the query plan.
// If any aggregate is present in the plan, the batch size is set to the configured value.
// The two exceptions where this cannot be done is if the aggregate is quantile, or
// when a binary expression precedes the aggregate.
//
// If EnableHighOverlapBatching is true, this optimizer also detects high-overlap queries
// and switches to high-overlap batching by setting StepsBatch to TotalSteps and reducing
// the series batch size.
func (m SelectorBatchSize) Optimize(plan Node, opts *query.Options) (Node, annotations.Annotations) {
	if m.EnableHighOverlapBatching && opts != nil {
		m.applyHighOverlapBatching(plan, opts)
	}

	m.applyStandardBatchSize(plan)

	return plan, nil
}

func (m SelectorBatchSize) applyStandardBatchSize(plan Node) {
	canBatch := false
	Traverse(&plan, func(current *Node) {
		switch e := (*current).(type) {
		case *FunctionCall:
			//TODO: calls can reduce the labelset of the input; think histogram_quantile reducing
			// multiple "le" labels into one output. We cannot handle this in batching. Revisit
			// what is safe here.
			canBatch = false
		case *Binary:
			canBatch = false
		case *Aggregation:
			if e.Op == parser.QUANTILE || e.Op == parser.TOPK || e.Op == parser.BOTTOMK || e.Op == parser.LIMITK || e.Op == parser.LIMIT_RATIO {
				canBatch = false
				return
			}
			canBatch = true
		case *VectorSelector:
			if canBatch {
				e.BatchSize = m.DefaultBatchSize
			}
			canBatch = false
		}
	})
}

func (m SelectorBatchSize) applyHighOverlapBatching(plan Node, opts *query.Options) {
	overlapThreshold := m.HighOverlapThreshold
	if overlapThreshold == 0 {
		overlapThreshold = 100
	}

	seriesBatchSize := m.HighOverlapBatchSize
	if seriesBatchSize == 0 {
		seriesBatchSize = 1000
	}

	vectorSelectors := make(map[*VectorSelector]bool)
	shouldBatch := false

	Traverse(&plan, func(current *Node) {
		ms, ok := (*current).(*MatrixSelector)
		if !ok {
			return
		}

		selectRangeMs := ms.Range.Milliseconds()
		stepMs := opts.Step.Milliseconds()
		if stepMs == 0 {
			stepMs = 1
		}

		overlap := (selectRangeMs-1)/stepMs + 1
		totalSteps := int64(opts.TotalSteps())
		if overlap > totalSteps {
			overlap = totalSteps
		}

		if overlap > overlapThreshold {
			vectorSelectors[ms.VectorSelector] = true
			shouldBatch = true
		}
	})

	if shouldBatch {
		opts.StepsBatch = opts.TotalSteps()

		for vs := range vectorSelectors {
			vs.BatchSize = seriesBatchSize
		}
	}
}
