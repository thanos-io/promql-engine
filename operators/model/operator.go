package model

import (
	"context"

	"github.com/prometheus/prometheus/model/labels"
)

type VectorOperator interface {
	Next(ctx context.Context) ([]StepVector, error)
	Series(ctx context.Context) ([]labels.Labels, error)
	GetPool() *VectorPool
}
