package query

import "time"

type Options struct {
	Start time.Time
	End   time.Time
	Step  time.Duration
	Range time.Duration
}
