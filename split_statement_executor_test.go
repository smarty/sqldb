package sqldb

import (
	"errors"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type SplitStatementExecutorFixture struct {
	*gunit.Fixture

	fakeInner *FakeExecutor
	executor  *SplitStatementExecutor
}

func (this *SplitStatementExecutorFixture) Setup() {
	this.fakeInner = &FakeExecutor{}
	this.executor = NewSplitStatementExecutor(this.fakeInner, "?")
}

///////////////////////////////////////////////////////////////

func (this *SplitStatementExecutorFixture) TestStatementAndParameterCountsDoNotMatch() {
	this.fakeInner.affected = 1

	affected, err := this.executor.Execute("? ? ?")

	this.So(affected, should.Equal, 0)
	this.So(err, should.NotBeNil)
	this.So(this.fakeInner.statements, should.BeEmpty)
}

func (this *SplitStatementExecutorFixture) TestSingleStatement() {
	this.fakeInner.affected = 1
	affected, err := this.executor.Execute("statement ? ?", 1, 2)

	this.So(affected, should.Equal, this.fakeInner.affected)
	this.So(err, should.BeNil)
	this.So(this.fakeInner.statements, should.Resemble, []string{"statement ? ?;"})
	this.So(this.fakeInner.parameters, should.Resemble, [][]interface{}{{1, 2}})
}

func (this *SplitStatementExecutorFixture) TestEmptyStatementsAreSkipped() {
	affected, err := this.executor.Execute(";;;;")

	this.So(affected, should.Equal, 0)
	this.So(err, should.BeNil)
	this.So(this.fakeInner.statements, should.BeEmpty)
	this.So(this.fakeInner.parameters, should.BeEmpty)
}

func (this *SplitStatementExecutorFixture) TestMultipleStatements() {
	this.fakeInner.affected = 2

	affected, err := this.executor.Execute("1 ?; 2 ? ?; 3 ? ? ?", 1, 2, 3, 4, 5, 6)

	this.So(affected, should.Equal, this.fakeInner.affected*3)
	this.So(err, should.BeNil)
	this.So(this.fakeInner.statements, should.Resemble, []string{
		"1 ?;",
		"2 ? ?;",
		"3 ? ? ?;",
	})
	this.So(this.fakeInner.parameters, should.Resemble, [][]interface{}{
		{1},
		{2, 3},
		{4, 5, 6},
	})
}

func (this *SplitStatementExecutorFixture) TestFailureAbortsAdditionalStatements() {
	this.fakeInner.affected = 10
	this.fakeInner.errorsToReturn = []error{nil, errors.New("")}

	affected, err := this.executor.Execute("1;2;3")

	this.So(affected, should.Equal, 0)
	this.So(err, should.Equal, this.fakeInner.errorsToReturn[1])
	this.So(this.fakeInner.statements, should.Resemble, []string{"1;", "2;"})
}
