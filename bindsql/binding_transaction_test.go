package bindsql

import (
	"errors"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/sqldb"
)

type BindingTransactionFixture struct {
	*gunit.Fixture

	transaction *BindingTransaction
	fakeInner   *FakeDriverTransaction
}

func (this *BindingTransactionFixture) Setup() {
	this.fakeInner = &FakeDriverTransaction{}
	this.transaction = NewDefaultTransaction(this.fakeInner)
}

///////////////////////////////////////////////////////////////

func (this *BindingTransactionFixture) TestCommit() {
	this.fakeInner.commitError = errors.New("")

	err := this.transaction.Commit()

	this.So(err, should.Equal, this.fakeInner.commitError)
	this.So(this.fakeInner.commit, should.Equal, 1)
}

func (this *BindingTransactionFixture) TestRollback() {
	this.fakeInner.rollbackError = errors.New("")

	err := this.transaction.Rollback()

	this.So(err, should.Equal, this.fakeInner.rollbackError)
	this.So(this.fakeInner.rollback, should.Equal, 1)
}

func (this *BindingTransactionFixture) TestExecute() {
	this.fakeInner.executeError = errors.New("")

	_, err := this.transaction.Execute("statement;")

	this.So(err, should.Equal, this.fakeInner.executeError)
	this.So(this.fakeInner.executes, should.Resemble, []string{"statement;"})
}

func (this *BindingTransactionFixture) TestMultiStatementExecute() {
	_, err := this.transaction.Execute("statement1;statement2;")

	this.So(err, should.BeNil)
	this.So(this.fakeInner.executes, should.Resemble, []string{"statement1;", "statement2;"})
}

func (this *BindingTransactionFixture) TestSelect() {
	this.fakeInner.queryError = errors.New("")

	err := this.transaction.Select(nil, "query")

	this.So(err, should.Equal, this.fakeInner.queryError)
	this.So(this.fakeInner.queries, should.Resemble, []string{"query"})
}

///////////////////////////////////////////////////////////////

type FakeDriverTransaction struct {
	commit   int
	rollback int

	queries  []string
	executes []string

	commitError   error
	rollbackError error
	executeError  error
	queryError    error
}

func (this *FakeDriverTransaction) Commit() error {
	this.commit++
	return this.commitError
}

func (this *FakeDriverTransaction) Rollback() error {
	this.rollback++
	return this.rollbackError
}

func (this *FakeDriverTransaction) Execute(statement string, parameters ...interface{}) (uint64, error) {
	this.executes = append(this.executes, statement)
	return 0, this.executeError
}

func (this *FakeDriverTransaction) Select(query string, parameters ...interface{}) (sqldb.SelectResult, error) {
	this.queries = append(this.queries, query)
	return nil, this.queryError
}