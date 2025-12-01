## kubecm clear

Clear lapsed context, cluster and user

### Synopsis

Clear lapsed context, cluster and user

```
kubecm clear
```

### Examples

```

# Clear lapsed context, cluster and user (default is $HOME/.kube/config)
kubecm clear
# Customised clear lapsed files
kubecm clear config.yaml test.yaml

```

### Options

```
  -h, --help   help for clear
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
* [kubecm clear docs](kubecm_clear_docs.md)	 - Open document website

