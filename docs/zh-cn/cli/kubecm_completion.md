## `kubecm` completion

生成自动补全脚本

### 简介

请根据下面的指南进行操作为普通的 `kubecm` 和 `kubectl kc`（`kubectl` 插件）配置基本的自动补全功能。

#### 配置 `kubecm` 的自动补全

Bash:

```shell
$ source <(kubecm completion bash)

# 可以通过运行一次下面的命令来为将来的每一个 Shell 会话激活自动补全：
# Linux:
$ kubecm completion bash > /etc/bash_completion.d/kubecm

# macOS:
$ kubecm completion bash > /usr/local/etc/bash_completion.d/kubecm

# 在执行完毕后，你需要重新打开一个新的终端才能使补全脚本生效。
```

Zsh:

```shell
# 如果你当前的 zsh 环境中尚未激活自动补全的功能，你需要先激活它。你
# 可以通过下面的命令来激活：

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# 可以通过运行一次下面的命令来为将来的每一个 Shell 会话激活自动补全：
$ kubecm completion zsh > "${fpath[1]}/_kubecm"

# 在执行完毕后，你需要重新打开一个新的终端才能使补全脚本生效。
```

fish:

```shell
# 可以通过运行一次下面的命令来为当前的 Shell 会话激活自动补全：
$ kubecm completion fish | source

# 可以通过运行一次下面的命令来为将来的每一个 Shell 会话激活自动补全：
$ kubecm completion fish > ~/.config/fish/completions/kubecm.fish

# 在执行完毕后，你需要重新打开一个新的终端才能使补全脚本生效。
```

PowerShell:

```shell
# 可以通过运行一次下面的命令来为当前的 Shell 会话激活自动补全：
PS> kubecm completion powershell | Out-String | Invoke-Expression

# 为将来的每一个 Shell 会话激活自动补全，请执行下面的命令：
PS> kubecm completion powershell > kubecm.ps1
# 然后在 PowerShell 的 profile 中引入这个文件。
```

#### 为作为 `kubectl` 插件的 `kubecm` 配置自动补全

与直接执行 `kubecm` 不同的是，在作为 `kubectl` 插件（即使用 [`krew`](https://krew.sigs.k8s.io/) 安装的 `kubecm`，或者以 `kubectl kc` 的形式调用）执行时，自动补全需要在[配置 `kubecm` 的自动补全](#配置-kubecm-的自动补全)的基础上进行一些额外的配置才能使 `kubectl kc` + <kbd>TAB</kbd> 时的自动补全生效。

> 原理是什么？
>
> 简而言之，`kubectl` 的插件会在 `$PATH` 环境变量中寻找遵循 `kubectl_complete-<插件名>` 命名规则的可执行文件作为**约定**，然后在调用 `kubectl <插件名>` 的时候由 `kubectl` 展开参数并传递参数给调用这个按约定配置的自动补全脚本。

因此，为了使 `kubectl kc` + <kbd>TAB</kbd> 的自动补全生效，下面有几种方案。

Note: **PowerShell 的自动补全无需额外配置，开箱即用。**

##### 使用预编写的自动补全脚本

你可以通过执行下面的命令来下载预编写的自动补全脚本并完成配置：

```shell
mkdir -p ~/.config/.kubectl-plugin-completions && \
  curl -LO "https://raw.githubusercontent.com/sunny0826/kubecm/master/hack/kubectl-plugin-completions/kubectl_complete-kc.sh" && \
  mv kubectl_complete-kc.sh ~/.config/.kubectl-plugin-completions/kubectl_complete-kc && \
  chmod +x ~/.config/.kubectl-plugin-completions/kubectl_complete-kc
```

将 `~/.config/.kubectl-plugin-completions` 目录追加到 `$PATH` 环境变量中，以满足 `kubectl` 插件的自动补全约定：

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

##### 手动创建

1. 创建用于放置自动补全脚本的目录 `.config/.kubectl-plugin-completions`：

```shell
mkdir -p ~/.config/.kubectl-plugin-completions
```

2. 填充下面的内容到该目录下一个名为 `kubectl_complete-kc` 的可执行文件中：

```shell
#!/usr/bin/env sh

# Call the __complete command passing it all arguments
kubectl kc __complete "$@"
```

你也可以通过下面的命令来一次性完成这一步骤：

```shell
cat <<EOF >~/.config/.kubectl-plugin-completions/kubectl_complete-kc
#!/usr/bin/env sh

# Call the __complete command passing it all arguments
kubectl kc __complete "\$@"
EOF
```

3. 为该文件添加可执行权限：

```shell
chmod +x ~/.config/.kubectl-plugin-completions/kubectl_complete-kc
```

4. 将该目录追加到 `$PATH` 环境变量中，以满足 `kubectl` 插件的自动补全约定：

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

##### 使用 `kubectl-plugin_completion` 插件自动生成和配置

请参考 [`kubectl-plugin_completion`](https://github.com/marckhouzam/kubectl-plugin_completion) 插件的文档进行配置。

---

```shell
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
