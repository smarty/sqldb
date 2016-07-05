package sqldb

import (
	"time"

	"github.com/smartystreets/clock"
)

type RetryBindingSelector struct {
	inner    BindingSelector
	sleep    *clock.Sleeper
	duration time.Duration
}

func NewRetryBindingSelector(actual BindingConnectionPool, duration time.Duration) *RetryBindingSelector {
	return &RetryBindingSelector{inner: actual, duration: duration}
}

func (this *RetryBindingSelector) BindSelect(binder Binder, statement string, parameters ...interface{}) error {
	for {
		if this.inner.BindSelect(binder, statement, parameters...) == nil {
			return nil
		}

		this.sleep.Sleep(this.duration)
	}
}
