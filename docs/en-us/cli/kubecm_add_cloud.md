## kubecm add cloud

Add kubeconfig from public cloud

### Synopsis

Add kubeconfig from public cloud

```
kubecm add cloud [flags]
```

### Examples

```bash

# Supports Ali Cloud and Tencent Cloud
# The AK/AS of the cloud platform will be retrieved directly 
# if it exists in the environment variable, 
# otherwise a prompt box will appear asking for it.

# Set env AliCloud secret key
export ACCESS_KEY_ID=xxx
export ACCESS_KEY_SECRET=xxx
# Set env Tencent secret key
export TENCENTCLOUD_SECRET_ID=xxx
export TENCENTCLOUD_SECRET_KEY=xxx
# Set env Rancher secret key
export RANCHER_SERVER_URL=https://xxx
export RANCHER_API_KEY=xxx
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
      --region_id string    cloud region id
```

### Options inherited from parent commands

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
  -c, --cover           Overwrite local kubeconfig files
      --ui-size int     number of list items to show in menu at once (default 4)
```

### SEE ALSO

* [kubecm add](kubecm_add.md)	 - Add KubeConfig to $HOME/.kube/config
