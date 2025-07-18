package sqldb

import (
	"context"
	"time"
)

type RetryBindingConnectionPool struct {
	BindingConnectionPool
	selector *RetryBindingSelector
}

func NewRetryBindingConnectionPool(inner BindingConnectionPool, sleep time.Duration) *RetryBindingConnectionPool {
	return &RetryBindingConnectionPool{
		BindingConnectionPool: inner,
		selector:              NewRetryBindingSelector(inner, sleep),
	}
}

// This does not implement an override of BeginTransaction because retry makes no sense within the concept of
// a transaction. If the tx fails, it's done.

func (this *RetryBindingConnectionPool) BindSelect(ctx context.Context, binder Binder, statement string, parameters ...any) error {
	return this.selector.BindSelect(ctx, binder, statement, parameters...)
}
