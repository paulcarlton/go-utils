// (c) Copyright 2018-2019 Hewlett Packard Enterprise Development LP

package testutils

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Test utilities

// FailTests determines if tests are set to failed so expected and results are emitted even if tests pass
var FailTests = false // Set to true to test failure output

// RemoveBottom removes items from caller list including and prior to testing.tRunner()
func RemoveBottom(callers []string) []string {
	ourCallers := []string{}
	for _, caller := range callers {
		if strings.HasPrefix(caller, "testing.tRunner()") {
			break
		}
		ourCallers = append(ourCallers, caller)
	}
	return ourCallers
}

// CompareWhereList compares a list of strings returned by GetCaller or Callers but ignores line numbers
func CompareWhereList(one, two []string) bool {
	if len(one) == 0 && len(two) == 0 { // Both empty return true
		return true
	}
	if len(two) != len(one) {
		return false
	}
	for index, field := range one {
		if !CompareWhere(field, two[index]) {
			return false
		}
	}
	return true
}

// CompareWhere compares strings returned by GetCaller or Callers but ignores line numbers
func CompareWhere(one, two string) bool {
	if strings.HasSuffix(one, "(NN)") || strings.HasSuffix(two, "(NN)") {
		return one[:strings.LastIndex(one, "(")] == two[:strings.LastIndex(two, "(")]
	}
	return one == two
}

// CompareItems compares two items, returning true if both are nil
func CompareItems(one, two interface{}) bool {
	if one == nil && two == nil {
		return true
	}
	if one == nil || two == nil {
		return false
	}
	return fmt.Sprintf("%+v", one) == fmt.Sprintf("%+v", two)
}

// ContainsStringArray checks items in two string arrays verifying that the second string array is
// contained in the first
func ContainsStringArray(one, two []string, first bool) bool {
	if len(one) == 0 && len(two) == 0 { // Both empty return true
		return true
	}
	if len(two) == 0 {
		return false
	}
	offset := 0
	for index, field := range one {
		if field == two[offset] {
			offset++
			index++
			if index < len(one) && offset < len(two) {
				return ContainsStringArray(one[index:], two[offset:], true)
			}
			if offset == len(two) {
				return true
			}
		}
		if first {
			return false
		}
	}
	return false
}

func init() {
	_, present := os.LookupEnv("FAILED_OUTPUT_TEST")
	if present {
		FailTests = true
	}
}

// ReadBuf reads lines from a buffer into a string array
func ReadBuf(out io.Reader) *[]string {
	output := []string{}
	for {
		p := make([]byte, 1024000)
		n, err := out.Read(p)
		if n > 0 {
			lines := strings.Split(string(p[:n]), "\n")
			output = append(output, lines...)
		}
		if err != nil {
			break
		}
	}
	// If last line is empty then remove it
	if len(output) > 0 && len(output[len(output)-1]) == 0 {
		output = output[:len(output)-1]
	}
	return &output
}

// DisplayStrings returns a []string as a single string, one element per line
func DisplayStrings(strs []string) string {
	output := ""
	newline := ""
	for index, str := range strs {
		output = fmt.Sprintf("%s%s%d - %s", output, newline, index, str)
		newline = "\n"
	}
	return output
}
