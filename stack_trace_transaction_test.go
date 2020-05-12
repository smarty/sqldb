package sqldb

import (
	"errors"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestStackTraceTransactionFixture(t *testing.T) {
	gunit.Run(new(StackTraceTransactionFixture), t)
}

type StackTraceTransactionFixture struct {
	*gunit.Fixture

	inner       *FakeTransaction
	transaction *StackTraceTransaction
}

func (this *StackTraceTransactionFixture) Setup() {
	this.inner = new(FakeTransaction)
	this.transaction = NewStackTraceTransaction(this.inner)
	this.transaction.StackTrace = ContrivedStackTrace("STACK")
}

func (this *StackTraceTransactionFixture) TestCommit() {
	this.inner.commitError = errors.New("ERROR")

	err := this.transaction.Commit()

	this.So(err.Error(), should.Equal, "ERROR\nStack Trace:\nSTACK")
	this.So(this.inner.commitCalls, should.Equal, 1)
}

func (this *StackTraceTransactionFixture) TestRollback() {
	this.inner.rollbackError = errors.New("ERROR")

	err := this.transaction.Rollback()

	this.So(err.Error(), should.Equal, "ERROR\nStack Trace:\nSTACK")
	this.So(this.inner.rollbackCalls, should.Equal, 1)
}

func (this *StackTraceTransactionFixture) TestExecute() {
	this.inner.executeError = errors.New("ERROR")
	this.inner.executeResult = 42

	rows, err := this.transaction.Execute("STATEMENT", 1, 2, 3)

	this.So(rows, should.Equal, 42)
	this.So(err.Error(), should.Equal, "ERROR\nStack Trace:\nSTACK")
	this.So(this.inner.executeCalls, should.Equal, 1)
	this.So(this.inner.executeStatement, should.Equal, "STATEMENT")
	this.So(this.inner.executeParameters, should.Resemble, []interface{}{1, 2, 3})
}

func (this *StackTraceTransactionFixture) TestSelect() {
	expectedResult := new(FakeSelectResult)
	this.inner.selectResult = expectedResult
	this.inner.selectError = errors.New("ERROR")

	result, err := this.transaction.Select("STATEMENT", 1, 2, 3)

	this.So(result, should.Equal, expectedResult)
	this.So(err.Error(), should.Equal, "ERROR\nStack Trace:\nSTACK")
	this.So(this.inner.selectCalls, should.Equal, 1)
	this.So(this.inner.selectStatement, should.Equal, "STATEMENT")
	this.So(this.inner.selectParameters, should.Resemble, []interface{}{1, 2, 3})
}
