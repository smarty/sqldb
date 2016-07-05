package sqldb

type BindingTransactionAdapter struct {
	actual   Transaction
	selector BindingSelector
}

func NewBindingTransactionAdapter(actual Transaction, panicOnBindError bool) *BindingTransactionAdapter {
	return &BindingTransactionAdapter{
		actual:   actual,
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
	return this.actual.Execute(statement, parameters...)
}

func (this *BindingTransactionAdapter) BindSelect(binder Binder, statement string, parameters ...interface{}) error {
	return this.selector.BindSelect(binder, statement, parameters...)
}
