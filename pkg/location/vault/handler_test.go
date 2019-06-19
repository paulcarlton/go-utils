//(c) Copyright 2019 Hewlett Packard Enterprise Development LP

package vault

import (
	"strings"
	"testing"

	"github.hpe.com/platform-core/utils/pkg/location"

	sidelib "github.hpe.com/ncs-security/sec-side-golib"
)

// Mocks

type MockSessionUserInfo struct {
	sidelib.SessionUserInfoI
}

func (m *MockSessionUserInfo) StoreSecretByPath(path string, data map[string]interface{}) error {
	return nil
}

// This function mocks a GetSecret call from sec-side-car and
// returns a map[string]interface{}
func (m *MockSessionUserInfo) GetSecret(path string) (map[string]interface{}, error) {
	d := make(map[string]string)
	d["two"] = "2"

	sd := make(map[string]interface{})
	sd[vaultLocationData] = d

	return sd, nil
}

type MockSidelibS struct {
	sidelibI
}

func (side *MockSidelibS) newSession(string, string, string) (sidelib.SessionUserInfoI, error) {
	return &MockSessionUserInfo{}, nil
}

// Tests

func TestVaultLocationHandlerVerifyScheme(t *testing.T) {
	var v vault
	err := v.VerifyScheme("test.local")

	// Since the uri scheme is not vault this should error out
	// So if no error we want to fail the test
	if err == nil {
		t.Errorf("%s", err)
	}

	// The error should tell what
	if !strings.Contains(err.Error(), location.ErrorStringURISchemeMismatch) {
		t.Errorf("Returned Error: %s", err)
	}
}

func TestVaultLocation_Handler_checkUserInfo(t *testing.T) {
	var v vault

	// Check it fails if userInfo is not of the form string,string
	_, _, err := v.parseUserInfo("test://user@/secret/location")
	if err == nil {
		t.Errorf("Expected error but got nil")
	}

	// Check it passes if userInfo is of the form string,string
	_, _, err = v.parseUserInfo("test://s,c@/secret/location")
	if err != nil {
		t.Errorf("%s", err)
	}

}

func TestVaultLocationHandler_getSession(t *testing.T) {
	var v vault
	v.side = &MockSidelibS{}

	uri := "vault://ms,mgmtsvc@"

	_, err := v.getSession(uri)
	if err != nil {
		t.Errorf("Session error: %s", err)
	}
}

func TestVaultLocationHandlerConnect(t *testing.T) {
	var v vault
	v.side = &MockSidelibS{}

	uri := "vault://ms,mgmtsvc@"
	if err := v.Connect(uri); err != nil {
		t.Errorf("Connect error: %s", err)
	}
}

func TestVaultLocationHandler_PutGetData(t *testing.T) {

	var v vault
	v.side = &MockSidelibS{}

	uri := "vault://ms,mgmtsvc@/secret/services/ms/ALL/ms/mgmtsvc/mgmt/test1"
	data1 := map[string]string{
		"one":  "1",
		"two":  "2",
		"four": "4",
		"five": "5",
	}

	if err := v.PutData(uri, data1); err != nil {
		t.Errorf("PutData error: %s", err)
	}

	uri = "vault://ms,mgmtsvc@/secret/services/ms/ALL/ms/mgmtsvc/mgmt/test1"
	data2, err := v.GetData(uri)
	if err != nil {
		t.Errorf("GetData error: %s", err)
	}

	d1, ok := data2.(map[string]string)
	if !ok {
		t.Errorf("unexpected data structure: %+v", d1)
		return
	}

	// Check if an arbitrary key exists
	if x, ok := d1["two"]; ok {
		if x != data1["two"] {
			t.Errorf("Expected:%s Got:%s ", data1["two"], x)
		}
	} else {
		t.Errorf("Expected: %+v Got: %+v", data1, d1)
	}

}
