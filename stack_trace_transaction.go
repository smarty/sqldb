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

func (this *StackTraceTransaction) Execute(statement string, args ...interface{}) (uint64, error) {
	result, err := this.inner.Execute(statement, args...)
	return result, this.Wrap(err)
}
func (this *StackTraceTransaction) ExecuteIdentity(statement string, parameters ...interface{}) (uint64, uint64, error) {
	panic("not implemented")
}

func (this *StackTraceTransaction) Select(statement string, args ...interface{}) (SelectResult, error) {
	result, err := this.inner.Select(statement, args...)
	return result, this.Wrap(err)
}
