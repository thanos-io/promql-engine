// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package prometheus

import (
	"context"
	"math"
	"time"

	"github.com/efficientgo/core/errors"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/promql/parser/posrange"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	"github.com/prometheus/prometheus/util/annotations"

	"github.com/thanos-io/promql-engine/execution/exchange"
	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/execution/parse"
	"github.com/thanos-io/promql-engine/execution/warnings"
	"github.com/thanos-io/promql-engine/logicalplan"
	"github.com/thanos-io/promql-engine/query"
)

type Scanners struct {
	selectors *SelectorPool

	querier storage.Querier
}

func (s *Scanners) Close() error {
	return s.querier.Close()
}

// subqueryTimes returns the sum of offsets and ranges of all subqueries in the path.
// If the @ modifier is used, then the offset and range is w.r.t. that timestamp
// (i.e. the sum is reset when we have @ modifier).
// The returned *int64 is the closest timestamp that was seen. nil for no @ modifier.
func subqueryTimes(path []*logicalplan.Node) (time.Duration, time.Duration, *int64) {
	var (
		subqOffset, subqRange time.Duration
		ts                    int64 = math.MaxInt64
	)
	for _, node := range path {
		switch n := (*node).(type) {
		case *logicalplan.Subquery:
			subqOffset += n.OriginalOffset
			subqRange += n.Range
			if n.Timestamp != nil {
				// The @ modifier on subquery invalidates all the offset and
				// range till now. Hence resetting it here.
				subqOffset = n.OriginalOffset
				subqRange = n.Range
				ts = *n.Timestamp
			}
		}
	}
	var tsp *int64
	if ts != math.MaxInt64 {
		tsp = &ts
	}
	return subqOffset, subqRange, tsp
}

func getTimeRangesForSelector(qOpts *query.Options, n *parser.VectorSelector, parents []*logicalplan.Node, evalRange time.Duration) (int64, int64) {
	start, end := qOpts.Start.UnixMilli(), qOpts.End.UnixMilli()
	subqOffset, subqRange, subqTs := subqueryTimes(parents)

	if subqTs != nil {
		// The timestamp on the subquery overrides the eval statement time ranges.
		start = *subqTs
		end = *subqTs
	}

	if n.Timestamp != nil {
		// The timestamp on the selector overrides everything.
		start = *n.Timestamp
		end = *n.Timestamp
	} else {
		offsetMilliseconds := subqOffset.Milliseconds()
		start = start - offsetMilliseconds - subqRange.Milliseconds()
		end -= offsetMilliseconds
	}

	if evalRange == 0 {
		start -= qOpts.LookbackDelta.Milliseconds()
	} else {
		start -= evalRange.Milliseconds()
	}

	start -= n.OriginalOffset.Milliseconds()
	end -= n.OriginalOffset.Milliseconds()

	if parse.IsExtFunction(extractFuncFromPath(parents)) {
		// Buffer more so that we could reasonably
		// inject a zero if there is only one point.
		start -= int64(qOpts.ExtLookbackDelta.Milliseconds())
	}

	return start, end
}

func extractFuncFromPath(p []*logicalplan.Node) string {
	if len(p) == 0 {
		return ""
	}
	switch n := (*(p[len(p)-1])).(type) {
	case *logicalplan.Aggregation:
		return n.Op.String()
	case *logicalplan.FunctionCall:
		return n.Func.Name
	case *logicalplan.Binary:
		// If we hit a binary expression we terminate since we only care about functions
		// or aggregations over a single metric.
		return ""
	}
	return extractFuncFromPath(p[:len(p)-1])
}

// findMinMaxTime returns the time in milliseconds of the earliest and latest point in time the statement will try to process.
// This takes into account offsets, @ modifiers, range selectors, and X functions.
// If the statement does not select series, then FindMinMaxTime returns (0, 0).
func findMinMaxTime(lplan logicalplan.Plan, qOpts *query.Options) (int64, int64) {
	var minTimestamp, maxTimestamp int64 = math.MaxInt64, math.MinInt64
	// Whenever a MatrixSelector is evaluated, evalRange is set to the corresponding range.
	// The evaluation of the VectorSelector inside then evaluates the given range and unsets
	// the variable.
	var evalRange time.Duration

	root := lplan.Root()

	logicalplan.TraverseWithParents(nil, &root, func(parents []*logicalplan.Node, node *logicalplan.Node) {
		switch n := (*node).(type) {
		case *logicalplan.VectorSelector:
			start, end := getTimeRangesForSelector(qOpts, n.VectorSelector, parents, evalRange)
			if start < minTimestamp {
				minTimestamp = start
			}
			if end > maxTimestamp {
				maxTimestamp = end
			}
			evalRange = 0
		case *logicalplan.MatrixSelector:
			evalRange = n.Range
		}
	})

	if maxTimestamp == math.MinInt64 {
		// This happens when there was no selector. Hence no time range to select.
		minTimestamp = 0
		maxTimestamp = 0
	}

	return minTimestamp, maxTimestamp
}

func NewPrometheusScanners(queryable storage.Queryable, qOpts *query.Options, lplan logicalplan.Plan) (*Scanners, error) {
	var min, max int64
	if lplan != nil {
		min, max = findMinMaxTime(lplan, qOpts)
	} else {
		min, max = qOpts.Start.UnixMilli(), qOpts.End.UnixMilli()
	}

	querier, err := queryable.Querier(min, max)
	if err != nil {
		return nil, errors.Wrap(err, "create querier")
	}
	return &Scanners{querier: querier, selectors: NewSelectorPool(querier)}, nil
}

func (p Scanners) NewVectorSelector(
	_ context.Context,
	opts *query.Options,
	hints storage.SelectHints,
	logicalNode logicalplan.VectorSelector,
) (model.VectorOperator, error) {
	selector := p.selectors.GetFilteredSelector(hints.Start, hints.End, opts.Step.Milliseconds(), logicalNode.VectorSelector.LabelMatchers, logicalNode.Filters, hints)
	if logicalNode.DecodeNativeHistogramStats {
		selector = newHistogramStatsSelector(selector)
	}

	operators := make([]model.VectorOperator, 0, opts.DecodingConcurrency)
	for i := 0; i < opts.DecodingConcurrency; i++ {
		operator := exchange.NewConcurrent(
			NewVectorSelector(
				model.NewVectorPool(opts.StepsBatch),
				selector,
				opts,
				logicalNode.Offset,
				logicalNode.BatchSize,
				logicalNode.SelectTimestamp,
				i,
				opts.DecodingConcurrency,
			), 2, opts)
		operators = append(operators, operator)
	}

	return exchange.NewCoalesce(model.NewVectorPool(opts.StepsBatch), opts, logicalNode.BatchSize*int64(opts.DecodingConcurrency), operators...), nil
}

func (p Scanners) NewMatrixSelector(
	ctx context.Context,
	opts *query.Options,
	hints storage.SelectHints,
	logicalNode logicalplan.MatrixSelector,
	call logicalplan.FunctionCall,
) (model.VectorOperator, error) {
	arg := 0.0
	switch call.Func.Name {
	case "quantile_over_time":
		unwrap, err := logicalplan.UnwrapFloat(call.Args[0])
		if err != nil {
			return nil, errors.Wrapf(parse.ErrNotSupportedExpr, "quantile_over_time with expression as first argument is not supported")
		}
		arg = unwrap
		if math.IsNaN(unwrap) || unwrap < 0 || unwrap > 1 {
			warnings.AddToContext(annotations.NewInvalidQuantileWarning(unwrap, posrange.PositionRange{}), ctx)
		}
	case "predict_linear":
		unwrap, err := logicalplan.UnwrapFloat(call.Args[1])
		if err != nil {
			return nil, errors.Wrapf(parse.ErrNotSupportedExpr, "predict_linear with expression as second argument is not supported")
		}
		arg = unwrap
	}

	vs := logicalNode.VectorSelector
	selector := p.selectors.GetFilteredSelector(hints.Start, hints.End, opts.Step.Milliseconds(), vs.LabelMatchers, vs.Filters, hints)
	if logicalNode.VectorSelector.DecodeNativeHistogramStats {
		selector = newHistogramStatsSelector(selector)
	}

	operators := make([]model.VectorOperator, 0, opts.DecodingConcurrency)
	for i := 0; i < opts.DecodingConcurrency; i++ {
		operator, err := NewMatrixSelector(
			model.NewVectorPool(opts.StepsBatch),
			selector,
			call.Func.Name,
			arg,
			opts,
			logicalNode.Range,
			vs.Offset,
			vs.BatchSize,
			i,
			opts.DecodingConcurrency,
		)
		if err != nil {
			return nil, err
		}
		operators = append(operators, exchange.NewConcurrent(operator, 2, opts))
	}

	return exchange.NewCoalesce(model.NewVectorPool(opts.StepsBatch), opts, vs.BatchSize*int64(opts.DecodingConcurrency), operators...), nil
}

type histogramStatsSelector struct {
	SeriesSelector
}

func newHistogramStatsSelector(seriesSelector SeriesSelector) histogramStatsSelector {
	return histogramStatsSelector{SeriesSelector: seriesSelector}
}

func (h histogramStatsSelector) GetSeries(ctx context.Context, shard, numShards int) ([]SignedSeries, error) {
	series, err := h.SeriesSelector.GetSeries(ctx, shard, numShards)
	if err != nil {
		return nil, err
	}
	for i := range series {
		series[i].Series = newHistogramStatsSeries(series[i].Series)
	}
	return series, nil
}

type histogramStatsSeries struct {
	storage.Series
}

func newHistogramStatsSeries(series storage.Series) histogramStatsSeries {
	return histogramStatsSeries{Series: series}
}

func (h histogramStatsSeries) Iterator(it chunkenc.Iterator) chunkenc.Iterator {
	return NewHistogramStatsIterator(h.Series.Iterator(it))
}
