## kubecm registry update

Update a registry's role, variables, or branch

### Synopsis

Modify a registry's configuration and optionally re-sync

```
kubecm registry update <name> [flags]
```

### Examples

```
# Change role
kubecm registry update rubix --role backend

# Update a variable
kubecm registry update rubix --var Username=new.user

# Change branch and sync
kubecm registry update rubix --ref develop
```

### Options

```
  -h, --help          help for update
      --ref string    new git ref/branch
      --role string   new role
      --var strings   set template variables as KEY=VALUE (repeatable)
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

