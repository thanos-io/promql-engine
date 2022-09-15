// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package query

import "time"

type Options struct {
	Start time.Time
	End   time.Time
	Step  time.Duration
	Range time.Duration
}
