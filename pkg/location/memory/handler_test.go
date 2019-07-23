//(c) Copyright 2019 Hewlett Packard Enterprise Development LP

package memory

import (
	"fmt"
	"strings"
	"testing"

	"github.com/paulcarlton/go-utils/pkg/core"
	"github.com/paulcarlton/go-utils/pkg/testutils"

	"github.com/paulcarlton/go-utils/pkg/location"
)

// Tests

func TestMemoryLocationHandlerVerifyScheme(t *testing.T) {
	var v memory
	err := v.VerifyScheme("test.local")

	// Since the uri scheme is not memory this should error out
	// So if no error we want to fail the test
	if err == nil {
		t.Errorf("%s", err)
	}

	// The error should tell what
	if !strings.Contains(err.Error(), location.ErrorStringURISchemeMismatch) {
		t.Errorf("Returned Error: %s", err)
	}
}

func TestMemoryLocationHandler_getSession(t *testing.T) {
	var v memory

	uri := "memory://"

	_, err := v.getSession(uri)
	if err != nil {
		t.Errorf("Session error: %s", err)
	}
}

func TestMemoryLocationHandlerConnect(t *testing.T) {
	var v memory

	uri := "memory://ms,mgmtsvc@"
	if err := v.Connect(uri); err != nil {
		t.Errorf("Connect error: %s", err)
	}
}

func TestMemoryLocationHandler_PutGetData(t *testing.T) {

	var v memory

	uri := "memory:///secret/services/ms/ALL/ms/mgmtsvc/mgmt/test1"
	data1 := map[string]string{
		"one":  "1",
		"two":  "2",
		"four": "4",
		"five": "5",
	}

	if err := v.PutData(uri, data1); err != nil {
		t.Errorf("PutData error: %s", err)
	}

	uri = "memory:///secret/services/ms/ALL/ms/mgmtsvc/mgmt/test1"
	data2, err := v.GetData(uri)
	if err != nil {
		t.Errorf("GetData error: %s", err)
	}

	d1, ok := data2.(map[string]string)
	if !ok {
		t.Errorf("unexpected data structure: %+v", d1)
		return
	}

	// Check if an arbitrary key exists
	if x, ok := d1["two"]; ok {
		if x != data1["two"] {
			t.Errorf("Expected:%s Got:%s ", data1["two"], x)
		}
	} else {
		t.Errorf("Expected: %+v Got: %+v", data1, d1)
	}

}

func TestGetData(t *testing.T) {
	type expected struct {
		result interface{}
		err    error
	}

	type test struct {
		testNum     int
		description string
		input       string
		expected    *expected
		setupFunc   func(t *testing.T, test *test) location.Handler
	}

	const (
		emptyPath = "/secret/services/ms/ALL/ms/mgmtsvc/mgmt/no-stuff"
		okPath    = "/secret/services/ms/ALL/ms/mgmtsvc/mgmt/stuff"
	)

	defaultSetup := func(t *testing.T, test *test) location.Handler {
		handler, err := GetHandler()
		if err != nil {
			t.Errorf("failed to get handler: %s", core.ErrorText(err))
			return nil
		}

		err = handler.PutData(test.input, test.expected.result)
		if err != nil {
			t.Errorf("failed to store data: %s", core.ErrorText(err))
			return nil
		}
		return handler
	}

	notFoundSetup := func(t *testing.T, test *test) location.Handler {
		handler, err := GetHandler()
		if err != nil {
			t.Errorf("failed to get handler: %s", core.ErrorText(err))
			return nil
		}
		return handler
	}

	var tests = []test{
		{testNum: 1,
			description: "successfully retrieve data",
			input:       fmt.Sprintf("memory://%s", okPath),
			expected: &expected{result: "hello",
				err: nil},
			setupFunc: defaultSetup},
		{testNum: 2,
			description: "data not present",
			input:       fmt.Sprintf("memory://%s", emptyPath),
			expected: &expected{result: nil,
				err: core.MakeErrorAt(HandlerID, core.ErrorNotFound,
					fmt.Sprintf("no data at: %s", emptyPath),
					"memory.(*memory).GetData() - handler.go(183)")},
			setupFunc: notFoundSetup},
	}

	for _, test := range tests {
		handler := test.setupFunc(t, &test)
		if handler == nil {
			continue
		}
		result, err := handler.GetData(test.input)
		if result != test.expected.result || core.ErrorText(err) != core.ErrorText(test.expected.err) || testutils.FailTests {
			t.Errorf("\nTest %d: %s\nInput1..:\n%s\nExpected:\n%+v\n%s\nGot.....:\n%+v\n%s",
				test.testNum, test.description, test.input, test.expected.result, core.ErrorText(test.expected.err), result, core.ErrorText(err))
		}
	}

}
