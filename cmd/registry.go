package cmd

import (
	"github.com/spf13/cobra"
)

// RegistryCommand registry command struct
type RegistryCommand struct {
	BaseCommand
}

// Init RegistryCommand
func (rc *RegistryCommand) Init() {
	rc.command = &cobra.Command{
		Use:   "registry [COMMANDS]",
		Short: "Manage kubeconfig registries (Git-backed distribution)",
		Long:  "Manage Git-backed kubeconfig registries for team-based cluster distribution",
	}
	rc.AddCommands(
		&RegistryAddCommand{},
		&RegistryListCommand{},
		&RegistrySyncCommand{},
		&RegistryRemoveCommand{},
		&RegistryUpdateCommand{},
	)
}
