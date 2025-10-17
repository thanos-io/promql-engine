// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package aggregate

import (
	"context"
	"fmt"
	"math"

	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/execution/parse"
	"github.com/thanos-io/promql-engine/execution/warnings"
	"github.com/thanos-io/promql-engine/extmath"
	"github.com/thanos-io/promql-engine/extwarnings"

	"github.com/efficientgo/core/errors"
	"github.com/prometheus/prometheus/model/histogram"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/promql/parser/posrange"
	"github.com/prometheus/prometheus/util/annotations"
)

// aggregateTable is a table that aggregates input samples into
// output samples for a single step.
type aggregateTable interface {
	// timestamp returns the timestamp of the table.
	// If the table is empty, it returns math.MinInt64.
	timestamp() int64
	// aggregate aggregates the given vector into the table.
	aggregate(vector model.StepVector) error
	// toVector writes out the accumulated result to the given vector and
	// resets the table.
	toVector(ctx context.Context, pool *model.VectorPool) model.StepVector
	// reset resets the table with a new aggregation argument.
	// The argument is currently used for quantile aggregation.
	reset(arg float64)
}

type scalarTable struct {
	ts           int64
	inputs       []uint64
	outputs      []*model.Series
	accumulators []extmath.Accumulator
}

func newScalarTables(stepsBatch int, inputCache []uint64, outputCache []*model.Series, aggregation parser.ItemType) ([]aggregateTable, error) {
	tables := make([]aggregateTable, stepsBatch)
	for i := range tables {
		table, err := newScalarTable(inputCache, outputCache, aggregation)
		if err != nil {
			return nil, err
		}
		tables[i] = table
	}
	return tables, nil
}

func (t *scalarTable) timestamp() int64 {
	return t.ts
}

func newScalarTable(inputSampleIDs []uint64, outputs []*model.Series, aggregation parser.ItemType) (*scalarTable, error) {
	accumulators := make([]extmath.Accumulator, len(outputs))
	for i := range accumulators {
		acc, err := newScalarAccumulator(aggregation)
		if err != nil {
			return nil, err
		}
		accumulators[i] = acc
	}
	return &scalarTable{
		ts:           math.MinInt64,
		inputs:       inputSampleIDs,
		outputs:      outputs,
		accumulators: accumulators,
	}, nil
}

func (t *scalarTable) aggregate(vector model.StepVector) error {
	t.ts = vector.T

	var err error
	for i := range vector.Samples {
		err = extwarnings.Coalesce(err, t.addSample(vector.SampleIDs[i], vector.Samples[i]))
	}
	for i := range vector.Histograms {
		err = extwarnings.Coalesce(err, t.addHistogram(vector.HistogramIDs[i], vector.Histograms[i]))
	}
	return err
}

func (t *scalarTable) addSample(sampleID uint64, sample float64) error {
	outputSampleID := t.inputs[sampleID]
	output := t.outputs[outputSampleID]

	return t.accumulators[output.ID].Add(sample, nil)
}

func (t *scalarTable) addHistogram(sampleID uint64, h *histogram.FloatHistogram) error {
	outputSampleID := t.inputs[sampleID]
	output := t.outputs[outputSampleID]

	return t.accumulators[output.ID].Add(0, h)
}

func (t *scalarTable) reset(arg float64) {
	for i := range t.outputs {
		t.accumulators[i].Reset(arg)
	}
	t.ts = math.MinInt64
}

func (t *scalarTable) toVector(ctx context.Context, pool *model.VectorPool) model.StepVector {
	result := pool.GetStepVector(t.ts)
	for i, v := range t.outputs {
		switch t.accumulators[i].ValueType() {
		case extmath.NoValue:
			continue
		case extmath.SingleTypeValue:
			f, h := t.accumulators[i].Value()
			if h == nil {
				result.AppendSample(pool, v.ID, f)
			} else {
				result.AppendHistogram(pool, v.ID, h)
			}
		case extmath.MixedTypeValue:
			warnings.AddToContext(annotations.NewMixedFloatsHistogramsAggWarning(posrange.PositionRange{}), ctx)
		}
	}
	return result
}

func hashMetric(
	builder labels.ScratchBuilder,
	metric labels.Labels,
	without bool,
	grouping []string,
	groupingSet map[string]struct{},
	buf []byte,
) (uint64, labels.Labels) {
	buf = buf[:0]
	builder.Reset()

	if without {
		metric.Range(func(lbl labels.Label) {
			if lbl.Name == labels.MetricName {
				return
			}
			if _, ok := groupingSet[lbl.Name]; ok {
				return
			}
			builder.Add(lbl.Name, lbl.Value)
		})
		key, _ := metric.HashWithoutLabels(buf, grouping...)
		return key, builder.Labels()
	}

	if len(grouping) == 0 {
		return 0, labels.Labels{}
	}

	metric.Range(func(lbl labels.Label) {
		if _, ok := groupingSet[lbl.Name]; !ok {
			return
		}
		builder.Add(lbl.Name, lbl.Value)
	})
	key, _ := metric.HashForLabels(buf, grouping...)
	return key, builder.Labels()
}

// doing it the prometheus way
// https://github.com/prometheus/prometheus/blob/f379e2eac7134dea12ae1d93ebdcb8109db3a5ef/promql/engine.go#L3809C1-L3833C2
// if ratioLimit > 0 and sampleOffset turns out to be < ratioLimit add sample to the result
// else if ratioLimit < 0 then do ratioLimit+1(switch to positive axis), therefore now we will be taking those samples whose sampleOffset >= 1+ratioLimit (inverting the logic from previous case).
func addRatioSample(ratioLimit float64, series labels.Labels) bool {
	sampleOffset := float64(series.Hash()) / float64(math.MaxUint64)

	return (ratioLimit >= 0 && sampleOffset < ratioLimit) ||
		(ratioLimit < 0 && sampleOffset >= (1.0+ratioLimit))
}

func newScalarAccumulator(expr parser.ItemType) (extmath.Accumulator, error) {
	t := parser.ItemTypeStr[expr]
	switch t {
	case "sum":
		return extmath.NewSumAcc(), nil
	case "max":
		return extmath.NewMaxAcc(), nil
	case "min":
		return extmath.NewMinAcc(), nil
	case "count":
		return extmath.NewCountAcc(), nil
	case "avg":
		return extmath.NewAvgAcc(), nil
	case "group":
		return extmath.NewGroupAcc(), nil
	case "stddev":
		return extmath.NewStdDevAcc(), nil
	case "stdvar":
		return extmath.NewStdVarAcc(), nil
	case "quantile":
		return extmath.NewQuantileAcc(), nil
	case "histogram_avg":
		return extmath.NewHistogramAvgAcc(), nil
	}

	msg := fmt.Sprintf("unknown aggregation function %s", t)
	return nil, errors.Wrap(parse.ErrNotSupportedExpr, msg)
}
