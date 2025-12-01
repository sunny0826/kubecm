## kubecm list

List KubeConfig

### Synopsis

List KubeConfig

```
kubecm list
```

### Examples

```

# List all the contexts in your KubeConfig file
kubecm list
# Aliases
kubecm ls
kubecm l
# Filter out keywords(Multi-keyword support)
kubecm ls kind k3s
# Useful environment variables
KUBECM_DISABLE_K8S_MORE_INFO: it will disable the k8s more info in the output

```

### Options

```
  -h, --help           help for list
      --no-server      Hide the server column
      --short-server   Shorten the server endpoint
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
* [kubecm list docs](kubecm_list_docs.md)	 - Open document website

