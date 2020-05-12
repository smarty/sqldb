package sqldb

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

func (this *StackTraceTransaction) Execute(statement string, parameters ...interface{}) (uint64, error) {
	affected, err := this.inner.Execute(statement, parameters...)
	return affected, this.Wrap(err)
}

func (this *StackTraceTransaction) Select(statement string, args ...interface{}) (SelectResult, error) {
	result, err := this.inner.Select(statement, args...)
	return result, this.Wrap(err)
}
