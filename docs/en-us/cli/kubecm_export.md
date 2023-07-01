## kubecm export

Export the specified context from the kubeconfig

### Synopsis

Export the specified context from the kubeconfig

```
kubecm export [flags]
```

### Examples

```

# Export context to myconfig.yaml file
kubecm export -f myconfig.yaml my-context1
# Export multiple contexts to myconfig.yaml file
kubecm export -f myconfig.yaml my-context1 my-context2

```

### Options

```
  -f, --file string   Path to export kubeconfig files
  -h, --help          help for export
```

### Options inherited from parent commands

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
  -m, --mac-notify      enable to display Mac notification banner
      --ui-size int     number of list items to show in menu at once (default 4)
```

### SEE ALSO

* [kubecm](kubecm.md)	 - KubeConfig Manager.

