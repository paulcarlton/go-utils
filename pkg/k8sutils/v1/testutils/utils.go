// (c) Copyright 2019 Hewlett Packard Enterprise Development LP

package testutils

import (
	"io/ioutil"
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/clientcmd"

	"github.hpe.com/platform-core/utils/pkg/core"
	k8sutilsv1 "github.hpe.com/platform-core/utils/pkg/k8sutils/v1"
)

var (
	// TestSecret is a test secret
	TestSecret = &v1.Secret{
		Type: v1.SecretTypeOpaque,
		ObjectMeta: metav1.ObjectMeta{
			Name:      "TestSecret",
			Namespace: "Test",
		},
		Data: map[string][]byte{
			"info": []byte("asdasdads"),
		},
	}
	// TestBlankSecret is an empty secret
	TestBlankSecret = &v1.Secret{}

	// TestConfigMap is a test config map
	TestConfigMap = &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "TestConfigMap",
			Namespace: "Test",
		},
		Data: map[string]string{
			"info": "asdasdads",
		},
	}

	// TestBlankConfigMap is an empty configmap
	TestBlankConfigMap = &v1.ConfigMap{}

	// TestSecretByte is test secret data
	TestSecretByte = []byte(`{
    "apiVersion": "v1",
    "kind": "Secret",
    "metadata": {
      "name": "TestSecret",
      "namespace": "Test"
    },
    "type": "Opaque",
    "data": {
      "info": "YWRtaW4="
    }
  }`)

	// TestInvalidSecretByte is invalid test secret data
	TestInvalidSecretByte = []byte(`{
    "apiVersion": "v1",
    "kind": "Secret",
    "metadata": {
      "name": "TestSecret",
      "namespace": "Test"
    },
    "type": "Opaque",
    "data": {
      "info": "hello"
    }
  }`)

	// TestConfigMapByte is test config map data
	TestConfigMapByte = []byte(`{
    "apiVersion": "v1",
    "kind": "ConfigMap",
    "metadata": {
      "name": "TestConfigMap",
      "namespace": "Test"
    },
    "type": "Opaque",
    "data": {
      "info": "asdasdads"
    }
  }`)

	// TestInvalidByte is invalid test data
	TestInvalidByte = []byte(`{
    "apiVersion": "v1",
    "kind": "Invalid",
    "metadata": {
      "name": "TestInvalid",
      "namespace": "Test"
    },
    "type": "Opaque",
    "data": {
      "info": "asdasdads"
    }
	}`)
	// TestInvalidJSONByte is test invalid JSON data
	TestInvalidJSONByte = []byte(`{
		"apiVersion": "v1",
		"kind": "Invalid",
		"metadata": {
		  "name": "TestInvalid",
		  "namespace": "Test"
		}
		"type": "Opaque",
		"data": {
		  "info": "asdasdads"
		}
	  }`)
)

// Tester is a structure that implements the K8sUtils interface
// It is used by this package to facilitate testing and avoid cyclic dependencies
type tester struct {
	k8sutilsv1.K8sUtils
	Client kubernetes.Interface
}

// SetClientset validates the kubernetes config then sets the client to a fake kubernetes client
func (tester *tester) SetClientset(configFileData []byte) error {
	config, err := clientcmd.RESTConfigFromKubeConfig(configFileData)
	if err != nil {
		return core.RaiseError("", core.ErrorUnknown, "failed trying to build config", err)
	}
	if _, err = kubernetes.NewForConfig(config); err != nil {
		return core.RaiseError("", core.ErrorUnknown, "failed trying to get clientset", err)
	}
	tester.Client = fake.NewSimpleClientset()
	return nil
}

// Setup is used to setup a test
func Setup(t *testing.T, utils k8sutilsv1.K8sUtils, configFile string) {
	if len(configFile) == 0 {
		return
	}
	testConfig, err := ioutil.ReadFile(configFile) // nolint: gosec
	if err != nil {
		t.Logf("failed to read fixture file, %s", err)
	}

	if err := utils.SetClientset(testConfig); err != nil {
		t.Logf("failed to get clientset, %s", err)
	}
}

const (
	// TestConfigFile is the file path to file holding test kubeconfig
	TestConfigFile = "../testutils/tests/fixtures/kubeconfig.yaml"
	// NoConfig is an empty string used to pass and empty kubeconfig
	NoConfig = ""
)
