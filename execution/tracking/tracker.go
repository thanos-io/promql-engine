package tracking

import (
	"github.com/efficientgo/core/errors"
	"go.uber.org/atomic"

	"github.com/thanos-io/promql-engine/query"
)

// Tracker can be used to track samples that enter the engine and limit them.
type Tracker struct {
	maxSamples      int
	samplesPerBatch []*atomic.Int64
}

func NewTracker(maxSamples int, opts *query.Options) *Tracker {
	res := &Tracker{
		maxSamples:      maxSamples,
		samplesPerBatch: make([]*atomic.Int64, opts.NumSteps()),
	}
	for i := range res.samplesPerBatch {
		res.samplesPerBatch[i] = atomic.NewInt64(0)
	}
	return res
}

func (t *Tracker) addSamplesAndCheckLimits(batch, n int) error {
	if t.maxSamples == 0 {
		return nil
	}

	if t.samplesPerBatch[batch].Load()+int64(n) > int64(t.maxSamples) {
		return errors.New("query processing would load too many samples into memory in query execution")
	}
	t.samplesPerBatch[batch].Add(int64(n))

	return nil
}

// We only check every 100 added samples if the limit is breached.
// Doing so for every sample would be prohibitively expensive.
const resolution = 100

// Limiter provides a Limiter for the operator to track samples with.
func (t *Tracker) Limiter() *Limiter {
	return &Limiter{
		tracker:    t,
		resolution: resolution,
	}
}

// Limiter is used to check limits in one batch. It will only
// check if the sample is safe to add every "resolution" samples.
// It is not safe for concurrent usage!
type Limiter struct {
	tracker *Tracker

	curBatch     int
	samplesAdded int
	resolution   int
}

func (l *Limiter) StartNewBatch() {
	l.curBatch++
	l.samplesAdded = 0
}

func (l *Limiter) AddSample() error {
	l.samplesAdded++
	if l.samplesAdded%l.resolution == 0 {
		if err := l.tracker.addSamplesAndCheckLimits(l.curBatch-1, l.samplesAdded); err != nil {
			// No need to reset samples here; if we return error; processing stops and no more
			// samples will be added.
			return err
		}
		l.samplesAdded = 0
	}
	return nil
}
