// (c) Copyright 2018-2019 Hewlett Packard Enterprise Development LP

package logutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.hpe.com/kronos/kelog"

	"github.com/paulcarlton/go-utils/pkg/testutils"
)

// Types

// MockLogger is a mock logger
type MockLogger struct {
	buffer     *bytes.Buffer
	logRecords *[]string
	realLogOut io.Writer // Used to old the original log output
}

// Receiver Methods

// MockLogger Methods

func (m *MockLogger) setLogOut(coreLog *log.Logger) {
	logger := logrus.New()
	m.buffer = bytes.NewBuffer(make([]byte, 0, 1024000))
	logger.Out = m.buffer
	m.realLogOut = coreLog.Out
	coreLog.Out = logger.Out
}

func (m *MockLogger) bufferToStr() {
	m.logRecords = testutils.ReadBuf(m.buffer)
}

func (m *MockLogger) restoreLogOut(coreLog *log.Logger) {
	coreLog.Out = m.realLogOut
}

// LogRec defines the log level and message
type LogRec struct {
	Level string
	Msg   string
}

type logData struct {
	File  string `json:"context,omitempty"`
	Level string `json:"level,omitempty"`
	Msg   string `json:"msg,omitempty"`
	Time  string `json:"@timestamp,omitempty"`
}

// StripLogRecords converts list of json format log records to a list of LogRec
func StripLogRecords(results []string) (stripped []LogRec) {
	for _, logLine := range results {
		text := strings.TrimLeft(logLine, " ")
		if !strings.HasPrefix(text, "{") {
			fmt.Printf("warning, invalid log record, missing {: %s", logLine)
			continue
		}
		logData := &logData{}
		if err := json.Unmarshal([]byte(text), logData); err != nil {
			fmt.Printf("warning, invalid log record: %s, failed to parse json, %s", logLine, err)
			continue
		}
		stripped = append(stripped, LogRec{logData.Level, logData.Msg})
	}
	return stripped
}

// ContainsLogRecords checks if the expected log records are in the log output
func ContainsLogRecords(results []LogRec, expected []LogRec) bool {
	if len(expected) == 0 {
		return true
	}
	offset := 0
	for _, logData := range results {
		if logData.Msg == expected[offset].Msg && logData.Level == expected[offset].Level {
			offset++
		}
		if offset == len(expected) {
			return true
		}
	}
	return false
}

// CloseWarn calls Close() on the io.Closer interface supplied and if an error occurs it emits a
// warning
func CloseWarn(log *log.Logger, toClose io.Closer) {
	if err := toClose.Close(); err != nil {
		log.Warnf("close failed: %s", err)
	}
}
