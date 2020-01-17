package cmd

func CmdInit() *baseCommand {
	cli := NewCli()
	baseCmd := &baseCommand{
		command: cli.rootCmd,
	}
	// add version command
	baseCmd.AddCommand(&VersionCommand{})
	// add add command
	addCommand := &AddCommand{}
	baseCmd.AddCommand(addCommand)
	// add completion command
	completionCommand := &CompletionCommand{}
	baseCmd.AddCommand(completionCommand)
	// add delete command
	deleteCommand := &DeleteCommand{}
	baseCmd.AddCommand(deleteCommand)
	// add merge command
	mergeCommand := &MergeCommand{}
	baseCmd.AddCommand(mergeCommand)
	// add rename command
	renameCommand := &RenameCommand{}
	baseCmd.AddCommand(renameCommand)
	// add switch command
	switchCommand := &SwitchCommand{}
	baseCmd.AddCommand(switchCommand)

	return baseCmd
}
