// (c) Copyright 2019 Hewlett Packard Enterprise Development LP

package core

import (
	"fmt"
	"net/http"
	"strings"

	"github.hpe.com/platform-core/utils/pkg/internal/common"
)

type (
	// Error should be used to report errors
	// A number of receiver methods and helper functions are available to simplify error creation and
	// reporting.
	Error interface {
		error
		ID() string
		SetCode() error
		Code() int
		SetMessage(message string) error
		Message() string
		AddDetails(details string) error
		Details() string
		AddRecommendedActions(actions ...string) error
		FullInfo() string
		Where() string
		Nested() error
		RecommendedActions() []string
	}

	cerror struct {
		// code: an opaque string uniquely identifying the error for programmatic use
		code int
		// message: clear and concise description of the error condition
		message string
		// details: optional verbose description of the error condition
		details string
		// recommendedActions: steps that a user can perform to correct the error condition
		recommendedActions []string
		// where: contains details of function and source line the error occurred at
		where string
		// id: contains the subject identifier related to the error
		id string
		// nestedError: subsidiary error that led to this error condition
		nestedError error
	}
)

const (
	// ErrorUnknown indicates the error cannot be categorized
	ErrorUnknown = 466
	// ErrorBadRequest indicates a problem with the request
	ErrorBadRequest = http.StatusBadRequest
	// ErrorDuplicateEntry  indicates a failure due to an unexpected duplicate
	ErrorDuplicateEntry = http.StatusConflict
	// ErrorInternal indicates an internal error
	ErrorInternal = http.StatusInternalServerError
	// ErrorInvalidInput indicates one or more input item is invalid
	ErrorInvalidInput = http.StatusUnprocessableEntity
	// ErrorNotFound indicates that the item specificed is not found
	ErrorNotFound = http.StatusNotFound
	// ErrorNotAllowed indicates that this operation is not allowed
	ErrorNotAllowed = http.StatusNotAcceptable
	// ErrorUnauthorized indicates that the caller is not authorized to perform the operation
	ErrorUnauthorized = http.StatusUnauthorized
	// ErrorServiceUnavailable indicates the service is unavailable at present
	ErrorServiceUnavailable = http.StatusServiceUnavailable
	// ErrorNotImplemented indicates the requested information or action is not implemented
	ErrorNotImplemented = http.StatusNotImplemented

	// private error constants
	nilErrorObjectPassed string = "called with a nil error object"
)

var statusText = map[int]string{
	ErrorUnknown: "Unknown Error",
}

// CodeText returns a text for the cor error code. It returns the empty string if the code is not defined
func CodeText(code int) string {
	if httpText := http.StatusText(code); len(httpText) > 0 {
		return httpText
	}
	return statusText[code]
}

// SetCode sets the core.Error Code based on search of error text
// The current implementation only checks for 'permission error' which vault emits
// The idea is that as we encounter other cases of error text we can use to set the error
// code they will be added
func (e *cerror) SetCode() error {
	if e == nil {
		return fmt.Errorf(nilErrorObjectPassed)
	}

	if e.code == ErrorUnknown &&
		strings.Contains(e.message, "permission error") {
		e.code = ErrorUnauthorized
	}

	return nil
}

// Code an opaque string uniquely identifying the error for programmatic or reference usage
func (e *cerror) Code() int {
	if e != nil {
		return e.code
	}
	return ErrorUnknown
}

// SetMessage adds message to a core.Error
func (e *cerror) SetMessage(message string) error {
	if e == nil {
		return fmt.Errorf(nilErrorObjectPassed)
	}
	e.message = message
	return nil
}

// Message gets Message of a core.Error
func (e *cerror) Message() string {
	if e != nil {
		return e.message
	}
	return ""
}

// AddDetails adds details to a core.Error
func (e *cerror) AddDetails(details string) error {
	if e == nil {
		return fmt.Errorf(nilErrorObjectPassed)
	}
	e.details = details
	return nil
}

// Details gets details of a core.Error
func (e *cerror) Details() string {
	if e != nil {
		return e.details
	}
	return ""
}

// AddActions adds recommended actions to a core.Error
func (e *cerror) AddRecommendedActions(actions ...string) error {
	if e == nil {
		return fmt.Errorf(nilErrorObjectPassed)
	}
	e.recommendedActions = append(e.recommendedActions, actions...)
	return nil
}

func (e *cerror) addNested(nested error) {
	if e != nil {
		e.nestedError = nested
	}
}

// Error returns a string representation of core.Error and thus implements std error interface
func (e *cerror) Error() string {
	if e == nil {
		return ""
	}

	sep := " "
	if len(e.id) == 0 {
		sep = ""
	}
	errorText := fmt.Sprintf("%s%s%s %s %s", e.where, sep, e.id, CodeText(e.code), e.message)

	return errorText
}

// FullInfo reports all details of core.Error
func (e *cerror) FullInfo() string {
	if e == nil {
		return ""
	}

	errorText := e.Error()

	if len(e.details) > 0 {
		errorText = fmt.Sprintf("%s\n%s", errorText, e.details)
	}

	if len(e.recommendedActions) > 0 {
		errorText = fmt.Sprintf("%s\nRecommended actions...", errorText)
		for _, rec := range e.recommendedActions {
			errorText = fmt.Sprintf("%s\n%s", errorText, rec)
		}
	}

	if e.nestedError != nil {
		errorText = fmt.Sprintf("%s\nNested Errors...", errorText)
		nested := e.nestedError
		errorText = fmt.Sprintf("%s\n%s", errorText, ErrorText(nested))
	}
	return errorText
}

// ID contains the subject identifier related to the error
func (e *cerror) ID() string {
	if e != nil {
		return e.id
	}
	return ""
}

// Where contains details of function and source line the error occurred at
func (e *cerror) Where() string {
	if e != nil {
		return e.where
	}
	return ""
}

// Nested subsidiary error that led to this error condition. Call GetNested
// recursively to walk the error chain
func (e *cerror) Nested() error {
	if e != nil && e.nestedError != nil {
		return e.nestedError
	}
	return nil
}

// RecommendedActions steps that a user can perform to correct the error condition
func (e *cerror) RecommendedActions() []string {
	if e != nil {
		return e.recommendedActions
	}
	return nil
}

// MakeError creates a core.Error
func MakeError(id string, code int, msg string) error {
	return makeError(id, code, msg, common.GetCaller(4, true))
}

// MakeErrorAt creates a core.Error, setting the function and file/line to the value provided
func MakeErrorAt(id string, code int, msg, where string) error {
	return makeError(id, code, msg, where)
}

// RaiseError creates a core.Error from a nested error
func RaiseError(id string, code int, msg string, nested interface{}) error {
	return raiseError(id, code, msg, common.GetCaller(4, true), nested)
}

// RaiseErrorAt creates a core.Error from a nested error and file/line to the value provided
func RaiseErrorAt(id string, code int, msg, where string, nested interface{}) error {
	return raiseError(id, code, msg, where, nested)
}

// makeError creates a core.Error
func makeError(id string, code int, msg, where string) error {
	result := &cerror{}
	result.code = code
	result.message = msg
	result.id = id
	result.where = where
	result.recommendedActions = []string{}
	result.nestedError = nil
	if err := result.SetCode(); err != nil {
		return nil
	}
	return result
}

// raiseError creates a core.Error from a nested error and file/line to the value provided
func raiseError(id string, code int, msg, where string, nested interface{}) error {
	// Check if nested error is of type core.Error and return a core.Error
	if nestedCoreError, ok := nested.(*cerror); ok {
		err := makeError(id, nestedCoreError.code, msg, where)
		if err != nil {
			if cerr, ok := err.(*cerror); ok {
				cerr.addNested(nestedCoreError)
				return cerr
			}
		}
	}

	var errMsg string

	// The nested error is probably a standard error
	if nestedErr, ok := nested.(error); ok {
		errMsg = fmt.Sprintf("%s, %s", msg, nestedErr)
	} else {
		// At this point the input parameter called nested is not something that implements std error.
		// So let fmt do it's magic
		errMsg = fmt.Sprintf("%s, %v", msg, nested)
	}

	return makeError(id, code, errMsg, where)
}
