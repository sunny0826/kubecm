## kubecm registry sync

Sync kubeconfig from registries

### Synopsis

Pull latest registry changes and sync kubeconfig contexts

```
kubecm registry sync [name] [flags]
```

### Examples

```
# Sync a specific registry
kubecm registry sync rubix

# Sync all registries
kubecm registry sync --all

# Dry-run to see what would change
kubecm registry sync rubix --dry-run
```

### Options

```
      --all       sync all registries
      --dry-run   show what would change without modifying kubeconfig
  -h, --help      help for sync
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

