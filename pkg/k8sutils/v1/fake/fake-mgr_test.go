// (c) Copyright 2018-2019 Hewlett Packard Enterprise Development LP

package fake

import (
	"testing"

	"k8s.io/client-go/kubernetes"

	"github.hpe.com/platform-core/utils/pkg/k8sutils/v1/testutils"
)

func TestGetClientset(t *testing.T) {
	// Test 1 - verify it works when passed vaild config data
	utils := &Fake{}
	testutils.Setup(t, utils, testutils.TestConfigFile)
	client := utils.GetClientset()
	if client == nil {
		t.Errorf("Unexpected error GetclientSet returned nil")
	}
	if _, ok := client.(kubernetes.Interface); !ok {
		t.Errorf("Unexpected error GetclientSet failed to return kubernetes.Interface")
	}

	// Test 2 - verify it fails with no config data provided
	utils = &Fake{}
	testutils.Setup(t, utils, testutils.NoConfig)
	if client = utils.GetClientset(); client != nil {
		t.Errorf("Unexpected error GetclientSet returned non nil")
	}
}
