// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

// Package testcases loads the engine test cases that are kept in YAML so that
// engine implementations in other languages can run them too. The YAML files
// are the source of truth; this package is only the Go binding for them.
package testcases

import (
	_ "embed"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

//go:embed range_queries.yaml
var rangeQueriesYAML []byte

// Case is a single range query test case with all defaults applied.
type Case struct {
	// Name identifies the case and is used as the Go subtest name.
	Name string
	// Description explains what the case covers. It is optional.
	Description string
	// Load is a Prometheus test script `load` block defining the input series.
	Load string
	// Query is the PromQL expression to execute.
	Query string
	// StartMS, EndMS and StepMS define the range query in milliseconds.
	StartMS int64
	EndMS   int64
	StepMS  int64
}

// Start returns the start of the query range.
func (c Case) Start() time.Time { return time.UnixMilli(c.StartMS) }

// End returns the end of the query range.
func (c Case) End() time.Time { return time.UnixMilli(c.EndMS) }

// Step returns the resolution step of the query.
func (c Case) Step() time.Duration { return time.Duration(c.StepMS) * time.Millisecond }

// RangeQueries returns the range query test cases from range_queries.yaml.
func RangeQueries() ([]Case, error) {
	return parse(rangeQueriesYAML)
}

type suite struct {
	Defaults struct {
		StartMS int64 `yaml:"start_ms"`
		EndMS   int64 `yaml:"end_ms"`
		StepMS  int64 `yaml:"step_ms"`
	} `yaml:"defaults"`
	Tests []struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Load        string `yaml:"load"`
		Query       string `yaml:"query"`
		StartMS     *int64 `yaml:"start_ms"`
		EndMS       *int64 `yaml:"end_ms"`
		StepMS      *int64 `yaml:"step_ms"`
	} `yaml:"tests"`
}

func parse(data []byte) ([]Case, error) {
	var s suite
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	if s.Defaults.StepMS <= 0 {
		return nil, fmt.Errorf("defaults.step_ms must be positive, got %d", s.Defaults.StepMS)
	}

	cases := make([]Case, 0, len(s.Tests))
	for i, t := range s.Tests {
		if t.Name == "" {
			return nil, fmt.Errorf("test %d has no name", i)
		}
		if t.Query == "" {
			return nil, fmt.Errorf("test %q has no query", t.Name)
		}
		cases = append(cases, Case{
			Name:        t.Name,
			Description: t.Description,
			Load:        t.Load,
			Query:       t.Query,
			StartMS:     orDefault(t.StartMS, s.Defaults.StartMS),
			EndMS:       orDefault(t.EndMS, s.Defaults.EndMS),
			StepMS:      orDefault(t.StepMS, s.Defaults.StepMS),
		})
	}
	return cases, nil
}

func orDefault(value *int64, def int64) int64 {
	if value == nil {
		return def
	}
	return *value
}
