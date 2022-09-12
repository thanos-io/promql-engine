package aggregate

import (
	"fmt"

	"github.com/thanos-community/promql-engine/physicalplan/model"

	"github.com/prometheus/prometheus/promql/parser"
	"gonum.org/v1/gonum/floats"
)

type vectorAccumulator func([]float64) float64

type vectorTable struct {
	timestamp   int64
	value       float64
	hasValue    bool
	accumulator vectorAccumulator
}

func newVectorizedTables(stepsBatch int, a parser.ItemType) ([]aggregateTable, error) {
	tables := make([]aggregateTable, stepsBatch)
	for i := 0; i < len(tables); i++ {
		accumulator, err := newVectorAccumulator(a)
		if err != nil {
			return nil, err
		}
		tables[i] = newVectorizedTable(accumulator)
	}

	return tables, nil
}

func newVectorizedTable(a vectorAccumulator) *vectorTable {
	return &vectorTable{
		accumulator: a,
	}
}

func (t *vectorTable) aggregate(vector model.StepVector) {
	if len(vector.SampleIDs) == 0 {
		t.hasValue = false
		return
	}
	t.hasValue = true
	t.timestamp = vector.T
	t.value = t.accumulator(vector.Samples)
}

func (t *vectorTable) toVector(pool *model.VectorPool) model.StepVector {
	result := pool.GetStepVector(t.timestamp)
	if !t.hasValue {
		return result
	}

	result.T = t.timestamp
	result.SampleIDs = append(result.SampleIDs, 0)
	result.Samples = append(result.Samples, t.value)
	return result
}

func (t *vectorTable) size() int {
	return 1
}

func newVectorAccumulator(expr parser.ItemType) (vectorAccumulator, error) {
	t := parser.ItemTypeStr[expr]
	switch t {
	case "sum":
		return floats.Sum, nil
	case "max":
		return floats.Max, nil
	case "min":
		return floats.Min, nil
	case "count":
		return func(float64s []float64) float64 {
			return float64(len(float64s))
		}, nil
	case "avg":
		return func(in []float64) float64 {
			return floats.Sum(in) / float64(len(in))
		}, nil
	}
	return nil, fmt.Errorf("unknown aggregation function %s", t)
}
