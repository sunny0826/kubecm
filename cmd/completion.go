package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// CompletionCommand completion cmd struct
type CompletionCommand struct {
	BaseCommand
}

// Init CompletionCommand
func (cc *CompletionCommand) Init() {
	cc.command = &cobra.Command{
		Use:       "completion [bash|zsh|fish|powershell]",
		Short:     "Generate completion script",
		Long:      longDetail(),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Args:      cobra.ExactValidArgs(1),
		Aliases:   []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				_ = cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				_ = cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				_ = cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				_ = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}
}

func longDetail() string {
	return `To load completions:

Bash:

  $ source <(kubecm completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ kubecm completion bash > /etc/bash_completion.d/kubecm
  # macOS:
  $ kubecm completion bash > /usr/local/etc/bash_completion.d/kubecm

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ kubecm completion zsh > "${fpath[1]}/_kubecm"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ kubecm completion fish | source

  # To load completions for each session, execute once:
  $ kubecm completion fish > ~/.config/fish/completions/kubecm.fish

PowerShell:

  PS> kubecm completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> kubecm completion powershell > kubecm.ps1
  # and source this file from your PowerShell profile.
`
}
