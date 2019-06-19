// (c) Copyright 2019 Hewlett Packard Enterprise Development LP

package v1

import (
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// K8sUtils is the interface for k8s cluster secrets and config-maps updater
type K8sUtils interface {
	// Name returns the name of the implemetantion
	Name() string
	// GetClientset gets the clientset connection for a given kubernetes config
	GetClientset() kubernetes.Interface

	// SetClientset sets the clientset connection for a given kubernetes config
	SetClientset(configFileData []byte) error

	// FindK8sSecret checks if a K8sCluster secret exists
	FindK8sSecret(secret *v1.Secret) (bool, error)

	// CreateK8sSecret creates an new K8sCluster secret with the data supplied
	CreateK8sSecret(secret *v1.Secret) error

	// UpdateK8sSecret updates an new K8sCluster secret with the data supplied
	UpdateK8sSecret(secret *v1.Secret) error

	// DeleteK8sSecret deletes an existing K8sCluster secret
	DeleteK8sSecret(secret *v1.Secret) error

	// GetK8sSecret gets an existing K8sCluster secret
	GetK8sSecret(secret *v1.Secret) (*v1.Secret, error)

	// FindK8sConfigMap checks if a K8sCluster configmap exists
	FindK8sConfigMap(configMap *v1.ConfigMap) (bool, error)

	// CreateK8sConfigMap creates an new K8sCluster configmap with the data supplied
	CreateK8sConfigMap(configMap *v1.ConfigMap) error

	// UpdateK8sConfigMap updates an new K8sCluster configmap with the data supplied
	UpdateK8sConfigMap(configMap *v1.ConfigMap) error

	// DeleteK8sConfigMap deletes a K8sCluster configmap
	DeleteK8sConfigMap(configMap *v1.ConfigMap) error
}

// K8sUtilsImpl is a structure that hold a kubernetes client and implements the K8sUtils interface
type K8sUtilsImpl struct {
	K8sUtils
	Client   kubernetes.Interface
	ImplName string
}
