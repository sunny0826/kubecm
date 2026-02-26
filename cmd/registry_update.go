package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sunny0826/kubecm/pkg/registry"
)

// RegistryUpdateCommand update a registry's settings
type RegistryUpdateCommand struct {
	BaseCommand
}

// Init RegistryUpdateCommand
func (c *RegistryUpdateCommand) Init() {
	c.command = &cobra.Command{
		Use:   "update <name>",
		Short: "Update a registry's role, variables, or branch",
		Long:  "Modify a registry's configuration and optionally re-sync",
		Example: `# Change role
kubecm registry update rubix --role backend

# Update a variable
kubecm registry update rubix --var Username=new.user

# Change branch and sync
kubecm registry update rubix --ref develop`,
		Args: cobra.ExactArgs(1),
		RunE: c.runUpdate,
	}
	c.command.Flags().String("role", "", "new role")
	c.command.Flags().String("ref", "", "new git ref/branch")
	c.command.Flags().StringSlice("var", nil, "set template variables as KEY=VALUE (repeatable)")
}

func (c *RegistryUpdateCommand) runUpdate(cmd *cobra.Command, args []string) error {
	name := args[0]

	cfg, err := registry.LoadConfig()
	if err != nil {
		return err
	}

	entry := cfg.GetRegistry(name)
	if entry == nil {
		return fmt.Errorf("registry %q not found", name)
	}

	changed := false

	if role, _ := cmd.Flags().GetString("role"); role != "" {
		// Validate role exists
		repoDir, err := registry.RegistryDir(name)
		if err != nil {
			return err
		}
		if _, err := registry.LoadRole(repoDir, role); err != nil {
			return err
		}
		entry.Role = role
		changed = true
	}

	if ref, _ := cmd.Flags().GetString("ref"); ref != "" {
		entry.Ref = ref
		changed = true
	}

	if varSlice, _ := cmd.Flags().GetStringSlice("var"); len(varSlice) > 0 {
		if entry.Variables == nil {
			entry.Variables = make(map[string]string)
		}
		for k, v := range parseVarSlice(varSlice) {
			entry.Variables[k] = v
		}
		changed = true
	}

	if !changed {
		return fmt.Errorf("nothing to update, use --role, --ref, or --var")
	}

	if err := registry.SaveConfig(cfg); err != nil {
		return err
	}

	fmt.Printf("Registry %q updated. Run 'kubecm registry sync %s' to apply changes.\n", name, name)
	return nil
}
