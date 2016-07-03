package bindsql

import (
	"errors"
	"reflect"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/sqldb"
)

type BindingConnectionAdapterFixture struct {
	*gunit.Fixture

	connection *BindingConnectionAdapter
	fakeInner  *FakeDriverConnection
}

func (this *BindingConnectionAdapterFixture) Setup() {
	this.fakeInner = &FakeDriverConnection{}
	this.connection = NewDefaultBindingConnection(this.fakeInner)
}

///////////////////////////////////////////////////////////////

func (this *BindingConnectionAdapterFixture) TestPing() {
	this.fakeInner.pingError = errors.New("")

	err := this.connection.Ping()

	this.So(err, should.Equal, this.fakeInner.pingError)
	this.So(this.fakeInner.ping, should.Equal, 1)
}

func (this *BindingConnectionAdapterFixture) TestBeginTransaction() {
	transaction, err := this.connection.BeginTransaction()

	this.So(transaction, should.NotBeNil)
	this.So(reflect.TypeOf(transaction), should.Equal, reflect.TypeOf(&BindingTransactionAdapter{}))
	this.So(err, should.BeNil)
}

func (this *BindingConnectionAdapterFixture) TestBeginFailedTransaction() {
	this.fakeInner.beginError = errors.New("")

	transaction, err := this.connection.BeginTransaction()

	this.So(transaction, should.BeNil)
	this.So(err, should.Equal, this.fakeInner.beginError)
}

func (this *BindingConnectionAdapterFixture) TestClose() {
	this.fakeInner.closeError = errors.New("")

	err := this.connection.Close()

	this.So(err, should.Equal, this.fakeInner.closeError)
	this.So(this.fakeInner.close, should.Equal, 1)
}

func (this *BindingConnectionAdapterFixture) TestExecute() {
	this.fakeInner.executeError = errors.New("")

	_, err := this.connection.Execute("statement;")

	this.So(err, should.Equal, this.fakeInner.executeError)
	this.So(this.fakeInner.executes, should.Resemble, []string{"statement;"})
}

func (this *BindingConnectionAdapterFixture) TestMultiStatementExecute() {
	_, err := this.connection.Execute("statement1;statement2;")

	this.So(err, should.BeNil)
	this.So(this.fakeInner.executes, should.Resemble, []string{"statement1;", "statement2;"})
}

func (this *BindingConnectionAdapterFixture) TestSelect() {
	this.fakeInner.queryError = errors.New("")

	err := this.connection.Select(nil, "query")

	this.So(err, should.Equal, this.fakeInner.queryError)
	this.So(this.fakeInner.queries, should.Resemble, []string{"query"})
}

///////////////////////////////////////////////////////////////

type FakeDriverConnection struct {
	ping  int
	close int
	begin int

	queries  []string
	executes []string

	pingError    error
	beginError   error
	closeError   error
	executeError error
	queryError   error
}

func (this *FakeDriverConnection) Ping() error {
	this.ping++
	return this.pingError
}

func (this *FakeDriverConnection) BeginTransaction() (sqldb.Transaction, error) {
	this.begin++
	return nil, this.beginError
}

func (this *FakeDriverConnection) Close() error {
	this.close++
	return this.closeError
}

func (this *FakeDriverConnection) Execute(statement string, parameters ...interface{}) (uint64, error) {
	this.executes = append(this.executes, statement)
	return 0, this.executeError
}

func (this *FakeDriverConnection) Select(query string, parameters ...interface{}) (sqldb.SelectResult, error) {
	this.queries = append(this.queries, query)
	return nil, this.queryError
}
