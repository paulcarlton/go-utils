// (c) Copyright 2018-2019 Hewlett Packard Enterprise Development LP

package core

import "strings"

// CompareErrors compares two core.Error objects or
// two error objects, the latter being compared for error string
func CompareErrors(one, two error) bool {
	if one == nil && two == nil {
		return true
	} else if one == nil || two == nil {
		return false
	}

	oneE, ok1 := one.(Error)
	twoE, ok2 := two.(Error)

	if ok1 && ok2 {
		return compareErrors(oneE, twoE)
	}

	return one.Error() == two.Error()
}

// compareErrors compares two core.Error
func compareErrors(one, two Error) bool {
	if one.Code() != two.Code() ||
		one.Message() != two.Message() ||
		one.ID() != two.ID() ||
		one.Details() != two.Details() ||
		!compareWhere(one.Where(), two.Where()) ||
		!compareStringArray(one.RecommendedActions(), two.RecommendedActions()) {
		return false
	}

	// Compare Nested errors

	nested1 := one.Nested()
	nested2 := two.Nested()

	if nested1 == nil && nested2 == nil {
		return true
	} else if nested1 == nil || nested2 == nil {
		return false
	}

	// Nested errors can be any structure that conforms
	// standard error. Currently we support comparing
	// core.Error and std error. More error types can
	// be added in the selector block below when needed.

	switch nested1.(type) {
	case Error:
		if _, ok := nested2.(Error); ok {
			return CompareErrors(nested1.(Error), nested2.(Error))
		}
		return false // Since the errors are of different types

	default:
		return nested1.Error() == nested2.Error()

	}

}

// compareWhere compares strings returned by GetCaller or Callers but ignores line numbers
func compareWhere(one, two string) bool {
	if strings.HasSuffix(one, "(NN)") || strings.HasSuffix(two, "(NN)") {
		return one[:strings.LastIndex(one, "(")] == two[:strings.LastIndex(two, "(")]
	}
	return one == two
}

// compareStringArray
func compareStringArray(one, two []string) bool {
	if (one == nil) != (two == nil) {
		return false
	}
	if len(one) != len(two) {
		return false
	}
	for i := range one {
		if one[i] != two[i] {
			return false
		}
	}
	return true
}
