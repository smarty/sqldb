package sqldb

type BindingTransactionAdapter struct {
	inner    Transaction
	selector BindingSelector
}

func NewBindingTransactionAdapter(actual Transaction, panicOnBindError bool) *BindingTransactionAdapter {
	return &BindingTransactionAdapter{
		inner:    actual,
		selector: NewBindingSelectorAdapter(actual, panicOnBindError),
	}
}

func (this *BindingTransactionAdapter) Commit() error {
	return this.inner.Commit()
}

func (this *BindingTransactionAdapter) Rollback() error {
	return this.inner.Rollback()
}

func (this *BindingTransactionAdapter) Execute(statement string, parameters ...interface{}) (uint64, error) {
	return this.inner.Execute(statement, parameters...)
}
func (this *BindingTransactionAdapter) ExecuteIdentity(statement string, parameters ...interface{}) (uint64, uint64, error) {
	return this.inner.ExecuteIdentity(statement, parameters...)
}

func (this *BindingTransactionAdapter) BindSelect(binder Binder, statement string, parameters ...interface{}) error {
	return this.selector.BindSelect(binder, statement, parameters...)
}
