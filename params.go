package sqldb

import (
	"iter"
	"strings"
)

// interleaveParameters splits the statements (on ';') and pairs each one with its corresponding parameters.
func interleaveParameters(statements string, parameters ...any) iter.Seq2[string, []any] {
	return func(yield func(string, []any) bool) {
		index := 0
		for statement := range strings.SplitSeq(statements, ";") {
			if len(strings.TrimSpace(statement)) == 0 {
				continue
			}
			statement += ";" // terminate the statement
			indexOffset := strings.Count(statement, "?")
			params := parameters[index : index+indexOffset]
			index += indexOffset
			if !yield(statement, params) {
				return
			}
		}
	}
}
