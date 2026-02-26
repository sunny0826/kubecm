package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupRegistryWithUsers creates a local git repository with a registry structure
// that includes users/, clusters/ dir, and roles using both fragment: and cluster: fields.
func setupRegistryWithUsers(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// registry.yaml
	writeRegistryFile(t, filepath.Join(dir, "registry.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: Registry
metadata:
  name: e2e-users
  description: "E2E test registry with users"
`)

	// fragments/ - legacy directory with static clusters
	os.MkdirAll(filepath.Join(dir, "fragments"), 0o755)

	writeRegistryFile(t, filepath.Join(dir, "fragments", "legacy-cluster.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: Fragment
metadata:
  name: legacy-cluster
provider: static
kubeconfig: |
  apiVersion: v1
  kind: Config
  clusters:
    - cluster:
        server: https://legacy.example.com:6443
      name: legacy
  contexts:
    - context:
        cluster: legacy
        user: legacy
      name: legacy
  users:
    - name: legacy
      user:
        token: legacy-token
`)

	writeRegistryFile(t, filepath.Join(dir, "fragments", "shared-cluster.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: Fragment
metadata:
  name: shared-cluster
provider: static
kubeconfig: |
  apiVersion: v1
  kind: Config
  clusters:
    - cluster:
        server: https://shared.example.com:6443
      name: shared
  contexts:
    - context:
        cluster: shared
        user: shared
      name: shared
  users:
    - name: shared
      user:
        token: default-token
`)

	// clusters/ - new directory
	os.MkdirAll(filepath.Join(dir, "clusters"), 0o755)

	writeRegistryFile(t, filepath.Join(dir, "clusters", "new-cluster.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: Cluster
metadata:
  name: new-cluster
provider: static
kubeconfig: |
  apiVersion: v1
  kind: Config
  clusters:
    - cluster:
        server: https://new.example.com:6443
      name: new
  contexts:
    - context:
        cluster: new
        user: new
      name: new
  users:
    - name: new
      user:
        token: new-token
`)

	// users/
	os.MkdirAll(filepath.Join(dir, "users"), 0o755)

	writeRegistryFile(t, filepath.Join(dir, "users", "admin.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: User
metadata:
  name: admin
provider: static
`)

	writeRegistryFile(t, filepath.Join(dir, "users", "readonly.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: User
metadata:
  name: readonly
provider: static
`)

	// roles/
	os.MkdirAll(filepath.Join(dir, "roles"), 0o755)

	// Legacy role using fragments: list (backward compat)
	writeRegistryFile(t, filepath.Join(dir, "roles", "legacy.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: Role
metadata:
  name: legacy
contextPrefix: "e2e"
fragments:
  - legacy-cluster
  - shared-cluster
`)

	// Role using contexts: with deprecated fragment: field
	writeRegistryFile(t, filepath.Join(dir, "roles", "contexts-fragment.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: Role
metadata:
  name: contexts-fragment
contextPrefix: "e2e"
contexts:
  - fragment: legacy-cluster
  - fragment: shared-cluster
`)

	// Role using contexts: with new cluster: field
	writeRegistryFile(t, filepath.Join(dir, "roles", "contexts-cluster.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: Role
metadata:
  name: contexts-cluster
contextPrefix: "e2e"
contexts:
  - cluster: legacy-cluster
  - cluster: new-cluster
`)

	// Role with user overrides
	writeRegistryFile(t, filepath.Join(dir, "roles", "with-users.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: Role
metadata:
  name: with-users
contextPrefix: "e2e"
contexts:
  - cluster: shared-cluster
    user: admin
    name: shared-admin
  - cluster: shared-cluster
    user: readonly
    name: shared-ro
  - cluster: legacy-cluster
`)

	// Mixed role: cluster: without user + cluster: with user + fragments/ fallback
	writeRegistryFile(t, filepath.Join(dir, "roles", "mixed.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: Role
metadata:
  name: mixed
contexts:
  - cluster: new-cluster
  - cluster: legacy-cluster
    user: admin
    name: legacy-admin
  - cluster: shared-cluster
`)

	// Initialize git repo
	runGitCommand(t, dir, "init", "-b", "main")
	runGitCommand(t, dir, "add", ".")
	runGitCommand(t, dir, "commit", "-m", "initial commit")

	return dir
}

// setupRegistryUsersTest creates a fresh isolated environment for a registry users e2e test.
func setupRegistryUsersTest(t *testing.T) (repoDir, kubeconfig, kubecmHome string, env map[string]string) {
	t.Helper()

	repoDir = setupRegistryWithUsers(t)
	kubeconfig = GetTempKubeconfig(t)
	kubecmHome = t.TempDir()
	env = registryEnv(kubeconfig, kubecmHome)

	initialConfig := CreateTestKubeconfig(t, "init.yaml", "initial-cluster", "initial-context")
	data, err := os.ReadFile(initialConfig)
	if err != nil {
		t.Fatalf("reading initial config: %v", err)
	}
	if err := os.WriteFile(kubeconfig, data, 0o644); err != nil {
		t.Fatalf("writing kubeconfig: %v", err)
	}

	return
}

// TestRegistryLegacyFragmentsFormat verifies backward compat with the legacy fragments: list.
func TestRegistryLegacyFragmentsFormat(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	repoDir, kubeconfig, _, env := setupRegistryUsersTest(t)

	_, err := RunKubecmWithEnv(t, env,
		"registry", "add", "--name", "acme", "--url", repoDir, "--role", "legacy")
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}

	output, err := RunKubecmWithEnv(t, env, "list")
	if err != nil {
		t.Fatalf("list failed: %v\nOutput: %s", err, output)
	}
	if !strings.Contains(output, "e2e-legacy-cluster") {
		t.Errorf("expected context 'e2e-legacy-cluster': %s", output)
	}
	if !strings.Contains(output, "e2e-shared-cluster") {
		t.Errorf("expected context 'e2e-shared-cluster': %s", output)
	}

	// Verify token from legacy fragment
	data, err := os.ReadFile(kubeconfig)
	if err != nil {
		t.Fatalf("reading kubeconfig: %v", err)
	}
	if !strings.Contains(string(data), "legacy-token") {
		t.Errorf("expected 'legacy-token' in kubeconfig")
	}
}

// TestRegistryContextsWithFragmentField verifies contexts: with deprecated fragment: field.
func TestRegistryContextsWithFragmentField(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	repoDir, _, _, env := setupRegistryUsersTest(t)

	_, err := RunKubecmWithEnv(t, env,
		"registry", "add", "--name", "acme", "--url", repoDir, "--role", "contexts-fragment")
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}

	output, err := RunKubecmWithEnv(t, env, "list")
	if err != nil {
		t.Fatalf("list failed: %v\nOutput: %s", err, output)
	}
	if !strings.Contains(output, "e2e-legacy-cluster") {
		t.Errorf("expected context 'e2e-legacy-cluster': %s", output)
	}
	if !strings.Contains(output, "e2e-shared-cluster") {
		t.Errorf("expected context 'e2e-shared-cluster': %s", output)
	}
}

// TestRegistryContextsWithClusterField verifies contexts: with new cluster: field,
// including loading from the new clusters/ directory.
func TestRegistryContextsWithClusterField(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	repoDir, kubeconfig, _, env := setupRegistryUsersTest(t)

	_, err := RunKubecmWithEnv(t, env,
		"registry", "add", "--name", "acme", "--url", repoDir, "--role", "contexts-cluster")
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}

	output, err := RunKubecmWithEnv(t, env, "list")
	if err != nil {
		t.Fatalf("list failed: %v\nOutput: %s", err, output)
	}
	// legacy-cluster loaded from fragments/ (fallback)
	if !strings.Contains(output, "e2e-legacy-cluster") {
		t.Errorf("expected context 'e2e-legacy-cluster' (from fragments/ fallback): %s", output)
	}
	// new-cluster loaded from clusters/ (new dir)
	if !strings.Contains(output, "e2e-new-cluster") {
		t.Errorf("expected context 'e2e-new-cluster' (from clusters/ dir): %s", output)
	}

	data, err := os.ReadFile(kubeconfig)
	if err != nil {
		t.Fatalf("reading kubeconfig: %v", err)
	}
	if !strings.Contains(string(data), "new-token") {
		t.Errorf("expected 'new-token' from clusters/ directory in kubeconfig")
	}
}

// TestRegistryWithUserOverrides verifies that user overrides work,
// including same cluster with different users producing different context names.
func TestRegistryWithUserOverrides(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	repoDir, _, _, env := setupRegistryUsersTest(t)

	_, err := RunKubecmWithEnv(t, env,
		"registry", "add", "--name", "acme", "--url", repoDir, "--role", "with-users")
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}

	output, err := RunKubecmWithEnv(t, env, "list")
	if err != nil {
		t.Fatalf("list failed: %v\nOutput: %s", err, output)
	}

	// Same cluster with two different users -> two distinct context names
	if !strings.Contains(output, "e2e-shared-admin") {
		t.Errorf("expected context 'e2e-shared-admin': %s", output)
	}
	if !strings.Contains(output, "e2e-shared-ro") {
		t.Errorf("expected context 'e2e-shared-ro': %s", output)
	}
	// Cluster without user override
	if !strings.Contains(output, "e2e-legacy-cluster") {
		t.Errorf("expected context 'e2e-legacy-cluster': %s", output)
	}
}

// TestRegistryMixedMode verifies a role mixing clusters/ dir, fragments/ fallback,
// and user overrides in a single sync.
func TestRegistryMixedMode(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	repoDir, _, _, env := setupRegistryUsersTest(t)

	_, err := RunKubecmWithEnv(t, env,
		"registry", "add", "--name", "acme", "--url", repoDir, "--role", "mixed")
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}

	output, err := RunKubecmWithEnv(t, env, "list")
	if err != nil {
		t.Fatalf("list failed: %v\nOutput: %s", err, output)
	}

	// new-cluster from clusters/ dir, no user
	if !strings.Contains(output, "new-cluster") {
		t.Errorf("expected context 'new-cluster': %s", output)
	}
	// legacy-cluster from fragments/ fallback, with user, custom name
	if !strings.Contains(output, "legacy-admin") {
		t.Errorf("expected context 'legacy-admin': %s", output)
	}
	// shared-cluster from fragments/, no user
	if !strings.Contains(output, "shared-cluster") {
		t.Errorf("expected context 'shared-cluster': %s", output)
	}
}

// TestRegistryUsersSwitchRole verifies switching from a legacy role to a user-based role
// correctly removes old contexts and adds new ones.
func TestRegistryUsersSwitchRole(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	repoDir, _, _, env := setupRegistryUsersTest(t)

	// Start with legacy role
	_, err := RunKubecmWithEnv(t, env,
		"registry", "add", "--name", "acme", "--url", repoDir, "--role", "legacy")
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}

	output, err := RunKubecmWithEnv(t, env, "list")
	if err != nil {
		t.Fatalf("list failed: %v\nOutput: %s", err, output)
	}
	if !strings.Contains(output, "e2e-legacy-cluster") {
		t.Errorf("expected e2e-legacy-cluster with legacy role: %s", output)
	}
	if !strings.Contains(output, "e2e-shared-cluster") {
		t.Errorf("expected e2e-shared-cluster with legacy role: %s", output)
	}

	// Switch to with-users role
	_, err = RunKubecmWithEnv(t, env, "registry", "update", "acme", "--role", "with-users")
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}

	_, err = RunKubecmWithEnv(t, env, "registry", "sync", "acme")
	if err != nil {
		t.Fatalf("sync failed: %v", err)
	}

	output, err = RunKubecmWithEnv(t, env, "list")
	if err != nil {
		t.Fatalf("list failed: %v\nOutput: %s", err, output)
	}

	// Old shared-cluster context should be removed (replaced by named ones)
	if strings.Contains(output, "e2e-shared-cluster") {
		t.Errorf("e2e-shared-cluster should have been removed after role switch: %s", output)
	}
	// New named contexts should exist
	if !strings.Contains(output, "e2e-shared-admin") {
		t.Errorf("expected e2e-shared-admin after role switch: %s", output)
	}
	if !strings.Contains(output, "e2e-shared-ro") {
		t.Errorf("expected e2e-shared-ro after role switch: %s", output)
	}
	// legacy-cluster is in both roles, should still exist
	if !strings.Contains(output, "e2e-legacy-cluster") {
		t.Errorf("expected e2e-legacy-cluster preserved after role switch: %s", output)
	}
}

// TestRegistrySyncDryRunWithUsers verifies dry-run with user overrides.
func TestRegistrySyncDryRunWithUsers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	repoDir, _, _, env := setupRegistryUsersTest(t)

	_, err := RunKubecmWithEnv(t, env,
		"registry", "add", "--name", "acme", "--url", repoDir, "--role", "with-users")
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}

	// Dry-run should show what would change
	output, err := RunKubecmWithEnv(t, env, "registry", "sync", "acme", "--dry-run")
	if err != nil {
		t.Fatalf("sync --dry-run failed: %v\nOutput: %s", err, output)
	}
	if !strings.Contains(output, "dry-run") {
		t.Errorf("expected dry-run notice: %s", output)
	}
	if !strings.Contains(output, "e2e-shared-admin") {
		t.Errorf("expected e2e-shared-admin in dry-run output: %s", output)
	}
}
