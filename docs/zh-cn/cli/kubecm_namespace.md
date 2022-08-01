## kubecm namespace

切换当前集群 namespace

### 简介


切换当前集群 namespace


```
kubecm namespace [flags]
```

![ns](../../static/ns.gif)

### 示例

```

# Switch Namespace interactively
kubecm namespace
# or
kubecm ns
# change to namespace of kube-system
kubecm ns kube-system

```

### 选项

```
  -h, --help   help for namespace
```

### 全局选项

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
      --ui-size int     number of list items to show in menu at once (default 4)
```
