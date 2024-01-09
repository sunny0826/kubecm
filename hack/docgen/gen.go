package main

import (
	"github.com/spf13/cobra"
	"log"

	"github.com/spf13/cobra/doc"
	"github.com/sunny0826/kubecm/cmd"
)

func main() {
	kubecm := cmd.NewBaseCommand().CobraCmd()
	// Docs of the completion command is not generated
	removeCommand(kubecm, "completion")
	// Docs of the namespace command is not generated
	removeCommand(kubecm, "namespace")
	err := doc.GenMarkdownTree(kubecm, "./docs/en-us/cli/")
	if err != nil {
		log.Fatal(err)
	}
}

// removeCommand
func removeCommand(root *cobra.Command, cmdToRemove string) {
	var newCommands []*cobra.Command
	for _, cmd := range root.Commands() {
		if cmd.Name() != cmdToRemove {
			newCommands = append(newCommands, cmd)
		}
	}
	root.ResetCommands()
	root.AddCommand(newCommands...)
}
