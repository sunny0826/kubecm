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

// Role is a roles/<name>.yaml file listing fragments for a team role.
type Role struct {
	APIVersion    string           `yaml:"apiVersion"`
	Kind          string           `yaml:"kind"`
	Metadata      RegistryMetadata `yaml:"metadata"`
	ContextPrefix string           `yaml:"contextPrefix,omitempty"`
	Fragments     []string         `yaml:"fragments"`
}

// Fragment is a fragments/<name>.yaml file describing one cluster.
type Fragment struct {
	APIVersion string           `yaml:"apiVersion"`
	Kind       string           `yaml:"kind"`
	Metadata   RegistryMetadata `yaml:"metadata"`
	Provider   string           `yaml:"provider"` // aws, azure, static
	AWS        *AWSFragment     `yaml:"aws,omitempty"`
	Azure      *AzureFragment   `yaml:"azure,omitempty"`
	Kubeconfig string           `yaml:"kubeconfig,omitempty"` // for static provider
}

// AWSFragment holds AWS EKS cluster reference.
type AWSFragment struct {
	Region  string `yaml:"region"`
	Cluster string `yaml:"cluster"`
	Profile string `yaml:"profile,omitempty"`
}

// AzureFragment holds Azure AKS cluster reference.
type AzureFragment struct {
	SubscriptionID string `yaml:"subscriptionId"`
	ResourceGroup  string `yaml:"resourceGroup"`
	Cluster        string `yaml:"cluster"`
	TenantID       string `yaml:"tenantId,omitempty"`
}

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
