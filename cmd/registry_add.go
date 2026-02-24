package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sunny0826/kubecm/pkg/registry"
)

// RegistryAddCommand add a registry
type RegistryAddCommand struct {
	BaseCommand
}

// Init RegistryAddCommand
func (c *RegistryAddCommand) Init() {
	c.command = &cobra.Command{
		Use:   "add",
		Short: "Add a new kubeconfig registry",
		Long:  "Clone a Git-backed kubeconfig registry and sync contexts",
		Example: `# Add a registry with inline variables
kubecm registry add --name rubix --url git@bitbucket.org:rubixdig/kubeconfig-registry.git --role devops --var Username=clark.n

# Add a registry (will prompt for required variables)
kubecm registry add --name rubix --url git@bitbucket.org:rubixdig/kubeconfig-registry.git --role devops`,
		RunE: c.runAdd,
	}
	c.command.Flags().String("name", "", "registry name (required)")
	c.command.Flags().String("url", "", "git repository URL (required)")
	c.command.Flags().String("role", "", "role to use (required)")
	c.command.Flags().String("ref", "main", "git branch/ref")
	c.command.Flags().StringSlice("var", nil, "template variables as KEY=VALUE (repeatable)")
	_ = c.command.MarkFlagRequired("name")
	_ = c.command.MarkFlagRequired("url")
	_ = c.command.MarkFlagRequired("role")
}

func (c *RegistryAddCommand) runAdd(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	url, _ := cmd.Flags().GetString("url")
	role, _ := cmd.Flags().GetString("role")
	ref, _ := cmd.Flags().GetString("ref")
	varSlice, _ := cmd.Flags().GetStringSlice("var")

	// Load config
	cfg, err := registry.LoadConfig()
	if err != nil {
		return err
	}

	// Check uniqueness
	if cfg.GetRegistry(name) != nil {
		return fmt.Errorf("registry %q already exists", name)
	}

	// Clone repo
	repoDir, err := registry.RegistryDir(name)
	if err != nil {
		return err
	}

	fmt.Printf("Cloning registry %q from %s...\n", name, url)
	if err := registry.GitClone(url, ref, repoDir); err != nil {
		return err
	}

	// Parse registry.yaml for variable specs
	meta, err := registry.LoadRegistryMeta(repoDir)
	if err != nil {
		// Clean up on failure
		os.RemoveAll(repoDir)
		return err
	}

	// Validate role exists
	if _, err := registry.LoadRole(repoDir, role); err != nil {
		os.RemoveAll(repoDir)
		return err
	}

	// Parse provided variables
	vars := parseVarSlice(varSlice)

	// Prompt for missing required variables
	for _, spec := range meta.Variables {
		if _, ok := vars[spec.Name]; !ok {
			if spec.Required {
				val := PromptUI(fmt.Sprintf("Variable %q (%s)", spec.Name, spec.Description), spec.Default)
				vars[spec.Name] = val
			} else if spec.Default != "" {
				vars[spec.Name] = spec.Default
			}
		}
	}

	// Create entry
	entry := registry.RegistryEntry{
		Name:      name,
		URL:       url,
		Ref:       ref,
		Role:      role,
		Variables: vars,
	}
	cfg.Registries = append(cfg.Registries, entry)

	// Save config before sync (so sync can update it)
	if err := registry.SaveConfig(cfg); err != nil {
		os.RemoveAll(repoDir)
		return err
	}

	// Run sync
	fmt.Printf("Syncing registry %q...\n", name)
	return runRegistrySync(cfg, &cfg.Registries[len(cfg.Registries)-1], repoDir, false)
}

// parseVarSlice parses ["KEY=VALUE", ...] into a map.
func parseVarSlice(vars []string) map[string]string {
	m := make(map[string]string)
	for _, v := range vars {
		for i := 0; i < len(v); i++ {
			if v[i] == '=' {
				m[v[:i]] = v[i+1:]
				break
			}
		}
	}
	return m
}
