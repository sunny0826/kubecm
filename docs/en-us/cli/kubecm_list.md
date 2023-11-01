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

```

### Options

```
  -h, --help   help for list
```

### Options inherited from parent commands

```
      --config string   path of kubeconfig (default "/Users/user/.kube/config")
  -m, --mac-notify      enable to display Mac notification banner
  -s, --silence-table   enable/disable output of context table on successful config update
      --ui-size int     number of list items to show in menu at once (default 4)
```

### SEE ALSO

* [kubecm](kubecm.md)	 - KubeConfig Manager.

