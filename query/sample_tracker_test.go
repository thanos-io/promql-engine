// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package query

import (
	"testing"
)

func TestSampleTracker_WithLimit(t *testing.T) {
	tracker := NewSampleTracker(100)

	tracker.Add(50)
	if err := tracker.CheckLimit(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	tracker.Add(60)
	if err := tracker.CheckLimit(); err == nil {
		t.Error("expected error when exceeding limit")
	}

	if tracker.Limit() != 100 {
		t.Errorf("expected limit 100, got %d", tracker.Limit())
	}
}

func TestSampleTracker_NoLimit(t *testing.T) {
	tracker := NewSampleTracker(0)

	tracker.Add(1000000)
	if err := tracker.CheckLimit(); err != nil {
		t.Errorf("nop tracker should never error: %v", err)
	}

	tracker.Add(1000000)
	if err := tracker.CheckLimit(); err != nil {
		t.Errorf("nop tracker should never error: %v", err)
	}
}

func TestSampleTracker_Remove(t *testing.T) {
	tracker := NewSampleTracker(100)

	tracker.Add(90)
	tracker.Remove(40)
	tracker.Add(40)

	if err := tracker.CheckLimit(); err != nil {
		t.Errorf("unexpected error after remove: %v", err)
	}
}
