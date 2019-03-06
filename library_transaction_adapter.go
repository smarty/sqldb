package sqldb

import "database/sql"

type LibraryTransactionAdapter struct {
	inner *sql.Tx
}

func NewLibraryTransactionAdapter(actual *sql.Tx) *LibraryTransactionAdapter {
	return &LibraryTransactionAdapter{inner: actual}
}

func (this *LibraryTransactionAdapter) Commit() error {
	return this.inner.Commit()
}
func (this *LibraryTransactionAdapter) Rollback() error {
	return this.inner.Rollback()
}

func (this *LibraryTransactionAdapter) Execute(query string, parameters ...interface{}) (uint64, error) {
	if result, err := this.inner.Exec(query, parameters...); err != nil {
		return 0, err
	} else {
		count, _ := result.RowsAffected()
		return uint64(count), nil
	}
}
func (this *LibraryTransactionAdapter) ExecuteIdentity(query string, parameters ...interface{}) (uint64, uint64, error) {
	if result, err := this.inner.Exec(query, parameters...); err != nil {
		return 0, 0, err
	} else {
		count, _ := result.RowsAffected()
		identity, _ := result.LastInsertId()
		return uint64(count), uint64(identity), nil
	}
}

func (this *LibraryTransactionAdapter) Select(query string, parameters ...interface{}) (SelectResult, error) {
	return this.inner.Query(query, parameters...)
}
