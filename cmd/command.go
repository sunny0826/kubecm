package cmd

import (
	"github.com/spf13/cobra"
)

// Command is cli command interface
type Command interface {
	// Init command
	Init()

	// CobraCmd
	CobraCmd() *cobra.Command
}

// BaseCommand struct
type BaseCommand struct {
	command *cobra.Command
}

// Init BaseCommand
func (bc *BaseCommand) Init() {}

// CobraCmd returns BaseCommand
func (bc *BaseCommand) CobraCmd() *cobra.Command {
	return bc.command
}

// Name returns name of BaseCommand
func (bc *BaseCommand) Name() string {
	return bc.command.Name()
}

// AddCommands adds children commands to the parent command
func (bc *BaseCommand) AddCommands(children ...Command) {
	for _, child := range children {
		child.Init()
		childCmd := child.CobraCmd()
		bc.CobraCmd().AddCommand(childCmd)
	}
}
