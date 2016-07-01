package sqldb

import (
	"fmt"
	"strings"
)

type SplitStatementExecutor struct {
	actual    Executor
	delimiter string
}

func NewSplitStatementExecutor(actual Executor, delimiter string) *SplitStatementExecutor {
	return &SplitStatementExecutor{actual: actual, delimiter: delimiter}
}

func (this *SplitStatementExecutor) Execute(statement string, parameters ...interface{}) error {
	if argumentCount := strings.Count(statement, this.delimiter); argumentCount != len(parameters) {
		return fmt.Errorf("Not enough arguments supplied for the statement. Expected: %d, received: %d", argumentCount, len(parameters))
	}

	index := 0
	for _, statement = range strings.Split(statement, ";") {
		if len(strings.TrimSpace(statement)) == 0 {
			continue
		}

		statement += ";" // terminate the statement
		indexOffset := strings.Count(statement, this.delimiter)
		if _, err := this.actual.Execute(statement, parameters[index:index+indexOffset]...); err != nil {
			return err
		}

		index += indexOffset
	}

	return nil
}
