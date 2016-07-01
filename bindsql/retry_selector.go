package bindsql

import (
	"time"

	"github.com/smartystreets/clock"
)

type RetrySelector struct {
	actual   Connection
	sleep    *clock.Sleeper
	duration time.Duration
}

func NewRetrySelector(actual Connection, duration time.Duration) *RetrySelector {
	return &RetrySelector{actual: actual, duration: duration}
}

func (this *RetrySelector) Select(binder Binder, statement string, parameters ...interface{}) error {
	for {
		if this.actual.Select(binder, statement, parameters...) == nil {
			return nil
		}

		this.sleep.Sleep(this.duration)
	}
}
