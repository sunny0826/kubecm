# Testing Guide

This directory contains the testing infrastructure for kubecm.

## Test Types

### Unit Tests
Unit tests are located in the `cmd/` and `pkg/` directories alongside the code they test. They test individual functions and components in isolation.

**Running unit tests:**
```bash
make test
```

### E2E Tests
End-to-end tests are located in `test/e2e/` and test the complete kubecm binary with realistic workflows. These tests validate that different components work together correctly.

**Running e2e tests:**
```bash
# Using make
make e2e-test

# Using the test runner script
./test/run-e2e.sh

# Manual execution
export KUBECM_BIN=$(pwd)/bin/kubecm
cd test/e2e
go test -v ./...
```

See [test/e2e/README.md](e2e/README.md) for detailed information about e2e tests.

## Test Organization

```
test/
├── e2e/                    # End-to-end tests
│   ├── README.md          # E2E testing documentation
│   ├── e2e_test.go        # Test setup/teardown
│   ├── utils.go           # Helper functions
│   ├── basic_test.go      # Tests for basic commands
│   └── advanced_test.go   # Tests for advanced commands
└── run-e2e.sh             # E2E test runner script
```

## CI/CD Integration

Tests run automatically in GitHub Actions:

- **Unit tests**: Run on every push and PR via `.github/workflows/go.yaml`
- **E2E tests (local)**: Run on every push and PR via `.github/workflows/e2e-test.yaml`
- **E2E tests (Kind)**: Integration tests with real Kubernetes clusters

## Writing Tests

### Unit Tests

Follow Go's standard testing conventions:

```go
func TestMyFunction(t *testing.T) {
    got := MyFunction("input")
    want := "expected output"
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}
```

### E2E Tests

Use table-driven tests and the provided helper functions:

```go
func TestMyCommand(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping e2e test in short mode")
    }

    tests := []struct {
        name    string
        args    []string
        wantErr bool
    }{
        // test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            config := CreateTestKubeconfig(t, "test.yaml", "cluster", "context")
            
            // Execute
            output, err := RunKubecmWithEnv(t, 
                map[string]string{"KUBECONFIG": config}, 
                tt.args...)
            
            // Verify
            // ...
        })
    }
}
```

## Best Practices

1. **Keep tests independent**: Each test should be able to run in isolation
2. **Use descriptive names**: Test names should clearly describe what they test
3. **Test both success and failure**: Include both positive and negative test cases
4. **Clean up resources**: Use `t.TempDir()` for temporary files
5. **Use table-driven tests**: Organize multiple test cases in a slice
6. **Log helpful information**: Use `t.Logf()` to provide context for failures

## Debugging Tests

```bash
# Run specific test
go test -v -run TestSpecificTest ./test/e2e/

# Enable verbose output
go test -v ./test/e2e/

# Run with increased timeout
go test -timeout 20m ./test/e2e/

# Skip long-running tests
go test -short ./test/e2e/
```

## Contributing

When adding new features:

1. Write unit tests for new functions
2. Add e2e tests for new commands or major workflows
3. Ensure all tests pass before submitting PR
4. Update documentation if needed

See [CONTRIBUTING.md](../CONTRIBUTING.md) for more information.
