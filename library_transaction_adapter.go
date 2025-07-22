package sqldb

import (
	"context"
	"database/sql"
)

type LibraryTransactionAdapter struct {
	*sql.Tx
	statements map[string]*sql.Stmt
}

func NewLibraryTransactionAdapter(actual *sql.Tx) *LibraryTransactionAdapter {
	return &LibraryTransactionAdapter{Tx: actual}
}

func (this *LibraryTransactionAdapter) ExecuteStatement(ctx context.Context, id, statement string, parameters ...any) (uint64, error) {
	preparedStatement, ok := this.statements[id]
	if !ok {
		preparedStatement, err := this.Tx.PrepareContext(ctx, statement)
		if err != nil {
			return 0, err
		}

		this.statements[id] = preparedStatement
	}

	result, err := preparedStatement.ExecContext(ctx, parameters...)
	if err != nil {
		return 0, err
	} else {
		count, _ := result.RowsAffected()
		return uint64(count), nil
	}
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

func (this *LibraryTransactionAdapter) CloseStatement(id string) {
	preparedStatement, ok := this.statements[id]
	if ok {
		_ = preparedStatement.Close()
		delete(this.statements, id)
	}
}
