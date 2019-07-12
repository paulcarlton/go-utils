// (c) Copyright 2019 Hewlett Packard Enterprise Development LP

package core

import (
	"fmt"
	"testing"

	"github.com/paulcarlton/go-utils/pkg/testutils"
)

var msg = "test error text"

func TestCaller(t *testing.T) {
	msg := "MakeError return"
	testFunc := func() error {
		return MakeError("test1", ErrorBadRequest, msg)
	}
	expected := &cerror{
		id:                 "test1",
		where:              "core.TestCaller.func1() - error_test.go(NN)",
		code:               ErrorBadRequest,
		message:            msg,
		recommendedActions: []string{},
	}
	newErr := testFunc()
	if !CompareErrors(newErr.(Error), expected) || testutils.FailTests {
		t.Errorf("\nTest: 1\nExpected:\n%s\nGot.....:\n%s", expected.FullInfo(), newErr.(Error).FullInfo())
	}

	err := fmt.Errorf("%s", "standard error")
	msg = fmt.Sprintf("RaiseError return, %s", err.Error())
	testFunc = func() error {
		return RaiseError("test1", ErrorBadRequest, "RaiseError return", err)
	}

	expected = &cerror{
		id:                 "test1",
		where:              "core.TestCaller.func2() - error_test.go(NN)",
		code:               ErrorBadRequest,
		message:            msg,
		recommendedActions: []string{},
	}
	newErr = testFunc()
	if !CompareErrors(newErr.(Error), expected) || testutils.FailTests {
		t.Errorf("\nTest: 2\nExpected:\n%s\nGot.....:\n%s", expected.FullInfo(), newErr.(Error).FullInfo())
	}
}

func TestMakeError(t *testing.T) {
	type ErrorParams struct {
		id   string
		msg  string
		code int
	}
	var tests = []struct {
		testNum  int
		params   *ErrorParams
		expected *cerror
	}{
		{
			testNum: 1,
			params: &ErrorParams{
				id:   "test1",
				code: ErrorBadRequest,
				msg:  msg,
			},
			expected: &cerror{
				id:                 "test1",
				where:              "core.TestMakeError() - error_test.go(NN)",
				code:               ErrorBadRequest,
				message:            msg,
				recommendedActions: []string{},
			},
		},
	}

	for _, test := range tests {
		newErr := MakeError(test.params.id, test.params.code, test.params.msg)
		if !CompareErrors(newErr.(Error), test.expected) || testutils.FailTests {
			t.Errorf("\nTest: %d\nInput1..:\n%+v\nExpected:\n%s\nGot.....:\n%s",
				test.testNum, test.params, test.expected.FullInfo(), newErr.(Error).FullInfo())
		}
	}
}

func TestRaiseErrorAt(t *testing.T) {
	type ErrorParams struct {
		id     string
		code   int
		msg    string
		where  string
		nested interface{}
	}
	var tests = []struct {
		testNum  int
		params   *ErrorParams
		expected *cerror
	}{
		{
			testNum: 1,
			params: &ErrorParams{
				id:     "test1",
				code:   ErrorBadRequest,
				msg:    "raising a Error here because a call generated",
				where:  "somewhere() - something.go (NN)",
				nested: fmt.Errorf("%s", "an instance of a standard error type"),
			},
			expected: &cerror{
				id:                 "test1",
				code:               ErrorBadRequest,
				message:            "raising a Error here because a call generated, an instance of a standard error type",
				where:              "somewhere() - something.go (NN)",
				nestedError:        nil,
				recommendedActions: []string{},
			},
		},
		{
			testNum: 2,
			params: &ErrorParams{
				id:     "test1",
				code:   ErrorUnknown,
				msg:    "failed to so something",
				where:  "somewhere() - something.go (NN)",
				nested: fmt.Errorf("%s", "something about a permission error etc"),
			},
			expected: &cerror{
				id:                 "test1",
				code:               ErrorUnauthorized,
				message:            "failed to so something, something about a permission error etc",
				where:              "somewhere() - something.go (NN)",
				nestedError:        nil,
				recommendedActions: []string{},
			},
		},

		{
			testNum: 3,
			params: &ErrorParams{
				id:    "test1",
				code:  ErrorUnknown,
				msg:   "failed to so something",
				where: "somewhere() - something.go (NN)",
				nested: &cerror{
					id:                 "test1",
					code:               ErrorUnauthorized,
					message:            "something about a permission error etc",
					where:              "core.TestRaiseError() - error_test.go(NN)",
					recommendedActions: []string{},
				},
			},
			expected: &cerror{
				id:                 "test1",
				code:               ErrorUnauthorized,
				message:            "failed to so something",
				where:              "somewhere() - something.go (NN)",
				recommendedActions: []string{},
				nestedError: &cerror{
					id:                 "test1",
					code:               ErrorUnauthorized,
					message:            "something about a permission error etc",
					where:              "core.TestRaiseError() - error_test.go(NN)",
					recommendedActions: []string{},
					nestedError:        nil,
				},
			},
		},
	}

	for _, test := range tests {
		newErr := RaiseErrorAt(test.params.id, test.params.code, test.params.msg, test.params.where, test.params.nested)
		if !CompareErrors(newErr, test.expected) || testutils.FailTests {
			t.Errorf("\nTest: %d\nInput1..:\n%+v\nExpected:\n%s\nGot.....:\n%s",
				test.testNum, test.params, test.expected.FullInfo(), newErr.(Error).FullInfo())
		}
	}
}
func TestRaiseError(t *testing.T) {
	type ErrorParams struct {
		id     string
		msg    string
		code   int
		nested interface{}
	}
	var tests = []struct {
		testNum  int
		params   *ErrorParams
		expected *cerror
	}{
		{
			testNum: 1,
			params: &ErrorParams{
				id:     "test1",
				code:   ErrorBadRequest,
				msg:    "raising a Error here because a call generated",
				nested: fmt.Errorf("%s", "an instance of a standard error type"),
			},
			expected: &cerror{
				id:                 "test1",
				where:              "core.TestRaiseError() - error_test.go(NN)",
				code:               ErrorBadRequest,
				message:            "raising a Error here because a call generated, an instance of a standard error type",
				nestedError:        nil,
				recommendedActions: []string{},
			},
		},

		{
			testNum: 2,
			params: &ErrorParams{
				id:   "test1",
				code: ErrorUnknown, msg: "failed to so something",
				nested: fmt.Errorf("%s", "something about a permission error etc"),
			},
			expected: &cerror{
				id:                 "test1",
				where:              "core.TestRaiseError() - error_test.go(NN)",
				code:               ErrorUnauthorized,
				message:            "failed to so something, something about a permission error etc",
				nestedError:        nil,
				recommendedActions: []string{},
			},
		},

		{
			testNum: 3,
			params: &ErrorParams{
				id:   "test1",
				code: ErrorUnknown,
				msg:  "failed to so something",
				nested: &cerror{
					id:                 "test1",
					where:              "core.TestRaiseError() - error_test.go(NN)",
					code:               ErrorUnauthorized,
					message:            "something about a permission error etc",
					recommendedActions: []string{},
				},
			},
			expected: &cerror{
				id:                 "test1",
				where:              "core.TestRaiseError() - error_test.go(NN)",
				code:               ErrorUnauthorized,
				message:            "failed to so something",
				recommendedActions: []string{},
				nestedError: &cerror{
					id:                 "test1",
					where:              "core.TestRaiseError() - error_test.go(NN)",
					code:               ErrorUnauthorized,
					message:            "something about a permission error etc",
					recommendedActions: []string{},
					nestedError:        nil,
				},
			},
		},
		{
			testNum: 4,
			params: &ErrorParams{
				id:     "test1",
				code:   ErrorBadRequest,
				msg:    "raising a Error here because a call generated",
				nested: "a generic error string"},
			expected: &cerror{
				id:                 "test1",
				where:              "core.TestRaiseError() - error_test.go(NN)",
				code:               ErrorBadRequest,
				message:            "raising a Error here because a call generated, a generic error string",
				recommendedActions: []string{},
				nestedError:        nil,
			},
		},
		{
			testNum: 5,
			params: &ErrorParams{
				id:     "test1",
				code:   ErrorBadRequest,
				msg:    "raising a Error here because a call generated",
				nested: 42,
			},
			expected: &cerror{
				id:                 "test1",
				where:              "core.TestRaiseError() - error_test.go(NN)",
				code:               ErrorBadRequest,
				message:            "raising a Error here because a call generated, 42",
				recommendedActions: []string{},
				nestedError:        nil,
			},
		},
		{
			testNum: 6,
			params: &ErrorParams{
				id:     "test1",
				code:   ErrorBadRequest,
				msg:    "raising a Error here because a call generated",
				nested: []float32{0.02, 0.111},
			},
			expected: &cerror{
				id:                 "test1",
				where:              "core.TestRaiseError() - error_test.go(NN)",
				code:               ErrorBadRequest,
				message:            "raising a Error here because a call generated, [0.02 0.111]",
				recommendedActions: []string{},
				nestedError:        nil,
			},
		},
	}

	for _, test := range tests {
		newErr := RaiseError(test.params.id, test.params.code, test.params.msg, test.params.nested)
		if !CompareErrors(newErr.(*cerror), test.expected) || testutils.FailTests {
			t.Errorf("\nTest: %d\nInput1..:\n%+v\nExpected:\n%s\nGot.....:\n%s",
				test.testNum, test.params, test.expected.FullInfo(), newErr.(Error).FullInfo())
		}
	}
}

type ErrorParams struct {
	id    string
	msg   string
	code  int
	where string
}

type ErrorTest struct {
	params   *ErrorParams
	expected *cerror
}

func tester(t *testing.T, tests []ErrorTest) {
	for _, test := range tests {
		newErr := makeError(test.params.id, test.params.code, test.params.msg, test.params.where)
		if newErr.(*cerror).Error() != test.expected.Error() || testutils.FailTests {
			t.Errorf("\nInput1..:\n%+v\nExpected:\n%s\nGot.....:\n%s",
				test.params, test.expected, newErr)
		}
	}
}

func TestMakeErrorAt(t *testing.T) {
	var tests = []ErrorTest{
		{
			params: &ErrorParams{
				id:    "test1",
				code:  ErrorBadRequest,
				msg:   msg,
				where: "core.TestMakeErrorAt() - error_test.go(NN)"},
			expected: &cerror{
				id:      "test1",
				where:   "core.TestMakeErrorAt() - error_test.go(NN)",
				code:    ErrorBadRequest,
				message: msg,
			},
		},
	}
	tester(t, tests)
}

func TestLocalMakeError(t *testing.T) {
	var tests = []ErrorTest{
		{
			params: &ErrorParams{
				id:    "test1",
				code:  ErrorBadRequest,
				msg:   msg,
				where: "core.TestmakeErrorAt() - error_test.go(NN)"},
			expected: &cerror{
				id:      "test1",
				where:   "core.TestmakeErrorAt() - error_test.go(NN)",
				code:    ErrorBadRequest,
				message: msg,
			},
		},
		{
			params: &ErrorParams{
				id:    "test1",
				code:  ErrorUnknown,
				msg:   "something about a permission error etc",
				where: "core.TestmakeErrorAt() - error_test.go(NN)",
			},
			expected: &cerror{
				id:      "test1",
				where:   "core.TestmakeErrorAt() - error_test.go(NN)",
				code:    ErrorUnauthorized,
				message: "something about a permission error etc",
			},
		},
	}
	tester(t, tests)
}

func TestErrorAddRecommendedAction(t *testing.T) {
	var tests = []struct {
		testNum  int
		coreErr  *cerror
		actions  []string
		expected string
	}{
		{
			testNum: 1,
			coreErr: &cerror{
				id:      "test1",
				where:   "not available",
				code:    ErrorBadRequest,
				message: msg,
			},
			actions: []string{"panic and run round like a headless chicken", "pray to whatever god you worship"},
			expected: fmt.Sprintf("%s %s %s %s\nRecommended actions...\npanic and run round like a headless chicken\npray to whatever god you worship",
				"not available", "test1", CodeText(ErrorBadRequest), msg),
		},
		{
			testNum: 2,
			coreErr: &cerror{
				id:      "test1",
				where:   "not available",
				code:    ErrorBadRequest,
				message: msg,
			},
			actions:  []string{},
			expected: fmt.Sprintf("%s %s %s %s", "not available", "test1", CodeText(ErrorBadRequest), msg)},
	}

	for _, test := range tests {
		if e := test.coreErr.AddRecommendedActions(test.actions...); e == nil {
			fullInfo := test.coreErr.FullInfo()
			if fullInfo != test.expected || testutils.FailTests {
				t.Errorf("\nTest: %d\nInput1..:\n%+v\nExpected:\n%s\nGot.....:\n%s",
					test.testNum, test.coreErr, test.expected, fullInfo)
			}
		} else {
			t.Errorf("\nTest: %d\nFailed when calling AddRecommendedActions", test.testNum)
		}
	}
}

func TestErrorAddDetails(t *testing.T) {
	var tests = []struct {
		testNum  int
		coreErr  *cerror
		details  string
		expected string
	}{
		{
			testNum: 1,
			coreErr: &cerror{
				id:      "test1",
				where:   "not available",
				code:    ErrorBadRequest,
				message: msg,
			},
			details: "this error occurs on a Tuesday if it is raining",
			expected: fmt.Sprintf("%s %s %s %s\nthis error occurs on a Tuesday if it is raining",
				"not available", "test1", CodeText(ErrorBadRequest), msg)},
		{
			testNum: 2,
			coreErr: &cerror{
				id:      "test1",
				where:   "not available",
				code:    ErrorBadRequest,
				message: msg,
			},
			details:  "",
			expected: fmt.Sprintf("%s %s %s %s", "not available", "test1", CodeText(ErrorBadRequest), msg)},
	}

	for _, test := range tests {
		if e := test.coreErr.AddDetails(test.details); e == nil {
			fullInfo := test.coreErr.FullInfo()
			if fullInfo != test.expected || testutils.FailTests {
				t.Errorf("\nTest: %d\nInput1..:\n%+v\nExpected:\n%s\nGot.....:\n%s",
					test.testNum, test.coreErr, test.expected, fullInfo)
			}
		} else {
			t.Errorf("\nTest: %d\nFailed when calling AddDetails", test.testNum)
		}
	}
}

func TestErrorAddNestedCoreError(t *testing.T) {
	var tests = []struct {
		testNum  int
		coreErr  *cerror
		nested   Error
		expected string
	}{
		{
			testNum: 1,
			coreErr: &cerror{
				id:      "test1",
				where:   "core.TestAddNested - error_test.go(NN)",
				code:    ErrorBadRequest,
				message: "very bad request",
			},
			nested: &cerror{
				id:      "test1",
				where:   "core.TestAddNested - error_test.go(NN)",
				code:    ErrorInvalidInput,
				message: "try giving me the right data",
			},
			expected: fmt.Sprintf("core.TestAddNested - error_test.go(NN) test1 %s very bad request\nNested Errors...\ncore.TestAddNested - "+
				"error_test.go(NN) test1 %s try giving me the right data", CodeText(ErrorBadRequest), CodeText(ErrorInvalidInput))},
		{
			testNum: 2,
			coreErr: &cerror{
				id:      "test1",
				where:   "core.TestAddNested - error_test.go(NN)",
				code:    ErrorBadRequest,
				message: "very bad request",
			},
			nested:   nil,
			expected: fmt.Sprintf("core.TestAddNested - error_test.go(NN) test1 %s very bad request", CodeText(ErrorBadRequest))},
	}

	for _, test := range tests {
		test.coreErr.addNested(test.nested)
		fullInfo := test.coreErr.FullInfo()
		if fullInfo != test.expected || testutils.FailTests {
			t.Errorf("\nTest: %d\nInput1..:\n%+v\nExpected:\n%s\nGot.....:\n%s",
				test.testNum, test.coreErr, test.expected, fullInfo)
		}
	}
}

func TestErrorAddNestedStdError(t *testing.T) {
	var tests = []struct {
		testNum  int
		coreErr  *cerror
		nested   error
		expected string
	}{
		{
			testNum: 1,
			coreErr: &cerror{
				id:      "test1",
				where:   "core.TestAddNested - error_test.go(NN)",
				code:    ErrorBadRequest,
				message: "very bad request",
			},
			nested: fmt.Errorf("%s", "an instance of a standard error type"),
			expected: fmt.Sprintf("core.TestAddNested - error_test.go(NN) test1 %s very bad request\nNested Errors...\n"+
				"an instance of a standard error type", CodeText(ErrorBadRequest))},
		{
			testNum: 2,
			coreErr: &cerror{
				id:      "test1",
				where:   "core.TestAddNested - error_test.go(NN)",
				code:    ErrorBadRequest,
				message: "very bad request",
			},
			nested: nil,
			expected: fmt.Sprintf("core.TestAddNested - error_test.go(NN) test1 %s very bad request",
				CodeText(ErrorBadRequest))},
	}

	for _, test := range tests {
		test.coreErr.addNested(test.nested)
		fullInfo := test.coreErr.FullInfo()
		if fullInfo != test.expected || testutils.FailTests {
			t.Errorf("\nTest: %d\nInput1..:\n%+v\nExpected:\n%s\nGot.....:\n%s",
				test.testNum, test.coreErr, test.expected, fullInfo)
		}
	}
}

func TestErrorFullInfo(t *testing.T) {
	var tests = []struct {
		testNum  int
		coreErr  *cerror
		nested   *cerror
		details  string
		actions  []string
		expected string
	}{
		{
			testNum: 1,
			coreErr: &cerror{
				id:      "test1",
				where:   "core.TestAddNested - error_test.go(NN)",
				code:    ErrorBadRequest,
				message: "very bad request",
			},
			nested: &cerror{
				id:      "test1",
				where:   "core.TestAddNested - error_test.go(NN)",
				code:    ErrorInvalidInput,
				message: "try giving me the right data",
			},
			details: "error details",
			actions: []string{"action1", "action2"},
			expected: fmt.Sprintf("core.TestAddNested - error_test.go(NN) test1 %s very bad request\nerror "+
				"details\nRecommended actions...\naction1\naction2\n"+
				"Nested Errors...\ncore.TestAddNested - error_test.go(NN) test1 %s try giving me the right data",
				CodeText(ErrorBadRequest), CodeText(ErrorInvalidInput))},
	}

	for _, test := range tests {
		test.coreErr.addNested(test.nested)
		e1 := test.coreErr.AddRecommendedActions(test.actions...)
		e2 := test.coreErr.AddDetails(test.details)
		if e1 != nil || e2 != nil {
			t.Errorf("\nTest: %d\nFailed during AddRecommendedActions or AddDetails",
				test.testNum)
		} else {
			fullInfo := test.coreErr.FullInfo()
			if fullInfo != test.expected || testutils.FailTests {
				t.Errorf("\nTest: %d\nInput1..:\n%+v\nExpected:\n%s\nGot.....:\n%s",
					test.testNum, test.coreErr, test.expected, fullInfo)
			}
		}
	}
}

func TestErrorError(t *testing.T) {
	var tests = []struct {
		testNum  int
		coreErr  *cerror
		expected string
	}{
		{
			testNum: 1,
			coreErr: &cerror{
				id:      "test1",
				where:   "not available",
				code:    ErrorBadRequest,
				message: msg,
			},
			expected: fmt.Sprintf("%s %s %s %s", "not available", "test1", CodeText(ErrorBadRequest), msg)},
		{
			testNum: 2,
			coreErr: &cerror{
				id:      "",
				where:   "not available",
				code:    ErrorBadRequest,
				message: msg,
			},
			expected: fmt.Sprintf("%s %s %s", "not available", CodeText(ErrorBadRequest), msg)},
		{
			testNum:  3,
			coreErr:  nil,
			expected: "",
		},
	}

	for _, test := range tests {
		result := test.coreErr.Error()
		if result != test.expected || testutils.FailTests {
			t.Errorf("\nTest: %d\nInput1..:\n%+v\nExpected:\n%s\nGot.....:\n%s",
				test.testNum, test.coreErr, test.expected, result)
		}
	}
}

func TestErrorSetCode(t *testing.T) {
	var tests = []struct {
		testNum  int
		coreErr  *cerror
		expected string
	}{
		{
			testNum: 1,
			coreErr: &cerror{
				id:      "test1",
				where:   "core.TestSetCode() - error_test.go(NN)",
				code:    ErrorBadRequest,
				message: msg,
			},
			expected: fmt.Sprintf("core.TestSetCode() - error_test.go(NN) test1 %s %s", CodeText(ErrorBadRequest), msg)},
		{
			testNum: 2,
			coreErr: &cerror{
				id:      "test1",
				where:   "core.TestSetCode() - error_test.go(NN)",
				code:    ErrorUnknown,
				message: "something about a permission error etc",
			},
			expected: fmt.Sprintf("core.TestSetCode() - error_test.go(NN) test1 %s %s", CodeText(ErrorUnauthorized),
				"something about a permission error etc")},
		{
			testNum:  3,
			coreErr:  nil,
			expected: "",
		},
	}

	for _, test := range tests {
		if e := test.coreErr.SetCode(); e == nil {
			if test.coreErr.Error() != test.expected || testutils.FailTests {
				t.Errorf("\nTest: %d\nInput1..:\n%+v\nExpected:\n%s\nGot.....:\n%s",
					test.testNum, test.coreErr, test.expected, test.coreErr.Error())
			}
		} else {
			if test.testNum != 3 { // testNum 3 is expected to return an error
				t.Errorf("\nTest: %d\nFailed when calling SetCode",
					test.testNum)
			}
		}
	}
}

func TestErrorSetMessage(t *testing.T) {
	var tests = []struct {
		testNum  int
		coreErr  *cerror
		message  string
		expected string
	}{
		{
			testNum: 1,
			coreErr: &cerror{
				id:      "test1",
				where:   "core.TestSetCode() - error_test.go(NN)",
				code:    ErrorBadRequest,
				message: msg,
			},
			message:  "something",
			expected: fmt.Sprintf("core.TestSetCode() - error_test.go(NN) test1 %s %s", CodeText(ErrorBadRequest), "something")},
		{
			testNum:  2,
			coreErr:  nil,
			message:  "something",
			expected: "",
		},
	}

	for _, test := range tests {
		if e := test.coreErr.SetMessage(test.message); e == nil {
			if test.coreErr.Error() != test.expected || testutils.FailTests {
				t.Errorf("\nTest: %d\nInput1..:\n%+v\nExpected:\n%s\nGot.....:\n%s",
					test.testNum, test.coreErr, test.expected, test.coreErr.Error())
			}
		} else {
			if test.testNum != 2 { // testNum 2 is expected to err
				t.Errorf("\nTest: %d\nFailed calling SetMessage",
					test.testNum)
			}
		}
	}
}
