## kubecm registry

Manage kubeconfig registries (Git-backed distribution)

### Synopsis

Manage Git-backed kubeconfig registries for team-based cluster distribution

### Options

```
  -h, --help   help for registry
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
* [kubecm registry add](kubecm_registry_add.md)	 - Add a new kubeconfig registry
* [kubecm registry list](kubecm_registry_list.md)	 - List configured registries
* [kubecm registry remove](kubecm_registry_remove.md)	 - Remove a kubeconfig registry
* [kubecm registry sync](kubecm_registry_sync.md)	 - Sync kubeconfig from registries
* [kubecm registry update](kubecm_registry_update.md)	 - Update a registry's role, variables, or branch

