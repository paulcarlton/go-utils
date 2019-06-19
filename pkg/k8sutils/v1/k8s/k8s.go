// (c) Copyright 2019 Hewlett Packard Enterprise Development LP

package k8s

import (
	"strings"

	"github.com/ugorji/go/codec"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.hpe.com/platform-core/utils/pkg/core"
	k8sutilsv1 "github.hpe.com/platform-core/utils/pkg/k8sutils/v1"
)

// K8s is a structure that hold a kubernetes client and implements the K8sUtils interface
type K8s struct {
	k8sutilsv1.K8sUtilsImpl
}

// Name returns the name of the implementation
func (k8s *K8s) Name() string {
	return k8s.ImplName
}

// GetClientset gets the clientset connection for a given kubernetes config
func (k8s *K8s) GetClientset() kubernetes.Interface {
	return k8s.Client
}

// SetClientset gets a clientset connection for a given kubernetes config
func (k8s *K8s) SetClientset(configFileData []byte) error {
	config, err := clientcmd.RESTConfigFromKubeConfig(configFileData)
	if err != nil {
		return core.RaiseError("", core.ErrorUnknown, "failed trying to build config", err)
	}
	k8s.Client, err = kubernetes.NewForConfig(config)
	if err != nil {
		return core.RaiseError("", core.ErrorUnknown, "failed trying to get clientset", err)
	}
	return nil
}

// FindK8sSecret checks if a K8sCluster secret exists
func (k8s *K8s) FindK8sSecret(secret *v1.Secret) (bool, error) {

	_, err := k8sGetSecret(k8s, secret)

	// If the error is "not found" we just return false
	if err != nil && strings.Contains(err.Error(), "not found") {
		return false, nil
	} else if err != nil {
		return false, core.RaiseError(secret.Name, core.ErrorUnknown, "failed trying to find secret", err)
	}

	return true, nil
}

// CreateK8sSecret creates an new K8sCluster secret with the data supplied
func (k8s *K8s) CreateK8sSecret(secret *v1.Secret) error {

	if _, err := k8s.Client.CoreV1().Secrets(secret.Namespace).Create(secret); err != nil {
		return core.RaiseError(secret.Name, core.ErrorUnknown, "failed trying to create secret", err)
	}
	return nil
}

// UpdateK8sSecret updates an new K8sCluster secret with the data supplied
func (k8s *K8s) UpdateK8sSecret(secret *v1.Secret) error {

	if _, err := k8s.Client.CoreV1().Secrets(secret.Namespace).Update(secret); err != nil {
		return core.RaiseError(secret.Name, core.ErrorUnknown, "failed trying to update secret", err)
	}
	return nil
}

// DeleteK8sSecret deletes an existing K8sCluster secret
func (k8s *K8s) DeleteK8sSecret(secret *v1.Secret) error {

	if err := k8s.Client.CoreV1().Secrets(secret.Namespace).Delete(secret.Name, &metav1.DeleteOptions{}); err != nil {
		return core.RaiseError(secret.Name, core.ErrorUnknown, "failed trying to delete secret", err)
	}
	return nil
}

// GetK8sSecret gets an existing K8sCluster secret
func (k8s *K8s) GetK8sSecret(secret *v1.Secret) (*v1.Secret, error) {

	foundSecret, err := k8sGetSecret(k8s, secret)
	if err != nil {
		return foundSecret, core.RaiseError(secret.Name, core.ErrorUnknown, "failed trying to find secret", err)
	}
	return foundSecret, nil
}

// FindK8sConfigMap checks if a K8sCluster configmap exists
func (k8s *K8s) FindK8sConfigMap(configMap *v1.ConfigMap) (bool, error) {

	if _, err := k8s.Client.CoreV1().ConfigMaps(configMap.Namespace).Get(configMap.Name, metav1.GetOptions{});
	// If the error is "not found" we just return false
	err != nil && strings.Contains(err.Error(), "not found") {
		return false, nil
	} else if err != nil {
		return false, core.RaiseError(configMap.Name, core.ErrorUnknown, "failed trying to find configmap", err)
	}

	return true, nil
}

// CreateK8sConfigMap creates an new K8sCluster configmap with the data supplied
func (k8s *K8s) CreateK8sConfigMap(configMap *v1.ConfigMap) error {

	if _, err := k8s.Client.CoreV1().ConfigMaps(configMap.Namespace).Create(configMap); err != nil {
		return core.RaiseError(configMap.Name, core.ErrorUnknown, "failed trying to create configmap", err)
	}

	return nil
}

// UpdateK8sConfigMap updates an new K8sCluster configmap with the data supplied
func (k8s *K8s) UpdateK8sConfigMap(configMap *v1.ConfigMap) error {

	if _, err := k8s.Client.CoreV1().ConfigMaps(configMap.Namespace).Update(configMap); err != nil {
		return core.RaiseError(configMap.Name, core.ErrorUnknown, "failed trying to update configmap", err)
	}

	return nil
}

// DeleteK8sConfigMap deletes a K8sCluster configmap
func (k8s *K8s) DeleteK8sConfigMap(configMap *v1.ConfigMap) error {

	if err := k8s.Client.CoreV1().ConfigMaps(configMap.Namespace).Delete(configMap.Name, &metav1.DeleteOptions{}); err != nil {
		return core.RaiseError(configMap.Name, core.ErrorUnknown, "failed trying to delete configmap", err)
	}

	return nil
}

// getHandle returns a new json handler
func getHandle() *codec.JsonHandle {
	h := new(codec.JsonHandle)
	h.InternString = true
	return h
}

// getK8sSecret wraps the k8s client secret Get
func k8sGetSecret(k8s *K8s, secret *v1.Secret) (*v1.Secret, error) {
	foundSecret, err := k8s.Client.CoreV1().Secrets(secret.Namespace).Get(secret.Name, metav1.GetOptions{})
	return foundSecret, err
}

// decodeK8s is a K8s object data decoder. This will determine the object type
// and decode the data into a data structure which is returned. If an error
// occurs the error is returned.
func decodeK8s(k8sData []byte) (interface{}, error) {
	var typeMeta metav1.TypeMeta
	var coreErr error
	dec := codec.NewDecoderBytes(k8sData, getHandle())
	err := dec.Decode(&typeMeta)
	if err != nil {
		return nil, core.RaiseError("", core.ErrorUnknown, "failed to decode k8s data", err)
	}
	dec.ResetBytes(k8sData)

	var data interface{}
	switch typeMeta.Kind {
	case "ConfigMap":
		var configMap v1.ConfigMap
		if err = dec.Decode(&configMap); err != nil {
			coreErr = core.RaiseError("", core.ErrorUnknown, "failed to decode k8s configmap", err)
		}
		data = &configMap
	case "Secret":
		var secret v1.Secret
		if err = dec.Decode(&secret); err != nil {
			coreErr = core.RaiseError("", core.ErrorUnknown, "failed to decode k8s secret", err)
		}
		data = &secret
	default:
		// We do not need to decode this type.
		return nil, nil
	}

	// One of the decoders failed.
	if coreErr != nil {
		return nil, coreErr
	}

	return data, nil
}

// DecodeK8sSecret is a wrapper for decodeK8s that takes a byte stream of a secret
// and returns a v1.Secret
func DecodeK8sSecret(k8sData []byte) (*v1.Secret, error) {
	secretIf, err := decodeK8s(k8sData)
	secret := secretIf.(*v1.Secret)
	return secret, err
}

// DecodeK8sConfigMap is a wrapper for decodeK8s that takes a byte stream of a configmap
// and returns a v1.ConfigMap
func DecodeK8sConfigMap(k8sData []byte) (*v1.ConfigMap, error) {
	configMapIf, err := decodeK8s(k8sData)
	configMap := configMapIf.(*v1.ConfigMap)
	return configMap, err
}
