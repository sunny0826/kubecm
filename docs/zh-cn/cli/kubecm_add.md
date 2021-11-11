## kubecm add

将 KubeConfig 加入到 `$HOME/.kube/config`

### 简介

将 KubeConfig 加入到 `$HOME/.kube/config`

```
kubecm add [flags]
```

### 示例

```

# Merge 1.yaml with $HOME/.kube/config
kubecm add -f 1.yaml 
# Interaction: select kubeconfig from the cloud
kubecm add cloud

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
```
