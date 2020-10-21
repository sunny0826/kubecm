package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

type ListCommand struct {
	baseCommand
}

func (lc *ListCommand) Init() {
	lc.command = &cobra.Command{
		Use:   "ls",
		Short: "List kubeconfig",
		Long:  "List kubeconfig",
		Aliases: []string{"l"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return lc.runList(cmd, args)
		},
		Example: listExample(),
	}
	lc.command.DisableFlagsInUseLine = true
}

func (lc *ListCommand) runList(command *cobra.Command, args []string) error {
	err := Formatable()
	if err != nil {
		return err
	}
	err = ClusterStatus()
	if err != nil {
		log.Fatalf("Cluster check failure!\n%v", err)
		return err
	}
	return nil
}

func listExample() string {
	return `
# List all the contexts in your kubeconfig file
kubecm ls
# Aliases
kubecm l
`
}