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
# Merge test.yaml with $HOME/.kube/config and add a prefix before context name
kubecm add -cf test.yaml --context-prefix test
# Merge test.yaml with $HOME/.kube/config and define the attributes used for composing the context name
kubecm add -f test.yaml --context-template user,cluster
# Merge test.yaml with $HOME/.kube/config, define the attributes used for composing the context name and add a prefix before context name
kubecm add -f test.yaml --context-template user,cluster --context-prefix demo
# Merge test.yaml with $HOME/.kube/config and override context name, it's useful if there is only one context in the kubeconfig file
kubecm add -f test.yaml --context-name test
# Merge test.yaml with $HOME/.kube/config and select the context to be added in interactive mode
kubecm add -f test.yaml --select-context
# Add kubeconfig from stdin
cat /etc/kubernetes/admin.conf | kubecm add -f -
```

### 选项

```
      --context-name string        override context name when add kubeconfig context, when context-name is set, context-prefix and context-template parameters will be ignored
      --context-prefix string      add a prefix before context name
      --context-template strings   define the attributes used for composing the context name, available values: filename, user, cluster, context, namespace (default [context])
  -c, --cover                      overwrite local kubeconfig files
  -f, --file string                path to merge kubeconfig files
  -h, --help                       help for add
      --select-context             select the context to be added
```

### 全局选项

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
  -m, --mac-notify      enable to display Mac notification banner
  -s, --silence-table   enable/disable output of context table on successful config update
      --ui-size int     number of list items to show in menu at once (default 4)
```
