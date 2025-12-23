package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// KubecmPath returns the path to the kubecm binary
func KubecmPath() string {
	path := os.Getenv("KUBECM_BIN")
	if path == "" {
		// Default to the binary in the bin directory
		path = "../../bin/kubecm"
	}
	return path
}

// RunKubecm executes kubecm with the given arguments
func RunKubecm(t *testing.T, args ...string) (string, error) {
	t.Helper()

	cmd := exec.Command(KubecmPath(), args...)
	output, err := cmd.CombinedOutput()

	t.Logf("Running: kubecm %s", strings.Join(args, " "))
	if len(output) > 0 {
		t.Logf("Output: %s", string(output))
	}

	return string(output), err
}

// RunKubecmWithEnv executes kubecm with the given arguments and environment variables
func RunKubecmWithEnv(t *testing.T, env map[string]string, args ...string) (string, error) {
	t.Helper()

	cmd := exec.Command(KubecmPath(), args...)

	// Set environment variables
	cmd.Env = os.Environ()
	for key, value := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	output, err := cmd.CombinedOutput()

	t.Logf("Running: kubecm %s (with env: %v)", strings.Join(args, " "), env)
	if len(output) > 0 {
		t.Logf("Output: %s", string(output))
	}

	return string(output), err
}

// CreateTestKubeconfig creates a test kubeconfig file with the given name and returns its path
func CreateTestKubeconfig(t *testing.T, name, clusterName, contextName string) string {
	t.Helper()

	tmpDir := t.TempDir()
	kubeconfigPath := filepath.Join(tmpDir, name)

	content := fmt.Sprintf(`apiVersion: v1
clusters:
- cluster:
    server: https://%s.example.com:6443
  name: %s
contexts:
- context:
    cluster: %s
    user: test-user
  name: %s
current-context: %s
kind: Config
preferences: {}
users:
- name: test-user
  user:
    token: test-token
`, clusterName, clusterName, clusterName, contextName, contextName)

	err := os.WriteFile(kubeconfigPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test kubeconfig: %v", err)
	}

	t.Logf("Created test kubeconfig at: %s", kubeconfigPath)
	return kubeconfigPath
}

// CreateMultiContextKubeconfig creates a test kubeconfig with multiple contexts
func CreateMultiContextKubeconfig(t *testing.T, name string, contexts []string) string {
	t.Helper()

	tmpDir := t.TempDir()
	kubeconfigPath := filepath.Join(tmpDir, name)

	var clusterSection strings.Builder
	var contextSection strings.Builder
	var userSection strings.Builder

	for i, ctx := range contexts {
		clusterName := fmt.Sprintf("cluster-%d", i)
		userName := fmt.Sprintf("user-%d", i)

		clusterSection.WriteString(fmt.Sprintf(`- cluster:
    server: https://%s.example.com:6443
  name: %s
`, clusterName, clusterName))

		contextSection.WriteString(fmt.Sprintf(`- context:
    cluster: %s
    user: %s
  name: %s
`, clusterName, userName, ctx))

		userSection.WriteString(fmt.Sprintf(`- name: %s
  user:
    token: token-%d
`, userName, i))
	}

	currentContext := contexts[0]
	if len(contexts) > 0 {
		currentContext = contexts[0]
	}

	content := fmt.Sprintf(`apiVersion: v1
clusters:
%s
contexts:
%s
current-context: %s
kind: Config
preferences: {}
users:
%s
`, clusterSection.String(), contextSection.String(), currentContext, userSection.String())

	err := os.WriteFile(kubeconfigPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create multi-context kubeconfig: %v", err)
	}

	t.Logf("Created multi-context kubeconfig at: %s with contexts: %v", kubeconfigPath, contexts)
	return kubeconfigPath
}

// GetTempKubeconfig returns a temporary kubeconfig path for testing
func GetTempKubeconfig(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	return filepath.Join(tmpDir, "kubeconfig")
}
