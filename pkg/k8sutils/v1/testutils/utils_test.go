// (c) Copyright 2019 Hewlett Packard Enterprise Development LP

package testutils

import (
	"testing"

	"k8s.io/client-go/kubernetes"
)

func TestSetup(t *testing.T) {
	// Test 1 - verify it works when passed vaild config data
	utils := &tester{}
	Setup(t, utils, TestConfigFile)
	if _, ok := utils.Client.(kubernetes.Interface); !ok {
		t.Errorf("Unexpected error Setup failed to set Client to kubernetes.Interface")
	}

	// Test 2 - verify it fails when passed invaild config data
	utils = &tester{}
	Setup(t, utils, NoConfig)
	if utils.Client != nil {
		t.Errorf("Unexpected error, Client shoud be set to nil")
	}
}
