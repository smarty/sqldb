package sqldb

import (
	"errors"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type SplitStatementTransactionFixture struct {
	*gunit.Fixture

	fakeInner   *FakeInnerTransaction
	transaction *SplitStatementTransaction
}

func (this *SplitStatementTransactionFixture) Setup() {
	this.fakeInner = &FakeInnerTransaction{}
	this.transaction = NewSplitStatementTransaction(this.fakeInner, "?")
}

///////////////////////////////////////////////////////////////

func (this *SplitStatementTransactionFixture) TestCommit() {
	this.fakeInner.commitError = errors.New("")

	err := this.transaction.Commit()

	this.So(err, should.Equal, this.fakeInner.commitError)
	this.So(this.fakeInner.commit, should.Equal, 1)
}

func (this *SplitStatementTransactionFixture) TestRollback() {
	this.fakeInner.rollbackError = errors.New("")

	err := this.transaction.Rollback()

	this.So(err, should.Equal, this.fakeInner.rollbackError)
	this.So(this.fakeInner.rollback, should.Equal, 1)
}

func (this *SplitStatementTransactionFixture) TestSelect() {
	this.fakeInner.selectError = errors.New("")
	this.fakeInner.selectResult = &FakeSelectResult{}

	result, err := this.transaction.Select("query", 1, 2, 3)

	this.So(result, should.Equal, this.fakeInner.selectResult)
	this.So(err, should.Equal, this.fakeInner.selectError)
	this.So(this.fakeInner.selects, should.Equal, 1)
	this.So(this.fakeInner.selectStatement, should.Equal, "query")
	this.So(this.fakeInner.selectParameters, should.Resemble, []interface{}{1, 2, 3})
}

///////////////////////////////////////////////////////////////
