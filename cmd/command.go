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
func (bc *BaseCommand) Init() {
}

// CobraCmd return BaseCommand
func (bc *BaseCommand) CobraCmd() *cobra.Command {
	return bc.command
}

// Name return name of BaseCommand
func (bc *BaseCommand) Name() string {
	return bc.command.Name()
}

//AddCommand is add child command to the parent command
func (bc *BaseCommand) AddCommand(child Command) {
	child.Init()
	childCmd := child.CobraCmd()
	bc.CobraCmd().AddCommand(childCmd)
}
