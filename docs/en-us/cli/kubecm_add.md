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
```

### SEE ALSO

* [kubecm](kubecm.md)	 - KubeConfig Manager.
* [kubecm add cloud](kubecm_add_cloud.md)	 - Manage kubeconfig with public cloud
