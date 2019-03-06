package sqldb

type SplitStatementConnectionPool struct {
	inner     ConnectionPool
	delimiter string
	executor  *SplitStatementExecutor
}

func NewSplitStatementConnectionPool(inner ConnectionPool, delimiter string) *SplitStatementConnectionPool {
	return &SplitStatementConnectionPool{
		inner:     inner,
		delimiter: delimiter,
		executor:  NewSplitStatementExecutor(inner, delimiter),
	}
}

func (this *SplitStatementConnectionPool) Ping() error {
	return this.inner.Ping()
}

func (this *SplitStatementConnectionPool) BeginTransaction() (Transaction, error) {
	if transaction, err := this.inner.BeginTransaction(); err == nil {
		return NewSplitStatementTransaction(transaction, this.delimiter), nil
	} else {
		return nil, err
	}
}

func (this *SplitStatementConnectionPool) Close() error {
	return this.inner.Close()
}

func (this *SplitStatementConnectionPool) Execute(statement string, parameters ...interface{}) (uint64, error) {
	return this.executor.Execute(statement, parameters...)
}
func (this *SplitStatementConnectionPool) ExecuteIdentity(statement string, parameters ...interface{}) (uint64, uint64, error) {
	panic("not implemented")
}

func (this *SplitStatementConnectionPool) Select(statement string, parameters ...interface{}) (SelectResult, error) {
	return this.inner.Select(statement, parameters...)
}
