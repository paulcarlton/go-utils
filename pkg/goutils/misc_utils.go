// (c) Copyright 2018-2019 Hewlett Packard Enterprise Development LP

package goutils

import (
	"github.hpe.com/platform-core/utils/pkg/core"
	"github.hpe.com/platform-core/utils/pkg/internal/common"
)

// Callers returns an array of strings containing the function name, source filename and line
// number for the caller of this function and its caller moving up the stack for as many levels as
// are available or the number of levels specified by the levels parameter.
// Set the short parameter to true to only return final element of Function and source file name.
func Callers(levels uint, short bool) ([]string, error) {
	return common.Callers(levels, short)
}

// GetCaller returns the caller of GetCaller 'skip' levels back
// Set the short parameter to true to only return final element of Function and source file name
func GetCaller(skip uint, short bool) string {
	return common.GetCaller(skip+1, short)
}

// ToJSON is used to convert a data structure into JSON format.
func ToJSON(data interface{}) (string, error) {
	return common.ToJSON(data)
}

// FindInStringSlice looks for a string in a slice of strings
// returns the index of the first instance of the string in the slice or -1 if not present
func FindInStringSlice(array []string, str string) int {
	return common.FindInStringSlice(array, str)
}

// CastToString casts an interface to a string if possible
func CastToString(i interface{}) (string, error) {
	if str, ok := i.(string); ok {
		return str, nil
	}
	return "", core.MakeError("", core.ErrorInvalidInput, "failed to cast to string")
}

// CompareAsJSON compares two interfaces by converting them to json and comparing json text
func CompareAsJSON(one, two interface{}) bool {
	return common.CompareAsJSON(one, two)
}

// CompareStringSlices compares two strings by sorting them and comparing results
func CompareStringSlices(one, two []string) bool {
	return common.CompareStringSlices(one, two)
}

// PrettyJSON is used to format JSON
func PrettyJSON(data string) (string, error) {
	return common.PrettyJSON(data)
}

// ExponentialDelay sleeps for wait time seconds
// it returns the wait time multiplied by 2 or 1 if the new wait is greater than max
// The intended usage is to pass a wait time and use the returned value in the next call
// This causes increasing waits up to max
func ExponentialDelay(waitTime, max uint) uint {
	return common.ExponentialDelay(waitTime, max)
}
