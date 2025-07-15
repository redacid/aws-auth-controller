package kube

import "k8s.io/client-go/rest"

var testConfig *rest.Config

// SetTestConfig встановлює конфігурацію для тестів
func SetTestConfig(cfg *rest.Config) {
	testConfig = cfg
}

// ResetTestConfig скидає тестову конфігурацію
func ResetTestConfig() {
	testConfig = nil
}
