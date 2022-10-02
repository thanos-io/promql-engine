// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package binary

import (
	"fmt"
	"math"

	"github.com/efficientgo/core/errors"
	"github.com/prometheus/prometheus/promql/parser"

	"github.com/thanos-community/promql-engine/physicalplan/model"
	"github.com/thanos-community/promql-engine/physicalplan/parse"
)

type sample struct {
	t int64
	v float64
}

type table struct {
	pool *model.VectorPool

	operation operation
	card      parser.VectorMatchCardinality

	outputValues []sample
	// highCardOutputIndex is a mapping from series ID of the high cardinality
	// operator to an output series ID.
	// During joins, each high cardinality series that has a matching
	// low cardinality series will map to exactly one output series.
	highCardOutputIndex outputIndex
	// lowCardOutputIndex is a mapping from series ID of the low cardinality
	// operator to an output series ID.
	// Each series from the low cardinality operator can join with many
	// series of the high cardinality operator.
	lowCardOutputIndex outputIndex
}

func newTable(
	pool *model.VectorPool,
	card parser.VectorMatchCardinality,
	operation operation,
	outputValues []sample,
	highCardOutputCache outputIndex,
	lowCardOutputCache outputIndex,
) *table {
	return &table{
		pool: pool,
		card: card,

		operation:           operation,
		outputValues:        outputValues,
		highCardOutputIndex: highCardOutputCache,
		lowCardOutputIndex:  lowCardOutputCache,
	}
}

func (t *table) execBinaryOperation(lhs model.StepVector, rhs model.StepVector) model.StepVector {
	ts := lhs.T
	step := t.pool.GetStepVector(ts)

	lhsIndex, rhsIndex := t.highCardOutputIndex, t.lowCardOutputIndex
	if t.card == parser.CardOneToMany {
		lhsIndex, rhsIndex = rhsIndex, lhsIndex
	}

	for i, sampleID := range lhs.SampleIDs {
		lhsVal := lhs.Samples[i]
		outputSampleIDs := lhsIndex.outputSamples(sampleID)
		for _, outputSampleID := range outputSampleIDs {
			t.outputValues[outputSampleID].t = lhs.T
			t.outputValues[outputSampleID].v = lhsVal
		}
	}

	for i, sampleID := range rhs.SampleIDs {
		rhVal := rhs.Samples[i]
		outputSampleIDs := rhsIndex.outputSamples(sampleID)
		for _, outputSampleID := range outputSampleIDs {
			lhSample := t.outputValues[outputSampleID]
			if rhs.T != lhSample.t {
				continue
			}

			outputVal, keep := t.operation(lhSample.v, rhVal)
			if !keep {
				continue
			}
			step.SampleIDs = append(step.SampleIDs, outputSampleID)
			step.Samples = append(step.Samples, outputVal)
		}
	}

	return step
}

type operation func(lhs float64, rhs float64) (float64, bool)

var operations = map[string]operation{
	"+":     func(lhs float64, rhs float64) (float64, bool) { return lhs + rhs, true },
	"-":     func(lhs float64, rhs float64) (float64, bool) { return lhs - rhs, true },
	"*":     func(lhs float64, rhs float64) (float64, bool) { return lhs * rhs, true },
	"/":     func(lhs float64, rhs float64) (float64, bool) { return lhs / rhs, true },
	"^":     func(lhs float64, rhs float64) (float64, bool) { return math.Pow(lhs, rhs), true },
	"%":     func(lhs float64, rhs float64) (float64, bool) { return math.Mod(lhs, rhs), true },
	"==":    func(lhs float64, rhs float64) (float64, bool) { return btos(lhs == rhs), true },
	"!=":    func(lhs float64, rhs float64) (float64, bool) { return btos(lhs != rhs), true },
	">":     func(lhs float64, rhs float64) (float64, bool) { return btos(lhs > rhs), true },
	"<":     func(lhs float64, rhs float64) (float64, bool) { return btos(lhs < rhs), true },
	">=":    func(lhs float64, rhs float64) (float64, bool) { return btos(lhs >= rhs), true },
	"<=":    func(lhs float64, rhs float64) (float64, bool) { return btos(lhs <= rhs), true },
	"atan2": func(lhs float64, rhs float64) (float64, bool) { return math.Atan2(lhs, rhs), true },
}

// For vector, those operations are handled differently to check whether to keep
// the value or not. https://github.com/prometheus/prometheus/blob/main/promql/engine.go#L2229
var vectorBinaryOperations = map[string]operation{
	"==": func(lhs float64, rhs float64) (float64, bool) { return lhs, lhs == rhs },
	"!=": func(lhs float64, rhs float64) (float64, bool) { return lhs, lhs != rhs },
	">":  func(lhs float64, rhs float64) (float64, bool) { return lhs, lhs > rhs },
	"<":  func(lhs float64, rhs float64) (float64, bool) { return lhs, lhs < rhs },
	">=": func(lhs float64, rhs float64) (float64, bool) { return lhs, lhs >= rhs },
	"<=": func(lhs float64, rhs float64) (float64, bool) { return lhs, lhs <= rhs },
}

func newOperation(expr parser.ItemType, vectorBinOp bool) (operation, error) {
	t := parser.ItemTypeStr[expr]
	if expr.IsComparisonOperator() && vectorBinOp {
		if o, ok := vectorBinaryOperations[t]; ok {
			return o, nil
		}
		return nil, parse.NotSupportedOperationErr(expr)
	}
	if o, ok := operations[t]; ok {
		return o, nil
	}
	return nil, parse.NotSupportedOperationErr(expr)
}

func unsupportedOperationErr(t string) error {
	msg := fmt.Sprintf("operation not supported: %s", t)
	return errors.Wrap(parse.ErrNotSupportedExpr, msg)
}

// btos returns 1 if b is true, 0 otherwise.
func btos(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
