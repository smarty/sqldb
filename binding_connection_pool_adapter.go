package sqldb

type BindingConnectionPoolAdapter struct {
	actual             ConnectionPool
	executor           Executor
	selector           BindingSelector
	parameterDelimiter string
	panicOnBindError   bool
}

func NewDefaultBindingConnectionPoolAdapter(actual ConnectionPool) *BindingConnectionPoolAdapter {
	return NewBindingConnectionPoolAdapter(actual, "?", true)
}
func NewBindingConnectionPoolAdapter(actual ConnectionPool, parameterDelimiter string, panicOnBindError bool) *BindingConnectionPoolAdapter {
	return &BindingConnectionPoolAdapter{
		actual:             actual,
		executor:           NewSplitStatementExecutor(actual, parameterDelimiter),
		selector:           NewBindingSelectorAdapter(actual, panicOnBindError),
		parameterDelimiter: parameterDelimiter,
		panicOnBindError:   panicOnBindError,
	}
}

func (this *BindingConnectionPoolAdapter) Ping() error {
	return this.actual.Ping()
}
func (this *BindingConnectionPoolAdapter) BeginTransaction() (BindingTransaction, error) {
	if tx, err := this.actual.BeginTransaction(); err == nil {
		return NewBindingTransactionAdapter(tx, this.parameterDelimiter, this.panicOnBindError), nil
	} else {
		return nil, err
	}
}
func (this *BindingConnectionPoolAdapter) Close() error {
	return this.actual.Close()
}

func (this *BindingConnectionPoolAdapter) Execute(statement string, parameters ...interface{}) (uint64, error) {
	return this.executor.Execute(statement, parameters...)
}

func (this *BindingConnectionPoolAdapter) Select(binder Binder, statement string, parameters ...interface{}) error {
	return this.selector.Select(binder, statement, parameters...)
}
