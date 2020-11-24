## kubecm add

将 kubeconfig 加入到 `$HOME/.kube/config`

### 简介

将 kubeconfig 加入到 `$HOME/.kube/config`

```
kubecm add [flags]
```

### 示例

```

# Merge 1.yaml with $HOME/.kube/config
kubecm add -f 1.yaml 

# Merge 1.yaml and name contexts test with $HOME/.kube/config
kubecm add -f 1.yaml -n test

# Overwrite the original kubeconfig file
kubecm add -f 1.yaml -c

```

### 选项

```
  -c, --cover         Overwrite the original kubeconfig file
  -f, --file string   Path to merge kubeconfig files
  -h, --help          help for add
  -n, --name string   The name of contexts. if this field is null,it will be named with file name.
```

### 全局选项

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
```
