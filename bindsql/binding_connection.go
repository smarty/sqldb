package bindsql

import "github.com/smartystreets/sqldb"

type BindingConnection struct {
	actual             sqldb.ConnectionPool
	executor           Executor
	selector           Selector
	parameterDelimiter string
	panicOnBindError   bool
}

func NewDefaultConnection(actual sqldb.ConnectionPool) *BindingConnection {
	return NewBindingConnection(actual, "?", true)
}
func NewBindingConnection(actual sqldb.ConnectionPool, parameterDelimiter string, panicOnBindError bool) *BindingConnection {
	return &BindingConnection{
		actual:             actual,
		executor:           sqldb.NewSplitStatementExecutor(actual, parameterDelimiter),
		selector:           NewBindingSelector(actual, panicOnBindError),
		parameterDelimiter: parameterDelimiter,
		panicOnBindError:   panicOnBindError,
	}
}

func (this *BindingConnection) Ping() error {
	return this.actual.Ping()
}
func (this *BindingConnection) BeginTransaction() (Transaction, error) {
	if tx, err := this.actual.BeginTransaction(); err == nil {
		return NewTransaction(tx, this.parameterDelimiter, this.panicOnBindError), nil
	} else {
		return nil, err
	}
}
func (this *BindingConnection) Close() error {
	return this.actual.Close()
}

func (this *BindingConnection) Execute(statement string, parameters ...interface{}) (uint64, error) {
	return this.executor.Execute(statement, parameters...)
}

func (this *BindingConnection) Select(binder Binder, statement string, parameters ...interface{}) error {
	return this.selector.Select(binder, statement, parameters...)
}
