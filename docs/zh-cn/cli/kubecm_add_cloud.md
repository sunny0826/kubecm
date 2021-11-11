## kubecm add cloud

获取公有云 K8s 服务的 `kubeconfig`

### 简介

获取公有云 K8s 服务的 `kubeconfig`

```
kubecm add cloud [flags]
```

### 示例

```bash

# Set env secret key
export ACCESS_KEY_ID=xxx
export ACCESS_KEY_SECRET=xxx
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
```

### 全局选项

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
  -c, --cover           Overwrite local kubeconfig files
```

### SEE ALSO

* [kubecm add](kubecm_add.md)	 - Add KubeConfig to $HOME/.kube/config
