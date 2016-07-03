package sqldb

type BindingTransactionAdapter struct {
	actual   Transaction
	executor Executor
	selector BindingSelector
}

func NewDefaultBindingTransactionAdapter(actual Transaction) *BindingTransactionAdapter {
	return NewBindingTransactionAdapter(actual, "?", true)
}
func NewBindingTransactionAdapter(actual Transaction, parameterDelimiter string, panicOnBindError bool) *BindingTransactionAdapter {
	return &BindingTransactionAdapter{
		actual:   actual,
		executor: NewSplitStatementExecutor(actual, parameterDelimiter),
		selector: NewBindingSelectorAdapter(actual, panicOnBindError),
	}
}

func (this *BindingTransactionAdapter) Commit() error {
	return this.actual.Commit()
}
func (this *BindingTransactionAdapter) Rollback() error {
	return this.actual.Rollback()
}

func (this *BindingTransactionAdapter) Execute(statement string, parameters ...interface{}) (uint64, error) {
	return this.executor.Execute(statement, parameters...)
}

func (this *BindingTransactionAdapter) Select(statement string, parameters ...interface{}) (SelectResult, error) {
	return this.actual.Select(statement, parameters...)
}
func (this *BindingTransactionAdapter) BindSelect(binder Binder, statement string, parameters ...interface{}) error {
	return this.selector.BindSelect(binder, statement, parameters...)
}
