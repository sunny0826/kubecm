## kubecm registry add

Add a new kubeconfig registry

### Synopsis

Clone a Git-backed kubeconfig registry and sync contexts

```
kubecm registry add [flags]
```

### Examples

```
# Add a registry with inline variables
kubecm registry add --name rubix --url git@bitbucket.org:rubixdig/kubeconfig-registry.git --role devops --var Username=clark.n

# Add a registry (will prompt for required variables)
kubecm registry add --name rubix --url git@bitbucket.org:rubixdig/kubeconfig-registry.git --role devops
```

### Options

```
  -h, --help          help for add
      --name string   registry name (required)
      --ref string    git branch/ref (default "main")
      --role string   role to use (required)
      --url string    git repository URL (required)
      --var strings   template variables as KEY=VALUE (repeatable)
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

