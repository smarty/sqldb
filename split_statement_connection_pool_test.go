package sqldb

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestSplitStatementConnectionPoolFixture(t *testing.T) {
	gunit.Run(new(SplitStatementConnectionPoolFixture), t)
}

type SplitStatementConnectionPoolFixture struct {
	*gunit.Fixture

	inner *FakeConnectionPool
	pool  *SplitStatementConnectionPool
}

func (this *SplitStatementConnectionPoolFixture) Setup() {
	this.inner = &FakeConnectionPool{}
	this.pool = NewSplitStatementConnectionPool(this.inner, "?")
}

///////////////////////////////////////////////////////////////

func (this *SplitStatementConnectionPoolFixture) TestPing() {
	this.inner.pingError = errors.New("")

	err := this.pool.Ping(context.Background())

	this.So(err, should.Equal, this.inner.pingError)
	this.So(this.inner.pingCalls, should.Equal, 1)
}

func (this *SplitStatementConnectionPoolFixture) TestBeginTransactionFails() {
	this.inner.transactionError = errors.New("")

	transaction, err := this.pool.BeginTransaction(context.Background())

	this.So(transaction, should.BeNil)
	this.So(err, should.Equal, this.inner.transactionError)
	this.So(this.inner.transactionCalls, should.Equal, 1)
}

func (this *SplitStatementConnectionPoolFixture) TestBeginTransactionSucceeds() {
	this.inner.transaction = &FakeTransaction{}

	transaction, err := this.pool.BeginTransaction(context.Background())

	this.So(reflect.TypeOf(transaction), should.Equal, reflect.TypeOf(&SplitStatementTransaction{}))
	this.So(err, should.BeNil)
	this.So(this.inner.transactionCalls, should.Equal, 1)
}

func (this *SplitStatementConnectionPoolFixture) TestClose() {
	this.inner.closeError = errors.New("")

	err := this.pool.Close()

	this.So(err, should.Equal, this.inner.closeError)
	this.So(this.inner.closeCalls, should.Equal, 1)
}

func (this *SplitStatementConnectionPoolFixture) TestExecute() {
	this.inner.executeResult = 5

	affected, err := this.pool.Execute(context.Background(), "statement1 ?; statement2 ? ?;", 1, 2, 3)

	this.So(affected, should.Equal, 10)
	this.So(err, should.BeNil)
	this.So(this.inner.executeCalls, should.Equal, 2)
	this.So(this.inner.executeParameters, should.Resemble, []interface{}{2, 3}) // last two parameters
}

func (this *SplitStatementConnectionPoolFixture) TestSelect() {
	this.inner.selectError = errors.New("")
	this.inner.selectResult = &FakeSelectResult{}

	result, err := this.pool.Select(context.Background(), "query", 1, 2, 3)

	this.So(result, should.Equal, this.inner.selectResult)
	this.So(err, should.Equal, this.inner.selectError)
	this.So(this.inner.selectCalls, should.Equal, 1)
	this.So(this.inner.selectStatement, should.Equal, "query")
	this.So(this.inner.selectParameters, should.Resemble, []interface{}{1, 2, 3})
}
