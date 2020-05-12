package sqldb

import (
	"context"
	"errors"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestSplitStatementTransactionFixture(t *testing.T) {
	gunit.Run(new(SplitStatementTransactionFixture), t)
}

type SplitStatementTransactionFixture struct {
	*gunit.Fixture

	inner       *FakeTransaction
	transaction *SplitStatementTransaction
}

func (this *SplitStatementTransactionFixture) Setup() {
	this.inner = &FakeTransaction{}
	this.transaction = NewSplitStatementTransaction(this.inner, "?")
}

///////////////////////////////////////////////////////////////

func (this *SplitStatementTransactionFixture) TestCommit() {
	this.inner.commitError = errors.New("")

	err := this.transaction.Commit()

	this.So(err, should.Equal, this.inner.commitError)
	this.So(this.inner.commitCalls, should.Equal, 1)
}

func (this *SplitStatementTransactionFixture) TestRollback() {
	this.inner.rollbackError = errors.New("")

	err := this.transaction.Rollback()

	this.So(err, should.Equal, this.inner.rollbackError)
	this.So(this.inner.rollbackCalls, should.Equal, 1)
}

func (this *SplitStatementTransactionFixture) TestSelect() {
	this.inner.selectError = errors.New("")
	this.inner.selectResult = &FakeSelectResult{}

	result, err := this.transaction.Select(context.Background(), "query", 1, 2, 3)

	this.So(result, should.Equal, this.inner.selectResult)
	this.So(err, should.Equal, this.inner.selectError)
	this.So(this.inner.selectCalls, should.Equal, 1)
	this.So(this.inner.selectStatement, should.Equal, "query")
	this.So(this.inner.selectParameters, should.Resemble, []interface{}{1, 2, 3})
}

func (this *SplitStatementTransactionFixture) TestExecute() {
	this.inner.executeResult = 5

	affected, err := this.transaction.Execute(context.Background(), "statement1 ?; statement2 ? ?;", 1, 2, 3)

	this.So(affected, should.Equal, 10)
	this.So(err, should.BeNil)
	this.So(this.inner.executeCalls, should.Equal, 2)
	this.So(this.inner.executeParameters, should.Resemble, []interface{}{2, 3}) // last two parameters
}
