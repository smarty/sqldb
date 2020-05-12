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

func (this *SplitStatementExecutor) Execute(statement string, parameters ...interface{}) (uint64, error) {
	// TODO
	if argumentCount := strings.Count(statement, this.delimiter); argumentCount != len(parameters) {
		return 0, fmt.Errorf("Not enough arguments supplied for the statement. Expected: %d, received: %d", argumentCount, len(parameters))
	}

	var count uint64
	index := 0
	for _, statement = range strings.Split(statement, ";") {
		if len(strings.TrimSpace(statement)) == 0 {
			continue
		}

		statement += ";" // terminate the statement
		indexOffset := strings.Count(statement, this.delimiter)
		if affected, err := this.actual.Execute(statement, parameters[index:index+indexOffset]...); err != nil {
			return 0, err
		} else {
			count += affected
		}

		index += indexOffset
	}

	return count, nil
}
