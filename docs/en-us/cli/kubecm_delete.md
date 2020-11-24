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
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
```
