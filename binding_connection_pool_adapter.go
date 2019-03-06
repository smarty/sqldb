package sqldb

type BindingConnectionPoolAdapter struct {
	inner            ConnectionPool
	selector         BindingSelector
	panicOnBindError bool
}

func NewBindingConnectionPoolAdapter(actual ConnectionPool, panicOnBindError bool) *BindingConnectionPoolAdapter {
	return &BindingConnectionPoolAdapter{
		inner:            actual,
		selector:         NewBindingSelectorAdapter(actual, panicOnBindError),
		panicOnBindError: panicOnBindError,
	}
}

func (this *BindingConnectionPoolAdapter) Ping() error {
	return this.inner.Ping()
}
func (this *BindingConnectionPoolAdapter) BeginTransaction() (BindingTransaction, error) {
	if tx, err := this.inner.BeginTransaction(); err == nil {
		return NewBindingTransactionAdapter(tx, this.panicOnBindError), nil
	} else {
		return nil, err
	}
}
func (this *BindingConnectionPoolAdapter) Close() error {
	return this.inner.Close()
}

func (this *BindingConnectionPoolAdapter) Execute(statement string, parameters ...interface{}) (uint64, error) {
	return this.inner.Execute(statement, parameters...)
}
func (this *BindingConnectionPoolAdapter) ExecuteIdentity(statement string, parameters ...interface{}) (uint64, uint64, error) {
	panic("not implemented")
}

func (this *BindingConnectionPoolAdapter) BindSelect(binder Binder, statement string, parameters ...interface{}) error {
	return this.selector.BindSelect(binder, statement, parameters...)
}
