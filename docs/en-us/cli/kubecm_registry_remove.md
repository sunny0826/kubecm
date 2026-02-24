## kubecm registry remove

Remove a kubeconfig registry

### Synopsis

Remove a registry and optionally its managed kubeconfig contexts

```
kubecm registry remove <name> [flags]
```

### Examples

```
# Remove registry and its contexts
kubecm registry remove rubix

# Remove registry but keep kubeconfig contexts
kubecm registry remove rubix --keep-contexts
```

### Options

```
  -h, --help            help for remove
      --keep-contexts   keep managed kubeconfig contexts
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

* [kubecm registry](kubecm_registry.md)	 - Manage kubeconfig registries (Git-backed distribution)

