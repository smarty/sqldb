package sqldb

import "database/sql"

type LibraryTransactionAdapter struct {
	actual *sql.Tx
}

func NewLibraryTransactionAdapter(actual *sql.Tx) *LibraryTransactionAdapter {
	return &LibraryTransactionAdapter{actual: actual}
}

func (this *LibraryTransactionAdapter) Commit() error {
	return this.actual.Commit()
}
func (this *LibraryTransactionAdapter) Rollback() error {
	return this.actual.Rollback()
}

func (this *LibraryTransactionAdapter) Execute(query string, parameters ...interface{}) (uint64, error) {
	if result, err := this.actual.Exec(query, parameters...); err != nil {
		return 0, err
	} else {
		count, _ := result.RowsAffected()
		return uint64(count), nil
	}
}

func (this *LibraryTransactionAdapter) Select(query string, parameters ...interface{}) (SelectResult, error) {
	return this.actual.Query(query, parameters...)
}
