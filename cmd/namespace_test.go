package cmd

import (
	"context"
	"errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testSelectNamespacePrompt struct {
	index int
	err   error
}

func (t *testSelectNamespacePrompt) Run() (int, string, error) {
	return t.index, "", t.err
}

func TestSelectNamespace(t *testing.T) {
	namespaces := []Namespaces{
		{Name: "Namespace1", Default: false},
		{Name: "Namespace2", Default: false},
		{Name: "Namespace3", Default: true},
		{Name: "<Exit>", Default: false},
	}

	tests := []struct {
		name         string
		selectPrompt SelectRunner
		expected     int
		expectErr    bool
	}{
		{
			name: "Select First Namespace",
			selectPrompt: &testSelectNamespacePrompt{
				index: 0,
				err:   nil,
			},
			expected:  0,
			expectErr: false,
		},
		{
			name: "Error Occurred",
			selectPrompt: &testSelectNamespacePrompt{
				index: 0,
				err:   errors.New("prompt error"),
			},
			expected:  0,
			expectErr: true,
		},
		{
			name: "Select Exit",
			selectPrompt: &testSelectNamespacePrompt{
				index: 3,
				err:   nil,
			},
			expected:  0,
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selectedNamespace, err := selectNamespaceWithRunner(namespaces, test.selectPrompt)
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, selectedNamespace)
			}
		})
	}
}

func TestCheckNamespaceExist(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		namespace   string
		namespaces  []string
		expectExist bool
		expectError bool
	}{
		{
			name:        "Namespace exists",
			namespace:   "test",
			namespaces:  []string{"default", "test"},
			expectExist: true,
			expectError: false,
		},
		{
			name:        "Namespace does not exist",
			namespace:   "nonexistent",
			namespaces:  []string{"default", "test"},
			expectExist: false,
			expectError: true,
		},
		{
			name:        "Error case",
			namespace:   "",
			namespaces:  nil, // Assuming this simulates an error condition
			expectExist: false,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clientset := mockKubernetesClientSet(tc.namespaces)
			exist, err := CheckNamespaceExist(tc.namespace, clientset)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if exist != tc.expectExist {
					t.Errorf("Expected existence to be %v, got %v", tc.expectExist, exist)
				}
			}
		})
	}
}

func TestGetNamespaceList(t *testing.T) {
	testCases := []struct {
		name             string
		currentNamespace string
		namespaces       []string
		expected         []Namespaces
		expectError      bool
	}{
		{
			name:             "Success with default namespace",
			currentNamespace: "default",
			namespaces:       []string{"default", "test"},
			expected: []Namespaces{
				{Name: "default", Default: true},
				{Name: "test", Default: false},
			},
			expectError: false,
		},
		{
			name:             "Success with specified namespace",
			currentNamespace: "test",
			namespaces:       []string{"default", "test"},
			expected: []Namespaces{
				{Name: "default", Default: false},
				{Name: "test", Default: true},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clientset := mockKubernetesClientSet(tc.namespaces)
			nss, err := GetNamespaceList(tc.currentNamespace, clientset)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, nss)
			}
		})
	}
}

// mockKubernetesClientSet creates a mock clientset that contains the provided namespaces.
func mockKubernetesClientSet(namespaces []string) kubernetes.Interface {
	clientset := fake.NewSimpleClientset()

	for _, ns := range namespaces {
		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: ns,
			},
		}
		clientset.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	}

	return clientset
}

func TestGetClientSet(t *testing.T) {
	// Create a fake clientset
	clientset := fake.NewSimpleClientset()

	// Create a namespace object
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
	}

	// Add the namespace to the fake clientset
	_, err := clientset.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Could not create namespace: %v", err)
	}

	// Use the clientset to get the list of namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Could not list namespaces: %v", err)
	}

	// Check that the "default" namespace exists
	for _, ns := range namespaces.Items {
		if ns.Name == "default" {
			return
		}
	}
	t.Fatal("Did not find 'default' namespace")
}
