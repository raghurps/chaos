package main

import (
	"testing"

	"chaosmonkey.monke/chaos/pkg/k8s"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MockK8sClient struct{}

func (mc *MockK8sClient) ListFilteredNamespaces(excludedNamespaces map[string]struct{}) ([]string, error) {
	return []string{"namespace1", "namespace2", "namespace3"}, nil
}

func (mc *MockK8sClient) ListDeploymentsWithFilters(namespaces []string, annotations map[string]string) ([]v1.Deployment, error) {
	return []v1.Deployment{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "deployment1",
				Namespace: "namespace1",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "deployment2",
				Namespace: "namespace2",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "deployment4",
				Namespace: "namespace3",
			},
		},
	}, nil
}

func (mc *MockK8sClient) ListPodsWithOpts(namespace string, opts metav1.ListOptions) ([]string, error) {
	return []string{"namespace1/pod1", "namespace2/pod2", "namespace3/pod3"}, nil
}

func (mc *MockK8sClient) KillPod(name, namespace string) error {
	return nil
}

func TestKillRandomPod(t *testing.T) {
	k8s.K8sClientInstance = &MockK8sClient{}
	err := killRandomPod([]string{"namespace1", "namespace2", "namespace3"}, []string{"namespace4"}, map[string]string{})
	assert.Nil(t, err)
}
