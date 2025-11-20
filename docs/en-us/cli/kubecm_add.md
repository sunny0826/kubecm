## kubecm add

Add KubeConfig to $HOME/.kube/config

### Synopsis

Add KubeConfig to $HOME/.kube/config

```
kubecm add [flags]
```

### Examples

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
# Merge test.yaml with $HOME/.kube/config and specify the context to be added
kubecm add -f test.yaml --context context1,context2
# Add kubeconfig from stdin
cat /etc/kubernetes/admin.conf | kubecm add -f -
# Merge test.yaml with $HOME/.kube/config and skip TLS certificate verification
kubecm add -f test.yaml --insecure-skip-tls-verify

```

### Options

```
      --context strings            specify the context to be added
      --context-name string        override context name when add kubeconfig context, when context-name is set, context-prefix and context-template parameters will be ignored
      --context-prefix string      add a prefix before context name
      --context-template strings   define the attributes used for composing the context name, available values: filename, user, cluster, context, namespace (default [context])
  -c, --cover                      overwrite local kubeconfig files
  -f, --file string                path to merge kubeconfig files
  -h, --help                       help for add
      --insecure-skip-tls-verify   if true, the server's certificate will not be checked for validity
      --select-context             select the context to be added in interactive mode
```

### Options inherited from parent commands

```
      --config string   path of kubeconfig (default "$HOME/.kube/config")
      --create          Create a new kubeconfig file if not exists
  -m, --mac-notify      enable to display Mac notification banner
  -s, --silence-table   enable/disable output of context table on successful config update
  -u, --ui-size int     number of list items to show in menu at once (default 10)
```

### SEE ALSO

* [kubecm](kubecm.md)	 - KubeConfig Manager.
* [kubecm add docs](kubecm_add_docs.md)	 - Open document website

