package sqldb

import (
	"context"
	"database/sql"
)

type LibraryTransactionAdapter struct {
	*sql.Tx
}

func NewLibraryTransactionAdapter(actual *sql.Tx) *LibraryTransactionAdapter {
	return &LibraryTransactionAdapter{Tx: actual}
}

func (this *LibraryTransactionAdapter) Execute(ctx context.Context, query string, parameters ...any) (uint64, error) {
	if result, err := this.Tx.ExecContext(ctx, query, parameters...); err != nil {
		return 0, err
	} else {
		count, _ := result.RowsAffected()
		return uint64(count), nil
	}
}

func (this *LibraryTransactionAdapter) Select(ctx context.Context, query string, parameters ...any) (SelectResult, error) {
	return this.Tx.QueryContext(ctx, query, parameters...)
}
