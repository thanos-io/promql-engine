package exchange

import (
	"context"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/thanos-io/promql-engine/execution/model"
	"github.com/thanos-io/promql-engine/execution/telemetry"
	"github.com/thanos-io/promql-engine/query"
	"sync"
	"time"
)

type removeSeriesHashOperator struct {
	telemetry.OperatorTelemetry

	once sync.Once
	next model.VectorOperator

	series []labels.Labels
}

func NewRemoveSeriesHashOperator(next model.VectorOperator, opts *query.Options) model.VectorOperator {
	oper := &removeSeriesHashOperator{
		next: next,
	}
	oper.OperatorTelemetry = telemetry.NewTelemetry(oper, opts)

	return oper
}

func (d *removeSeriesHashOperator) Next(ctx context.Context) ([]model.StepVector, error) {
	return d.next.Next(ctx)
}

func (d *removeSeriesHashOperator) Series(ctx context.Context) ([]labels.Labels, error) {
	start := time.Now()
	defer func() { d.AddExecutionTimeTaken(time.Since(start)) }()

	var err error
	d.once.Do(func() { err = d.loadSeries(ctx) })
	return d.series, err
}

func (d *removeSeriesHashOperator) loadSeries(ctx context.Context) (err error) {
	series, err := d.next.Series(ctx)
	if err != nil {
		return err
	}
	d.series = make([]labels.Labels, len(series))
	for i, lbls := range series {
		lb := labels.NewBuilder(lbls)
		lb = lb.Del("__series_hash__")
		d.series[i] = lb.Labels()
	}
	return nil
}

func (d *removeSeriesHashOperator) GetPool() *model.VectorPool {
	return d.next.GetPool()
}

func (d *removeSeriesHashOperator) Explain() (next []model.VectorOperator) {
	return []model.VectorOperator{d.next}
}

func (d *removeSeriesHashOperator) String() string {
	return "[removeSeriesHash]"
}
