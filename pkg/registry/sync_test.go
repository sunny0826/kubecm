package registry

import (
	"os"
	"path/filepath"
	"testing"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// setupTestRegistry creates a temporary registry repo structure for testing.
func setupTestRegistry(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// registry.yaml
	writeFile(t, filepath.Join(dir, "registry.yaml"), `
apiVersion: kubecm.io/v1alpha1
kind: Registry
metadata:
  name: test
  description: "Test registry"
variables:
  - name: Username
    description: "Test username"
    required: true
`)

	// roles/
	os.MkdirAll(filepath.Join(dir, "roles"), 0o755)
	writeFile(t, filepath.Join(dir, "roles", "devops.yaml"), `
apiVersion: kubecm.io/v1alpha1
kind: Role
metadata:
  name: devops
  description: "DevOps role"
contextPrefix: "test"
fragments:
  - onprem-dc1
  - onprem-dc2
`)

	// fragments/
	os.MkdirAll(filepath.Join(dir, "fragments"), 0o755)
	writeFile(t, filepath.Join(dir, "fragments", "onprem-dc1.yaml"), `
apiVersion: kubecm.io/v1alpha1
kind: Fragment
metadata:
  name: onprem-dc1
provider: static
kubeconfig: |
  apiVersion: v1
  kind: Config
  clusters:
    - cluster:
        server: https://k8s-dc1.internal:6443
      name: dc1
  contexts:
    - context:
        cluster: dc1
        user: dc1
      name: dc1
  users:
    - name: dc1
      user:
        token: "{{ .Username }}-token"
`)

	writeFile(t, filepath.Join(dir, "fragments", "onprem-dc2.yaml"), `
apiVersion: kubecm.io/v1alpha1
kind: Fragment
metadata:
  name: onprem-dc2
provider: static
kubeconfig: |
  apiVersion: v1
  kind: Config
  clusters:
    - cluster:
        server: https://k8s-dc2.internal:6443
      name: dc2
  contexts:
    - context:
        cluster: dc2
        user: dc2
      name: dc2
  users:
    - name: dc2
      user:
        token: default-token
`)

	return dir
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}

func TestSync_AddContexts(t *testing.T) {
	repoDir := setupTestRegistry(t)
	entry := &RegistryEntry{
		Name: "test",
		Role: "devops",
		Variables: map[string]string{
			"Username": "clark",
		},
	}
	currentConfig := clientcmdapi.NewConfig()

	result, err := Sync(repoDir, entry, currentConfig, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Added) != 2 {
		t.Errorf("expected 2 added contexts, got %d: %v", len(result.Added), result.Added)
	}

	// Check contexts exist with prefix
	if _, ok := currentConfig.Contexts["test-onprem-dc1"]; !ok {
		t.Error("expected context 'test-onprem-dc1'")
	}
	if _, ok := currentConfig.Contexts["test-onprem-dc2"]; !ok {
		t.Error("expected context 'test-onprem-dc2'")
	}

	// Check template was applied
	if user, ok := currentConfig.AuthInfos["test-onprem-dc1"]; ok {
		if user.Token != "clark-token" {
			t.Errorf("token = %q, want %q", user.Token, "clark-token")
		}
	} else {
		t.Error("expected authinfo 'test-onprem-dc1'")
	}

	// Check managed contexts updated
	if len(entry.ManagedContexts) != 2 {
		t.Errorf("expected 2 managed contexts, got %d", len(entry.ManagedContexts))
	}
	if entry.LastSync == nil {
		t.Error("expected lastSync to be set")
	}
}

func TestSync_RemoveStaleContexts(t *testing.T) {
	repoDir := setupTestRegistry(t)

	// Setup: entry has a stale context that no longer appears in the role
	entry := &RegistryEntry{
		Name: "test",
		Role: "devops",
		Variables: map[string]string{
			"Username": "clark",
		},
		ManagedContexts: []string{"test-onprem-dc1", "test-onprem-dc2", "test-removed-cluster"},
	}

	currentConfig := clientcmdapi.NewConfig()
	// Add the stale context to the config
	currentConfig.Contexts["test-removed-cluster"] = &clientcmdapi.Context{
		Cluster:  "test-removed-cluster",
		AuthInfo: "test-removed-cluster",
	}
	currentConfig.Clusters["test-removed-cluster"] = &clientcmdapi.Cluster{
		Server: "https://removed:6443",
	}
	currentConfig.AuthInfos["test-removed-cluster"] = &clientcmdapi.AuthInfo{
		Token: "old",
	}

	result, err := Sync(repoDir, entry, currentConfig, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Removed) != 1 || result.Removed[0] != "test-removed-cluster" {
		t.Errorf("expected 1 removed context 'test-removed-cluster', got %v", result.Removed)
	}

	if _, ok := currentConfig.Contexts["test-removed-cluster"]; ok {
		t.Error("stale context should have been removed")
	}
}

func TestSync_SkipConflict(t *testing.T) {
	repoDir := setupTestRegistry(t)

	// Create a role with only one fragment to avoid complexity
	writeFile(t, filepath.Join(repoDir, "roles", "single.yaml"), `
apiVersion: kubecm.io/v1alpha1
kind: Role
metadata:
  name: single
contextPrefix: "test"
fragments:
  - onprem-dc1
`)

	entry := &RegistryEntry{
		Name: "test",
		Role: "single",
		Variables: map[string]string{
			"Username": "clark",
		},
		ManagedContexts: []string{}, // empty - nothing managed yet
	}

	currentConfig := clientcmdapi.NewConfig()
	// Pre-existing context with same name but not managed
	currentConfig.Contexts["test-onprem-dc1"] = &clientcmdapi.Context{
		Cluster:  "existing-cluster",
		AuthInfo: "existing-user",
	}

	result, err := Sync(repoDir, entry, currentConfig, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d: %v", len(result.Skipped), result.Skipped)
	}

	// Original context should be preserved
	if currentConfig.Contexts["test-onprem-dc1"].Cluster != "existing-cluster" {
		t.Error("existing context should not be overwritten")
	}
}

func TestSync_DryRun(t *testing.T) {
	repoDir := setupTestRegistry(t)
	entry := &RegistryEntry{
		Name: "test",
		Role: "devops",
		Variables: map[string]string{
			"Username": "clark",
		},
	}
	currentConfig := clientcmdapi.NewConfig()

	result, err := Sync(repoDir, entry, currentConfig, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Added) != 2 {
		t.Errorf("dry-run should report 2 added, got %d", len(result.Added))
	}

	// Config should NOT be modified in dry-run
	if len(currentConfig.Contexts) != 0 {
		t.Error("dry-run should not modify config")
	}
	if entry.LastSync != nil {
		t.Error("dry-run should not set lastSync")
	}
}

func TestSync_UpdateManagedContexts(t *testing.T) {
	repoDir := setupTestRegistry(t)

	entry := &RegistryEntry{
		Name: "test",
		Role: "devops",
		Variables: map[string]string{
			"Username": "clark",
		},
		ManagedContexts: []string{"test-onprem-dc1", "test-onprem-dc2"},
	}

	currentConfig := clientcmdapi.NewConfig()
	// Pre-existing managed contexts
	currentConfig.Contexts["test-onprem-dc1"] = &clientcmdapi.Context{
		Cluster:  "test-onprem-dc1",
		AuthInfo: "test-onprem-dc1",
	}
	currentConfig.Clusters["test-onprem-dc1"] = &clientcmdapi.Cluster{Server: "https://old:6443"}
	currentConfig.AuthInfos["test-onprem-dc1"] = &clientcmdapi.AuthInfo{Token: "old"}

	currentConfig.Contexts["test-onprem-dc2"] = &clientcmdapi.Context{
		Cluster:  "test-onprem-dc2",
		AuthInfo: "test-onprem-dc2",
	}
	currentConfig.Clusters["test-onprem-dc2"] = &clientcmdapi.Cluster{Server: "https://old2:6443"}
	currentConfig.AuthInfos["test-onprem-dc2"] = &clientcmdapi.AuthInfo{Token: "old2"}

	result, err := Sync(repoDir, entry, currentConfig, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Updated) != 2 {
		t.Errorf("expected 2 updated, got %d: %v", len(result.Updated), result.Updated)
	}

	// Server should be updated
	if currentConfig.Clusters["test-onprem-dc1"].Server != "https://k8s-dc1.internal:6443" {
		t.Errorf("cluster server not updated: %s", currentConfig.Clusters["test-onprem-dc1"].Server)
	}
}

func TestSync_NoPrefix(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "registry.yaml"), `
apiVersion: kubecm.io/v1alpha1
kind: Registry
metadata:
  name: acme
`)

	os.MkdirAll(filepath.Join(dir, "roles"), 0o755)
	writeFile(t, filepath.Join(dir, "roles", "devops.yaml"), `
apiVersion: kubecm.io/v1alpha1
kind: Role
metadata:
  name: devops
fragments:
  - eks-prod
  - eks-staging
`)

	os.MkdirAll(filepath.Join(dir, "fragments"), 0o755)
	writeFile(t, filepath.Join(dir, "fragments", "eks-prod.yaml"), `
apiVersion: kubecm.io/v1alpha1
kind: Fragment
metadata:
  name: eks-prod
provider: static
kubeconfig: |
  apiVersion: v1
  kind: Config
  clusters:
    - cluster:
        server: https://eks-prod.example.com:6443
      name: eks-prod
  contexts:
    - context:
        cluster: eks-prod
        user: eks-prod
      name: eks-prod
  users:
    - name: eks-prod
      user:
        token: tok
`)
	writeFile(t, filepath.Join(dir, "fragments", "eks-staging.yaml"), `
apiVersion: kubecm.io/v1alpha1
kind: Fragment
metadata:
  name: eks-staging
provider: static
kubeconfig: |
  apiVersion: v1
  kind: Config
  clusters:
    - cluster:
        server: https://eks-staging.example.com:6443
      name: eks-staging
  contexts:
    - context:
        cluster: eks-staging
        user: eks-staging
      name: eks-staging
  users:
    - name: eks-staging
      user:
        token: tok
`)

	entry := &RegistryEntry{
		Name: "acme",
		Role: "devops",
	}
	currentConfig := clientcmdapi.NewConfig()

	result, err := Sync(dir, entry, currentConfig, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Without prefix, context names = fragment names
	if _, ok := currentConfig.Contexts["eks-prod"]; !ok {
		t.Errorf("expected context 'eks-prod', got: %v", contextNames(currentConfig))
	}
	if _, ok := currentConfig.Contexts["eks-staging"]; !ok {
		t.Errorf("expected context 'eks-staging', got: %v", contextNames(currentConfig))
	}
	if len(result.Added) != 2 {
		t.Errorf("expected 2 added, got %d: %v", len(result.Added), result.Added)
	}
}

func contextNames(cfg *clientcmdapi.Config) []string {
	var names []string
	for k := range cfg.Contexts {
		names = append(names, k)
	}
	return names
}

func TestFormatSyncResult(t *testing.T) {
	r := &SyncResult{
		Added:   []string{"ctx1"},
		Updated: []string{"ctx2"},
		Removed: []string{"ctx3"},
		Skipped: []string{"ctx4"},
		Errors:  []string{"something failed"},
	}
	out := FormatSyncResult(r)
	if out == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormatSyncResult_NoChanges(t *testing.T) {
	r := &SyncResult{}
	out := FormatSyncResult(r)
	if out != "  No changes.\n" {
		t.Errorf("expected 'No changes.', got %q", out)
	}
}
