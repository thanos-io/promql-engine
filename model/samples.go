package model

import "github.com/prometheus/prometheus/model/labels"

type Sample struct {
	ID uint64

	Metric labels.Labels
	V      float64
}

type Vector struct {
	T       int64
	Samples []Sample
}
