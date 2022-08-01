## kubecm add

Add KubeConfig to $HOME/.kube/config

### Synopsis

Add KubeConfig to $HOME/.kube/config

```
kubecm add [flags]
```

### Examples

```

# Merge test.yaml with $HOME/.kube/config
kubecm add -f test.yaml 
# Add kubeconfig from stdin
cat /etc/kubernetes/admin.conf |  kubecm add -f -
# Interaction: select kubeconfig from the cloud
kubecm add cloud

```

### Options

```
  -c, --cover         Overwrite local kubeconfig files
  -f, --file string   Path to merge kubeconfig files
  -h, --help          help for add
```

### Options inherited from parent commands

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
      --ui-size int     number of list items to show in menu at once (default 4)
```

