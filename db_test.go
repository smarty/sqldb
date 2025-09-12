package sqldb

import (
	"reflect"
	"testing"
)

func TestInterleaveParameters(t *testing.T) {
	actual := make(map[string][]any)
	for statement, args := range interleaveParameters("?,?;"+"?;;;"+"?,?,?", 1, 2, 3, 4, 5, 6) {
		actual[statement] = args
	}
	expected := map[string][]any{
		"?,?;":   {1, 2},
		"?;":     {3},
		"?,?,?;": {4, 5, 6},
	}
	assertEqual(t, expected, actual)
}

func assertEqual(t *testing.T, expected, actual any) {
	if reflect.DeepEqual(expected, actual) {
		return
	}
	t.Helper()
	t.Errorf("\n"+
		"expected: %v\n"+
		"actual:   %v",
		expected,
		actual,
	)
}
