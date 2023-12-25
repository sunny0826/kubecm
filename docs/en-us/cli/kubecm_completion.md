## kubecm completion

Generate completion script

### Synopsis

To load completions:

Bash:
  ```
  $ source <(kubecm completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ kubecm completion bash > /etc/bash_completion.d/kubecm
  # macOS:
  $ kubecm completion bash > /usr/local/etc/bash_completion.d/kubecm
  ```
Zsh:
  ```
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ kubecm completion zsh > "${fpath[1]}/_kubecm"

  # You will need to start a new shell for this setup to take effect.
  ```
fish:
  ```
  $ kubecm completion fish | source

  # To load completions for each session, execute once:
  $ kubecm completion fish > ~/.config/fish/completions/kubecm.fish
  ```
PowerShell:
  ```
  PS> kubecm completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> kubecm completion powershell > kubecm.ps1
  # and source this file from your PowerShell profile.
  ```
---


```
kubecm completion [bash|zsh|fish|powershell] [flags]
```

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
  -m, --mac-notify      enable to display Mac notification banner
  -s, --silence-table   enable/disable output of context table on successful config update
      --ui-size int     number of list items to show in menu at once (default 4)
```

### SEE ALSO

* [kubecm](kubecm.md)	 - KubeConfig Manager.

