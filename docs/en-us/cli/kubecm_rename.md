## kubecm rename

Rename the contexts of kubeconfig

### Synopsis

Rename the contexts of kubeconfig

```
kubecm rename [flags]
```

### Examples

```

# Renamed the context interactively
kubecm rename
# Renamed the context non-interactively
kubecm rename <kube-context-name> <new-kube-context-name>

```

### Options

```
  -h, --help   help for rename
```

### Options inherited from parent commands

```
      --config string   path of kubeconfig (default "$HOME/.kube/config")
      --create          Create a new kubeconfig file if not exists
  -m, --mac-notify      enable to display Mac notification banner
  -s, --silence-table   enable/disable output of context table on successful config update
  -u, --ui-size int     number of list items to show in menu at once (default 4)
```

### SEE ALSO

* [kubecm](kubecm.md)	 - KubeConfig Manager.
* [kubecm rename docs](kubecm_rename_docs.md)	 - Open document website

