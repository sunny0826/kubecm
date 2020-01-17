package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

// Command is cli command interface
type Command interface {
	// Init command
	Init()

	// CobraCmd
	CobraCmd() *cobra.Command

}

// baseCommand
type baseCommand struct {
	command *cobra.Command
}

func (bc *baseCommand) Init() {
}

func (bc *baseCommand) CobraCmd() *cobra.Command {
	return bc.command
}

func (bc *baseCommand) Name() string {
	return bc.command.Name()
}

//AddCommand is add child command to the parent command
func (bc *baseCommand) AddCommand(child Command) {
	child.Init()
	childCmd := child.CobraCmd()
	childCmd.PreRun = func(cmd *cobra.Command, args []string) {
		Info = log.New(os.Stdout, "Info:", log.Ldate|log.Ltime|log.Lshortfile)
		Warning = log.New(os.Stdout, "Warning:", log.Ldate|log.Ltime|log.Lshortfile)
		Error = log.New(os.Stderr, "Error:", log.Ldate|log.Ltime|log.Lshortfile)
	}
	bc.CobraCmd().AddCommand(childCmd)
}
