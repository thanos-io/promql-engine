package limits

import (
	"github.com/efficientgo/core/errors"
	"go.uber.org/atomic"

	"github.com/thanos-io/promql-engine/query"
)

// We only check every 100 added samples if the limit is breached.
// Doing so for every sample would be prohibitively expensive.
const resolution = 100

type Limits struct {
	maxSamples      int
	samplesPerBatch []*atomic.Int64
}

// NewLimits returns a pointer to a Limits struct. It can be used to
// track samples that enter the engine in one batch and limit it
// to a maximum number.
func NewLimits(maxSamples int, opts *query.Options) *Limits {
	limits := &Limits{
		maxSamples:      maxSamples,
		samplesPerBatch: make([]*atomic.Int64, opts.NumSteps()),
	}
	for i := range limits.samplesPerBatch {
		limits.samplesPerBatch[i] = atomic.NewInt64(0)
	}
	return limits
}

// AccountSamplesForTimestamp keeps track of the samples used for the batch for timestamp t.
// It will return an error if a batch wants to add use more samples then the configured
// maxSamples value.
func (l *Limits) addSamplesAndCheckLimits(batch, n int) error {
	if l.maxSamples == 0 {
		return nil
	}

	if l.samplesPerBatch[batch].Load()+int64(n) > int64(l.maxSamples) {
		return errors.New("query processing would load too many samples into memory in query execution")
	}
	l.samplesPerBatch[batch].Add(int64(n))

	return nil
}

func (l *Limits) Accounter() *Accounter {
	return &Accounter{
		limits:     l,
		resolution: resolution,
	}
}

// Accounter is used to check limits in one batch. It will only
// check if the sample is safe to add every "resolution" samples.
// It is not safe for concurrent usage!
type Accounter struct {
	limits *Limits

	curBatch     int
	samplesAdded int
	resolution   int
}

func (acc *Accounter) StartNewBatch() {
	acc.curBatch++
	acc.samplesAdded = 0
}

func (acc *Accounter) AddSample() error {
	acc.samplesAdded++
	if acc.samplesAdded%acc.resolution == 0 {
		if err := acc.limits.addSamplesAndCheckLimits(acc.curBatch-1, acc.samplesAdded); err != nil {
			// No need to reset samples here; if we return error; processing stops and no more
			// samples will be added.
			return err
		}
		acc.samplesAdded = 0
	}
	return nil
}
