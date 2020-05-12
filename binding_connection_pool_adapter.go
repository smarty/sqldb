package sqldb

import (
	"context"
)

type BindingConnectionPoolAdapter struct {
	inner            ConnectionPool
	selector         BindingSelector
	panicOnBindError bool
}

func NewBindingConnectionPoolAdapter(actual ConnectionPool, panicOnBindError bool) *BindingConnectionPoolAdapter {
	return &BindingConnectionPoolAdapter{
		inner:            actual,
		selector:         NewBindingSelectorAdapter(actual, panicOnBindError),
		panicOnBindError: panicOnBindError,
	}
}

func (this *BindingConnectionPoolAdapter) Ping(ctx context.Context) error {
	return this.inner.Ping(ctx)
}
func (this *BindingConnectionPoolAdapter) BeginTransaction(ctx context.Context) (BindingTransaction, error) {
	if tx, err := this.inner.BeginTransaction(ctx); err == nil {
		return NewBindingTransactionAdapter(tx, this.panicOnBindError), nil
	} else {
		return nil, err
	}
}
func (this *BindingConnectionPoolAdapter) Close() error {
	return this.inner.Close()
}

func (this *BindingConnectionPoolAdapter) Execute(ctx context.Context, statement string, parameters ...interface{}) (uint64, error) {
	return this.inner.Execute(ctx, statement, parameters...)
}

func (this *BindingConnectionPoolAdapter) BindSelect(ctx context.Context, binder Binder, statement string, parameters ...interface{}) error {
	return this.selector.BindSelect(ctx, binder, statement, parameters...)
}
