package sqldb

import (
	"context"
	"time"
)

type RetryBindingSelector struct {
	BindingSelector
	duration time.Duration
}

func NewRetryBindingSelector(actual BindingConnectionPool, duration time.Duration) *RetryBindingSelector {
	return &RetryBindingSelector{BindingSelector: actual, duration: duration}
}

func (this *RetryBindingSelector) BindSelect(ctx context.Context, binder Binder, statement string, parameters ...any) error {
	for {
		if this.BindingSelector.BindSelect(ctx, binder, statement, parameters...) == nil {
			return nil
		}

		time.Sleep(this.duration) // TODO: context.WithTimeout()
	}
}
