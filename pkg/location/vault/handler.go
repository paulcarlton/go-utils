//(c) Copyright 2019 Hewlett Packard Enterprise Development LP

// Package vault implements a location handler interface that uses
// NCS Secmgr vault as a backend. It wraps sec-side-car and uses
// it for authentication against the secmgr. An instance of a vault
// struct can be obtained by calling location/factory/SelectHandler
// with a URI that contains 'vault://" scheme.
// Required fields in the URI:
// Scheme: should be equal to location.HandlerSchemeVault
// Userinfo: should be of the form '<service>,<component>'. Eg 'ms,mgmtsvc'
// Path: should be of the form "/secret/..."
package vault

import (
	"fmt"
	"net/url"
	"strings"

	sidelib "github.hpe.com/ncs-security/sec-side-golib"
	"github.hpe.com/platform-core/utils/pkg/core"
	"github.hpe.com/platform-core/utils/pkg/location"
)

const (
	// HandlerScheme The scheme for the vault handler
	HandlerScheme string = "vault"
	// HandlerID The ID for the vault handler
	HandlerID string = "vault location handler"

	// Name of the key with which the data will be stored if PutData
	// is used to write the data
	vaultLocationData = "VAULT_LOC_DATA"

	// ErrorConnectFail Failed to connect to vault
	ErrorConnectFail string = "failed to connect to vault"
	// ErrorInvalidUserInfo Wrong userinfo supplied
	ErrorInvalidUserInfo string = "the supplied user info is incorrect. Should be of the form <service>,<component>"
)

type sidelibI interface {
	sidelib.SessionUserInfoI
	newSession(name, service, cn string) (sidelib.SessionUserInfoI, error)
}

type sidelibS struct {
	sidelibI
}

func (side *sidelibS) newSession(name, service, cn string) (sidelib.SessionUserInfoI, error) {
	s, err := sidelib.NewSession(name, service, cn)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// implements a Location interface
type vault struct {
	location.Handler
	session map[string]sidelib.SessionUserInfoI // Default nil. Will be initialized on first getSession
	side    sidelibI
}

// ID id
func (vault *vault) ID() string {
	return HandlerID
}

// Scheme scheme
func (vault *vault) Scheme() string {
	return HandlerScheme
}

//GetHandler A factory method to return a vault handler object
func GetHandler() (location.Handler, error) {
	var v vault
	return &v, nil
}

// VerifyScheme verify scheme
func (vault *vault) VerifyScheme(uri string) error {
	uriParts, err := url.Parse(uri)
	if err != nil {
		return core.RaiseError("", core.ErrorUnknown, fmt.Sprintf("%s %s", location.ErrorStringURIParseFail, uri), err)
	}

	if uriParts.Scheme != vault.Scheme() {
		return core.MakeError(vault.ID(), core.ErrorInvalidInput, fmt.Sprintf("%s %s:", location.ErrorStringURISchemeMismatch, vault.Scheme()))
	}
	return nil
}

// parseUserInfo parses a well-formed URI and extracts the
// user field and splits it into service and component
func (vault *vault) parseUserInfo(uri string) (service, component string, err error) {
	uriParts, err := url.Parse(uri)
	if err != nil {
		return "", "", core.RaiseError("", core.ErrorUnknown, location.ErrorStringURIParseFail, err)
	}

	//The userinfo field is a composite field separated by comma. Eg. 'ms,mgmtsvc'
	userInfo := strings.Split(uriParts.User.Username(), ",")

	if len(userInfo) < 2 {
		return "", "", core.MakeError(vault.ID(), core.ErrorInvalidInput, fmt.Sprintf("%s %s:", ErrorInvalidUserInfo, uri))
	}

	return userInfo[0], userInfo[1], nil
}

func (vault *vault) getSession(uri string) (sidelib.SessionUserInfoI, error) {
	if err := vault.VerifyScheme(uri); err != nil {
		return nil, err
	}

	if vault.side == nil {
		vault.side = &sidelibS{}
	}

	service, component, err := vault.parseUserInfo(uri)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("%s,%s", service, component) //Session cache key

	// Check if session map is initialized and if so check if session for this user is cached.
	if vault.session == nil {
		vault.session = make(map[string]sidelib.SessionUserInfoI)
	}

	if _, found := vault.session[u]; !found {
		session, err := vault.side.newSession("", service, component)
		if err != nil {
			return nil, core.RaiseError(vault.ID(), core.ErrorUnknown, ErrorConnectFail, err)
		}

		vault.session[u] = session
	}

	return vault.session[u], nil
}

// Connect performs a vault backend connection and sets up
// a session for a uri. Any subsequent calls to other operations
// such as GetData and PutData with the same user credentials
// will reuse this session if it hasn't expired.
func (vault *vault) Connect(uri string) error {
	if _, err := vault.getSession(uri); err != nil {
		return core.RaiseError(vault.ID(), core.ErrorUnknown, "failed to connect", err)
	}

	return nil
}

// ListData lists data at uri from the vault backend
func (vault *vault) ListData(uri string) ([]string, error) {
	uriParts, err := url.Parse(uri)
	if err != nil {
		return nil, core.RaiseError(vault.ID(), core.ErrorUnknown, ErrorConnectFail, err)
	}

	session, err := vault.getSession(uri)
	if err != nil {
		return nil, core.RaiseError(vault.ID(), core.ErrorUnknown, ErrorConnectFail, err)
	}

	secretsInfo, err := session.GetSecrets(uriParts.Path)
	if err != nil {
		return nil, core.RaiseError(vault.ID(), core.ErrorUnknown, location.ErrorStringDeleteDataFail, err)
	}
	return secretsInfo.SecretList, nil
}

// DeleteData deletes data for a uri from the vault backend
func (vault *vault) DeleteData(uri string) error {
	uriParts, err := url.Parse(uri)
	if err != nil {
		return core.RaiseError(vault.ID(), core.ErrorUnknown, ErrorConnectFail, err)
	}

	session, err := vault.getSession(uri)
	if err != nil {
		return core.RaiseError(vault.ID(), core.ErrorUnknown, ErrorConnectFail, err)
	}

	if err := session.DeleteSecret(uriParts.Path); err != nil {
		return core.RaiseError(vault.ID(), core.ErrorUnknown, location.ErrorStringDeleteDataFail, err)
	}

	return nil
}

// GetData returns data for a uri from the vault backend
func (vault *vault) GetData(uri string) (interface{}, error) {
	uriParts, err := url.Parse(uri)
	if err != nil {
		return nil, core.RaiseError(vault.ID(), core.ErrorUnknown, ErrorConnectFail, err)
	}

	session, err := vault.getSession(uri)
	if err != nil {
		return nil, core.RaiseError(vault.ID(), core.ErrorUnknown, ErrorConnectFail, err)
	}

	data, err := session.GetSecret(uriParts.Path)
	if err != nil {
		return nil, core.RaiseError(vault.ID(), core.ErrorUnknown, location.ErrorStringGetDataFail, err)
	}

	// If this data was written by us the data stored should have been of the form data.(map[string]interface{})
	// and a key vaultLocationData would have existed. Check that and strip it off before sending back to the
	// caller. If this data wasn't written by us, send the data as it is.
	d1, ok := data[vaultLocationData]
	if ok {
		return d1, nil
	}

	return data, nil
}

// PutData sets data value for a uri into the vault backend
func (vault *vault) PutData(uri string, data interface{}) error {
	uriParts, err := url.Parse(uri)
	if err != nil {
		return core.RaiseError(vault.ID(), core.ErrorUnknown, location.ErrorStringURIParseFail, err)
	}

	session, err := vault.getSession(uri)
	if err != nil {
		return core.RaiseError(vault.ID(), core.ErrorUnknown, ErrorConnectFail, err)
	}

	data1 := make(map[string]interface{})

	// The data will be stored in Vault under key VAULT_LOC_DATA in
	// the path obtained from uri
	data1[vaultLocationData] = data

	if err := session.StoreSecretByPath(uriParts.Path, data1); err != nil {
		return core.RaiseError(vault.ID(), core.ErrorUnknown, location.ErrorStringPutDataFail, err)
	}

	return nil
}
