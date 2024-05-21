## kubecm cloud list

List Cluster from cloud

### Synopsis

List Cluster from cloud

```
kubecm cloud list [flags]
```

### Examples

```

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
# Interaction: list kubeconfig from cloud
kubecm cloud list
# Add kubeconfig from cloud
kubecm cloud list --provider alibabacloud --cluster_id=xxxxxx

```

### Options

```
  -h, --help   help for list
```

### Options inherited from parent commands

```
      --cluster_id string   kubernetes cluster id
      --config string       path of kubeconfig (default "$HOME/.kube/config")
  -m, --mac-notify          enable to display Mac notification banner
      --provider string     public cloud
      --region_id string    cloud region id
  -s, --silence-table       enable/disable output of context table on successful config update
  -u, --ui-size int         number of list items to show in menu at once (default 4)
```

### SEE ALSO

* [kubecm cloud](kubecm_cloud.md)	 - Manage kubeconfig from cloud

