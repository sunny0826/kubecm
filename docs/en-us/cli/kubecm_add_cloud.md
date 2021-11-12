## kubecm add cloud

Manage kubeconfig with public cloud

### Synopsis

Manage kubeconfig with public cloud

```
kubecm add cloud [flags]
```

Supports [Ali Cloud](https://www.alibabacloud.com/) and [Tencent Cloud](https://cloud.tencent.com/). 
The `AK/AS` of the cloud platform will be retrieved directly if it exists in the environment variable, otherwise a prompt box will appear asking for it.


### Examples

```bash

# Set env AliCloud secret key
export ACCESS_KEY_ID=xxx
export ACCESS_KEY_SECRET=xxx
# Set env Tencent secret key
export TENCENTCLOUD_SECRET_ID=xxx
export TENCENTCLOUD_SECRET_KEY=xxx
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
