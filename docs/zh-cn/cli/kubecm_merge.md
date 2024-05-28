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
# Merge test.yaml with $HOME/.kube/config and add a prefix before context name
kubecm merge test.yaml --context-prefix test
# Merge test.yaml with $HOME/.kube/config and define the attributes used for composing the context name
kubecm merge test.yaml --context-template user,cluster
# Merge test.yaml with $HOME/.kube/config, define the attributes used for composing the context name and add a prefix before context name
kubecm merge test.yaml --context-template user,cluster --context-prefix demo
# Merge test.yaml with $HOME/.kube/config and select the context to be added in interactive mode
kubecm merge test.yaml --select-context
# Merge test.yaml with $HOME/.kube/config and specify the context to be added
kubecm merge test.yaml --context context1,context2
```

### 选项

```
  -y, --assumeyes                  skip interactive file overwrite confirmation
      --context strings            specify the context to be merged
      --context-prefix string      add a prefix before context name
      --context-template strings   define the attributes used for composing the context name, available values: filename, user, cluster, context, namespace (default [context])
  -f, --folder string              KubeConfig folder
  -h, --help                       help for merge
      --select-context             select the context to be merged in interactive mode
```

### 全局选项

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
  -m, --mac-notify      enable to display Mac notification banner
  -s, --silence-table   enable/disable output of context table on successful config update
      --ui-size int     number of list items to show in menu at once (default 4)
```
