package bindsql

import "github.com/smartystreets/sqldb"

type BindingTransaction struct {
	actual   sqldb.Transaction
	executor Executor
	selector Selector
}

func NewDefaultTransaction(actual sqldb.Transaction) *BindingTransaction {
	return NewTransaction(actual, "?", true)
}

func NewTransaction(actual sqldb.Transaction, parameterDelimiter string, panicOnBindError bool) *BindingTransaction {
	return &BindingTransaction{
		actual:   actual,
		executor: sqldb.NewSplitStatementExecutor(actual, parameterDelimiter),
		selector: NewBindingSelector(actual, panicOnBindError),
	}
}

func (this *BindingTransaction) Commit() error {
	return this.actual.Commit()
}
func (this *BindingTransaction) Rollback() error {
	return this.actual.Rollback()
}

func (this *BindingTransaction) Execute(statement string, parameters ...interface{}) (uint64, error) {
	return this.executor.Execute(statement, parameters...)
}

func (this *BindingTransaction) Select(binder Binder, statement string, parameters ...interface{}) error {
	return this.selector.Select(binder, statement, parameters...)
}
