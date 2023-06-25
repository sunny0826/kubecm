## kubecm export

导出指定 context

### 简介

导出指定 context

```
kubecm export [flags]
```

### 示例

```
# Export context to myconfig.yaml file
kubecm export -f myconfig.yaml my-context1
# Export multiple contexts to myconfig.yaml file
kubecm export -f myconfig.yaml my-context1 my-context2
```

### 选项

```
  -f, --file string   Path to export kubeconfig files 
  -h, --help   help for export
```

### 全局选项

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
  -m, --mac-notify      enable to display Mac notification banner
      --ui-size int     number of list items to show in menu at once (default 4)
```
