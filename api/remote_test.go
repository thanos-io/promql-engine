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

func TestMemoizedRemoteEngine(t *testing.T) {
	originalLabels := []labels.Labels{labels.FromStrings("region", "east")}
	engine := newEngineMock(10, 20, originalLabels)
	memoized := NewMemoizedRemoteEngine(engine)

	engine.minT = 30
	engine.maxT = 40
	engine.labelSets = []labels.Labels{labels.FromStrings("region", "west")}
	engine.partitionLabelSets = engine.labelSets

	testutil.Equals(t, int64(10), memoized.MinT())
	testutil.Equals(t, int64(20), memoized.MaxT())
	testutil.Equals(t, originalLabels, memoized.LabelSets())
	testutil.Equals(t, originalLabels, memoized.PartitionLabelSets())
}

func TestMemoizedEndpointsRefreshesSnapshots(t *testing.T) {
	engine := newEngineMock(10, 20, nil)
	endpoints := NewMemoizedEndpoints(NewStaticEndpoints([]RemoteEngine{engine}))

	first := endpoints.Engines(0, 100)
	engine.minT = 30
	engine.maxT = 40
	second := endpoints.Engines(0, 100)

	testutil.Equals(t, int64(10), first[0].MinT())
	testutil.Equals(t, int64(20), first[0].MaxT())
	testutil.Equals(t, int64(30), second[0].MinT())
	testutil.Equals(t, int64(40), second[0].MaxT())
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
