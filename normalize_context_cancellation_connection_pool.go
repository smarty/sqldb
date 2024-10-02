package sqldb

import (
	"context"
	"fmt"
	"strings"
)

type NormalizeContextCancellationConnectionPool struct {
	inner ConnectionPool
}

func NewNormalizeContextCancellationConnectionPool(inner ConnectionPool) *NormalizeContextCancellationConnectionPool {
	return &NormalizeContextCancellationConnectionPool{inner: inner}
}

func (this *NormalizeContextCancellationConnectionPool) Ping(ctx context.Context) error {
	return this.normalizeContextCancellationError(this.inner.Ping(ctx))
}

func (this *NormalizeContextCancellationConnectionPool) BeginTransaction(ctx context.Context) (Transaction, error) {
	if tx, err := this.inner.BeginTransaction(ctx); err == nil {
		return NewStackTraceTransaction(tx), nil
	} else {
		return nil, this.normalizeContextCancellationError(err)
	}
}

func (this *NormalizeContextCancellationConnectionPool) Close() error {
	return this.normalizeContextCancellationError(this.inner.Close())
}

func (this *NormalizeContextCancellationConnectionPool) Execute(ctx context.Context, statement string, parameters ...interface{}) (uint64, error) {
	affected, err := this.inner.Execute(ctx, statement, parameters...)
	return affected, this.normalizeContextCancellationError(err)
}

func (this *NormalizeContextCancellationConnectionPool) Select(ctx context.Context, query string, parameters ...interface{}) (SelectResult, error) {
	result, err := this.inner.Select(ctx, query, parameters...)
	return result, this.normalizeContextCancellationError(err)
}

// TODO remove manual check of "use of closed network connection" with release of https://github.com/go-sql-driver/mysql/pull/1615
func (this *NormalizeContextCancellationConnectionPool) normalizeContextCancellationError(err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "operation was canceled") {
		return fmt.Errorf("%w: %w", context.Canceled, err)
	}
	if strings.Contains(err.Error(), "use of closed network connection") {
		return fmt.Errorf("%w: %w", context.Canceled, err)
	}
	return err
}
