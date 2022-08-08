## kubecm merge

合并选中的 KubeConfig

### 简介

合并选中的 KubeConfig

```
kubecm merge [flags]
```

### 示例

```

# Merge multiple kubeconfig
kubecm merge 1st.yaml 2nd.yaml 3rd.yaml
# Merge KubeConfig in the dir directory
kubecm merge -f dir
# Merge KubeConfig in the dir directory to the specified file.
kubecm merge -f dir --config kubecm.config

```

### 选项

```
  -y, --assumeyes       skip interactive file overwrite confirmation
  -f, --folder string   KubeConfig folder
  -h, --help            help for merge
```

### 全局选项

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
      --ui-size int     number of list items to show in menu at once (default 4)
```
