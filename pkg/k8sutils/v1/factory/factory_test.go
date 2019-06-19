// (c) Copyright 2019 Hewlett Packard Enterprise Development LP

package factory

import (
	"testing"

	k8sutilsv1 "github.hpe.com/platform-core/utils/pkg/k8sutils/v1"
)

func TestGetk8sUtils(t *testing.T) {
	// Test 1 - get Fake implementation
	utils, err := Getk8sUtils(FakeImpl)
	if err != nil {
		t.Errorf("Unexpected an error from GetK8sUtils")
	}
	result, ok := utils.(k8sutilsv1.K8sUtils)
	if !ok {
		t.Errorf("Unexpected error GetK8sUtils failed return implementation")
	}
	if result.Name() != FakeImpl {
		t.Errorf("Unexpected implementation type returned by GetK8sUtils")
	}
	// Test 2 - get K8s implementation
	utils, err = Getk8sUtils(K8sImpl)
	if err != nil {
		t.Errorf("Unexpected an error from GetK8sUtils")
	}
	result, ok = utils.(k8sutilsv1.K8sUtils)
	if !ok {
		t.Errorf("Unexpected error GetK8sUtils failed return implementation")
	}
	if result.Name() != K8sImpl {
		t.Errorf("Unexpected implementation type returned by GetK8sUtils")
	}

	// Test 3 - get invalid implementation
	if _, err = Getk8sUtils("none"); err == nil {
		t.Errorf("Expected an error from GetK8sUtils when called with invalid implementation type")
	}
}
