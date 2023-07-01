## kubecm cloud add

Add kubeconfig from cloud

### Synopsis

Add kubeconfig from cloud

```
kubecm cloud add [flags]
```

### Examples

```

# Supports AWS, Ali Cloud, Tencent Cloud and Rancher
# The AK/AS of the cloud platform will be retrieved directly 
# if it exists in the environment variable, 
# otherwise a prompt box will appear asking for it.

# Set env AliCloud secret key
export ACCESS_KEY_ID=YOUR_AKID
export ACCESS_KEY_SECRET=YOUR_SECRET_KEY

# Set env Tencent secret key
export TENCENTCLOUD_SECRET_ID=YOUR_SECRET_ID
export TENCENTCLOUD_SECRET_KEY=YOUR_SECRET_KEY

# Set env Rancher secret key
export RANCHER_SERVER_URL=https://xxx
export RANCHER_API_KEY=YOUR_API_KEY

# Set env AWS secret key
export AWS_ACCESS_KEY_ID=YOUR_AKID
export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_KEY

# Interaction: select kubeconfig from the cloud
kubecm cloud add
# Add kubeconfig from cloud
kubecm cloud add --provider alibabacloud --cluster_id=xxxxxx

```

### Options

```
  -h, --help   help for add
```

### Options inherited from parent commands

```
      --cluster_id string   kubernetes cluster id
      --config string       path of kubeconfig (default "/Users/guoxudong/.kube/config")
  -m, --mac-notify          enable to display Mac notification banner
      --provider string     public cloud
      --region_id string    cloud region id
      --ui-size int         number of list items to show in menu at once (default 4)
```

### SEE ALSO

* [kubecm cloud](kubecm_cloud.md)	 - Manage kubeconfig from cloud

