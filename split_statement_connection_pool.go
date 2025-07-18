package sqldb

import "context"

type SplitStatementConnectionPool struct {
	ConnectionPool
	delimiter string
	executor  *SplitStatementExecutor
}

func NewSplitStatementConnectionPool(inner ConnectionPool, delimiter string) *SplitStatementConnectionPool {
	return &SplitStatementConnectionPool{
		ConnectionPool: inner,
		delimiter:      delimiter,
		executor:       NewSplitStatementExecutor(inner, delimiter),
	}
}

func (this *SplitStatementConnectionPool) BeginTransaction(ctx context.Context) (Transaction, error) {
	if transaction, err := this.ConnectionPool.BeginTransaction(ctx); err == nil {
		return NewSplitStatementTransaction(transaction, this.delimiter), nil
	} else {
		return nil, err
	}
}

func (this *SplitStatementConnectionPool) Execute(ctx context.Context, statement string, parameters ...any) (uint64, error) {
	return this.executor.Execute(ctx, statement, parameters...)
}
