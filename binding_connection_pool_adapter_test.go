package sqldb

import (
	"errors"
	"reflect"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type BindingConnectionPoolAdapterFixture struct {
	*gunit.Fixture

	inner      *FakeConnectionPool
	connection *BindingConnectionPoolAdapter
}

func (this *BindingConnectionPoolAdapterFixture) Setup() {
	this.inner = &FakeConnectionPool{}
	this.connection = NewBindingConnectionPoolAdapter(this.inner, false)
}

///////////////////////////////////////////////////////////////

func (this *BindingConnectionPoolAdapterFixture) TestPing() {
	this.inner.pingError = errors.New("")

	err := this.connection.Ping()

	this.So(err, should.Equal, this.inner.pingError)
	this.So(this.inner.pingCalls, should.Equal, 1)
}

func (this *BindingConnectionPoolAdapterFixture) TestBeginTransaction() {
	transaction, err := this.connection.BeginTransaction()

	this.So(transaction, should.NotBeNil)
	this.So(reflect.TypeOf(transaction), should.Equal, reflect.TypeOf(&BindingTransactionAdapter{}))
	this.So(err, should.BeNil)
}

func (this *BindingConnectionPoolAdapterFixture) TestBeginFailedTransaction() {
	this.inner.transactionError = errors.New("")

	transaction, err := this.connection.BeginTransaction()

	this.So(transaction, should.BeNil)
	this.So(err, should.Equal, this.inner.transactionError)
}

func (this *BindingConnectionPoolAdapterFixture) TestClose() {
	this.inner.closeError = errors.New("")

	err := this.connection.Close()

	this.So(err, should.Equal, this.inner.closeError)
	this.So(this.inner.closeCalls, should.Equal, 1)
}

func (this *BindingConnectionPoolAdapterFixture) TestExecute() {
	this.inner.executeResult = 42
	this.inner.executeError = errors.New("")

	affected, err := this.connection.Execute("statement")

	this.So(affected, should.Equal, this.inner.executeResult)
	this.So(err, should.Equal, this.inner.executeError)
	this.So(this.inner.executeStatement, should.Equal, "statement")
	this.So(this.inner.executeCalls, should.Equal, 1)
}

func (this *BindingConnectionPoolAdapterFixture) TestBindSelect() {
	this.inner.selectError = errors.New("")

	err := this.connection.BindSelect(nil, "query", 1, 2, 3)

	this.So(err, should.Equal, this.inner.selectError)
	this.So(this.inner.selectCalls, should.Equal, 1)
	this.So(this.inner.selectStatement, should.Equal, "query")
	this.So(this.inner.selectParameters, should.Resemble, []interface{}{1, 2, 3})
}
