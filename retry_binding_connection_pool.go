package sqldb

import "time"

type RetryBindingConnectionPool struct {
	inner    BindingConnectionPool
	selector *RetryBindingSelector
}

func NewRetryBindingConnectionPool(inner BindingConnectionPool, sleep time.Duration) *RetryBindingConnectionPool {
	return &RetryBindingConnectionPool{
		inner:    inner,
		selector: NewRetryBindingSelector(inner, sleep),
	}
}

func (this *RetryBindingConnectionPool) Ping() error {
	return this.inner.Ping()
}

func (this *RetryBindingConnectionPool) BeginTransaction() (BindingTransaction, error) {
	return this.inner.BeginTransaction()
}

func (this *RetryBindingConnectionPool) Close() error {
	return this.inner.Close()
}

func (this *RetryBindingConnectionPool) Execute(statement string, parameters ...interface{}) (uint64, error) {
	return this.inner.Execute(statement, parameters...)
}
func (this *RetryBindingConnectionPool) ExecuteIdentity(statement string, parameters ...interface{}) (uint64, uint64, error) {
	return this.inner.ExecuteIdentity(statement, parameters...)
}

func (this *RetryBindingConnectionPool) BindSelect(binder Binder, statement string, parameters ...interface{}) error {
	return this.selector.BindSelect(binder, statement, parameters...)
}
