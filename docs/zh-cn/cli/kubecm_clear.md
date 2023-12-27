## kubecm clear

清理 kubeconfig 中失效的 `context`, `cluster` 和 `user`

### 简介

默认情况下会清理 `~/.kube/confg` 中失效的 `context`, `cluster` 和 `user`，后跟非默认的 kubeconfig 路径，可以清理指定文件，支持指定多文件。

```
kubecm clear
```

### 示例

```

# Clear lapsed context, cluster and user (default is /Users/saybot/.kube/config)
kubecm clear
# 自定义清理
kubecm clear config.yaml test.yaml

```

### Options

```
  -h, --help   help for clear
```

### Options inherited from parent commands

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
  -m, --mac-notify      enable to display Mac notification banner
  -s, --silence-table   enable/disable output of context table on successful config update
      --ui-size int     number of list items to show in menu at once (default 4)
```
