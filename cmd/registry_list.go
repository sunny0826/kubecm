package cmd

import (
	"fmt"
	"os"

	"github.com/bndr/gotabulate"
	"github.com/spf13/cobra"
	"github.com/sunny0826/kubecm/pkg/registry"
)

// RegistryListCommand list registries
type RegistryListCommand struct {
	BaseCommand
}

// Init RegistryListCommand
func (c *RegistryListCommand) Init() {
	c.command = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List configured registries",
		Long:    "Show all configured kubeconfig registries and their status",
		RunE:    c.runList,
	}
}

func (c *RegistryListCommand) runList(cmd *cobra.Command, args []string) error {
	cfg, err := registry.LoadConfig()
	if err != nil {
		return err
	}

	if len(cfg.Registries) == 0 {
		fmt.Println("No registries configured. Use 'kubecm registry add' to add one.")
		return nil
	}

	var table [][]string
	for _, r := range cfg.Registries {
		lastSync := "never"
		if r.LastSync != nil {
			lastSync = r.LastSync.Format("2006-01-02 15:04:05")
		}
		table = append(table, []string{
			r.Name,
			r.URL,
			r.Ref,
			r.Role,
			fmt.Sprintf("%d", len(r.ManagedContexts)),
			lastSync,
		})
	}

	tabulate := gotabulate.Create(table)
	tabulate.SetHeaders([]string{"NAME", "URL", "REF", "ROLE", "CONTEXTS", "LAST SYNC"})
	tabulate.SetWrapStrings(false)
	tabulate.SetAlign("left")
	fmt.Fprintln(os.Stdout, tabulate.Render("grid", "left"))
	return nil
}
