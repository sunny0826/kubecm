## kubecm alias

为所有 contexts 生成 alias，可自动将其注入到 `.zshrc` 或 `.bash_profile` 中。

### 简介

为所有 contexts 生成 alias，可自动将其注入到 `.zshrc` 或 `.bash_profile` 中。

```
kubecm alias
```

### 示例

```

$ kubecm alias -o zsh
# add context to ~/.zshrc
$ kubecm alias -o bash
# add context to ~/.bash_profile

```

### 选项

```
  -h, --help         help for alias
  -o, --out string   output to ~/.zshrc or ~/.bash_profile
```

### 全局选项

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
  -m, --mac-notify      enable to display Mac notification banner
  -s, --silence-table   enable/disable output of context table on successful config update
      --ui-size int     number of list items to show in menu at once (default 4)
```
