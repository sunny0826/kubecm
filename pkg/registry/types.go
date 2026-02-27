package registry

import "time"

// RegistryMeta is the top-level registry.yaml in a registry repo.
type RegistryMeta struct {
	APIVersion string           `yaml:"apiVersion"`
	Kind       string           `yaml:"kind"`
	Metadata   RegistryMetadata `yaml:"metadata"`
	Variables  []VariableSpec   `yaml:"variables,omitempty"`
}

// RegistryMetadata holds name and description.
type RegistryMetadata struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
}

// VariableSpec defines a template variable.
type VariableSpec struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Required    bool   `yaml:"required,omitempty"`
	Default     string `yaml:"default,omitempty"`
}

// Role is a roles/<name>.yaml file listing clusters for a team role.
type Role struct {
	APIVersion    string           `yaml:"apiVersion"`
	Kind          string           `yaml:"kind"`
	Metadata      RegistryMetadata `yaml:"metadata"`
	ContextPrefix string           `yaml:"contextPrefix,omitempty"`
	Fragments     []string         `yaml:"fragments,omitempty"` // legacy format
	Contexts      []RoleContext    `yaml:"contexts,omitempty"`
}

// RoleContext associates a cluster with an optional user override.
type RoleContext struct {
	Cluster  string `yaml:"cluster,omitempty"`
	Fragment string `yaml:"fragment,omitempty"` // deprecated: use cluster
	User     string `yaml:"user,omitempty"`
	Name     string `yaml:"name,omitempty"`
}

// ClusterRef returns the cluster name, supporting both the new cluster:
// field and the deprecated fragment: field for backward compatibility.
func (rc *RoleContext) ClusterRef() string {
	if rc.Cluster != "" {
		return rc.Cluster
	}
	return rc.Fragment
}

// NormalizedContexts returns the role's contexts list.
// If the new contexts: format is used, it is returned as-is.
// Otherwise the legacy fragments: list is converted.
func (r *Role) NormalizedContexts() []RoleContext {
	if len(r.Contexts) > 0 {
		return r.Contexts
	}
	result := make([]RoleContext, len(r.Fragments))
	for i, f := range r.Fragments {
		result[i] = RoleContext{Cluster: f}
	}
	return result
}

// User is a users/<name>.yaml file describing authentication credentials.
type User struct {
	APIVersion string           `yaml:"apiVersion"`
	Kind       string           `yaml:"kind"`
	Metadata   RegistryMetadata `yaml:"metadata"`
	Provider   string           `yaml:"provider"`
	AWS        *AWSUserConfig   `yaml:"aws,omitempty"`
	Azure      *AzureUserConfig `yaml:"azure,omitempty"`
}

// AWSUserConfig holds AWS-specific user settings.
type AWSUserConfig struct {
	Profile string `yaml:"profile"`
}

// AzureUserConfig holds Azure-specific user settings.
type AzureUserConfig struct {
	TenantID string `yaml:"tenantId,omitempty"`
}

// Cluster is a clusters/<name>.yaml (or fragments/<name>.yaml) file describing one cluster.
type Cluster struct {
	APIVersion string           `yaml:"apiVersion"`
	Kind       string           `yaml:"kind"`
	Metadata   RegistryMetadata `yaml:"metadata"`
	Provider   string           `yaml:"provider"` // aws, azure, static
	AWS        *AWSClusterConfig  `yaml:"aws,omitempty"`
	Azure      *AzureClusterConfig `yaml:"azure,omitempty"`
	Kubeconfig string           `yaml:"kubeconfig,omitempty"` // for static provider
}

// AWSClusterConfig holds AWS EKS cluster reference.
type AWSClusterConfig struct {
	Region  string `yaml:"region"`
	Cluster string `yaml:"cluster"`
	Profile string `yaml:"profile,omitempty"`
}

// AzureClusterConfig holds Azure AKS cluster reference.
type AzureClusterConfig struct {
	SubscriptionID string `yaml:"subscriptionId"`
	ResourceGroup  string `yaml:"resourceGroup"`
	Cluster        string `yaml:"cluster"`
	TenantID       string `yaml:"tenantId,omitempty"`
}

// Deprecated aliases for backward compatibility with external code.
type Fragment = Cluster
type AWSFragment = AWSClusterConfig
type AzureFragment = AzureClusterConfig

// KubecmConfig is the local ~/.kubecm/config.yaml state file.
type KubecmConfig struct {
	APIVersion string           `yaml:"apiVersion"`
	Kind       string           `yaml:"kind"`
	Registries []RegistryEntry  `yaml:"registries"`
}

// RegistryEntry tracks one configured registry.
type RegistryEntry struct {
	Name            string            `yaml:"name"`
	URL             string            `yaml:"url"`
	Ref             string            `yaml:"ref"`
	Role            string            `yaml:"role"`
	Variables       map[string]string `yaml:"variables,omitempty"`
	LastSync        *time.Time        `yaml:"lastSync,omitempty"`
	ManagedContexts []string          `yaml:"managedContexts,omitempty"`
}
