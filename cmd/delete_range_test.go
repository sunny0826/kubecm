package cmd

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"testing"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// create test config
func createTestContexts() map[string]*clientcmdapi.Context {
	return map[string]*clientcmdapi.Context{
		"dev-cluster1": {Cluster: "cluster1", AuthInfo: "user1"},
		"dev-cluster2": {Cluster: "cluster2", AuthInfo: "user2"},
		"prod-cluster": {Cluster: "cluster3", AuthInfo: "user3"},
		"test-staging": {Cluster: "cluster4", AuthInfo: "user4"},
	}
}

func TestMatchContexts(t *testing.T) {
	tests := []struct {
		name          string
		contexts      map[string]*clientcmdapi.Context
		pattern       string
		matchMode     string
		expected      []string
		expectedError error
	}{
		{
			name:      "prefix dev-",
			contexts:  createTestContexts(),
			pattern:   "dev-",
			matchMode: "prefix",
			expected:  []string{"dev-cluster1", "dev-cluster2"},
		},
		{
			name:      "suffix -cluster",
			contexts:  createTestContexts(),
			pattern:   "-cluster",
			matchMode: "suffix",
			expected:  []string{"prod-cluster"},
		},
		{
			name:      "contains staging",
			contexts:  createTestContexts(),
			pattern:   "staging",
			matchMode: "contains",
			expected:  []string{"test-staging"},
		},
		{
			name:          "no matching contexts",
			contexts:      createTestContexts(),
			pattern:       "nonexistent",
			matchMode:     "contains",
			expected:      nil,
			expectedError: nil,
		},
		{
			name:          "invalid match mode",
			contexts:      createTestContexts(),
			pattern:       "dev-",
			matchMode:     "invalid",
			expected:      nil,
			expectedError: fmt.Errorf("invalid match mode: %s, must be one of: prefix, suffix, contains", "invalid"),
		},
		{
			name:          "empty pattern",
			contexts:      createTestContexts(),
			pattern:       "",
			matchMode:     "prefix",
			expected:      nil,
			expectedError: errors.New("pattern cannot be empty"),
		},
		{
			name:          "empty contexts",
			contexts:      map[string]*clientcmdapi.Context{},
			pattern:       "dev-",
			matchMode:     "prefix",
			expected:      nil,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchContexts(tt.contexts, tt.pattern, tt.matchMode)

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				if got != nil {
					t.Errorf("expected nil result, got %v", got)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			sort.Strings(got)
			sort.Strings(tt.expected)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected matches %v, got %v", tt.expected, got)
			}
		})
	}
}
