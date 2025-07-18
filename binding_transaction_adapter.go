package sqldb

import "context"

type BindingTransactionAdapter struct {
	Transaction
	selector BindingSelector
}

func NewBindingTransactionAdapter(actual Transaction, panicOnBindError bool) *BindingTransactionAdapter {
	return &BindingTransactionAdapter{
		Transaction: actual,
		selector:    NewBindingSelectorAdapter(actual, panicOnBindError),
	}
}

func (this *BindingTransactionAdapter) BindSelect(ctx context.Context, binder Binder, statement string, parameters ...any) error {
	return this.selector.BindSelect(ctx, binder, statement, parameters...)
}
