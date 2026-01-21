// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package api

import (
	"testing"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/model/labels"
)

func TestCachedEndpoints(t *testing.T) {
	engines := remoteEndpointsFunc(func(mint, maxt int64) []RemoteEngine {
		testutil.Equals(t, int64(10), mint)
		testutil.Equals(t, int64(20), maxt)
		return []RemoteEngine{newEngineMock(0, 1, nil)}
	})
	endpoints := NewCachedEndpoints(engines)

	es := endpoints.Engines(10, 20)
	testutil.Equals(t, 1, len(es))
}

func TestCachedEndpointsCachesEngines(t *testing.T) {
	var calls int
	engines := remoteEndpointsFunc(func(mint, maxt int64) []RemoteEngine {
		calls++
		return []RemoteEngine{
			newEngineMock(100*int64(calls), 1000*int64(calls), nil),
			newEngineMock(200*int64(calls), 2000*int64(calls), nil),
		}
	})
	endpoints := NewCachedEndpoints(engines)

	es1 := endpoints.Engines(10, 10000)
	testutil.Equals(t, 2, len(es1))

	es2 := endpoints.Engines(20, 20000)
	testutil.Equals(t, 2, len(es2))

	testutil.Equals(t, 1, calls)
	testutil.Equals(t, es1, es2)

	// Engines must be mutable.
	es1[0].(*engineMock).maxT = 1337
	testutil.Equals(t, int64(1337), es1[0].MaxT())
	testutil.Equals(t, int64(1337), es2[0].MaxT())
}

type remoteEndpointsFunc func(mint, maxt int64) []RemoteEngine

func (f remoteEndpointsFunc) Engines(mint, maxt int64) []RemoteEngine {
	return f(mint, maxt)
}

type engineMock struct {
	RemoteEngine
	minT               int64
	maxT               int64
	labelSets          []labels.Labels
	partitionLabelSets []labels.Labels
}

func (e engineMock) MaxT() int64                         { return e.maxT }
func (e engineMock) MinT() int64                         { return e.minT }
func (e engineMock) LabelSets() []labels.Labels          { return e.labelSets }
func (e engineMock) PartitionLabelSets() []labels.Labels { return e.partitionLabelSets }

func newEngineMock(mint, maxt int64, labelSets []labels.Labels) *engineMock {
	return &engineMock{minT: mint, maxT: maxt, labelSets: labelSets, partitionLabelSets: labelSets}
}
