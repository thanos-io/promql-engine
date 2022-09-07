package model

import "github.com/prometheus/prometheus/promql"

type Sample struct {
	promql.Sample

	Signature string
}

type Vector []Sample
