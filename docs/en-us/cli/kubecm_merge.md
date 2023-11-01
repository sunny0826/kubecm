## kubecm merge

Merge multiple kubeconfig files into one

### Synopsis

Merge multiple kubeconfig files into one

```
kubecm merge [flags]
```

### Examples

```

# Merge multiple kubeconfig
kubecm merge 1st.yaml 2nd.yaml 3rd.yaml
# Merge KubeConfig in the dir directory
kubecm merge -f dir
# Merge KubeConfig in the dir directory to the specified file.
kubecm merge -f dir --config kubecm.config

```

### Options

```
  -y, --assumeyes       skip interactive file overwrite confirmation
  -f, --folder string   KubeConfig folder
  -h, --help            help for merge
```

### Options inherited from parent commands

```
      --config string   path of kubeconfig (default "/Users/user/.kube/config")
  -m, --mac-notify      enable to display Mac notification banner
  -s, --silence-table   enable/disable output of context table on successful config update
      --ui-size int     number of list items to show in menu at once (default 4)
```

### SEE ALSO

* [kubecm](kubecm.md)	 - KubeConfig Manager.

