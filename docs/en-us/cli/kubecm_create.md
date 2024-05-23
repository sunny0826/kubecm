## kubecm create

Create new KubeConfig(experiment)

### Synopsis

Create new KubeConfig(experiment)

WARNING: This command is experimental and this feature is only supported in kubernates v1.24 and later.


```
kubecm create [flags]
```

### Examples

```

# Create new KubeConfig(experiment)
kubecm create
# Create new KubeConfig(experiment) with flags
kubecm create --user test --namespace default --cluster-role view --context-name kind-kind

```

### Options

```
      --cluster-role string   cluster role for user
      --context-name string   context name for kubeconfig
  -h, --help                  help for create
  -n, --namespace string      namespace for user
      --print-clean-up        print clean up command
      --user string           user name for kubeconfig
```

### Options inherited from parent commands

```
      --config string   path of kubeconfig (default "$HOME/.kube/config")
  -m, --mac-notify      enable to display Mac notification banner
  -s, --silence-table   enable/disable output of context table on successful config update
  -u, --ui-size int     number of list items to show in menu at once (default 4)
```

### SEE ALSO

* [kubecm](kubecm.md)	 - KubeConfig Manager.

