package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// setupRegistryGitRepo creates a local git repository with a valid registry structure
// containing only static fragments (no cloud API calls needed).
func setupRegistryGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// registry.yaml
	writeRegistryFile(t, filepath.Join(dir, "registry.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: Registry
metadata:
  name: e2e-test
  description: "E2E test registry"
variables:
  - name: Username
    description: "Test username"
    required: true
`)

	// roles/
	os.MkdirAll(filepath.Join(dir, "roles"), 0o755)
	writeRegistryFile(t, filepath.Join(dir, "roles", "devops.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: Role
metadata:
  name: devops
  description: "DevOps - all clusters"
contextPrefix: "e2e"
fragments:
  - cluster-a
  - cluster-b
`)

	writeRegistryFile(t, filepath.Join(dir, "roles", "readonly.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: Role
metadata:
  name: readonly
  description: "Read-only - single cluster"
contextPrefix: "e2e"
fragments:
  - cluster-a
`)

	// fragments/
	os.MkdirAll(filepath.Join(dir, "fragments"), 0o755)
	writeRegistryFile(t, filepath.Join(dir, "fragments", "cluster-a.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: Fragment
metadata:
  name: cluster-a
provider: static
kubeconfig: |
  apiVersion: v1
  kind: Config
  clusters:
    - cluster:
        server: https://cluster-a.example.com:6443
      name: cluster-a
  contexts:
    - context:
        cluster: cluster-a
        user: cluster-a
      name: cluster-a
  users:
    - name: cluster-a
      user:
        token: "{{ .Username }}-token-a"
`)

	writeRegistryFile(t, filepath.Join(dir, "fragments", "cluster-b.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: Fragment
metadata:
  name: cluster-b
provider: static
kubeconfig: |
  apiVersion: v1
  kind: Config
  clusters:
    - cluster:
        server: https://cluster-b.example.com:6443
      name: cluster-b
  contexts:
    - context:
        cluster: cluster-b
        user: cluster-b
      name: cluster-b
  users:
    - name: cluster-b
      user:
        token: default-token-b
`)

	// Initialize git repo
	runGitCommand(t, dir, "init", "-b", "main")
	runGitCommand(t, dir, "add", ".")
	runGitCommand(t, dir, "commit", "-m", "initial commit")

	return dir
}

func writeRegistryFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}

func runGitCommand(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=E2E Test",
		"GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=E2E Test",
		"GIT_COMMITTER_EMAIL=test@example.com",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
}

// registryEnv returns environment variables for isolated registry e2e tests.
// KUBECM_HOME isolates the registry state, KUBECONFIG isolates the kubeconfig.
func registryEnv(kubeconfigPath, kubecmHome string) map[string]string {
	return map[string]string{
		"KUBECONFIG":  kubeconfigPath,
		"KUBECM_HOME": kubecmHome,
	}
}

// setupRegistryTest creates a fresh isolated environment for a registry e2e test.
// Returns the git repo path, kubeconfig path, kubecm home dir, and env map.
func setupRegistryTest(t *testing.T) (repoDir, kubeconfig, kubecmHome string, env map[string]string) {
	t.Helper()

	repoDir = setupRegistryGitRepo(t)
	kubeconfig = GetTempKubeconfig(t)
	kubecmHome = t.TempDir()
	env = registryEnv(kubeconfig, kubecmHome)

	// Create an initial kubeconfig so kubecm can load it
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

// TestRegistryFullLifecycle tests the complete registry workflow:
// add -> list -> sync dry-run -> sync -> update role -> sync -> remove
func TestRegistryFullLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	repoDir, _, _, env := setupRegistryTest(t)

	// Step 1: registry add
	t.Run("add", func(t *testing.T) {
		output, err := RunKubecmWithEnv(t, env,
			"registry", "add",
			"--name", "acme",
			"--url", repoDir,
			"--role", "devops",
			"--var", "Username=testuser",
		)
		if err != nil {
			t.Fatalf("registry add failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "e2e-cluster-a") || !strings.Contains(output, "e2e-cluster-b") {
			t.Errorf("expected both contexts in add output: %s", output)
		}
	})

	// Step 2: verify contexts in kubeconfig
	t.Run("list-contexts", func(t *testing.T) {
		output, err := RunKubecmWithEnv(t, env, "list")
		if err != nil {
			t.Fatalf("list failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "e2e-cluster-a") {
			t.Errorf("context e2e-cluster-a not found: %s", output)
		}
		if !strings.Contains(output, "e2e-cluster-b") {
			t.Errorf("context e2e-cluster-b not found: %s", output)
		}
	})

	// Step 3: registry list
	t.Run("registry-list", func(t *testing.T) {
		output, err := RunKubecmWithEnv(t, env, "registry", "list")
		if err != nil {
			t.Fatalf("registry list failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "acme") {
			t.Errorf("registry 'acme' not in list: %s", output)
		}
		if !strings.Contains(output, "devops") {
			t.Errorf("role 'devops' not in list: %s", output)
		}
	})

	// Step 4: sync dry-run
	t.Run("sync-dry-run", func(t *testing.T) {
		output, err := RunKubecmWithEnv(t, env, "registry", "sync", "acme", "--dry-run")
		if err != nil {
			t.Fatalf("sync --dry-run failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "dry-run") {
			t.Errorf("expected dry-run notice: %s", output)
		}
	})

	// Step 5: sync
	t.Run("sync", func(t *testing.T) {
		output, err := RunKubecmWithEnv(t, env, "registry", "sync", "acme")
		if err != nil {
			t.Fatalf("sync failed: %v\nOutput: %s", err, output)
		}
	})

	// Step 6: update role to readonly (only cluster-a)
	t.Run("update-role", func(t *testing.T) {
		output, err := RunKubecmWithEnv(t, env, "registry", "update", "acme", "--role", "readonly")
		if err != nil {
			t.Fatalf("update failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "updated") {
			t.Errorf("expected 'updated' in output: %s", output)
		}
	})

	// Step 7: sync after role change -> cluster-b should be removed
	t.Run("sync-after-role-change", func(t *testing.T) {
		output, err := RunKubecmWithEnv(t, env, "registry", "sync", "acme")
		if err != nil {
			t.Fatalf("sync failed: %v\nOutput: %s", err, output)
		}
		t.Logf("sync output: %s", output)
	})

	// Step 8: verify cluster-b is gone
	t.Run("verify-role-change", func(t *testing.T) {
		output, err := RunKubecmWithEnv(t, env, "list")
		if err != nil {
			t.Fatalf("list failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "e2e-cluster-a") {
			t.Errorf("e2e-cluster-a should still exist: %s", output)
		}
		if strings.Contains(output, "e2e-cluster-b") {
			t.Errorf("e2e-cluster-b should have been removed: %s", output)
		}
	})

	// Step 9: remove
	t.Run("remove", func(t *testing.T) {
		output, err := RunKubecmWithEnv(t, env, "registry", "remove", "acme")
		if err != nil {
			t.Fatalf("remove failed: %v\nOutput: %s", err, output)
		}
	})

	// Step 10: verify contexts removed
	t.Run("verify-remove", func(t *testing.T) {
		output, err := RunKubecmWithEnv(t, env, "list")
		if err != nil {
			t.Logf("list after remove (may fail if only initial-context remains): %v", err)
		}
		if strings.Contains(output, "e2e-cluster-a") {
			t.Errorf("e2e-cluster-a should have been removed: %s", output)
		}
	})

	// Step 11: registry list should be empty
	t.Run("registry-list-empty", func(t *testing.T) {
		output, err := RunKubecmWithEnv(t, env, "registry", "list")
		if err != nil {
			t.Fatalf("registry list failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "No registries configured") {
			t.Errorf("expected empty registry list: %s", output)
		}
	})
}

// TestRegistryAddDuplicate verifies that adding a registry with a duplicate name fails.
func TestRegistryAddDuplicate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	repoDir, _, _, env := setupRegistryTest(t)

	// First add should succeed
	_, err := RunKubecmWithEnv(t, env,
		"registry", "add", "--name", "acme", "--url", repoDir,
		"--role", "devops", "--var", "Username=testuser")
	if err != nil {
		t.Fatalf("first add should succeed: %v", err)
	}

	// Second add with same name should fail
	output, err := RunKubecmWithEnv(t, env,
		"registry", "add", "--name", "acme", "--url", repoDir,
		"--role", "devops", "--var", "Username=testuser")
	if err == nil {
		t.Error("expected error for duplicate registry name")
	}
	if !strings.Contains(output, "already exists") {
		t.Errorf("expected 'already exists' error, got: %s", output)
	}
}

// TestRegistryAddMissingRole verifies that adding a registry with a nonexistent role fails.
func TestRegistryAddMissingRole(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	repoDir, _, _, env := setupRegistryTest(t)

	output, err := RunKubecmWithEnv(t, env,
		"registry", "add", "--name", "acme", "--url", repoDir,
		"--role", "nonexistent", "--var", "Username=testuser")
	if err == nil {
		t.Error("expected error for missing role")
	}
	if !strings.Contains(output, "nonexistent") {
		t.Errorf("expected error mentioning nonexistent role, got: %s", output)
	}
}

// TestRegistrySyncNoRegistries verifies that syncing with no registries configured fails.
func TestRegistrySyncNoRegistries(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	kubeconfig := GetTempKubeconfig(t)
	kubecmHome := t.TempDir()
	env := registryEnv(kubeconfig, kubecmHome)

	output, err := RunKubecmWithEnv(t, env, "registry", "sync", "--all")
	if err == nil {
		t.Error("expected error when no registries configured")
	}
	if !strings.Contains(output, "no registries configured") {
		t.Errorf("expected 'no registries configured' error, got: %s", output)
	}
}

// TestRegistryRemoveNotFound verifies that removing a nonexistent registry fails.
func TestRegistryRemoveNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	_, _, _, env := setupRegistryTest(t)

	output, err := RunKubecmWithEnv(t, env, "registry", "remove", "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent registry")
	}
	if !strings.Contains(output, "not found") {
		t.Errorf("expected 'not found' error, got: %s", output)
	}
}

// TestRegistryRemoveKeepContexts verifies that --keep-contexts preserves kubeconfig entries.
func TestRegistryRemoveKeepContexts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	repoDir, _, _, env := setupRegistryTest(t)

	// Add registry
	_, err := RunKubecmWithEnv(t, env,
		"registry", "add", "--name", "acme", "--url", repoDir,
		"--role", "devops", "--var", "Username=testuser")
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}

	// Remove with --keep-contexts
	_, err = RunKubecmWithEnv(t, env, "registry", "remove", "acme", "--keep-contexts")
	if err != nil {
		t.Fatalf("remove --keep-contexts failed: %v", err)
	}

	// Contexts should still exist in kubeconfig
	output, err := RunKubecmWithEnv(t, env, "list")
	if err != nil {
		t.Fatalf("list failed: %v\nOutput: %s", err, output)
	}
	if !strings.Contains(output, "e2e-cluster-a") {
		t.Errorf("e2e-cluster-a should be preserved with --keep-contexts: %s", output)
	}
	if !strings.Contains(output, "e2e-cluster-b") {
		t.Errorf("e2e-cluster-b should be preserved with --keep-contexts: %s", output)
	}

	// Registry list should be empty
	output, err = RunKubecmWithEnv(t, env, "registry", "list")
	if err != nil {
		t.Fatalf("registry list failed: %v", err)
	}
	if !strings.Contains(output, "No registries configured") {
		t.Errorf("registry should be removed from config: %s", output)
	}
}

// TestRegistrySyncAll verifies that --all flag syncs all configured registries.
func TestRegistrySyncAll(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	repoDir, _, _, env := setupRegistryTest(t)

	// Add registry
	_, err := RunKubecmWithEnv(t, env,
		"registry", "add", "--name", "acme", "--url", repoDir,
		"--role", "devops", "--var", "Username=testuser")
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}

	// Sync all
	output, err := RunKubecmWithEnv(t, env, "registry", "sync", "--all")
	if err != nil {
		t.Fatalf("sync --all failed: %v\nOutput: %s", err, output)
	}
	if !strings.Contains(output, "acme") {
		t.Errorf("expected 'acme' in sync output: %s", output)
	}
}

// TestRegistryTemplateVariables verifies that template variables are correctly applied.
func TestRegistryTemplateVariables(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	repoDir, kubeconfig, _, env := setupRegistryTest(t)

	// Add with a specific username
	_, err := RunKubecmWithEnv(t, env,
		"registry", "add", "--name", "acme", "--url", repoDir,
		"--role", "devops", "--var", "Username=e2euser")
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}

	// Read the kubeconfig and verify the token was templated
	data, err := os.ReadFile(kubeconfig)
	if err != nil {
		t.Fatalf("reading kubeconfig: %v", err)
	}
	content := string(data)

	// cluster-a fragment uses {{ .Username }}-token-a
	if !strings.Contains(content, "e2euser-token-a") {
		t.Errorf("expected templated token 'e2euser-token-a' in kubeconfig:\n%s", content)
	}
}

// TestRegistryNoPrefix verifies that registries work without a contextPrefix.
func TestRegistryNoPrefix(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	// Create a repo with a role that has no contextPrefix
	dir := t.TempDir()

	writeRegistryFile(t, filepath.Join(dir, "registry.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: Registry
metadata:
  name: noprefix
`)

	os.MkdirAll(filepath.Join(dir, "roles"), 0o755)
	writeRegistryFile(t, filepath.Join(dir, "roles", "plain.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: Role
metadata:
  name: plain
fragments:
  - server-alpha
`)

	os.MkdirAll(filepath.Join(dir, "fragments"), 0o755)
	writeRegistryFile(t, filepath.Join(dir, "fragments", "server-alpha.yaml"), `apiVersion: kubecm.io/v1alpha1
kind: Fragment
metadata:
  name: server-alpha
provider: static
kubeconfig: |
  apiVersion: v1
  kind: Config
  clusters:
    - cluster:
        server: https://alpha.example.com:6443
      name: server-alpha
  contexts:
    - context:
        cluster: server-alpha
        user: server-alpha
      name: server-alpha
  users:
    - name: server-alpha
      user:
        token: alpha-token
`)

	runGitCommand(t, dir, "init", "-b", "main")
	runGitCommand(t, dir, "add", ".")
	runGitCommand(t, dir, "commit", "-m", "initial commit")

	kubeconfig := GetTempKubeconfig(t)
	kubecmHome := t.TempDir()
	env := registryEnv(kubeconfig, kubecmHome)

	initialConfig := CreateTestKubeconfig(t, "init.yaml", "initial-cluster", "initial-context")
	data, _ := os.ReadFile(initialConfig)
	os.WriteFile(kubeconfig, data, 0o644)

	// Add registry with no-prefix role
	_, err := RunKubecmWithEnv(t, env,
		"registry", "add", "--name", "noprefix", "--url", dir,
		"--role", "plain")
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}

	// Context should use fragment name directly (no prefix)
	output, err := RunKubecmWithEnv(t, env, "list")
	if err != nil {
		t.Fatalf("list failed: %v\nOutput: %s", err, output)
	}
	if !strings.Contains(output, "server-alpha") {
		t.Errorf("expected context 'server-alpha' (no prefix): %s", output)
	}
}

// TestRegistryUpdateVariable verifies that updating a variable works.
func TestRegistryUpdateVariable(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	repoDir, _, _, env := setupRegistryTest(t)

	// Add registry
	_, err := RunKubecmWithEnv(t, env,
		"registry", "add", "--name", "acme", "--url", repoDir,
		"--role", "devops", "--var", "Username=olduser")
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}

	// Update variable
	output, err := RunKubecmWithEnv(t, env, "registry", "update", "acme", "--var", "Username=newuser")
	if err != nil {
		t.Fatalf("update failed: %v\nOutput: %s", err, output)
	}
	if !strings.Contains(output, "updated") {
		t.Errorf("expected 'updated' in output: %s", output)
	}
}
