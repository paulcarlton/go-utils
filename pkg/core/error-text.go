// (c) Copyright 2019 Hewlett Packard Enterprise Development LP

package core

import (
	"fmt"
)

// ErrorText returns full details for core.Error and tries to send
// stringable data for other types of errors or types. This function
// can be extended to handle other types of error objects when needed.
func ErrorText(e interface{}) string {
	// Is this a core.Error?
	if err, ok := e.(Error); ok {
		return err.FullInfo()
	}

	// Is this a std error?
	if err, ok := e.(error); ok {
		return err.Error()
	}

	// Is this a stringable object?
	if text, ok := e.(fmt.Stringer); ok {
		return text.String()
	}

	// Format the object and fields
	return fmt.Sprintf("%+v", e)
}
