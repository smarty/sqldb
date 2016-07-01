package sqldb

import (
	"errors"
	"strings"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type SplitStatementExecutorFixture struct {
	*gunit.Fixture

	fakeInner *FakeDriverExecutor
	executor  *SplitStatementExecutor
}

func (this *SplitStatementExecutorFixture) Setup() {
	this.fakeInner = &FakeDriverExecutor{}
	this.executor = NewSplitStatementExecutor(this.fakeInner, "?")
}

///////////////////////////////////////////////////////////////

func (this *SplitStatementExecutorFixture) TestStatementAndParameterCountsDoNotMatch() {
	err := this.executor.Execute("? ? ?")

	this.So(err, should.NotBeNil)
	this.So(this.fakeInner.statements, should.BeEmpty)
}

func (this *SplitStatementExecutorFixture) TestSingleStatement() {
	err := this.executor.Execute("statement ? ?", 1, 2)

	this.So(err, should.BeNil)
	this.So(this.fakeInner.statements, should.Resemble, []string{"statement ? ?;"})
	this.So(this.fakeInner.parameters, should.Resemble, [][]interface{}{{1, 2}})
}

func (this *SplitStatementExecutorFixture) TestEmptyStatementsAreSkipped() {
	err := this.executor.Execute(";;;;")

	this.So(err, should.BeNil)
	this.So(this.fakeInner.statements, should.BeEmpty)
	this.So(this.fakeInner.parameters, should.BeEmpty)
}

func (this *SplitStatementExecutorFixture) TestMultipleStatements() {
	err := this.executor.Execute("1 ?; 2 ? ?; 3 ? ? ?", 1, 2, 3, 4, 5, 6)

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
	this.fakeInner.errorsToReturn = []error{nil, errors.New("")}

	err := this.executor.Execute("1;2;3")

	this.So(err, should.Equal, this.fakeInner.errorsToReturn[1])
	this.So(this.fakeInner.statements, should.Resemble, []string{"1;", "2;"})
}

///////////////////////////////////////////////////////////////

type FakeDriverExecutor struct {
	errorsToReturn []error
	statements     []string
	parameters     [][]interface{}
}

func (this *FakeDriverExecutor) Execute(statement string, parameters ...interface{}) (uint64, error) {
	this.statements = append(this.statements, strings.TrimSpace(statement))
	this.parameters = append(this.parameters, parameters)

	if len(this.statements) <= len(this.errorsToReturn) {
		return 1, this.errorsToReturn[len(this.statements)-1]
	}

	return 0, nil
}
