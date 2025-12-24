package e2e

import (
	"os"
	"testing"
)

// TestMain handles setup and teardown for all e2e tests
func TestMain(m *testing.M) {
	// Setup
	if err := setup(); err != nil {
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Teardown
	teardown()

	os.Exit(code)
}

func setup() error {
	// Setup will be called before running e2e tests
	return nil
}

func teardown() {
	// Teardown will be called after running e2e tests
}
