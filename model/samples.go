package model

import "github.com/prometheus/prometheus/model/labels"

type StepSample struct {
	ID uint64

	Metric labels.Labels
	V      float64
}

type StepVector struct {
	T       int64
	Samples []StepSample
}
