package sqldb

import "context"

type SplitStatementTransaction struct {
	Transaction
	executor *SplitStatementExecutor
}

func NewSplitStatementTransaction(inner Transaction, delimiter string) *SplitStatementTransaction {
	return &SplitStatementTransaction{
		Transaction: inner,
		executor:    NewSplitStatementExecutor(inner, delimiter),
	}
}

func (this *SplitStatementTransaction) Execute(ctx context.Context, statement string, parameters ...any) (uint64, error) {
	return this.executor.Execute(ctx, statement, parameters...)
}
