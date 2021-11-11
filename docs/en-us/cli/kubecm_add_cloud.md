## kubecm add cloud

Manage kubeconfig with public cloud

### Synopsis

Manage kubeconfig with public cloud

```
kubecm add cloud [flags]
```

### Examples

```bash

# Set env secret key
export ACCESS_KEY_ID=xxx
export ACCESS_KEY_SECRET=xxx
# Interaction: select kubeconfig from the cloud
kubecm add cloud
# Add kubeconfig from cloud
kubecm add cloud --provider alibabacloud --cluster_id=xxxxxx

```

### Options

```
      --cluster_id string   kubernetes cluster id
  -h, --help                help for cloud
      --provider string     public cloud
```

### Options inherited from parent commands

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
  -c, --cover           Overwrite local kubeconfig files
```

### SEE ALSO

* [kubecm add](kubecm_add.md)	 - Add KubeConfig to $HOME/.kube/config
