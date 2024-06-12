## kubecm delete

Delete the specified context from the kubeconfig

### Synopsis

Delete the specified context from the kubeconfig

```
kubecm delete [flags]
```

### Examples

```

# Delete the context interactively
kubecm delete
# Delete the context
kubecm delete my-context
# Deleting multiple contexts
kubecm delete my-context1 my-context2

```

### Options

```
  -h, --help   help for delete
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
* [kubecm delete docs](kubecm_delete_docs.md)	 - Open document website

