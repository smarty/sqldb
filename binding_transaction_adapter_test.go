package sqldb

import (
	"context"
	"errors"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestBindingTransactionAdapterFixture(t *testing.T) {
	gunit.Run(new(BindingTransactionAdapterFixture), t)
}

type BindingTransactionAdapterFixture struct {
	*gunit.Fixture

	inner       *FakeTransaction
	transaction *BindingTransactionAdapter
}

func (this *BindingTransactionAdapterFixture) Setup() {
	this.inner = &FakeTransaction{}
	this.transaction = NewBindingTransactionAdapter(this.inner, false)
}

///////////////////////////////////////////////////////////////

func (this *BindingTransactionAdapterFixture) TestCommit() {
	this.inner.commitError = errors.New("")

	err := this.transaction.Commit()

	this.So(err, should.Equal, this.inner.commitError)
	this.So(this.inner.commitCalls, should.Equal, 1)
}

func (this *BindingTransactionAdapterFixture) TestRollback() {
	this.inner.rollbackError = errors.New("")

	err := this.transaction.Rollback()

	this.So(err, should.Equal, this.inner.rollbackError)
	this.So(this.inner.rollbackCalls, should.Equal, 1)
}

func (this *BindingTransactionAdapterFixture) TestExecute() {
	this.inner.executeResult = 42
	this.inner.executeError = errors.New("")

	affected, err := this.transaction.Execute(context.Background(), "statement")

	this.So(affected, should.Equal, this.inner.executeResult)
	this.So(err, should.Equal, this.inner.executeError)
	this.So(this.inner.executeStatement, should.Equal, "statement")
	this.So(this.inner.executeCalls, should.Equal, 1)
}

func (this *BindingTransactionAdapterFixture) TestBindSelect() {
	this.inner.selectError = errors.New("")

	err := this.transaction.BindSelect(context.Background(), nil, "query", 1, 2, 3)

	this.So(err, should.Equal, this.inner.selectError)
	this.So(this.inner.selectCalls, should.Equal, 1)
	this.So(this.inner.selectStatement, should.Equal, "query")
	this.So(this.inner.selectParameters, should.Resemble, []any{1, 2, 3})
}
