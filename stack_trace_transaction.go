package sqldb

import "context"

type StackTraceTransaction struct {
	*StackTrace
	inner Transaction
}

func NewStackTraceTransaction(inner Transaction) *StackTraceTransaction {
	return &StackTraceTransaction{inner: inner}
}

func (this *StackTraceTransaction) Commit() error {
	return this.Wrap(this.inner.Commit())
}

func (this *StackTraceTransaction) Rollback() error {
	return this.Wrap(this.inner.Rollback())
}

func (this *StackTraceTransaction) Execute(ctx context.Context, statement string, parameters ...any) (uint64, error) {
	affected, err := this.inner.Execute(ctx, statement, parameters...)
	return affected, this.Wrap(err)
}

func (this *StackTraceTransaction) ExecuteStatement(ctx context.Context, id, statement string, parameters ...any) (uint64, error) {
	affected, err := this.inner.ExecuteStatement(ctx, id, statement, parameters...)
	return affected, this.Wrap(err)
}

func (this *StackTraceTransaction) Select(ctx context.Context, statement string, args ...any) (SelectResult, error) {
	result, err := this.inner.Select(ctx, statement, args...)
	return result, this.Wrap(err)
}
