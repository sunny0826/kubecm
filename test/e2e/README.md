# E2E Testing Guide for KubeCM

This guide explains how to run and extend end-to-end (e2e) tests for kubecm.

## Overview

The e2e tests validate the complete workflow of kubecm by testing the actual binary with real command executions. These tests ensure that different components work together correctly and catch issues that unit tests might miss.

## Running E2E Tests Locally

### Prerequisites

- Go 1.24.2 or later
- Make (optional, for using Makefile targets)

### Quick Start

The easiest way to run e2e tests is using the provided script:

```bash
# Run e2e tests (builds kubecm automatically)
./test/run-e2e.sh

# Run with verbose output
./test/run-e2e.sh -v

# Skip building kubecm (use existing binary)
./test/run-e2e.sh --no-build
```

### Manual Execution

You can also run e2e tests manually:

```bash
# 1. Build kubecm
go build -o bin/kubecm .

# 2. Set the binary path
export KUBECM_BIN=$(pwd)/bin/kubecm

# 3. Run e2e tests
cd test/e2e
go test -v ./...
```

### Running Specific Tests

```bash
# Run a specific test
cd test/e2e
go test -v -run TestAddCommand

# Run tests matching a pattern
go test -v -run TestList

# Skip e2e tests with short flag
go test -short ./...
```

## Test Structure

### Directory Layout

```
test/
├── e2e/
│   ├── e2e_test.go        # Test main with setup/teardown
│   ├── utils.go           # Helper functions and utilities
│   ├── basic_test.go      # Basic command tests (add, list, switch, delete)
│   └── advanced_test.go   # Advanced command tests (merge, rename, namespace, clear)
└── run-e2e.sh            # E2E test runner script
```

### Test Categories

#### Basic Commands (`basic_test.go`)
- **TestAddCommand**: Tests adding kubeconfig contexts
- **TestListCommand**: Tests listing contexts (including `ls` alias)
- **TestSwitchCommand**: Tests switching between contexts
- **TestDeleteCommand**: Tests deleting contexts

#### Advanced Commands (`advanced_test.go`)
- **TestMergeCommand**: Tests merging multiple kubeconfig files
- **TestRenameCommand**: Tests renaming contexts
- **TestNamespaceCommand**: Tests setting namespaces
- **TestClearCommand**: Tests clearing all contexts

## Writing New E2E Tests

### Basic Test Template

```go
package e2e

import (
	"strings"
	"testing"
)

func TestYourCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "your test case",
			args:    []string{"your", "command", "args"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			testConfig := CreateTestKubeconfig(t, "test.yaml", "test-cluster", "test-context")

			// Run kubecm command
			output, err := RunKubecmWithEnv(t, 
				map[string]string{"KUBECONFIG": testConfig}, 
				tt.args...)

			// Verify results
			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v, output: %s", err, output)
			}

			// Additional assertions
			if !strings.Contains(output, "expected string") {
				t.Errorf("Output does not contain expected string: %s", output)
			}
		})
	}
}
```

### Helper Functions

The `utils.go` file provides several helper functions:

#### Running Commands

```go
// Run kubecm with arguments
output, err := RunKubecm(t, "list")

// Run kubecm with environment variables
output, err := RunKubecmWithEnv(t, 
	map[string]string{"KUBECONFIG": "/path/to/config"}, 
	"switch", "context-name")
```

#### Creating Test Kubeconfigs

```go
// Create a single-context kubeconfig
configPath := CreateTestKubeconfig(t, "test.yaml", "cluster-name", "context-name")

// Create a multi-context kubeconfig
configPath := CreateMultiContextKubeconfig(t, "multi.yaml", 
	[]string{"context-1", "context-2", "context-3"})

// Get a temporary kubeconfig path
tempPath := GetTempKubeconfig(t)
```

### Best Practices

1. **Use table-driven tests**: Define multiple test cases in a slice for better organization
2. **Clean up resources**: Use `t.TempDir()` for temporary files (automatically cleaned up)
3. **Test error cases**: Include both success and failure scenarios
4. **Verify side effects**: Check that commands produce expected changes
5. **Use meaningful names**: Make test names descriptive of what they test
6. **Skip in short mode**: Add `if testing.Short() { t.Skip(...) }` for expensive tests
7. **Log command output**: Use `t.Logf()` to log command execution and output

### Example: Testing a New Command

Here's a complete example for testing a hypothetical `export` command:

```go
func TestExportCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	tests := []struct {
		name       string
		contexts   []string
		exportCtx  string
		outputFile string
		wantErr    bool
	}{
		{
			name:       "export existing context",
			contexts:   []string{"dev", "prod"},
			exportCtx:  "dev",
			outputFile: "exported.yaml",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test kubeconfig
			testConfig := CreateMultiContextKubeconfig(t, "test.yaml", tt.contexts)
			
			// Prepare output path
			tmpDir := t.TempDir()
			outputPath := tmpDir + "/" + tt.outputFile

			// Run export command
			output, err := RunKubecmWithEnv(t, 
				map[string]string{"KUBECONFIG": testConfig},
				"export", tt.exportCtx, "-o", outputPath)

			// Verify results
			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v, output: %s", err, output)
			}

			// Verify exported file exists
			if !tt.wantErr {
				if _, err := os.Stat(outputPath); os.IsNotExist(err) {
					t.Errorf("Exported file not created at: %s", outputPath)
				}
			}
		})
	}
}
```

## CI Integration

E2E tests run automatically in CI through GitHub Actions. See `.github/workflows/e2e-test.yaml` for the workflow configuration.

The CI workflow:
1. Sets up Go environment
2. Builds kubecm
3. Creates test Kubernetes clusters using Kind
4. Runs the e2e test suite
5. Reports results

## Troubleshooting

### Common Issues

**Issue**: `kubecm binary not found`
```bash
# Solution: Build kubecm first
go build -o bin/kubecm .
```

**Issue**: Tests timeout
```bash
# Solution: Increase timeout
go test -timeout 20m ./...
```

**Issue**: Permission denied on test script
```bash
# Solution: Make script executable
chmod +x test/run-e2e.sh
```

### Debug Tips

1. **Enable verbose output**: Run with `-v` flag
2. **Check test logs**: Look for `t.Logf()` output showing command execution
3. **Run single test**: Use `-run` flag to isolate issues
4. **Examine temp files**: Tests use `t.TempDir()` - add debugging to see paths

## Contributing

When adding new features to kubecm:

1. Write e2e tests that cover the new functionality
2. Ensure tests follow the existing patterns
3. Update this documentation if adding new test categories
4. Run the full e2e test suite before submitting PR

For more information, see the [Contributing Guide](../../CONTRIBUTING.md).

## Further Reading

- [Go Testing Package](https://pkg.go.dev/testing)
- [Table Driven Tests in Go](https://go.dev/wiki/TableDrivenTests)
- [KubeCM Documentation](https://kubecm.cloud)
