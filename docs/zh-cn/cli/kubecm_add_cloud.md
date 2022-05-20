## kubecm add cloud

获取公有云 K8s 服务的 `kubeconfig`

### 简介

获取公有云 K8s 服务的 `kubeconfig`

```
kubecm add cloud [flags]
```

目前支持阿里云和腾讯云，相应平台的 AK/AS 如果存在于环境变量中将会直接获取，否则将出现提示框询问。

### 示例

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

### 选项

```
      --cluster_id string   kubernetes cluster id
  -h, --help                help for cloud
      --provider string     public cloud
      --region_id string    cloud region id
```

### 全局选项

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
  -c, --cover           Overwrite local kubeconfig files
      --ui-size int     number of list items to show in menu at once (default 4)
```

### SEE ALSO

* [kubecm add](kubecm_add.md)	 - Add KubeConfig to $HOME/.kube/config
