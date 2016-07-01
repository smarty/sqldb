package sqldb

import "database/sql"

type Transaction struct {
	actual *sql.Tx
}

func NewTransaction(actual *sql.Tx) *Transaction {
	return &Transaction{actual: actual}
}

func (this *Transaction) Commit() error {
	return this.actual.Commit()
}

func (this *Transaction) Rollback() error {
	return this.actual.Rollback()
}

func (this *Transaction) Execute(query string, parameters ...interface{}) (uint64, error) {
	if result, err := this.actual.Exec(query, parameters...); err != nil {
		return 0, err
	} else {
		count, _ := result.RowsAffected()
		return uint64(count), nil
	}
}

func (this *Transaction) Select(query string, parameters ...interface{}) (SelectResult, error) {
	return this.actual.Query(query, parameters...)
}
