// (c) Copyright 2019 Hewlett Packard Enterprise Development LP

package core

import (
	"fmt"
	"testing"

	"github.hpe.com/platform-core/utils/pkg/testutils"
)

func TestErrorText(t *testing.T) {
	type testStruct struct {
		f float32
		s string
	}
	var tests = []struct {
		testNum  int
		theErr   interface{}
		expected string
	}{
		{
			testNum: 1,
			theErr: &cerror{
				id:      "test1",
				where:   "core.TestErrorText() - error_test.go(NN)",
				code:    ErrorBadRequest,
				message: "very bad request",
			},
			expected: fmt.Sprintf("core.TestErrorText() - error_test.go(NN) test1 %s very bad request", CodeText(ErrorBadRequest))},
		{
			testNum:  2,
			theErr:   fmt.Errorf("%s", "An instance of a standard error type"),
			expected: "An instance of a standard error type",
		},
		{
			testNum:  3,
			theErr:   "",
			expected: "",
		},
		{
			testNum:  4,
			theErr:   &cerror{},
			expected: "  ",
		},
		{
			testNum:  5,
			theErr:   int(2),
			expected: "2",
		},
		{
			testNum:  6,
			theErr:   nil,
			expected: "<nil>",
		},
		{
			testNum:  7,
			theErr:   &testStruct{f: 0, s: ""},
			expected: "&{f:0 s:}",
		},
		{
			testNum:  8,
			theErr:   testStruct{f: 0, s: ""},
			expected: "{f:0 s:}",
		},
	}

	for _, test := range tests {
		result := ErrorText(test.theErr)
		if result != test.expected || testutils.FailTests {
			t.Errorf("\nTest: %d\nInput1..:\n%+v\nExpected:\n'%s'\nGot.....:\n'%s'",
				test.testNum, test.theErr, test.expected, result)
		}
	}
}
