## kubecm add

Merge configuration file with $HOME/.kube/config

### Synopsis

Merge configuration file with $HOME/.kube/config

```
kubecm add [flags]
```

### Examples

```

# Merge 1.yaml with $HOME/.kube/config
kubecm add -f 1.yaml 

# Merge 1.yaml and name contexts test with $HOME/.kube/config
kubecm add -f 1.yaml -n test

# Overwrite the original kubeconfig file
kubecm add -f 1.yaml -c

```

### Options

```
  -c, --cover         Overwrite the original kubeconfig file
  -f, --file string   Path to merge kubeconfig files
  -h, --help          help for add
  -n, --name string   The name of contexts. if this field is null,it will be named with file name.
```

### Options inherited from parent commands

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
```
