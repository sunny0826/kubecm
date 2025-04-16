package cmd

// NewBaseCommand cmd struct
func NewBaseCommand() *BaseCommand {
	cli := NewCli()
	baseCmd := &BaseCommand{
		command: cli.rootCmd,
	}
	baseCmd.AddCommands(
		&VersionCommand{},     // version command
		&AddCommand{},         // add command
		&CompletionCommand{},  // completion command
		&DeleteCommand{},      // delete command
		&RangeDeleteCommand{}, // range delete command
		&MergeCommand{},       // merge command
		&RenameCommand{},      // rename command
		&SwitchCommand{},      // switch command
		&NamespaceCommand{},   // namespace command
		&ListCommand{},        // list command
		&AliasCommand{},       // alias command
		&ClearCommand{},       // clear command
		&CreateCommand{},      // create command
		&CloudCommand{},       // cloud command
		&ExportCommand{},      // export command
		&DocsCommand{},        // docs command
	)

	return baseCmd
}
