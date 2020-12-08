## kubecm add

Add kubeconfig to $HOME/.kube/config

### Synopsis

Add kubeconfig to $HOME/.kube/config

```
kubecm add [flags]
```

### Examples

```

# Merge test.yaml with $HOME/.kube/config
kubecm add -f test.yaml 

```

### Options

```
  -f, --file string   Path to merge kubeconfig files
  -h, --help          help for add
```

### Options inherited from parent commands

```
      --config string   path of kubeconfig (default "/Users/saybot/.kube/config")
```
