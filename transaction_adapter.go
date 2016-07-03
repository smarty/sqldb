package sqldb

import "database/sql"

type TransactionAdapter struct {
	actual *sql.Tx
}

func NewTransactionAdapter(actual *sql.Tx) *TransactionAdapter {
	return &TransactionAdapter{actual: actual}
}

func (this *TransactionAdapter) Commit() error {
	return this.actual.Commit()
}

func (this *TransactionAdapter) Rollback() error {
	return this.actual.Rollback()
}

func (this *TransactionAdapter) Execute(query string, parameters ...interface{}) (uint64, error) {
	if result, err := this.actual.Exec(query, parameters...); err != nil {
		return 0, err
	} else {
		count, _ := result.RowsAffected()
		return uint64(count), nil
	}
}

func (this *TransactionAdapter) Select(query string, parameters ...interface{}) (SelectResult, error) {
	return this.actual.Query(query, parameters...)
}
