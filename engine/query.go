package engine

import (
	"context"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/util/stats"
)

type query struct {
	statement *parser.EvalStmt
	queryable storage.Queryable

	stats       *stats.QueryTimers
	sampleStats *stats.QuerySamples
}

func (q *query) Exec(ctx context.Context) *promql.Result {
	return nil
}

func (q *query) Close() {}

func (q *query) Statement() parser.Statement { return q.statement }

func (q *query) Stats() *stats.Statistics { panic("implement me") }

func (q *query) Cancel() { panic("implement me") }

func (q *query) String() string { panic("implement me") }
