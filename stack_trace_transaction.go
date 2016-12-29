package sqldb

type StackTraceTransaction struct {
	inner    Transaction
}

func NewStackTraceTransaction(inner Transaction) *StackTraceTransaction {
	return &StackTraceTransaction{inner: inner}
}

func (this *StackTraceTransaction) Commit() error {
	panic("implement me")
}

func (this *StackTraceTransaction) Rollback() error {
	panic("implement me")
}

func (this *StackTraceTransaction) Execute(string, ...interface{}) (uint64, error) {
	panic("implement me")
}

func (this *StackTraceTransaction) Select(string, ...interface{}) (SelectResult, error) {
	panic("implement me")
}
