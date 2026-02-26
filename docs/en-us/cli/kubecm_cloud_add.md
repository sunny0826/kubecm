## kubecm cloud add

Add kubeconfig from cloud

### Synopsis

Add kubeconfig from cloud

```
kubecm cloud add [flags]
```

### Examples

```

# Supports AWS, Ali Cloud, Tencent Cloud, Rancher and Azure
# The AK/AS of the cloud platform will be retrieved directly
# if it exists in the environment variable,
# otherwise a prompt box will appear asking for it.

# Interaction: select kubeconfig from the cloud
kubecm cloud add

# AlibabaCloud
export ACCESS_KEY_ID=YOUR_AKID
export ACCESS_KEY_SECRET=YOUR_SECRET_KEY
kubecm cloud add --provider alibabacloud --cluster_id=xxxxxx

# Tencent Cloud
export TENCENTCLOUD_SECRET_ID=YOUR_SECRET_ID
export TENCENTCLOUD_SECRET_KEY=YOUR_SECRET_KEY
kubecm cloud add --provider tencent --region_id=ap-guangzhou

# Rancher
export RANCHER_SERVER_URL=https://xxx
export RANCHER_API_KEY=YOUR_API_KEY
kubecm cloud add --provider rancher

# AWS with named profile (recommended)
kubecm cloud add --provider aws --aws_profile my-profile --region_id us-east-1

# AWS with profile from environment
export AWS_PROFILE=my-profile
kubecm cloud add --provider aws

# AWS with static credentials (backward compatible)
export AWS_ACCESS_KEY_ID=YOUR_AKID
export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_KEY
kubecm cloud add --provider aws --region_id us-east-1

# Azure with default SDK auth
kubecm cloud add --provider azure

# Azure with service principal
export AZURE_CLIENT_ID=YOUR_CLIENT_ID
export AZURE_CLIENT_SECRET=YOUR_CLIENT_SECRET
export AZURE_TENANT_ID=YOUR_TENANT_ID
kubecm cloud add --provider azure

```

### Options

```
  -h, --help   help for add
```

### Options inherited from parent commands

```
      --aws_profile string   AWS profile name (from ~/.aws/config)
      --cluster_id string    kubernetes cluster id
      --config string        path of kubeconfig (default "$HOME/.kube/config")
      --create               Create a new kubeconfig file if not exists
  -m, --mac-notify           enable to display Mac notification banner
      --provider string      public cloud
      --region_id string     cloud region id
  -s, --silence-table        enable/disable output of context table on successful config update
  -u, --ui-size int          number of list items to show in menu at once (default 10)
```

### SEE ALSO

* [kubecm cloud](kubecm_cloud.md)	 - Manage kubeconfig from cloud

