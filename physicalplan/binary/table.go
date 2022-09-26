// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package binary

import (
	"fmt"

	"github.com/efficientgo/core/errors"

	"github.com/thanos-community/promql-engine/physicalplan/parse"

	"github.com/prometheus/prometheus/promql/parser"

	"github.com/thanos-community/promql-engine/physicalplan/model"
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
	expr parser.ItemType,
	outputValues []sample,
	highCardOutputCache outputIndex,
	lowCardOutputCache outputIndex,
) (*table, error) {
	op, err := newOperation(expr)
	if err != nil {
		return nil, err
	}
	return &table{
		pool: pool,
		card: card,

		operation:           op,
		outputValues:        outputValues,
		highCardOutputIndex: highCardOutputCache,
		lowCardOutputIndex:  lowCardOutputCache,
	}, nil
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

			outputVal := t.operation(lhSample.v, rhVal)
			step.SampleIDs = append(step.SampleIDs, outputSampleID)
			step.Samples = append(step.Samples, outputVal)
		}
	}

	return step
}

type operation func(lhs float64, rhs float64) float64

var operations = map[string]operation{
	"+": func(lhs float64, rhs float64) float64 { return lhs + rhs },
	"-": func(lhs float64, rhs float64) float64 { return lhs - rhs },
	"*": func(lhs float64, rhs float64) float64 { return lhs * rhs },
	"/": func(lhs float64, rhs float64) float64 { return lhs / rhs },
}

func newOperation(expr parser.ItemType) (operation, error) {
	t := parser.ItemTypeStr[expr]
	if o, ok := operations[t]; ok {
		return o, nil
	}
	msg := fmt.Sprintf("operation not supported: %s", t)
	return nil, errors.Wrap(parse.ErrNotSupportedExpr, msg)
}
