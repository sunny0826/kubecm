## kubecm add

将 KubeConfig 加入到 `$HOME/.kube/config`

### 简介

将 KubeConfig 加入到 `$HOME/.kube/config`

```
kubecm add [flags]
```

### 示例

```

# Merge test.yaml with $HOME/.kube/config
kubecm add -f test.yaml 
# Add kubeconfig from stdin
cat /etc/kubernetes/admin.conf |  kubecm add -f -

```

### 选项

```
  -c, --cover         Overwrite local kubeconfig files
  -f, --file string   Path to merge kubeconfig files
  -h, --help          help for add
```

### 全局选项

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
      --ui-size int     number of list items to show in menu at once (default 4)
```
