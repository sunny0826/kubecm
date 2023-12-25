## `kubecm` completion

Generate completion script

### Synopsis

Please follow the guide below to configure auto-completion for both regular `kubecm` and `kubectl kc` (as `kubectl` plugin).

#### Configuration for auto-completion of `kubecm`

Bash:

```shell
# To load completions for the current session, execute once:
$ source <(kubecm completion bash)

# To load completions for each session, execute once:
# Linux:
$ kubecm completion bash > /etc/bash_completion.d/kubecm

# macOS:
$ kubecm completion bash > /usr/local/etc/bash_completion.d/kubecm

# You will need to start a new shell for this setup to take effect.
```

Zsh:

```shell
# If shell completion is not already enabled in your environment,
# you will need to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ kubecm completion zsh > "${fpath[1]}/_kubecm"

# You will need to start a new shell for this setup to take effect.
```

fish:

```shell
# To load completions for the current session, execute once:
$ kubecm completion fish | source

# To load completions for each session, execute once:
$ kubecm completion fish > ~/.config/fish/completions/kubecm.fish

# You will need to start a new shell for this setup to take effect.
```

PowerShell:

```powershell
# To load completions for the current session, execute once:
PS> kubecm completion powershell | Out-String | Invoke-Expression

# To load completions for every new session, run:
PS> kubecm completion powershell > kubecm.ps1
# and source this file from your PowerShell profile.
```

#### Configuration for auto-completion of `kubectl kc` (as `kubectl` plugin)

Unlike direct execution of `kubecm`, when executed as a `kubectl` plugin (i.e., using `kubecm` installed with [`krew`](https://krew.sigs.k8s.io/) or invoked as `kubectl kc`), auto-completion requires some additional configuration based on [Configuration for auto-completion of `kubecm`](#configuration-for-auto-completion-of-kubecm) to enable auto-completion when use `kubectl kc` + <kbd>TAB</kbd>.

> How does this work?
>
> In a nut shell, `kubectl` looks for a executable script or binary that named with the `kubectl_complete-<plugin-name>` pattern as a **convention** in the `$PATH` environment variable. If it does find, `kubectl` will pass the arguments that calls from `kubectl <plugin-name>` to the conventionally configured auto-completion script to seek out for completion candidates list.

Therefore, in order for the auto-completion of `kubectl kc` + <kbd>TAB</kbd> to work, there are several options for you to choose from and configure.

Note: **PowerShell's auto-completion requires no additional configuration and works right out of the box.**

##### Use pre-written auto-completion script

You can download the pre-written auto-completion script and configure it automatically by executing the following command:

```shell
mkdir -p ~/.config/.kubectl-plugin-completions && \
  curl -LO "https://raw.githubusercontent.com/sunny0826/kubecm/master/hack/kubectl-plugin-completions/kubectl_complete-kc.sh" && \
  mv kubectl_complete-kc.sh ~/.config/.kubectl-plugin-completions/kubectl_complete-kc && \
  chmod +x ~/.config/.kubectl-plugin-completions/kubectl_complete-kc
```

Append the `~/.config/.kubectl-plugin-completions` directory we've created just now to the `$PATH` environment variable to satisfy the auto-completion convention of the `kubectl` plugin:

Bash:

```shell
echo 'export PATH=$PATH:~/.config/.kubectl-plugin-completions' >> ~/.bashrc
```

Zsh:

```shell
echo 'export PATH=$PATH:~/.config/.kubectl-plugin-completions' >> ~/.zshrc
```

Fish:

```shell
fish_add_path ~/.config/.kubectl-plugin-completions
```

##### Manually configure auto-completion script

1. Create the directory `.config/.kubectl-plugin-completions` for the auto-completion scripts:

```shell
mkdir -p ~/.config/.kubectl-plugin-completions
```

2. Write the following content into an executable shell script named `kubectl_complete-kc` in this directory:

```shell
#!/usr/bin/env sh

# Call the __complete command passing it all arguments
kubectl kc __complete "$@"
```

You can also complete this step as one-liner using the following command:

```shell
cat <<EOF >~/.config/.kubectl-plugin-completions/kubectl_complete-kc
#!/usr/bin/env sh

# Call the __complete command passing it all arguments
kubectl kc __complete "\$@"
EOF
```

3. Add executable permissions to the file:

```shell
chmod +x ~/.config/.kubectl-plugin-completions/kubectl_complete-kc
```

Append the `~/.config/.kubectl-plugin-completion` directory we've created just now to the `$PATH` environment variable to satisfy the auto-completion convention of the `kubectl` plugin:

Bash:

```shell
echo 'export PATH=$PATH:~/.config/.kubectl-plugin-completions' >> ~/.bashrc
```

Zsh:

```shell
echo 'export PATH=$PATH:~/.config/.kubectl-plugin-completions' >> ~/.zshrc
```

Fish:

```shell
fish_add_path ~/.config/.kubectl-plugin-completions
```

##### Use a `kubectl` plugin `kubectl-plugin_completion` to automatically generate and configure

Please refer to the [`kubectl-plugin_completion`](https://github.com/marckhouzam/kubectl-plugin_completion) plugin's documentation for configuration.

---


```shell
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
