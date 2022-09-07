package model

import "github.com/prometheus/prometheus/promql"

type Sample struct {
	promql.Sample

	ID uint64
}

type Vector []Sample
