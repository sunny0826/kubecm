## kubecm merge

选定目录，合并其中所有的 kubeconfig

### 简介

选定目录，合并其中所有的 kubeconfig

```
kubecm merge [flags]
```

### 示例

```

# Merge kubeconfig in the dir directory
kubecm merge -f dir

# Merge kubeconfig in the directory and overwrite the original kubeconfig file
kubecm merge -f dir -c

```

### 选项

```
  -c, --cover           Overwrite the original kubeconfig file
  -f, --folder string   Kubeconfig folder
  -h, --help            help for merge
```

### 全局选项

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
```
