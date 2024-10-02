package sqldb

import (
	"context"
	"errors"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestNormalizeContextCancellationConnectionPoolFixture(t *testing.T) {
	gunit.Run(new(NormalizeContextCancellationConnectionPoolFixture), t)
}

type NormalizeContextCancellationConnectionPoolFixture struct {
	*gunit.Fixture

	inner   *FakeConnectionPool
	adapter *NormalizeContextCancellationConnectionPool
}

func (this *NormalizeContextCancellationConnectionPoolFixture) Setup() {
	this.inner = &FakeConnectionPool{}
	this.adapter = NewNormalizeContextCancellationConnectionPool(this.inner)
}

func (this *NormalizeContextCancellationConnectionPoolFixture) TestPing_Successful() {
	err := this.adapter.Ping(context.Background())

	this.So(err, should.BeNil)
	this.So(this.inner.pingCalls, should.Equal, 1)
}
func (this *NormalizeContextCancellationConnectionPoolFixture) TestPing_Failed() {
	pingErr := errors.New("PING ERROR")
	this.inner.pingError = pingErr

	err := this.adapter.Ping(context.Background())

	this.So(this.inner.pingCalls, should.Equal, 1)
	this.So(err, should.Equal, pingErr)
}
func (this *NormalizeContextCancellationConnectionPoolFixture) TestPing_AdaptContextCancelled() {
	this.inner.pingError = operationCanceledErr

	err := this.adapter.Ping(context.Background())

	this.So(this.inner.pingCalls, should.Equal, 1)
	this.So(errors.Is(err, operationCanceledErr), should.BeTrue)
	this.So(errors.Is(err, context.Canceled), should.BeTrue)
}

func (this *NormalizeContextCancellationConnectionPoolFixture) TestBeginTransaction_Successful() {
	transaction := new(FakeTransaction)
	this.inner.transaction = transaction

	tx, err := this.adapter.BeginTransaction(context.Background())

	this.So(err, should.BeNil)
	this.So(this.inner.transactionCalls, should.Equal, 1)
	this.So(tx, should.NotBeNil)
}
func (this *NormalizeContextCancellationConnectionPoolFixture) TestBeginTransaction_Failed() {
	transactionErr := errors.New("BEGIN TRANSACTION ERROR")
	this.inner.transactionError = transactionErr

	tx, err := this.adapter.BeginTransaction(context.Background())

	this.So(tx, should.BeNil)
	this.So(err, should.Equal, transactionErr)
}
func (this *NormalizeContextCancellationConnectionPoolFixture) TestBeginTransaction_AdaptContextCancelled() {
	this.inner.transactionError = operationCanceledErr

	tx, err := this.adapter.BeginTransaction(context.Background())

	this.So(tx, should.BeNil)
	this.So(errors.Is(err, operationCanceledErr), should.BeTrue)
	this.So(errors.Is(err, context.Canceled), should.BeTrue)
}

func (this *NormalizeContextCancellationConnectionPoolFixture) TestClose_Successful() {
	err := this.adapter.Close()

	this.So(err, should.BeNil)
	this.So(this.inner.closeCalls, should.Equal, 1)
}
func (this *NormalizeContextCancellationConnectionPoolFixture) TestClose_Failed() {
	closeErr := errors.New("CLOSE ERROR")
	this.inner.closeError = closeErr

	err := this.adapter.Close()

	this.So(this.inner.closeCalls, should.Equal, 1)
	this.So(err, should.Equal, closeErr)
}
func (this *NormalizeContextCancellationConnectionPoolFixture) TestClose_AdaptContextCancelled() {
	this.inner.closeError = operationCanceledErr

	err := this.adapter.Close()

	this.So(this.inner.closeCalls, should.Equal, 1)
	this.So(errors.Is(err, operationCanceledErr), should.BeTrue)
	this.So(errors.Is(err, context.Canceled), should.BeTrue)
}

func (this *NormalizeContextCancellationConnectionPoolFixture) TestExecute_Successful() {
	this.inner.executeResult = 42

	result, err := this.adapter.Execute(context.Background(), "statement")

	this.So(result, should.Equal, 42)
	this.So(err, should.BeNil)
	this.So(this.inner.executeCalls, should.Equal, 1)
	this.So(this.inner.executeStatement, should.Equal, "statement")
}
func (this *NormalizeContextCancellationConnectionPoolFixture) TestExecute_Failed() {
	this.inner.executeResult = 42
	executeErr := errors.New("EXECUTE ERROR")
	this.inner.executeError = executeErr

	result, err := this.adapter.Execute(context.Background(), "statement")

	this.So(result, should.Equal, 42)
	this.So(err, should.Equal, executeErr)
	this.So(this.inner.executeCalls, should.Equal, 1)
	this.So(this.inner.executeStatement, should.Equal, "statement")
}
func (this *NormalizeContextCancellationConnectionPoolFixture) TestExecute_AdaptContextCancelled() {
	this.inner.executeResult = 42
	this.inner.executeError = operationCanceledErr

	result, err := this.adapter.Execute(context.Background(), "statement")

	this.So(result, should.Equal, 42)
	this.So(this.inner.executeCalls, should.Equal, 1)
	this.So(this.inner.executeStatement, should.Equal, "statement")
	this.So(errors.Is(err, operationCanceledErr), should.BeTrue)
	this.So(errors.Is(err, context.Canceled), should.BeTrue)
}

func (this *NormalizeContextCancellationConnectionPoolFixture) TestSelect_Successful() {
	expectedResult := new(FakeSelectResult)
	this.inner.selectResult = expectedResult

	result, err := this.adapter.Select(context.Background(), "query", 1, 2, 3)

	this.So(result, should.Equal, expectedResult)
	this.So(err, should.BeNil)
	this.So(this.inner.selectCalls, should.Equal, 1)
	this.So(this.inner.selectStatement, should.Equal, "query")
	this.So(this.inner.selectParameters, should.Equal, []any{1, 2, 3})
}
func (this *NormalizeContextCancellationConnectionPoolFixture) TestSelect_Failed() {
	expectedResult := new(FakeSelectResult)
	this.inner.selectResult = expectedResult
	selectErr := errors.New("SELECT ERROR")
	this.inner.selectError = selectErr

	result, err := this.adapter.Select(context.Background(), "query", 1, 2, 3)

	this.So(result, should.Equal, expectedResult)
	this.So(err, should.Equal, selectErr)
	this.So(this.inner.selectCalls, should.Equal, 1)
	this.So(this.inner.selectStatement, should.Equal, "query")
	this.So(this.inner.selectParameters, should.Equal, []any{1, 2, 3})
}
func (this *NormalizeContextCancellationConnectionPoolFixture) TestSelect_AdaptContextCancelled() {
	expectedResult := new(FakeSelectResult)
	this.inner.selectResult = expectedResult
	this.inner.selectError = operationCanceledErr

	result, err := this.adapter.Select(context.Background(), "query", 1, 2, 3)

	this.So(result, should.Equal, expectedResult)
	this.So(this.inner.selectCalls, should.Equal, 1)
	this.So(this.inner.selectStatement, should.Equal, "query")
	this.So(this.inner.selectParameters, should.Equal, []any{1, 2, 3})
	this.So(errors.Is(err, operationCanceledErr), should.BeTrue)
	this.So(errors.Is(err, context.Canceled), should.BeTrue)
}

func (this *NormalizeContextCancellationConnectionPoolFixture) TestContextCancellationErrorAdapter_NilError() {
	err := this.adapter.normalizeContextCancellationError(nil)
	this.So(err, should.BeNil)
}
func (this *NormalizeContextCancellationConnectionPoolFixture) TestContextCancellationErrorAdapter_GenericError() {
	genericErr := errors.New("generic error")
	err := this.adapter.normalizeContextCancellationError(genericErr)
	this.So(err, should.Equal, genericErr)
}
func (this *NormalizeContextCancellationConnectionPoolFixture) TestContextCancellationErrorAdapter_OperationCanceledError() {
	err := this.adapter.normalizeContextCancellationError(operationCanceledErr)
	this.So(errors.Is(err, operationCanceledErr), should.BeTrue)
	this.So(errors.Is(err, context.Canceled), should.BeTrue)
}
func (this *NormalizeContextCancellationConnectionPoolFixture) TestContextCancellationErrorAdapter_ClosedConnectionError() {
	err := this.adapter.normalizeContextCancellationError(closedNetworkConnectionErr)
	this.So(errors.Is(err, closedNetworkConnectionErr), should.BeTrue)
	this.So(errors.Is(err, context.Canceled), should.BeTrue)
}

var (
	operationCanceledErr       = errors.New("operation was canceled")
	closedNetworkConnectionErr = errors.New("use of closed network connection")
)
