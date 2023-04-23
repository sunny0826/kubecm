## kubecm completion

生成 bash/zsh 自动补全脚本

### 简介

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

### 选项

```
  -h, --help   help for completion
```

### 全局选项

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
  -m, --mac-notify      enable to display Mac notification banner
      --ui-size int     number of list items to show in menu at once (default 4)
```
