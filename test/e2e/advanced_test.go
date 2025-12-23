package e2e

import (
	"os"
	"strings"
	"testing"
)

// TestMergeCommand tests the 'kubecm merge' command
func TestMergeCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	tests := []struct {
		name           string
		configs        []string // config names with their contexts
		outputFile     string
		wantContexts   []string
		wantErr        bool
	}{
		{
			name:         "merge two kubeconfigs",
			configs:      []string{"config1:ctx1", "config2:ctx2"},
			outputFile:   "merged.yaml",
			wantContexts: []string{"ctx1", "ctx2"},
			wantErr:      false,
		},
		{
			name:         "merge three kubeconfigs",
			configs:      []string{"config1:ctx1", "config2:ctx2", "config3:ctx3"},
			outputFile:   "merged-three.yaml",
			wantContexts: []string{"ctx1", "ctx2", "ctx3"},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			
			// Create test kubeconfigs
			var configPaths []string
			for _, cfg := range tt.configs {
				parts := strings.Split(cfg, ":")
				configName := parts[0]
				contextName := parts[1]
				
				configPath := CreateTestKubeconfig(t, configName, 
					"cluster-"+contextName, contextName)
				configPaths = append(configPaths, configPath)
			}

			// Prepare output path
			outputPath := tmpDir + "/" + tt.outputFile

			// Build merge command arguments
			args := []string{"merge", "-y"}
			args = append(args, configPaths...)
			args = append(args, "--config", outputPath)

			// Run kubecm merge
			output, err := RunKubecm(t, args...)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v, output: %s", err, output)
			}

			// Verify the merged config contains all expected contexts
			if !tt.wantErr {
				// Check if output file exists
				if _, err := os.Stat(outputPath); os.IsNotExist(err) {
					t.Errorf("Merged config file not created at: %s", outputPath)
					return
				}

				// List contexts from merged config
				listOutput, err := RunKubecmWithEnv(t, 
					map[string]string{"KUBECONFIG": outputPath}, "list")
				if err != nil {
					t.Errorf("Failed to list contexts from merged config: %v", err)
				}

				// Verify all expected contexts are present
				for _, expectedCtx := range tt.wantContexts {
					if !strings.Contains(listOutput, expectedCtx) {
						t.Errorf("Expected context %s not found in merged config: %s", 
							expectedCtx, listOutput)
					}
				}
			}
		})
	}
}

// TestRenameCommand tests the 'kubecm rename' command
func TestRenameCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	tests := []struct {
		name        string
		contexts    []string
		oldName     string
		newName     string
		wantErr     bool
	}{
		{
			name:     "rename existing context",
			contexts: []string{"old-context", "other-context"},
			oldName:  "old-context",
			newName:  "new-context",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test kubeconfig with multiple contexts
			testConfig := CreateMultiContextKubeconfig(t, "test.yaml", tt.contexts)

			// Run kubecm rename
			output, err := RunKubecmWithEnv(t, map[string]string{"KUBECONFIG": testConfig}, 
				"rename", tt.oldName, tt.newName)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v, output: %s", err, output)
			}

			// Verify the context was renamed
			if !tt.wantErr {
				listOutput, err := RunKubecmWithEnv(t, map[string]string{"KUBECONFIG": testConfig}, "list")
				if err != nil {
					t.Errorf("Failed to list contexts after rename: %v", err)
				}
				
				// Old name should not exist
				if strings.Contains(listOutput, tt.oldName) {
					t.Errorf("Old context name %s still found after rename: %s", tt.oldName, listOutput)
				}
				
				// New name should exist
				if !strings.Contains(listOutput, tt.newName) {
					t.Errorf("New context name %s not found after rename: %s", tt.newName, listOutput)
				}
			}
		})
	}
}

// TestNamespaceCommand tests the 'kubecm namespace' command
// Note: This test requires connection to a real cluster, so it's skipped in basic e2e tests
func TestNamespaceCommand(t *testing.T) {
	t.Skip("Namespace command requires connection to real Kubernetes cluster - skipped in basic e2e tests")

	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	tests := []struct {
		name      string
		contexts  []string
		namespace string
		wantErr   bool
	}{
		{
			name:      "set namespace for current context",
			contexts:  []string{"test-context"},
			namespace: "kube-system",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test kubeconfig
			testConfig := CreateMultiContextKubeconfig(t, "test.yaml", tt.contexts)

			// Run kubecm namespace
			output, err := RunKubecmWithEnv(t, map[string]string{"KUBECONFIG": testConfig}, 
				"namespace", tt.namespace)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v, output: %s", err, output)
			}
		})
	}
}

// TestClearCommand tests the 'kubecm clear' command
// Note: Clear command removes lapsed contexts/clusters which requires validation against real clusters
func TestClearCommand(t *testing.T) {
	t.Skip("Clear command validates against real clusters - skipped in basic e2e tests")

	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	tests := []struct {
		name     string
		contexts []string
		wantErr  bool
	}{
		{
			name:     "clear all contexts",
			contexts: []string{"context-1", "context-2", "context-3"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test kubeconfig with multiple contexts
			testConfig := CreateMultiContextKubeconfig(t, "test.yaml", tt.contexts)

			// Run kubecm clear (passes file path as argument to avoid interactive prompt)
			output, err := RunKubecmWithEnv(t, map[string]string{"KUBECONFIG": testConfig}, 
				"clear", testConfig)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v, output: %s", err, output)
			}
		})
	}
}
