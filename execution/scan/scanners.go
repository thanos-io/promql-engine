package scan

import (
	"runtime"

	"github.com/efficientgo/core/errors"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"

	"github.com/thanos-io/promql-engine/execution/exchange"
	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/execution/parse"
	"github.com/thanos-io/promql-engine/logicalplan"
	"github.com/thanos-io/promql-engine/query"
	promstorage "github.com/thanos-io/promql-engine/storage/prometheus"
)

type prometheusScanners struct {
	selectors *promstorage.SelectorPool
}

func NewPrometheusScanners(queryable storage.Queryable) *prometheusScanners {
	return &prometheusScanners{selectors: promstorage.NewSelectorPool(queryable)}
}

func (p prometheusScanners) NewVectorSelector(
	opts *query.Options,
	hints storage.SelectHints,
	logicalNode logicalplan.VectorSelector,
) (model.VectorOperator, error) {
	numShards := runtime.GOMAXPROCS(0) / 2
	if numShards < 1 {
		numShards = 1
	}
	selector := p.selectors.GetFilteredSelector(hints.Start, hints.End, opts.Step.Milliseconds(), logicalNode.VectorSelector.LabelMatchers, logicalNode.Filters, hints)

	operators := make([]model.VectorOperator, 0, numShards)
	for i := 0; i < numShards; i++ {
		operator := exchange.NewConcurrent(
			NewVectorSelector(
				model.NewVectorPool(opts.StepsBatch),
				selector,
				opts,
				logicalNode.Offset,
				logicalNode.BatchSize,
				logicalNode.SelectTimestamp,
				i,
				numShards,
			), 2, opts)
		operators = append(operators, operator)
	}

	return exchange.NewCoalesce(model.NewVectorPool(opts.StepsBatch), opts, logicalNode.BatchSize*int64(numShards), operators...), nil
}

func (p prometheusScanners) NewMatrixSelector(
	opts *query.Options,
	hints storage.SelectHints,
	logicalNode logicalplan.MatrixSelector,
	call parser.Call,
) (model.VectorOperator, error) {
	numShards := runtime.GOMAXPROCS(0) / 2
	if numShards < 1 {
		numShards = 1
	}
	var arg float64
	if call.Func.Name == "quantile_over_time" {
		constVal, err := unwrapConstVal(call.Args[0])
		if err != nil {
			return nil, err
		}
		arg = constVal
	}

	vs := logicalNode.VectorSelector.(*logicalplan.VectorSelector)
	filter := p.selectors.GetFilteredSelector(hints.Start, hints.End, opts.Step.Milliseconds(), vs.LabelMatchers, vs.Filters, hints)

	operators := make([]model.VectorOperator, 0, numShards)
	for i := 0; i < numShards; i++ {
		operator, err := NewMatrixSelector(
			model.NewVectorPool(opts.StepsBatch),
			filter,
			call.Func.Name,
			arg,
			opts,
			logicalNode.Range,
			vs.Offset,
			vs.BatchSize,
			i,
			numShards,
		)
		if err != nil {
			return nil, err
		}
		operators = append(operators, exchange.NewConcurrent(operator, 2, opts))
	}

	return exchange.NewCoalesce(model.NewVectorPool(opts.StepsBatch), opts, vs.BatchSize*int64(numShards), operators...), nil
}

func unwrapConstVal(e parser.Expr) (float64, error) {
	switch c := e.(type) {
	case *parser.NumberLiteral:
		return c.Val, nil
	case *parser.StepInvariantExpr:
		return unwrapConstVal(c.Expr)
	}

	return 0, errors.Wrap(parse.ErrNotSupportedExpr, "matrix selector argument must be a constant")
}
