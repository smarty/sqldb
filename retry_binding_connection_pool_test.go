package sqldb

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestRetryBindingConnectionPoolFixture(t *testing.T) {
	gunit.Run(new(RetryBindingConnectionPoolFixture), t)
}

type RetryBindingConnectionPoolFixture struct {
	*gunit.Fixture

	inner *FakeBindingConnectionPool
	pool  *RetryBindingConnectionPool
}

func (this *RetryBindingConnectionPoolFixture) Setup() {
	this.inner = &FakeBindingConnectionPool{}
	this.pool = NewRetryBindingConnectionPool(this.inner, time.Second)
}

///////////////////////////////////////////////////////////////

func (this *RetryBindingConnectionPoolFixture) TestPing() {
	this.inner.pingError = errors.New("")

	err := this.pool.Ping(context.Background())

	this.So(err, should.Equal, this.inner.pingError)
	this.So(this.inner.pingCalls, should.Equal, 1)
}

func (this *RetryBindingConnectionPoolFixture) TestBeginTransaction() {
	this.inner.transaction = &FakeBindingTransaction{}

	transaction, err := this.pool.BeginTransaction(context.Background())

	this.So(transaction, should.Equal, this.inner.transaction)
	this.So(err, should.BeNil)
	this.So(this.inner.transactionCalls, should.Equal, 1)
}

func (this *RetryBindingConnectionPoolFixture) TestClose() {
	this.inner.closeError = errors.New("")

	err := this.pool.Close()

	this.So(err, should.Equal, this.inner.closeError)
	this.So(this.inner.closeCalls, should.Equal, 1)
}

func (this *RetryBindingConnectionPoolFixture) TestExecute() {
	this.inner.executeResult = 42
	this.inner.executeError = errors.New("")

	affected, err := this.pool.Execute(context.Background(), "statement", 1, 2, 3)

	this.So(affected, should.Equal, 42)
	this.So(err, should.Equal, this.inner.executeError)
	this.So(this.inner.executeCalls, should.Equal, 1)
	this.So(this.inner.executeParameters, should.Resemble, []any{1, 2, 3})
}

func (this *RetryBindingConnectionPoolFixture) TestBindSelect() {
	this.inner.selectResult = &FakeSelectResult{}

	err := this.pool.BindSelect(context.Background(), func(Scanner) error {
		return nil
	}, "query", 1, 2, 3)

	this.So(err, should.BeNil)
	this.So(this.inner.selectBinder, should.NotBeNil)
	this.So(this.inner.selectCalls, should.Equal, 1)
	this.So(this.inner.selectStatement, should.Equal, "query")
	this.So(this.inner.selectParameters, should.Resemble, []any{1, 2, 3})
}
