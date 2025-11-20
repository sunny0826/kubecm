## kubecm delete range

Delete contexts matching a pattern

### Synopsis

Delete all contexts that match a specified pattern from the kubeconfig

```
kubecm delete range [flags]
```

### Examples

```

# Delete all contexts with prefix "dev-"
kubecm delete range dev-
or 
kubecm delete range --mode prefix dev-

# Delete all contexts with suffix "prod"
kubecm delete range --mode suffix prod

# Delete all contexts containing "staging"
kubecm delete range --mode contains staging

# Force delete all contexts with prefix "dev-" (skip confirmation)
kubecm delete range dev- -y

```

### Options

```
  -h, --help          help for range
      --mode string   Match mode: prefix, suffix, or contains (default "prefix")
  -y, --yes           Skip confirmation prompt
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

* [kubecm delete](kubecm_delete.md)	 - Delete the specified context from the kubeconfig
* [kubecm delete range docs](kubecm_delete_range_docs.md)	 - Open document website

