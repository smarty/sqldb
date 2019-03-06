package sqldb

type StackTraceConnectionPool struct {
	*StackTrace
	inner ConnectionPool
}

func NewStackTraceConnectionPool(inner ConnectionPool) *StackTraceConnectionPool {
	return &StackTraceConnectionPool{inner: inner}
}

func (this *StackTraceConnectionPool) Ping() error {
	return this.Wrap(this.inner.Ping())
}

func (this *StackTraceConnectionPool) BeginTransaction() (Transaction, error) {
	if tx, err := this.inner.BeginTransaction(); err == nil {
		return NewStackTraceTransaction(tx), nil
	} else {
		return nil, this.Wrap(err)
	}
}

func (this *StackTraceConnectionPool) Close() error {
	return this.Wrap(this.inner.Close())
}

func (this *StackTraceConnectionPool) Execute(query string, parameters ...interface{}) (uint64, error) {
	rows, err := this.inner.Execute(query, parameters...)
	return rows, this.Wrap(err)
}
func (this *StackTraceConnectionPool) ExecuteIdentity(statement string, parameters ...interface{}) (uint64, uint64, error) {
	panic("not implemented")
}

func (this *StackTraceConnectionPool) Select(query string, parameters ...interface{}) (SelectResult, error) {
	result, err := this.inner.Select(query, parameters...)
	return result, this.Wrap(err)
}
