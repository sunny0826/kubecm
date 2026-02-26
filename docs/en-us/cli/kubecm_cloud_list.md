## kubecm cloud list

List Cluster from cloud

### Synopsis

List Cluster from cloud

```
kubecm cloud list [flags]
```

### Examples

```

# Supports AlibabaCloud, Tencent Cloud, Rancher, AWS and Azure

# Interaction: list clusters from cloud
kubecm cloud list

# AlibabaCloud
kubecm cloud list --provider alibabacloud

# AWS with named profile
kubecm cloud list --provider aws --aws_profile my-profile --region_id eu-west-3

# AWS with default credential chain
kubecm cloud list --provider aws --region_id us-east-1

# Azure
kubecm cloud list --provider azure

```

### Options

```
  -h, --help   help for list
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

