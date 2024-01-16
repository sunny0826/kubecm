package cmd

import (
	"context"
	"testing"

	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCreateRoleBinding(t *testing.T) {
	// Create a fake clientset
	clientset := fake.NewSimpleClientset()

	// Create a CreateOptions instance
	co := &CreateOptions{
		clientSet: clientset,
		userName:  "test-user",
		role:      "test-role",
		namespace: "test-namespace",
	}

	// Call the function
	err := co.createRoleBinding()
	if err != nil {
		t.Fatalf("createRoleBinding() error = %v", err)
	}

	// Get the role binding
	rb, err := clientset.RbacV1().RoleBindings(co.namespace).Get(context.TODO(), co.userName+"-"+co.role, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Failed to get role binding: %v", err)
	}

	// Check the role binding
	if rb.Name != co.userName+"-"+co.role {
		t.Errorf("Unexpected role binding name: got %v, want %v", rb.Name, co.userName+"-"+co.role)
	}
	if rb.Namespace != co.namespace {
		t.Errorf("Unexpected namespace: got %v, want %v", rb.Namespace, co.namespace)
	}
	if rb.RoleRef.Name != co.role {
		t.Errorf("Unexpected role ref: got %v, want %v", rb.RoleRef.Name, co.role)
	}
	if len(rb.Subjects) != 1 || rb.Subjects[0].Name != co.userName {
		t.Errorf("Unexpected subjects: got %v, want %v", rb.Subjects, []rbacV1.Subject{{Name: co.userName}})
	}
}
