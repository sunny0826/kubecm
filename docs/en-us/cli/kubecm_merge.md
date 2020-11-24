## kubecm merge

Merge the kubeconfig files in the specified directory

### Synopsis

Merge the kubeconfig files in the specified directory

```
kubecm merge [flags]
```

### Examples

```

# Merge kubeconfig in the dir directory
kubecm merge -f dir

# Merge kubeconfig in the directory and overwrite the original kubeconfig file
kubecm merge -f dir -c

```

### Options

```
  -c, --cover           Overwrite the original kubeconfig file
  -f, --folder string   Kubeconfig folder
  -h, --help            help for merge
```

### Options inherited from parent commands

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
```
