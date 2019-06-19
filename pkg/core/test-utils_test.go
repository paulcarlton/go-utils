// (c) Copyright 2018-2019 Hewlett Packard Enterprise Development LP

package core

import (
	"fmt"
	"testing"

	"github.hpe.com/platform-core/utils/pkg/testutils"
)

func TestCompareErrors(t *testing.T) {
	type Params struct {
		one *cerror
		two *cerror
	}
	coreError := &cerror{
		id:          "test1",
		where:       "testutils.TestRaiseError() - core_error.go(NN)",
		code:        ErrorBadRequest,
		message:     "raising a Error here because a call generated, An instance of a standard error type",
		nestedError: nil,
	}
	nestedError := &cerror{
		id:      "test1",
		where:   "testutils.TestRaiseError() - core_error.go(NN)",
		code:    ErrorUnauthorized,
		message: "failed to so something",
		nestedError: &cerror{
			id:          "test1",
			where:       "testutils.TestRaiseError() - core_error.go(NN)",
			code:        ErrorUnauthorized,
			message:     "something about a permission error etc",
			nestedError: nil,
		},
	}
	nestedError2 := &cerror{
		id:      "test1",
		where:   "testutils.TestRaiseError() - core_error.go(NN)",
		code:    ErrorUnauthorized,
		message: "failed to so something",
		nestedError: &cerror{
			id:          "test1",
			where:       "testutils.TestRaiseError() - core_error.go(NN)",
			code:        ErrorUnauthorized,
			message:     "something about a permission error",
			nestedError: nil,
		},
	}
	nestedError3 := &cerror{
		id:          "test1",
		where:       "testutils.TestRaiseError() - core_error.go(NN)",
		code:        ErrorUnauthorized,
		message:     "failed to so something",
		nestedError: fmt.Errorf("std error as nested"),
	}
	nestedError4 := &cerror{
		id:          "test1",
		where:       "testutils.TestRaiseError() - core_error.go(NN)",
		code:        ErrorUnauthorized,
		message:     "failed to so something",
		nestedError: fmt.Errorf("std error as nested"),
	}
	coreErrorTest := &cerror{
		id:          "test1",
		where:       "testutils.TestRaiseError() - core_error._test.go(NN)",
		code:        ErrorBadRequest,
		message:     "raising a Error here because a call generated, An instance of a standard error type",
		nestedError: nil,
	}

	var tests = []struct {
		testNum  int
		params   Params
		expected bool
	}{
		{
			testNum: 1,
			params: Params{
				one: coreError,
				two: coreError,
			},
			expected: true,
		},
		{
			testNum: 2,
			params: Params{
				one: nil,
				two: coreError,
			},
			expected: false,
		},
		{
			testNum: 3,
			params: Params{
				one: coreError,
				two: nil,
			},
			expected: false,
		},
		{
			testNum: 4,
			params: Params{
				one: nil,
				two: nil,
			},
			expected: true,
		},
		{
			testNum: 5,
			params: Params{
				one: nestedError,
				two: nestedError,
			},
			expected: true,
		},
		{
			testNum: 6,
			params: Params{
				one: coreError,
				two: nestedError,
			},
			expected: false,
		},
		{
			testNum: 7,
			params: Params{
				one: nestedError,
				two: nestedError2,
			},
			expected: false,
		},
		{
			testNum: 8,
			params: Params{
				one: coreErrorTest,
				two: coreErrorTest,
			}, expected: true,
		},
		{
			testNum: 9,
			params: Params{
				one: coreErrorTest,
				two: nestedError3,
			},
			expected: false,
		},
		{
			testNum: 10,
			params: Params{
				one: nestedError4,
				two: nestedError3,
			},
			expected: true,
		},
	}

	for _, test := range tests {
		result := CompareErrors(Error(test.params.one), Error(test.params.two))
		if result != test.expected || testutils.FailTests {
			t.Errorf("\nTest: %d\nInput1..:\n%s\n\n%s\n\nExpected:\n%t\nGot.....:\n%t",
				test.testNum, test.params.one.FullInfo(), test.params.two.FullInfo(), test.expected, result)
		}
	}
}
