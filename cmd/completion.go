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
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
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
	return `To load completions for kubecm cli:

For Bash users:

  # To load completions for the current session, execute once:
  source <(kubecm completion bash)

  # To load completions for each session, execute once (you will need to start a new
  # shell for this setup to take effect.):

  # Linux:
  kubecm completion bash > /etc/bash_completion.d/kubecm

  # macOS:
  kubecm completion bash > /usr/local/etc/bash_completion.d/kubecm

For Zsh users, execute:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once (you will need to start a new
  # shell for this setup to take effect.):
  kubecm completion zsh > "${fpath[1]}/_kubecm"

For fish users, execute:

  # Load completions for the current session, execute once:
  kubecm completion fish | source

  # Load completions for each session, execute once (you will need to start a new
  # shell for this setup to take effect.):
  kubecm completion fish > ~/.config/fish/completions/kubecm.fish

For PowerShell users, execute:

  # To load completions for the current session, execute once:
  kubecm completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  kubecm completion powershell > kubecm.ps1
  # and source this file from your PowerShell profile.

---

To load completions for kubectl kc when used as kubectl plugin:

1. Use the following command as one-liner to add the completion to the current shell session:

  mkdir -p ~/.config/.kubectl-plugin-completions
  cat <<EOF >~/.config/.kubectl-plugin-completions/kubectl_complete-kc
  #!/usr/bin/env sh

  # Call the __complete command passing it all arguments
  kubectl kc __complete "\$@"
  EOF
  chmod +x ~/.config/.kubectl-plugin-completions/kubectl_complete-kc

2. Append the directory to the $PATH environment variable.

For Bash users, execute:

  echo 'export PATH=$PATH:~/.config/.kubectl-plugin-completions' >> ~/.bashrc

For Zsh users, execute:

  echo 'export PATH=$PATH:~/.config/.kubectl-plugin-completions' >> ~/.zshrc

For fish users, execute:

  fish_add_path ~/.config/.kubectl-plugin-completions

---
`
}
