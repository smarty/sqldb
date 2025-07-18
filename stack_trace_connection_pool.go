package sqldb

import "context"

type StackTraceConnectionPool struct {
	inner ConnectionPool
	*StackTrace
}

func NewStackTraceConnectionPool(inner ConnectionPool) *StackTraceConnectionPool {
	return &StackTraceConnectionPool{inner: inner}
}

func (this *StackTraceConnectionPool) Ping(ctx context.Context) error {
	return this.Wrap(this.inner.Ping(ctx))
}

func (this *StackTraceConnectionPool) BeginTransaction(ctx context.Context) (Transaction, error) {
	if tx, err := this.inner.BeginTransaction(ctx); err == nil {
		return NewStackTraceTransaction(tx), nil
	} else {
		return nil, this.Wrap(err)
	}
}

func (this *StackTraceConnectionPool) Close() error {
	return this.Wrap(this.inner.Close())
}

func (this *StackTraceConnectionPool) Execute(ctx context.Context, statement string, parameters ...any) (uint64, error) {
	affected, err := this.inner.Execute(ctx, statement, parameters...)
	return affected, this.Wrap(err)
}

func (this *StackTraceConnectionPool) Select(ctx context.Context, query string, parameters ...any) (SelectResult, error) {
	result, err := this.inner.Select(ctx, query, parameters...)
	return result, this.Wrap(err)
}
