package binary

import (
	"fmt"

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

	outputValues        []sample
	highCardOutputCache []*uint64
	lowCardOutputCache  [][]uint64
}

func newTable(
	pool *model.VectorPool,
	card parser.VectorMatchCardinality,
	expr parser.ItemType,
	outputValues []sample,
	highCardOutputCache []*uint64,
	lowCardOutputCache [][]uint64,
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
		highCardOutputCache: highCardOutputCache,
		lowCardOutputCache:  lowCardOutputCache,
	}, nil
}

func (t *table) execBinaryOperation(lhs model.StepVector, rhs model.StepVector) model.StepVector {
	ts := lhs.T
	step := t.pool.GetStepVector(ts)

	for i, sampleID := range lhs.SampleIDs {
		lhsVal := lhs.Samples[i]
		if t.card == parser.CardOneToMany {
			outputSampleIDs := t.lowCardOutputCache[sampleID]
			for _, outputSampleID := range outputSampleIDs {
				t.outputValues[outputSampleID].t = lhs.T
				t.outputValues[outputSampleID].v = lhsVal
			}
		} else {
			outputSampleID := t.highCardOutputCache[sampleID]
			if outputSampleID == nil {
				continue
			}
			t.outputValues[*outputSampleID].t = lhs.T
			t.outputValues[*outputSampleID].v = lhsVal
		}
	}

	for i, sampleID := range rhs.SampleIDs {
		rhVal := rhs.Samples[i]
		if t.card == parser.CardManyToOne {
			outputSampleIDs := t.lowCardOutputCache[sampleID]
			for _, outputSampleID := range outputSampleIDs {
				lhSample := t.outputValues[outputSampleID]
				if rhs.T != lhSample.t {
					continue
				}

				outputVal := t.operation(lhSample.v, rhVal)
				step.SampleIDs = append(step.SampleIDs, outputSampleID)
				step.Samples = append(step.Samples, outputVal)
			}
		} else {
			outputSampleID := t.highCardOutputCache[sampleID]
			if outputSampleID == nil {
				continue
			}
			lhSample := t.outputValues[*outputSampleID]
			if rhs.T != lhSample.t {
				continue
			}

			outputVal := t.operation(t.outputValues[*outputSampleID].v, rhVal)
			step.SampleIDs = append(step.SampleIDs, *outputSampleID)
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
	return nil, fmt.Errorf("operation not supported")
}
