package sqldb

import (
	"errors"
	"runtime/debug"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestStackTraceConnectionPoolAdapterFixture(t *testing.T) {
	gunit.Run(new(StackTraceConnectionPoolAdapterFixture), t)
}

type StackTraceConnectionPoolAdapterFixture struct {
	*gunit.Fixture

	pool    *FakeConnectionPool
	adapter *StackTraceConnectionPoolAdapter
}

func (this *StackTraceConnectionPoolAdapterFixture) Setup() {
	this.pool = &FakeConnectionPool{}
	this.adapter = NewStackTraceConnectionPoolAdapter(this.pool)
	this.adapter.stack = ContrivedStackTrace("HELLO, WORLD!")
}

func (this *StackTraceConnectionPoolAdapterFixture) TestPing_WhenSuccessful_NoStackTraceIncluded() {
	err := this.adapter.Ping()

	this.So(err, should.BeNil)
	this.So(this.pool.pingCalls, should.Equal, 1)
}

func (this *StackTraceConnectionPoolAdapterFixture) TestPing_WhenFails_StackTraceAppendedToErr() {
	this.pool.pingError = errors.New("PING ERROR")

	err := this.adapter.Ping()

	this.So(this.pool.pingCalls, should.Equal, 1)
	this.So(err, should.NotBeNil)
	this.So(err.Error(), should.Equal, "PING ERROR\nStack Trace:\nHELLO, WORLD!")
}

func (this *StackTraceConnectionPoolAdapterFixture) TestClose_WhenSuccessful_NoStackTraceIncluded() {
	err := this.adapter.Close()

	this.So(err, should.BeNil)
	this.So(this.pool.closeCalls, should.Equal, 1)
}

func (this *StackTraceConnectionPoolAdapterFixture) TestClose_WhenFails_StackTraceAppendedToErr() {
	this.pool.closeError = errors.New("CLOSE ERROR")

	err := this.adapter.Close()

	this.So(this.pool.closeCalls, should.Equal, 1)
	this.So(err, should.NotBeNil)
	this.So(err.Error(), should.Equal, "CLOSE ERROR\nStack Trace:\nHELLO, WORLD!")
}

func (this *StackTraceConnectionPoolAdapterFixture) TestBeginTransaction_WhenSuccessful_NoStackTraceIncluded() {
	transaction := new(FakeTransaction)
	this.pool.transaction = transaction

	tx, err := this.adapter.BeginTransaction()

	this.So(err, should.BeNil)
	this.So(this.pool.transactionCalls, should.Equal, 1)
	this.So(tx, should.Equal, transaction)
}

func (this *StackTraceConnectionPoolAdapterFixture) TestBeginTransaction_WhenFails_StackTraceAppendedToErr() {
	transaction := new(FakeTransaction)
	this.pool.transaction = transaction
	this.pool.transactionError = errors.New("TX ERROR")

	tx, err := this.adapter.BeginTransaction()

	this.So(this.pool.transactionCalls, should.Equal, 1)
	this.So(tx, should.Equal, transaction)
	this.So(err, should.NotBeNil)
	this.So(err.Error(), should.Equal, "TX ERROR\nStack Trace:\nHELLO, WORLD!")
}

func (this *StackTraceConnectionPoolAdapterFixture) TestExecute_WhenSuccessful_NoStackTraceIncluded() {
	this.pool.executeResult = 42

	result, err := this.adapter.Execute("QUERY", 1, 2, 3)

	this.So(result, should.Equal, 42)
	this.So(err, should.BeNil)
	this.So(this.pool.executeCalls, should.Equal, 1)
	this.So(this.pool.executeStatement, should.Equal, "QUERY")
	this.So(this.pool.executeParameters, should.Resemble, []interface{}{1, 2, 3})
}

func (this *StackTraceConnectionPoolAdapterFixture) TestExecute_WhenFails_StackTraceAppendedToErr() {
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

func (this *StackTraceConnectionPoolAdapterFixture) TestSelect_WhenSuccessful_NoStackTraceIncluded() {
	expectedResult := new(FakeSelectResult)
	this.pool.selectResult = expectedResult

	result, err := this.adapter.Select("QUERY", 1, 2, 3)

	this.So(result, should.Equal, expectedResult)
	this.So(err, should.BeNil)
	this.So(this.pool.selectCalls, should.Equal, 1)
	this.So(this.pool.selectStatement, should.Equal, "QUERY")
	this.So(this.pool.selectParameters, should.Resemble, []interface{}{1, 2, 3})
}

func (this *StackTraceConnectionPoolAdapterFixture) TestSelect_WhenFails_StackTraceAppendedToErr() {
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

/**************************************************************************/

func TestStackTraceFixture(t *testing.T) {
	gunit.Run(new(StackTraceFixture), t)
}

type StackTraceFixture struct {
	*gunit.Fixture

	stack *StackTrace
}

func (this *StackTraceFixture) TestWhenNil_ReturnsActualStackTrace() {
	actual := this.stack.StackTrace()
	expected := string(debug.Stack())

	actual = actual[len(actual)-1000:]       // last 1000 characters
	expected = expected[len(expected)-1000:] // last 1000 characters

	this.So(actual, should.Equal, expected)
}

func (this *StackTraceFixture) TestWhenNonNil_ReturnsPreSetMessage() {
	this.stack = ContrivedStackTrace("HELLO")
	this.So(this.stack.StackTrace(), should.Equal, "HELLO")
}
