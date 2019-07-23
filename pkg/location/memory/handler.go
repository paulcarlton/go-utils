//(c) Copyright 2019 Hewlett Packard Enterprise Development LP

// Package memory implements a location handler interface that uses
// a memory backend and can be instantiated by calling
// core/location/factory.SelectHandler
// with a URI that contains 'memory://" scheme.
// Required fields in the URI:
// Scheme: should be equal to "memory"
// Path: should be of the form "/secret/..."
package memory

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/paulcarlton/go-utils/pkg/core"
	"github.com/paulcarlton/go-utils/pkg/goutils"
	"github.com/paulcarlton/go-utils/pkg/location"
)

const (
	// HandlerScheme The scheme for the memory handler
	HandlerScheme string = "memory"
	// HandlerID The ID for the memory handler
	HandlerID string = "memory location handler"

	// ErrorConnectFail Failed to connect to memory
	ErrorConnectFail string = "failed to connect to memory"
)

// implements a Location interface
type memory struct {
	data
}

var v memory

type data map[string]interface{}

// PathInfo holds information about path contents
type PathInfo struct {
	Path     string   `json:"path"`
	PathList []string `json:"subpath-list"`
	ItemList []string `json:"item-list"`
}

// Init initializes
func (pathInfo *PathInfo) Init() {
	pathInfo.ItemList = []string{}
	pathInfo.PathList = []string{}
}

// ID id
func (memory *memory) ID() string {
	return HandlerID
}

// Scheme scheme
func (memory *memory) Scheme() string {
	return HandlerScheme
}

//GetHandler A factory method to return a memory handler object
func GetHandler() (location.Handler, error) {
	return &v, nil
}

// VerifyScheme verify scheme
func (memory *memory) VerifyScheme(uri string) error {
	uriParts, err := url.Parse(uri)
	if err != nil {
		return core.RaiseError("", core.ErrorInvalidInput, fmt.Sprintf("%s %s", location.ErrorStringURIParseFail, uri), err)
	}

	if uriParts.Scheme != memory.Scheme() {
		return core.MakeError(memory.ID(), core.ErrorInvalidInput, fmt.Sprintf("%s %s:", location.ErrorStringURISchemeMismatch, memory.Scheme()))
	}
	return nil
}

// getSession resuses an existing session or gets a new one
func (memory *memory) getSession(uri string) (data, error) {
	if err := memory.VerifyScheme(uri); err != nil {
		return nil, err
	}

	// Check if session map is initialized and if so check if session for this user is cached.
	if memory.data == nil {
		memory.data = make(data)
	}

	return memory.data, nil
}

// Connect performs a memory backend connection and sets up
// a session for a uri. Any subsequent calls to other operations
// such as GetData and PutData with the same user credentials
// will reuse this session if it hasn't expired.
func (memory *memory) Connect(uri string) error {
	if _, err := memory.getSession(uri); err != nil {
		return core.RaiseError(memory.ID(), core.ErrorUnknown, "failed to connect", err)
	}

	return nil
}

// list returns PathInfo for the supplied path in the memory store
func list(data *data, path string) *PathInfo {
	pathInfo := &PathInfo{}
	pathInfo.Init()
	for k := range *data {
		if strings.HasPrefix(k, path) {
			if len(k) == len(path) { // shouldn't occur since we don't store empty 'directories'
				continue
			}
			endPath := k[len(path)+1:] // Get path without path prefix
			index := strings.Index(endPath, "/")
			if index > 0 { // It is a 'directory'
				if goutils.FindInStringSlice(pathInfo.PathList, endPath[:index]) < 0 {
					pathInfo.PathList = append(pathInfo.PathList, endPath[:index])
				}
				continue
			}
			pathInfo.ItemList = append(pathInfo.ItemList, endPath)
		}
	}
	pathInfo.Path = path
	return pathInfo
}

// ListData lists data at uri from the memory backend
func (memory *memory) ListData(uri string) ([]string, error) {
	uriParts, err := url.Parse(uri)
	if err != nil {
		return nil, core.RaiseError(memory.ID(), core.ErrorInvalidInput, ErrorConnectFail, err)
	}

	session, err := memory.getSession(uri)
	if err != nil {
		return nil, core.RaiseError(memory.ID(), core.ErrorUnknown, ErrorConnectFail, err)
	}

	return list(&session, uriParts.Path).ItemList, nil
}

// DeleteData deletes data for a uri from the memory backend
func (memory *memory) DeleteData(uri string) error {
	uriParts, err := url.Parse(uri)
	if err != nil {
		return core.RaiseError(memory.ID(), core.ErrorInvalidInput, ErrorConnectFail, err)
	}

	session, err := memory.getSession(uri)
	if err != nil {
		return core.RaiseError(memory.ID(), core.ErrorUnknown, ErrorConnectFail, err)
	}

	if _, ok := session[uriParts.Path]; !ok {
		return nil
	}
	delete(session, uriParts.Path)

	return nil
}

// GetData returns data for a uri from the memory backend
func (memory *memory) GetData(uri string) (interface{}, error) {
	uriParts, err := url.Parse(uri)
	if err != nil {
		return nil, core.RaiseError(memory.ID(), core.ErrorInvalidInput, ErrorConnectFail, err)
	}

	session, err := memory.getSession(uri)
	if err != nil {
		return nil, core.RaiseError(memory.ID(), core.ErrorUnknown, ErrorConnectFail, err)
	}

	if value, ok := session[uriParts.Path]; ok {
		return value, nil
	}

	return nil, core.MakeError(memory.ID(), core.ErrorNotFound, fmt.Sprintf("no data at: %s", uriParts.Path))
}

// PutData sets data value for a uri into the memory backend
func (memory *memory) PutData(uri string, data interface{}) error {
	uriParts, err := url.Parse(uri)
	if err != nil {
		return core.RaiseError(memory.ID(), core.ErrorUnknown, location.ErrorStringURIParseFail, err)
	}

	session, err := memory.getSession(uri)
	if err != nil {
		return core.RaiseError(memory.ID(), core.ErrorUnknown, ErrorConnectFail, err)
	}

	session[uriParts.Path] = data

	return nil
}
