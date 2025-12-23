package e2e

import (
	"os"
	"strings"
	"testing"
)

// TestAddCommand tests the 'kubecm add' command
func TestAddCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	tests := []struct {
		name            string
		setupKubeconfig bool
		contextName     string
		wantErr         bool
	}{
		{
			name:            "add context from valid kubeconfig",
			setupKubeconfig: true,
			contextName:     "test-context",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary kubeconfig for the main config
			mainConfig := GetTempKubeconfig(t)

			// Create initial config with one context
			initialConfig := CreateTestKubeconfig(t, "initial.yaml", "initial-cluster", "initial-context")

			// Copy initial config to main config
			data, err := os.ReadFile(initialConfig)
			if err != nil {
				t.Fatalf("Failed to read initial config: %v", err)
			}
			err = os.WriteFile(mainConfig, data, 0644)
			if err != nil {
				t.Fatalf("Failed to write main config: %v", err)
			}

			// Create a test kubeconfig to add
			testConfig := CreateTestKubeconfig(t, "test.yaml", "test-cluster", tt.contextName)

			// Run kubecm add
			output, err := RunKubecmWithEnv(t, map[string]string{"KUBECONFIG": mainConfig},
				"add", "-cf", testConfig)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v, output: %s", err, output)
			}

			// Verify the context was added by listing contexts
			if !tt.wantErr {
				listOutput, err := RunKubecmWithEnv(t, map[string]string{"KUBECONFIG": mainConfig}, "list")
				if err != nil {
					t.Errorf("Failed to list contexts after add: %v", err)
				}
				if !strings.Contains(listOutput, tt.contextName) {
					t.Errorf("Added context %s not found in list output: %s", tt.contextName, listOutput)
				}
			}
		})
	}
}

// TestListCommand tests the 'kubecm list' command
func TestListCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	tests := []struct {
		name         string
		contexts     []string
		args         []string
		wantContains []string
		wantErr      bool
	}{
		{
			name:         "list single context",
			contexts:     []string{"context-1"},
			args:         []string{"list"},
			wantContains: []string{"context-1"},
			wantErr:      false,
		},
		{
			name:         "list multiple contexts",
			contexts:     []string{"dev-context", "prod-context", "staging-context"},
			args:         []string{"list"},
			wantContains: []string{"dev-context", "prod-context", "staging-context"},
			wantErr:      false,
		},
		{
			name:         "list with alias ls",
			contexts:     []string{"context-1"},
			args:         []string{"ls"},
			wantContains: []string{"context-1"},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test kubeconfig with multiple contexts
			testConfig := CreateMultiContextKubeconfig(t, "test.yaml", tt.contexts)

			// Run kubecm list
			output, err := RunKubecmWithEnv(t, map[string]string{"KUBECONFIG": testConfig}, tt.args...)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v, output: %s", err, output)
			}

			// Verify all expected contexts are in the output
			if !tt.wantErr {
				for _, expected := range tt.wantContains {
					if !strings.Contains(output, expected) {
						t.Errorf("Expected context %s not found in output: %s", expected, output)
					}
				}
			}
		})
	}
}

// TestSwitchCommand tests the 'kubecm switch' command
func TestSwitchCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	tests := []struct {
		name     string
		contexts []string
		switchTo string
		wantErr  bool
	}{
		{
			name:     "switch to existing context",
			contexts: []string{"context-1", "context-2"},
			switchTo: "context-2",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test kubeconfig with multiple contexts
			testConfig := CreateMultiContextKubeconfig(t, "test.yaml", tt.contexts)

			// Run kubecm switch
			output, err := RunKubecmWithEnv(t, map[string]string{"KUBECONFIG": testConfig},
				"switch", tt.switchTo)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v, output: %s", err, output)
			}
		})
	}
}

// TestDeleteCommand tests the 'kubecm delete' command
func TestDeleteCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	tests := []struct {
		name      string
		contexts  []string
		deleteCtx string
		wantErr   bool
	}{
		{
			name:      "delete existing context",
			contexts:  []string{"context-1", "context-2", "context-3"},
			deleteCtx: "context-2",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test kubeconfig with multiple contexts
			testConfig := CreateMultiContextKubeconfig(t, "test.yaml", tt.contexts)

			// Run kubecm delete (no -y flag needed when context name is provided)
			output, err := RunKubecmWithEnv(t, map[string]string{"KUBECONFIG": testConfig},
				"delete", tt.deleteCtx)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v, output: %s", err, output)
			}

			// Verify the context was deleted
			if !tt.wantErr {
				listOutput, err := RunKubecmWithEnv(t, map[string]string{"KUBECONFIG": testConfig}, "list")
				if err != nil {
					t.Logf("Note: list command may fail if no contexts remain")
				} else if strings.Contains(listOutput, tt.deleteCtx) {
					t.Errorf("Deleted context %s still found in list output: %s", tt.deleteCtx, listOutput)
				}
			}
		})
	}
}
