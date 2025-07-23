package sqldb

import (
	"context"
	"database/sql"
)

type LibraryConnectionPoolAdapter struct {
	*sql.DB
	txOptions  *sql.TxOptions
	statements map[string]*sql.Stmt
}

func NewLibraryConnectionPoolAdapter(actual *sql.DB, txOptions *sql.TxOptions) *LibraryConnectionPoolAdapter {
	return &LibraryConnectionPoolAdapter{
		DB:         actual,
		txOptions:  txOptions,
		statements: make(map[string]*sql.Stmt),
	}
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

func (this *LibraryConnectionPoolAdapter) ExecuteStatement(ctx context.Context, id, statement string, parameters ...any) (uint64, error) {
	preparedStatement, ok := this.statements[id]
	if !ok {
		preparedStatement, err := this.DB.PrepareContext(ctx, statement)
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

func (this *LibraryConnectionPoolAdapter) CloseStatement(id string) {
	preparedStatement, ok := this.statements[id]
	if ok {
		_ = preparedStatement.Close()
		delete(this.statements, id)
	}
}
