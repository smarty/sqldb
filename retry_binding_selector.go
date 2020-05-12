package sqldb

import (
	"time"
)

type RetryBindingSelector struct {
	inner    BindingSelector
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

		time.Sleep(this.duration)
	}
}
