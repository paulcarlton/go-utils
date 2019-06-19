// (c) Copyright 2019 Hewlett Packard Enterprise Development LP

package factory

import (
	"fmt"

	"github.hpe.com/platform-core/utils/pkg/core"
	k8sutilsv1 "github.hpe.com/platform-core/utils/pkg/k8sutils/v1"
	"github.hpe.com/platform-core/utils/pkg/k8sutils/v1/fake"
	"github.hpe.com/platform-core/utils/pkg/k8sutils/v1/k8s"
)

const (
	// K8sImpl is the name of the k8s implementation
	K8sImpl = "k8s"
	// FakeImpl is the name of the fake implementation
	FakeImpl = "fake"
)

// Getk8sUtils returns a k8sutils v1 implementation
func Getk8sUtils(implType string) (k8sutilsv1.K8sUtils, error) {
	switch implType {
	case K8sImpl:
		return &k8s.K8s{K8sUtilsImpl: k8sutilsv1.K8sUtilsImpl{ImplName: K8sImpl}}, nil
	case FakeImpl:
		return &fake.Fake{K8s: k8s.K8s{K8sUtilsImpl: k8sutilsv1.K8sUtilsImpl{ImplName: FakeImpl}}}, nil
	default:
		return nil, core.MakeError("", core.ErrorInvalidInput, fmt.Sprintf("%s is not a vaild k8sUtils.v1 implementation type", implType))
	}
}
