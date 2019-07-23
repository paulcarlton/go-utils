// (c) Copyright 2018-2019 Hewlett Packard Enterprise Development LP

package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

// JSONtext generates a string containing a json representation of an interface
func JSONtext(i interface{}) string {
	details := fmt.Sprintf("json for %+v...\n", i)
	if jsonData, err := json.Marshal(i); err != nil {
		details = details + fmt.Sprintf("json marshal error: %s\n", err)
	} else {
		if jsonText, err := PrettyJSON(string(jsonData)); err != nil {
			details = details + fmt.Sprintf("json format error: %s\n", err)
		} else {
			details = details + fmt.Sprintf("%s\n", jsonText)
		}
	}
	return details
}

// RequestDebug generates a string containing details of a request
func RequestDebug(r *http.Request) string {
	debugText := fmt.Sprintf("URL: %+v\n", r.URL)
	buf, _ := ioutil.ReadAll(r.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
	data, e := ioutil.ReadAll(rdr1)
	if e != nil {
		debugText = debugText + fmt.Sprintf("error reading body, %s", e)
	} else {
		debugText = debugText + fmt.Sprintf("Body..\n%s\n", string(data))
	}
	r.Body = rdr2 // OK since rdr2 implements the io.ReadCloser interface
	return debugText
}

// Callers returns an array of strings containing the function name, source filename and line
// number for the caller of this function and its caller moving up the stack for as many levels as
// are available or the number of levels specified by the levels parameter.
// Set the short parameter to true to only return final element of Function and source file name.
func Callers(levels uint, short bool) ([]string, error) {
	var callers []string
	if levels == 0 {
		return callers, nil
	}
	// we get the callers as uintptrs
	fpcs := make([]uintptr, levels)

	// skip 1 levels to get to the caller of whoever called Callers()
	n := runtime.Callers(1, fpcs)
	if n == 0 {
		return nil, fmt.Errorf("caller not availalble")
	}

	frames := runtime.CallersFrames(fpcs)
	for {
		frame, more := frames.Next()
		if frame.Line == 0 {
			break
		}
		funcName := frame.Function
		sourceFile := frame.File
		lineNumber := frame.Line
		if short {
			funcName = filepath.Base(funcName)
			sourceFile = filepath.Base(sourceFile)
		}
		caller := fmt.Sprintf("%s() - %s(%d)", funcName, sourceFile, lineNumber)
		callers = append(callers, caller)
		if !more {
			break
		}
	}
	return callers, nil
}

// GetCaller returns the caller of GetCaller 'skip' levels back
// Set the short parameter to true to only return final element of Function and source file name
func GetCaller(skip uint, short bool) string {
	callers, err := Callers(skip, short)
	if err != nil {
		return fmt.Sprintf("not available, %s", err)
	}
	if skip == 0 {
		return fmt.Sprintf("not available")
	}
	if int(skip) > len(callers) {
		return fmt.Sprintf("not available")
	}
	return callers[skip-1]
}

// ToJSON is used to convert a data structure into JSON format.
func ToJSON(data interface{}) (string, error) {

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, jsonData, "", "\t")
	if err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

// FindInStringSlice looks for a string in a slice of strings
// returns the index of the first instance of the string in the slice or -1 if not present
func FindInStringSlice(array []string, str string) int {
	for index, data := range array {
		if data == str {
			return index
		}
	}
	return -1
}

// CompareAsJSON compares two interfaces by converting them to json and comparing json text
func CompareAsJSON(one, two interface{}) bool {
	if one == nil && two == nil {
		return true
	}
	jsonOne, err := ToJSON(one)
	if err != nil {
		return false
	}

	jsonTwo, err := ToJSON(two)
	if err != nil {
		return false
	}
	return jsonOne == jsonTwo
}

// CompareStringSlices compares two strings by sorting them and comparing results
func CompareStringSlices(one, two []string) bool {
	if len(one) != len(two) {
		return false
	}
	sort.Strings(one)
	sort.Strings(two)
	return strings.Join(one[:], "") == strings.Join(two[:], "")
}

// PrettyJSON is used to format JSON
func PrettyJSON(data string) (string, error) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(data), "", "\t")
	if err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

var sleepSecsFunc = sleepSecs

func sleepSecs(seconds uint) {
	time.Sleep(time.Duration(int64(time.Second) * int64(seconds)))
}

// ExponentialDelay sleeps for wait time seconds
// it returns the wait time multiplied by 2 or 1 if the new wait is greater than max
// The intended usage is to pass a wait time and use the returned value in the next call
// This causes increasing waits up to max
func ExponentialDelay(waitTime, max uint) uint {
	sleepSecsFunc(waitTime)
	waitTime += waitTime
	if waitTime > max {
		waitTime = 1
	}
	return waitTime
}
