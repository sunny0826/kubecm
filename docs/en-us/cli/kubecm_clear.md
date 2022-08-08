## kubecm clear

Clear lapsed context, cluster and user

### Synopsis

Clear lapsed context, cluster and user

```
kubecm clear
```

### Examples

```

# Clear lapsed context, cluster and user (default is /Users/guoxudong/.kube/config)
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
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
      --ui-size int     number of list items to show in menu at once (default 4)
```

### SEE ALSO

* [kubecm](kubecm.md)	 - KubeConfig Manager.

