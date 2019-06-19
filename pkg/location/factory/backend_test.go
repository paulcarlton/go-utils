// (c) Copyright 2019 Hewlett Packard Enterprise Development LP

package factory

import (
	"strings"
	"testing"

	"github.hpe.com/platform-core/utils/pkg/core"
	"github.hpe.com/platform-core/utils/pkg/location"
)

func TestSelectHandlerErrors(t *testing.T) {
	l, _ := SelectHandler("vault://")
	if l.Scheme() != "vault" {
		t.Errorf("Got %s", l.Scheme())
	}

	_, err := SelectHandler("something://")
	if err == nil {
		t.Errorf("Expected error but got nil")
		return
	}

	// First run a generic compare
	if !strings.Contains(err.Error(), core.CodeText(core.ErrorNotImplemented)) {
		t.Errorf("Expected %s but got %s", location.ErrorNotImplemented, err)
	}

	// Test if the returned error type is as expected
	if _, ok := err.(core.Error); !ok {
		t.Errorf("Expected error Type core.Error but received error object didn't match this type")
		return
	}

}

func TestSelectHandlerSchemes(t *testing.T) {
	l, _ := SelectHandler("vault://")
	if l.Scheme() != "vault" {
		t.Errorf("Expected: vault Got:%s", l.Scheme())
	}
	l, _ = SelectHandler("memory://")
	if l.Scheme() != "memory" {
		t.Errorf("Expected: memory Got: %s", l.Scheme())
	}
}
