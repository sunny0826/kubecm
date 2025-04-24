## kubecm range-delete

Delete contexts matching a pattern

### Synopsis

Delete all contexts that match a specified pattern from the kubeconfig

```
kubecm range-delete [flags]
```

### Examples

```

# Delete all contexts with prefix "dev-"
kubecm range-delete dev-
or 
kubecm range-delete -m prefix dev-

# Delete all contexts with suffix "-prod"
kubecm range-delete -m suffix -prod

# Delete all contexts containing "staging"
kubecm range-delete -m contains staging

```

### Options

```
  -h, --help          help for range-delete
      --mode string   Match mode: prefix, suffix, or contains (default "prefix")
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
* [kubecm range-delete docs](kubecm_range-delete_docs.md)	 - Open document website

