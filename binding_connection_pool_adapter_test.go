package sqldb

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestBindingConnectionPoolAdapterFixture(t *testing.T) {
	gunit.Run(new(BindingConnectionPoolAdapterFixture), t)
}

type BindingConnectionPoolAdapterFixture struct {
	*gunit.Fixture

	inner *FakeConnectionPool
	pool  *BindingConnectionPoolAdapter
}

func (this *BindingConnectionPoolAdapterFixture) Setup() {
	this.inner = &FakeConnectionPool{}
	this.pool = NewBindingConnectionPoolAdapter(this.inner, false)
}

///////////////////////////////////////////////////////////////

func (this *BindingConnectionPoolAdapterFixture) TestPing() {
	this.inner.pingError = errors.New("")

	err := this.pool.Ping(context.Background())

	this.So(err, should.Equal, this.inner.pingError)
	this.So(this.inner.pingCalls, should.Equal, 1)
}

func (this *BindingConnectionPoolAdapterFixture) TestBeginTransaction() {
	transaction, err := this.pool.BeginTransaction(context.Background())

	this.So(transaction, should.NotBeNil)
	this.So(reflect.TypeOf(transaction), should.Equal, reflect.TypeOf(&BindingTransactionAdapter{}))
	this.So(err, should.BeNil)
}

func (this *BindingConnectionPoolAdapterFixture) TestBeginFailedTransaction() {
	this.inner.transactionError = errors.New("")

	transaction, err := this.pool.BeginTransaction(context.Background())

	this.So(transaction, should.BeNil)
	this.So(err, should.Equal, this.inner.transactionError)
}

func (this *BindingConnectionPoolAdapterFixture) TestClose() {
	this.inner.closeError = errors.New("")

	err := this.pool.Close()

	this.So(err, should.Equal, this.inner.closeError)
	this.So(this.inner.closeCalls, should.Equal, 1)
}

func (this *BindingConnectionPoolAdapterFixture) TestExecute() {
	this.inner.executeResult = 42
	this.inner.executeError = errors.New("")

	affected, err := this.pool.Execute(context.Background(), "statement")

	this.So(affected, should.Equal, this.inner.executeResult)
	this.So(err, should.Equal, this.inner.executeError)
	this.So(this.inner.executeStatement, should.Equal, "statement")
	this.So(this.inner.executeCalls, should.Equal, 1)
}

func (this *BindingConnectionPoolAdapterFixture) TestBindSelect() {
	this.inner.selectError = errors.New("")

	err := this.pool.BindSelect(context.Background(), nil, "query", 1, 2, 3)

	this.So(err, should.Equal, this.inner.selectError)
	this.So(this.inner.selectCalls, should.Equal, 1)
	this.So(this.inner.selectStatement, should.Equal, "query")
	this.So(this.inner.selectParameters, should.Resemble, []any{1, 2, 3})
}
