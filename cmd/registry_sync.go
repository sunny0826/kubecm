package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sunny0826/kubecm/pkg/registry"
	"k8s.io/client-go/tools/clientcmd"
)

// RegistrySyncCommand sync registries
type RegistrySyncCommand struct {
	BaseCommand
}

// Init RegistrySyncCommand
func (c *RegistrySyncCommand) Init() {
	c.command = &cobra.Command{
		Use:   "sync [name]",
		Short: "Sync kubeconfig from registries",
		Long:  "Pull latest registry changes and sync kubeconfig contexts",
		Example: `# Sync a specific registry
kubecm registry sync rubix

# Sync all registries
kubecm registry sync --all

# Dry-run to see what would change
kubecm registry sync rubix --dry-run`,
		RunE: c.runSync,
	}
	c.command.Flags().Bool("all", false, "sync all registries")
	c.command.Flags().Bool("dry-run", false, "show what would change without modifying kubeconfig")
}

func (c *RegistrySyncCommand) runSync(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	cfg, err := registry.LoadConfig()
	if err != nil {
		return err
	}

	if len(cfg.Registries) == 0 {
		return fmt.Errorf("no registries configured, use 'kubecm registry add' first")
	}

	if all {
		for i := range cfg.Registries {
			repoDir, err := registry.RegistryDir(cfg.Registries[i].Name)
			if err != nil {
				return err
			}
			fmt.Printf("Syncing registry %q...\n", cfg.Registries[i].Name)
			if err := runRegistrySync(cfg, &cfg.Registries[i], repoDir, dryRun); err != nil {
				fmt.Printf("  Error syncing %q: %v\n", cfg.Registries[i].Name, err)
			}
		}
		return nil
	}

	if len(args) == 0 {
		return fmt.Errorf("specify a registry name or use --all")
	}

	name := args[0]
	entry := cfg.GetRegistry(name)
	if entry == nil {
		return fmt.Errorf("registry %q not found", name)
	}

	repoDir, err := registry.RegistryDir(name)
	if err != nil {
		return err
	}

	fmt.Printf("Syncing registry %q...\n", name)
	return runRegistrySync(cfg, entry, repoDir, dryRun)
}

// runRegistrySync is the shared sync logic used by add and sync commands.
func runRegistrySync(cfg *registry.KubecmConfig, entry *registry.RegistryEntry, repoDir string, dryRun bool) error {
	// Git pull
	if err := registry.GitPull(repoDir); err != nil {
		fmt.Printf("  Warning: git pull failed: %v (using cached copy)\n", err)
	}

	// Load current kubeconfig
	kubeConfig, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return fmt.Errorf("loading kubeconfig: %w", err)
	}

	// Sync
	result, err := registry.Sync(repoDir, entry, kubeConfig, dryRun)
	if err != nil {
		return err
	}

	// Print result
	fmt.Print(registry.FormatSyncResult(result))

	if dryRun {
		fmt.Println("  (dry-run, no changes applied)")
		return nil
	}

	// Write kubeconfig
	if err := clientcmd.WriteToFile(*kubeConfig, cfgFile); err != nil {
		return fmt.Errorf("writing kubeconfig: %w", err)
	}

	// Save registry config
	if err := registry.SaveConfig(cfg); err != nil {
		return fmt.Errorf("saving registry config: %w", err)
	}

	return nil
}
