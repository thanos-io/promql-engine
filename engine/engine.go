package engine

import (
	"fmt"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/util/stats"
	"time"
)

type engine struct {
	logger               promql.QueryLogger
	lookbackDelta        time.Duration
	enablePerStepStats   bool
	enableAtModifier     bool
	enableNegativeOffset bool
}

func (e *engine) SetQueryLogger(l promql.QueryLogger) {
	e.logger = l
}

func (e *engine) NewInstantQuery(q storage.Queryable, opts *promql.QueryOpts, qs string, ts time.Time) (promql.Query, error) {
	expr, err := parser.ParseExpr(qs)
	if err != nil {
		return nil, err
	}
	qry, err := e.newQuery(q, opts, expr, ts, ts, 0)
	if err != nil {
		return nil, err
	}

	return qry, nil
}

func (e *engine) NewRangeQuery(q storage.Queryable, opts *promql.QueryOpts, qs string, start, end time.Time, interval time.Duration) (promql.Query, error) {
	expr, err := parser.ParseExpr(qs)
	if err != nil {
		return nil, err
	}
	if expr.Type() != parser.ValueTypeVector && expr.Type() != parser.ValueTypeScalar {
		return nil, fmt.Errorf("invalid expression type %q for range query, must be Scalar or instant Vector", parser.DocumentedType(expr.Type()))
	}
	qry, err := e.newQuery(q, opts, expr, start, end, interval)
	if err != nil {
		return nil, err
	}

	return qry, nil
}

func (e *engine) newQuery(q storage.Queryable, opts *promql.QueryOpts, expr parser.Expr, start time.Time, end time.Time, interval time.Duration) (promql.Query, error) {
	if err := e.validateOpts(expr); err != nil {
		return nil, err
	}

	if opts == nil {
		opts = &promql.QueryOpts{}
	}

	lookbackDelta := opts.LookbackDelta
	if lookbackDelta <= 0 {
		lookbackDelta = e.lookbackDelta
	}

	es := &parser.EvalStmt{
		Expr:          expr,
		Start:         start,
		End:           end,
		Interval:      interval,
		LookbackDelta: lookbackDelta,
	}
	qry := &query{
		statement:   es,
		queryable:   q,
		stats:       stats.NewQueryTimers(),
		sampleStats: stats.NewQuerySamples(e.enablePerStepStats && opts.EnablePerStepStats),
	}
	return qry, nil
}

func (e *engine) validateOpts(expr parser.Expr) error {
	if e.enableAtModifier && e.enableNegativeOffset {
		return nil
	}

	var atModifierUsed, negativeOffsetUsed bool

	var validationErr error
	parser.Inspect(expr, func(node parser.Node, path []parser.Node) error {
		switch n := node.(type) {
		case *parser.VectorSelector:
			if n.Timestamp != nil || n.StartOrEnd == parser.START || n.StartOrEnd == parser.END {
				atModifierUsed = true
			}
			if n.OriginalOffset < 0 {
				negativeOffsetUsed = true
			}

		case *parser.MatrixSelector:
			vs := n.VectorSelector.(*parser.VectorSelector)
			if vs.Timestamp != nil || vs.StartOrEnd == parser.START || vs.StartOrEnd == parser.END {
				atModifierUsed = true
			}
			if vs.OriginalOffset < 0 {
				negativeOffsetUsed = true
			}

		case *parser.SubqueryExpr:
			if n.Timestamp != nil || n.StartOrEnd == parser.START || n.StartOrEnd == parser.END {
				atModifierUsed = true
			}
			if n.OriginalOffset < 0 {
				negativeOffsetUsed = true
			}
		}

		if atModifierUsed && !e.enableAtModifier {
			validationErr = promql.ErrValidationAtModifierDisabled
			return validationErr
		}
		if negativeOffsetUsed && !e.enableNegativeOffset {
			validationErr = promql.ErrValidationNegativeOffsetDisabled
			return validationErr
		}

		return nil
	})

	return validationErr
}
