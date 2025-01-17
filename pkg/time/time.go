package time

import (
	"time"
)

type Time interface {
	Since(sinceTime time.Time) time.Duration
}

type SystemTime struct{}

func (t *SystemTime) Since(sinceTime time.Time) time.Duration {
	return time.Since(sinceTime)
}

type StubTime struct {
	NowTime time.Time
}

func (t *StubTime) Since(sinceTime time.Time) time.Duration {
	return t.NowTime.Sub(sinceTime)
}
