// (c) Copyright 2018-2019 Hewlett Packard Enterprise Development LP

// Package location provides an interface to location handlers
// Location is a URI of the format...
// URI = scheme://userinfo@host:port/path?query#fragment
// Where scheme:// and path are required but other components are optional
// or mandatory based on the scheme. For example: vault scheme requires
// a userinfo of the form '<service>,<cn>'
package location

type (
	// Handler defines the interface of the Location Handler
	Handler interface {
		// Connect setups up connection to persistentStore
		Connect(string) error
		// PutData stores data in the location
		PutData(string, interface{}) error
		// GetData retrieves from the location
		GetData(string) (interface{}, error)
		// DeleteData deletes from the location
		DeleteData(string) error
		// ListData lists items at a location
		ListData(string) ([]string, error)
		// ID returns a text representation of the handler
		ID() string
		// Scheme returns the scheme the handler manages
		Scheme() string
		// VerifyScheme verifies a URI string against URI schema spec
		VerifyScheme(uri string) error
	}
)

const (
	// ErrorStringURISchemeMismatch The scheme provided in the URI doesn't
	// match the scheme implemented by the handler.
	ErrorStringURISchemeMismatch string = "wrong uri scheme"
	// ErrorStringURIParseFail The provided URI is malformed and
	// couldn't be parsed
	ErrorStringURIParseFail string = "failed to parse uri"
	// ErrorStringGetDataFail Reading of data from the backend failed
	ErrorStringGetDataFail string = "failed to get data from backend"
	// ErrorStringDeleteDataFail Deletion of data from the backend failed
	ErrorStringDeleteDataFail string = "failed to delete data from backend"
	// ErrorStringListDataFail List of data from the backend failed
	ErrorStringListDataFail string = "failed to list data from backend"
	// ErrorStringPutDataFail Writing data to the backend failed
	ErrorStringPutDataFail string = "failed to put data to backend"
	// ErrorNotImplemented The handler scheme is not yet implemented
	ErrorNotImplemented string = "handler scheme is not yet implemented"
)
