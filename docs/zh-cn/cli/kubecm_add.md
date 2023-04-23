## kubecm add

将 KubeConfig 加入到 `$HOME/.kube/config`

### 简介

将 KubeConfig 加入到 `$HOME/.kube/config`

```
kubecm add [flags]
```

>注意：如果 `-c` 被设置，且添加的 kubeconfig 文件中有**超过一个** context，会出现如下情况：
>- 如果设置了 `--context-name`，则 context 会以 `<context-name-0>`, `<context-name-1>` 的形式产生
>- 如果没有设置 `--context-name`，则会以 `<file-name-{hash}>` 的方式展示，其中 `{hash}` 是文件名的 MD5 哈希值

### 示例

```

# Merge test.yaml with $HOME/.kube/config
kubecm add -f test.yaml 
# Merge test.yaml with $HOME/.kube/config and rename context name
kubecm add -cf test.yaml --context-name test
# Add kubeconfig from stdin
cat /etc/kubernetes/admin.conf |  kubecm add -f -

```

### 选项

```
      --context-name string   override context name when add kubeconfig context
  -c, --cover         Overwrite local kubeconfig files
  -f, --file string   Path to merge kubeconfig files
  
  -h, --help          help for add
```

### 全局选项

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
  -m, --mac-notify      enable to display Mac notification banner
      --ui-size int     number of list items to show in menu at once (default 4)
```
