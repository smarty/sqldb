package sqldb

import (
	"errors"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/clock"
	"github.com/smartystreets/gunit"
)

func TestRetryBindingConnectionPoolFixture(t *testing.T) {
	gunit.Run(new(RetryBindingConnectionPoolFixture), t)
}

type RetryBindingConnectionPoolFixture struct {
	*gunit.Fixture

	sleeper *clock.Sleeper
	inner   *FakeBindingConnectionPool
	pool    *RetryBindingConnectionPool
}

func (this *RetryBindingConnectionPoolFixture) Setup() {
	this.sleeper = clock.StayAwake()
	this.inner = &FakeBindingConnectionPool{}
	this.pool = NewRetryBindingConnectionPool(this.inner, time.Second)
	this.pool.selector.sleep = this.sleeper
}

///////////////////////////////////////////////////////////////

func (this *RetryBindingConnectionPoolFixture) TestPing() {
	this.inner.pingError = errors.New("")

	err := this.pool.Ping()

	this.So(err, should.Equal, this.inner.pingError)
	this.So(this.inner.pingCalls, should.Equal, 1)
}

func (this *RetryBindingConnectionPoolFixture) TestBeginTransaction() {
	this.inner.transaction = &FakeBindingTransaction{}

	transaction, err := this.pool.BeginTransaction()

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

	affected, err := this.pool.Execute("statement", 1, 2, 3)

	this.So(affected, should.Equal, 42)
	this.So(err, should.Equal, this.inner.executeError)
	this.So(this.inner.executeCalls, should.Equal, 1)
	this.So(this.inner.executeParameters, should.Resemble, []interface{}{1, 2, 3})
}

func (this *RetryBindingConnectionPoolFixture) TestBindSelect() {
	this.inner.selectResult = &FakeSelectResult{}

	err := this.pool.BindSelect(func(Scanner) error {
		return nil
	}, "query", 1, 2, 3)

	this.So(err, should.BeNil)
	this.So(this.inner.selectBinder, should.NotBeNil)
	this.So(this.inner.selectCalls, should.Equal, 1)
	this.So(this.inner.selectStatement, should.Equal, "query")
	this.So(this.inner.selectParameters, should.Resemble, []interface{}{1, 2, 3})
}
