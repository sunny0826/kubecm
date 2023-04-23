## kubecm delete

删除指定 context

### 简介

删除指定 context

```
kubecm delete [flags]
```

### 示例

```

# Delete the context interactively
kubecm delete
# Delete the context
kubecm delete my-context
# Deleting multiple contexts
kubecm delete my-context1 my-context2

```

### 选项

```
  -h, --help   help for delete
```

### 全局选项

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
  -m, --mac-notify      enable to display Mac notification banner
      --ui-size int     number of list items to show in menu at once (default 4)
```
