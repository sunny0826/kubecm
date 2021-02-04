## kubecm ls

列举 Context

### 简介

列举 Context，可以通过后跟关键字，筛选出包含该关键字的 context，支持多关键字筛选，关键字之间是 `或` 的关系。

```
kubecm list
```

### 示例

```

# List all the contexts in your KubeConfig file
kubecm list
# Aliases
kubecm ls
kubecm l
# Filter out keywords(Multi-keyword support)
kubecm ls kind k3s

```

### 选项

```
  -h, --help   help for ls
```

### 全局选项

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
```
