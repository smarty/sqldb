package sqldb

import (
	"time"

	"github.com/smartystreets/clock"
)

type RetrySelector struct {
	inner    BindingSelector
	sleep    *clock.Sleeper
	duration time.Duration
}

func NewRetrySelector(actual BindingConnectionPool, duration time.Duration) *RetrySelector {
	return &RetrySelector{inner: actual, duration: duration}
}

func (this *RetrySelector) BindSelect(binder Binder, statement string, parameters ...interface{}) error {
	for {
		if this.inner.BindSelect(binder, statement, parameters...) == nil {
			return nil
		}

		this.sleep.Sleep(this.duration)
	}
}
