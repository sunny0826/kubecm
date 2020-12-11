package cmd

// NewBaseCommand cmd struct
func NewBaseCommand() *BaseCommand {
	cli := NewCli()
	baseCmd := &BaseCommand{
		command: cli.rootCmd,
	}
	baseCmd.AddCommands(
		&VersionCommand{},    // version command
		&AddCommand{},        // add command
		&CompletionCommand{}, // completion command
		&DeleteCommand{},     // delete command
		&MergeCommand{},      // merge command
		&RenameCommand{},     // rename command
		&SwitchCommand{},     // switch command
		&NamespaceCommand{},  // namespace command
		&ListCommand{},       // list command
		&AliasCommand{},      // alias command
		&ClearCommand{},      // clear command
	)

	return baseCmd
}
