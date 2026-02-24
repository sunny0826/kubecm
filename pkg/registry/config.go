package registry

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

const (
	kubecmDir  = ".kubecm"
	configFile = "config.yaml"
	registriesDir = "registries"
)

// ConfigDir returns ~/.kubecm/
func ConfigDir() (string, error) {
	home, err := homeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, kubecmDir), nil
}

// ConfigFilePath returns ~/.kubecm/config.yaml
func ConfigFilePath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFile), nil
}

// RegistryDir returns ~/.kubecm/registries/<name>/
func RegistryDir(name string) (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, registriesDir, name), nil
}

// LoadConfig reads ~/.kubecm/config.yaml. Returns empty config if file doesn't exist.
func LoadConfig() (*KubecmConfig, error) {
	path, err := ConfigFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &KubecmConfig{
				APIVersion: "kubecm.io/v1alpha1",
				Kind:       "KubecmConfig",
			}, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg KubecmConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return &cfg, nil
}

// SaveConfig writes the config to ~/.kubecm/config.yaml.
func SaveConfig(cfg *KubecmConfig) error {
	path, err := ConfigFilePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}

// GetRegistry returns the named registry entry, or nil if not found.
func (cfg *KubecmConfig) GetRegistry(name string) *RegistryEntry {
	for i := range cfg.Registries {
		if cfg.Registries[i].Name == name {
			return &cfg.Registries[i]
		}
	}
	return nil
}

// RemoveRegistry removes a registry entry by name. Returns false if not found.
func (cfg *KubecmConfig) RemoveRegistry(name string) bool {
	for i, r := range cfg.Registries {
		if r.Name == name {
			cfg.Registries = append(cfg.Registries[:i], cfg.Registries[i+1:]...)
			return true
		}
	}
	return false
}

// LoadRegistryMeta reads registry.yaml from a cloned registry repo.
func LoadRegistryMeta(repoDir string) (*RegistryMeta, error) {
	path := filepath.Join(repoDir, "registry.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading registry.yaml: %w", err)
	}
	var meta RegistryMeta
	if err := yaml.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("parsing registry.yaml: %w", err)
	}
	return &meta, nil
}

// LoadRole reads roles/<name>.yaml from a cloned registry repo.
func LoadRole(repoDir, roleName string) (*Role, error) {
	path := filepath.Join(repoDir, "roles", roleName+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading role %q: %w", roleName, err)
	}
	var role Role
	if err := yaml.Unmarshal(data, &role); err != nil {
		return nil, fmt.Errorf("parsing role %q: %w", roleName, err)
	}
	return &role, nil
}

// LoadFragment reads fragments/<name>.yaml from a cloned registry repo.
func LoadFragment(repoDir, fragmentName string) (*Fragment, error) {
	path := filepath.Join(repoDir, "fragments", fragmentName+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading fragment %q: %w", fragmentName, err)
	}
	var frag Fragment
	if err := yaml.Unmarshal(data, &frag); err != nil {
		return nil, fmt.Errorf("parsing fragment %q: %w", fragmentName, err)
	}
	return &frag, nil
}

func homeDir() (string, error) {
	u, err := user.Current()
	if err == nil {
		return u.HomeDir, nil
	}
	if runtime.GOOS == "windows" {
		drive := os.Getenv("HOMEDRIVE")
		path := os.Getenv("HOMEPATH")
		if drive != "" && path != "" {
			return drive + path, nil
		}
		if home := os.Getenv("USERPROFILE"); home != "" {
			return home, nil
		}
		return "", fmt.Errorf("cannot determine home directory")
	}
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}
	return "", fmt.Errorf("cannot determine home directory")
}
