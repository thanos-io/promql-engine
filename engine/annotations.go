package engine

import (
	"context"
	"fmt"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/util/annotations"
)

var (
	//lint:ignore faillint - this format is expected for caching and other consumers
	TriggeredPromQLFallback = fmt.Errorf("%w: query triggered fallback in thanos engine", annotations.PromQLInfo)
)

type annotatedQuery struct {
	promql.Query

	annotations annotations.Annotations
}

func (q annotatedQuery) Exec(ctx context.Context) *promql.Result {
	res := q.Query.Exec(ctx)

	res.Warnings = res.Warnings.Merge(q.annotations)

	return res
}

func fallbackAnnotatedQuery(q promql.Query) promql.Query {
	return annotatedQuery{
		Query:       q,
		annotations: annotations.New().Add(TriggeredPromQLFallback),
	}
}
