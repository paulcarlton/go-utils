// (c) Copyright 2018-2019 Hewlett Packard Enterprise Development LP

package logutils

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	_ "github.com/sirupsen/logrus"
	log "github.hpe.com/kronos/kelog"
	"github.com/paulcarlton/go-utils/pkg/testutils"
)

type LogRecord struct {
	Level  string
	Format string
	Data   []interface{}
}

type LogTest struct {
	*MockLogger
	logRecords []LogRecord
}

// MockCloser is a mock closer
type MockCloser struct {
	io.Closer
	err error
}

func (mockCloser MockCloser) Close() error {
	return mockCloser.err
}

// Log writes a log record at correct level
func (logData *LogTest) Log(log *log.Logger) {
	var logFunc = func(string, ...interface{}) {}
	for _, logDat := range logData.logRecords {
		switch logDat.Level {
		case "debug":
			logFunc = log.Debugf
		case "info":
			logFunc = log.Infof
		case "warning":
			logFunc = log.Warnf
		case "error":
			logFunc = log.Errorf
		default:
			fmt.Printf("Invalid log level: %s", logDat.Level)
			os.Exit(1)
		}
		logFunc(logDat.Format, logDat.Data...)
	}
}

// String returns a text representation of
func (logData *LogTest) String() string {
	var text string
	for _, logData := range logData.logRecords {
		text = fmt.Sprintf("%s%s, %s\n", text, logData.Level, fmt.Sprintf(logData.Format, logData.Data...))
	}
	return text[:len(text)-1]
}

// Expected generates the log data expected to be emitted
func (logData *LogTest) Expected() [][]string {
	var records [][]string
	for _, logData := range logData.logRecords {
		record := []string{fmt.Sprintf("\"level\":\"%s\"", logData.Level),
			fmt.Sprintf("\"msg\":\"%s\"}", fmt.Sprintf(logData.Format, logData.Data...))}
		records = append(records, record)
	}
	return records
}

func testMockLog(t *testing.T, testData *LogTest, log *log.Logger) {
	testData.setLogOut(log)
	defer testData.restoreLogOut(log)

	testData.Log(log)

	expected := testData.Expected()

	testData.MockLogger.bufferToStr()
	results := *testData.MockLogger.logRecords
	for index, record := range results {
		result := strings.Split(record, ",")
		if !testutils.ContainsStringArray(result, expected[index], false) || testutils.FailTests {
			t.Errorf("\nInput...:\n%s\nExpected:\n%s\nGot.....:\n%s",
				testData.String(), expected[index], result)
			if !testutils.FailTests {
				return
			}
		}
	}
}

func TestMockLog(t *testing.T) {
	testLog := log.NewLogger()
	var tests = []LogTest{{MockLogger: nil, logRecords: []LogRecord{{Level: "info", Format: "%s", Data: []interface{}{"this is a info msg"}},
		{Level: "debug", Format: "%s", Data: []interface{}{"this is a debug msg"}}}},
		{MockLogger: nil, logRecords: []LogRecord{{Level: "warning", Format: "%s", Data: []interface{}{"this is a warn msg"}},
			{Level: "error", Format: "%s", Data: []interface{}{"this is a error msg"}}}},
	}

	for _, test := range tests {
		test.MockLogger = &MockLogger{}
		testMockLog(t, &test, testLog)
	}
}

type CloseWarnTest struct {
	testNum int
	*MockLogger
	closer   io.Closer
	expected []LogRec
}

func testCloseWarnLog(t *testing.T, test *CloseWarnTest) {
	test.setLogOut(log.GetLogger())
	defer test.restoreLogOut(log.GetLogger())

	CloseWarn(log.GetLogger(), test.closer)

	test.MockLogger.bufferToStr()
	result := test.MockLogger.logRecords
	if !ContainsLogRecords(StripLogRecords(*result), test.expected) {
		t.Errorf("\nTest: %d\nInput...: %+v\nExpected: %s\nGot.....: %s",
			test.testNum, test.closer, test.expected, *result)
	}
}

func TestCloseWarn(t *testing.T) {
	var tests = []CloseWarnTest{
		{1, &MockLogger{}, MockCloser{err: nil}, []LogRec{}},
		{2, &MockLogger{}, MockCloser{err: fmt.Errorf("because")}, []LogRec{{"warning", "close failed: because"}}},
	}

	for _, test := range tests {
		testCloseWarnLog(t, &test)
	}
}
