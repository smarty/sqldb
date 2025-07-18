package sqldb

import (
	"context"
	"database/sql"
)

type LibraryConnectionPoolAdapter struct {
	*sql.DB
	txOptions *sql.TxOptions
}

func NewLibraryConnectionPoolAdapter(actual *sql.DB, txOptions *sql.TxOptions) *LibraryConnectionPoolAdapter {
	return &LibraryConnectionPoolAdapter{DB: actual, txOptions: txOptions}
}

func (this *LibraryConnectionPoolAdapter) Ping(ctx context.Context) error {
	return this.DB.PingContext(ctx)
}

func (this *LibraryConnectionPoolAdapter) BeginTransaction(ctx context.Context) (Transaction, error) {
	if tx, err := this.DB.BeginTx(ctx, this.txOptions); err == nil {
		return NewLibraryTransactionAdapter(tx), nil
	} else {
		return nil, err
	}
}

func (this *LibraryConnectionPoolAdapter) Execute(ctx context.Context, query string, parameters ...any) (uint64, error) {
	if result, err := this.DB.ExecContext(ctx, query, parameters...); err != nil {
		return 0, err
	} else {
		count, _ := result.RowsAffected()
		return uint64(count), nil
	}
}

func (this *LibraryConnectionPoolAdapter) Select(ctx context.Context, query string, parameters ...any) (SelectResult, error) {
	return this.DB.QueryContext(ctx, query, parameters...)
}
