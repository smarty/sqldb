package sqldb

type SplitStatementTransaction struct {
	inner    Transaction
	executor *SplitStatementExecutor
}

func NewSplitStatementTransaction(inner Transaction, delimiter string) *SplitStatementTransaction {
	return &SplitStatementTransaction{
		inner:    inner,
		executor: NewSplitStatementExecutor(inner, delimiter),
	}
}

func (this *SplitStatementTransaction) Commit() error {
	return this.inner.Commit()
}

func (this *SplitStatementTransaction) Rollback() error {
	return this.inner.Rollback()
}

func (this *SplitStatementTransaction) Execute(statement string, parameters ...interface{}) (uint64, error) {
	return this.executor.Execute(statement, parameters...)
}
func (this *SplitStatementTransaction) ExecuteIdentity(statement string, parameters ...interface{}) (uint64, uint64, error) {
	panic("not implemented")
}

func (this *SplitStatementTransaction) Select(statement string, parameters ...interface{}) (SelectResult, error) {
	return this.inner.Select(statement, parameters...)
}
