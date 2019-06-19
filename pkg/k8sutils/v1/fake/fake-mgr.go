// (c) Copyright 2018-2019 Hewlett Packard Enterprise Development LP

package fake

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/clientcmd"

	"github.hpe.com/platform-core/utils/pkg/core"
	"github.hpe.com/platform-core/utils/pkg/k8sutils/v1/k8s"
)

// Fake is a structure that hold a fake kubernetes client and implements the K8sUtils interface
type Fake struct {
	k8s.K8s
}

// SetClientset validates the kubernetes config then sets the client to a fake kubernetes client
func (k8s *Fake) SetClientset(configFileData []byte) error {
	config, err := clientcmd.RESTConfigFromKubeConfig(configFileData)
	if err != nil {
		return core.RaiseError("", core.ErrorUnknown, "failed trying to build config", err)
	}
	if _, err = kubernetes.NewForConfig(config); err != nil {
		return core.RaiseError("", core.ErrorUnknown, "failed trying to get clientset", err)
	}
	k8s.Client = fake.NewSimpleClientset()
	return nil
}
