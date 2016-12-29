package sqldb

import "fmt"

type StackTraceConnectionPool struct {
	inner ConnectionPool
	stack *StackTrace
}

func NewStackTraceConnectionPool(inner ConnectionPool) *StackTraceConnectionPool {
	return &StackTraceConnectionPool{inner: inner}
}

func (this *StackTraceConnectionPool) Ping() error {
	return this.wrap(this.inner.Ping())
}

func (this *StackTraceConnectionPool) BeginTransaction() (Transaction, error) {
	if tx, err := this.inner.BeginTransaction(); err == nil {
		return NewStackTraceTransaction(tx), nil
	} else {
		return nil, this.wrap(err)
	}
}

func (this *StackTraceConnectionPool) Close() error {
	return this.wrap(this.inner.Close())
}

func (this *StackTraceConnectionPool) Execute(query string, parameters ...interface{}) (uint64, error) {
	rows, err := this.inner.Execute(query, parameters...)
	return rows, this.wrap(err)
}

func (this *StackTraceConnectionPool) Select(query string, parameters ...interface{}) (SelectResult, error) {
	result, err := this.inner.Select(query, parameters...)
	return result, this.wrap(err)
}

func (this *StackTraceConnectionPool) wrap(err error) error {
	if err != nil {
		err = fmt.Errorf("%s\nStack Trace:\n%s", err.Error(), this.stack.StackTrace())
	}
	return err
}
