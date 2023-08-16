package limits

import (
	"sync"

	"github.com/efficientgo/core/errors"
	"go.uber.org/atomic"
)

type Limits struct {
	maxSamples int

	curSamplesPerTimestamp sync.Map
}

func NewLimits(maxSamples int) *Limits {
	return &Limits{
		maxSamples: maxSamples,
	}
}

func (l *Limits) AccountSamplesForTimestamp(t int64, n int) error {
	if l.maxSamples == 0 {
		return nil
	}
	v, _ := l.curSamplesPerTimestamp.LoadOrStore(t, atomic.NewInt64(0))
	av := v.(*atomic.Int64)

	if cur := av.Load(); cur+int64(n) > int64(l.maxSamples) {
		return errors.New("query processing would load too many samples into memory in query execution")
	}

	av.Add(int64(n))
	return nil
}
