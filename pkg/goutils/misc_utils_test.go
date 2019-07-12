// (c) Copyright 2018-2019 Hewlett Packard Enterprise Development LP

package goutils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/paulcarlton/go-utils/pkg/core"
	"github.com/paulcarlton/go-utils/pkg/testutils"
)

func TestCallers(t *testing.T) {
	type callerInfo struct {
		testNum  int
		levels   uint
		short    bool
		expected []string
	}
	pwd, err := os.Getwd()
	if err != nil {
		t.Errorf("failed to get current directory: %s", err)
	}
	var tests = []callerInfo{
		{testNum: 1, levels: 10, short: false, expected: []string{
			fmt.Sprintf("%s%s%s",
				"github.com/paulcarlton/go-utils/pkg/internal/common.Callers() - ", filepath.Dir(pwd), "/internal/common/misc_utils.go(NN)"),
			fmt.Sprintf("%s%s%s",
				"github.com/paulcarlton/go-utils/pkg/goutils.Callers() - ", pwd, "/misc_utils.go(NN)"),
			fmt.Sprintf("%s%s%s",
				"github.com/paulcarlton/go-utils/pkg/goutils.TestCallers() - ", pwd, "/misc_utils_test.go(NN)")}},
		{testNum: 2, levels: 10, short: true, expected: []string{
			"common.Callers() - misc_utils.go(NN)",
			"goutils.Callers() - misc_utils.go(NN)",
			"goutils.TestCallers() - misc_utils_test.go(NN)"}},
		{testNum: 3, levels: 1, short: true, expected: []string{
			"common.Callers() - misc_utils.go(NN)"}},
		{testNum: 4, levels: 0, short: true, expected: []string{}},
	}

	for _, test := range tests {
		callers, err := Callers(test.levels, test.short)
		if err != nil {
			t.Errorf("%s", err)
		}
		callers = testutils.RemoveBottom(callers)
		if !testutils.CompareWhereList(test.expected, callers) || testutils.FailTests {
			t.Errorf("\nTest: %d\nExpected:\n%s\nGot:\n%s", test.testNum, testutils.DisplayStrings(test.expected), testutils.DisplayStrings(callers))
		}
	}
}

func TestGetCaller(t *testing.T) {
	type callerInfo struct {
		testNum  int
		skip     uint
		short    bool
		expected string
	}
	pwd, err := os.Getwd()
	if err != nil {
		t.Errorf("failed to get current directory: %s", err)
	}
	var tests = []callerInfo{
		{testNum: 1, skip: 1, short: false,
			expected: fmt.Sprintf("github.com/paulcarlton/go-utils/pkg/internal/common.GetCaller() - %s/internal/common/misc_utils.go(NN)", filepath.Dir(pwd))},
		{testNum: 2, skip: 1, short: true,
			expected: "common.GetCaller() - misc_utils.go(NN)"},
		{testNum: 3, skip: 2, short: true,
			expected: "goutils.GetCaller() - misc_utils.go(NN)"},
		{testNum: 4, skip: 3, short: true,
			expected: "goutils.TestGetCaller() - misc_utils_test.go(NN)"},
		{testNum: 5, skip: 7, short: true,
			expected: "not available"},
	}

	for _, test := range tests {
		caller := GetCaller(test.skip, test.short)
		if !testutils.CompareWhere(caller, test.expected) || testutils.FailTests {
			t.Errorf("\nTest: %d\nExpected:\n%s\nGot....:\n%s", test.testNum, test.expected, caller)
		}
	}
}

func TestToJSON(t *testing.T) {
	type TestInfo struct {
		object   interface{}
		expected string
	}
	var tests = []TestInfo{{object: []string{"one", "two"}, expected: "[\n\t\"one\",\n\t\"two\"\n]"},
		{object: "one", expected: "\"one\""},
		{object: map[string]string{"1": "one", "2": "two"}, expected: "{\n\t\"1\": \"one\",\n\t\"2\": \"two\"\n}"},
		{object: "", expected: "\"\""},
		{object: nil, expected: "null"},
	}

	for _, test := range tests {
		json, err := ToJSON(test.object)
		if err != nil {
			t.Errorf(err.Error())
		}
		if json != test.expected || testutils.FailTests {
			t.Errorf("\nExpected:\n%s\nGot....:\n%s", test.expected, json)
		}
	}
}

func TestFindInStringSlice(t *testing.T) {
	var tests = []struct {
		testNum  int
		array    []string
		str      string
		expected int
	}{
		{1, []string{"a", "b", "c", "d"}, "b", 1},
		{2, []string{"a", "b", "c", "d"}, "x", -1},
		{3, []string{"a", "b", "c", "d"}, "d", 3},
		{4, []string{"a", "d", "c", "d"}, "d", 1},
		{5, []string{"a"}, "", -1},
		{6, []string{}, "", -1},
		{7, []string{}, "s", -1},
	}

	for _, test := range tests {
		result := FindInStringSlice(test.array, test.str)
		if result != test.expected || testutils.FailTests {
			t.Errorf("\nTest: %d\narray:\n%+v\nstr:\n%s\nExpected: %d\nGot.....: %d",
				test.testNum, test.array, test.str, test.expected, result)
		}
	}
}

func TestCastToString(t *testing.T) {
	type expected struct {
		result  string
		coreErr error
	}
	type TestInfo struct {
		testNum  int
		object   interface{}
		expected expected
	}
	coreErr := core.MakeErrorAt("", core.ErrorInvalidInput, "failed to cast to string", "goutils.CastToString() - misc_utils.go(NN)")
	var tests = []TestInfo{
		{testNum: 1, object: []string{"one", "two"}, expected: expected{"", coreErr}},
		{testNum: 2, object: "one", expected: expected{"one", nil}},
		{testNum: 3, object: expected{"str", nil}, expected: expected{"", coreErr}},
		{testNum: 4, object: 2, expected: expected{"", coreErr}},
	}

	for _, test := range tests {
		result, err := CastToString(test.object)
		if result != test.expected.result || !core.CompareErrors(err, test.expected.coreErr) || testutils.FailTests {
			t.Errorf("Test: %d\nExpected:\n%s\n%+v\nGot....:\n%s\n%+v\n", test.testNum, test.expected.result, test.expected.coreErr, result, err)
		}
	}
}

func TestCompareAsJSON(t *testing.T) {
	type TestInfo struct {
		testNum  int
		objects  []interface{}
		expected bool
	}

	type subData struct {
		S string
		A []int
	}
	type testData struct {
		B bool
		I int
		F float64
		X interface{}
		E subData
		subData
	}

	var tests = []TestInfo{
		{testNum: 1, objects: []interface{}{"one", "two"}, expected: false},
		{testNum: 2, objects: []interface{}{"one", "one"}, expected: true},
		{testNum: 3, objects: []interface{}{1, 2}, expected: false},
		{testNum: 4, objects: []interface{}{1, 1}, expected: true},
		{testNum: 5, objects: []interface{}{
			testData{B: true, I: 1, F: 12.43,
				X:       subData{S: "interface"},
				E:       subData{S: "sub struct", A: []int{1, 2, 3}},
				subData: subData{S: "embedded", A: []int{9, 8, 7}}},
			testData{B: true, I: 1, F: 12.43,
				X:       subData{S: "interface"},
				E:       subData{S: "sub struct", A: []int{1, 2, 3}},
				subData: subData{S: "embedded", A: []int{9, 11, 7}}}}, expected: false},
		{testNum: 6, objects: []interface{}{
			testData{B: true, I: 1, F: 12.43,
				X:       subData{S: "interface"},
				E:       subData{S: "sub struct", A: []int{1, 2, 3}},
				subData: subData{S: "embedded", A: []int{9, 8, 7}}},
			testData{B: true, I: 1, F: 12.43,
				X:       subData{S: "interface"},
				E:       subData{S: "sub struct", A: []int{1, 2, 3}},
				subData: subData{S: "embedded", A: []int{9, 8, 7}}}}, expected: true},
	}

	for _, test := range tests {
		result := CompareAsJSON(test.objects[0], test.objects[1])
		if result != test.expected || testutils.FailTests {
			oneJSON, err := ToJSON(test.objects[0])
			if err != nil {
				t.Errorf("failed to convert test data to json, %s", err)
			}
			twoJSON, err := ToJSON(test.objects[1])
			if err != nil {
				t.Errorf("failed to convert test data to json, %s", err)
			}
			t.Errorf("Test: %d\nExpected:\n%t\nGot....:\n%t\nInput Data:\n%s\n%s\n", test.testNum, result, test.expected, oneJSON, twoJSON)
		}
	}
}

func TestCompareStringSlices(t *testing.T) {
	type TestInfo struct {
		testNum   int
		strSlice1 []string
		strSlice2 []string
		expected  bool
	}

	var tests = []TestInfo{
		{testNum: 1, strSlice1: []string{"one", "two"}, strSlice2: []string{"one", "two"}, expected: true},
		{testNum: 2, strSlice1: []string{"one", "two"}, strSlice2: []string{"two", "one"}, expected: true},
		{testNum: 3, strSlice1: []string{"one", "two"}, strSlice2: []string{"one", "one"}, expected: false},
		{testNum: 4, strSlice1: nil, strSlice2: []string{"two", "one"}, expected: false},
		{testNum: 5, strSlice1: []string{}, strSlice2: []string{}, expected: true},
		{testNum: 6, strSlice1: nil, strSlice2: nil, expected: true}}

	for _, test := range tests {
		result := CompareStringSlices(test.strSlice1, test.strSlice2)
		if result != test.expected || testutils.FailTests {
			t.Errorf("Test: %d\nExpected:\n%t\nGot....:\n%t\nInput Data:\n%+v\n%+v\n", test.testNum, test.expected, result, test.strSlice1, test.strSlice2)
		}
	}
}

func TestPrettyJSON(t *testing.T) {
	type expected struct {
		result string
		err    error
	}

	type TestInfo struct {
		testNum  int
		input    string
		expected expected
	}

	var tests = []TestInfo{
		{testNum: 1, input: "", expected: expected{"", fmt.Errorf("unexpected end of JSON input")}},
		{testNum: 2, input: "{\"key\":\"data\"}", expected: expected{"{\n\t\"key\": \"data\"\n}", nil}}}

	for _, test := range tests {
		result, err := PrettyJSON(test.input)
		if result != test.expected.result || !testutils.CompareItems(test.expected.err, err) || testutils.FailTests {
			t.Errorf("Test: %d\nExpected:\n%s\n%+v\nGot....:\n%s\n%+v\n", test.testNum, test.expected.result, test.expected.err, result, err)
		}
	}
}
