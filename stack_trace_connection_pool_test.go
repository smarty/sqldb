package sqldb

import (
	"errors"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestStackTraceConnectionPoolFixture(t *testing.T) {
	gunit.Run(new(StackTraceConnectionPoolFixture), t)
}

type StackTraceConnectionPoolFixture struct {
	*gunit.Fixture

	pool    *FakeConnectionPool
	adapter *StackTraceConnectionPool
}

func (this *StackTraceConnectionPoolFixture) Setup() {
	this.pool = &FakeConnectionPool{}
	this.adapter = NewStackTraceConnectionPool(this.pool)
	this.adapter.StackTrace = ContrivedStackTrace("HELLO, WORLD!")
}

func (this *StackTraceConnectionPoolFixture) TestPing_WhenSuccessful_NoStackTraceIncluded() {
	err := this.adapter.Ping()

	this.So(err, should.BeNil)
	this.So(this.pool.pingCalls, should.Equal, 1)
}

func (this *StackTraceConnectionPoolFixture) TestPing_WhenFails_StackTraceAppendedToErr() {
	this.pool.pingError = errors.New("PING ERROR")

	err := this.adapter.Ping()

	this.So(this.pool.pingCalls, should.Equal, 1)
	this.So(err, should.NotBeNil)
	this.So(err.Error(), should.Equal, "PING ERROR\nStack Trace:\nHELLO, WORLD!")
}

func (this *StackTraceConnectionPoolFixture) TestClose_WhenSuccessful_NoStackTraceIncluded() {
	err := this.adapter.Close()

	this.So(err, should.BeNil)
	this.So(this.pool.closeCalls, should.Equal, 1)
}

func (this *StackTraceConnectionPoolFixture) TestClose_WhenFails_StackTraceAppendedToErr() {
	this.pool.closeError = errors.New("CLOSE ERROR")

	err := this.adapter.Close()

	this.So(this.pool.closeCalls, should.Equal, 1)
	this.So(err, should.NotBeNil)
	this.So(err.Error(), should.Equal, "CLOSE ERROR\nStack Trace:\nHELLO, WORLD!")
}

func (this *StackTraceConnectionPoolFixture) TestBeginTransaction_WhenSuccessful_NoStackTraceIncluded() {
	transaction := new(FakeTransaction)
	this.pool.transaction = transaction

	tx, err := this.adapter.BeginTransaction()

	this.So(err, should.BeNil)
	this.So(this.pool.transactionCalls, should.Equal, 1)
	this.So(tx.(*StackTraceTransaction).inner, should.Equal, transaction)
}

func (this *StackTraceConnectionPoolFixture) TestBeginTransaction_WhenFails_StackTraceAppendedToErr() {
	transaction := new(FakeTransaction)
	this.pool.transaction = transaction
	this.pool.transactionError = errors.New("TX ERROR")

	tx, err := this.adapter.BeginTransaction()

	this.So(this.pool.transactionCalls, should.Equal, 1)
	this.So(tx, should.BeNil)
	this.So(err, should.NotBeNil)
	this.So(err.Error(), should.Equal, "TX ERROR\nStack Trace:\nHELLO, WORLD!")
}

func (this *StackTraceConnectionPoolFixture) TestExecute_WhenSuccessful_NoStackTraceIncluded() {
	this.pool.executeResult = 42

	result, err := this.adapter.Execute("QUERY", 1, 2, 3)

	this.So(result, should.Equal, 42)
	this.So(err, should.BeNil)
	this.So(this.pool.executeCalls, should.Equal, 1)
	this.So(this.pool.executeStatement, should.Equal, "QUERY")
	this.So(this.pool.executeParameters, should.Resemble, []interface{}{1, 2, 3})
}

func (this *StackTraceConnectionPoolFixture) TestExecute_WhenFails_StackTraceAppendedToErr() {
	this.pool.executeError = errors.New("EXECUTE ERROR")
	this.pool.executeResult = 42

	result, err := this.adapter.Execute("QUERY", 1, 2, 3)

	this.So(result, should.Equal, 42)
	this.So(err, should.NotBeNil)
	this.So(err.Error(), should.Equal, "EXECUTE ERROR\nStack Trace:\nHELLO, WORLD!")
	this.So(this.pool.executeCalls, should.Equal, 1)
	this.So(this.pool.executeStatement, should.Equal, "QUERY")
	this.So(this.pool.executeParameters, should.Resemble, []interface{}{1, 2, 3})
}

func (this *StackTraceConnectionPoolFixture) TestSelect_WhenSuccessful_NoStackTraceIncluded() {
	expectedResult := new(FakeSelectResult)
	this.pool.selectResult = expectedResult

	result, err := this.adapter.Select("QUERY", 1, 2, 3)

	this.So(result, should.Equal, expectedResult)
	this.So(err, should.BeNil)
	this.So(this.pool.selectCalls, should.Equal, 1)
	this.So(this.pool.selectStatement, should.Equal, "QUERY")
	this.So(this.pool.selectParameters, should.Resemble, []interface{}{1, 2, 3})
}

func (this *StackTraceConnectionPoolFixture) TestSelect_WhenFails_StackTraceAppendedToErr() {
	expectedResult := new(FakeSelectResult)
	this.pool.selectResult = expectedResult
	this.pool.selectError = errors.New("SELECT ERROR")

	result, err := this.adapter.Select("QUERY", 1, 2, 3)

	this.So(result, should.Equal, expectedResult)
	this.So(err, should.NotBeNil)
	this.So(err.Error(), should.Equal, "SELECT ERROR\nStack Trace:\nHELLO, WORLD!")
	this.So(this.pool.selectCalls, should.Equal, 1)
	this.So(this.pool.selectStatement, should.Equal, "QUERY")
	this.So(this.pool.selectParameters, should.Resemble, []interface{}{1, 2, 3})
}
