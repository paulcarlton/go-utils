// (c) Copyright 2019 Hewlett Packard Enterprise Development LP

package k8s

import (
	"io/ioutil"
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	"github.hpe.com/platform-core/utils/pkg/k8sutils/v1/testutils"
)

func TestGetClientset(t *testing.T) {
	// Test 1 - verify it works when passed vaild config data
	utils := &K8s{}
	testutils.Setup(t, utils, testutils.TestConfigFile)
	client := utils.GetClientset()
	if client == nil {
		t.Errorf("Unexpected error GetclientSet returned nil")
	}
	if _, ok := client.(kubernetes.Interface); !ok {
		t.Errorf("Unexpected error GetclientSet failed to return kubernetes.Interface")
	}

	// Test 2 - verify it failes when passed invaild config data
	utils = &K8s{}
	testutils.Setup(t, utils, testutils.NoConfig)
	client = utils.GetClientset()
	if client != nil {
		t.Errorf("Unexpected error GetclientSet returned non nil")
	}
}

func TestSetClientset(t *testing.T) {
	utils := &K8s{}
	testutils.Setup(t, utils, testutils.TestConfigFile)
	utils.Client = fake.NewSimpleClientset()

	testConfig, err := ioutil.ReadFile("../testutils/tests/fixtures/kubeconfig.yaml")
	if err != nil {
		t.Errorf("failed to read fixture file, %s", err)
	}
	cases := []struct {
		config        []byte
		errorExpected bool
	}{
		{
			config:        testConfig,
			errorExpected: false,
		},
		{
			config:        []byte(""),
			errorExpected: true,
		},
	}
	for _, c := range cases {
		err := utils.SetClientset(c.config)
		if err == nil && c.errorExpected {
			t.Errorf("Expected error for config %v, got no error", c.config)
		}
		if err != nil && !c.errorExpected {
			t.Errorf("Unexpected error for config %v, error: %s", c.config, err)
		}
	}
}

func TestFindK8sSecret(t *testing.T) {
	utils := &K8s{}
	testutils.Setup(t, utils, testutils.TestConfigFile)
	utils.Client = fake.NewSimpleClientset()

	cases := []struct {
		namespace string
		expected  bool
	}{
		{
			namespace: "Test",
			expected:  true,
		},
		{
			namespace: "Wrong",
			expected:  false,
		},
	}
	secret := *testutils.TestSecret
	_, err := utils.Client.CoreV1().Secrets(secret.Namespace).Create(&secret)
	if err != nil {
		t.Errorf("failed to create test secret")
	}
	for _, c := range cases {
		secret.Namespace = c.namespace
		result, _ := utils.FindK8sSecret(&secret)

		if result != c.expected {
			t.Errorf("Expected result for namespace \"%v\" was: %v, got: %v.", c.namespace, c.expected, result)
		}
	}

}

func TestCreateK8sSecret(t *testing.T) {
	utils := &K8s{}
	testutils.Setup(t, utils, testutils.TestConfigFile)
	utils.Client = fake.NewSimpleClientset()

	cases := []struct {
		secret        *v1.Secret
		errorExpected bool
	}{
		{
			secret:        testutils.TestSecret,
			errorExpected: false,
		},
		{
			secret:        testutils.TestBlankSecret,
			errorExpected: true,
		},
	}

	for _, c := range cases {
		coreErr := utils.CreateK8sSecret(c.secret)
		if coreErr != nil && !c.errorExpected {
			t.Error(coreErr.Error())
		}

		_, err := utils.Client.CoreV1().Secrets(testutils.TestSecret.Namespace).Get(testutils.TestSecret.Name, metav1.GetOptions{})
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestUpdateK8sSecret(t *testing.T) {
	utils := &K8s{}
	testutils.Setup(t, utils, testutils.TestConfigFile)
	utils.Client = fake.NewSimpleClientset()

	cases := []struct {
		namespace     string
		errorExpected bool
	}{
		{
			namespace:     "Test",
			errorExpected: false,
		},
		{
			namespace:     "Wrong",
			errorExpected: true,
		},
	}
	secret := *testutils.TestSecret
	_, err := utils.Client.CoreV1().Secrets(secret.Namespace).Create(&secret)
	if err != nil {
		t.Errorf("failed to create test secret")
	}

	for _, c := range cases {
		secret.Namespace = c.namespace
		err := utils.UpdateK8sSecret(&secret)

		if err != nil && !c.errorExpected {
			t.Errorf("Unexpected error for namespace \"%v\": %v.", c.namespace, err)

		}
	}

}

func TestDeleteK8sSecret(t *testing.T) {
	utils := &K8s{}
	testutils.Setup(t, utils, testutils.TestConfigFile)
	utils.Client = fake.NewSimpleClientset()

	cases := []struct {
		errorExpected bool
	}{
		{
			errorExpected: false,
		},
		{
			errorExpected: true,
		},
	}

	_, err := utils.Client.CoreV1().Secrets(testutils.TestSecret.Namespace).Create(testutils.TestSecret)
	if err != nil {
		t.Errorf("failed to create test secret")
	}

	for _, c := range cases {
		err := utils.DeleteK8sSecret(testutils.TestSecret)

		if err != nil && !c.errorExpected {
			t.Errorf("Unexpected error: %v.", err)
		}
	}

}

func TestGetK8sSecret(t *testing.T) {
	utils := &K8s{}
	testutils.Setup(t, utils, testutils.TestConfigFile)
	utils.Client = fake.NewSimpleClientset()

	cases := []struct {
		namespace     string
		errorExpected bool
	}{
		{
			namespace:     "Test",
			errorExpected: false,
		},
		{
			namespace:     "Wrong",
			errorExpected: true,
		},
	}
	secret := *testutils.TestSecret
	_, err := utils.Client.CoreV1().Secrets(secret.Namespace).Create(&secret)
	if err != nil {
		t.Error("failed to create test secret")
	}
	for _, c := range cases {
		secret.Namespace = c.namespace
		gotSecret, err := utils.GetK8sSecret(&secret)

		if err != nil && !c.errorExpected {
			t.Errorf("Unexpected error: %v", err)
		}

		if err == nil && gotSecret.Namespace != c.namespace {
			t.Errorf("Secret has namespace %v, exp[ected %v", gotSecret.Namespace, c.namespace)
		}
	}
}

func TestFindK8sConfigMap(t *testing.T) {
	utils := &K8s{}
	testutils.Setup(t, utils, testutils.TestConfigFile)
	utils.Client = fake.NewSimpleClientset()
	cases := []struct {
		namespace string
		expected  bool
	}{
		{
			namespace: "Test",
			expected:  true,
		},
		{
			namespace: "Wrong",
			expected:  false,
		},
	}
	configMap := *testutils.TestConfigMap
	_, err := utils.Client.CoreV1().ConfigMaps(configMap.Namespace).Create(&configMap)
	if err != nil {
		t.Errorf("failed to create test secret")
	}

	for _, c := range cases {
		configMap.Namespace = c.namespace
		result, _ := utils.FindK8sConfigMap(&configMap)

		if result != c.expected {
			t.Errorf("Expected result for namespace \"%v\" was: %v, got: %v.", c.namespace, c.expected, result)
		}
	}

}

func TestCreateK8sConfigmap(t *testing.T) {
	utils := &K8s{}
	testutils.Setup(t, utils, testutils.TestConfigFile)
	utils.Client = fake.NewSimpleClientset()

	cases := []struct {
		configMap     *v1.ConfigMap
		errorExpected bool
	}{
		{
			configMap:     testutils.TestConfigMap,
			errorExpected: false,
		},
		{
			configMap:     testutils.TestBlankConfigMap,
			errorExpected: true,
		},
	}

	for _, c := range cases {
		coreErr := utils.CreateK8sConfigMap(c.configMap)
		if coreErr != nil && !c.errorExpected {
			t.Error(coreErr.Error())
		}

		_, err := utils.Client.CoreV1().ConfigMaps(testutils.TestConfigMap.Namespace).Get(testutils.TestConfigMap.Name, metav1.GetOptions{})
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestUpdateK8sConfigMap(t *testing.T) {
	utils := &K8s{}
	testutils.Setup(t, utils, testutils.TestConfigFile)
	utils.Client = fake.NewSimpleClientset()

	cases := []struct {
		namespace     string
		errorExpected bool
	}{
		{
			namespace:     "Test",
			errorExpected: false,
		},
		{
			namespace:     "Wrong",
			errorExpected: true,
		},
	}
	configMap := *testutils.TestConfigMap
	_, err := utils.Client.CoreV1().ConfigMaps(configMap.Namespace).Create(&configMap)
	if err != nil {
		t.Errorf("failed to create test config map")
	}

	for _, c := range cases {
		configMap.Namespace = c.namespace
		err := utils.UpdateK8sConfigMap(&configMap)

		if err != nil && !c.errorExpected {
			t.Errorf("Unexpected error for namespace \"%v\": %v.", c.namespace, err)
		}
	}

}

func TestDeleteK8sConfigMap(t *testing.T) {
	utils := &K8s{}
	testutils.Setup(t, utils, testutils.TestConfigFile)
	utils.Client = fake.NewSimpleClientset()
	cases := []struct {
		errorExpected bool
	}{
		{
			errorExpected: false,
		},
		{
			errorExpected: true,
		},
	}

	_, err := utils.Client.CoreV1().ConfigMaps(testutils.TestConfigMap.Namespace).Create(testutils.TestConfigMap)
	if err != nil {
		t.Errorf("failed to create test secret")
	}
	for _, c := range cases {
		err := utils.DeleteK8sConfigMap(testutils.TestConfigMap)

		if err != nil && !c.errorExpected {
			t.Errorf("Unexpected error: %v.", err)
		}
	}

}

func TestK8sGetSecret(t *testing.T) {
	utils := &K8s{}
	testutils.Setup(t, utils, testutils.TestConfigFile)
	utils.Client = fake.NewSimpleClientset()

	cases := []struct {
		namespace     string
		errorExpected bool
	}{
		{
			namespace:     "Test",
			errorExpected: false,
		},
		{
			namespace:     "Wrong",
			errorExpected: true,
		},
	}
	secret := *testutils.TestSecret
	_, err := utils.Client.CoreV1().Secrets(secret.Namespace).Create(&secret)
	if err != nil {
		t.Error("failed to create test secret")
	}
	for _, c := range cases {
		secret.Namespace = c.namespace
		_, err := k8sGetSecret(utils, &secret)

		if err != nil && !c.errorExpected {
			t.Errorf("Unexpected error: %v", err)
		}

	}
}

func TestDecodeK8s(t *testing.T) {
	cases := []struct {
		testType      string
		inputData     []byte
		handled       bool
		errorExpected bool
	}{
		{
			testType:      "secret",
			inputData:     testutils.TestSecretByte,
			handled:       false,
			errorExpected: false,
		},
		{
			testType:      "secret",
			inputData:     testutils.TestInvalidSecretByte,
			handled:       false,
			errorExpected: true,
		},
		{
			testType:      "config map",
			inputData:     testutils.TestConfigMapByte,
			handled:       false,
			errorExpected: false,
		},
		{
			testType:      "invalid",
			inputData:     testutils.TestInvalidByte,
			handled:       false,
			errorExpected: false,
		},
		{
			testType:      "invalid",
			inputData:     testutils.TestInvalidJSONByte,
			handled:       false,
			errorExpected: true,
		},
	}
	for _, c := range cases {
		data, err := decodeK8s(c.inputData)

		if err != nil && !c.errorExpected {
			t.Error(err.Error())
		}

		switch data.(type) {
		case *v1.Secret:
			c.handled = true

		case *v1.ConfigMap:
			c.handled = true
		default:
			if data == nil {
				c.handled = true
			}
		}

		if !c.handled {
			t.Errorf("Input data \"%v\" not handled correctly.", c.testType)
		}
	}

}

func TestDecodeK8sSecret(t *testing.T) {
	cases := []struct {
		inputData     []byte
		panicExpected bool
	}{
		{
			inputData:     testutils.TestSecretByte,
			panicExpected: false,
		},
		{
			inputData:     testutils.TestInvalidByte,
			panicExpected: true,
		},
	}
	for _, c := range cases {
		defer func() {
			r := recover()
			if r != nil && !c.panicExpected {
				t.Errorf("Recovered from unexpected panic: %v", r)
			}
		}()
		_, err := DecodeK8sSecret(c.inputData)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestDecodeK8sConfigMap(t *testing.T) {
	cases := []struct {
		inputData     []byte
		panicExpected bool
	}{
		{
			inputData:     testutils.TestConfigMapByte,
			panicExpected: false,
		},
		{
			inputData:     testutils.TestInvalidByte,
			panicExpected: true,
		},
	}
	for _, c := range cases {
		defer func() {
			r := recover()
			if r != nil && !c.panicExpected {
				t.Errorf("Recovered from unexpected panic: %v", r)
			}
		}()
		_, err := DecodeK8sConfigMap(c.inputData)
		if err != nil {
			t.Error(err.Error())
		}
	}
}
