package bindsql

import "github.com/smartystreets/sqldb"

type BindingConnectionAdapter struct {
	actual             sqldb.ConnectionPool
	executor           Executor
	selector           Selector
	parameterDelimiter string
	panicOnBindError   bool
}

func NewDefaultBindingConnection(actual sqldb.ConnectionPool) *BindingConnectionAdapter {
	return NewBindingConnectionAdapter(actual, "?", true)
}
func NewBindingConnectionAdapter(actual sqldb.ConnectionPool, parameterDelimiter string, panicOnBindError bool) *BindingConnectionAdapter {
	return &BindingConnectionAdapter{
		actual:             actual,
		executor:           sqldb.NewSplitStatementExecutor(actual, parameterDelimiter),
		selector:           NewBindingSelectorAdapter(actual, panicOnBindError),
		parameterDelimiter: parameterDelimiter,
		panicOnBindError:   panicOnBindError,
	}
}

func (this *BindingConnectionAdapter) Ping() error {
	return this.actual.Ping()
}
func (this *BindingConnectionAdapter) BeginTransaction() (BindingTransaction, error) {
	if tx, err := this.actual.BeginTransaction(); err == nil {
		return NewBindingTransaction(tx, this.parameterDelimiter, this.panicOnBindError), nil
	} else {
		return nil, err
	}
}
func (this *BindingConnectionAdapter) Close() error {
	return this.actual.Close()
}

func (this *BindingConnectionAdapter) Execute(statement string, parameters ...interface{}) (uint64, error) {
	return this.executor.Execute(statement, parameters...)
}

func (this *BindingConnectionAdapter) Select(binder Binder, statement string, parameters ...interface{}) error {
	return this.selector.Select(binder, statement, parameters...)
}
