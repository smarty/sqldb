package sqldb

type BindingConnectionPoolAdapter struct {
	actual           ConnectionPool
	selector         BindingSelector
	panicOnBindError bool
}

func NewBindingConnectionPoolAdapter(actual ConnectionPool, panicOnBindError bool) *BindingConnectionPoolAdapter {
	return &BindingConnectionPoolAdapter{
		actual:           actual,
		selector:         NewBindingSelectorAdapter(actual, panicOnBindError),
		panicOnBindError: panicOnBindError,
	}
}

func (this *BindingConnectionPoolAdapter) Ping() error {
	return this.actual.Ping()
}
func (this *BindingConnectionPoolAdapter) BeginTransaction() (BindingTransaction, error) {
	if tx, err := this.actual.BeginTransaction(); err == nil {
		return NewBindingTransactionAdapter(tx, this.panicOnBindError), nil
	} else {
		return nil, err
	}
}
func (this *BindingConnectionPoolAdapter) Close() error {
	return this.actual.Close()
}

func (this *BindingConnectionPoolAdapter) Execute(statement string, parameters ...interface{}) (uint64, error) {
	return this.actual.Execute(statement, parameters...)
}

func (this *BindingConnectionPoolAdapter) BindSelect(binder Binder, statement string, parameters ...interface{}) error {
	return this.selector.BindSelect(binder, statement, parameters...)
}
