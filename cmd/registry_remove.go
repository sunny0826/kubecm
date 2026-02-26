package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sunny0826/kubecm/pkg/registry"
	"k8s.io/client-go/tools/clientcmd"
)

// RegistryRemoveCommand remove a registry
type RegistryRemoveCommand struct {
	BaseCommand
}

// Init RegistryRemoveCommand
func (c *RegistryRemoveCommand) Init() {
	c.command = &cobra.Command{
		Use:     "remove <name>",
		Aliases: []string{"rm"},
		Short:   "Remove a kubeconfig registry",
		Long:    "Remove a registry and optionally its managed kubeconfig contexts",
		Example: `# Remove registry and its contexts
kubecm registry remove rubix

# Remove registry but keep kubeconfig contexts
kubecm registry remove rubix --keep-contexts`,
		Args: cobra.ExactArgs(1),
		RunE: c.runRemove,
	}
	c.command.Flags().Bool("keep-contexts", false, "keep managed kubeconfig contexts")
}

func (c *RegistryRemoveCommand) runRemove(cmd *cobra.Command, args []string) error {
	name := args[0]
	keepContexts, _ := cmd.Flags().GetBool("keep-contexts")

	cfg, err := registry.LoadConfig()
	if err != nil {
		return err
	}

	entry := cfg.GetRegistry(name)
	if entry == nil {
		return fmt.Errorf("registry %q not found", name)
	}

	// Remove managed contexts from kubeconfig
	if !keepContexts && len(entry.ManagedContexts) > 0 {
		kubeConfig, err := clientcmd.LoadFromFile(cfgFile)
		if err != nil {
			return fmt.Errorf("loading kubeconfig: %w", err)
		}

		for _, ctx := range entry.ManagedContexts {
			if err := deleteContext([]string{ctx}, kubeConfig); err != nil {
				fmt.Printf("  Warning: %v\n", err)
			}
		}

		if err := clientcmd.WriteToFile(*kubeConfig, cfgFile); err != nil {
			return fmt.Errorf("writing kubeconfig: %w", err)
		}
	}

	// Remove cloned repo
	repoDir, err := registry.RegistryDir(name)
	if err != nil {
		return err
	}
	if err := os.RemoveAll(repoDir); err != nil {
		fmt.Printf("  Warning: removing repo dir: %v\n", err)
	}

	// Remove from config
	cfg.RemoveRegistry(name)
	if err := registry.SaveConfig(cfg); err != nil {
		return err
	}

	fmt.Printf("Registry %q removed.\n", name)
	return nil
}
