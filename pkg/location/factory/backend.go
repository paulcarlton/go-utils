// (c) Copyright 2019 Hewlett Packard Enterprise Development LP

package factory

import (
	"fmt"
	"net/url"

	"github.hpe.com/platform-core/utils/pkg/core"
	"github.hpe.com/platform-core/utils/pkg/location"
	"github.hpe.com/platform-core/utils/pkg/location/memory"
	"github.hpe.com/platform-core/utils/pkg/location/vault"
)

const (
	id string = "location factory"
)

// SelectHandler returns the appropriate location
// handler that implements the scheme used in the URI.
// Currently only Vault handler is implemented but
// there can be others.
func SelectHandler(uri string) (location.Handler, error) {
	uriParts, err := url.Parse(uri)
	if err != nil {
		return nil, core.RaiseError(id, core.ErrorUnknown, fmt.Sprintf("%s %s:", uri, location.ErrorStringURIParseFail), err)
	}

	switch uriParts.Scheme {
	case vault.HandlerScheme:
		return vault.GetHandler()
	case memory.HandlerScheme:
		return memory.GetHandler()
	default:
		return nil, core.MakeError(id, core.ErrorInvalidInput, fmt.Sprintf("%s %s:", core.CodeText(core.ErrorNotImplemented), uriParts.Scheme))
	}
}
