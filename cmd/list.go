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
		Run: func(cmd *cobra.Command, args []string) {
			lc.runList(cmd, args)
		},
		Example: listExample(),
	}
	lc.command.DisableFlagsInUseLine = true
}

func (lc *ListCommand) runList(command *cobra.Command, args []string) {
	if len(args) == 0 {
		err := Formatable(nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := Formatable(args)
		if err != nil {
			log.Fatal(err)
		}
	}
	err := ClusterStatus()
	if err != nil {
		log.Fatalf("Cluster check failure!\n%v", err)
	}
}

func listExample() string {
	return `
# List all the contexts in your kubeconfig file
kubecm ls
# Aliases
kubecm l
`
}